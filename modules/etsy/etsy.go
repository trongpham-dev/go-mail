package etsy

import (
	"go-mail/component"
	"io"
	"log"

	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
)

type etSy struct {
	idx                    int
	orderId                string
	orderStr               strings.Builder
	orderDetailStr         strings.Builder
	arrEtsyOrder           []EtsyFieldOrder
	arrEtsyOrderDetail     []EtsyFieldOrderDetail
	arrEtsyOrderShipping   []EtsyFieldOrderShipping
	etsyOrder              EtsyOrder
	etsyFieldOrder         EtsyFieldOrder
	etsyOrderDetail        EtsyOrderDetail
	etsyFieldOrderDetail   EtsyFieldOrderDetail
	etsyOrderShipping      EtsyOrderShipping
	etsyFieldOrderShipping EtsyFieldOrderShipping
	count                  int
}

func NewEtsy() *etSy {
	return &etSy{}
}

var address = ""
var custName = ""
var road = ""
var city = ""
var state = ""
var zip = ""
var country = ""
var cusMail = ""
var shopName = ""

func (e *etSy) CrawlEtsy(appCtx component.AppContext, mr *mail.Reader, mailTo string, recievedAt string) error {
	e.count = 0
	index := 0
	var myErr error
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
			//log.Fatal(err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:

			if index > 0 {
				b, _ := io.ReadAll(p.Body)
				text := string(b)

				doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))

				//log.Println(text)

				if err != nil {
					return err
				}

				shopName = getShopName(mailTo)

				doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
					href, _ := item.Attr("href")
					if strings.Contains(href, "mailto:") {
						cusMail = strings.Replace(href, "mailto:", "", -1)
					}
				})

				doc.Find(`address[style="font-style: normal;"]`).Each(func(index int, s *goquery.Selection) {
					address = strings.TrimSpace(s.Text())
				})

				//rs := ""
				doc.Find(`div[style="font-family: arial, helvetica, sans-serif; color: #444444; font-size: 16px; line-height: 24px;"]`).Each(func(i int, s *goquery.Selection) {
					pattern := regexp.MustCompile("Your order number is:\\s+(\\S+)")
					match := pattern.FindStringSubmatch(s.Text())
					if len(match) > 0 {
						e.orderId = match[1]
						e.orderStr.WriteString(e.orderId + "\n")
					}
				})

				doc.Find(`td[class="right"][style="line-height: 0px"]`).Each(func(i int, s *goquery.Selection) {
					e.orderStr = strings.Builder{}
					e.etsyOrder = EtsyOrder{}
					i = 0
					if len(s.Nodes) > 0 {
						e.arrEtsyOrder = make([]EtsyFieldOrder, len(s.Nodes))

						sl := s.Find(`div[style*="font-family: arial, helvetica, sans-serif;"]`)
						for idx := range sl.Nodes {
							if sl.Eq(idx).Find(`a[style="text-decoration: none; color: #222222"]`).Text() != "" {
								return
							}

							if sl.Eq(idx).Find(`a[style="text-decoration: none; color: #444444;"]`).Text() != "" {
								// log.Println("Item:", strings.Replace(strings.ReplaceAll(sl.Eq(idx).Text(), "\n", ""), "  ", "", -1))
								e.orderStr.WriteString("Item: " + strings.Replace(strings.ReplaceAll(sl.Eq(idx).Text(), "\n", ""), "  ", "", -1))
								continue
							}

							rs := strings.Replace(strings.ReplaceAll(sl.Eq(idx).Text(), "\n", ""), "  ", "", -1)
							e.orderStr.WriteString(rs + "\n")

							e.etsyOrder.OrderId = e.orderId
						}
						ExtractEtsyOrder(e.orderStr.String(), &e.etsyOrder)
						e.etsyOrder.ShopName = shopName
						e.etsyOrder.OrderDate = recievedAt
						// e.etsyOrder.Address = address
						// e.etsyOrder.CustMail = cusMail
						if e.etsyOrder.TransactionId != "" {

							e.arrEtsyOrder[i] = NewEtsyFieldOrder(e.etsyOrder)
							i = i + 1

							etsyOrderRecords := NewEtsyOrderRecord(e.arrEtsyOrder)
							myErr = CreateEtsyOrder(appCtx, etsyOrderRecords)
							if myErr != nil {
								return
							}

						}
					}
				})

				//create order detail
				e.orderDetailStr = strings.Builder{}
				e.etsyOrderDetail = EtsyOrderDetail{}

				doc.Find(`td[style="border-collapse:collapse; text-align:left;"]`).Each(func(i int, s *goquery.Selection) {
					tdRight := s.Next()
					field := strings.ReplaceAll(strings.ReplaceAll(s.Text(), "\n", ""), "  ", "")
					val := strings.ReplaceAll(strings.ReplaceAll(tdRight.Text(), "\n", ""), "  ", "")
					rs := field + val
					e.orderDetailStr.WriteString(rs + "\n")
					t := e.orderDetailStr.String()
					ExtractEtsyOrderDetail(t, &e.etsyOrderDetail)
				})

				e.etsyOrderDetail.OrderDate = recievedAt
				e.arrEtsyOrderDetail = make([]EtsyFieldOrderDetail, 1)
				e.etsyOrderDetail.OrderId = e.orderId
				e.arrEtsyOrderDetail[0] = NewEtsyFieldOrderDetail(e.etsyOrderDetail)
				etsyOrderDetailRecords := NewEtsyOrderDetailRecord(e.arrEtsyOrderDetail)
				myErr = CreateEtsyOrderDetail(appCtx, etsyOrderDetailRecords)

				e.count = e.count + 1

				// custName
				doc.Find(`span[class='name']`).Each(func(i int, s *goquery.Selection) {
					custName = s.Text()
				})
				// road
				doc.Find(`span[class='first-line']`).Each(func(i int, s *goquery.Selection) {
					road = s.Text()
				})
				//city
				doc.Find(`span[class='city']`).Each(func(i int, s *goquery.Selection) {
					city = s.Text()
				})
				//state
				doc.Find(`span[class='state']`).Each(func(i int, s *goquery.Selection) {
					state = s.Text()
				})
				//zip
				doc.Find(`span[class='zip']`).Each(func(i int, s *goquery.Selection) {
					zip = s.Text()
				})
				//country
				doc.Find(`span[class='country-name']`).Each(func(i int, s *goquery.Selection) {
					country = s.Text()
				})

				//create order shipping
				e.etsyOrderShipping = EtsyOrderShipping{
					OrderId:  e.orderId,
					ShopName: shopName,
					CustName: custName,
					CustMail: cusMail,
					Road:     road,
					City:     city,
					State:    state,
					Zip:      zip,
					Country:  country,
				}
				e.arrEtsyOrderShipping = make([]EtsyFieldOrderShipping, 1)
				e.arrEtsyOrderShipping[0] = NewEtsyFieldOrderShipping(e.etsyOrderShipping)
				etsyOrderShippingRecords := NewEtsyOrderShippingRecord(e.arrEtsyOrderShipping)
				myErr = CreateEtsyOrderShipping(appCtx, etsyOrderShippingRecords)

			}
			index++

		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			log.Println("Got attachment: %v", filename)
		}

	}

	if myErr != nil {
		return myErr
	}

	return nil
}

func getShopName(mail string) string {
	shopName := ""
	switch mail {
	case "khanhlinhquachvietnam@gmail.com":
		shopName = "LIQ"
	case "ecoprimevietnam@gmail.com":
		shopName = "Paracel"
	case "oanhle310794@gmail.com":
		shopName = "Wong"
	case "newecoprime@gmail.com":
		shopName = "Kidscrafters"
	case "hangle254301@gmail.com":
		shopName = "Clackee"
	case "ngocanhleexyz@gmail.com":
		shopName = "Tatsuhouse"
	case "henryhousestore@gmail.com":
		shopName = "Henryhouses"
	}
	return shopName
}
