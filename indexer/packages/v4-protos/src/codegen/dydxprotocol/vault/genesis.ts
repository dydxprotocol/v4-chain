import { Params, ParamsSDKType } from "./params";
import { VaultId, VaultIdSDKType, NumShares, NumSharesSDKType, OwnerShare, OwnerShareSDKType, VaultParams, VaultParamsSDKType } from "./vault";
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
  /** The individual parameters of the vault. */

  vaultParams?: VaultParams;
  /** The client IDs of the most recently placed orders of the vault. */

  mostRecentClientIds: number[];
}
/** Vault defines the total shares and owner shares of a vault. */

export interface VaultSDKType {
  /** The ID of the vault. */
  vault_id?: VaultIdSDKType;
  /** The total number of shares in the vault. */

  total_shares?: NumSharesSDKType;
  /** The shares of each owner in the vault. */

  owner_shares: OwnerShareSDKType[];
  /** The individual parameters of the vault. */

  vault_params?: VaultParamsSDKType;
  /** The client IDs of the most recently placed orders of the vault. */

  most_recent_client_ids: number[];
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
    ownerShares: [],
    vaultParams: undefined,
    mostRecentClientIds: []
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

    if (message.vaultParams !== undefined) {
      VaultParams.encode(message.vaultParams, writer.uint32(34).fork()).ldelim();
    }

    writer.uint32(42).fork();

    for (const v of message.mostRecentClientIds) {
      writer.uint32(v);
    }

    writer.ldelim();
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

        case 4:
          message.vaultParams = VaultParams.decode(reader, reader.uint32());
          break;

        case 5:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;

            while (reader.pos < end2) {
              message.mostRecentClientIds.push(reader.uint32());
            }
          } else {
            message.mostRecentClientIds.push(reader.uint32());
          }

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
    message.vaultParams = object.vaultParams !== undefined && object.vaultParams !== null ? VaultParams.fromPartial(object.vaultParams) : undefined;
    message.mostRecentClientIds = object.mostRecentClientIds?.map(e => e) || [];
    return message;
  }

};