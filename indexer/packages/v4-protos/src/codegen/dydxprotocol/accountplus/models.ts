import * as _m0 from "protobufjs/minimal";
import { Long, DeepPartial } from "../../helpers";
/**
 * AccountAuthenticator represents a foundational model for all authenticators.
 * It provides extensibility by allowing concrete types to interpret and
 * validate transactions based on the encapsulated data.
 */

export interface AccountAuthenticator {
  /** ID uniquely identifies the authenticator instance. */
  id: Long;
  /**
   * Type specifies the category of the AccountAuthenticator.
   * This type information is essential for differentiating authenticators
   * and ensuring precise data retrieval from the storage layer.
   */

  type: string;
  /**
   * Config is a versatile field used in conjunction with the specific type of
   * account authenticator to facilitate complex authentication processes.
   * The interpretation of this field is overloaded, enabling multiple
   * authenticators to utilize it for their respective purposes.
   */

  config: Uint8Array;
}
/**
 * AccountAuthenticator represents a foundational model for all authenticators.
 * It provides extensibility by allowing concrete types to interpret and
 * validate transactions based on the encapsulated data.
 */

export interface AccountAuthenticatorSDKType {
  /** ID uniquely identifies the authenticator instance. */
  id: Long;
  /**
   * Type specifies the category of the AccountAuthenticator.
   * This type information is essential for differentiating authenticators
   * and ensuring precise data retrieval from the storage layer.
   */

  type: string;
  /**
   * Config is a versatile field used in conjunction with the specific type of
   * account authenticator to facilitate complex authentication processes.
   * The interpretation of this field is overloaded, enabling multiple
   * authenticators to utilize it for their respective purposes.
   */

  config: Uint8Array;
}

function createBaseAccountAuthenticator(): AccountAuthenticator {
  return {
    id: Long.UZERO,
    type: "",
    config: new Uint8Array()
  };
}

export const AccountAuthenticator = {
  encode(message: AccountAuthenticator, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.id.isZero()) {
      writer.uint32(8).uint64(message.id);
    }

    if (message.type !== "") {
      writer.uint32(18).string(message.type);
    }

    if (message.config.length !== 0) {
      writer.uint32(26).bytes(message.config);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AccountAuthenticator {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAccountAuthenticator();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = (reader.uint64() as Long);
          break;

        case 2:
          message.type = reader.string();
          break;

        case 3:
          message.config = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AccountAuthenticator>): AccountAuthenticator {
    const message = createBaseAccountAuthenticator();
    message.id = object.id !== undefined && object.id !== null ? Long.fromValue(object.id) : Long.UZERO;
    message.type = object.type ?? "";
    message.config = object.config ?? new Uint8Array();
    return message;
  }

};