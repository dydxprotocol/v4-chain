"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Post = void 0;
const math_1 = require("@cosmjs/math");
const stargate_1 = require("@cosmjs/stargate");
const tendermint_rpc_1 = require("@cosmjs/tendermint-rpc");
const lodash_1 = __importDefault(require("lodash"));
const long_1 = __importDefault(require("long"));
const protobufjs_1 = __importDefault(require("protobufjs"));
const constants_1 = require("../constants");
const errors_1 = require("../lib/errors");
const registry_1 = require("../lib/registry");
const types_1 = require("../types");
const composer_1 = require("./composer");
const proto_includes_1 = require("./proto-includes");
// Required for encoding and decoding queries that are of type Long.
// Must be done once but since the individal modules should be usable
// - must be set in each module that encounters encoding/decoding Longs.
// Reference: https://github.com/protobufjs/protobuf.js/issues/921
protobufjs_1.default.util.Long = long_1.default;
protobufjs_1.default.configure();
class Post {
    constructor(get, chainId, denoms) {
        this.accountNumberCache = new Map();
        this.get = get;
        this.chainId = chainId;
        this.registry = (0, registry_1.generateRegistry)();
        this.composer = new composer_1.Composer();
        this.denoms = denoms;
        this.defaultGasPrice = stargate_1.GasPrice
            .fromString(`0.025${denoms.USDC_GAS_DENOM !== undefined ? denoms.USDC_GAS_DENOM : denoms.USDC_DENOM}`);
        this.defaultDydxGasPrice = stargate_1.GasPrice
            .fromString(`25000000000${denoms.CHAINTOKEN_GAS_DENOM !== undefined ? denoms.CHAINTOKEN_GAS_DENOM : denoms.CHAINTOKEN_DENOM}`);
    }
    /**
     * @description Simulate a transaction
     * the calling function is responsible for creating the messages.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The Fee for broadcasting a transaction.
     */
    async simulate(wallet, messaging, gasPrice = this.defaultGasPrice, memo, account) {
        const msgsPromise = messaging();
        const accountPromise = account ? (await account()) : this.account(wallet.address);
        const msgsAndAccount = await Promise.all([msgsPromise, accountPromise]);
        const msgs = msgsAndAccount[0];
        return this.simulateTransaction(wallet.pubKey, msgsAndAccount[1].sequence, msgs, gasPrice, memo);
    }
    /**
     * @description Sign a transaction
     * the calling function is responsible for creating the messages.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The Signature.
     */
    async sign(wallet, messaging, zeroFee, gasPrice = this.defaultGasPrice, memo, account) {
        const msgsPromise = await messaging();
        const accountPromise = account ? (await account()) : this.account(wallet.address);
        const msgsAndAccount = await Promise.all([msgsPromise, accountPromise]);
        const msgs = msgsAndAccount[0];
        return this.signTransaction(wallet, msgs, msgsAndAccount[1], zeroFee, gasPrice, memo);
    }
    /**
     * @description Send a transaction
     * the calling function is responsible for creating the messages.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The Tx Hash.
     */
    async send(wallet, messaging, zeroFee, gasPrice = this.defaultGasPrice, memo, broadcastMode, account) {
        const msgsPromise = messaging();
        const accountPromise = account ? (await account()) : this.account(wallet.address);
        const msgsAndAccount = await Promise.all([msgsPromise, accountPromise]);
        const msgs = msgsAndAccount[0];
        return this.signAndSendTransaction(wallet, msgsAndAccount[1], msgs, zeroFee, gasPrice, memo, broadcastMode !== null && broadcastMode !== void 0 ? broadcastMode : this.defaultBroadcastMode(msgs));
    }
    /**
     * @description Calculate the default broadcast mode.
     */
    defaultBroadcastMode(msgs) {
        var _a, _b;
        if (msgs.length === 1 && msgs[0].typeUrl === '/dydxprotocol.clob.MsgPlaceOrder') {
            const msg = msgs[0].value;
            const orderFlags = (_b = (_a = msg.order) === null || _a === void 0 ? void 0 : _a.orderId) === null || _b === void 0 ? void 0 : _b.orderFlags;
            switch (orderFlags) {
                case types_1.OrderFlags.SHORT_TERM:
                    return tendermint_rpc_1.Method.BroadcastTxSync;
                case types_1.OrderFlags.LONG_TERM:
                    return tendermint_rpc_1.Method.BroadcastTxCommit;
                default:
                    break;
            }
        }
        return tendermint_rpc_1.Method.BroadcastTxSync;
    }
    /**
     * @description Sign and send a message
     *
     * @returns The Tx Response.
     */
    async signTransaction(wallet, messages, account, zeroFee, gasPrice = this.defaultGasPrice, memo) {
        // Simulate transaction if no fee is specified.
        const fee = zeroFee ? {
            amount: [],
            gas: '1000000',
        } : await this.simulateTransaction(wallet.pubKey, account.sequence, messages, gasPrice, memo);
        const txOptions = {
            sequence: account.sequence,
            accountNumber: account.accountNumber,
            chainId: this.chainId,
        };
        // Generate signed transaction.
        return wallet.signTransaction(messages, txOptions, fee, memo);
    }
    /**
     * @description Retrieve an account structure for transactions.
     * For short term orders, the sequence doesn't matter. Use cached if available.
     * For long term and conditional orders, a round trip to validator must be made.
     */
    async account(address, orderFlags) {
        if (orderFlags === types_1.OrderFlags.SHORT_TERM) {
            if (this.accountNumberCache.has(address)) {
                // For SHORT_TERM orders, the sequence doesn't matter
                return this.accountNumberCache.get(address);
            }
        }
        const account = await this.get.getAccount(address);
        this.accountNumberCache.set(address, account);
        return account;
    }
    /**
     * @description Sign and send a message
     *
     * @returns The Tx Response.
     */
    async signAndSendTransaction(wallet, account, messages, zeroFee, gasPrice = this.defaultGasPrice, memo, broadcastMode) {
        const signedTransaction = await this.signTransaction(wallet, messages, account, zeroFee, gasPrice, memo);
        return this.sendSignedTransaction(signedTransaction, broadcastMode);
    }
    /**
     * @description Send signed transaction.
     *
     * @returns The Tx Response.
     */
    async sendSignedTransaction(signedTransaction, broadcastMode) {
        return this.get.tendermintClient.broadcastTransaction(signedTransaction, broadcastMode !== undefined
            ? broadcastMode
            : tendermint_rpc_1.Method.BroadcastTxSync);
    }
    /**
     * @description Simulate broadcasting a transaction.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The Fee for broadcasting a transaction.
     */
    async simulateTransaction(pubKey, sequence, messages, gasPrice = this.defaultGasPrice, memo) {
        // Get simulated response.
        const encodedMessages = messages.map((message) => this.registry.encodeAsAny(message));
        const simulationResponse = await this.get.stargateQueryClient.tx.simulate(encodedMessages, memo, pubKey, sequence);
        // The promise should have been rejected if the gasInfo was undefined.
        if (simulationResponse.gasInfo === undefined) {
            throw new errors_1.UnexpectedClientError();
        }
        // Calculate and return fee from gasEstimate.
        const gasEstimate = math_1.Uint53.fromString(simulationResponse.gasInfo.gasUsed.toString()).toNumber();
        const fee = (0, stargate_1.calculateFee)(Math.floor(gasEstimate * constants_1.GAS_MULTIPLIER), gasPrice);
        // TODO(TRCL-2550): Temporary workaround before IBC denom is supported in '@cosmjs/stargate'.
        // The '@cosmjs/stargate' does not support denom with '/', so currently GAS_PRICE is
        // represented in 'uusdc', and the output of `calculateFee` is in '', which is replaced
        // below by USDC_DENOM string.
        const amount = lodash_1.default.map(fee.amount, (coin) => {
            if (coin.denom === 'uusdc') {
                return {
                    amount: coin.amount,
                    denom: this.denoms.USDC_DENOM,
                };
            }
            return coin;
        });
        return {
            ...fee,
            amount,
        };
    }
    // ------ State-Changing Requests ------ //
    async placeOrder(subaccount, clientId, clobPairId, side, quantums, subticks, timeInForce, orderFlags, reduceOnly, goodTilBlock, goodTilBlockTime, clientMetadata = 0, conditionType = proto_includes_1.Order_ConditionType.CONDITION_TYPE_UNSPECIFIED, conditionalOrderTriggerSubticks = long_1.default.fromInt(0), broadcastMode) {
        const msgs = new Promise((resolve) => {
            const msg = this.composer.composeMsgPlaceOrder(subaccount.address, subaccount.subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock !== null && goodTilBlock !== void 0 ? goodTilBlock : 0, goodTilBlockTime !== null && goodTilBlockTime !== void 0 ? goodTilBlockTime : 0, side, quantums, subticks, timeInForce, reduceOnly, clientMetadata, conditionType, conditionalOrderTriggerSubticks);
            resolve([msg]);
        });
        const account = this.account(subaccount.address, orderFlags);
        return this.send(subaccount.wallet, () => msgs, true, undefined, undefined, broadcastMode, () => account);
    }
    async placeOrderObject(subaccount, placeOrder, broadcastMode) {
        var _a, _b;
        return this.placeOrder(subaccount, placeOrder.clientId, placeOrder.clobPairId, placeOrder.side, placeOrder.quantums, placeOrder.subticks, placeOrder.timeInForce, placeOrder.orderFlags, placeOrder.reduceOnly, placeOrder.goodTilBlock, placeOrder.goodTilBlockTime, placeOrder.clientMetadata, (_a = placeOrder.conditionType) !== null && _a !== void 0 ? _a : proto_includes_1.Order_ConditionType.CONDITION_TYPE_UNSPECIFIED, (_b = placeOrder.conditionalOrderTriggerSubticks) !== null && _b !== void 0 ? _b : long_1.default.fromInt(0), broadcastMode);
    }
    async cancelOrder(subaccount, clientId, orderFlags, clobPairId, goodTilBlock, goodTilBlockTime, broadcastMode) {
        const msgs = new Promise((resolve) => {
            const msg = this.composer.composeMsgCancelOrder(subaccount.address, subaccount.subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock !== null && goodTilBlock !== void 0 ? goodTilBlock : 0, goodTilBlockTime !== null && goodTilBlockTime !== void 0 ? goodTilBlockTime : 0);
            resolve([msg]);
        });
        return this.send(subaccount.wallet, () => msgs, true, undefined, undefined, broadcastMode);
    }
    async cancelOrderObject(subaccount, cancelOrder, broadcastMode) {
        return this.cancelOrder(subaccount, cancelOrder.clientId, cancelOrder.orderFlags, cancelOrder.clobPairId, cancelOrder.goodTilBlock, cancelOrder.goodTilBlockTime, broadcastMode);
    }
    async transfer(subaccount, recipientAddress, recipientSubaccountNumber, assetId, amount, broadcastMode) {
        const msgs = new Promise((resolve) => {
            const msg = this.composer.composeMsgTransfer(subaccount.address, subaccount.subaccountNumber, recipientAddress, recipientSubaccountNumber, assetId, amount);
            resolve([msg]);
        });
        return this.send(subaccount.wallet, () => msgs, false, undefined, undefined, broadcastMode);
    }
    async deposit(subaccount, assetId, quantums, broadcastMode) {
        const msgs = new Promise((resolve) => {
            const msg = this.composer.composeMsgDepositToSubaccount(subaccount.address, subaccount.subaccountNumber, assetId, quantums);
            resolve([msg]);
        });
        return this.send(subaccount.wallet, () => msgs, false, undefined, undefined, broadcastMode);
    }
    async withdraw(subaccount, assetId, quantums, recipient, broadcastMode) {
        const msgs = new Promise((resolve) => {
            const msg = this.composer.composeMsgWithdrawFromSubaccount(subaccount.address, subaccount.subaccountNumber, assetId, quantums, recipient);
            resolve([msg]);
        });
        return this.send(subaccount.wallet, () => msgs, false, undefined, undefined, broadcastMode);
    }
    async sendToken(subaccount, recipient, coinDenom, quantums, zeroFee = true, broadcastMode) {
        if (coinDenom !== this.denoms.CHAINTOKEN_DENOM && coinDenom !== this.denoms.USDC_DENOM) {
            throw new Error('Unsupported coinDenom');
        }
        const msgs = new Promise((resolve) => {
            const msg = this.composer.composeMsgSendToken(subaccount.address, recipient, coinDenom, quantums);
            resolve([msg]);
        });
        return this.send(subaccount.wallet, () => msgs, zeroFee, coinDenom === this.denoms.CHAINTOKEN_DENOM
            ? this.defaultDydxGasPrice
            : this.defaultGasPrice, undefined, broadcastMode);
    }
}
exports.Post = Post;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicG9zdC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9jbGllbnRzL21vZHVsZXMvcG9zdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFDQSx1Q0FBc0M7QUFLdEMsK0NBTTBCO0FBQzFCLDJEQUVnQztBQUtoQyxvREFBdUI7QUFDdkIsZ0RBQXdCO0FBQ3hCLDREQUFrQztBQUVsQyw0Q0FBOEM7QUFDOUMsMENBQXNEO0FBQ3RELDhDQUFtRDtBQUVuRCxvQ0FPa0I7QUFDbEIseUNBQXNDO0FBR3RDLHFEQUUwQjtBQUUxQixvRUFBb0U7QUFDcEUscUVBQXFFO0FBQ3JFLHdFQUF3RTtBQUN4RSxrRUFBa0U7QUFDbEUsb0JBQVEsQ0FBQyxJQUFJLENBQUMsSUFBSSxHQUFHLGNBQUksQ0FBQztBQUMxQixvQkFBUSxDQUFDLFNBQVMsRUFBRSxDQUFDO0FBRXJCLE1BQWEsSUFBSTtJQVliLFlBQ0UsR0FBUSxFQUNSLE9BQWUsRUFDZixNQUFtQjtRQUxiLHVCQUFrQixHQUF5QixJQUFJLEdBQUcsRUFBRSxDQUFDO1FBTzNELElBQUksQ0FBQyxHQUFHLEdBQUcsR0FBRyxDQUFDO1FBQ2YsSUFBSSxDQUFDLE9BQU8sR0FBRyxPQUFPLENBQUM7UUFDdkIsSUFBSSxDQUFDLFFBQVEsR0FBRyxJQUFBLDJCQUFnQixHQUFFLENBQUM7UUFDbkMsSUFBSSxDQUFDLFFBQVEsR0FBRyxJQUFJLG1CQUFRLEVBQUUsQ0FBQztRQUMvQixJQUFJLENBQUMsTUFBTSxHQUFHLE1BQU0sQ0FBQztRQUNyQixJQUFJLENBQUMsZUFBZSxHQUFHLG1CQUFRO2FBQzVCLFVBQVUsQ0FBQyxRQUFRLE1BQU0sQ0FBQyxjQUFjLEtBQUssU0FBUyxDQUFDLENBQUMsQ0FBQyxNQUFNLENBQUMsY0FBYyxDQUFDLENBQUMsQ0FBQyxNQUFNLENBQUMsVUFBVSxFQUFFLENBQUMsQ0FBQztRQUN6RyxJQUFJLENBQUMsbUJBQW1CLEdBQUcsbUJBQVE7YUFDaEMs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