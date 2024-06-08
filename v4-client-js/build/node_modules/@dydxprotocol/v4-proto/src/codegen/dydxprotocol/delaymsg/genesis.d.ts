import { DelayedMessage, DelayedMessageSDKType } from "./delayed_message";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the delaymsg module's genesis state. */
export interface GenesisState {
    /** delayed_messages is a list of delayed messages. */
    delayedMessages: DelayedMessage[];
    /** next_delayed_message_id is the id to be assigned to next delayed message. */
    nextDelayedMessageId: number;
}
/** GenesisState defines the delaymsg module's genesis state. */
export interface GenesisStateSDKType {
    delayed_messages: DelayedMessageSDKType[];
    next_delayed_message_id: number;
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
