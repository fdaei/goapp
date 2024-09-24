package event

import "log"

type handlerFunc func(event Event) error
type Router map[Topic]handlerFunc

type EventConsumer struct {
	Consumers []Consumer
	Router    Router
}

func (c EventConsumer) Start() {
	eventStream := make(chan Event, 1024)
	for _, consumer := range c.Consumers {
		err := consumer.Consume(eventStream)
		if err != nil {
			log.Printf("can't start consuming events: %v", err)
		}
	}
	for e := range eventStream {
		go func() {
			err := c.Router[e.Topic](e)
			if err != nil {
				log.Printf("can't handle event with topic %s, : %v", e.Topic, err)
			}
		}()

	}
}
