import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * BlockMessageIds stores the id of each message that should be processed at a
 * given block height.
 */
export interface BlockMessageIds {
    /**
     * ids stores a list of DelayedMessage ids that should be processed at a given
     * block height.
     */
    ids: number[];
}
/**
 * BlockMessageIds stores the id of each message that should be processed at a
 * given block height.
 */
export interface BlockMessageIdsSDKType {
    ids: number[];
}
export declare const BlockMessageIds: {
    encode(message: BlockMessageIds, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): BlockMessageIds;
    fromPartial(object: DeepPartial<BlockMessageIds>): BlockMessageIds;
};
