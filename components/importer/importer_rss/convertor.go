package importer_rss

import (
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/importer"
	"github.com/pavlo67/workshop/components/tagger"
)

var _ data.Convertor = &Item{}

type Item struct {
	sourceTime time.Time
	sourceURL  string
	feedItem   *gofeed.Item
}

func (item *Item) GetData() (*data.Item, error) {
	if item == nil || item.feedItem == nil {
		return nil, nil
	}

	feedItem := item.feedItem

	originalID := feedItem.GUID
	if originalID == "" {
		originalID = feedItem.Link
	}

	sourceTime := item.sourceTime
	if feedItem.PublishedParsed != nil {
		sourceTime = *feedItem.PublishedParsed
	}

	var embedded []data.Item
	if strings.TrimSpace(feedItem.Link) != "" {
		embedded = append(embedded, data.Item{
			URL: importer.ImportedHREF + feedItem.Link,
		})
	}

	if feedItem.Image != nil {
		embedded = append(embedded, data.Item{
			URL:   importer.ImportedHREF + feedItem.Image.URL,
			Title: feedItem.Image.Title,
		})
	}

	if len(feedItem.Enclosures) > 0 {
		for _, p := range feedItem.Enclosures {
			embedded = append(embedded, data.Item{
				URL:   importer.ImportedHREF + p.URL,
				Title: p.Type + ": " + p.Length,
			})
		}
	}

	var tags []tagger.Tag
	for _, c := range feedItem.Categories {
		tags = append(tags, tagger.Tag(c))
	}

	return &data.Item{
		Title:    feedItem.Title,
		Summary:  feedItem.Description,
		Embedded: embedded,
		Tags:     tags,
		Details:  feedItem.Content,
		Status:   crud.Status{CreatedAt: item.sourceTime},
		Origin: flow.Origin{
			Source: item.sourceURL,
			Key:    originalID,
			Time:   &sourceTime,
			Data:   &item,
		},
	}, nil

}

func (*Item) SaveData(data.Item) error {
	return common.ErrNotImplemented
}

// language := feed.Language

// type Census struct {
// 	Label        string                `json:"title,omitempty"`
// 	Description    string                `json:"description,omitempty"`
// 	Contentus        string                `json:"content,omitempty"`
// 	Link        string                `json:"link,omitempty"`
// 	Updated        string                `json:"updated,omitempty"`
// 	UpdatedParsed    *time.Time            `json:"updatedParsed,omitempty"`
// 	Published    string                `json:"published,omitempty"`
// 	PublishedParsed    *time.Time            `json:"publishedParsed,omitempty"`
// 	Author        *Person                `json:"author,omitempty"`
// 	GUID        string                `json:"guid,omitempty"`
// 	Image        *Image                `json:"image,omitempty"`
// 	Categories    []string            `json:"categories,omitempty"`
// 	Enclosures    []*Enclosure            `json:"enclosures,omitempty"`
// 	DublinCoreExt    *ext.DublinCoreExtension    `json:"dcExt,omitempty"`
// 	ITunesExt    *ext.ITunesItemExtension    `json:"itunesExt,omitempty"`
// 	Extensions    ext.Extensions            `json:"extensions,omitempty"`
// 	Custom        map[string]string        `json:"custom,omitempty"`
// }