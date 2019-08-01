package demo_starters

import (
	"github.com/pavlo67/associatio/auth/auth_ecdsa"
	"github.com/pavlo67/associatio/basis"
	"github.com/pavlo67/associatio/basis/filelib"
	"github.com/pavlo67/associatio/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/associatio/starter"

	"github.com/pavlo67/associatio/_demo_main/demo_server_http"
)

func Starters() ([]starter.Starter, string) {
	paramsServerStatic := basis.Info{
		"static_path":   filelib.CurrentPath() + "../demo_server_http/static/",
		"template_path": filelib.CurrentPath() + "../demo_server_http/static/demo_server.html",
	}

	var starters []starter.Starter

	starters = append(starters, starter.Starter{auth_ecdsa.Starter(), nil})
	starters = append(starters, starter.Starter{server_http_jschmhr.Starter(), paramsServerStatic})
	starters = append(starters, starter.Starter{demo_server_http.Starter(), nil})

	return starters, "PUNCTUM DEMO BUILD"
}
