# Overview

This directory contains structs and functions for sending data from the V4 application to the
[Indexer](https://github.com/dydxprotocol/indexer).

The base package contains definitions of Indexer specific command-line flags which are added to the
`start` command of the V4 application.

## events

The `events` package contains definitions of on-chain event structs the V4 application emits to the 
Indexer along with helper functions to instantiate instances of the events.

## msgsender

The `msgsender` package contains structs used to send both off-chain and on-chain data to the
Indexer. Currently, all data is sent to the Indexer via Kafka.

## off_chain_updates

The `off_chain_updates` package contains definitions of off-chain update structs the V4 application
emits to the Indexer.
