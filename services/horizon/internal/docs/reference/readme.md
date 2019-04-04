---
title: Overview
---

Horizon is an API server for the Kin ecosystem.  It acts as the interface between [core](https://github.com/kinecosystem/stellar-core) and applications that want to access the Kin network. It allows you to submit transactions to the network, check the status of accounts, subscribe to event streams, etc. You can also watch a [talk on Horizon](https://www.youtube.com/watch?v=AtJ-f6Ih4A4) by Stellar.org developer Scott Fleckenstein:

[![Horizon: API webserver for the Stellar network](https://img.youtube.com/vi/AtJ-f6Ih4A4/sddefault.jpg "Horizon: API webserver for the Stellar network")](https://www.youtube.com/watch?v=AtJ-f6Ih4A4)

Horizon provides a RESTful API to allow client applications to interact with the Stellar network. You can communicate with Horizon using cURL or just your web browser. However, if you're building a client application, you'll likely want to use a Stellar SDK in the language of your client.
Kin Foundation provides [SDKs](https://kin.org/developers) for clients to use to interact with Horizon.

Kin Foundation runs an instance of Horizon that is connected to the test net: [https://horizon-testnet.kinfederation.com/](https://horizon-testnet.kinfederation.com/) and one that is connected to the public Stellar network:
[https://horizon.kinfederation.com/](https://horizon.kinfederation.com/).

