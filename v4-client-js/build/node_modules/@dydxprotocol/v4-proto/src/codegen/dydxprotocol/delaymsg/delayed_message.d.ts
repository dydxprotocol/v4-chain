import { Any, AnySDKType } from "../../google/protobuf/any";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** DelayedMessage is a message that is delayed until a certain block height. */
export interface DelayedMessage {
    /** The ID of the delayed message. */
    id: number;
    /** The message to be executed. */
    msg?: Any;
    /** The block height at which the message should be executed. */
    blockHeight: number;
}
/** DelayedMessage is a message that is delayed until a certain block height. */
export interface DelayedMessageSDKType {
    id: number;
    msg?: AnySDKType;
    block_height: number;
}
export declare const DelayedMessage: {
    encode(message: DelayedMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DelayedMessage;
    fromPartial(object: DeepPartial<DelayedMessage>): DelayedMessage;
};
