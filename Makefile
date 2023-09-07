DOCKER := $(shell which docker)
protoVer=0.13.3
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-format:
	@$(protoImage) find ./proto -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./protocol/scripts/protocgen.sh

proto-check-bc-breaking:
	@rm -rf ./.proto-export
	@$(protoImage) buf breaking --against .git#branch=$$(git merge-base HEAD origin/main)

proto-export:
	@rm -rf proto/.proto-export && cd proto && buf export --config ./buf.yaml --output ../.proto-export

proto-export-deps:
	@rm -rf ./.proto-export-deps
	@cd proto && buf export --config ./buf.yaml --output ../.proto-export-deps --exclude-imports && buf export buf.build/cosmos/cosmos-sdk:v0.47.0 --output ../.proto-export-deps

PROTO_DIRS=$(shell find .proto-export-deps -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

v4-proto-py-gen: proto-export-deps
	@rm -rf ./v4-proto-py/v4-proto-py
	@mkdir -p ./v4-proto-py/v4-proto-py
	@for dir in $(PROTO_DIRS); do \
		python3 -m grpc_tools.protoc \
		-I .proto-export-deps \
		--python_out=./v4-proto-py/v4-proto-py \
		--pyi_out=./v4-proto-py/v4-proto-py \
		--grpc_python_out=./v4-proto-py/v4-proto-py \
		$$(find ./$${dir} -type f -name '*.proto'); \
	done; \
	touch v4-proto-py/v4-proto-py/__init__.py

.PHONY: proto-format proto-lint proto-check-bc-breaking proto-export proto-export-deps v4-proto-py-gen
