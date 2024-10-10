import { EncodeObject } from '@cosmjs/proto-signing';
import {
  Account, GasPrice, IndexedTx, StdFee,
} from '@cosmjs/stargate';
import { BroadcastTxAsyncResponse, BroadcastTxSyncResponse } from '@cosmjs/tendermint-rpc/build/tendermint37';
import { Order_ConditionType, Order_TimeInForce } from '@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/order';
import { parseUnits } from 'ethers';
import Long from 'long';
import protobuf from 'protobufjs';

import { isStatefulOrder, verifyOrderFlags } from '../lib/validation';
import { OrderFlags } from '../types';
import {
  Network,
  OrderExecution,
  OrderSide,
  OrderTimeInForce,
  OrderType,
  SHORT_BLOCK_FORWARD,
  SHORT_BLOCK_WINDOW,
} from './constants';
import {
  calculateQuantums,
  calculateSubticks,
  calculateSide,
  calculateTimeInForce,
  calculateOrderFlags,
  calculateClientMetadata,
  calculateConditionType,
  calculateConditionalOrderTriggerSubticks,
} from './helpers/chain-helpers';
import { IndexerClient } from './indexer-client';
import { UserError } from './lib/errors';
import LocalWallet from './modules/local-wallet';
import { SubaccountInfo } from './subaccount';
import { ValidatorClient } from './validator-client';

// Required for encoding and decoding queries that are of type Long.
// Must be done once but since the individal modules should be usable
// - must be set in each module that encounters encoding/decoding Longs.
// Reference: https://github.com/protobufjs/protobuf.js/issues/921
protobuf.util.Long = Long;
protobuf.configure();

export interface MarketInfo {
  clobPairId: number;
  atomicResolution: number;
  stepBaseQuantums: number;
  quantumConversionExponent: number;
  subticksPerTick: number;
}

export class CompositeClient {
  public readonly network: Network;
  private _indexerClient: IndexerClient;
  private _validatorClient?: ValidatorClient;

  static async connect(network: Network): Promise<CompositeClient> {
    const client = new CompositeClient(network);
    await client.initialize();
    return client;
  }

  private constructor(
    network: Network,
    apiTimeout?: number,
  ) {
    this.network = network;
    this._indexerClient = new IndexerClient(
      network.indexerConfig,
      apiTimeout,
    );
  }

  private async initialize(): Promise<void> {
    this._validatorClient = await ValidatorClient.connect(this.network.validatorConfig);
  }

  get indexerClient(): IndexerClient {
    /**
     * Get the validator client
     */
    return this._indexerClient!;
  }

  get validatorClient(): ValidatorClient {
    /**
     * Get the validator client
     */
    return this._validatorClient!;
  }

  /**
     * @description Sign a list of messages with a wallet.
     * the calling function is responsible for creating the messages.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The Signature.
     */
  async sign(
    wallet: LocalWallet,
    messaging: () => Promise<EncodeObject[]>,
    zeroFee: boolean,
    gasPrice?: GasPrice,
    memo?: string,
    account?: () => Promise<Account>,
  ): Promise<Uint8Array> {
    return this.validatorClient.post.sign(
      wallet,
      messaging,
      zeroFee,
      gasPrice,
      memo,
      account,
    );
  }

  /**
     * @description Send a list of messages with a wallet.
     * the calling function is responsible for creating the messages.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The Transaction Hash.
     */
  async send(
    wallet: LocalWallet,
    messaging: () => Promise<EncodeObject[]>,
    zeroFee: boolean,
    gasPrice?: GasPrice,
    memo?: string,
    account?: () => Promise<Account>,
  ): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx> {
    return this.validatorClient.post.send(
      wallet,
      messaging,
      zeroFee,
      gasPrice,
      memo,
      undefined,
      account,
    );
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
  async sendSignedTransaction(
    signedTransaction: Uint8Array,
  ): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx> {
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
  async simulate(
    wallet: LocalWallet,
    messaging: () => Promise<EncodeObject[]>,
    gasPrice?: GasPrice,
    memo?: string,
    account?: () => Promise<Account>,
  ): Promise<StdFee> {
    return this.validatorClient.post.simulate(
      wallet,
      messaging,
      gasPrice,
      memo,
      account,
    );
  }

  /**
     * @description Calculate the goodTilBlock value for a SHORT_TERM order
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The goodTilBlock value
     */

  private async calculateGoodTilBlock(
    orderFlags: OrderFlags,
    currentHeight?: number,
  ): Promise<number> {
    if (orderFlags === OrderFlags.SHORT_TERM) {
      const height = currentHeight ?? await this.validatorClient.get.latestBlockHeight();
      return height + SHORT_BLOCK_FORWARD;
    } else {
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
  private async validateGoodTilBlock(goodTilBlock: number): Promise<void> {
    const height = await this.validatorClient.get.latestBlockHeight();
    const nextValidBlockHeight = height + 1;
    const lowerBound = nextValidBlockHeight;
    const upperBound = nextValidBlockHeight + SHORT_BLOCK_WINDOW;
    if (goodTilBlock < lowerBound || goodTilBlock > upperBound) {
      throw new UserError(`Invalid Short-Term order GoodTilBlock.
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
  private calculateGoodTilBlockTime(goodTilTimeInSeconds: number): number {
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
  async placeShortTermOrder(
    subaccount: SubaccountInfo,
    marketId: string,
    side: OrderSide,
    price: number,
    size: number,
    clientId: number,
    goodTilBlock: number,
    timeInForce: Order_TimeInForce,
    reduceOnly: boolean,
  ): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx> {
    const msgs: Promise<EncodeObject[]> = new Promise((resolve) => {
      const msg = this.placeShortTermOrderMessage(
        subaccount,
        marketId,
        side,
        price,
        size,
        clientId,
        goodTilBlock,
        timeInForce,
        reduceOnly,
      );
      msg.then((it) => resolve([it])).catch((err) => {
        console.log(err);
      });
    });
    const account: Promise<Account> = this.validatorClient.post.account(
      subaccount.address,
      undefined,
    );
    return this.send(
      subaccount.wallet,
      () => msgs,
      true,
      undefined,
      undefined,
      () => account,
    );
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
  async placeOrder(
    subaccount: SubaccountInfo,
    marketId: string,
    type: OrderType,
    side: OrderSide,
    price: number,
    size: number,
    clientId: number,
    timeInForce?: OrderTimeInForce,
    goodTilTimeInSeconds?: number,
    execution?: OrderExecution,
    postOnly?: boolean,
    reduceOnly?: boolean,
    triggerPrice?: number,
    marketInfo?: MarketInfo,
    currentHeight?: number,
  ): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx> {
    const msgs: Promise<EncodeObject[]> = new Promise((resolve) => {
      const msg = this.placeOrderMessage(
        subaccount,
        marketId,
        type,
        side,
        price,
        // trigger_price: number,   // not used for MARKET and LIMIT
        size,
        clientId,
        timeInForce,
        goodTilTimeInSeconds,
        execution,
        postOnly,
        reduceOnly,
        triggerPrice,
        marketInfo,
        currentHeight,
      );
      msg.then((it) => resolve([it])).catch((err) => {
        console.log(err);
      });
    });
    const orderFlags = calculateOrderFlags(type, timeInForce);
    const account: Promise<Account> = this.validatorClient.post.account(
      subaccount.address,
      orderFlags,
    );
    return this.send(
      subaccount.wallet,
      () => msgs,
      true,
      undefined,
      undefined,
      () => account,
    );
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
  private async placeOrderMessage(
    subaccount: SubaccountInfo,
    marketId: string,
    type: OrderType,
    side: OrderSide,
    price: number,
    // trigger_price: number,   // not used for MARKET and LIMIT
    size: number,
    clientId: number,
    timeInForce?: OrderTimeInForce,
    goodTilTimeInSeconds?: number,
    execution?: OrderExecution,
    postOnly?: boolean,
    reduceOnly?: boolean,
    triggerPrice?: number,
    marketInfo?: MarketInfo,
    currentHeight?: number,
  ): Promise<EncodeObject> {
    const orderFlags = calculateOrderFlags(type, timeInForce);

    const result = await Promise.all([
      this.calculateGoodTilBlock(orderFlags, currentHeight),
      this.retrieveMarketInfo(marketId, marketInfo),
    ],
    );
    const goodTilBlock = result[0];
    const clobPairId = result[1].clobPairId;
    const atomicResolution = result[1].atomicResolution;
    const stepBaseQuantums = result[1].stepBaseQuantums;
    const quantumConversionExponent = result[1].quantumConversionExponent;
    const subticksPerTick = result[1].subticksPerTick;
    const orderSide = calculateSide(side);
    const quantums = calculateQuantums(
      size,
      atomicResolution,
      stepBaseQuantums,
    );
    const subticks = calculateSubticks(
      price,
      atomicResolution,
      quantumConversionExponent,
      subticksPerTick,
    );
    const orderTimeInForce = calculateTimeInForce(type, timeInForce, execution, postOnly);
    let goodTilBlockTime = 0;
    if (orderFlags === OrderFlags.LONG_TERM || orderFlags === OrderFlags.CONDITIONAL) {
      if (goodTilTimeInSeconds == null) {
        throw new Error('goodTilTimeInSeconds must be set for LONG_TERM or CONDITIONAL order');
      } else {
        goodTilBlockTime = this.calculateGoodTilBlockTime(goodTilTimeInSeconds);
      }
    }
    const clientMetadata = calculateClientMetadata(type);
    const conditionalType = calculateConditionType(type);
    const conditionalOrderTriggerSubticks = calculateConditionalOrderTriggerSubticks(
      type,
      atomicResolution,
      quantumConversionExponent,
      subticksPerTick,
      triggerPrice);
    return this.validatorClient.post.composer.composeMsgPlaceOrder(
      subaccount.address,
      subaccount.subaccountNumber,
      clientId,
      clobPairId,
      orderFlags,
      goodTilBlock,
      goodTilBlockTime,
      orderSide,
      quantums,
      subticks,
      orderTimeInForce,
      reduceOnly ?? false,
      clientMetadata,
      conditionalType,
      conditionalOrderTriggerSubticks,
    );
  }

  private async retrieveMarketInfo(marketId: string, marketInfo?:MarketInfo): Promise<MarketInfo> {
    if (marketInfo) {
      return Promise.resolve(marketInfo);
    } else {
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
  private async placeShortTermOrderMessage(
    subaccount: SubaccountInfo,
    marketId: string,
    side: OrderSide,
    price: number,
    size: number,
    clientId: number,
    goodTilBlock: number,
    timeInForce: Order_TimeInForce,
    reduceOnly: boolean,
  ): Promise<EncodeObject> {
    await this.validateGoodTilBlock(goodTilBlock);

    const marketsResponse = await this.indexerClient.markets.getPerpetualMarkets(marketId);
    const market = marketsResponse.markets[marketId];
    const clobPairId = market.clobPairId;
    const atomicResolution = market.atomicResolution;
    const stepBaseQuantums = market.stepBaseQuantums;
    const quantumConversionExponent = market.quantumConversionExponent;
    const subticksPerTick = market.subticksPerTick;
    const orderSide = calculateSide(side);
    const quantums = calculateQuantums(
      size,
      atomicResolution,
      stepBaseQuantums,
    );
    const subticks = calculateSubticks(
      price,
      atomicResolution,
      quantumConversionExponent,
      subticksPerTick,
    );
    const orderFlags = OrderFlags.SHORT_TERM;
    return this.validatorClient.post.composer.composeMsgPlaceOrder(
      subaccount.address,
      subaccount.subaccountNumber,
      clientId,
      clobPairId,
      orderFlags,
      goodTilBlock,
      0, // Short term orders use goodTilBlock.
      orderSide,
      quantums,
      subticks,
      timeInForce,
      reduceOnly,
      0, // Client metadata is 0 for short term orders.
      Order_ConditionType.CONDITION_TYPE_UNSPECIFIED, // Short term orders cannot be conditional.
      Long.fromInt(0), // Short term orders cannot be conditional.
    );
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
  async cancelRawOrder(
    subaccount: SubaccountInfo,
    clientId: number,
    orderFlags: OrderFlags,
    clobPairId: number,
    goodTilBlock?: number,
    goodTilBlockTime?: number,
  ): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx> {
    return this.validatorClient.post.cancelOrder(
      subaccount,
      clientId,
      orderFlags,
      clobPairId,
      goodTilBlock,
      goodTilBlockTime,
    );
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
  async cancelOrder(
    subaccount: SubaccountInfo,
    clientId: number,
    orderFlags: OrderFlags,
    marketId: string,
    goodTilBlock?: number,
    goodTilTimeInSeconds?: number,
  ): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx> {

    const marketsResponse = await this.indexerClient.markets.getPerpetualMarkets(marketId);
    const market = marketsResponse.markets[marketId];
    const clobPairId = market.clobPairId;

    if (!verifyOrderFlags(orderFlags)) {
      throw new Error(`Invalid order flags: ${orderFlags}`);
    }

    let goodTilBlockTime;
    if (isStatefulOrder(orderFlags)) {
      if (goodTilTimeInSeconds === undefined || goodTilTimeInSeconds === 0) {
        throw new Error('goodTilTimeInSeconds must be set for LONG_TERM or CONDITIONAL order');
      }
      if (goodTilBlock !== 0) {
        throw new Error(
          'goodTilBlock should be zero since LONG_TERM or CONDITIONAL orders ' +
          'use goodTilTimeInSeconds instead of goodTilBlock.',
        );
      }
      goodTilBlockTime = this.calculateGoodTilBlockTime(goodTilTimeInSeconds);
    } else {
      if (goodTilBlock === undefined || goodTilBlock === 0) {
        throw new Error('goodTilBlock must be non-zero for SHORT_TERM orders');
      }
      if (goodTilTimeInSeconds !== undefined && goodTilTimeInSeconds !== 0) {
        throw new Error('goodTilTimeInSeconds should be zero since SHORT_TERM orders use goodTilBlock instead of goodTilTimeInSeconds.');
      }
    }

    return this.validatorClient.post.cancelOrder(
      subaccount,
      clientId,
      orderFlags,
      clobPairId,
      goodTilBlock,
      goodTilBlockTime,
    );
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
  async transferToSubaccount(
    subaccount: SubaccountInfo,
    recipientAddress: string,
    recipientSubaccountNumber: number,
    amount: string,
  ): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx> {
    const msgs: Promise<EncodeObject[]> = new Promise((resolve) => {
      const msg = this.transferToSubaccountMessage(
        subaccount,
        recipientAddress,
        recipientSubaccountNumber,
        amount,
      );
      resolve([msg]);
    });
    return this.send(
      subaccount.wallet,
      () => msgs,
      true);
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
  transferToSubaccountMessage(
    subaccount: SubaccountInfo,
    recipientAddress: string,
    recipientSubaccountNumber: number,
    amount: string,
  ): EncodeObject {
    const validatorClient = this._validatorClient;
    if (validatorClient === undefined) {
      throw new Error('validatorClient not set');
    }
    const quantums = parseUnits(amount, validatorClient.config.denoms.TDAI_DECIMALS);
    if (quantums > BigInt(Long.MAX_VALUE.toString())) {
      throw new Error('amount to large');
    }
    if (quantums < 0) {
      throw new Error('amount must be positive');
    }

    return this.validatorClient.post.composer.composeMsgTransfer(
      subaccount.address,
      subaccount.subaccountNumber,
      recipientAddress,
      recipientSubaccountNumber,
      0,
      Long.fromString(quantums.toString()),
    );
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
  async depositToSubaccount(
    subaccount: SubaccountInfo,
    amount: string,
  ): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx> {
    const msgs: Promise<EncodeObject[]> = new Promise((resolve) => {
      const msg = this.depositToSubaccountMessage(
        subaccount,
        amount,
      );
      resolve([msg]);
    });
    return this.validatorClient.post.send(subaccount.wallet,
      () => msgs,
      false);
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
  depositToSubaccountMessage(
    subaccount: SubaccountInfo,
    amount: string,
  ): EncodeObject {
    const validatorClient = this._validatorClient;
    if (validatorClient === undefined) {
      throw new Error('validatorClient not set');
    }
    const quantums = parseUnits(amount, validatorClient.config.denoms.TDAI_DECIMALS);
    if (quantums > BigInt(Long.MAX_VALUE.toString())) {
      throw new Error('amount to large');
    }
    if (quantums < 0) {
      throw new Error('amount must be positive');
    }

    return this.validatorClient.post.composer.composeMsgDepositToSubaccount(
      subaccount.address,
      subaccount.subaccountNumber,
      0,
      Long.fromString(quantums.toString()),
    );
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
  async withdrawFromSubaccount(
    subaccount: SubaccountInfo,
    amount: string,
    recipient?: string,
  ): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx> {
    const msgs: Promise<EncodeObject[]> = new Promise((resolve) => {
      const msg = this.withdrawFromSubaccountMessage(
        subaccount,
        amount,
        recipient,
      );
      resolve([msg]);
    });
    return this.send(
      subaccount.wallet,
      () => msgs,
      false);
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
  withdrawFromSubaccountMessage(
    subaccount: SubaccountInfo,
    amount: string,
    recipient?: string,
  ): EncodeObject {
    const validatorClient = this._validatorClient;
    if (validatorClient === undefined) {
      throw new Error('validatorClient not set');
    }
    const quantums = parseUnits(amount, validatorClient.config.denoms.TDAI_DECIMALS);
    if (quantums > BigInt(Long.MAX_VALUE.toString())) {
      throw new Error('amount to large');
    }
    if (quantums < 0) {
      throw new Error('amount must be positive');
    }

    return this.validatorClient.post.composer.composeMsgWithdrawFromSubaccount(
      subaccount.address,
      subaccount.subaccountNumber,
      0,
      Long.fromString(quantums.toString()),
      recipient,
    );
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
  sendTokenMessage(
    wallet: LocalWallet,
    amount: string,
    recipient: string,
  ): EncodeObject {
    const address = wallet.address;
    if (address === undefined) {
      throw new UserError('wallet address is not set. Call connectWallet() first');
    }
    const {
      CHAINTOKEN_DENOM: chainTokenDenom,
      CHAINTOKEN_DECIMALS: chainTokenDecimals,
    } = this._validatorClient?.config.denoms || {};

    if (chainTokenDenom === undefined || chainTokenDecimals === undefined) {
      throw new Error('Chain token denom not set in validator config');
    }

    const quantums = parseUnits(amount, chainTokenDecimals);

    return this.validatorClient.post.composer.composeMsgSendToken(
      address,
      recipient,
      chainTokenDenom,
      quantums.toString(),
    );
  }

  async signPlaceOrder(
    subaccount: SubaccountInfo,
    marketId: string,
    type: OrderType,
    side: OrderSide,
    price: number,
    // trigger_price: number,   // not used for MARKET and LIMIT
    size: number,
    clientId: number,
    timeInForce: OrderTimeInForce,
    goodTilTimeInSeconds: number,
    execution: OrderExecution,
    postOnly: boolean,
    reduceOnly: boolean,
  ): Promise<string> {
    const msgs: Promise<EncodeObject[]> = new Promise((resolve) => {
      const msg = this.placeOrderMessage(
        subaccount,
        marketId,
        type,
        side,
        price,
        // trigger_price: number,   // not used for MARKET and LIMIT
        size,
        clientId,
        timeInForce,
        goodTilTimeInSeconds,
        execution,
        postOnly,
        reduceOnly,
      );
      msg.then((it) => resolve([it])).catch((err) => {
        console.log(err);
      });
    });
    const signature = await this.sign(
      wallet,
      () => msgs,
      true,
    );

    return Buffer.from(signature).toString('base64');
  }

  async signCancelOrder(
    subaccount: SubaccountInfo,
    clientId: number,
    orderFlags: OrderFlags,
    clobPairId: number,
    goodTilBlock: number,
    goodTilBlockTime: number,
  ): Promise<string> {
    const msgs: Promise<EncodeObject[]> = new Promise((resolve) => {
      const msg = this.validatorClient.post.composer.composeMsgCancelOrder(
        subaccount.address,
        subaccount.subaccountNumber,
        clientId,
        clobPairId,
        orderFlags,
        goodTilBlock,
        goodTilBlockTime,
      );
      resolve([msg]);
    });
    const signature = await this.sign(
      subaccount.wallet,
      () => msgs,
      true,
    );

    return Buffer.from(signature).toString('base64');
  }
}
