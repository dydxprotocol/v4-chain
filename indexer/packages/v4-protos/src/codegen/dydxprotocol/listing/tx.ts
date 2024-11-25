import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { ListingVaultDepositParams, ListingVaultDepositParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * MsgSetMarketsHardCap is used to set a hard cap on the number of markets
 * listed
 */

export interface MsgSetMarketsHardCap {
  authority: string;
  /** Hard cap for the total number of markets listed */

  hardCapForMarkets: number;
}
/**
 * MsgSetMarketsHardCap is used to set a hard cap on the number of markets
 * listed
 */

export interface MsgSetMarketsHardCapSDKType {
  authority: string;
  /** Hard cap for the total number of markets listed */

  hard_cap_for_markets: number;
}
/** MsgSetMarketsHardCapResponse defines the MsgSetMarketsHardCap response */

export interface MsgSetMarketsHardCapResponse {}
/** MsgSetMarketsHardCapResponse defines the MsgSetMarketsHardCap response */

export interface MsgSetMarketsHardCapResponseSDKType {}
/**
 * MsgCreateMarketPermissionless is a message used to create new markets without
 * going through x/gov
 */

export interface MsgCreateMarketPermissionless {
  /** The name of the `Perpetual` (e.g. `BTC-USD`). */
  ticker: string;
  /** The subaccount to deposit from. */

  subaccountId?: SubaccountId;
}
/**
 * MsgCreateMarketPermissionless is a message used to create new markets without
 * going through x/gov
 */

export interface MsgCreateMarketPermissionlessSDKType {
  /** The name of the `Perpetual` (e.g. `BTC-USD`). */
  ticker: string;
  /** The subaccount to deposit from. */

  subaccount_id?: SubaccountIdSDKType;
}
/**
 * MsgCreateMarketPermissionlessResponse defines the
 * MsgCreateMarketPermissionless response
 */

export interface MsgCreateMarketPermissionlessResponse {}
/**
 * MsgCreateMarketPermissionlessResponse defines the
 * MsgCreateMarketPermissionless response
 */

export interface MsgCreateMarketPermissionlessResponseSDKType {}
/**
 * MsgSetListingVaultDepositParams is a message used to set PML megavault
 * deposit params
 */

export interface MsgSetListingVaultDepositParams {
  authority: string;
  /** Params which define the vault deposit for market listing */

  params?: ListingVaultDepositParams;
}
/**
 * MsgSetListingVaultDepositParams is a message used to set PML megavault
 * deposit params
 */

export interface MsgSetListingVaultDepositParamsSDKType {
  authority: string;
  /** Params which define the vault deposit for market listing */

  params?: ListingVaultDepositParamsSDKType;
}
/**
 * MsgSetListingVaultDepositParamsResponse defines the
 * MsgSetListingVaultDepositParams response
 */

export interface MsgSetListingVaultDepositParamsResponse {}
/**
 * MsgSetListingVaultDepositParamsResponse defines the
 * MsgSetListingVaultDepositParams response
 */

export interface MsgSetListingVaultDepositParamsResponseSDKType {}
/**
 * MsgUpgradeIsolatedPerpetualToCross is used to upgrade a market from
 * isolated margin to cross margin.
 */

export interface MsgUpgradeIsolatedPerpetualToCross {
  authority: string;
  /** ID of the perpetual to be upgraded to CROSS */

  perpetualId: number;
}
/**
 * MsgUpgradeIsolatedPerpetualToCross is used to upgrade a market from
 * isolated margin to cross margin.
 */

export interface MsgUpgradeIsolatedPerpetualToCrossSDKType {
  authority: string;
  /** ID of the perpetual to be upgraded to CROSS */

  perpetual_id: number;
}
/**
 * MsgUpgradeIsolatedPerpetualToCrossResponse defines the
 * UpgradeIsolatedPerpetualToCross response type.
 */

export interface MsgUpgradeIsolatedPerpetualToCrossResponse {}
/**
 * MsgUpgradeIsolatedPerpetualToCrossResponse defines the
 * UpgradeIsolatedPerpetualToCross response type.
 */

export interface MsgUpgradeIsolatedPerpetualToCrossResponseSDKType {}

function createBaseMsgSetMarketsHardCap(): MsgSetMarketsHardCap {
  return {
    authority: "",
    hardCapForMarkets: 0
  };
}

export const MsgSetMarketsHardCap = {
  encode(message: MsgSetMarketsHardCap, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.hardCapForMarkets !== 0) {
      writer.uint32(16).uint32(message.hardCapForMarkets);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMarketsHardCap {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetMarketsHardCap();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.hardCapForMarkets = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetMarketsHardCap>): MsgSetMarketsHardCap {
    const message = createBaseMsgSetMarketsHardCap();
    message.authority = object.authority ?? "";
    message.hardCapForMarkets = object.hardCapForMarkets ?? 0;
    return message;
  }

};

function createBaseMsgSetMarketsHardCapResponse(): MsgSetMarketsHardCapResponse {
  return {};
}

export const MsgSetMarketsHardCapResponse = {
  encode(_: MsgSetMarketsHardCapResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMarketsHardCapResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetMarketsHardCapResponse();

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

  fromPartial(_: DeepPartial<MsgSetMarketsHardCapResponse>): MsgSetMarketsHardCapResponse {
    const message = createBaseMsgSetMarketsHardCapResponse();
    return message;
  }

};

function createBaseMsgCreateMarketPermissionless(): MsgCreateMarketPermissionless {
  return {
    ticker: "",
    subaccountId: undefined
  };
}

export const MsgCreateMarketPermissionless = {
  encode(message: MsgCreateMarketPermissionless, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ticker !== "") {
      writer.uint32(10).string(message.ticker);
    }

    if (message.subaccountId !== undefined) {
      SubaccountId.encode(message.subaccountId, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateMarketPermissionless {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateMarketPermissionless();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.ticker = reader.string();
          break;

        case 2:
          message.subaccountId = SubaccountId.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgCreateMarketPermissionless>): MsgCreateMarketPermissionless {
    const message = createBaseMsgCreateMarketPermissionless();
    message.ticker = object.ticker ?? "";
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? SubaccountId.fromPartial(object.subaccountId) : undefined;
    return message;
  }

};

function createBaseMsgCreateMarketPermissionlessResponse(): MsgCreateMarketPermissionlessResponse {
  return {};
}

export const MsgCreateMarketPermissionlessResponse = {
  encode(_: MsgCreateMarketPermissionlessResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateMarketPermissionlessResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateMarketPermissionlessResponse();

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

  fromPartial(_: DeepPartial<MsgCreateMarketPermissionlessResponse>): MsgCreateMarketPermissionlessResponse {
    const message = createBaseMsgCreateMarketPermissionlessResponse();
    return message;
  }

};

function createBaseMsgSetListingVaultDepositParams(): MsgSetListingVaultDepositParams {
  return {
    authority: "",
    params: undefined
  };
}

export const MsgSetListingVaultDepositParams = {
  encode(message: MsgSetListingVaultDepositParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.params !== undefined) {
      ListingVaultDepositParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetListingVaultDepositParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetListingVaultDepositParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params = ListingVaultDepositParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetListingVaultDepositParams>): MsgSetListingVaultDepositParams {
    const message = createBaseMsgSetListingVaultDepositParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? ListingVaultDepositParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgSetListingVaultDepositParamsResponse(): MsgSetListingVaultDepositParamsResponse {
  return {};
}

export const MsgSetListingVaultDepositParamsResponse = {
  encode(_: MsgSetListingVaultDepositParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetListingVaultDepositParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetListingVaultDepositParamsResponse();

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

  fromPartial(_: DeepPartial<MsgSetListingVaultDepositParamsResponse>): MsgSetListingVaultDepositParamsResponse {
    const message = createBaseMsgSetListingVaultDepositParamsResponse();
    return message;
  }

};

function createBaseMsgUpgradeIsolatedPerpetualToCross(): MsgUpgradeIsolatedPerpetualToCross {
  return {
    authority: "",
    perpetualId: 0
  };
}

export const MsgUpgradeIsolatedPerpetualToCross = {
  encode(message: MsgUpgradeIsolatedPerpetualToCross, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.perpetualId !== 0) {
      writer.uint32(16).uint32(message.perpetualId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpgradeIsolatedPerpetualToCross {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpgradeIsolatedPerpetualToCross();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.perpetualId = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpgradeIsolatedPerpetualToCross>): MsgUpgradeIsolatedPerpetualToCross {
    const message = createBaseMsgUpgradeIsolatedPerpetualToCross();
    message.authority = object.authority ?? "";
    message.perpetualId = object.perpetualId ?? 0;
    return message;
  }

};

function createBaseMsgUpgradeIsolatedPerpetualToCrossResponse(): MsgUpgradeIsolatedPerpetualToCrossResponse {
  return {};
}

export const MsgUpgradeIsolatedPerpetualToCrossResponse = {
  encode(_: MsgUpgradeIsolatedPerpetualToCrossResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpgradeIsolatedPerpetualToCrossResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpgradeIsolatedPerpetualToCrossResponse();

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

  fromPartial(_: DeepPartial<MsgUpgradeIsolatedPerpetualToCrossResponse>): MsgUpgradeIsolatedPerpetualToCrossResponse {
    const message = createBaseMsgUpgradeIsolatedPerpetualToCrossResponse();
    return message;
  }

};