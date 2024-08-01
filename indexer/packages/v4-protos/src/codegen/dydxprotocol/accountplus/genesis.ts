import { AccountState, AccountStateSDKType } from "./accountplus";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Module genesis state */

export interface GenesisState {
  accounts: AccountState[];
}
/** Module genesis state */

export interface GenesisStateSDKType {
  accounts: AccountStateSDKType[];
}

function createBaseGenesisState(): GenesisState {
  return {
    accounts: []
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.accounts) {
      AccountState.encode(v!, writer.uint32(10).fork()).ldelim();
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
          message.accounts.push(AccountState.decode(reader, reader.uint32()));
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
    message.accounts = object.accounts?.map(e => AccountState.fromPartial(e)) || [];
    return message;
  }

};