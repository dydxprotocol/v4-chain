#!/bin/bash

# Builds and pushes Docker container image to ECR for "dev" or "staging" AWS accounts.
# User's current AWS_PROFILE must match the desired environment.
# Dependencies:
#  - jq
#  - aws-cli
#  - docker
#  - git

make clean

if [ "$1" == "dev" ]; then
	aws_account="329916310755"
	ecr_repo="dev-validator"
elif [ "$1" == "dev2" ]; then
	aws_account="220669270420"
	ecr_repo="dev2-validator"
elif [ "$1" == "dev3" ]; then
	aws_account="295746761472"
	ecr_repo="dev3-validator"
elif [ "$1" == "dev4" ]; then
	aws_account="525975847385"
	ecr_repo="dev4-validator"
elif [ "$1" == "staging" ]; then
	aws_account="677285201534"
	ecr_repo="staging-validator"
else
	echo "Usage: build-push-ecr.sh (dev|dev2|dev3|dev4|staging)"
	exit 1
fi

if ! command -v jq &>/dev/null; then
	echo "Program 'jq' required. Run 'brew install jq'"
	exit
fi

caller_identity=$(aws sts get-caller-identity | jq -r .Account)

if [ "$caller_identity" != "$aws_account" ]; then
	echo "Your AWS caller identity doesn't match the desired environment. Check your ~/.aws/credentials file to ensure you have the correct AWS credentials for the environment. Use the AWS_PROFILE env var to switch AWS profiles between dev/staging before running this script again."
	exit 1
fi

read -p "Locally build and push Docker container image to ECR for AWS account ${aws_account}? y/n" -n 1 -r
echo # new line
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
	exit 1
fi

ecr="$aws_account.dkr.ecr.us-east-2.amazonaws.com/$ecr_repo"

current_time=$(date '+%Y%m%d%H%M')
commit_hash=$(git rev-parse --short=7 HEAD)

DOCKER_BUILDKIT=1 docker build \
	--platform linux/amd64 \
	-t dydxprotocol-base \
	-f Dockerfile .

if [ "$1" == "dev" ]; then
	docker build \
		--platform linux/amd64 \
		-t "$ecr:$current_time-$commit_hash-test-build" \
		-f testing/testnet-dev/Dockerfile \
		--progress plain .
fi

if [ "$1" == "staging" ]; then
	docker build \
		--platform linux/amd64 \
		-t "$ecr:$current_time-$commit_hash-test-build" \
		-f testing/testnet-staging/Dockerfile \
		--progress plain .
fi

# Note: the following is based on this [AWS user guide](https://docs.aws.amazon.com/AmazonECR/latest/userguide/docker-push-ecr-image.html).
aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin $aws_account.dkr.ecr.us-east-2.amazonaws.com

docker push "$ecr:$current_time-$commit_hash-test-build"
