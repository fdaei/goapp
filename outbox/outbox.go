package outbox

/*
This solution only works with single container
We should add some distributed synchronization solution to handle workload
*/

import (
	"context"
	"fmt"
	"log"
	"time"

	"git.gocasts.ir/remenu/beehive/event"
	"github.com/go-co-op/gocron"
)

type Repository interface {
	UpdatePublished(ctx context.Context, eventIDs []uint64, publishedAt time.Time) (int64, error)
	UpdateUnpublished(ctx context.Context, eventIDs []uint64, lastRetryAt time.Time) (int64, error)
	//update `basket_outbox set retry_count = retry_count+1, last_retry_at = ? where id in(?)`
	UnpublishedCount(ctx context.Context, retryThreshold int64) (int64, error)
	//select count(*) from basket_outbox where retry_count < ?
	GetUnPublished(ctx context.Context, offset, limit, retryThreshold int) ([]Event, error)
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
	unPublishedOutBoxEvents, err := s.repository.GetUnPublished(context.Background(),
		0, s.config.BatchSize, s.config.RetryThreshold)
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
		_, err = s.repository.UpdatePublished(context.Background(), []uint64{}, time.Now())
		if err != nil {
			log.Printf("Error updating event: %s", err)
		}

		log.Printf("Published event: %s successfully", eventMessage.Topic)
	}
}
