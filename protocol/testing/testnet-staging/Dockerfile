FROM dydxprotocol-base

COPY ./testing/testnet-staging/staging.sh /dydxprotocol/staging.sh
COPY ./testing/genesis.sh /dydxprotocol/genesis.sh

RUN /dydxprotocol/staging.sh

ENTRYPOINT ["dydxprotocold"]
