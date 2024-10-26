package order

import (
	"git.gocasts.ir/remenu/beehive/event"
	"git.gocasts.ir/remenu/beehive/service/payment"
)

func NewEventConsumer(consumers ...event.Consumer) event.EventConsumer {
	s := New()
	return event.EventConsumer{Consumers: consumers, Router: event.Router{
		payment.PurchaseSucceedTopic: s.PurchaseSucceedHandler,
		payment.PurchaseFailedTopic:  s.PurchaseFailedHandler,
	}}
}
