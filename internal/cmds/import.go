package cmds

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/chrisrob11/oms/internal/oms/models"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	ErrMissingParameter    = errors.New("missing parameter")
	errUnexpectedJSONToken = errors.New("Unexpected json token")
)

var Import = &cli.Command{
	Name:    "import",
	Aliases: []string{"i"},
	Usage:   "Import data into oms",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}
		filePath := c.String("file")
		if filePath == "" {
			return NewMissingError("file")
		}

		cmd := newImportCommand(url, filePath)
		return cmd.Run(c)
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "file",
		},
		&cli.BoolFlag{
			Name: "generateInvoices",
		},
	},
}

func NewMissingError(parameterName string) error {
	return fmt.Errorf("%w: %s", ErrMissingParameter, parameterName)
}

type importCommand struct {
	serviceURL string
	filePath   string
}

func newImportCommand(serviceURL, filePath string) *importCommand {
	return &importCommand{serviceURL: serviceURL, filePath: filePath}
}

func (i *importCommand) Run(c *cli.Context) error {
	file, err := os.Open(i.filePath)
	if err != nil {
		return errors.Wrapf(err, "Cannot open file at path %s", i.filePath)
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			fmt.Printf("Error closing import file: %s\n", err)
		}
	}()

	campaignsCreatedIDMap := map[int]struct{}{}
	campaignLineItemsCount := 0
	omsClient := client.NewClient(i.serviceURL)

	decoder, err := i.initializeJSONDecoder(file)
	if err != nil {
		return err
	}

	genInvoices := c.Bool("generateInvoices")

	for decoder.More() {
		var item map[string]interface{}
		if decodeErr := decoder.Decode(&item); decodeErr != nil {
			return errors.Wrap(err, "Cannot decode json")
		}

		uploadErr := uploadData(omsClient, campaignsCreatedIDMap, item)
		if uploadErr != nil {
			return err
		}
		campaignLineItemsCount++
	}

	createdInvoices := 0

	if genInvoices {
		createdInvoices, err = generateInvoices(omsClient, campaignsCreatedIDMap)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Created %d campaigns\n", len(campaignsCreatedIDMap))
	fmt.Printf("Created %d campaign order lines\n", campaignLineItemsCount)
	fmt.Printf("Created %d invoices\n", createdInvoices)

	return nil
}

func uploadData(omsClient *client.Client, campaignsCreatedIDMap map[int]struct{}, item map[string]interface{}) error {
	campaign, campaignLineItem := processLine(item)
	campaignID := campaign.ID

	if _, exists := campaignsCreatedIDMap[campaign.ID]; !exists {
		_, err := omsClient.CreateCampaign(campaign)
		if err != nil {
			return errors.Wrapf(err, "Cannot create a campaign")
		}

		campaignsCreatedIDMap[campaignID] = struct{}{}
	}

	_, err := omsClient.CreateCampaignOrderLine(campaignLineItem)
	if err != nil {
		return errors.Wrapf(err, "Cannot create a campaign line item")
	}

	return nil
}

func generateInvoices(omsClient *client.Client, campaignsCreatedIDMap map[int]struct{}) (int, error) {
	createdInvoices := 0

	for campaignID := range campaignsCreatedIDMap {
		_, err := omsClient.GenerateInvoiceFromCampaign(&client.GenerateInvoiceFromCampaignRequest{ID: campaignID})
		if err != nil {
			return createdInvoices, errors.Wrapf(err, "Cannot create a invoice for campaign line: %d", campaignID)
		}
		createdInvoices++
	}

	return createdInvoices, nil
}

func (i *importCommand) initializeJSONDecoder(reader io.Reader) (*json.Decoder, error) {
	decoder := json.NewDecoder(reader)

	token, err := decoder.Token()
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading JSON token")
	}

	// Check if the token is an array token
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return nil, errUnexpectedJSONToken
	}

	return decoder, nil
}

func processLine(item map[string]interface{}) (*models.Campaign, *models.CampaignLineItem) {
	// Process the line item
	id := item["id"].(float64)
	campaignID := item["campaign_id"].(float64)
	campaignName := item["campaign_name"].(string)
	lineItemName := item["line_item_name"].(string)
	bookedAmount := item["booked_amount"].(float64)
	actualAmount := item["actual_amount"].(float64)
	adjustments := item["adjustments"].(float64)

	// ideally should have more of a bulk apis to do this.
	// but just doing this "the longer way" at the moment,
	// will perf optimize if bring in 10k is slower
	campaign := &models.Campaign{
		ID:   int(campaignID),
		Name: campaignName,
	}

	campaignLine := &models.CampaignLineItem{
		ID:          int(id),
		CampaignID:  int(campaignID),
		Name:        lineItemName,
		Booked:      bookedAmount,
		Actual:      actualAmount,
		Adjustments: adjustments,
	}

	return campaign, campaignLine
}
