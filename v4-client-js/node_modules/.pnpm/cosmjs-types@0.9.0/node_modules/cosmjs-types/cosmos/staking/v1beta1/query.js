"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.QueryClientImpl = exports.QueryParamsResponse = exports.QueryParamsRequest = exports.QueryPoolResponse = exports.QueryPoolRequest = exports.QueryHistoricalInfoResponse = exports.QueryHistoricalInfoRequest = exports.QueryDelegatorValidatorResponse = exports.QueryDelegatorValidatorRequest = exports.QueryDelegatorValidatorsResponse = exports.QueryDelegatorValidatorsRequest = exports.QueryRedelegationsResponse = exports.QueryRedelegationsRequest = exports.QueryDelegatorUnbondingDelegationsResponse = exports.QueryDelegatorUnbondingDelegationsRequest = exports.QueryDelegatorDelegationsResponse = exports.QueryDelegatorDelegationsRequest = exports.QueryUnbondingDelegationResponse = exports.QueryUnbondingDelegationRequest = exports.QueryDelegationResponse = exports.QueryDelegationRequest = exports.QueryValidatorUnbondingDelegationsResponse = exports.QueryValidatorUnbondingDelegationsRequest = exports.QueryValidatorDelegationsResponse = exports.QueryValidatorDelegationsRequest = exports.QueryValidatorResponse = exports.QueryValidatorRequest = exports.QueryValidatorsResponse = exports.QueryValidatorsRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const pagination_1 = require("../../base/query/v1beta1/pagination");
const staking_1 = require("./staking");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.staking.v1beta1";
function createBaseQueryValidatorsRequest() {
    return {
        status: "",
        pagination: undefined,
    };
}
exports.QueryValidatorsRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryValidatorsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.status !== "") {
            writer.uint32(10).string(message.status);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.status = reader.string();
                    break;
                case 2:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorsRequest();
        if ((0, helpers_1.isSet)(object.status))
            obj.status = String(object.status);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.status !== undefined && (obj.status = message.status);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorsRequest();
        message.status = object.status ?? "";
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryValidatorsResponse() {
    return {
        validators: [],
        pagination: undefined,
    };
}
exports.QueryValidatorsResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryValidatorsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.validators) {
            staking_1.Validator.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validators.push(staking_1.Validator.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorsResponse();
        if (Array.isArray(object?.validators))
            obj.validators = object.validators.map((e) => staking_1.Validator.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.validators) {
            obj.validators = message.validators.map((e) => (e ? staking_1.Validator.toJSON(e) : undefined));
        }
        else {
            obj.validators = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorsResponse();
        message.validators = object.validators?.map((e) => staking_1.Validator.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryValidatorRequest() {
    return {
        validatorAddr: "",
    };
}
exports.QueryValidatorRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryValidatorRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validatorAddr !== "") {
            writer.uint32(10).string(message.validatorAddr);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddr = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorRequest();
        if ((0, helpers_1.isSet)(object.validatorAddr))
            obj.validatorAddr = String(object.validatorAddr);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validatorAddr !== undefined && (obj.validatorAddr = message.validatorAddr);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorRequest();
        message.validatorAddr = object.validatorAddr ?? "";
        return message;
    },
};
function createBaseQueryValidatorResponse() {
    return {
        validator: staking_1.Validator.fromPartial({}),
    };
}
exports.QueryValidatorResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryValidatorResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validator !== undefined) {
            staking_1.Validator.encode(message.validator, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validator = staking_1.Validator.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorResponse();
        if ((0, helpers_1.isSet)(object.validator))
            obj.validator = staking_1.Validator.fromJSON(object.validator);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validator !== undefined &&
            (obj.validator = message.validator ? staking_1.Validator.toJSON(message.validator) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorResponse();
        if (object.validator !== undefined && object.validator !== null) {
            message.validator = staking_1.Validator.fromPartial(object.validator);
        }
        return message;
    },
};
function createBaseQueryValidatorDelegationsRequest() {
    return {
        validatorAddr: "",
        pagination: undefined,
    };
}
exports.QueryValidatorDelegationsRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryValidatorDelegationsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validatorAddr !== "") {
            writer.uint32(10).string(message.validatorAddr);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorDelegationsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddr = reader.string();
                    break;
                case 2:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorDelegationsRequest();
        if ((0, helpers_1.isSet)(object.validatorAddr))
            obj.validatorAddr = String(object.validatorAddr);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validatorAddr !== undefined && (obj.validatorAddr = message.validatorAddr);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorDelegationsRequest();
        message.validatorAddr = object.validatorAddr ?? "";
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryValidatorDelegationsResponse() {
    return {
        delegationResponses: [],
        pagination: undefined,
    };
}
exports.QueryValidatorDelegationsResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryValidatorDelegationsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.delegationResponses) {
            staking_1.DelegationResponse.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorDelegationsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegationResponses.push(staking_1.DelegationResponse.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorDelegationsResponse();
        if (Array.isArray(object?.delegationResponses))
            obj.delegationResponses = object.delegationResponses.map((e) => staking_1.DelegationResponse.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.delegationResponses) {
            obj.delegationResponses = message.delegationResponses.map((e) => e ? staking_1.DelegationResponse.toJSON(e) : undefined);
        }
        else {
            obj.delegationResponses = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorDelegationsResponse();
        message.delegationResponses =
            object.delegationResponses?.map((e) => staking_1.DelegationResponse.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryValidatorUnbondingDelegationsRequest() {
    return {
        validatorAddr: "",
        pagination: undefined,
    };
}
exports.QueryValidatorUnbondingDelegationsRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validatorAddr !== "") {
            writer.uint32(10).string(message.validatorAddr);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorUnbondingDelegationsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddr = reader.string();
                    break;
                case 2:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorUnbondingDelegationsRequest();
        if ((0, helpers_1.isSet)(object.validatorAddr))
            obj.validatorAddr = String(object.validatorAddr);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validatorAddr !== undefined && (obj.validatorAddr = message.validatorAddr);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorUnbondingDelegationsRequest();
        message.validatorAddr = object.validatorAddr ?? "";
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryValidatorUnbondingDelegationsResponse() {
    return {
        unbondingResponses: [],
        pagination: undefined,
    };
}
exports.QueryValidatorUnbondingDelegationsResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.unbondingResponses) {
            staking_1.UnbondingDelegation.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorUnbondingDelegationsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.unbondingResponses.push(staking_1.UnbondingDelegation.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorUnbondingDelegationsResponse();
        if (Array.isArray(object?.unbondingResponses))
            obj.unbondingResponses = object.unbondingResponses.map((e) => staking_1.UnbondingDelegation.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.unbondingResponses) {
            obj.unbondingResponses = message.unbondingResponses.map((e) => e ? staking_1.UnbondingDelegation.toJSON(e) : undefined);
        }
        else {
            obj.unbondingResponses = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorUnbondingDelegationsResponse();
        message.unbondingResponses =
            object.unbondingResponses?.map((e) => staking_1.UnbondingDelegation.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryDelegationRequest() {
    return {
        delegatorAddr: "",
        validatorAddr: "",
    };
}
exports.QueryDelegationRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryDelegationRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddr !== "") {
            writer.uint32(10).string(message.delegatorAddr);
        }
        if (message.validatorAddr !== "") {
            writer.uint32(18).string(message.validatorAddr);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegationRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddr = reader.string();
                    break;
                case 2:
                    message.validatorAddr = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegationRequest();
        if ((0, helpers_1.isSet)(object.delegatorAddr))
            obj.delegatorAddr = String(object.delegatorAddr);
        if ((0, helpers_1.isSet)(object.validatorAddr))
            obj.validatorAddr = String(object.validatorAddr);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddr !== undefined && (obj.delegatorAddr = message.delegatorAddr);
        message.validatorAddr !== undefined && (obj.validatorAddr = message.validatorAddr);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegationRequest();
        message.delegatorAddr = object.delegatorAddr ?? "";
        message.validatorAddr = object.validatorAddr ?? "";
        return message;
    },
};
function createBaseQueryDelegationResponse() {
    return {
        delegationResponse: undefined,
    };
}
exports.QueryDelegationResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryDelegationResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegationResponse !== undefined) {
            staking_1.DelegationResponse.encode(message.delegationResponse, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegationResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegationResponse = staking_1.DelegationResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegationResponse();
        if ((0, helpers_1.isSet)(object.delegationResponse))
            obj.delegationResponse = staking_1.DelegationResponse.fromJSON(object.delegationResponse);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegationResponse !== undefined &&
            (obj.delegationResponse = message.delegationResponse
                ? staking_1.DelegationResponse.toJSON(message.delegationResponse)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegationResponse();
        if (object.delegationResponse !== undefined && object.delegationResponse !== null) {
            message.delegationResponse = staking_1.DelegationResponse.fromPartial(object.delegationResponse);
        }
        return message;
    },
};
function createBaseQueryUnbondingDelegationRequest() {
    return {
        delegatorAddr: "",
        validatorAddr: "",
    };
}
exports.QueryUnbondingDelegationRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryUnbondingDelegationRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddr !== "") {
            writer.uint32(10).string(message.delegatorAddr);
        }
        if (message.validatorAddr !== "") {
            writer.uint32(18).string(message.validatorAddr);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryUnbondingDelegationRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddr = reader.string();
                    break;
                case 2:
                    message.validatorAddr = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryUnbondingDelegationRequest();
        if ((0, helpers_1.isSet)(object.delegatorAddr))
            obj.delegatorAddr = String(object.delegatorAddr);
        if ((0, helpers_1.isSet)(object.validatorAddr))
            obj.validatorAddr = String(object.validatorAddr);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddr !== undefined && (obj.delegatorAddr = message.delegatorAddr);
        message.validatorAddr !== undefined && (obj.validatorAddr = message.validatorAddr);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryUnbondingDelegationRequest();
        message.delegatorAddr = object.delegatorAddr ?? "";
        message.validatorAddr = object.validatorAddr ?? "";
        return message;
    },
};
function createBaseQueryUnbondingDelegationResponse() {
    return {
        unbond: staking_1.UnbondingDelegation.fromPartial({}),
    };
}
exports.QueryUnbondingDelegationResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryUnbondingDelegationResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.unbond !== undefined) {
            staking_1.UnbondingDelegation.encode(message.unbond, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryUnbondingDelegationResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.unbond = staking_1.UnbondingDelegation.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryUnbondingDelegationResponse();
        if ((0, helpers_1.isSet)(object.unbond))
            obj.unbond = staking_1.UnbondingDelegation.fromJSON(object.unbond);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.unbond !== undefined &&
            (obj.unbond = message.unbond ? staking_1.UnbondingDelegation.toJSON(message.unbond) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryUnbondingDelegationResponse();
        if (object.unbond !== undefined && object.unbond !== null) {
            message.unbond = staking_1.UnbondingDelegation.fromPartial(object.unbond);
        }
        return message;
    },
};
function createBaseQueryDelegatorDelegationsRequest() {
    return {
        delegatorAddr: "",
        pagination: undefined,
    };
}
exports.QueryDelegatorDelegationsRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryDelegatorDelegationsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddr !== "") {
            writer.uint32(10).string(message.delegatorAddr);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegatorDelegationsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddr = reader.string();
                    break;
                case 2:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegatorDelegationsRequest();
        if ((0, helpers_1.isSet)(object.delegatorAddr))
            obj.delegatorAddr = String(object.delegatorAddr);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddr !== undefined && (obj.delegatorAddr = message.delegatorAddr);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorDelegationsRequest();
        message.delegatorAddr = object.delegatorAddr ?? "";
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryDelegatorDelegationsResponse() {
    return {
        delegationResponses: [],
        pagination: undefined,
    };
}
exports.QueryDelegatorDelegationsResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryDelegatorDelegationsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.delegationResponses) {
            staking_1.DelegationResponse.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegatorDelegationsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegationResponses.push(staking_1.DelegationResponse.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegatorDelegationsResponse();
        if (Array.isArray(object?.delegationResponses))
            obj.delegationResponses = object.delegationResponses.map((e) => staking_1.DelegationResponse.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.delegationResponses) {
            obj.delegationResponses = message.delegationResponses.map((e) => e ? staking_1.DelegationResponse.toJSON(e) : undefined);
        }
        else {
            obj.delegationResponses = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorDelegationsResponse();
        message.delegationResponses =
            object.delegationResponses?.map((e) => staking_1.DelegationResponse.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryDelegatorUnbondingDelegationsRequest() {
    return {
        delegatorAddr: "",
        pagination: undefined,
    };
}
exports.QueryDelegatorUnbondingDelegationsRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddr !== "") {
            writer.uint32(10).string(message.delegatorAddr);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegatorUnbondingDelegationsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddr = reader.string();
                    break;
                case 2:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegatorUnbondingDelegationsRequest();
        if ((0, helpers_1.isSet)(object.delegatorAddr))
            obj.delegatorAddr = String(object.delegatorAddr);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddr !== undefined && (obj.delegatorAddr = message.delegatorAddr);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorUnbondingDelegationsRequest();
        message.delegatorAddr = object.delegatorAddr ?? "";
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryDelegatorUnbondingDelegationsResponse() {
    return {
        unbondingResponses: [],
        pagination: undefined,
    };
}
exports.QueryDelegatorUnbondingDelegationsResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.unbondingResponses) {
            staking_1.UnbondingDelegation.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegatorUnbondingDelegationsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.unbondingResponses.push(staking_1.UnbondingDelegation.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegatorUnbondingDelegationsResponse();
        if (Array.isArray(object?.unbondingResponses))
            obj.unbondingResponses = object.unbondingResponses.map((e) => staking_1.UnbondingDelegation.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.unbondingResponses) {
            obj.unbondingResponses = message.unbondingResponses.map((e) => e ? staking_1.UnbondingDelegation.toJSON(e) : undefined);
        }
        else {
            obj.unbondingResponses = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorUnbondingDelegationsResponse();
        message.unbondingResponses =
            object.unbondingResponses?.map((e) => staking_1.UnbondingDelegation.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryRedelegationsRequest() {
    return {
        delegatorAddr: "",
        srcValidatorAddr: "",
        dstValidatorAddr: "",
        pagination: undefined,
    };
}
exports.QueryRedelegationsRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryRedelegationsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddr !== "") {
            writer.uint32(10).string(message.delegatorAddr);
        }
        if (message.srcValidatorAddr !== "") {
            writer.uint32(18).string(message.srcValidatorAddr);
        }
        if (message.dstValidatorAddr !== "") {
            writer.uint32(26).string(message.dstValidatorAddr);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryRedelegationsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddr = reader.string();
                    break;
                case 2:
                    message.srcValidatorAddr = reader.string();
                    break;
                case 3:
                    message.dstValidatorAddr = reader.string();
                    break;
                case 4:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryRedelegationsRequest();
        if ((0, helpers_1.isSet)(object.delegatorAddr))
            obj.delegatorAddr = String(object.delegatorAddr);
        if ((0, helpers_1.isSet)(object.srcValidatorAddr))
            obj.srcValidatorAddr = String(object.srcValidatorAddr);
        if ((0, helpers_1.isSet)(object.dstValidatorAddr))
            obj.dstValidatorAddr = String(object.dstValidatorAddr);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddr !== undefined && (obj.delegatorAddr = message.delegatorAddr);
        message.srcValidatorAddr !== undefined && (obj.srcValidatorAddr = message.srcValidatorAddr);
        message.dstValidatorAddr !== undefined && (obj.dstValidatorAddr = message.dstValidatorAddr);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryRedelegationsRequest();
        message.delegatorAddr = object.delegatorAddr ?? "";
        message.srcValidatorAddr = object.srcValidatorAddr ?? "";
        message.dstValidatorAddr = object.dstValidatorAddr ?? "";
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryRedelegationsResponse() {
    return {
        redelegationResponses: [],
        pagination: undefined,
    };
}
exports.QueryRedelegationsResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryRedelegationsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.redelegationResponses) {
            staking_1.RedelegationResponse.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryRedelegationsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.redelegationResponses.push(staking_1.RedelegationResponse.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryRedelegationsResponse();
        if (Array.isArray(object?.redelegationResponses))
            obj.redelegationResponses = object.redelegationResponses.map((e) => staking_1.RedelegationResponse.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.redelegationResponses) {
            obj.redelegationResponses = message.redelegationResponses.map((e) => e ? staking_1.RedelegationResponse.toJSON(e) : undefined);
        }
        else {
            obj.redelegationResponses = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryRedelegationsResponse();
        message.redelegationResponses =
            object.redelegationResponses?.map((e) => staking_1.RedelegationResponse.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryDelegatorValidatorsRequest() {
    return {
        delegatorAddr: "",
        pagination: undefined,
    };
}
exports.QueryDelegatorValidatorsRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryDelegatorValidatorsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddr !== "") {
            writer.uint32(10).string(message.delegatorAddr);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegatorValidatorsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddr = reader.string();
                    break;
                case 2:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegatorValidatorsRequest();
        if ((0, helpers_1.isSet)(object.delegatorAddr))
            obj.delegatorAddr = String(object.delegatorAddr);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddr !== undefined && (obj.delegatorAddr = message.delegatorAddr);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorValidatorsRequest();
        message.delegatorAddr = object.delegatorAddr ?? "";
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryDelegatorValidatorsResponse() {
    return {
        validators: [],
        pagination: undefined,
    };
}
exports.QueryDelegatorValidatorsResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryDelegatorValidatorsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.validators) {
            staking_1.Validator.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegatorValidatorsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validators.push(staking_1.Validator.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegatorValidatorsResponse();
        if (Array.isArray(object?.validators))
            obj.validators = object.validators.map((e) => staking_1.Validator.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.validators) {
            obj.validators = message.validators.map((e) => (e ? staking_1.Validator.toJSON(e) : undefined));
        }
        else {
            obj.validators = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorValidatorsResponse();
        message.validators = object.validators?.map((e) => staking_1.Validator.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryDelegatorValidatorRequest() {
    return {
        delegatorAddr: "",
        validatorAddr: "",
    };
}
exports.QueryDelegatorValidatorRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryDelegatorValidatorRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddr !== "") {
            writer.uint32(10).string(message.delegatorAddr);
        }
        if (message.validatorAddr !== "") {
            writer.uint32(18).string(message.validatorAddr);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegatorValidatorRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddr = reader.string();
                    break;
                case 2:
                    message.validatorAddr = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegatorValidatorRequest();
        if ((0, helpers_1.isSet)(object.delegatorAddr))
            obj.delegatorAddr = String(object.delegatorAddr);
        if ((0, helpers_1.isSet)(object.validatorAddr))
            obj.validatorAddr = String(object.validatorAddr);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddr !== undefined && (obj.delegatorAddr = message.delegatorAddr);
        message.validatorAddr !== undefined && (obj.validatorAddr = message.validatorAddr);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorValidatorRequest();
        message.delegatorAddr = object.delegatorAddr ?? "";
        message.validatorAddr = object.validatorAddr ?? "";
        return message;
    },
};
function createBaseQueryDelegatorValidatorResponse() {
    return {
        validator: staking_1.Validator.fromPartial({}),
    };
}
exports.QueryDelegatorValidatorResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryDelegatorValidatorResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validator !== undefined) {
            staking_1.Validator.encode(message.validator, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegatorValidatorResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validator = staking_1.Validator.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegatorValidatorResponse();
        if ((0, helpers_1.isSet)(object.validator))
            obj.validator = staking_1.Validator.fromJSON(object.validator);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validator !== undefined &&
            (obj.validator = message.validator ? staking_1.Validator.toJSON(message.validator) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorValidatorResponse();
        if (object.validator !== undefined && object.validator !== null) {
            message.validator = staking_1.Validator.fromPartial(object.validator);
        }
        return message;
    },
};
function createBaseQueryHistoricalInfoRequest() {
    return {
        height: BigInt(0),
    };
}
exports.QueryHistoricalInfoRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryHistoricalInfoRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== BigInt(0)) {
            writer.uint32(8).int64(message.height);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryHistoricalInfoRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryHistoricalInfoRequest();
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryHistoricalInfoRequest();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        return message;
    },
};
function createBaseQueryHistoricalInfoResponse() {
    return {
        hist: undefined,
    };
}
exports.QueryHistoricalInfoResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryHistoricalInfoResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.hist !== undefined) {
            staking_1.HistoricalInfo.encode(message.hist, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryHistoricalInfoResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.hist = staking_1.HistoricalInfo.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryHistoricalInfoResponse();
        if ((0, helpers_1.isSet)(object.hist))
            obj.hist = staking_1.HistoricalInfo.fromJSON(object.hist);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.hist !== undefined && (obj.hist = message.hist ? staking_1.HistoricalInfo.toJSON(message.hist) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryHistoricalInfoResponse();
        if (object.hist !== undefined && object.hist !== null) {
            message.hist = staking_1.HistoricalInfo.fromPartial(object.hist);
        }
        return message;
    },
};
function createBaseQueryPoolRequest() {
    return {};
}
exports.QueryPoolRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryPoolRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryPoolRequest();
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
        const obj = createBaseQueryPoolRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseQueryPoolRequest();
        return message;
    },
};
function createBaseQueryPoolResponse() {
    return {
        pool: staking_1.Pool.fromPartial({}),
    };
}
exports.QueryPoolResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryPoolResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.pool !== undefined) {
            staking_1.Pool.encode(message.pool, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryPoolResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pool = staking_1.Pool.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryPoolResponse();
        if ((0, helpers_1.isSet)(object.pool))
            obj.pool = staking_1.Pool.fromJSON(object.pool);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.pool !== undefined && (obj.pool = message.pool ? staking_1.Pool.toJSON(message.pool) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryPoolResponse();
        if (object.pool !== undefined && object.pool !== null) {
            message.pool = staking_1.Pool.fromPartial(object.pool);
        }
        return message;
    },
};
function createBaseQueryParamsRequest() {
    return {};
}
exports.QueryParamsRequest = {
    typeUrl: "/cosmos.staking.v1beta1.QueryParamsRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryParamsRequest();
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
        const obj = createBaseQueryParamsRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseQueryParamsRequest();
        return message;
    },
};
function createBaseQueryParamsResponse() {
    return {
        params: staking_1.Params.fromPartial({}),
    };
}
exports.QueryParamsResponse = {
    typeUrl: "/cosmos.staking.v1beta1.QueryParamsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.params !== undefined) {
            staking_1.Params.encode(message.params, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryParamsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseQueryParamsResponse();
        if ((0, helpers_1.isSet)(object.params))
            obj.params = staking_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.params !== undefined && (obj.params = message.params ? staking_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryParamsResponse();
        if (object.params !== undefined && object.params !== null) {
            message.params = staking_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.Validators = this.Validators.bind(this);
        this.Validator = this.Validator.bind(this);
        this.ValidatorDelegations = this.ValidatorDelegations.bind(this);
        this.ValidatorUnbondingDelegations = this.ValidatorUnbondingDelegations.bind(this);
        this.Delegation = this.Delegation.bind(this);
        this.UnbondingDelegation = this.UnbondingDelegation.bind(this);
        this.DelegatorDelegations = this.DelegatorDelegations.bind(this);
        this.DelegatorUnbondingDelegations = this.DelegatorUnbondingDelegations.bind(this);
        this.Redelegations = this.Redelegations.bind(this);
        this.DelegatorValidators = this.DelegatorValidators.bind(this);
        this.DelegatorValidator = this.DelegatorValidator.bind(this);
        this.HistoricalInfo = this.HistoricalInfo.bind(this);
        this.Pool = this.Pool.bind(this);
        this.Params = this.Params.bind(this);
    }
    Validators(request) {
        const data = exports.QueryValidatorsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "Validators", data);
        return promise.then((data) => exports.QueryValidatorsResponse.decode(new binary_1.BinaryReader(data)));
    }
    Validator(request) {
        const data = exports.QueryValidatorRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "Validator", data);
        return promise.then((data) => exports.QueryValidatorResponse.decode(new binary_1.BinaryReader(data)));
    }
    ValidatorDelegations(request) {
        const data = exports.QueryValidatorDelegationsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "ValidatorDelegations", data);
        return promise.then((data) => exports.QueryValidatorDelegationsResponse.decode(new binary_1.BinaryReader(data)));
    }
    ValidatorUnbondingDelegations(request) {
        const data = exports.QueryValidatorUnbondingDelegationsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "ValidatorUnbondingDelegations", data);
        return promise.then((data) => exports.QueryValidatorUnbondingDelegationsResponse.decode(new binary_1.BinaryReader(data)));
    }
    Delegation(request) {
        const data = exports.QueryDelegationRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "Delegation", data);
        return promise.then((data) => exports.QueryDelegationResponse.decode(new binary_1.BinaryReader(data)));
    }
    UnbondingDelegation(request) {
        const data = exports.QueryUnbondingDelegationRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "UnbondingDelegation", data);
        return promise.then((data) => exports.QueryUnbondingDelegationResponse.decode(new binary_1.BinaryReader(data)));
    }
    DelegatorDelegations(request) {
        const data = exports.QueryDelegatorDelegationsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "DelegatorDelegations", data);
        return promise.then((data) => exports.QueryDelegatorDelegationsResponse.decode(new binary_1.BinaryReader(data)));
    }
    DelegatorUnbondingDelegations(request) {
        const data = exports.QueryDelegatorUnbondingDelegationsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "DelegatorUnbondingDelegations", data);
        return promise.then((data) => exports.QueryDelegatorUnbondingDelegationsResponse.decode(new binary_1.BinaryReader(data)));
    }
    Redelegations(request) {
        const data = exports.QueryRedelegationsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "Redelegations", data);
        return promise.then((data) => exports.QueryRedelegationsResponse.decode(new binary_1.BinaryReader(data)));
    }
    DelegatorValidators(request) {
        const data = exports.QueryDelegatorValidatorsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "DelegatorValidators", data);
        return promise.then((data) => exports.QueryDelegatorValidatorsResponse.decode(new binary_1.BinaryReader(data)));
    }
    DelegatorValidator(request) {
        const data = exports.QueryDelegatorValidatorRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "DelegatorValidator", data);
        return promise.then((data) => exports.QueryDelegatorValidatorResponse.decode(new binary_1.BinaryReader(data)));
    }
    HistoricalInfo(request) {
        const data = exports.QueryHistoricalInfoRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "HistoricalInfo", data);
        return promise.then((data) => exports.QueryHistoricalInfoResponse.decode(new binary_1.BinaryReader(data)));
    }
    Pool(request = {}) {
        const data = exports.QueryPoolRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "Pool", data);
        return promise.then((data) => exports.QueryPoolResponse.decode(new binary_1.BinaryReader(data)));
    }
    Params(request = {}) {
        const data = exports.QueryParamsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Query", "Params", data);
        return promise.then((data) => exports.QueryParamsResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.QueryClientImpl = QueryClientImpl;
//# sourceMappingURL=query.js.map