# E2E testing

E2E testing framework for Indexer/protocol. Read more about potential
issues in the E2E testing notion.

### Instructions

Spin up a containerized environment running both the network and Indexer services:

In one terminal, run

```
./run-containerized-env.sh
```

Wait until the chain that was launched has reached height at least 50. 
Then, in another terminal, from the e2e-testing directory, run:
```
pnpm build && pnpm test
```

#### Quickest way to reset the network/clear all Indexer data sources without rebuilding from scratch

```
./reset-network.sh
```