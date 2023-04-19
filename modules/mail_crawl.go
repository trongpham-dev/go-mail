package mailcrawl

import (
	"go-mail/modules/amazon"
	"go-mail/modules/etsy"
	"log"
	"regexp"

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

func Crawl(c *client.Client, ids []uint32) error {
	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)

		// Get the whole message body
		var section imap.BodySectionName
		items := []imap.FetchItem{section.FetchItem()}

		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)

		go func() {
			done <- c.Fetch(seqset, items, messages)
		}()

		if err := <-done; err != nil {
			log.Fatal(err)
			return err
		}

		log.Println("Unseen messages:")

		for msg := range messages {

			if msg == nil {
				log.Fatal("Server didn't returned message")
			}

			r := msg.GetBody(&section)
			if r == nil {
				log.Fatal("Server didn't returned message body")
			}

			// Create a new mail reader
			mr, err := mail.CreateReader(r)

			if err != nil {
				log.Fatal(err)
				return err
			}

			header := mr.Header
			from, err := header.AddressList("From")

			if err != nil {
				log.Fatal(err)
				return err
			}

			if contains(from, "uylongpham2910@gmail.com") {
				etsy := etsy.NewEtsy()
				etsy.CrawlEtsy(mr)
			} else {
				amazon.CrawlAmazonn(mr)
			}
		}
	}
	return nil
}

func FindOrderInfor(t string) *amazonOrderInfo {
	var rs = amazonOrderInfo{}

	//extracting Ship By
	pattern := regexp.MustCompile("Ship by:\\s+([\\d/]+)")
	rs.ShipBy = pattern.FindString(t)
	rs.ShipBy = rs.ShipBy[8:len(rs.ShipBy)]

	//extracting item name
	pattern = regexp.MustCompile("Item:\\s+(\\S+)")
	rs.Item = pattern.FindString(t)

	//extracting condition
	pattern = regexp.MustCompile("Condition:\\s+(\\S+)")
	rs.Condition = pattern.FindString(t)

	//extracting SKU
	pattern = regexp.MustCompile("SKU:\\s+(\\S+)")
	rs.SKU = pattern.FindString(t)

	//extracting quantity
	pattern = regexp.MustCompile("Quantity:\\s+(\\d+)")
	rs.Quantity = pattern.FindString(t)

	//extracting orderdate
	pattern = regexp.MustCompile("Order date:\\s+([\\d/]+)")
	rs.OrderDate = pattern.FindString(t)

	//extracting price
	pattern = regexp.MustCompile("Price:\\s+\\$(\\d+\\.\\d+)")
	rs.Price = pattern.FindString(t)

	//extracting Tax
	pattern = regexp.MustCompile("Tax:\\s+\\$(\\d+\\.\\d+)")
	rs.Tax = pattern.FindString(t)

	//extracting Shipping
	pattern = regexp.MustCompile("Shipping:\\s+\\$(\\d+\\.\\d+)")
	rs.Shipping = pattern.FindString(t)

	//extracting Promotion
	pattern = regexp.MustCompile("Promotions:\\s+-\\$(\\d+\\.\\d+)")
	rs.Promotions = pattern.FindString(t)

	//extracting Amazon fee
	pattern = regexp.MustCompile("Amazon fees:\\s+-\\$(\\d+\\.\\d+)")
	rs.AmazonFee = pattern.FindString(t)

	//extracting Marketplace Facilitator Tax
	pattern = regexp.MustCompile("Marketplace Facilitator Tax:\\s+-\\$(\\d+\\.\\d+)")
	rs.MarketPlaceFacilitatorTax = pattern.FindString(t)

	//extracting Your earnings
	pattern = regexp.MustCompile("Your earnings:\\s+\\$(\\d+\\.\\d+)")
	rs.YourEarning = pattern.FindString(t)

	return &rs
}
