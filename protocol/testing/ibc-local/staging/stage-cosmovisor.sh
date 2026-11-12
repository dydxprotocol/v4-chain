#!/bin/bash
# Stages the testnet-validator image's dydxprotocold under cosmovisor upgrades/v9.9 on a
# testnet host, next to the running node container, without restarting it. Run over
# SSM RunShellScript (as a command list) on every validator and full-node host.
set -e
IMG=013339450148.dkr.ecr.us-east-2.amazonaws.com/testnet-validator:f9d4939
c=$(docker ps --format '{{.ID}} {{.Image}}' | grep 'testnet-validator' | head -1 | cut -d' ' -f1)
[ -n "$c" ] || { echo "no node container"; exit 1; }
home=$(docker inspect "$c" --format '{{range .Mounts}}{{.Destination}} {{end}}' | tr ' ' '\n' | grep '^/dydxprotocol/chain/' | head -1)
vol=$(docker inspect "$c" --format "{{range .Mounts}}{{if eq .Destination \"$home\"}}{{.Source}}{{end}}{{end}}")
echo "container=$c home=$home vol=$vol current=$(readlink $vol/cosmovisor/current)"
command -v aws >/dev/null || { echo "installing awscli"; yum install -y -q awscli >/dev/null 2>&1 || amazon-linux-extras install -y -q awscli >/dev/null 2>&1; }
aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin 013339450148.dkr.ecr.us-east-2.amazonaws.com >/dev/null
docker pull -q "$IMG" >/dev/null
tmp=$(docker create "$IMG")
docker cp "$tmp:/bin/dydxprotocold" /tmp/dydxprotocold-v9.9
docker rm "$tmp" >/dev/null
mkdir -p "$vol/cosmovisor/upgrades/v9.9/bin"
install -m 755 /tmp/dydxprotocold-v9.9 "$vol/cosmovisor/upgrades/v9.9/bin/dydxprotocold"
rm -f /tmp/dydxprotocold-v9.9
echo "staged: $(ls -la $vol/cosmovisor/upgrades/v9.9/bin/dydxprotocold)"
echo "sha256: $(sha256sum $vol/cosmovisor/upgrades/v9.9/bin/dydxprotocold | cut -c1-16)"
echo "runs in container: $(docker exec "$c" "$home/cosmovisor/upgrades/v9.9/bin/dydxprotocold" version 2>&1 | tail -1)"
echo "current still: $(docker exec "$c" readlink "$home/cosmovisor/current")"
echo "running: $(docker exec "$c" "$home/cosmovisor/current/bin/dydxprotocold" version 2>&1 | tail -1)"
