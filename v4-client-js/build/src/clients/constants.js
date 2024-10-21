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
//export const LOCAL_CHAIN_ID = 'localdydxprotocol';
exports.LOCAL_CHAIN_ID = 'consu';
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29uc3RhbnRzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NsaWVudHMvY29uc3RhbnRzLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQ0EsZ0RBQXdCO0FBSXhCLG1EQUFpQztBQUVqQyxXQUFXO0FBQ0UsUUFBQSxZQUFZLEdBQUcsc0JBQXNCLENBQUM7QUFDdEMsUUFBQSxnQkFBZ0IsR0FBRyxzQkFBc0IsQ0FBQztBQUMxQyxRQUFBLGdCQUFnQixHQUFHLGdCQUFnQixDQUFDO0FBQ2pELG9EQUFvRDtBQUN2QyxRQUFBLGNBQWMsR0FBRyxPQUFPLENBQUM7QUFFdEMscUNBQXFDO0FBQ3JDLElBQVksY0FJWDtBQUpELFdBQVksY0FBYztJQUN4QiwrREFBNkMsQ0FBQTtJQUM3QyxpREFBK0IsQ0FBQTtJQUMvQixvQkFBb0I7QUFDdEIsQ0FBQyxFQUpXLGNBQWMsOEJBQWQsY0FBYyxRQUl6QjtBQUVELElBQVksYUFJWDtBQUpELFdBQVksYUFBYTtJQUN2QixrRUFBaUQsQ0FBQTtJQUNqRCxvQkFBb0I7SUFDcEIsOENBQTZCLENBQUE7QUFDL0IsQ0FBQyxFQUpXLGFBQWEsNkJBQWIsYUFBYSxRQUl4QjtBQUVELElBQVksYUFFWDtBQUZELFdBQVksYUFBYTtJQUN2QixtRUFBa0QsQ0FBQTtBQUNwRCxDQUFDLEVBRlcsYUFBYSw2QkFBYixhQUFhLFFBRXhCO0FBRUQsSUFBWSxnQkFJWDtBQUpELFdBQVksZ0JBQWdCO0lBQzFCLCtEQUEyQyxDQUFBO0lBQzNDLG9CQUFvQjtJQUNwQixvREFBZ0MsQ0FBQTtBQUNsQyxDQUFDLEVBSlcsZ0JBQWdCLGdDQUFoQixnQkFBZ0IsUUFJM0I7QUFFRCx3Q0FBd0M7QUFFeEMsSUFBWSxTQUdYO0FBSEQsV0FBWSxTQUFTO0lBQ25CLHVDQUEwQixDQUFBO0lBQzFCLG9CQUFvQjtBQUN0QixDQUFDLEVBSFcsU0FBUyx5QkFBVCxTQUFTLFFBR3BCO0FBQ1ks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