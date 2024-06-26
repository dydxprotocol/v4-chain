import { Params, ParamsSDKType } from "./params";
import { VaultId, VaultIdSDKType, NumShares, NumSharesSDKType } from "./vault";
import { OwnerShare, OwnerShareSDKType } from "./query";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines `x/vault`'s genesis state. */

export interface GenesisState {
  /** The parameters of the module. */
  params?: Params;
  /** The vaults. */

  vaults: Vault[];
}
/** GenesisState defines `x/vault`'s genesis state. */

export interface GenesisStateSDKType {
  /** The parameters of the module. */
  params?: ParamsSDKType;
  /** The vaults. */

  vaults: VaultSDKType[];
}
/** Vault defines the total shares and owner shares of a vault. */

export interface Vault {
  /** The ID of the vault. */
  vaultId?: VaultId;
  /** The total number of shares in the vault. */

  totalShares?: NumShares;
  /** The shares of each owner in the vault. */

  ownerShares: OwnerShare[];
}
/** Vault defines the total shares and owner shares of a vault. */

export interface VaultSDKType {
  /** The ID of the vault. */
  vault_id?: VaultIdSDKType;
  /** The total number of shares in the vault. */

  total_shares?: NumSharesSDKType;
  /** The shares of each owner in the vault. */

  owner_shares: OwnerShareSDKType[];
}

function createBaseGenesisState(): GenesisState {
  return {
    params: undefined,
    vaults: []
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.vaults) {
      Vault.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;

        case 2:
          message.vaults.push(Vault.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    message.vaults = object.vaults?.map(e => Vault.fromPartial(e)) || [];
    return message;
  }

};

function createBaseVault(): Vault {
  return {
    vaultId: undefined,
    totalShares: undefined,
    ownerShares: []
  };
}

export const Vault = {
  encode(message: Vault, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.vaultId !== undefined) {
      VaultId.encode(message.vaultId, writer.uint32(10).fork()).ldelim();
    }

    if (message.totalShares !== undefined) {
      NumShares.encode(message.totalShares, writer.uint32(18).fork()).ldelim();
    }

    for (const v of message.ownerShares) {
      OwnerShare.encode(v!, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Vault {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVault();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.vaultId = VaultId.decode(reader, reader.uint32());
          break;

        case 2:
          message.totalShares = NumShares.decode(reader, reader.uint32());
          break;

        case 3:
          message.ownerShares.push(OwnerShare.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Vault>): Vault {
    const message = createBaseVault();
    message.vaultId = object.vaultId !== undefined && object.vaultId !== null ? VaultId.fromPartial(object.vaultId) : undefined;
    message.totalShares = object.totalShares !== undefined && object.totalShares !== null ? NumShares.fromPartial(object.totalShares) : undefined;
    message.ownerShares = object.ownerShares?.map(e => OwnerShare.fromPartial(e)) || [];
    return message;
  }

};