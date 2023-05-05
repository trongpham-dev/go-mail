package pubsub

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap/client"
)

type MailData struct {
	Client *client.Client
	Ids    []uint32
	Mail   string
}

// description of message that pushed to message queue by publishers then recieved by subscribers
type Message struct {
	id        string
	channel   Topic // can be ignore
	data      MailData
	createdAt time.Time
}

func NewMessage(data MailData) *Message {
	now := time.Now().UTC()

	return &Message{
		id:        fmt.Sprintf("%d", now.UnixNano()),
		data:      data,
		createdAt: now,
	}
}

// getter
func (msg *Message) String() string {
	return fmt.Sprintf("Message %s", msg.channel)
}

// getter
func (msg *Message) Channel() Topic {
	return msg.channel
}

// setter
func (msg *Message) SetChannel(channel Topic) {
	msg.channel = channel
}

// getter
func (msg *Message) Data() MailData {
	return msg.data
}
