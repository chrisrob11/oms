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

	req := buildListCampaignLineItemRequest(limit, token)

	resp, err := omsClient.ListCampaignLineItems(req)
	if err != nil {
		return errors.Wrap(err, "failed to list campaign line items")
	}

	printCampaignLineItems(resp.Items, true)

	if pageThrough {
		err = paginateCampaignLineItems(omsClient, resp.NextPageToken)
		if err != nil {
			return errors.Wrap(err, "failed to paginate campaign line items")
		}
	}

	return nil
}

func buildListCampaignLineItemRequest(limit int, token string) *client.ListCampaignLineItemRequest {
	req := &client.ListCampaignLineItemRequest{
		Size: limit,
	}

	if token != "" {
		req.Token = &token
	}

	return req
}

func printCampaignLineItems(lineItems []*models.CampaignLineItem, writeHeader bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() {
		err := w.Flush()
		if err != nil {
			fmt.Printf("Unexpected flush error: %v", err)
		}
	}()

	if writeHeader {
		fmt.Fprintf(w, "ID\tCampaignID\tName\n")
	}

	for _, c := range lineItems {
		fmt.Fprintf(w, "%d\t%d\t%s\n", c.ID, c.CampaignID, c.Name)
	}
}

func paginateCampaignLineItems(omsClient *client.Client, nextPageToken string) error {
	for nextPageToken != "" {
		resp, err := omsClient.ListCampaignLineItems(&client.ListCampaignLineItemRequest{Token: &nextPageToken})
		if err != nil {
			return err
		}

		printCampaignLineItems(resp.Items, false)
		nextPageToken = resp.NextPageToken
	}

	return nil
}
