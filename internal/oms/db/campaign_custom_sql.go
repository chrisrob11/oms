package db

import (
	"context"
	"fmt"
)

type BuildCampaignInvoiceParams struct {
	CampaignID int32
}

// NOTE: not using sqlc as encountered an error, creating by hand in the same pattern
// Probably in future wouldn't use sqlc for this type of app, too much boilerplate and
// not flexible enough.
//
//nolint:lll //Why: valid sql statement that is required to be long
const buildCampaignInvoice = `select oms.campaigns.ID as campaign_id, sum(actual) as total_actual_amount, sum(booked) as total_booked_amount, sum(adjustments) as total_adjustments_amount from oms.campaigns
LEFT JOIN oms.campaign_line_items on oms.campaigns.ID = oms.campaign_line_items.campaign_id
where oms.campaigns.archiving = false and oms.campaigns.ID = %d
group by oms.campaigns.ID
order by oms.campaigns.ID`

func (q *Queries) BuildCampaignInvoice(ctx context.Context,
	arg BuildCampaignInvoiceParams) (CreateInvoiceParams, error) {
	p := CreateInvoiceParams{}
	query := fmt.Sprintf(buildCampaignInvoice, arg.CampaignID)

	row := q.db.QueryRowContext(ctx, query)

	if err := row.Scan(
		&p.CampaignID,
		&p.TotalActualAmount,
		&p.TotalBookedAmount,
		&p.TotalAdjustments,
	); err != nil {
		return CreateInvoiceParams{}, err
	}

	return p, nil
}

// nolint:lll // Why this is a long sql statement that is ok to be long
const resetCampaignSerialID = "SELECT setval(pg_get_serial_sequence('oms.campaigns', 'id'), (SELECT MAX(id) FROM oms.campaigns) + 1);"

func (q *Queries) ResetCampaignID(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetCampaignSerialID)
	return err
}

// nolint:lll // Why this is a long sql statement that is ok to be long
const resetCampaignLineItemSerialID = "SELECT setval(pg_get_serial_sequence('oms.campaign_line_items', 'id'), (SELECT MAX(id) FROM oms.campaign_line_items) + 1);"

func (q *Queries) ResetCampaignLineItemID(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetCampaignLineItemSerialID)
	return err
}
