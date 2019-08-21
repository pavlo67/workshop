package news_starters

import (
	"github.com/pavlo67/constructor/components/common/starter"

	"github.com/pavlo67/constructor/applications/flow/news_leveldb"
	"github.com/pavlo67/constructor/components/auth/auth_ecdsa"
	"github.com/pavlo67/constructor/components/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/constructor/processor/founts/founts_leveldb"
	"github.com/pavlo67/constructor/processor/router_news"
)

func Starters(routerStarters ...starter.Starter) ([]starter.Starter, string) {
	var starters = []starter.Starter{
		// 1. operational interfaces
		{founts_leveldb.Starter(), nil},
		{news_leveldb.Starter(), nil},
	}

	return append(starters, routerStarters...), "NEWS BUILD"
}

func RouterHTTPStarters() []starter.Starter {
	return []starter.Starter{
		{auth_ecdsa.Starter(), nil},
		{server_http_jschmhr.Starter(), nil},
		{router_news.Starter(), nil},
	}
}
