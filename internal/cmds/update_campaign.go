package cmds

import (
	"fmt"
	"time"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var UpdateCampaign = &cli.Command{
	Name:    "update-campaign",
	Aliases: []string{"uc"},
	Usage:   "Update a campaign",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newUpdateCampaignCommand(url)
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
	},
}

type updateCampaignCommand struct {
	serviceURL string
}

func newUpdateCampaignCommand(serviceURL string) *updateCampaignCommand {
	return &updateCampaignCommand{serviceURL: serviceURL}
}

func (i *updateCampaignCommand) Run(c *cli.Context) error {
	omsClient := client.NewClient(i.serviceURL)

	id := c.Int("id")
	if id == 0 {
		return errMissingID
	}

	foundCampaign, err := omsClient.ShowCampaign(&client.ShowCampaignRequest{
		ID: id,
	})

	if err != nil {
		return errors.Wrap(err, "Error finding Campaign")
	}

	name := c.String("name")
	if name != "" {
		foundCampaign.Name = name
	}

	startedAt := c.Timestamp("startedAt")
	if startedAt != nil {
		foundCampaign.StartedAt = startedAt
	}

	endedAt := c.Timestamp("endedAt")
	if endedAt != nil {
		foundCampaign.EndedAt = endedAt
	}

	err = omsClient.UpdateCampaign(*foundCampaign)
	if err != nil {
		return errors.Wrap(err, "Cannot update campaign")
	}

	fmt.Println("Update processed.")

	return nil
}
