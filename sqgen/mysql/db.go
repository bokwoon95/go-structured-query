package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func openAndPing(database string) (*sql.DB, error) {
	db, err := sql.Open("postgres", database)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf(
			"Could not ping the database, is the database reachable via %s? %w",
			database,
			err,
		)
	}

	return db, nil
}
