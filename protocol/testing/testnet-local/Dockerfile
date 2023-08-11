FROM dydxprotocol-base

COPY ./testing/testnet-local/local.sh /dydxprotocol/local.sh
COPY ./testing/genesis.sh /dydxprotocol/genesis.sh

RUN /dydxprotocol/local.sh

ENTRYPOINT ["dydxprotocold"]
