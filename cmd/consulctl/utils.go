package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/armon/consul-api"
	"github.com/codegangsta/cli"
)

var (
	dcFlag = cli.StringFlag{
		Name:  "dc",
		Usage: "consul datacenter",
	}
	consistentFlag = cli.BoolFlag{
		Name:  "consistent",
		Usage: "require consistent consistency",
	}
	staleFlag = cli.BoolFlag{
		Name:  "stale",
		Usage: "allow stale consistency",
	}
	waitIndexFlag = cli.IntFlag{
		Name:  "wait.index",
		Usage: "specify wait index",
	}
	waitFlag = cli.DurationFlag{
		Name:  "wait",
		Usage: "specify time to wait",
	}
)

func prettyPrint(data interface{}) {
	if data == nil {
		log.Fatal("no results returned.")
	}
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("failed to pretty print %v: %v", data, err)
	}
	fmt.Println(string(out))
}

func consulClient(c *cli.Context) *consulapi.Client {
	conf := consulapi.DefaultConfig()
	conf.Address = c.GlobalString(apiFlag.Name)
	client, err := consulapi.NewClient(conf)
	if err != nil {
		log.Fatalf("failed to build client: %v", err)
	}
	return client
}

func queryOptionFlags() []cli.Flag {
	return []cli.Flag{
		dcFlag,
		consistentFlag,
		staleFlag,
		waitFlag,
		waitIndexFlag,
	}
}

func queryOptions(c *cli.Context) *consulapi.QueryOptions {
	consistent := c.Bool(consistentFlag.Name)
	stale := c.Bool(staleFlag.Name)
	if consistent && stale {
		cli.ShowSubcommandHelp(c)
		log.Fatalf("only --stale or --consistent may be set, not both")
	}
	return &consulapi.QueryOptions{
		Datacenter:        c.String(dcFlag.Name),
		AllowStale:        stale,
		RequireConsistent: consistent,
		WaitTime:          c.Duration(waitFlag.Name),
		WaitIndex:         uint64(c.Int(waitIndexFlag.Name)),
	}
}

func nodeName(c *cli.Context) string {
	name := c.String(nodeFlag.Name)
	if name == "" {
		cli.ShowSubcommandHelp(c)
		log.Fatal("node is required")
	}
	return name
}

func serviceName(c *cli.Context) string {
	name := c.String(serviceFlag.Name)
	if name == "" {
		cli.ShowSubcommandHelp(c)
		log.Fatal("service is required")
	}
	return name
}
