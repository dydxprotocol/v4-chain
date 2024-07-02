#!/bin/bash

# This script takes in 3 arguments: service, commit, and environment.
# It pushes an ECR image with tag local_<commit> and deploys it with Orb to our $env environment.

# To run the script, use the following command:
#   scripts/deploy-commit-to-env.sh <service> <commit> <dev|dev2|dev3|dev4|dev5|staging>
# For example:
#   scripts/deploy-commit-to-env.sh comlink 875aecd staging

service=$1
commit=$2
env=$3

tag=local_$commit

# set the account number depending on the environment
account=329916310755
case $env in
    "dev") account=329916310755;;
    "dev2") account=220669270420;;
    "dev3") account=295746761472;;
    "dev4") account=525975847385;;
    "dev5") account=917958511744;;
    "staging") account=677285201534;;
    "testnet") account=013339450148;; # public testnet
    "mainnet") account=332066407361;; # mainnet
    *) account=329916310755;;
esac

printf "account: %s\n" $account

dockerfile=Dockerfile.service.remote
case $service in
    "bazooka") dockerfile=Dockerfile.bazooka.remote;;
    "auxo") dockerfile=Dockerfile.auxo.remote;;
esac

printf "dockerfile: %s\n" $dockerfile

cd services/$service
pnpm build
cd -

AWS_PROFILE=dydx-v4-$env aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $account.dkr.ecr.ap-northeast-1.amazonaws.com

DOCKER_BUILDKIT=1 docker build \
            --platform linux/amd64 \
            -t $account.dkr.ecr.ap-northeast-1.amazonaws.com/$env-indexer-$service:$tag \
            -f $dockerfile \
            --build-arg service=$service \
            --build-arg NPM_TOKEN=$NPM_TOKEN .

docker push $account.dkr.ecr.ap-northeast-1.amazonaws.com/$env-indexer-$service:$tag

AWS_PROFILE=dydx-v4-$env orb deploy $service $tag
