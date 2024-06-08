"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SimpleValidator = exports.Validator = exports.ValidatorSet = exports.protobufPackage = void 0;
/* eslint-disable */
const keys_1 = require("../crypto/keys");
const binary_1 = require("../../binary");
const helpers_1 = require("../../helpers");
exports.protobufPackage = "tendermint.types";
function createBaseValidatorSet() {
    return {
        validators: [],
        proposer: undefined,
        totalVotingPower: BigInt(0),
    };
}
exports.ValidatorSet = {
    typeUrl: "/tendermint.types.ValidatorSet",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.validators) {
            exports.Validator.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.proposer !== undefined) {
            exports.Validator.encode(message.proposer, writer.uint32(18).fork()).ldelim();
        }
        if (message.totalVotingPower !== BigInt(0)) {
            writer.uint32(24).int64(message.totalVotingPower);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseValidatorSet();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.validators.push(exports.Validator.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.proposer = exports.Validator.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.totalVotingPower = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseValidatorSet();
        if (Array.isArray(object?.validators))
            obj.validators = object.validators.map((e) => exports.Validator.fromJSON(e));
        if ((0, helpers_1.isSet)(object.proposer))
            obj.proposer = exports.Validator.fromJSON(object.proposer);
        if ((0, helpers_1.isSet)(object.totalVotingPower))
            obj.totalVotingPower = BigInt(object.totalVotingPower.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.validators) {
            obj.validators = message.validators.map((e) => (e ? exports.Validator.toJSON(e) : undefined));
        }
        else {
            obj.validators = [];
        }
        message.proposer !== undefined &&
            (obj.proposer = message.proposer ? exports.Validator.toJSON(message.proposer) : undefined);
        message.totalVotingPower !== undefined &&
            (obj.totalVotingPower = (message.totalVotingPower || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseValidatorSet();
        message.validators = object.validators?.map((e) => exports.Validator.fromPartial(e)) || [];
        if (object.proposer !== undefined && object.proposer !== null) {
            message.proposer = exports.Validator.fromPartial(object.proposer);
        }
        if (object.totalVotingPower !== undefined && object.totalVotingPower !== null) {
            message.totalVotingPower = BigInt(object.totalVotingPower.toString());
        }
        return message;
    },
};
function createBaseValidator() {
    return {
        address: new Uint8Array(),
        pubKey: keys_1.PublicKey.fromPartial({}),
        votingPower: BigInt(0),
        proposerPriority: BigInt(0),
    };
}
exports.Validator = {
    typeUrl: "/tendermint.types.Validator",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address.length !== 0) {
            writer.uint32(10).bytes(message.address);
        }
        if (message.pubKey !== undefined) {
            keys_1.PublicKey.encode(message.pubKey, writer.uint32(18).fork()).ldelim();
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
                    message.address = reader.bytes();
                    break;
                case 2:
                    message.pubKey = keys_1.PublicKey.decode(reader, reader.uint32());
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
            obj.address = (0, helpers_1.bytesFromBase64)(object.address);
        if ((0, helpers_1.isSet)(object.pubKey))
            obj.pubKey = keys_1.PublicKey.fromJSON(object.pubKey);
        if ((0, helpers_1.isSet)(object.votingPower))
            obj.votingPower = BigInt(object.votingPower.toString());
        if ((0, helpers_1.isSet)(object.proposerPriority))
            obj.proposerPriority = BigInt(object.proposerPriority.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined &&
            (obj.address = (0, helpers_1.base64FromBytes)(message.address !== undefined ? message.address : new Uint8Array()));
        message.pubKey !== undefined &&
            (obj.pubKey = message.pubKey ? keys_1.PublicKey.toJSON(message.pubKey) : undefined);
        message.votingPower !== undefined && (obj.votingPower = (message.votingPower || BigInt(0)).toString());
        message.proposerPriority !== undefined &&
            (obj.proposerPriority = (message.proposerPriority || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseValidator();
        message.address = object.address ?? new Uint8Array();
        if (object.pubKey !== undefined && object.pubKey !== null) {
            message.pubKey = keys_1.PublicKey.fromPartial(object.pubKey);
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
function createBaseSimpleValidator() {
    return {
        pubKey: undefined,
        votingPower: BigInt(0),
    };
}
exports.SimpleValidator = {
    typeUrl: "/tendermint.types.SimpleValidator",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.pubKey !== undefined) {
            keys_1.PublicKey.encode(message.pubKey, writer.uint32(10).fork()).ldelim();
        }
        if (message.votingPower !== BigInt(0)) {
            writer.uint32(16).int64(message.votingPower);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSimpleValidator();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pubKey = keys_1.PublicKey.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.votingPower = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSimpleValidator();
        if ((0, helpers_1.isSet)(object.pubKey))
            obj.pubKey = keys_1.PublicKey.fromJSON(object.pubKey);
        if ((0, helpers_1.isSet)(object.votingPower))
            obj.votingPower = BigInt(object.votingPower.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.pubKey !== undefined &&
            (obj.pubKey = message.pubKey ? keys_1.PublicKey.toJSON(message.pubKey) : undefined);
        message.votingPower !== undefined && (obj.votingPower = (message.votingPower || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSimpleValidator();
        if (object.pubKey !== undefined && object.pubKey !== null) {
            message.pubKey = keys_1.PublicKey.fromPartial(object.pubKey);
        }
        if (object.votingPower !== undefined && object.votingPower !== null) {
            message.votingPower = BigInt(object.votingPower.toString());
        }
        return message;
    },
};
//# sourceMappingURL=validator.js.map