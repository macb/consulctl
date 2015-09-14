package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func catalogCommand() cli.Command {
	return cli.Command{
		Name:    "catalog",
		Aliases: []string{"c"},
		Usage:   "catalog related actions",
		Action: func(c *cli.Context) {
			cli.ShowSubcommandHelp(c)
			os.Exit(1)
		},
		Subcommands: []cli.Command{
			cli.Command{
				Name:    "datacenters",
				Aliases: []string{"d"},
				Usage:   "list of known datacenters",
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					dcs, err := cc.Catalog().Datacenters()
					if err != nil {
						log.Fatalf("failed to fetch datacenters: %v", err)
					}
					prettyPrint(dcs)
				},
			},
			cli.Command{
				Name:    "nodes",
				Aliases: []string{"ns"},
				Usage:   "list of nodes",
				Flags:   append([]cli.Flag{}, queryOptionFlags()...),
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					nodes, _, err := cc.Catalog().Nodes(queryOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch nodes: %v", err)
					}
					prettyPrint(nodes)
				},
			},
			cli.Command{
				Name:    "node",
				Aliases: []string{"n"},
				Usage:   "retrieve information for a given node",
				Flags:   append([]cli.Flag{nodeFlag}, queryOptionFlags()...),
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					name := nodeName(c)
					node, _, err := cc.Catalog().Node(name, queryOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch nodes: %v", err)
					}
					prettyPrint(node)
				},
			},
			cli.Command{
				Name:    "service",
				Aliases: []string{"s"},
				Usage:   "retrieve information for a service",
				Flags:   append([]cli.Flag{serviceFlag, tagFlag}, queryOptionFlags()...),
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					name := serviceName(c)
					tag := c.String(tagFlag.Name)
					service, _, err := cc.Catalog().Service(name, tag, queryOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch service information: %v", err)
					}
					prettyPrint(service)
				},
			},
			cli.Command{
				Name:    "services",
				Aliases: []string{"ss"},
				Usage:   "list of known services",
				Flags:   append([]cli.Flag{}, queryOptionFlags()...),
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					services, _, err := cc.Catalog().Services(queryOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch services: %v", err)
					}
					prettyPrint(services)
				},
			},
		},
	}
}
