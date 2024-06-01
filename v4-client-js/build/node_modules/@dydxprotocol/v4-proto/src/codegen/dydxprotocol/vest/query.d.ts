import { VestEntry, VestEntrySDKType } from "./vest_entry";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QueryVestEntryRequest is a request type for the VestEntry RPC method. */
export interface QueryVestEntryRequest {
    /** QueryVestEntryRequest is a request type for the VestEntry RPC method. */
    vesterAccount: string;
}
/** QueryVestEntryRequest is a request type for the VestEntry RPC method. */
export interface QueryVestEntryRequestSDKType {
    vester_account: string;
}
/** QueryVestEntryResponse is a response type for the VestEntry RPC method. */
export interface QueryVestEntryResponse {
    entry?: VestEntry;
}
/** QueryVestEntryResponse is a response type for the VestEntry RPC method. */
export interface QueryVestEntryResponseSDKType {
    entry?: VestEntrySDKType;
}
export declare const QueryVestEntryRequest: {
    encode(message: QueryVestEntryRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryVestEntryRequest;
    fromPartial(object: DeepPartial<QueryVestEntryRequest>): QueryVestEntryRequest;
};
export declare const QueryVestEntryResponse: {
    encode(message: QueryVestEntryResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryVestEntryResponse;
    fromPartial(object: DeepPartial<QueryVestEntryResponse>): QueryVestEntryResponse;
};
