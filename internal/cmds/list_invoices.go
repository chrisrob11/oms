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

var ListInvoices = &cli.Command{
	Name:    "list-invoices",
	Aliases: []string{"li"},
	Usage:   "List invoices",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newListInvoicesCommand(url)
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

type listInvoicesCommand struct {
	serviceURL string
}

func newListInvoicesCommand(serviceURL string) *listInvoicesCommand {
	return &listInvoicesCommand{serviceURL: serviceURL}
}

func (i *listInvoicesCommand) Run(c *cli.Context) error {
	omsClient := client.NewClient(i.serviceURL)

	limit := c.Int("limit")
	token := c.String("token")
	pageThrough := c.Bool("followNextPage")
	req := &client.ListInvoicesRequest{
		Size: limit,
	}

	if token != "" {
		req.Token = &token
	}

	resp, err := omsClient.ListInvoices(req)
	if err != nil {
		return errors.Wrap(err, "Cannot create campaign line item")
	}

	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print header row
	fmt.Fprintf(w, "ID\tCampaignID\tTotalAdjustments\n")

	// Print data rows
	for _, c := range resp.Items {
		fmt.Fprintf(w, "%d\t%d\t%f\n", c.ID, c.CampaignID, c.TotalAdjustments)
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
			pagingReq := &client.ListInvoicesRequest{
				Token: &resp.NextPageToken,
			}

			resp, err = omsClient.ListInvoices(pagingReq)
			if err != nil {
				return errors.Wrap(err, "Cannot create campaign line item")
			}

			for _, c := range resp.Items {
				fmt.Fprintf(w, "%d\t%d\t%f\n", c.ID, c.CampaignID, c.TotalAdjustments)
			}

			flushErr := w.Flush()
			if flushErr != nil {
				fmt.Printf("Error flushing to stdout: %v", flushErr)
			}
		}
	}

	return nil
}
