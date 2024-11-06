FROM klyraprotocol-base

COPY ./testing/testnet-local/local.sh /klyraprotocol/local.sh
COPY ./testing/genesis.sh /klyraprotocol/genesis.sh
COPY ./testing/start.sh /klyraprotocol/start.sh
COPY ./daemons/pricefeed/client/constants/testdata /klyraprotocol/exchange_config
COPY ./testing/delaymsg_config /klyraprotocol/delaymsg_config

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

RUN /klyraprotocol/local.sh

ENV DAEMON_NAME=klyraprotocold
ENTRYPOINT ["cosmovisor", "run"]
