import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * RewardShare stores the relative weight of rewards that each address is
 * entitled to.
 */
export interface RewardShare {
    address: string;
    weight: Uint8Array;
}
/**
 * RewardShare stores the relative weight of rewards that each address is
 * entitled to.
 */
export interface RewardShareSDKType {
    address: string;
    weight: Uint8Array;
}
export declare const RewardShare: {
    encode(message: RewardShare, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): RewardShare;
    fromPartial(object: DeepPartial<RewardShare>): RewardShare;
};
