FROM dydxprotocol-base

COPY ./testing/testnet-dev/dev.sh /dydxprotocol/dev.sh
COPY ./testing/genesis.sh /dydxprotocol/genesis.sh
COPY ./daemons/pricefeed/client/constants/testdata /dydxprotocol/exchange_config

RUN /dydxprotocol/dev.sh

ENTRYPOINT ["dydxprotocold"]
