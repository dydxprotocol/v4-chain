import * as _m0 from "protobufjs/minimal";
import { Long, DeepPartial } from "../../helpers";
/**
 * PendingSendPacket contains the channel_id and sequence pair to identify a
 * pending packet
 */

export interface PendingSendPacket {
  channelId: string;
  sequence: Long;
}
/**
 * PendingSendPacket contains the channel_id and sequence pair to identify a
 * pending packet
 */

export interface PendingSendPacketSDKType {
  channel_id: string;
  sequence: Long;
}

function createBasePendingSendPacket(): PendingSendPacket {
  return {
    channelId: "",
    sequence: Long.UZERO
  };
}

export const PendingSendPacket = {
  encode(message: PendingSendPacket, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }

    if (!message.sequence.isZero()) {
      writer.uint32(16).uint64(message.sequence);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PendingSendPacket {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePendingSendPacket();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;

        case 2:
          message.sequence = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<PendingSendPacket>): PendingSendPacket {
    const message = createBasePendingSendPacket();
    message.channelId = object.channelId ?? "";
    message.sequence = object.sequence !== undefined && object.sequence !== null ? Long.fromValue(object.sequence) : Long.UZERO;
    return message;
  }

};