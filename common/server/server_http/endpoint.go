package server_http

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
)

type Endpoint struct {
	Method     string   `json:"method,omitempty"`
	Path       string   `json:"path"`
	ParamNames []string `json:"param_names,omitempty"`
	// AllowedIDs []common.ID `json:"allowed_ids,omitempty"`

	WorkerHTTP

	// DataItem   interface{} `json:"data_item,omitempty"` // for Interface

	SwaggerDescription string
}

func InitEndpoint(epPtrs *[]Endpoint, method, path string, paramNames []string, workerHTTP WorkerHTTP, swaggerDescription string, l logger.Operator) int {
	if epPtrs == nil {
		l.Errorf("no epPtrs for InitEndpoint() in %s", filelib.CurrentPath())
		return 1
	}

	*epPtrs = append(*epPtrs, Endpoint{
		Method:             method,
		Path:               path,
		ParamNames:         paramNames,
		WorkerHTTP:         workerHTTP,
		SwaggerDescription: swaggerDescription,
	})

	return 0
}

var rePathParam = regexp.MustCompile(":[^/]+")

func (ep Endpoint) PathWithParams(params ...string) string {
	matches := rePathParam.FindAllStringSubmatchIndex(ep.Path, -1)

	numMatches := len(matches)
	if len(params) < numMatches {
		numMatches = len(params)
	}

	path := ep.Path
	for nm := numMatches - 1; nm >= 0; nm-- {
		path = path[:matches[nm][0]] + url.PathEscape(strings.Replace(params[nm], "/", "%2F", -1)) + path[matches[nm][1]:]
	}

	return path
}

func (ep Endpoint) PathTemplate() string {
	path := ep.Path
	if len(path) == 0 || path[0] != '/' {
		path = "/" + path
	}

	if len(ep.ParamNames) < 1 {
		return path
	}

	return path + "/:" + strings.Join(ep.ParamNames, "/:")
}
