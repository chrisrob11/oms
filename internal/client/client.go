// Package client is a strongly typed client for interacting with
// the service to do operations
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/chrisrob11/oms/internal/oms/models"
	"github.com/pkg/errors"
)

var errUnexpectedStatusCode = fmt.Errorf("unexpected status code")

func newErrUnexpectedStatusCode(resp *http.Response) error {
	return errors.Wrapf(errUnexpectedStatusCode, "status code: %s", resp.Status)
}

// Client represents the HTTP client for the API.
type Client struct {
	BaseURL string
	logger  *slog.Logger
}

// NewClient creates a new instance of the API client.
func NewClient(baseURL string) *Client {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	return &Client{BaseURL: baseURL, logger: logger}
}

// createResource sends a POST request to create a new resource.
func (c *Client) createResource(endpoint string, data interface{}, out interface{}) error {
	reqBody, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	resp, err := http.Post(c.BaseURL+endpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}

	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			c.logger.Warn("Error closing body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return newErrUnexpectedStatusCode(resp)
	}

	outData, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read body")
	}

	err = json.Unmarshal(outData, out)
	if err != nil {
		return errors.Wrap(err, "cannot unmarshal body to output value")
	}

	return nil
}

// executeAction sends a POST request to do an action on a resource.
func (c *Client) executeAction(endpoint string, id int, action string, data interface{}, out interface{}) error {
	var dataReader io.Reader

	if data != nil {
		reqBody, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error encoding JSON: %w", err)
		}

		dataReader = bytes.NewBuffer(reqBody)
	}

	resp, err := http.Post(c.BaseURL+endpoint+"/"+strconv.Itoa(id)+"/"+action, "application/json", dataReader)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}

	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			c.logger.Warn("Error closing body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return newErrUnexpectedStatusCode(resp)
	}

	outData, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read body")
	}

	err = json.Unmarshal(outData, out)
	if err != nil {
		return errors.Wrap(err, "cannot unmarshal body to output value")
	}

	return nil
}

// listResources sends a Get request to list resources.
func (c *Client) listResources(endpoint string, token *string, limit *int, out interface{}) error {
	queryValues := url.Values{}
	if token != nil {
		queryValues.Add("$token", *token)
	}

	if limit != nil {
		queryValues.Add("$limit", strconv.Itoa(*limit))
	}

	query := c.BaseURL + endpoint
	if len(queryValues) > 0 {
		query = query + "?" + queryValues.Encode()
	}

	//nolint:gosec // Why: controlled input above.
	resp, err := http.Get(query)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}

	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			c.logger.Warn("Error closing body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return newErrUnexpectedStatusCode(resp)
	}

	outData, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read body")
	}

	err = json.Unmarshal(outData, out)
	if err != nil {
		return errors.Wrap(err, "cannot unmarshal body to output value")
	}

	return nil
}

// listResources sends a Get request to list resources.
func (c *Client) showResources(endpoint string, id int, out interface{}) error {
	resp, err := http.Get(c.BaseURL + endpoint + "/" + strconv.Itoa(id))
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}

	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			c.logger.Warn("Error closing body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return newErrUnexpectedStatusCode(resp)
	}

	outData, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read body")
	}

	err = json.Unmarshal(outData, out)
	if err != nil {
		return errors.Wrap(err, "cannot unmarshal body to output value")
	}

	return nil
}

// CreateCampaign sends a POST request to create a new campaign.
func (c *Client) CreateCampaign(campaign *models.Campaign) (int, error) {
	outCampaign := models.Campaign{}

	err := c.createResource("/campaigns", campaign, &outCampaign)
	if err != nil {
		return 0, err
	}

	return outCampaign.ID, nil
}

// CreateCampaignOrderLine sends a POST request to create a new campaign order line.
func (c *Client) CreateCampaignOrderLine(orderLine *models.CampaignLineItem) (int, error) {
	var id int

	err := c.createResource("/campaignLineItems", orderLine, &id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

type ListCampaignRequest struct {
	Size  int
	Token *string
}

// ListCampaigns sends a Get request to get a list of campaigns.
func (c *Client) ListCampaigns(req *ListCampaignRequest) (*models.List[models.Campaign], error) {
	if req.Size == 0 {
		req.Size = 100
	}

	items := &models.List[models.Campaign]{}

	err := c.listResources("/campaigns", req.Token, &req.Size, items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type ListCampaignLineItemRequest struct {
	Size  int
	Token *string
}

func (c *Client) ListCampaignLineItems(
	req *ListCampaignLineItemRequest) (*models.List[models.CampaignLineItem], error) {
	if req.Size == 0 {
		req.Size = 100
	}

	items := &models.List[models.CampaignLineItem]{}

	err := c.listResources("/campaignLineItems", req.Token, &req.Size, items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type ShowCampaignRequest struct {
	ID int
}

// ShowCampaign sends a Get request to get a specific campaigns.
func (c *Client) ShowCampaign(req *ShowCampaignRequest) (*models.Campaign, error) {
	campaign := models.Campaign{}

	err := c.showResources("/campaigns", req.ID, &campaign)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

type ShowCampaignOrderLineRequest struct {
	ID int
}

func (c *Client) ShowCampaignLineItem(req *ShowCampaignOrderLineRequest) (*models.CampaignLineItem, error) {
	campaignLineItem := models.CampaignLineItem{}

	err := c.showResources("/campaignLineItems", req.ID, &campaignLineItem)
	if err != nil {
		return nil, err
	}

	return &campaignLineItem, nil
}

type ShowInvoiceRequest struct {
	ID int
}

// ShowInvoice sends a Get request to get a specific invoice.
func (c *Client) ShowInvoice(req *ShowInvoiceRequest) (*models.Invoice, error) {
	invoice := models.Invoice{}

	err := c.showResources("/invoices", req.ID, &invoice)
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

type AdjustInvoiceRequest struct {
	ID               int
	TotalAdjustments float64
}

// AdjustInvoice sends an adjustment for an invoice.
func (c *Client) AdjustInvoice(req *AdjustInvoiceRequest) (int, error) {
	var outID int

	err := c.executeAction("/invoices", req.ID, "adjust", &req.TotalAdjustments, &outID)
	if err != nil {
		return 0, err
	}

	return outID, nil
}

type GenerateInvoiceFromCampaignRequest struct {
	ID int
}

func (c *Client) GenerateInvoiceFromCampaign(req *GenerateInvoiceFromCampaignRequest) (int, error) {
	var intValue int

	err := c.executeAction("/campaigns", req.ID, "generateInvoice", nil, &intValue)
	if err != nil {
		return intValue, err
	}

	return intValue, nil
}

type ListInvoicesRequest struct {
	Size  int
	Token *string
}

func (c *Client) ListInvoices(req *ListInvoicesRequest) (*models.List[models.Invoice], error) {
	if req.Size == 0 {
		req.Size = 100
	}

	items := &models.List[models.Invoice]{}

	err := c.listResources("/invoices", req.Token, &req.Size, items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
