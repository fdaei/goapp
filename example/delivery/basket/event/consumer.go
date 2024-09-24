package event

import (
	"git.gocasts.ir/remenu/beehive/event"
	"git.gocasts.ir/remenu/beehive/example/service/basket"
	"git.gocasts.ir/remenu/beehive/example/service/payment"
)

func New(service basket.Service, consumers ...event.Consumer) event.EventConsumer {
	h := handler{service: service}
	return event.EventConsumer{Consumers: consumers, Router: event.Router{
		payment.PurchaseSucceedTopic: h.PurchaseSucceedHandler,
		payment.PurchaseFailedTopic:  h.PurchaseFailedHandler,
	}}
}
