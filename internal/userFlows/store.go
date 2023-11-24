package userFlows

import (
	"database/sql"

	"github.com/cdipaolo/sentiment"
)

// Store represents a userFlows store.
type Store struct {
	db *sql.DB

	sentimentAnalysisModel *sentiment.Models
}

// NewStore creates a new userFlows store.
func NewStore(db *sql.DB, model *sentiment.Models) *Store {
	return &Store{
		db:                     db,
		sentimentAnalysisModel: model,
	}
}
