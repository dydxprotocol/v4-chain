"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CapabilityOwners = exports.Owner = exports.Capability = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.capability.v1beta1";
function createBaseCapability() {
    return {
        index: BigInt(0),
    };
}
exports.Capability = {
    typeUrl: "/cosmos.capability.v1beta1.Capability",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.index !== BigInt(0)) {
            writer.uint32(8).uint64(message.index);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCapability();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.index = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseCapability();
        if ((0, helpers_1.isSet)(object.index))
            obj.index = BigInt(object.index.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.index !== undefined && (obj.index = (message.index || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseCapability();
        if (object.index !== undefined && object.index !== null) {
            message.index = BigInt(object.index.toString());
        }
        return message;
    },
};
function createBaseOwner() {
    return {
        module: "",
        name: "",
    };
}
exports.Owner = {
    typeUrl: "/cosmos.capability.v1beta1.Owner",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.module !== "") {
            writer.uint32(10).string(message.module);
        }
        if (message.name !== "") {
            writer.uint32(18).string(message.name);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseOwner();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.module = reader.string();
                    break;
                case 2:
                    message.name = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseOwner();
        if ((0, helpers_1.isSet)(object.module))
            obj.module = String(object.module);
        if ((0, helpers_1.isSet)(object.name))
            obj.name = String(object.name);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.module !== undefined && (obj.module = message.module);
        message.name !== undefined && (obj.name = message.name);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseOwner();
        message.module = object.module ?? "";
        message.name = object.name ?? "";
        return message;
    },
};
function createBaseCapabilityOwners() {
    return {
        owners: [],
    };
}
exports.CapabilityOwners = {
    typeUrl: "/cosmos.capability.v1beta1.CapabilityOwners",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.owners) {
            exports.Owner.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCapabilityOwners();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.owners.push(exports.Owner.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseCapabilityOwners();
        if (Array.isArray(object?.owners))
            obj.owners = object.owners.map((e) => exports.Owner.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.owners) {
            obj.owners = message.owners.map((e) => (e ? exports.Owner.toJSON(e) : undefined));
        }
        else {
            obj.owners = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseCapabilityOwners();
        message.owners = object.owners?.map((e) => exports.Owner.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=capability.js.map