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
    async placeShortTermOrder(subaccount, marketId, side, price, size, clientId, goodTilBlock, timeInForce, reduceOnly, routerFeePpm = 0, routerFeeSubaccountOwner = '', routerFeeSubaccountNumber = 0) {
        const msgs = new Promise((resolve) => {
            const msg = this.placeShortTermOrderMessage(subaccount, marketId, side, price, size, clientId, goodTilBlock, timeInForce, reduceOnly, routerFeePpm, routerFeeSubaccountOwner, routerFeeSubaccountNumber);
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
    async placeOrder(subaccount, marketId, type, side, price, size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, triggerPrice, marketInfo, currentHeight, routerFeePpm, routerFeeSubaccountOwner, routerFeeSubaccountNumber) {
        const msgs = new Promise((resolve) => {
            const msg = this.placeOrderMessage(subaccount, marketId, type, side, price, 
            // trigger_price: number,   // not used for MARKET and LIMIT
            size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, triggerPrice, marketInfo, currentHeight, routerFeePpm, routerFeeSubaccountOwner, routerFeeSubaccountNumber);
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
    size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, triggerPrice, marketInfo, currentHeight, routerFeePpm, routerFeeSubaccountOwner, routerFeeSubaccountNumber) {
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
        const finalRouterFeePpm = routerFeePpm !== null && routerFeePpm !== void 0 ? routerFeePpm : 0;
        const finalRouterFeeSubaccountOwner = routerFeeSubaccountOwner !== null && routerFeeSubaccountOwner !== void 0 ? routerFeeSubaccountOwner : '';
        const finalRouterFeeSubaccountNumber = routerFeeSubaccountNumber !== null && routerFeeSubaccountNumber !== void 0 ? routerFeeSubaccountNumber : 0;
        const clientMetadata = (0, chain_helpers_1.calculateClientMetadata)(type);
        const conditionalType = (0, chain_helpers_1.calculateConditionType)(type);
        const conditionalOrderTriggerSubticks = (0, chain_helpers_1.calculateConditionalOrderTriggerSubticks)(type, atomicResolution, quantumConversionExponent, subticksPerTick, triggerPrice);
        return this.validatorClient.post.composer.composeMsgPlaceOrder(subaccount.address, subaccount.subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock, goodTilBlockTime, orderSide, quantums, subticks, orderTimeInForce, reduceOnly !== null && reduceOnly !== void 0 ? reduceOnly : false, clientMetadata, conditionalType, conditionalOrderTriggerSubticks, finalRouterFeePpm, finalRouterFeeSubaccountOwner, finalRouterFeeSubaccountNumber);
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
    async placeShortTermOrderMessage(subaccount, marketId, side, price, size, clientId, goodTilBlock, timeInForce, reduceOnly, routerFeePpm = 0, routerFeeSubaccountOwner = '', routerFeeSubaccountNumber = 0) {
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
        long_1.default.fromInt(0), // Short term orders cannot be conditional.
        routerFeePpm, routerFeeSubaccountOwner, routerFeeSubaccountNumber);
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
    size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, routerFeePpm, routerFeeSubaccountOwner, routerFeeSubaccountNumber) {
        const msgs = new Promise((resolve) => {
            const msg = this.placeOrderMessage(subaccount, marketId, type, side, price, 
            // trigger_price: number,   // not used for MARKET and LIMIT
            size, clientId, timeInForce, goodTilTimeInSeconds, execution, postOnly, reduceOnly, routerFeePpm, undefined, undefined, undefined, routerFeeSubaccountOwner, routerFeeSubaccountNumber);
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