import { VaultId, VaultIdSDKType } from "./vault";
import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { QuotingParams, QuotingParamsSDKType, Params, ParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * MsgDepositToVault deposits the specified asset from the subaccount to the
 * vault.
 */

export interface MsgDepositToVault {
  /** The vault to deposit into. */
  vaultId?: VaultId;
  /** The subaccount to deposit from. */

  subaccountId?: SubaccountId;
  /** Number of quote quantums to deposit. */

  quoteQuantums: Uint8Array;
}
/**
 * MsgDepositToVault deposits the specified asset from the subaccount to the
 * vault.
 */

export interface MsgDepositToVaultSDKType {
  /** The vault to deposit into. */
  vault_id?: VaultIdSDKType;
  /** The subaccount to deposit from. */

  subaccount_id?: SubaccountIdSDKType;
  /** Number of quote quantums to deposit. */

  quote_quantums: Uint8Array;
}
/** MsgDepositToVaultResponse is the Msg/DepositToVault response type. */

export interface MsgDepositToVaultResponse {}
/** MsgDepositToVaultResponse is the Msg/DepositToVault response type. */

export interface MsgDepositToVaultResponseSDKType {}
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
/** MsgSetVaultQuotingParams is the Msg/SetVaultQuotingParams request type. */

export interface MsgSetVaultQuotingParams {
  authority: string;
  /** The vault to set quoting params of. */

  vaultId?: VaultId;
  /** The quoting parameters to set. Each field must be set. */

  quotingParams?: QuotingParams;
}
/** MsgSetVaultQuotingParams is the Msg/SetVaultQuotingParams request type. */

export interface MsgSetVaultQuotingParamsSDKType {
  authority: string;
  /** The vault to set quoting params of. */

  vault_id?: VaultIdSDKType;
  /** The quoting parameters to set. Each field must be set. */

  quoting_params?: QuotingParamsSDKType;
}
/**
 * MsgSetVaultQuotingParamsResponse is the Msg/SetVaultQuotingParams response
 * type.
 */

export interface MsgSetVaultQuotingParamsResponse {}
/**
 * MsgSetVaultQuotingParamsResponse is the Msg/SetVaultQuotingParams response
 * type.
 */

export interface MsgSetVaultQuotingParamsResponseSDKType {}
/**
 * MsgUpdateParams is the Msg/UpdateParams request type.
 * Deprecated since v6.x as is replaced by MsgUpdateDefaultQuotingParams.
 */

export interface MsgUpdateParams {
  authority: string;
  /** The parameters to update. Each field must be set. */

  params?: Params;
}
/**
 * MsgUpdateParams is the Msg/UpdateParams request type.
 * Deprecated since v6.x as is replaced by MsgUpdateDefaultQuotingParams.
 */

export interface MsgUpdateParamsSDKType {
  authority: string;
  /** The parameters to update. Each field must be set. */

  params?: ParamsSDKType;
}

function createBaseMsgDepositToVault(): MsgDepositToVault {
  return {
    vaultId: undefined,
    subaccountId: undefined,
    quoteQuantums: new Uint8Array()
  };
}

export const MsgDepositToVault = {
  encode(message: MsgDepositToVault, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.vaultId !== undefined) {
      VaultId.encode(message.vaultId, writer.uint32(10).fork()).ldelim();
    }

    if (message.subaccountId !== undefined) {
      SubaccountId.encode(message.subaccountId, writer.uint32(18).fork()).ldelim();
    }

    if (message.quoteQuantums.length !== 0) {
      writer.uint32(26).bytes(message.quoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositToVault {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDepositToVault();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.vaultId = VaultId.decode(reader, reader.uint32());
          break;

        case 2:
          message.subaccountId = SubaccountId.decode(reader, reader.uint32());
          break;

        case 3:
          message.quoteQuantums = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgDepositToVault>): MsgDepositToVault {
    const message = createBaseMsgDepositToVault();
    message.vaultId = object.vaultId !== undefined && object.vaultId !== null ? VaultId.fromPartial(object.vaultId) : undefined;
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? SubaccountId.fromPartial(object.subaccountId) : undefined;
    message.quoteQuantums = object.quoteQuantums ?? new Uint8Array();
    return message;
  }

};

function createBaseMsgDepositToVaultResponse(): MsgDepositToVaultResponse {
  return {};
}

export const MsgDepositToVaultResponse = {
  encode(_: MsgDepositToVaultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDepositToVaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDepositToVaultResponse();

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

  fromPartial(_: DeepPartial<MsgDepositToVaultResponse>): MsgDepositToVaultResponse {
    const message = createBaseMsgDepositToVaultResponse();
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

function createBaseMsgSetVaultQuotingParams(): MsgSetVaultQuotingParams {
  return {
    authority: "",
    vaultId: undefined,
    quotingParams: undefined
  };
}

export const MsgSetVaultQuotingParams = {
  encode(message: MsgSetVaultQuotingParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.vaultId !== undefined) {
      VaultId.encode(message.vaultId, writer.uint32(18).fork()).ldelim();
    }

    if (message.quotingParams !== undefined) {
      QuotingParams.encode(message.quotingParams, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetVaultQuotingParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetVaultQuotingParams();

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
          message.quotingParams = QuotingParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetVaultQuotingParams>): MsgSetVaultQuotingParams {
    const message = createBaseMsgSetVaultQuotingParams();
    message.authority = object.authority ?? "";
    message.vaultId = object.vaultId !== undefined && object.vaultId !== null ? VaultId.fromPartial(object.vaultId) : undefined;
    message.quotingParams = object.quotingParams !== undefined && object.quotingParams !== null ? QuotingParams.fromPartial(object.quotingParams) : undefined;
    return message;
  }

};

function createBaseMsgSetVaultQuotingParamsResponse(): MsgSetVaultQuotingParamsResponse {
  return {};
}

export const MsgSetVaultQuotingParamsResponse = {
  encode(_: MsgSetVaultQuotingParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetVaultQuotingParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetVaultQuotingParamsResponse();

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

  fromPartial(_: DeepPartial<MsgSetVaultQuotingParamsResponse>): MsgSetVaultQuotingParamsResponse {
    const message = createBaseMsgSetVaultQuotingParamsResponse();
    return message;
  }

};

function createBaseMsgUpdateParams(): MsgUpdateParams {
  return {
    authority: "",
    params: undefined
  };
}

export const MsgUpdateParams = {
  encode(message: MsgUpdateParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params = Params.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdateParams>): MsgUpdateParams {
    const message = createBaseMsgUpdateParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    return message;
  }

};