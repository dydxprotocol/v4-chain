"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenesisState = exports.GenesisOwners = exports.protobufPackage = void 0;
/* eslint-disable */
const capability_1 = require("./capability");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.capability.v1beta1";
function createBaseGenesisOwners() {
    return {
        index: BigInt(0),
        indexOwners: capability_1.CapabilityOwners.fromPartial({}),
    };
}
exports.GenesisOwners = {
    typeUrl: "/cosmos.capability.v1beta1.GenesisOwners",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.index !== BigInt(0)) {
            writer.uint32(8).uint64(message.index);
        }
        if (message.indexOwners !== undefined) {
            capability_1.CapabilityOwners.encode(message.indexOwners, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGenesisOwners();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.index = reader.uint64();
                    break;
                case 2:
                    message.indexOwners = capability_1.CapabilityOwners.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGenesisOwners();
        if ((0, helpers_1.isSet)(object.index))
            obj.index = BigInt(object.index.toString());
        if ((0, helpers_1.isSet)(object.indexOwners))
            obj.indexOwners = capability_1.CapabilityOwners.fromJSON(object.indexOwners);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.index !== undefined && (obj.index = (message.index || BigInt(0)).toString());
        message.indexOwners !== undefined &&
            (obj.indexOwners = message.indexOwners ? capability_1.CapabilityOwners.toJSON(message.indexOwners) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisOwners();
        if (object.index !== undefined && object.index !== null) {
            message.index = BigInt(object.index.toString());
        }
        if (object.indexOwners !== undefined && object.indexOwners !== null) {
            message.indexOwners = capability_1.CapabilityOwners.fromPartial(object.indexOwners);
        }
        return message;
    },
};
function createBaseGenesisState() {
    return {
        index: BigInt(0),
        owners: [],
    };
}
exports.GenesisState = {
    typeUrl: "/cosmos.capability.v1beta1.GenesisState",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.index !== BigInt(0)) {
            writer.uint32(8).uint64(message.index);
        }
        for (const v of message.owners) {
            exports.GenesisOwners.encode(v, writer.uint32(18).fork()).ldelim();
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
                    message.index = reader.uint64();
                    break;
                case 2:
                    message.owners.push(exports.GenesisOwners.decode(reader, reader.uint32()));
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
        if ((0, helpers_1.isSet)(object.index))
            obj.index = BigInt(object.index.toString());
        if (Array.isArray(object?.owners))
            obj.owners = object.owners.map((e) => exports.GenesisOwners.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.index !== undefined && (obj.index = (message.index || BigInt(0)).toString());
        if (message.owners) {
            obj.owners = message.owners.map((e) => (e ? exports.GenesisOwners.toJSON(e) : undefined));
        }
        else {
            obj.owners = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGenesisState();
        if (object.index !== undefined && object.index !== null) {
            message.index = BigInt(object.index.toString());
        }
        message.owners = object.owners?.map((e) => exports.GenesisOwners.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=genesis.js.map