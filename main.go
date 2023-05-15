package main

import (
	"context"
	"go-mail/common"
	"go-mail/component"
	"go-mail/memcache"
	mailcrawler "go-mail/modules"
	"go-mail/pubsub"
	"go-mail/pubsub/pblocal"
	"go-mail/subscriber"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
		return
	}

	m := mailcrawler.NewMailCrawl()

	email1 := os.Getenv("MAIL_USER")
	password1 := os.Getenv("MAIL_PASSWORD")

	email2 := os.Getenv("MAIL_USER2")
	password2 := os.Getenv("MAIL_PASSWORD2")

	email3 := os.Getenv("MAIL_USER3")
	password3 := os.Getenv("MAIL_PASSWORD3")

	email4 := os.Getenv("MAIL_USER4")
	password4 := os.Getenv("MAIL_PASSWORD4")

	email5 := os.Getenv("MAIL_USER5")
	password5 := os.Getenv("MAIL_PASSWORD5")

	email6 := os.Getenv("MAIL_USER6")
	password6 := os.Getenv("MAIL_PASSWORD6")

	email7 := os.Getenv("MAIL_USER7")
	password7 := os.Getenv("MAIL_PASSWORD7")

	c1, err := m.MailConnection(email1, password1)

	if err != nil {
		log.Fatal(err)
	}

	c2, err := m.MailConnection(email2, password2)

	if err != nil {
		log.Fatal(err)
	}
	c3, err := m.MailConnection(email3, password3)

	if err != nil {
		log.Fatal(err)
	}

	c4, err := m.MailConnection(email4, password4)

	if err != nil {
		log.Fatal(err)
	}

	c5, err := m.MailConnection(email5, password5)

	if err != nil {
		log.Fatal(err)
	}

	c6, err := m.MailConnection(email6, password6)

	if err != nil {
		log.Fatal(err)
	}

	c7, err := m.MailConnection(email7, password7)

	if err != nil {
		log.Fatal(err)
	}

	appCtx := component.NewAppContext(pblocal.NewPubSub(), memcache.NewCaching(), memcache.NewAppToken())

	//subscriber.Setup(appCtx)
	if err := subscriber.NewEngine(appCtx).Start(); err != nil {
		log.Fatalln(err)
	}

	// //publish ids mail to subscriber
	var wg sync.WaitGroup
	for {
		wg.Add(7)
		go func() {
			defer wg.Done()
			err := c1.Noop()
			if err != nil {
				log.Println("Connection lost")
				c1, err = m.MailConnection(email1, password1)

				if err != nil {
					log.Println(err)
				}
			}
			ids := m.FindUnseenMail(c1)
			if len(ids) > 0 {
				if len(appCtx.GetCaching().Read(c1.Mailbox().Name)) > 0 {
					return
				}
				appCtx.GetCaching().Write(c1.Mailbox().Name, ids)
				appCtx.GetPubsub().Publish(context.Background(), common.TopicCrawlMail, pubsub.NewMessage(pubsub.MailData{Client: c1, Ids: ids, Mail: email1}))
			}
		}()

		go func() {
			defer wg.Done()
			err := c2.Noop()
			if err != nil {
				log.Println("Connection lost")
				c2, err = m.MailConnection(email2, password2)

				if err != nil {
					log.Println(err)
				}
			}
			ids2 := m.FindUnseenMail(c2)
			if len(ids2) > 0 {
				if len(appCtx.GetCaching().Read(c2.Mailbox().Name)) > 0 {
					return
				}
				appCtx.GetCaching().Write(c2.Mailbox().Name, ids2)
				appCtx.GetPubsub().Publish(context.Background(), common.TopicCrawlMail, pubsub.NewMessage(pubsub.MailData{Client: c2, Ids: ids2, Mail: email2}))
			}
		}()

		go func() {
			defer wg.Done()
			err := c3.Noop()
			if err != nil {
				log.Println("Connection lost")
				c3, err = m.MailConnection(email3, password3)

				if err != nil {
					log.Println(err)
				}
			}

			ids2 := m.FindUnseenMail(c3)
			if len(ids2) > 0 {
				if len(appCtx.GetCaching().Read(c3.Mailbox().Name)) > 0 {
					return
				}
				appCtx.GetCaching().Write(c3.Mailbox().Name, ids2)
				appCtx.GetPubsub().Publish(context.Background(), common.TopicCrawlMail, pubsub.NewMessage(pubsub.MailData{Client: c3, Ids: ids2, Mail: email3}))
			}
		}()

		go func() {
			defer wg.Done()
			err := c4.Noop()
			if err != nil {
				log.Println("Connection lost")
				c4, err = m.MailConnection(email4, password4)

				if err != nil {
					log.Println(err)
				}
			}
			ids2 := m.FindUnseenMail(c4)
			if len(ids2) > 0 {
				if len(appCtx.GetCaching().Read(c4.Mailbox().Name)) > 0 {
					return
				}
				appCtx.GetCaching().Write(c4.Mailbox().Name, ids2)
				appCtx.GetPubsub().Publish(context.Background(), common.TopicCrawlMail, pubsub.NewMessage(pubsub.MailData{Client: c4, Ids: ids2, Mail: email4}))
			}
		}()

		go func() {
			defer wg.Done()
			err := c5.Noop()
			if err != nil {
				log.Println("Connection lost")
				c5, err = m.MailConnection(email5, password5)

				if err != nil {
					log.Println(err)
				}
			}
			ids2 := m.FindUnseenMail(c5)
			if len(ids2) > 0 {
				if len(appCtx.GetCaching().Read(c5.Mailbox().Name)) > 0 {
					return
				}
				appCtx.GetCaching().Write(c5.Mailbox().Name, ids2)
				appCtx.GetPubsub().Publish(context.Background(), common.TopicCrawlMail, pubsub.NewMessage(pubsub.MailData{Client: c5, Ids: ids2, Mail: email5}))
			}
		}()

		go func() {
			defer wg.Done()
			err := c6.Noop()
			if err != nil {
				log.Println("Connection lost")
				c6, err = m.MailConnection(email6, password6)

				if err != nil {
					log.Println(err)
				}
			}
			ids2 := m.FindUnseenMail(c6)
			if len(ids2) > 0 {
				if len(appCtx.GetCaching().Read(c6.Mailbox().Name)) > 0 {
					return
				}
				appCtx.GetCaching().Write(c6.Mailbox().Name, ids2)
				appCtx.GetPubsub().Publish(context.Background(), common.TopicCrawlMail, pubsub.NewMessage(pubsub.MailData{Client: c6, Ids: ids2, Mail: email6}))
			}
		}()

		go func() {
			defer wg.Done()
			err := c7.Noop()
			if err != nil {
				log.Println("Connection lost")
				c7, err = m.MailConnection(email7, password7)

				if err != nil {
					log.Println(err)
				}
			}
			ids2 := m.FindUnseenMail(c7)
			if len(ids2) > 0 {
				if len(appCtx.GetCaching().Read(c7.Mailbox().Name)) > 0 {
					return
				}
				appCtx.GetCaching().Write(c7.Mailbox().Name, ids2)
				appCtx.GetPubsub().Publish(context.Background(), common.TopicCrawlMail, pubsub.NewMessage(pubsub.MailData{Client: c7, Ids: ids2, Mail: email7}))
			}
		}()

		<-time.After(1 * time.Minute)

		wg.Wait()
	}
}
