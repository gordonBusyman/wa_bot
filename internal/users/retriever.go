package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
)

// Retrieve retrieves a user by id.
func (s *Store) Retrieve(ctx context.Context, id int) (*Resource, error) {
	user := &Resource{
		ID: id,
	}

	err := squirrel.Select("name", "chat_id").
		From("users").
		Where(squirrel.Eq{
			"users.id": id,
		}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(s.db).
		QueryRowContext(ctx).
		Scan(&user.Name, &user.ChatID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
