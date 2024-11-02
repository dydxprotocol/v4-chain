FROM klyraprotocol-base

COPY ./testing/testnet-dev/dev.sh /klyraprotocol/dev.sh
COPY ./testing/genesis.sh /klyraprotocol/genesis.sh
COPY ./testing/start.sh /klyraprotocol/start.sh
COPY ./daemons/pricefeed/client/constants/testdata /klyraprotocol/exchange_config
COPY ./testing/delaymsg_config /klyraprotocol/delaymsg_config

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

RUN /klyraprotocol/dev.sh

ENV DAEMON_NAME=klyraprotocold
ENV DAEMON_HOME=/klyraprotocol/chain/.alice

ENTRYPOINT ["cosmovisor", "run"]
