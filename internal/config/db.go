package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

type ConfigDB struct {
	ConnStrDB string
}

func NewConfig() *ConfigDB {

	return &ConfigDB{
		ConnStrDB: os.Getenv("DATABASE_URL"),
	}
}

func (c *ConfigDB) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.ConnStrDB)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(nil, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}
