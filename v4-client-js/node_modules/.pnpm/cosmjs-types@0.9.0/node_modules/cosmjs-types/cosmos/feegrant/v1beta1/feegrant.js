"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Grant = exports.AllowedMsgAllowance = exports.PeriodicAllowance = exports.BasicAllowance = exports.protobufPackage = void 0;
/* eslint-disable */
const coin_1 = require("../../base/v1beta1/coin");
const timestamp_1 = require("../../../google/protobuf/timestamp");
const duration_1 = require("../../../google/protobuf/duration");
const any_1 = require("../../../google/protobuf/any");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.feegrant.v1beta1";
function createBaseBasicAllowance() {
    return {
        spendLimit: [],
        expiration: undefined,
    };
}
exports.BasicAllowance = {
    typeUrl: "/cosmos.feegrant.v1beta1.BasicAllowance",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.spendLimit) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.expiration !== undefined) {
            timestamp_1.Timestamp.encode(message.expiration, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseBasicAllowance();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.spendLimit.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.expiration = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseBasicAllowance();
        if (Array.isArray(object?.spendLimit))
            obj.spendLimit = object.spendLimit.map((e) => coin_1.Coin.fromJSON(e));
        if ((0, helpers_1.isSet)(object.expiration))
            obj.expiration = (0, helpers_1.fromJsonTimestamp)(object.expiration);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.spendLimit) {
            obj.spendLimit = message.spendLimit.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.spendLimit = [];
        }
        message.expiration !== undefined && (obj.expiration = (0, helpers_1.fromTimestamp)(message.expiration).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseBasicAllowance();
        message.spendLimit = object.spendLimit?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        if (object.expiration !== undefined && object.expiration !== null) {
            message.expiration = timestamp_1.Timestamp.fromPartial(object.expiration);
        }
        return message;
    },
};
function createBasePeriodicAllowance() {
    return {
        basic: exports.BasicAllowance.fromPartial({}),
        period: duration_1.Duration.fromPartial({}),
        periodSpendLimit: [],
        periodCanSpend: [],
        periodReset: timestamp_1.Timestamp.fromPartial({}),
    };
}
exports.PeriodicAllowance = {
    typeUrl: "/cosmos.feegrant.v1beta1.PeriodicAllowance",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.basic !== undefined) {
            exports.BasicAllowance.encode(message.basic, writer.uint32(10).fork()).ldelim();
        }
        if (message.period !== undefined) {
            duration_1.Duration.encode(message.period, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.periodSpendLimit) {
            coin_1.Coin.encode(v, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.periodCanSpend) {
            coin_1.Coin.encode(v, writer.uint32(34).fork()).ldelim();
        }
        if (message.periodReset !== undefined) {
            timestamp_1.Timestamp.encode(message.periodReset, writer.uint32(42).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePeriodicAllowance();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.basic = exports.BasicAllowance.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.period = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.periodSpendLimit.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.periodCanSpend.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 5:
                    message.periodReset = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePeriodicAllowance();
        if ((0, helpers_1.isSet)(object.basic))
            obj.basic = exports.BasicAllowance.fromJSON(object.basic);
        if ((0, helpers_1.isSet)(object.period))
            obj.period = duration_1.Duration.fromJSON(object.period);
        if (Array.isArray(object?.periodSpendLimit))
            obj.periodSpendLimit = object.periodSpendLimit.map((e) => coin_1.Coin.fromJSON(e));
        if (Array.isArray(object?.periodCanSpend))
            obj.periodCanSpend = object.periodCanSpend.map((e) => coin_1.Coin.fromJSON(e));
        if ((0, helpers_1.isSet)(object.periodReset))
            obj.periodReset = (0, helpers_1.fromJsonTimestamp)(object.periodReset);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.basic !== undefined &&
            (obj.basic = message.basic ? exports.BasicAllowance.toJSON(message.basic) : undefined);
        message.period !== undefined &&
            (obj.period = message.period ? duration_1.Duration.toJSON(message.period) : undefined);
        if (message.periodSpendLimit) {
            obj.periodSpendLimit = message.periodSpendLimit.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.periodSpendLimit = [];
        }
        if (message.periodCanSpend) {
            obj.periodCanSpend = message.periodCanSpend.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.periodCanSpend = [];
        }
        message.periodReset !== undefined && (obj.periodReset = (0, helpers_1.fromTimestamp)(message.periodReset).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBasePeriodicAllowance();
        if (object.basic !== undefined && object.basic !== null) {
            message.basic = exports.BasicAllowance.fromPartial(object.basic);
        }
        if (object.period !== undefined && object.period !== null) {
            message.period = duration_1.Duration.fromPartial(object.period);
        }
        message.periodSpendLimit = object.periodSpendLimit?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        message.periodCanSpend = object.periodCanSpend?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        if (object.periodReset !== undefined && object.periodReset !== null) {
            message.periodReset = timestamp_1.Timestamp.fromPartial(object.periodReset);
        }
        return message;
    },
};
function createBaseAllowedMsgAllowance() {
    return {
        allowance: undefined,
        allowedMessages: [],
    };
}
exports.AllowedMsgAllowance = {
    typeUrl: "/cosmos.feegrant.v1beta1.AllowedMsgAllowance",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.allowance !== undefined) {
            any_1.Any.encode(message.allowance, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.allowedMessages) {
            writer.uint32(18).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseAllowedMsgAllowance();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.allowance = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.allowedMessages.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseAllowedMsgAllowance();
        if ((0, helpers_1.isSet)(object.allowance))
            obj.allowance = any_1.Any.fromJSON(object.allowance);
        if (Array.isArray(object?.allowedMessages))
            obj.allowedMessages = object.allowedMessages.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.allowance !== undefined &&
            (obj.allowance = message.allowance ? any_1.Any.toJSON(message.allowance) : undefined);
        if (message.allowedMessages) {
            obj.allowedMessages = message.allowedMessages.map((e) => e);
        }
        else {
            obj.allowedMessages = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseAllowedMsgAllowance();
        if (object.allowance !== undefined && object.allowance !== null) {
            message.allowance = any_1.Any.fromPartial(object.allowance);
        }
        message.allowedMessages = object.allowedMessages?.map((e) => e) || [];
        return message;
    },
};
function createBaseGrant() {
    return {
        granter: "",
        grantee: "",
        allowance: undefined,
    };
}
exports.Grant = {
    typeUrl: "/cosmos.feegrant.v1beta1.Grant",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.granter !== "") {
            writer.uint32(10).string(message.granter);
        }
        if (message.grantee !== "") {
            writer.uint32(18).string(message.grantee);
        }
        if (message.allowance !== undefined) {
            any_1.Any.encode(message.allowance, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGrant();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.granter = reader.string();
                    break;
                case 2:
                    message.grantee = reader.string();
                    break;
                case 3:
                    message.allowance = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGrant();
        if ((0, helpers_1.isSet)(object.granter))
            obj.granter = String(object.granter);
        if ((0, helpers_1.isSet)(object.grantee))
            obj.grantee = String(object.grantee);
        if ((0, helpers_1.isSet)(object.allowance))
            obj.allowance = any_1.Any.fromJSON(object.allowance);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.granter !== undefined && (obj.granter = message.granter);
        message.grantee !== undefined && (obj.grantee = message.grantee);
        message.allowance !== undefined &&
            (obj.allowance = message.allowance ? any_1.Any.toJSON(message.allowance) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGrant();
        message.granter = object.granter ?? "";
        message.grantee = object.grantee ?? "";
        if (object.allowance !== undefined && object.allowance !== null) {
            message.allowance = any_1.Any.fromPartial(object.allowance);
        }
        return message;
    },
};
//# sourceMappingURL=feegrant.js.map