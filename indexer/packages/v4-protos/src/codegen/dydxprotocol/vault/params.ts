import { VaultStatus, VaultStatusSDKType } from "./vault";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QuotingParams stores vault quoting parameters. */

export interface QuotingParams {
  /**
   * The number of layers of orders a vault places. For example if
   * `layers=2`, a vault places 2 asks and 2 bids.
   */
  layers: number;
  /** The minimum base spread when a vault quotes around reservation price. */

  spreadMinPpm: number;
  /**
   * The buffer amount to add to min_price_change_ppm to arrive at `spread`
   * according to formula:
   * `spread = max(spread_min_ppm, min_price_change_ppm + spread_buffer_ppm)`.
   */

  spreadBufferPpm: number;
  /** The factor that determines how aggressive a vault skews its orders. */

  skewFactorPpm: number;
  /** The percentage of vault equity that each order is sized at. */

  orderSizePctPpm: number;
  /** The duration that a vault's orders are valid for. */

  orderExpirationSeconds: number;
  /**
   * The number of quote quantums in quote asset that a vault with no perpetual
   * positions must have to activate, i.e. if a vault has no perpetual positions
   * and has strictly less than this amount of quote asset, it will not
   * activate.
   */

  activationThresholdQuoteQuantums: Uint8Array;
}
/** QuotingParams stores vault quoting parameters. */

export interface QuotingParamsSDKType {
  /**
   * The number of layers of orders a vault places. For example if
   * `layers=2`, a vault places 2 asks and 2 bids.
   */
  layers: number;
  /** The minimum base spread when a vault quotes around reservation price. */

  spread_min_ppm: number;
  /**
   * The buffer amount to add to min_price_change_ppm to arrive at `spread`
   * according to formula:
   * `spread = max(spread_min_ppm, min_price_change_ppm + spread_buffer_ppm)`.
   */

  spread_buffer_ppm: number;
  /** The factor that determines how aggressive a vault skews its orders. */

  skew_factor_ppm: number;
  /** The percentage of vault equity that each order is sized at. */

  order_size_pct_ppm: number;
  /** The duration that a vault's orders are valid for. */

  order_expiration_seconds: number;
  /**
   * The number of quote quantums in quote asset that a vault with no perpetual
   * positions must have to activate, i.e. if a vault has no perpetual positions
   * and has strictly less than this amount of quote asset, it will not
   * activate.
   */

  activation_threshold_quote_quantums: Uint8Array;
}
/** VaultParams stores vault parameters. */

export interface VaultParams {
  /** Status of the vault. */
  status: VaultStatus;
  /** Quoting parameters of the vault. */

  quotingParams?: QuotingParams;
}
/** VaultParams stores vault parameters. */

export interface VaultParamsSDKType {
  /** Status of the vault. */
  status: VaultStatusSDKType;
  /** Quoting parameters of the vault. */

  quoting_params?: QuotingParamsSDKType;
}
/** OperatorParams stores parameters regarding megavault operator. */

export interface OperatorParams {
  /** Address of the operator. */
  operator: string;
  /** Metadata of the operator. */

  metadata?: OperatorMetadata;
}
/** OperatorParams stores parameters regarding megavault operator. */

export interface OperatorParamsSDKType {
  /** Address of the operator. */
  operator: string;
  /** Metadata of the operator. */

  metadata?: OperatorMetadataSDKType;
}
/** OperatorMetadata stores metadata regarding megavault operator. */

export interface OperatorMetadata {
  /** Name of the operator. */
  name: string;
  /** Description of the operator. */

  description: string;
}
/** OperatorMetadata stores metadata regarding megavault operator. */

export interface OperatorMetadataSDKType {
  /** Name of the operator. */
  name: string;
  /** Description of the operator. */

  description: string;
}
/**
 * Deprecated: Params stores `x/vault` parameters.
 * Deprecated since v6.x as is replaced by QuotingParams.
 */

export interface Params {
  /**
   * The number of layers of orders a vault places. For example if
   * `layers=2`, a vault places 2 asks and 2 bids.
   */
  layers: number;
  /** The minimum base spread when a vault quotes around reservation price. */

  spreadMinPpm: number;
  /**
   * The buffer amount to add to min_price_change_ppm to arrive at `spread`
   * according to formula:
   * `spread = max(spread_min_ppm, min_price_change_ppm + spread_buffer_ppm)`.
   */

  spreadBufferPpm: number;
  /** The factor that determines how aggressive a vault skews its orders. */

  skewFactorPpm: number;
  /** The percentage of vault equity that each order is sized at. */

  orderSizePctPpm: number;
  /** The duration that a vault's orders are valid for. */

  orderExpirationSeconds: number;
  /**
   * The number of quote quantums in quote asset that a vault with no perpetual
   * positions must have to activate, i.e. if a vault has no perpetual positions
   * and has strictly less than this amount of quote asset, it will not
   * activate.
   */

  activationThresholdQuoteQuantums: Uint8Array;
}
/**
 * Deprecated: Params stores `x/vault` parameters.
 * Deprecated since v6.x as is replaced by QuotingParams.
 */

export interface ParamsSDKType {
  /**
   * The number of layers of orders a vault places. For example if
   * `layers=2`, a vault places 2 asks and 2 bids.
   */
  layers: number;
  /** The minimum base spread when a vault quotes around reservation price. */

  spread_min_ppm: number;
  /**
   * The buffer amount to add to min_price_change_ppm to arrive at `spread`
   * according to formula:
   * `spread = max(spread_min_ppm, min_price_change_ppm + spread_buffer_ppm)`.
   */

  spread_buffer_ppm: number;
  /** The factor that determines how aggressive a vault skews its orders. */

  skew_factor_ppm: number;
  /** The percentage of vault equity that each order is sized at. */

  order_size_pct_ppm: number;
  /** The duration that a vault's orders are valid for. */

  order_expiration_seconds: number;
  /**
   * The number of quote quantums in quote asset that a vault with no perpetual
   * positions must have to activate, i.e. if a vault has no perpetual positions
   * and has strictly less than this amount of quote asset, it will not
   * activate.
   */

  activation_threshold_quote_quantums: Uint8Array;
}

function createBaseQuotingParams(): QuotingParams {
  return {
    layers: 0,
    spreadMinPpm: 0,
    spreadBufferPpm: 0,
    skewFactorPpm: 0,
    orderSizePctPpm: 0,
    orderExpirationSeconds: 0,
    activationThresholdQuoteQuantums: new Uint8Array()
  };
}

export const QuotingParams = {
  encode(message: QuotingParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.layers !== 0) {
      writer.uint32(8).uint32(message.layers);
    }

    if (message.spreadMinPpm !== 0) {
      writer.uint32(16).uint32(message.spreadMinPpm);
    }

    if (message.spreadBufferPpm !== 0) {
      writer.uint32(24).uint32(message.spreadBufferPpm);
    }

    if (message.skewFactorPpm !== 0) {
      writer.uint32(32).uint32(message.skewFactorPpm);
    }

    if (message.orderSizePctPpm !== 0) {
      writer.uint32(40).uint32(message.orderSizePctPpm);
    }

    if (message.orderExpirationSeconds !== 0) {
      writer.uint32(48).uint32(message.orderExpirationSeconds);
    }

    if (message.activationThresholdQuoteQuantums.length !== 0) {
      writer.uint32(58).bytes(message.activationThresholdQuoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuotingParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuotingParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.layers = reader.uint32();
          break;

        case 2:
          message.spreadMinPpm = reader.uint32();
          break;

        case 3:
          message.spreadBufferPpm = reader.uint32();
          break;

        case 4:
          message.skewFactorPpm = reader.uint32();
          break;

        case 5:
          message.orderSizePctPpm = reader.uint32();
          break;

        case 6:
          message.orderExpirationSeconds = reader.uint32();
          break;

        case 7:
          message.activationThresholdQuoteQuantums = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QuotingParams>): QuotingParams {
    const message = createBaseQuotingParams();
    message.layers = object.layers ?? 0;
    message.spreadMinPpm = object.spreadMinPpm ?? 0;
    message.spreadBufferPpm = object.spreadBufferPpm ?? 0;
    message.skewFactorPpm = object.skewFactorPpm ?? 0;
    message.orderSizePctPpm = object.orderSizePctPpm ?? 0;
    message.orderExpirationSeconds = object.orderExpirationSeconds ?? 0;
    message.activationThresholdQuoteQuantums = object.activationThresholdQuoteQuantums ?? new Uint8Array();
    return message;
  }

};

function createBaseVaultParams(): VaultParams {
  return {
    status: 0,
    quotingParams: undefined
  };
}

export const VaultParams = {
  encode(message: VaultParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }

    if (message.quotingParams !== undefined) {
      QuotingParams.encode(message.quotingParams, writer.uint32(18).fork()).ldelim();
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
          message.status = (reader.int32() as any);
          break;

        case 2:
          message.quotingParams = QuotingParams.decode(reader, reader.uint32());
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
    message.status = object.status ?? 0;
    message.quotingParams = object.quotingParams !== undefined && object.quotingParams !== null ? QuotingParams.fromPartial(object.quotingParams) : undefined;
    return message;
  }

};

function createBaseOperatorParams(): OperatorParams {
  return {
    operator: "",
    metadata: undefined
  };
}

export const OperatorParams = {
  encode(message: OperatorParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.operator !== "") {
      writer.uint32(10).string(message.operator);
    }

    if (message.metadata !== undefined) {
      OperatorMetadata.encode(message.metadata, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OperatorParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOperatorParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.operator = reader.string();
          break;

        case 2:
          message.metadata = OperatorMetadata.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OperatorParams>): OperatorParams {
    const message = createBaseOperatorParams();
    message.operator = object.operator ?? "";
    message.metadata = object.metadata !== undefined && object.metadata !== null ? OperatorMetadata.fromPartial(object.metadata) : undefined;
    return message;
  }

};

function createBaseOperatorMetadata(): OperatorMetadata {
  return {
    name: "",
    description: ""
  };
}

export const OperatorMetadata = {
  encode(message: OperatorMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }

    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OperatorMetadata {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOperatorMetadata();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;

        case 2:
          message.description = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OperatorMetadata>): OperatorMetadata {
    const message = createBaseOperatorMetadata();
    message.name = object.name ?? "";
    message.description = object.description ?? "";
    return message;
  }

};

function createBaseParams(): Params {
  return {
    layers: 0,
    spreadMinPpm: 0,
    spreadBufferPpm: 0,
    skewFactorPpm: 0,
    orderSizePctPpm: 0,
    orderExpirationSeconds: 0,
    activationThresholdQuoteQuantums: new Uint8Array()
  };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.layers !== 0) {
      writer.uint32(8).uint32(message.layers);
    }

    if (message.spreadMinPpm !== 0) {
      writer.uint32(16).uint32(message.spreadMinPpm);
    }

    if (message.spreadBufferPpm !== 0) {
      writer.uint32(24).uint32(message.spreadBufferPpm);
    }

    if (message.skewFactorPpm !== 0) {
      writer.uint32(32).uint32(message.skewFactorPpm);
    }

    if (message.orderSizePctPpm !== 0) {
      writer.uint32(40).uint32(message.orderSizePctPpm);
    }

    if (message.orderExpirationSeconds !== 0) {
      writer.uint32(48).uint32(message.orderExpirationSeconds);
    }

    if (message.activationThresholdQuoteQuantums.length !== 0) {
      writer.uint32(58).bytes(message.activationThresholdQuoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.layers = reader.uint32();
          break;

        case 2:
          message.spreadMinPpm = reader.uint32();
          break;

        case 3:
          message.spreadBufferPpm = reader.uint32();
          break;

        case 4:
          message.skewFactorPpm = reader.uint32();
          break;

        case 5:
          message.orderSizePctPpm = reader.uint32();
          break;

        case 6:
          message.orderExpirationSeconds = reader.uint32();
          break;

        case 7:
          message.activationThresholdQuoteQuantums = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Params>): Params {
    const message = createBaseParams();
    message.layers = object.layers ?? 0;
    message.spreadMinPpm = object.spreadMinPpm ?? 0;
    message.spreadBufferPpm = object.spreadBufferPpm ?? 0;
    message.skewFactorPpm = object.skewFactorPpm ?? 0;
    message.orderSizePctPpm = object.orderSizePctPpm ?? 0;
    message.orderExpirationSeconds = object.orderExpirationSeconds ?? 0;
    message.activationThresholdQuoteQuantums = object.activationThresholdQuoteQuantums ?? new Uint8Array();
    return message;
  }

};