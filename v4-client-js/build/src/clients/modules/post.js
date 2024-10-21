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
            .fromString(`0.025${denoms.TDAI_GAS_DENOM !== undefined ? denoms.TDAI_GAS_DENOM : denoms.TDAI_DENOM}`);
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
        // represented in 'utdai', and the output of `calculateFee` is in '', which is replaced
        // below by TDAI_DENOM string.
        const amount = lodash_1.default.map(fee.amount, (coin) => {
            if (coin.denom === 'utdai') {
                return {
                    amount: coin.amount,
                    denom: this.denoms.TDAI_DENOM,
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
    async placeOrder(subaccount, clientId, clobPairId, side, quantums, subticks, timeInForce, orderFlags, reduceOnly, goodTilBlock, goodTilBlockTime, clientMetadata = 0, conditionType = proto_includes_1.Order_ConditionType.CONDITION_TYPE_UNSPECIFIED, conditionalOrderTriggerSubticks = long_1.default.fromInt(0), routerFeePpm = 0, routerFeeSubaccountOwner = '', routerFeeSubaccountNumber = 0, broadcastMode) {
        const msgs = new Promise((resolve) => {
            const msg = this.composer.composeMsgPlaceOrder(subaccount.address, subaccount.subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock !== null && goodTilBlock !== void 0 ? goodTilBlock : 0, goodTilBlockTime !== null && goodTilBlockTime !== void 0 ? goodTilBlockTime : 0, side, quantums, subticks, timeInForce, reduceOnly, clientMetadata, conditionType, conditionalOrderTriggerSubticks, routerFeePpm, routerFeeSubaccountOwner, routerFeeSubaccountNumber);
            resolve([msg]);
        });
        const account = this.account(subaccount.address, orderFlags);
        return this.send(subaccount.wallet, () => msgs, true, undefined, undefined, broadcastMode, () => account);
    }
    async placeOrderObject(subaccount, placeOrder, broadcastMode) {
        var _a, _b;
        return this.placeOrder(subaccount, placeOrder.clientId, placeOrder.clobPairId, placeOrder.side, placeOrder.quantums, placeOrder.subticks, placeOrder.timeInForce, placeOrder.orderFlags, placeOrder.reduceOnly, placeOrder.goodTilBlock, placeOrder.goodTilBlockTime, placeOrder.clientMetadata, (_a = placeOrder.conditionType) !== null && _a !== void 0 ? _a : proto_includes_1.Order_ConditionType.CONDITION_TYPE_UNSPECIFIED, (_b = placeOrder.conditionalOrderTriggerSubticks) !== null && _b !== void 0 ? _b : long_1.default.fromInt(0), placeOrder.routerFeePpm, placeOrder.routerFeeSubaccountOwner, placeOrder.routerFeeSubaccountNumber, broadcastMode);
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
        if (coinDenom !== this.denoms.CHAINTOKEN_DENOM && coinDenom !== this.denoms.TDAI_DENOM) {
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