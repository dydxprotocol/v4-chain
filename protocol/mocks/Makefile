# needs to be called from the protocol directory using this command make -f mocks/Makefile mock-gen

GOPATH=$(shell go env GOPATH)

COSMOS_VERSION=$(shell go list -m all | grep "github.com/StreamFinance-Protocol/cosmos-sdk[^/]" | awk '{print $$NF}')
COSMOS_STORE_VERSION=$(shell go list -m all | grep "cosmossdk.io/store[^/]" | awk '{print $$NF}')
COSMOS_LOG_VERSION=$(shell go list -m all | grep "cosmossdk.io/log[^/]" | awk '{print $$NF}')
COSMOS_GOGOPROTO_VERSION=$(shell go list -m all | grep "github.com/cosmos/gogoproto[^/]" | awk '{print $$NF}')


mock-clean:
	@rm -f ./mocks/*.go

mock-gen:
	@echo "COSMOS_VERSION=$(COSMOS_VERSION)"
	@echo "COSMOS_STORE_VERSION=$(COSMOS_STORE_VERSION)"
	@echo "COSMOS_LOG_VERSION=$(COSMOS_LOG_VERSION)"
	@echo "COSMOS_GOGOPROTO_VERSION=$(COSMOS_GOGOPROTO_VERSION)"
	@echo "GOPATH=$(GOPATH)"
	@go run github.com/vektra/mockery/v2 --name=Configurator --dir=$(GOPATH)/pkg/mod/github.com/!stream!finance-!protocol/cosmos-sdk@$(COSMOS_VERSION)/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=MultiStore --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk/store@$(COSMOS_STORE_VERSION)/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=CacheMultiStore --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk/store@$(COSMOS_STORE_VERSION)/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=AnteDecorator --dir=$(GOPATH)/pkg/mod/github.com/!stream!finance-!protocol/cosmos-sdk@$(COSMOS_VERSION)/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=TxConfig --dir=$(GOPATH)/pkg/mod/github.com/!stream!finance-!protocol/cosmos-sdk@$(COSMOS_VERSION)/client --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=TxBuilder --dir=$(GOPATH)/pkg/mod/github.com/!stream!finance-!protocol/cosmos-sdk@$(COSMOS_VERSION)/client --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Keeper --dir=$(GOPATH)/pkg/mod/github.com/!stream!finance-!protocol/cosmos-sdk@$(COSMOS_VERSION)/x/bank/keeper --filename=BankKeeper.go --structname=BankKeeper --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Server --dir=$(GOPATH)/pkg/mod/github.com/cosmos/gogoproto@$(COSMOS_GOGOPROTO_VERSION)/grpc --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=AppOptions --dir=$(GOPATH)/pkg/mod/github.com/!stream!finance-!protocol/cosmos-sdk@$(COSMOS_VERSION)/server/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Logger --dir=$(GOPATH)/pkg/mod/cosmossdk.io/log@$(COSMOS_LOG_VERSION) --filename=logger.go --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=MsgRouter --dir=./lib --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=HealthCheckable --dir=./daemons/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PrepareClobKeeper --dir=./app/prepare --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PreparePerpetualsKeeper --dir=./app/prepare --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PreBlockExecPricesKeeper --dir=./app/ve --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PrepareConsumerKeeper --dir=./app/prepare/ --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ProcessClobKeeper --dir=./app/process --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ProcessStakingKeeper --dir=./app/process --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ProcessPerpetualKeeper --dir=./app/process --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ProcessConsumerKeeper --dir=./app/process --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ProcessProposalVEApplier --dir=./app/process --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ProcessConsumerKeeper --dir=./app/process --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ExtendVotePerpetualsKeeper --dir=./app/ve --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ExtendVoteClobKeeper --dir=./app/ve --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=CCValidatorStore --dir=./app/ve/math --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=VEApplierPricesKeeper --dir=./app/ve/applier --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=VEApplierRatelimitKeeper --dir=./app/ve/applier --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=VoteAggregator --dir=./app/ve/aggregator --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=VEApplierInterface --dir=./app/ve --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=VoteExtensionRateLimitKeeper --dir=./app/ve --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=MemClob --dir=./x/clob/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=DelayMsgKeeper --dir=./x/delaymsg/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ClobKeeper --dir=./x/clob/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=MemClobKeeper --dir=./x/clob/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PricesKeeper --dir=./x/prices/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PerpetualsKeeper --dir=./x/perpetuals/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=SendingKeeper --dir=./x/sending/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=SubaccountsKeeper --dir=./x/subaccounts/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=FileHandler --dir=./daemons/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=GrpcServer --dir=./daemons/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=GrpcClient --dir=./daemons/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=TimeProvider --dir=./lib/time --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=QueryClient --dir=./testutil/grpc --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=QueryServer --dir=./testutil/grpc --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ExchangeQueryHandler --dir=./daemons/pricefeed/client/handler --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=RequestHandler --dir=./daemons/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ExchangeToMarketPrices --dir=./daemons/pricefeed/client/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=IndexerMessageSender --dir=./indexer --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PerpetualsClobKeeper --dir=x/perpetuals/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=IndexerEventManager --dir=./indexer/indexer_manager --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PricefeedMutableMarketConfigs --dir=./daemons/pricefeed/client/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ExchangeConfigUpdater --dir=./daemons/pricefeed/client/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=EthQueryClient --dir=./daemons/sdaioracle/client/eth_query_client --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=SDAIServiceClient --dir=./daemons/sdaioracle/api --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ICS4Wrapper --dir=x/ratelimit/types --recursive --output=./mocks --replace-type github.com/cosmos/ibc-go/v8/modules/core/02-client/types=clienttypes:github.com/cosmos/ibc-go/v8/modules/core/02-client/types