FROM dydxprotocol-base

COPY ./testing/testnet-local/local.sh /dydxprotocol/local.sh
COPY ./testing/genesis.sh /dydxprotocol/genesis.sh
COPY ./daemons/pricefeed/client/constants/testdata /dydxprotocol/exchange_config

RUN /dydxprotocol/local.sh

ENTRYPOINT ["dydxprotocold"]
