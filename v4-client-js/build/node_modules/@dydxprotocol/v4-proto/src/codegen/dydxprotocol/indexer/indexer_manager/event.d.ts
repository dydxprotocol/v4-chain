import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/** enum to specify that the IndexerTendermintEvent is a block event. */
export declare enum IndexerTendermintEvent_BlockEvent {
    /** BLOCK_EVENT_UNSPECIFIED - Default value. This value is invalid and unused. */
    BLOCK_EVENT_UNSPECIFIED = 0,
    /**
     * BLOCK_EVENT_BEGIN_BLOCK - BLOCK_EVENT_BEGIN_BLOCK indicates that the event was generated during
     * BeginBlock.
     */
    BLOCK_EVENT_BEGIN_BLOCK = 1,
    /**
     * BLOCK_EVENT_END_BLOCK - BLOCK_EVENT_END_BLOCK indicates that the event was generated during
     * EndBlock.
     */
    BLOCK_EVENT_END_BLOCK = 2,
    UNRECOGNIZED = -1
}
export declare const IndexerTendermintEvent_BlockEventSDKType: typeof IndexerTendermintEvent_BlockEvent;
export declare function indexerTendermintEvent_BlockEventFromJSON(object: any): IndexerTendermintEvent_BlockEvent;
export declare function indexerTendermintEvent_BlockEventToJSON(object: IndexerTendermintEvent_BlockEvent): string;
/**
 * IndexerTendermintEventWrapper is a wrapper around IndexerTendermintEvent,
 * with an additional txn_hash field.
 */
export interface IndexerTendermintEventWrapper {
    event?: IndexerTendermintEvent;
    txnHash: string;
}
/**
 * IndexerTendermintEventWrapper is a wrapper around IndexerTendermintEvent,
 * with an additional txn_hash field.
 */
export interface IndexerTendermintEventWrapperSDKType {
    event?: IndexerTendermintEventSDKType;
    txn_hash: string;
}
/**
 * IndexerEventsStoreValue represents the type of the value of the
 * `IndexerEventsStore` in state.
 */
export interface IndexerEventsStoreValue {
    events: IndexerTendermintEventWrapper[];
}
/**
 * IndexerEventsStoreValue represents the type of the value of the
 * `IndexerEventsStore` in state.
 */
export interface IndexerEventsStoreValueSDKType {
    events: IndexerTendermintEventWrapperSDKType[];
}
/**
 * IndexerTendermintEvent contains the base64 encoded event proto emitted from
 * the dYdX application as well as additional metadata to determine the ordering
 * of the event within the block and the subtype of the event.
 */
export interface IndexerTendermintEvent {
    /** Subtype of the event e.g. "order_fill", "subaccount_update", etc. */
    subtype: string;
    transactionIndex?: number;
    blockEvent?: IndexerTendermintEvent_BlockEvent;
    /**
     * Index of the event within the list of events that happened either during a
     * transaction or during processing of a block.
     * TODO(DEC-537): Deprecate this field because events are already ordered.
     */
    eventIndex: number;
    /** Version of the event. */
    version: number;
    /** Tendermint event bytes. */
    dataBytes: Uint8Array;
}
/**
 * IndexerTendermintEvent contains the base64 encoded event proto emitted from
 * the dYdX application as well as additional metadata to determine the ordering
 * of the event within the block and the subtype of the event.
 */
export interface IndexerTendermintEventSDKType {
    subtype: string;
    transaction_index?: number;
    block_event?: IndexerTendermintEvent_BlockEvent;
    event_index: number;
    version: number;
    data_bytes: Uint8Array;
}
/**
 * IndexerTendermintBlock contains all the events for the block along with
 * metadata for the block height, timestamp of the block and a list of all the
 * hashes of the transactions within the block. The transaction hashes follow
 * the ordering of the transactions as they appear within the block.
 */
export interface IndexerTendermintBlock {
    height: number;
    time?: Date;
    events: IndexerTendermintEvent[];
    txHashes: string[];
}
/**
 * IndexerTendermintBlock contains all the events for the block along with
 * metadata for the block height, timestamp of the block and a list of all the
 * hashes of the transactions within the block. The transaction hashes follow
 * the ordering of the transactions as they appear within the block.
 */
export interface IndexerTendermintBlockSDKType {
    height: number;
    time?: Date;
    events: IndexerTendermintEventSDKType[];
    tx_hashes: string[];
}
export declare const IndexerTendermintEventWrapper: {
    encode(message: IndexerTendermintEventWrapper, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): IndexerTendermintEventWrapper;
    fromPartial(object: DeepPartial<IndexerTendermintEventWrapper>): IndexerTendermintEventWrapper;
};
export declare const IndexerEventsStoreValue: {
    encode(message: IndexerEventsStoreValue, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): IndexerEventsStoreValue;
    fromPartial(object: DeepPartial<IndexerEventsStoreValue>): IndexerEventsStoreValue;
};
export declare const IndexerTendermintEvent: {
    encode(message: IndexerTendermintEvent, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): IndexerTendermintEvent;
    fromPartial(object: DeepPartial<IndexerTendermintEvent>): IndexerTendermintEvent;
};
export declare const IndexerTendermintBlock: {
    encode(message: IndexerTendermintBlock, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): IndexerTendermintBlock;
    fromPartial(object: DeepPartial<IndexerTendermintBlock>): IndexerTendermintBlock;
};
