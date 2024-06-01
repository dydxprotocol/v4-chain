import { Block, IndexedTx } from '@cosmjs/stargate';
import { Tendermint37Client } from '@cosmjs/tendermint-rpc';
import { BroadcastTxAsyncResponse, BroadcastTxSyncResponse } from '@cosmjs/tendermint-rpc/build/tendermint37';
import { BroadcastMode, BroadcastOptions } from '../types';
export declare class TendermintClient {
    readonly baseClient: Tendermint37Client;
    broadcastOptions: BroadcastOptions;
    constructor(baseClient: Tendermint37Client, broadcastOptions: BroadcastOptions);
    /**
     * @description Get a specific block if height is specified. Otherwise, get the most recent block.
     *
     * @returns Information about the block queried.
     */
    getBlock(height?: number): Promise<Block>;
    /**
      * @description Broadcast a signed transaction with a specific mode.
      * @throws BroadcastErrorObject when result code is not zero. TypeError when mode is invalid.
      * @returns Differs depending on the BroadcastMode used.
      * See https://docs.cosmos.network/master/run-node/txs.html for more information.
      */
    broadcastTransaction(tx: Uint8Array, mode: BroadcastMode): Promise<BroadcastTxAsyncResponse | BroadcastTxSyncResponse | IndexedTx>;
    /**
     * @description Broadcast a signed transaction.
     * @returns The transaction hash.
     */
    broadcastTransactionAsync(tx: Uint8Array): Promise<BroadcastTxAsyncResponse>;
    /**
     * @description Broadcast a signed transaction and await the response.
     * @throws BroadcastErrorObject when result code is not zero.
     * @returns The response from the node once the transaction is processed by `CheckTx`.
     */
    broadcastTransactionSync(tx: Uint8Array): Promise<BroadcastTxSyncResponse>;
    /**
     * @description Broadcast a signed transaction and await for it to be included in the blockchain.
     * @throws BroadcastErrorObject when result code is not zero.
     * @returns The result of the transaction once included in the blockchain.
     */
    broadcastTransactionCommit(tx: Uint8Array): Promise<IndexedTx>;
    /**
     * @description Using tx method, query for a transaction on-chain with retries specified by
     * the client BroadcastOptions.
     *
     * @throws TimeoutError if the transaction is not committed on-chain within the timeout limit.
     * @returns An indexed transaction containing information about the transaction when committed.
     */
    queryHash(hash: Uint8Array, time?: number): Promise<IndexedTx>;
    /**
     * @description Set the broadcast options for this module.
     */
    setBroadcastOptions(broadcastOptions: BroadcastOptions): void;
}
