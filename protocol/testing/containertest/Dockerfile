FROM dydxprotocol-base

COPY ./testing/containertest/containertest.sh /dydxprotocol/containertest.sh
COPY ./testing/genesis.sh /dydxprotocol/genesis.sh
COPY ./daemons/pricefeed/client/constants/testdata /dydxprotocol/exchange_config
COPY ./testing/delaymsg_config /dydxprotocol/delaymsg_config

RUN /dydxprotocol/containertest.sh

COPY ./testing/containertest/config/pricefeed_exchange_config.toml /dydxprotocol/chain/.alice/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /dydxprotocol/chain/.bob/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /dydxprotocol/chain/.carl/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /dydxprotocol/chain/.dave/config/pricefeed_exchange_config.toml

ENTRYPOINT ["dydxprotocold"]
