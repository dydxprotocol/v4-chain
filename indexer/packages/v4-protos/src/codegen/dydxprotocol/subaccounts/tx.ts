import { SubaccountId, SubaccountIdSDKType } from "./subaccount";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgClaimYieldForSubaccount is the Msg/ClaimYieldForSubaccount request type. */

export interface MsgClaimYieldForSubaccount {
  id?: SubaccountId;
}
/** MsgClaimYieldForSubaccount is the Msg/ClaimYieldForSubaccount request type. */

export interface MsgClaimYieldForSubaccountSDKType {
  id?: SubaccountIdSDKType;
}
/** MsgClaimYieldForSubaccountResponse is the Msg/ClaimYieldForSubaccount response type. */

export interface MsgClaimYieldForSubaccountResponse {}
/** MsgClaimYieldForSubaccountResponse is the Msg/ClaimYieldForSubaccount response type. */

export interface MsgClaimYieldForSubaccountResponseSDKType {}

function createBaseMsgClaimYieldForSubaccount(): MsgClaimYieldForSubaccount {
  return {
    id: undefined
  };
}

export const MsgClaimYieldForSubaccount = {
  encode(message: MsgClaimYieldForSubaccount, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined) {
      SubaccountId.encode(message.id, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaimYieldForSubaccount {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgClaimYieldForSubaccount();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = SubaccountId.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgClaimYieldForSubaccount>): MsgClaimYieldForSubaccount {
    const message = createBaseMsgClaimYieldForSubaccount();
    message.id = object.id !== undefined && object.id !== null ? SubaccountId.fromPartial(object.id) : undefined;
    return message;
  }

};

function createBaseMsgClaimYieldForSubaccountResponse(): MsgClaimYieldForSubaccountResponse {
  return {};
}

export const MsgClaimYieldForSubaccountResponse = {
  encode(_: MsgClaimYieldForSubaccountResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaimYieldForSubaccountResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgClaimYieldForSubaccountResponse();

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

  fromPartial(_: DeepPartial<MsgClaimYieldForSubaccountResponse>): MsgClaimYieldForSubaccountResponse {
    const message = createBaseMsgClaimYieldForSubaccountResponse();
    return message;
  }

};