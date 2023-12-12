import { Timestamp } from "../../google/protobuf/timestamp";
import * as _m0 from "protobufjs/minimal";
import { toTimestamp, fromTimestamp, DeepPartial } from "../../helpers";
/**
 * VestEntry specifies a Vester Account and the rate at which tokens are
 * dripped into the corresponding Treasury Account.
 */

export interface VestEntry {
  /**
   * The module account to vest tokens from.
   * This is also the key to this `VestEntry` in state.
   */
  vesterAccount: string;
  /** The module account to vest tokens to. */

  treasuryAccount: string;
  /** The denom of the token to vest. */

  denom: string;
  /** The start time of vest. Before this time, no vest will occur. */

  startTime?: Date;
  /**
   * The end time of vest. At this target date, all funds should be in the
   * Treasury Account and none left in the Vester Account.
   */

  endTime?: Date;
}
/**
 * VestEntry specifies a Vester Account and the rate at which tokens are
 * dripped into the corresponding Treasury Account.
 */

export interface VestEntrySDKType {
  /**
   * The module account to vest tokens from.
   * This is also the key to this `VestEntry` in state.
   */
  vester_account: string;
  /** The module account to vest tokens to. */

  treasury_account: string;
  /** The denom of the token to vest. */

  denom: string;
  /** The start time of vest. Before this time, no vest will occur. */

  start_time?: Date;
  /**
   * The end time of vest. At this target date, all funds should be in the
   * Treasury Account and none left in the Vester Account.
   */

  end_time?: Date;
}

function createBaseVestEntry(): VestEntry {
  return {
    vesterAccount: "",
    treasuryAccount: "",
    denom: "",
    startTime: undefined,
    endTime: undefined
  };
}

export const VestEntry = {
  encode(message: VestEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.vesterAccount !== "") {
      writer.uint32(10).string(message.vesterAccount);
    }

    if (message.treasuryAccount !== "") {
      writer.uint32(18).string(message.treasuryAccount);
    }

    if (message.denom !== "") {
      writer.uint32(26).string(message.denom);
    }

    if (message.startTime !== undefined) {
      Timestamp.encode(toTimestamp(message.startTime), writer.uint32(34).fork()).ldelim();
    }

    if (message.endTime !== undefined) {
      Timestamp.encode(toTimestamp(message.endTime), writer.uint32(42).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VestEntry {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVestEntry();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.vesterAccount = reader.string();
          break;

        case 2:
          message.treasuryAccount = reader.string();
          break;

        case 3:
          message.denom = reader.string();
          break;

        case 4:
          message.startTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;

        case 5:
          message.endTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<VestEntry>): VestEntry {
    const message = createBaseVestEntry();
    message.vesterAccount = object.vesterAccount ?? "";
    message.treasuryAccount = object.treasuryAccount ?? "";
    message.denom = object.denom ?? "";
    message.startTime = object.startTime ?? undefined;
    message.endTime = object.endTime ?? undefined;
    return message;
  }

};