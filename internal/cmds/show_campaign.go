package cmds

import (
	"fmt"
	"strconv"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var ShowCampaign = &cli.Command{
	Name:    "show-campaign",
	Aliases: []string{"sc"},
	Usage:   "Show campaign by id of oms",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newShowCampaignCommand(url)
		return cmd.Run(c)
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "id",
			Usage: "Id of the campaign",
		},
	},
}

var errMissingID = errors.New("Missing error id")

type showCampaignCommand struct {
	serviceURL string
}

func newShowCampaignCommand(serviceURL string) *showCampaignCommand {
	return &showCampaignCommand{serviceURL: serviceURL}
}

func (i *showCampaignCommand) Run(c *cli.Context) error {
	omsClient := client.NewClient(i.serviceURL)

	id := c.Int("id")
	if id == 0 {
		return errMissingID
	}

	resp, err := omsClient.ShowCampaign(&client.ShowCampaignRequest{ID: id})
	if err != nil {
		return errors.Wrap(err, "Cannot show campaign")
	}

	fmt.Printf("Campaign\n")
	fmt.Printf("ID:\t\t%d\n", resp.ID)
	fmt.Printf("Archiving:\t%s\n", strconv.FormatBool(resp.Archiving))
	fmt.Printf("CreatedAt:\t%s\n", toCompactTime(&resp.CreatedAt))
	fmt.Printf("EndedAt:\t%s\n", toCompactTime(resp.EndedAt))
	fmt.Printf("Name:\t\t%s\n", resp.Name)
	fmt.Printf("StartedAt:\t%s\n", toCompactTime(resp.StartedAt))
	fmt.Printf("UpdatedAt:\t%s\n", toCompactTime(&resp.UpdatedAt))

	return nil
}
