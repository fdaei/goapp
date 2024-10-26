package basketmysql

import (
	"database/sql"

	"git.gocasts.ir/remenu/beehive/pkg/postgresql"
)

type Repository struct {
	db *sql.DB
}

func New(config postgresql.Config) Repository {
	// create db connection
	// test ping request
	return sql.DB{}
}
