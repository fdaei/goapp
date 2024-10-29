package outbox

import (
	"database/sql"
	"fmt"
	"git.gocasts.ir/remenu/beehive/event"
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

type Repository interface {
	Update(outboxEvent Event) (sql.Result, error)
	GetByIsPublished(isPublished bool) ([]Event, error)
}

type Scheduler struct {
	repository Repository
	sch        *gocron.Scheduler
	publisher  event.Publisher
	config     Config
}

func New(repository Repository, pub event.Publisher, cfg Config) Scheduler {
	return Scheduler{
		repository: repository,
		sch:        gocron.NewScheduler(time.UTC),
		publisher:  pub,
		config:     cfg,
	}
}

func (s Scheduler) Start(done <-chan bool) {
	_, err := s.sch.Every(s.config.IntervalInSeconds).Second().Do(s.PublishOutBoxEvents)
	if err != nil {
		log.Printf("Error starting outbox: %s", err)
	}
	s.sch.StartAsync()

	<-done
	// wait to finish job
	fmt.Println("stop scheduler..")
	s.sch.Stop()
}

func (s Scheduler) PublishOutBoxEvents() {
	log.Println("starting outbox publisher..")
	unPublishedOutBoxEvents, err := s.repository.GetByIsPublished(false)
	if err != nil {
		log.Printf("Error fetching OutBoxEvents: %s", err)
	}

	for _, eventMessage := range unPublishedOutBoxEvents {
		err := s.publisher.Publish(event.Event{
			Topic:   eventMessage.Topic,
			Payload: eventMessage.Payload,
		})
		if err != nil {
			log.Printf("Error publishing event: %s", err)
		}

		eventMessage.IsPublished = true
		_, err = s.repository.Update(eventMessage)
		if err != nil {
			log.Printf("Error updating event: %s", err)
		}

		log.Printf("Published event: %s successfully", eventMessage.Topic)
	}
}
