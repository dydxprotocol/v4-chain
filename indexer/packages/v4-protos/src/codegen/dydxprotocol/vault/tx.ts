import { VaultId, VaultIdSDKType, NumShares, NumSharesSDKType } from "./vault";
import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { Params, ParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * MsgDepositToVault deposits the specified asset from the subaccout to the
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
 * MsgDepositToVault deposits the specified asset from the subaccout to the
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
 * MsgWithdrawFromVault withdraws the specified amount of shares from the vault
 * to the subaccount.
 */

export interface MsgWithdrawFromVault {
  /** The vault to withdraw from. */
  vaultId?: VaultId;
  /** The subaccount to withdraw to. */

  subaccountId?: SubaccountId;
  /**
   * The target amount of quote quantums to withdraw.
   * The target amount should account for slippage. For example:
   *    * user equity = $1020
   *    * target amount = $1000
   *    * slippage = $50
   * then, the amount withdrawn is $970 (= 1020 - (1000 + 50 - 1020))
   */

  targetQuoteQuantums: Uint8Array;
}
/**
 * MsgWithdrawFromVault withdraws the specified amount of shares from the vault
 * to the subaccount.
 */

export interface MsgWithdrawFromVaultSDKType {
  /** The vault to withdraw from. */
  vault_id?: VaultIdSDKType;
  /** The subaccount to withdraw to. */

  subaccount_id?: SubaccountIdSDKType;
  /**
   * The target amount of quote quantums to withdraw.
   * The target amount should account for slippage. For example:
   *    * user equity = $1020
   *    * target amount = $1000
   *    * slippage = $50
   * then, the amount withdrawn is $970 (= 1020 - (1000 + 50 - 1020))
   */

  target_quote_quantums: Uint8Array;
}
/** MsgWithdrawFromVaultResponse is the Msg/WithdrawFromVault response type. */

export interface MsgWithdrawFromVaultResponse {
  /** Number of shares that have been redeemed as part of the withdrawal. */
  redeemedShares?: NumShares;
  /** Number of shares remaining after the withdrawal. */

  remainingShares?: NumShares;
  /** Number of shares for the vault after the withdrawal. */

  totalShares?: NumShares;
}
/** MsgWithdrawFromVaultResponse is the Msg/WithdrawFromVault response type. */

export interface MsgWithdrawFromVaultResponseSDKType {
  /** Number of shares that have been redeemed as part of the withdrawal. */
  redeemed_shares?: NumSharesSDKType;
  /** Number of shares remaining after the withdrawal. */

  remaining_shares?: NumSharesSDKType;
  /** Number of shares for the vault after the withdrawal. */

  total_shares?: NumSharesSDKType;
}
/** MsgUpdateParams is the Msg/UpdateParams request type. */

export interface MsgUpdateParams {
  authority: string;
  /** The parameters to update. Each field must be set. */

  params?: Params;
}
/** MsgUpdateParams is the Msg/UpdateParams request type. */

export interface MsgUpdateParamsSDKType {
  authority: string;
  /** The parameters to update. Each field must be set. */

  params?: ParamsSDKType;
}
/** MsgUpdateParamsResponse is the Msg/UpdateParams response type. */

export interface MsgUpdateParamsResponse {}
/** MsgUpdateParamsResponse is the Msg/UpdateParams response type. */

export interface MsgUpdateParamsResponseSDKType {}

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

function createBaseMsgWithdrawFromVault(): MsgWithdrawFromVault {
  return {
    vaultId: undefined,
    subaccountId: undefined,
    targetQuoteQuantums: new Uint8Array()
  };
}

export const MsgWithdrawFromVault = {
  encode(message: MsgWithdrawFromVault, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.vaultId !== undefined) {
      VaultId.encode(message.vaultId, writer.uint32(10).fork()).ldelim();
    }

    if (message.subaccountId !== undefined) {
      SubaccountId.encode(message.subaccountId, writer.uint32(18).fork()).ldelim();
    }

    if (message.targetQuoteQuantums.length !== 0) {
      writer.uint32(26).bytes(message.targetQuoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawFromVault {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFromVault();

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
          message.targetQuoteQuantums = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgWithdrawFromVault>): MsgWithdrawFromVault {
    const message = createBaseMsgWithdrawFromVault();
    message.vaultId = object.vaultId !== undefined && object.vaultId !== null ? VaultId.fromPartial(object.vaultId) : undefined;
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? SubaccountId.fromPartial(object.subaccountId) : undefined;
    message.targetQuoteQuantums = object.targetQuoteQuantums ?? new Uint8Array();
    return message;
  }

};

function createBaseMsgWithdrawFromVaultResponse(): MsgWithdrawFromVaultResponse {
  return {
    redeemedShares: undefined,
    remainingShares: undefined,
    totalShares: undefined
  };
}

export const MsgWithdrawFromVaultResponse = {
  encode(message: MsgWithdrawFromVaultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.redeemedShares !== undefined) {
      NumShares.encode(message.redeemedShares, writer.uint32(10).fork()).ldelim();
    }

    if (message.remainingShares !== undefined) {
      NumShares.encode(message.remainingShares, writer.uint32(18).fork()).ldelim();
    }

    if (message.totalShares !== undefined) {
      NumShares.encode(message.totalShares, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawFromVaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFromVaultResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.redeemedShares = NumShares.decode(reader, reader.uint32());
          break;

        case 2:
          message.remainingShares = NumShares.decode(reader, reader.uint32());
          break;

        case 3:
          message.totalShares = NumShares.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgWithdrawFromVaultResponse>): MsgWithdrawFromVaultResponse {
    const message = createBaseMsgWithdrawFromVaultResponse();
    message.redeemedShares = object.redeemedShares !== undefined && object.redeemedShares !== null ? NumShares.fromPartial(object.redeemedShares) : undefined;
    message.remainingShares = object.remainingShares !== undefined && object.remainingShares !== null ? NumShares.fromPartial(object.remainingShares) : undefined;
    message.totalShares = object.totalShares !== undefined && object.totalShares !== null ? NumShares.fromPartial(object.totalShares) : undefined;
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

function createBaseMsgUpdateParamsResponse(): MsgUpdateParamsResponse {
  return {};
}

export const MsgUpdateParamsResponse = {
  encode(_: MsgUpdateParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateParamsResponse();

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

  fromPartial(_: DeepPartial<MsgUpdateParamsResponse>): MsgUpdateParamsResponse {
    const message = createBaseMsgUpdateParamsResponse();
    return message;
  }

};