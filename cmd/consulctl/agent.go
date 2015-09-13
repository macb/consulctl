package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func agentCommand() cli.Command {
	enableFlag := cli.BoolFlag{
		Name:  "enable",
		Usage: "enable maintenance mode",
	}
	disableFlag := cli.BoolFlag{
		Name:  "disable",
		Usage: "disable maintenance mode",
	}
	reasonFlag := cli.StringFlag{
		Name:  "reason",
		Usage: "reason for maintenance mode",
	}

	return cli.Command{
		Name:    "agent",
		Aliases: []string{"a"},
		Usage:   "Agent endpoint",
		Action: func(c *cli.Context) {
			cli.ShowSubcommandHelp(c)
			os.Exit(1)
		},
		Subcommands: []cli.Command{
			cli.Command{
				Name:    "checks",
				Aliases: []string{"c"},
				Usage:   "Agent checks",
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					checks, err := cc.Agent().Checks()
					if err != nil {
						log.Fatalf("failed to fetch agent checks: %v", err)
					}
					prettyPrint(checks)
				},
			},
			cli.Command{
				Name:    "services",
				Aliases: []string{"srv"},
				Usage:   "Agent services",
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					services, err := cc.Agent().Services()
					if err != nil {
						log.Fatalf("failed to fetch agent services: %v", err)
					}
					prettyPrint(services)
				},
			},
			cli.Command{
				Name:    "members",
				Aliases: []string{"m"},
				Usage:   "Agent members",
				Flags:   []cli.Flag{wanFlag},
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					wan := c.Bool(wanFlag.Name)
					members, err := cc.Agent().Members(wan)
					if err != nil {
						log.Fatalf("failed to fetch members: %v", err)
					}
					prettyPrint(members)
				},
			},
			cli.Command{
				Name:    "self",
				Aliases: []string{"s"},
				Usage:   "agent specific information",
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					self, err := cc.Agent().Self()
					if err != nil {
						log.Fatalf("failed to fetch self: %v", err)
					}
					prettyPrint(self)
				},
			},
			cli.Command{
				Name:    "maintenance",
				Aliases: []string{"maint"},
				Usage:   "agent maintenance",
				Action: func(c *cli.Context) {
					cli.ShowSubcommandHelp(c)
					os.Exit(1)
				},
				Subcommands: []cli.Command{
					cli.Command{
						Name:    "node",
						Aliases: []string{"n"},
						Usage:   "node maintenance",
						Flags:   []cli.Flag{enableFlag, disableFlag, reasonFlag},
						Action: func(c *cli.Context) {
							cc := consulClient(c)
							enable := c.Bool(enableFlag.Name)
							disable := c.Bool(disableFlag.Name)
							var err error
							switch {
							case enable && disable:
								cli.ShowSubcommandHelp(c)
								log.Fatal("enable or disable must be specified, not both")
							case enable:
								reason := c.String(reasonFlag.Name)
								if reason == "" {
									cli.ShowSubcommandHelp(c)
									log.Fatal("reason must be specified when enabling maintenance mode")
								}
								err = cc.Agent().EnableNodeMaintenance(reason)
							case disable:
								err = cc.Agent().DisableNodeMaintenance()
							default:
								cli.ShowSubcommandHelp(c)
								log.Fatal("enable or disable must be specified")
							}
							if err != nil {
								log.Fatalf("failed to enable/disable node maintence: %v", err)
							}
							prettyPrint("OK")
						},
					},
					cli.Command{
						Name:    "service",
						Aliases: []string{"s"},
						Usage:   "service maintenance",
						Flags:   []cli.Flag{enableFlag, disableFlag, reasonFlag, serviceFlag},
						Action: func(c *cli.Context) {
							cc := consulClient(c)
							enable := c.Bool(enableFlag.Name)
							disable := c.Bool(disableFlag.Name)
							name := serviceName(c)

							var err error
							switch {
							case enable && disable:
								cli.ShowSubcommandHelp(c)
								log.Fatal("enable or disable must be specified, not both")
							case enable:
								reason := c.String(reasonFlag.Name)
								if reason == "" {
									cli.ShowSubcommandHelp(c)
									log.Fatal("reason must be specified when enabling maintenance mode")
								}
								err = cc.Agent().EnableServiceMaintenance(name, reason)
							case disable:
								err = cc.Agent().DisableServiceMaintenance(name)
							default:
								cli.ShowSubcommandHelp(c)
								log.Fatal("enable or disable must be specified")
							}

							if err != nil {
								log.Fatalf("failed to enable/disable service maintence: %v", err)
							}
							prettyPrint("OK")
						},
					},
				},
			},
		},
	}
}
