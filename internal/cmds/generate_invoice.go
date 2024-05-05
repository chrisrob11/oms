package cmds

import (
	"fmt"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var GenerateInvoice = &cli.Command{
	Name:    "generate-invoice",
	Aliases: []string{"gi"},
	Usage:   "create an invoice from a specified campaign",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newGenerateInvoice(url)
		return cmd.Run(c)
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name: "id",
		},
	},
}

type generateInvoice struct {
	serviceURL string
}

func newGenerateInvoice(serviceURL string) *generateInvoice {
	return &generateInvoice{serviceURL: serviceURL}
}

func (i *generateInvoice) Run(c *cli.Context) error {
	id := c.Int("id")
	if id == 0 {
		return ErrMissingParameter
	}

	omsClient := client.NewClient(i.serviceURL)

	invoiceID, err := omsClient.GenerateInvoiceFromCampaign(&client.GenerateInvoiceFromCampaignRequest{
		ID: id,
	})
	if err != nil {
		return errors.Wrap(err, "Cannot generate invoice")
	}

	fmt.Printf("Invoice with ID %d was created\n", invoiceID)

	return nil
}
