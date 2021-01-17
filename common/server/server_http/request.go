package server_http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server"
)

const OperatorJWTKey = "_operator"
const bodyLogLimit = 2048

const onRequest = "on server_http.Request()"

const ReAuthOpKey = "re_auth_operator"
const ReAuthSuffix = "re_auth_suffix"

type ResponseBinary struct {
	MIMEType string
	Data     []byte
}

func Request(serverURL string, ep EndpointConfig, requestData, responseData interface{}, identity *auth.Identity, logfile string) error {
	client := &http.Client{}
	if ep.Handler == nil {
		return errors.Errorf("no ep.Handler: %#v", ep)
	}
	method := ep.Handler.Method

	var reAuthTries []bool
	if identity != nil {
		reAuthTries = []bool{true, false}
	} else {
		reAuthTries = []bool{false}
	}

	var err error
	for _, doReAuth := range reAuthTries {

		// start of single try

		var requestBody []byte
		var requestBodyReader io.Reader

		if requestData != nil {
			switch v := requestData.(type) {
			case []byte:
				requestBody = v
			case *[]byte:
				requestBody = *v
			case string:
				requestBody = []byte(v)
			case *string:
				requestBody = []byte(*v)
			default:
				if requestBody, err = json.Marshal(requestData); err != nil {
					return errors.Wrapf(err, onRequest+": can't marshal request responseData (%#v)", requestData)
				}
			}

			// must be checked for nil instead direct write
			// the external for GET requests expected nil body, but nil-requestData after json.Marshal return not empty responseData

			requestBodyReader = bytes.NewBuffer(requestBody)
		}

		req, err := http.NewRequest(method, serverURL, requestBodyReader)
		if err != nil || req == nil {
			logger.LogRequest(logfile, nil, method, serverURL, nil, requestBody, nil, nil, err, 0)
			return errors.Errorf("can't create request %s %s, got %#v, %s", method, serverURL, req, err)
		} else if req.Body != nil {
			defer common.Close(req.Body, client, nil)
		}

		if identity != nil && identity.JWT != "" {
			req.Header.Add("Authorization", identity.JWT)
		}

		var responseBody []byte

		resp, err := client.Do(req)
		if resp != nil && resp.Body != nil {
			defer common.Close(resp.Body, client, nil)
		}

		if err != nil {
			var statusCode int
			var responseHeaders http.Header
			if resp != nil {
				statusCode = resp.StatusCode
				responseHeaders = resp.Header
				responseBody, _ = ioutil.ReadAll(resp.Body)
			}

			logger.LogRequest(logfile, nil, method, serverURL, req.Header, requestBody, responseHeaders, responseBody, err, statusCode)
			return errors.Wrapf(err, "can't %s %s", method, serverURL)
		}

		responseBody, err = ioutil.ReadAll(resp.Body)
		logger.LogRequest(logfile, nil, method, serverURL, req.Header, requestBody, resp.Header, responseBody, err, resp.StatusCode)
		if err != nil {
			return errors.Wrapf(err, "can't read body from %s %s", method, serverURL)
		}

		if resp.StatusCode == http.StatusUnauthorized && doReAuth {
			if identity.JWT = reAuthJWT(*identity); identity.JWT != "" {
				continue
			}
		}

		if resp.StatusCode != http.StatusOK {
			// TODO!!! be careful writing server_http handlers, http.StatusOK is the only success code accepted here

			if len(responseBody) > bodyLogLimit {
				responseBody = responseBody[:bodyLogLimit]
			}

			var data common.Map
			if err = json.Unmarshal(responseBody, &data); err != nil {
				if len(responseBody) > bodyLogLimit {
					responseBody = responseBody[:bodyLogLimit]
				}
				return errors.Wrapf(err, "can't unmarshal body from %s %s: status = %d, body = %s", method, serverURL, resp.StatusCode, responseBody)
			}

			errorKey := common.ErrorKey(data.StringDefault(server.ErrorKey, ""))
			return common.KeyableError(errorKey, data, errors.Errorf("can't %s %s: status = %d, body = %s", method, serverURL, resp.StatusCode, responseBody))
		}

		if dataBytes, ok := responseData.(*[]byte); ok {
			*dataBytes = responseBody
		} else if dataBytes, ok := responseData.(*string); ok {
			*dataBytes = string(responseBody)
		} else if responseBinary, ok := responseData.(*ResponseBinary); ok {
			responseBinary.MIMEType = resp.Header.Get("Content-Type")
			responseBinary.Data = responseBody
		} else if responseData != nil {
			if err = json.Unmarshal(responseBody, responseData); err != nil {
				if len(responseBody) > bodyLogLimit {
					responseBody = responseBody[:bodyLogLimit]
				}
				return errors.Wrapf(err, "can't unmarshal body from %s %s: %s", method, serverURL, responseBody)
			}
		}

		break // end of each try means the end of all tries if something other wasn't managed before
	}

	return nil
}

func reAuthJWT(identity auth.Identity) string {
	authOp, _ := identity.ReAuthData[ReAuthOpKey].(auth.Operator)
	nickname := identity.ReAuthData.StringDefault(auth.CredsNickname, "")
	password := identity.ReAuthData.StringDefault(auth.CredsPassword, "")

	if authOp == nil { // || nickname == "" || password == ""
		return ""
	}

	creds := auth.Creds{auth.CredsNickname: nickname, auth.CredsPassword: password}
	identityNew, err := authOp.Authenticate(creds)
	if err != nil || identityNew == nil {
		// TODO: do it prettily
		log.Printf("on authOp.Authenticate(%#v): got %#v / %s", creds, identityNew, err)
		return ""
	}

	return identityNew.JWT + identity.ReAuthData.StringDefault(ReAuthSuffix, "")
}

//TRIES_ON_OVERLOAD:
//	for n := 1; n <= maxTriesOnOverload; n++ {
//
//		if statusCode == http.StatusTooManyRequests {
//			time.Sleep(delayOnOverload)
//			continue TRIES_ON_OVERLOAD
//		}
//	}