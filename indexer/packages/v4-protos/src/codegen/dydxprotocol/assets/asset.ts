import { BinaryReader, BinaryWriter } from "../../binary";
/** Asset defines a single exchangable asset. */
export interface Asset {
  /** Unique, sequentially-generated. */
  id: number;
  /**
   * The human readable symbol of the `Asset` (e.g. `USDC`, `ATOM`).
   * Must be uppercase, unique and correspond to the canonical symbol of the
   * full coin.
   */
  symbol: string;
  /**
   * The name of base denomination unit of the `Asset` (e.g. `uatom`,
   * 'ibc/xxxxx'). Must be unique and match the `denom` used in the `sdk.Coin`
   * type in the `x/bank` module.
   */
  denom: string;
  /**
   * The exponent of converting one unit of `denom` to a full coin.
   * For example, `name=USDC, denom=uusdc, denom_exponent=-6` defines that
   * `1 uusdc = 10^(-6) USDC`. Note that `uusdc` refers to a `Coin` type in
   * `x/bank`, where the prefix `u` means `micro` by convetion. `uusdc` is
   * a different concept from a "quantum" defined by `atomic_resolution` below.
   * To convert from an amount of `denom` to quantums:
   * `quantums = denom_amount * 10^(denom_exponent - atomic_resolution)`
   */
  denomExponent: number;
  /** `true` if this `Asset` has a valid `MarketId` value. */
  hasMarket: boolean;
  /**
   * The `Id` of the `Market` associated with this `Asset`. It acts as the
   * oracle price for the purposes of calculating collateral
   * and margin requirements.
   */
  marketId: number;
  /**
   * The exponent for converting an atomic amount (1 'quantum')
   * to a full coin. For example, if `atomic_resolution = -8`
   * then an `asset_position` with `base_quantums = 1e8` is equivalent to
   * a position size of one full coin.
   */
  atomicResolution: number;
}
export interface AssetProtoMsg {
  typeUrl: "/dydxprotocol.assets.Asset";
  value: Uint8Array;
}
/** Asset defines a single exchangable asset. */
export interface AssetAmino {
  /** Unique, sequentially-generated. */
  id?: number;
  /**
   * The human readable symbol of the `Asset` (e.g. `USDC`, `ATOM`).
   * Must be uppercase, unique and correspond to the canonical symbol of the
   * full coin.
   */
  symbol?: string;
  /**
   * The name of base denomination unit of the `Asset` (e.g. `uatom`,
   * 'ibc/xxxxx'). Must be unique and match the `denom` used in the `sdk.Coin`
   * type in the `x/bank` module.
   */
  denom?: string;
  /**
   * The exponent of converting one unit of `denom` to a full coin.
   * For example, `name=USDC, denom=uusdc, denom_exponent=-6` defines that
   * `1 uusdc = 10^(-6) USDC`. Note that `uusdc` refers to a `Coin` type in
   * `x/bank`, where the prefix `u` means `micro` by convetion. `uusdc` is
   * a different concept from a "quantum" defined by `atomic_resolution` below.
   * To convert from an amount of `denom` to quantums:
   * `quantums = denom_amount * 10^(denom_exponent - atomic_resolution)`
   */
  denom_exponent?: number;
  /** `true` if this `Asset` has a valid `MarketId` value. */
  has_market?: boolean;
  /**
   * The `Id` of the `Market` associated with this `Asset`. It acts as the
   * oracle price for the purposes of calculating collateral
   * and margin requirements.
   */
  market_id?: number;
  /**
   * The exponent for converting an atomic amount (1 'quantum')
   * to a full coin. For example, if `atomic_resolution = -8`
   * then an `asset_position` with `base_quantums = 1e8` is equivalent to
   * a position size of one full coin.
   */
  atomic_resolution?: number;
}
export interface AssetAminoMsg {
  type: "/dydxprotocol.assets.Asset";
  value: AssetAmino;
}
/** Asset defines a single exchangable asset. */
export interface AssetSDKType {
  id: number;
  symbol: string;
  denom: string;
  denom_exponent: number;
  has_market: boolean;
  market_id: number;
  atomic_resolution: number;
}
function createBaseAsset(): Asset {
  return {
    id: 0,
    symbol: "",
    denom: "",
    denomExponent: 0,
    hasMarket: false,
    marketId: 0,
    atomicResolution: 0
  };
}
export const Asset = {
  typeUrl: "/dydxprotocol.assets.Asset",
  encode(message: Asset, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }
    if (message.symbol !== "") {
      writer.uint32(18).string(message.symbol);
    }
    if (message.denom !== "") {
      writer.uint32(26).string(message.denom);
    }
    if (message.denomExponent !== 0) {
      writer.uint32(32).sint32(message.denomExponent);
    }
    if (message.hasMarket === true) {
      writer.uint32(40).bool(message.hasMarket);
    }
    if (message.marketId !== 0) {
      writer.uint32(48).uint32(message.marketId);
    }
    if (message.atomicResolution !== 0) {
      writer.uint32(56).sint32(message.atomicResolution);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Asset {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAsset();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;
        case 2:
          message.symbol = reader.string();
          break;
        case 3:
          message.denom = reader.string();
          break;
        case 4:
          message.denomExponent = reader.sint32();
          break;
        case 5:
          message.hasMarket = reader.bool();
          break;
        case 6:
          message.marketId = reader.uint32();
          break;
        case 7:
          message.atomicResolution = reader.sint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Asset>): Asset {
    const message = createBaseAsset();
    message.id = object.id ?? 0;
    message.symbol = object.symbol ?? "";
    message.denom = object.denom ?? "";
    message.denomExponent = object.denomExponent ?? 0;
    message.hasMarket = object.hasMarket ?? false;
    message.marketId = object.marketId ?? 0;
    message.atomicResolution = object.atomicResolution ?? 0;
    return message;
  },
  fromAmino(object: AssetAmino): Asset {
    const message = createBaseAsset();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    if (object.symbol !== undefined && object.symbol !== null) {
      message.symbol = object.symbol;
    }
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = object.denom;
    }
    if (object.denom_exponent !== undefined && object.denom_exponent !== null) {
      message.denomExponent = object.denom_exponent;
    }
    if (object.has_market !== undefined && object.has_market !== null) {
      message.hasMarket = object.has_market;
    }
    if (object.market_id !== undefined && object.market_id !== null) {
      message.marketId = object.market_id;
    }
    if (object.atomic_resolution !== undefined && object.atomic_resolution !== null) {
      message.atomicResolution = object.atomic_resolution;
    }
    return message;
  },
  toAmino(message: Asset): AssetAmino {
    const obj: any = {};
    obj.id = message.id;
    obj.symbol = message.symbol;
    obj.denom = message.denom;
    obj.denom_exponent = message.denomExponent;
    obj.has_market = message.hasMarket;
    obj.market_id = message.marketId;
    obj.atomic_resolution = message.atomicResolution;
    return obj;
  },
  fromAminoMsg(object: AssetAminoMsg): Asset {
    return Asset.fromAmino(object.value);
  },
  fromProtoMsg(message: AssetProtoMsg): Asset {
    return Asset.decode(message.value);
  },
  toProto(message: Asset): Uint8Array {
    return Asset.encode(message).finish();
  },
  toProtoMsg(message: Asset): AssetProtoMsg {
    return {
      typeUrl: "/dydxprotocol.assets.Asset",
      value: Asset.encode(message).finish()
    };
  }
};