FROM dydxprotocol-base

COPY ./testing/containertest/dydxprotocold_preupgrade* /bin/
COPY ./testing/containertest/containertest.sh /dydxprotocol/containertest.sh
COPY ./testing/containertest/preupgrade_entrypoint.sh /dydxprotocol/preupgrade_entrypoint.sh
COPY ./testing/containertest/preupgrade_genesis.json /dydxprotocol/preupgrade_genesis.json
COPY ./testing/genesis.sh /dydxprotocol/genesis.sh
COPY ./daemons/pricefeed/client/constants/testdata /dydxprotocol/exchange_config
COPY ./testing/delaymsg_config /dydxprotocol/delaymsg_config
COPY ./testing/version/. /dydxprotocol/

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

RUN /dydxprotocol/containertest.sh

COPY ./testing/containertest/config/pricefeed_exchange_config.toml /dydxprotocol/chain/.alice/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /dydxprotocol/chain/.bob/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /dydxprotocol/chain/.carl/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /dydxprotocol/chain/.dave/config/pricefeed_exchange_config.toml

ENTRYPOINT ["dydxprotocold"]
