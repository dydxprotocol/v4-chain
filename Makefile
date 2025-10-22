DOCKER := $(shell which docker)
protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-format:
	@$(protoImage) find ./proto -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./protocol/scripts/protocgen.sh

proto-clean:
	@rm -rf ./proto/.proto-export
	@rm -rf ./.proto-export

proto-gen-clean:
	@echo "Cleaning old artifacts"
	@rm -rf ./proto/.proto-export
	@rm -rf ./.proto-export
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./protocol/scripts/protocgen.sh
	@cd proto && make proto-export-v4-proto-js

proto-check-bc-breaking:
	@rm -rf ./.proto-export
	@$(protoImage) buf breaking --against .git#branch=$$(git merge-base HEAD origin/main)

proto-export:
	@rm -rf ./.proto-export && cd proto && buf export --config ./buf.yaml --output ../.proto-export

proto-export-deps:
	@rm -rf ./.proto-export-deps
	@cd proto && buf export --config ./buf.yaml --output ../.proto-export-deps --exclude-imports && buf export buf.build/cosmos/cosmos-sdk:v0.50.0 --output ../.proto-export-deps

PROTO_DIRS=$(shell find .proto-export-deps -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

# The perl replace script is used to fix import statements (hacking around https://github.com/protocolbuffers/protobuf/issues/2283).
# We exclude google.protobuf because it is a part of the protobuf python package.
# We can't use sed here because we're using negative lookahead to exclude google.protobuf.
v4-proto-py-gen: proto-export-deps
	@rm -rf ./v4-proto-py/v4_proto
	@mkdir -p ./v4-proto-py/v4_proto
	@for dir in $(PROTO_DIRS); do \
		python3 -m grpc_tools.protoc \
		-I .proto-export-deps \
		--python_out=./v4-proto-py/v4_proto \
		--pyi_out=./v4-proto-py/v4_proto \
		--grpc_python_out=./v4-proto-py/v4_proto \
		$$(find ./$${dir} -type f -name '*.proto'); \
	done;
	perl -i -pe 's/^from (?!google\.protobuf)([^ ]*) import ([^ ]*)_pb2 as ([^ ]*)$$/from v4_proto.\1 import \2_pb2 as \3/' $$(find ./v4-proto-py/v4_proto -type f \( -name '*_pb2.py' -o -name '*_pb2_grpc.py' -o -name '*_pb2.pyi' -o -name '*_pb2_grpc.pyi' \))

.PHONY: proto-format proto-lint proto-check-bc-breaking proto-export proto-export-deps v4-proto-py-gen
