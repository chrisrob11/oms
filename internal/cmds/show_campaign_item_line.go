package cmds

import (
	"fmt"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var ShowCampaignItemLine = &cli.Command{
	Name:    "show-campaign-line-item",
	Aliases: []string{"scli"},
	Usage:   "Show campaign line item by id",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newShowCampaignItemLineCommand(url)
		return cmd.Run(c)
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "id",
			Usage: "Id of the campaign",
		},
	},
}

type showCampaignItemLineCommand struct {
	serviceURL string
}

func newShowCampaignItemLineCommand(serviceURL string) *showCampaignItemLineCommand {
	return &showCampaignItemLineCommand{serviceURL: serviceURL}
}

func (i *showCampaignItemLineCommand) Run(c *cli.Context) error {
	omsClient := client.NewClient(i.serviceURL)

	id := c.Int("id")
	if id == 0 {
		return errMissingID
	}

	resp, err := omsClient.ShowCampaignLineItem(&client.ShowCampaignOrderLineRequest{ID: id})
	if err != nil {
		return errors.Wrap(err, "Cannot show campaign line item")
	}

	fmt.Printf("CampaignLineItem\n")
	fmt.Printf("ID:\t\t%d\n", resp.ID)
	fmt.Printf("CampaignID:\t%d\n", resp.CampaignID)
	fmt.Printf("Name:\t\t%s\n", resp.Name)
	fmt.Printf("Actual:\t\t%f\n", resp.Actual)
	fmt.Printf("Booked:\t\t%f\n", resp.Booked)
	fmt.Printf("Adjustments:\t%f\n", resp.Adjustments)
	fmt.Printf("CreatedAt:\t%s\n", toCompactTime(&resp.CreatedAt))
	fmt.Printf("EndedAt:\t%s\n", toCompactTime(resp.EndedAt))
	fmt.Printf("StartedAt:\t%s\n", toCompactTime(resp.StartedAt))
	fmt.Printf("UpdatedAt:\t%s\n", toCompactTime(&resp.UpdatedAt))

	return nil
}
