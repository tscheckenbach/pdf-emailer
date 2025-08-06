package database

import (
	"context"
)

// DBInterface defines the interface for database operations
type DBInterface interface {
	Close(ctx context.Context)
	Pool() interface{}
	Ping(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...interface{}) (interface{}, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) interface{}
	Exec(ctx context.Context, sql string, args ...interface{}) (interface{}, error)
	Begin(ctx context.Context) (interface{}, error)
}
