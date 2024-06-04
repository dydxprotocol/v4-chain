import { Duration, DurationSDKType } from "../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** BlockInfo stores information about a block */
export interface BlockInfo {
    height: number;
    timestamp?: Date;
}
/** BlockInfo stores information about a block */
export interface BlockInfoSDKType {
    height: number;
    timestamp?: Date;
}
/** AllDowntimeInfo stores information for all downtime durations. */
export interface AllDowntimeInfo {
    /**
     * The downtime information for each tracked duration. Sorted by duration,
     * ascending. (i.e. the same order as they appear in DowntimeParams).
     */
    infos: AllDowntimeInfo_DowntimeInfo[];
}
/** AllDowntimeInfo stores information for all downtime durations. */
export interface AllDowntimeInfoSDKType {
    infos: AllDowntimeInfo_DowntimeInfoSDKType[];
}
/**
 * Stores information about downtime. block_info corresponds to the most
 * recent block at which a downtime occurred.
 */
export interface AllDowntimeInfo_DowntimeInfo {
    duration?: Duration;
    blockInfo?: BlockInfo;
}
/**
 * Stores information about downtime. block_info corresponds to the most
 * recent block at which a downtime occurred.
 */
export interface AllDowntimeInfo_DowntimeInfoSDKType {
    duration?: DurationSDKType;
    block_info?: BlockInfoSDKType;
}
export declare const BlockInfo: {
    encode(message: BlockInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): BlockInfo;
    fromPartial(object: DeepPartial<BlockInfo>): BlockInfo;
};
export declare const AllDowntimeInfo: {
    encode(message: AllDowntimeInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): AllDowntimeInfo;
    fromPartial(object: DeepPartial<AllDowntimeInfo>): AllDowntimeInfo;
};
export declare const AllDowntimeInfo_DowntimeInfo: {
    encode(message: AllDowntimeInfo_DowntimeInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): AllDowntimeInfo_DowntimeInfo;
    fromPartial(object: DeepPartial<AllDowntimeInfo_DowntimeInfo>): AllDowntimeInfo_DowntimeInfo;
};
