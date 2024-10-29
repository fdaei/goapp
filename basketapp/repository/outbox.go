package repository

import (
	"database/sql"
	"fmt"
	"git.gocasts.ir/remenu/beehive/outbox"
)

type OutBoxRepo struct {
	PostgreSQL *sql.DB
}

func NewOutBoxRepo(db *sql.DB) OutBoxRepo {
	return OutBoxRepo{
		PostgreSQL: db,
	}
}

func (repo OutBoxRepo) Create(outBoxEvent outbox.Event) (sql.Result, error) {

	// TODO: use from sqlc
	query := "INSERT INTO outbox_events (id, topic, payload, is_published) VALUES ($1, $2, $3, $4)"
	result, err := repo.PostgreSQL.Exec(query, outBoxEvent.ID, outBoxEvent.Topic, outBoxEvent.Payload, outBoxEvent.IsPublished)
	if err != nil {
		return nil, fmt.Errorf("error creating basket: %v", err)
	}
	return result, nil
}

func (repo OutBoxRepo) Update(outBoxEvent outbox.Event) (sql.Result, error) {
	query := "UPDATE outbox_events SET is_published=$1 WHERE id=$2"
	result, err := repo.PostgreSQL.Exec(query, outBoxEvent.IsPublished, outBoxEvent.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating basket: %v", err)
	}
	return result, nil
}

func (repo OutBoxRepo) GetByIsPublished(isPublished bool) ([]outbox.Event, error) {
	query := "SELECT id, topic, payload, is_published FROM outbox_events WHERE is_published=$1"
	rows, err := repo.PostgreSQL.Query(query, isPublished)
	if err != nil {
		return nil, fmt.Errorf("error retrieving baskets: %v", err)
	}
	defer rows.Close()

	var eventMessages []outbox.Event
	for rows.Next() {
		var e outbox.Event
		if err := rows.Scan(&e.ID, &e.Topic, &e.Payload, &e.IsPublished); err != nil {
			return nil, fmt.Errorf("error scanning basket: %v", err)
		}
		eventMessages = append(eventMessages, e)
	}

	return eventMessages, nil
}
