import { Params, ParamsSDKType } from "./params";
import { UserStats, UserStatsSDKType } from "./stats";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the stats module's genesis state. */

export interface GenesisState {
  /** The parameters of the module. */
  params?: Params;
  addressToUserStats: AddressToUserStats[];
}
/** GenesisState defines the stats module's genesis state. */

export interface GenesisStateSDKType {
  /** The parameters of the module. */
  params?: ParamsSDKType;
  address_to_user_stats: AddressToUserStatsSDKType[];
}
/** AddressToUserStats is a struct that contains the user stats for an address. */

export interface AddressToUserStats {
  /** The address of the user. */
  address: string;
  /** The user stats for the address. */

  userStats?: UserStats;
}
/** AddressToUserStats is a struct that contains the user stats for an address. */

export interface AddressToUserStatsSDKType {
  /** The address of the user. */
  address: string;
  /** The user stats for the address. */

  user_stats?: UserStatsSDKType;
}

function createBaseGenesisState(): GenesisState {
  return {
    params: undefined,
    addressToUserStats: []
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.addressToUserStats) {
      AddressToUserStats.encode(v!, writer.uint32(18).fork()).ldelim();
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
          message.addressToUserStats.push(AddressToUserStats.decode(reader, reader.uint32()));
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
    message.addressToUserStats = object.addressToUserStats?.map(e => AddressToUserStats.fromPartial(e)) || [];
    return message;
  }

};

function createBaseAddressToUserStats(): AddressToUserStats {
  return {
    address: "",
    userStats: undefined
  };
}

export const AddressToUserStats = {
  encode(message: AddressToUserStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    if (message.userStats !== undefined) {
      UserStats.encode(message.userStats, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddressToUserStats {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddressToUserStats();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        case 2:
          message.userStats = UserStats.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AddressToUserStats>): AddressToUserStats {
    const message = createBaseAddressToUserStats();
    message.address = object.address ?? "";
    message.userStats = object.userStats !== undefined && object.userStats !== null ? UserStats.fromPartial(object.userStats) : undefined;
    return message;
  }

};