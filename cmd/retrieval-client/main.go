// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/client"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli"
)

var (
	log = logging.Logger("main")

	bootnodesFlag = cli.StringFlag{
		Name:  "bootnodes",
		Usage: "comma-separated list of peer addresses",
		Required: true,
	}

	pieceCIDFlag = cli.StringFlag{
		Name: "pieceCID",
		Usage: "specifies a piece CID to query",
	}

	timeoutFlag = cli.Float64Flag{
		Name:  "timeout",
		Usage: "Specify how long to listen for requests (seconds)",
		Value: defaultResponseTimeout,
	}

	flags = []cli.Flag{
		bootnodesFlag,
		pieceCIDFlag,
		timeoutFlag,
	}

	app = cli.NewApp()

	defaultResponseTimeout = time.Minute.Seconds()
)

func init() {
	app.Action = run
	app.Name = "retrieval-client"
	app.Flags = flags
	app.Usage = "Client for secondary retrieval markets"
	app.UsageText = "retrieval-client [options] <CID>"
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
	err = logging.SetLogLevel("main", "debug")
	if err != nil {
		return err
	}

	cidStr := ctx.Args().First()
	pieceCIDStr := ctx.String(pieceCIDFlag.Name)
	bootnodesStr := ctx.String(bootnodesFlag.Name)
	timeout := ctx.Float64(timeoutFlag.Name)

	n, err := newNetwork(bootnodesStr)
	if err != nil {
		return err
	}

	c := client.NewClient(n)
	var payloadCID, pieceCID cid.Cid
	payloadCID, err = cid.Decode(cidStr)
	if err != nil {
		return err
	}

	if pieceCIDStr != "" {
		pieceCID, err = cid.Decode(pieceCIDStr)
		if err != nil {
			return err
		}
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

	// TODO: update cli to allow specifying Selector
	params := shared.Params{
		PayloadCID: payloadCID,
		PieceCID: &pieceCID,
	}

	if pieceCIDStr != "" {
		log.Infof("Querying for payload %s and piece %s", payloadCID, pieceCIDStr)
	} else {
		log.Infof("Querying for payload %s", payloadCID)
	}

	h := newResponseHandler()
	unsubscribe := c.SubscribeToQueryResponses(h.handleResponse, params)
	defer unsubscribe()

	err = c.SubmitQuery(context.Background(), params)
	if err != nil {
		return err
	}
	for {
		select {
		case resp := <-h.respCh:
			log.Info("got response! ", resp)
		case <-time.After(time.Duration(timeout)):
			log.Info("no responses received by timeout")
			return nil
		}
	}
}
