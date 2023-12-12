import { OrderId, OrderIdAmino, OrderIdSDKType } from "./order";
import { SubaccountId, SubaccountIdAmino, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { BinaryReader, BinaryWriter } from "../../binary";
/**
 * ClobMatch represents an operations queue entry around all different types
 * of matches, specifically regular matches, liquidation matches, and
 * deleveraging matches.
 */
export interface ClobMatch {
  matchOrders?: MatchOrders;
  matchPerpetualLiquidation?: MatchPerpetualLiquidation;
  matchPerpetualDeleveraging?: MatchPerpetualDeleveraging;
}
export interface ClobMatchProtoMsg {
  typeUrl: "/dydxprotocol.clob.ClobMatch";
  value: Uint8Array;
}
/**
 * ClobMatch represents an operations queue entry around all different types
 * of matches, specifically regular matches, liquidation matches, and
 * deleveraging matches.
 */
export interface ClobMatchAmino {
  match_orders?: MatchOrdersAmino;
  match_perpetual_liquidation?: MatchPerpetualLiquidationAmino;
  match_perpetual_deleveraging?: MatchPerpetualDeleveragingAmino;
}
export interface ClobMatchAminoMsg {
  type: "/dydxprotocol.clob.ClobMatch";
  value: ClobMatchAmino;
}
/**
 * ClobMatch represents an operations queue entry around all different types
 * of matches, specifically regular matches, liquidation matches, and
 * deleveraging matches.
 */
export interface ClobMatchSDKType {
  match_orders?: MatchOrdersSDKType;
  match_perpetual_liquidation?: MatchPerpetualLiquidationSDKType;
  match_perpetual_deleveraging?: MatchPerpetualDeleveragingSDKType;
}
/** MakerFill represents the filled amount of a matched maker order. */
export interface MakerFill {
  /**
   * The filled amount of the matched maker order, in base quantums.
   * TODO(CLOB-571): update to use SerializableInt.
   */
  fillAmount: bigint;
  /** The `OrderId` of the matched maker order. */
  makerOrderId: OrderId;
}
export interface MakerFillProtoMsg {
  typeUrl: "/dydxprotocol.clob.MakerFill";
  value: Uint8Array;
}
/** MakerFill represents the filled amount of a matched maker order. */
export interface MakerFillAmino {
  /**
   * The filled amount of the matched maker order, in base quantums.
   * TODO(CLOB-571): update to use SerializableInt.
   */
  fill_amount?: string;
  /** The `OrderId` of the matched maker order. */
  maker_order_id?: OrderIdAmino;
}
export interface MakerFillAminoMsg {
  type: "/dydxprotocol.clob.MakerFill";
  value: MakerFillAmino;
}
/** MakerFill represents the filled amount of a matched maker order. */
export interface MakerFillSDKType {
  fill_amount: bigint;
  maker_order_id: OrderIdSDKType;
}
/** MatchOrders is an injected message used for matching orders. */
export interface MatchOrders {
  /** The `OrderId` of the taker order. */
  takerOrderId: OrderId;
  /** An ordered list of fills created by this taker order. */
  fills: MakerFill[];
}
export interface MatchOrdersProtoMsg {
  typeUrl: "/dydxprotocol.clob.MatchOrders";
  value: Uint8Array;
}
/** MatchOrders is an injected message used for matching orders. */
export interface MatchOrdersAmino {
  /** The `OrderId` of the taker order. */
  taker_order_id?: OrderIdAmino;
  /** An ordered list of fills created by this taker order. */
  fills?: MakerFillAmino[];
}
export interface MatchOrdersAminoMsg {
  type: "/dydxprotocol.clob.MatchOrders";
  value: MatchOrdersAmino;
}
/** MatchOrders is an injected message used for matching orders. */
export interface MatchOrdersSDKType {
  taker_order_id: OrderIdSDKType;
  fills: MakerFillSDKType[];
}
/**
 * MatchPerpetualLiquidation is an injected message used for liquidating a
 * subaccount.
 */
export interface MatchPerpetualLiquidation {
  /** ID of the subaccount that was liquidated. */
  liquidated: SubaccountId;
  /** The ID of the clob pair involved in the liquidation. */
  clobPairId: number;
  /** The ID of the perpetual involved in the liquidation. */
  perpetualId: number;
  /** The total size of the liquidation order including any unfilled size. */
  totalSize: bigint;
  /** `true` if liquidating a short position, `false` otherwise. */
  isBuy: boolean;
  /** An ordered list of fills created by this liquidation. */
  fills: MakerFill[];
}
export interface MatchPerpetualLiquidationProtoMsg {
  typeUrl: "/dydxprotocol.clob.MatchPerpetualLiquidation";
  value: Uint8Array;
}
/**
 * MatchPerpetualLiquidation is an injected message used for liquidating a
 * subaccount.
 */
export interface MatchPerpetualLiquidationAmino {
  /** ID of the subaccount that was liquidated. */
  liquidated?: SubaccountIdAmino;
  /** The ID of the clob pair involved in the liquidation. */
  clob_pair_id?: number;
  /** The ID of the perpetual involved in the liquidation. */
  perpetual_id?: number;
  /** The total size of the liquidation order including any unfilled size. */
  total_size?: string;
  /** `true` if liquidating a short position, `false` otherwise. */
  is_buy?: boolean;
  /** An ordered list of fills created by this liquidation. */
  fills?: MakerFillAmino[];
}
export interface MatchPerpetualLiquidationAminoMsg {
  type: "/dydxprotocol.clob.MatchPerpetualLiquidation";
  value: MatchPerpetualLiquidationAmino;
}
/**
 * MatchPerpetualLiquidation is an injected message used for liquidating a
 * subaccount.
 */
export interface MatchPerpetualLiquidationSDKType {
  liquidated: SubaccountIdSDKType;
  clob_pair_id: number;
  perpetual_id: number;
  total_size: bigint;
  is_buy: boolean;
  fills: MakerFillSDKType[];
}
/**
 * MatchPerpetualDeleveraging is an injected message used for deleveraging a
 * subaccount.
 */
export interface MatchPerpetualDeleveraging {
  /** ID of the subaccount that was liquidated. */
  liquidated: SubaccountId;
  /** The ID of the perpetual that was liquidated. */
  perpetualId: number;
  /** An ordered list of fills created by this liquidation. */
  fills: MatchPerpetualDeleveraging_Fill[];
  /**
   * Flag denoting whether the deleveraging operation was for the purpose
   * of final settlement. Final settlement matches are at the oracle price,
   * whereas deleveraging happens at the bankruptcy price of the deleveraged
   * subaccount.
   */
  isFinalSettlement: boolean;
}
export interface MatchPerpetualDeleveragingProtoMsg {
  typeUrl: "/dydxprotocol.clob.MatchPerpetualDeleveraging";
  value: Uint8Array;
}
/**
 * MatchPerpetualDeleveraging is an injected message used for deleveraging a
 * subaccount.
 */
export interface MatchPerpetualDeleveragingAmino {
  /** ID of the subaccount that was liquidated. */
  liquidated?: SubaccountIdAmino;
  /** The ID of the perpetual that was liquidated. */
  perpetual_id?: number;
  /** An ordered list of fills created by this liquidation. */
  fills?: MatchPerpetualDeleveraging_FillAmino[];
  /**
   * Flag denoting whether the deleveraging operation was for the purpose
   * of final settlement. Final settlement matches are at the oracle price,
   * whereas deleveraging happens at the bankruptcy price of the deleveraged
   * subaccount.
   */
  is_final_settlement?: boolean;
}
export interface MatchPerpetualDeleveragingAminoMsg {
  type: "/dydxprotocol.clob.MatchPerpetualDeleveraging";
  value: MatchPerpetualDeleveragingAmino;
}
/**
 * MatchPerpetualDeleveraging is an injected message used for deleveraging a
 * subaccount.
 */
export interface MatchPerpetualDeleveragingSDKType {
  liquidated: SubaccountIdSDKType;
  perpetual_id: number;
  fills: MatchPerpetualDeleveraging_FillSDKType[];
  is_final_settlement: boolean;
}
/** Fill represents a fill between the liquidated and offsetting subaccount. */
export interface MatchPerpetualDeleveraging_Fill {
  /**
   * ID of the subaccount that was used to offset the liquidated subaccount's
   * position.
   */
  offsettingSubaccountId: SubaccountId;
  /**
   * The amount filled between the liquidated and offsetting position, in
   * base quantums.
   * TODO(CLOB-571): update to use SerializableInt.
   */
  fillAmount: bigint;
}
export interface MatchPerpetualDeleveraging_FillProtoMsg {
  typeUrl: "/dydxprotocol.clob.Fill";
  value: Uint8Array;
}
/** Fill represents a fill between the liquidated and offsetting subaccount. */
export interface MatchPerpetualDeleveraging_FillAmino {
  /**
   * ID of the subaccount that was used to offset the liquidated subaccount's
   * position.
   */
  offsetting_subaccount_id?: SubaccountIdAmino;
  /**
   * The amount filled between the liquidated and offsetting position, in
   * base quantums.
   * TODO(CLOB-571): update to use SerializableInt.
   */
  fill_amount?: string;
}
export interface MatchPerpetualDeleveraging_FillAminoMsg {
  type: "/dydxprotocol.clob.Fill";
  value: MatchPerpetualDeleveraging_FillAmino;
}
/** Fill represents a fill between the liquidated and offsetting subaccount. */
export interface MatchPerpetualDeleveraging_FillSDKType {
  offsetting_subaccount_id: SubaccountIdSDKType;
  fill_amount: bigint;
}
function createBaseClobMatch(): ClobMatch {
  return {
    matchOrders: undefined,
    matchPerpetualLiquidation: undefined,
    matchPerpetualDeleveraging: undefined
  };
}
export const ClobMatch = {
  typeUrl: "/dydxprotocol.clob.ClobMatch",
  encode(message: ClobMatch, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.matchOrders !== undefined) {
      MatchOrders.encode(message.matchOrders, writer.uint32(10).fork()).ldelim();
    }
    if (message.matchPerpetualLiquidation !== undefined) {
      MatchPerpetualLiquidation.encode(message.matchPerpetualLiquidation, writer.uint32(18).fork()).ldelim();
    }
    if (message.matchPerpetualDeleveraging !== undefined) {
      MatchPerpetualDeleveraging.encode(message.matchPerpetualDeleveraging, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ClobMatch {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClobMatch();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.matchOrders = MatchOrders.decode(reader, reader.uint32());
          break;
        case 2:
          message.matchPerpetualLiquidation = MatchPerpetualLiquidation.decode(reader, reader.uint32());
          break;
        case 3:
          message.matchPerpetualDeleveraging = MatchPerpetualDeleveraging.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ClobMatch>): ClobMatch {
    const message = createBaseClobMatch();
    message.matchOrders = object.matchOrders !== undefined && object.matchOrders !== null ? MatchOrders.fromPartial(object.matchOrders) : undefined;
    message.matchPerpetualLiquidation = object.matchPerpetualLiquidation !== undefined && object.matchPerpetualLiquidation !== null ? MatchPerpetualLiquidation.fromPartial(object.matchPerpetualLiquidation) : undefined;
    message.matchPerpetualDeleveraging = object.matchPerpetualDeleveraging !== undefined && object.matchPerpetualDeleveraging !== null ? MatchPerpetualDeleveraging.fromPartial(object.matchPerpetualDeleveraging) : undefined;
    return message;
  },
  fromAmino(object: ClobMatchAmino): ClobMatch {
    const message = createBaseClobMatch();
    if (object.match_orders !== undefined && object.match_orders !== null) {
      message.matchOrders = MatchOrders.fromAmino(object.match_orders);
    }
    if (object.match_perpetual_liquidation !== undefined && object.match_perpetual_liquidation !== null) {
      message.matchPerpetualLiquidation = MatchPerpetualLiquidation.fromAmino(object.match_perpetual_liquidation);
    }
    if (object.match_perpetual_deleveraging !== undefined && object.match_perpetual_deleveraging !== null) {
      message.matchPerpetualDeleveraging = MatchPerpetualDeleveraging.fromAmino(object.match_perpetual_deleveraging);
    }
    return message;
  },
  toAmino(message: ClobMatch): ClobMatchAmino {
    const obj: any = {};
    obj.match_orders = message.matchOrders ? MatchOrders.toAmino(message.matchOrders) : undefined;
    obj.match_perpetual_liquidation = message.matchPerpetualLiquidation ? MatchPerpetualLiquidation.toAmino(message.matchPerpetualLiquidation) : undefined;
    obj.match_perpetual_deleveraging = message.matchPerpetualDeleveraging ? MatchPerpetualDeleveraging.toAmino(message.matchPerpetualDeleveraging) : undefined;
    return obj;
  },
  fromAminoMsg(object: ClobMatchAminoMsg): ClobMatch {
    return ClobMatch.fromAmino(object.value);
  },
  fromProtoMsg(message: ClobMatchProtoMsg): ClobMatch {
    return ClobMatch.decode(message.value);
  },
  toProto(message: ClobMatch): Uint8Array {
    return ClobMatch.encode(message).finish();
  },
  toProtoMsg(message: ClobMatch): ClobMatchProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.ClobMatch",
      value: ClobMatch.encode(message).finish()
    };
  }
};
function createBaseMakerFill(): MakerFill {
  return {
    fillAmount: BigInt(0),
    makerOrderId: OrderId.fromPartial({})
  };
}
export const MakerFill = {
  typeUrl: "/dydxprotocol.clob.MakerFill",
  encode(message: MakerFill, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.fillAmount !== BigInt(0)) {
      writer.uint32(8).uint64(message.fillAmount);
    }
    if (message.makerOrderId !== undefined) {
      OrderId.encode(message.makerOrderId, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MakerFill {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMakerFill();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fillAmount = reader.uint64();
          break;
        case 2:
          message.makerOrderId = OrderId.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MakerFill>): MakerFill {
    const message = createBaseMakerFill();
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? BigInt(object.fillAmount.toString()) : BigInt(0);
    message.makerOrderId = object.makerOrderId !== undefined && object.makerOrderId !== null ? OrderId.fromPartial(object.makerOrderId) : undefined;
    return message;
  },
  fromAmino(object: MakerFillAmino): MakerFill {
    const message = createBaseMakerFill();
    if (object.fill_amount !== undefined && object.fill_amount !== null) {
      message.fillAmount = BigInt(object.fill_amount);
    }
    if (object.maker_order_id !== undefined && object.maker_order_id !== null) {
      message.makerOrderId = OrderId.fromAmino(object.maker_order_id);
    }
    return message;
  },
  toAmino(message: MakerFill): MakerFillAmino {
    const obj: any = {};
    obj.fill_amount = message.fillAmount ? message.fillAmount.toString() : undefined;
    obj.maker_order_id = message.makerOrderId ? OrderId.toAmino(message.makerOrderId) : undefined;
    return obj;
  },
  fromAminoMsg(object: MakerFillAminoMsg): MakerFill {
    return MakerFill.fromAmino(object.value);
  },
  fromProtoMsg(message: MakerFillProtoMsg): MakerFill {
    return MakerFill.decode(message.value);
  },
  toProto(message: MakerFill): Uint8Array {
    return MakerFill.encode(message).finish();
  },
  toProtoMsg(message: MakerFill): MakerFillProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MakerFill",
      value: MakerFill.encode(message).finish()
    };
  }
};
function createBaseMatchOrders(): MatchOrders {
  return {
    takerOrderId: OrderId.fromPartial({}),
    fills: []
  };
}
export const MatchOrders = {
  typeUrl: "/dydxprotocol.clob.MatchOrders",
  encode(message: MatchOrders, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.takerOrderId !== undefined) {
      OrderId.encode(message.takerOrderId, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.fills) {
      MakerFill.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MatchOrders {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMatchOrders();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.takerOrderId = OrderId.decode(reader, reader.uint32());
          break;
        case 2:
          message.fills.push(MakerFill.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MatchOrders>): MatchOrders {
    const message = createBaseMatchOrders();
    message.takerOrderId = object.takerOrderId !== undefined && object.takerOrderId !== null ? OrderId.fromPartial(object.takerOrderId) : undefined;
    message.fills = object.fills?.map(e => MakerFill.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: MatchOrdersAmino): MatchOrders {
    const message = createBaseMatchOrders();
    if (object.taker_order_id !== undefined && object.taker_order_id !== null) {
      message.takerOrderId = OrderId.fromAmino(object.taker_order_id);
    }
    message.fills = object.fills?.map(e => MakerFill.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: MatchOrders): MatchOrdersAmino {
    const obj: any = {};
    obj.taker_order_id = message.takerOrderId ? OrderId.toAmino(message.takerOrderId) : undefined;
    if (message.fills) {
      obj.fills = message.fills.map(e => e ? MakerFill.toAmino(e) : undefined);
    } else {
      obj.fills = [];
    }
    return obj;
  },
  fromAminoMsg(object: MatchOrdersAminoMsg): MatchOrders {
    return MatchOrders.fromAmino(object.value);
  },
  fromProtoMsg(message: MatchOrdersProtoMsg): MatchOrders {
    return MatchOrders.decode(message.value);
  },
  toProto(message: MatchOrders): Uint8Array {
    return MatchOrders.encode(message).finish();
  },
  toProtoMsg(message: MatchOrders): MatchOrdersProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MatchOrders",
      value: MatchOrders.encode(message).finish()
    };
  }
};
function createBaseMatchPerpetualLiquidation(): MatchPerpetualLiquidation {
  return {
    liquidated: SubaccountId.fromPartial({}),
    clobPairId: 0,
    perpetualId: 0,
    totalSize: BigInt(0),
    isBuy: false,
    fills: []
  };
}
export const MatchPerpetualLiquidation = {
  typeUrl: "/dydxprotocol.clob.MatchPerpetualLiquidation",
  encode(message: MatchPerpetualLiquidation, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.liquidated !== undefined) {
      SubaccountId.encode(message.liquidated, writer.uint32(10).fork()).ldelim();
    }
    if (message.clobPairId !== 0) {
      writer.uint32(16).uint32(message.clobPairId);
    }
    if (message.perpetualId !== 0) {
      writer.uint32(24).uint32(message.perpetualId);
    }
    if (message.totalSize !== BigInt(0)) {
      writer.uint32(32).uint64(message.totalSize);
    }
    if (message.isBuy === true) {
      writer.uint32(40).bool(message.isBuy);
    }
    for (const v of message.fills) {
      MakerFill.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MatchPerpetualLiquidation {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMatchPerpetualLiquidation();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.liquidated = SubaccountId.decode(reader, reader.uint32());
          break;
        case 2:
          message.clobPairId = reader.uint32();
          break;
        case 3:
          message.perpetualId = reader.uint32();
          break;
        case 4:
          message.totalSize = reader.uint64();
          break;
        case 5:
          message.isBuy = reader.bool();
          break;
        case 6:
          message.fills.push(MakerFill.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MatchPerpetualLiquidation>): MatchPerpetualLiquidation {
    const message = createBaseMatchPerpetualLiquidation();
    message.liquidated = object.liquidated !== undefined && object.liquidated !== null ? SubaccountId.fromPartial(object.liquidated) : undefined;
    message.clobPairId = object.clobPairId ?? 0;
    message.perpetualId = object.perpetualId ?? 0;
    message.totalSize = object.totalSize !== undefined && object.totalSize !== null ? BigInt(object.totalSize.toString()) : BigInt(0);
    message.isBuy = object.isBuy ?? false;
    message.fills = object.fills?.map(e => MakerFill.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: MatchPerpetualLiquidationAmino): MatchPerpetualLiquidation {
    const message = createBaseMatchPerpetualLiquidation();
    if (object.liquidated !== undefined && object.liquidated !== null) {
      message.liquidated = SubaccountId.fromAmino(object.liquidated);
    }
    if (object.clob_pair_id !== undefined && object.clob_pair_id !== null) {
      message.clobPairId = object.clob_pair_id;
    }
    if (object.perpetual_id !== undefined && object.perpetual_id !== null) {
      message.perpetualId = object.perpetual_id;
    }
    if (object.total_size !== undefined && object.total_size !== null) {
      message.totalSize = BigInt(object.total_size);
    }
    if (object.is_buy !== undefined && object.is_buy !== null) {
      message.isBuy = object.is_buy;
    }
    message.fills = object.fills?.map(e => MakerFill.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: MatchPerpetualLiquidation): MatchPerpetualLiquidationAmino {
    const obj: any = {};
    obj.liquidated = message.liquidated ? SubaccountId.toAmino(message.liquidated) : undefined;
    obj.clob_pair_id = message.clobPairId;
    obj.perpetual_id = message.perpetualId;
    obj.total_size = message.totalSize ? message.totalSize.toString() : undefined;
    obj.is_buy = message.isBuy;
    if (message.fills) {
      obj.fills = message.fills.map(e => e ? MakerFill.toAmino(e) : undefined);
    } else {
      obj.fills = [];
    }
    return obj;
  },
  fromAminoMsg(object: MatchPerpetualLiquidationAminoMsg): MatchPerpetualLiquidation {
    return MatchPerpetualLiquidation.fromAmino(object.value);
  },
  fromProtoMsg(message: MatchPerpetualLiquidationProtoMsg): MatchPerpetualLiquidation {
    return MatchPerpetualLiquidation.decode(message.value);
  },
  toProto(message: MatchPerpetualLiquidation): Uint8Array {
    return MatchPerpetualLiquidation.encode(message).finish();
  },
  toProtoMsg(message: MatchPerpetualLiquidation): MatchPerpetualLiquidationProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MatchPerpetualLiquidation",
      value: MatchPerpetualLiquidation.encode(message).finish()
    };
  }
};
function createBaseMatchPerpetualDeleveraging(): MatchPerpetualDeleveraging {
  return {
    liquidated: SubaccountId.fromPartial({}),
    perpetualId: 0,
    fills: [],
    isFinalSettlement: false
  };
}
export const MatchPerpetualDeleveraging = {
  typeUrl: "/dydxprotocol.clob.MatchPerpetualDeleveraging",
  encode(message: MatchPerpetualDeleveraging, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.liquidated !== undefined) {
      SubaccountId.encode(message.liquidated, writer.uint32(10).fork()).ldelim();
    }
    if (message.perpetualId !== 0) {
      writer.uint32(16).uint32(message.perpetualId);
    }
    for (const v of message.fills) {
      MatchPerpetualDeleveraging_Fill.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.isFinalSettlement === true) {
      writer.uint32(32).bool(message.isFinalSettlement);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MatchPerpetualDeleveraging {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMatchPerpetualDeleveraging();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.liquidated = SubaccountId.decode(reader, reader.uint32());
          break;
        case 2:
          message.perpetualId = reader.uint32();
          break;
        case 3:
          message.fills.push(MatchPerpetualDeleveraging_Fill.decode(reader, reader.uint32()));
          break;
        case 4:
          message.isFinalSettlement = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MatchPerpetualDeleveraging>): MatchPerpetualDeleveraging {
    const message = createBaseMatchPerpetualDeleveraging();
    message.liquidated = object.liquidated !== undefined && object.liquidated !== null ? SubaccountId.fromPartial(object.liquidated) : undefined;
    message.perpetualId = object.perpetualId ?? 0;
    message.fills = object.fills?.map(e => MatchPerpetualDeleveraging_Fill.fromPartial(e)) || [];
    message.isFinalSettlement = object.isFinalSettlement ?? false;
    return message;
  },
  fromAmino(object: MatchPerpetualDeleveragingAmino): MatchPerpetualDeleveraging {
    const message = createBaseMatchPerpetualDeleveraging();
    if (object.liquidated !== undefined && object.liquidated !== null) {
      message.liquidated = SubaccountId.fromAmino(object.liquidated);
    }
    if (object.perpetual_id !== undefined && object.perpetual_id !== null) {
      message.perpetualId = object.perpetual_id;
    }
    message.fills = object.fills?.map(e => MatchPerpetualDeleveraging_Fill.fromAmino(e)) || [];
    if (object.is_final_settlement !== undefined && object.is_final_settlement !== null) {
      message.isFinalSettlement = object.is_final_settlement;
    }
    return message;
  },
  toAmino(message: MatchPerpetualDeleveraging): MatchPerpetualDeleveragingAmino {
    const obj: any = {};
    obj.liquidated = message.liquidated ? SubaccountId.toAmino(message.liquidated) : undefined;
    obj.perpetual_id = message.perpetualId;
    if (message.fills) {
      obj.fills = message.fills.map(e => e ? MatchPerpetualDeleveraging_Fill.toAmino(e) : undefined);
    } else {
      obj.fills = [];
    }
    obj.is_final_settlement = message.isFinalSettlement;
    return obj;
  },
  fromAminoMsg(object: MatchPerpetualDeleveragingAminoMsg): MatchPerpetualDeleveraging {
    return MatchPerpetualDeleveraging.fromAmino(object.value);
  },
  fromProtoMsg(message: MatchPerpetualDeleveragingProtoMsg): MatchPerpetualDeleveraging {
    return MatchPerpetualDeleveraging.decode(message.value);
  },
  toProto(message: MatchPerpetualDeleveraging): Uint8Array {
    return MatchPerpetualDeleveraging.encode(message).finish();
  },
  toProtoMsg(message: MatchPerpetualDeleveraging): MatchPerpetualDeleveragingProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MatchPerpetualDeleveraging",
      value: MatchPerpetualDeleveraging.encode(message).finish()
    };
  }
};
function createBaseMatchPerpetualDeleveraging_Fill(): MatchPerpetualDeleveraging_Fill {
  return {
    offsettingSubaccountId: SubaccountId.fromPartial({}),
    fillAmount: BigInt(0)
  };
}
export const MatchPerpetualDeleveraging_Fill = {
  typeUrl: "/dydxprotocol.clob.Fill",
  encode(message: MatchPerpetualDeleveraging_Fill, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.offsettingSubaccountId !== undefined) {
      SubaccountId.encode(message.offsettingSubaccountId, writer.uint32(10).fork()).ldelim();
    }
    if (message.fillAmount !== BigInt(0)) {
      writer.uint32(16).uint64(message.fillAmount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MatchPerpetualDeleveraging_Fill {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMatchPerpetualDeleveraging_Fill();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.offsettingSubaccountId = SubaccountId.decode(reader, reader.uint32());
          break;
        case 2:
          message.fillAmount = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MatchPerpetualDeleveraging_Fill>): MatchPerpetualDeleveraging_Fill {
    const message = createBaseMatchPerpetualDeleveraging_Fill();
    message.offsettingSubaccountId = object.offsettingSubaccountId !== undefined && object.offsettingSubaccountId !== null ? SubaccountId.fromPartial(object.offsettingSubaccountId) : undefined;
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? BigInt(object.fillAmount.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MatchPerpetualDeleveraging_FillAmino): MatchPerpetualDeleveraging_Fill {
    const message = createBaseMatchPerpetualDeleveraging_Fill();
    if (object.offsetting_subaccount_id !== undefined && object.offsetting_subaccount_id !== null) {
      message.offsettingSubaccountId = SubaccountId.fromAmino(object.offsetting_subaccount_id);
    }
    if (object.fill_amount !== undefined && object.fill_amount !== null) {
      message.fillAmount = BigInt(object.fill_amount);
    }
    return message;
  },
  toAmino(message: MatchPerpetualDeleveraging_Fill): MatchPerpetualDeleveraging_FillAmino {
    const obj: any = {};
    obj.offsetting_subaccount_id = message.offsettingSubaccountId ? SubaccountId.toAmino(message.offsettingSubaccountId) : undefined;
    obj.fill_amount = message.fillAmount ? message.fillAmount.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MatchPerpetualDeleveraging_FillAminoMsg): MatchPerpetualDeleveraging_Fill {
    return MatchPerpetualDeleveraging_Fill.fromAmino(object.value);
  },
  fromProtoMsg(message: MatchPerpetualDeleveraging_FillProtoMsg): MatchPerpetualDeleveraging_Fill {
    return MatchPerpetualDeleveraging_Fill.decode(message.value);
  },
  toProto(message: MatchPerpetualDeleveraging_Fill): Uint8Array {
    return MatchPerpetualDeleveraging_Fill.encode(message).finish();
  },
  toProtoMsg(message: MatchPerpetualDeleveraging_Fill): MatchPerpetualDeleveraging_FillProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.Fill",
      value: MatchPerpetualDeleveraging_Fill.encode(message).finish()
    };
  }
};