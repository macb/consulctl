package main

import (
	"log"

	"github.com/codegangsta/cli"
)

var (
	apiFlag = cli.StringFlag{
		Name:   "api.addr",
		EnvVar: "CONSUL_API_ADDR",
		Value:  "localhost:8500",
		Usage:  "address for the Consul API to access.",
	}

	serviceFlag = cli.StringFlag{
		Name:  "service",
		Usage: "specify a service (required)",
	}
	nodeFlag = cli.StringFlag{
		Name:  "node",
		Usage: "specify a node (required)",
	}
	tagFlag = cli.StringFlag{
		Name:  "tag",
		Usage: "specify a service tag",
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "consulctl"
	app.Author = "Mac Browning"
	app.Email = "mac@macasaurus.com"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		apiFlag,
	}

	app.Commands = []cli.Command{
		catalogCommand(),
		statusCommand(),
		healthCommand(),
	}

	log.SetPrefix("consulctl: ")
	log.SetFlags(0)

	app.RunAndExitOnError()
}
