"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgUpdateParamsResponse = exports.MsgUpdateParams = exports.MsgCancelUnbondingDelegationResponse = exports.MsgCancelUnbondingDelegation = exports.MsgUndelegateResponse = exports.MsgUndelegate = exports.MsgBeginRedelegateResponse = exports.MsgBeginRedelegate = exports.MsgDelegateResponse = exports.MsgDelegate = exports.MsgEditValidatorResponse = exports.MsgEditValidator = exports.MsgCreateValidatorResponse = exports.MsgCreateValidator = exports.protobufPackage = void 0;
/* eslint-disable */
const staking_1 = require("./staking");
const any_1 = require("../../../google/protobuf/any");
const coin_1 = require("../../base/v1beta1/coin");
const timestamp_1 = require("../../../google/protobuf/timestamp");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.staking.v1beta1";
function createBaseMsgCreateValidator() {
    return {
        description: staking_1.Description.fromPartial({}),
        commission: staking_1.CommissionRates.fromPartial({}),
        minSelfDelegation: "",
        delegatorAddress: "",
        validatorAddress: "",
        pubkey: undefined,
        value: coin_1.Coin.fromPartial({}),
    };
}
exports.MsgCreateValidator = {
    typeUrl: "/cosmos.staking.v1beta1.MsgCreateValidator",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.description !== undefined) {
            staking_1.Description.encode(message.description, writer.uint32(10).fork()).ldelim();
        }
        if (message.commission !== undefined) {
            staking_1.CommissionRates.encode(message.commission, writer.uint32(18).fork()).ldelim();
        }
        if (message.minSelfDelegation !== "") {
            writer.uint32(26).string(message.minSelfDelegation);
        }
        if (message.delegatorAddress !== "") {
            writer.uint32(34).string(message.delegatorAddress);
        }
        if (message.validatorAddress !== "") {
            writer.uint32(42).string(message.validatorAddress);
        }
        if (message.pubkey !== undefined) {
            any_1.Any.encode(message.pubkey, writer.uint32(50).fork()).ldelim();
        }
        if (message.value !== undefined) {
            coin_1.Coin.encode(message.value, writer.uint32(58).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgCreateValidator();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.description = staking_1.Description.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.commission = staking_1.CommissionRates.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.minSelfDelegation = reader.string();
                    break;
                case 4:
                    message.delegatorAddress = reader.string();
                    break;
                case 5:
                    message.validatorAddress = reader.string();
                    break;
                case 6:
                    message.pubkey = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.value = coin_1.Coin.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgCreateValidator();
        if ((0, helpers_1.isSet)(object.description))
            obj.description = staking_1.Description.fromJSON(object.description);
        if ((0, helpers_1.isSet)(object.commission))
            obj.commission = staking_1.CommissionRates.fromJSON(object.commission);
        if ((0, helpers_1.isSet)(object.minSelfDelegation))
            obj.minSelfDelegation = String(object.minSelfDelegation);
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        if ((0, helpers_1.isSet)(object.pubkey))
            obj.pubkey = any_1.Any.fromJSON(object.pubkey);
        if ((0, helpers_1.isSet)(object.value))
            obj.value = coin_1.Coin.fromJSON(object.value);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.description !== undefined &&
            (obj.description = message.description ? staking_1.Description.toJSON(message.description) : undefined);
        message.commission !== undefined &&
            (obj.commission = message.commission ? staking_1.CommissionRates.toJSON(message.commission) : undefined);
        message.minSelfDelegation !== undefined && (obj.minSelfDelegation = message.minSelfDelegation);
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        message.pubkey !== undefined && (obj.pubkey = message.pubkey ? any_1.Any.toJSON(message.pubkey) : undefined);
        message.value !== undefined && (obj.value = message.value ? coin_1.Coin.toJSON(message.value) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgCreateValidator();
        if (object.description !== undefined && object.description !== null) {
            message.description = staking_1.Description.fromPartial(object.description);
        }
        if (object.commission !== undefined && object.commission !== null) {
            message.commission = staking_1.CommissionRates.fromPartial(object.commission);
        }
        message.minSelfDelegation = object.minSelfDelegation ?? "";
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorAddress = object.validatorAddress ?? "";
        if (object.pubkey !== undefined && object.pubkey !== null) {
            message.pubkey = any_1.Any.fromPartial(object.pubkey);
        }
        if (object.value !== undefined && object.value !== null) {
            message.value = coin_1.Coin.fromPartial(object.value);
        }
        return message;
    },
};
function createBaseMsgCreateValidatorResponse() {
    return {};
}
exports.MsgCreateValidatorResponse = {
    typeUrl: "/cosmos.staking.v1beta1.MsgCreateValidatorResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgCreateValidatorResponse();
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
    fromJSON(_) {
        const obj = createBaseMsgCreateValidatorResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgCreateValidatorResponse();
        return message;
    },
};
function createBaseMsgEditValidator() {
    return {
        description: staking_1.Description.fromPartial({}),
        validatorAddress: "",
        commissionRate: "",
        minSelfDelegation: "",
    };
}
exports.MsgEditValidator = {
    typeUrl: "/cosmos.staking.v1beta1.MsgEditValidator",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.description !== undefined) {
            staking_1.Description.encode(message.description, writer.uint32(10).fork()).ldelim();
        }
        if (message.validatorAddress !== "") {
            writer.uint32(18).string(message.validatorAddress);
        }
        if (message.commissionRate !== "") {
            writer.uint32(26).string(message.commissionRate);
        }
        if (message.minSelfDelegation !== "") {
            writer.uint32(34).string(message.minSelfDelegation);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgEditValidator();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.description = staking_1.Description.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.validatorAddress = reader.string();
                    break;
                case 3:
                    message.commissionRate = reader.string();
                    break;
                case 4:
                    message.minSelfDelegation = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgEditValidator();
        if ((0, helpers_1.isSet)(object.description))
            obj.description = staking_1.Description.fromJSON(object.description);
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        if ((0, helpers_1.isSet)(object.commissionRate))
            obj.commissionRate = String(object.commissionRate);
        if ((0, helpers_1.isSet)(object.minSelfDelegation))
            obj.minSelfDelegation = String(object.minSelfDelegation);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.description !== undefined &&
            (obj.description = message.description ? staking_1.Description.toJSON(message.description) : undefined);
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        message.commissionRate !== undefined && (obj.commissionRate = message.commissionRate);
        message.minSelfDelegation !== undefined && (obj.minSelfDelegation = message.minSelfDelegation);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgEditValidator();
        if (object.description !== undefined && object.description !== null) {
            message.description = staking_1.Description.fromPartial(object.description);
        }
        message.validatorAddress = object.validatorAddress ?? "";
        message.commissionRate = object.commissionRate ?? "";
        message.minSelfDelegation = object.minSelfDelegation ?? "";
        return message;
    },
};
function createBaseMsgEditValidatorResponse() {
    return {};
}
exports.MsgEditValidatorResponse = {
    typeUrl: "/cosmos.staking.v1beta1.MsgEditValidatorResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgEditValidatorResponse();
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
    fromJSON(_) {
        const obj = createBaseMsgEditValidatorResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgEditValidatorResponse();
        return message;
    },
};
function createBaseMsgDelegate() {
    return {
        delegatorAddress: "",
        validatorAddress: "",
        amount: coin_1.Coin.fromPartial({}),
    };
}
exports.MsgDelegate = {
    typeUrl: "/cosmos.staking.v1beta1.MsgDelegate",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        if (message.validatorAddress !== "") {
            writer.uint32(18).string(message.validatorAddress);
        }
        if (message.amount !== undefined) {
            coin_1.Coin.encode(message.amount, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgDelegate();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                case 2:
                    message.validatorAddress = reader.string();
                    break;
                case 3:
                    message.amount = coin_1.Coin.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgDelegate();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        if ((0, helpers_1.isSet)(object.amount))
            obj.amount = coin_1.Coin.fromJSON(object.amount);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        message.amount !== undefined && (obj.amount = message.amount ? coin_1.Coin.toJSON(message.amount) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgDelegate();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorAddress = object.validatorAddress ?? "";
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = coin_1.Coin.fromPartial(object.amount);
        }
        return message;
    },
};
function createBaseMsgDelegateResponse() {
    return {};
}
exports.MsgDelegateResponse = {
    typeUrl: "/cosmos.staking.v1beta1.MsgDelegateResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgDelegateResponse();
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
    fromJSON(_) {
        const obj = createBaseMsgDelegateResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgDelegateResponse();
        return message;
    },
};
function createBaseMsgBeginRedelegate() {
    return {
        delegatorAddress: "",
        validatorSrcAddress: "",
        validatorDstAddress: "",
        amount: coin_1.Coin.fromPartial({}),
    };
}
exports.MsgBeginRedelegate = {
    typeUrl: "/cosmos.staking.v1beta1.MsgBeginRedelegate",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        if (message.validatorSrcAddress !== "") {
            writer.uint32(18).string(message.validatorSrcAddress);
        }
        if (message.validatorDstAddress !== "") {
            writer.uint32(26).string(message.validatorDstAddress);
        }
        if (message.amount !== undefined) {
            coin_1.Coin.encode(message.amount, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgBeginRedelegate();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                case 2:
                    message.validatorSrcAddress = reader.string();
                    break;
                case 3:
                    message.validatorDstAddress = reader.string();
                    break;
                case 4:
                    message.amount = coin_1.Coin.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgBeginRedelegate();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.validatorSrcAddress))
            obj.validatorSrcAddress = String(object.validatorSrcAddress);
        if ((0, helpers_1.isSet)(object.validatorDstAddress))
            obj.validatorDstAddress = String(object.validatorDstAddress);
        if ((0, helpers_1.isSet)(object.amount))
            obj.amount = coin_1.Coin.fromJSON(object.amount);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.validatorSrcAddress !== undefined && (obj.validatorSrcAddress = message.validatorSrcAddress);
        message.validatorDstAddress !== undefined && (obj.validatorDstAddress = message.validatorDstAddress);
        message.amount !== undefined && (obj.amount = message.amount ? coin_1.Coin.toJSON(message.amount) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgBeginRedelegate();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorSrcAddress = object.validatorSrcAddress ?? "";
        message.validatorDstAddress = object.validatorDstAddress ?? "";
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = coin_1.Coin.fromPartial(object.amount);
        }
        return message;
    },
};
function createBaseMsgBeginRedelegateResponse() {
    return {
        completionTime: timestamp_1.Timestamp.fromPartial({}),
    };
}
exports.MsgBeginRedelegateResponse = {
    typeUrl: "/cosmos.staking.v1beta1.MsgBeginRedelegateResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.completionTime !== undefined) {
            timestamp_1.Timestamp.encode(message.completionTime, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgBeginRedelegateResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.completionTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgBeginRedelegateResponse();
        if ((0, helpers_1.isSet)(object.completionTime))
            obj.completionTime = (0, helpers_1.fromJsonTimestamp)(object.completionTime);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.completionTime !== undefined &&
            (obj.completionTime = (0, helpers_1.fromTimestamp)(message.completionTime).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgBeginRedelegateResponse();
        if (object.completionTime !== undefined && object.completionTime !== null) {
            message.completionTime = timestamp_1.Timestamp.fromPartial(object.completionTime);
        }
        return message;
    },
};
function createBaseMsgUndelegate() {
    return {
        delegatorAddress: "",
        validatorAddress: "",
        amount: coin_1.Coin.fromPartial({}),
    };
}
exports.MsgUndelegate = {
    typeUrl: "/cosmos.staking.v1beta1.MsgUndelegate",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        if (message.validatorAddress !== "") {
            writer.uint32(18).string(message.validatorAddress);
        }
        if (message.amount !== undefined) {
            coin_1.Coin.encode(message.amount, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUndelegate();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                case 2:
                    message.validatorAddress = reader.string();
                    break;
                case 3:
                    message.amount = coin_1.Coin.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgUndelegate();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        if ((0, helpers_1.isSet)(object.amount))
            obj.amount = coin_1.Coin.fromJSON(object.amount);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        message.amount !== undefined && (obj.amount = message.amount ? coin_1.Coin.toJSON(message.amount) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUndelegate();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorAddress = object.validatorAddress ?? "";
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = coin_1.Coin.fromPartial(object.amount);
        }
        return message;
    },
};
function createBaseMsgUndelegateResponse() {
    return {
        completionTime: timestamp_1.Timestamp.fromPartial({}),
    };
}
exports.MsgUndelegateResponse = {
    typeUrl: "/cosmos.staking.v1beta1.MsgUndelegateResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.completionTime !== undefined) {
            timestamp_1.Timestamp.encode(message.completionTime, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUndelegateResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.completionTime = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgUndelegateResponse();
        if ((0, helpers_1.isSet)(object.completionTime))
            obj.completionTime = (0, helpers_1.fromJsonTimestamp)(object.completionTime);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.completionTime !== undefined &&
            (obj.completionTime = (0, helpers_1.fromTimestamp)(message.completionTime).toISOString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUndelegateResponse();
        if (object.completionTime !== undefined && object.completionTime !== null) {
            message.completionTime = timestamp_1.Timestamp.fromPartial(object.completionTime);
        }
        return message;
    },
};
function createBaseMsgCancelUnbondingDelegation() {
    return {
        delegatorAddress: "",
        validatorAddress: "",
        amount: coin_1.Coin.fromPartial({}),
        creationHeight: BigInt(0),
    };
}
exports.MsgCancelUnbondingDelegation = {
    typeUrl: "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        if (message.validatorAddress !== "") {
            writer.uint32(18).string(message.validatorAddress);
        }
        if (message.amount !== undefined) {
            coin_1.Coin.encode(message.amount, writer.uint32(26).fork()).ldelim();
        }
        if (message.creationHeight !== BigInt(0)) {
            writer.uint32(32).int64(message.creationHeight);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgCancelUnbondingDelegation();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                case 2:
                    message.validatorAddress = reader.string();
                    break;
                case 3:
                    message.amount = coin_1.Coin.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.creationHeight = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgCancelUnbondingDelegation();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        if ((0, helpers_1.isSet)(object.amount))
            obj.amount = coin_1.Coin.fromJSON(object.amount);
        if ((0, helpers_1.isSet)(object.creationHeight))
            obj.creationHeight = BigInt(object.creationHeight.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        message.amount !== undefined && (obj.amount = message.amount ? coin_1.Coin.toJSON(message.amount) : undefined);
        message.creationHeight !== undefined &&
            (obj.creationHeight = (message.creationHeight || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgCancelUnbondingDelegation();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorAddress = object.validatorAddress ?? "";
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = coin_1.Coin.fromPartial(object.amount);
        }
        if (object.creationHeight !== undefined && object.creationHeight !== null) {
            message.creationHeight = BigInt(object.creationHeight.toString());
        }
        return message;
    },
};
function createBaseMsgCancelUnbondingDelegationResponse() {
    return {};
}
exports.MsgCancelUnbondingDelegationResponse = {
    typeUrl: "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegationResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgCancelUnbondingDelegationResponse();
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
    fromJSON(_) {
        const obj = createBaseMsgCancelUnbondingDelegationResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgCancelUnbondingDelegationResponse();
        return message;
    },
};
function createBaseMsgUpdateParams() {
    return {
        authority: "",
        params: staking_1.Params.fromPartial({}),
    };
}
exports.MsgUpdateParams = {
    typeUrl: "/cosmos.staking.v1beta1.MsgUpdateParams",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.authority !== "") {
            writer.uint32(10).string(message.authority);
        }
        if (message.params !== undefined) {
            staking_1.Params.encode(message.params, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpdateParams();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.authority = reader.string();
                    break;
                case 2:
                    message.params = staking_1.Params.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgUpdateParams();
        if ((0, helpers_1.isSet)(object.authority))
            obj.authority = String(object.authority);
        if ((0, helpers_1.isSet)(object.params))
            obj.params = staking_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.authority !== undefined && (obj.authority = message.authority);
        message.params !== undefined && (obj.params = message.params ? staking_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUpdateParams();
        message.authority = object.authority ?? "";
        if (object.params !== undefined && object.params !== null) {
            message.params = staking_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
function createBaseMsgUpdateParamsResponse() {
    return {};
}
exports.MsgUpdateParamsResponse = {
    typeUrl: "/cosmos.staking.v1beta1.MsgUpdateParamsResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpdateParamsResponse();
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
    fromJSON(_) {
        const obj = createBaseMsgUpdateParamsResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgUpdateParamsResponse();
        return message;
    },
};
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.CreateValidator = this.CreateValidator.bind(this);
        this.EditValidator = this.EditValidator.bind(this);
        this.Delegate = this.Delegate.bind(this);
        this.BeginRedelegate = this.BeginRedelegate.bind(this);
        this.Undelegate = this.Undelegate.bind(this);
        this.CancelUnbondingDelegation = this.CancelUnbondingDelegation.bind(this);
        this.UpdateParams = this.UpdateParams.bind(this);
    }
    CreateValidator(request) {
        const data = exports.MsgCreateValidator.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "CreateValidator", data);
        return promise.then((data) => exports.MsgCreateValidatorResponse.decode(new binary_1.BinaryReader(data)));
    }
    EditValidator(request) {
        const data = exports.MsgEditValidator.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "EditValidator", data);
        return promise.then((data) => exports.MsgEditValidatorResponse.decode(new binary_1.BinaryReader(data)));
    }
    Delegate(request) {
        const data = exports.MsgDelegate.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "Delegate", data);
        return promise.then((data) => exports.MsgDelegateResponse.decode(new binary_1.BinaryReader(data)));
    }
    BeginRedelegate(request) {
        const data = exports.MsgBeginRedelegate.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "BeginRedelegate", data);
        return promise.then((data) => exports.MsgBeginRedelegateResponse.decode(new binary_1.BinaryReader(data)));
    }
    Undelegate(request) {
        const data = exports.MsgUndelegate.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "Undelegate", data);
        return promise.then((data) => exports.MsgUndelegateResponse.decode(new binary_1.BinaryReader(data)));
    }
    CancelUnbondingDelegation(request) {
        const data = exports.MsgCancelUnbondingDelegation.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "CancelUnbondingDelegation", data);
        return promise.then((data) => exports.MsgCancelUnbondingDelegationResponse.decode(new binary_1.BinaryReader(data)));
    }
    UpdateParams(request) {
        const data = exports.MsgUpdateParams.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "UpdateParams", data);
        return promise.then((data) => exports.MsgUpdateParamsResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map