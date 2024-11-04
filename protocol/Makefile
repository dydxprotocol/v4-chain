export VERSION := $(shell echo $(shell git describe --tags --always --match "protocol/v*") | sed 's/^protocol\/v//')
BUILDDIR ?= $(CURDIR)/build
DOCKER := $(shell which docker)
GO_VERSION := $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f 2)
PROJECT_NAME := $(shell git remote get-url origin | xargs basename -s .git)
CGO_ENABLED=1
export GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
export COMMIT=$(shell git rev-parse HEAD)

export GO111MODULE = on

# process build tags
build_tags = netgo

ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += rocksdb
endif
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += boltdb
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=klyraprotocol \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=klyraprotocold \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  $(info ################################################################)
  $(info To use rocksdb, you need to install rocksdb first)
  $(info Please follow this guide https://github.com/rockset/rocksdb-cloud/blob/master/INSTALL.md)
  $(info ################################################################)
  CGO_ENABLED=1
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=boltdb
endif

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ifeq ($(LINK_STATICALLY),true)
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

all: tools install lint test

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/klyraprotocold

aws-push-dev:
	./scripts/build-push-ecr.sh dev

aws-push-dev2:
	./scripts/build-push-ecr.sh dev2

aws-push-dev3:
	./scripts/build-push-ecr.sh dev3

aws-push-dev4:
	./scripts/build-push-ecr.sh dev4

aws-push-dev5:
	./scripts/build-push-ecr.sh dev5

aws-push-staging:
	./scripts/build-push-ecr.sh staging

###############################################################################
###                                  Build                                  ###
###############################################################################

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

build: go.sum
ifeq ($(OS),Windows_NT)
	exit 1
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/klyraprotocold ./cmd/klyraprotocold
endif

# TODO(DEC-1743): `build-reproducible` is taken from Osmosis. This currently takes a long time,
# since cross-building arm64 from amd64 (or vice-versa) is slow due to QEMU virtualization.
# We should look into making this faster. Gaia and other Cosmos SDK chains use `rbuilder`,
# which doesn't support go 1.19 yet.
build-reproducible: build-reproducible-linux-amd64 build-reproducible-linux-arm64
build-reproducible-linux-amd64: go.sum $(BUILDDIR)/
	$(DOCKER) buildx create --name klyrav4builder || true
	$(DOCKER) buildx use klyrav4builder
	$(DOCKER) buildx build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg RUNNER_IMAGE=alpine:3.16 \
		--platform linux/amd64 \
		-t klyra:local-linux-amd64 \
		--load \
		-f Dockerfile .
	$(DOCKER) rm -f klyrav4binary || true
	$(DOCKER) create -ti --name klyrav4binary klyra:local-linux-amd64
	$(DOCKER) cp klyrav4binary:/bin/klyraprotocold $(BUILDDIR)/klyraprotocold:linux-amd64
	$(DOCKER) rm -f klyrav4binary
build-reproducible-linux-arm64: go.sum $(BUILDDIR)/
	$(DOCKER) buildx create --name klyrav4builder || true
	$(DOCKER) buildx use klyrav4builder
	$(DOCKER) buildx build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg RUNNER_IMAGE=alpine:3.16 \
		--platform linux/arm64 \
		-t klyra:local-linux-arm64 \
		--load \
		-f Dockerfile .
	$(DOCKER) rm -f klyrav4binary || true
	$(DOCKER) create -ti --name klyrav4binary klyra:local-linux-arm64
	$(DOCKER) cp klyrav4binary:/bin/klyraprotocold $(BUILDDIR)/klyraprotocold:linux-arm64
	$(DOCKER) rm -f klyrav4binary
# TODO(DEC-1743): Make this also a reproducible build.
build-darwin-arm64: go.sum
	LEDGER_ENABLED=false GOOS=darwin GOARCH=arm64 $(MAKE) build -ldflags="${ldflags}"
draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/klyraprotocold -d 2 | dot -Tpng -o dependency-graph.png

################################################################################
### 						Tools & dependencies							 ###
################################################################################

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

distclean: clean tools-clean
clean:
	rm -rf \
    $(BUILDDIR)/ \
    artifacts/ \
    tmp-swagger-gen/ \
	swagger-proto

.PHONY: distclean clean

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

include sims.mk
include ./mocks/Makefile

remove-container:
	@echo "removing containers"
	cd ../e2e-testing && pnpm clean

test: test-unit

test-ci: test-unit-and-integration

test-all: test-unit test-race test-cover test-unit-and-integration

test-slow:
	$(MAKE) test-unit-and-integration
	$(MAKE) test-race-slow
	$(MAKE) test-cover-slow
	$(MAKE) test-exchange-slow

test-unit:
	@VERSION=$(VERSION) go test -count=1 -mod=readonly -timeout 20m -tags='test_ledger_mock $(build_tags)' ./...

test-unit-specific:
	@echo "Running specific test: $(TESTNAME) in file: $(FILE)"
	@VERSION=$(VERSION) go test -mod=readonly -tags='test_ledger_mock $(build_tags)' -timeout 20m -run "$(TESTNAME)" $(FILE)

test-unit-file:
	@echo "Running specific test: $(TESTNAME) in file: $(FILE)"
	@VERSION=$(VERSION) go test -mod=readonly -tags='test_ledger_mock $(build_tags)' -timeout 20m $(FILE)

test-unit-verbose:
	@VERSION=$(VERSION) go test -mod=readonly -timeout 20m -tags='test_ledger_mock $(build_tags)' -v -count=1 -failfast -cover ./...

specific-unit: 
	@VERSION=$(VERSION) go test -mod=readonly -tags='test_ledger_mock $(build_tags)' -p 1 -timeout 20m -run "TestGetNetCollateralAndMarginRequirements" ./x/subaccounts/keeper/...

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -timeout 20m -race -tags='test_ledger_mock $(build_tags)' ./...

test-race-slow:
	@VERSION=$(VERSION) go test -mod=readonly -timeout 20m -race -tags='test_ledger_mock $(build_tags)' -p 1 -count=1 ./...

test-cover:
	@VERSION=$(VERSION) go test -mod=readonly -timeout 12m -coverprofile=coverage.out -covermode=atomic -coverpkg=github.com/StreamFinance-Protocol/stream-chain/protocol/... -tags='test_ledger_mock $(build_tags)' ./...

test-cover-slow:
	@VERSION=$(VERSION) go test -mod=readonly -timeout 12m -coverprofile=coverage.out -covermode=atomic -coverpkg=github.com/StreamFinance-Protocol/stream-chain/protocol/... -tags='test_ledger_mock $(build_tags)' -p 1 -count=1 ./...

test-exchange-slow:
	@VERSION=$(VERSION) go test -mod=readonly -tags='test_ledger_mock exchange_tests $(build_tags)' -p 1 -count=1 github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/handler -run "TestQueryingActualExchanges"

test-unit-and-integration:
	@VERSION=$(VERSION) go test -mod=readonly -tags='integration_test test_ledger_mock $(build_tags)' ./...

specific-test: 
	@VERSION=$(VERSION) go test -mod=readonly -tags='integration_test test_ledger_mock $(build_tags)' -p 1 -timeout 80m ./x/ratelimit/client/cli/...

benchmark:
	@VERSION=$(VERSION) go test -mod=readonly -tags='$(build_tags)' -bench=. ./...

test-container:
	@SKIP_DISABLED=true VERSION=$(VERSION) go test -mod=readonly -tags='container_test $(build_tags)' ./testing/containertest

test-container-accept:
	@SKIP_DISABLED=true VERSION=$(VERSION) go test -mod=readonly -tags='container_test $(build_tags)' ./testing/containertest -args -accept

test-container-build:
	$(MAKE) localnet-build
	@docker build . -t klyraprotocol-container-test -f testing/containertest/Dockerfile --no-cache

.PHONY: test test-all test-cover test-unit test-race test-container test-slow
	\ container-test-build-base container-test-build-local

###############################################################################
###                                Coverage                                  ###
###############################################################################

get-coverage:
	@echo "Running tests and generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"
	@open coverage.html

get-file-coverage:
	@echo "Running tests and generating coverage report for directory: $(DIR)"
	@go test -coverprofile=coverage.out $(DIR)
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"
	@open coverage.html

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --build-tags='$(build_tags) all' --timeout=10m

lint-fix:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --build-tags='$(build_tags) all' --fix --out-format=tab --issues-exit-code=0
.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name '*.pb.go' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name '*.pb.go' | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name '*.pb.go' | xargs goimports -w -local github.com/cosmos/cosmos-sdk
.PHONY: format

###############################################################################
###                                Protobuf                                 ###
###############################################################################

# Note: update github action workflow `proto.yaml` when updating version or images.
protoVer=0.11.5
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --user root --rm -v $(CURDIR)/..:/workspace --workdir /workspace $(protoImageName)

proto-swagger-gen:
	@make clean
	@echo "Downloading Protobuf dependencies"
	@make proto-download-deps
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./protocol/scripts/protoc-swagger-gen.sh
	@echo "Generating static files for swagger docs"
	$(MAKE) update-swagger-docs

SWAGGER_DIR=./swagger-proto
THIRD_PARTY_DIR=$(SWAGGER_DIR)/third_party
proto-download-deps:
	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-sdk.git" && \
	git config core.sparseCheckout true && \
	printf "proto\nthird_party\n" > .git/info/sparse-checkout && \
	git pull origin main && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/ibc-go.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git pull origin main && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/ibc_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-proto.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git pull origin main && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_proto_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/gogoproto" && \
	curl -SSL https://raw.githubusercontent.com/cosmos/gogoproto/main/gogoproto/gogo.proto > "$(THIRD_PARTY_DIR)/gogoproto/gogo.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/google/api" && \
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > "$(THIRD_PARTY_DIR)/google/api/annotations.proto"
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > "$(THIRD_PARTY_DIR)/google/api/http.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos/ics23/v1" && \
	curl -sSL https://raw.githubusercontent.com/cosmos/ics23/master/proto/cosmos/ics23/v1/proofs.proto > "$(THIRD_PARTY_DIR)/cosmos/ics23/v1/proofs.proto"

.PHONY: proto-swagger-gen


###############################################################################
###                                Localnet                                 ###
###############################################################################

localnet-new-validator:
	@docker run --rm -it $(shell docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) . -t klyraprotocol-base && docker build . -q -f ./testing/util/Dockerfile)

localnet-build:
	docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) . -t klyraprotocol-base --no-cache

localnet-build-amd64:
	DOCKER_BUILDKIT=1 docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) -f Dockerfile . -t klyraprotocol-base --no-cache --platform linux/amd64

reset-chain:
	rm -rf localnet
	# Creates the directories which store the chain data for each node.
	mkdir -p localnet/klyraprotocol0
	mkdir -p localnet/klyraprotocol1
	mkdir -p localnet/klyraprotocol2
	mkdir -p localnet/klyraprotocol3
	# These files are created as part of `klyraprotocold init`, and are required to initialize the chain,
	# but mounting the volumes removes them. This is a workaround to create them again.
	echo '{"height": "0", "round": 0, "step": 0 }' > localnet/klyraprotocol0/priv_validator_state.json
	echo '{"height": "0", "round": 0, "step": 0 }' > localnet/klyraprotocol1/priv_validator_state.json
	echo '{"height": "0", "round": 0, "step": 0 }' > localnet/klyraprotocol2/priv_validator_state.json
	echo '{"height": "0", "round": 0, "step": 0 }' > localnet/klyraprotocol3/priv_validator_state.json

localnet-init:
	make reset-chain
	# Build the base image.
	make localnet-build

localnet-compose-up:
	@echo "Launching localnet at commit ${GIT_COMMIT_HASH}"
	@docker build . -t local:klyraprotocol -f testing/testnet-local/Dockerfile --no-cache
	@docker-compose -f docker-compose.yml up --force-recreate $(ARGS)

localnet-compose-upd:
	@echo "Launching localnet at commit ${GIT_COMMIT_HASH}"
	@docker build . -t local:klyraprotocol -f testing/testnet-local/Dockerfile --no-cache
	@docker-compose -f docker-compose.yml up --force-recreate -d $(ARGS)

build-e2etest-image:
	@echo "Build e2e test image at commit ${GIT_COMMIT_HASH}"
	@docker build . -t local:e2etest-klyraprotocol -f testing/e2etest-local/Dockerfile --no-cache

localnet-start: localnet-init localnet-compose-up
localnet-startd: localnet-init localnet-compose-upd

e2etest-build-image: localnet-init build-e2etest-image

# Continue the localnet with the same chain state.
localnet-continue:
	@docker-compose -f docker-compose.yml up $(ARGS)

localnet-stop:
	@docker-compose -f docker-compose.yml down

.PHONY: all build-linux install format lint \
	go-mod-cache draw-deps clean build build-contract-tests-hooks \
	test test-all test-build test-cover test-unit test-race benchmark

###############################################################################
###                              Documentation                              ###
###############################################################################

update-swagger-docs: statik
	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
		echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
		exit 1;\
	else \
		echo "\033[92mSwagger docs are in sync\033[0m";\
	fi

.PHONY: update-swagger-docs

###############################################################################
###                           Sample Pregenesis                             ###
###############################################################################

# Run at `./protocol` directory.
update-sample-pregenesis:
	$(MAKE) localnet-build
	@docker build . -t check-sample-pregenesis -f scripts/genesis/Dockerfile --no-cache
	@docker run --entrypoint ./scripts/genesis/copy_sample_pregenesis.sh  -v $(CURDIR):/workspace check-sample-pregenesis
	@echo "Updated ./scripts/genesis/sample_pregenesis.json"

check-sample-pregenesis-up-to-date:
	$(MAKE) localnet-build
	@docker build . -t check-sample-pregenesis -f scripts/genesis/Dockerfile --no-cache
	@docker run -v $(CURDIR):/workspace check-sample-pregenesis 

.PHONY: update-sample-pregenesis check-sample-pregenesis-up-to-date