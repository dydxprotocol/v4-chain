FROM dydxprotocol-base

COPY ./testing/testnet-staging/staging.sh /dydxprotocol/staging.sh
COPY ./testing/genesis.sh /dydxprotocol/genesis.sh
COPY ./testing/start.sh /dydxprotocol/start.sh
COPY ./daemons/pricefeed/client/constants/testdata /dydxprotocol/exchange_config
COPY ./testing/delaymsg_config /dydxprotocol/delaymsg_config

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

RUN /dydxprotocol/staging.sh

ENV DAEMON_NAME=dydxprotocold
ENTRYPOINT ["cosmovisor", "run"]
