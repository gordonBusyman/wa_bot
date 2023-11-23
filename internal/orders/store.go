package orders

import "database/sql"

// Store represents orders store.
type Store struct {
	db *sql.DB
}

// NewStore creates a new orders store.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
