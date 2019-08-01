package server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/associatio/auth"
	"github.com/pavlo67/associatio/basis"
)

var errNoIdentityOpsMap = errors.New("no map[CredsType]identity.Operator")

func UserWithRequest(r *http.Request, identOpsMap map[auth.CredsType][]auth.Operator) (*auth.User, error) {
	if identOpsMap == nil {
		return nil, errNoIdentityOpsMap
	}

	var errs basis.Errors
	var user *auth.User

	// TOKEN CHECK
	token := r.Header.Get("Token")
	if token != "" {
		user, errs = auth.GetUser([]auth.Creds{{Type: auth.CredsToken, Value: token}}, identOpsMap, errs)
		if user != nil {
			return user, errs.Err()
		}
		// previous errs is added by auth.GetUser()
	}

	// COOKIE CHECK
	c, _ := r.Cookie("Token") // ErrNoCookie only
	if c != nil && c.Value != "" {
		user, errs = auth.GetUser([]auth.Creds{{Type: auth.CredsToken, Value: c.Value}}, identOpsMap, errs)
		if user != nil {
			return user, errs.Err()
		}
		// previous errs is added by auth.GetUser()
	}

	// SIGNATURE CHECK
	signature := r.Header.Get("Signature")
	if signature != "" && r.URL != nil {
		publicKeyAddress := r.Header.Get("Public-Key-Address")
		numberToSignature := r.Header.Get("Number-To-Signature")

		credsSignature := []auth.Creds{
			{Type: auth.CredsPublicKeyAddress, Value: publicKeyAddress},
			{Type: auth.CredsContentToSignature, Value: r.URL.Path + "?" + r.URL.RawQuery},
			{Type: auth.CredsNumberToSignature, Value: numberToSignature},
			{Type: auth.CredsSignature, Value: signature},
		}

		user, errs = auth.GetUser(credsSignature, identOpsMap, errs)
		// previous errs is added by auth.GetUser()
	}

	return user, errs.Err()
}
