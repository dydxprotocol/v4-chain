# E2E testing

E2E testing framework for Indexer/protocol.

### Instructions

Spin up a containerized environment running both the network and Indexer services:

```
cd indexer
docker compose -f docker-compose-e2e-test.yml up
```

```
cd indexer/services/e2e-testing
pnpm build && pnpm test
```
