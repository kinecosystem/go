---
title: Overview
---

The Go SDK contains packages for interacting with most aspects of the Kin ecosystem.  In addition to generally useful, low-level packages such as [`keypair`](https://godoc.org/github.com/kinecosystem/go/keypair) (used for creating Kin-compliant public/secret key pairs), the Go SDK also contains code for the server applications and client tools written in go.

## Godoc reference

The most accurate and up-to-date reference information on the Go SDK is found within godoc.  The godoc.org service automatically updates the documentation for the Go SDK everytime github is updated.  The godoc for all of our packages can be found at (https://godoc.org/github.com/kinecosystem/go).

## Client Packages

The Go SDK contains packages for interacting with the various Kin services:

- [`horizon`](https://godoc.org/github.com/kinecosystem/go/clients/horizon) provides client access to a horizon server, allowing you to load account information, stream payments, post transactions and more.
- [`stellartoml`](https://godoc.org/github.com/kinecosystem/go/clients/stellartoml) provides the ability to resolve Stellar.toml files from the internet.  You can read about [Stellar.toml concepts here](../../guides/concepts/stellar-toml.md).

