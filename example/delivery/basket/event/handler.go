package event

import (
	"fmt"
	"git.gocasts.ir/remenu/beehive/event"
	"git.gocasts.ir/remenu/beehive/example/service/basket"
)

type handler struct {
	service basket.Service
}

func (h handler) PurchaseSucceedHandler(event event.Event) error {
	fmt.Println("PurchaseSucceedHandler")
	return nil
}

func (h handler) PurchaseFailedHandler(event event.Event) error {
	fmt.Println("PurchaseFailedHandler")
	return nil
}
