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
exports.connectClient = connectClient;
exports.connectNetwork = connectNetwork;
exports.connectWallet = connectWallet;
exports.connect = connect;
exports.deriveMnemomicFromEthereumSignature = deriveMnemomicFromEthereumSignature;
exports.getHeight = getHeight;
exports.getFeeTiers = getFeeTiers;
exports.getUserFeeTier = getUserFeeTier;
exports.getEquityTiers = getEquityTiers;
exports.getPerpetualMarkets = getPerpetualMarkets;
exports.placeOrder = placeOrder;
exports.wrappedError = wrappedError;
exports.cancelOrder = cancelOrder;
exports.deposit = deposit;
exports.withdraw = withdraw;
exports.faucet = faucet;
exports.withdrawToIBC = withdrawToIBC;
exports.transferNativeToken = transferNativeToken;
exports.getAccountBalance = getAccountBalance;
exports.getAccountBalances = getAccountBalances;
exports.getUserStats = getUserStats;
exports.simulateDeposit = simulateDeposit;
exports.simulateWithdraw = simulateWithdraw;
exports.simulateTransferNativeToken = simulateTransferNativeToken;
exports.signRawPlaceOrder = signRawPlaceOrder;
exports.signPlaceOrder = signPlaceOrder;
exports.signCancelOrder = signCancelOrder;
exports.encodeAccountRequestData = encodeAccountRequestData;
exports.decodeAccountResponseValue = decodeAccountResponseValue;
exports.getOptimalNode = getOptimalNode;
exports.getOptimalIndexer = getOptimalIndexer;
exports.getRewardsParams = getRewardsParams;
exports.getDelegatorDelegations = getDelegatorDelegations;
exports.getDelegatorUnbondingDelegations = getDelegatorUnbondingDelegations;
exports.getMarketPrice = getMarketPrice;
exports.getNobleBalance = getNobleBalance;
exports.sendNobleIBC = sendNobleIBC;
exports.withdrawToNobleIBC = withdrawToNobleIBC;
exports.cctpWithdraw = cctpWithdraw;
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
async function connect(network, mnemonic) {
    try {
        await connectClient(network);
        return connectWallet(mnemonic);
    }
    catch (e) {
        return wrappedError(e);
    }
}
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
async function placeOrder(payload) {
    var _a, _b, _c;
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
        const routerFeePpm = (_c = json.routerFeePpm) !== null && _c !== void 0 ? _c : 0;
        const routerFeeSubaccountOwner = json.routerFeeSubaccountOwner;
        const routerFeeSubaccountNumber = json.routerFeeSubaccountNumber;
        const subaccount = new subaccount_1.SubaccountInfo(wallet, subaccountNumber);
        const tx = await client.placeOrder(subaccount, marketId, type, side, price, size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, triggerPrice, marketInfo, currentHeight, routerFeePpm, routerFeeSubaccountOwner, routerFeeSubaccountNumber);
        return (0, helpers_1.encodeJson)(tx);
    }
    catch (error) {
        return wrappedError(error);
    }
}
function wrappedError(error) {
    const text = JSON.stringify(error, Object.getOwnPropertyNames(error));
    return `{"error": ${text}}`;
}
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
            typeUrl: json.msgTypeUrl, // '/ibc.applications.transfer.v1.MsgTransfer',
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
async function signRawPlaceOrder(subaccountNumber, clientId, clobPairId, side, quantums, subticks, timeInForce, orderFlags, reduceOnly, goodTilBlock, goodTilBlockTime, clientMetadata, routerFeePpm = 0, routerFeeSubaccountOwner = '', routerFeeSubaccountNumber = 0) {
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
            const msg = client.validatorClient.post.composer.composeMsgPlaceOrder(wallet.address, subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock, goodTilBlockTime, side, quantums, subticks, timeInForce, reduceOnly, clientMetadata !== null && clientMetadata !== void 0 ? clientMetadata : 0, undefined, undefined, routerFeePpm, routerFeeSubaccountOwner, routerFeeSubaccountNumber);
            resolve([msg]);
        });
        const signed = await client.sign(wallet, () => msgs, true);
        return Buffer.from(signed).toString('base64');
    }
    catch (error) {
        return wrappedError(error);
    }
}
async function signPlaceOrder(subaccountNumber, marketId, type, side, price, 
// trigger_price: number,   // not used for MARKET and LIMIT
size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, routerFeePpm = 0, routerFeeSubaccountOwner = '', routerFeeSubaccountNumber = 0) {
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
        const signed = await client.signPlaceOrder(subaccount, marketId, type, side, price, size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, routerFeePpm, routerFeeSubaccountOwner, routerFeeSubaccountNumber);
        return signed;
    }
    catch (error) {
        return wrappedError(error);
    }
}
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
async function sendNobleIBC(squidPayload) {
    try {
        const client = globalThis.nobleClient;
        if (client === undefined || !client.isConnected) {
            throw new errors_1.UserError('client is not connected.');
        }
        const json = JSON.parse(squidPayload);
        const ibcMsg = {
            typeUrl: json.msgTypeUrl, // '/ibc.applications.transfer.v1.MsgTransfer',
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
async function cctpWithdraw(squidPayload) {
    try {
        const client = globalThis.nobleClient;
        if (client === undefined || !client.isConnected) {
            throw new errors_1.UserError('client is not connected.');
        }
        const json = JSON.parse(squidPayload);
        const ibcMsg = {
            typeUrl: json.typeUrl, // '/circle.cctp.v1.MsgDepositForBurn',
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibmF0aXZlLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NsaWVudHMvbmF0aXZlLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7QUFBQTs7RUFFRTs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQXFDRixzQ0FTQztBQUVELHdDQXdEQztBQUVELHNDQWdCQztBQUVELDBCQVVDO0FBRUQsa0ZBV0M7QUFFRCw4QkFXQztBQUVELGtDQVdDO0FBRUQsd0NBV0M7QUFFRCx3Q0FZQztBQUVELGtEQVdDO0FBRUQsZ0NBa0ZDO0FBRUQsb0NBR0M7QUFFRCxrQ0E4Q0M7QUFFRCwwQkFnQ0M7QUFFRCw0QkFpQ0M7QUFFRCx3QkE2QkM7QUFFRCxzQ0E0Q0M7QUFFRCxrREFzQ0M7QUFFRCw4Q0FrQkM7QUFFRCxnREFpQkM7QUFFRCxvQ0FtQkM7QUFFRCwwQ0F3Q0M7QUFFRCw0Q0F5Q0M7QUFFRCxrRUF3Q0M7QUFFRCw4Q0EyREM7QUFFRCx3Q0FrREM7QUFFRCwwQ0ErQkM7QUFFRCw0REFXQztBQUVELGdFQWVDO0FBRUQsd0NBd0JDO0FBRUQsOENBc0JDO0FBRUQsNENBV0M7QUFFRCwwREFtQkM7QUFFRCw0RUFtQkM7QUFFRCx3Q0FrQkM7QUFFRCwwQ0FhQztBQUVELG9DQStCQztBQUVELGdEQXdDQztBQUVELG9DQWlDQztBQTFuQ0QsK0NBQWtEO0FBRWxELG1GQUFxRTtBQUdyRSxnREFBc0Y7QUFDdEYsMENBQTBDO0FBQzFDLDRDQUErRDtBQUMvRCxrREFBcUU7QUFDckUsNERBQXdEO0FBQ3hELHlEQUFpRTtBQUNqRSwyQ0FFcUI7QUFDckIsbURBQStDO0FBQy9DLDBFQUFpRDtBQUNqRCxpREFBNkM7QUFDN0MsNkNBQThDO0FBaUJ2QyxLQUFLLFVBQVUsYUFBYSxDQUNqQyxPQUFnQjtJQUVoQixJQUFJLENBQUM7UUFDSCxVQUFVLENBQUMsTUFBTSxHQUFHLE1BQU0sa0NBQWUsQ0FBQyxPQUFPLENBQUMsT0FBTyxDQUFDLENBQUM7UUFD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