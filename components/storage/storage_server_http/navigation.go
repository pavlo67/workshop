package storage_server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
)

const interfaceKeyParamName = "key"
const tagLabelParamName = "tag"

var countTagsEndpoint = server_http.Endpoint{Method: "GET", QueryParams: []string{interfaceKeyParamName}, WorkerHTTP: CountTags}

func CountTags(_ *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	var interfaceKeyPtr *joiner.InterfaceKey
	if key := req.URL.Query().Get(interfaceKeyParamName); key != "" {
		interfaceKey := joiner.InterfaceKey(key)
		interfaceKeyPtr = &interfaceKey
	}

	counter, err := dataTaggedOp.CountTags(interfaceKeyPtr, nil)
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...CountTags (%#v): %s", req.URL.Query(), err))
	}

	return server.ResponseRESTOk(counter)
}

var listWithTagEndpoint = server_http.Endpoint{Method: "GET", QueryParams: []string{interfaceKeyParamName, tagLabelParamName}, WorkerHTTP: ListWithTag}

func ListWithTag(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	var interfaceKeyPtr *joiner.InterfaceKey
	if key := req.URL.Query().Get(interfaceKeyParamName); key != "" {
		interfaceKey := joiner.InterfaceKey(key)
		interfaceKeyPtr = &interfaceKey
	}

	tagLabel := req.URL.Query().Get(tagLabelParamName)

	items, err := dataTaggedOp.ListWithTag(interfaceKeyPtr, tagLabel, nil, nil)

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...ListWithTag (%#v): %s", req.URL.Query(), err))
	}

	return server.ResponseRESTOk(items)
}