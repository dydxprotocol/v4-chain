import { AccountAuthenticator, AccountAuthenticatorSDKType } from "./models";
import { AccountState, AccountStateSDKType } from "./accountplus";
import { Params, ParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/**
 * AuthenticatorData represents a genesis exported account with Authenticators.
 * The address is used as the key, and the account authenticators are stored in
 * the authenticators field.
 */

export interface AuthenticatorData {
  /** address is an account address, one address can have many authenticators */
  address: string;
  /**
   * authenticators are the account's authenticators, these can be multiple
   * types including SignatureVerification, AllOfs, CosmWasmAuthenticators, etc
   */

  authenticators: AccountAuthenticator[];
}
/**
 * AuthenticatorData represents a genesis exported account with Authenticators.
 * The address is used as the key, and the account authenticators are stored in
 * the authenticators field.
 */

export interface AuthenticatorDataSDKType {
  /** address is an account address, one address can have many authenticators */
  address: string;
  /**
   * authenticators are the account's authenticators, these can be multiple
   * types including SignatureVerification, AllOfs, CosmWasmAuthenticators, etc
   */

  authenticators: AccountAuthenticatorSDKType[];
}
/** Module genesis state */

export interface GenesisState {
  accounts: AccountState[];
  /** params define the parameters for the authenticator module. */

  params?: Params;
  /** next_authenticator_id is the next available authenticator ID. */

  nextAuthenticatorId: Long;
  /**
   * authenticator_data contains the data for multiple accounts, each with their
   * authenticators.
   */

  authenticatorData: AuthenticatorData[];
}
/** Module genesis state */

export interface GenesisStateSDKType {
  accounts: AccountStateSDKType[];
  /** params define the parameters for the authenticator module. */

  params?: ParamsSDKType;
  /** next_authenticator_id is the next available authenticator ID. */

  next_authenticator_id: Long;
  /**
   * authenticator_data contains the data for multiple accounts, each with their
   * authenticators.
   */

  authenticator_data: AuthenticatorDataSDKType[];
}

function createBaseAuthenticatorData(): AuthenticatorData {
  return {
    address: "",
    authenticators: []
  };
}

export const AuthenticatorData = {
  encode(message: AuthenticatorData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    for (const v of message.authenticators) {
      AccountAuthenticator.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthenticatorData {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticatorData();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        case 2:
          message.authenticators.push(AccountAuthenticator.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AuthenticatorData>): AuthenticatorData {
    const message = createBaseAuthenticatorData();
    message.address = object.address ?? "";
    message.authenticators = object.authenticators?.map(e => AccountAuthenticator.fromPartial(e)) || [];
    return message;
  }

};

function createBaseGenesisState(): GenesisState {
  return {
    accounts: [],
    params: undefined,
    nextAuthenticatorId: Long.UZERO,
    authenticatorData: []
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.accounts) {
      AccountState.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    if (!message.nextAuthenticatorId.isZero()) {
      writer.uint32(24).uint64(message.nextAuthenticatorId);
    }

    for (const v of message.authenticatorData) {
      AuthenticatorData.encode(v!, writer.uint32(34).fork()).ldelim();
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

        case 2:
          message.params = Params.decode(reader, reader.uint32());
          break;

        case 3:
          message.nextAuthenticatorId = (reader.uint64() as Long);
          break;

        case 4:
          message.authenticatorData.push(AuthenticatorData.decode(reader, reader.uint32()));
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
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    message.nextAuthenticatorId = object.nextAuthenticatorId !== undefined && object.nextAuthenticatorId !== null ? Long.fromValue(object.nextAuthenticatorId) : Long.UZERO;
    message.authenticatorData = object.authenticatorData?.map(e => AuthenticatorData.fromPartial(e)) || [];
    return message;
  }

};