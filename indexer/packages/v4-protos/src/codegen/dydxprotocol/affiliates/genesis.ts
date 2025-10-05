import { AffiliateTiers, AffiliateTiersSDKType, AffiliateParameters, AffiliateParametersSDKType } from "./affiliates";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines generis state of `x/affiliates` */

export interface GenesisState {
  /** The list of affiliate tiers */
  affiliateTiers?: AffiliateTiers;
  /** The affiliate parameters */

  affiliateParameters?: AffiliateParameters;
}
/** GenesisState defines generis state of `x/affiliates` */

export interface GenesisStateSDKType {
  /** The list of affiliate tiers */
  affiliate_tiers?: AffiliateTiersSDKType;
  /** The affiliate parameters */

  affiliate_parameters?: AffiliateParametersSDKType;
}

function createBaseGenesisState(): GenesisState {
  return {
    affiliateTiers: undefined,
    affiliateParameters: undefined
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.affiliateTiers !== undefined) {
      AffiliateTiers.encode(message.affiliateTiers, writer.uint32(10).fork()).ldelim();
    }

    if (message.affiliateParameters !== undefined) {
      AffiliateParameters.encode(message.affiliateParameters, writer.uint32(18).fork()).ldelim();
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
          message.affiliateTiers = AffiliateTiers.decode(reader, reader.uint32());
          break;

        case 2:
          message.affiliateParameters = AffiliateParameters.decode(reader, reader.uint32());
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
    message.affiliateTiers = object.affiliateTiers !== undefined && object.affiliateTiers !== null ? AffiliateTiers.fromPartial(object.affiliateTiers) : undefined;
    message.affiliateParameters = object.affiliateParameters !== undefined && object.affiliateParameters !== null ? AffiliateParameters.fromPartial(object.affiliateParameters) : undefined;
    return message;
  }

};