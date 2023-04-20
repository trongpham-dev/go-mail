package subscriber

import (
	"context"
	"go-mail/component"
	mailcrawl "go-mail/modules"
	"go-mail/pubsub"
)

func RunCrawlEtsyMailData(appCtx component.AppContext) consumerJob {
	return consumerJob{
		Title: "Crawling Etsy Mail Data",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			m := mailcrawl.NewMailCrawl()
			return m.Crawl(appCtx, message.Data().Client, message.Data().Ids)
		},
	}
}
