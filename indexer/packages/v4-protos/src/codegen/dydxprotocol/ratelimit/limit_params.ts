import { Duration, DurationSDKType } from "../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** LimitParams defines rate limit params on a denom. */

export interface LimitParams {
  /**
   * denom is the denomination of the token being rate limited.
   * e.g. ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5
   */
  denom: string;
  /**
   * limiters is a list of rate-limiters on this denom. All limiters
   * must be satified for a withdrawal to proceed.
   */

  limiters: Limiter[];
}
/** LimitParams defines rate limit params on a denom. */

export interface LimitParamsSDKType {
  /**
   * denom is the denomination of the token being rate limited.
   * e.g. ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5
   */
  denom: string;
  /**
   * limiters is a list of rate-limiters on this denom. All limiters
   * must be satified for a withdrawal to proceed.
   */

  limiters: LimiterSDKType[];
}
/** Limiter defines one rate-limiter on a specfic denom. */

export interface Limiter {
  /**
   * period is the rolling time period for which the limit applies
   * e.g. 3600 (an hour)
   */
  period?: Duration;
  /**
   * baseline_minimum is the minimum maximum withdrawal coin amount within the
   * time period.
   * e.g. 100_000_000_000 uusdc for 100k USDC; 5e22 adv4tnt for 50k DV4TNT
   */

  baselineMinimum: Uint8Array;
  /**
   * baseline_tvl_ppm is the maximum ratio of TVL withdrawable in
   * the time period, in part-per-million.
   * e.g. 100_000 (10%)
   */

  baselineTvlPpm: number;
}
/** Limiter defines one rate-limiter on a specfic denom. */

export interface LimiterSDKType {
  /**
   * period is the rolling time period for which the limit applies
   * e.g. 3600 (an hour)
   */
  period?: DurationSDKType;
  /**
   * baseline_minimum is the minimum maximum withdrawal coin amount within the
   * time period.
   * e.g. 100_000_000_000 uusdc for 100k USDC; 5e22 adv4tnt for 50k DV4TNT
   */

  baseline_minimum: Uint8Array;
  /**
   * baseline_tvl_ppm is the maximum ratio of TVL withdrawable in
   * the time period, in part-per-million.
   * e.g. 100_000 (10%)
   */

  baseline_tvl_ppm: number;
}

function createBaseLimitParams(): LimitParams {
  return {
    denom: "",
    limiters: []
  };
}

export const LimitParams = {
  encode(message: LimitParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }

    for (const v of message.limiters) {
      Limiter.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LimitParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLimitParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.denom = reader.string();
          break;

        case 2:
          message.limiters.push(Limiter.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<LimitParams>): LimitParams {
    const message = createBaseLimitParams();
    message.denom = object.denom ?? "";
    message.limiters = object.limiters?.map(e => Limiter.fromPartial(e)) || [];
    return message;
  }

};

function createBaseLimiter(): Limiter {
  return {
    period: undefined,
    baselineMinimum: new Uint8Array(),
    baselineTvlPpm: 0
  };
}

export const Limiter = {
  encode(message: Limiter, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.period !== undefined) {
      Duration.encode(message.period, writer.uint32(10).fork()).ldelim();
    }

    if (message.baselineMinimum.length !== 0) {
      writer.uint32(26).bytes(message.baselineMinimum);
    }

    if (message.baselineTvlPpm !== 0) {
      writer.uint32(32).uint32(message.baselineTvlPpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Limiter {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLimiter();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.period = Duration.decode(reader, reader.uint32());
          break;

        case 3:
          message.baselineMinimum = reader.bytes();
          break;

        case 4:
          message.baselineTvlPpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Limiter>): Limiter {
    const message = createBaseLimiter();
    message.period = object.period !== undefined && object.period !== null ? Duration.fromPartial(object.period) : undefined;
    message.baselineMinimum = object.baselineMinimum ?? new Uint8Array();
    message.baselineTvlPpm = object.baselineTvlPpm ?? 0;
    return message;
  }

};