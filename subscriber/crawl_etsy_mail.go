package subscriber

import (
	"context"
	"go-mail/component"
	mailcrawl "go-mail/modules"
	"go-mail/pubsub"
	"log"
)

func RunCrawlEtsyMailData(appCtx component.AppContext) consumerJob {
	return consumerJob{
		Title: "Crawling Etsy Mail Data",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			
			log.Println(message.Data().Ids)
			return mailcrawl.Crawl(message.Data().Client, message.Data().Ids)
		},
	}
}