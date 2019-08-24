package server_http_jschmhr

import (
	"fmt"

	"github.com/pkg/errors"

	"log"

	"github.com/pavlo67/workshop/basis/auth"
	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/config"
	"github.com/pavlo67/workshop/basis/joiner"
	"github.com/pavlo67/workshop/basis/logger"
	"github.com/pavlo67/workshop/basis/server/server_http"
	"github.com/pavlo67/workshop/basis/starter"
)

func Starter() starter.Operator {
	return &server_http_jschmhrStarter{}
}

var l logger.Operator
var _ starter.Operator = &server_http_jschmhrStarter{}

type server_http_jschmhrStarter struct {
	interfaceKey joiner.InterfaceKey
	// interfaceKeyRouter joiner.InterfaceKey
	config config.ServerTLS

	staticPaths map[string]string
}

func (ss *server_http_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *server_http_jschmhrStarter) Init(conf *config.Config, options common.Info) (info []common.Info, err error) {
	var errs common.Errors
	l = conf.Logger

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))
	// ss.interfaceKeyRouter = joiner.InterfaceKey(options.StringDefault("interface_key_router", string(controller.InterfaceKey)))

	ss.config = conf.Server
	if ss.config.Port <= 0 {
		errs = append(errs, fmt.Errorf("wrong port for serverOp: %#v", ss.config))
	}

	// TODO: use more then one static path
	if staticPath, ok := options.String("static_path"); ok {
		ss.staticPaths = map[string]string{"static": staticPath}
	}

	return nil, errs.Err()
}

func (ss *server_http_jschmhrStarter) Setup() error {
	return nil
}

func (ss *server_http_jschmhrStarter) Run(joinerOp joiner.Operator) error {

	authOp, ok := joinerOp.Interface(auth.InterfaceKey).(auth.Operator)
	if !ok {
		log.Fatalf("no auth.Operator with key %s", auth.InterfaceKey)
	}

	srvOp, err := New(ss.config.Port, ss.config.TLSCertFile, ss.config.TLSKeyFile, authOp)
	if err != nil {
		return errors.Wrap(err, "can't init serverHTTPJschmhr.Operator")
	}

	for path, staticPath := range ss.staticPaths {
		srvOp.HandleFiles("/"+path+"/*filepath", staticPath, nil)
	}

	err = joinerOp.Join(srvOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join serverHTTPJschmhr srvOp as server.Operator with key '%s'", ss.interfaceKey)
	}

	return nil
}