package storagepg

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PostgresDB represents PostgreSQL connection
type PostgresDB struct {
	database *sqlx.DB
}

// PostgresDBConn init PostgreSQL connection by config information
func PostgresDBConn(config string) *PostgresDB {
	db, err := sqlx.Connect("postgres", config)
	db.MustExec(ddl)
	if err != nil {
		log.Fatalln(err)
	}
	return &PostgresDB{
		database: db,
	}
}

// Ping checks PostgreSQL connection
func (d *PostgresDB) Ping() bool {
	return d.database.Ping() == nil
}
