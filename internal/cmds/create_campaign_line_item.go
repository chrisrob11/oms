package cmds

import (
	"fmt"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/chrisrob11/oms/internal/oms/models"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var CreateCampaignLineItem = &cli.Command{
	Name:    "create-campaign-line-item",
	Aliases: []string{"ccli"},
	Usage:   "create a campaign line item into oms",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newCreateCampaignLineItemCommand(url)
		return cmd.Run(c)
	},
	Flags: []cli.Flag{
		&cli.Float64Flag{
			Name: "actual",
		},
		&cli.Float64Flag{
			Name: "adjustments",
		},
		&cli.Float64Flag{
			Name: "booked",
		},
		&cli.IntFlag{
			Name: "id",
		},
		&cli.IntFlag{
			Name: "campaignId",
		},
	},
}

type createCampaignLineItemCommand struct {
	serviceURL string
}

func newCreateCampaignLineItemCommand(serviceURL string) *createCampaignLineItemCommand {
	return &createCampaignLineItemCommand{serviceURL: serviceURL}
}

func (i *createCampaignLineItemCommand) Run(c *cli.Context) error {
	campaignLineItem := &models.CampaignLineItem{}
	campaignLineItem.Actual = c.Float64("actual")
	campaignLineItem.Adjustments = c.Float64("adjustments")
	campaignLineItem.Booked = c.Float64("booked")
	campaignLineItem.ID = c.Int("id")
	campaignLineItem.CampaignID = c.Int("campaignId")

	omsClient := client.NewClient(i.serviceURL)

	id, err := omsClient.CreateCampaignOrderLine(campaignLineItem)
	if err != nil {
		return errors.Wrap(err, "Cannot create campaign")
	}

	fmt.Printf("CampaignOrderLine with ID %d was created\n", id)

	return nil
}
