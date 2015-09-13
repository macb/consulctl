package main

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var (
	keyFlag = cli.StringFlag{
		Name:  "key",
		Usage: "desired key",
	}
	valueFlag = cli.StringFlag{
		Name:  "value",
		Usage: "desired value",
	}
)

func kvCommand() cli.Command {
	encodedFlag := cli.BoolFlag{
		Name:  "encoded",
		Usage: "returns the value encoded",
	}

	return cli.Command{
		Name:    "key-value",
		Aliases: []string{"kv"},
		Usage:   "key-value endpoint",
		Action: func(c *cli.Context) {
			cli.ShowSubcommandHelp(c)
			os.Exit(1)
		},
		Subcommands: []cli.Command{
			cli.Command{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "retrieve the value for a given key",
				Flags:   []cli.Flag{keyFlag, encodedFlag},
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					key := requiredKey(c)
					kv, _, err := cc.KV().Get(key, queryOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch key %s: %v", key, err)
					}
					if !c.Bool(encodedFlag.Name) {
						decoded := make([]byte, base64.URLEncoding.DecodedLen(len(kv.Value)))
						base64.URLEncoding.Decode(decoded, kv.Value)
						kv.Value = decoded
					}
					prettyPrint(kv)
				},
			},
			cli.Command{
				Name:    "put",
				Aliases: []string{"p"},
				Usage:   "put the value for a given key",
				Flags:   append([]cli.Flag{keyFlag, valueFlag}, writeOptionFlags()...),
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					key := requiredKey(c)
					value := requiredValue(c)
					pair := kvpair(key, []byte(value))
					_, err := cc.KV().Put(pair, writeOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch key %s: %v", key, err)
					}
					prettyPrint("OK")
				},
			},
		},
	}
}

func requiredKey(c *cli.Context) string {
	key := c.String(keyFlag.Name)
	if key == "" {
		log.Fatal("key is required")
	}

	return key
}

func requiredValue(c *cli.Context) string {
	value := c.String(valueFlag.Name)
	if value == "" {
		log.Fatal("value is required")
	}

	return value
}
