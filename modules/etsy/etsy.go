package etsy

import (
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/emersion/go-message/mail"
)

func CrawlEtsy(mr *mail.Reader){
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
					rs :=  pattern.FindString(s.Text())
					if rs != ""{
		
						log.Println(s.Text())
					}
				})

				doc.Find(`td[valign="top"][style="line-height:0px"]`).Each(func(i int, s *goquery.Selection) {
					if len(s.Nodes) > 0 {
						// arr := make([]interface{}, len(s.Nodes))
						
						sl := s.Find(`div[style*="font-family:arial,helvetica,sans-serif;"]`)
						for idx := range sl.Nodes{
							if sl.Eq(idx).Find(`a[style="text-decoration:none;color:#222222"]`).Text() != ""{
								break
							}

							if sl.Eq(idx).Find(`a[style="text-decoration:none;color:#444444"]`).Text() != ""{
								log.Println("Item:", strings.Replace(strings.ReplaceAll(sl.Eq(idx).Text(), "\n", ""),"  ", "",-1))
								continue
							}

							rs := strings.Replace(strings.ReplaceAll(sl.Eq(idx).Text(), "\n", ""),"  ", "",-1)
								log.Println(rs)
						}
					}
				})

				//
				doc.Find(`td[style="border-collapse:collapse;text-align:left"]`).Each(func(i int, s *goquery.Selection) {
					tdRight := s.Next()
					log.Println(strings.Replace(strings.ReplaceAll(s.Text(), "\n", ""),"  ", "",-1), strings.Replace(strings.ReplaceAll(tdRight.Text(), "\n", ""),"  ", "",-1))
				})

			case *mail.AttachmentHeader:
				// This is an attachment
				filename, _ := h.Filename()
				log.Println("Got attachment: %v", filename)
		}

	}
}