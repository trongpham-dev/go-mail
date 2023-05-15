package mailcrawl

import (
	"go-mail/component"
	"go-mail/modules/etsy"
	"log"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

type mailCrawl struct {
}

func NewMailCrawl() *mailCrawl {
	return &mailCrawl{}
}

func (m *mailCrawl) MailConnection(email, password string) (*client.Client, error) {
	log.Println(email, " is connecting to server!")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Login
	if err := c.Login(email, password); err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println(email, "logged in!")

	return c, nil
}

func (m *mailCrawl) FindUnseenMail(c *client.Client) []uint32 {
	_, err := c.Select("INBOX", false)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Set search criteria
	criteria1 := imap.NewSearchCriteria()
	criteria2 := imap.NewSearchCriteria()
	criteria3 := imap.NewSearchCriteria()

	criteria1.WithoutFlags = []string{imap.SeenFlag}
	criteria1.Text = []string{"Congratulations! You just sold an item on Amazon!"}

	criteria2.WithoutFlags = []string{imap.SeenFlag}
	criteria2.Text = []string{"Congratulations on your Etsy"}
	criteria3.Or = [][2]*imap.SearchCriteria{
		{criteria1, criteria2},
	}

	ids, err := c.Search(criteria3)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return ids
}

type amazonOrderInfo struct {
	ShipBy                    string `json:ship_by`
	Item                      string `json:item`
	Condition                 string `json:condition`
	SKU                       string `json:sku`
	Quantity                  string `json:quantity`
	OrderDate                 string `json:order_date`
	Price                     string `json:price`
	Tax                       string `json:tax`
	Shipping                  string `json:shipping`
	Promotions                string `json:promotions`
	AmazonFee                 string `json:amazon_fee`
	MarketPlaceFacilitatorTax string `json:marketplace_facilitator_tax`
	YourEarning               string `json:your_earning`
}

func contains(s []*mail.Address, str string) bool {
	for _, v := range s {
		log.Println(v.Address)
		if v.Address == str {

			return true
		}
	}

	return false
}

func (m *mailCrawl) Crawl(appCtx component.AppContext, c *client.Client, ids []uint32, mailTo string) error {
	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)

		// Get the whole message body
		var section imap.BodySectionName
		items := []imap.FetchItem{section.FetchItem()}

		messages := make(chan *imap.Message, 10000)
		done := make(chan error, 1)

		go func() {
			done <- c.Fetch(seqset, items, messages)
		}()

		if err := <-done; err != nil {
			return err
		}

		log.Println("Unseen messages:")

		for msg := range messages {
			if msg == nil {
				markMailAsUnseen(c, msg.SeqNum)
			}

			r := msg.GetBody(&section)
			if r == nil {
				markMailAsUnseen(c, msg.SeqNum)
			}

			// Create a new mail reader
			mr, err := mail.CreateReader(r)

			if err != nil {
				markMailAsUnseen(c, msg.SeqNum)
				return err
			}

			header := mr.Header
			from, err := header.AddressList("From")

			if err != nil {
				markMailAsUnseen(c, msg.SeqNum)
				return err
			}

			recivedAt, err := header.Date()

			if err != nil {
				markMailAsUnseen(c, msg.SeqNum)
				return err
			}

			// Load American timezone
			loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
			if err != nil {
				markMailAsUnseen(c, msg.SeqNum)
				return err
			}

			recivedAt = recivedAt.In(loc)

			// log.Println(recivedAt.Format("2006/01/02 15:04"))
			// log.Println(recivedAt.Unix())
			// log.Println(recivedAt)

			if contains(from, "transaction@etsy.com") {
				etsy := etsy.NewEtsy()
				if err = etsy.CrawlEtsy(appCtx, mr, mailTo, recivedAt.Format("2006/01/02 15:04")); err != nil {
					markMailAsUnseen(c, msg.SeqNum)
					return err
				}

			}
		}
	}

	//remove all mails from cache.
	appCtx.GetCaching().Write(c.Mailbox().Name, []uint32{})
	return nil
}

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
