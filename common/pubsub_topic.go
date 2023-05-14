package common

import (
	"go-mail/pubsub"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

const (
	TopicCrawlAmazonMail pubsub.Topic = "TopicCrawlAmazonMail"
	TopicCrawlEtsyMail   pubsub.Topic = "TopicCrawlEtsyMail"
	TopicCrawlMail       pubsub.Topic = "TopicCrawlMail"
)

func markMailAsUnseen(c *client.Client, uid uint32) error {
	// mark mail as unseen
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uid)

	item := imap.FormatFlagsOp(imap.RemoveFlags, true)
	flags := []interface{}{imap.SeenFlag}
	err := c.Store(seqSet, item, flags, nil)

	if err != nil {
		return err
	}

	return nil
}
