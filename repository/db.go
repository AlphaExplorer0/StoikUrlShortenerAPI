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

func (dbc dbClient) PutNewURL(ctx context.Context, originalUrl, shortUrl string) (string, error) {

	tx, err := dbc.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback() }()

	var count int
	err = dbc.db.QueryRowContext(ctx,
		`SELECT count(1)
		 FROM shortUrls WHERE 
		 shortUrls.short_url = $1`,
		shortUrl,
	).Scan(&count)

	if err != nil {
		return "", fmt.Errorf("failed to get if url exists for %s in DB: %w", shortUrl, err)
	}

	if count > 0 {
		return "", ErrURLAlreadyExists
	}

	var ansUrl string
	err = dbc.db.QueryRowContext(ctx,
		`WITH ins AS (
		 	INSERT INTO shortUrls(
				base_url,
				short_url,
				created_at
			)
		 	VALUES ($1,$2,$3)
			ON CONFLICT (base_url)
		 	DO NOTHING
		 	RETURNING short_url
			)
		 SELECT short_url FROM ins
		 UNION  ALL
		 SELECT short_url FROM shortUrls
		 WHERE  base_url = $1`,
		originalUrl, shortUrl, time.Now()).Scan(&ansUrl)

	if err != nil {
		return "", fmt.Errorf("failed to insert new url for %s in DB: %w", originalUrl, err)
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return ansUrl, nil
}
