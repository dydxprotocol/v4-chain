#!/bin/bash

# Kafka Zookeeper address
ZOOKEEPER="localhost:2181"

# Topics to process
declare -a topics=("to-ender" "to-vulcan" "to-websockets-orderbooks" "to-websockets-subaccounts" "to-websockets-trades" "to-websockets-markets" "to-websockets-candles")

for topic in "${topics[@]}"
do
    echo "Deleting topic: $topic"
    kafka-topics.sh --zookeeper $ZOOKEEPER --delete --topic $topic
    echo "Creating topic: $topic"
    kafka-topics.sh --create --zookeeper $ZOOKEEPER --replication-factor 1 --partitions 1 --topic $topic
done

echo "Topic processing completed."
