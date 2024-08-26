import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { QuotingParams, QuotingParamsSDKType, VaultParams, VaultParamsSDKType } from "./params";
import { VaultId, VaultIdSDKType } from "./vault";
import { NumShares, NumSharesSDKType } from "./share";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * MsgDepositToMegavault deposits the specified asset from the subaccount to
 * megavault.
 */

export interface MsgDepositToMegavault {
  /** The subaccount to deposit from. */
  subaccountId?: SubaccountId;
  /** Number of quote quantums to deposit. */

  quoteQuantums: Uint8Array;
}
/**
 * MsgDepositToMegavault deposits the specified asset from the subaccount to
 * megavault.
 */

export interface MsgDepositToMegavaultSDKType {
  /** The subaccount to deposit from. */
  subaccount_id?: SubaccountIdSDKType;
  /** Number of quote quantums to deposit. */

  quote_quantums: Uint8Array;
}
/** MsgDepositToMegavaultResponse is the Msg/DepositToMegavault response type. */

export interface MsgDepositToMegavaultResponse {
  /** The number of shares minted from the deposit. */
  mintedShares?: NumShares;
}
/** MsgDepositToMegavaultResponse is the Msg/DepositToMegavault response type. */

export interface MsgDepositToMegavaultResponseSDKType {
  /** The number of shares minted from the deposit. */
  minted_shares?: NumSharesSDKType;
}
/**
 * MsgUpdateDefaultQuotingParams is the Msg/UpdateDefaultQuotingParams request
 * type.
 */

export interface MsgUpdateDefaultQuotingParams {
  authority: string;
  /** The quoting parameters to update to. Every field must be set. */

  defaultQuotingParams?: QuotingParams;
}
/**
 * MsgUpdateDefaultQuotingParams is the Msg/UpdateDefaultQuotingParams request
 * type.
 */

export interface MsgUpdateDefaultQuotingParamsSDKType {
  authority: string;
  /** The quoting parameters to update to. Every field must be set. */

  default_quoting_params?: QuotingParamsSDKType;
}
/**
 * MsgUpdateDefaultQuotingParamsResponse is the Msg/UpdateDefaultQuotingParams
 * response type.
 */

export interface MsgUpdateDefaultQuotingParamsResponse {}
/**
 * MsgUpdateDefaultQuotingParamsResponse is the Msg/UpdateDefaultQuotingParams
 * response type.
 */

export interface MsgUpdateDefaultQuotingParamsResponseSDKType {}
/** MsgSetVaultParams is the Msg/SetVaultParams request type. */

export interface MsgSetVaultParams {
  authority: string;
  /** The vault to set params of. */

  vaultId?: VaultId;
  /** The parameters to set. */

  vaultParams?: VaultParams;
}
/** MsgSetVaultParams is the Msg/SetVaultParams request type. */

export interface MsgSetVaultParamsSDKType {
  authority: string;
  /** The vault to set params of. */

  vault_id?: VaultIdSDKType;
  /** The parameters to set. */

  vault_params?: VaultParamsSDKType;
}
/** MsgSetVaultParamsResponse is the Msg/SetVaultParams response type. */

export interface MsgSetVaultParamsResponse {}
/** MsgSetVaultParamsResponse is the Msg/SetVaultParams response type. */

export interface MsgSetVaultParamsResponseSDKType {}

function createBaseMsgDepositToMegavault(): MsgDepositToMegavault {
  return {
    subaccountId: undefined,
    quoteQuantums: new Uint8Array()
  };
}

export const MsgDepositToMegavault = {
  encode(message: MsgDepositToMegavault, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.subaccountId !== undefined) {
      SubaccountId.encode(message.subaccountId, writer.uint32(10).fork()).ldelim();
    }

    if (message.quoteQuantums.length !== 0) {
      writer.uint32(18).bytes(message.quoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositToMegavault {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDepositToMegavault();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountId = SubaccountId.decode(reader, reader.uint32());
          break;

        case 2:
          message.quoteQuantums = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgDepositToMegavault>): MsgDepositToMegavault {
    const message = createBaseMsgDepositToMegavault();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? SubaccountId.fromPartial(object.subaccountId) : undefined;
    message.quoteQuantums = object.quoteQuantums ?? new Uint8Array();
    return message;
  }

};

function createBaseMsgDepositToMegavaultResponse(): MsgDepositToMegavaultResponse {
  return {
    mintedShares: undefined
  };
}

export const MsgDepositToMegavaultResponse = {
  encode(message: MsgDepositToMegavaultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.mintedShares !== undefined) {
      NumShares.encode(message.mintedShares, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositToMegavaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDepositToMegavaultResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.mintedShares = NumShares.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgDepositToMegavaultResponse>): MsgDepositToMegavaultResponse {
    const message = createBaseMsgDepositToMegavaultResponse();
    message.mintedShares = object.mintedShares !== undefined && object.mintedShares !== null ? NumShares.fromPartial(object.mintedShares) : undefined;
    return message;
  }

};

function createBaseMsgUpdateDefaultQuotingParams(): MsgUpdateDefaultQuotingParams {
  return {
    authority: "",
    defaultQuotingParams: undefined
  };
}

export const MsgUpdateDefaultQuotingParams = {
  encode(message: MsgUpdateDefaultQuotingParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.defaultQuotingParams !== undefined) {
      QuotingParams.encode(message.defaultQuotingParams, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDefaultQuotingParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDefaultQuotingParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.defaultQuotingParams = QuotingParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdateDefaultQuotingParams>): MsgUpdateDefaultQuotingParams {
    const message = createBaseMsgUpdateDefaultQuotingParams();
    message.authority = object.authority ?? "";
    message.defaultQuotingParams = object.defaultQuotingParams !== undefined && object.defaultQuotingParams !== null ? QuotingParams.fromPartial(object.defaultQuotingParams) : undefined;
    return message;
  }

};

function createBaseMsgUpdateDefaultQuotingParamsResponse(): MsgUpdateDefaultQuotingParamsResponse {
  return {};
}

export const MsgUpdateDefaultQuotingParamsResponse = {
  encode(_: MsgUpdateDefaultQuotingParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDefaultQuotingParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDefaultQuotingParamsResponse();

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

  fromPartial(_: DeepPartial<MsgUpdateDefaultQuotingParamsResponse>): MsgUpdateDefaultQuotingParamsResponse {
    const message = createBaseMsgUpdateDefaultQuotingParamsResponse();
    return message;
  }

};

function createBaseMsgSetVaultParams(): MsgSetVaultParams {
  return {
    authority: "",
    vaultId: undefined,
    vaultParams: undefined
  };
}

export const MsgSetVaultParams = {
  encode(message: MsgSetVaultParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.vaultId !== undefined) {
      VaultId.encode(message.vaultId, writer.uint32(18).fork()).ldelim();
    }

    if (message.vaultParams !== undefined) {
      VaultParams.encode(message.vaultParams, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetVaultParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetVaultParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.vaultId = VaultId.decode(reader, reader.uint32());
          break;

        case 3:
          message.vaultParams = VaultParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetVaultParams>): MsgSetVaultParams {
    const message = createBaseMsgSetVaultParams();
    message.authority = object.authority ?? "";
    message.vaultId = object.vaultId !== undefined && object.vaultId !== null ? VaultId.fromPartial(object.vaultId) : undefined;
    message.vaultParams = object.vaultParams !== undefined && object.vaultParams !== null ? VaultParams.fromPartial(object.vaultParams) : undefined;
    return message;
  }

};

function createBaseMsgSetVaultParamsResponse(): MsgSetVaultParamsResponse {
  return {};
}

export const MsgSetVaultParamsResponse = {
  encode(_: MsgSetVaultParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetVaultParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetVaultParamsResponse();

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

  fromPartial(_: DeepPartial<MsgSetVaultParamsResponse>): MsgSetVaultParamsResponse {
    const message = createBaseMsgSetVaultParamsResponse();
    return message;
  }

};