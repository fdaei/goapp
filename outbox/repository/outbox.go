package repository

import (
	"context"
	"database/sql"
	"fmt"
	"git.gocasts.ir/remenu/beehive/outbox"
	"git.gocasts.ir/remenu/beehive/types"
	pg "github.com/lib/pq"
	"time"
)

type OutBoxRepo struct {
	PostgreSQL *sql.DB
}

func NewOutBoxRepo(db *sql.DB) OutBoxRepo {
	return OutBoxRepo{
		PostgreSQL: db,
	}
}

func (repo OutBoxRepo) UpdatePublished(ctx context.Context, eventIDs []types.ID, publishedAt time.Time) error {
	query := "UPDATE outbox_events set retried_count = retried_count+1, last_retried_at = $1, published_at=$2 ,is_published=TRUE where id = ANY($3::bigint[])"
	_, err := repo.PostgreSQL.ExecContext(ctx, query, publishedAt, publishedAt, pg.Array(eventIDs))
	if err != nil {
		return fmt.Errorf("error updating basket: %v", err)
	}
	return nil
}

func (repo OutBoxRepo) UpdateUnpublished(ctx context.Context, eventIDs []types.ID, lastRetriedAt time.Time) error {
	query := "UPDATE outbox_events set retried_count = retried_count+1, last_retried_at = $1, is_published=TRUE WHERE id IN (?) "
	_, err := repo.PostgreSQL.ExecContext(ctx, query, lastRetriedAt, eventIDs)
	if err != nil {
		return fmt.Errorf("error updating basket: %v", err)
	}
	return nil
}

func (repo OutBoxRepo) GetUnPublished(ctx context.Context, offset, limit, retryThreshold int) ([]outbox.Event, error) {
	query := `
		SELECT id, topic, payload, is_published 
		FROM outbox_events 
		WHERE is_published = FALSE AND retried_count < $1
		ORDER BY id 
		LIMIT $2 OFFSET $3`

	// Execute the query with limit and offset as parameters
	rows, err := repo.PostgreSQL.QueryContext(ctx, query, retryThreshold, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving outbox events: %v", err)
	}
	defer rows.Close()

	// Collect results without predefining a limit on slice size
	outBoxEvents := make([]outbox.Event, 0)
	for rows.Next() {
		var e outbox.Event
		if err := rows.Scan(&e.ID, &e.Topic, &e.Payload, &e.IsPublished); err != nil {
			return nil, fmt.Errorf("error scanning outbox event: %v", err)
		}
		outBoxEvents = append(outBoxEvents, e)
	}

	return outBoxEvents, nil
}

func (repo OutBoxRepo) UnpublishedCount(ctx context.Context, retryThreshold int64) (int64, error) {
	var unPublishedCount int64
	query := "SELECT count(*) from outbox_events where retried_count < $1"
	rows, err := repo.PostgreSQL.QueryContext(ctx, query, retryThreshold)
	if err != nil {
		return 0, fmt.Errorf("error retrieving baskets: %v", err)
	}
	defer rows.Close()
	rows.Scan(&unPublishedCount)
	return unPublishedCount, nil
}
