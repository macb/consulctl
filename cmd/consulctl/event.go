package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func eventCommand() cli.Command {
	eventFlag := cli.StringFlag{
		Name:  "event.name",
		Usage: "event name",
	}
	return cli.Command{
		Name:    "event",
		Aliases: []string{"e"},
		Usage:   "event endpoint",
		Action: func(c *cli.Context) {
			cli.ShowSubcommandHelp(c)
			os.Exit(1)
		},
		Subcommands: []cli.Command{
			cli.Command{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list latest events an agent has seen",
				Flags:   append([]cli.Flag{eventFlag}, queryOptionFlags()...),
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					name := c.String(eventFlag.Name)
					events, _, err := cc.Event().List(name, queryOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch agent events: %v", err)
					}
					prettyPrint(events)
				},
			},
		},
	}
}
