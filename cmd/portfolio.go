package main

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

var (
	portfolioCmd = &cli.Command{
		Name:    "portfolio",
		Aliases: []string{"p"},
		Usage:   "portfolio related commands",
		Subcommands: []*cli.Command{
			holdingsCmd,
			positionsCmd,
		},
	}

	holdingsCmd = &cli.Command{
		Name:    "holdings",
		Aliases: []string{"h"},
		Usage:   "show the holdings",
		Action: func(c *cli.Context) error {
			m, err := kc.GetHoldings()
			if err != nil {
				log.Fatalf("err: %v", err)
			}
			fmt.Println(m)
			return nil
		},
	}

	positionsCmd = &cli.Command{
		Name:    "positions",
		Aliases: []string{"p"},
		Usage:   "show the positions",
		Action: func(c *cli.Context) error {
			m, err := kc.GetPositions()
			if err != nil {
				log.Fatalf("err: %v", err)
			}
			fmt.Println(m)
			return nil
		},
	}
)
