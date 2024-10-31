package basket

import (
	"git.gocasts.ir/remenu/beehive/event"
	payment "git.gocasts.ir/remenu/beehive/paymentapp/service"
)

func NewEventConsumer(repo Repository, consumers ...event.Consumer) event.EventConsumer {
	s := NewService(repo)

	return event.EventConsumer{
		Consumers: consumers,
		Router: event.Router{
			payment.PurchaseSucceedTopic: s.PurchaseSucceedHandler,
			payment.PurchaseFailedTopic:  s.PurchaseFailedHandler,
		},
	}
}
