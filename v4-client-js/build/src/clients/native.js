"use strict";
/*
    Native app can call JS functions with primitives.
*/
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
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.cctpWithdraw = exports.withdrawToNobleIBC = exports.sendNobleIBC = exports.getNobleBalance = exports.getMarketPrice = exports.getDelegatorUnbondingDelegations = exports.getDelegatorDelegations = exports.getRewardsParams = exports.getOptimalIndexer = exports.getOptimalNode = exports.decodeAccountResponseValue = exports.encodeAccountRequestData = exports.signCancelOrder = exports.signPlaceOrder = exports.signRawPlaceOrder = exports.simulateTransferNativeToken = exports.simulateWithdraw = exports.simulateDeposit = exports.getUserStats = exports.getAccountBalances = exports.getAccountBalance = exports.transferNativeToken = exports.withdrawToIBC = exports.faucet = exports.withdraw = exports.deposit = exports.cancelOrder = exports.wrappedError = exports.placeOrder = exports.getPerpetualMarkets = exports.getEquityTiers = exports.getUserFeeTier = exports.getFeeTiers = exports.getHeight = exports.deriveMnemomicFromEthereumSignature = exports.connect = exports.connectWallet = exports.connectNetwork = exports.connectClient = void 0;
const stargate_1 = require("@cosmjs/stargate");
const AuthModule = __importStar(require("cosmjs-types/cosmos/auth/v1beta1/query"));
const constants_1 = require("../lib/constants");
const errors_1 = require("../lib/errors");
const helpers_1 = require("../lib/helpers");
const onboarding_1 = require("../lib/onboarding");
const network_optimizer_1 = require("../network_optimizer");
const composite_client_1 = require("./composite-client");
const constants_2 = require("./constants");
const faucet_client_1 = require("./faucet-client");
const local_wallet_1 = __importDefault(require("./modules/local-wallet"));
const noble_client_1 = require("./noble-client");
const subaccount_1 = require("./subaccount");
async function connectClient(network) {
    try {
        globalThis.client = await composite_client_1.CompositeClient.connect(network);
        return (0, helpers_1.encodeJson)(network);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.connectClient = connectClient;
async function connectNetwork(paramsJSON) {
    try {
        const params = JSON.parse(paramsJSON);
        const { indexerUrl, websocketUrl, validatorUrl, chainId, faucetUrl, nobleValidatorUrl, TDAI_DENOM, TDAI_DECIMALS, TDAI_GAS_DENOM, CHAINTOKEN_DENOM, CHAINTOKEN_DECIMALS, CHAINTOKEN_GAS_DENOM, } = params;
        if (indexerUrl === undefined ||
            websocketUrl === undefined ||
            validatorUrl === undefined ||
            chainId === undefined) {
            throw new errors_1.UserError('Missing required network params');
        }
        if (TDAI_DENOM === undefined ||
            TDAI_DECIMALS === undefined ||
            CHAINTOKEN_DENOM === undefined ||
            CHAINTOKEN_DECIMALS === undefined) {
            throw new errors_1.UserError('Missing required token params');
        }
        const indexerConfig = new constants_2.IndexerConfig(indexerUrl, websocketUrl);
        const validatorConfig = new constants_2.ValidatorConfig(validatorUrl, chainId, {
            TDAI_DENOM,
            TDAI_DECIMALS,
            TDAI_GAS_DENOM,
            CHAINTOKEN_DENOM,
            CHAINTOKEN_DECIMALS,
            CHAINTOKEN_GAS_DENOM,
        });
        const config = new constants_2.Network('native', indexerConfig, validatorConfig);
        globalThis.client = await composite_client_1.CompositeClient.connect(config);
        if (faucetUrl !== undefined) {
            globalThis.faucetClient = new faucet_client_1.FaucetClient(faucetUrl);
        }
        else {
            globalThis.faucetClient = null;
        }
        globalThis.nobleClient = new noble_client_1.NobleClient(nobleValidatorUrl);
        if (globalThis.nobleWallet)
            await globalThis.nobleClient.connect(globalThis.nobleWallet);
        return (0, helpers_1.encodeJson)(config);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.connectNetwork = connectNetwork;
async function connectWallet(mnemonic) {
    var _a;
    try {
        globalThis.wallet = await local_wallet_1.default.fromMnemonic(mnemonic, constants_1.BECH32_PREFIX);
        globalThis.nobleWallet = await local_wallet_1.default.fromMnemonic(mnemonic, constants_1.NOBLE_BECH32_PREFIX);
        await ((_a = globalThis.nobleClient) === null || _a === void 0 ? void 0 : _a.connect(globalThis.nobleWallet));
        const address = globalThis.wallet.address;
        return (0, helpers_1.encodeJson)({ address });
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.connectWallet = connectWallet;
async function connect(network, mnemonic) {
    try {
        await connectClient(network);
        return connectWallet(mnemonic);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.connect = connect;
async function deriveMnemomicFromEthereumSignature(signature) {
    try {
        const { mnemonic } = (0, onboarding_1.deriveHDKeyFromEthereumSignature)(signature);
        const wallet = await local_wallet_1.default.fromMnemonic(mnemonic, constants_1.BECH32_PREFIX);
        const result = { mnemonic, address: wallet.address };
        return new Promise((resolve) => {
            resolve((0, helpers_1.encodeJson)(result));
        });
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.deriveMnemomicFromEthereumSignature = deriveMnemomicFromEthereumSignature;
async function getHeight() {
    var _a;
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const block = await ((_a = globalThis.client) === null || _a === void 0 ? void 0 : _a.validatorClient.get.latestBlock());
        return (0, helpers_1.encodeJson)(block);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.getHeight = getHeight;
async function getFeeTiers() {
    var _a;
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const feeTiers = await ((_a = globalThis.client) === null || _a === void 0 ? void 0 : _a.validatorClient.get.getFeeTiers());
        return (0, helpers_1.encodeJson)(feeTiers);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.getFeeTiers = getFeeTiers;
async function getUserFeeTier(address) {
    var _a;
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const feeTiers = await ((_a = globalThis.client) === null || _a === void 0 ? void 0 : _a.validatorClient.get.getUserFeeTier(address));
        return (0, helpers_1.encodeJson)(feeTiers);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.getUserFeeTier = getUserFeeTier;
async function getEquityTiers() {
    var _a;
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const equityTiers = await ((_a = globalThis.client) === null || _a === void 0 ? void 0 : _a.validatorClient.get.getEquityTierLimitConfiguration());
        return (0, helpers_1.encodeJson)(equityTiers, helpers_1.ByteArrayEncoding.BIGINT);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.getEquityTiers = getEquityTiers;
async function getPerpetualMarkets() {
    var _a;
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const markets = await ((_a = globalThis.client) === null || _a === void 0 ? void 0 : _a.indexerClient.markets.getPerpetualMarkets());
        return (0, helpers_1.encodeJson)(markets);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.getPerpetualMarkets = getPerpetualMarkets;
async function placeOrder(payload) {
    var _a, _b;
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const json = JSON.parse(payload);
        const subaccountNumber = json.subaccountNumber;
        if (subaccountNumber === undefined) {
            throw new errors_1.UserError('subaccountNumber is not set');
        }
        const marketId = json.marketId;
        if (marketId === undefined) {
            throw new errors_1.UserError('marketId is not set');
        }
        const type = json.type;
        if (type === undefined) {
            throw new errors_1.UserError('type is not set');
        }
        const side = json.side;
        if (side === undefined) {
            throw new errors_1.UserError('side is not set');
        }
        const price = json.price;
        if (price === undefined) {
            throw new errors_1.UserError('price is not set');
        }
        // trigger_price: number,   // not used for MARKET and LIMIT
        const size = json.size;
        if (size === undefined) {
            throw new errors_1.UserError('size is not set');
        }
        const clientId = json.clientId;
        if (clientId === undefined) {
            throw new errors_1.UserError('clientId is not set');
        }
        const timeInForce = json.timeInForce;
        const goodTilTimeInSeconds = json.goodTilTimeInSeconds;
        const execution = json.execution;
        const postOnly = (_a = json.postOnly) !== null && _a !== void 0 ? _a : false;
        const reduceOnly = (_b = json.reduceOnly) !== null && _b !== void 0 ? _b : false;
        const triggerPrice = json.triggerPrice;
        const marketInfo = json.marketInfo;
        const currentHeight = json.currentHeight;
        const subaccount = new subaccount_1.SubaccountInfo(wallet, subaccountNumber);
        const tx = await client.placeOrder(subaccount, marketId, type, side, price, size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, triggerPrice, marketInfo, currentHeight);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.placeOrder = placeOrder;
function wrappedError(error) {
    const text = JSON.stringify(error, Object.getOwnPropertyNames(error));
    return `{"error": ${text}}`;
}
exports.wrappedError = wrappedError;
async function cancelOrder(payload) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectNetwork() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const json = JSON.parse(payload);
        const subaccountNumber = json.subaccountNumber;
        if (subaccountNumber === undefined) {
            throw new errors_1.UserError('subaccountNumber is not set');
        }
        const clientId = json.clientId;
        if (clientId === undefined) {
            throw new errors_1.UserError('clientId is not set');
        }
        const orderFlags = json.orderFlags;
        if (orderFlags === undefined) {
            throw new errors_1.UserError('orderFlags is not set');
        }
        const clobPairId = json.clobPairId;
        if (clobPairId === undefined) {
            throw new errors_1.UserError('clobPairId is not set');
        }
        const goodTilBlock = json.goodTilBlock;
        const goodTilBlockTime = json.goodTilBlockTime;
        const subaccount = new subaccount_1.SubaccountInfo(wallet, subaccountNumber);
        const tx = await client.cancelRawOrder(subaccount, clientId, orderFlags, clobPairId, goodTilBlock !== 0 ? goodTilBlock : undefined, goodTilBlockTime !== 0 ? goodTilBlockTime : undefined);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.cancelOrder = cancelOrder;
async function deposit(payload) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectNetwork() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const json = JSON.parse(payload);
        const subaccountNumber = json.subaccountNumber;
        if (subaccountNumber === undefined) {
            throw new errors_1.UserError('subaccountNumber is not set');
        }
        const amount = json.amount;
        if (amount === undefined) {
            throw new errors_1.UserError('amount is not set');
        }
        const subaccount = new subaccount_1.SubaccountInfo(wallet, subaccountNumber);
        const tx = await client.depositToSubaccount(subaccount, amount);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.deposit = deposit;
async function withdraw(payload) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectNetwork() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const json = JSON.parse(payload);
        const subaccountNumber = json.subaccountNumber;
        if (subaccountNumber === undefined) {
            throw new errors_1.UserError('subaccountNumber is not set');
        }
        const amount = json.amount;
        if (amount === undefined) {
            throw new errors_1.UserError('amount is not set');
        }
        const subaccount = new subaccount_1.SubaccountInfo(wallet, subaccountNumber);
        const tx = await client.withdrawFromSubaccount(subaccount, amount, json.recipient);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.withdraw = withdraw;
async function faucet(payload) {
    try {
        const faucetClient = globalThis.faucetClient;
        if (!faucetClient) {
            throw new errors_1.UserError('faucetClient is not connected. Call connectNetwork() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const json = JSON.parse(payload);
        const subaccountNumber = json.subaccountNumber;
        if (subaccountNumber === undefined) {
            throw new errors_1.UserError('subaccountNumber is not set');
        }
        const amount = json.amount;
        if (amount === undefined) {
            throw new errors_1.UserError('amount is not set');
        }
        const response = await faucetClient.fill(wallet.address, subaccountNumber, amount);
        return (0, helpers_1.encodeJson)(response);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.faucet = faucet;
async function withdrawToIBC(subaccountNumber, amount, payload) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const decode = (str) => Buffer.from(str, 'base64').toString('binary');
        const decoded = decode(payload);
        const json = JSON.parse(decoded);
        const ibcMsg = {
            typeUrl: json.msgTypeUrl,
            value: json.msg,
        };
        const subaccount = new subaccount_1.SubaccountInfo(wallet, subaccountNumber);
        const subaccountMsg = client.withdrawFromSubaccountMessage(subaccount, amount);
        const msgs = [subaccountMsg, ibcMsg];
        const encodeObjects = new Promise((resolve) => resolve(msgs));
        const tx = await client.send(wallet, () => {
            return encodeObjects;
        }, false, undefined, undefined);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.withdrawToIBC = withdrawToIBC;
async function transferNativeToken(payload) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const json = JSON.parse(payload);
        const amount = json.amount;
        if (amount === undefined) {
            throw new errors_1.UserError('amount is not set');
        }
        const msg = client.sendTokenMessage(wallet, amount, json.recipient);
        const msgs = [msg];
        const encodeObjects = new Promise((resolve) => resolve(msgs));
        const tx = await client.send(wallet, () => {
            return encodeObjects;
        }, false);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.transferNativeToken = transferNativeToken;
async function getAccountBalance() {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const address = globalThis.wallet.address;
        const tx = await client.validatorClient.get
            .getAccountBalance(address, client.validatorClient.config.denoms.TDAI_DENOM);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.getAccountBalance = getAccountBalance;
async function getAccountBalances() {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const address = globalThis.wallet.address;
        const tx = await client.validatorClient.get.getAccountBalances(address);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.getAccountBalances = getAccountBalances;
async function getUserStats(payload) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const json = JSON.parse(payload);
        const address = json.address;
        if (address === undefined) {
            throw new errors_1.UserError('address is not set');
        }
        const tx = await client.validatorClient.get.getUserStats(address);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.getUserStats = getUserStats;
async function simulateDeposit(payload) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const json = JSON.parse(payload);
        const subaccountNumber = json.subaccountNumber;
        if (subaccountNumber === undefined) {
            throw new errors_1.UserError('subaccountNumber is not set');
        }
        const amount = json.amount;
        if (amount === undefined) {
            throw new errors_1.UserError('amount is not set');
        }
        const subaccount = new subaccount_1.SubaccountInfo(wallet, subaccountNumber);
        const msg = client.depositToSubaccountMessage(subaccount, amount);
        const msgs = [msg];
        const encodeObjects = new Promise((resolve) => resolve(msgs));
        const stdFee = await client.simulate(globalThis.wallet, () => {
            return encodeObjects;
        });
        return JSON.stringify(stdFee);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.simulateDeposit = simulateDeposit;
async function simulateWithdraw(payload) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const json = JSON.parse(payload);
        const subaccountNumber = json.subaccountNumber;
        if (subaccountNumber === undefined) {
            throw new errors_1.UserError('subaccountNumber is not set');
        }
        const amount = json.amount;
        if (amount === undefined) {
            throw new errors_1.UserError('amount is not set');
        }
        const subaccount = new subaccount_1.SubaccountInfo(wallet, subaccountNumber);
        const msg = client.withdrawFromSubaccountMessage(subaccount, amount, json.recipient);
        const msgs = [msg];
        const encodeObjects = new Promise((resolve) => resolve(msgs));
        const stdFee = await client.simulate(globalThis.wallet, () => {
            return encodeObjects;
        });
        return (0, helpers_1.encodeJson)(stdFee);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.simulateWithdraw = simulateWithdraw;
async function simulateTransferNativeToken(payload) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const json = JSON.parse(payload);
        const recipient = json.recipient;
        if (recipient === undefined) {
            throw new errors_1.UserError('recipient is not set');
        }
        const amount = json.amount;
        if (amount === undefined) {
            throw new errors_1.UserError('amount is not set');
        }
        const msg = client.sendTokenMessage(wallet, amount, json.recipient);
        const msgs = [msg];
        const encodeObjects = new Promise((resolve) => resolve(msgs));
        const stdFee = await client.simulate(globalThis.wallet, () => {
            return encodeObjects;
        });
        return (0, helpers_1.encodeJson)(stdFee);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.simulateTransferNativeToken = simulateTransferNativeToken;
async function signRawPlaceOrder(subaccountNumber, clientId, clobPairId, side, quantums, subticks, timeInForce, orderFlags, reduceOnly, goodTilBlock, goodTilBlockTime, clientMetadata) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const msgs = new Promise((resolve) => {
            const msg = client.validatorClient.post.composer.composeMsgPlaceOrder(wallet.address, subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock, goodTilBlockTime, side, quantums, subticks, timeInForce, reduceOnly, clientMetadata !== null && clientMetadata !== void 0 ? clientMetadata : 0);
            resolve([msg]);
        });
        const signed = await client.sign(wallet, () => msgs, true);
        return Buffer.from(signed).toString('base64');
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.signRawPlaceOrder = signRawPlaceOrder;
async function signPlaceOrder(subaccountNumber, marketId, type, side, price, 
// trigger_price: number,   // not used for MARKET and LIMIT
size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const subaccount = new subaccount_1.SubaccountInfo(wallet, subaccountNumber);
        const signed = await client.signPlaceOrder(subaccount, marketId, type, side, price, size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly);
        return signed;
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.signPlaceOrder = signPlaceOrder;
async function signCancelOrder(subaccountNumber, clientId, orderFlags, clobPairId, goodTilBlock, goodTilBlockTime) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const subaccount = new subaccount_1.SubaccountInfo(wallet, subaccountNumber);
        const signed = await client.signCancelOrder(subaccount, clientId, orderFlags, clobPairId, goodTilBlock, goodTilBlockTime);
        return signed;
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.signCancelOrder = signCancelOrder;
async function encodeAccountRequestData(address) {
    return new Promise((resolve, reject) => {
        try {
            const requestData = Uint8Array.from(AuthModule.QueryAccountRequest.encode({ address }).finish());
            resolve(Buffer.from(requestData).toString('hex'));
        }
        catch (error) {
            reject(error);
        }
    });
}
exports.encodeAccountRequestData = encodeAccountRequestData;
async function decodeAccountResponseValue(value) {
    return new Promise((resolve, reject) => {
        try {
            const rawData = Buffer.from(value, 'base64');
            const rawAccount = AuthModule.QueryAccountResponse.decode(rawData).account;
            // The promise should have been rejected if the rawAccount was undefined.
            if (rawAccount === undefined) {
                throw Error('rawAccount is undefined');
            }
            const account = (0, stargate_1.accountFromAny)(rawAccount);
            resolve((0, helpers_1.encodeJson)(account));
        }
        catch (error) {
            reject(error);
        }
    });
}
exports.decodeAccountResponseValue = decodeAccountResponseValue;
async function getOptimalNode(endpointUrlsAsJson) {
    /*
      param:
        endpointUrlsAsJson:
        {
          "endpointUrls": [
            "https://rpc.testnet.near.org"
          ],
          "chainId": "testnet"
        }
    */
    try {
        const param = JSON.parse(endpointUrlsAsJson);
        const endpointUrls = param.endpointUrls;
        const chainId = param.chainId;
        const networkOptimizer = new network_optimizer_1.NetworkOptimizer();
        const optimalUrl = await networkOptimizer.findOptimalNode(endpointUrls, chainId);
        const url = {
            url: optimalUrl,
        };
        return (0, helpers_1.encodeJson)(url);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.getOptimalNode = getOptimalNode;
async function getOptimalIndexer(endpointUrlsAsJson) {
    /*
      param:
        endpointUrlsAsJson:
        {
          "endpointUrls": [
            "https://api.example.org"
          ]
        }
    */
    try {
        const param = JSON.parse(endpointUrlsAsJson);
        const endpointUrls = param.endpointUrls;
        const networkOptimizer = new network_optimizer_1.NetworkOptimizer();
        const optimalUrl = await networkOptimizer.findOptimalIndexer(endpointUrls);
        const url = {
            url: optimalUrl,
        };
        return (0, helpers_1.encodeJson)(url);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.getOptimalIndexer = getOptimalIndexer;
async function getRewardsParams() {
    var _a;
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const rewardsParams = await ((_a = globalThis.client) === null || _a === void 0 ? void 0 : _a.validatorClient.get.getRewardsParams());
        return (0, helpers_1.encodeJson)(rewardsParams);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.getRewardsParams = getRewardsParams;
async function getDelegatorDelegations(payload) {
    var _a;
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const json = JSON.parse(payload);
        const address = json.address;
        if (address === undefined) {
            throw new errors_1.UserError('address is not set');
        }
        const delegations = await ((_a = globalThis
            .client) === null || _a === void 0 ? void 0 : _a.validatorClient.get.getDelegatorDelegations(address));
        return (0, helpers_1.encodeJson)(delegations);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.getDelegatorDelegations = getDelegatorDelegations;
async function getDelegatorUnbondingDelegations(payload) {
    var _a;
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const json = JSON.parse(payload);
        const address = json.address;
        if (address === undefined) {
            throw new errors_1.UserError('address is not set');
        }
        const delegations = await ((_a = globalThis
            .client) === null || _a === void 0 ? void 0 : _a.validatorClient.get.getDelegatorUnbondingDelegations(address));
        return (0, helpers_1.encodeJson)(delegations);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.getDelegatorUnbondingDelegations = getDelegatorUnbondingDelegations;
async function getMarketPrice(payload) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const json = JSON.parse(payload);
        const marketId = json.marketId;
        if (marketId === undefined) {
            throw new errors_1.UserError('marketId is not set');
        }
        const marketPrice = await client.validatorClient.get.getPrice(marketId);
        return (0, helpers_1.encodeJson)(marketPrice);
    }
    catch (e) {
        return wrappedError(e);
    }
}
exports.getMarketPrice = getMarketPrice;
async function getNobleBalance() {
    try {
        const client = globalThis.nobleClient;
        if (client === undefined || !client.isConnected) {
            throw new errors_1.UserError('client is not connected.');
        }
        const coin = await client.getAccountBalance('utdai');
        return (0, helpers_1.encodeJson)(coin);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.getNobleBalance = getNobleBalance;
async function sendNobleIBC(squidPayload) {
    try {
        const client = globalThis.nobleClient;
        if (client === undefined || !client.isConnected) {
            throw new errors_1.UserError('client is not connected.');
        }
        const json = JSON.parse(squidPayload);
        const ibcMsg = {
            typeUrl: json.msgTypeUrl,
            value: json.msg,
        };
        const fee = await client.simulateTransaction([ibcMsg]);
        // take out fee from amount before sweeping
        const amount = parseInt(ibcMsg.value.token.amount, 10) -
            Math.floor(parseInt(fee.amount[0].amount, 10) * constants_1.GAS_MULTIPLIER);
        if (amount <= 0) {
            throw new errors_1.UserError('noble balance does not cover fees');
        }
        ibcMsg.value.token.amount = amount.toString();
        const tx = await client.send([ibcMsg]);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.sendNobleIBC = sendNobleIBC;
async function withdrawToNobleIBC(payload) {
    try {
        const client = globalThis.client;
        if (client === undefined) {
            throw new errors_1.UserError('client is not connected. Call connectClient() first');
        }
        const wallet = globalThis.wallet;
        if (wallet === undefined) {
            throw new errors_1.UserError('wallet is not set. Call connectWallet() first');
        }
        const json = JSON.parse(payload);
        const { subaccountNumber, amount, ibcPayload } = json !== null && json !== void 0 ? json : {};
        const decode = (str) => Buffer.from(str, 'base64').toString('binary');
        const decoded = decode(ibcPayload);
        const parsedIbcPayload = JSON.parse(decoded);
        const msg = client.withdrawFromSubaccountMessage(new subaccount_1.SubaccountInfo(wallet, subaccountNumber), parseFloat(amount).toFixed(client.validatorClient.config.denoms.TDAI_DECIMALS));
        const ibcMsg = {
            typeUrl: parsedIbcPayload.msgTypeUrl,
            value: parsedIbcPayload.msg,
        };
        const tx = await client.send(wallet, () => Promise.resolve([msg, ibcMsg]), false);
        return (0, helpers_1.encodeJson)({
            txHash: `0x${Buffer.from(tx === null || tx === void 0 ? void 0 : tx.hash).toString('hex')}`,
        });
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.withdrawToNobleIBC = withdrawToNobleIBC;
async function cctpWithdraw(squidPayload) {
    try {
        const client = globalThis.nobleClient;
        if (client === undefined || !client.isConnected) {
            throw new errors_1.UserError('client is not connected.');
        }
        const json = JSON.parse(squidPayload);
        const ibcMsg = {
            typeUrl: json.typeUrl,
            value: json.value,
        };
        const fee = await client.simulateTransaction([ibcMsg]);
        // take out fee from amount before sweeping
        const amount = parseInt(ibcMsg.value.amount, 10) -
            Math.floor(parseInt(fee.amount[0].amount, 10) * constants_1.GAS_MULTIPLIER);
        if (amount <= 0) {
            throw new Error('noble balance does not cover fees');
        }
        ibcMsg.value.amount = amount.toString();
        const tx = await client.send([ibcMsg]);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
exports.cctpWithdraw = cctpWithdraw;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibmF0aXZlLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NsaWVudHMvbmF0aXZlLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7QUFBQTs7RUFFRTs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFHRiwrQ0FBa0Q7QUFFbEQsbUZBQXFFO0FBR3JFLGdEQUFzRjtBQUN0RiwwQ0FBMEM7QUFDMUMsNENBQStEO0FBQy9ELGtEQUFxRTtBQUNyRSw0REFBd0Q7QUFDeEQseURBQWlFO0FBQ2pFLDJDQUVxQjtBQUNyQixtREFBK0M7QUFDL0MsMEVBQWlEO0FBQ2pELGlEQUE2QztBQUM3Qyw2Q0FBOEM7QUFpQnZDLEtBQUssVUFBVSxhQUFhLENBQ2pDLE9BQWdCO0lBRWhCLElBQUk7UUFDRixVQUFVLENBQUMsTUFBTSxHQUFHLE1BQU0sa0NBQWUsQ0FBQyxPQUFPLENBQUMsT0FBTyxDQUFDLENBQUM7UUFDM0QsT0FBTyxJQUFBLG9CQUFVLEVBQUMsT0FBTyxDQUFDLENBQUM7S0FDNUI7SUFBQyxPQUFPLENBQUMsRUFBRTtRQUNWLE9BQU8sWUFBWSxDQUFDLENBQUMsQ0FBQyxDQUFDO0tBQ3hCO0FBQ0gsQ0FBQztBQVRELHNDQVNDO0FBRU0sS0FBSyxVQUFVLGNBQWMsQ0FDbEMsVUFBa0I7SUFFbEIsSUFBSTtRQUNGLE1BQU0sTUFBTSxHQUFHLElBQUksQ0FBQyxLQUFLLENBQUMsVUFBVSxDQUFDLENBQUM7UUFDdEMsTUFBTSxFQUNKLFVBQVUsRUFDVixZQUFZLEVBQ1osWUFBWSxFQUNaLE9BQU8sRUFDUCxTQUFTLEVBQ1QsaUJBQWlCLEVBQ2pCLFVBQVUsRUFDVixhQUFhLEVBQ2IsY0FBYyxFQUNkLGdCQUFnQixFQUNoQixtQkFBbUIsRUFDbkIsb0JBQW9CLEdBQ3JCLEdBQUcsTUFBTSxDQUFDO1FBRVgsSUFBSSxVQUFVLEtBQUssU0FBUztZQUMxQixZQUFZLEtBQUssU0FBUztZQUMxQixZQUFZLEtBQUssU0FBUztZQUMxQixP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