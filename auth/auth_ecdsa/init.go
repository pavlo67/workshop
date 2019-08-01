package auth_ecdsa

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/associatio/auth"
	"github.com/pavlo67/associatio/basis"
	"github.com/pavlo67/associatio/starter"
	"github.com/pavlo67/associatio/starter/config"
	"github.com/pavlo67/associatio/starter/joiner"
	"github.com/pavlo67/associatio/starter/logger"
)

func Starter() starter.Operator {
	return &identity_btcStarter{}
}

var l logger.Operator
var _ starter.Operator = &identity_btcStarter{}

type identity_btcStarter struct {
	interfaceKey joiner.InterfaceKey
}

func (ss *identity_btcStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *identity_btcStarter) Prepare(conf *config.Config, options, runtimeOptions basis.Info) error {
	l = logger.Get()

	// var errs basis.Errors

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(auth.InterfaceKey)))

	return nil
}

func (ss *identity_btcStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (ss *identity_btcStarter) Setup() error {
	return nil
}

func (ss *identity_btcStarter) Init(joiner joiner.Operator) error {
	identOp, err := New(nil)
	if err != nil {
		return errors.Wrap(err, "can't init identity_ecdsa.Operator")
	}

	err = joiner.JoinInterface(identOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join identity_ecdsa identOp as identity.Operator with key '%s'", ss.interfaceKey)
	}

	return nil
}
