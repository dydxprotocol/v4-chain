import { ListingVaultDepositParams, ListingVaultDepositParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines `x/listing`'s genesis state. */

export interface GenesisState {
  /**
   * hard_cap_for_markets is the hard cap for the number of markets that can be
   * listed
   */
  hardCapForMarkets: number;
  /** listing_vault_deposit_params is the params for PML megavault deposits */

  listingVaultDepositParams?: ListingVaultDepositParams;
}
/** GenesisState defines `x/listing`'s genesis state. */

export interface GenesisStateSDKType {
  /**
   * hard_cap_for_markets is the hard cap for the number of markets that can be
   * listed
   */
  hard_cap_for_markets: number;
  /** listing_vault_deposit_params is the params for PML megavault deposits */

  listing_vault_deposit_params?: ListingVaultDepositParamsSDKType;
}

function createBaseGenesisState(): GenesisState {
  return {
    hardCapForMarkets: 0,
    listingVaultDepositParams: undefined
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.hardCapForMarkets !== 0) {
      writer.uint32(8).uint32(message.hardCapForMarkets);
    }

    if (message.listingVaultDepositParams !== undefined) {
      ListingVaultDepositParams.encode(message.listingVaultDepositParams, writer.uint32(18).fork()).ldelim();
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
          message.hardCapForMarkets = reader.uint32();
          break;

        case 2:
          message.listingVaultDepositParams = ListingVaultDepositParams.decode(reader, reader.uint32());
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
    message.hardCapForMarkets = object.hardCapForMarkets ?? 0;
    message.listingVaultDepositParams = object.listingVaultDepositParams !== undefined && object.listingVaultDepositParams !== null ? ListingVaultDepositParams.fromPartial(object.listingVaultDepositParams) : undefined;
    return message;
  }

};