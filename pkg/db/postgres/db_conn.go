package postgres

import (
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

// Return new Postgresql db instance
func NewPsqlDB(c *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlPort,
		c.Postgres.PostgresqlUser,
		c.Postgres.PostgresqlDbname,
		c.Postgres.PostgresqlPassword,
	)

	db, err := sqlx.Connect(c.Postgres.PgDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(30)
	db.SetConnMaxLifetime(120 * time.Second)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(20 * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
