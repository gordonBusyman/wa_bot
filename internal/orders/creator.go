package orders

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

// Create creates a new user.
func (s *Store) Create(ctx context.Context, userID int, items []Item) (*Resource, error) {

	order := &Resource{
		UserID: userID,
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %v", err)
	}

	defer tx.Rollback()

	err = squirrel.Insert("orders").
		SetMap(map[string]any{
			"user_id": userID,
		}).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		QueryRowContext(ctx).
		Scan(&order.ID)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Order ID INSERTED: %d\n", order.ID)

	for _, item := range items {
		err = squirrel.Insert("order_items").
			SetMap(map[string]any{
				"order_id":   order.ID,
				"product_id": item.ProductID,
				"quantity":   item.Quantity,
			}).
			Suffix("RETURNING id").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(tx).
			QueryRowContext(ctx).
			Scan(&item.ID)

		order.Items = append(order.Items, item)

		if err != nil {
			return nil, err
		}
	}

	tx.Commit()

	return order, nil
}
