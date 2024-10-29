package outbox

import "git.gocasts.ir/remenu/beehive/event"

type Event struct {
	ID          uint        `json:"id"`
	Topic       event.Topic `json:"topic"`
	Payload     []byte      `json:"payload"`
	IsPublished bool        `json:"is_published"`
}
