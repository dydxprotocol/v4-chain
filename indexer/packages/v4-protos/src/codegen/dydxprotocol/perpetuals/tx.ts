import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * FundingPremium represents a funding premium value for a perpetual
 * market. Can be used to represent a premium vote or a premium sample.
 */

export interface FundingPremium {
  /** The id of the perpetual market. */
  perpetualId: number;
  /** The sampled premium rate. In parts-per-million. */

  premiumPpm: number;
}
/**
 * FundingPremium represents a funding premium value for a perpetual
 * market. Can be used to represent a premium vote or a premium sample.
 */

export interface FundingPremiumSDKType {
  /** The id of the perpetual market. */
  perpetual_id: number;
  /** The sampled premium rate. In parts-per-million. */

  premium_ppm: number;
}
/** MsgAddPremiumVotes is a request type for the AddPremiumVotes method. */

export interface MsgAddPremiumVotes {
  votes: FundingPremium[];
}
/** MsgAddPremiumVotes is a request type for the AddPremiumVotes method. */

export interface MsgAddPremiumVotesSDKType {
  votes: FundingPremiumSDKType[];
}
/**
 * MsgAddPremiumVotesResponse defines the AddPremiumVotes
 * response type.
 */

export interface MsgAddPremiumVotesResponse {}
/**
 * MsgAddPremiumVotesResponse defines the AddPremiumVotes
 * response type.
 */

export interface MsgAddPremiumVotesResponseSDKType {}

function createBaseFundingPremium(): FundingPremium {
  return {
    perpetualId: 0,
    premiumPpm: 0
  };
}

export const FundingPremium = {
  encode(message: FundingPremium, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.perpetualId !== 0) {
      writer.uint32(8).uint32(message.perpetualId);
    }

    if (message.premiumPpm !== 0) {
      writer.uint32(16).int32(message.premiumPpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FundingPremium {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFundingPremium();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.perpetualId = reader.uint32();
          break;

        case 2:
          message.premiumPpm = reader.int32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<FundingPremium>): FundingPremium {
    const message = createBaseFundingPremium();
    message.perpetualId = object.perpetualId ?? 0;
    message.premiumPpm = object.premiumPpm ?? 0;
    return message;
  }

};

function createBaseMsgAddPremiumVotes(): MsgAddPremiumVotes {
  return {
    votes: []
  };
}

export const MsgAddPremiumVotes = {
  encode(message: MsgAddPremiumVotes, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.votes) {
      FundingPremium.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgAddPremiumVotes {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddPremiumVotes();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.votes.push(FundingPremium.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgAddPremiumVotes>): MsgAddPremiumVotes {
    const message = createBaseMsgAddPremiumVotes();
    message.votes = object.votes?.map(e => FundingPremium.fromPartial(e)) || [];
    return message;
  }

};

function createBaseMsgAddPremiumVotesResponse(): MsgAddPremiumVotesResponse {
  return {};
}

export const MsgAddPremiumVotesResponse = {
  encode(_: MsgAddPremiumVotesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgAddPremiumVotesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddPremiumVotesResponse();

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

  fromPartial(_: DeepPartial<MsgAddPremiumVotesResponse>): MsgAddPremiumVotesResponse {
    const message = createBaseMsgAddPremiumVotesResponse();
    return message;
  }

};