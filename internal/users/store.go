package users

import "database/sql"

// Store represents a user store.
type Store struct {
	db *sql.DB
}

// NewStore creates a new user store.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
