"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Record_Offline = exports.Record_Multi = exports.Record_Ledger = exports.Record_Local = exports.Record = exports.protobufPackage = void 0;
/* eslint-disable */
const any_1 = require("../../../../google/protobuf/any");
const hd_1 = require("../../hd/v1/hd");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.crypto.keyring.v1";
function createBaseRecord() {
    return {
        name: "",
        pubKey: undefined,
        local: undefined,
        ledger: undefined,
        multi: undefined,
        offline: undefined,
    };
}
exports.Record = {
    typeUrl: "/cosmos.crypto.keyring.v1.Record",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        if (message.pubKey !== undefined) {
            any_1.Any.encode(message.pubKey, writer.uint32(18).fork()).ldelim();
        }
        if (message.local !== undefined) {
            exports.Record_Local.encode(message.local, writer.uint32(26).fork()).ldelim();
        }
        if (message.ledger !== undefined) {
            exports.Record_Ledger.encode(message.ledger, writer.uint32(34).fork()).ldelim();
        }
        if (message.multi !== undefined) {
            exports.Record_Multi.encode(message.multi, writer.uint32(42).fork()).ldelim();
        }
        if (message.offline !== undefined) {
            exports.Record_Offline.encode(message.offline, writer.uint32(50).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRecord();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.pubKey = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.local = exports.Record_Local.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.ledger = exports.Record_Ledger.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.multi = exports.Record_Multi.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.offline = exports.Record_Offline.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRecord();
        if ((0, helpers_1.isSet)(object.name))
            obj.name = String(object.name);
        if ((0, helpers_1.isSet)(object.pubKey))
            obj.pubKey = any_1.Any.fromJSON(object.pubKey);
        if ((0, helpers_1.isSet)(object.local))
            obj.local = exports.Record_Local.fromJSON(object.local);
        if ((0, helpers_1.isSet)(object.ledger))
            obj.ledger = exports.Record_Ledger.fromJSON(object.ledger);
        if ((0, helpers_1.isSet)(object.multi))
            obj.multi = exports.Record_Multi.fromJSON(object.multi);
        if ((0, helpers_1.isSet)(object.offline))
            obj.offline = exports.Record_Offline.fromJSON(object.offline);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.pubKey !== undefined && (obj.pubKey = message.pubKey ? any_1.Any.toJSON(message.pubKey) : undefined);
        message.local !== undefined &&
            (obj.local = message.local ? exports.Record_Local.toJSON(message.local) : undefined);
        message.ledger !== undefined &&
            (obj.ledger = message.ledger ? exports.Record_Ledger.toJSON(message.ledger) : undefined);
        message.multi !== undefined &&
            (obj.multi = message.multi ? exports.Record_Multi.toJSON(message.multi) : undefined);
        message.offline !== undefined &&
            (obj.offline = message.offline ? exports.Record_Offline.toJSON(message.offline) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRecord();
        message.name = object.name ?? "";
        if (object.pubKey !== undefined && object.pubKey !== null) {
            message.pubKey = any_1.Any.fromPartial(object.pubKey);
        }
        if (object.local !== undefined && object.local !== null) {
            message.local = exports.Record_Local.fromPartial(object.local);
        }
        if (object.ledger !== undefined && object.ledger !== null) {
            message.ledger = exports.Record_Ledger.fromPartial(object.ledger);
        }
        if (object.multi !== undefined && object.multi !== null) {
            message.multi = exports.Record_Multi.fromPartial(object.multi);
        }
        if (object.offline !== undefined && object.offline !== null) {
            message.offline = exports.Record_Offline.fromPartial(object.offline);
        }
        return message;
    },
};
function createBaseRecord_Local() {
    return {
        privKey: undefined,
    };
}
exports.Record_Local = {
    typeUrl: "/cosmos.crypto.keyring.v1.Local",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.privKey !== undefined) {
            any_1.Any.encode(message.privKey, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRecord_Local();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.privKey = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRecord_Local();
        if ((0, helpers_1.isSet)(object.privKey))
            obj.privKey = any_1.Any.fromJSON(object.privKey);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.privKey !== undefined &&
            (obj.privKey = message.privKey ? any_1.Any.toJSON(message.privKey) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRecord_Local();
        if (object.privKey !== undefined && object.privKey !== null) {
            message.privKey = any_1.Any.fromPartial(object.privKey);
        }
        return message;
    },
};
function createBaseRecord_Ledger() {
    return {
        path: undefined,
    };
}
exports.Record_Ledger = {
    typeUrl: "/cosmos.crypto.keyring.v1.Ledger",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.path !== undefined) {
            hd_1.BIP44Params.encode(message.path, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRecord_Ledger();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.path = hd_1.BIP44Params.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseRecord_Ledger();
        if ((0, helpers_1.isSet)(object.path))
            obj.path = hd_1.BIP44Params.fromJSON(object.path);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.path !== undefined && (obj.path = message.path ? hd_1.BIP44Params.toJSON(message.path) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseRecord_Ledger();
        if (object.path !== undefined && object.path !== null) {
            message.path = hd_1.BIP44Params.fromPartial(object.path);
        }
        return message;
    },
};
function createBaseRecord_Multi() {
    return {};
}
exports.Record_Multi = {
    typeUrl: "/cosmos.crypto.keyring.v1.Multi",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRecord_Multi();
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
        const obj = createBaseRecord_Multi();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseRecord_Multi();
        return message;
    },
};
function createBaseRecord_Offline() {
    return {};
}
exports.Record_Offline = {
    typeUrl: "/cosmos.crypto.keyring.v1.Offline",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRecord_Offline();
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
        const obj = createBaseRecord_Offline();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseRecord_Offline();
        return message;
    },
};
//# sourceMappingURL=record.js.map