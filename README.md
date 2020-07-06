# Secondary Filecoin Retrieval Markets

❗**Current development should be considered a work in progress.**

The aim of this project is to enable secondary retrieval markets in [Filecoin](https://filecoin.io/). 

Presently, the core Filecoin protocol directly enables those who are incentivized to store files (Storage Providers) to also be incentivized to provide those files to users (Retrieval Providers). The protocol provides extensibility such that a secondary market can be built on top of this. By enabling a secondary market we can reduce reliance on miners, improve network speed and reliability, and allow the network to scale further.

The main feature of the secondary market is an improved discovery mechanism that allows anyone to participate as a Retrieval Provider. This software allows them to listen for and respond to requests for data. 

## Dependencies

Install go version `>=1.14`

## Installation

```
go get -u github.com/ChainSafe/fil-secondary-retrieval-markets
```

## License

This repo is dual licensed under [MIT](/LICENSE-MIT) and [Apache 2.0](/LICENSE-APACHE).
