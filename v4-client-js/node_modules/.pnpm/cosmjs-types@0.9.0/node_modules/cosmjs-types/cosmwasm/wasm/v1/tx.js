"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MsgClientImpl = exports.MsgUpdateInstantiateConfigResponse = exports.MsgUpdateInstantiateConfig = exports.MsgClearAdminResponse = exports.MsgClearAdmin = exports.MsgUpdateAdminResponse = exports.MsgUpdateAdmin = exports.MsgMigrateContractResponse = exports.MsgMigrateContract = exports.MsgExecuteContractResponse = exports.MsgExecuteContract = exports.MsgInstantiateContract2Response = exports.MsgInstantiateContractResponse = exports.MsgInstantiateContract2 = exports.MsgInstantiateContract = exports.MsgStoreCodeResponse = exports.MsgStoreCode = exports.protobufPackage = void 0;
/* eslint-disable */
const types_1 = require("./types");
const coin_1 = require("../../../cosmos/base/v1beta1/coin");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmwasm.wasm.v1";
function createBaseMsgStoreCode() {
    return {
        sender: "",
        wasmByteCode: new Uint8Array(),
        instantiatePermission: undefined,
    };
}
exports.MsgStoreCode = {
    typeUrl: "/cosmwasm.wasm.v1.MsgStoreCode",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.wasmByteCode.length !== 0) {
            writer.uint32(18).bytes(message.wasmByteCode);
        }
        if (message.instantiatePermission !== undefined) {
            types_1.AccessConfig.encode(message.instantiatePermission, writer.uint32(42).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgStoreCode();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 2:
                    message.wasmByteCode = reader.bytes();
                    break;
                case 5:
                    message.instantiatePermission = types_1.AccessConfig.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgStoreCode();
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        if ((0, helpers_1.isSet)(object.wasmByteCode))
            obj.wasmByteCode = (0, helpers_1.bytesFromBase64)(object.wasmByteCode);
        if ((0, helpers_1.isSet)(object.instantiatePermission))
            obj.instantiatePermission = types_1.AccessConfig.fromJSON(object.instantiatePermission);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.wasmByteCode !== undefined &&
            (obj.wasmByteCode = (0, helpers_1.base64FromBytes)(message.wasmByteCode !== undefined ? message.wasmByteCode : new Uint8Array()));
        message.instantiatePermission !== undefined &&
            (obj.instantiatePermission = message.instantiatePermission
                ? types_1.AccessConfig.toJSON(message.instantiatePermission)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgStoreCode();
        message.sender = object.sender ?? "";
        message.wasmByteCode = object.wasmByteCode ?? new Uint8Array();
        if (object.instantiatePermission !== undefined && object.instantiatePermission !== null) {
            message.instantiatePermission = types_1.AccessConfig.fromPartial(object.instantiatePermission);
        }
        return message;
    },
};
function createBaseMsgStoreCodeResponse() {
    return {
        codeId: BigInt(0),
        checksum: new Uint8Array(),
    };
}
exports.MsgStoreCodeResponse = {
    typeUrl: "/cosmwasm.wasm.v1.MsgStoreCodeResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.codeId !== BigInt(0)) {
            writer.uint32(8).uint64(message.codeId);
        }
        if (message.checksum.length !== 0) {
            writer.uint32(18).bytes(message.checksum);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgStoreCodeResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.codeId = reader.uint64();
                    break;
                case 2:
                    message.checksum = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgStoreCodeResponse();
        if ((0, helpers_1.isSet)(object.codeId))
            obj.codeId = BigInt(object.codeId.toString());
        if ((0, helpers_1.isSet)(object.checksum))
            obj.checksum = (0, helpers_1.bytesFromBase64)(object.checksum);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.codeId !== undefined && (obj.codeId = (message.codeId || BigInt(0)).toString());
        message.checksum !== undefined &&
            (obj.checksum = (0, helpers_1.base64FromBytes)(message.checksum !== undefined ? message.checksum : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgStoreCodeResponse();
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = BigInt(object.codeId.toString());
        }
        message.checksum = object.checksum ?? new Uint8Array();
        return message;
    },
};
function createBaseMsgInstantiateContract() {
    return {
        sender: "",
        admin: "",
        codeId: BigInt(0),
        label: "",
        msg: new Uint8Array(),
        funds: [],
    };
}
exports.MsgInstantiateContract = {
    typeUrl: "/cosmwasm.wasm.v1.MsgInstantiateContract",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.admin !== "") {
            writer.uint32(18).string(message.admin);
        }
        if (message.codeId !== BigInt(0)) {
            writer.uint32(24).uint64(message.codeId);
        }
        if (message.label !== "") {
            writer.uint32(34).string(message.label);
        }
        if (message.msg.length !== 0) {
            writer.uint32(42).bytes(message.msg);
        }
        for (const v of message.funds) {
            coin_1.Coin.encode(v, writer.uint32(50).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgInstantiateContract();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 2:
                    message.admin = reader.string();
                    break;
                case 3:
                    message.codeId = reader.uint64();
                    break;
                case 4:
                    message.label = reader.string();
                    break;
                case 5:
                    message.msg = reader.bytes();
                    break;
                case 6:
                    message.funds.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgInstantiateContract();
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        if ((0, helpers_1.isSet)(object.admin))
            obj.admin = String(object.admin);
        if ((0, helpers_1.isSet)(object.codeId))
            obj.codeId = BigInt(object.codeId.toString());
        if ((0, helpers_1.isSet)(object.label))
            obj.label = String(object.label);
        if ((0, helpers_1.isSet)(object.msg))
            obj.msg = (0, helpers_1.bytesFromBase64)(object.msg);
        if (Array.isArray(object?.funds))
            obj.funds = object.funds.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.admin !== undefined && (obj.admin = message.admin);
        message.codeId !== undefined && (obj.codeId = (message.codeId || BigInt(0)).toString());
        message.label !== undefined && (obj.label = message.label);
        message.msg !== undefined &&
            (obj.msg = (0, helpers_1.base64FromBytes)(message.msg !== undefined ? message.msg : new Uint8Array()));
        if (message.funds) {
            obj.funds = message.funds.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.funds = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgInstantiateContract();
        message.sender = object.sender ?? "";
        message.admin = object.admin ?? "";
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = BigInt(object.codeId.toString());
        }
        message.label = object.label ?? "";
        message.msg = object.msg ?? new Uint8Array();
        message.funds = object.funds?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseMsgInstantiateContract2() {
    return {
        sender: "",
        admin: "",
        codeId: BigInt(0),
        label: "",
        msg: new Uint8Array(),
        funds: [],
        salt: new Uint8Array(),
        fixMsg: false,
    };
}
exports.MsgInstantiateContract2 = {
    typeUrl: "/cosmwasm.wasm.v1.MsgInstantiateContract2",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.admin !== "") {
            writer.uint32(18).string(message.admin);
        }
        if (message.codeId !== BigInt(0)) {
            writer.uint32(24).uint64(message.codeId);
        }
        if (message.label !== "") {
            writer.uint32(34).string(message.label);
        }
        if (message.msg.length !== 0) {
            writer.uint32(42).bytes(message.msg);
        }
        for (const v of message.funds) {
            coin_1.Coin.encode(v, writer.uint32(50).fork()).ldelim();
        }
        if (message.salt.length !== 0) {
            writer.uint32(58).bytes(message.salt);
        }
        if (message.fixMsg === true) {
            writer.uint32(64).bool(message.fixMsg);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgInstantiateContract2();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 2:
                    message.admin = reader.string();
                    break;
                case 3:
                    message.codeId = reader.uint64();
                    break;
                case 4:
                    message.label = reader.string();
                    break;
                case 5:
                    message.msg = reader.bytes();
                    break;
                case 6:
                    message.funds.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 7:
                    message.salt = reader.bytes();
                    break;
                case 8:
                    message.fixMsg = reader.bool();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgInstantiateContract2();
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        if ((0, helpers_1.isSet)(object.admin))
            obj.admin = String(object.admin);
        if ((0, helpers_1.isSet)(object.codeId))
            obj.codeId = BigInt(object.codeId.toString());
        if ((0, helpers_1.isSet)(object.label))
            obj.label = String(object.label);
        if ((0, helpers_1.isSet)(object.msg))
            obj.msg = (0, helpers_1.bytesFromBase64)(object.msg);
        if (Array.isArray(object?.funds))
            obj.funds = object.funds.map((e) => coin_1.Coin.fromJSON(e));
        if ((0, helpers_1.isSet)(object.salt))
            obj.salt = (0, helpers_1.bytesFromBase64)(object.salt);
        if ((0, helpers_1.isSet)(object.fixMsg))
            obj.fixMsg = Boolean(object.fixMsg);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.admin !== undefined && (obj.admin = message.admin);
        message.codeId !== undefined && (obj.codeId = (message.codeId || BigInt(0)).toString());
        message.label !== undefined && (obj.label = message.label);
        message.msg !== undefined &&
            (obj.msg = (0, helpers_1.base64FromBytes)(message.msg !== undefined ? message.msg : new Uint8Array()));
        if (message.funds) {
            obj.funds = message.funds.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.funds = [];
        }
        message.salt !== undefined &&
            (obj.salt = (0, helpers_1.base64FromBytes)(message.salt !== undefined ? message.salt : new Uint8Array()));
        message.fixMsg !== undefined && (obj.fixMsg = message.fixMsg);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgInstantiateContract2();
        message.sender = object.sender ?? "";
        message.admin = object.admin ?? "";
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = BigInt(object.codeId.toString());
        }
        message.label = object.label ?? "";
        message.msg = object.msg ?? new Uint8Array();
        message.funds = object.funds?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        message.salt = object.salt ?? new Uint8Array();
        message.fixMsg = object.fixMsg ?? false;
        return message;
    },
};
function createBaseMsgInstantiateContractResponse() {
    return {
        address: "",
        data: new Uint8Array(),
    };
}
exports.MsgInstantiateContractResponse = {
    typeUrl: "/cosmwasm.wasm.v1.MsgInstantiateContractResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        if (message.data.length !== 0) {
            writer.uint32(18).bytes(message.data);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgInstantiateContractResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
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
        const obj = createBaseMsgInstantiateContractResponse();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgInstantiateContractResponse();
        message.address = object.address ?? "";
        message.data = object.data ?? new Uint8Array();
        return message;
    },
};
function createBaseMsgInstantiateContract2Response() {
    return {
        address: "",
        data: new Uint8Array(),
    };
}
exports.MsgInstantiateContract2Response = {
    typeUrl: "/cosmwasm.wasm.v1.MsgInstantiateContract2Response",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        if (message.data.length !== 0) {
            writer.uint32(18).bytes(message.data);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgInstantiateContract2Response();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
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
        const obj = createBaseMsgInstantiateContract2Response();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgInstantiateContract2Response();
        message.address = object.address ?? "";
        message.data = object.data ?? new Uint8Array();
        return message;
    },
};
function createBaseMsgExecuteContract() {
    return {
        sender: "",
        contract: "",
        msg: new Uint8Array(),
        funds: [],
    };
}
exports.MsgExecuteContract = {
    typeUrl: "/cosmwasm.wasm.v1.MsgExecuteContract",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.contract !== "") {
            writer.uint32(18).string(message.contract);
        }
        if (message.msg.length !== 0) {
            writer.uint32(26).bytes(message.msg);
        }
        for (const v of message.funds) {
            coin_1.Coin.encode(v, writer.uint32(42).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgExecuteContract();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 2:
                    message.contract = reader.string();
                    break;
                case 3:
                    message.msg = reader.bytes();
                    break;
                case 5:
                    message.funds.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgExecuteContract();
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        if ((0, helpers_1.isSet)(object.contract))
            obj.contract = String(object.contract);
        if ((0, helpers_1.isSet)(object.msg))
            obj.msg = (0, helpers_1.bytesFromBase64)(object.msg);
        if (Array.isArray(object?.funds))
            obj.funds = object.funds.map((e) => coin_1.Coin.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.contract !== undefined && (obj.contract = message.contract);
        message.msg !== undefined &&
            (obj.msg = (0, helpers_1.base64FromBytes)(message.msg !== undefined ? message.msg : new Uint8Array()));
        if (message.funds) {
            obj.funds = message.funds.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.funds = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgExecuteContract();
        message.sender = object.sender ?? "";
        message.contract = object.contract ?? "";
        message.msg = object.msg ?? new Uint8Array();
        message.funds = object.funds?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        return message;
    },
};
function createBaseMsgExecuteContractResponse() {
    return {
        data: new Uint8Array(),
    };
}
exports.MsgExecuteContractResponse = {
    typeUrl: "/cosmwasm.wasm.v1.MsgExecuteContractResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.data.length !== 0) {
            writer.uint32(10).bytes(message.data);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgExecuteContractResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseMsgExecuteContractResponse();
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgExecuteContractResponse();
        message.data = object.data ?? new Uint8Array();
        return message;
    },
};
function createBaseMsgMigrateContract() {
    return {
        sender: "",
        contract: "",
        codeId: BigInt(0),
        msg: new Uint8Array(),
    };
}
exports.MsgMigrateContract = {
    typeUrl: "/cosmwasm.wasm.v1.MsgMigrateContract",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.contract !== "") {
            writer.uint32(18).string(message.contract);
        }
        if (message.codeId !== BigInt(0)) {
            writer.uint32(24).uint64(message.codeId);
        }
        if (message.msg.length !== 0) {
            writer.uint32(34).bytes(message.msg);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgMigrateContract();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 2:
                    message.contract = reader.string();
                    break;
                case 3:
                    message.codeId = reader.uint64();
                    break;
                case 4:
                    message.msg = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgMigrateContract();
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        if ((0, helpers_1.isSet)(object.contract))
            obj.contract = String(object.contract);
        if ((0, helpers_1.isSet)(object.codeId))
            obj.codeId = BigInt(object.codeId.toString());
        if ((0, helpers_1.isSet)(object.msg))
            obj.msg = (0, helpers_1.bytesFromBase64)(object.msg);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.contract !== undefined && (obj.contract = message.contract);
        message.codeId !== undefined && (obj.codeId = (message.codeId || BigInt(0)).toString());
        message.msg !== undefined &&
            (obj.msg = (0, helpers_1.base64FromBytes)(message.msg !== undefined ? message.msg : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgMigrateContract();
        message.sender = object.sender ?? "";
        message.contract = object.contract ?? "";
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = BigInt(object.codeId.toString());
        }
        message.msg = object.msg ?? new Uint8Array();
        return message;
    },
};
function createBaseMsgMigrateContractResponse() {
    return {
        data: new Uint8Array(),
    };
}
exports.MsgMigrateContractResponse = {
    typeUrl: "/cosmwasm.wasm.v1.MsgMigrateContractResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.data.length !== 0) {
            writer.uint32(10).bytes(message.data);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgMigrateContractResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseMsgMigrateContractResponse();
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgMigrateContractResponse();
        message.data = object.data ?? new Uint8Array();
        return message;
    },
};
function createBaseMsgUpdateAdmin() {
    return {
        sender: "",
        newAdmin: "",
        contract: "",
    };
}
exports.MsgUpdateAdmin = {
    typeUrl: "/cosmwasm.wasm.v1.MsgUpdateAdmin",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.newAdmin !== "") {
            writer.uint32(18).string(message.newAdmin);
        }
        if (message.contract !== "") {
            writer.uint32(26).string(message.contract);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpdateAdmin();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 2:
                    message.newAdmin = reader.string();
                    break;
                case 3:
                    message.contract = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgUpdateAdmin();
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        if ((0, helpers_1.isSet)(object.newAdmin))
            obj.newAdmin = String(object.newAdmin);
        if ((0, helpers_1.isSet)(object.contract))
            obj.contract = String(object.contract);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.newAdmin !== undefined && (obj.newAdmin = message.newAdmin);
        message.contract !== undefined && (obj.contract = message.contract);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUpdateAdmin();
        message.sender = object.sender ?? "";
        message.newAdmin = object.newAdmin ?? "";
        message.contract = object.contract ?? "";
        return message;
    },
};
function createBaseMsgUpdateAdminResponse() {
    return {};
}
exports.MsgUpdateAdminResponse = {
    typeUrl: "/cosmwasm.wasm.v1.MsgUpdateAdminResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpdateAdminResponse();
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
        const obj = createBaseMsgUpdateAdminResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgUpdateAdminResponse();
        return message;
    },
};
function createBaseMsgClearAdmin() {
    return {
        sender: "",
        contract: "",
    };
}
exports.MsgClearAdmin = {
    typeUrl: "/cosmwasm.wasm.v1.MsgClearAdmin",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.contract !== "") {
            writer.uint32(26).string(message.contract);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgClearAdmin();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 3:
                    message.contract = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgClearAdmin();
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        if ((0, helpers_1.isSet)(object.contract))
            obj.contract = String(object.contract);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.contract !== undefined && (obj.contract = message.contract);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgClearAdmin();
        message.sender = object.sender ?? "";
        message.contract = object.contract ?? "";
        return message;
    },
};
function createBaseMsgClearAdminResponse() {
    return {};
}
exports.MsgClearAdminResponse = {
    typeUrl: "/cosmwasm.wasm.v1.MsgClearAdminResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgClearAdminResponse();
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
        const obj = createBaseMsgClearAdminResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgClearAdminResponse();
        return message;
    },
};
function createBaseMsgUpdateInstantiateConfig() {
    return {
        sender: "",
        codeId: BigInt(0),
        newInstantiatePermission: undefined,
    };
}
exports.MsgUpdateInstantiateConfig = {
    typeUrl: "/cosmwasm.wasm.v1.MsgUpdateInstantiateConfig",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.codeId !== BigInt(0)) {
            writer.uint32(16).uint64(message.codeId);
        }
        if (message.newInstantiatePermission !== undefined) {
            types_1.AccessConfig.encode(message.newInstantiatePermission, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpdateInstantiateConfig();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 2:
                    message.codeId = reader.uint64();
                    break;
                case 3:
                    message.newInstantiatePermission = types_1.AccessConfig.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseMsgUpdateInstantiateConfig();
        if ((0, helpers_1.isSet)(object.sender))
            obj.sender = String(object.sender);
        if ((0, helpers_1.isSet)(object.codeId))
            obj.codeId = BigInt(object.codeId.toString());
        if ((0, helpers_1.isSet)(object.newInstantiatePermission))
            obj.newInstantiatePermission = types_1.AccessConfig.fromJSON(object.newInstantiatePermission);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.codeId !== undefined && (obj.codeId = (message.codeId || BigInt(0)).toString());
        message.newInstantiatePermission !== undefined &&
            (obj.newInstantiatePermission = message.newInstantiatePermission
                ? types_1.AccessConfig.toJSON(message.newInstantiatePermission)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgUpdateInstantiateConfig();
        message.sender = object.sender ?? "";
        if (object.codeId !== undefined && object.codeId !== null) {
            message.codeId = BigInt(object.codeId.toString());
        }
        if (object.newInstantiatePermission !== undefined && object.newInstantiatePermission !== null) {
            message.newInstantiatePermission = types_1.AccessConfig.fromPartial(object.newInstantiatePermission);
        }
        return message;
    },
};
function createBaseMsgUpdateInstantiateConfigResponse() {
    return {};
}
exports.MsgUpdateInstantiateConfigResponse = {
    typeUrl: "/cosmwasm.wasm.v1.MsgUpdateInstantiateConfigResponse",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgUpdateInstantiateConfigResponse();
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
        const obj = createBaseMsgUpdateInstantiateConfigResponse();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseMsgUpdateInstantiateConfigResponse();
        return message;
    },
};
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.StoreCode = this.StoreCode.bind(this);
        this.InstantiateContract = this.InstantiateContract.bind(this);
        this.InstantiateContract2 = this.InstantiateContract2.bind(this);
        this.ExecuteContract = this.ExecuteContract.bind(this);
        this.MigrateContract = this.MigrateContract.bind(this);
        this.UpdateAdmin = this.UpdateAdmin.bind(this);
        this.ClearAdmin = this.ClearAdmin.bind(this);
        this.UpdateInstantiateConfig = this.UpdateInstantiateConfig.bind(this);
    }
    StoreCode(request) {
        const data = exports.MsgStoreCode.encode(request).finish();
        const promise = this.rpc.request("cosmwasm.wasm.v1.Msg", "StoreCode", data);
        return promise.then((data) => exports.MsgStoreCodeResponse.decode(new binary_1.BinaryReader(data)));
    }
    InstantiateContract(request) {
        const data = exports.MsgInstantiateContract.encode(request).finish();
        const promise = this.rpc.request("cosmwasm.wasm.v1.Msg", "InstantiateContract", data);
        return promise.then((data) => exports.MsgInstantiateContractResponse.decode(new binary_1.BinaryReader(data)));
    }
    InstantiateContract2(request) {
        const data = exports.MsgInstantiateContract2.encode(request).finish();
        const promise = this.rpc.request("cosmwasm.wasm.v1.Msg", "InstantiateContract2", data);
        return promise.then((data) => exports.MsgInstantiateContract2Response.decode(new binary_1.BinaryReader(data)));
    }
    ExecuteContract(request) {
        const data = exports.MsgExecuteContract.encode(request).finish();
        const promise = this.rpc.request("cosmwasm.wasm.v1.Msg", "ExecuteContract", data);
        return promise.then((data) => exports.MsgExecuteContractResponse.decode(new binary_1.BinaryReader(data)));
    }
    MigrateContract(request) {
        const data = exports.MsgMigrateContract.encode(request).finish();
        const promise = this.rpc.request("cosmwasm.wasm.v1.Msg", "MigrateContract", data);
        return promise.then((data) => exports.MsgMigrateContractResponse.decode(new binary_1.BinaryReader(data)));
    }
    UpdateAdmin(request) {
        const data = exports.MsgUpdateAdmin.encode(request).finish();
        const promise = this.rpc.request("cosmwasm.wasm.v1.Msg", "UpdateAdmin", data);
        return promise.then((data) => exports.MsgUpdateAdminResponse.decode(new binary_1.BinaryReader(data)));
    }
    ClearAdmin(request) {
        const data = exports.MsgClearAdmin.encode(request).finish();
        const promise = this.rpc.request("cosmwasm.wasm.v1.Msg", "ClearAdmin", data);
        return promise.then((data) => exports.MsgClearAdminResponse.decode(new binary_1.BinaryReader(data)));
    }
    UpdateInstantiateConfig(request) {
        const data = exports.MsgUpdateInstantiateConfig.encode(request).finish();
        const promise = this.rpc.request("cosmwasm.wasm.v1.Msg", "UpdateInstantiateConfig", data);
        return promise.then((data) => exports.MsgUpdateInstantiateConfigResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=tx.js.map