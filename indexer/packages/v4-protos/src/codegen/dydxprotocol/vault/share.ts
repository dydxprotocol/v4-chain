import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** NumShares represents the number of shares. */

export interface NumShares {
  /** Number of shares. */
  numShares: Uint8Array;
}
/** NumShares represents the number of shares. */

export interface NumSharesSDKType {
  /** Number of shares. */
  num_shares: Uint8Array;
}
/** OwnerShare is a type for owner shares. */

export interface OwnerShare {
  owner: string;
  shares?: NumShares;
}
/** OwnerShare is a type for owner shares. */

export interface OwnerShareSDKType {
  owner: string;
  shares?: NumSharesSDKType;
}

function createBaseNumShares(): NumShares {
  return {
    numShares: new Uint8Array()
  };
}

export const NumShares = {
  encode(message: NumShares, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.numShares.length !== 0) {
      writer.uint32(18).bytes(message.numShares);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): NumShares {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNumShares();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 2:
          message.numShares = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<NumShares>): NumShares {
    const message = createBaseNumShares();
    message.numShares = object.numShares ?? new Uint8Array();
    return message;
  }

};

function createBaseOwnerShare(): OwnerShare {
  return {
    owner: "",
    shares: undefined
  };
}

export const OwnerShare = {
  encode(message: OwnerShare, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.owner !== "") {
      writer.uint32(10).string(message.owner);
    }

    if (message.shares !== undefined) {
      NumShares.encode(message.shares, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OwnerShare {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOwnerShare();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.owner = reader.string();
          break;

        case 2:
          message.shares = NumShares.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OwnerShare>): OwnerShare {
    const message = createBaseOwnerShare();
    message.owner = object.owner ?? "";
    message.shares = object.shares !== undefined && object.shares !== null ? NumShares.fromPartial(object.shares) : undefined;
    return message;
  }

};