package etsy

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/emersion/go-message/mail"
	"io"
	"log"
	"regexp"
	"strings"
)

type etSy struct {
	idx                  int
	orderId              string
	orderStr             strings.Builder
	orderDetailStr       strings.Builder
	arrEtsyOrder         []EtsyFieldOrder
	arrEtsyOrderDetail   []EtsyFieldOrderDetail
	etsyOrder            EtsyOrder
	etsyFieldOrder       EtsyFieldOrder
	etsyOrderDetail      EtsyOrderDetail
	etsyFieldOrderDetail EtsyFieldOrderDetail
	count                int
}

func NewEtsy() *etSy {
	return &etSy{}
}

func (e *etSy) CrawlEtsy(mr *mail.Reader) {
	e.count = 0
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			doc, err := goquery.NewDocumentFromReader(p.Body)

			if err != nil {
				panic(err)
			}
			// rs := ""
			doc.Find(`div[style="font-family:arial,helvetica,sans-serif;color:#444444;font-size:16px;line-height:24px"]`).Each(func(i int, s *goquery.Selection) {
				pattern := regexp.MustCompile("Your order number is:\\s+(\\S+)")
				match := pattern.FindStringSubmatch(s.Text())
				if len(match) > 0 {
					e.orderId = match[1]
					e.orderStr.WriteString(e.orderId + "\n")
				}
			})

			doc.Find(`td[valign="top"][style="line-height:0px"]`).Each(func(i int, s *goquery.Selection) {
				e.orderStr = strings.Builder{}
				e.etsyOrder = EtsyOrder{}
				i = 0
				if len(s.Nodes) > 0 {
					e.arrEtsyOrder = make([]EtsyFieldOrder, len(s.Nodes))

					sl := s.Find(`div[style*="font-family:arial,helvetica,sans-serif;"]`)
					for idx := range sl.Nodes {
						if sl.Eq(idx).Find(`a[style="text-decoration:none;color:#222222"]`).Text() != "" {
							return
						}

						if sl.Eq(idx).Find(`a[style="text-decoration:none;color:#444444"]`).Text() != "" {
							// log.Println("Item:", strings.Replace(strings.ReplaceAll(sl.Eq(idx).Text(), "\n", ""), "  ", "", -1))
							e.orderStr.WriteString("Item: " + strings.Replace(strings.ReplaceAll(sl.Eq(idx).Text(), "\n", ""), "  ", "", -1))
							continue
						}

						rs := strings.Replace(strings.ReplaceAll(sl.Eq(idx).Text(), "\n", ""), "  ", "", -1)
						e.orderStr.WriteString(rs + "\n")

						e.etsyOrder.OrderId = e.orderId
					}
					ExtractEtsyOrder(e.orderStr.String(), &e.etsyOrder)
					if e.etsyOrder.TransactionId != "" {

						e.arrEtsyOrder[i] = NewEtsyFieldOrder(e.etsyOrder)
						i = i + 1

						etsyOrderRecords := NewEtsyOrderRecord(e.arrEtsyOrder)
						CreateEtsyOrder(etsyOrderRecords)
					}
				}
			})

			//
			e.orderDetailStr = strings.Builder{}
			e.etsyOrderDetail = EtsyOrderDetail{}
			doc.Find(`td[style="border-collapse:collapse;text-align:left"]`).Each(func(i int, s *goquery.Selection) {
				tdRight := s.Next()
				field := strings.ReplaceAll(strings.ReplaceAll(s.Text(), "\n", ""), "  ", "")
				val := strings.ReplaceAll(strings.ReplaceAll(tdRight.Text(), "\n", ""), "  ", "")
				rs := field + val
				e.orderDetailStr.WriteString(rs + "\n")
				t := e.orderDetailStr.String()
				ExtractEtsyOrderDetail(t, &e.etsyOrderDetail)
			})

			if e.count > 0 {
				e.arrEtsyOrderDetail = make([]EtsyFieldOrderDetail, 1)
				e.etsyOrderDetail.OrderId = e.orderId
				e.arrEtsyOrderDetail[0] = NewEtsyFieldOrderDetail(e.etsyOrderDetail)
				etsyOrderDetailRecords := NewEtsyOrderDetailRecord(e.arrEtsyOrderDetail)
				CreateEtsyOrderDetail(etsyOrderDetailRecords)
			}
			e.count = e.count + 1

		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			log.Println("Got attachment: %v", filename)
		}

	}
}
