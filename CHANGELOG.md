# Changelog

All notable changes to this project will be documented in this
file.  This project adheres to [Semantic Versioning](http://semver.org/).

As this project is pre 1.0, breaking changes may happen for minor version
bumps.  A breaking change will get clearly notified in this log.

NOTE:  this changelog represents the changes that are associated with the library code in this repo (rather than the tools or services in this repo).  

## [Unreleased]

### Added

- xdr: added support for new signer types
- build: `Signer` learned support for new signer types
- strkey: added support for new signer types
- network:  Added the `HashTransaction` helper func to get the hash of a transaction targeted to a specific stellar network.
- trades: Added Server-Sent Events endpoint to support streaming of trades
- trades: add `base_offer_id` and `counter_offer_id` to trade resources.
- trade aggregation: Added an optional `offset` parameter that lets you offset the bucket timestamps in hour-long increments. Can only be used if the `resolution` parameter is greater than 1 hour. `offset` must also be in whole-hours and less than 24 hours.


### Changed:

- build: _BREAKING CHANGE_:  A transaction built and signed using the `build` package no longer default to the test network.
- trades for offer endpoint will query for trades that match the given offer on either side of trades, rather than just the "sell" offer.

[Unreleased]: https://github.com/stellar/go/commits/master