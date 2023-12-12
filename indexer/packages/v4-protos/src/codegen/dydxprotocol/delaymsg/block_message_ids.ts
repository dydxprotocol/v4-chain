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
  /**
   * ids stores a list of DelayedMessage ids that should be processed at a given
   * block height.
   */
  ids: number[];
}

function createBaseBlockMessageIds(): BlockMessageIds {
  return {
    ids: []
  };
}

export const BlockMessageIds = {
  encode(message: BlockMessageIds, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    writer.uint32(10).fork();

    for (const v of message.ids) {
      writer.uint32(v);
    }

    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BlockMessageIds {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockMessageIds();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;

            while (reader.pos < end2) {
              message.ids.push(reader.uint32());
            }
          } else {
            message.ids.push(reader.uint32());
          }

          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<BlockMessageIds>): BlockMessageIds {
    const message = createBaseBlockMessageIds();
    message.ids = object.ids?.map(e => e) || [];
    return message;
  }

};