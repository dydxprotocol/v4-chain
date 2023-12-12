# E2E testing

E2E testing framework for Indexer/protocol.

### Instructions

Spin up a containerized environment running both the network and Indexer services:

In one terminal, run

```
./run-containerized-env.sh
```

In another terminal, run
```
pnpm build && pnpm test
```
