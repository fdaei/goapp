package order

import (
	"git.gocasts.ir/remenu/beehive/event"
	payment "git.gocasts.ir/remenu/beehive/paymentapp/service"
)

func NewEventConsumer(consumers ...event.Consumer) event.EventConsumer {
	s := New()
	return event.EventConsumer{Consumers: consumers, Router: event.Router{
		payment.PurchaseSucceedTopic: s.PurchaseSucceedHandler,
		payment.PurchaseFailedTopic:  s.PurchaseFailedHandler,
	}}
}
