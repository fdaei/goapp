package payment

import "git.gocasts.ir/remenu/beehive/event"

const (
	PurchaseSucceedTopic = event.Topic("payment.purchase_succeed")
	PurchaseFailedTopic  = event.Topic("payment.purchase_failed")
)
