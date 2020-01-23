package main

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

var (
	userCmd = &cli.Command{
		Name:    "user",
		Aliases: []string{"u"},
		Usage:   "user related commands",
		Subcommands: []*cli.Command{
			marginCmd,
			profileCmd,
		},
	}

	marginCmd = &cli.Command{
		Name:    "margin",
		Aliases: []string{"m"},
		Usage:   "show the margins",
		Action: func(c *cli.Context) error {
			m, err := kc.GetUserMargins()
			if err != nil {
				log.Fatalf("err: %v", err)
			}
			fmt.Println(m)
			return nil
		},
	}

	profileCmd = &cli.Command{
		Name:    "profile",
		Aliases: []string{"p"},
		Usage:   "show the profile",
		Action: func(c *cli.Context) error {
			m, err := kc.GetUserProfile()
			if err != nil {
				log.Fatalf("err: %v", err)
			}
			fmt.Println(m)
			return nil
		},
	}
)
