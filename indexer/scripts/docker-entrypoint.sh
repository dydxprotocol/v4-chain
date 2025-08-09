#!/bin/sh
dnsmasq &
echo "nameserver 127.0.0.1" > /etc/resolv.conf
chronyd &
exec "$@"