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
  /** delayed_messages is a list of delayed messages. */
  delayed_messages: DelayedMessageSDKType[];
  /** next_delayed_message_id is the id to be assigned to next delayed message. */

  next_delayed_message_id: number;
}

function createBaseGenesisState(): GenesisState {
  return {
    delayedMessages: [],
    nextDelayedMessageId: 0
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.delayedMessages) {
      DelayedMessage.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.nextDelayedMessageId !== 0) {
      writer.uint32(16).uint32(message.nextDelayedMessageId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.delayedMessages.push(DelayedMessage.decode(reader, reader.uint32()));
          break;

        case 2:
          message.nextDelayedMessageId = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.delayedMessages = object.delayedMessages?.map(e => DelayedMessage.fromPartial(e)) || [];
    message.nextDelayedMessageId = object.nextDelayedMessageId ?? 0;
    return message;
  }

};