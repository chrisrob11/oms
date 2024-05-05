package cmds

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/chrisrob11/oms/internal/oms"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var ListCampaignItemLine = &cli.Command{
	Name:    "list-campaign-line-items",
	Aliases: []string{"lcli"},
	Usage:   "List campaigns of oms",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newListCampaignItemLineCommand(url)
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

type listCampaignItemLineCommand struct {
	serviceURL string
}

func newListCampaignItemLineCommand(serviceURL string) *listCampaignItemLineCommand {
	return &listCampaignItemLineCommand{serviceURL: serviceURL}
}

func (i *listCampaignItemLineCommand) Run(c *cli.Context) error {
	omsClient := client.NewClient(i.serviceURL)

	limit := c.Int("limit")
	token := c.String("token")
	pageThrough := c.Bool("followNextPage")
	req := &client.ListCampaignLineItemRequest{
		Size: limit,
	}

	if token != "" {
		req.Token = &token
	}

	resp, err := omsClient.ListCampaignLineItems(req)
	if err != nil {
		return errors.Wrap(err, "Cannot create campaign line item")
	}

	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print header row
	fmt.Fprintf(w, "ID\tCampaignID\tName\n")

	// Print data rows
	for _, c := range resp.Items {
		fmt.Fprintf(w, "%d\t%d\t%s\n", c.ID, c.CampaignID, c.Name)
	}

	flushErr := w.Flush()
	if flushErr != nil {
		fmt.Printf("Error flushing to stdout: %v", flushErr)
	}

	if resp.NextPageToken != "" {
		decodedToken, decodeErr := oms.DecodeToken(resp.NextPageToken)
		if decodeErr == nil {
			fmt.Printf("PagingToken: Token: %s, Size: %d, StartID: %d\n",
				resp.NextPageToken, decodedToken.Size, decodedToken.StartID)
		}
	}

	if pageThrough {
		for resp.NextPageToken != "" {
			pagingReq := &client.ListCampaignLineItemRequest{
				Token: &resp.NextPageToken,
			}

			resp, err = omsClient.ListCampaignLineItems(pagingReq)
			if err != nil {
				return errors.Wrap(err, "Cannot create campaign line item")
			}

			for _, c := range resp.Items {
				fmt.Fprintf(w, "%d\t%d\t%s\n", c.ID, c.CampaignID, c.Name)
			}

			flushErr := w.Flush()
			if flushErr != nil {
				fmt.Printf("Error flushing to stdout: %v", flushErr)
			}
		}
	}

	return nil
}
