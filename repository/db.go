package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

var ErrURLAlreadyExists = fmt.Errorf("short URL already exists")

type dbClient struct {
	db *sql.DB
}

func NewUrlStorage(db *sql.DB) dbClient {
	return dbClient{
		db: db,
	}
}

func (dbc dbClient) PutNewURL(ctx context.Context, originalUrl, shortUrl string) error {

	tx, err := dbc.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = dbc.db.ExecContext(ctx,
		`INSERT INTO shortUrls(
		 	base_url,
		 	short_url,
		 	created_at
		 )
		 VALUES ($1,$2,$3)`,
		originalUrl, shortUrl, time.Now())

	if err != nil {
		return fmt.Errorf("failed to insert new url for %s in DB: %w", shortUrl, err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
