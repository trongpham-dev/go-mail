package main

import (
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

func main() {
	// err := godotenv.Load(".env")

	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// m := mc.NewMailCrawl()

	// email1 := os.Getenv("MAIL_USER")
	// password1 := os.Getenv("MAIL_PASSWORD")

	// email2 := os.Getenv("MAIL_USER2")
	// password2 := os.Getenv("MAIL_PASSWORD2")

	// c1, err := m.MailConnection(email1, password1)

	// if err != nil{
	// 	log.Fatal(err)
	// }

	// c2, err := m.MailConnection(email2, password2)

	// if err != nil{
	// 	log.Fatal(err)
	// }

	// appCtx := component.NewAppContext(pblocal.NewPubSub())

	// //subscriber.Setup(appCtx)
	// if err := subscriber.NewEngine(appCtx).Start(); err != nil {
	// 	log.Fatalln(err)
	// }

	// //publish ids mail to subscriber
	// for{			
	// 		ids := am.FindUnseenMail(c1)
	// 		if len(ids) > 0{
	// 			appCtx.GetPubsub().Publish(context.Background(), common.TopicCrawlAmazonMail, pubsub.NewMessage(pubsub.MailData{Client:c1, Ids: ids}))
	// 		}
		
	// 		ids2 := am.FindUnseenMail(c2)
	// 		if len(ids2) > 0{ 
	// 			appCtx.GetPubsub().Publish(context.Background(), common.TopicCrawlAmazonMail, pubsub.NewMessage(pubsub.MailData{Client:c2, Ids: ids2}))
	// 		}
	// }
	start()
}

func start(){
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)

	if err != nil {
		log.Fatal(err)
		return
	}

	// Don't forget to logout
	defer c.Logout()

	err = godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
		return
	}

	user := os.Getenv("MAIL_USER")
	password := os.Getenv("MAIL_PASSWORD")

	// Login
	if err := c.Login(user, password); err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Logged in")

	// List mailboxes
	// mailboxes := make(chan *imap.MailboxInfo, 10)
	// done := make(chan error, 1)
	// go func() {
	// 	done <- c.List("", "*", mailboxes)
	// }()

	// log.Println("Mailboxes:")
	// for m := range mailboxes {
	// 	log.Println("* " + m.Name)
	// }

	// if err := <-done; err != nil {
	// 	log.Fatal(err)
	// }

	// // Select INBOX
	for {
	_, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	// Set search criteria
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	criteria.Text = []string{"Congratulations on your Etsy"}
	ids, err := c.Search(criteria)
	
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mails found:", ids)

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
			}

				// Print some info about the message
			// header := mr.Header
			// if date, err := header.Date(); err == nil {
			// 	log.Println("Date:", date)
			// }
			// if from, err := header.AddressList("From"); err == nil {
			// 	log.Println("From:", from)
			// }
			// if to, err := header.AddressList("To"); err == nil {
			// 	log.Println("To:", to)
			// }
			// if subject, err := header.Subject(); err == nil {
			// 	log.Println("Subject:", subject)
			// }

			// Process each message's part
			// i := 0
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
						// b, _ := ioutil.ReadAll(p.Body)

						// Load the HTML document
  						doc, err := goquery.NewDocumentFromReader(p.Body)

						if err != nil {
							panic(err)
						}

						doc.Find(`td[valign="top"][style="line-height:0px"]`).Each(func(i int, s *goquery.Selection) {

							sl := s.Find(`div[style*="font-family:arial,helvetica,sans-serif;"]`)
								for idx := range sl.Nodes{
									if sl.Eq(idx).Find(`a[style="text-decoration:none;color:#222222"]`).Text() != ""{
										break
									}

									rs := strings.Replace(strings.ReplaceAll(sl.Eq(idx).Text(), "\n", ""),"  ", "",-1)
										log.Println(rs)
								}
						})

						// fmt.Println(doc)
						
						// text := string(b)
						// text2, err := html2text.FromString(text, html2text.Options{PrettyTables: true})
						// if err != nil {
						// 	panic(err)
						// }
						// fmt.Println(text2)
						// if i == 0 {
						// 	rs := findOrderInfor(text)
						// 	log.Println("=====================================")
						// 	log.Println(rs.ShipBy)	
						// 	log.Println(rs.Item)	
						// 	log.Println(rs.Condition)
						// 	log.Println(rs.SKU)
						// 	log.Println(rs.Quantity)
						// 	log.Println(rs.OrderDate)
						// 	log.Println(rs.Price)
						// 	log.Println(rs.Tax)
						// 	log.Println(rs.Promotions)
						// 	log.Println(rs.AmazonFee)
						// 	log.Println(rs.MarketPlaceFacilitatorTax)
						// 	log.Println(rs.YourEarning)
						// 	log.Println("=====================================")
						// 	i++
						// }		
					case *mail.AttachmentHeader:
						// This is an attachment
						filename, _ := h.Filename()
						log.Println("Got attachment: %v", filename)
				}

			}
		}
	}
	}
}

type amazonOrderInfo struct{
	ShipBy string `json:ship_by`
	Item string `json:item`
	Condition string `json:condition`
	SKU string `json:sku`
	Quantity string `json:quantity`
	OrderDate string `json:order_date`
	Price string `json:price`
	Tax string `json:tax`
	Shipping string `json:shipping`
	Promotions string `json:promotions`
	AmazonFee string `json:amazon_fee`
	MarketPlaceFacilitatorTax string `json:marketplace_facilitator_tax`
	YourEarning string `json:your_earning`
}

func findOrderInfor(t string) *amazonOrderInfo {
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

func getSubstring(s string, indices []int) string {
    return string(s[indices[0]:indices[1]])
}
