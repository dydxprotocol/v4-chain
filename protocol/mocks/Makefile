GOPATH=$(shell go env GOPATH)

COSMOS_VERSION=v0.47.2-0.20230501215803-6477dcf56934
COMETBFT_VERSION=v0.37.0-rc2.0.20230501193110-4d4d3b0c0995

mock-clean:
	@rm -f ./mocks/*.go

mock-gen:
	@go run github.com/vektra/mockery/v2 --name=InterfaceRegistry --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/codec --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Configurator --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=AnteDecorator --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=TxConfig --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/client --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=TxBuilder --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/client --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Keeper --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/x/bank/keeper --filename=BankKeeper.go --structname=BankKeeper --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Server --dir=$(GOPATH)/pkg/mod/github.com/cosmos/gogoproto@v1.4.6/grpc --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=AppOptions --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cosmos-sdk@$(COSMOS_VERSION)/server/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Logger --dir=$(GOPATH)/pkg/mod/github.com/dydxprotocol/cometbft@$(COMETBFT_VERSION)/libs/log --filename=logger.go --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PrepareClobKeeper --dir=./app/prepare --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PreparePerpetualsKeeper --dir=./app/prepare --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PreparePricesKeeper --dir=./app/prepare --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ContextHelper --dir=./app/prepare --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=MemClob --dir=./x/clob/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ClobKeeper --dir=./x/clob/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=MemClobKeeper --dir=./x/clob/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PricesKeeper --dir=./x/prices/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PerpetualsKeeper --dir=./x/perpetuals/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=SendingKeeper --dir=./x/sending/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=SubaccountsKeeper --dir=./x/subaccounts/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=FileHandler --dir=./lib --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=GrpcServer --dir=./lib --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=GrpcClient --dir=./lib --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=TimeProvider --dir=./lib --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Marshaler --dir=./indexer/common --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=QueryClient --dir=./testutil/grpc --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=ExchangeQueryHandler --dir=./daemons/pricefeed/client/handler --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=RequestHandler --dir=./lib --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=Medianizer --dir=./lib --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=IndexerMessageSender --dir=./indexer --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=PricePremiumGetter --dir=x/perpetuals/types --recursive --output=./mocks
	@go run github.com/vektra/mockery/v2 --name=IndexerEventManager --dir=./indexer/indexer_manager --recursive --output=./mocks
