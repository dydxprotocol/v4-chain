import { EncodeObject } from '@cosmjs/proto-signing';
import { Account, GasPrice, IndexedTx, StdFee } from '@cosmjs/stargate';
import { BroadcastTxAsyncResponse, BroadcastTxSyncResponse } from '@cosmjs/tendermint-rpc/build/tendermint37';
import Long from 'long';
import { SubaccountInfo } from '../subaccount';
import { OrderFlags, BroadcastMode, IPlaceOrder, ICancelOrder, DenomConfig } from '../types';
import { Composer } from './composer';
import { Get } from './get';
import LocalWallet from './local-wallet';
import { Order_Side, Order_TimeInForce, Order_ConditionType } from './proto-includes';
export declare class Post {
    readonly composer: Composer;
    private readonly registry;
    private readonly chainId;
    readonly get: Get;
    readonly denoms: DenomConfig;
    readonly defaultGasPrice: GasPrice;
    readonly defaultDydxGasPrice: GasPrice;
    private accountNumberCache;
    constructor(get: Get, chainId: string, denoms: DenomConfig);
    /**
     * @description Simulate a transaction
     * the calling function is responsible for creating the messages.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The Fee for broadcasting a transaction.
     */
    simulate(wallet: LocalWallet, messaging: () => Promise<EncodeObject[]>, gasPrice?: GasPrice, memo?: string, account?: () => Promise<Account>): Promise<StdFee>;
    /**
     * @description Sign a transaction
     * the calling function is responsible for creating the messages.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The Signature.
     */
    sign(wallet: LocalWallet, messaging: () => Promise<EncodeObject[]>, zeroFee: boolean, gasPrice?: GasPrice, memo?: string, account?: () => Promise<Account>): Promise<Uint8Array>;
    /**
     * @description Send a transaction
     * the calling function is responsible for creating the messages.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The Tx Hash.
     */
    send(wallet: LocalWallet, messaging: () => Promise<EncodeObject[]>, zeroFee: boolean, gasPrice?: GasPrice, memo?: string, broadcastMode?: BroadcastMode, account?: () => Promise<Account>): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
    /**
     * @description Calculate the default broadcast mode.
     */
    private defaultBroadcastMode;
    /**
     * @description Sign and send a message
     *
     * @returns The Tx Response.
     */
    private signTransaction;
    /**
     * @description Retrieve an account structure for transactions.
     * For short term orders, the sequence doesn't matter. Use cached if available.
     * For long term and conditional orders, a round trip to validator must be made.
     */
    account(address: string, orderFlags?: OrderFlags): Promise<Account>;
    /**
     * @description Sign and send a message
     *
     * @returns The Tx Response.
     */
    private signAndSendTransaction;
    /**
     * @description Send signed transaction.
     *
     * @returns The Tx Response.
     */
    sendSignedTransaction(signedTransaction: Uint8Array, broadcastMode?: BroadcastMode): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
    /**
     * @description Simulate broadcasting a transaction.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error
     * at any point.
     * @returns The Fee for broadcasting a transaction.
     */
    private simulateTransaction;
    placeOrder(subaccount: SubaccountInfo, clientId: number, clobPairId: number, side: Order_Side, quantums: Long, subticks: Long, timeInForce: Order_TimeInForce, orderFlags: number, reduceOnly: boolean, goodTilBlock?: number, goodTilBlockTime?: number, clientMetadata?: number, conditionType?: Order_ConditionType, conditionalOrderTriggerSubticks?: Long, broadcastMode?: BroadcastMode): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
    placeOrderObject(subaccount: SubaccountInfo, placeOrder: IPlaceOrder, broadcastMode?: BroadcastMode): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
    cancelOrder(subaccount: SubaccountInfo, clientId: number, orderFlags: OrderFlags, clobPairId: number, goodTilBlock?: number, goodTilBlockTime?: number, broadcastMode?: BroadcastMode): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
    cancelOrderObject(subaccount: SubaccountInfo, cancelOrder: ICancelOrder, broadcastMode?: BroadcastMode): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
    transfer(subaccount: SubaccountInfo, recipientAddress: string, recipientSubaccountNumber: number, assetId: number, amount: Long, broadcastMode?: BroadcastMode): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
    deposit(subaccount: SubaccountInfo, assetId: number, quantums: Long, broadcastMode?: BroadcastMode): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
    withdraw(subaccount: SubaccountInfo, assetId: number, quantums: Long, recipient?: string, broadcastMode?: BroadcastMode): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
    sendToken(subaccount: SubaccountInfo, recipient: string, coinDenom: string, quantums: string, zeroFee?: boolean, broadcastMode?: BroadcastMode): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
}
