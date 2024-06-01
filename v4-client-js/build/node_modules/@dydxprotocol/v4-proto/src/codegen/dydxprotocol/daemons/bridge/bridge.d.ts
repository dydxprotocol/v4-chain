import { BridgeEvent, BridgeEventSDKType } from "../../bridge/bridge_event";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/**
 * AddBridgeEventsRequest is a request message that contains a list of new
 * bridge events. The events should be contiguous and sorted by (unique) id.
 */
export interface AddBridgeEventsRequest {
    bridgeEvents: BridgeEvent[];
}
/**
 * AddBridgeEventsRequest is a request message that contains a list of new
 * bridge events. The events should be contiguous and sorted by (unique) id.
 */
export interface AddBridgeEventsRequestSDKType {
    bridge_events: BridgeEventSDKType[];
}
/** AddBridgeEventsResponse is a response message for BridgeEventRequest. */
export interface AddBridgeEventsResponse {
}
/** AddBridgeEventsResponse is a response message for BridgeEventRequest. */
export interface AddBridgeEventsResponseSDKType {
}
export declare const AddBridgeEventsRequest: {
    encode(message: AddBridgeEventsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): AddBridgeEventsRequest;
    fromPartial(object: DeepPartial<AddBridgeEventsRequest>): AddBridgeEventsRequest;
};
export declare const AddBridgeEventsResponse: {
    encode(_: AddBridgeEventsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): AddBridgeEventsResponse;
    fromPartial(_: DeepPartial<AddBridgeEventsResponse>): AddBridgeEventsResponse;
};
