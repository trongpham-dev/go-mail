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

func CrawlAmazonn(mr *mail.Reader){
	i := 0
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		switch h := p.Header.(type) {
			case *mail.InlineHeader:
				// This is the message's text (can be plain-text or HTML)
				b, _ := ioutil.ReadAll(p.Body)
				text := string(b)
				if i == 0 {
					rs := FindOrderInfor(text)
					log.Println("=====================================")
					log.Println(rs.OrderId)	
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

// func Crawl(c *client.Client, ids []uint32) error{
// 	if len(ids) > 0 {
// 		seqset := new(imap.SeqSet)
// 		seqset.AddNum(ids...)

// 		// Get the whole message body
// 		var section imap.BodySectionName
// 		items := []imap.FetchItem{section.FetchItem()}

// 		messages := make(chan *imap.Message, 10)
// 		done := make(chan error, 1)

// 		go func() {
// 			done <- c.Fetch(seqset, items, messages)
// 		}()

// 		if err := <-done; err != nil {
// 			log.Fatal(err)
// 			return err
// 		}

// 		log.Println("Unseen messages:")

// 		for msg := range messages {

// 			if msg == nil {
// 				log.Fatal("Server didn't returned message")
// 			}
	
// 			r := msg.GetBody(&section)
// 			if r == nil {
// 				log.Fatal("Server didn't returned message body")
// 			}
			
// 			// Create a new mail reader
// 			mr, err := mail.CreateReader(r)
	
// 			if err != nil {
// 				log.Fatal(err)
// 				return err
// 			}
			
// 			i := 0
// 			for {
// 				p, err := mr.NextPart()
// 				if err == io.EOF {
// 					break
// 				} else if err != nil {
// 					log.Fatal(err)
// 					return err
// 				}

// 				switch h := p.Header.(type) {
// 					case *mail.InlineHeader:
// 						// This is the message's text (can be plain-text or HTML)
// 						b, _ := ioutil.ReadAll(p.Body)
// 						text := string(b)
// 						if i == 0 {
// 							rs := FindOrderInfor(text)
// 							log.Println("=====================================")
// 							log.Println(rs.OrderId)	
// 							log.Println(rs.ShipBy)	
// 							log.Println(rs.Item)	
// 							log.Println(rs.Condition)
// 							log.Println(rs.SKU)
// 							log.Println(rs.Quantity)
// 							log.Println(rs.OrderDate)
// 							log.Println(rs.Price)
// 							log.Println(rs.Tax)
// 							log.Println(rs.Promotions)
// 							log.Println(rs.AmazonFee)
// 							log.Println(rs.MarketPlaceFacilitatorTax)
// 							log.Println(rs.YourEarning)
// 							log.Println("=====================================")
// 							i++
// 						}		
// 					case *mail.AttachmentHeader:
// 						// This is an attachment
// 						filename, _ := h.Filename()
// 						log.Println("Got attachment: %v", filename)
// 				}

// 			}
// 		}
// 	}
// 	return nil
// }


func FindOrderInfor(t string) *amazonOrderInfo {
	var rs = amazonOrderInfo{}

	//extracting OrderId
	pattern := regexp.MustCompile("Order ID:\\s+(\\S+)")
	rs.OrderId =  pattern.FindString(t)
	rs.OrderId = rs.OrderId[10:len(rs.OrderId)]

	//extracting Ship By
	pattern = regexp.MustCompile("Ship by:\\s+([\\d/]+)")
	rs.ShipBy =  pattern.FindString(t)
	rs.ShipBy = rs.ShipBy[9:len(rs.ShipBy)]

	//extracting item name
	pattern = regexp.MustCompile("Item:\\s+(\\S+)")
	rs.Item = pattern.FindString(t)
	rs.Item = rs.Item[6:len(rs.Item)]

	//extracting condition
	pattern = regexp.MustCompile("Condition:\\s+(\\S+)")
	rs.Condition = pattern.FindString(t)
	rs.Condition = rs.Condition[11:len(rs.Condition)]

	//extracting SKU
	pattern = regexp.MustCompile("SKU:\\s+(\\S+)")
	rs.SKU = pattern.FindString(t)
	rs.SKU = rs.SKU[5:len(rs.SKU)]

	//extracting quantity
	pattern = regexp.MustCompile("Quantity:\\s+(\\d+)")
	rs.Quantity = pattern.FindString(t)
	rs.Quantity = rs.Quantity[10:len(rs.Quantity)]

	//extracting orderdate
	pattern = regexp.MustCompile("Order date:\\s+([\\d/]+)")
	rs.OrderDate = pattern.FindString(t)
	rs.OrderDate = rs.OrderDate[12:len(rs.OrderDate)]

	//extracting price
	pattern = regexp.MustCompile("Price:\\s+\\$(\\d+\\.\\d+)")
	rs.Price = pattern.FindString(t)
	rs.Price = rs.Price[7:len(rs.Price)]

	//extracting Tax
	pattern = regexp.MustCompile("Tax:\\s+\\$(\\d+\\.\\d+)")
	rs.Tax = pattern.FindString(t)
	rs.Tax = rs.Tax[5:len(rs.Tax)]

	//extracting Shipping
	pattern = regexp.MustCompile("Shipping:\\s+\\$(\\d+\\.\\d+)")
	rs.Shipping = pattern.FindString(t)
	rs.Shipping = rs.Shipping[10:len(rs.Shipping)]

	//extracting Promotion
	pattern = regexp.MustCompile("Promotions:\\s+-\\$(\\d+\\.\\d+)")
	rs.Promotions = pattern.FindString(t)
	rs.Promotions = rs.Promotions[12:len(rs.Promotions)]

	//extracting Amazon fee
	pattern = regexp.MustCompile("Amazon fees:\\s+-\\$(\\d+\\.\\d+)")
	rs.AmazonFee = pattern.FindString(t)
	rs.AmazonFee = rs.AmazonFee[13:len(rs.AmazonFee)]

	//extracting Marketplace Facilitator Tax
	pattern = regexp.MustCompile("Marketplace Facilitator Tax:\\s+-\\$(\\d+\\.\\d+)")
	rs.MarketPlaceFacilitatorTax = pattern.FindString(t)
	rs.MarketPlaceFacilitatorTax = rs.MarketPlaceFacilitatorTax[29:len(rs.MarketPlaceFacilitatorTax)]

	//extracting Your earnings
	pattern = regexp.MustCompile("Your earnings:\\s+\\$(\\d+\\.\\d+)")
	rs.YourEarning = pattern.FindString(t)
	rs.YourEarning = rs.YourEarning[15:len(rs.YourEarning)]

	return &rs
}