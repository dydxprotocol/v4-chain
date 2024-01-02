#!/bin/bash

# Quickest way to reset the network/clear all Indexer data sources without rebuilding from scratch

echo "Stopping all Docker containers..."
docker stop $(docker ps -a | grep 'e2e-testing' | awk '{print $1}')

echo "Deleting all dydxprotocold* containers..."
docker rm $(docker ps -a | grep dydxprotocold | awk '{print $1}')

echo "Resetting the protocol..."
cd ../protocol
make reset-chain

echo "Deleting the postgres container..."
docker rm $(docker ps -a | grep postgres | awk '{print $1}')

echo "Restarting the Kafka container..."
KAFKA_CONTAINER=$(docker ps -a | grep 'e2e-testing' | grep 'kafka' | awk '{print $1}')
docker start $KAFKA_CONTAINER

echo "Clearing all Kafka topics..."
cd ../e2e-testing
docker cp clear-all-kafa-topics.sh $KAFKA_CONTAINER:/opt/kafka
docker exec -it $KAFKA_CONTAINER /bin/bash -c "./clear-all-kafa-topics.sh"

echo "Restarting all containers..."
docker-compose -f docker-compose-e2e-test.yml up -d
