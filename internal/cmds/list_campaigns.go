package cmds

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var ListCampaign = &cli.Command{
	Name:    "list-campaign",
	Aliases: []string{"lc"},
	Usage:   "List campaigns of oms",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newListCampaignCommand(url)
		return cmd.Run(c)
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "limit",
			Value: 500,
		},
		&cli.StringFlag{
			Name: "token",
		},
		&cli.BoolFlag{
			Name:    "followNextPage",
			Aliases: []string{"fnp"},
		},
	},
}

type listCampaignCommand struct {
	serviceURL string
}

func newListCampaignCommand(serviceURL string) *listCampaignCommand {
	return &listCampaignCommand{serviceURL: serviceURL}
}

func (i *listCampaignCommand) Run(c *cli.Context) error {
	omsClient := client.NewClient(i.serviceURL)

	limit := c.Int("limit")
	token := c.String("token")
	pageThrough := c.Bool("followNextPage")
	req := &client.ListCampaignRequest{
		Size: limit,
	}

	if token != "" {
		req.Token = &token
	}

	resp, err := omsClient.ListCampaigns(req)
	if err != nil {
		return errors.Wrap(err, "Cannot create campaign")
	}

	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print header row
	fmt.Fprintf(w, "ID\tName\tArchiving\n")

	// Print data rows
	for _, c := range resp.Items {
		fmt.Fprintf(w, "%d\t%s\t%s\n", c.ID, c.Name, strconv.FormatBool(c.Archiving))
	}

	err = w.Flush()
	if err != nil {
		return errors.Wrap(err, "unexpected flush error")
	}

	if pageThrough {
		for resp.NextPageToken != "" {
			pagingReq := &client.ListCampaignRequest{
				Token: &resp.NextPageToken,
			}

			resp, err = omsClient.ListCampaigns(pagingReq)
			if err != nil {
				return errors.Wrap(err, "Cannot create campaign line item")
			}

			for _, c := range resp.Items {
				fmt.Fprintf(w, "%d\t%s\t%s\n", c.ID, c.Name, strconv.FormatBool(c.Archiving))
			}

			flushErr := w.Flush()
			if flushErr != nil {
				return errors.Wrapf(err, "unable to flush data")
			}
		}
	}

	return nil
}

func toCompactTime(t *time.Time) string {
	if t == nil {
		return ""
	}

	return t.Format(time.DateTime)
}
