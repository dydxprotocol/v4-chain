"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BlockMetadata_DeliverTx = exports.BlockMetadata = exports.StoreKVPair = exports.protobufPackage = void 0;
/* eslint-disable */
const types_1 = require("../../../../tendermint/abci/types");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.base.store.v1beta1";
function createBaseStoreKVPair() {
    return {
        storeKey: "",
        delete: false,
        key: new Uint8Array(),
        value: new Uint8Array(),
    };
}
exports.StoreKVPair = {
    typeUrl: "/cosmos.base.store.v1beta1.StoreKVPair",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.storeKey !== "") {
            writer.uint32(10).string(message.storeKey);
        }
        if (message.delete === true) {
            writer.uint32(16).bool(message.delete);
        }
        if (message.key.length !== 0) {
            writer.uint32(26).bytes(message.key);
        }
        if (message.value.length !== 0) {
            writer.uint32(34).bytes(message.value);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseStoreKVPair();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.storeKey = reader.string();
                    break;
                case 2:
                    message.delete = reader.bool();
                    break;
                case 3:
                    message.key = reader.bytes();
                    break;
                case 4:
                    message.value = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseStoreKVPair();
        if ((0, helpers_1.isSet)(object.storeKey))
            obj.storeKey = String(object.storeKey);
        if ((0, helpers_1.isSet)(object.delete))
            obj.delete = Boolean(object.delete);
        if ((0, helpers_1.isSet)(object.key))
            obj.key = (0, helpers_1.bytesFromBase64)(object.key);
        if ((0, helpers_1.isSet)(object.value))
            obj.value = (0, helpers_1.bytesFromBase64)(object.value);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.storeKey !== undefined && (obj.storeKey = message.storeKey);
        message.delete !== undefined && (obj.delete = message.delete);
        message.key !== undefined &&
            (obj.key = (0, helpers_1.base64FromBytes)(message.key !== undefined ? message.key : new Uint8Array()));
        message.value !== undefined &&
            (obj.value = (0, helpers_1.base64FromBytes)(message.value !== undefined ? message.value : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseStoreKVPair();
        message.storeKey = object.storeKey ?? "";
        message.delete = object.delete ?? false;
        message.key = object.key ?? new Uint8Array();
        message.value = object.value ?? new Uint8Array();
        return message;
    },
};
function createBaseBlockMetadata() {
    return {
        requestBeginBlock: undefined,
        responseBeginBlock: undefined,
        deliverTxs: [],
        requestEndBlock: undefined,
        responseEndBlock: undefined,
        responseCommit: undefined,
    };
}
exports.BlockMetadata = {
    typeUrl: "/cosmos.base.store.v1beta1.BlockMetadata",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.requestBeginBlock !== undefined) {
            types_1.RequestBeginBlock.encode(message.requestBeginBlock, writer.uint32(10).fork()).ldelim();
        }
        if (message.responseBeginBlock !== undefined) {
            types_1.ResponseBeginBlock.encode(message.responseBeginBlock, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.deliverTxs) {
            exports.BlockMetadata_DeliverTx.encode(v, writer.uint32(26).fork()).ldelim();
        }
        if (message.requestEndBlock !== undefined) {
            types_1.RequestEndBlock.encode(message.requestEndBlock, writer.uint32(34).fork()).ldelim();
        }
        if (message.responseEndBlock !== undefined) {
            types_1.ResponseEndBlock.encode(message.responseEndBlock, writer.uint32(42).fork()).ldelim();
        }
        if (message.responseCommit !== undefined) {
            types_1.ResponseCommit.encode(message.responseCommit, writer.uint32(50).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseBlockMetadata();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.requestBeginBlock = types_1.RequestBeginBlock.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.responseBeginBlock = types_1.ResponseBeginBlock.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.deliverTxs.push(exports.BlockMetadata_DeliverTx.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.requestEndBlock = types_1.RequestEndBlock.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.responseEndBlock = types_1.ResponseEndBlock.decode(reader, reader.uint32());
                    break;
                case 6:
                    message.responseCommit = types_1.ResponseCommit.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseBlockMetadata();
        if ((0, helpers_1.isSet)(object.requestBeginBlock))
            obj.requestBeginBlock = types_1.RequestBeginBlock.fromJSON(object.requestBeginBlock);
        if ((0, helpers_1.isSet)(object.responseBeginBlock))
            obj.responseBeginBlock = types_1.ResponseBeginBlock.fromJSON(object.responseBeginBlock);
        if (Array.isArray(object?.deliverTxs))
            obj.deliverTxs = object.deliverTxs.map((e) => exports.BlockMetadata_DeliverTx.fromJSON(e));
        if ((0, helpers_1.isSet)(object.requestEndBlock))
            obj.requestEndBlock = types_1.RequestEndBlock.fromJSON(object.requestEndBlock);
        if ((0, helpers_1.isSet)(object.responseEndBlock))
            obj.responseEndBlock = types_1.ResponseEndBlock.fromJSON(object.responseEndBlock);
        if ((0, helpers_1.isSet)(object.responseCommit))
            obj.responseCommit = types_1.ResponseCommit.fromJSON(object.responseCommit);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.requestBeginBlock !== undefined &&
            (obj.requestBeginBlock = message.requestBeginBlock
                ? types_1.RequestBeginBlock.toJSON(message.requestBeginBlock)
                : undefined);
        message.responseBeginBlock !== undefined &&
            (obj.responseBeginBlock = message.responseBeginBlock
                ? types_1.ResponseBeginBlock.toJSON(message.responseBeginBlock)
                : undefined);
        if (message.deliverTxs) {
            obj.deliverTxs = message.deliverTxs.map((e) => (e ? exports.BlockMetadata_DeliverTx.toJSON(e) : undefined));
        }
        else {
            obj.deliverTxs = [];
        }
        message.requestEndBlock !== undefined &&
            (obj.requestEndBlock = message.requestEndBlock
                ? types_1.RequestEndBlock.toJSON(message.requestEndBlock)
                : undefined);
        message.responseEndBlock !== undefined &&
            (obj.responseEndBlock = message.responseEndBlock
                ? types_1.ResponseEndBlock.toJSON(message.responseEndBlock)
                : undefined);
        message.responseCommit !== undefined &&
            (obj.responseCommit = message.responseCommit
                ? types_1.ResponseCommit.toJSON(message.responseCommit)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseBlockMetadata();
        if (object.requestBeginBlock !== undefined && object.requestBeginBlock !== null) {
            message.requestBeginBlock = types_1.RequestBeginBlock.fromPartial(object.requestBeginBlock);
        }
        if (object.responseBeginBlock !== undefined && object.responseBeginBlock !== null) {
            message.responseBeginBlock = types_1.ResponseBeginBlock.fromPartial(object.responseBeginBlock);
        }
        message.deliverTxs = object.deliverTxs?.map((e) => exports.BlockMetadata_DeliverTx.fromPartial(e)) || [];
        if (object.requestEndBlock !== undefined && object.requestEndBlock !== null) {
            message.requestEndBlock = types_1.RequestEndBlock.fromPartial(object.requestEndBlock);
        }
        if (object.responseEndBlock !== undefined && object.responseEndBlock !== null) {
            message.responseEndBlock = types_1.ResponseEndBlock.fromPartial(object.responseEndBlock);
        }
        if (object.responseCommit !== undefined && object.responseCommit !== null) {
            message.responseCommit = types_1.ResponseCommit.fromPartial(object.responseCommit);
        }
        return message;
    },
};
function createBaseBlockMetadata_DeliverTx() {
    return {
        request: undefined,
        response: undefined,
    };
}
exports.BlockMetadata_DeliverTx = {
    typeUrl: "/cosmos.base.store.v1beta1.DeliverTx",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.request !== undefined) {
            types_1.RequestDeliverTx.encode(message.request, writer.uint32(10).fork()).ldelim();
        }
        if (message.response !== undefined) {
            types_1.ResponseDeliverTx.encode(message.response, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseBlockMetadata_DeliverTx();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.request = types_1.RequestDeliverTx.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.response = types_1.ResponseDeliverTx.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseBlockMetadata_DeliverTx();
        if ((0, helpers_1.isSet)(object.request))
            obj.request = types_1.RequestDeliverTx.fromJSON(object.request);
        if ((0, helpers_1.isSet)(object.response))
            obj.response = types_1.ResponseDeliverTx.fromJSON(object.response);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.request !== undefined &&
            (obj.request = message.request ? types_1.RequestDeliverTx.toJSON(message.request) : undefined);
        message.response !== undefined &&
            (obj.response = message.response ? types_1.ResponseDeliverTx.toJSON(message.response) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseBlockMetadata_DeliverTx();
        if (object.request !== undefined && object.request !== null) {
            message.request = types_1.RequestDeliverTx.fromPartial(object.request);
        }
        if (object.response !== undefined && object.response !== null) {
            message.response = types_1.ResponseDeliverTx.fromPartial(object.response);
        }
        return message;
    },
};
//# sourceMappingURL=listening.js.map