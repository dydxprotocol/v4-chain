import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { ValidatorMevMatches, ValidatorMevMatchesSDKType, MevNodeToNodeMetrics, MevNodeToNodeMetricsSDKType } from "./mev";
import { ClobPair, ClobPairSDKType } from "./clob_pair";
import { EquityTierLimitConfiguration, EquityTierLimitConfigurationSDKType } from "./equity_tier_limit_config";
import { Order, OrderSDKType } from "./order";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** QueryGetClobPairRequest is request type for the ClobPair method. */

export interface QueryGetClobPairRequest {
  /** QueryGetClobPairRequest is request type for the ClobPair method. */
  id: number;
}
/** QueryGetClobPairRequest is request type for the ClobPair method. */

export interface QueryGetClobPairRequestSDKType {
  /** QueryGetClobPairRequest is request type for the ClobPair method. */
  id: number;
}
/** QueryClobPairResponse is response type for the ClobPair method. */

export interface QueryClobPairResponse {
  clobPair?: ClobPair;
}
/** QueryClobPairResponse is response type for the ClobPair method. */

export interface QueryClobPairResponseSDKType {
  clob_pair?: ClobPairSDKType;
}
/** QueryAllClobPairRequest is request type for the ClobPairAll method. */

export interface QueryAllClobPairRequest {
  pagination?: PageRequest;
}
/** QueryAllClobPairRequest is request type for the ClobPairAll method. */

export interface QueryAllClobPairRequestSDKType {
  pagination?: PageRequestSDKType;
}
/** QueryClobPairAllResponse is response type for the ClobPairAll method. */

export interface QueryClobPairAllResponse {
  clobPair: ClobPair[];
  pagination?: PageResponse;
}
/** QueryClobPairAllResponse is response type for the ClobPairAll method. */

export interface QueryClobPairAllResponseSDKType {
  clob_pair: ClobPairSDKType[];
  pagination?: PageResponseSDKType;
}
/**
 * AreSubaccountsLiquidatableRequest is a request message used to check whether
 * the given subaccounts are liquidatable.
 * The subaccount ids should not contain duplicates.
 */

export interface AreSubaccountsLiquidatableRequest {
  subaccountIds: SubaccountId[];
}
/**
 * AreSubaccountsLiquidatableRequest is a request message used to check whether
 * the given subaccounts are liquidatable.
 * The subaccount ids should not contain duplicates.
 */

export interface AreSubaccountsLiquidatableRequestSDKType {
  subaccount_ids: SubaccountIdSDKType[];
}
/**
 * AreSubaccountsLiquidatableResponse is a response message that contains the
 * liquidation status for each subaccount.
 */

export interface AreSubaccountsLiquidatableResponse {
  results: AreSubaccountsLiquidatableResponse_Result[];
}
/**
 * AreSubaccountsLiquidatableResponse is a response message that contains the
 * liquidation status for each subaccount.
 */

export interface AreSubaccountsLiquidatableResponseSDKType {
  results: AreSubaccountsLiquidatableResponse_ResultSDKType[];
}
/** Result returns whether a subaccount should be liquidated. */

export interface AreSubaccountsLiquidatableResponse_Result {
  subaccountId?: SubaccountId;
  isLiquidatable: boolean;
}
/** Result returns whether a subaccount should be liquidated. */

export interface AreSubaccountsLiquidatableResponse_ResultSDKType {
  subaccount_id?: SubaccountIdSDKType;
  is_liquidatable: boolean;
}
/**
 * MevNodeToNodeCalculationRequest is a request message used to run the
 * MEV node <> node calculation.
 */

export interface MevNodeToNodeCalculationRequest {
  /**
   * Represents the matches on the "block proposer". Note that this field
   * does not need to be the actual block proposer's matches for a block, since
   * the MEV calculation logic is run with this nodes matches as the "block
   * proposer" matches.
   */
  blockProposerMatches?: ValidatorMevMatches;
  /** Represents the matches and mid-prices on the validator. */

  validatorMevMetrics?: MevNodeToNodeMetrics;
}
/**
 * MevNodeToNodeCalculationRequest is a request message used to run the
 * MEV node <> node calculation.
 */

export interface MevNodeToNodeCalculationRequestSDKType {
  /**
   * Represents the matches on the "block proposer". Note that this field
   * does not need to be the actual block proposer's matches for a block, since
   * the MEV calculation logic is run with this nodes matches as the "block
   * proposer" matches.
   */
  block_proposer_matches?: ValidatorMevMatchesSDKType;
  /** Represents the matches and mid-prices on the validator. */

  validator_mev_metrics?: MevNodeToNodeMetricsSDKType;
}
/**
 * MevNodeToNodeCalculationResponse is a response message that contains the
 * MEV node <> node calculation result.
 */

export interface MevNodeToNodeCalculationResponse {
  results: MevNodeToNodeCalculationResponse_MevAndVolumePerClob[];
}
/**
 * MevNodeToNodeCalculationResponse is a response message that contains the
 * MEV node <> node calculation result.
 */

export interface MevNodeToNodeCalculationResponseSDKType {
  results: MevNodeToNodeCalculationResponse_MevAndVolumePerClobSDKType[];
}
/** MevAndVolumePerClob contains information about the MEV and volume per CLOB. */

export interface MevNodeToNodeCalculationResponse_MevAndVolumePerClob {
  clobPairId: number;
  mev: number;
  volume: Long;
}
/** MevAndVolumePerClob contains information about the MEV and volume per CLOB. */

export interface MevNodeToNodeCalculationResponse_MevAndVolumePerClobSDKType {
  clob_pair_id: number;
  mev: number;
  volume: Long;
}
/**
 * QueryEquityTierLimitConfigurationRequest is a request message for
 * EquityTierLimitConfiguration.
 */

export interface QueryEquityTierLimitConfigurationRequest {}
/**
 * QueryEquityTierLimitConfigurationRequest is a request message for
 * EquityTierLimitConfiguration.
 */

export interface QueryEquityTierLimitConfigurationRequestSDKType {}
/**
 * QueryEquityTierLimitConfigurationResponse is a response message that contains
 * the EquityTierLimitConfiguration.
 */

export interface QueryEquityTierLimitConfigurationResponse {
  equityTierLimitConfig?: EquityTierLimitConfiguration;
}
/**
 * QueryEquityTierLimitConfigurationResponse is a response message that contains
 * the EquityTierLimitConfiguration.
 */

export interface QueryEquityTierLimitConfigurationResponseSDKType {
  equity_tier_limit_config?: EquityTierLimitConfigurationSDKType;
}
/** QueryAllStatefulOrdersRequest is a request message for AllStatefulOrders. */

export interface QueryAllStatefulOrdersRequest {
  pagination?: PageRequest;
}
/** QueryAllStatefulOrdersRequest is a request message for AllStatefulOrders. */

export interface QueryAllStatefulOrdersRequestSDKType {
  pagination?: PageRequestSDKType;
}
/**
 * QueryAllStateOrdersResponse is a response message that contains all stateful
 * orders.
 */

export interface QueryAllStatefulOrdersResponse {
  statefulOrders: Order[];
  pagination?: PageResponse;
}
/**
 * QueryAllStateOrdersResponse is a response message that contains all stateful
 * orders.
 */

export interface QueryAllStatefulOrdersResponseSDKType {
  stateful_orders: OrderSDKType[];
  pagination?: PageResponseSDKType;
}
/** QueryStatefulOrderCountRequest is a request message for StatefulOrderCount. */

export interface QueryStatefulOrderCountRequest {
  subaccountId?: SubaccountId;
}
/** QueryStatefulOrderCountRequest is a request message for StatefulOrderCount. */

export interface QueryStatefulOrderCountRequestSDKType {
  subaccount_id?: SubaccountIdSDKType;
}
/**
 * QueryStatefulOrderCountResponse is a response message for StatefulOrderCount
 * that contains the count of all stateful orders.
 */

export interface QueryStatefulOrderCountResponse {
  /**
   * QueryStatefulOrderCountResponse is a response message for StatefulOrderCount
   * that contains the count of all stateful orders.
   */
  count: number;
}
/**
 * QueryStatefulOrderCountResponse is a response message for StatefulOrderCount
 * that contains the count of all stateful orders.
 */

export interface QueryStatefulOrderCountResponseSDKType {
  /**
   * QueryStatefulOrderCountResponse is a response message for StatefulOrderCount
   * that contains the count of all stateful orders.
   */
  count: number;
}

function createBaseQueryGetClobPairRequest(): QueryGetClobPairRequest {
  return {
    id: 0
  };
}

export const QueryGetClobPairRequest = {
  encode(message: QueryGetClobPairRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetClobPairRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetClobPairRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryGetClobPairRequest>): QueryGetClobPairRequest {
    const message = createBaseQueryGetClobPairRequest();
    message.id = object.id ?? 0;
    return message;
  }

};

function createBaseQueryClobPairResponse(): QueryClobPairResponse {
  return {
    clobPair: undefined
  };
}

export const QueryClobPairResponse = {
  encode(message: QueryClobPairResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobPair !== undefined) {
      ClobPair.encode(message.clobPair, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClobPairResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClobPairResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPair = ClobPair.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryClobPairResponse>): QueryClobPairResponse {
    const message = createBaseQueryClobPairResponse();
    message.clobPair = object.clobPair !== undefined && object.clobPair !== null ? ClobPair.fromPartial(object.clobPair) : undefined;
    return message;
  }

};

function createBaseQueryAllClobPairRequest(): QueryAllClobPairRequest {
  return {
    pagination: undefined
  };
}

export const QueryAllClobPairRequest = {
  encode(message: QueryAllClobPairRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllClobPairRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllClobPairRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllClobPairRequest>): QueryAllClobPairRequest {
    const message = createBaseQueryAllClobPairRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryClobPairAllResponse(): QueryClobPairAllResponse {
  return {
    clobPair: [],
    pagination: undefined
  };
}

export const QueryClobPairAllResponse = {
  encode(message: QueryClobPairAllResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.clobPair) {
      ClobPair.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClobPairAllResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClobPairAllResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPair.push(ClobPair.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryClobPairAllResponse>): QueryClobPairAllResponse {
    const message = createBaseQueryClobPairAllResponse();
    message.clobPair = object.clobPair?.map(e => ClobPair.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseAreSubaccountsLiquidatableRequest(): AreSubaccountsLiquidatableRequest {
  return {
    subaccountIds: []
  };
}

export const AreSubaccountsLiquidatableRequest = {
  encode(message: AreSubaccountsLiquidatableRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.subaccountIds) {
      SubaccountId.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AreSubaccountsLiquidatableRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAreSubaccountsLiquidatableRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountIds.push(SubaccountId.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AreSubaccountsLiquidatableRequest>): AreSubaccountsLiquidatableRequest {
    const message = createBaseAreSubaccountsLiquidatableRequest();
    message.subaccountIds = object.subaccountIds?.map(e => SubaccountId.fromPartial(e)) || [];
    return message;
  }

};

function createBaseAreSubaccountsLiquidatableResponse(): AreSubaccountsLiquidatableResponse {
  return {
    results: []
  };
}

export const AreSubaccountsLiquidatableResponse = {
  encode(message: AreSubaccountsLiquidatableResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.results) {
      AreSubaccountsLiquidatableResponse_Result.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AreSubaccountsLiquidatableResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAreSubaccountsLiquidatableResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.results.push(AreSubaccountsLiquidatableResponse_Result.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AreSubaccountsLiquidatableResponse>): AreSubaccountsLiquidatableResponse {
    const message = createBaseAreSubaccountsLiquidatableResponse();
    message.results = object.results?.map(e => AreSubaccountsLiquidatableResponse_Result.fromPartial(e)) || [];
    return message;
  }

};

function createBaseAreSubaccountsLiquidatableResponse_Result(): AreSubaccountsLiquidatableResponse_Result {
  return {
    subaccountId: undefined,
    isLiquidatable: false
  };
}

export const AreSubaccountsLiquidatableResponse_Result = {
  encode(message: AreSubaccountsLiquidatableResponse_Result, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.subaccountId !== undefined) {
      SubaccountId.encode(message.subaccountId, writer.uint32(10).fork()).ldelim();
    }

    if (message.isLiquidatable === true) {
      writer.uint32(16).bool(message.isLiquidatable);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AreSubaccountsLiquidatableResponse_Result {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAreSubaccountsLiquidatableResponse_Result();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountId = SubaccountId.decode(reader, reader.uint32());
          break;

        case 2:
          message.isLiquidatable = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AreSubaccountsLiquidatableResponse_Result>): AreSubaccountsLiquidatableResponse_Result {
    const message = createBaseAreSubaccountsLiquidatableResponse_Result();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? SubaccountId.fromPartial(object.subaccountId) : undefined;
    message.isLiquidatable = object.isLiquidatable ?? false;
    return message;
  }

};

function createBaseMevNodeToNodeCalculationRequest(): MevNodeToNodeCalculationRequest {
  return {
    blockProposerMatches: undefined,
    validatorMevMetrics: undefined
  };
}

export const MevNodeToNodeCalculationRequest = {
  encode(message: MevNodeToNodeCalculationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockProposerMatches !== undefined) {
      ValidatorMevMatches.encode(message.blockProposerMatches, writer.uint32(10).fork()).ldelim();
    }

    if (message.validatorMevMetrics !== undefined) {
      MevNodeToNodeMetrics.encode(message.validatorMevMetrics, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MevNodeToNodeCalculationRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMevNodeToNodeCalculationRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.blockProposerMatches = ValidatorMevMatches.decode(reader, reader.uint32());
          break;

        case 2:
          message.validatorMevMetrics = MevNodeToNodeMetrics.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MevNodeToNodeCalculationRequest>): MevNodeToNodeCalculationRequest {
    const message = createBaseMevNodeToNodeCalculationRequest();
    message.blockProposerMatches = object.blockProposerMatches !== undefined && object.blockProposerMatches !== null ? ValidatorMevMatches.fromPartial(object.blockProposerMatches) : undefined;
    message.validatorMevMetrics = object.validatorMevMetrics !== undefined && object.validatorMevMetrics !== null ? MevNodeToNodeMetrics.fromPartial(object.validatorMevMetrics) : undefined;
    return message;
  }

};

function createBaseMevNodeToNodeCalculationResponse(): MevNodeToNodeCalculationResponse {
  return {
    results: []
  };
}

export const MevNodeToNodeCalculationResponse = {
  encode(message: MevNodeToNodeCalculationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.results) {
      MevNodeToNodeCalculationResponse_MevAndVolumePerClob.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MevNodeToNodeCalculationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMevNodeToNodeCalculationResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.results.push(MevNodeToNodeCalculationResponse_MevAndVolumePerClob.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MevNodeToNodeCalculationResponse>): MevNodeToNodeCalculationResponse {
    const message = createBaseMevNodeToNodeCalculationResponse();
    message.results = object.results?.map(e => MevNodeToNodeCalculationResponse_MevAndVolumePerClob.fromPartial(e)) || [];
    return message;
  }

};

function createBaseMevNodeToNodeCalculationResponse_MevAndVolumePerClob(): MevNodeToNodeCalculationResponse_MevAndVolumePerClob {
  return {
    clobPairId: 0,
    mev: 0,
    volume: Long.UZERO
  };
}

export const MevNodeToNodeCalculationResponse_MevAndVolumePerClob = {
  encode(message: MevNodeToNodeCalculationResponse_MevAndVolumePerClob, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobPairId !== 0) {
      writer.uint32(8).uint32(message.clobPairId);
    }

    if (message.mev !== 0) {
      writer.uint32(21).float(message.mev);
    }

    if (!message.volume.isZero()) {
      writer.uint32(24).uint64(message.volume);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MevNodeToNodeCalculationResponse_MevAndVolumePerClob {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMevNodeToNodeCalculationResponse_MevAndVolumePerClob();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPairId = reader.uint32();
          break;

        case 2:
          message.mev = reader.float();
          break;

        case 3:
          message.volume = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MevNodeToNodeCalculationResponse_MevAndVolumePerClob>): MevNodeToNodeCalculationResponse_MevAndVolumePerClob {
    const message = createBaseMevNodeToNodeCalculationResponse_MevAndVolumePerClob();
    message.clobPairId = object.clobPairId ?? 0;
    message.mev = object.mev ?? 0;
    message.volume = object.volume !== undefined && object.volume !== null ? Long.fromValue(object.volume) : Long.UZERO;
    return message;
  }

};

function createBaseQueryEquityTierLimitConfigurationRequest(): QueryEquityTierLimitConfigurationRequest {
  return {};
}

export const QueryEquityTierLimitConfigurationRequest = {
  encode(_: QueryEquityTierLimitConfigurationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryEquityTierLimitConfigurationRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryEquityTierLimitConfigurationRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<QueryEquityTierLimitConfigurationRequest>): QueryEquityTierLimitConfigurationRequest {
    const message = createBaseQueryEquityTierLimitConfigurationRequest();
    return message;
  }

};

function createBaseQueryEquityTierLimitConfigurationResponse(): QueryEquityTierLimitConfigurationResponse {
  return {
    equityTierLimitConfig: undefined
  };
}

export const QueryEquityTierLimitConfigurationResponse = {
  encode(message: QueryEquityTierLimitConfigurationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.equityTierLimitConfig !== undefined) {
      EquityTierLimitConfiguration.encode(message.equityTierLimitConfig, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryEquityTierLimitConfigurationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryEquityTierLimitConfigurationResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.equityTierLimitConfig = EquityTierLimitConfiguration.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryEquityTierLimitConfigurationResponse>): QueryEquityTierLimitConfigurationResponse {
    const message = createBaseQueryEquityTierLimitConfigurationResponse();
    message.equityTierLimitConfig = object.equityTierLimitConfig !== undefined && object.equityTierLimitConfig !== null ? EquityTierLimitConfiguration.fromPartial(object.equityTierLimitConfig) : undefined;
    return message;
  }

};

function createBaseQueryAllStatefulOrdersRequest(): QueryAllStatefulOrdersRequest {
  return {
    pagination: undefined
  };
}

export const QueryAllStatefulOrdersRequest = {
  encode(message: QueryAllStatefulOrdersRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllStatefulOrdersRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllStatefulOrdersRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllStatefulOrdersRequest>): QueryAllStatefulOrdersRequest {
    const message = createBaseQueryAllStatefulOrdersRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryAllStatefulOrdersResponse(): QueryAllStatefulOrdersResponse {
  return {
    statefulOrders: [],
    pagination: undefined
  };
}

export const QueryAllStatefulOrdersResponse = {
  encode(message: QueryAllStatefulOrdersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.statefulOrders) {
      Order.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllStatefulOrdersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllStatefulOrdersResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.statefulOrders.push(Order.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllStatefulOrdersResponse>): QueryAllStatefulOrdersResponse {
    const message = createBaseQueryAllStatefulOrdersResponse();
    message.statefulOrders = object.statefulOrders?.map(e => Order.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryStatefulOrderCountRequest(): QueryStatefulOrderCountRequest {
  return {
    subaccountId: undefined
  };
}

export const QueryStatefulOrderCountRequest = {
  encode(message: QueryStatefulOrderCountRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.subaccountId !== undefined) {
      SubaccountId.encode(message.subaccountId, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryStatefulOrderCountRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryStatefulOrderCountRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountId = SubaccountId.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryStatefulOrderCountRequest>): QueryStatefulOrderCountRequest {
    const message = createBaseQueryStatefulOrderCountRequest();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? SubaccountId.fromPartial(object.subaccountId) : undefined;
    return message;
  }

};

function createBaseQueryStatefulOrderCountResponse(): QueryStatefulOrderCountResponse {
  return {
    count: 0
  };
}

export const QueryStatefulOrderCountResponse = {
  encode(message: QueryStatefulOrderCountResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.count !== 0) {
      writer.uint32(8).uint32(message.count);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryStatefulOrderCountResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryStatefulOrderCountResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.count = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryStatefulOrderCountResponse>): QueryStatefulOrderCountResponse {
    const message = createBaseQueryStatefulOrderCountResponse();
    message.count = object.count ?? 0;
    return message;
  }

};