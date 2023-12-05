import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { ClobPair, ClobPairSDKType } from "./clob_pair";
import * as _m0 from "protobufjs/minimal";
import { Long, DeepPartial } from "../../helpers";
/** MEVMatch represents all necessary data to calculate MEV for a regular match. */

export interface MEVMatch {
  takerOrderSubaccountId?: SubaccountId;
  takerFeePpm: number;
  makerOrderSubaccountId?: SubaccountId;
  makerOrderSubticks: Long;
  makerOrderIsBuy: boolean;
  makerFeePpm: number;
  clobPairId: number;
  fillAmount: Long;
}
/** MEVMatch represents all necessary data to calculate MEV for a regular match. */

export interface MEVMatchSDKType {
  taker_order_subaccount_id?: SubaccountIdSDKType;
  taker_fee_ppm: number;
  maker_order_subaccount_id?: SubaccountIdSDKType;
  maker_order_subticks: Long;
  maker_order_is_buy: boolean;
  maker_fee_ppm: number;
  clob_pair_id: number;
  fill_amount: Long;
}
/**
 * MEVLiquidationMatch represents all necessary data to calculate MEV for a
 * liquidation.
 */

export interface MEVLiquidationMatch {
  liquidatedSubaccountId?: SubaccountId;
  insuranceFundDeltaQuoteQuantums: Long;
  makerOrderSubaccountId?: SubaccountId;
  makerOrderSubticks: Long;
  makerOrderIsBuy: boolean;
  makerFeePpm: number;
  clobPairId: number;
  fillAmount: Long;
}
/**
 * MEVLiquidationMatch represents all necessary data to calculate MEV for a
 * liquidation.
 */

export interface MEVLiquidationMatchSDKType {
  liquidated_subaccount_id?: SubaccountIdSDKType;
  insurance_fund_delta_quote_quantums: Long;
  maker_order_subaccount_id?: SubaccountIdSDKType;
  maker_order_subticks: Long;
  maker_order_is_buy: boolean;
  maker_fee_ppm: number;
  clob_pair_id: number;
  fill_amount: Long;
}
/** ClobMidPrice contains the mid price of a CLOB pair, represented by it's ID. */

export interface ClobMidPrice {
  clobPair?: ClobPair;
  subticks: Long;
}
/** ClobMidPrice contains the mid price of a CLOB pair, represented by it's ID. */

export interface ClobMidPriceSDKType {
  clob_pair?: ClobPairSDKType;
  subticks: Long;
}
/**
 * ValidatorMevMatches contains all matches from the validator's local
 * operations queue.
 */

export interface ValidatorMevMatches {
  matches: MEVMatch[];
  liquidationMatches: MEVLiquidationMatch[];
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
  proposalReceiveTime: Long;
}
/**
 * MevNodeToNodeMetrics is a data structure for encapsulating all MEV node <>
 * node metrics.
 */

export interface MevNodeToNodeMetricsSDKType {
  validator_mev_matches?: ValidatorMevMatchesSDKType;
  clob_mid_prices: ClobMidPriceSDKType[];
  bp_mev_matches?: ValidatorMevMatchesSDKType;
  proposal_receive_time: Long;
}

function createBaseMEVMatch(): MEVMatch {
  return {
    takerOrderSubaccountId: undefined,
    takerFeePpm: 0,
    makerOrderSubaccountId: undefined,
    makerOrderSubticks: Long.UZERO,
    makerOrderIsBuy: false,
    makerFeePpm: 0,
    clobPairId: 0,
    fillAmount: Long.UZERO
  };
}

export const MEVMatch = {
  encode(message: MEVMatch, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.takerOrderSubaccountId !== undefined) {
      SubaccountId.encode(message.takerOrderSubaccountId, writer.uint32(10).fork()).ldelim();
    }

    if (message.takerFeePpm !== 0) {
      writer.uint32(16).int32(message.takerFeePpm);
    }

    if (message.makerOrderSubaccountId !== undefined) {
      SubaccountId.encode(message.makerOrderSubaccountId, writer.uint32(26).fork()).ldelim();
    }

    if (!message.makerOrderSubticks.isZero()) {
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

    if (!message.fillAmount.isZero()) {
      writer.uint32(64).uint64(message.fillAmount);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MEVMatch {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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
          message.makerOrderSubticks = (reader.uint64() as Long);
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
          message.fillAmount = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MEVMatch>): MEVMatch {
    const message = createBaseMEVMatch();
    message.takerOrderSubaccountId = object.takerOrderSubaccountId !== undefined && object.takerOrderSubaccountId !== null ? SubaccountId.fromPartial(object.takerOrderSubaccountId) : undefined;
    message.takerFeePpm = object.takerFeePpm ?? 0;
    message.makerOrderSubaccountId = object.makerOrderSubaccountId !== undefined && object.makerOrderSubaccountId !== null ? SubaccountId.fromPartial(object.makerOrderSubaccountId) : undefined;
    message.makerOrderSubticks = object.makerOrderSubticks !== undefined && object.makerOrderSubticks !== null ? Long.fromValue(object.makerOrderSubticks) : Long.UZERO;
    message.makerOrderIsBuy = object.makerOrderIsBuy ?? false;
    message.makerFeePpm = object.makerFeePpm ?? 0;
    message.clobPairId = object.clobPairId ?? 0;
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? Long.fromValue(object.fillAmount) : Long.UZERO;
    return message;
  }

};

function createBaseMEVLiquidationMatch(): MEVLiquidationMatch {
  return {
    liquidatedSubaccountId: undefined,
    insuranceFundDeltaQuoteQuantums: Long.ZERO,
    makerOrderSubaccountId: undefined,
    makerOrderSubticks: Long.UZERO,
    makerOrderIsBuy: false,
    makerFeePpm: 0,
    clobPairId: 0,
    fillAmount: Long.UZERO
  };
}

export const MEVLiquidationMatch = {
  encode(message: MEVLiquidationMatch, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.liquidatedSubaccountId !== undefined) {
      SubaccountId.encode(message.liquidatedSubaccountId, writer.uint32(10).fork()).ldelim();
    }

    if (!message.insuranceFundDeltaQuoteQuantums.isZero()) {
      writer.uint32(16).int64(message.insuranceFundDeltaQuoteQuantums);
    }

    if (message.makerOrderSubaccountId !== undefined) {
      SubaccountId.encode(message.makerOrderSubaccountId, writer.uint32(26).fork()).ldelim();
    }

    if (!message.makerOrderSubticks.isZero()) {
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

    if (!message.fillAmount.isZero()) {
      writer.uint32(64).uint64(message.fillAmount);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MEVLiquidationMatch {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMEVLiquidationMatch();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.liquidatedSubaccountId = SubaccountId.decode(reader, reader.uint32());
          break;

        case 2:
          message.insuranceFundDeltaQuoteQuantums = (reader.int64() as Long);
          break;

        case 3:
          message.makerOrderSubaccountId = SubaccountId.decode(reader, reader.uint32());
          break;

        case 4:
          message.makerOrderSubticks = (reader.uint64() as Long);
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
          message.fillAmount = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MEVLiquidationMatch>): MEVLiquidationMatch {
    const message = createBaseMEVLiquidationMatch();
    message.liquidatedSubaccountId = object.liquidatedSubaccountId !== undefined && object.liquidatedSubaccountId !== null ? SubaccountId.fromPartial(object.liquidatedSubaccountId) : undefined;
    message.insuranceFundDeltaQuoteQuantums = object.insuranceFundDeltaQuoteQuantums !== undefined && object.insuranceFundDeltaQuoteQuantums !== null ? Long.fromValue(object.insuranceFundDeltaQuoteQuantums) : Long.ZERO;
    message.makerOrderSubaccountId = object.makerOrderSubaccountId !== undefined && object.makerOrderSubaccountId !== null ? SubaccountId.fromPartial(object.makerOrderSubaccountId) : undefined;
    message.makerOrderSubticks = object.makerOrderSubticks !== undefined && object.makerOrderSubticks !== null ? Long.fromValue(object.makerOrderSubticks) : Long.UZERO;
    message.makerOrderIsBuy = object.makerOrderIsBuy ?? false;
    message.makerFeePpm = object.makerFeePpm ?? 0;
    message.clobPairId = object.clobPairId ?? 0;
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? Long.fromValue(object.fillAmount) : Long.UZERO;
    return message;
  }

};

function createBaseClobMidPrice(): ClobMidPrice {
  return {
    clobPair: undefined,
    subticks: Long.UZERO
  };
}

export const ClobMidPrice = {
  encode(message: ClobMidPrice, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobPair !== undefined) {
      ClobPair.encode(message.clobPair, writer.uint32(10).fork()).ldelim();
    }

    if (!message.subticks.isZero()) {
      writer.uint32(16).uint64(message.subticks);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ClobMidPrice {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClobMidPrice();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPair = ClobPair.decode(reader, reader.uint32());
          break;

        case 2:
          message.subticks = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ClobMidPrice>): ClobMidPrice {
    const message = createBaseClobMidPrice();
    message.clobPair = object.clobPair !== undefined && object.clobPair !== null ? ClobPair.fromPartial(object.clobPair) : undefined;
    message.subticks = object.subticks !== undefined && object.subticks !== null ? Long.fromValue(object.subticks) : Long.UZERO;
    return message;
  }

};

function createBaseValidatorMevMatches(): ValidatorMevMatches {
  return {
    matches: [],
    liquidationMatches: []
  };
}

export const ValidatorMevMatches = {
  encode(message: ValidatorMevMatches, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.matches) {
      MEVMatch.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.liquidationMatches) {
      MEVLiquidationMatch.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorMevMatches {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<ValidatorMevMatches>): ValidatorMevMatches {
    const message = createBaseValidatorMevMatches();
    message.matches = object.matches?.map(e => MEVMatch.fromPartial(e)) || [];
    message.liquidationMatches = object.liquidationMatches?.map(e => MEVLiquidationMatch.fromPartial(e)) || [];
    return message;
  }

};

function createBaseMevNodeToNodeMetrics(): MevNodeToNodeMetrics {
  return {
    validatorMevMatches: undefined,
    clobMidPrices: [],
    bpMevMatches: undefined,
    proposalReceiveTime: Long.UZERO
  };
}

export const MevNodeToNodeMetrics = {
  encode(message: MevNodeToNodeMetrics, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.validatorMevMatches !== undefined) {
      ValidatorMevMatches.encode(message.validatorMevMatches, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.clobMidPrices) {
      ClobMidPrice.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    if (message.bpMevMatches !== undefined) {
      ValidatorMevMatches.encode(message.bpMevMatches, writer.uint32(26).fork()).ldelim();
    }

    if (!message.proposalReceiveTime.isZero()) {
      writer.uint32(32).uint64(message.proposalReceiveTime);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MevNodeToNodeMetrics {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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
          message.proposalReceiveTime = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MevNodeToNodeMetrics>): MevNodeToNodeMetrics {
    const message = createBaseMevNodeToNodeMetrics();
    message.validatorMevMatches = object.validatorMevMatches !== undefined && object.validatorMevMatches !== null ? ValidatorMevMatches.fromPartial(object.validatorMevMatches) : undefined;
    message.clobMidPrices = object.clobMidPrices?.map(e => ClobMidPrice.fromPartial(e)) || [];
    message.bpMevMatches = object.bpMevMatches !== undefined && object.bpMevMatches !== null ? ValidatorMevMatches.fromPartial(object.bpMevMatches) : undefined;
    message.proposalReceiveTime = object.proposalReceiveTime !== undefined && object.proposalReceiveTime !== null ? Long.fromValue(object.proposalReceiveTime) : Long.UZERO;
    return message;
  }

};