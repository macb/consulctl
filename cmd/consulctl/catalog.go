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
		Usage:   "Catalog endpoint",
		Action: func(c *cli.Context) {
			cli.ShowSubcommandHelp(c)
			os.Exit(1)
		},
		Subcommands: []cli.Command{
			cli.Command{
				Name:    "datacaenters",
				Aliases: []string{"d"},
				Usage:   "Datacenter list",
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
				Usage:   "List of nodes",
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
				Usage:   "Information for a single node",
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
				Usage:   "Information for a single service",
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
				Usage:   "List of services",
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
			cli.Command{
				Name:    "register",
				Aliases: []string{"r"},
				Usage:   "Register a service",
				//Flags:   append([]cli.Flag{}, queryOptionFlags()...),
				Action: func(c *cli.Context) {
					log.Fatal("TODO: implement")
					//cc := consulClient(c)

					//if err != nil {
					//log.Fatalf("failed to fetch services: %v", err)
					//}
					//prettyPrint(services)
				},
			},
			cli.Command{
				Name:    "deregister",
				Aliases: []string{"dr"},
				Usage:   "Deregister a service",
				//Flags:   append([]cli.Flag{}, queryOptionFlags()...),
				Action: func(c *cli.Context) {
					log.Fatal("TODO: implement")
					//cc := consulClient(c)

					//if err != nil {
					//log.Fatalf("failed to fetch services: %v", err)
					//}
					//prettyPrint(services)
				},
			},
		},
	}
}
