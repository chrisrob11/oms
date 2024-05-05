package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/chrisrob11/oms/internal/cmds"
	"github.com/urfave/cli/v2"
)

var ErrMissingParameter = errors.New("missing parameter")

func main() {
	app := &cli.App{
		Name:  "omsclient",
		Usage: "Interact with the oms service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
				Value:   "http://localhost:8080",
				Usage:   "URL of the oms service",
				EnvVars: []string{"OMS_URL"},
			},
		},
		Commands: []*cli.Command{
			cmds.Import,
			cmds.CreateCampaign,
			cmds.ListCampaign,
			cmds.ShowCampaign,
			cmds.CreateCampaignLineItem,
			cmds.ListCampaignItemLine,
			cmds.ShowCampaignItemLine,
			cmds.GenerateInvoice,
			cmds.ListInvoices,
			cmds.ShowInvoice,
			cmds.AdjustInvoice,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func NewMissingError(parameterName string) error {
	return fmt.Errorf("%w: %s", ErrMissingParameter, parameterName)
}
