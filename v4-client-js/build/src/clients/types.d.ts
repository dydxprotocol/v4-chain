import { Method } from '@cosmjs/tendermint-rpc';
import { Order_ConditionType, Order_Side, Order_TimeInForce } from '@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/order';
import BigNumber from 'bignumber.js';
import Long from 'long';
export type Integer = BigNumber;
export type GenericParams = {
    [name: string]: any;
};
export type Data = any;
export type PartialBy<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>;
export interface PartialTransactionOptions {
    accountNumber: number;
    chainId: string;
}
export interface TransactionOptions extends PartialTransactionOptions {
    sequence: number;
}
export declare enum OrderFlags {
    SHORT_TERM = 0,
    LONG_TERM = 64,
    CONDITIONAL = 32
}
export interface IBasicOrder {
    clientId: number;
    orderFlags: OrderFlags;
    clobPairId: number;
    goodTilBlock?: number;
    goodTilBlockTime?: number;
}
export interface IPlaceOrder extends IBasicOrder {
    side: Order_Side;
    quantums: Long;
    subticks: Long;
    timeInForce: Order_TimeInForce;
    reduceOnly: boolean;
    clientMetadata: number;
    conditionType?: Order_ConditionType;
    conditionalOrderTriggerSubticks?: Long;
}
export interface ICancelOrder extends IBasicOrder {
}
export interface BroadcastOptions {
    broadcastPollIntervalMs: number;
    broadcastTimeoutMs: number;
}
export interface DenomConfig {
    TDAI_DENOM: string;
    TDAI_DECIMALS: number;
    TDAI_GAS_DENOM?: string;
    CHAINTOKEN_DENOM: string;
    CHAINTOKEN_DECIMALS: number;
    CHAINTOKEN_GAS_DENOM?: string;
}
export type BroadcastMode = (Method.BroadcastTxAsync | Method.BroadcastTxSync | Method.BroadcastTxCommit);
export interface TimeResponse {
    iso: string;
    epoch: number;
}
export interface HeightResponse {
    height: number;
    time: string;
}
export interface ComplianceResponse {
    restricted: boolean;
}
export * from './modules/proto-includes';
