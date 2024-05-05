// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: campaign_lines.sql

package db

import (
	"context"
	"database/sql"
)

const createCampaignLine = `-- name: CreateCampaignLine :one
INSERT INTO oms.campaign_line_items (campaign_id, name, booked, actual, adjustments, started_at, ended_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id
`

type CreateCampaignLineParams struct {
	CampaignID  int32
	Name        string
	Booked      string
	Actual      sql.NullString
	Adjustments sql.NullString
	StartedAt   sql.NullTime
	EndedAt     sql.NullTime
}

func (q *Queries) CreateCampaignLine(ctx context.Context, arg CreateCampaignLineParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, createCampaignLine,
		arg.CampaignID,
		arg.Name,
		arg.Booked,
		arg.Actual,
		arg.Adjustments,
		arg.StartedAt,
		arg.EndedAt,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createCampaignLineWithID = `-- name: CreateCampaignLineWithID :one

INSERT INTO oms.campaign_line_items (id, campaign_id, name, booked, actual, adjustments, started_at, ended_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id
`

type CreateCampaignLineWithIDParams struct {
	ID          int32
	CampaignID  int32
	Name        string
	Booked      string
	Actual      sql.NullString
	Adjustments sql.NullString
	StartedAt   sql.NullTime
	EndedAt     sql.NullTime
}

// campaign_line.sql
func (q *Queries) CreateCampaignLineWithID(ctx context.Context, arg CreateCampaignLineWithIDParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, createCampaignLineWithID,
		arg.ID,
		arg.CampaignID,
		arg.Name,
		arg.Booked,
		arg.Actual,
		arg.Adjustments,
		arg.StartedAt,
		arg.EndedAt,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteCampaignLine = `-- name: DeleteCampaignLine :exec
DELETE FROM oms.campaign_line_items WHERE id = $1
`

func (q *Queries) DeleteCampaignLine(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteCampaignLine, id)
	return err
}

const getCampaignLine = `-- name: GetCampaignLine :one
SELECT id, campaign_id, name, booked, actual, adjustments, started_at, ended_at, created_at, updated_at FROM oms.campaign_line_items WHERE id = $1
`

func (q *Queries) GetCampaignLine(ctx context.Context, id int32) (OmsCampaignLineItem, error) {
	row := q.db.QueryRowContext(ctx, getCampaignLine, id)
	var i OmsCampaignLineItem
	err := row.Scan(
		&i.ID,
		&i.CampaignID,
		&i.Name,
		&i.Booked,
		&i.Actual,
		&i.Adjustments,
		&i.StartedAt,
		&i.EndedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listCampaignLineItems = `-- name: ListCampaignLineItems :many
SELECT id, campaign_id, name, booked, actual, adjustments, started_at, ended_at, created_at, updated_at FROM oms.campaign_line_items 
WHERE id > $1
Order by id
LIMIT $2
`

type ListCampaignLineItemsParams struct {
	ID    int32
	Limit int32
}

func (q *Queries) ListCampaignLineItems(ctx context.Context, arg ListCampaignLineItemsParams) ([]OmsCampaignLineItem, error) {
	rows, err := q.db.QueryContext(ctx, listCampaignLineItems, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OmsCampaignLineItem
	for rows.Next() {
		var i OmsCampaignLineItem
		if err := rows.Scan(
			&i.ID,
			&i.CampaignID,
			&i.Name,
			&i.Booked,
			&i.Actual,
			&i.Adjustments,
			&i.StartedAt,
			&i.EndedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCampaignLine = `-- name: UpdateCampaignLine :exec
UPDATE oms.campaign_line_items
SET name = $2, booked = $3, actual = $4, adjustments = $5, started_at = $6, ended_at = $7
WHERE id = $1
`

type UpdateCampaignLineParams struct {
	ID          int32
	Name        string
	Booked      string
	Actual      sql.NullString
	Adjustments sql.NullString
	StartedAt   sql.NullTime
	EndedAt     sql.NullTime
}

func (q *Queries) UpdateCampaignLine(ctx context.Context, arg UpdateCampaignLineParams) error {
	_, err := q.db.ExecContext(ctx, updateCampaignLine,
		arg.ID,
		arg.Name,
		arg.Booked,
		arg.Actual,
		arg.Adjustments,
		arg.StartedAt,
		arg.EndedAt,
	)
	return err
}
