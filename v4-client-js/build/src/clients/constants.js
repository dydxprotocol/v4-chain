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
})(IndexerApiHost = exports.IndexerApiHost || (exports.IndexerApiHost = {}));
var IndexerWSHost;
(function (IndexerWSHost) {
    IndexerWSHost["TESTNET"] = "wss://dydx-testnet.imperator.co/v4/ws";
    // TODO: Add MAINNET
    IndexerWSHost["LOCAL"] = "ws://localhost:3003";
})(IndexerWSHost = exports.IndexerWSHost || (exports.IndexerWSHost = {}));
var FaucetApiHost;
(function (FaucetApiHost) {
    FaucetApiHost["TESTNET"] = "https://faucet.v4testnet.dydx.exchange";
})(FaucetApiHost = exports.FaucetApiHost || (exports.FaucetApiHost = {}));
var ValidatorApiHost;
(function (ValidatorApiHost) {
    ValidatorApiHost["TESTNET"] = "https://test-dydx.kingnodes.com";
    // TODO: Add MAINNET
    ValidatorApiHost["LOCAL"] = "http://localhost:26657";
})(ValidatorApiHost = exports.ValidatorApiHost || (exports.ValidatorApiHost = {}));
// ------------ Network IDs ------------
var NetworkId;
(function (NetworkId) {
    NetworkId["TESTNET"] = "dydx-testnet-4";
    // TODO: Add MAINNET
})(NetworkId = exports.NetworkId || (exports.NetworkId = {}));
exports.NETWORK_ID_MAINNET = null;
exports.NETWORK_ID_TESTNET = 'dydxprotocol-testnet';
// ------------ Market Statistic Day Types ------------
var MarketStatisticDay;
(function (MarketStatisticDay) {
    MarketStatisticDay["ONE"] = "1";
    MarketStatisticDay["SEVEN"] = "7";
    MarketStatisticDay["THIRTY"] = "30";
})(MarketStatisticDay = exports.MarketStatisticDay || (exports.MarketStatisticDay = {}));
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
})(OrderType = exports.OrderType || (exports.OrderType = {}));
// ------------ Order Side ------------
// This should match OrderSide in Abacus
var OrderSide;
(function (OrderSide) {
    OrderSide["BUY"] = "BUY";
    OrderSide["SELL"] = "SELL";
})(OrderSide = exports.OrderSide || (exports.OrderSide = {}));
// ------------ Order TimeInForce ------------
// This should match OrderTimeInForce in Abacus
var OrderTimeInForce;
(function (OrderTimeInForce) {
    OrderTimeInForce["GTT"] = "GTT";
    OrderTimeInForce["IOC"] = "IOC";
    OrderTimeInForce["FOK"] = "FOK";
})(OrderTimeInForce = exports.OrderTimeInForce || (exports.OrderTimeInForce = {}));
// ------------ Order Execution ------------
// This should match OrderExecution in Abacus
var OrderExecution;
(function (OrderExecution) {
    OrderExecution["DEFAULT"] = "DEFAULT";
    OrderExecution["IOC"] = "IOC";
    OrderExecution["FOK"] = "FOK";
    OrderExecution["POST_ONLY"] = "POST_ONLY";
})(OrderExecution = exports.OrderExecution || (exports.OrderExecution = {}));
// ------------ Order Status ------------
// This should match OrderStatus in Abacus
var OrderStatus;
(function (OrderStatus) {
    OrderStatus["BEST_EFFORT_OPENED"] = "BEST_EFFORT_OPENED";
    OrderStatus["OPEN"] = "OPEN";
    OrderStatus["FILLED"] = "FILLED";
    OrderStatus["BEST_EFFORT_CANCELED"] = "BEST_EFFORT_CANCELED";
    OrderStatus["CANCELED"] = "CANCELED";
})(OrderStatus = exports.OrderStatus || (exports.OrderStatus = {}));
var TickerType;
(function (TickerType) {
    TickerType["PERPETUAL"] = "PERPETUAL";
})(TickerType = exports.TickerType || (exports.TickerType = {}));
var PositionStatus;
(function (PositionStatus) {
    PositionStatus["OPEN"] = "OPEN";
    PositionStatus["CLOSED"] = "CLOSED";
    PositionStatus["LIQUIDATED"] = "LIQUIDATED";
})(PositionStatus = exports.PositionStatus || (exports.PositionStatus = {}));
// ----------- Time Period for Sparklines -------------
var TimePeriod;
(function (TimePeriod) {
    TimePeriod["ONE_DAY"] = "ONE_DAY";
    TimePeriod["SEVEN_DAYS"] = "SEVEN_DAYS";
})(TimePeriod = exports.TimePeriod || (exports.TimePeriod = {}));
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29uc3RhbnRzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NsaWVudHMvY29uc3RhbnRzLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQ0EsZ0RBQXdCO0FBSXhCLG1EQUFpQztBQUVqQyxXQUFXO0FBQ0UsUUFBQSxZQUFZLEdBQUcsc0JBQXNCLENBQUM7QUFDdEMsUUFBQSxnQkFBZ0IsR0FBRyxzQkFBc0IsQ0FBQztBQUMxQyxRQUFBLGdCQUFnQixHQUFHLGdCQUFnQixDQUFDO0FBQ3BDLFFBQUEsY0FBYyxHQUFHLG1CQUFtQixDQUFDO0FBRWxELHFDQUFxQztBQUNyQyxJQUFZLGNBSVg7QUFKRCxXQUFZLGNBQWM7SUFDeEIsK0RBQTZDLENBQUE7SUFDN0MsaURBQStCLENBQUE7SUFDL0Isb0JBQW9CO0FBQ3RCLENBQUMsRUFKVyxjQUFjLEdBQWQsc0JBQWMsS0FBZCxzQkFBYyxRQUl6QjtBQUVELElBQVksYUFJWDtBQUpELFdBQVksYUFBYTtJQUN2QixrRUFBaUQsQ0FBQTtJQUNqRCxvQkFBb0I7SUFDcEIsOENBQTZCLENBQUE7QUFDL0IsQ0FBQyxFQUpXLGFBQWEsR0FBYixxQkFBYSxLQUFiLHFCQUFhLFFBSXhCO0FBRUQsSUFBWSxhQUVYO0FBRkQsV0FBWSxhQUFhO0lBQ3ZCLG1FQUFrRCxDQUFBO0FBQ3BELENBQUMsRUFGVyxhQUFhLEdBQWIscUJBQWEsS0FBYixxQkFBYSxRQUV4QjtBQUVELElBQVksZ0JBSVg7QUFKRCxXQUFZLGdCQUFnQjtJQUMxQiwrREFBMkMsQ0FBQTtJQUMzQyxvQkFBb0I7SUFDcEIsb0RBQWdDLENBQUE7QUFDbEMsQ0FBQyxFQUpXLGdCQUFnQixHQUFoQix3QkFBZ0IsS0FBaEIsd0JBQWdCLFFBSTNCO0FBRUQsd0NBQXdDO0FBRXhDLElBQVksU0FHWDtBQUhELFdBQVksU0FBUztJQUNuQix1Q0FBMEIsQ0FBQTtJQUMxQixvQkFBb0I7QUFDdEIsQ0FB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