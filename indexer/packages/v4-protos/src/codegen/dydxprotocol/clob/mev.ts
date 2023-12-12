import { SubaccountId, SubaccountIdAmino, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { ClobPair, ClobPairAmino, ClobPairSDKType } from "./clob_pair";
import { BinaryReader, BinaryWriter } from "../../binary";
/** MEVMatch represents all necessary data to calculate MEV for a regular match. */
export interface MEVMatch {
  takerOrderSubaccountId?: SubaccountId;
  takerFeePpm: number;
  makerOrderSubaccountId?: SubaccountId;
  makerOrderSubticks: bigint;
  makerOrderIsBuy: boolean;
  makerFeePpm: number;
  clobPairId: number;
  fillAmount: bigint;
}
export interface MEVMatchProtoMsg {
  typeUrl: "/dydxprotocol.clob.MEVMatch";
  value: Uint8Array;
}
/** MEVMatch represents all necessary data to calculate MEV for a regular match. */
export interface MEVMatchAmino {
  taker_order_subaccount_id?: SubaccountIdAmino;
  taker_fee_ppm?: number;
  maker_order_subaccount_id?: SubaccountIdAmino;
  maker_order_subticks?: string;
  maker_order_is_buy?: boolean;
  maker_fee_ppm?: number;
  clob_pair_id?: number;
  fill_amount?: string;
}
export interface MEVMatchAminoMsg {
  type: "/dydxprotocol.clob.MEVMatch";
  value: MEVMatchAmino;
}
/** MEVMatch represents all necessary data to calculate MEV for a regular match. */
export interface MEVMatchSDKType {
  taker_order_subaccount_id?: SubaccountIdSDKType;
  taker_fee_ppm: number;
  maker_order_subaccount_id?: SubaccountIdSDKType;
  maker_order_subticks: bigint;
  maker_order_is_buy: boolean;
  maker_fee_ppm: number;
  clob_pair_id: number;
  fill_amount: bigint;
}
/**
 * MEVLiquidationMatch represents all necessary data to calculate MEV for a
 * liquidation.
 */
export interface MEVLiquidationMatch {
  liquidatedSubaccountId: SubaccountId;
  insuranceFundDeltaQuoteQuantums: bigint;
  makerOrderSubaccountId: SubaccountId;
  makerOrderSubticks: bigint;
  makerOrderIsBuy: boolean;
  makerFeePpm: number;
  clobPairId: number;
  fillAmount: bigint;
}
export interface MEVLiquidationMatchProtoMsg {
  typeUrl: "/dydxprotocol.clob.MEVLiquidationMatch";
  value: Uint8Array;
}
/**
 * MEVLiquidationMatch represents all necessary data to calculate MEV for a
 * liquidation.
 */
export interface MEVLiquidationMatchAmino {
  liquidated_subaccount_id?: SubaccountIdAmino;
  insurance_fund_delta_quote_quantums?: string;
  maker_order_subaccount_id?: SubaccountIdAmino;
  maker_order_subticks?: string;
  maker_order_is_buy?: boolean;
  maker_fee_ppm?: number;
  clob_pair_id?: number;
  fill_amount?: string;
}
export interface MEVLiquidationMatchAminoMsg {
  type: "/dydxprotocol.clob.MEVLiquidationMatch";
  value: MEVLiquidationMatchAmino;
}
/**
 * MEVLiquidationMatch represents all necessary data to calculate MEV for a
 * liquidation.
 */
export interface MEVLiquidationMatchSDKType {
  liquidated_subaccount_id: SubaccountIdSDKType;
  insurance_fund_delta_quote_quantums: bigint;
  maker_order_subaccount_id: SubaccountIdSDKType;
  maker_order_subticks: bigint;
  maker_order_is_buy: boolean;
  maker_fee_ppm: number;
  clob_pair_id: number;
  fill_amount: bigint;
}
/** ClobMidPrice contains the mid price of a CLOB pair, represented by it's ID. */
export interface ClobMidPrice {
  clobPair: ClobPair;
  subticks: bigint;
}
export interface ClobMidPriceProtoMsg {
  typeUrl: "/dydxprotocol.clob.ClobMidPrice";
  value: Uint8Array;
}
/** ClobMidPrice contains the mid price of a CLOB pair, represented by it's ID. */
export interface ClobMidPriceAmino {
  clob_pair?: ClobPairAmino;
  subticks?: string;
}
export interface ClobMidPriceAminoMsg {
  type: "/dydxprotocol.clob.ClobMidPrice";
  value: ClobMidPriceAmino;
}
/** ClobMidPrice contains the mid price of a CLOB pair, represented by it's ID. */
export interface ClobMidPriceSDKType {
  clob_pair: ClobPairSDKType;
  subticks: bigint;
}
/**
 * ValidatorMevMatches contains all matches from the validator's local
 * operations queue.
 */
export interface ValidatorMevMatches {
  matches: MEVMatch[];
  liquidationMatches: MEVLiquidationMatch[];
}
export interface ValidatorMevMatchesProtoMsg {
  typeUrl: "/dydxprotocol.clob.ValidatorMevMatches";
  value: Uint8Array;
}
/**
 * ValidatorMevMatches contains all matches from the validator's local
 * operations queue.
 */
export interface ValidatorMevMatchesAmino {
  matches?: MEVMatchAmino[];
  liquidation_matches?: MEVLiquidationMatchAmino[];
}
export interface ValidatorMevMatchesAminoMsg {
  type: "/dydxprotocol.clob.ValidatorMevMatches";
  value: ValidatorMevMatchesAmino;
}
/**
 * ValidatorMevMatches contains all matches from the validator's local
 * operations queue.
 */
export interface ValidatorMevMatchesSDKType {
  matches: MEVMatchSDKType[];
  liquidation_matches: MEVLiquidationMatchSDKType[];
}
/**
 * MevNodeToNodeMetrics is a data structure for encapsulating all MEV node <>
 * node metrics.
 */
export interface MevNodeToNodeMetrics {
  validatorMevMatches?: ValidatorMevMatches;
  clobMidPrices: ClobMidPrice[];
  bpMevMatches?: ValidatorMevMatches;
  proposalReceiveTime: bigint;
}
export interface MevNodeToNodeMetricsProtoMsg {
  typeUrl: "/dydxprotocol.clob.MevNodeToNodeMetrics";
  value: Uint8Array;
}
/**
 * MevNodeToNodeMetrics is a data structure for encapsulating all MEV node <>
 * node metrics.
 */
export interface MevNodeToNodeMetricsAmino {
  validator_mev_matches?: ValidatorMevMatchesAmino;
  clob_mid_prices?: ClobMidPriceAmino[];
  bp_mev_matches?: ValidatorMevMatchesAmino;
  proposal_receive_time?: string;
}
export interface MevNodeToNodeMetricsAminoMsg {
  type: "/dydxprotocol.clob.MevNodeToNodeMetrics";
  value: MevNodeToNodeMetricsAmino;
}
/**
 * MevNodeToNodeMetrics is a data structure for encapsulating all MEV node <>
 * node metrics.
 */
export interface MevNodeToNodeMetricsSDKType {
  validator_mev_matches?: ValidatorMevMatchesSDKType;
  clob_mid_prices: ClobMidPriceSDKType[];
  bp_mev_matches?: ValidatorMevMatchesSDKType;
  proposal_receive_time: bigint;
}
function createBaseMEVMatch(): MEVMatch {
  return {
    takerOrderSubaccountId: undefined,
    takerFeePpm: 0,
    makerOrderSubaccountId: undefined,
    makerOrderSubticks: BigInt(0),
    makerOrderIsBuy: false,
    makerFeePpm: 0,
    clobPairId: 0,
    fillAmount: BigInt(0)
  };
}
export const MEVMatch = {
  typeUrl: "/dydxprotocol.clob.MEVMatch",
  encode(message: MEVMatch, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.takerOrderSubaccountId !== undefined) {
      SubaccountId.encode(message.takerOrderSubaccountId, writer.uint32(10).fork()).ldelim();
    }
    if (message.takerFeePpm !== 0) {
      writer.uint32(16).int32(message.takerFeePpm);
    }
    if (message.makerOrderSubaccountId !== undefined) {
      SubaccountId.encode(message.makerOrderSubaccountId, writer.uint32(26).fork()).ldelim();
    }
    if (message.makerOrderSubticks !== BigInt(0)) {
      writer.uint32(32).uint64(message.makerOrderSubticks);
    }
    if (message.makerOrderIsBuy === true) {
      writer.uint32(40).bool(message.makerOrderIsBuy);
    }
    if (message.makerFeePpm !== 0) {
      writer.uint32(48).int32(message.makerFeePpm);
    }
    if (message.clobPairId !== 0) {
      writer.uint32(56).uint32(message.clobPairId);
    }
    if (message.fillAmount !== BigInt(0)) {
      writer.uint32(64).uint64(message.fillAmount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MEVMatch {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMEVMatch();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.takerOrderSubaccountId = SubaccountId.decode(reader, reader.uint32());
          break;
        case 2:
          message.takerFeePpm = reader.int32();
          break;
        case 3:
          message.makerOrderSubaccountId = SubaccountId.decode(reader, reader.uint32());
          break;
        case 4:
          message.makerOrderSubticks = reader.uint64();
          break;
        case 5:
          message.makerOrderIsBuy = reader.bool();
          break;
        case 6:
          message.makerFeePpm = reader.int32();
          break;
        case 7:
          message.clobPairId = reader.uint32();
          break;
        case 8:
          message.fillAmount = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MEVMatch>): MEVMatch {
    const message = createBaseMEVMatch();
    message.takerOrderSubaccountId = object.takerOrderSubaccountId !== undefined && object.takerOrderSubaccountId !== null ? SubaccountId.fromPartial(object.takerOrderSubaccountId) : undefined;
    message.takerFeePpm = object.takerFeePpm ?? 0;
    message.makerOrderSubaccountId = object.makerOrderSubaccountId !== undefined && object.makerOrderSubaccountId !== null ? SubaccountId.fromPartial(object.makerOrderSubaccountId) : undefined;
    message.makerOrderSubticks = object.makerOrderSubticks !== undefined && object.makerOrderSubticks !== null ? BigInt(object.makerOrderSubticks.toString()) : BigInt(0);
    message.makerOrderIsBuy = object.makerOrderIsBuy ?? false;
    message.makerFeePpm = object.makerFeePpm ?? 0;
    message.clobPairId = object.clobPairId ?? 0;
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? BigInt(object.fillAmount.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MEVMatchAmino): MEVMatch {
    const message = createBaseMEVMatch();
    if (object.taker_order_subaccount_id !== undefined && object.taker_order_subaccount_id !== null) {
      message.takerOrderSubaccountId = SubaccountId.fromAmino(object.taker_order_subaccount_id);
    }
    if (object.taker_fee_ppm !== undefined && object.taker_fee_ppm !== null) {
      message.takerFeePpm = object.taker_fee_ppm;
    }
    if (object.maker_order_subaccount_id !== undefined && object.maker_order_subaccount_id !== null) {
      message.makerOrderSubaccountId = SubaccountId.fromAmino(object.maker_order_subaccount_id);
    }
    if (object.maker_order_subticks !== undefined && object.maker_order_subticks !== null) {
      message.makerOrderSubticks = BigInt(object.maker_order_subticks);
    }
    if (object.maker_order_is_buy !== undefined && object.maker_order_is_buy !== null) {
      message.makerOrderIsBuy = object.maker_order_is_buy;
    }
    if (object.maker_fee_ppm !== undefined && object.maker_fee_ppm !== null) {
      message.makerFeePpm = object.maker_fee_ppm;
    }
    if (object.clob_pair_id !== undefined && object.clob_pair_id !== null) {
      message.clobPairId = object.clob_pair_id;
    }
    if (object.fill_amount !== undefined && object.fill_amount !== null) {
      message.fillAmount = BigInt(object.fill_amount);
    }
    return message;
  },
  toAmino(message: MEVMatch): MEVMatchAmino {
    const obj: any = {};
    obj.taker_order_subaccount_id = message.takerOrderSubaccountId ? SubaccountId.toAmino(message.takerOrderSubaccountId) : undefined;
    obj.taker_fee_ppm = message.takerFeePpm;
    obj.maker_order_subaccount_id = message.makerOrderSubaccountId ? SubaccountId.toAmino(message.makerOrderSubaccountId) : undefined;
    obj.maker_order_subticks = message.makerOrderSubticks ? message.makerOrderSubticks.toString() : undefined;
    obj.maker_order_is_buy = message.makerOrderIsBuy;
    obj.maker_fee_ppm = message.makerFeePpm;
    obj.clob_pair_id = message.clobPairId;
    obj.fill_amount = message.fillAmount ? message.fillAmount.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MEVMatchAminoMsg): MEVMatch {
    return MEVMatch.fromAmino(object.value);
  },
  fromProtoMsg(message: MEVMatchProtoMsg): MEVMatch {
    return MEVMatch.decode(message.value);
  },
  toProto(message: MEVMatch): Uint8Array {
    return MEVMatch.encode(message).finish();
  },
  toProtoMsg(message: MEVMatch): MEVMatchProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MEVMatch",
      value: MEVMatch.encode(message).finish()
    };
  }
};
function createBaseMEVLiquidationMatch(): MEVLiquidationMatch {
  return {
    liquidatedSubaccountId: SubaccountId.fromPartial({}),
    insuranceFundDeltaQuoteQuantums: BigInt(0),
    makerOrderSubaccountId: SubaccountId.fromPartial({}),
    makerOrderSubticks: BigInt(0),
    makerOrderIsBuy: false,
    makerFeePpm: 0,
    clobPairId: 0,
    fillAmount: BigInt(0)
  };
}
export const MEVLiquidationMatch = {
  typeUrl: "/dydxprotocol.clob.MEVLiquidationMatch",
  encode(message: MEVLiquidationMatch, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.liquidatedSubaccountId !== undefined) {
      SubaccountId.encode(message.liquidatedSubaccountId, writer.uint32(10).fork()).ldelim();
    }
    if (message.insuranceFundDeltaQuoteQuantums !== BigInt(0)) {
      writer.uint32(16).int64(message.insuranceFundDeltaQuoteQuantums);
    }
    if (message.makerOrderSubaccountId !== undefined) {
      SubaccountId.encode(message.makerOrderSubaccountId, writer.uint32(26).fork()).ldelim();
    }
    if (message.makerOrderSubticks !== BigInt(0)) {
      writer.uint32(32).uint64(message.makerOrderSubticks);
    }
    if (message.makerOrderIsBuy === true) {
      writer.uint32(40).bool(message.makerOrderIsBuy);
    }
    if (message.makerFeePpm !== 0) {
      writer.uint32(48).int32(message.makerFeePpm);
    }
    if (message.clobPairId !== 0) {
      writer.uint32(56).uint32(message.clobPairId);
    }
    if (message.fillAmount !== BigInt(0)) {
      writer.uint32(64).uint64(message.fillAmount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MEVLiquidationMatch {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMEVLiquidationMatch();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.liquidatedSubaccountId = SubaccountId.decode(reader, reader.uint32());
          break;
        case 2:
          message.insuranceFundDeltaQuoteQuantums = reader.int64();
          break;
        case 3:
          message.makerOrderSubaccountId = SubaccountId.decode(reader, reader.uint32());
          break;
        case 4:
          message.makerOrderSubticks = reader.uint64();
          break;
        case 5:
          message.makerOrderIsBuy = reader.bool();
          break;
        case 6:
          message.makerFeePpm = reader.int32();
          break;
        case 7:
          message.clobPairId = reader.uint32();
          break;
        case 8:
          message.fillAmount = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MEVLiquidationMatch>): MEVLiquidationMatch {
    const message = createBaseMEVLiquidationMatch();
    message.liquidatedSubaccountId = object.liquidatedSubaccountId !== undefined && object.liquidatedSubaccountId !== null ? SubaccountId.fromPartial(object.liquidatedSubaccountId) : undefined;
    message.insuranceFundDeltaQuoteQuantums = object.insuranceFundDeltaQuoteQuantums !== undefined && object.insuranceFundDeltaQuoteQuantums !== null ? BigInt(object.insuranceFundDeltaQuoteQuantums.toString()) : BigInt(0);
    message.makerOrderSubaccountId = object.makerOrderSubaccountId !== undefined && object.makerOrderSubaccountId !== null ? SubaccountId.fromPartial(object.makerOrderSubaccountId) : undefined;
    message.makerOrderSubticks = object.makerOrderSubticks !== undefined && object.makerOrderSubticks !== null ? BigInt(object.makerOrderSubticks.toString()) : BigInt(0);
    message.makerOrderIsBuy = object.makerOrderIsBuy ?? false;
    message.makerFeePpm = object.makerFeePpm ?? 0;
    message.clobPairId = object.clobPairId ?? 0;
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? BigInt(object.fillAmount.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MEVLiquidationMatchAmino): MEVLiquidationMatch {
    const message = createBaseMEVLiquidationMatch();
    if (object.liquidated_subaccount_id !== undefined && object.liquidated_subaccount_id !== null) {
      message.liquidatedSubaccountId = SubaccountId.fromAmino(object.liquidated_subaccount_id);
    }
    if (object.insurance_fund_delta_quote_quantums !== undefined && object.insurance_fund_delta_quote_quantums !== null) {
      message.insuranceFundDeltaQuoteQuantums = BigInt(object.insurance_fund_delta_quote_quantums);
    }
    if (object.maker_order_subaccount_id !== undefined && object.maker_order_subaccount_id !== null) {
      message.makerOrderSubaccountId = SubaccountId.fromAmino(object.maker_order_subaccount_id);
    }
    if (object.maker_order_subticks !== undefined && object.maker_order_subticks !== null) {
      message.makerOrderSubticks = BigInt(object.maker_order_subticks);
    }
    if (object.maker_order_is_buy !== undefined && object.maker_order_is_buy !== null) {
      message.makerOrderIsBuy = object.maker_order_is_buy;
    }
    if (object.maker_fee_ppm !== undefined && object.maker_fee_ppm !== null) {
      message.makerFeePpm = object.maker_fee_ppm;
    }
    if (object.clob_pair_id !== undefined && object.clob_pair_id !== null) {
      message.clobPairId = object.clob_pair_id;
    }
    if (object.fill_amount !== undefined && object.fill_amount !== null) {
      message.fillAmount = BigInt(object.fill_amount);
    }
    return message;
  },
  toAmino(message: MEVLiquidationMatch): MEVLiquidationMatchAmino {
    const obj: any = {};
    obj.liquidated_subaccount_id = message.liquidatedSubaccountId ? SubaccountId.toAmino(message.liquidatedSubaccountId) : undefined;
    obj.insurance_fund_delta_quote_quantums = message.insuranceFundDeltaQuoteQuantums ? message.insuranceFundDeltaQuoteQuantums.toString() : undefined;
    obj.maker_order_subaccount_id = message.makerOrderSubaccountId ? SubaccountId.toAmino(message.makerOrderSubaccountId) : undefined;
    obj.maker_order_subticks = message.makerOrderSubticks ? message.makerOrderSubticks.toString() : undefined;
    obj.maker_order_is_buy = message.makerOrderIsBuy;
    obj.maker_fee_ppm = message.makerFeePpm;
    obj.clob_pair_id = message.clobPairId;
    obj.fill_amount = message.fillAmount ? message.fillAmount.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MEVLiquidationMatchAminoMsg): MEVLiquidationMatch {
    return MEVLiquidationMatch.fromAmino(object.value);
  },
  fromProtoMsg(message: MEVLiquidationMatchProtoMsg): MEVLiquidationMatch {
    return MEVLiquidationMatch.decode(message.value);
  },
  toProto(message: MEVLiquidationMatch): Uint8Array {
    return MEVLiquidationMatch.encode(message).finish();
  },
  toProtoMsg(message: MEVLiquidationMatch): MEVLiquidationMatchProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MEVLiquidationMatch",
      value: MEVLiquidationMatch.encode(message).finish()
    };
  }
};
function createBaseClobMidPrice(): ClobMidPrice {
  return {
    clobPair: ClobPair.fromPartial({}),
    subticks: BigInt(0)
  };
}
export const ClobMidPrice = {
  typeUrl: "/dydxprotocol.clob.ClobMidPrice",
  encode(message: ClobMidPrice, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.clobPair !== undefined) {
      ClobPair.encode(message.clobPair, writer.uint32(10).fork()).ldelim();
    }
    if (message.subticks !== BigInt(0)) {
      writer.uint32(16).uint64(message.subticks);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ClobMidPrice {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClobMidPrice();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.clobPair = ClobPair.decode(reader, reader.uint32());
          break;
        case 2:
          message.subticks = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ClobMidPrice>): ClobMidPrice {
    const message = createBaseClobMidPrice();
    message.clobPair = object.clobPair !== undefined && object.clobPair !== null ? ClobPair.fromPartial(object.clobPair) : undefined;
    message.subticks = object.subticks !== undefined && object.subticks !== null ? BigInt(object.subticks.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: ClobMidPriceAmino): ClobMidPrice {
    const message = createBaseClobMidPrice();
    if (object.clob_pair !== undefined && object.clob_pair !== null) {
      message.clobPair = ClobPair.fromAmino(object.clob_pair);
    }
    if (object.subticks !== undefined && object.subticks !== null) {
      message.subticks = BigInt(object.subticks);
    }
    return message;
  },
  toAmino(message: ClobMidPrice): ClobMidPriceAmino {
    const obj: any = {};
    obj.clob_pair = message.clobPair ? ClobPair.toAmino(message.clobPair) : undefined;
    obj.subticks = message.subticks ? message.subticks.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: ClobMidPriceAminoMsg): ClobMidPrice {
    return ClobMidPrice.fromAmino(object.value);
  },
  fromProtoMsg(message: ClobMidPriceProtoMsg): ClobMidPrice {
    return ClobMidPrice.decode(message.value);
  },
  toProto(message: ClobMidPrice): Uint8Array {
    return ClobMidPrice.encode(message).finish();
  },
  toProtoMsg(message: ClobMidPrice): ClobMidPriceProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.ClobMidPrice",
      value: ClobMidPrice.encode(message).finish()
    };
  }
};
function createBaseValidatorMevMatches(): ValidatorMevMatches {
  return {
    matches: [],
    liquidationMatches: []
  };
}
export const ValidatorMevMatches = {
  typeUrl: "/dydxprotocol.clob.ValidatorMevMatches",
  encode(message: ValidatorMevMatches, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.matches) {
      MEVMatch.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.liquidationMatches) {
      MEVLiquidationMatch.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ValidatorMevMatches {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValidatorMevMatches();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.matches.push(MEVMatch.decode(reader, reader.uint32()));
          break;
        case 2:
          message.liquidationMatches.push(MEVLiquidationMatch.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ValidatorMevMatches>): ValidatorMevMatches {
    const message = createBaseValidatorMevMatches();
    message.matches = object.matches?.map(e => MEVMatch.fromPartial(e)) || [];
    message.liquidationMatches = object.liquidationMatches?.map(e => MEVLiquidationMatch.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: ValidatorMevMatchesAmino): ValidatorMevMatches {
    const message = createBaseValidatorMevMatches();
    message.matches = object.matches?.map(e => MEVMatch.fromAmino(e)) || [];
    message.liquidationMatches = object.liquidation_matches?.map(e => MEVLiquidationMatch.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: ValidatorMevMatches): ValidatorMevMatchesAmino {
    const obj: any = {};
    if (message.matches) {
      obj.matches = message.matches.map(e => e ? MEVMatch.toAmino(e) : undefined);
    } else {
      obj.matches = [];
    }
    if (message.liquidationMatches) {
      obj.liquidation_matches = message.liquidationMatches.map(e => e ? MEVLiquidationMatch.toAmino(e) : undefined);
    } else {
      obj.liquidation_matches = [];
    }
    return obj;
  },
  fromAminoMsg(object: ValidatorMevMatchesAminoMsg): ValidatorMevMatches {
    return ValidatorMevMatches.fromAmino(object.value);
  },
  fromProtoMsg(message: ValidatorMevMatchesProtoMsg): ValidatorMevMatches {
    return ValidatorMevMatches.decode(message.value);
  },
  toProto(message: ValidatorMevMatches): Uint8Array {
    return ValidatorMevMatches.encode(message).finish();
  },
  toProtoMsg(message: ValidatorMevMatches): ValidatorMevMatchesProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.ValidatorMevMatches",
      value: ValidatorMevMatches.encode(message).finish()
    };
  }
};
function createBaseMevNodeToNodeMetrics(): MevNodeToNodeMetrics {
  return {
    validatorMevMatches: undefined,
    clobMidPrices: [],
    bpMevMatches: undefined,
    proposalReceiveTime: BigInt(0)
  };
}
export const MevNodeToNodeMetrics = {
  typeUrl: "/dydxprotocol.clob.MevNodeToNodeMetrics",
  encode(message: MevNodeToNodeMetrics, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.validatorMevMatches !== undefined) {
      ValidatorMevMatches.encode(message.validatorMevMatches, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.clobMidPrices) {
      ClobMidPrice.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.bpMevMatches !== undefined) {
      ValidatorMevMatches.encode(message.bpMevMatches, writer.uint32(26).fork()).ldelim();
    }
    if (message.proposalReceiveTime !== BigInt(0)) {
      writer.uint32(32).uint64(message.proposalReceiveTime);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MevNodeToNodeMetrics {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMevNodeToNodeMetrics();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.validatorMevMatches = ValidatorMevMatches.decode(reader, reader.uint32());
          break;
        case 2:
          message.clobMidPrices.push(ClobMidPrice.decode(reader, reader.uint32()));
          break;
        case 3:
          message.bpMevMatches = ValidatorMevMatches.decode(reader, reader.uint32());
          break;
        case 4:
          message.proposalReceiveTime = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MevNodeToNodeMetrics>): MevNodeToNodeMetrics {
    const message = createBaseMevNodeToNodeMetrics();
    message.validatorMevMatches = object.validatorMevMatches !== undefined && object.validatorMevMatches !== null ? ValidatorMevMatches.fromPartial(object.validatorMevMatches) : undefined;
    message.clobMidPrices = object.clobMidPrices?.map(e => ClobMidPrice.fromPartial(e)) || [];
    message.bpMevMatches = object.bpMevMatches !== undefined && object.bpMevMatches !== null ? ValidatorMevMatches.fromPartial(object.bpMevMatches) : undefined;
    message.proposalReceiveTime = object.proposalReceiveTime !== undefined && object.proposalReceiveTime !== null ? BigInt(object.proposalReceiveTime.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MevNodeToNodeMetricsAmino): MevNodeToNodeMetrics {
    const message = createBaseMevNodeToNodeMetrics();
    if (object.validator_mev_matches !== undefined && object.validator_mev_matches !== null) {
      message.validatorMevMatches = ValidatorMevMatches.fromAmino(object.validator_mev_matches);
    }
    message.clobMidPrices = object.clob_mid_prices?.map(e => ClobMidPrice.fromAmino(e)) || [];
    if (object.bp_mev_matches !== undefined && object.bp_mev_matches !== null) {
      message.bpMevMatches = ValidatorMevMatches.fromAmino(object.bp_mev_matches);
    }
    if (object.proposal_receive_time !== undefined && object.proposal_receive_time !== null) {
      message.proposalReceiveTime = BigInt(object.proposal_receive_time);
    }
    return message;
  },
  toAmino(message: MevNodeToNodeMetrics): MevNodeToNodeMetricsAmino {
    const obj: any = {};
    obj.validator_mev_matches = message.validatorMevMatches ? ValidatorMevMatches.toAmino(message.validatorMevMatches) : undefined;
    if (message.clobMidPrices) {
      obj.clob_mid_prices = message.clobMidPrices.map(e => e ? ClobMidPrice.toAmino(e) : undefined);
    } else {
      obj.clob_mid_prices = [];
    }
    obj.bp_mev_matches = message.bpMevMatches ? ValidatorMevMatches.toAmino(message.bpMevMatches) : undefined;
    obj.proposal_receive_time = message.proposalReceiveTime ? message.proposalReceiveTime.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MevNodeToNodeMetricsAminoMsg): MevNodeToNodeMetrics {
    return MevNodeToNodeMetrics.fromAmino(object.value);
  },
  fromProtoMsg(message: MevNodeToNodeMetricsProtoMsg): MevNodeToNodeMetrics {
    return MevNodeToNodeMetrics.decode(message.value);
  },
  toProto(message: MevNodeToNodeMetrics): Uint8Array {
    return MevNodeToNodeMetrics.encode(message).finish();
  },
  toProtoMsg(message: MevNodeToNodeMetrics): MevNodeToNodeMetricsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MevNodeToNodeMetrics",
      value: MevNodeToNodeMetrics.encode(message).finish()
    };
  }
};