package userFlows

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"strings"
)

// RetrieveCurrent retrieves a user flow by user id.
func (s *Store) RetrieveCurrent(ctx context.Context, userID int) (*Resource, error) {
	flowStep := &Resource{
		UserID:  userID,
		Step:    &Step{},
		Product: &Product{},
	}
	var (
		options     sql.NullString
		productName sql.NullString
		orderItemID sql.NullInt32
		orderID     sql.NullInt32
	)

	err := squirrel.Select(
		"user_flows.id",
		"flow_steps.name",
		"flow_steps.details",
		"flow_steps.options",
		"order_items.id",
		"products.name",
		"user_flows.order_id",
	).
		From("user_flows").
		LeftJoin("flow_steps ON user_flows.step_id = flow_steps.id").
		LeftJoin("order_items ON user_flows.order_item_id = order_items.id").
		LeftJoin("products ON order_items.product_id = products.id").
		Where(squirrel.And{
			squirrel.Eq{"user_flows.user_id": userID},
			squirrel.Eq{"user_flows.complete": false},
		}).
		OrderBy("flow_steps.order ASC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(s.db).
		QueryRowContext(ctx).
		Scan(
			&flowStep.ID,
			&flowStep.Step.Name,
			&flowStep.Step.Details,
			&options,
			&orderItemID,
			&productName,
			&orderID,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if options.Valid {
		flowStep.Step.Options = strings.Split(options.String, ",")
	}

	if orderItemID.Valid {
		flowStep.OrderItemID = int(orderItemID.Int32)
	}

	if productName.Valid {
		flowStep.Product.Name = productName.String
	}

	if orderID.Valid {
		flowStep.OrderID = int(orderID.Int32)
	}

	return flowStep, nil
}
