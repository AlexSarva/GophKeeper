package storagepg

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	database *sqlx.DB
}

func PostgresDBConn(config string) *PostgresDB {
	db, err := sqlx.Connect("postgres", config)
	if err != nil {
		log.Fatalln(err)
	}
	return &PostgresDB{
		database: db,
	}
}

func (d *PostgresDB) Ping() bool {
	return d.database.Ping() == nil
}
