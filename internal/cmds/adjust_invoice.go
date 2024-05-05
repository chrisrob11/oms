package cmds

import (
	"fmt"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var AdjustInvoice = &cli.Command{
	Name:    "adjust-invoice",
	Aliases: []string{"ai"},
	Usage:   "Adjust invoice",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newAdjustInvoiceCommand(url)
		return cmd.Run(c)
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "id",
			Usage: "Id of the campaign",
		},
		&cli.Float64Flag{
			Name: "totalAdjustments",
		},
	},
}

type adjustInvoiceCommand struct {
	serviceURL string
}

func newAdjustInvoiceCommand(serviceURL string) *adjustInvoiceCommand {
	return &adjustInvoiceCommand{serviceURL: serviceURL}
}

func (i *adjustInvoiceCommand) Run(c *cli.Context) error {
	omsClient := client.NewClient(i.serviceURL)

	id := c.Int("id")
	if id == 0 {
		return errMissingID
	}

	totalAdjustments := c.Float64("totalAdjustments")

	if id == 0 {
		return errMissingID
	}

	resp, err := omsClient.AdjustInvoice(&client.AdjustInvoiceRequest{ID: id, TotalAdjustments: totalAdjustments})
	if err != nil {
		return errors.Wrap(err, "Cannot show invoice")
	}

	fmt.Printf("Invoice %d was adjusted\n", resp)

	return nil
}
