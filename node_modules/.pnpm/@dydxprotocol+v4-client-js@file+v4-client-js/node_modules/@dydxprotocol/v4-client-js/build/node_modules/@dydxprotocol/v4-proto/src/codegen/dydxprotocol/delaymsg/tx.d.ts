/// <reference types="long" />
import { Any, AnySDKType } from "../../google/protobuf/any";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** MsgDelayMessage is a request type for the DelayMessage method. */
export interface MsgDelayMessage {
    authority: string;
    /** The message to be delayed. */
    msg?: Any;
    /** The number of blocks to delay the message for. */
    delayBlocks: number;
}
/** MsgDelayMessage is a request type for the DelayMessage method. */
export interface MsgDelayMessageSDKType {
    authority: string;
    msg?: AnySDKType;
    delay_blocks: number;
}
/** MsgDelayMessageResponse is a response type for the DelayMessage method. */
export interface MsgDelayMessageResponse {
    /** The id of the created delayed message. */
    id: Long;
}
/** MsgDelayMessageResponse is a response type for the DelayMessage method. */
export interface MsgDelayMessageResponseSDKType {
    id: Long;
}
export declare const MsgDelayMessage: {
    encode(message: MsgDelayMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDelayMessage;
    fromPartial(object: DeepPartial<MsgDelayMessage>): MsgDelayMessage;
};
export declare const MsgDelayMessageResponse: {
    encode(message: MsgDelayMessageResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDelayMessageResponse;
    fromPartial(object: DeepPartial<MsgDelayMessageResponse>): MsgDelayMessageResponse;
};
