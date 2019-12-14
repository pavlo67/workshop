package importer_tasks

import (
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/scheduler"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/selectors/logic"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/importer/importer_rss"
)

func NewLoader(dataOp data.Operator) (scheduler.Task, error) {

	if dataOp == nil {
		return nil, errors.New("on importer_task.NewLoader(): data.Operator == nil")
	}

	return &loadTask{dataOp}, nil
}

var _ scheduler.Task = &loadTask{}

type loadTask struct {
	dataOp data.Operator
}

func (it *loadTask) Name() string {
	return "loader"
}

func (it *loadTask) Run(timeSheduled time.Time) error {
	if it == nil {
		return errors.New("on importer_task.Run(): loadTask == nil")
	}

	return LoadAll(it.dataOp)
}

func LoadAll(dataOp data.Operator) error {
	urls := []string{
		"https://rss.unian.net/site/news_ukr.rss",
		"https://censor.net.ua/includes/news_uk.xml",
		"http://texty.org.ua/mod/news/?view=rss",
		"http://texty.org.ua/mod/article/?view=rss&ed=1",
		"http://texty.org.ua/mod/blog/blog_list.php?view=rss",
		"https://www.pravda.com.ua/rss/",
		"http://k.img.com.ua/rss/ua/all_news2.0.xml",
		"https://www.obozrevatel.com/rss.xml",
		"https://lenta.ru/rss",
		"https://www.gazeta.ru/export/rss/first.xml",
		"https://www.gazeta.ru/export/rss/lenta.xml",
	}

	for i, url := range urls {
		l.Infof("%d f %d: %s", i+1, len(urls), url)

		numAll, numProcessed, numNew, err := Load(url, dataOp)
		l.Infof("numAll = %d, numProcessed = %d, numNew = %d", numAll, numProcessed, numNew)

		if err != nil {
			l.Error(err)
		}

	}

	return nil
}

func Load(url string, dataOp data.Operator) (int, int, int, error) {
	impOp, err := importer_rss.NewRSS(l)
	if err != nil {
		return 0, 0, 0, errors.Errorf("can't importer_rss.NewRSS(%#v): %s", l, err)
	}

	series, err := impOp.Get(url)
	if err != nil {
		return 0, 0, 0, errors.Errorf("can't impOp.Get('%s', nil): %s", url, err)
	}
	if series == nil {
		return 0, 0, 0, errors.Errorf("no series from impOp.Get('%s', nil)", url)
	}

	numAll := len(series.Data)
	var numProcessed, numNew int

	for _, item := range series.Data {
		var cnt uint64

		numProcessed++

		term := logic.AND(
			selectors.Binary(selectors.Eq, "source", selectors.Value{item.Origin.Source}),
			selectors.Binary(selectors.Eq, "source_key", selectors.Value{item.Origin.Key}),
		)

		//itemStr, _ := json.Marshal(item)
		//l.Infof("%s ", itemStr)

		//termStr, _ := json.Marshal(term)
		//l.Infof("%s", termStr)

		cnt, err = dataOp.Count(term, nil)
		if err != nil {
			err = errors.Errorf("can't dataOp.Count(%#v): %s", term, err)
			break

		} else if cnt > 0 {
			// already exists!
			continue
		}

		item.ID = ""

		_, err = dataOp.Save([]data.Item{item}, nil)
		if err != nil {
			err = errors.Errorf("can't adminOp.Save(%#v): %s", item, err)
			break

		} else {
			numNew++
		}
	}

	return numAll, numProcessed, numNew, err
}
