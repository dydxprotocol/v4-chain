#!/bin/bash

# Quickest way to reset the network/clear all Indexer data sources without rebuilding from scratch

echo "Stopping all Docker containers..."
docker stop $(docker ps -a -q)

echo "Deleting all dydxprotocold* containers..."
docker rm $(docker ps -a | grep dydxprotocold | awk '{print $1}')

echo "Resetting the protocol..."
cd ../protocol
make reset-chain

echo "Deleting the postgres container..."
docker rm $(docker ps -a | grep postgres | awk '{print $1}')

echo "Restarting the Kafka container..."
docker start $(docker ps -a | grep kafka | awk '{print $1}')

echo "Clearing all Kafka topics..."
KAFKA_CONTAINER=$(docker ps -a | grep kafka | awk '{print $1}')
docker cp remove-all-kafka-msgs.sh $KAFKA_CONTAINER:/opt/kafka
docker exec -it $KAFKA_CONTAINER /bin/bash -c "./remove-all-kafka-msgs.sh"

echo "Restarting all containers..."
docker-compose -f docker-compose-e2e-test.yml up
