import { Duration, DurationAmino, DurationSDKType } from "../../google/protobuf/duration";
import { BinaryReader, BinaryWriter } from "../../binary";
/**
 * EventParams stores parameters about which events to recognize and which
 * tokens to mint.
 */
export interface EventParams {
  /** The denom of the token to mint. */
  denom: string;
  /** The numerical chain ID of the Ethereum chain to query. */
  ethChainId: bigint;
  /** The address of the Ethereum contract to monitor for logs. */
  ethAddress: string;
}
export interface EventParamsProtoMsg {
  typeUrl: "/dydxprotocol.bridge.EventParams";
  value: Uint8Array;
}
/**
 * EventParams stores parameters about which events to recognize and which
 * tokens to mint.
 */
export interface EventParamsAmino {
  /** The denom of the token to mint. */
  denom?: string;
  /** The numerical chain ID of the Ethereum chain to query. */
  eth_chain_id?: string;
  /** The address of the Ethereum contract to monitor for logs. */
  eth_address?: string;
}
export interface EventParamsAminoMsg {
  type: "/dydxprotocol.bridge.EventParams";
  value: EventParamsAmino;
}
/**
 * EventParams stores parameters about which events to recognize and which
 * tokens to mint.
 */
export interface EventParamsSDKType {
  denom: string;
  eth_chain_id: bigint;
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
  proposeDelayDuration: Duration;
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
  skipIfBlockDelayedByDuration: Duration;
}
export interface ProposeParamsProtoMsg {
  typeUrl: "/dydxprotocol.bridge.ProposeParams";
  value: Uint8Array;
}
/** ProposeParams stores parameters for proposing to the module. */
export interface ProposeParamsAmino {
  /**
   * The maximum number of bridge events to propose per block.
   * Limits the number of events to propose in a single block
   * in-order to smooth out the flow of events.
   */
  max_bridges_per_block?: number;
  /**
   * The minimum duration to wait between a finalized bridge and
   * proposing it. This allows other validators to have enough time to
   * also recognize its occurence. Therefore the bridge daemon should
   * pool for new finalized events at least as often as this parameter.
   */
  propose_delay_duration?: DurationAmino;
  /**
   * Do not propose any events if a [0, 1_000_000) random number generator
   * generates a number smaller than this number.
   * Setting this parameter to 1_000_000 means always skipping proposing events.
   */
  skip_rate_ppm?: number;
  /**
   * Do not propose any events if the timestamp of the proposal block is
   * behind the proposers' wall-clock by at least this duration.
   */
  skip_if_block_delayed_by_duration?: DurationAmino;
}
export interface ProposeParamsAminoMsg {
  type: "/dydxprotocol.bridge.ProposeParams";
  value: ProposeParamsAmino;
}
/** ProposeParams stores parameters for proposing to the module. */
export interface ProposeParamsSDKType {
  max_bridges_per_block: number;
  propose_delay_duration: DurationSDKType;
  skip_rate_ppm: number;
  skip_if_block_delayed_by_duration: DurationSDKType;
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
export interface SafetyParamsProtoMsg {
  typeUrl: "/dydxprotocol.bridge.SafetyParams";
  value: Uint8Array;
}
/** SafetyParams stores safety parameters for the module. */
export interface SafetyParamsAmino {
  /** True if bridging is disabled. */
  is_disabled?: boolean;
  /**
   * The number of blocks that bridges accepted in-consensus will be pending
   * until the minted tokens are granted.
   */
  delay_blocks?: number;
}
export interface SafetyParamsAminoMsg {
  type: "/dydxprotocol.bridge.SafetyParams";
  value: SafetyParamsAmino;
}
/** SafetyParams stores safety parameters for the module. */
export interface SafetyParamsSDKType {
  is_disabled: boolean;
  delay_blocks: number;
}
function createBaseEventParams(): EventParams {
  return {
    denom: "",
    ethChainId: BigInt(0),
    ethAddress: ""
  };
}
export const EventParams = {
  typeUrl: "/dydxprotocol.bridge.EventParams",
  encode(message: EventParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }
    if (message.ethChainId !== BigInt(0)) {
      writer.uint32(16).uint64(message.ethChainId);
    }
    if (message.ethAddress !== "") {
      writer.uint32(26).string(message.ethAddress);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.denom = reader.string();
          break;
        case 2:
          message.ethChainId = reader.uint64();
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
  fromPartial(object: Partial<EventParams>): EventParams {
    const message = createBaseEventParams();
    message.denom = object.denom ?? "";
    message.ethChainId = object.ethChainId !== undefined && object.ethChainId !== null ? BigInt(object.ethChainId.toString()) : BigInt(0);
    message.ethAddress = object.ethAddress ?? "";
    return message;
  },
  fromAmino(object: EventParamsAmino): EventParams {
    const message = createBaseEventParams();
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = object.denom;
    }
    if (object.eth_chain_id !== undefined && object.eth_chain_id !== null) {
      message.ethChainId = BigInt(object.eth_chain_id);
    }
    if (object.eth_address !== undefined && object.eth_address !== null) {
      message.ethAddress = object.eth_address;
    }
    return message;
  },
  toAmino(message: EventParams): EventParamsAmino {
    const obj: any = {};
    obj.denom = message.denom;
    obj.eth_chain_id = message.ethChainId ? message.ethChainId.toString() : undefined;
    obj.eth_address = message.ethAddress;
    return obj;
  },
  fromAminoMsg(object: EventParamsAminoMsg): EventParams {
    return EventParams.fromAmino(object.value);
  },
  fromProtoMsg(message: EventParamsProtoMsg): EventParams {
    return EventParams.decode(message.value);
  },
  toProto(message: EventParams): Uint8Array {
    return EventParams.encode(message).finish();
  },
  toProtoMsg(message: EventParams): EventParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.EventParams",
      value: EventParams.encode(message).finish()
    };
  }
};
function createBaseProposeParams(): ProposeParams {
  return {
    maxBridgesPerBlock: 0,
    proposeDelayDuration: Duration.fromPartial({}),
    skipRatePpm: 0,
    skipIfBlockDelayedByDuration: Duration.fromPartial({})
  };
}
export const ProposeParams = {
  typeUrl: "/dydxprotocol.bridge.ProposeParams",
  encode(message: ProposeParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): ProposeParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<ProposeParams>): ProposeParams {
    const message = createBaseProposeParams();
    message.maxBridgesPerBlock = object.maxBridgesPerBlock ?? 0;
    message.proposeDelayDuration = object.proposeDelayDuration !== undefined && object.proposeDelayDuration !== null ? Duration.fromPartial(object.proposeDelayDuration) : undefined;
    message.skipRatePpm = object.skipRatePpm ?? 0;
    message.skipIfBlockDelayedByDuration = object.skipIfBlockDelayedByDuration !== undefined && object.skipIfBlockDelayedByDuration !== null ? Duration.fromPartial(object.skipIfBlockDelayedByDuration) : undefined;
    return message;
  },
  fromAmino(object: ProposeParamsAmino): ProposeParams {
    const message = createBaseProposeParams();
    if (object.max_bridges_per_block !== undefined && object.max_bridges_per_block !== null) {
      message.maxBridgesPerBlock = object.max_bridges_per_block;
    }
    if (object.propose_delay_duration !== undefined && object.propose_delay_duration !== null) {
      message.proposeDelayDuration = Duration.fromAmino(object.propose_delay_duration);
    }
    if (object.skip_rate_ppm !== undefined && object.skip_rate_ppm !== null) {
      message.skipRatePpm = object.skip_rate_ppm;
    }
    if (object.skip_if_block_delayed_by_duration !== undefined && object.skip_if_block_delayed_by_duration !== null) {
      message.skipIfBlockDelayedByDuration = Duration.fromAmino(object.skip_if_block_delayed_by_duration);
    }
    return message;
  },
  toAmino(message: ProposeParams): ProposeParamsAmino {
    const obj: any = {};
    obj.max_bridges_per_block = message.maxBridgesPerBlock;
    obj.propose_delay_duration = message.proposeDelayDuration ? Duration.toAmino(message.proposeDelayDuration) : undefined;
    obj.skip_rate_ppm = message.skipRatePpm;
    obj.skip_if_block_delayed_by_duration = message.skipIfBlockDelayedByDuration ? Duration.toAmino(message.skipIfBlockDelayedByDuration) : undefined;
    return obj;
  },
  fromAminoMsg(object: ProposeParamsAminoMsg): ProposeParams {
    return ProposeParams.fromAmino(object.value);
  },
  fromProtoMsg(message: ProposeParamsProtoMsg): ProposeParams {
    return ProposeParams.decode(message.value);
  },
  toProto(message: ProposeParams): Uint8Array {
    return ProposeParams.encode(message).finish();
  },
  toProtoMsg(message: ProposeParams): ProposeParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.ProposeParams",
      value: ProposeParams.encode(message).finish()
    };
  }
};
function createBaseSafetyParams(): SafetyParams {
  return {
    isDisabled: false,
    delayBlocks: 0
  };
}
export const SafetyParams = {
  typeUrl: "/dydxprotocol.bridge.SafetyParams",
  encode(message: SafetyParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.isDisabled === true) {
      writer.uint32(8).bool(message.isDisabled);
    }
    if (message.delayBlocks !== 0) {
      writer.uint32(16).uint32(message.delayBlocks);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): SafetyParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<SafetyParams>): SafetyParams {
    const message = createBaseSafetyParams();
    message.isDisabled = object.isDisabled ?? false;
    message.delayBlocks = object.delayBlocks ?? 0;
    return message;
  },
  fromAmino(object: SafetyParamsAmino): SafetyParams {
    const message = createBaseSafetyParams();
    if (object.is_disabled !== undefined && object.is_disabled !== null) {
      message.isDisabled = object.is_disabled;
    }
    if (object.delay_blocks !== undefined && object.delay_blocks !== null) {
      message.delayBlocks = object.delay_blocks;
    }
    return message;
  },
  toAmino(message: SafetyParams): SafetyParamsAmino {
    const obj: any = {};
    obj.is_disabled = message.isDisabled;
    obj.delay_blocks = message.delayBlocks;
    return obj;
  },
  fromAminoMsg(object: SafetyParamsAminoMsg): SafetyParams {
    return SafetyParams.fromAmino(object.value);
  },
  fromProtoMsg(message: SafetyParamsProtoMsg): SafetyParams {
    return SafetyParams.decode(message.value);
  },
  toProto(message: SafetyParams): Uint8Array {
    return SafetyParams.encode(message).finish();
  },
  toProtoMsg(message: SafetyParams): SafetyParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.SafetyParams",
      value: SafetyParams.encode(message).finish()
    };
  }
};