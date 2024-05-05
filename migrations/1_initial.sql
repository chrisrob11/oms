-- +migrate Up

-- Create the schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS oms;

-- Create campaigns table if it does not exist
CREATE TABLE IF NOT EXISTS oms.campaigns (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    archiving BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create campaign_line_items table if it does not exist
CREATE TABLE IF NOT EXISTS oms.campaign_line_items (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER NOT NULL REFERENCES oms.campaigns(id),
    name VARCHAR(255) NOT NULL,
    booked NUMERIC NOT NULL,
    actual NUMERIC,
    adjustments NUMERIC,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for campaign_id on campaign_line_items if it does not exist
CREATE INDEX IF NOT EXISTS idx_line_item_campaign_id ON oms.campaign_line_items(campaign_id);

-- Create invoices table if it does not exist
CREATE TABLE IF NOT EXISTS oms.invoices (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER NOT NULL REFERENCES oms.campaigns(id),
    total_booked_amount NUMERIC,
    total_actual_amount NUMERIC,
    total_adjustments NUMERIC,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    issued_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for campaign_id on invoices if it does not exist
CREATE INDEX IF NOT EXISTS idx_invoice_campaign_id ON oms.invoices(campaign_id);





