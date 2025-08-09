#!/bin/sh
set -e
dnsmasq
echo "nameserver 127.0.0.1" > /etc/resolv.conf
chronyd
exec su dydx -s /bin/sh -c 'exec "$@"' -- "$@"