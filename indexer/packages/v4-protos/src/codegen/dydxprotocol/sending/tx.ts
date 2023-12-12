import { Transfer, TransferAmino, TransferSDKType } from "./transfer";
import { BinaryReader, BinaryWriter } from "../../binary";
/** MsgCreateTransfer is a request type used for initiating new transfers. */
export interface MsgCreateTransfer {
  /** MsgCreateTransfer is a request type used for initiating new transfers. */
  transfer?: Transfer;
}
export interface MsgCreateTransferProtoMsg {
  typeUrl: "/dydxprotocol.sending.MsgCreateTransfer";
  value: Uint8Array;
}
/** MsgCreateTransfer is a request type used for initiating new transfers. */
export interface MsgCreateTransferAmino {
  /** MsgCreateTransfer is a request type used for initiating new transfers. */
  transfer?: TransferAmino;
}
export interface MsgCreateTransferAminoMsg {
  type: "/dydxprotocol.sending.MsgCreateTransfer";
  value: MsgCreateTransferAmino;
}
/** MsgCreateTransfer is a request type used for initiating new transfers. */
export interface MsgCreateTransferSDKType {
  transfer?: TransferSDKType;
}
/** MsgCreateTransferResponse is a response type used for new transfers. */
export interface MsgCreateTransferResponse {}
export interface MsgCreateTransferResponseProtoMsg {
  typeUrl: "/dydxprotocol.sending.MsgCreateTransferResponse";
  value: Uint8Array;
}
/** MsgCreateTransferResponse is a response type used for new transfers. */
export interface MsgCreateTransferResponseAmino {}
export interface MsgCreateTransferResponseAminoMsg {
  type: "/dydxprotocol.sending.MsgCreateTransferResponse";
  value: MsgCreateTransferResponseAmino;
}
/** MsgCreateTransferResponse is a response type used for new transfers. */
export interface MsgCreateTransferResponseSDKType {}
/**
 * MsgDepositToSubaccountResponse is a response type used for new
 * account-to-subaccount transfers.
 */
export interface MsgDepositToSubaccountResponse {}
export interface MsgDepositToSubaccountResponseProtoMsg {
  typeUrl: "/dydxprotocol.sending.MsgDepositToSubaccountResponse";
  value: Uint8Array;
}
/**
 * MsgDepositToSubaccountResponse is a response type used for new
 * account-to-subaccount transfers.
 */
export interface MsgDepositToSubaccountResponseAmino {}
export interface MsgDepositToSubaccountResponseAminoMsg {
  type: "/dydxprotocol.sending.MsgDepositToSubaccountResponse";
  value: MsgDepositToSubaccountResponseAmino;
}
/**
 * MsgDepositToSubaccountResponse is a response type used for new
 * account-to-subaccount transfers.
 */
export interface MsgDepositToSubaccountResponseSDKType {}
/**
 * MsgWithdrawFromSubaccountResponse is a response type used for new
 * subaccount-to-account transfers.
 */
export interface MsgWithdrawFromSubaccountResponse {}
export interface MsgWithdrawFromSubaccountResponseProtoMsg {
  typeUrl: "/dydxprotocol.sending.MsgWithdrawFromSubaccountResponse";
  value: Uint8Array;
}
/**
 * MsgWithdrawFromSubaccountResponse is a response type used for new
 * subaccount-to-account transfers.
 */
export interface MsgWithdrawFromSubaccountResponseAmino {}
export interface MsgWithdrawFromSubaccountResponseAminoMsg {
  type: "/dydxprotocol.sending.MsgWithdrawFromSubaccountResponse";
  value: MsgWithdrawFromSubaccountResponseAmino;
}
/**
 * MsgWithdrawFromSubaccountResponse is a response type used for new
 * subaccount-to-account transfers.
 */
export interface MsgWithdrawFromSubaccountResponseSDKType {}
/**
 * MsgSendFromModuleToAccountResponse is a response type used for new
 * module-to-account transfers.
 */
export interface MsgSendFromModuleToAccountResponse {}
export interface MsgSendFromModuleToAccountResponseProtoMsg {
  typeUrl: "/dydxprotocol.sending.MsgSendFromModuleToAccountResponse";
  value: Uint8Array;
}
/**
 * MsgSendFromModuleToAccountResponse is a response type used for new
 * module-to-account transfers.
 */
export interface MsgSendFromModuleToAccountResponseAmino {}
export interface MsgSendFromModuleToAccountResponseAminoMsg {
  type: "/dydxprotocol.sending.MsgSendFromModuleToAccountResponse";
  value: MsgSendFromModuleToAccountResponseAmino;
}
/**
 * MsgSendFromModuleToAccountResponse is a response type used for new
 * module-to-account transfers.
 */
export interface MsgSendFromModuleToAccountResponseSDKType {}
function createBaseMsgCreateTransfer(): MsgCreateTransfer {
  return {
    transfer: undefined
  };
}
export const MsgCreateTransfer = {
  typeUrl: "/dydxprotocol.sending.MsgCreateTransfer",
  encode(message: MsgCreateTransfer, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.transfer !== undefined) {
      Transfer.encode(message.transfer, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateTransfer {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateTransfer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.transfer = Transfer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgCreateTransfer>): MsgCreateTransfer {
    const message = createBaseMsgCreateTransfer();
    message.transfer = object.transfer !== undefined && object.transfer !== null ? Transfer.fromPartial(object.transfer) : undefined;
    return message;
  },
  fromAmino(object: MsgCreateTransferAmino): MsgCreateTransfer {
    const message = createBaseMsgCreateTransfer();
    if (object.transfer !== undefined && object.transfer !== null) {
      message.transfer = Transfer.fromAmino(object.transfer);
    }
    return message;
  },
  toAmino(message: MsgCreateTransfer): MsgCreateTransferAmino {
    const obj: any = {};
    obj.transfer = message.transfer ? Transfer.toAmino(message.transfer) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgCreateTransferAminoMsg): MsgCreateTransfer {
    return MsgCreateTransfer.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCreateTransferProtoMsg): MsgCreateTransfer {
    return MsgCreateTransfer.decode(message.value);
  },
  toProto(message: MsgCreateTransfer): Uint8Array {
    return MsgCreateTransfer.encode(message).finish();
  },
  toProtoMsg(message: MsgCreateTransfer): MsgCreateTransferProtoMsg {
    return {
      typeUrl: "/dydxprotocol.sending.MsgCreateTransfer",
      value: MsgCreateTransfer.encode(message).finish()
    };
  }
};
function createBaseMsgCreateTransferResponse(): MsgCreateTransferResponse {
  return {};
}
export const MsgCreateTransferResponse = {
  typeUrl: "/dydxprotocol.sending.MsgCreateTransferResponse",
  encode(_: MsgCreateTransferResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateTransferResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateTransferResponse();
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
  fromPartial(_: Partial<MsgCreateTransferResponse>): MsgCreateTransferResponse {
    const message = createBaseMsgCreateTransferResponse();
    return message;
  },
  fromAmino(_: MsgCreateTransferResponseAmino): MsgCreateTransferResponse {
    const message = createBaseMsgCreateTransferResponse();
    return message;
  },
  toAmino(_: MsgCreateTransferResponse): MsgCreateTransferResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgCreateTransferResponseAminoMsg): MsgCreateTransferResponse {
    return MsgCreateTransferResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCreateTransferResponseProtoMsg): MsgCreateTransferResponse {
    return MsgCreateTransferResponse.decode(message.value);
  },
  toProto(message: MsgCreateTransferResponse): Uint8Array {
    return MsgCreateTransferResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgCreateTransferResponse): MsgCreateTransferResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.sending.MsgCreateTransferResponse",
      value: MsgCreateTransferResponse.encode(message).finish()
    };
  }
};
function createBaseMsgDepositToSubaccountResponse(): MsgDepositToSubaccountResponse {
  return {};
}
export const MsgDepositToSubaccountResponse = {
  typeUrl: "/dydxprotocol.sending.MsgDepositToSubaccountResponse",
  encode(_: MsgDepositToSubaccountResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgDepositToSubaccountResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDepositToSubaccountResponse();
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
  fromPartial(_: Partial<MsgDepositToSubaccountResponse>): MsgDepositToSubaccountResponse {
    const message = createBaseMsgDepositToSubaccountResponse();
    return message;
  },
  fromAmino(_: MsgDepositToSubaccountResponseAmino): MsgDepositToSubaccountResponse {
    const message = createBaseMsgDepositToSubaccountResponse();
    return message;
  },
  toAmino(_: MsgDepositToSubaccountResponse): MsgDepositToSubaccountResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgDepositToSubaccountResponseAminoMsg): MsgDepositToSubaccountResponse {
    return MsgDepositToSubaccountResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgDepositToSubaccountResponseProtoMsg): MsgDepositToSubaccountResponse {
    return MsgDepositToSubaccountResponse.decode(message.value);
  },
  toProto(message: MsgDepositToSubaccountResponse): Uint8Array {
    return MsgDepositToSubaccountResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgDepositToSubaccountResponse): MsgDepositToSubaccountResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.sending.MsgDepositToSubaccountResponse",
      value: MsgDepositToSubaccountResponse.encode(message).finish()
    };
  }
};
function createBaseMsgWithdrawFromSubaccountResponse(): MsgWithdrawFromSubaccountResponse {
  return {};
}
export const MsgWithdrawFromSubaccountResponse = {
  typeUrl: "/dydxprotocol.sending.MsgWithdrawFromSubaccountResponse",
  encode(_: MsgWithdrawFromSubaccountResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgWithdrawFromSubaccountResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFromSubaccountResponse();
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
  fromPartial(_: Partial<MsgWithdrawFromSubaccountResponse>): MsgWithdrawFromSubaccountResponse {
    const message = createBaseMsgWithdrawFromSubaccountResponse();
    return message;
  },
  fromAmino(_: MsgWithdrawFromSubaccountResponseAmino): MsgWithdrawFromSubaccountResponse {
    const message = createBaseMsgWithdrawFromSubaccountResponse();
    return message;
  },
  toAmino(_: MsgWithdrawFromSubaccountResponse): MsgWithdrawFromSubaccountResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgWithdrawFromSubaccountResponseAminoMsg): MsgWithdrawFromSubaccountResponse {
    return MsgWithdrawFromSubaccountResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgWithdrawFromSubaccountResponseProtoMsg): MsgWithdrawFromSubaccountResponse {
    return MsgWithdrawFromSubaccountResponse.decode(message.value);
  },
  toProto(message: MsgWithdrawFromSubaccountResponse): Uint8Array {
    return MsgWithdrawFromSubaccountResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgWithdrawFromSubaccountResponse): MsgWithdrawFromSubaccountResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.sending.MsgWithdrawFromSubaccountResponse",
      value: MsgWithdrawFromSubaccountResponse.encode(message).finish()
    };
  }
};
function createBaseMsgSendFromModuleToAccountResponse(): MsgSendFromModuleToAccountResponse {
  return {};
}
export const MsgSendFromModuleToAccountResponse = {
  typeUrl: "/dydxprotocol.sending.MsgSendFromModuleToAccountResponse",
  encode(_: MsgSendFromModuleToAccountResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgSendFromModuleToAccountResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSendFromModuleToAccountResponse();
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
  fromPartial(_: Partial<MsgSendFromModuleToAccountResponse>): MsgSendFromModuleToAccountResponse {
    const message = createBaseMsgSendFromModuleToAccountResponse();
    return message;
  },
  fromAmino(_: MsgSendFromModuleToAccountResponseAmino): MsgSendFromModuleToAccountResponse {
    const message = createBaseMsgSendFromModuleToAccountResponse();
    return message;
  },
  toAmino(_: MsgSendFromModuleToAccountResponse): MsgSendFromModuleToAccountResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgSendFromModuleToAccountResponseAminoMsg): MsgSendFromModuleToAccountResponse {
    return MsgSendFromModuleToAccountResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSendFromModuleToAccountResponseProtoMsg): MsgSendFromModuleToAccountResponse {
    return MsgSendFromModuleToAccountResponse.decode(message.value);
  },
  toProto(message: MsgSendFromModuleToAccountResponse): Uint8Array {
    return MsgSendFromModuleToAccountResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgSendFromModuleToAccountResponse): MsgSendFromModuleToAccountResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.sending.MsgSendFromModuleToAccountResponse",
      value: MsgSendFromModuleToAccountResponse.encode(message).finish()
    };
  }
};