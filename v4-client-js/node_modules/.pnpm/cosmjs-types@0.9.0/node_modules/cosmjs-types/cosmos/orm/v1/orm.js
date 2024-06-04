"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SingletonDescriptor = exports.SecondaryIndexDescriptor = exports.PrimaryKeyDescriptor = exports.TableDescriptor = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.orm.v1";
function createBaseTableDescriptor() {
    return {
        primaryKey: undefined,
        index: [],
        id: 0,
    };
}
exports.TableDescriptor = {
    typeUrl: "/cosmos.orm.v1.TableDescriptor",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.primaryKey !== undefined) {
            exports.PrimaryKeyDescriptor.encode(message.primaryKey, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.index) {
            exports.SecondaryIndexDescriptor.encode(v, writer.uint32(18).fork()).ldelim();
        }
        if (message.id !== 0) {
            writer.uint32(24).uint32(message.id);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTableDescriptor();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.primaryKey = exports.PrimaryKeyDescriptor.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.index.push(exports.SecondaryIndexDescriptor.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.id = reader.uint32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTableDescriptor();
        if ((0, helpers_1.isSet)(object.primaryKey))
            obj.primaryKey = exports.PrimaryKeyDescriptor.fromJSON(object.primaryKey);
        if (Array.isArray(object?.index))
            obj.index = object.index.map((e) => exports.SecondaryIndexDescriptor.fromJSON(e));
        if ((0, helpers_1.isSet)(object.id))
            obj.id = Number(object.id);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.primaryKey !== undefined &&
            (obj.primaryKey = message.primaryKey ? exports.PrimaryKeyDescriptor.toJSON(message.primaryKey) : undefined);
        if (message.index) {
            obj.index = message.index.map((e) => (e ? exports.SecondaryIndexDescriptor.toJSON(e) : undefined));
        }
        else {
            obj.index = [];
        }
        message.id !== undefined && (obj.id = Math.round(message.id));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTableDescriptor();
        if (object.primaryKey !== undefined && object.primaryKey !== null) {
            message.primaryKey = exports.PrimaryKeyDescriptor.fromPartial(object.primaryKey);
        }
        message.index = object.index?.map((e) => exports.SecondaryIndexDescriptor.fromPartial(e)) || [];
        message.id = object.id ?? 0;
        return message;
    },
};
function createBasePrimaryKeyDescriptor() {
    return {
        fields: "",
        autoIncrement: false,
    };
}
exports.PrimaryKeyDescriptor = {
    typeUrl: "/cosmos.orm.v1.PrimaryKeyDescriptor",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.fields !== "") {
            writer.uint32(10).string(message.fields);
        }
        if (message.autoIncrement === true) {
            writer.uint32(16).bool(message.autoIncrement);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePrimaryKeyDescriptor();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.fields = reader.string();
                    break;
                case 2:
                    message.autoIncrement = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBasePrimaryKeyDescriptor();
        if ((0, helpers_1.isSet)(object.fields))
            obj.fields = String(object.fields);
        if ((0, helpers_1.isSet)(object.autoIncrement))
            obj.autoIncrement = Boolean(object.autoIncrement);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.fields !== undefined && (obj.fields = message.fields);
        message.autoIncrement !== undefined && (obj.autoIncrement = message.autoIncrement);
        return obj;
    },
    fromPartial(object) {
        const message = createBasePrimaryKeyDescriptor();
        message.fields = object.fields ?? "";
        message.autoIncrement = object.autoIncrement ?? false;
        return message;
    },
};
function createBaseSecondaryIndexDescriptor() {
    return {
        fields: "",
        id: 0,
        unique: false,
    };
}
exports.SecondaryIndexDescriptor = {
    typeUrl: "/cosmos.orm.v1.SecondaryIndexDescriptor",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.fields !== "") {
            writer.uint32(10).string(message.fields);
        }
        if (message.id !== 0) {
            writer.uint32(16).uint32(message.id);
        }
        if (message.unique === true) {
            writer.uint32(24).bool(message.unique);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSecondaryIndexDescriptor();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.fields = reader.string();
                    break;
                case 2:
                    message.id = reader.uint32();
                    break;
                case 3:
                    message.unique = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSecondaryIndexDescriptor();
        if ((0, helpers_1.isSet)(object.fields))
            obj.fields = String(object.fields);
        if ((0, helpers_1.isSet)(object.id))
            obj.id = Number(object.id);
        if ((0, helpers_1.isSet)(object.unique))
            obj.unique = Boolean(object.unique);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.fields !== undefined && (obj.fields = message.fields);
        message.id !== undefined && (obj.id = Math.round(message.id));
        message.unique !== undefined && (obj.unique = message.unique);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSecondaryIndexDescriptor();
        message.fields = object.fields ?? "";
        message.id = object.id ?? 0;
        message.unique = object.unique ?? false;
        return message;
    },
};
function createBaseSingletonDescriptor() {
    return {
        id: 0,
    };
}
exports.SingletonDescriptor = {
    typeUrl: "/cosmos.orm.v1.SingletonDescriptor",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.id !== 0) {
            writer.uint32(8).uint32(message.id);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSingletonDescriptor();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.uint32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSingletonDescriptor();
        if ((0, helpers_1.isSet)(object.id))
            obj.id = Number(object.id);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.id !== undefined && (obj.id = Math.round(message.id));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSingletonDescriptor();
        message.id = object.id ?? 0;
        return message;
    },
};
//# sourceMappingURL=orm.js.map