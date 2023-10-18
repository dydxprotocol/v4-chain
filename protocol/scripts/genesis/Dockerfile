FROM dydxprotocol-base

RUN apk add --no-cache bash dasel jq

WORKDIR /workspace

ENTRYPOINT ["sh", "/workspace/scripts/genesis/check_sample_pregenesis_uptodate.sh"]
