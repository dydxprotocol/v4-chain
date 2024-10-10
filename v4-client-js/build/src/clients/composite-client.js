"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.CompositeClient = void 0;
const order_1 = require("@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/order");
const ethers_1 = require("ethers");
const long_1 = __importDefault(require("long"));
const protobufjs_1 = __importDefault(require("protobufjs"));
const validation_1 = require("../lib/validation");
const types_1 = require("../types");
const constants_1 = require("./constants");
const chain_helpers_1 = require("./helpers/chain-helpers");
const indexer_client_1 = require("./indexer-client");
const errors_1 = require("./lib/errors");
const validator_client_1 = require("./validator-client");
// Required for encoding and decoding queries that are of type Long.
// Must be done once but since the individal modules should be usable
// - must be set in each module that encounters encoding/decoding Longs.
// Reference: https://github.com/protobufjs/protobuf.js/issues/921
protobufjs_1.default.util.Long = long_1.default;
protobufjs_1.default.configure();
class CompositeClient {
    static async connect(network) {
        const client = new CompositeClient(network);
        await client.initialize();
        return client;
    }
    constructor(network, apiTimeout) {
        this.network = network;
        this._indexerClient = new indexer_client_1.IndexerClient(network.indexerConfig, apiTimeout);
    }
    async initialize() {
        this._validatorClient = await validator_client_1.ValidatorClient.connect(this.network.validatorConfig);
    }
    get indexerClient() {
        /**
         * Get the validator client
         */
        return this._indexerClient;
    }
    get validatorClient() {
        /**
         * Get the validator client
         */
        return this._validatorClient;
    }
    /**
       * @description Sign a list of messages with a wallet.
       * the calling function is responsible for creating the messages.
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The Signature.
       */
    async sign(wallet, messaging, zeroFee, gasPrice, memo, account) {
        return this.validatorClient.post.sign(wallet, messaging, zeroFee, gasPrice, memo, account);
    }
    /**
       * @description Send a list of messages with a wallet.
       * the calling function is responsible for creating the messages.
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The Transaction Hash.
       */
    async send(wallet, messaging, zeroFee, gasPrice, memo, account) {
        return this.validatorClient.post.send(wallet, messaging, zeroFee, gasPrice, memo, undefined, account);
    }
    /**
       * @description Send a signed transaction.
       *
       * @param signedTransaction The signed transaction to send.
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The Transaction Hash.
       */
    async sendSignedTransaction(signedTransaction) {
        return this.validatorClient.post.sendSignedTransaction(signedTransaction);
    }
    /**
       * @description Simulate a list of messages with a wallet.
       * the calling function is responsible for creating the messages.
       *
       * To send multiple messages with gas estimate:
       * 1. Client is responsible for creating the messages.
       * 2. Call simulate() to get the gas estimate.
       * 3. Call send() to send the messages.
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The gas estimate.
       */
    async simulate(wallet, messaging, gasPrice, memo, account) {
        return this.validatorClient.post.simulate(wallet, messaging, gasPrice, memo, account);
    }
    /**
       * @description Calculate the goodTilBlock value for a SHORT_TERM order
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The goodTilBlock value
       */
    async calculateGoodTilBlock(orderFlags, currentHeight) {
        if (orderFlags === types_1.OrderFlags.SHORT_TERM) {
            const height = currentHeight !== null && currentHeight !== void 0 ? currentHeight : await this.validatorClient.get.latestBlockHeight();
            return height + constants_1.SHORT_BLOCK_FORWARD;
        }
        else {
            return Promise.resolve(0);
        }
    }
    /**
     * @description Validate the goodTilBlock value for a SHORT_TERM order
     *
     * @param goodTilBlock Number of blocks from the current block height the order will
     * be valid for.
     *
     * @throws UserError if the goodTilBlock value is not valid given latest block height and
     * SHORT_BLOCK_WINDOW.
     */
    async validateGoodTilBlock(goodTilBlock) {
        const height = await this.validatorClient.get.latestBlockHeight();
        const nextValidBlockHeight = height + 1;
        const lowerBound = nextValidBlockHeight;
        const upperBound = nextValidBlockHeight + constants_1.SHORT_BLOCK_WINDOW;
        if (goodTilBlock < lowerBound || goodTilBlock > upperBound) {
            throw new errors_1.UserError(`Invalid Short-Term order GoodTilBlock.
        Should be greater-than-or-equal-to ${lowerBound} and less-than-or-equal-to ${upperBound}.
        Provided good til block: ${goodTilBlock}`);
        }
    }
    /**
       * @description Calculate the goodTilBlockTime value for a LONG_TERM order
       * the calling function is responsible for creating the messages.
       *
       * @param goodTilTimeInSeconds The goodTilTimeInSeconds of the order to place.
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The goodTilBlockTime value
       */
    calculateGoodTilBlockTime(goodTilTimeInSeconds) {
        const now = new Date();
        const millisecondsPerSecond = 1000;
        const interval = goodTilTimeInSeconds * millisecondsPerSecond;
        const future = new Date(now.valueOf() + interval);
        return Math.round(future.getTime() / 1000);
    }
    /**
     * @description Place a short term order with human readable input.
     *
     * Use human readable form of input, including price and size
     * The quantum and subticks are calculated and submitted
     *
     * @param subaccount The subaccount to place the order under
     * @param marketId The market to place the order on
     * @param side The side of the order to place
     * @param price The price of the order to place
     * @param size The size of the order to place
     * @param clientId The client id of the order to place
     * @param timeInForce The time in force of the order to place
     * @param goodTilBlock The goodTilBlock of the order to place
     * @param reduceOnly The reduceOnly of the order to place
     *
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The transaction hash.
     */
    async placeShortTermOrder(subaccount, marketId, side, price, size, clientId, goodTilBlock, timeInForce, reduceOnly) {
        const msgs = new Promise((resolve) => {
            const msg = this.placeShortTermOrderMessage(subaccount, marketId, side, price, size, clientId, goodTilBlock, timeInForce, reduceOnly);
            msg.then((it) => resolve([it])).catch((err) => {
                console.log(err);
            });
        });
        const account = this.validatorClient.post.account(subaccount.address, undefined);
        return this.send(subaccount.wallet, () => msgs, true, undefined, undefined, () => account);
    }
    /**
       * @description Place an order with human readable input.
       *
       * Only MARKET and LIMIT types are supported right now
       * Use human readable form of input, including price and size
       * The quantum and subticks are calculated and submitted
       *
       * @param subaccount The subaccount to place the order on.
       * @param marketId The market to place the order on.
       * @param type The type of order to place.
       * @param side The side of the order to place.
       * @param price The price of the order to place.
       * @param size The size of the order to place.
       * @param clientId The client id of the order to place.
       * @param timeInForce The time in force of the order to place.
       * @param goodTilTimeInSeconds The goodTilTimeInSeconds of the order to place.
       * @param execution The execution of the order to place.
       * @param postOnly The postOnly of the order to place.
       * @param reduceOnly The reduceOnly of the order to place.
       * @param triggerPrice The trigger price of conditional orders.
       * @param marketInfo optional market information for calculating quantums and subticks.
       *        This can be constructed from Indexer API. If set to null, additional round
       *        trip to Indexer API will be made.
       * @param currentHeight Current block height. This can be obtained from ValidatorClient.
       *        If set to null, additional round trip to ValidatorClient will be made.
       *
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The transaction hash.
       */
    async placeOrder(subaccount, marketId, type, side, price, size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, triggerPrice, marketInfo, currentHeight) {
        const msgs = new Promise((resolve) => {
            const msg = this.placeOrderMessage(subaccount, marketId, type, side, price, 
            // trigger_price: number,   // not used for MARKET and LIMIT
            size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, triggerPrice, marketInfo, currentHeight);
            msg.then((it) => resolve([it])).catch((err) => {
                console.log(err);
            });
        });
        const orderFlags = (0, chain_helpers_1.calculateOrderFlags)(type, timeInForce);
        const account = this.validatorClient.post.account(subaccount.address, orderFlags);
        return this.send(subaccount.wallet, () => msgs, true, undefined, undefined, () => account);
    }
    /**
       * @description Calculate and create the place order message
       *
       * Only MARKET and LIMIT types are supported right now
       * Use human readable form of input, including price and size
       * The quantum and subticks are calculated and submitted
       *
       * @param subaccount The subaccount to place the order under
       * @param marketId The market to place the order on
       * @param type The type of order to place
       * @param side The side of the order to place
       * @param price The price of the order to place
       * @param size The size of the order to place
       * @param clientId The client id of the order to place
       * @param timeInForce The time in force of the order to place
       * @param goodTilTimeInSeconds The goodTilTimeInSeconds of the order to place
       * @param execution The execution of the order to place
       * @param postOnly The postOnly of the order to place
       * @param reduceOnly The reduceOnly of the order to place
       *
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The message to be passed into the protocol
       */
    async placeOrderMessage(subaccount, marketId, type, side, price, 
    // trigger_price: number,   // not used for MARKET and LIMIT
    size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, triggerPrice, marketInfo, currentHeight) {
        const orderFlags = (0, chain_helpers_1.calculateOrderFlags)(type, timeInForce);
        const result = await Promise.all([
            this.calculateGoodTilBlock(orderFlags, currentHeight),
            this.retrieveMarketInfo(marketId, marketInfo),
        ]);
        const goodTilBlock = result[0];
        const clobPairId = result[1].clobPairId;
        const atomicResolution = result[1].atomicResolution;
        const stepBaseQuantums = result[1].stepBaseQuantums;
        const quantumConversionExponent = result[1].quantumConversionExponent;
        const subticksPerTick = result[1].subticksPerTick;
        const orderSide = (0, chain_helpers_1.calculateSide)(side);
        const quantums = (0, chain_helpers_1.calculateQuantums)(size, atomicResolution, stepBaseQuantums);
        const subticks = (0, chain_helpers_1.calculateSubticks)(price, atomicResolution, quantumConversionExponent, subticksPerTick);
        const orderTimeInForce = (0, chain_helpers_1.calculateTimeInForce)(type, timeInForce, execution, postOnly);
        let goodTilBlockTime = 0;
        if (orderFlags === types_1.OrderFlags.LONG_TERM || orderFlags === types_1.OrderFlags.CONDITIONAL) {
            if (goodTilTimeInSeconds == null) {
                throw new Error('goodTilTimeInSeconds must be set for LONG_TERM or CONDITIONAL order');
            }
            else {
                goodTilBlockTime = this.calculateGoodTilBlockTime(goodTilTimeInSeconds);
            }
        }
        const clientMetadata = (0, chain_helpers_1.calculateClientMetadata)(type);
        const conditionalType = (0, chain_helpers_1.calculateConditionType)(type);
        const conditionalOrderTriggerSubticks = (0, chain_helpers_1.calculateConditionalOrderTriggerSubticks)(type, atomicResolution, quantumConversionExponent, subticksPerTick, triggerPrice);
        return this.validatorClient.post.composer.composeMsgPlaceOrder(subaccount.address, subaccount.subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock, goodTilBlockTime, orderSide, quantums, subticks, orderTimeInForce, reduceOnly !== null && reduceOnly !== void 0 ? reduceOnly : false, clientMetadata, conditionalType, conditionalOrderTriggerSubticks);
    }
    async retrieveMarketInfo(marketId, marketInfo) {
        if (marketInfo) {
            return Promise.resolve(marketInfo);
        }
        else {
            const marketsResponse = await this.indexerClient.markets.getPerpetualMarkets(marketId);
            const market = marketsResponse.markets[marketId];
            const clobPairId = market.clobPairId;
            const atomicResolution = market.atomicResolution;
            const stepBaseQuantums = market.stepBaseQuantums;
            const quantumConversionExponent = market.quantumConversionExponent;
            const subticksPerTick = market.subticksPerTick;
            return {
                clobPairId,
                atomicResolution,
                stepBaseQuantums,
                quantumConversionExponent,
                subticksPerTick,
            };
        }
    }
    /**
       * @description Calculate and create the short term place order message
       *
       * Use human readable form of input, including price and size
       * The quantum and subticks are calculated and submitted
       *
       * @param subaccount The subaccount to place the order under
       * @param marketId The market to place the order on
       * @param side The side of the order to place
       * @param price The price of the order to place
       * @param size The size of the order to place
       * @param clientId The client id of the order to place
       * @param timeInForce The time in force of the order to place
       * @param goodTilBlock The goodTilBlock of the order to place
       * @param reduceOnly The reduceOnly of the order to place
       *
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The message to be passed into the protocol
       */
    async placeShortTermOrderMessage(subaccount, marketId, side, price, size, clientId, goodTilBlock, timeInForce, reduceOnly) {
        await this.validateGoodTilBlock(goodTilBlock);
        const marketsResponse = await this.indexerClient.markets.getPerpetualMarkets(marketId);
        const market = marketsResponse.markets[marketId];
        const clobPairId = market.clobPairId;
        const atomicResolution = market.atomicResolution;
        const stepBaseQuantums = market.stepBaseQuantums;
        const quantumConversionExponent = market.quantumConversionExponent;
        const subticksPerTick = market.subticksPerTick;
        const orderSide = (0, chain_helpers_1.calculateSide)(side);
        const quantums = (0, chain_helpers_1.calculateQuantums)(size, atomicResolution, stepBaseQuantums);
        const subticks = (0, chain_helpers_1.calculateSubticks)(price, atomicResolution, quantumConversionExponent, subticksPerTick);
        const orderFlags = types_1.OrderFlags.SHORT_TERM;
        return this.validatorClient.post.composer.composeMsgPlaceOrder(subaccount.address, subaccount.subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock, 0, // Short term orders use goodTilBlock.
        orderSide, quantums, subticks, timeInForce, reduceOnly, 0, // Client metadata is 0 for short term orders.
        order_1.Order_ConditionType.CONDITION_TYPE_UNSPECIFIED, // Short term orders cannot be conditional.
        long_1.default.fromInt(0));
    }
    /**
       * @description Cancel an order with order information from web socket or REST.
       *
       * @param subaccount The subaccount to cancel the order from
       * @param clientId The client id of the order to cancel
       * @param orderFlags The order flags of the order to cancel
       * @param clobPairId The clob pair id of the order to cancel
       * @param goodTilBlock The goodTilBlock of the order to cancel
       * @param goodTilBlockTime The goodTilBlockTime of the order to cancel
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The transaction hash.
       */
    async cancelRawOrder(subaccount, clientId, orderFlags, clobPairId, goodTilBlock, goodTilBlockTime) {
        return this.validatorClient.post.cancelOrder(subaccount, clientId, orderFlags, clobPairId, goodTilBlock, goodTilBlockTime);
    }
    /**
       * @description Cancel an order with human readable input.
       *
       * @param subaccount The subaccount to cancel the order from
       * @param clientId The client id of the order to cancel
       * @param orderFlags The order flags of the order to cancel
       * @param marketId The market to cancel the order on
       * @param goodTilBlock The goodTilBlock of the order to cancel
       * @param goodTilBlockTime The goodTilBlockTime of the order to cancel
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The transaction hash.
       */
    async cancelOrder(subaccount, clientId, orderFlags, marketId, goodTilBlock, goodTilTimeInSeconds) {
        const marketsResponse = await this.indexerClient.markets.getPerpetualMarkets(marketId);
        const market = marketsResponse.markets[marketId];
        const clobPairId = market.clobPairId;
        if (!(0, validation_1.verifyOrderFlags)(orderFlags)) {
            throw new Error(`Invalid order flags: ${orderFlags}`);
        }
        let goodTilBlockTime;
        if ((0, validation_1.isStatefulOrder)(orderFlags)) {
            if (goodTilTimeInSeconds === undefined || goodTilTimeInSeconds === 0) {
                throw new Error('goodTilTimeInSeconds must be set for LONG_TERM or CONDITIONAL order');
            }
            if (goodTilBlock !== 0) {
                throw new Error('goodTilBlock should be zero since LONG_TERM or CONDITIONAL orders ' +
                    'use goodTilTimeInSeconds instead of goodTilBlock.');
            }
            goodTilBlockTime = this.calculateGoodTilBlockTime(goodTilTimeInSeconds);
        }
        else {
            if (goodTilBlock === undefined || goodTilBlock === 0) {
                throw new Error('goodTilBlock must be non-zero for SHORT_TERM orders');
            }
            if (goodTilTimeInSeconds !== undefined && goodTilTimeInSeconds !== 0) {
                throw new Error('goodTilTimeInSeconds should be zero since SHORT_TERM orders use goodTilBlock instead of goodTilTimeInSeconds.');
            }
        }
        return this.validatorClient.post.cancelOrder(subaccount, clientId, orderFlags, clobPairId, goodTilBlock, goodTilBlockTime);
    }
    /**
       * @description Transfer from a subaccount to another subaccount
       *
       * @param subaccount The subaccount to transfer from
       * @param recipientAddress The recipient address
       * @param recipientSubaccountNumber The recipient subaccount number
       * @param amount The amount to transfer
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The transaction hash.
       */
    async transferToSubaccount(subaccount, recipientAddress, recipientSubaccountNumber, amount) {
        const msgs = new Promise((resolve) => {
            const msg = this.transferToSubaccountMessage(subaccount, recipientAddress, recipientSubaccountNumber, amount);
            resolve([msg]);
        });
        return this.send(subaccount.wallet, () => msgs, true);
    }
    /**
       * @description Create message to transfer from a subaccount to another subaccount
       *
       * @param subaccount The subaccount to transfer from
       * @param recipientAddress The recipient address
       * @param recipientSubaccountNumber The recipient subaccount number
       * @param amount The amount to transfer
       *
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The message
       */
    transferToSubaccountMessage(subaccount, recipientAddress, recipientSubaccountNumber, amount) {
        const validatorClient = this._validatorClient;
        if (validatorClient === undefined) {
            throw new Error('validatorClient not set');
        }
        const quantums = (0, ethers_1.parseUnits)(amount, validatorClient.config.denoms.TDAI_DECIMALS);
        if (quantums > BigInt(long_1.default.MAX_VALUE.toString())) {
            throw new Error('amount to large');
        }
        if (quantums < 0) {
            throw new Error('amount must be positive');
        }
        return this.validatorClient.post.composer.composeMsgTransfer(subaccount.address, subaccount.subaccountNumber, recipientAddress, recipientSubaccountNumber, 0, long_1.default.fromString(quantums.toString()));
    }
    /**
       * @description Deposit from wallet to subaccount
       *
       * @param subaccount The subaccount to deposit to
       * @param amount The amount to deposit
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The transaction hash.
       */
    async depositToSubaccount(subaccount, amount) {
        const msgs = new Promise((resolve) => {
            const msg = this.depositToSubaccountMessage(subaccount, amount);
            resolve([msg]);
        });
        return this.validatorClient.post.send(subaccount.wallet, () => msgs, false);
    }
    /**
       * @description Create message to deposit from wallet to subaccount
       *
       * @param subaccount The subaccount to deposit to
       * @param amount The amount to deposit
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The message
       */
    depositToSubaccountMessage(subaccount, amount) {
        const validatorClient = this._validatorClient;
        if (validatorClient === undefined) {
            throw new Error('validatorClient not set');
        }
        const quantums = (0, ethers_1.parseUnits)(amount, validatorClient.config.denoms.TDAI_DECIMALS);
        if (quantums > BigInt(long_1.default.MAX_VALUE.toString())) {
            throw new Error('amount to large');
        }
        if (quantums < 0) {
            throw new Error('amount must be positive');
        }
        return this.validatorClient.post.composer.composeMsgDepositToSubaccount(subaccount.address, subaccount.subaccountNumber, 0, long_1.default.fromString(quantums.toString()));
    }
    /**
       * @description Withdraw from subaccount to wallet
       *
       * @param subaccount The subaccount to withdraw from
       * @param amount The amount to withdraw
       * @param recipient The recipient address, default to subaccount address
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The transaction hash
       */
    async withdrawFromSubaccount(subaccount, amount, recipient) {
        const msgs = new Promise((resolve) => {
            const msg = this.withdrawFromSubaccountMessage(subaccount, amount, recipient);
            resolve([msg]);
        });
        return this.send(subaccount.wallet, () => msgs, false);
    }
    /**
       * @description Create message to withdraw from subaccount to wallet
       * with human readable input.
       *
       * @param subaccount The subaccount to withdraw from
       * @param amount The amount to withdraw
       * @param recipient The recipient address
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The message
       */
    withdrawFromSubaccountMessage(subaccount, amount, recipient) {
        const validatorClient = this._validatorClient;
        if (validatorClient === undefined) {
            throw new Error('validatorClient not set');
        }
        const quantums = (0, ethers_1.parseUnits)(amount, validatorClient.config.denoms.TDAI_DECIMALS);
        if (quantums > BigInt(long_1.default.MAX_VALUE.toString())) {
            throw new Error('amount to large');
        }
        if (quantums < 0) {
            throw new Error('amount must be positive');
        }
        return this.validatorClient.post.composer.composeMsgWithdrawFromSubaccount(subaccount.address, subaccount.subaccountNumber, 0, long_1.default.fromString(quantums.toString()), recipient);
    }
    /**
       * @description Create message to send chain token from subaccount to wallet
       * with human readable input.
       *
       * @param subaccount The subaccount to withdraw from
       * @param amount The amount to withdraw
       * @param recipient The recipient address
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The message
       */
    sendTokenMessage(wallet, amount, recipient) {
        var _a;
        const address = wallet.address;
        if (address === undefined) {
            throw new errors_1.UserError('wallet address is not set. Call connectWallet() first');
        }
        const { CHAINTOKEN_DENOM: chainTokenDenom, CHAINTOKEN_DECIMALS: chainTokenDecimals, } = ((_a = this._validatorClient) === null || _a === void 0 ? void 0 : _a.config.denoms) || {};
        if (chainTokenDenom === undefined || chainTokenDecimals === undefined) {
            throw new Error('Chain token denom not set in validator config');
        }
        const quantums = (0, ethers_1.parseUnits)(amount, chainTokenDecimals);
        return this.validatorClient.post.composer.composeMsgSendToken(address, recipient, chainTokenDenom, quantums.toString());
    }
    async signPlaceOrder(subaccount, marketId, type, side, price, 
    // trigger_price: number,   // not used for MARKET and LIMIT
    size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly) {
        const msgs = new Promise((resolve) => {
            const msg = this.placeOrderMessage(subaccount, marketId, type, side, price, 
            // trigger_price: number,   // not used for MARKET and LIMIT
            size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly);
            msg.then((it) => resolve([it])).catch((err) => {
                console.log(err);
            });
        });
        const signature = await this.sign(wallet, () => msgs, true);
        return Buffer.from(signature).toString('base64');
    }
    async signCancelOrder(subaccount, clientId, orderFlags, clobPairId, goodTilBlock, goodTilBlockTime) {
        const msgs = new Promise((resolve) => {
            const msg = this.validatorClient.post.composer.composeMsgCancelOrder(subaccount.address, subaccount.subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock, goodTilBlockTime);
            resolve([msg]);
        });
        const signature = await this.sign(subaccount.wallet, () => msgs, true);
        return Buffer.from(signature).toString('base64');
    }
}
exports.CompositeClient = CompositeClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29tcG9zaXRlLWNsaWVudC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9jbGllbnRzL2NvbXBvc2l0ZS1jbGllbnQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7O0FBS0Esc0ZBQW9IO0FBQ3BILG1DQUFvQztBQUNwQyxnREFBd0I7QUFDeEIsNERBQWtDO0FBRWxDLGtEQUFzRTtBQUN0RSxvQ0FBc0M7QUFDdEMsMkNBUXFCO0FBQ3JCLDJEQVNpQztBQUNqQyxxREFBaUQ7QUFDakQseUNBQXlDO0FBR3pDLHlEQUFxRDtBQUVyRCxvRUFBb0U7QUFDcEUscUVBQXFFO0FBQ3JFLHdFQUF3RTtBQUN4RSxrRUFBa0U7QUFDbEUsb0JBQVEsQ0FBQyxJQUFJLENBQUMsSUFBSSxHQUFHLGNBQUksQ0FBQztBQUMxQixvQkFBUSxDQUFDLFNBQVMsRUFBRSxDQUFDO0FBVXJCLE1BQWEsZUFBZTtJQUsxQixNQUFNLENBQUMsS0FBSyxDQUFDLE9BQU8sQ0FBQyxPQUFnQjtRQUNuQyxNQUFNLE1BQU0sR0FBRyxJQUFJLGVBQWUsQ0FBQyxPQUFPLENBQUMsQ0FBQztRQUM1QyxNQUFNLE1BQU0sQ0FBQyxVQUFVLEVBQUUsQ0FBQztRQUMxQixPQUFPLE1BQU0sQ0FBQztJQUNoQixDQUFDO0lBRUQsWUFDRSxPQUFnQixFQUNoQixVQUFtQjtRQUVuQixJQUFJLENBQUMsT0FBTyxHQUFHLE9BQU8sQ0FBQztRQUN2QixJQUFJLENBQUMsY0FBYyxHQUFHLElBQUksOEJBQWEsQ0FDckMsT0FBTyxDQUFDLGFBQWEsRUFDckIsVUFBVSxDQUNYLENBQUM7SUFDSixDQUFDO0lBRU8sS0FBSyxDQUFDLFVBQVU7UUFDdEIsSUFBSSxDQUFDLGdCQUFnQixHQUFHLE1BQU0sa0NBQWUsQ0FBQyxPQUFPLENBQUMsSUFBSSxDQUFDLE9BQU8sQ0FBQyxlQUFlLENBQUMsQ0FBQztJQUN0RixDQUFDO0lBRUQsSUFBSSxhQUFhO1FBQ2Y7O1dBRUc7UUFDSCxPQUFP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