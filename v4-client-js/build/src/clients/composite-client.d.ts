import { EncodeObject } from '@cosmjs/proto-signing';
import { Account, GasPrice, IndexedTx, StdFee } from '@cosmjs/stargate';
import { BroadcastTxAsyncResponse, BroadcastTxSyncResponse } from '@cosmjs/tendermint-rpc/build/tendermint37';
import { Order_TimeInForce } from '@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/order';
import { OrderFlags } from '../types';
import { Network, OrderExecution, OrderSide, OrderTimeInForce, OrderType } from './constants';
import { IndexerClient } from './indexer-client';
import LocalWallet from './modules/local-wallet';
import { SubaccountInfo } from './subaccount';
import { ValidatorClient } from './validator-client';
export interface MarketInfo {
    clobPairId: number;
    atomicResolution: number;
    stepBaseQuantums: number;
    quantumConversionExponent: number;
    subticksPerTick: number;
}
export declare class CompositeClient {
    readonly network: Network;
    private _indexerClient;
    private _validatorClient?;
    static connect(network: Network): Promise<CompositeClient>;
    private constructor();
    private initialize;
    get indexerClient(): IndexerClient;
    get validatorClient(): ValidatorClient;
    /**
       * @description Sign a list of messages with a wallet.
       * the calling function is responsible for creating the messages.
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The Signature.
       */
    sign(wallet: LocalWallet, messaging: () => Promise<EncodeObject[]>, zeroFee: boolean, gasPrice?: GasPrice, memo?: string, account?: () => Promise<Account>): Promise<Uint8Array>;
    /**
       * @description Send a list of messages with a wallet.
       * the calling function is responsible for creating the messages.
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The Transaction Hash.
       */
    send(wallet: LocalWallet, messaging: () => Promise<EncodeObject[]>, zeroFee: boolean, gasPrice?: GasPrice, memo?: string, account?: () => Promise<Account>): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
    /**
       * @description Send a signed transaction.
       *
       * @param signedTransaction The signed transaction to send.
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The Transaction Hash.
       */
    sendSignedTransaction(signedTransaction: Uint8Array): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
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
    simulate(wallet: LocalWallet, messaging: () => Promise<EncodeObject[]>, gasPrice?: GasPrice, memo?: string, account?: () => Promise<Account>): Promise<StdFee>;
    /**
       * @description Calculate the goodTilBlock value for a SHORT_TERM order
       *
       * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
       * at any point.
       * @returns The goodTilBlock value
       */
    private calculateGoodTilBlock;
    /**
     * @description Validate the goodTilBlock value for a SHORT_TERM order
     *
     * @param goodTilBlock Number of blocks from the current block height the order will
     * be valid for.
     *
     * @throws UserError if the goodTilBlock value is not valid given latest block height and
     * SHORT_BLOCK_WINDOW.
     */
    private validateGoodTilBlock;
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
    private calculateGoodTilBlockTime;
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
    placeShortTermOrder(subaccount: SubaccountInfo, marketId: string, side: OrderSide, price: number, size: number, clientId: number, goodTilBlock: number, timeInForce: Order_TimeInForce, reduceOnly: boolean): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
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
    placeOrder(subaccount: SubaccountInfo, marketId: string, type: OrderType, side: OrderSide, price: number, size: number, clientId: number, timeInForce?: OrderTimeInForce, goodTilTimeInSeconds?: number, execution?: OrderExecution, postOnly?: boolean, reduceOnly?: boolean, triggerPrice?: number, marketInfo?: MarketInfo, currentHeight?: number): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
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
    private placeOrderMessage;
    private retrieveMarketInfo;
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
    private placeShortTermOrderMessage;
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
    cancelRawOrder(subaccount: SubaccountInfo, clientId: number, orderFlags: OrderFlags, clobPairId: number, goodTilBlock?: number, goodTilBlockTime?: number): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
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
    cancelOrder(subaccount: SubaccountInfo, clientId: number, orderFlags: OrderFlags, marketId: string, goodTilBlock?: number, goodTilTimeInSeconds?: number): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
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
    transferToSubaccount(subaccount: SubaccountInfo, recipientAddress: string, recipientSubaccountNumber: number, amount: string): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
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
    transferToSubaccountMessage(subaccount: SubaccountInfo, recipientAddress: string, recipientSubaccountNumber: number, amount: string): EncodeObject;
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
    depositToSubaccount(subaccount: SubaccountInfo, amount: string): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
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
    depositToSubaccountMessage(subaccount: SubaccountInfo, amount: string): EncodeObject;
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
    withdrawFromSubaccount(subaccount: SubaccountInfo, amount: string, recipient?: string): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
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
    withdrawFromSubaccountMessage(subaccount: SubaccountInfo, amount: string, recipient?: string): EncodeObject;
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
    sendTokenMessage(wallet: LocalWallet, amount: string, recipient: string): EncodeObject;
    signPlaceOrder(subaccount: SubaccountInfo, marketId: string, type: OrderType, side: OrderSide, price: number, size: number, clientId: number, timeInForce: OrderTimeInForce, goodTilTimeInSeconds: number, execution: OrderExecution, postOnly: boolean, reduceOnly: boolean): Promise<string>;
    signCancelOrder(subaccount: SubaccountInfo, clientId: number, orderFlags: OrderFlags, clobPairId: number, goodTilBlock: number, goodTilBlockTime: number): Promise<string>;
}
