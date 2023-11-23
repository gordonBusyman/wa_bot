package userFlows

import "database/sql"

// Store represents a userFlows store.
type Store struct {
	db *sql.DB
}

// NewStore creates a new userFlows store.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
