-- campaign.sql

-- name: CreateCampaignWithID :one
INSERT INTO oms.campaigns (name, id, started_at, ended_at, archiving)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: CreateCampaign :one
INSERT INTO oms.campaigns (name, started_at, ended_at, archiving)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCampaign :one
SELECT * FROM oms.campaigns WHERE id = $1;

-- name: ListCampaigns :many
SELECT * FROM oms.campaigns 
WHERE archiving = false AND id > $1
Order by id
LIMIT $2;

-- name: UpdateCampaign :exec
UPDATE oms.campaigns
SET name = $1, started_at = $2, ended_at = $3, archiving = $4
WHERE id = $5;

-- name: DeleteCampaign :exec 
DELETE FROM oms.campaigns WHERE id = $1;