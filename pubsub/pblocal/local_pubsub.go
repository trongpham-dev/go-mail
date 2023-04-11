package pblocal

import (
	"context"
	"go-mail/common"
	"go-mail/pubsub"
	"log"
	"sync"
)

// A pb run locally (in-mem)
// It has a queue (buffer channel) at it's core and many group of subscribers.
// Because we want to send a message with a specific topic for many subscribers in a group can handle.

// message broker
type localPubSub struct {
	messageQueue chan *pubsub.Message                    // used by publisher
	mapChannel   map[pubsub.Topic][]chan *pubsub.Message // topic : [] chanel message that is used by subscribers.
	locker       *sync.RWMutex
}

func NewPubSub() *localPubSub {
	pb := &localPubSub{
		messageQueue: make(chan *pubsub.Message, 10000),
		mapChannel:   make(map[pubsub.Topic][]chan *pubsub.Message),
		locker:       new(sync.RWMutex),
	}

	pb.run()

	return pb
}

func (ps *localPubSub) Publish(ctx context.Context, topic pubsub.Topic, data *pubsub.Message) error {
	data.SetChannel(topic) // set topic name

	go func() {
		defer common.AppRecover()
		ps.messageQueue <- data // push message to message queue
		log.Println("New event published:", data.String(), "with data", data.Data())
	}()
	return nil
}

// ch <-chan *pubsub.Message -> Subscriber communicates with this chanel to recieve message
// close func() -> for unsubscribing
func (ps *localPubSub) Subscribe(ctx context.Context, topic pubsub.Topic) (ch <-chan *pubsub.Message, close func()) {
	c := make(chan *pubsub.Message)

	ps.locker.Lock()

	val, ok := ps.mapChannel[topic] // get list of message chanel according to the topic

	// if list of message chanel is existed
	if ok {
		val = append(ps.mapChannel[topic], c) // append a msg chanel to an existing list of msg chanel
		ps.mapChannel[topic] = val
	} else {
		ps.mapChannel[topic] = []chan *pubsub.Message{c} // create an array of message chanel with c is the first element
	}

	ps.locker.Unlock()

	return c, func() {
		log.Println("Unsubscribe")

		if chans, ok := ps.mapChannel[topic]; ok {
			for i := range chans {
				if chans[i] == c {
					// remove element at index in chans
					chans = append(chans[:i], chans[i+1:]...)

					ps.locker.Lock()
					ps.mapChannel[topic] = chans
					ps.locker.Unlock()
					break
				}
			}
		}
	}

}

// ditributes message
func (ps *localPubSub) run() error {
	log.Println("Pubsub started")

	go func() {
		for {
			mess := <-ps.messageQueue // get message from message queue
			log.Println("Message dequeue:", mess)

			subs, ok := ps.mapChannel[mess.Channel()] // get list message chanel from topic

			if ok {
				for i := range subs {
					go func(c chan *pubsub.Message) {
						c <- mess // push message to topic message chanel
						//f(mess)
					}(subs[i])
				}
			}
			//else {
			//	ps.messageQueue <- mess
			//}
		}
	}()

	return nil
}
