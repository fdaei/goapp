package order

import (
	"fmt"

	"git.gocasts.ir/remenu/beehive/event"
)

type Service struct {
}

func New() Service {
	return Service{}
}

func (s Service) PurchaseSucceedHandler(event event.Event) error {
	fmt.Println("PurchaseSucceedHandler")
	return nil
}

func (s Service) PurchaseFailedHandler(event event.Event) error {
	fmt.Println("PurchaseFailedHandler")
	return nil
}
