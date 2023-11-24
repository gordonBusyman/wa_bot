package users

import (
	"context"

	"github.com/Masterminds/squirrel"
)

// Create creates a new user.
func (s *Store) Create(ctx context.Context, res *Resource) (*Resource, error) {
	_, err := squirrel.Insert("users").
		SetMap(map[string]any{
			"id":      res.ID,
			"chat_id": res.ChatID,
			"name":    res.Name,
		}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(s.db).
		ExecContext(ctx)

	if err != nil {
		return nil, err
	}

	return s.Retrieve(ctx, res.ID)
}
