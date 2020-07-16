// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/client"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli"
)

var (
	log = logging.Logger("cli")

	queryFlag = cli.StringFlag{
		Name:  "query",
		Usage: "submit query for a CID",
	}
	bootnodesFlag = cli.StringFlag{
		Name:  "bootnodes",
		Usage: "comma-separated list of peer addresses",
	}

	flags = []cli.Flag{
		queryFlag,
		bootnodesFlag,
	}

	app = cli.NewApp()

	responseTimeout = time.Minute
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
	err := logging.SetLogLevel("client", "debug")
	if err != nil {
		return err
	}

	cidStr := ctx.String(queryFlag.Name)
	bootnodesStr := ctx.String(bootnodesFlag.Name)

	n, err := newNetwork(bootnodesStr)
	if err != nil {
		return err
	}

	c := client.NewClient(n)
	cid, err := cid.Decode(cidStr)
	if err != nil {
		return err
	}

	err = c.Start()
	if err != nil {
		return err
	}

	defer func() {
		err = c.Stop()
		if err != nil {
			log.Error("failed to stop client", err)
		}
	}()

	h := newResponseHandler()
	unsubscribe := c.SubscribeToQueryResponses(h.handleResponse, cid)
	defer unsubscribe()

	err = c.SubmitQuery(context.Background(), cid)
	if err != nil {
		return err
	}

	log.Info("submit query ", cid)

	for {
		select {
		case resp := <-h.respCh:
			log.Info("got response! ", resp)
		case <-time.After(responseTimeout):
			log.Info("no responses received by timeout")
			return nil
		}
	}
}
