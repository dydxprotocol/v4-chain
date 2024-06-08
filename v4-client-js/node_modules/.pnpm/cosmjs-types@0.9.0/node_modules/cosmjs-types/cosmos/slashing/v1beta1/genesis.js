"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MissedBlock = exports.ValidatorMissedBlocks = exports.SigningInfo = exports.GenesisState = exports.protobufPackage = void 0;
/* eslint-disable */
const slashing_1 = require("./slashing");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.slashing.v1beta1";
function createBaseGenesisState() {
    return {
        params: slashing_1.Params.fromPartial({}),
        signingInfos: [],
        missedBlocks: [],
    };
}
exports.GenesisState = {
    typeUrl: "/cosmos.slashing.v1beta1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.params !== undefined) {
            slashing_1.Params.encode(message.params, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.signingInfos) {
            exports.SigningInfo.encode(v, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.missedBlocks) {
            exports.ValidatorMissedBlocks.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGenesisState();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.params = slashing_1.Params.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.signingInfos.push(exports.SigningInfo.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.missedBlocks.push(exports.ValidatorMissedBlocks.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGenesisState();
        if ((0, helpers_1.isSet)(object.params))
            obj.params = slashing_1.Params.fromJSON(object.params);
        if (Array.isArray(object?.signingInfos))
            obj.signingInfos = object.signingInfos.map((e) => exports.SigningInfo.fromJSON(e));
        if (Array.isArray(object?.missedBlocks))
            obj.missedBlocks = object.missedBlocks.map((e) => exports.ValidatorMissedBlocks.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.params !== undefined && (obj.params = message.params ? slashing_1.Params.toJSON(message.params) : undefined);
        if (message.signingInfos) {
            obj.signingInfos = message.signingInfos.map((e) => (e ? exports.SigningInfo.toJSON(e) : undefined));
        }
        else {
            obj.signingInfos = [];
        }
        if (message.missedBlocks) {
            obj.missedBlocks = message.missedBlocks.map((e) => (e ? exports.ValidatorMissedBlocks.toJSON(e) : undefined));
        }
        else {
            obj.missedBlocks = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        if (object.params !== undefined && object.params !== null) {
            message.params = slashing_1.Params.fromPartial(object.params);
        }
        message.signingInfos = object.signingInfos?.map((e) => exports.SigningInfo.fromPartial(e)) || [];
        message.missedBlocks = object.missedBlocks?.map((e) => exports.ValidatorMissedBlocks.fromPartial(e)) || [];
        return message;
    },
};
function createBaseSigningInfo() {
    return {
        address: "",
        validatorSigningInfo: slashing_1.ValidatorSigningInfo.fromPartial({}),
    };
}
exports.SigningInfo = {
    typeUrl: "/cosmos.slashing.v1beta1.SigningInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        if (message.validatorSigningInfo !== undefined) {
            slashing_1.ValidatorSigningInfo.encode(message.validatorSigningInfo, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSigningInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
                    message.validatorSigningInfo = slashing_1.ValidatorSigningInfo.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSigningInfo();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.validatorSigningInfo))
            obj.validatorSigningInfo = slashing_1.ValidatorSigningInfo.fromJSON(object.validatorSigningInfo);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        message.validatorSigningInfo !== undefined &&
            (obj.validatorSigningInfo = message.validatorSigningInfo
                ? slashing_1.ValidatorSigningInfo.toJSON(message.validatorSigningInfo)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSigningInfo();
        message.address = object.address ?? "";
        if (object.validatorSigningInfo !== undefined && object.validatorSigningInfo !== null) {
            message.validatorSigningInfo = slashing_1.ValidatorSigningInfo.fromPartial(object.validatorSigningInfo);
        }
        return message;
    },
};
function createBaseValidatorMissedBlocks() {
    return {
        address: "",
        missedBlocks: [],
    };
}
exports.ValidatorMissedBlocks = {
    typeUrl: "/cosmos.slashing.v1beta1.ValidatorMissedBlocks",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        for (const v of message.missedBlocks) {
            exports.MissedBlock.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseValidatorMissedBlocks();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
                    message.missedBlocks.push(exports.MissedBlock.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseValidatorMissedBlocks();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if (Array.isArray(object?.missedBlocks))
            obj.missedBlocks = object.missedBlocks.map((e) => exports.MissedBlock.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        if (message.missedBlocks) {
            obj.missedBlocks = message.missedBlocks.map((e) => (e ? exports.MissedBlock.toJSON(e) : undefined));
        }
        else {
            obj.missedBlocks = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseValidatorMissedBlocks();
        message.address = object.address ?? "";
        message.missedBlocks = object.missedBlocks?.map((e) => exports.MissedBlock.fromPartial(e)) || [];
        return message;
    },
};
function createBaseMissedBlock() {
    return {
        index: BigInt(0),
        missed: false,
    };
}
exports.MissedBlock = {
    typeUrl: "/cosmos.slashing.v1beta1.MissedBlock",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.index !== BigInt(0)) {
            writer.uint32(8).int64(message.index);
        }
        if (message.missed === true) {
            writer.uint32(16).bool(message.missed);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMissedBlock();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.index = reader.int64();
                    break;
                case 2:
                    message.missed = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMissedBlock();
        if ((0, helpers_1.isSet)(object.index))
            obj.index = BigInt(object.index.toString());
        if ((0, helpers_1.isSet)(object.missed))
            obj.missed = Boolean(object.missed);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.index !== undefined && (obj.index = (message.index || BigInt(0)).toString());
        message.missed !== undefined && (obj.missed = message.missed);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMissedBlock();
        if (object.index !== undefined && object.index !== null) {
            message.index = BigInt(object.index.toString());
        }
        message.missed = object.missed ?? false;
        return message;
    },
};
//# sourceMappingURL=genesis.js.map