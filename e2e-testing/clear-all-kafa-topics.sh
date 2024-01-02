#!/bin/bash

# Kafka Zookeeper address
ZOOKEEPER="localhost:2181"

# Topics to process
declare -a topics=("to-ender" "to-vulcan" "to-websockets-orderbooks" "to-websockets-subaccounts" "to-websockets-trades" "to-websockets-markets" "to-websockets-candles")

for topic in "${topics[@]}"
do
    echo "Deleting topic: $topic"
    kafka-topics.sh --zookeeper $ZOOKEEPER --delete --topic $topic || { echo "Failed to delete topic: $topic"; exit 1; }
    echo "Creating topic: $topic"
    kafka-topics.sh --create --zookeeper $ZOOKEEPER --replication-factor 1 --partitions 1 --topic $topic || { echo "Failed to create topic: $topic"; exit 1; }
done

echo "Topic processing completed."
