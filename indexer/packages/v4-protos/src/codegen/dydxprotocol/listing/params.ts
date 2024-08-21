import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** ListingVaultDepositParams represents the params for PML megavault deposits */

export interface ListingVaultDepositParams {
  /** Amount that will be deposited into the new market vault exclusively */
  newVaultDepositAmount: Uint8Array;
  /**
   * Amount deposited into the main vault exclusively. This amount does not
   * include the amount deposited into the new vault.
   */

  mainVaultDepositAmount: Uint8Array;
  /** Lockup period for this deposit */

  numBlocksToLockShares: number;
}
/** ListingVaultDepositParams represents the params for PML megavault deposits */

export interface ListingVaultDepositParamsSDKType {
  /** Amount that will be deposited into the new market vault exclusively */
  new_vault_deposit_amount: Uint8Array;
  /**
   * Amount deposited into the main vault exclusively. This amount does not
   * include the amount deposited into the new vault.
   */

  main_vault_deposit_amount: Uint8Array;
  /** Lockup period for this deposit */

  num_blocks_to_lock_shares: number;
}

function createBaseListingVaultDepositParams(): ListingVaultDepositParams {
  return {
    newVaultDepositAmount: new Uint8Array(),
    mainVaultDepositAmount: new Uint8Array(),
    numBlocksToLockShares: 0
  };
}

export const ListingVaultDepositParams = {
  encode(message: ListingVaultDepositParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.newVaultDepositAmount.length !== 0) {
      writer.uint32(10).bytes(message.newVaultDepositAmount);
    }

    if (message.mainVaultDepositAmount.length !== 0) {
      writer.uint32(18).bytes(message.mainVaultDepositAmount);
    }

    if (message.numBlocksToLockShares !== 0) {
      writer.uint32(24).uint32(message.numBlocksToLockShares);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListingVaultDepositParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListingVaultDepositParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.newVaultDepositAmount = reader.bytes();
          break;

        case 2:
          message.mainVaultDepositAmount = reader.bytes();
          break;

        case 3:
          message.numBlocksToLockShares = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ListingVaultDepositParams>): ListingVaultDepositParams {
    const message = createBaseListingVaultDepositParams();
    message.newVaultDepositAmount = object.newVaultDepositAmount ?? new Uint8Array();
    message.mainVaultDepositAmount = object.mainVaultDepositAmount ?? new Uint8Array();
    message.numBlocksToLockShares = object.numBlocksToLockShares ?? 0;
    return message;
  }

};