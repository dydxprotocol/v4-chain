import { Duration, DurationAmino, DurationSDKType } from "../../google/protobuf/duration";
import { BinaryReader, BinaryWriter } from "../../binary";
/** DowntimeParams defines the parameters for downtime. */
export interface DowntimeParams {
  /**
   * Durations tracked for downtime. The durations must be sorted from
   * shortest to longest and must all be positive.
   */
  durations: Duration[];
}
export interface DowntimeParamsProtoMsg {
  typeUrl: "/dydxprotocol.blocktime.DowntimeParams";
  value: Uint8Array;
}
/** DowntimeParams defines the parameters for downtime. */
export interface DowntimeParamsAmino {
  /**
   * Durations tracked for downtime. The durations must be sorted from
   * shortest to longest and must all be positive.
   */
  durations?: DurationAmino[];
}
export interface DowntimeParamsAminoMsg {
  type: "/dydxprotocol.blocktime.DowntimeParams";
  value: DowntimeParamsAmino;
}
/** DowntimeParams defines the parameters for downtime. */
export interface DowntimeParamsSDKType {
  durations: DurationSDKType[];
}
function createBaseDowntimeParams(): DowntimeParams {
  return {
    durations: []
  };
}
export const DowntimeParams = {
  typeUrl: "/dydxprotocol.blocktime.DowntimeParams",
  encode(message: DowntimeParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.durations) {
      Duration.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): DowntimeParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<DowntimeParams>): DowntimeParams {
    const message = createBaseDowntimeParams();
    message.durations = object.durations?.map(e => Duration.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: DowntimeParamsAmino): DowntimeParams {
    const message = createBaseDowntimeParams();
    message.durations = object.durations?.map(e => Duration.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: DowntimeParams): DowntimeParamsAmino {
    const obj: any = {};
    if (message.durations) {
      obj.durations = message.durations.map(e => e ? Duration.toAmino(e) : undefined);
    } else {
      obj.durations = [];
    }
    return obj;
  },
  fromAminoMsg(object: DowntimeParamsAminoMsg): DowntimeParams {
    return DowntimeParams.fromAmino(object.value);
  },
  fromProtoMsg(message: DowntimeParamsProtoMsg): DowntimeParams {
    return DowntimeParams.decode(message.value);
  },
  toProto(message: DowntimeParams): Uint8Array {
    return DowntimeParams.encode(message).finish();
  },
  toProtoMsg(message: DowntimeParams): DowntimeParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.blocktime.DowntimeParams",
      value: DowntimeParams.encode(message).finish()
    };
  }
};