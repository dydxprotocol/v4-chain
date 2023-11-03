import { Duration, DurationSDKType } from "../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { Long, DeepPartial } from "../../helpers";
/**
 * EventParams stores parameters about which events to recognize and which
 * tokens to mint.
 */

export interface EventParams {
  /** The denom of the token to mint. */
  denom: string;
  /** The numerical chain ID of the Ethereum chain to query. */

  ethChainId: Long;
  /** The address of the Ethereum contract to monitor for logs. */

  ethAddress: string;
}
/**
 * EventParams stores parameters about which events to recognize and which
 * tokens to mint.
 */

export interface EventParamsSDKType {
  /** The denom of the token to mint. */
  denom: string;
  /** The numerical chain ID of the Ethereum chain to query. */

  eth_chain_id: Long;
  /** The address of the Ethereum contract to monitor for logs. */

  eth_address: string;
}
/** ProposeParams stores parameters for proposing to the module. */

export interface ProposeParams {
  /**
   * The maximum number of bridge events to propose per block.
   * Limits the number of events to propose in a single block
   * in-order to smooth out the flow of events.
   */
  maxBridgesPerBlock: number;
  /**
   * The minimum duration to wait between a finalized bridge and
   * proposing it. This allows other validators to have enough time to
   * also recognize its occurence. Therefore the bridge daemon should
   * pool for new finalized events at least as often as this parameter.
   */

  proposeDelayDuration?: Duration;
  /**
   * Do not propose any events if a [0, 1_000_000) random number generator
   * generates a number smaller than this number.
   * Setting this parameter to 1_000_000 means always skipping proposing events.
   */

  skipRatePpm: number;
  /**
   * Do not propose any events if the timestamp of the proposal block is
   * behind the proposers' wall-clock by at least this duration.
   */

  skipIfBlockDelayedByDuration?: Duration;
}
/** ProposeParams stores parameters for proposing to the module. */

export interface ProposeParamsSDKType {
  /**
   * The maximum number of bridge events to propose per block.
   * Limits the number of events to propose in a single block
   * in-order to smooth out the flow of events.
   */
  max_bridges_per_block: number;
  /**
   * The minimum duration to wait between a finalized bridge and
   * proposing it. This allows other validators to have enough time to
   * also recognize its occurence. Therefore the bridge daemon should
   * pool for new finalized events at least as often as this parameter.
   */

  propose_delay_duration?: DurationSDKType;
  /**
   * Do not propose any events if a [0, 1_000_000) random number generator
   * generates a number smaller than this number.
   * Setting this parameter to 1_000_000 means always skipping proposing events.
   */

  skip_rate_ppm: number;
  /**
   * Do not propose any events if the timestamp of the proposal block is
   * behind the proposers' wall-clock by at least this duration.
   */

  skip_if_block_delayed_by_duration?: DurationSDKType;
}
/** SafetyParams stores safety parameters for the module. */

export interface SafetyParams {
  /** True if bridging is disabled. */
  isDisabled: boolean;
  /**
   * The number of blocks that bridges accepted in-consensus will be pending
   * until the minted tokens are granted.
   */

  delayBlocks: number;
}
/** SafetyParams stores safety parameters for the module. */

export interface SafetyParamsSDKType {
  /** True if bridging is disabled. */
  is_disabled: boolean;
  /**
   * The number of blocks that bridges accepted in-consensus will be pending
   * until the minted tokens are granted.
   */

  delay_blocks: number;
}

function createBaseEventParams(): EventParams {
  return {
    denom: "",
    ethChainId: Long.UZERO,
    ethAddress: ""
  };
}

export const EventParams = {
  encode(message: EventParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }

    if (!message.ethChainId.isZero()) {
      writer.uint32(16).uint64(message.ethChainId);
    }

    if (message.ethAddress !== "") {
      writer.uint32(26).string(message.ethAddress);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.denom = reader.string();
          break;

        case 2:
          message.ethChainId = (reader.uint64() as Long);
          break;

        case 3:
          message.ethAddress = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<EventParams>): EventParams {
    const message = createBaseEventParams();
    message.denom = object.denom ?? "";
    message.ethChainId = object.ethChainId !== undefined && object.ethChainId !== null ? Long.fromValue(object.ethChainId) : Long.UZERO;
    message.ethAddress = object.ethAddress ?? "";
    return message;
  }

};

function createBaseProposeParams(): ProposeParams {
  return {
    maxBridgesPerBlock: 0,
    proposeDelayDuration: undefined,
    skipRatePpm: 0,
    skipIfBlockDelayedByDuration: undefined
  };
}

export const ProposeParams = {
  encode(message: ProposeParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.maxBridgesPerBlock !== 0) {
      writer.uint32(8).uint32(message.maxBridgesPerBlock);
    }

    if (message.proposeDelayDuration !== undefined) {
      Duration.encode(message.proposeDelayDuration, writer.uint32(18).fork()).ldelim();
    }

    if (message.skipRatePpm !== 0) {
      writer.uint32(24).uint32(message.skipRatePpm);
    }

    if (message.skipIfBlockDelayedByDuration !== undefined) {
      Duration.encode(message.skipIfBlockDelayedByDuration, writer.uint32(34).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProposeParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProposeParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.maxBridgesPerBlock = reader.uint32();
          break;

        case 2:
          message.proposeDelayDuration = Duration.decode(reader, reader.uint32());
          break;

        case 3:
          message.skipRatePpm = reader.uint32();
          break;

        case 4:
          message.skipIfBlockDelayedByDuration = Duration.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ProposeParams>): ProposeParams {
    const message = createBaseProposeParams();
    message.maxBridgesPerBlock = object.maxBridgesPerBlock ?? 0;
    message.proposeDelayDuration = object.proposeDelayDuration !== undefined && object.proposeDelayDuration !== null ? Duration.fromPartial(object.proposeDelayDuration) : undefined;
    message.skipRatePpm = object.skipRatePpm ?? 0;
    message.skipIfBlockDelayedByDuration = object.skipIfBlockDelayedByDuration !== undefined && object.skipIfBlockDelayedByDuration !== null ? Duration.fromPartial(object.skipIfBlockDelayedByDuration) : undefined;
    return message;
  }

};

function createBaseSafetyParams(): SafetyParams {
  return {
    isDisabled: false,
    delayBlocks: 0
  };
}

export const SafetyParams = {
  encode(message: SafetyParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.isDisabled === true) {
      writer.uint32(8).bool(message.isDisabled);
    }

    if (message.delayBlocks !== 0) {
      writer.uint32(16).uint32(message.delayBlocks);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SafetyParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSafetyParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.isDisabled = reader.bool();
          break;

        case 2:
          message.delayBlocks = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<SafetyParams>): SafetyParams {
    const message = createBaseSafetyParams();
    message.isDisabled = object.isDisabled ?? false;
    message.delayBlocks = object.delayBlocks ?? 0;
    return message;
  }

};