import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { ValidatorMevMatches, ValidatorMevMatchesSDKType, MevNodeToNodeMetrics, MevNodeToNodeMetricsSDKType } from "./mev";
import { ClobPair, ClobPairSDKType } from "./clob_pair";
import { EquityTierLimitConfiguration, EquityTierLimitConfigurationSDKType } from "./equity_tier_limit_config";
import { BlockRateLimitConfiguration, BlockRateLimitConfigurationSDKType } from "./block_rate_limit_config";
import { LiquidationsConfig, LiquidationsConfigSDKType } from "./liquidations_config";
import { OffChainUpdateV1, OffChainUpdateV1SDKType } from "../indexer/off_chain_updates/off_chain_updates";
import { ClobMatch, ClobMatchSDKType } from "./matches";
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
/**
 * QueryBlockRateLimitConfigurationRequest is a request message for
 * BlockRateLimitConfiguration.
 */

export interface QueryBlockRateLimitConfigurationRequest {}
/**
 * QueryBlockRateLimitConfigurationRequest is a request message for
 * BlockRateLimitConfiguration.
 */

export interface QueryBlockRateLimitConfigurationRequestSDKType {}
/**
 * QueryBlockRateLimitConfigurationResponse is a response message that contains
 * the BlockRateLimitConfiguration.
 */

export interface QueryBlockRateLimitConfigurationResponse {
  blockRateLimitConfig?: BlockRateLimitConfiguration;
}
/**
 * QueryBlockRateLimitConfigurationResponse is a response message that contains
 * the BlockRateLimitConfiguration.
 */

export interface QueryBlockRateLimitConfigurationResponseSDKType {
  block_rate_limit_config?: BlockRateLimitConfigurationSDKType;
}
/**
 * QueryLiquidationsConfigurationRequest is a request message for
 * LiquidationsConfiguration.
 */

export interface QueryLiquidationsConfigurationRequest {}
/**
 * QueryLiquidationsConfigurationRequest is a request message for
 * LiquidationsConfiguration.
 */

export interface QueryLiquidationsConfigurationRequestSDKType {}
/**
 * QueryLiquidationsConfigurationResponse is a response message that contains
 * the LiquidationsConfiguration.
 */

export interface QueryLiquidationsConfigurationResponse {
  liquidationsConfig?: LiquidationsConfig;
}
/**
 * QueryLiquidationsConfigurationResponse is a response message that contains
 * the LiquidationsConfiguration.
 */

export interface QueryLiquidationsConfigurationResponseSDKType {
  liquidations_config?: LiquidationsConfigSDKType;
}
/**
 * StreamOrderbookUpdatesRequest is a request message for the
 * StreamOrderbookUpdates method.
 */

export interface StreamOrderbookUpdatesRequest {
  /** Clob pair ids to stream orderbook updates for. */
  clobPairId: number[];
}
/**
 * StreamOrderbookUpdatesRequest is a request message for the
 * StreamOrderbookUpdates method.
 */

export interface StreamOrderbookUpdatesRequestSDKType {
  /** Clob pair ids to stream orderbook updates for. */
  clob_pair_id: number[];
}
/**
 * StreamOrderbookUpdatesResponse is a response message for the
 * StreamOrderbookUpdates method.
 */

export interface StreamOrderbookUpdatesResponse {
  /** Orderbook updates for the clob pair. */
  updates: StreamUpdate[];
  /**
   * ---Additional fields used to debug issues---
   * Block height of the updates.
   */

  blockHeight: number;
  /** Exec mode of the updates. */

  execMode: number;
}
/**
 * StreamOrderbookUpdatesResponse is a response message for the
 * StreamOrderbookUpdates method.
 */

export interface StreamOrderbookUpdatesResponseSDKType {
  /** Orderbook updates for the clob pair. */
  updates: StreamUpdateSDKType[];
  /**
   * ---Additional fields used to debug issues---
   * Block height of the updates.
   */

  block_height: number;
  /** Exec mode of the updates. */

  exec_mode: number;
}
/**
 * StreamUpdate is an update that will be pushed through the
 * GRPC stream.
 */

export interface StreamUpdate {
  orderbookUpdate?: StreamOrderbookUpdate;
  orderFill?: StreamOrderbookFill;
}
/**
 * StreamUpdate is an update that will be pushed through the
 * GRPC stream.
 */

export interface StreamUpdateSDKType {
  orderbook_update?: StreamOrderbookUpdateSDKType;
  order_fill?: StreamOrderbookFillSDKType;
}
/**
 * StreamOrderbookUpdate provides information on an orderbook update. Used in
 * the full node GRPC stream.
 */

export interface StreamOrderbookUpdate {
  /**
   * Orderbook updates for the clob pair. Can contain order place, removals,
   * or updates.
   */
  updates: OffChainUpdateV1[];
  /**
   * Snapshot indicates if the response is from a snapshot of the orderbook.
   * This is true for the initial response and false for all subsequent updates.
   * Note that if the snapshot is true, then all previous entries should be
   * discarded and the orderbook should be resynced.
   */

  snapshot: boolean;
}
/**
 * StreamOrderbookUpdate provides information on an orderbook update. Used in
 * the full node GRPC stream.
 */

export interface StreamOrderbookUpdateSDKType {
  /**
   * Orderbook updates for the clob pair. Can contain order place, removals,
   * or updates.
   */
  updates: OffChainUpdateV1SDKType[];
  /**
   * Snapshot indicates if the response is from a snapshot of the orderbook.
   * This is true for the initial response and false for all subsequent updates.
   * Note that if the snapshot is true, then all previous entries should be
   * discarded and the orderbook should be resynced.
   */

  snapshot: boolean;
}
/**
 * StreamOrderbookFill provides information on an orderbook fill. Used in
 * the full node GRPC stream.
 */

export interface StreamOrderbookFill {
  /**
   * Clob match. Provides information on which orders were matched
   * and the type of order. Fill amounts here are relative.
   */
  clobMatch?: ClobMatch;
  /**
   * All orders involved in the specified clob match. Used to look up
   * price of a match through a given maker order id.
   */

  orders: Order[];
  /** Resulting fill amounts for each order in the orders array. */

  fillAmounts: Long[];
}
/**
 * StreamOrderbookFill provides information on an orderbook fill. Used in
 * the full node GRPC stream.
 */

export interface StreamOrderbookFillSDKType {
  /**
   * Clob match. Provides information on which orders were matched
   * and the type of order. Fill amounts here are relative.
   */
  clob_match?: ClobMatchSDKType;
  /**
   * All orders involved in the specified clob match. Used to look up
   * price of a match through a given maker order id.
   */

  orders: OrderSDKType[];
  /** Resulting fill amounts for each order in the orders array. */

  fill_amounts: Long[];
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

function createBaseQueryBlockRateLimitConfigurationRequest(): QueryBlockRateLimitConfigurationRequest {
  return {};
}

export const QueryBlockRateLimitConfigurationRequest = {
  encode(_: QueryBlockRateLimitConfigurationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryBlockRateLimitConfigurationRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryBlockRateLimitConfigurationRequest();

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

  fromPartial(_: DeepPartial<QueryBlockRateLimitConfigurationRequest>): QueryBlockRateLimitConfigurationRequest {
    const message = createBaseQueryBlockRateLimitConfigurationRequest();
    return message;
  }

};

function createBaseQueryBlockRateLimitConfigurationResponse(): QueryBlockRateLimitConfigurationResponse {
  return {
    blockRateLimitConfig: undefined
  };
}

export const QueryBlockRateLimitConfigurationResponse = {
  encode(message: QueryBlockRateLimitConfigurationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockRateLimitConfig !== undefined) {
      BlockRateLimitConfiguration.encode(message.blockRateLimitConfig, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryBlockRateLimitConfigurationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryBlockRateLimitConfigurationResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.blockRateLimitConfig = BlockRateLimitConfiguration.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryBlockRateLimitConfigurationResponse>): QueryBlockRateLimitConfigurationResponse {
    const message = createBaseQueryBlockRateLimitConfigurationResponse();
    message.blockRateLimitConfig = object.blockRateLimitConfig !== undefined && object.blockRateLimitConfig !== null ? BlockRateLimitConfiguration.fromPartial(object.blockRateLimitConfig) : undefined;
    return message;
  }

};

function createBaseQueryLiquidationsConfigurationRequest(): QueryLiquidationsConfigurationRequest {
  return {};
}

export const QueryLiquidationsConfigurationRequest = {
  encode(_: QueryLiquidationsConfigurationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryLiquidationsConfigurationRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryLiquidationsConfigurationRequest();

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

  fromPartial(_: DeepPartial<QueryLiquidationsConfigurationRequest>): QueryLiquidationsConfigurationRequest {
    const message = createBaseQueryLiquidationsConfigurationRequest();
    return message;
  }

};

function createBaseQueryLiquidationsConfigurationResponse(): QueryLiquidationsConfigurationResponse {
  return {
    liquidationsConfig: undefined
  };
}

export const QueryLiquidationsConfigurationResponse = {
  encode(message: QueryLiquidationsConfigurationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.liquidationsConfig !== undefined) {
      LiquidationsConfig.encode(message.liquidationsConfig, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryLiquidationsConfigurationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryLiquidationsConfigurationResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.liquidationsConfig = LiquidationsConfig.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryLiquidationsConfigurationResponse>): QueryLiquidationsConfigurationResponse {
    const message = createBaseQueryLiquidationsConfigurationResponse();
    message.liquidationsConfig = object.liquidationsConfig !== undefined && object.liquidationsConfig !== null ? LiquidationsConfig.fromPartial(object.liquidationsConfig) : undefined;
    return message;
  }

};

function createBaseStreamOrderbookUpdatesRequest(): StreamOrderbookUpdatesRequest {
  return {
    clobPairId: []
  };
}

export const StreamOrderbookUpdatesRequest = {
  encode(message: StreamOrderbookUpdatesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    writer.uint32(10).fork();

    for (const v of message.clobPairId) {
      writer.uint32(v);
    }

    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamOrderbookUpdatesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamOrderbookUpdatesRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;

            while (reader.pos < end2) {
              message.clobPairId.push(reader.uint32());
            }
          } else {
            message.clobPairId.push(reader.uint32());
          }

          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StreamOrderbookUpdatesRequest>): StreamOrderbookUpdatesRequest {
    const message = createBaseStreamOrderbookUpdatesRequest();
    message.clobPairId = object.clobPairId?.map(e => e) || [];
    return message;
  }

};

function createBaseStreamOrderbookUpdatesResponse(): StreamOrderbookUpdatesResponse {
  return {
    updates: [],
    blockHeight: 0,
    execMode: 0
  };
}

export const StreamOrderbookUpdatesResponse = {
  encode(message: StreamOrderbookUpdatesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.updates) {
      StreamUpdate.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.blockHeight !== 0) {
      writer.uint32(16).uint32(message.blockHeight);
    }

    if (message.execMode !== 0) {
      writer.uint32(24).uint32(message.execMode);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamOrderbookUpdatesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamOrderbookUpdatesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.updates.push(StreamUpdate.decode(reader, reader.uint32()));
          break;

        case 2:
          message.blockHeight = reader.uint32();
          break;

        case 3:
          message.execMode = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StreamOrderbookUpdatesResponse>): StreamOrderbookUpdatesResponse {
    const message = createBaseStreamOrderbookUpdatesResponse();
    message.updates = object.updates?.map(e => StreamUpdate.fromPartial(e)) || [];
    message.blockHeight = object.blockHeight ?? 0;
    message.execMode = object.execMode ?? 0;
    return message;
  }

};

function createBaseStreamUpdate(): StreamUpdate {
  return {
    orderbookUpdate: undefined,
    orderFill: undefined
  };
}

export const StreamUpdate = {
  encode(message: StreamUpdate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orderbookUpdate !== undefined) {
      StreamOrderbookUpdate.encode(message.orderbookUpdate, writer.uint32(10).fork()).ldelim();
    }

    if (message.orderFill !== undefined) {
      StreamOrderbookFill.encode(message.orderFill, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamUpdate {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamUpdate();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.orderbookUpdate = StreamOrderbookUpdate.decode(reader, reader.uint32());
          break;

        case 2:
          message.orderFill = StreamOrderbookFill.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StreamUpdate>): StreamUpdate {
    const message = createBaseStreamUpdate();
    message.orderbookUpdate = object.orderbookUpdate !== undefined && object.orderbookUpdate !== null ? StreamOrderbookUpdate.fromPartial(object.orderbookUpdate) : undefined;
    message.orderFill = object.orderFill !== undefined && object.orderFill !== null ? StreamOrderbookFill.fromPartial(object.orderFill) : undefined;
    return message;
  }

};

function createBaseStreamOrderbookUpdate(): StreamOrderbookUpdate {
  return {
    updates: [],
    snapshot: false
  };
}

export const StreamOrderbookUpdate = {
  encode(message: StreamOrderbookUpdate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.updates) {
      OffChainUpdateV1.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.snapshot === true) {
      writer.uint32(16).bool(message.snapshot);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamOrderbookUpdate {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamOrderbookUpdate();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.updates.push(OffChainUpdateV1.decode(reader, reader.uint32()));
          break;

        case 2:
          message.snapshot = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StreamOrderbookUpdate>): StreamOrderbookUpdate {
    const message = createBaseStreamOrderbookUpdate();
    message.updates = object.updates?.map(e => OffChainUpdateV1.fromPartial(e)) || [];
    message.snapshot = object.snapshot ?? false;
    return message;
  }

};

function createBaseStreamOrderbookFill(): StreamOrderbookFill {
  return {
    clobMatch: undefined,
    orders: [],
    fillAmounts: []
  };
}

export const StreamOrderbookFill = {
  encode(message: StreamOrderbookFill, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobMatch !== undefined) {
      ClobMatch.encode(message.clobMatch, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.orders) {
      Order.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    writer.uint32(26).fork();

    for (const v of message.fillAmounts) {
      writer.uint64(v);
    }

    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamOrderbookFill {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamOrderbookFill();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobMatch = ClobMatch.decode(reader, reader.uint32());
          break;

        case 2:
          message.orders.push(Order.decode(reader, reader.uint32()));
          break;

        case 3:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;

            while (reader.pos < end2) {
              message.fillAmounts.push((reader.uint64() as Long));
            }
          } else {
            message.fillAmounts.push((reader.uint64() as Long));
          }

          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StreamOrderbookFill>): StreamOrderbookFill {
    const message = createBaseStreamOrderbookFill();
    message.clobMatch = object.clobMatch !== undefined && object.clobMatch !== null ? ClobMatch.fromPartial(object.clobMatch) : undefined;
    message.orders = object.orders?.map(e => Order.fromPartial(e)) || [];
    message.fillAmounts = object.fillAmounts?.map(e => Long.fromValue(e)) || [];
    return message;
  }

};