import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { ValidatorMevMatches, ValidatorMevMatchesSDKType, MevNodeToNodeMetrics, MevNodeToNodeMetricsSDKType } from "./mev";
import { OrderId, OrderIdSDKType, LongTermOrderPlacement, LongTermOrderPlacementSDKType, Order, OrderSDKType, StreamLiquidationOrder, StreamLiquidationOrderSDKType } from "./order";
import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { ClobPair, ClobPairSDKType } from "./clob_pair";
import { EquityTierLimitConfiguration, EquityTierLimitConfigurationSDKType } from "./equity_tier_limit_config";
import { BlockRateLimitConfiguration, BlockRateLimitConfigurationSDKType } from "./block_rate_limit_config";
import { LiquidationsConfig, LiquidationsConfigSDKType } from "./liquidations_config";
import { StreamSubaccountUpdate, StreamSubaccountUpdateSDKType } from "../subaccounts/streaming";
import { StreamPriceUpdate, StreamPriceUpdateSDKType } from "../prices/streaming";
import { OffChainUpdateV1, OffChainUpdateV1SDKType } from "../indexer/off_chain_updates/off_chain_updates";
import { ClobMatch, ClobMatchSDKType } from "./matches";
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
/** QueryStatefulOrderRequest is a request message for StatefulOrder. */

export interface QueryStatefulOrderRequest {
  /** Order id to query. */
  orderId?: OrderId;
}
/** QueryStatefulOrderRequest is a request message for StatefulOrder. */

export interface QueryStatefulOrderRequestSDKType {
  /** Order id to query. */
  order_id?: OrderIdSDKType;
}
/**
 * QueryStatefulOrderResponse is a response message that contains the stateful
 * order.
 */

export interface QueryStatefulOrderResponse {
  /** Stateful order placement. */
  orderPlacement?: LongTermOrderPlacement;
  /** Fill amounts. */

  fillAmount: Long;
  /** Triggered status. */

  triggered: boolean;
}
/**
 * QueryStatefulOrderResponse is a response message that contains the stateful
 * order.
 */

export interface QueryStatefulOrderResponseSDKType {
  /** Stateful order placement. */
  order_placement?: LongTermOrderPlacementSDKType;
  /** Fill amounts. */

  fill_amount: Long;
  /** Triggered status. */

  triggered: boolean;
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
/** QueryNextClobPairIdRequest is a request message for the next clob pair id */

export interface QueryNextClobPairIdRequest {}
/** QueryNextClobPairIdRequest is a request message for the next clob pair id */

export interface QueryNextClobPairIdRequestSDKType {}
/** QueryNextClobPairIdResponse is a response message for the next clob pair id */

export interface QueryNextClobPairIdResponse {
  /** QueryNextClobPairIdResponse is a response message for the next clob pair id */
  nextClobPairId: number;
}
/** QueryNextClobPairIdResponse is a response message for the next clob pair id */

export interface QueryNextClobPairIdResponseSDKType {
  /** QueryNextClobPairIdResponse is a response message for the next clob pair id */
  next_clob_pair_id: number;
}
/** QueryLeverageRequest is a request message for Leverage. */

export interface QueryLeverageRequest {
  /** The address of the wallet that owns the subaccount. */
  owner: string;
  /** The unique number of the subaccount for the owner. */

  number: number;
}
/** QueryLeverageRequest is a request message for Leverage. */

export interface QueryLeverageRequestSDKType {
  /** The address of the wallet that owns the subaccount. */
  owner: string;
  /** The unique number of the subaccount for the owner. */

  number: number;
}
/** QueryLeverageResponse is a response message that contains the leverage map. */

export interface QueryLeverageResponse {
  /** List of clob pair leverage settings. */
  clobPairLeverage: ClobPairLeverageInfo[];
}
/** QueryLeverageResponse is a response message that contains the leverage map. */

export interface QueryLeverageResponseSDKType {
  /** List of clob pair leverage settings. */
  clob_pair_leverage: ClobPairLeverageInfoSDKType[];
}
/** ClobPairLeverageInfo represents the leverage setting for a single clob pair. */

export interface ClobPairLeverageInfo {
  /** The clob pair ID. */
  clobPairId: number;
  /** The user selected imf. */

  customImfPpm: number;
}
/** ClobPairLeverageInfo represents the leverage setting for a single clob pair. */

export interface ClobPairLeverageInfoSDKType {
  /** The clob pair ID. */
  clob_pair_id: number;
  /** The user selected imf. */

  custom_imf_ppm: number;
}
/**
 * StreamOrderbookUpdatesRequest is a request message for the
 * StreamOrderbookUpdates method.
 */

export interface StreamOrderbookUpdatesRequest {
  /** Clob pair ids to stream orderbook updates for. */
  clobPairId: number[];
  /** Subaccount ids to stream subaccount updates for. */

  subaccountIds: SubaccountId[];
  /** Market ids for price updates. */

  marketIds: number[];
  /**
   * Filter order updates by subaccount IDs.
   * If true, the orderbook updates only include orders from provided subaccount
   * IDs.
   */

  filterOrdersBySubaccountId: boolean;
}
/**
 * StreamOrderbookUpdatesRequest is a request message for the
 * StreamOrderbookUpdates method.
 */

export interface StreamOrderbookUpdatesRequestSDKType {
  /** Clob pair ids to stream orderbook updates for. */
  clob_pair_id: number[];
  /** Subaccount ids to stream subaccount updates for. */

  subaccount_ids: SubaccountIdSDKType[];
  /** Market ids for price updates. */

  market_ids: number[];
  /**
   * Filter order updates by subaccount IDs.
   * If true, the orderbook updates only include orders from provided subaccount
   * IDs.
   */

  filter_orders_by_subaccount_id: boolean;
}
/**
 * StreamOrderbookUpdatesResponse is a response message for the
 * StreamOrderbookUpdates method.
 */

export interface StreamOrderbookUpdatesResponse {
  /** Batch of updates for the clob pair. */
  updates: StreamUpdate[];
}
/**
 * StreamOrderbookUpdatesResponse is a response message for the
 * StreamOrderbookUpdates method.
 */

export interface StreamOrderbookUpdatesResponseSDKType {
  /** Batch of updates for the clob pair. */
  updates: StreamUpdateSDKType[];
}
/**
 * StreamUpdate is an update that will be pushed through the
 * GRPC stream.
 */

export interface StreamUpdate {
  /** Block height of the update. */
  blockHeight: number;
  /** Exec mode of the update. */

  execMode: number;
  orderbookUpdate?: StreamOrderbookUpdate;
  orderFill?: StreamOrderbookFill;
  takerOrder?: StreamTakerOrder;
  subaccountUpdate?: StreamSubaccountUpdate;
  priceUpdate?: StreamPriceUpdate;
}
/**
 * StreamUpdate is an update that will be pushed through the
 * GRPC stream.
 */

export interface StreamUpdateSDKType {
  /** Block height of the update. */
  block_height: number;
  /** Exec mode of the update. */

  exec_mode: number;
  orderbook_update?: StreamOrderbookUpdateSDKType;
  order_fill?: StreamOrderbookFillSDKType;
  taker_order?: StreamTakerOrderSDKType;
  subaccount_update?: StreamSubaccountUpdateSDKType;
  price_update?: StreamPriceUpdateSDKType;
}
/**
 * StreamOrderbookUpdate provides information on an orderbook update. Used in
 * the full node GRPC stream.
 */

export interface StreamOrderbookUpdate {
  /**
   * Snapshot indicates if the response is from a snapshot of the orderbook.
   * All updates should be ignored until snapshot is recieved.
   * If the snapshot is true, then all previous entries should be
   * discarded and the orderbook should be resynced.
   */
  snapshot: boolean;
  /**
   * Orderbook updates for the clob pair. Can contain order place, removals,
   * or updates.
   */

  updates: OffChainUpdateV1[];
}
/**
 * StreamOrderbookUpdate provides information on an orderbook update. Used in
 * the full node GRPC stream.
 */

export interface StreamOrderbookUpdateSDKType {
  /**
   * Snapshot indicates if the response is from a snapshot of the orderbook.
   * All updates should be ignored until snapshot is recieved.
   * If the snapshot is true, then all previous entries should be
   * discarded and the orderbook should be resynced.
   */
  snapshot: boolean;
  /**
   * Orderbook updates for the clob pair. Can contain order place, removals,
   * or updates.
   */

  updates: OffChainUpdateV1SDKType[];
}
/**
 * StreamOrderbookFill provides information on an orderbook fill. Used in
 * the full node GRPC stream.
 */

export interface StreamOrderbookFill {
  /**
   * Clob match. Provides information on which orders were matched
   * and the type of order.
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
   * and the type of order.
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
/**
 * StreamTakerOrder provides information on a taker order that was attempted
 * to be matched on the orderbook.
 * It is intended to be used only in full node streaming.
 */

export interface StreamTakerOrder {
  order?: Order;
  liquidationOrder?: StreamLiquidationOrder;
  /**
   * Information on the taker order after it is matched on the book,
   * either successfully or unsuccessfully.
   */

  takerOrderStatus?: StreamTakerOrderStatus;
}
/**
 * StreamTakerOrder provides information on a taker order that was attempted
 * to be matched on the orderbook.
 * It is intended to be used only in full node streaming.
 */

export interface StreamTakerOrderSDKType {
  order?: OrderSDKType;
  liquidation_order?: StreamLiquidationOrderSDKType;
  /**
   * Information on the taker order after it is matched on the book,
   * either successfully or unsuccessfully.
   */

  taker_order_status?: StreamTakerOrderStatusSDKType;
}
/**
 * StreamTakerOrderStatus is a representation of a taker order
 * after it is attempted to be matched on the orderbook.
 * It is intended to be used only in full node streaming.
 */

export interface StreamTakerOrderStatus {
  /**
   * The state of the taker order after attempting to match it against the
   * orderbook. Possible enum values can be found here:
   * https://github.com/dydxprotocol/v4-chain/blob/main/protocol/x/clob/types/orderbook.go#L105
   */
  orderStatus: number;
  /** The amount of remaining (non-matched) base quantums of this taker order. */

  remainingQuantums: Long;
  /**
   * The amount of base quantums that were *optimistically* filled for this
   * taker order when the order is matched against the orderbook. Note that if
   * any quantums of this order were optimistically filled or filled in state
   * before this invocation of the matching loop, this value will not include
   * them.
   */

  optimisticallyFilledQuantums: Long;
}
/**
 * StreamTakerOrderStatus is a representation of a taker order
 * after it is attempted to be matched on the orderbook.
 * It is intended to be used only in full node streaming.
 */

export interface StreamTakerOrderStatusSDKType {
  /**
   * The state of the taker order after attempting to match it against the
   * orderbook. Possible enum values can be found here:
   * https://github.com/dydxprotocol/v4-chain/blob/main/protocol/x/clob/types/orderbook.go#L105
   */
  order_status: number;
  /** The amount of remaining (non-matched) base quantums of this taker order. */

  remaining_quantums: Long;
  /**
   * The amount of base quantums that were *optimistically* filled for this
   * taker order when the order is matched against the orderbook. Note that if
   * any quantums of this order were optimistically filled or filled in state
   * before this invocation of the matching loop, this value will not include
   * them.
   */

  optimistically_filled_quantums: Long;
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

function createBaseQueryStatefulOrderRequest(): QueryStatefulOrderRequest {
  return {
    orderId: undefined
  };
}

export const QueryStatefulOrderRequest = {
  encode(message: QueryStatefulOrderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orderId !== undefined) {
      OrderId.encode(message.orderId, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryStatefulOrderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryStatefulOrderRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.orderId = OrderId.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryStatefulOrderRequest>): QueryStatefulOrderRequest {
    const message = createBaseQueryStatefulOrderRequest();
    message.orderId = object.orderId !== undefined && object.orderId !== null ? OrderId.fromPartial(object.orderId) : undefined;
    return message;
  }

};

function createBaseQueryStatefulOrderResponse(): QueryStatefulOrderResponse {
  return {
    orderPlacement: undefined,
    fillAmount: Long.UZERO,
    triggered: false
  };
}

export const QueryStatefulOrderResponse = {
  encode(message: QueryStatefulOrderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orderPlacement !== undefined) {
      LongTermOrderPlacement.encode(message.orderPlacement, writer.uint32(10).fork()).ldelim();
    }

    if (!message.fillAmount.isZero()) {
      writer.uint32(16).uint64(message.fillAmount);
    }

    if (message.triggered === true) {
      writer.uint32(24).bool(message.triggered);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryStatefulOrderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryStatefulOrderResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.orderPlacement = LongTermOrderPlacement.decode(reader, reader.uint32());
          break;

        case 2:
          message.fillAmount = (reader.uint64() as Long);
          break;

        case 3:
          message.triggered = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryStatefulOrderResponse>): QueryStatefulOrderResponse {
    const message = createBaseQueryStatefulOrderResponse();
    message.orderPlacement = object.orderPlacement !== undefined && object.orderPlacement !== null ? LongTermOrderPlacement.fromPartial(object.orderPlacement) : undefined;
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? Long.fromValue(object.fillAmount) : Long.UZERO;
    message.triggered = object.triggered ?? false;
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

function createBaseQueryNextClobPairIdRequest(): QueryNextClobPairIdRequest {
  return {};
}

export const QueryNextClobPairIdRequest = {
  encode(_: QueryNextClobPairIdRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNextClobPairIdRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNextClobPairIdRequest();

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

  fromPartial(_: DeepPartial<QueryNextClobPairIdRequest>): QueryNextClobPairIdRequest {
    const message = createBaseQueryNextClobPairIdRequest();
    return message;
  }

};

function createBaseQueryNextClobPairIdResponse(): QueryNextClobPairIdResponse {
  return {
    nextClobPairId: 0
  };
}

export const QueryNextClobPairIdResponse = {
  encode(message: QueryNextClobPairIdResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nextClobPairId !== 0) {
      writer.uint32(8).uint32(message.nextClobPairId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNextClobPairIdResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNextClobPairIdResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.nextClobPairId = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryNextClobPairIdResponse>): QueryNextClobPairIdResponse {
    const message = createBaseQueryNextClobPairIdResponse();
    message.nextClobPairId = object.nextClobPairId ?? 0;
    return message;
  }

};

function createBaseQueryLeverageRequest(): QueryLeverageRequest {
  return {
    owner: "",
    number: 0
  };
}

export const QueryLeverageRequest = {
  encode(message: QueryLeverageRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.owner !== "") {
      writer.uint32(10).string(message.owner);
    }

    if (message.number !== 0) {
      writer.uint32(16).uint32(message.number);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryLeverageRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryLeverageRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.owner = reader.string();
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

  fromPartial(object: DeepPartial<QueryLeverageRequest>): QueryLeverageRequest {
    const message = createBaseQueryLeverageRequest();
    message.owner = object.owner ?? "";
    message.number = object.number ?? 0;
    return message;
  }

};

function createBaseQueryLeverageResponse(): QueryLeverageResponse {
  return {
    clobPairLeverage: []
  };
}

export const QueryLeverageResponse = {
  encode(message: QueryLeverageResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.clobPairLeverage) {
      ClobPairLeverageInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryLeverageResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryLeverageResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPairLeverage.push(ClobPairLeverageInfo.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryLeverageResponse>): QueryLeverageResponse {
    const message = createBaseQueryLeverageResponse();
    message.clobPairLeverage = object.clobPairLeverage?.map(e => ClobPairLeverageInfo.fromPartial(e)) || [];
    return message;
  }

};

function createBaseClobPairLeverageInfo(): ClobPairLeverageInfo {
  return {
    clobPairId: 0,
    customImfPpm: 0
  };
}

export const ClobPairLeverageInfo = {
  encode(message: ClobPairLeverageInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobPairId !== 0) {
      writer.uint32(8).uint32(message.clobPairId);
    }

    if (message.customImfPpm !== 0) {
      writer.uint32(16).uint32(message.customImfPpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ClobPairLeverageInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClobPairLeverageInfo();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPairId = reader.uint32();
          break;

        case 2:
          message.customImfPpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ClobPairLeverageInfo>): ClobPairLeverageInfo {
    const message = createBaseClobPairLeverageInfo();
    message.clobPairId = object.clobPairId ?? 0;
    message.customImfPpm = object.customImfPpm ?? 0;
    return message;
  }

};

function createBaseStreamOrderbookUpdatesRequest(): StreamOrderbookUpdatesRequest {
  return {
    clobPairId: [],
    subaccountIds: [],
    marketIds: [],
    filterOrdersBySubaccountId: false
  };
}

export const StreamOrderbookUpdatesRequest = {
  encode(message: StreamOrderbookUpdatesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    writer.uint32(10).fork();

    for (const v of message.clobPairId) {
      writer.uint32(v);
    }

    writer.ldelim();

    for (const v of message.subaccountIds) {
      SubaccountId.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    writer.uint32(26).fork();

    for (const v of message.marketIds) {
      writer.uint32(v);
    }

    writer.ldelim();

    if (message.filterOrdersBySubaccountId === true) {
      writer.uint32(32).bool(message.filterOrdersBySubaccountId);
    }

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

        case 2:
          message.subaccountIds.push(SubaccountId.decode(reader, reader.uint32()));
          break;

        case 3:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;

            while (reader.pos < end2) {
              message.marketIds.push(reader.uint32());
            }
          } else {
            message.marketIds.push(reader.uint32());
          }

          break;

        case 4:
          message.filterOrdersBySubaccountId = reader.bool();
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
    message.subaccountIds = object.subaccountIds?.map(e => SubaccountId.fromPartial(e)) || [];
    message.marketIds = object.marketIds?.map(e => e) || [];
    message.filterOrdersBySubaccountId = object.filterOrdersBySubaccountId ?? false;
    return message;
  }

};

function createBaseStreamOrderbookUpdatesResponse(): StreamOrderbookUpdatesResponse {
  return {
    updates: []
  };
}

export const StreamOrderbookUpdatesResponse = {
  encode(message: StreamOrderbookUpdatesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.updates) {
      StreamUpdate.encode(v!, writer.uint32(10).fork()).ldelim();
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
    return message;
  }

};

function createBaseStreamUpdate(): StreamUpdate {
  return {
    blockHeight: 0,
    execMode: 0,
    orderbookUpdate: undefined,
    orderFill: undefined,
    takerOrder: undefined,
    subaccountUpdate: undefined,
    priceUpdate: undefined
  };
}

export const StreamUpdate = {
  encode(message: StreamUpdate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockHeight !== 0) {
      writer.uint32(8).uint32(message.blockHeight);
    }

    if (message.execMode !== 0) {
      writer.uint32(16).uint32(message.execMode);
    }

    if (message.orderbookUpdate !== undefined) {
      StreamOrderbookUpdate.encode(message.orderbookUpdate, writer.uint32(26).fork()).ldelim();
    }

    if (message.orderFill !== undefined) {
      StreamOrderbookFill.encode(message.orderFill, writer.uint32(34).fork()).ldelim();
    }

    if (message.takerOrder !== undefined) {
      StreamTakerOrder.encode(message.takerOrder, writer.uint32(42).fork()).ldelim();
    }

    if (message.subaccountUpdate !== undefined) {
      StreamSubaccountUpdate.encode(message.subaccountUpdate, writer.uint32(50).fork()).ldelim();
    }

    if (message.priceUpdate !== undefined) {
      StreamPriceUpdate.encode(message.priceUpdate, writer.uint32(58).fork()).ldelim();
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
          message.blockHeight = reader.uint32();
          break;

        case 2:
          message.execMode = reader.uint32();
          break;

        case 3:
          message.orderbookUpdate = StreamOrderbookUpdate.decode(reader, reader.uint32());
          break;

        case 4:
          message.orderFill = StreamOrderbookFill.decode(reader, reader.uint32());
          break;

        case 5:
          message.takerOrder = StreamTakerOrder.decode(reader, reader.uint32());
          break;

        case 6:
          message.subaccountUpdate = StreamSubaccountUpdate.decode(reader, reader.uint32());
          break;

        case 7:
          message.priceUpdate = StreamPriceUpdate.decode(reader, reader.uint32());
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
    message.blockHeight = object.blockHeight ?? 0;
    message.execMode = object.execMode ?? 0;
    message.orderbookUpdate = object.orderbookUpdate !== undefined && object.orderbookUpdate !== null ? StreamOrderbookUpdate.fromPartial(object.orderbookUpdate) : undefined;
    message.orderFill = object.orderFill !== undefined && object.orderFill !== null ? StreamOrderbookFill.fromPartial(object.orderFill) : undefined;
    message.takerOrder = object.takerOrder !== undefined && object.takerOrder !== null ? StreamTakerOrder.fromPartial(object.takerOrder) : undefined;
    message.subaccountUpdate = object.subaccountUpdate !== undefined && object.subaccountUpdate !== null ? StreamSubaccountUpdate.fromPartial(object.subaccountUpdate) : undefined;
    message.priceUpdate = object.priceUpdate !== undefined && object.priceUpdate !== null ? StreamPriceUpdate.fromPartial(object.priceUpdate) : undefined;
    return message;
  }

};

function createBaseStreamOrderbookUpdate(): StreamOrderbookUpdate {
  return {
    snapshot: false,
    updates: []
  };
}

export const StreamOrderbookUpdate = {
  encode(message: StreamOrderbookUpdate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.snapshot === true) {
      writer.uint32(8).bool(message.snapshot);
    }

    for (const v of message.updates) {
      OffChainUpdateV1.encode(v!, writer.uint32(18).fork()).ldelim();
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
          message.snapshot = reader.bool();
          break;

        case 2:
          message.updates.push(OffChainUpdateV1.decode(reader, reader.uint32()));
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
    message.snapshot = object.snapshot ?? false;
    message.updates = object.updates?.map(e => OffChainUpdateV1.fromPartial(e)) || [];
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

function createBaseStreamTakerOrder(): StreamTakerOrder {
  return {
    order: undefined,
    liquidationOrder: undefined,
    takerOrderStatus: undefined
  };
}

export const StreamTakerOrder = {
  encode(message: StreamTakerOrder, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.order !== undefined) {
      Order.encode(message.order, writer.uint32(10).fork()).ldelim();
    }

    if (message.liquidationOrder !== undefined) {
      StreamLiquidationOrder.encode(message.liquidationOrder, writer.uint32(18).fork()).ldelim();
    }

    if (message.takerOrderStatus !== undefined) {
      StreamTakerOrderStatus.encode(message.takerOrderStatus, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamTakerOrder {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamTakerOrder();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.order = Order.decode(reader, reader.uint32());
          break;

        case 2:
          message.liquidationOrder = StreamLiquidationOrder.decode(reader, reader.uint32());
          break;

        case 3:
          message.takerOrderStatus = StreamTakerOrderStatus.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StreamTakerOrder>): StreamTakerOrder {
    const message = createBaseStreamTakerOrder();
    message.order = object.order !== undefined && object.order !== null ? Order.fromPartial(object.order) : undefined;
    message.liquidationOrder = object.liquidationOrder !== undefined && object.liquidationOrder !== null ? StreamLiquidationOrder.fromPartial(object.liquidationOrder) : undefined;
    message.takerOrderStatus = object.takerOrderStatus !== undefined && object.takerOrderStatus !== null ? StreamTakerOrderStatus.fromPartial(object.takerOrderStatus) : undefined;
    return message;
  }

};

function createBaseStreamTakerOrderStatus(): StreamTakerOrderStatus {
  return {
    orderStatus: 0,
    remainingQuantums: Long.UZERO,
    optimisticallyFilledQuantums: Long.UZERO
  };
}

export const StreamTakerOrderStatus = {
  encode(message: StreamTakerOrderStatus, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orderStatus !== 0) {
      writer.uint32(8).uint32(message.orderStatus);
    }

    if (!message.remainingQuantums.isZero()) {
      writer.uint32(16).uint64(message.remainingQuantums);
    }

    if (!message.optimisticallyFilledQuantums.isZero()) {
      writer.uint32(24).uint64(message.optimisticallyFilledQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamTakerOrderStatus {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamTakerOrderStatus();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.orderStatus = reader.uint32();
          break;

        case 2:
          message.remainingQuantums = (reader.uint64() as Long);
          break;

        case 3:
          message.optimisticallyFilledQuantums = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StreamTakerOrderStatus>): StreamTakerOrderStatus {
    const message = createBaseStreamTakerOrderStatus();
    message.orderStatus = object.orderStatus ?? 0;
    message.remainingQuantums = object.remainingQuantums !== undefined && object.remainingQuantums !== null ? Long.fromValue(object.remainingQuantums) : Long.UZERO;
    message.optimisticallyFilledQuantums = object.optimisticallyFilledQuantums !== undefined && object.optimisticallyFilledQuantums !== null ? Long.fromValue(object.optimisticallyFilledQuantums) : Long.UZERO;
    return message;
  }

};