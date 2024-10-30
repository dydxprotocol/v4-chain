"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __exportStar = (this && this.__exportStar) || function(m, exports) {
    for (var p in m) if (p !== "default" && !Object.prototype.hasOwnProperty.call(exports, p)) __createBinding(exports, m, p);
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Network = exports.ValidatorConfig = exports.IndexerConfig = exports.PAGE_REQUEST = exports.SHORT_BLOCK_FORWARD = exports.SHORT_BLOCK_WINDOW = exports.MAX_MEMO_CHARACTERS = exports.DEFAULT_API_TIMEOUT = exports.TimePeriod = exports.PositionStatus = exports.TickerType = exports.OrderStatus = exports.OrderExecution = exports.OrderTimeInForce = exports.OrderSide = exports.OrderType = exports.MarketStatisticDay = exports.NETWORK_ID_TESTNET = exports.NETWORK_ID_MAINNET = exports.NetworkId = exports.ValidatorApiHost = exports.FaucetApiHost = exports.IndexerWSHost = exports.IndexerApiHost = exports.LOCAL_CHAIN_ID = exports.TESTNET_CHAIN_ID = exports.STAGING_CHAIN_ID = exports.DEV_CHAIN_ID = void 0;
const long_1 = __importDefault(require("long"));
__exportStar(require("../lib/constants"), exports);
// Chain ID
exports.DEV_CHAIN_ID = 'dydxprotocol-testnet';
exports.STAGING_CHAIN_ID = 'dydxprotocol-testnet';
exports.TESTNET_CHAIN_ID = 'dydx-testnet-4';
exports.LOCAL_CHAIN_ID = 'localdydxprotocol';
// ------------ API URLs ------------
var IndexerApiHost;
(function (IndexerApiHost) {
    IndexerApiHost["TESTNET"] = "https://dydx-testnet.imperator.co";
    IndexerApiHost["LOCAL"] = "http://localhost:3002";
    // TODO: Add MAINNET
})(IndexerApiHost || (exports.IndexerApiHost = IndexerApiHost = {}));
var IndexerWSHost;
(function (IndexerWSHost) {
    IndexerWSHost["TESTNET"] = "wss://dydx-testnet.imperator.co/v4/ws";
    // TODO: Add MAINNET
    IndexerWSHost["LOCAL"] = "ws://localhost:3003";
})(IndexerWSHost || (exports.IndexerWSHost = IndexerWSHost = {}));
var FaucetApiHost;
(function (FaucetApiHost) {
    FaucetApiHost["TESTNET"] = "https://faucet.v4testnet.dydx.exchange";
})(FaucetApiHost || (exports.FaucetApiHost = FaucetApiHost = {}));
var ValidatorApiHost;
(function (ValidatorApiHost) {
    ValidatorApiHost["TESTNET"] = "https://test-dydx.kingnodes.com";
    // TODO: Add MAINNET
    ValidatorApiHost["LOCAL"] = "http://localhost:26657";
})(ValidatorApiHost || (exports.ValidatorApiHost = ValidatorApiHost = {}));
// ------------ Network IDs ------------
var NetworkId;
(function (NetworkId) {
    NetworkId["TESTNET"] = "dydx-testnet-4";
    // TODO: Add MAINNET
})(NetworkId || (exports.NetworkId = NetworkId = {}));
exports.NETWORK_ID_MAINNET = null;
exports.NETWORK_ID_TESTNET = 'dydxprotocol-testnet';
// ------------ Market Statistic Day Types ------------
var MarketStatisticDay;
(function (MarketStatisticDay) {
    MarketStatisticDay["ONE"] = "1";
    MarketStatisticDay["SEVEN"] = "7";
    MarketStatisticDay["THIRTY"] = "30";
})(MarketStatisticDay || (exports.MarketStatisticDay = MarketStatisticDay = {}));
// ------------ Order Types ------------
// This should match OrderType in Abacus
var OrderType;
(function (OrderType) {
    OrderType["LIMIT"] = "LIMIT";
    OrderType["MARKET"] = "MARKET";
    OrderType["STOP_LIMIT"] = "STOP_LIMIT";
    OrderType["TAKE_PROFIT_LIMIT"] = "TAKE_PROFIT";
    OrderType["STOP_MARKET"] = "STOP_MARKET";
    OrderType["TAKE_PROFIT_MARKET"] = "TAKE_PROFIT_MARKET";
})(OrderType || (exports.OrderType = OrderType = {}));
// ------------ Order Side ------------
// This should match OrderSide in Abacus
var OrderSide;
(function (OrderSide) {
    OrderSide["BUY"] = "BUY";
    OrderSide["SELL"] = "SELL";
})(OrderSide || (exports.OrderSide = OrderSide = {}));
// ------------ Order TimeInForce ------------
// This should match OrderTimeInForce in Abacus
var OrderTimeInForce;
(function (OrderTimeInForce) {
    OrderTimeInForce["GTT"] = "GTT";
    OrderTimeInForce["IOC"] = "IOC";
    OrderTimeInForce["FOK"] = "FOK";
})(OrderTimeInForce || (exports.OrderTimeInForce = OrderTimeInForce = {}));
// ------------ Order Execution ------------
// This should match OrderExecution in Abacus
var OrderExecution;
(function (OrderExecution) {
    OrderExecution["DEFAULT"] = "DEFAULT";
    OrderExecution["IOC"] = "IOC";
    OrderExecution["FOK"] = "FOK";
    OrderExecution["POST_ONLY"] = "POST_ONLY";
})(OrderExecution || (exports.OrderExecution = OrderExecution = {}));
// ------------ Order Status ------------
// This should match OrderStatus in Abacus
var OrderStatus;
(function (OrderStatus) {
    OrderStatus["BEST_EFFORT_OPENED"] = "BEST_EFFORT_OPENED";
    OrderStatus["OPEN"] = "OPEN";
    OrderStatus["FILLED"] = "FILLED";
    OrderStatus["BEST_EFFORT_CANCELED"] = "BEST_EFFORT_CANCELED";
    OrderStatus["CANCELED"] = "CANCELED";
})(OrderStatus || (exports.OrderStatus = OrderStatus = {}));
var TickerType;
(function (TickerType) {
    TickerType["PERPETUAL"] = "PERPETUAL";
})(TickerType || (exports.TickerType = TickerType = {}));
var PositionStatus;
(function (PositionStatus) {
    PositionStatus["OPEN"] = "OPEN";
    PositionStatus["CLOSED"] = "CLOSED";
    PositionStatus["LIQUIDATED"] = "LIQUIDATED";
})(PositionStatus || (exports.PositionStatus = PositionStatus = {}));
// ----------- Time Period for Sparklines -------------
var TimePeriod;
(function (TimePeriod) {
    TimePeriod["ONE_DAY"] = "ONE_DAY";
    TimePeriod["SEVEN_DAYS"] = "SEVEN_DAYS";
})(TimePeriod || (exports.TimePeriod = TimePeriod = {}));
// ------------ API Defaults ------------
exports.DEFAULT_API_TIMEOUT = 3000;
exports.MAX_MEMO_CHARACTERS = 256;
exports.SHORT_BLOCK_WINDOW = 20;
exports.SHORT_BLOCK_FORWARD = 3;
// Querying
exports.PAGE_REQUEST = {
    key: new Uint8Array(),
    offset: long_1.default.UZERO,
    limit: long_1.default.MAX_UNSIGNED_VALUE,
    countTotal: true,
    reverse: false,
};
class IndexerConfig {
    constructor(restEndpoint, websocketEndpoint) {
        this.restEndpoint = restEndpoint;
        this.websocketEndpoint = websocketEndpoint;
    }
}
exports.IndexerConfig = IndexerConfig;
class ValidatorConfig {
    constructor(restEndpoint, chainId, denoms, broadcastOptions) {
        this.restEndpoint = (restEndpoint === null || restEndpoint === void 0 ? void 0 : restEndpoint.endsWith('/')) ? restEndpoint.slice(0, -1) : restEndpoint;
        this.chainId = chainId;
        this.denoms = denoms;
        this.broadcastOptions = broadcastOptions;
    }
}
exports.ValidatorConfig = ValidatorConfig;
class Network {
    constructor(env, indexerConfig, validatorConfig) {
        this.env = env;
        this.indexerConfig = indexerConfig;
        this.validatorConfig = validatorConfig;
    }
    static testnet() {
        const indexerConfig = new IndexerConfig(IndexerApiHost.TESTNET, IndexerWSHost.TESTNET);
        const validatorConfig = new ValidatorConfig(ValidatorApiHost.TESTNET, exports.TESTNET_CHAIN_ID, {
            CHAINTOKEN_DENOM: 'adv4tnt',
            TDAI_DENOM: 'utdai',
            TDAI_GAS_DENOM: 'utdai',
            TDAI_DECIMALS: 6,
            CHAINTOKEN_DECIMALS: 18,
        });
        return new Network('testnet', indexerConfig, validatorConfig);
    }
    static local() {
        const indexerConfig = new IndexerConfig(IndexerApiHost.LOCAL, IndexerWSHost.LOCAL);
        const validatorConfig = new ValidatorConfig(ValidatorApiHost.LOCAL, exports.LOCAL_CHAIN_ID, {
            CHAINTOKEN_DENOM: 'adv4tnt',
            TDAI_DENOM: 'utdai',
            TDAI_GAS_DENOM: 'utdai',
            TDAI_DECIMALS: 6,
            CHAINTOKEN_DECIMALS: 18,
        });
        console.log("LOGGING NEW NETWORK");
        return new Network('local', indexerConfig, validatorConfig);
    }
    // TODO: Add mainnet(): Network
    getString() {
        return this.env;
    }
}
exports.Network = Network;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29uc3RhbnRzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NsaWVudHMvY29uc3RhbnRzLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQ0EsZ0RBQXdCO0FBSXhCLG1EQUFpQztBQUVqQyxXQUFXO0FBQ0UsUUFBQSxZQUFZLEdBQUcsc0JBQXNCLENBQUM7QUFDdEMsUUFBQSxnQkFBZ0IsR0FBRyxzQkFBc0IsQ0FBQztBQUMxQyxRQUFBLGdCQUFnQixHQUFHLGdCQUFnQixDQUFDO0FBQ3BDLFFBQUEsY0FBYyxHQUFHLG1CQUFtQixDQUFDO0FBRWxELHFDQUFxQztBQUNyQyxJQUFZLGNBSVg7QUFKRCxXQUFZLGNBQWM7SUFDeEIsK0RBQTZDLENBQUE7SUFDN0MsaURBQStCLENBQUE7SUFDL0Isb0JBQW9CO0FBQ3RCLENBQUMsRUFKVyxjQUFjLDhCQUFkLGNBQWMsUUFJekI7QUFFRCxJQUFZLGFBSVg7QUFKRCxXQUFZLGFBQWE7SUFDdkIsa0VBQWlELENBQUE7SUFDakQsb0JBQW9CO0lBQ3BCLDhDQUE2QixDQUFBO0FBQy9CLENBQUMsRUFKVyxhQUFhLDZCQUFiLGFBQWEsUUFJeEI7QUFFRCxJQUFZLGFBRVg7QUFGRCxXQUFZLGFBQWE7SUFDdkIsbUVBQWtELENBQUE7QUFDcEQsQ0FBQyxFQUZXLGFBQWEsNkJBQWIsYUFBYSxRQUV4QjtBQUVELElBQVksZ0JBSVg7QUFKRCxXQUFZLGdCQUFnQjtJQUMxQiwrREFBMkMsQ0FBQTtJQUMzQyxvQkFBb0I7SUFDcEIsb0RBQWdDLENBQUE7QUFDbEMsQ0FBQyxFQUpXLGdCQUFnQixnQ0FBaEIsZ0JBQWdCLFFBSTNCO0FBRUQsd0NBQXdDO0FBRXhDLElBQVksU0FHWDtBQUhELFdBQVksU0FBUztJQUNuQix1Q0FBMEIsQ0FBQTtJQUMxQixvQkFBb0I7QUFDdEIsQ0FBQyxFQUhXLFNBQVMseUJBQVQsU0FBUyxRQUdwQjtBQUNZLFFBQUEsa0JBQWtC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