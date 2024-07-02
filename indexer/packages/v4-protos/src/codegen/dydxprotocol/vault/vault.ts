import { MarketPrice, MarketPriceSDKType } from "../prices/market_price";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** VaultType represents different types of vaults. */

export enum VaultType {
  /** VAULT_TYPE_UNSPECIFIED - Default value, invalid and unused. */
  VAULT_TYPE_UNSPECIFIED = 0,

  /** VAULT_TYPE_CLOB - Vault is associated with a CLOB pair. */
  VAULT_TYPE_CLOB = 1,
  UNRECOGNIZED = -1,
}
/** VaultType represents different types of vaults. */

export enum VaultTypeSDKType {
  /** VAULT_TYPE_UNSPECIFIED - Default value, invalid and unused. */
  VAULT_TYPE_UNSPECIFIED = 0,

  /** VAULT_TYPE_CLOB - Vault is associated with a CLOB pair. */
  VAULT_TYPE_CLOB = 1,
  UNRECOGNIZED = -1,
}
export function vaultTypeFromJSON(object: any): VaultType {
  switch (object) {
    case 0:
    case "VAULT_TYPE_UNSPECIFIED":
      return VaultType.VAULT_TYPE_UNSPECIFIED;

    case 1:
    case "VAULT_TYPE_CLOB":
      return VaultType.VAULT_TYPE_CLOB;

    case -1:
    case "UNRECOGNIZED":
    default:
      return VaultType.UNRECOGNIZED;
  }
}
export function vaultTypeToJSON(object: VaultType): string {
  switch (object) {
    case VaultType.VAULT_TYPE_UNSPECIFIED:
      return "VAULT_TYPE_UNSPECIFIED";

    case VaultType.VAULT_TYPE_CLOB:
      return "VAULT_TYPE_CLOB";

    case VaultType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/** VaultId uniquely identifies a vault by its type and number. */

export interface VaultId {
  /** Type of the vault. */
  type: VaultType;
  /** Unique ID of the vault within above type. */

  number: number;
}
/** VaultId uniquely identifies a vault by its type and number. */

export interface VaultIdSDKType {
  /** Type of the vault. */
  type: VaultTypeSDKType;
  /** Unique ID of the vault within above type. */

  number: number;
}
/** NumShares represents the number of shares in a vault. */

export interface NumShares {
  /** Number of shares. */
  numShares: Uint8Array;
}
/** NumShares represents the number of shares in a vault. */

export interface NumSharesSDKType {
  /** Number of shares. */
  num_shares: Uint8Array;
}
/** OwnerShare is a type for owner shares in a vault. */

export interface OwnerShare {
  owner: string;
  shares?: NumShares;
}
/** OwnerShare is a type for owner shares in a vault. */

export interface OwnerShareSDKType {
  owner: string;
  shares?: NumSharesSDKType;
}
/** VaultParams is the individual parameters of a vault. */

export interface VaultParams {
  /** Lagged price that the vault quotes at. */
  laggedPrice?: MarketPrice;
}
/** VaultParams is the individual parameters of a vault. */

export interface VaultParamsSDKType {
  /** Lagged price that the vault quotes at. */
  lagged_price?: MarketPriceSDKType;
}

function createBaseVaultId(): VaultId {
  return {
    type: 0,
    number: 0
  };
}

export const VaultId = {
  encode(message: VaultId, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }

    if (message.number !== 0) {
      writer.uint32(16).uint32(message.number);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VaultId {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVaultId();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.type = (reader.int32() as any);
          break;

        case 2:
          message.number = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<VaultId>): VaultId {
    const message = createBaseVaultId();
    message.type = object.type ?? 0;
    message.number = object.number ?? 0;
    return message;
  }

};

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

function createBaseVaultParams(): VaultParams {
  return {
    laggedPrice: undefined
  };
}

export const VaultParams = {
  encode(message: VaultParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.laggedPrice !== undefined) {
      MarketPrice.encode(message.laggedPrice, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VaultParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVaultParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.laggedPrice = MarketPrice.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<VaultParams>): VaultParams {
    const message = createBaseVaultParams();
    message.laggedPrice = object.laggedPrice !== undefined && object.laggedPrice !== null ? MarketPrice.fromPartial(object.laggedPrice) : undefined;
    return message;
  }

};