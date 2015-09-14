package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var (
	stateUnknownFlag = cli.BoolFlag{
		Name:  "unknown",
		Usage: "list services in the 'unknown' state",
	}
	statePassingFlag = cli.BoolFlag{
		Name:  "passing",
		Usage: "list services in the 'passing' state",
	}
	stateWarningFlag = cli.BoolFlag{
		Name:  "warning",
		Usage: "list services in the 'warning' state",
	}
	stateCriticalFlag = cli.BoolFlag{
		Name:  "critical",
		Usage: "list services in the 'critical' state",
	}

	stateFlags = []cli.Flag{
		stateUnknownFlag,
		statePassingFlag,
		stateWarningFlag,
		stateCriticalFlag,
	}
)

func healthCommand() cli.Command {
	return cli.Command{
		Name:    "health",
		Aliases: []string{"h"},
		Usage:   "health related actions",
		Action: func(c *cli.Context) {
			cli.ShowSubcommandHelp(c)
			os.Exit(1)
		},
		Subcommands: []cli.Command{
			cli.Command{
				Name:    "checks",
				Aliases: []string{"c"},
				Usage:   "list registered health checks for a service",
				Flags:   append([]cli.Flag{serviceFlag}, queryOptionFlags()...),
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					name := serviceName(c)
					checks, _, err := cc.Health().Checks(name, queryOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch health checks: %v", err)
					}
					prettyPrint(checks)
				},
			},
			cli.Command{
				Name:    "node",
				Aliases: []string{"n"},
				Usage:   "retrieve health of a given node",
				Flags:   append([]cli.Flag{nodeFlag}, queryOptionFlags()...),
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					name := nodeName(c)
					checks, _, err := cc.Health().Node(name, queryOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch node health: %v", err)
					}
					prettyPrint(checks)
				},
			},
			cli.Command{
				Name:    "service",
				Aliases: []string{"s"},
				Usage:   "retrieve health of a given service",
				Flags: append([]cli.Flag{
					serviceFlag,
					cli.BoolFlag{
						Name:  "passing",
						Usage: "only show passing services",
					},
				}, queryOptionFlags()...),
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					name := serviceName(c)
					passing := c.Bool("passing")
					tag := c.String(tagFlag.Name)
					checks, _, err := cc.Health().Service(name, tag, passing, queryOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch service health: %v", err)
					}
					prettyPrint(checks)
				},
			},
			cli.Command{
				Name:    "state",
				Aliases: []string{"st"},
				Usage:   "list services based on check state, 'any' by default.",
				Flags:   append(stateFlags, queryOptionFlags()...),
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					state := serviceState(c)
					services, _, err := cc.Health().State(state, queryOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch services: %v", err)
					}
					prettyPrint(services)
				},
			},
		},
	}
}

func serviceState(c *cli.Context) string {
	states := make([]string, 0)

	for _, state := range stateFlags {
		bs := state.(cli.BoolFlag)
		if c.Bool(bs.Name) {
			states = append(states, bs.Name)
		}
	}

	switch len(states) {
	case 0:
		return "any"
	case 1:
		return states[0]
	default:
		cli.ShowSubcommandHelp(c)
		log.Fatal("specifying more than 1 state is not supported")
		return ""
	}
}
