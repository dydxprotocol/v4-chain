import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { NumShares, NumSharesSDKType } from "./share";
import { QuotingParams, QuotingParamsSDKType, OperatorParams, OperatorParamsSDKType, VaultParams, VaultParamsSDKType, Params, ParamsSDKType } from "./params";
import { VaultId, VaultIdSDKType } from "./vault";
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
 * MsgWithdrawFromMegavault withdraws the specified shares from megavault to
 * a subaccount.
 */

export interface MsgWithdrawFromMegavault {
  /** The subaccount to withdraw to. */
  subaccountId?: SubaccountId;
  /** Number of shares to withdraw. */

  shares?: NumShares;
  /**
   * The minimum number of quote quantums above shares should redeem, i.e.
   * transaction fails if above shares redeem less than min_quote_quantums.
   */

  minQuoteQuantums: Uint8Array;
}
/**
 * MsgWithdrawFromMegavault withdraws the specified shares from megavault to
 * a subaccount.
 */

export interface MsgWithdrawFromMegavaultSDKType {
  /** The subaccount to withdraw to. */
  subaccount_id?: SubaccountIdSDKType;
  /** Number of shares to withdraw. */

  shares?: NumSharesSDKType;
  /**
   * The minimum number of quote quantums above shares should redeem, i.e.
   * transaction fails if above shares redeem less than min_quote_quantums.
   */

  min_quote_quantums: Uint8Array;
}
/**
 * MsgWithdrawFromMegavaultResponse is the Msg/WithdrawFromMegavault response
 * type.
 */

export interface MsgWithdrawFromMegavaultResponse {
  /** The number of quote quantums redeemed from the withdrawal. */
  quoteQuantums: Uint8Array;
}
/**
 * MsgWithdrawFromMegavaultResponse is the Msg/WithdrawFromMegavault response
 * type.
 */

export interface MsgWithdrawFromMegavaultResponseSDKType {
  /** The number of quote quantums redeemed from the withdrawal. */
  quote_quantums: Uint8Array;
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
/** MsgUnlockShares is the Msg/UnlockShares request type. */

export interface MsgUnlockShares {
  authority: string;
  /** Address of the owner to unlock shares of. */

  ownerAddress: string;
}
/** MsgUnlockShares is the Msg/UnlockShares request type. */

export interface MsgUnlockSharesSDKType {
  authority: string;
  /** Address of the owner to unlock shares of. */

  owner_address: string;
}
/** MsgUnlockSharesResponse is the Msg/UnlockShares response type. */

export interface MsgUnlockSharesResponse {
  /** The number of shares unlocked. */
  unlockedShares?: NumShares;
}
/** MsgUnlockSharesResponse is the Msg/UnlockShares response type. */

export interface MsgUnlockSharesResponseSDKType {
  /** The number of shares unlocked. */
  unlocked_shares?: NumSharesSDKType;
}
/** MsgUpdateOperatorParams is the Msg/UpdateOperatorParams request type. */

export interface MsgUpdateOperatorParams {
  authority: string;
  /** Operator parameters to set. */

  params?: OperatorParams;
}
/** MsgUpdateOperatorParams is the Msg/UpdateOperatorParams request type. */

export interface MsgUpdateOperatorParamsSDKType {
  authority: string;
  /** Operator parameters to set. */

  params?: OperatorParamsSDKType;
}
/** MsgUpdateVaultParamsResponse is the Msg/UpdateOperatorParams response type. */

export interface MsgUpdateOperatorParamsResponse {}
/** MsgUpdateVaultParamsResponse is the Msg/UpdateOperatorParams response type. */

export interface MsgUpdateOperatorParamsResponseSDKType {}
/** MsgAllocateToVault is the Msg/AllocateToVault request type. */

export interface MsgAllocateToVault {
  authority: string;
  /** The vault to allocate to. */

  vaultId?: VaultId;
  /** Number of quote quantums to allocate. */

  quoteQuantums: Uint8Array;
}
/** MsgAllocateToVault is the Msg/AllocateToVault request type. */

export interface MsgAllocateToVaultSDKType {
  authority: string;
  /** The vault to allocate to. */

  vault_id?: VaultIdSDKType;
  /** Number of quote quantums to allocate. */

  quote_quantums: Uint8Array;
}
/** MsgAllocateToVaultResponse is the Msg/AllocateToVault response type. */

export interface MsgAllocateToVaultResponse {}
/** MsgAllocateToVaultResponse is the Msg/AllocateToVault response type. */

export interface MsgAllocateToVaultResponseSDKType {}
/** MsgRetrieveFromVault is the Msg/RetrieveFromVault request type. */

export interface MsgRetrieveFromVault {
  authority: string;
  /** The vault to retrieve from. */

  vaultId?: VaultId;
  /** Number of quote quantums to retrieve. */

  quoteQuantums: Uint8Array;
}
/** MsgRetrieveFromVault is the Msg/RetrieveFromVault request type. */

export interface MsgRetrieveFromVaultSDKType {
  authority: string;
  /** The vault to retrieve from. */

  vault_id?: VaultIdSDKType;
  /** Number of quote quantums to retrieve. */

  quote_quantums: Uint8Array;
}
/** MsgRetrieveFromVaultResponse is the Msg/RetrieveFromVault response type. */

export interface MsgRetrieveFromVaultResponse {}
/** MsgRetrieveFromVaultResponse is the Msg/RetrieveFromVault response type. */

export interface MsgRetrieveFromVaultResponseSDKType {}
/**
 * MsgUpdateParams is the Msg/UpdateParams request type.
 * Deprecated since v6.x as is replaced by MsgUpdateDefaultQuotingParams.
 */

/** @deprecated */

export interface MsgUpdateParams {
  authority: string;
  /** The parameters to update. Each field must be set. */

  params?: Params;
}
/**
 * MsgUpdateParams is the Msg/UpdateParams request type.
 * Deprecated since v6.x as is replaced by MsgUpdateDefaultQuotingParams.
 */

/** @deprecated */

export interface MsgUpdateParamsSDKType {
  authority: string;
  /** The parameters to update. Each field must be set. */

  params?: ParamsSDKType;
}
/**
 * MsgSetVaultQuotingParams is the Msg/SetVaultQuotingParams request type.
 * Deprecated since v6.x as is replaced by MsgSetVaultParams.
 */

/** @deprecated */

export interface MsgSetVaultQuotingParams {
  authority: string;
  /** The vault to set quoting params of. */

  vaultId?: VaultId;
  /** The quoting parameters to set. Each field must be set. */

  quotingParams?: QuotingParams;
}
/**
 * MsgSetVaultQuotingParams is the Msg/SetVaultQuotingParams request type.
 * Deprecated since v6.x as is replaced by MsgSetVaultParams.
 */

/** @deprecated */

export interface MsgSetVaultQuotingParamsSDKType {
  authority: string;
  /** The vault to set quoting params of. */

  vault_id?: VaultIdSDKType;
  /** The quoting parameters to set. Each field must be set. */

  quoting_params?: QuotingParamsSDKType;
}

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

function createBaseMsgWithdrawFromMegavault(): MsgWithdrawFromMegavault {
  return {
    subaccountId: undefined,
    shares: undefined,
    minQuoteQuantums: new Uint8Array()
  };
}

export const MsgWithdrawFromMegavault = {
  encode(message: MsgWithdrawFromMegavault, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.subaccountId !== undefined) {
      SubaccountId.encode(message.subaccountId, writer.uint32(10).fork()).ldelim();
    }

    if (message.shares !== undefined) {
      NumShares.encode(message.shares, writer.uint32(18).fork()).ldelim();
    }

    if (message.minQuoteQuantums.length !== 0) {
      writer.uint32(26).bytes(message.minQuoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawFromMegavault {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFromMegavault();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountId = SubaccountId.decode(reader, reader.uint32());
          break;

        case 2:
          message.shares = NumShares.decode(reader, reader.uint32());
          break;

        case 3:
          message.minQuoteQuantums = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgWithdrawFromMegavault>): MsgWithdrawFromMegavault {
    const message = createBaseMsgWithdrawFromMegavault();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? SubaccountId.fromPartial(object.subaccountId) : undefined;
    message.shares = object.shares !== undefined && object.shares !== null ? NumShares.fromPartial(object.shares) : undefined;
    message.minQuoteQuantums = object.minQuoteQuantums ?? new Uint8Array();
    return message;
  }

};

function createBaseMsgWithdrawFromMegavaultResponse(): MsgWithdrawFromMegavaultResponse {
  return {
    quoteQuantums: new Uint8Array()
  };
}

export const MsgWithdrawFromMegavaultResponse = {
  encode(message: MsgWithdrawFromMegavaultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.quoteQuantums.length !== 0) {
      writer.uint32(10).bytes(message.quoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgWithdrawFromMegavaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawFromMegavaultResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.quoteQuantums = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgWithdrawFromMegavaultResponse>): MsgWithdrawFromMegavaultResponse {
    const message = createBaseMsgWithdrawFromMegavaultResponse();
    message.quoteQuantums = object.quoteQuantums ?? new Uint8Array();
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

function createBaseMsgUnlockShares(): MsgUnlockShares {
  return {
    authority: "",
    ownerAddress: ""
  };
}

export const MsgUnlockShares = {
  encode(message: MsgUnlockShares, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.ownerAddress !== "") {
      writer.uint32(18).string(message.ownerAddress);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnlockShares {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnlockShares();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.ownerAddress = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUnlockShares>): MsgUnlockShares {
    const message = createBaseMsgUnlockShares();
    message.authority = object.authority ?? "";
    message.ownerAddress = object.ownerAddress ?? "";
    return message;
  }

};

function createBaseMsgUnlockSharesResponse(): MsgUnlockSharesResponse {
  return {
    unlockedShares: undefined
  };
}

export const MsgUnlockSharesResponse = {
  encode(message: MsgUnlockSharesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.unlockedShares !== undefined) {
      NumShares.encode(message.unlockedShares, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnlockSharesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnlockSharesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.unlockedShares = NumShares.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUnlockSharesResponse>): MsgUnlockSharesResponse {
    const message = createBaseMsgUnlockSharesResponse();
    message.unlockedShares = object.unlockedShares !== undefined && object.unlockedShares !== null ? NumShares.fromPartial(object.unlockedShares) : undefined;
    return message;
  }

};

function createBaseMsgUpdateOperatorParams(): MsgUpdateOperatorParams {
  return {
    authority: "",
    params: undefined
  };
}

export const MsgUpdateOperatorParams = {
  encode(message: MsgUpdateOperatorParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.params !== undefined) {
      OperatorParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateOperatorParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateOperatorParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params = OperatorParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdateOperatorParams>): MsgUpdateOperatorParams {
    const message = createBaseMsgUpdateOperatorParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? OperatorParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgUpdateOperatorParamsResponse(): MsgUpdateOperatorParamsResponse {
  return {};
}

export const MsgUpdateOperatorParamsResponse = {
  encode(_: MsgUpdateOperatorParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateOperatorParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateOperatorParamsResponse();

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

  fromPartial(_: DeepPartial<MsgUpdateOperatorParamsResponse>): MsgUpdateOperatorParamsResponse {
    const message = createBaseMsgUpdateOperatorParamsResponse();
    return message;
  }

};

function createBaseMsgAllocateToVault(): MsgAllocateToVault {
  return {
    authority: "",
    vaultId: undefined,
    quoteQuantums: new Uint8Array()
  };
}

export const MsgAllocateToVault = {
  encode(message: MsgAllocateToVault, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.vaultId !== undefined) {
      VaultId.encode(message.vaultId, writer.uint32(18).fork()).ldelim();
    }

    if (message.quoteQuantums.length !== 0) {
      writer.uint32(26).bytes(message.quoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgAllocateToVault {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAllocateToVault();

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
          message.quoteQuantums = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgAllocateToVault>): MsgAllocateToVault {
    const message = createBaseMsgAllocateToVault();
    message.authority = object.authority ?? "";
    message.vaultId = object.vaultId !== undefined && object.vaultId !== null ? VaultId.fromPartial(object.vaultId) : undefined;
    message.quoteQuantums = object.quoteQuantums ?? new Uint8Array();
    return message;
  }

};

function createBaseMsgAllocateToVaultResponse(): MsgAllocateToVaultResponse {
  return {};
}

export const MsgAllocateToVaultResponse = {
  encode(_: MsgAllocateToVaultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgAllocateToVaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAllocateToVaultResponse();

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

  fromPartial(_: DeepPartial<MsgAllocateToVaultResponse>): MsgAllocateToVaultResponse {
    const message = createBaseMsgAllocateToVaultResponse();
    return message;
  }

};

function createBaseMsgRetrieveFromVault(): MsgRetrieveFromVault {
  return {
    authority: "",
    vaultId: undefined,
    quoteQuantums: new Uint8Array()
  };
}

export const MsgRetrieveFromVault = {
  encode(message: MsgRetrieveFromVault, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.vaultId !== undefined) {
      VaultId.encode(message.vaultId, writer.uint32(18).fork()).ldelim();
    }

    if (message.quoteQuantums.length !== 0) {
      writer.uint32(26).bytes(message.quoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRetrieveFromVault {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRetrieveFromVault();

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
          message.quoteQuantums = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgRetrieveFromVault>): MsgRetrieveFromVault {
    const message = createBaseMsgRetrieveFromVault();
    message.authority = object.authority ?? "";
    message.vaultId = object.vaultId !== undefined && object.vaultId !== null ? VaultId.fromPartial(object.vaultId) : undefined;
    message.quoteQuantums = object.quoteQuantums ?? new Uint8Array();
    return message;
  }

};

function createBaseMsgRetrieveFromVaultResponse(): MsgRetrieveFromVaultResponse {
  return {};
}

export const MsgRetrieveFromVaultResponse = {
  encode(_: MsgRetrieveFromVaultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgRetrieveFromVaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRetrieveFromVaultResponse();

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

  fromPartial(_: DeepPartial<MsgRetrieveFromVaultResponse>): MsgRetrieveFromVaultResponse {
    const message = createBaseMsgRetrieveFromVaultResponse();
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