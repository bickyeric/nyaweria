package repository

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/bickyeric/nyaweria/entity"
	"github.com/doug-martin/goqu/v9"
)

type SummaryRequest struct {
	RecipientID        string
	Limit              int
	StartTime, EndTime time.Time
}

type Donate interface {
	Create(ctx context.Context, record *entity.Donation) error
	Summary(ctx context.Context, req SummaryRequest) ([]*entity.DonationSummary, error)
}

type donate struct {
	db *sql.DB
}

// Summary implements Donate.
func (u *donate) Summary(ctx context.Context, req SummaryRequest) ([]*entity.DonationSummary, error) {
	query, args, err := goqu.From("donations").
		Select("sender", goqu.SUM(goqu.C("amount").Cast("INTEGER")).As("sum")).
		Where(goqu.C("created_at").Gte(req.StartTime), goqu.C("created_at").Lte(req.EndTime), goqu.C("recipient_id").Eq(req.RecipientID)).
		GroupBy("sender").
		Order(goqu.I("sum").Desc()).
		Limit(uint(req.Limit)).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := u.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var summaries []*entity.DonationSummary
	for rows.Next() {
		var summary entity.DonationSummary

		if err := rows.Scan(&summary.Sender, &summary.Sum); err != nil {
			return nil, err
		}

		summaries = append(summaries, &summary)
	}

	return summaries, nil
}

func (u *donate) Create(ctx context.Context, record *entity.Donation) error {
	query, args, err := goqu.Insert("donations").
		Rows(goqu.Record{
			"sender":       record.From,
			"recipient_id": record.RecipientID,
			"currency":     "IDR",
			"amount":       record.Amount,
			"message":      record.Message,
		}).
		ToSQL()
	if err != nil {
		return err
	}

	result, err := u.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowID, _ := result.LastInsertId()
	record.ID = strconv.FormatInt(rowID, 10)
	return nil
}

func NewDonate(db *sql.DB) Donate {
	return &donate{
		db: db,
	}
}
