import { Duration, DurationSDKType } from "../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** DowntimeParams defines the parameters for downtime. */

export interface DowntimeParams {
  /**
   * Durations tracked for downtime. The durations must be sorted from
   * shortest to longest and must all be positive.
   */
  durations: Duration[];
}
/** DowntimeParams defines the parameters for downtime. */

export interface DowntimeParamsSDKType {
  /**
   * Durations tracked for downtime. The durations must be sorted from
   * shortest to longest and must all be positive.
   */
  durations: DurationSDKType[];
}
/** SynchronyParams defines the parameters for block synchrony. */

export interface SynchronyParams {
  /**
   * next_block_delay replaces the locally configured timeout_commit in
   * CometBFT. It determines the amount of time the CometBFT waits after the
   * `CommitTime` (subjective time when +2/3 precommits were received), before
   * moving to next height.
   * If the application sends next_block_delay = 0 to the consensus engine, the
   * latter defaults back to using timeout_commit.
   */
  nextBlockDelay?: Duration;
}
/** SynchronyParams defines the parameters for block synchrony. */

export interface SynchronyParamsSDKType {
  /**
   * next_block_delay replaces the locally configured timeout_commit in
   * CometBFT. It determines the amount of time the CometBFT waits after the
   * `CommitTime` (subjective time when +2/3 precommits were received), before
   * moving to next height.
   * If the application sends next_block_delay = 0 to the consensus engine, the
   * latter defaults back to using timeout_commit.
   */
  next_block_delay?: DurationSDKType;
}

function createBaseDowntimeParams(): DowntimeParams {
  return {
    durations: []
  };
}

export const DowntimeParams = {
  encode(message: DowntimeParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.durations) {
      Duration.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DowntimeParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDowntimeParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.durations.push(Duration.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<DowntimeParams>): DowntimeParams {
    const message = createBaseDowntimeParams();
    message.durations = object.durations?.map(e => Duration.fromPartial(e)) || [];
    return message;
  }

};

function createBaseSynchronyParams(): SynchronyParams {
  return {
    nextBlockDelay: undefined
  };
}

export const SynchronyParams = {
  encode(message: SynchronyParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nextBlockDelay !== undefined) {
      Duration.encode(message.nextBlockDelay, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SynchronyParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSynchronyParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.nextBlockDelay = Duration.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<SynchronyParams>): SynchronyParams {
    const message = createBaseSynchronyParams();
    message.nextBlockDelay = object.nextBlockDelay !== undefined && object.nextBlockDelay !== null ? Duration.fromPartial(object.nextBlockDelay) : undefined;
    return message;
  }

};