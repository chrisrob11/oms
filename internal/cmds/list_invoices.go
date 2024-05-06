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
		&cli.BoolFlag{
			Name: "allFields",
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

	req := i.buildListInvoicesRequest(limit, token)

	resp, err := omsClient.ListInvoices(req)
	if err != nil {
		return errors.Wrap(err, "failed to list invoices")
	}

	allFields := c.Bool("allFields")

	i.printInvoices(resp.Items, true, allFields)

	if pageThrough {
		err = i.paginateInvoices(omsClient, resp.NextPageToken, allFields)
		if err != nil {
			return errors.Wrap(err, "failed to paginate invoices")
		}
	}

	return nil
}

func (i *listInvoicesCommand) buildListInvoicesRequest(limit int, token string) *client.ListInvoicesRequest {
	req := &client.ListInvoicesRequest{
		Size: limit,
	}

	if token != "" {
		req.Token = &token
	}

	return req
}

func (i *listInvoicesCommand) printInvoices(invoices []*models.Invoice, writeHeader, allFields bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() {
		err := w.Flush()
		if err != nil {
			fmt.Printf("Unexpected flush error: %v", err)
		}
	}()

	if writeHeader {
		if allFields {
			//nolint:lll //Why: this is the required headers
			fmt.Fprintf(w, "ID\tCampaignID\tTotalActual\tTotalBooked\tTotalAdjustments\tIssuedAt\tCreatedAt\tUpdatedAt\tStartedAt\tEndedAt\n")
		} else {
			fmt.Fprintf(w, "ID\tCampaignID\tTotalAdjustments\n")
		}
	}

	for _, inv := range invoices {
		if allFields {
			fmt.Fprintf(w, "%d\t%d\t%f\t%f\t%f\t%s\t%s\t%s\t%s\t%s\n", inv.ID, inv.CampaignID, inv.TotalActualAmount,
				inv.TotalBookedAmount, inv.TotalActualAmount, toCompactTime(&inv.IssuedAt),
				toCompactTime(&inv.CreatedAt), toCompactTime(&inv.UpdatedAt),
				toCompactTime(inv.StartedAt), toCompactTime(inv.EndedAt))
		} else {
			fmt.Fprintf(w, "%d\t%d\t%f\n", inv.ID, inv.CampaignID, inv.TotalAdjustments)
		}
	}
}

func (i *listInvoicesCommand) paginateInvoices(omsClient *client.Client, nextPageToken string, allFields bool) error {
	for nextPageToken != "" {
		resp, err := omsClient.ListInvoices(&client.ListInvoicesRequest{Token: &nextPageToken})
		if err != nil {
			return err
		}

		i.printInvoices(resp.Items, false, allFields)
		nextPageToken = resp.NextPageToken
	}

	return nil
}
