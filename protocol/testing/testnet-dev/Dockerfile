FROM dydxprotocol-base

COPY ./testing/testnet-dev/dev.sh /dydxprotocol/dev.sh
COPY ./testing/genesis.sh /dydxprotocol/genesis.sh

RUN /dydxprotocol/dev.sh

ENTRYPOINT ["dydxprotocold"]
