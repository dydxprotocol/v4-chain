FROM dydxprotocol-base

RUN apk add --no-cache bash jq aws-cli
RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

COPY ./testing/current_version.sh /dydxprotocol/
COPY ./testing/mainnet/. /dydxprotocol/

ENV HOME /dydxprotocol
WORKDIR $HOME

RUN /dydxprotocol/mainnet.sh

ENTRYPOINT ["/dydxprotocol/start.sh"]
