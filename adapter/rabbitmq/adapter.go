package rabbitmq

import (
	"context"
	"fmt"
	"git.gocasts.ir/remenu/beehive/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type Config struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

type Adapter struct {
	config Config
	queue  string
	conn   *amqp.Connection
	ch     *amqp.Channel
}

func New(config Config, queue string, topics []event.Topic) *Adapter {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", config.Username, config.Password, config.Host, config.Port))
	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: %s", "Failed to open a channel", err)
	}
	err = ch.ExchangeDeclare(
		"main",  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Panicf("%s: %s", "Failed to declare an exchange", err)
	}
	_, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Panicf("%s: %s", "Failed to declare a queue", err)
	}
	err = ch.ExchangeDeclare(
		"main",  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Panicf("%s: %s", "Failed to declare an exchange", err)
	}
	for _, t := range topics {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			queue, "main", t)
		err = ch.QueueBind(
			queue,     // queue name
			string(t), // routing key
			"main",    // exchange
			false,
			nil)
		if err != nil {
			log.Panicf("%s: %s", "Failed to bind a queue", err)
		}
	}
	return &Adapter{config: config, ch: ch, conn: conn}
}
func (a *Adapter) Publish(event event.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.ch.PublishWithContext(ctx,
		"main",              // exchange
		string(event.Topic), // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        event.Payload,
		})
	if err != nil {
		log.Panicf("%s: %s", "Failed to publish a message", err)
	}
	log.Printf(" [x] Sent %s", event.Payload)

	return err
}
func (a *Adapter) Consume(eventStream chan<- event.Event) error {

	messages, err := a.ch.Consume(
		a.queue, // queue
		"",      // basket
		true,    // auto ack
		false,   // exclusive
		false,   // no local

		false, // no wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("%s: %s", "Failed to register a basket", err)
	}

	go func() {
		for d := range messages {
			log.Println("arrived", string(d.Body))
			eventStream <- event.Event{Topic: event.Topic(d.RoutingKey), Payload: d.Body}
		}
	}()

	return nil
}
func (a *Adapter) Close() {
	err := a.ch.Close()
	if err != nil {
		log.Panicf("unexpected error: %v", err)
	}
	err = a.conn.Close()
	if err != nil {
		log.Panicf("unexpected error: %v", err)
	}
}
