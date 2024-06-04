#!/bin/bash

# Ensure the script exits on any command failure
set -e
LISTEN_ADDRESS=$1
sudo ip netns exec consu-query yum install -y dnsmasq
CONFIG_CONTENT="interface=eth0
listen-address=${LISTEN_ADDRESS}
bind-interfaces
server=127.0.0.11"

# Write the content to /etc/dnsmasq.conf
echo "$CONFIG_CONTENT" | sudo tee /etc/dnsmasq.conf > /dev/null

# Restart dnsmasq to apply the changes
sudo pkill dnsmasq
sudo dnsmasq --conf-file=/etc/dnsmasq.conf

echo "dnsmasq configuration updated and service restarted."

ip netns exec consu-query sh -c "echo 'nameserver ${LISTEN_ADDRESS}' > /etc/resolv.conf"
