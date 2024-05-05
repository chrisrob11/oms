-- invoice.sql

-- name: CreateInvoice :one
INSERT INTO oms.invoices (campaign_id, total_booked_amount, total_actual_amount, total_adjustments, started_at, ended_at, issued_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: GetInvoice :one
SELECT * FROM oms.invoices WHERE id = $1;

-- name: ListInvoices :many
SELECT * FROM oms.invoices 
WHERE id > $1
Order by id
LIMIT $2;

-- name: AdjustInvoice :exec
UPDATE oms.invoices
SET total_adjustments = $2
WHERE id = $1;

-- name: DeleteInvoice :exec
DELETE FROM oms.invoices WHERE id = $1;