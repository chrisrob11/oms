package cmds

import (
	"fmt"

	"github.com/chrisrob11/oms/internal/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var ShowInvoice = &cli.Command{
	Name:    "show-invoice",
	Aliases: []string{"si"},
	Usage:   "Show an invoice",
	Action: func(c *cli.Context) error {
		url := c.String("url")
		if url == "" {
			return NewMissingError("url")
		}

		cmd := newShowInvoiceCommand(url)
		return cmd.Run(c)
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "id",
			Usage: "Id of the campaign",
		},
	},
}

type showInvoiceCommand struct {
	serviceURL string
}

func newShowInvoiceCommand(serviceURL string) *showInvoiceCommand {
	return &showInvoiceCommand{serviceURL: serviceURL}
}

func (i *showInvoiceCommand) Run(c *cli.Context) error {
	omsClient := client.NewClient(i.serviceURL)

	id := c.Int("id")
	if id == 0 {
		return errMissingID
	}

	resp, err := omsClient.ShowInvoice(&client.ShowInvoiceRequest{ID: id})
	if err != nil {
		return errors.Wrap(err, "Cannot show invoice")
	}

	fmt.Printf("Invoice\n")
	fmt.Printf("ID:\t\t\t%d\n", resp.ID)
	fmt.Printf("CreatedAt:\t\t%s\n", toCompactTime(&resp.CreatedAt))
	fmt.Printf("EndedAt:\t\t%s\n", toCompactTime(resp.EndedAt))
	fmt.Printf("StartedAt:\t\t%s\n", toCompactTime(resp.StartedAt))
	fmt.Printf("UpdatedAt:\t\t%s\n", toCompactTime(&resp.UpdatedAt))
	fmt.Printf("TotalActual:\t\t%f\n", resp.TotalActualAmount)
	fmt.Printf("TotalBooked:\t\t%f\n", resp.TotalBookedAmount)
	fmt.Printf("TotalAdjustments:\t%f\n", resp.TotalAdjustments)

	return nil
}
