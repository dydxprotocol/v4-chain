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
elif [ "$1" == "dev5" ]; then
	aws_account="917958511744"
	ecr_repo="dev5-validator"
elif [ "$1" == "staging" ]; then
	aws_account="677285201534"
	ecr_repo="staging-validator"
elif [ "$1" == "testnet" ]; then
	aws_account="419937869548"
	ecr_repo="testnet-validator"
else
	echo "Usage: build-push-ecr.sh (dev|dev2|dev3|dev4|dev5|staging)"
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

# Note: the following is based on this [AWS user guide](https://docs.aws.amazon.com/AmazonECR/latest/userguide/docker-push-ecr-image.html).
aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin $aws_account.dkr.ecr.us-east-2.amazonaws.com

ecr="$aws_account.dkr.ecr.us-east-2.amazonaws.com/$ecr_repo"

current_time=$(date '+%Y%m%d%H%M')
commit_hash=$(git rev-parse --short=7 HEAD)

make localnet-build-amd64

docker_tag="$ecr:$current_time-$commit_hash-test-build"
docker_file="testing/testnet-dev/Dockerfile"

if [ "$1" == "staging" ]; then
	docker_file="testing/testnet-staging/Dockerfile"
fi

docker build \
	--platform linux/amd64 \
	-t $docker_tag \
	-f $docker_file \
	--progress plain .

docker push $docker_tag

# Build and push the snapshot image with the commit hash
docker_tag="$ecr-snapshot:$current_time-$commit_hash-test-build"
docker_file="testing/snapshotting/Dockerfile.snapshot"

docker build \
	--platform linux/amd64 \
	-t $docker_tag \
	-f $docker_file \
	--progress plain .

docker push $docker_tag
