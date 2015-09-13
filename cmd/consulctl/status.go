package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func statusCommand() cli.Command {
	return cli.Command{
		Name:    "status",
		Aliases: []string{"st"},
		Usage:   "Status endpoint",
		Subcommands: []cli.Command{
			cli.Command{
				Name:    "peer",
				Aliases: []string{"p"},
				Usage:   "Raft peer set",
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					peers, err := cc.Status().Peers()
					if err != nil {
						log.Fatalf("failed to fetch peers: %v", err)
					}
					prettyPrint(peers)
				},
			},
			cli.Command{
				Name:    "leader",
				Aliases: []string{"l"},
				Usage:   "Raft leader",
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					leader, err := cc.Status().Leader()
					if err != nil {
						log.Fatalf("failed to fetch leader: %v", err)
					}
					prettyPrint(leader)
				},
			},
		},
		Action: func(c *cli.Context) {
			cli.ShowSubcommandHelp(c)
			os.Exit(1)
		},
	}
}
