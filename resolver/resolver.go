package resolver

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"time"
)

type resolver struct {
	onMessageEvents     chan *messageEventResolver
	onMessageSubscriber chan *onMessageSubscriber

	client *redis.Client
}

func NewResolver(client *redis.Client) *resolver {
	r := &resolver{
		onMessageEvents:     make(chan *messageEventResolver),
		onMessageSubscriber: make(chan *onMessageSubscriber),

		client: client,
	}

	go r.broadcastMessage()
	go r.broadcastRedis()

	return r
}

type onMessageSubscriber struct {
	stop   <-chan struct{}
	events chan<- *messageEventResolver
}

func (r *resolver) broadcastMessage() {
	subscribers := map[string]*onMessageSubscriber{}
	unsubscribe := make(chan string)

	// NOTE: subscribing and sending events are at odds.
	for {
		select {
		case id := <-unsubscribe:
			delete(subscribers, id)
		case s := <-r.onMessageSubscriber:
			id := randomID()
			log.Println("subscribe ", id)
			subscribers[id] = s
		case e := <-r.onMessageEvents:
			for id, s := range subscribers {
				go func(id string, s *onMessageSubscriber) {
					select {
					case <-s.stop:
						log.Println("unsubscribe ", id)
						unsubscribe <- id
					case s.events <- e:
					case <-time.After(time.Second):
					}
				}(id, s)
			}
		}
	}
}

func (r *resolver) broadcastRedis() {
	pubsub := r.client.Subscribe("mychannel1")
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()

	// Consume messages.
	for msg := range ch {
		var e messageEvent
		if err := json.Unmarshal([]byte(msg.Payload), &e); err != nil {
			log.Fatalln("Invalid message payload: ", msg.Payload)
		}

		r.onMessageEvents <- &messageEventResolver{e: e}
	}
}

func randomID() string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, 16)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}