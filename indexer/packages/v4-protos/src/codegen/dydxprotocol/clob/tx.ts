import { Order, OrderAmino, OrderSDKType, OrderId, OrderIdAmino, OrderIdSDKType } from "./order";
import { ClobPair, ClobPairAmino, ClobPairSDKType } from "./clob_pair";
import { EquityTierLimitConfiguration, EquityTierLimitConfigurationAmino, EquityTierLimitConfigurationSDKType } from "./equity_tier_limit_config";
import { BlockRateLimitConfiguration, BlockRateLimitConfigurationAmino, BlockRateLimitConfigurationSDKType } from "./block_rate_limit_config";
import { LiquidationsConfig, LiquidationsConfigAmino, LiquidationsConfigSDKType } from "./liquidations_config";
import { ClobMatch, ClobMatchAmino, ClobMatchSDKType } from "./matches";
import { OrderRemoval, OrderRemovalAmino, OrderRemovalSDKType } from "./order_removals";
import { BinaryReader, BinaryWriter } from "../../binary";
import { bytesFromBase64, base64FromBytes } from "../../helpers";
/** MsgCreateClobPair is a message used by x/gov for creating a new clob pair. */
export interface MsgCreateClobPair {
  /** The address that controls the module. */
  authority: string;
  /** `clob_pair` defines parameters for the new clob pair. */
  clobPair: ClobPair;
}
export interface MsgCreateClobPairProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgCreateClobPair";
  value: Uint8Array;
}
/** MsgCreateClobPair is a message used by x/gov for creating a new clob pair. */
export interface MsgCreateClobPairAmino {
  /** The address that controls the module. */
  authority?: string;
  /** `clob_pair` defines parameters for the new clob pair. */
  clob_pair?: ClobPairAmino;
}
export interface MsgCreateClobPairAminoMsg {
  type: "/dydxprotocol.clob.MsgCreateClobPair";
  value: MsgCreateClobPairAmino;
}
/** MsgCreateClobPair is a message used by x/gov for creating a new clob pair. */
export interface MsgCreateClobPairSDKType {
  authority: string;
  clob_pair: ClobPairSDKType;
}
/** MsgCreateClobPairResponse defines the CreateClobPair response type. */
export interface MsgCreateClobPairResponse {}
export interface MsgCreateClobPairResponseProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgCreateClobPairResponse";
  value: Uint8Array;
}
/** MsgCreateClobPairResponse defines the CreateClobPair response type. */
export interface MsgCreateClobPairResponseAmino {}
export interface MsgCreateClobPairResponseAminoMsg {
  type: "/dydxprotocol.clob.MsgCreateClobPairResponse";
  value: MsgCreateClobPairResponseAmino;
}
/** MsgCreateClobPairResponse defines the CreateClobPair response type. */
export interface MsgCreateClobPairResponseSDKType {}
/**
 * MsgProposedOperations is a message injected by block proposers to
 * specify the operations that occurred in a block.
 */
export interface MsgProposedOperations {
  /** The list of operations proposed by the block proposer. */
  operationsQueue: OperationRaw[];
}
export interface MsgProposedOperationsProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgProposedOperations";
  value: Uint8Array;
}
/**
 * MsgProposedOperations is a message injected by block proposers to
 * specify the operations that occurred in a block.
 */
export interface MsgProposedOperationsAmino {
  /** The list of operations proposed by the block proposer. */
  operations_queue?: OperationRawAmino[];
}
export interface MsgProposedOperationsAminoMsg {
  type: "/dydxprotocol.clob.MsgProposedOperations";
  value: MsgProposedOperationsAmino;
}
/**
 * MsgProposedOperations is a message injected by block proposers to
 * specify the operations that occurred in a block.
 */
export interface MsgProposedOperationsSDKType {
  operations_queue: OperationRawSDKType[];
}
/**
 * MsgProposedOperationsResponse is the response type of the message injected
 * by block proposers to specify the operations that occurred in a block.
 */
export interface MsgProposedOperationsResponse {}
export interface MsgProposedOperationsResponseProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgProposedOperationsResponse";
  value: Uint8Array;
}
/**
 * MsgProposedOperationsResponse is the response type of the message injected
 * by block proposers to specify the operations that occurred in a block.
 */
export interface MsgProposedOperationsResponseAmino {}
export interface MsgProposedOperationsResponseAminoMsg {
  type: "/dydxprotocol.clob.MsgProposedOperationsResponse";
  value: MsgProposedOperationsResponseAmino;
}
/**
 * MsgProposedOperationsResponse is the response type of the message injected
 * by block proposers to specify the operations that occurred in a block.
 */
export interface MsgProposedOperationsResponseSDKType {}
/** MsgPlaceOrder is a request type used for placing orders. */
export interface MsgPlaceOrder {
  /** MsgPlaceOrder is a request type used for placing orders. */
  order: Order;
}
export interface MsgPlaceOrderProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgPlaceOrder";
  value: Uint8Array;
}
/** MsgPlaceOrder is a request type used for placing orders. */
export interface MsgPlaceOrderAmino {
  /** MsgPlaceOrder is a request type used for placing orders. */
  order?: OrderAmino;
}
export interface MsgPlaceOrderAminoMsg {
  type: "/dydxprotocol.clob.MsgPlaceOrder";
  value: MsgPlaceOrderAmino;
}
/** MsgPlaceOrder is a request type used for placing orders. */
export interface MsgPlaceOrderSDKType {
  order: OrderSDKType;
}
/** MsgPlaceOrderResponse is a response type used for placing orders. */
export interface MsgPlaceOrderResponse {}
export interface MsgPlaceOrderResponseProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgPlaceOrderResponse";
  value: Uint8Array;
}
/** MsgPlaceOrderResponse is a response type used for placing orders. */
export interface MsgPlaceOrderResponseAmino {}
export interface MsgPlaceOrderResponseAminoMsg {
  type: "/dydxprotocol.clob.MsgPlaceOrderResponse";
  value: MsgPlaceOrderResponseAmino;
}
/** MsgPlaceOrderResponse is a response type used for placing orders. */
export interface MsgPlaceOrderResponseSDKType {}
/** MsgCancelOrder is a request type used for canceling orders. */
export interface MsgCancelOrder {
  orderId: OrderId;
  /**
   * The last block this order cancellation can be executed at.
   * Used only for Short-Term orders and must be zero for stateful orders.
   */
  goodTilBlock?: number;
  /**
   * good_til_block_time represents the unix timestamp (in seconds) at which a
   * stateful order cancellation will be considered expired. The
   * good_til_block_time is always evaluated against the previous block's
   * `BlockTime` instead of the block in which the order is committed.
   * This value must be zero for Short-Term orders.
   */
  goodTilBlockTime?: number;
}
export interface MsgCancelOrderProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgCancelOrder";
  value: Uint8Array;
}
/** MsgCancelOrder is a request type used for canceling orders. */
export interface MsgCancelOrderAmino {
  order_id?: OrderIdAmino;
  /**
   * The last block this order cancellation can be executed at.
   * Used only for Short-Term orders and must be zero for stateful orders.
   */
  good_til_block?: number;
  /**
   * good_til_block_time represents the unix timestamp (in seconds) at which a
   * stateful order cancellation will be considered expired. The
   * good_til_block_time is always evaluated against the previous block's
   * `BlockTime` instead of the block in which the order is committed.
   * This value must be zero for Short-Term orders.
   */
  good_til_block_time?: number;
}
export interface MsgCancelOrderAminoMsg {
  type: "/dydxprotocol.clob.MsgCancelOrder";
  value: MsgCancelOrderAmino;
}
/** MsgCancelOrder is a request type used for canceling orders. */
export interface MsgCancelOrderSDKType {
  order_id: OrderIdSDKType;
  good_til_block?: number;
  good_til_block_time?: number;
}
/** MsgCancelOrderResponse is a response type used for canceling orders. */
export interface MsgCancelOrderResponse {}
export interface MsgCancelOrderResponseProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgCancelOrderResponse";
  value: Uint8Array;
}
/** MsgCancelOrderResponse is a response type used for canceling orders. */
export interface MsgCancelOrderResponseAmino {}
export interface MsgCancelOrderResponseAminoMsg {
  type: "/dydxprotocol.clob.MsgCancelOrderResponse";
  value: MsgCancelOrderResponseAmino;
}
/** MsgCancelOrderResponse is a response type used for canceling orders. */
export interface MsgCancelOrderResponseSDKType {}
/** MsgUpdateClobPair is a request type used for updating a ClobPair in state. */
export interface MsgUpdateClobPair {
  /** Authority is the address that may send this message. */
  authority: string;
  /** `clob_pair` is the ClobPair to write to state. */
  clobPair: ClobPair;
}
export interface MsgUpdateClobPairProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgUpdateClobPair";
  value: Uint8Array;
}
/** MsgUpdateClobPair is a request type used for updating a ClobPair in state. */
export interface MsgUpdateClobPairAmino {
  /** Authority is the address that may send this message. */
  authority?: string;
  /** `clob_pair` is the ClobPair to write to state. */
  clob_pair?: ClobPairAmino;
}
export interface MsgUpdateClobPairAminoMsg {
  type: "/dydxprotocol.clob.MsgUpdateClobPair";
  value: MsgUpdateClobPairAmino;
}
/** MsgUpdateClobPair is a request type used for updating a ClobPair in state. */
export interface MsgUpdateClobPairSDKType {
  authority: string;
  clob_pair: ClobPairSDKType;
}
/**
 * MsgUpdateClobPairResponse is a response type used for setting a ClobPair's
 * status.
 */
export interface MsgUpdateClobPairResponse {}
export interface MsgUpdateClobPairResponseProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgUpdateClobPairResponse";
  value: Uint8Array;
}
/**
 * MsgUpdateClobPairResponse is a response type used for setting a ClobPair's
 * status.
 */
export interface MsgUpdateClobPairResponseAmino {}
export interface MsgUpdateClobPairResponseAminoMsg {
  type: "/dydxprotocol.clob.MsgUpdateClobPairResponse";
  value: MsgUpdateClobPairResponseAmino;
}
/**
 * MsgUpdateClobPairResponse is a response type used for setting a ClobPair's
 * status.
 */
export interface MsgUpdateClobPairResponseSDKType {}
/**
 * OperationRaw represents an operation in the proposed operations.
 * Note that the `order_placement` operation is a signed message.
 */
export interface OperationRaw {
  match?: ClobMatch;
  shortTermOrderPlacement?: Uint8Array;
  orderRemoval?: OrderRemoval;
}
export interface OperationRawProtoMsg {
  typeUrl: "/dydxprotocol.clob.OperationRaw";
  value: Uint8Array;
}
/**
 * OperationRaw represents an operation in the proposed operations.
 * Note that the `order_placement` operation is a signed message.
 */
export interface OperationRawAmino {
  match?: ClobMatchAmino;
  short_term_order_placement?: string;
  order_removal?: OrderRemovalAmino;
}
export interface OperationRawAminoMsg {
  type: "/dydxprotocol.clob.OperationRaw";
  value: OperationRawAmino;
}
/**
 * OperationRaw represents an operation in the proposed operations.
 * Note that the `order_placement` operation is a signed message.
 */
export interface OperationRawSDKType {
  match?: ClobMatchSDKType;
  short_term_order_placement?: Uint8Array;
  order_removal?: OrderRemovalSDKType;
}
/**
 * MsgUpdateEquityTierLimitConfiguration is the Msg/EquityTierLimitConfiguration
 * request type.
 */
export interface MsgUpdateEquityTierLimitConfiguration {
  authority: string;
  /**
   * Defines the equity tier limit configuration to update to. All fields must
   * be set.
   */
  equityTierLimitConfig: EquityTierLimitConfiguration;
}
export interface MsgUpdateEquityTierLimitConfigurationProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration";
  value: Uint8Array;
}
/**
 * MsgUpdateEquityTierLimitConfiguration is the Msg/EquityTierLimitConfiguration
 * request type.
 */
export interface MsgUpdateEquityTierLimitConfigurationAmino {
  authority?: string;
  /**
   * Defines the equity tier limit configuration to update to. All fields must
   * be set.
   */
  equity_tier_limit_config?: EquityTierLimitConfigurationAmino;
}
export interface MsgUpdateEquityTierLimitConfigurationAminoMsg {
  type: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration";
  value: MsgUpdateEquityTierLimitConfigurationAmino;
}
/**
 * MsgUpdateEquityTierLimitConfiguration is the Msg/EquityTierLimitConfiguration
 * request type.
 */
export interface MsgUpdateEquityTierLimitConfigurationSDKType {
  authority: string;
  equity_tier_limit_config: EquityTierLimitConfigurationSDKType;
}
/**
 * MsgUpdateEquityTierLimitConfiguration is the Msg/EquityTierLimitConfiguration
 * response type.
 */
export interface MsgUpdateEquityTierLimitConfigurationResponse {}
export interface MsgUpdateEquityTierLimitConfigurationResponseProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfigurationResponse";
  value: Uint8Array;
}
/**
 * MsgUpdateEquityTierLimitConfiguration is the Msg/EquityTierLimitConfiguration
 * response type.
 */
export interface MsgUpdateEquityTierLimitConfigurationResponseAmino {}
export interface MsgUpdateEquityTierLimitConfigurationResponseAminoMsg {
  type: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfigurationResponse";
  value: MsgUpdateEquityTierLimitConfigurationResponseAmino;
}
/**
 * MsgUpdateEquityTierLimitConfiguration is the Msg/EquityTierLimitConfiguration
 * response type.
 */
export interface MsgUpdateEquityTierLimitConfigurationResponseSDKType {}
/**
 * MsgUpdateBlockRateLimitConfiguration is the Msg/BlockRateLimitConfiguration
 * request type.
 */
export interface MsgUpdateBlockRateLimitConfiguration {
  authority: string;
  /**
   * Defines the block rate limit configuration to update to. All fields must be
   * set.
   */
  blockRateLimitConfig: BlockRateLimitConfiguration;
}
export interface MsgUpdateBlockRateLimitConfigurationProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration";
  value: Uint8Array;
}
/**
 * MsgUpdateBlockRateLimitConfiguration is the Msg/BlockRateLimitConfiguration
 * request type.
 */
export interface MsgUpdateBlockRateLimitConfigurationAmino {
  authority?: string;
  /**
   * Defines the block rate limit configuration to update to. All fields must be
   * set.
   */
  block_rate_limit_config?: BlockRateLimitConfigurationAmino;
}
export interface MsgUpdateBlockRateLimitConfigurationAminoMsg {
  type: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration";
  value: MsgUpdateBlockRateLimitConfigurationAmino;
}
/**
 * MsgUpdateBlockRateLimitConfiguration is the Msg/BlockRateLimitConfiguration
 * request type.
 */
export interface MsgUpdateBlockRateLimitConfigurationSDKType {
  authority: string;
  block_rate_limit_config: BlockRateLimitConfigurationSDKType;
}
/**
 * MsgUpdateBlockRateLimitConfiguration is a response type for updating the
 * liquidations config.
 */
export interface MsgUpdateBlockRateLimitConfigurationResponse {}
export interface MsgUpdateBlockRateLimitConfigurationResponseProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfigurationResponse";
  value: Uint8Array;
}
/**
 * MsgUpdateBlockRateLimitConfiguration is a response type for updating the
 * liquidations config.
 */
export interface MsgUpdateBlockRateLimitConfigurationResponseAmino {}
export interface MsgUpdateBlockRateLimitConfigurationResponseAminoMsg {
  type: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfigurationResponse";
  value: MsgUpdateBlockRateLimitConfigurationResponseAmino;
}
/**
 * MsgUpdateBlockRateLimitConfiguration is a response type for updating the
 * liquidations config.
 */
export interface MsgUpdateBlockRateLimitConfigurationResponseSDKType {}
/**
 * MsgUpdateLiquidationsConfig is a request type for updating the liquidations
 * config.
 */
export interface MsgUpdateLiquidationsConfig {
  /** Authority is the address that may send this message. */
  authority: string;
  /**
   * Defines the liquidations configuration to update to. All fields must
   * be set.
   */
  liquidationsConfig: LiquidationsConfig;
}
export interface MsgUpdateLiquidationsConfigProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgUpdateLiquidationsConfig";
  value: Uint8Array;
}
/**
 * MsgUpdateLiquidationsConfig is a request type for updating the liquidations
 * config.
 */
export interface MsgUpdateLiquidationsConfigAmino {
  /** Authority is the address that may send this message. */
  authority?: string;
  /**
   * Defines the liquidations configuration to update to. All fields must
   * be set.
   */
  liquidations_config?: LiquidationsConfigAmino;
}
export interface MsgUpdateLiquidationsConfigAminoMsg {
  type: "/dydxprotocol.clob.MsgUpdateLiquidationsConfig";
  value: MsgUpdateLiquidationsConfigAmino;
}
/**
 * MsgUpdateLiquidationsConfig is a request type for updating the liquidations
 * config.
 */
export interface MsgUpdateLiquidationsConfigSDKType {
  authority: string;
  liquidations_config: LiquidationsConfigSDKType;
}
/** MsgUpdateLiquidationsConfig is the Msg/LiquidationsConfig response type. */
export interface MsgUpdateLiquidationsConfigResponse {}
export interface MsgUpdateLiquidationsConfigResponseProtoMsg {
  typeUrl: "/dydxprotocol.clob.MsgUpdateLiquidationsConfigResponse";
  value: Uint8Array;
}
/** MsgUpdateLiquidationsConfig is the Msg/LiquidationsConfig response type. */
export interface MsgUpdateLiquidationsConfigResponseAmino {}
export interface MsgUpdateLiquidationsConfigResponseAminoMsg {
  type: "/dydxprotocol.clob.MsgUpdateLiquidationsConfigResponse";
  value: MsgUpdateLiquidationsConfigResponseAmino;
}
/** MsgUpdateLiquidationsConfig is the Msg/LiquidationsConfig response type. */
export interface MsgUpdateLiquidationsConfigResponseSDKType {}
function createBaseMsgCreateClobPair(): MsgCreateClobPair {
  return {
    authority: "",
    clobPair: ClobPair.fromPartial({})
  };
}
export const MsgCreateClobPair = {
  typeUrl: "/dydxprotocol.clob.MsgCreateClobPair",
  encode(message: MsgCreateClobPair, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.clobPair !== undefined) {
      ClobPair.encode(message.clobPair, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateClobPair {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateClobPair();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.clobPair = ClobPair.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgCreateClobPair>): MsgCreateClobPair {
    const message = createBaseMsgCreateClobPair();
    message.authority = object.authority ?? "";
    message.clobPair = object.clobPair !== undefined && object.clobPair !== null ? ClobPair.fromPartial(object.clobPair) : undefined;
    return message;
  },
  fromAmino(object: MsgCreateClobPairAmino): MsgCreateClobPair {
    const message = createBaseMsgCreateClobPair();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.clob_pair !== undefined && object.clob_pair !== null) {
      message.clobPair = ClobPair.fromAmino(object.clob_pair);
    }
    return message;
  },
  toAmino(message: MsgCreateClobPair): MsgCreateClobPairAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.clob_pair = message.clobPair ? ClobPair.toAmino(message.clobPair) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgCreateClobPairAminoMsg): MsgCreateClobPair {
    return MsgCreateClobPair.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCreateClobPairProtoMsg): MsgCreateClobPair {
    return MsgCreateClobPair.decode(message.value);
  },
  toProto(message: MsgCreateClobPair): Uint8Array {
    return MsgCreateClobPair.encode(message).finish();
  },
  toProtoMsg(message: MsgCreateClobPair): MsgCreateClobPairProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgCreateClobPair",
      value: MsgCreateClobPair.encode(message).finish()
    };
  }
};
function createBaseMsgCreateClobPairResponse(): MsgCreateClobPairResponse {
  return {};
}
export const MsgCreateClobPairResponse = {
  typeUrl: "/dydxprotocol.clob.MsgCreateClobPairResponse",
  encode(_: MsgCreateClobPairResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateClobPairResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateClobPairResponse();
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
  fromPartial(_: Partial<MsgCreateClobPairResponse>): MsgCreateClobPairResponse {
    const message = createBaseMsgCreateClobPairResponse();
    return message;
  },
  fromAmino(_: MsgCreateClobPairResponseAmino): MsgCreateClobPairResponse {
    const message = createBaseMsgCreateClobPairResponse();
    return message;
  },
  toAmino(_: MsgCreateClobPairResponse): MsgCreateClobPairResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgCreateClobPairResponseAminoMsg): MsgCreateClobPairResponse {
    return MsgCreateClobPairResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCreateClobPairResponseProtoMsg): MsgCreateClobPairResponse {
    return MsgCreateClobPairResponse.decode(message.value);
  },
  toProto(message: MsgCreateClobPairResponse): Uint8Array {
    return MsgCreateClobPairResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgCreateClobPairResponse): MsgCreateClobPairResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgCreateClobPairResponse",
      value: MsgCreateClobPairResponse.encode(message).finish()
    };
  }
};
function createBaseMsgProposedOperations(): MsgProposedOperations {
  return {
    operationsQueue: []
  };
}
export const MsgProposedOperations = {
  typeUrl: "/dydxprotocol.clob.MsgProposedOperations",
  encode(message: MsgProposedOperations, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.operationsQueue) {
      OperationRaw.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgProposedOperations {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgProposedOperations();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.operationsQueue.push(OperationRaw.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgProposedOperations>): MsgProposedOperations {
    const message = createBaseMsgProposedOperations();
    message.operationsQueue = object.operationsQueue?.map(e => OperationRaw.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: MsgProposedOperationsAmino): MsgProposedOperations {
    const message = createBaseMsgProposedOperations();
    message.operationsQueue = object.operations_queue?.map(e => OperationRaw.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: MsgProposedOperations): MsgProposedOperationsAmino {
    const obj: any = {};
    if (message.operationsQueue) {
      obj.operations_queue = message.operationsQueue.map(e => e ? OperationRaw.toAmino(e) : undefined);
    } else {
      obj.operations_queue = [];
    }
    return obj;
  },
  fromAminoMsg(object: MsgProposedOperationsAminoMsg): MsgProposedOperations {
    return MsgProposedOperations.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgProposedOperationsProtoMsg): MsgProposedOperations {
    return MsgProposedOperations.decode(message.value);
  },
  toProto(message: MsgProposedOperations): Uint8Array {
    return MsgProposedOperations.encode(message).finish();
  },
  toProtoMsg(message: MsgProposedOperations): MsgProposedOperationsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgProposedOperations",
      value: MsgProposedOperations.encode(message).finish()
    };
  }
};
function createBaseMsgProposedOperationsResponse(): MsgProposedOperationsResponse {
  return {};
}
export const MsgProposedOperationsResponse = {
  typeUrl: "/dydxprotocol.clob.MsgProposedOperationsResponse",
  encode(_: MsgProposedOperationsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgProposedOperationsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgProposedOperationsResponse();
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
  fromPartial(_: Partial<MsgProposedOperationsResponse>): MsgProposedOperationsResponse {
    const message = createBaseMsgProposedOperationsResponse();
    return message;
  },
  fromAmino(_: MsgProposedOperationsResponseAmino): MsgProposedOperationsResponse {
    const message = createBaseMsgProposedOperationsResponse();
    return message;
  },
  toAmino(_: MsgProposedOperationsResponse): MsgProposedOperationsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgProposedOperationsResponseAminoMsg): MsgProposedOperationsResponse {
    return MsgProposedOperationsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgProposedOperationsResponseProtoMsg): MsgProposedOperationsResponse {
    return MsgProposedOperationsResponse.decode(message.value);
  },
  toProto(message: MsgProposedOperationsResponse): Uint8Array {
    return MsgProposedOperationsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgProposedOperationsResponse): MsgProposedOperationsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgProposedOperationsResponse",
      value: MsgProposedOperationsResponse.encode(message).finish()
    };
  }
};
function createBaseMsgPlaceOrder(): MsgPlaceOrder {
  return {
    order: Order.fromPartial({})
  };
}
export const MsgPlaceOrder = {
  typeUrl: "/dydxprotocol.clob.MsgPlaceOrder",
  encode(message: MsgPlaceOrder, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.order !== undefined) {
      Order.encode(message.order, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgPlaceOrder {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgPlaceOrder();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.order = Order.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgPlaceOrder>): MsgPlaceOrder {
    const message = createBaseMsgPlaceOrder();
    message.order = object.order !== undefined && object.order !== null ? Order.fromPartial(object.order) : undefined;
    return message;
  },
  fromAmino(object: MsgPlaceOrderAmino): MsgPlaceOrder {
    const message = createBaseMsgPlaceOrder();
    if (object.order !== undefined && object.order !== null) {
      message.order = Order.fromAmino(object.order);
    }
    return message;
  },
  toAmino(message: MsgPlaceOrder): MsgPlaceOrderAmino {
    const obj: any = {};
    obj.order = message.order ? Order.toAmino(message.order) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgPlaceOrderAminoMsg): MsgPlaceOrder {
    return MsgPlaceOrder.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgPlaceOrderProtoMsg): MsgPlaceOrder {
    return MsgPlaceOrder.decode(message.value);
  },
  toProto(message: MsgPlaceOrder): Uint8Array {
    return MsgPlaceOrder.encode(message).finish();
  },
  toProtoMsg(message: MsgPlaceOrder): MsgPlaceOrderProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgPlaceOrder",
      value: MsgPlaceOrder.encode(message).finish()
    };
  }
};
function createBaseMsgPlaceOrderResponse(): MsgPlaceOrderResponse {
  return {};
}
export const MsgPlaceOrderResponse = {
  typeUrl: "/dydxprotocol.clob.MsgPlaceOrderResponse",
  encode(_: MsgPlaceOrderResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgPlaceOrderResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgPlaceOrderResponse();
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
  fromPartial(_: Partial<MsgPlaceOrderResponse>): MsgPlaceOrderResponse {
    const message = createBaseMsgPlaceOrderResponse();
    return message;
  },
  fromAmino(_: MsgPlaceOrderResponseAmino): MsgPlaceOrderResponse {
    const message = createBaseMsgPlaceOrderResponse();
    return message;
  },
  toAmino(_: MsgPlaceOrderResponse): MsgPlaceOrderResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgPlaceOrderResponseAminoMsg): MsgPlaceOrderResponse {
    return MsgPlaceOrderResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgPlaceOrderResponseProtoMsg): MsgPlaceOrderResponse {
    return MsgPlaceOrderResponse.decode(message.value);
  },
  toProto(message: MsgPlaceOrderResponse): Uint8Array {
    return MsgPlaceOrderResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgPlaceOrderResponse): MsgPlaceOrderResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgPlaceOrderResponse",
      value: MsgPlaceOrderResponse.encode(message).finish()
    };
  }
};
function createBaseMsgCancelOrder(): MsgCancelOrder {
  return {
    orderId: OrderId.fromPartial({}),
    goodTilBlock: undefined,
    goodTilBlockTime: undefined
  };
}
export const MsgCancelOrder = {
  typeUrl: "/dydxprotocol.clob.MsgCancelOrder",
  encode(message: MsgCancelOrder, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.orderId !== undefined) {
      OrderId.encode(message.orderId, writer.uint32(10).fork()).ldelim();
    }
    if (message.goodTilBlock !== undefined) {
      writer.uint32(16).uint32(message.goodTilBlock);
    }
    if (message.goodTilBlockTime !== undefined) {
      writer.uint32(29).fixed32(message.goodTilBlockTime);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCancelOrder {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCancelOrder();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orderId = OrderId.decode(reader, reader.uint32());
          break;
        case 2:
          message.goodTilBlock = reader.uint32();
          break;
        case 3:
          message.goodTilBlockTime = reader.fixed32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgCancelOrder>): MsgCancelOrder {
    const message = createBaseMsgCancelOrder();
    message.orderId = object.orderId !== undefined && object.orderId !== null ? OrderId.fromPartial(object.orderId) : undefined;
    message.goodTilBlock = object.goodTilBlock ?? undefined;
    message.goodTilBlockTime = object.goodTilBlockTime ?? undefined;
    return message;
  },
  fromAmino(object: MsgCancelOrderAmino): MsgCancelOrder {
    const message = createBaseMsgCancelOrder();
    if (object.order_id !== undefined && object.order_id !== null) {
      message.orderId = OrderId.fromAmino(object.order_id);
    }
    if (object.good_til_block !== undefined && object.good_til_block !== null) {
      message.goodTilBlock = object.good_til_block;
    }
    if (object.good_til_block_time !== undefined && object.good_til_block_time !== null) {
      message.goodTilBlockTime = object.good_til_block_time;
    }
    return message;
  },
  toAmino(message: MsgCancelOrder): MsgCancelOrderAmino {
    const obj: any = {};
    obj.order_id = message.orderId ? OrderId.toAmino(message.orderId) : undefined;
    obj.good_til_block = message.goodTilBlock;
    obj.good_til_block_time = message.goodTilBlockTime;
    return obj;
  },
  fromAminoMsg(object: MsgCancelOrderAminoMsg): MsgCancelOrder {
    return MsgCancelOrder.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCancelOrderProtoMsg): MsgCancelOrder {
    return MsgCancelOrder.decode(message.value);
  },
  toProto(message: MsgCancelOrder): Uint8Array {
    return MsgCancelOrder.encode(message).finish();
  },
  toProtoMsg(message: MsgCancelOrder): MsgCancelOrderProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgCancelOrder",
      value: MsgCancelOrder.encode(message).finish()
    };
  }
};
function createBaseMsgCancelOrderResponse(): MsgCancelOrderResponse {
  return {};
}
export const MsgCancelOrderResponse = {
  typeUrl: "/dydxprotocol.clob.MsgCancelOrderResponse",
  encode(_: MsgCancelOrderResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCancelOrderResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCancelOrderResponse();
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
  fromPartial(_: Partial<MsgCancelOrderResponse>): MsgCancelOrderResponse {
    const message = createBaseMsgCancelOrderResponse();
    return message;
  },
  fromAmino(_: MsgCancelOrderResponseAmino): MsgCancelOrderResponse {
    const message = createBaseMsgCancelOrderResponse();
    return message;
  },
  toAmino(_: MsgCancelOrderResponse): MsgCancelOrderResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgCancelOrderResponseAminoMsg): MsgCancelOrderResponse {
    return MsgCancelOrderResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCancelOrderResponseProtoMsg): MsgCancelOrderResponse {
    return MsgCancelOrderResponse.decode(message.value);
  },
  toProto(message: MsgCancelOrderResponse): Uint8Array {
    return MsgCancelOrderResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgCancelOrderResponse): MsgCancelOrderResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgCancelOrderResponse",
      value: MsgCancelOrderResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateClobPair(): MsgUpdateClobPair {
  return {
    authority: "",
    clobPair: ClobPair.fromPartial({})
  };
}
export const MsgUpdateClobPair = {
  typeUrl: "/dydxprotocol.clob.MsgUpdateClobPair",
  encode(message: MsgUpdateClobPair, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.clobPair !== undefined) {
      ClobPair.encode(message.clobPair, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateClobPair {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateClobPair();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.clobPair = ClobPair.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateClobPair>): MsgUpdateClobPair {
    const message = createBaseMsgUpdateClobPair();
    message.authority = object.authority ?? "";
    message.clobPair = object.clobPair !== undefined && object.clobPair !== null ? ClobPair.fromPartial(object.clobPair) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateClobPairAmino): MsgUpdateClobPair {
    const message = createBaseMsgUpdateClobPair();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.clob_pair !== undefined && object.clob_pair !== null) {
      message.clobPair = ClobPair.fromAmino(object.clob_pair);
    }
    return message;
  },
  toAmino(message: MsgUpdateClobPair): MsgUpdateClobPairAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.clob_pair = message.clobPair ? ClobPair.toAmino(message.clobPair) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateClobPairAminoMsg): MsgUpdateClobPair {
    return MsgUpdateClobPair.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateClobPairProtoMsg): MsgUpdateClobPair {
    return MsgUpdateClobPair.decode(message.value);
  },
  toProto(message: MsgUpdateClobPair): Uint8Array {
    return MsgUpdateClobPair.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateClobPair): MsgUpdateClobPairProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgUpdateClobPair",
      value: MsgUpdateClobPair.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateClobPairResponse(): MsgUpdateClobPairResponse {
  return {};
}
export const MsgUpdateClobPairResponse = {
  typeUrl: "/dydxprotocol.clob.MsgUpdateClobPairResponse",
  encode(_: MsgUpdateClobPairResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateClobPairResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateClobPairResponse();
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
  fromPartial(_: Partial<MsgUpdateClobPairResponse>): MsgUpdateClobPairResponse {
    const message = createBaseMsgUpdateClobPairResponse();
    return message;
  },
  fromAmino(_: MsgUpdateClobPairResponseAmino): MsgUpdateClobPairResponse {
    const message = createBaseMsgUpdateClobPairResponse();
    return message;
  },
  toAmino(_: MsgUpdateClobPairResponse): MsgUpdateClobPairResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateClobPairResponseAminoMsg): MsgUpdateClobPairResponse {
    return MsgUpdateClobPairResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateClobPairResponseProtoMsg): MsgUpdateClobPairResponse {
    return MsgUpdateClobPairResponse.decode(message.value);
  },
  toProto(message: MsgUpdateClobPairResponse): Uint8Array {
    return MsgUpdateClobPairResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateClobPairResponse): MsgUpdateClobPairResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgUpdateClobPairResponse",
      value: MsgUpdateClobPairResponse.encode(message).finish()
    };
  }
};
function createBaseOperationRaw(): OperationRaw {
  return {
    match: undefined,
    shortTermOrderPlacement: undefined,
    orderRemoval: undefined
  };
}
export const OperationRaw = {
  typeUrl: "/dydxprotocol.clob.OperationRaw",
  encode(message: OperationRaw, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.match !== undefined) {
      ClobMatch.encode(message.match, writer.uint32(10).fork()).ldelim();
    }
    if (message.shortTermOrderPlacement !== undefined) {
      writer.uint32(18).bytes(message.shortTermOrderPlacement);
    }
    if (message.orderRemoval !== undefined) {
      OrderRemoval.encode(message.orderRemoval, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): OperationRaw {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOperationRaw();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.match = ClobMatch.decode(reader, reader.uint32());
          break;
        case 2:
          message.shortTermOrderPlacement = reader.bytes();
          break;
        case 3:
          message.orderRemoval = OrderRemoval.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<OperationRaw>): OperationRaw {
    const message = createBaseOperationRaw();
    message.match = object.match !== undefined && object.match !== null ? ClobMatch.fromPartial(object.match) : undefined;
    message.shortTermOrderPlacement = object.shortTermOrderPlacement ?? undefined;
    message.orderRemoval = object.orderRemoval !== undefined && object.orderRemoval !== null ? OrderRemoval.fromPartial(object.orderRemoval) : undefined;
    return message;
  },
  fromAmino(object: OperationRawAmino): OperationRaw {
    const message = createBaseOperationRaw();
    if (object.match !== undefined && object.match !== null) {
      message.match = ClobMatch.fromAmino(object.match);
    }
    if (object.short_term_order_placement !== undefined && object.short_term_order_placement !== null) {
      message.shortTermOrderPlacement = bytesFromBase64(object.short_term_order_placement);
    }
    if (object.order_removal !== undefined && object.order_removal !== null) {
      message.orderRemoval = OrderRemoval.fromAmino(object.order_removal);
    }
    return message;
  },
  toAmino(message: OperationRaw): OperationRawAmino {
    const obj: any = {};
    obj.match = message.match ? ClobMatch.toAmino(message.match) : undefined;
    obj.short_term_order_placement = message.shortTermOrderPlacement ? base64FromBytes(message.shortTermOrderPlacement) : undefined;
    obj.order_removal = message.orderRemoval ? OrderRemoval.toAmino(message.orderRemoval) : undefined;
    return obj;
  },
  fromAminoMsg(object: OperationRawAminoMsg): OperationRaw {
    return OperationRaw.fromAmino(object.value);
  },
  fromProtoMsg(message: OperationRawProtoMsg): OperationRaw {
    return OperationRaw.decode(message.value);
  },
  toProto(message: OperationRaw): Uint8Array {
    return OperationRaw.encode(message).finish();
  },
  toProtoMsg(message: OperationRaw): OperationRawProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.OperationRaw",
      value: OperationRaw.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateEquityTierLimitConfiguration(): MsgUpdateEquityTierLimitConfiguration {
  return {
    authority: "",
    equityTierLimitConfig: EquityTierLimitConfiguration.fromPartial({})
  };
}
export const MsgUpdateEquityTierLimitConfiguration = {
  typeUrl: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration",
  encode(message: MsgUpdateEquityTierLimitConfiguration, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.equityTierLimitConfig !== undefined) {
      EquityTierLimitConfiguration.encode(message.equityTierLimitConfig, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateEquityTierLimitConfiguration {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateEquityTierLimitConfiguration();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.equityTierLimitConfig = EquityTierLimitConfiguration.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateEquityTierLimitConfiguration>): MsgUpdateEquityTierLimitConfiguration {
    const message = createBaseMsgUpdateEquityTierLimitConfiguration();
    message.authority = object.authority ?? "";
    message.equityTierLimitConfig = object.equityTierLimitConfig !== undefined && object.equityTierLimitConfig !== null ? EquityTierLimitConfiguration.fromPartial(object.equityTierLimitConfig) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateEquityTierLimitConfigurationAmino): MsgUpdateEquityTierLimitConfiguration {
    const message = createBaseMsgUpdateEquityTierLimitConfiguration();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.equity_tier_limit_config !== undefined && object.equity_tier_limit_config !== null) {
      message.equityTierLimitConfig = EquityTierLimitConfiguration.fromAmino(object.equity_tier_limit_config);
    }
    return message;
  },
  toAmino(message: MsgUpdateEquityTierLimitConfiguration): MsgUpdateEquityTierLimitConfigurationAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.equity_tier_limit_config = message.equityTierLimitConfig ? EquityTierLimitConfiguration.toAmino(message.equityTierLimitConfig) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateEquityTierLimitConfigurationAminoMsg): MsgUpdateEquityTierLimitConfiguration {
    return MsgUpdateEquityTierLimitConfiguration.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateEquityTierLimitConfigurationProtoMsg): MsgUpdateEquityTierLimitConfiguration {
    return MsgUpdateEquityTierLimitConfiguration.decode(message.value);
  },
  toProto(message: MsgUpdateEquityTierLimitConfiguration): Uint8Array {
    return MsgUpdateEquityTierLimitConfiguration.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateEquityTierLimitConfiguration): MsgUpdateEquityTierLimitConfigurationProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration",
      value: MsgUpdateEquityTierLimitConfiguration.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateEquityTierLimitConfigurationResponse(): MsgUpdateEquityTierLimitConfigurationResponse {
  return {};
}
export const MsgUpdateEquityTierLimitConfigurationResponse = {
  typeUrl: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfigurationResponse",
  encode(_: MsgUpdateEquityTierLimitConfigurationResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateEquityTierLimitConfigurationResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateEquityTierLimitConfigurationResponse();
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
  fromPartial(_: Partial<MsgUpdateEquityTierLimitConfigurationResponse>): MsgUpdateEquityTierLimitConfigurationResponse {
    const message = createBaseMsgUpdateEquityTierLimitConfigurationResponse();
    return message;
  },
  fromAmino(_: MsgUpdateEquityTierLimitConfigurationResponseAmino): MsgUpdateEquityTierLimitConfigurationResponse {
    const message = createBaseMsgUpdateEquityTierLimitConfigurationResponse();
    return message;
  },
  toAmino(_: MsgUpdateEquityTierLimitConfigurationResponse): MsgUpdateEquityTierLimitConfigurationResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateEquityTierLimitConfigurationResponseAminoMsg): MsgUpdateEquityTierLimitConfigurationResponse {
    return MsgUpdateEquityTierLimitConfigurationResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateEquityTierLimitConfigurationResponseProtoMsg): MsgUpdateEquityTierLimitConfigurationResponse {
    return MsgUpdateEquityTierLimitConfigurationResponse.decode(message.value);
  },
  toProto(message: MsgUpdateEquityTierLimitConfigurationResponse): Uint8Array {
    return MsgUpdateEquityTierLimitConfigurationResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateEquityTierLimitConfigurationResponse): MsgUpdateEquityTierLimitConfigurationResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgUpdateEquityTierLimitConfigurationResponse",
      value: MsgUpdateEquityTierLimitConfigurationResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateBlockRateLimitConfiguration(): MsgUpdateBlockRateLimitConfiguration {
  return {
    authority: "",
    blockRateLimitConfig: BlockRateLimitConfiguration.fromPartial({})
  };
}
export const MsgUpdateBlockRateLimitConfiguration = {
  typeUrl: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration",
  encode(message: MsgUpdateBlockRateLimitConfiguration, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.blockRateLimitConfig !== undefined) {
      BlockRateLimitConfiguration.encode(message.blockRateLimitConfig, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateBlockRateLimitConfiguration {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateBlockRateLimitConfiguration();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 3:
          message.blockRateLimitConfig = BlockRateLimitConfiguration.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateBlockRateLimitConfiguration>): MsgUpdateBlockRateLimitConfiguration {
    const message = createBaseMsgUpdateBlockRateLimitConfiguration();
    message.authority = object.authority ?? "";
    message.blockRateLimitConfig = object.blockRateLimitConfig !== undefined && object.blockRateLimitConfig !== null ? BlockRateLimitConfiguration.fromPartial(object.blockRateLimitConfig) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateBlockRateLimitConfigurationAmino): MsgUpdateBlockRateLimitConfiguration {
    const message = createBaseMsgUpdateBlockRateLimitConfiguration();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.block_rate_limit_config !== undefined && object.block_rate_limit_config !== null) {
      message.blockRateLimitConfig = BlockRateLimitConfiguration.fromAmino(object.block_rate_limit_config);
    }
    return message;
  },
  toAmino(message: MsgUpdateBlockRateLimitConfiguration): MsgUpdateBlockRateLimitConfigurationAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.block_rate_limit_config = message.blockRateLimitConfig ? BlockRateLimitConfiguration.toAmino(message.blockRateLimitConfig) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateBlockRateLimitConfigurationAminoMsg): MsgUpdateBlockRateLimitConfiguration {
    return MsgUpdateBlockRateLimitConfiguration.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateBlockRateLimitConfigurationProtoMsg): MsgUpdateBlockRateLimitConfiguration {
    return MsgUpdateBlockRateLimitConfiguration.decode(message.value);
  },
  toProto(message: MsgUpdateBlockRateLimitConfiguration): Uint8Array {
    return MsgUpdateBlockRateLimitConfiguration.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateBlockRateLimitConfiguration): MsgUpdateBlockRateLimitConfigurationProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration",
      value: MsgUpdateBlockRateLimitConfiguration.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateBlockRateLimitConfigurationResponse(): MsgUpdateBlockRateLimitConfigurationResponse {
  return {};
}
export const MsgUpdateBlockRateLimitConfigurationResponse = {
  typeUrl: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfigurationResponse",
  encode(_: MsgUpdateBlockRateLimitConfigurationResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateBlockRateLimitConfigurationResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateBlockRateLimitConfigurationResponse();
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
  fromPartial(_: Partial<MsgUpdateBlockRateLimitConfigurationResponse>): MsgUpdateBlockRateLimitConfigurationResponse {
    const message = createBaseMsgUpdateBlockRateLimitConfigurationResponse();
    return message;
  },
  fromAmino(_: MsgUpdateBlockRateLimitConfigurationResponseAmino): MsgUpdateBlockRateLimitConfigurationResponse {
    const message = createBaseMsgUpdateBlockRateLimitConfigurationResponse();
    return message;
  },
  toAmino(_: MsgUpdateBlockRateLimitConfigurationResponse): MsgUpdateBlockRateLimitConfigurationResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateBlockRateLimitConfigurationResponseAminoMsg): MsgUpdateBlockRateLimitConfigurationResponse {
    return MsgUpdateBlockRateLimitConfigurationResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateBlockRateLimitConfigurationResponseProtoMsg): MsgUpdateBlockRateLimitConfigurationResponse {
    return MsgUpdateBlockRateLimitConfigurationResponse.decode(message.value);
  },
  toProto(message: MsgUpdateBlockRateLimitConfigurationResponse): Uint8Array {
    return MsgUpdateBlockRateLimitConfigurationResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateBlockRateLimitConfigurationResponse): MsgUpdateBlockRateLimitConfigurationResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgUpdateBlockRateLimitConfigurationResponse",
      value: MsgUpdateBlockRateLimitConfigurationResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateLiquidationsConfig(): MsgUpdateLiquidationsConfig {
  return {
    authority: "",
    liquidationsConfig: LiquidationsConfig.fromPartial({})
  };
}
export const MsgUpdateLiquidationsConfig = {
  typeUrl: "/dydxprotocol.clob.MsgUpdateLiquidationsConfig",
  encode(message: MsgUpdateLiquidationsConfig, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.liquidationsConfig !== undefined) {
      LiquidationsConfig.encode(message.liquidationsConfig, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateLiquidationsConfig {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateLiquidationsConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.liquidationsConfig = LiquidationsConfig.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateLiquidationsConfig>): MsgUpdateLiquidationsConfig {
    const message = createBaseMsgUpdateLiquidationsConfig();
    message.authority = object.authority ?? "";
    message.liquidationsConfig = object.liquidationsConfig !== undefined && object.liquidationsConfig !== null ? LiquidationsConfig.fromPartial(object.liquidationsConfig) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateLiquidationsConfigAmino): MsgUpdateLiquidationsConfig {
    const message = createBaseMsgUpdateLiquidationsConfig();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.liquidations_config !== undefined && object.liquidations_config !== null) {
      message.liquidationsConfig = LiquidationsConfig.fromAmino(object.liquidations_config);
    }
    return message;
  },
  toAmino(message: MsgUpdateLiquidationsConfig): MsgUpdateLiquidationsConfigAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.liquidations_config = message.liquidationsConfig ? LiquidationsConfig.toAmino(message.liquidationsConfig) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateLiquidationsConfigAminoMsg): MsgUpdateLiquidationsConfig {
    return MsgUpdateLiquidationsConfig.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateLiquidationsConfigProtoMsg): MsgUpdateLiquidationsConfig {
    return MsgUpdateLiquidationsConfig.decode(message.value);
  },
  toProto(message: MsgUpdateLiquidationsConfig): Uint8Array {
    return MsgUpdateLiquidationsConfig.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateLiquidationsConfig): MsgUpdateLiquidationsConfigProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgUpdateLiquidationsConfig",
      value: MsgUpdateLiquidationsConfig.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateLiquidationsConfigResponse(): MsgUpdateLiquidationsConfigResponse {
  return {};
}
export const MsgUpdateLiquidationsConfigResponse = {
  typeUrl: "/dydxprotocol.clob.MsgUpdateLiquidationsConfigResponse",
  encode(_: MsgUpdateLiquidationsConfigResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateLiquidationsConfigResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateLiquidationsConfigResponse();
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
  fromPartial(_: Partial<MsgUpdateLiquidationsConfigResponse>): MsgUpdateLiquidationsConfigResponse {
    const message = createBaseMsgUpdateLiquidationsConfigResponse();
    return message;
  },
  fromAmino(_: MsgUpdateLiquidationsConfigResponseAmino): MsgUpdateLiquidationsConfigResponse {
    const message = createBaseMsgUpdateLiquidationsConfigResponse();
    return message;
  },
  toAmino(_: MsgUpdateLiquidationsConfigResponse): MsgUpdateLiquidationsConfigResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateLiquidationsConfigResponseAminoMsg): MsgUpdateLiquidationsConfigResponse {
    return MsgUpdateLiquidationsConfigResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateLiquidationsConfigResponseProtoMsg): MsgUpdateLiquidationsConfigResponse {
    return MsgUpdateLiquidationsConfigResponse.decode(message.value);
  },
  toProto(message: MsgUpdateLiquidationsConfigResponse): Uint8Array {
    return MsgUpdateLiquidationsConfigResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateLiquidationsConfigResponse): MsgUpdateLiquidationsConfigResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MsgUpdateLiquidationsConfigResponse",
      value: MsgUpdateLiquidationsConfigResponse.encode(message).finish()
    };
  }
};