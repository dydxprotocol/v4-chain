import { ResponseCommit, ResponseCommitSDKType, RequestFinalizeBlock, RequestFinalizeBlockSDKType, ResponseFinalizeBlock, ResponseFinalizeBlockSDKType } from "../../../tendermint/abci/types";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/**
 * StoreKVPair is a KVStore KVPair used for listening to state changes (Sets and Deletes)
 * It optionally includes the StoreKey for the originating KVStore and a Boolean flag to distinguish between Sets and
 * Deletes
 *
 * Since: cosmos-sdk 0.43
 */
export interface StoreKVPair {
    /** the store key for the KVStore this pair originates from */
    storeKey: string;
    /** true indicates a delete operation, false indicates a set operation */
    delete: boolean;
    key: Uint8Array;
    value: Uint8Array;
}
/**
 * StoreKVPair is a KVStore KVPair used for listening to state changes (Sets and Deletes)
 * It optionally includes the StoreKey for the originating KVStore and a Boolean flag to distinguish between Sets and
 * Deletes
 *
 * Since: cosmos-sdk 0.43
 */
export interface StoreKVPairSDKType {
    store_key: string;
    delete: boolean;
    key: Uint8Array;
    value: Uint8Array;
}
/**
 * BlockMetadata contains all the abci event data of a block
 * the file streamer dump them into files together with the state changes.
 */
export interface BlockMetadata {
    responseCommit?: ResponseCommit;
    requestFinalizeBlock?: RequestFinalizeBlock;
    /** TODO: should we renumber this? */
    responseFinalizeBlock?: ResponseFinalizeBlock;
}
/**
 * BlockMetadata contains all the abci event data of a block
 * the file streamer dump them into files together with the state changes.
 */
export interface BlockMetadataSDKType {
    response_commit?: ResponseCommitSDKType;
    request_finalize_block?: RequestFinalizeBlockSDKType;
    response_finalize_block?: ResponseFinalizeBlockSDKType;
}
export declare const StoreKVPair: {
    encode(message: StoreKVPair, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): StoreKVPair;
    fromPartial(object: DeepPartial<StoreKVPair>): StoreKVPair;
};
export declare const BlockMetadata: {
    encode(message: BlockMetadata, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): BlockMetadata;
    fromPartial(object: DeepPartial<BlockMetadata>): BlockMetadata;
};
