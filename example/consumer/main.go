package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"git.gocasts.ir/remenu/beehive/adapter/rabbitmq"
	"git.gocasts.ir/remenu/beehive/config"
	"git.gocasts.ir/remenu/beehive/event"
	payment "git.gocasts.ir/remenu/beehive/payment_app/service"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	cfg := config.Load("config.yml")
	queue := "basket"

	topics1 := []event.Topic{
		payment.PurchaseSucceedTopic,
	}
	rabbitMQ1 := rabbitmq.New(cfg.RabbitMQ, queue, topics1)

	topics2 := []event.Topic{
		payment.PurchaseFailedTopic,
	}
	rabbitMQ2 := rabbitmq.New(cfg.RabbitMQ, queue, topics2)

	go func() {
		for i := range 10 {
			data := map[string]interface{}{"message": "hello", "count": i}
			payload, err := json.Marshal(data)
			err = rabbitMQ1.Publish(event.Event{
				Topic:   payment.PurchaseSucceedTopic,
				Payload: payload,
			})
			err = rabbitMQ2.Publish(event.Event{
				Topic:   payment.PurchaseFailedTopic,
				Payload: payload,
			})
			if err != nil {
				fmt.Println("err", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		// eventConsumer := basket.NewEventConsumer(rabbitMQ1, rabbitMQ2)
		// eventConsumer.Start()
		wg.Done()
	}()

	wg.Wait()
}
