"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ServiceClientImpl = exports.ProofOps = exports.ProofOp = exports.ABCIQueryResponse = exports.ABCIQueryRequest = exports.Module = exports.VersionInfo = exports.GetNodeInfoResponse = exports.GetNodeInfoRequest = exports.GetSyncingResponse = exports.GetSyncingRequest = exports.GetLatestBlockResponse = exports.GetLatestBlockRequest = exports.GetBlockByHeightResponse = exports.GetBlockByHeightRequest = exports.Validator = exports.GetLatestValidatorSetResponse = exports.GetLatestValidatorSetRequest = exports.GetValidatorSetByHeightResponse = exports.GetValidatorSetByHeightRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const pagination_1 = require("../../query/v1beta1/pagination");
const any_1 = require("../../../../google/protobuf/any");
const types_1 = require("../../../../tendermint/types/types");
const block_1 = require("../../../../tendermint/types/block");
const types_2 = require("./types");
const types_3 = require("../../../../tendermint/p2p/types");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.base.tendermint.v1beta1";
function createBaseGetValidatorSetByHeightRequest() {
    return {
        height: BigInt(0),
        pagination: undefined,
    };
}
exports.GetValidatorSetByHeightRequest = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetValidatorSetByHeightRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== BigInt(0)) {
            writer.uint32(8).int64(message.height);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetValidatorSetByHeightRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = reader.int64();
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
        const obj = createBaseGetValidatorSetByHeightRequest();
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetValidatorSetByHeightRequest();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseGetValidatorSetByHeightResponse() {
    return {
        blockHeight: BigInt(0),
        validators: [],
        pagination: undefined,
    };
}
exports.GetValidatorSetByHeightResponse = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetValidatorSetByHeightResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.blockHeight !== BigInt(0)) {
            writer.uint32(8).int64(message.blockHeight);
        }
        for (const v of message.validators) {
            exports.Validator.encode(v, writer.uint32(18).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetValidatorSetByHeightResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.blockHeight = reader.int64();
                    break;
                case 2:
                    message.validators.push(exports.Validator.decode(reader, reader.uint32()));
                    break;
                case 3:
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
        const obj = createBaseGetValidatorSetByHeightResponse();
        if ((0, helpers_1.isSet)(object.blockHeight))
            obj.blockHeight = BigInt(object.blockHeight.toString());
        if (Array.isArray(object?.validators))
            obj.validators = object.validators.map((e) => exports.Validator.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.blockHeight !== undefined && (obj.blockHeight = (message.blockHeight || BigInt(0)).toString());
        if (message.validators) {
            obj.validators = message.validators.map((e) => (e ? exports.Validator.toJSON(e) : undefined));
        }
        else {
            obj.validators = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetValidatorSetByHeightResponse();
        if (object.blockHeight !== undefined && object.blockHeight !== null) {
            message.blockHeight = BigInt(object.blockHeight.toString());
        }
        message.validators = object.validators?.map((e) => exports.Validator.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseGetLatestValidatorSetRequest() {
    return {
        pagination: undefined,
    };
}
exports.GetLatestValidatorSetRequest = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetLatestValidatorSetRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetLatestValidatorSetRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseGetLatestValidatorSetRequest();
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetLatestValidatorSetRequest();
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseGetLatestValidatorSetResponse() {
    return {
        blockHeight: BigInt(0),
        validators: [],
        pagination: undefined,
    };
}
exports.GetLatestValidatorSetResponse = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetLatestValidatorSetResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.blockHeight !== BigInt(0)) {
            writer.uint32(8).int64(message.blockHeight);
        }
        for (const v of message.validators) {
            exports.Validator.encode(v, writer.uint32(18).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetLatestValidatorSetResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.blockHeight = reader.int64();
                    break;
                case 2:
                    message.validators.push(exports.Validator.decode(reader, reader.uint32()));
                    break;
                case 3:
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
        const obj = createBaseGetLatestValidatorSetResponse();
        if ((0, helpers_1.isSet)(object.blockHeight))
            obj.blockHeight = BigInt(object.blockHeight.toString());
        if (Array.isArray(object?.validators))
            obj.validators = object.validators.map((e) => exports.Validator.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.blockHeight !== undefined && (obj.blockHeight = (message.blockHeight || BigInt(0)).toString());
        if (message.validators) {
            obj.validators = message.validators.map((e) => (e ? exports.Validator.toJSON(e) : undefined));
        }
        else {
            obj.validators = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetLatestValidatorSetResponse();
        if (object.blockHeight !== undefined && object.blockHeight !== null) {
            message.blockHeight = BigInt(object.blockHeight.toString());
        }
        message.validators = object.validators?.map((e) => exports.Validator.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseValidator() {
    return {
        address: "",
        pubKey: undefined,
        votingPower: BigInt(0),
        proposerPriority: BigInt(0),
    };
}
exports.Validator = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.Validator",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        if (message.pubKey !== undefined) {
            any_1.Any.encode(message.pubKey, writer.uint32(18).fork()).ldelim();
        }
        if (message.votingPower !== BigInt(0)) {
            writer.uint32(24).int64(message.votingPower);
        }
        if (message.proposerPriority !== BigInt(0)) {
            writer.uint32(32).int64(message.proposerPriority);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseValidator();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
                    message.pubKey = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.votingPower = reader.int64();
                    break;
                case 4:
                    message.proposerPriority = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseValidator();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.pubKey))
            obj.pubKey = any_1.Any.fromJSON(object.pubKey);
        if ((0, helpers_1.isSet)(object.votingPower))
            obj.votingPower = BigInt(object.votingPower.toString());
        if ((0, helpers_1.isSet)(object.proposerPriority))
            obj.proposerPriority = BigInt(object.proposerPriority.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        message.pubKey !== undefined && (obj.pubKey = message.pubKey ? any_1.Any.toJSON(message.pubKey) : undefined);
        message.votingPower !== undefined && (obj.votingPower = (message.votingPower || BigInt(0)).toString());
        message.proposerPriority !== undefined &&
            (obj.proposerPriority = (message.proposerPriority || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseValidator();
        message.address = object.address ?? "";
        if (object.pubKey !== undefined && object.pubKey !== null) {
            message.pubKey = any_1.Any.fromPartial(object.pubKey);
        }
        if (object.votingPower !== undefined && object.votingPower !== null) {
            message.votingPower = BigInt(object.votingPower.toString());
        }
        if (object.proposerPriority !== undefined && object.proposerPriority !== null) {
            message.proposerPriority = BigInt(object.proposerPriority.toString());
        }
        return message;
    },
};
function createBaseGetBlockByHeightRequest() {
    return {
        height: BigInt(0),
    };
}
exports.GetBlockByHeightRequest = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetBlockByHeightRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== BigInt(0)) {
            writer.uint32(8).int64(message.height);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetBlockByHeightRequest();
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
        const obj = createBaseGetBlockByHeightRequest();
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
        const message = createBaseGetBlockByHeightRequest();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        return message;
    },
};
function createBaseGetBlockByHeightResponse() {
    return {
        blockId: undefined,
        block: undefined,
        sdkBlock: undefined,
    };
}
exports.GetBlockByHeightResponse = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetBlockByHeightResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.blockId !== undefined) {
            types_1.BlockID.encode(message.blockId, writer.uint32(10).fork()).ldelim();
        }
        if (message.block !== undefined) {
            block_1.Block.encode(message.block, writer.uint32(18).fork()).ldelim();
        }
        if (message.sdkBlock !== undefined) {
            types_2.Block.encode(message.sdkBlock, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetBlockByHeightResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.blockId = types_1.BlockID.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.block = block_1.Block.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.sdkBlock = types_2.Block.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetBlockByHeightResponse();
        if ((0, helpers_1.isSet)(object.blockId))
            obj.blockId = types_1.BlockID.fromJSON(object.blockId);
        if ((0, helpers_1.isSet)(object.block))
            obj.block = block_1.Block.fromJSON(object.block);
        if ((0, helpers_1.isSet)(object.sdkBlock))
            obj.sdkBlock = types_2.Block.fromJSON(object.sdkBlock);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.blockId !== undefined &&
            (obj.blockId = message.blockId ? types_1.BlockID.toJSON(message.blockId) : undefined);
        message.block !== undefined && (obj.block = message.block ? block_1.Block.toJSON(message.block) : undefined);
        message.sdkBlock !== undefined &&
            (obj.sdkBlock = message.sdkBlock ? types_2.Block.toJSON(message.sdkBlock) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetBlockByHeightResponse();
        if (object.blockId !== undefined && object.blockId !== null) {
            message.blockId = types_1.BlockID.fromPartial(object.blockId);
        }
        if (object.block !== undefined && object.block !== null) {
            message.block = block_1.Block.fromPartial(object.block);
        }
        if (object.sdkBlock !== undefined && object.sdkBlock !== null) {
            message.sdkBlock = types_2.Block.fromPartial(object.sdkBlock);
        }
        return message;
    },
};
function createBaseGetLatestBlockRequest() {
    return {};
}
exports.GetLatestBlockRequest = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetLatestBlockRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetLatestBlockRequest();
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
        const obj = createBaseGetLatestBlockRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseGetLatestBlockRequest();
        return message;
    },
};
function createBaseGetLatestBlockResponse() {
    return {
        blockId: undefined,
        block: undefined,
        sdkBlock: undefined,
    };
}
exports.GetLatestBlockResponse = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetLatestBlockResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.blockId !== undefined) {
            types_1.BlockID.encode(message.blockId, writer.uint32(10).fork()).ldelim();
        }
        if (message.block !== undefined) {
            block_1.Block.encode(message.block, writer.uint32(18).fork()).ldelim();
        }
        if (message.sdkBlock !== undefined) {
            types_2.Block.encode(message.sdkBlock, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetLatestBlockResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.blockId = types_1.BlockID.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.block = block_1.Block.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.sdkBlock = types_2.Block.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetLatestBlockResponse();
        if ((0, helpers_1.isSet)(object.blockId))
            obj.blockId = types_1.BlockID.fromJSON(object.blockId);
        if ((0, helpers_1.isSet)(object.block))
            obj.block = block_1.Block.fromJSON(object.block);
        if ((0, helpers_1.isSet)(object.sdkBlock))
            obj.sdkBlock = types_2.Block.fromJSON(object.sdkBlock);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.blockId !== undefined &&
            (obj.blockId = message.blockId ? types_1.BlockID.toJSON(message.blockId) : undefined);
        message.block !== undefined && (obj.block = message.block ? block_1.Block.toJSON(message.block) : undefined);
        message.sdkBlock !== undefined &&
            (obj.sdkBlock = message.sdkBlock ? types_2.Block.toJSON(message.sdkBlock) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetLatestBlockResponse();
        if (object.blockId !== undefined && object.blockId !== null) {
            message.blockId = types_1.BlockID.fromPartial(object.blockId);
        }
        if (object.block !== undefined && object.block !== null) {
            message.block = block_1.Block.fromPartial(object.block);
        }
        if (object.sdkBlock !== undefined && object.sdkBlock !== null) {
            message.sdkBlock = types_2.Block.fromPartial(object.sdkBlock);
        }
        return message;
    },
};
function createBaseGetSyncingRequest() {
    return {};
}
exports.GetSyncingRequest = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetSyncingRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetSyncingRequest();
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
        const obj = createBaseGetSyncingRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseGetSyncingRequest();
        return message;
    },
};
function createBaseGetSyncingResponse() {
    return {
        syncing: false,
    };
}
exports.GetSyncingResponse = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetSyncingResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.syncing === true) {
            writer.uint32(8).bool(message.syncing);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetSyncingResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.syncing = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetSyncingResponse();
        if ((0, helpers_1.isSet)(object.syncing))
            obj.syncing = Boolean(object.syncing);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.syncing !== undefined && (obj.syncing = message.syncing);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetSyncingResponse();
        message.syncing = object.syncing ?? false;
        return message;
    },
};
function createBaseGetNodeInfoRequest() {
    return {};
}
exports.GetNodeInfoRequest = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetNodeInfoRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetNodeInfoRequest();
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
        const obj = createBaseGetNodeInfoRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseGetNodeInfoRequest();
        return message;
    },
};
function createBaseGetNodeInfoResponse() {
    return {
        defaultNodeInfo: undefined,
        applicationVersion: undefined,
    };
}
exports.GetNodeInfoResponse = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.GetNodeInfoResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.defaultNodeInfo !== undefined) {
            types_3.DefaultNodeInfo.encode(message.defaultNodeInfo, writer.uint32(10).fork()).ldelim();
        }
        if (message.applicationVersion !== undefined) {
            exports.VersionInfo.encode(message.applicationVersion, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetNodeInfoResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.defaultNodeInfo = types_3.DefaultNodeInfo.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.applicationVersion = exports.VersionInfo.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetNodeInfoResponse();
        if ((0, helpers_1.isSet)(object.defaultNodeInfo))
            obj.defaultNodeInfo = types_3.DefaultNodeInfo.fromJSON(object.defaultNodeInfo);
        if ((0, helpers_1.isSet)(object.applicationVersion))
            obj.applicationVersion = exports.VersionInfo.fromJSON(object.applicationVersion);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.defaultNodeInfo !== undefined &&
            (obj.defaultNodeInfo = message.defaultNodeInfo
                ? types_3.DefaultNodeInfo.toJSON(message.defaultNodeInfo)
                : undefined);
        message.applicationVersion !== undefined &&
            (obj.applicationVersion = message.applicationVersion
                ? exports.VersionInfo.toJSON(message.applicationVersion)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetNodeInfoResponse();
        if (object.defaultNodeInfo !== undefined && object.defaultNodeInfo !== null) {
            message.defaultNodeInfo = types_3.DefaultNodeInfo.fromPartial(object.defaultNodeInfo);
        }
        if (object.applicationVersion !== undefined && object.applicationVersion !== null) {
            message.applicationVersion = exports.VersionInfo.fromPartial(object.applicationVersion);
        }
        return message;
    },
};
function createBaseVersionInfo() {
    return {
        name: "",
        appName: "",
        version: "",
        gitCommit: "",
        buildTags: "",
        goVersion: "",
        buildDeps: [],
        cosmosSdkVersion: "",
    };
}
exports.VersionInfo = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.VersionInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        if (message.appName !== "") {
            writer.uint32(18).string(message.appName);
        }
        if (message.version !== "") {
            writer.uint32(26).string(message.version);
        }
        if (message.gitCommit !== "") {
            writer.uint32(34).string(message.gitCommit);
        }
        if (message.buildTags !== "") {
            writer.uint32(42).string(message.buildTags);
        }
        if (message.goVersion !== "") {
            writer.uint32(50).string(message.goVersion);
        }
        for (const v of message.buildDeps) {
            exports.Module.encode(v, writer.uint32(58).fork()).ldelim();
        }
        if (message.cosmosSdkVersion !== "") {
            writer.uint32(66).string(message.cosmosSdkVersion);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseVersionInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.appName = reader.string();
                    break;
                case 3:
                    message.version = reader.string();
                    break;
                case 4:
                    message.gitCommit = reader.string();
                    break;
                case 5:
                    message.buildTags = reader.string();
                    break;
                case 6:
                    message.goVersion = reader.string();
                    break;
                case 7:
                    message.buildDeps.push(exports.Module.decode(reader, reader.uint32()));
                    break;
                case 8:
                    message.cosmosSdkVersion = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseVersionInfo();
        if ((0, helpers_1.isSet)(object.name))
            obj.name = String(object.name);
        if ((0, helpers_1.isSet)(object.appName))
            obj.appName = String(object.appName);
        if ((0, helpers_1.isSet)(object.version))
            obj.version = String(object.version);
        if ((0, helpers_1.isSet)(object.gitCommit))
            obj.gitCommit = String(object.gitCommit);
        if ((0, helpers_1.isSet)(object.buildTags))
            obj.buildTags = String(object.buildTags);
        if ((0, helpers_1.isSet)(object.goVersion))
            obj.goVersion = String(object.goVersion);
        if (Array.isArray(object?.buildDeps))
            obj.buildDeps = object.buildDeps.map((e) => exports.Module.fromJSON(e));
        if ((0, helpers_1.isSet)(object.cosmosSdkVersion))
            obj.cosmosSdkVersion = String(object.cosmosSdkVersion);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.appName !== undefined && (obj.appName = message.appName);
        message.version !== undefined && (obj.version = message.version);
        message.gitCommit !== undefined && (obj.gitCommit = message.gitCommit);
        message.buildTags !== undefined && (obj.buildTags = message.buildTags);
        message.goVersion !== undefined && (obj.goVersion = message.goVersion);
        if (message.buildDeps) {
            obj.buildDeps = message.buildDeps.map((e) => (e ? exports.Module.toJSON(e) : undefined));
        }
        else {
            obj.buildDeps = [];
        }
        message.cosmosSdkVersion !== undefined && (obj.cosmosSdkVersion = message.cosmosSdkVersion);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseVersionInfo();
        message.name = object.name ?? "";
        message.appName = object.appName ?? "";
        message.version = object.version ?? "";
        message.gitCommit = object.gitCommit ?? "";
        message.buildTags = object.buildTags ?? "";
        message.goVersion = object.goVersion ?? "";
        message.buildDeps = object.buildDeps?.map((e) => exports.Module.fromPartial(e)) || [];
        message.cosmosSdkVersion = object.cosmosSdkVersion ?? "";
        return message;
    },
};
function createBaseModule() {
    return {
        path: "",
        version: "",
        sum: "",
    };
}
exports.Module = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.Module",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.path !== "") {
            writer.uint32(10).string(message.path);
        }
        if (message.version !== "") {
            writer.uint32(18).string(message.version);
        }
        if (message.sum !== "") {
            writer.uint32(26).string(message.sum);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModule();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.path = reader.string();
                    break;
                case 2:
                    message.version = reader.string();
                    break;
                case 3:
                    message.sum = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseModule();
        if ((0, helpers_1.isSet)(object.path))
            obj.path = String(object.path);
        if ((0, helpers_1.isSet)(object.version))
            obj.version = String(object.version);
        if ((0, helpers_1.isSet)(object.sum))
            obj.sum = String(object.sum);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.path !== undefined && (obj.path = message.path);
        message.version !== undefined && (obj.version = message.version);
        message.sum !== undefined && (obj.sum = message.sum);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModule();
        message.path = object.path ?? "";
        message.version = object.version ?? "";
        message.sum = object.sum ?? "";
        return message;
    },
};
function createBaseABCIQueryRequest() {
    return {
        data: new Uint8Array(),
        path: "",
        height: BigInt(0),
        prove: false,
    };
}
exports.ABCIQueryRequest = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.ABCIQueryRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.data.length !== 0) {
            writer.uint32(10).bytes(message.data);
        }
        if (message.path !== "") {
            writer.uint32(18).string(message.path);
        }
        if (message.height !== BigInt(0)) {
            writer.uint32(24).int64(message.height);
        }
        if (message.prove === true) {
            writer.uint32(32).bool(message.prove);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseABCIQueryRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.data = reader.bytes();
                    break;
                case 2:
                    message.path = reader.string();
                    break;
                case 3:
                    message.height = reader.int64();
                    break;
                case 4:
                    message.prove = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseABCIQueryRequest();
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        if ((0, helpers_1.isSet)(object.path))
            obj.path = String(object.path);
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.prove))
            obj.prove = Boolean(object.prove);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        message.path !== undefined && (obj.path = message.path);
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.prove !== undefined && (obj.prove = message.prove);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseABCIQueryRequest();
        message.data = object.data ?? new Uint8Array();
        message.path = object.path ?? "";
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        message.prove = object.prove ?? false;
        return message;
    },
};
function createBaseABCIQueryResponse() {
    return {
        code: 0,
        log: "",
        info: "",
        index: BigInt(0),
        key: new Uint8Array(),
        value: new Uint8Array(),
        proofOps: undefined,
        height: BigInt(0),
        codespace: "",
    };
}
exports.ABCIQueryResponse = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.ABCIQueryResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.code !== 0) {
            writer.uint32(8).uint32(message.code);
        }
        if (message.log !== "") {
            writer.uint32(26).string(message.log);
        }
        if (message.info !== "") {
            writer.uint32(34).string(message.info);
        }
        if (message.index !== BigInt(0)) {
            writer.uint32(40).int64(message.index);
        }
        if (message.key.length !== 0) {
            writer.uint32(50).bytes(message.key);
        }
        if (message.value.length !== 0) {
            writer.uint32(58).bytes(message.value);
        }
        if (message.proofOps !== undefined) {
            exports.ProofOps.encode(message.proofOps, writer.uint32(66).fork()).ldelim();
        }
        if (message.height !== BigInt(0)) {
            writer.uint32(72).int64(message.height);
        }
        if (message.codespace !== "") {
            writer.uint32(82).string(message.codespace);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseABCIQueryResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.code = reader.uint32();
                    break;
                case 3:
                    message.log = reader.string();
                    break;
                case 4:
                    message.info = reader.string();
                    break;
                case 5:
                    message.index = reader.int64();
                    break;
                case 6:
                    message.key = reader.bytes();
                    break;
                case 7:
                    message.value = reader.bytes();
                    break;
                case 8:
                    message.proofOps = exports.ProofOps.decode(reader, reader.uint32());
                    break;
                case 9:
                    message.height = reader.int64();
                    break;
                case 10:
                    message.codespace = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseABCIQueryResponse();
        if ((0, helpers_1.isSet)(object.code))
            obj.code = Number(object.code);
        if ((0, helpers_1.isSet)(object.log))
            obj.log = String(object.log);
        if ((0, helpers_1.isSet)(object.info))
            obj.info = String(object.info);
        if ((0, helpers_1.isSet)(object.index))
            obj.index = BigInt(object.index.toString());
        if ((0, helpers_1.isSet)(object.key))
            obj.key = (0, helpers_1.bytesFromBase64)(object.key);
        if ((0, helpers_1.isSet)(object.value))
            obj.value = (0, helpers_1.bytesFromBase64)(object.value);
        if ((0, helpers_1.isSet)(object.proofOps))
            obj.proofOps = exports.ProofOps.fromJSON(object.proofOps);
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.codespace))
            obj.codespace = String(object.codespace);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.code !== undefined && (obj.code = Math.round(message.code));
        message.log !== undefined && (obj.log = message.log);
        message.info !== undefined && (obj.info = message.info);
        message.index !== undefined && (obj.index = (message.index || BigInt(0)).toString());
        message.key !== undefined &&
            (obj.key = (0, helpers_1.base64FromBytes)(message.key !== undefined ? message.key : new Uint8Array()));
        message.value !== undefined &&
            (obj.value = (0, helpers_1.base64FromBytes)(message.value !== undefined ? message.value : new Uint8Array()));
        message.proofOps !== undefined &&
            (obj.proofOps = message.proofOps ? exports.ProofOps.toJSON(message.proofOps) : undefined);
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.codespace !== undefined && (obj.codespace = message.codespace);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseABCIQueryResponse();
        message.code = object.code ?? 0;
        message.log = object.log ?? "";
        message.info = object.info ?? "";
        if (object.index !== undefined && object.index !== null) {
            message.index = BigInt(object.index.toString());
        }
        message.key = object.key ?? new Uint8Array();
        message.value = object.value ?? new Uint8Array();
        if (object.proofOps !== undefined && object.proofOps !== null) {
            message.proofOps = exports.ProofOps.fromPartial(object.proofOps);
        }
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        message.codespace = object.codespace ?? "";
        return message;
    },
};
function createBaseProofOp() {
    return {
        type: "",
        key: new Uint8Array(),
        data: new Uint8Array(),
    };
}
exports.ProofOp = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.ProofOp",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.type !== "") {
            writer.uint32(10).string(message.type);
        }
        if (message.key.length !== 0) {
            writer.uint32(18).bytes(message.key);
        }
        if (message.data.length !== 0) {
            writer.uint32(26).bytes(message.data);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseProofOp();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.type = reader.string();
                    break;
                case 2:
                    message.key = reader.bytes();
                    break;
                case 3:
                    message.data = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseProofOp();
        if ((0, helpers_1.isSet)(object.type))
            obj.type = String(object.type);
        if ((0, helpers_1.isSet)(object.key))
            obj.key = (0, helpers_1.bytesFromBase64)(object.key);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.type !== undefined && (obj.type = message.type);
        message.key !== undefined &&
            (obj.key = (0, helpers_1.base64FromBytes)(message.key !== undefined ? message.key : new Uint8Array()));
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseProofOp();
        message.type = object.type ?? "";
        message.key = object.key ?? new Uint8Array();
        message.data = object.data ?? new Uint8Array();
        return message;
    },
};
function createBaseProofOps() {
    return {
        ops: [],
    };
}
exports.ProofOps = {
    typeUrl: "/cosmos.base.tendermint.v1beta1.ProofOps",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.ops) {
            exports.ProofOp.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseProofOps();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.ops.push(exports.ProofOp.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseProofOps();
        if (Array.isArray(object?.ops))
            obj.ops = object.ops.map((e) => exports.ProofOp.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.ops) {
            obj.ops = message.ops.map((e) => (e ? exports.ProofOp.toJSON(e) : undefined));
        }
        else {
            obj.ops = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseProofOps();
        message.ops = object.ops?.map((e) => exports.ProofOp.fromPartial(e)) || [];
        return message;
    },
};
class ServiceClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.GetNodeInfo = this.GetNodeInfo.bind(this);
        this.GetSyncing = this.GetSyncing.bind(this);
        this.GetLatestBlock = this.GetLatestBlock.bind(this);
        this.GetBlockByHeight = this.GetBlockByHeight.bind(this);
        this.GetLatestValidatorSet = this.GetLatestValidatorSet.bind(this);
        this.GetValidatorSetByHeight = this.GetValidatorSetByHeight.bind(this);
        this.ABCIQuery = this.ABCIQuery.bind(this);
    }
    GetNodeInfo(request = {}) {
        const data = exports.GetNodeInfoRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.base.tendermint.v1beta1.Service", "GetNodeInfo", data);
        return promise.then((data) => exports.GetNodeInfoResponse.decode(new binary_1.BinaryReader(data)));
    }
    GetSyncing(request = {}) {
        const data = exports.GetSyncingRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.base.tendermint.v1beta1.Service", "GetSyncing", data);
        return promise.then((data) => exports.GetSyncingResponse.decode(new binary_1.BinaryReader(data)));
    }
    GetLatestBlock(request = {}) {
        const data = exports.GetLatestBlockRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.base.tendermint.v1beta1.Service", "GetLatestBlock", data);
        return promise.then((data) => exports.GetLatestBlockResponse.decode(new binary_1.BinaryReader(data)));
    }
    GetBlockByHeight(request) {
        const data = exports.GetBlockByHeightRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.base.tendermint.v1beta1.Service", "GetBlockByHeight", data);
        return promise.then((data) => exports.GetBlockByHeightResponse.decode(new binary_1.BinaryReader(data)));
    }
    GetLatestValidatorSet(request = {
        pagination: pagination_1.PageRequest.fromPartial({}),
    }) {
        const data = exports.GetLatestValidatorSetRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.base.tendermint.v1beta1.Service", "GetLatestValidatorSet", data);
        return promise.then((data) => exports.GetLatestValidatorSetResponse.decode(new binary_1.BinaryReader(data)));
    }
    GetValidatorSetByHeight(request) {
        const data = exports.GetValidatorSetByHeightRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.base.tendermint.v1beta1.Service", "GetValidatorSetByHeight", data);
        return promise.then((data) => exports.GetValidatorSetByHeightResponse.decode(new binary_1.BinaryReader(data)));
    }
    ABCIQuery(request) {
        const data = exports.ABCIQueryRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.base.tendermint.v1beta1.Service", "ABCIQuery", data);
        return promise.then((data) => exports.ABCIQueryResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.ServiceClientImpl = ServiceClientImpl;
//# sourceMappingURL=query.js.map