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
/**
 * LockedShares stores for an owner their total number of locked shares
 * and a schedule of share unlockings.
 */

export interface LockedShares {
  /** Address of the owner of below shares. */
  ownerAddress: string;
  /** Total number of locked shares. */

  totalLockedShares?: NumShares;
  /** Details of each unlock. */

  unlockDetails: UnlockDetail[];
}
/**
 * LockedShares stores for an owner their total number of locked shares
 * and a schedule of share unlockings.
 */

export interface LockedSharesSDKType {
  /** Address of the owner of below shares. */
  owner_address: string;
  /** Total number of locked shares. */

  total_locked_shares?: NumSharesSDKType;
  /** Details of each unlock. */

  unlock_details: UnlockDetailSDKType[];
}
/** UnlockDetail stores how many shares unlock at which block height. */

export interface UnlockDetail {
  /** Number of shares to unlock. */
  shares?: NumShares;
  /** Block height at which above shares unlock. */

  unlockBlockHeight: number;
}
/** UnlockDetail stores how many shares unlock at which block height. */

export interface UnlockDetailSDKType {
  /** Number of shares to unlock. */
  shares?: NumSharesSDKType;
  /** Block height at which above shares unlock. */

  unlock_block_height: number;
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

function createBaseLockedShares(): LockedShares {
  return {
    ownerAddress: "",
    totalLockedShares: undefined,
    unlockDetails: []
  };
}

export const LockedShares = {
  encode(message: LockedShares, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ownerAddress !== "") {
      writer.uint32(10).string(message.ownerAddress);
    }

    if (message.totalLockedShares !== undefined) {
      NumShares.encode(message.totalLockedShares, writer.uint32(18).fork()).ldelim();
    }

    for (const v of message.unlockDetails) {
      UnlockDetail.encode(v!, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LockedShares {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLockedShares();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.ownerAddress = reader.string();
          break;

        case 2:
          message.totalLockedShares = NumShares.decode(reader, reader.uint32());
          break;

        case 3:
          message.unlockDetails.push(UnlockDetail.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<LockedShares>): LockedShares {
    const message = createBaseLockedShares();
    message.ownerAddress = object.ownerAddress ?? "";
    message.totalLockedShares = object.totalLockedShares !== undefined && object.totalLockedShares !== null ? NumShares.fromPartial(object.totalLockedShares) : undefined;
    message.unlockDetails = object.unlockDetails?.map(e => UnlockDetail.fromPartial(e)) || [];
    return message;
  }

};

function createBaseUnlockDetail(): UnlockDetail {
  return {
    shares: undefined,
    unlockBlockHeight: 0
  };
}

export const UnlockDetail = {
  encode(message: UnlockDetail, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.shares !== undefined) {
      NumShares.encode(message.shares, writer.uint32(10).fork()).ldelim();
    }

    if (message.unlockBlockHeight !== 0) {
      writer.uint32(16).uint32(message.unlockBlockHeight);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UnlockDetail {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnlockDetail();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.shares = NumShares.decode(reader, reader.uint32());
          break;

        case 2:
          message.unlockBlockHeight = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<UnlockDetail>): UnlockDetail {
    const message = createBaseUnlockDetail();
    message.shares = object.shares !== undefined && object.shares !== null ? NumShares.fromPartial(object.shares) : undefined;
    message.unlockBlockHeight = object.unlockBlockHeight ?? 0;
    return message;
  }

};