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

#### Quickest way to reset the network without rebuilding from scratch

Go to Docker Desktop

Stop all containers

Delete all dydxprotocold* containers

Reset the protocol by doing the following:
```
cd ../protocol
make reset-chain
```

Delete the postgres container.

Restart the Kafka container.

Clear all Kafka topics:
```
docker cp remove-all-kafka-msgs.sh <container_id>:/opt/kafka
docker exec -it <container_id> /bin/bash
./remove-all-kafka-msgs.sh
```

Restart all containers:
```
docker compose -f docker-compose-e2e-test.yml up
```
