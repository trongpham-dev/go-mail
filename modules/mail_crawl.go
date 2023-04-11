package mailcrawl

import (
	"log"

	"github.com/emersion/go-imap/client"
)

type mailCrawl struct{

}

func NewMailCrawl() *mailCrawl{
	return &mailCrawl{}
}

func (m *mailCrawl) MailConnection(email, password string)(*client.Client, error){
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