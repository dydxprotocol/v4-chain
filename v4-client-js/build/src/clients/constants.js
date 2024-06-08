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
            USDC_DENOM: 'ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5',
            USDC_GAS_DENOM: 'uusdc',
            USDC_DECIMALS: 6,
            CHAINTOKEN_DECIMALS: 18,
        });
        return new Network('testnet', indexerConfig, validatorConfig);
    }
    static local() {
        const indexerConfig = new IndexerConfig(IndexerApiHost.LOCAL, IndexerWSHost.LOCAL);
        const validatorConfig = new ValidatorConfig(ValidatorApiHost.LOCAL, exports.LOCAL_CHAIN_ID, {
            CHAINTOKEN_DENOM: 'adv4tnt',
            USDC_DENOM: 'ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5',
            USDC_GAS_DENOM: 'uusdc',
            USDC_DECIMALS: 6,
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29uc3RhbnRzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NsaWVudHMvY29uc3RhbnRzLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQ0EsZ0RBQXdCO0FBSXhCLG1EQUFpQztBQUVqQyxXQUFXO0FBQ0UsUUFBQSxZQUFZLEdBQUcsc0JBQXNCLENBQUM7QUFDdEMsUUFBQSxnQkFBZ0IsR0FBRyxzQkFBc0IsQ0FBQztBQUMxQyxRQUFBLGdCQUFnQixHQUFHLGdCQUFnQixDQUFDO0FBQ2pELG9EQUFvRDtBQUN2QyxRQUFBLGNBQWMsR0FBRyxPQUFPLENBQUM7QUFFdEMscUNBQXFDO0FBQ3JDLElBQVksY0FJWDtBQUpELFdBQVksY0FBYztJQUN4QiwrREFBNkMsQ0FBQTtJQUM3QyxpREFBK0IsQ0FBQTtJQUMvQixvQkFBb0I7QUFDdEIsQ0FBQyxFQUpXLGNBQWMsR0FBZCxzQkFBYyxLQUFkLHNCQUFjLFFBSXpCO0FBRUQsSUFBWSxhQUlYO0FBSkQsV0FBWSxhQUFhO0lBQ3ZCLGtFQUFpRCxDQUFBO0lBQ2pELG9CQUFvQjtJQUNwQiw4Q0FBNkIsQ0FBQTtBQUMvQixDQUFDLEVBSlcsYUFBYSxHQUFiLHFCQUFhLEtBQWIscUJBQWEsUUFJeEI7QUFFRCxJQUFZLGFBRVg7QUFGRCxXQUFZLGFBQWE7SUFDdkIsbUVBQWtELENBQUE7QUFDcEQsQ0FBQyxFQUZXLGFBQWEsR0FBYixxQkFBYSxLQUFiLHFCQUFhLFFBRXhCO0FBRUQsSUFBWSxnQkFJWDtBQUpELFdBQVksZ0JBQWdCO0lBQzFCLCtEQUEyQyxDQUFBO0lBQzNDLG9CQUFvQjtJQUNwQixvREFBZ0MsQ0FBQTtBQUNsQyxDQUFDLEVBSlcsZ0JBQWdCLEdBQWhCLHdCQUFnQixLQUFoQix3QkFBZ0IsUUFJM0I7QUFFRCx3Q0FBd0M7QUFFeEMsSUFBWSxTQUdYO0FBSEQsV0FBWSxTQUFTO0lBQ25CLHVDQUEwQixDQUFBO0lBQzFCLG9CQUFv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