# Secondary Filecoin Retrieval Markets

[![Build Status](https://travis-ci.com/ChainSafe/fil-secondary-retrieval-markets.svg?branch=main)](https://travis-ci.com/ChainSafe/fil-secondary-retrieval-markets)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/ChainSafe/fil-secondary-retrieval-markets)

❗**Current development should be considered a work in progress.**

The aim of this project is to enable secondary retrieval markets in [Filecoin](https://filecoin.io/). 

## Summary
Presently, the core Filecoin protocol directly enables those who are incentivized to store files (Storage Providers) to also be incentivized to provide those files to users (Retrieval Providers). The protocol provides extensibility such that a secondary market can be built on top of this. By enabling a secondary market we can reduce reliance on storage miners, improve network speed and reliability, and allow the network to scale further.

The key requirement for a secondary market is an improved discovery mechanism that allows anyone to participate as a Retrieval Provider. This software allows them to listen for and respond to requests for data. 

## Dependencies
`go >=1.14`

## Installation

To install the client, run:
```
git clone https://github.com/ChainSafe/fil-secondary-retrieval-markets && cd fil-secondary-retrieval-markets
make install
```

## Usage

The client requires a list of bootnodes and a payload CID to query:
```
retrieval-client --bootnodes <bootnodes> <CID>
```

Eg. 

```
retrieval-client --bootnodes "/dns4/some.network/tcp/1347/p2p/12D3KooWBEDQ5Xwh3JC67yxjNf91pZcpavrAwaqprNzbquC1yj6t,/dns4/some.network/tcp/1347/p2p/12D3KooWKbUF17McnN516w8TjmbkVNkcAZS9LnE5yJwH7pVDYPUJ" bafybeierhgbz4zp2x2u67urqrgfnrnlukciupzenpqpipiz5nwtq7uxpx4
```

## License

This repo is dual licensed under [MIT](/LICENSE-MIT) and [Apache 2.0](/LICENSE-APACHE).
