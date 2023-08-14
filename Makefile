DOCKER := $(shell which docker)
protoVer=0.13.3
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-format:
	@$(protoImage) find ./proto -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-check-bc-breaking:
	@rm -rf ./.proto-export
	@$(protoImage) buf breaking --against .git#branch=$$(git merge-base HEAD origin/main)

.PHONY: proto-format proto-lint proto-check-bc-breaking
