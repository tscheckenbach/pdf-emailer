package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps the pgx pool for easier testing
type DB struct {
	pool *pgxpool.Pool
}

// NewConnection creates a new database connection pool
func NewConnection(databaseURL string) (*DB, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{pool: pool}, nil
}

// Close closes the database connection pool
func (db *DB) Close(ctx context.Context) {
	if db.pool != nil {
		db.pool.Close()
	}
}

// Pool returns the underlying pgx pool (for testing)
func (db *DB) Pool() interface{} {
	return db.pool
}

// Ping tests the database connection
func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Query executes a query and returns rows
func (db *DB) Query(ctx context.Context, sql string, args ...interface{}) (interface{}, error) {
	return db.pool.Query(ctx, sql, args...)
}

// QueryRow executes a query and returns a single row
func (db *DB) QueryRow(ctx context.Context, sql string, args ...interface{}) interface{} {
	return db.pool.QueryRow(ctx, sql, args...)
}

// Exec executes a query without returning rows
func (db *DB) Exec(ctx context.Context, sql string, args ...interface{}) (interface{}, error) {
	return db.pool.Exec(ctx, sql, args...)
}

// Begin starts a transaction
func (db *DB) Begin(ctx context.Context) (interface{}, error) {
	return db.pool.Begin(ctx)
}
