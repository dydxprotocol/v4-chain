import { SubaccountId, SubaccountIdAmino, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { Coin, CoinAmino, CoinSDKType } from "../../cosmos/base/v1beta1/coin";
import { BinaryReader, BinaryWriter } from "../../binary";
/** Transfer represents a single transfer between two subaccounts. */
export interface Transfer {
  /** The sender subaccount ID. */
  sender: SubaccountId;
  /** The recipient subaccount ID. */
  recipient: SubaccountId;
  /** Id of the asset to transfer. */
  assetId: number;
  /** The amount of asset to transfer */
  amount: bigint;
}
export interface TransferProtoMsg {
  typeUrl: "/dydxprotocol.sending.Transfer";
  value: Uint8Array;
}
/** Transfer represents a single transfer between two subaccounts. */
export interface TransferAmino {
  /** The sender subaccount ID. */
  sender?: SubaccountIdAmino;
  /** The recipient subaccount ID. */
  recipient?: SubaccountIdAmino;
  /** Id of the asset to transfer. */
  asset_id?: number;
  /** The amount of asset to transfer */
  amount?: string;
}
export interface TransferAminoMsg {
  type: "/dydxprotocol.sending.Transfer";
  value: TransferAmino;
}
/** Transfer represents a single transfer between two subaccounts. */
export interface TransferSDKType {
  sender: SubaccountIdSDKType;
  recipient: SubaccountIdSDKType;
  asset_id: number;
  amount: bigint;
}
/**
 * MsgDepositToSubaccount represents a single transfer from an `x/bank`
 * account to an `x/subaccounts` subaccount.
 */
export interface MsgDepositToSubaccount {
  /** The sender wallet address. */
  sender: string;
  /** The recipient subaccount ID. */
  recipient: SubaccountId;
  /** Id of the asset to transfer. */
  assetId: number;
  /** The number of quantums of asset to transfer. */
  quantums: bigint;
}
export interface MsgDepositToSubaccountProtoMsg {
  typeUrl: "/dydxprotocol.sending.MsgDepositToSubaccount";
  value: Uint8Array;
}
/**
 * MsgDepositToSubaccount represents a single transfer from an `x/bank`
 * account to an `x/subaccounts` subaccount.
 */
export interface MsgDepositToSubaccountAmino {
  /** The sender wallet address. */
  sender?: string;
  /** The recipient subaccount ID. */
  recipient?: SubaccountIdAmino;
  /** Id of the asset to transfer. */
  asset_id?: number;
  /** The number of quantums of asset to transfer. */
  quantums?: string;
}
export interface MsgDepositToSubaccountAminoMsg {
  type: "/dydxprotocol.sending.MsgDepositToSubaccount";
  value: MsgDepositToSubaccountAmino;
}
/**
 * MsgDepositToSubaccount represents a single transfer from an `x/bank`
 * account to an `x/subaccounts` subaccount.
 */
export interface MsgDepositToSubaccountSDKType {
  sender: string;
  recipient: SubaccountIdSDKType;
  asset_id: number;
  quantums: bigint;
}
/**
 * MsgWithdrawFromSubaccount represents a single transfer from an
 * `x/subaccounts` subaccount to an `x/bank` account.
 */
export interface MsgWithdrawFromSubaccount {
  /** The sender subaccount ID. */
  sender: SubaccountId;
  /** The recipient wallet address. */
  recipient: string;
  /** Id of the asset to transfer. */
  assetId: number;
  /** The number of quantums of asset to transfer. */
  quantums: bigint;
}
export interface MsgWithdrawFromSubaccountProtoMsg {
  typeUrl: "/dydxprotocol.sending.MsgWithdrawFromSubaccount";
  value: Uint8Array;
}
/**
 * MsgWithdrawFromSubaccount represents a single transfer from an
 * `x/subaccounts` subaccount to an `x/bank` account.
 */
export interface MsgWithdrawFromSubaccountAmino {
  /** The sender subaccount ID. */
  sender?: SubaccountIdAmino;
  /** The recipient wallet address. */
  recipient?: string;
  /** Id of the asset to transfer. */
  asset_id?: number;
  /** The number of quantums of asset to transfer. */
  quantums?: string;
}
export interface MsgWithdrawFromSubaccountAminoMsg {
  type: "/dydxprotocol.sending.MsgWithdrawFromSubaccount";
  value: MsgWithdrawFromSubaccountAmino;
}
/**
 * MsgWithdrawFromSubaccount represents a single transfer from an
 * `x/subaccounts` subaccount to an `x/bank` account.
 */
export interface MsgWithdrawFromSubaccountSDKType {
  sender: SubaccountIdSDKType;
  recipient: string;
  asset_id: number;
  quantums: bigint;
}
/**
 * MsgSendFromModuleToAccount represents a single transfer from a module
 * to an `x/bank` account (can be either a module account address or a user
 * account address).
 * Should only be executed by governance.
 */
export interface MsgSendFromModuleToAccount {
  authority: string;
  /** The sender module name. */
  senderModuleName: string;
  /**
   * The recipient account address (can be either a module account address
   * or a user account address).
   */
  recipient: string;
  /** The coin to transfer, which specifies both denom and amount. */
  coin: Coin;
}
export interface MsgSendFromModuleToAccountProtoMsg {
  typeUrl: "/dydxprotocol.sending.MsgSendFromModuleToAccount";
  value: Uint8Array;
}
/**
 * MsgSendFromModuleToAccount represents a single transfer from a module
 * to an `x/bank` account (can be either a module account address or a user
 * account address).
 * Should only be executed by governance.
 */
export interface MsgSendFromModuleToAccountAmino {
  authority?: string;
  /** The sender module name. */
  sender_module_name?: string;
  /**
   * The recipient account address (can be either a module account address
   * or a user account address).
   */
  recipient?: string;
  /** The coin to transfer, which specifies both denom and amount. */
  coin?: CoinAmino;
}
export interface MsgSendFromModuleToAccountAminoMsg {
  type: "/dydxprotocol.sending.MsgSendFromModuleToAccount";
  value: MsgSendFromModuleToAccountAmino;
}
/**
 * MsgSendFromModuleToAccount represents a single transfer from a module
 * to an `x/bank` account (can be either a module account address or a user
 * account address).
 * Should only be executed by governance.
 */
export interface MsgSendFromModuleToAccountSDKType {
  authority: string;
  sender_module_name: string;
  recipient: string;
  coin: CoinSDKType;
}
function createBaseTransfer(): Transfer {
  return {
    sender: SubaccountId.fromPartial({}),
    recipient: SubaccountId.fromPartial({}),
    assetId: 0,
    amount: BigInt(0)
  };
}
export const Transfer = {
  typeUrl: "/dydxprotocol.sending.Transfer",
  encode(message: Transfer, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.sender !== undefined) {
      SubaccountId.encode(message.sender, writer.uint32(10).fork()).ldelim();
    }
    if (message.recipient !== undefined) {
      SubaccountId.encode(message.recipient, writer.uint32(18).fork()).ldelim();
    }
    if (message.assetId !== 0) {
      writer.uint32(24).uint32(message.assetId);
    }
    if (message.amount !== BigInt(0)) {
      writer.uint32(32).uint64(message.amount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Transfer {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTransfer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sender = SubaccountId.decode(reader, reader.uint32());
          break;
        case 2:
          message.recipient = SubaccountId.decode(reader, reader.uint32());
          break;
        case 3:
          message.assetId = reader.uint32();
          break;
        case 4:
          message.amount = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Transfer>): Transfer {
    const message = createBaseTransfer();
    message.sender = object.sender !== undefined && object.sender !== null ? SubaccountId.fromPartial(object.sender) : undefined;
    message.recipient = object.recipient !== undefined && object.recipient !== null ? SubaccountId.fromPartial(object.recipient) : undefined;
    message.assetId = object.assetId ?? 0;
    message.amount = object.amount !== undefined && object.amount !== null ? BigInt(object.amount.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: TransferAmino): Transfer {
    const message = createBaseTransfer();
    if (object.sender !== undefined && object.sender !== null) {
      message.sender = SubaccountId.fromAmino(object.sender);
    }
    if (object.recipient !== undefined && object.recipient !== null) {
      message.recipient = SubaccountId.fromAmino(object.recipient);
    }
    if (object.asset_id !== undefined && object.asset_id !== null) {
      message.assetId = object.asset_id;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = BigInt(object.amount);
    }
    return message;
  },
  toAmino(message: Transfer): TransferAmino {
    const obj: any = {};
    obj.sender = message.sender ? SubaccountId.toAmino(message.sender) : undefined;
    obj.recipient = message.recipient ? SubaccountId.toAmino(message.recipient) : undefined;
    obj.asset_id = message.assetId;
    obj.amount = message.amount ? message.amount.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: TransferAminoMsg): Transfer {
    return Transfer.fromAmino(object.value);
  },
  fromProtoMsg(message: TransferProtoMsg): Transfer {
    return Transfer.decode(message.value);
  },
  toProto(message: Transfer): Uint8Array {
    return Transfer.encode(message).finish();
  },
  toProtoMsg(message: Transfer): TransferProtoMsg {
    return {
      typeUrl: "/dydxprotocol.sending.Transfer",
      value: Transfer.encode(message).finish()
    };
  }
};
function createBaseMsgDepositToSubaccount(): MsgDepositToSubaccount {
  return {
    sender: "",
    recipient: SubaccountId.fromPartial({}),
    assetId: 0,
    quantums: BigInt(0)
  };
}
export const MsgDepositToSubaccount = {
  typeUrl: "/dydxprotocol.sending.MsgDepositToSubaccount",
  encode(message: MsgDepositToSubaccount, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.sender !== "") {
      writer.uint32(10).string(message.sender);
    }
    if (message.recipient !== undefined) {
      SubaccountId.encode(message.recipient, writer.uint32(18).fork()).ldelim();
    }
    if (message.assetId !== 0) {
      writer.uint32(24).uint32(message.assetId);
    }
    if (message.quantums !== BigInt(0)) {
      writer.uint32(32).uint64(message.quantums);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgDepositToSubaccount {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDepositToSubaccount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sender = reader.string();
          break;
        case 2:
          message.recipient = SubaccountId.decode(reader, reader.uint32());
          break;
        case 3:
          message.assetId = reader.uint32();
          break;
        case 4:
          message.quantums = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgDepositToSubaccount>): MsgDepositToSubaccount {
    const message = createBaseMsgDepositToSubaccount();
    message.sender = object.sender ?? "";
    message.recipient = object.recipient !== undefined && object.recipient !== null ? SubaccountId.fromPartial(object.recipient) : undefined;
    message.assetId = object.assetId ?? 0;
    message.quantums = object.quantums !== undefined && object.quantums !== null ? BigInt(object.quantums.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MsgDepositToSubaccountAmino): MsgDepositToSubaccount {
    const message = createBaseMsgDepositToSubaccount();
    if (object.sender !== undefined && object.sender !== null) {
      message.sender = object.sender;
    }
    if (object.recipient !== undefined && object.recipient !== null) {
      message.recipient = SubaccountId.fromAmino(object.recipient);
    }
    if (object.asset_id !== undefined && object.asset_id !== null) {
      message.assetId = object.asset_id;
    }
    if (object.quantums !== undefined && object.quantums !== null) {
      message.quantums = BigInt(object.quantums);
    }
    return message;
  },
  toAmino(message: MsgDepositToSubaccount): MsgDepositToSubaccountAmino {
    const obj: any = {};
    obj.sender = message.sender;
    obj.recipient = message.recipient ? SubaccountId.toAmino(message.recipient) : undefined;
    obj.asset_id = message.assetId;
    obj.quantums = message.quantums ? message.quantums.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgDepositToSubaccountAminoMsg): MsgDepositToSubaccount {
    return MsgDepositToSubaccount.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgDepositToSubaccountProtoMsg): MsgDepositToSubaccount {
    return MsgDepositToSubaccount.decode(message.value);
  },
  toProto(message: MsgDepositToSubaccount): Uint8Array {
    return MsgDepositToSubaccount.encode(message).finish();
  },
  toProtoMsg(message: MsgDepositToSubaccount): MsgDepositToSubaccountProtoMsg {
    return {
      typeUrl: "/dydxprotocol.sending.MsgDepositToSubaccount",
      value: MsgDepositToSubaccount.encode(message).finish()
    };
  }
};
function createBaseMsgWithdrawFromSubaccount(): MsgWithdrawFromSubaccount {
  return {
    sender: SubaccountId.fromPartial({}),
    recipient: "",
    assetId: 0,
    quantums: BigInt(0)
  };
}
export const MsgWithdrawFromSubaccount = {
  typeUrl: "/dydxprotocol.sending.MsgWithdrawFromSubaccount",
  encode(message: MsgWithdrawFromSubaccount, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.sender !== undefined) {
      SubaccountId.encode(message.sender, writer.uint32(18).fork()).ldelim();
    }
    if (message.recipient !== "") {
      writer.uint32(10).string(message.recipient);
    }
    if (message.assetId !== 0) {
      writer.uint32(24).uint32(message.assetId);
    }
    if (message.quantums !== BigInt(0)) {
      writer.uint32(32).uint64(message.quantums);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgWithdrawFromSubaccount {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFromSubaccount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.sender = SubaccountId.decode(reader, reader.uint32());
          break;
        case 1:
          message.recipient = reader.string();
          break;
        case 3:
          message.assetId = reader.uint32();
          break;
        case 4:
          message.quantums = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgWithdrawFromSubaccount>): MsgWithdrawFromSubaccount {
    const message = createBaseMsgWithdrawFromSubaccount();
    message.sender = object.sender !== undefined && object.sender !== null ? SubaccountId.fromPartial(object.sender) : undefined;
    message.recipient = object.recipient ?? "";
    message.assetId = object.assetId ?? 0;
    message.quantums = object.quantums !== undefined && object.quantums !== null ? BigInt(object.quantums.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MsgWithdrawFromSubaccountAmino): MsgWithdrawFromSubaccount {
    const message = createBaseMsgWithdrawFromSubaccount();
    if (object.sender !== undefined && object.sender !== null) {
      message.sender = SubaccountId.fromAmino(object.sender);
    }
    if (object.recipient !== undefined && object.recipient !== null) {
      message.recipient = object.recipient;
    }
    if (object.asset_id !== undefined && object.asset_id !== null) {
      message.assetId = object.asset_id;
    }
    if (object.quantums !== undefined && object.quantums !== null) {
      message.quantums = BigInt(object.quantums);
    }
    return message;
  },
  toAmino(message: MsgWithdrawFromSubaccount): MsgWithdrawFromSubaccountAmino {
    const obj: any = {};
    obj.sender = message.sender ? SubaccountId.toAmino(message.sender) : undefined;
    obj.recipient = message.recipient;
    obj.asset_id = message.assetId;
    obj.quantums = message.quantums ? message.quantums.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgWithdrawFromSubaccountAminoMsg): MsgWithdrawFromSubaccount {
    return MsgWithdrawFromSubaccount.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgWithdrawFromSubaccountProtoMsg): MsgWithdrawFromSubaccount {
    return MsgWithdrawFromSubaccount.decode(message.value);
  },
  toProto(message: MsgWithdrawFromSubaccount): Uint8Array {
    return MsgWithdrawFromSubaccount.encode(message).finish();
  },
  toProtoMsg(message: MsgWithdrawFromSubaccount): MsgWithdrawFromSubaccountProtoMsg {
    return {
      typeUrl: "/dydxprotocol.sending.MsgWithdrawFromSubaccount",
      value: MsgWithdrawFromSubaccount.encode(message).finish()
    };
  }
};
function createBaseMsgSendFromModuleToAccount(): MsgSendFromModuleToAccount {
  return {
    authority: "",
    senderModuleName: "",
    recipient: "",
    coin: Coin.fromPartial({})
  };
}
export const MsgSendFromModuleToAccount = {
  typeUrl: "/dydxprotocol.sending.MsgSendFromModuleToAccount",
  encode(message: MsgSendFromModuleToAccount, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.senderModuleName !== "") {
      writer.uint32(18).string(message.senderModuleName);
    }
    if (message.recipient !== "") {
      writer.uint32(26).string(message.recipient);
    }
    if (message.coin !== undefined) {
      Coin.encode(message.coin, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgSendFromModuleToAccount {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSendFromModuleToAccount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.senderModuleName = reader.string();
          break;
        case 3:
          message.recipient = reader.string();
          break;
        case 4:
          message.coin = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgSendFromModuleToAccount>): MsgSendFromModuleToAccount {
    const message = createBaseMsgSendFromModuleToAccount();
    message.authority = object.authority ?? "";
    message.senderModuleName = object.senderModuleName ?? "";
    message.recipient = object.recipient ?? "";
    message.coin = object.coin !== undefined && object.coin !== null ? Coin.fromPartial(object.coin) : undefined;
    return message;
  },
  fromAmino(object: MsgSendFromModuleToAccountAmino): MsgSendFromModuleToAccount {
    const message = createBaseMsgSendFromModuleToAccount();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.sender_module_name !== undefined && object.sender_module_name !== null) {
      message.senderModuleName = object.sender_module_name;
    }
    if (object.recipient !== undefined && object.recipient !== null) {
      message.recipient = object.recipient;
    }
    if (object.coin !== undefined && object.coin !== null) {
      message.coin = Coin.fromAmino(object.coin);
    }
    return message;
  },
  toAmino(message: MsgSendFromModuleToAccount): MsgSendFromModuleToAccountAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.sender_module_name = message.senderModuleName;
    obj.recipient = message.recipient;
    obj.coin = message.coin ? Coin.toAmino(message.coin) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgSendFromModuleToAccountAminoMsg): MsgSendFromModuleToAccount {
    return MsgSendFromModuleToAccount.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSendFromModuleToAccountProtoMsg): MsgSendFromModuleToAccount {
    return MsgSendFromModuleToAccount.decode(message.value);
  },
  toProto(message: MsgSendFromModuleToAccount): Uint8Array {
    return MsgSendFromModuleToAccount.encode(message).finish();
  },
  toProtoMsg(message: MsgSendFromModuleToAccount): MsgSendFromModuleToAccountProtoMsg {
    return {
      typeUrl: "/dydxprotocol.sending.MsgSendFromModuleToAccount",
      value: MsgSendFromModuleToAccount.encode(message).finish()
    };
  }
};