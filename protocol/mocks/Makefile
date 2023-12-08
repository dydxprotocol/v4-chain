GOPATH=$(shell go env GOPATH)

COSMOS_VERSION=$(shell go list -m all | grep "github.com/dydxprotocol/cosmos-sdk" | awk '{print $$NF}')
COMETBFT_VERSION=$(shell go list -m all | grep "github.com/dydxprotocol/cometbft" | awk '{print $$NF}')
COSMOS_GOGOPROTO_VERSION=$(shell go list -m all | grep "github.com/cosmos/gogoproto" | awk '{print $$NF}')

mock-clean:
	@rm -f ./mocks/*.go

mock-gen:
	@go run github.com/vektra/mockery/v2 --name=InterfaceRegistry --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/codec --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Configurator --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=MultiStore --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/store/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=CacheMultiStore --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/store/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=AnteDecorator --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=TxConfig --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/client --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Msg --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=TxBuilder --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/client --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Keeper --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/x/bank/keeper --filename=BankKeeper.go --structname=BankKeeper --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Server --dir=$(GOPATH)/pkg/mod/github.com/cosmos/gogoproto@$(COSMOS_GOGOPROTO_VERSION)/grpc --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=AppOptions --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/server/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Logger --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cometbft@$(COMETBFT_VERSION)/libs/log --filename=logger.go --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=MsgRouter --dir=./lib --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=HealthCheckable --dir=./daemons/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PrepareBridgeKeeper --dir=./app/prepare --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PrepareClobKeeper --dir=./app/prepare --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PreparePerpetualsKeeper --dir=./app/prepare --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PreparePricesKeeper --dir=./app/prepare --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ProcessBridgeKeeper --dir=./app/process --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ProcessClobKeeper --dir=./app/process --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ProcessStakingKeeper --dir=./app/process --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ProcessPerpetualKeeper --dir=./app/process --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=MemClob --dir=./x/clob/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=BridgeKeeper --dir=./x/bridge/types --recursive --output=./mocks
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
	@go run github.com/vektra/mockery/v2 --name=Marshaler --dir=./indexer/common --recursive --output=./mocks
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
	@go run github.com/vektra/mockery/v2 --name=BridgeServiceClient --dir=./daemons/bridge/api --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=BridgeQueryClient --dir=./daemons/bridge/client/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=EthClient --dir=./daemons/bridge/client/types --recursive --output=./mocks
