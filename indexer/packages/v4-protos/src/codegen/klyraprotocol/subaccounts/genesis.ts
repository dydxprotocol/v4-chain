import { Subaccount, SubaccountSDKType } from "./subaccount";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the subaccounts module's genesis state. */

export interface GenesisState {
  subaccounts: Subaccount[];
}
/** GenesisState defines the subaccounts module's genesis state. */

export interface GenesisStateSDKType {
  subaccounts: SubaccountSDKType[];
}

function createBaseGenesisState(): GenesisState {
  return {
    subaccounts: []
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.subaccounts) {
      Subaccount.encode(v!, writer.uint32(10).fork()).ldelim();
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
          message.subaccounts.push(Subaccount.decode(reader, reader.uint32()));
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
    message.subaccounts = object.subaccounts?.map(e => Subaccount.fromPartial(e)) || [];
    return message;
  }

};