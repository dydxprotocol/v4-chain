import { Controls, ControlsSDKType } from "./canonical_usdc";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgUpdateControls updates all controls atomically. */

export interface MsgUpdateControls {
  /** authority is the x/gov or x/delaymsg module address. */
  authority: string;
  /** controls is the full replacement set of canonical routing controls. */

  controls?: Controls;
}
/** MsgUpdateControls updates all controls atomically. */

export interface MsgUpdateControlsSDKType {
  /** authority is the x/gov or x/delaymsg module address. */
  authority: string;
  /** controls is the full replacement set of canonical routing controls. */

  controls?: ControlsSDKType;
}
/** MsgUpdateControlsResponse is the controls update response. */

export interface MsgUpdateControlsResponse {}
/** MsgUpdateControlsResponse is the controls update response. */

export interface MsgUpdateControlsResponseSDKType {}

function createBaseMsgUpdateControls(): MsgUpdateControls {
  return {
    authority: "",
    controls: undefined
  };
}

export const MsgUpdateControls = {
  encode(message: MsgUpdateControls, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.controls !== undefined) {
      Controls.encode(message.controls, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateControls {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateControls();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.controls = Controls.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdateControls>): MsgUpdateControls {
    const message = createBaseMsgUpdateControls();
    message.authority = object.authority ?? "";
    message.controls = object.controls !== undefined && object.controls !== null ? Controls.fromPartial(object.controls) : undefined;
    return message;
  }

};

function createBaseMsgUpdateControlsResponse(): MsgUpdateControlsResponse {
  return {};
}

export const MsgUpdateControlsResponse = {
  encode(_: MsgUpdateControlsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateControlsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateControlsResponse();

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

  fromPartial(_: DeepPartial<MsgUpdateControlsResponse>): MsgUpdateControlsResponse {
    const message = createBaseMsgUpdateControlsResponse();
    return message;
  }

};