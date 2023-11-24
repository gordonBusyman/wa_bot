package userFlows

import (
	"context"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/cdipaolo/sentiment"
)

// Update updates a user flow element.
func (s *Store) Update(ctx context.Context, step *Resource, response string) error {
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

	if step.Step.Options != nil {
		rating, err := strconv.Atoi(response)
		if err != nil {
			return nil
		}

		_, err = squirrel.Update("order_items").
			Set("rating", rating).
			Where(squirrel.Eq{
				"id": step.OrderItemID,
			}).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(s.db).
			ExecContext(ctx)
		if err != nil {
			return err
		}
	} else {
		//
		// Kind of sentiment "analysis"
		//
		analysis := s.sentimentAnalysisModel.SentimentAnalysis(response, sentiment.English)

		_, err = squirrel.Update("orders").
			Set("feedback", response).
			Set("score", int(analysis.Score)).
			Where(squirrel.Eq{"id": step.OrderID}).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(s.db).
			ExecContext(ctx)

		if err != nil {
			return err
		}
	}

	return nil
}

// AnalyzeSentiment returns the sentiment score for the given text
func AnalyzeSentiment(text string) int {
	model, err := sentiment.Restore()
	if err != nil {
		panic(err)
	}

	analysis := model.SentimentAnalysis(text, sentiment.English)
	return int(analysis.Score)
}
