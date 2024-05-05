// Package models has the main models and how to convert them to the db models back and forth
package models

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/chrisrob11/oms/internal/oms/db"
	"github.com/pkg/errors"
)

// Campaign represents the data structure for a campaign.
type Campaign struct {
	ID        int `json:"omitempty"`
	Name      string
	StartedAt *time.Time
	EndedAt   *time.Time
	Archiving bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCampaignFromDB(c *db.OmsCampaign) *Campaign {
	return &Campaign{
		ID:        int(c.ID),
		Name:      c.Name,
		StartedAt: toTime(c.StartedAt),
		EndedAt:   toTime(c.EndedAt),
		Archiving: c.Archiving.Bool,
		CreatedAt: *toTime(c.CreatedAt),
		UpdatedAt: *toTime(c.UpdatedAt),
	}
}

func (c *Campaign) ToCreateCampaignWithID() *db.CreateCampaignWithIDParams {
	return &db.CreateCampaignWithIDParams{
		ID:        int32(c.ID),
		Name:      c.Name,
		StartedAt: toSQLTime(c.StartedAt),
		EndedAt:   toSQLTime(c.EndedAt),
		Archiving: sql.NullBool{Valid: true, Bool: c.Archiving},
	}
}

func (c *Campaign) ToCreateCampaign() *db.CreateCampaignParams {
	return &db.CreateCampaignParams{
		Name:      c.Name,
		StartedAt: toSQLTime(c.StartedAt),
		EndedAt:   toSQLTime(c.EndedAt),
		Archiving: sql.NullBool{Valid: true, Bool: c.Archiving},
	}
}

// CampaignLineItem represents the data structure for a campaign order line.
type CampaignLineItem struct {
	ID          int
	CampaignID  int
	Name        string
	Booked      float64
	Actual      float64
	Adjustments float64
	StartedAt   *time.Time
	EndedAt     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCampaignLineItemFromDB(c *db.OmsCampaignLineItem) (*CampaignLineItem, error) {
	booked, err := strconv.ParseFloat(c.Booked, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "string booked not convertable to float: %s", c.Booked)
	}

	var actual float64
	if c.Actual.Valid {
		actual, err = strconv.ParseFloat(c.Actual.String, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "string actual not convertable to float: %s", c.Actual.String)
		}
	}

	var adjustments float64

	if c.Adjustments.Valid {
		actual, err = strconv.ParseFloat(c.Adjustments.String, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "string adjustments not convertable to float: %s", c.Adjustments.String)
		}
	}

	var createdAt time.Time
	if c.CreatedAt.Valid {
		createdAt = c.CreatedAt.Time
	}

	var updatedAt time.Time

	if c.UpdatedAt.Valid {
		createdAt = c.UpdatedAt.Time
	}

	return &CampaignLineItem{
		ID:          int(c.ID),
		CampaignID:  int(c.CampaignID),
		Name:        c.Name,
		Booked:      booked,
		Actual:      actual,
		Adjustments: adjustments,
		StartedAt:   toTime(c.StartedAt),
		EndedAt:     toTime(c.EndedAt),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (c *CampaignLineItem) ToCreateCampaignLineItemWithID() db.CreateCampaignLineWithIDParams {
	return db.CreateCampaignLineWithIDParams{
		ID:          int32(c.ID),
		CampaignID:  int32(c.CampaignID),
		Name:        c.Name,
		Booked:      strconv.FormatFloat(c.Booked, 'f', 16, 64),
		Actual:      sql.NullString{String: strconv.FormatFloat(c.Actual, 'f', 16, 64), Valid: true},
		Adjustments: sql.NullString{String: strconv.FormatFloat(c.Adjustments, 'f', 16, 64), Valid: true},
		StartedAt:   toSQLTime(c.StartedAt),
		EndedAt:     toSQLTime(c.EndedAt),
	}
}

func (c *CampaignLineItem) ToCreateCampaignLineItem() db.CreateCampaignLineParams {
	return db.CreateCampaignLineParams{
		CampaignID:  int32(c.CampaignID),
		Name:        c.Name,
		Booked:      strconv.FormatFloat(c.Booked, 'f', 16, 64),
		Actual:      sql.NullString{String: strconv.FormatFloat(c.Actual, 'f', 16, 64), Valid: true},
		Adjustments: sql.NullString{String: strconv.FormatFloat(c.Adjustments, 'f', 16, 64), Valid: true},
		StartedAt:   toSQLTime(c.StartedAt),
		EndedAt:     toSQLTime(c.EndedAt),
	}
}

type Invoice struct {
	ID                int
	CampaignID        int
	TotalBookedAmount float64
	TotalActualAmount float64
	TotalAdjustments  float64
	StartedAt         *time.Time
	EndedAt           *time.Time
	IssuedAt          time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (i *Invoice) ToCreateInvoiceParams() db.CreateInvoiceParams {
	return db.CreateInvoiceParams{
		CampaignID:        int32(i.CampaignID),
		TotalActualAmount: toSQLStringFromFloat64(i.TotalActualAmount),
		TotalBookedAmount: toSQLStringFromFloat64(i.TotalBookedAmount),
		TotalAdjustments:  toSQLStringFromFloat64(i.TotalAdjustments),
		StartedAt:         toSQLTime(i.StartedAt),
		EndedAt:           toSQLTime(i.EndedAt),
	}
}

func NewInvoiceFromDB(i db.OmsInvoice) (*Invoice, error) {
	var err error

	var actual float64

	if i.TotalActualAmount.Valid {
		actual, err = strconv.ParseFloat(i.TotalActualAmount.String, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "string TotalActualAmount not convertable to float: %s", i.TotalActualAmount.String)
		}
	}

	var booked float64

	if i.TotalBookedAmount.Valid {
		booked, err = strconv.ParseFloat(i.TotalBookedAmount.String, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "string TotalBookedAmount not convertable to float: %s", i.TotalBookedAmount.String)
		}
	}

	var adjustment float64

	if i.TotalAdjustments.Valid {
		adjustment, err = strconv.ParseFloat(i.TotalAdjustments.String, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "string TotalAdjustments not convertable to float: %s", i.TotalAdjustments.String)
		}
	}

	return &Invoice{
		ID:                int(i.ID),
		CampaignID:        int(i.CampaignID),
		TotalActualAmount: actual,
		TotalBookedAmount: booked,
		TotalAdjustments:  adjustment,
		StartedAt:         toTime(i.StartedAt),
		EndedAt:           toTime(i.EndedAt),
		CreatedAt:         i.CreatedAt.Time,
		UpdatedAt:         i.EndedAt.Time,
	}, nil
}

func toSQLStringFromFloat64(t float64) sql.NullString {
	strValue := strconv.FormatFloat(t, 'b', 15, 64)
	return sql.NullString{Valid: true, String: strValue}
}

func toSQLTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Valid: true, Time: *t}
}

func toTime(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}

	return nil
}

type Paging struct {
	Token string
	Size  int
}

type List[T any] struct {
	Items         []*T
	NextPageToken string `json:"omitempty"`
}
