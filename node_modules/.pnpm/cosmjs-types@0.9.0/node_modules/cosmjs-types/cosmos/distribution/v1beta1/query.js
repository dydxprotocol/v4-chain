"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.QueryClientImpl = exports.QueryCommunityPoolResponse = exports.QueryCommunityPoolRequest = exports.QueryDelegatorWithdrawAddressResponse = exports.QueryDelegatorWithdrawAddressRequest = exports.QueryDelegatorValidatorsResponse = exports.QueryDelegatorValidatorsRequest = exports.QueryDelegationTotalRewardsResponse = exports.QueryDelegationTotalRewardsRequest = exports.QueryDelegationRewardsResponse = exports.QueryDelegationRewardsRequest = exports.QueryValidatorSlashesResponse = exports.QueryValidatorSlashesRequest = exports.QueryValidatorCommissionResponse = exports.QueryValidatorCommissionRequest = exports.QueryValidatorOutstandingRewardsResponse = exports.QueryValidatorOutstandingRewardsRequest = exports.QueryValidatorDistributionInfoResponse = exports.QueryValidatorDistributionInfoRequest = exports.QueryParamsResponse = exports.QueryParamsRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const pagination_1 = require("../../base/query/v1beta1/pagination");
const distribution_1 = require("./distribution");
const coin_1 = require("../../base/v1beta1/coin");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.distribution.v1beta1";
function createBaseQueryParamsRequest() {
    return {};
}
exports.QueryParamsRequest = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryParamsRequest",
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
        params: distribution_1.Params.fromPartial({}),
    };
}
exports.QueryParamsResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryParamsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.params !== undefined) {
            distribution_1.Params.encode(message.params, writer.uint32(10).fork()).ldelim();
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
                    message.params = distribution_1.Params.decode(reader, reader.uint32());
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
            obj.params = distribution_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.params !== undefined && (obj.params = message.params ? distribution_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryParamsResponse();
        if (object.params !== undefined && object.params !== null) {
            message.params = distribution_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
function createBaseQueryValidatorDistributionInfoRequest() {
    return {
        validatorAddress: "",
    };
}
exports.QueryValidatorDistributionInfoRequest = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryValidatorDistributionInfoRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validatorAddress !== "") {
            writer.uint32(10).string(message.validatorAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorDistributionInfoRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorDistributionInfoRequest();
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorDistributionInfoRequest();
        message.validatorAddress = object.validatorAddress ?? "";
        return message;
    },
};
function createBaseQueryValidatorDistributionInfoResponse() {
    return {
        operatorAddress: "",
        selfBondRewards: [],
        commission: [],
    };
}
exports.QueryValidatorDistributionInfoResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryValidatorDistributionInfoResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.operatorAddress !== "") {
            writer.uint32(10).string(message.operatorAddress);
        }
        for (const v of message.selfBondRewards) {
            coin_1.DecCoin.encode(v, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.commission) {
            coin_1.DecCoin.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorDistributionInfoResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.operatorAddress = reader.string();
                    break;
                case 2:
                    message.selfBondRewards.push(coin_1.DecCoin.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.commission.push(coin_1.DecCoin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorDistributionInfoResponse();
        if ((0, helpers_1.isSet)(object.operatorAddress))
            obj.operatorAddress = String(object.operatorAddress);
        if (Array.isArray(object?.selfBondRewards))
            obj.selfBondRewards = object.selfBondRewards.map((e) => coin_1.DecCoin.fromJSON(e));
        if (Array.isArray(object?.commission))
            obj.commission = object.commission.map((e) => coin_1.DecCoin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.operatorAddress !== undefined && (obj.operatorAddress = message.operatorAddress);
        if (message.selfBondRewards) {
            obj.selfBondRewards = message.selfBondRewards.map((e) => (e ? coin_1.DecCoin.toJSON(e) : undefined));
        }
        else {
            obj.selfBondRewards = [];
        }
        if (message.commission) {
            obj.commission = message.commission.map((e) => (e ? coin_1.DecCoin.toJSON(e) : undefined));
        }
        else {
            obj.commission = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorDistributionInfoResponse();
        message.operatorAddress = object.operatorAddress ?? "";
        message.selfBondRewards = object.selfBondRewards?.map((e) => coin_1.DecCoin.fromPartial(e)) || [];
        message.commission = object.commission?.map((e) => coin_1.DecCoin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseQueryValidatorOutstandingRewardsRequest() {
    return {
        validatorAddress: "",
    };
}
exports.QueryValidatorOutstandingRewardsRequest = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validatorAddress !== "") {
            writer.uint32(10).string(message.validatorAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorOutstandingRewardsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorOutstandingRewardsRequest();
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorOutstandingRewardsRequest();
        message.validatorAddress = object.validatorAddress ?? "";
        return message;
    },
};
function createBaseQueryValidatorOutstandingRewardsResponse() {
    return {
        rewards: distribution_1.ValidatorOutstandingRewards.fromPartial({}),
    };
}
exports.QueryValidatorOutstandingRewardsResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.rewards !== undefined) {
            distribution_1.ValidatorOutstandingRewards.encode(message.rewards, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorOutstandingRewardsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.rewards = distribution_1.ValidatorOutstandingRewards.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorOutstandingRewardsResponse();
        if ((0, helpers_1.isSet)(object.rewards))
            obj.rewards = distribution_1.ValidatorOutstandingRewards.fromJSON(object.rewards);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.rewards !== undefined &&
            (obj.rewards = message.rewards ? distribution_1.ValidatorOutstandingRewards.toJSON(message.rewards) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorOutstandingRewardsResponse();
        if (object.rewards !== undefined && object.rewards !== null) {
            message.rewards = distribution_1.ValidatorOutstandingRewards.fromPartial(object.rewards);
        }
        return message;
    },
};
function createBaseQueryValidatorCommissionRequest() {
    return {
        validatorAddress: "",
    };
}
exports.QueryValidatorCommissionRequest = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryValidatorCommissionRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validatorAddress !== "") {
            writer.uint32(10).string(message.validatorAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorCommissionRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorCommissionRequest();
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorCommissionRequest();
        message.validatorAddress = object.validatorAddress ?? "";
        return message;
    },
};
function createBaseQueryValidatorCommissionResponse() {
    return {
        commission: distribution_1.ValidatorAccumulatedCommission.fromPartial({}),
    };
}
exports.QueryValidatorCommissionResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryValidatorCommissionResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.commission !== undefined) {
            distribution_1.ValidatorAccumulatedCommission.encode(message.commission, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorCommissionResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.commission = distribution_1.ValidatorAccumulatedCommission.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryValidatorCommissionResponse();
        if ((0, helpers_1.isSet)(object.commission))
            obj.commission = distribution_1.ValidatorAccumulatedCommission.fromJSON(object.commission);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.commission !== undefined &&
            (obj.commission = message.commission
                ? distribution_1.ValidatorAccumulatedCommission.toJSON(message.commission)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorCommissionResponse();
        if (object.commission !== undefined && object.commission !== null) {
            message.commission = distribution_1.ValidatorAccumulatedCommission.fromPartial(object.commission);
        }
        return message;
    },
};
function createBaseQueryValidatorSlashesRequest() {
    return {
        validatorAddress: "",
        startingHeight: BigInt(0),
        endingHeight: BigInt(0),
        pagination: undefined,
    };
}
exports.QueryValidatorSlashesRequest = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryValidatorSlashesRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.validatorAddress !== "") {
            writer.uint32(10).string(message.validatorAddress);
        }
        if (message.startingHeight !== BigInt(0)) {
            writer.uint32(16).uint64(message.startingHeight);
        }
        if (message.endingHeight !== BigInt(0)) {
            writer.uint32(24).uint64(message.endingHeight);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorSlashesRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validatorAddress = reader.string();
                    break;
                case 2:
                    message.startingHeight = reader.uint64();
                    break;
                case 3:
                    message.endingHeight = reader.uint64();
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
        const obj = createBaseQueryValidatorSlashesRequest();
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        if ((0, helpers_1.isSet)(object.startingHeight))
            obj.startingHeight = BigInt(object.startingHeight.toString());
        if ((0, helpers_1.isSet)(object.endingHeight))
            obj.endingHeight = BigInt(object.endingHeight.toString());
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        message.startingHeight !== undefined &&
            (obj.startingHeight = (message.startingHeight || BigInt(0)).toString());
        message.endingHeight !== undefined && (obj.endingHeight = (message.endingHeight || BigInt(0)).toString());
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorSlashesRequest();
        message.validatorAddress = object.validatorAddress ?? "";
        if (object.startingHeight !== undefined && object.startingHeight !== null) {
            message.startingHeight = BigInt(object.startingHeight.toString());
        }
        if (object.endingHeight !== undefined && object.endingHeight !== null) {
            message.endingHeight = BigInt(object.endingHeight.toString());
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryValidatorSlashesResponse() {
    return {
        slashes: [],
        pagination: undefined,
    };
}
exports.QueryValidatorSlashesResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryValidatorSlashesResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.slashes) {
            distribution_1.ValidatorSlashEvent.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryValidatorSlashesResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.slashes.push(distribution_1.ValidatorSlashEvent.decode(reader, reader.uint32()));
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
        const obj = createBaseQueryValidatorSlashesResponse();
        if (Array.isArray(object?.slashes))
            obj.slashes = object.slashes.map((e) => distribution_1.ValidatorSlashEvent.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.slashes) {
            obj.slashes = message.slashes.map((e) => (e ? distribution_1.ValidatorSlashEvent.toJSON(e) : undefined));
        }
        else {
            obj.slashes = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryValidatorSlashesResponse();
        message.slashes = object.slashes?.map((e) => distribution_1.ValidatorSlashEvent.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryDelegationRewardsRequest() {
    return {
        delegatorAddress: "",
        validatorAddress: "",
    };
}
exports.QueryDelegationRewardsRequest = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryDelegationRewardsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        if (message.validatorAddress !== "") {
            writer.uint32(18).string(message.validatorAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegationRewardsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                case 2:
                    message.validatorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegationRewardsRequest();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        if ((0, helpers_1.isSet)(object.validatorAddress))
            obj.validatorAddress = String(object.validatorAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegationRewardsRequest();
        message.delegatorAddress = object.delegatorAddress ?? "";
        message.validatorAddress = object.validatorAddress ?? "";
        return message;
    },
};
function createBaseQueryDelegationRewardsResponse() {
    return {
        rewards: [],
    };
}
exports.QueryDelegationRewardsResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryDelegationRewardsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.rewards) {
            coin_1.DecCoin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegationRewardsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.rewards.push(coin_1.DecCoin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegationRewardsResponse();
        if (Array.isArray(object?.rewards))
            obj.rewards = object.rewards.map((e) => coin_1.DecCoin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.rewards) {
            obj.rewards = message.rewards.map((e) => (e ? coin_1.DecCoin.toJSON(e) : undefined));
        }
        else {
            obj.rewards = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegationRewardsResponse();
        message.rewards = object.rewards?.map((e) => coin_1.DecCoin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseQueryDelegationTotalRewardsRequest() {
    return {
        delegatorAddress: "",
    };
}
exports.QueryDelegationTotalRewardsRequest = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryDelegationTotalRewardsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegationTotalRewardsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegationTotalRewardsRequest();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegationTotalRewardsRequest();
        message.delegatorAddress = object.delegatorAddress ?? "";
        return message;
    },
};
function createBaseQueryDelegationTotalRewardsResponse() {
    return {
        rewards: [],
        total: [],
    };
}
exports.QueryDelegationTotalRewardsResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryDelegationTotalRewardsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.rewards) {
            distribution_1.DelegationDelegatorReward.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.total) {
            coin_1.DecCoin.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegationTotalRewardsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.rewards.push(distribution_1.DelegationDelegatorReward.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.total.push(coin_1.DecCoin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegationTotalRewardsResponse();
        if (Array.isArray(object?.rewards))
            obj.rewards = object.rewards.map((e) => distribution_1.DelegationDelegatorReward.fromJSON(e));
        if (Array.isArray(object?.total))
            obj.total = object.total.map((e) => coin_1.DecCoin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.rewards) {
            obj.rewards = message.rewards.map((e) => (e ? distribution_1.DelegationDelegatorReward.toJSON(e) : undefined));
        }
        else {
            obj.rewards = [];
        }
        if (message.total) {
            obj.total = message.total.map((e) => (e ? coin_1.DecCoin.toJSON(e) : undefined));
        }
        else {
            obj.total = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegationTotalRewardsResponse();
        message.rewards = object.rewards?.map((e) => distribution_1.DelegationDelegatorReward.fromPartial(e)) || [];
        message.total = object.total?.map((e) => coin_1.DecCoin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseQueryDelegatorValidatorsRequest() {
    return {
        delegatorAddress: "",
    };
}
exports.QueryDelegatorValidatorsRequest = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryDelegatorValidatorsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
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
                    message.delegatorAddress = reader.string();
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
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorValidatorsRequest();
        message.delegatorAddress = object.delegatorAddress ?? "";
        return message;
    },
};
function createBaseQueryDelegatorValidatorsResponse() {
    return {
        validators: [],
    };
}
exports.QueryDelegatorValidatorsResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryDelegatorValidatorsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.validators) {
            writer.uint32(10).string(v);
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
                    message.validators.push(reader.string());
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
            obj.validators = object.validators.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.validators) {
            obj.validators = message.validators.map((e) => e);
        }
        else {
            obj.validators = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorValidatorsResponse();
        message.validators = object.validators?.map((e) => e) || [];
        return message;
    },
};
function createBaseQueryDelegatorWithdrawAddressRequest() {
    return {
        delegatorAddress: "",
    };
}
exports.QueryDelegatorWithdrawAddressRequest = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.delegatorAddress !== "") {
            writer.uint32(10).string(message.delegatorAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegatorWithdrawAddressRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.delegatorAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegatorWithdrawAddressRequest();
        if ((0, helpers_1.isSet)(object.delegatorAddress))
            obj.delegatorAddress = String(object.delegatorAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.delegatorAddress !== undefined && (obj.delegatorAddress = message.delegatorAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorWithdrawAddressRequest();
        message.delegatorAddress = object.delegatorAddress ?? "";
        return message;
    },
};
function createBaseQueryDelegatorWithdrawAddressResponse() {
    return {
        withdrawAddress: "",
    };
}
exports.QueryDelegatorWithdrawAddressResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.withdrawAddress !== "") {
            writer.uint32(10).string(message.withdrawAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryDelegatorWithdrawAddressResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.withdrawAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryDelegatorWithdrawAddressResponse();
        if ((0, helpers_1.isSet)(object.withdrawAddress))
            obj.withdrawAddress = String(object.withdrawAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.withdrawAddress !== undefined && (obj.withdrawAddress = message.withdrawAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryDelegatorWithdrawAddressResponse();
        message.withdrawAddress = object.withdrawAddress ?? "";
        return message;
    },
};
function createBaseQueryCommunityPoolRequest() {
    return {};
}
exports.QueryCommunityPoolRequest = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryCommunityPoolRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryCommunityPoolRequest();
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
        const obj = createBaseQueryCommunityPoolRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseQueryCommunityPoolRequest();
        return message;
    },
};
function createBaseQueryCommunityPoolResponse() {
    return {
        pool: [],
    };
}
exports.QueryCommunityPoolResponse = {
    typeUrl: "/cosmos.distribution.v1beta1.QueryCommunityPoolResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.pool) {
            coin_1.DecCoin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryCommunityPoolResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pool.push(coin_1.DecCoin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryCommunityPoolResponse();
        if (Array.isArray(object?.pool))
            obj.pool = object.pool.map((e) => coin_1.DecCoin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.pool) {
            obj.pool = message.pool.map((e) => (e ? coin_1.DecCoin.toJSON(e) : undefined));
        }
        else {
            obj.pool = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryCommunityPoolResponse();
        message.pool = object.pool?.map((e) => coin_1.DecCoin.fromPartial(e)) || [];
        return message;
    },
};
class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.Params = this.Params.bind(this);
        this.ValidatorDistributionInfo = this.ValidatorDistributionInfo.bind(this);
        this.ValidatorOutstandingRewards = this.ValidatorOutstandingRewards.bind(this);
        this.ValidatorCommission = this.ValidatorCommission.bind(this);
        this.ValidatorSlashes = this.ValidatorSlashes.bind(this);
        this.DelegationRewards = this.DelegationRewards.bind(this);
        this.DelegationTotalRewards = this.DelegationTotalRewards.bind(this);
        this.DelegatorValidators = this.DelegatorValidators.bind(this);
        this.DelegatorWithdrawAddress = this.DelegatorWithdrawAddress.bind(this);
        this.CommunityPool = this.CommunityPool.bind(this);
    }
    Params(request = {}) {
        const data = exports.QueryParamsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Query", "Params", data);
        return promise.then((data) => exports.QueryParamsResponse.decode(new binary_1.BinaryReader(data)));
    }
    ValidatorDistributionInfo(request) {
        const data = exports.QueryValidatorDistributionInfoRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Query", "ValidatorDistributionInfo", data);
        return promise.then((data) => exports.QueryValidatorDistributionInfoResponse.decode(new binary_1.BinaryReader(data)));
    }
    ValidatorOutstandingRewards(request) {
        const data = exports.QueryValidatorOutstandingRewardsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Query", "ValidatorOutstandingRewards", data);
        return promise.then((data) => exports.QueryValidatorOutstandingRewardsResponse.decode(new binary_1.BinaryReader(data)));
    }
    ValidatorCommission(request) {
        const data = exports.QueryValidatorCommissionRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Query", "ValidatorCommission", data);
        return promise.then((data) => exports.QueryValidatorCommissionResponse.decode(new binary_1.BinaryReader(data)));
    }
    ValidatorSlashes(request) {
        const data = exports.QueryValidatorSlashesRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Query", "ValidatorSlashes", data);
        return promise.then((data) => exports.QueryValidatorSlashesResponse.decode(new binary_1.BinaryReader(data)));
    }
    DelegationRewards(request) {
        const data = exports.QueryDelegationRewardsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Query", "DelegationRewards", data);
        return promise.then((data) => exports.QueryDelegationRewardsResponse.decode(new binary_1.BinaryReader(data)));
    }
    DelegationTotalRewards(request) {
        const data = exports.QueryDelegationTotalRewardsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Query", "DelegationTotalRewards", data);
        return promise.then((data) => exports.QueryDelegationTotalRewardsResponse.decode(new binary_1.BinaryReader(data)));
    }
    DelegatorValidators(request) {
        const data = exports.QueryDelegatorValidatorsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Query", "DelegatorValidators", data);
        return promise.then((data) => exports.QueryDelegatorValidatorsResponse.decode(new binary_1.BinaryReader(data)));
    }
    DelegatorWithdrawAddress(request) {
        const data = exports.QueryDelegatorWithdrawAddressRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Query", "DelegatorWithdrawAddress", data);
        return promise.then((data) => exports.QueryDelegatorWithdrawAddressResponse.decode(new binary_1.BinaryReader(data)));
    }
    CommunityPool(request = {}) {
        const data = exports.QueryCommunityPoolRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.distribution.v1beta1.Query", "CommunityPool", data);
        return promise.then((data) => exports.QueryCommunityPoolResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.QueryClientImpl = QueryClientImpl;
//# sourceMappingURL=query.js.map