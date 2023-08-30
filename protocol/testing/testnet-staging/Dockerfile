FROM dydxprotocol-base

COPY ./testing/testnet-staging/staging.sh /dydxprotocol/staging.sh
COPY ./testing/genesis.sh /dydxprotocol/genesis.sh
COPY ./daemons/pricefeed/client/constants/testdata /dydxprotocol/exchange_config
COPY ./testing/delaymsg_config /dydxprotocol/delaymsg_config

RUN /dydxprotocol/staging.sh

ENTRYPOINT ["dydxprotocold"]
