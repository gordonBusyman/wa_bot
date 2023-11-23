package users

import (
	"context"
	"github.com/Masterminds/squirrel"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Create creates a new user.
func (s *Store) Create(ctx context.Context, update *tgbotapi.Message) (*Resource, error) {
	_, err := squirrel.Insert("users").
		SetMap(map[string]any{
			"id":      update.From.ID,
			"chat_id": update.Chat.ID,
			"name":    update.From.FirstName,
		}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(s.db).
		ExecContext(ctx)

	if err != nil {
		return nil, err
	}

	return s.Retrieve(ctx, update.From.ID)
}
