-- campaign_line.sql

-- name: CreateCampaignLineWithID :one
INSERT INTO oms.campaign_line_items (id, campaign_id, name, booked, actual, adjustments, started_at, ended_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: CreateCampaignLine :one
INSERT INTO oms.campaign_line_items (campaign_id, name, booked, actual, adjustments, started_at, ended_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: GetCampaignLine :one
SELECT * FROM oms.campaign_line_items WHERE id = $1;

-- name: UpdateCampaignLine :exec
UPDATE oms.campaign_line_items
SET name = $2, booked = $3, actual = $4, adjustments = $5, started_at = $6, ended_at = $7
WHERE id = $1;

-- name: DeleteCampaignLine :exec
DELETE FROM oms.campaign_line_items WHERE id = $1;

-- name: ListCampaignLineItems :many
SELECT * FROM oms.campaign_line_items 
WHERE id > $1
Order by id
LIMIT $2;