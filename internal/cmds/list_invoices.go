package cmds

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/chrisrob11/oms/internal/oms/models"
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

	req := buildListInvoicesRequest(limit, token)

	resp, err := omsClient.ListInvoices(req)
	if err != nil {
		return errors.Wrap(err, "failed to list invoices")
	}

	printInvoices(resp.Items, true)

	if pageThrough {
		err = paginateInvoices(omsClient, resp.NextPageToken)
		if err != nil {
			return errors.Wrap(err, "failed to paginate invoices")
		}
	}

	return nil
}

func buildListInvoicesRequest(limit int, token string) *client.ListInvoicesRequest {
	req := &client.ListInvoicesRequest{
		Size: limit,
	}

	if token != "" {
		req.Token = &token
	}

	return req
}

func printInvoices(invoices []*models.Invoice, writeHeader bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() {
		err := w.Flush()
		if err != nil {
			fmt.Printf("Unexpected flush error: %v", err)
		}
	}()

	if writeHeader {
		fmt.Fprintf(w, "ID\tCampaignID\tTotalAdjustments\n")
	}

	for _, invoice := range invoices {
		fmt.Fprintf(w, "%d\t%d\t%f\n", invoice.ID, invoice.CampaignID, invoice.TotalAdjustments)
	}
}

func paginateInvoices(omsClient *client.Client, nextPageToken string) error {
	for nextPageToken != "" {
		resp, err := omsClient.ListInvoices(&client.ListInvoicesRequest{Token: &nextPageToken})
		if err != nil {
			return err
		}

		printInvoices(resp.Items, false)
		nextPageToken = resp.NextPageToken
	}

	return nil
}
