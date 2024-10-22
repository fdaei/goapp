package basket

import (
	"git.gocasts.ir/remenu/beehive/event"
	basketrepo "git.gocasts.ir/remenu/beehive/service/basket/repository"
	"git.gocasts.ir/remenu/beehive/service/payment"
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
