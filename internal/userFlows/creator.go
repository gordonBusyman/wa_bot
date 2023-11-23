package userFlows

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
)

// CreateMany creates a user flow elements for a flow_id
func (s *Store) CreateMany(ctx context.Context, userID int, flowID int, orderID int) error {
	//
	// Get all steps for a flow_id
	//
	fmt.Println("STEP 1")
	rows, err := squirrel.Select("id", "options").
		From("flow_steps").
		Where(squirrel.Eq{
			"flow_id": flowID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(s.db).
		QueryContext(ctx)
	if err != nil {
		return err
	}
	fmt.Println("STEP 2")

	step := &Step{}

	var options sql.NullString

	for rows.Next() {
		fmt.Println("STEP 3")

		if err := rows.Scan(&step.ID, &options); err != nil {
			return err
		}
		fmt.Println("STEP 4")

		//
		// If step has options, get all products for an order_id
		//
		if options.Valid {
			productRows, err := squirrel.Select("id").
				From("order_items").
				Where(squirrel.Eq{
					"order_id": orderID,
				}).
				PlaceholderFormat(squirrel.Dollar).
				RunWith(s.db).
				QueryContext(ctx)
			if err != nil {
				return err
			}

			var orderItemID int
			for productRows.Next() {
				if err := productRows.Scan(&orderItemID); err != nil {
					return err
				}

				fmt.Println("STEP ID", step.ID)
				//
				// Create a user flow element for each product
				//
				_, err = squirrel.Insert("user_flows").
					SetMap(map[string]any{
						"user_id":       userID,
						"step_id":       step.ID,
						"order_item_id": orderItemID,
					}).
					PlaceholderFormat(squirrel.Dollar).
					RunWith(s.db).
					ExecContext(ctx)
			}
		} else {
			//
			// If step has no options, create a user flow element for each step
			//
			_, err = squirrel.Insert("user_flows").
				SetMap(map[string]any{
					"user_id": userID,
					"step_id": step.ID,
				}).
				PlaceholderFormat(squirrel.Dollar).
				RunWith(s.db).
				ExecContext(ctx)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
