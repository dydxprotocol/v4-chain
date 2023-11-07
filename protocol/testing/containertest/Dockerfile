FROM dydxprotocol-base

COPY ./testing/containertest/containertest.sh /dydxprotocol/containertest.sh
COPY ./testing/genesis.sh /dydxprotocol/genesis.sh
COPY ./daemons/pricefeed/client/constants/testdata /dydxprotocol/exchange_config
COPY ./testing/delaymsg_config /dydxprotocol/delaymsg_config

RUN /dydxprotocol/containertest.sh

ENTRYPOINT ["dydxprotocold"]
