DOCKER := $(shell which docker)
protoVer=0.13.3
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./protocgen.sh

# Note: this should mirror `proto-gen`.
proto-gen-without-container:
	@echo "Generating Protobuf files"
	sh ./protocgen.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-export:
	rm -rf ./.proto-export && buf export --config ./buf.yaml --output ./.proto-export

# dydxjs needs cosmos protos that aren't depended on by dydx protos.
# Running the export command above will only get a subset of cosmos protos.
# For now, we will export cosmos protos separately. Note that this doesn't verify the digest.
# TODO(DEC-2027): Do this in a more reproducible way.
proto-export-dydxjs:
	rm -rf ./.proto-export
	buf export --config ./buf.yaml --output ./.proto-export --exclude-imports
	buf export buf.build/cosmos/cosmos-sdk:v0.47.0 --output ./.proto-export

proto-check-bc-breaking:
	@rm -rf ./.proto-export
	@$(protoImage) buf breaking --against .git#branch=$$(git merge-base HEAD origin/main)

.PHONY: proto-all proto-gen proto-format proto-lint proto-export proto-check-bc-breaking
