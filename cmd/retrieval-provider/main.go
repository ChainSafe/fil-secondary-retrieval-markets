// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/cache"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/cmd/utils"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/provider"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli"
)

var (
	log = logging.Logger("provider")

	dataFlag = cli.StringFlag{
		Name:  "data",
		Usage: "JSON file of available CIDs",
	}
	bootnodesFlag = cli.StringFlag{
		Name:  "bootnodes",
		Usage: "comma-separated list of peer addresses",
	}

	flags = []cli.Flag{
		dataFlag,
		bootnodesFlag,
	}

	app = cli.NewApp()
)

func init() {
	app.Action = run
	app.Flags = flags
}

func main() {
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	err := logging.SetLogLevel("provider", "debug")
	if err != nil {
		return err
	}

	dataStr := ctx.String(dataFlag.Name)
	bootnodesStr := ctx.String(bootnodesFlag.Name)

	psJSON := new(ProviderStoreJSON)

	if dataStr != "" {
		data, err := ioutil.ReadFile(dataStr)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &psJSON.cids)
		if err != nil {
			return err
		}
	}

	ps := psJSON.ToProviderStore()

	net, err := utils.NewNetwork(bootnodesStr)
	if err != nil {
		return err
	}

	log.Info("provider has ", ps.cids)

	p := provider.NewProvider(net, ps, cache.NewLFUCache(1024))
	err = p.Start()
	if err != nil {
		return err
	}

	log.Info("provider listening at ", net.MultiAddrs())
	select {}
}
