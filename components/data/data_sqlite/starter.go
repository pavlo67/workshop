package data_sqlite

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/tags"
)

func Starter() starter.Operator {
	return &dataSQLiteStarter{}
}

var l logger.Operator
var _ starter.Operator = &dataSQLiteStarter{}

type dataSQLiteStarter struct {
	config       config.Access
	table        string
	interfaceKey joiner.InterfaceKey

	noTagger bool
}

func (ts *dataSQLiteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ts *dataSQLiteStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	var cfgSQLite config.Access
	err := cfg.Value("sqlite", &cfgSQLite)
	if err != nil {
		return nil, err
	}

	ts.config = cfgSQLite
	ts.table, _ = options.String("table")
	ts.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(data.InterfaceKey)))

	ts.noTagger, _ = options.Bool("no_tagger")

	// sqllib.CheckTables

	return nil, nil
}

func (ts *dataSQLiteStarter) Setup() error {
	return nil

	//return sqllib.SetupTables(
	//	sm.mysqlConfig,
	//	sm.index.MySQL,
	//	[]config.Table{{ID: "table", Title: sm.table}},
	//)
}

func (ts *dataSQLiteStarter) Run(joinerOp joiner.Operator) error {
	var ok bool
	var taggerOp tags.Operator
	var cleanerOp crud.Cleaner

	if !ts.noTagger {
		taggerOp, ok = joinerOp.Interface(tags.InterfaceKey).(tags.Operator)
		if !ok {
			return errors.Errorf("no tagger.Operator with key %s", tags.InterfaceKey)
		}

		cleanerOp, ok = joinerOp.Interface(tags.CleanerInterfaceKey).(crud.Cleaner)
		if !ok {
			return errors.Errorf("no tagger.Cleaner with key %s", tags.InterfaceKey)
		}
	}

	dataOp, _, err := New(ts.config, ts.table, ts.interfaceKey, taggerOp, cleanerOp)
	if err != nil {
		return errors.Wrap(err, "can't init data.Operator")
	}

	err = joinerOp.Join(dataOp, ts.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *dataSQLite as data.Operator with key '%s'", ts.interfaceKey)
	}

	return nil
}
