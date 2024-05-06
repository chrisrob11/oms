package cmds

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/chrisrob11/oms/internal/oms/models"
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

	req := buildListCampaignRequest(limit, token)
	resp, err := omsClient.ListCampaigns(req)

	if err != nil {
		return errors.Wrap(err, "failed to list campaigns")
	}

	printCampaigns(resp.Items, true)

	if pageThrough {
		err = paginateCampaigns(omsClient, resp.NextPageToken)
		if err != nil {
			return errors.Wrap(err, "failed to paginate campaigns")
		}
	}

	return nil
}

func buildListCampaignRequest(limit int, token string) *client.ListCampaignRequest {
	req := &client.ListCampaignRequest{
		Size: limit,
	}

	if token != "" {
		req.Token = &token
	}

	return req
}

func printCampaigns(campaigns []*models.Campaign, writeHeader bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() {
		err := w.Flush()
		if err != nil {
			fmt.Printf("Unexpected flush error: %v", err)
		}
	}()

	if writeHeader {
		fmt.Fprintf(w, "ID\tName\tArchiving\n")
	}

	for _, c := range campaigns {
		fmt.Fprintf(w, "%d\t%s\t%t\n", c.ID, c.Name, c.Archiving)
	}
}

func paginateCampaigns(omsClient *client.Client, nextPageToken string) error {
	for nextPageToken != "" {
		resp, err := omsClient.ListCampaigns(&client.ListCampaignRequest{Token: &nextPageToken})
		if err != nil {
			return err
		}

		printCampaigns(resp.Items, false)
		nextPageToken = resp.NextPageToken
	}

	return nil
}

func toCompactTime(t *time.Time) string {
	if t == nil {
		return ""
	}

	if t.IsZero() {
		return ""
	}

	return t.Format(time.DateTime)
}
