package data_sqlite

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/tagger"
)

func Starter() starter.Operator {
	return &dataSQLiteStarter{}
}

var l logger.Operator
var _ starter.Operator = &dataSQLiteStarter{}

type dataSQLiteStarter struct {
	config       config.Access
	interfaceKey joiner.InterfaceKey
}

func (ts *dataSQLiteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ts *dataSQLiteStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Options) ([]common.Options, error) {
	l = lCommon

	var cfgSQLite config.Access
	err := cfg.Value("sqlite", &cfgSQLite)
	if err != nil {
		return nil, err
	}

	ts.config = cfgSQLite
	ts.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(data.InterfaceKey)))

	// sqllib.CheckTables

	return nil, nil
}

func (ts *dataSQLiteStarter) Setup() error {
	return nil

	//return sqllib.SetupTables(
	//	sm.mysqlConfig,
	//	sm.index.MySQL,
	//	[]config.Table{{Key: "table", Title: sm.table}},
	//)
}

func (ts *dataSQLiteStarter) Run(joinerOp joiner.Operator) error {
	taggerOp, ok := joinerOp.Interface(tagger.InterfaceKey).(tagger.Operator)
	if !ok {
		return errors.Errorf("no tagger.Operator with key %s", tagger.InterfaceKey)
	}

	dataOp, _, err := NewData(ts.config, "", taggerOp, ts.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't init data.Operator")
	}

	err = joinerOp.Join(dataOp, ts.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *dataSQLite as data.Operator with key '%s'", ts.interfaceKey)
	}

	return nil
}
