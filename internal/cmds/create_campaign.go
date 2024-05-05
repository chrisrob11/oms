package cmds

import (
	"fmt"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/chrisrob11/oms/internal/oms/models"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var CreateCampaign = &cli.Command{
	Name:    "create-campaign",
	Aliases: []string{"cc"},
	Usage:   "create a campaign into oms",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newCreateCampaignCommand(url)
		return cmd.Run(c)
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "archiving",
		},
		&cli.IntFlag{
			Name: "id",
		},
		&cli.StringFlag{
			Name: "name",
		},
	},
}

type createCampaignCommand struct {
	serviceURL string
}

func newCreateCampaignCommand(serviceURL string) *createCampaignCommand {
	return &createCampaignCommand{serviceURL: serviceURL}
}

func (i *createCampaignCommand) Run(c *cli.Context) error {
	campaign := &models.Campaign{}
	campaign.Archiving = c.Bool("archiving")
	campaign.ID = c.Int("id")
	campaign.Name = c.String("name")
	campaign.StartedAt = c.Timestamp("startedAt")
	campaign.EndedAt = c.Timestamp("endedAt")

	omsClient := client.NewClient(i.serviceURL)

	id, err := omsClient.CreateCampaign(campaign)
	if err != nil {
		return errors.Wrap(err, "Cannot create campaign")
	}

	fmt.Printf("Campaign with ID %d was created\n", id)

	return nil
}
