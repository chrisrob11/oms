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
