package userFlows

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"strconv"
)

// Update updates a user flow element.
func (s *Store) Update(ctx context.Context, step *Resource, response string) error {

	fmt.Println("Update: ", step, response)
	_, err := squirrel.Update("user_flows").
		Set("complete", true).
		Where(squirrel.Eq{
			"id": step.ID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Update the order_item ID: ", step.OrderItemID)

	if step.Step.Options != nil {
		fmt.Println("Update the order_item ID (rating): ", step.OrderItemID)

		rating, err := strconv.Atoi(response)
		if err != nil {
			fmt.Println("RATING ERROR : ", response)

			return nil
		}

		fmt.Println("RATING: ", rating)

		_, err = squirrel.Update("order_items").
			Set("rating", rating).
			Where(squirrel.Eq{"id": step.OrderItemID}).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(s.db).
			ExecContext(ctx)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Update the order_item ID (feedback): ", step.OrderItemID)
		fmt.Println("FEEDBACK: ", response)

		_, err = squirrel.Update("order_items").
			Set("feedback", response).
			Where(squirrel.Eq{"id": step.OrderItemID}).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(s.db).
			ExecContext(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
