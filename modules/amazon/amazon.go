package amazon

import (
	"io"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

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

func FindUnseenMail(c *client.Client) []uint32 {
	_, err := c.Select("INBOX", false)

	if err != nil {
		log.Fatal(err)
		return nil
	}
	
	// Set search criteria
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	criteria.Text = []string{"Congratulations! You just sold an item on Amazon!"}
	ids, err := c.Search(criteria)
	
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return ids
}

func Crawl(c *client.Client, ids []uint32) error{
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
			
			i := 0
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
					return err
				}

				switch h := p.Header.(type) {
					case *mail.InlineHeader:
						// This is the message's text (can be plain-text or HTML)
						b, _ := ioutil.ReadAll(p.Body)
						text := string(b)
						if i == 0 {
							rs := FindOrderInfor(text)
							log.Println("=====================================")
							log.Println(rs.ShipBy)	
							log.Println(rs.Item)	
							log.Println(rs.Condition)
							log.Println(rs.SKU)
							log.Println(rs.Quantity)
							log.Println(rs.OrderDate)
							log.Println(rs.Price)
							log.Println(rs.Tax)
							log.Println(rs.Promotions)
							log.Println(rs.AmazonFee)
							log.Println(rs.MarketPlaceFacilitatorTax)
							log.Println(rs.YourEarning)
							log.Println("=====================================")
							i++
						}		
					case *mail.AttachmentHeader:
						// This is an attachment
						filename, _ := h.Filename()
						log.Println("Got attachment: %v", filename)
				}

			}
		}
	}
	return nil
}


func FindOrderInfor(t string) *amazonOrderInfo {
	var rs = amazonOrderInfo{}

	//extracting Ship By
	pattern := regexp.MustCompile("Ship by:\\s+([\\d/]+)")
	rs.ShipBy =  pattern.FindString(t)

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