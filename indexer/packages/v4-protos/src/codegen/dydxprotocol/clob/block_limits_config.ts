import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** BlockLimitsConfig stores global per-block limits for the CLOB module. */

export interface BlockLimitsConfig {
  /**
   * The maximum number of expired stateful orders that can be removed from
   * state in a single block. This prevents performance degradation when
   * processing a large number of expired orders.
   */
  maxStatefulOrderRemovalsPerBlock: number;
}
/** BlockLimitsConfig stores global per-block limits for the CLOB module. */

export interface BlockLimitsConfigSDKType {
  /**
   * The maximum number of expired stateful orders that can be removed from
   * state in a single block. This prevents performance degradation when
   * processing a large number of expired orders.
   */
  max_stateful_order_removals_per_block: number;
}

function createBaseBlockLimitsConfig(): BlockLimitsConfig {
  return {
    maxStatefulOrderRemovalsPerBlock: 0
  };
}

export const BlockLimitsConfig = {
  encode(message: BlockLimitsConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.maxStatefulOrderRemovalsPerBlock !== 0) {
      writer.uint32(8).uint32(message.maxStatefulOrderRemovalsPerBlock);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BlockLimitsConfig {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockLimitsConfig();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.maxStatefulOrderRemovalsPerBlock = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<BlockLimitsConfig>): BlockLimitsConfig {
    const message = createBaseBlockLimitsConfig();
    message.maxStatefulOrderRemovalsPerBlock = object.maxStatefulOrderRemovalsPerBlock ?? 0;
    return message;
  }

};