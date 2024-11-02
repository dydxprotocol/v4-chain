FROM klyraprotocol-base

COPY ./testing/containertest/klyraprotocold_preupgrade* /bin/
COPY ./testing/containertest/containertest.sh /klyraprotocol/containertest.sh
COPY ./testing/containertest/preupgrade_entrypoint.sh /klyraprotocol/preupgrade_entrypoint.sh
COPY ./testing/containertest/preupgrade_genesis.json /klyraprotocol/preupgrade_genesis.json
COPY ./testing/genesis.sh /klyraprotocol/genesis.sh
COPY ./daemons/pricefeed/client/constants/testdata /klyraprotocol/exchange_config
COPY ./testing/delaymsg_config /klyraprotocol/delaymsg_config

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

RUN /klyraprotocol/containertest.sh

COPY ./testing/containertest/config/pricefeed_exchange_config.toml /klyraprotocol/chain/.alice/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /klyraprotocol/chain/.bob/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /klyraprotocol/chain/.carl/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /klyraprotocol/chain/.dave/config/pricefeed_exchange_config.toml

ENTRYPOINT ["klyraprotocold"]
