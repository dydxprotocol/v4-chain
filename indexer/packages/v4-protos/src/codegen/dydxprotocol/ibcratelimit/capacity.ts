import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** DenomCapacity stores a list of rate limit capacity for a denom. */

export interface DenomCapacity {
  /**
   * denom is the denomination of the token being rate limited.
   * e.g. ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5
   */
  denom: string;
  /**
   * capacity_list is a list of capacity amount tracked for each `Limiter`
   * on the denom. This list has a 1:1 mapping to `limiter` list under
   * `LimitParams`.
   */

  capacityList: Uint8Array[];
}
/** DenomCapacity stores a list of rate limit capacity for a denom. */

export interface DenomCapacitySDKType {
  /**
   * denom is the denomination of the token being rate limited.
   * e.g. ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5
   */
  denom: string;
  /**
   * capacity_list is a list of capacity amount tracked for each `Limiter`
   * on the denom. This list has a 1:1 mapping to `limiter` list under
   * `LimitParams`.
   */

  capacity_list: Uint8Array[];
}

function createBaseDenomCapacity(): DenomCapacity {
  return {
    denom: "",
    capacityList: []
  };
}

export const DenomCapacity = {
  encode(message: DenomCapacity, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }

    for (const v of message.capacityList) {
      writer.uint32(18).bytes(v!);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DenomCapacity {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDenomCapacity();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.denom = reader.string();
          break;

        case 2:
          message.capacityList.push(reader.bytes());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<DenomCapacity>): DenomCapacity {
    const message = createBaseDenomCapacity();
    message.denom = object.denom ?? "";
    message.capacityList = object.capacityList?.map(e => e) || [];
    return message;
  }

};