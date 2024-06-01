import { StdFee } from '@cosmjs/amino';
import { Method } from '@cosmjs/tendermint-rpc';
import { TransactionOptions } from './clients/types';
export * from './clients/types';
export * from './clients/constants';
export interface BroadcastOptions {
    broadcastPollIntervalMs: number;
    broadcastTimeoutMs: number;
}
export interface ApiOptions {
    faucetHost?: string;
    indexerHost?: string;
    timeout?: number;
}
export type BroadcastMode = (Method.BroadcastTxAsync | Method.BroadcastTxSync | Method.BroadcastTxCommit);
export interface Options {
    transactionOptions?: TransactionOptions;
    memo?: string;
    broadcastMode?: BroadcastMode;
    fee?: StdFee;
}
export declare enum ClobPairId {
    PERPETUAL_PAIR_BTC_USD = 0
}
