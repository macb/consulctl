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
	modifyIndexFlag := cli.IntFlag{
		Name:  "modify-index",
		Usage: "modify index for the KV pair",
	}

	return cli.Command{
		Name:    "key-value",
		Aliases: []string{"kv"},
		Usage:   "key-value related actions",
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
						_, err := base64.URLEncoding.Decode(decoded, kv.Value)
						if err != nil {
							log.Printf("failed to base64 decode: %v", err)
						} else {
							kv.Value = decoded
						}
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
			cli.Command{
				Name:    "compare-and-set",
				Aliases: []string{"cas"},
				Usage:   "set the value for a given key if it doesn't exist or if the modify index matches",
				Flags:   append([]cli.Flag{keyFlag, valueFlag, modifyIndexFlag}, writeOptionFlags()...),
				Action: func(c *cli.Context) {
					cc := consulClient(c)
					key := requiredKey(c)
					value := requiredValue(c)
					pair := kvpair(key, []byte(value))
					pair.ModifyIndex = uint64(c.Int(modifyIndexFlag.Name))
					set, _, err := cc.KV().CAS(pair, writeOptions(c))
					if err != nil {
						log.Fatalf("failed to fetch key %s: %v", key, err)
					}
					resp := "SET"
					if !set {
						resp = "NOT SET"
					}
					prettyPrint(resp)
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
