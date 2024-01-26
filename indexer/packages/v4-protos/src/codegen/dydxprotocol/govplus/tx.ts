import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgSlashValidator is the Msg/SlashValidator request type. */

export interface MsgSlashValidator {
  authority: string;
  /** Address of the validator to slash */

  validatorAddress: string;
  /**
   * Colloquially, the height at which the validator is deemed to have
   * misbehaved. In practice, this is the height used to determine the targets
   * of the slash. For example, undelegating after this height will not escape
   * slashing. This height should be set to a recent height at the time of the
   * proposal to prevent delegators from undelegating during the vote period.
   * 
   * NB: At the time this message is applied, this height must have occured
   * equal to or less than an unbonding period in the past in order for the
   * slash to be effective.
   */

  infractionHeight: number;
  /**
   * Power of the validator at the specified height. Used to compute the slash
   * amount. Note that this is not directly the token amount; it must be
   * adjusted by power reduction to get consensus power.
   */

  powerAtInfractionHeight: Uint8Array;
  /**
   * Multiplier for how much of the validator's stake should be slashed.
   * slash_factor * power = amount slashed
   */

  slashFactor: string;
}
/** MsgSlashValidator is the Msg/SlashValidator request type. */

export interface MsgSlashValidatorSDKType {
  authority: string;
  /** Address of the validator to slash */

  validator_address: string;
  /**
   * Colloquially, the height at which the validator is deemed to have
   * misbehaved. In practice, this is the height used to determine the targets
   * of the slash. For example, undelegating after this height will not escape
   * slashing. This height should be set to a recent height at the time of the
   * proposal to prevent delegators from undelegating during the vote period.
   * 
   * NB: At the time this message is applied, this height must have occured
   * equal to or less than an unbonding period in the past in order for the
   * slash to be effective.
   */

  infraction_height: number;
  /**
   * Power of the validator at the specified height. Used to compute the slash
   * amount. Note that this is not directly the token amount; it must be
   * adjusted by power reduction to get consensus power.
   */

  power_at_infraction_height: Uint8Array;
  /**
   * Multiplier for how much of the validator's stake should be slashed.
   * slash_factor * power = amount slashed
   */

  slash_factor: string;
}
/** MsgSlashValidatorResponse is the Msg/SlashValidator response type. */

export interface MsgSlashValidatorResponse {}
/** MsgSlashValidatorResponse is the Msg/SlashValidator response type. */

export interface MsgSlashValidatorResponseSDKType {}

function createBaseMsgSlashValidator(): MsgSlashValidator {
  return {
    authority: "",
    validatorAddress: "",
    infractionHeight: 0,
    powerAtInfractionHeight: new Uint8Array(),
    slashFactor: ""
  };
}

export const MsgSlashValidator = {
  encode(message: MsgSlashValidator, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }

    if (message.infractionHeight !== 0) {
      writer.uint32(24).uint32(message.infractionHeight);
    }

    if (message.powerAtInfractionHeight.length !== 0) {
      writer.uint32(34).bytes(message.powerAtInfractionHeight);
    }

    if (message.slashFactor !== "") {
      writer.uint32(42).string(message.slashFactor);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSlashValidator {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSlashValidator();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.validatorAddress = reader.string();
          break;

        case 3:
          message.infractionHeight = reader.uint32();
          break;

        case 4:
          message.powerAtInfractionHeight = reader.bytes();
          break;

        case 5:
          message.slashFactor = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSlashValidator>): MsgSlashValidator {
    const message = createBaseMsgSlashValidator();
    message.authority = object.authority ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    message.infractionHeight = object.infractionHeight ?? 0;
    message.powerAtInfractionHeight = object.powerAtInfractionHeight ?? new Uint8Array();
    message.slashFactor = object.slashFactor ?? "";
    return message;
  }

};

function createBaseMsgSlashValidatorResponse(): MsgSlashValidatorResponse {
  return {};
}

export const MsgSlashValidatorResponse = {
  encode(_: MsgSlashValidatorResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSlashValidatorResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSlashValidatorResponse();

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

  fromPartial(_: DeepPartial<MsgSlashValidatorResponse>): MsgSlashValidatorResponse {
    const message = createBaseMsgSlashValidatorResponse();
    return message;
  }

};