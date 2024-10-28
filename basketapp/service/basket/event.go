package basket

import (
	basketrepo "git.gocasts.ir/remenu/beehive/basket_app/repository/basket"
	"git.gocasts.ir/remenu/beehive/event"
	payment "git.gocasts.ir/remenu/beehive/payment_app/service"
)

func NewEventConsumer(repo basketrepo.Repository, consumers ...event.Consumer) event.EventConsumer {
	s := NewBasketService(repo)

	return event.EventConsumer{
		Consumers: consumers,
		Router: event.Router{
			payment.PurchaseSucceedTopic: s.PurchaseSucceedHandler,
			payment.PurchaseFailedTopic:  s.PurchaseFailedHandler,
		},
	}
}
