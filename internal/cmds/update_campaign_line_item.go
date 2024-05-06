package cmds

import (
	"fmt"
	"time"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var UpdateCampaignLineItem = &cli.Command{
	Name:    "update-campaign-line-item",
	Aliases: []string{"ucli"},
	Usage:   "Update a campaign line item",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newUpdateCampaignLineItemCommand(url)
		return cmd.Run(c)
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "id",
			Usage: "Id of the campaign",
		},
		&cli.StringFlag{
			Name: "name",
		},
		&cli.TimestampFlag{
			Name:   "startedAt",
			Layout: time.DateTime,
		},
		&cli.TimestampFlag{
			Name:   "endedAt",
			Layout: time.DateTime,
		},
		&cli.Float64Flag{
			Name: "booked",
		},
		&cli.Float64Flag{
			Name: "actual",
		},
		&cli.Float64Flag{
			Name: "adjustments",
		},
	},
}

type updateCampaignLineItemCommand struct {
	serviceURL string
}

func newUpdateCampaignLineItemCommand(serviceURL string) *updateCampaignLineItemCommand {
	return &updateCampaignLineItemCommand{serviceURL: serviceURL}
}

func (i *updateCampaignLineItemCommand) Run(c *cli.Context) error {
	omsClient := client.NewClient(i.serviceURL)

	id := c.Int("id")
	if id == 0 {
		return errMissingID
	}

	foundCampaignLineItem, err := omsClient.ShowCampaignLineItem(&client.ShowCampaignOrderLineRequest{
		ID: id,
	})

	if err != nil {
		return errors.Wrap(err, "Error finding Campaign")
	}

	name := c.String("name")
	if name != "" {
		foundCampaignLineItem.Name = name
	}

	startedAt := c.Timestamp("startedAt")
	if startedAt != nil {
		foundCampaignLineItem.StartedAt = startedAt
	}

	endedAt := c.Timestamp("endedAt")
	if endedAt != nil {
		foundCampaignLineItem.EndedAt = endedAt
	}

	booked := c.Float64("booked")
	// NOTE: this is kind of a bug, we might want any of these to actually be 0
	// just noting it. Probably would need to add another flag to handle this
	if booked != 0 {
		foundCampaignLineItem.Booked = booked
	}

	actual := c.Float64("actual")
	if actual != 0 {
		foundCampaignLineItem.Actual = actual
	}

	adjustments := c.Float64("adjustments")
	if adjustments != 0 {
		foundCampaignLineItem.Adjustments = adjustments
	}

	err = omsClient.UpdateCampaignLineItem(*foundCampaignLineItem)
	if err != nil {
		return errors.Wrap(err, "Cannot update campaign")
	}

	fmt.Println("Update processed.")

	return nil
}
