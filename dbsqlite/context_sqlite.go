package dbsqlite

import (
	"context"
	"time"
)

// NewDBContext creates a new context with a timeout for database operations.
// This ensures that long-running queries are canceled after the specified duration.
func NewDBContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 2*time.Second)
}
