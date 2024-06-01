"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AuxSignerData = exports.Tip = exports.Fee = exports.ModeInfo_Multi = exports.ModeInfo_Single = exports.ModeInfo = exports.SignerInfo = exports.AuthInfo = exports.TxBody = exports.SignDocDirectAux = exports.SignDoc = exports.TxRaw = exports.Tx = exports.protobufPackage = void 0;
/* eslint-disable */
const any_1 = require("../../../google/protobuf/any");
const signing_1 = require("../signing/v1beta1/signing");
const multisig_1 = require("../../crypto/multisig/v1beta1/multisig");
const coin_1 = require("../../base/v1beta1/coin");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.tx.v1beta1";
function createBaseTx() {
    return {
        body: undefined,
        authInfo: undefined,
        signatures: [],
    };
}
exports.Tx = {
    typeUrl: "/cosmos.tx.v1beta1.Tx",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.body !== undefined) {
            exports.TxBody.encode(message.body, writer.uint32(10).fork()).ldelim();
        }
        if (message.authInfo !== undefined) {
            exports.AuthInfo.encode(message.authInfo, writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.signatures) {
            writer.uint32(26).bytes(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTx();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.body = exports.TxBody.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.authInfo = exports.AuthInfo.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.signatures.push(reader.bytes());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTx();
        if ((0, helpers_1.isSet)(object.body))
            obj.body = exports.TxBody.fromJSON(object.body);
        if ((0, helpers_1.isSet)(object.authInfo))
            obj.authInfo = exports.AuthInfo.fromJSON(object.authInfo);
        if (Array.isArray(object?.signatures))
            obj.signatures = object.signatures.map((e) => (0, helpers_1.bytesFromBase64)(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.body !== undefined && (obj.body = message.body ? exports.TxBody.toJSON(message.body) : undefined);
        message.authInfo !== undefined &&
            (obj.authInfo = message.authInfo ? exports.AuthInfo.toJSON(message.authInfo) : undefined);
        if (message.signatures) {
            obj.signatures = message.signatures.map((e) => (0, helpers_1.base64FromBytes)(e !== undefined ? e : new Uint8Array()));
        }
        else {
            obj.signatures = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTx();
        if (object.body !== undefined && object.body !== null) {
            message.body = exports.TxBody.fromPartial(object.body);
        }
        if (object.authInfo !== undefined && object.authInfo !== null) {
            message.authInfo = exports.AuthInfo.fromPartial(object.authInfo);
        }
        message.signatures = object.signatures?.map((e) => e) || [];
        return message;
    },
};
function createBaseTxRaw() {
    return {
        bodyBytes: new Uint8Array(),
        authInfoBytes: new Uint8Array(),
        signatures: [],
    };
}
exports.TxRaw = {
    typeUrl: "/cosmos.tx.v1beta1.TxRaw",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.bodyBytes.length !== 0) {
            writer.uint32(10).bytes(message.bodyBytes);
        }
        if (message.authInfoBytes.length !== 0) {
            writer.uint32(18).bytes(message.authInfoBytes);
        }
        for (const v of message.signatures) {
            writer.uint32(26).bytes(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxRaw();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.bodyBytes = reader.bytes();
                    break;
                case 2:
                    message.authInfoBytes = reader.bytes();
                    break;
                case 3:
                    message.signatures.push(reader.bytes());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxRaw();
        if ((0, helpers_1.isSet)(object.bodyBytes))
            obj.bodyBytes = (0, helpers_1.bytesFromBase64)(object.bodyBytes);
        if ((0, helpers_1.isSet)(object.authInfoBytes))
            obj.authInfoBytes = (0, helpers_1.bytesFromBase64)(object.authInfoBytes);
        if (Array.isArray(object?.signatures))
            obj.signatures = object.signatures.map((e) => (0, helpers_1.bytesFromBase64)(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.bodyBytes !== undefined &&
            (obj.bodyBytes = (0, helpers_1.base64FromBytes)(message.bodyBytes !== undefined ? message.bodyBytes : new Uint8Array()));
        message.authInfoBytes !== undefined &&
            (obj.authInfoBytes = (0, helpers_1.base64FromBytes)(message.authInfoBytes !== undefined ? message.authInfoBytes : new Uint8Array()));
        if (message.signatures) {
            obj.signatures = message.signatures.map((e) => (0, helpers_1.base64FromBytes)(e !== undefined ? e : new Uint8Array()));
        }
        else {
            obj.signatures = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxRaw();
        message.bodyBytes = object.bodyBytes ?? new Uint8Array();
        message.authInfoBytes = object.authInfoBytes ?? new Uint8Array();
        message.signatures = object.signatures?.map((e) => e) || [];
        return message;
    },
};
function createBaseSignDoc() {
    return {
        bodyBytes: new Uint8Array(),
        authInfoBytes: new Uint8Array(),
        chainId: "",
        accountNumber: BigInt(0),
    };
}
exports.SignDoc = {
    typeUrl: "/cosmos.tx.v1beta1.SignDoc",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.bodyBytes.length !== 0) {
            writer.uint32(10).bytes(message.bodyBytes);
        }
        if (message.authInfoBytes.length !== 0) {
            writer.uint32(18).bytes(message.authInfoBytes);
        }
        if (message.chainId !== "") {
            writer.uint32(26).string(message.chainId);
        }
        if (message.accountNumber !== BigInt(0)) {
            writer.uint32(32).uint64(message.accountNumber);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSignDoc();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.bodyBytes = reader.bytes();
                    break;
                case 2:
                    message.authInfoBytes = reader.bytes();
                    break;
                case 3:
                    message.chainId = reader.string();
                    break;
                case 4:
                    message.accountNumber = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSignDoc();
        if ((0, helpers_1.isSet)(object.bodyBytes))
            obj.bodyBytes = (0, helpers_1.bytesFromBase64)(object.bodyBytes);
        if ((0, helpers_1.isSet)(object.authInfoBytes))
            obj.authInfoBytes = (0, helpers_1.bytesFromBase64)(object.authInfoBytes);
        if ((0, helpers_1.isSet)(object.chainId))
            obj.chainId = String(object.chainId);
        if ((0, helpers_1.isSet)(object.accountNumber))
            obj.accountNumber = BigInt(object.accountNumber.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.bodyBytes !== undefined &&
            (obj.bodyBytes = (0, helpers_1.base64FromBytes)(message.bodyBytes !== undefined ? message.bodyBytes : new Uint8Array()));
        message.authInfoBytes !== undefined &&
            (obj.authInfoBytes = (0, helpers_1.base64FromBytes)(message.authInfoBytes !== undefined ? message.authInfoBytes : new Uint8Array()));
        message.chainId !== undefined && (obj.chainId = message.chainId);
        message.accountNumber !== undefined &&
            (obj.accountNumber = (message.accountNumber || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSignDoc();
        message.bodyBytes = object.bodyBytes ?? new Uint8Array();
        message.authInfoBytes = object.authInfoBytes ?? new Uint8Array();
        message.chainId = object.chainId ?? "";
        if (object.accountNumber !== undefined && object.accountNumber !== null) {
            message.accountNumber = BigInt(object.accountNumber.toString());
        }
        return message;
    },
};
function createBaseSignDocDirectAux() {
    return {
        bodyBytes: new Uint8Array(),
        publicKey: undefined,
        chainId: "",
        accountNumber: BigInt(0),
        sequence: BigInt(0),
        tip: undefined,
    };
}
exports.SignDocDirectAux = {
    typeUrl: "/cosmos.tx.v1beta1.SignDocDirectAux",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.bodyBytes.length !== 0) {
            writer.uint32(10).bytes(message.bodyBytes);
        }
        if (message.publicKey !== undefined) {
            any_1.Any.encode(message.publicKey, writer.uint32(18).fork()).ldelim();
        }
        if (message.chainId !== "") {
            writer.uint32(26).string(message.chainId);
        }
        if (message.accountNumber !== BigInt(0)) {
            writer.uint32(32).uint64(message.accountNumber);
        }
        if (message.sequence !== BigInt(0)) {
            writer.uint32(40).uint64(message.sequence);
        }
        if (message.tip !== undefined) {
            exports.Tip.encode(message.tip, writer.uint32(50).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSignDocDirectAux();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.bodyBytes = reader.bytes();
                    break;
                case 2:
                    message.publicKey = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.chainId = reader.string();
                    break;
                case 4:
                    message.accountNumber = reader.uint64();
                    break;
                case 5:
                    message.sequence = reader.uint64();
                    break;
                case 6:
                    message.tip = exports.Tip.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSignDocDirectAux();
        if ((0, helpers_1.isSet)(object.bodyBytes))
            obj.bodyBytes = (0, helpers_1.bytesFromBase64)(object.bodyBytes);
        if ((0, helpers_1.isSet)(object.publicKey))
            obj.publicKey = any_1.Any.fromJSON(object.publicKey);
        if ((0, helpers_1.isSet)(object.chainId))
            obj.chainId = String(object.chainId);
        if ((0, helpers_1.isSet)(object.accountNumber))
            obj.accountNumber = BigInt(object.accountNumber.toString());
        if ((0, helpers_1.isSet)(object.sequence))
            obj.sequence = BigInt(object.sequence.toString());
        if ((0, helpers_1.isSet)(object.tip))
            obj.tip = exports.Tip.fromJSON(object.tip);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.bodyBytes !== undefined &&
            (obj.bodyBytes = (0, helpers_1.base64FromBytes)(message.bodyBytes !== undefined ? message.bodyBytes : new Uint8Array()));
        message.publicKey !== undefined &&
            (obj.publicKey = message.publicKey ? any_1.Any.toJSON(message.publicKey) : undefined);
        message.chainId !== undefined && (obj.chainId = message.chainId);
        message.accountNumber !== undefined &&
            (obj.accountNumber = (message.accountNumber || BigInt(0)).toString());
        message.sequence !== undefined && (obj.sequence = (message.sequence || BigInt(0)).toString());
        message.tip !== undefined && (obj.tip = message.tip ? exports.Tip.toJSON(message.tip) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSignDocDirectAux();
        message.bodyBytes = object.bodyBytes ?? new Uint8Array();
        if (object.publicKey !== undefined && object.publicKey !== null) {
            message.publicKey = any_1.Any.fromPartial(object.publicKey);
        }
        message.chainId = object.chainId ?? "";
        if (object.accountNumber !== undefined && object.accountNumber !== null) {
            message.accountNumber = BigInt(object.accountNumber.toString());
        }
        if (object.sequence !== undefined && object.sequence !== null) {
            message.sequence = BigInt(object.sequence.toString());
        }
        if (object.tip !== undefined && object.tip !== null) {
            message.tip = exports.Tip.fromPartial(object.tip);
        }
        return message;
    },
};
function createBaseTxBody() {
    return {
        messages: [],
        memo: "",
        timeoutHeight: BigInt(0),
        extensionOptions: [],
        nonCriticalExtensionOptions: [],
    };
}
exports.TxBody = {
    typeUrl: "/cosmos.tx.v1beta1.TxBody",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.messages) {
            any_1.Any.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.memo !== "") {
            writer.uint32(18).string(message.memo);
        }
        if (message.timeoutHeight !== BigInt(0)) {
            writer.uint32(24).uint64(message.timeoutHeight);
        }
        for (const v of message.extensionOptions) {
            any_1.Any.encode(v, writer.uint32(8186).fork()).ldelim();
        }
        for (const v of message.nonCriticalExtensionOptions) {
            any_1.Any.encode(v, writer.uint32(16378).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxBody();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.messages.push(any_1.Any.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.memo = reader.string();
                    break;
                case 3:
                    message.timeoutHeight = reader.uint64();
                    break;
                case 1023:
                    message.extensionOptions.push(any_1.Any.decode(reader, reader.uint32()));
                    break;
                case 2047:
                    message.nonCriticalExtensionOptions.push(any_1.Any.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxBody();
        if (Array.isArray(object?.messages))
            obj.messages = object.messages.map((e) => any_1.Any.fromJSON(e));
        if ((0, helpers_1.isSet)(object.memo))
            obj.memo = String(object.memo);
        if ((0, helpers_1.isSet)(object.timeoutHeight))
            obj.timeoutHeight = BigInt(object.timeoutHeight.toString());
        if (Array.isArray(object?.extensionOptions))
            obj.extensionOptions = object.extensionOptions.map((e) => any_1.Any.fromJSON(e));
        if (Array.isArray(object?.nonCriticalExtensionOptions))
            obj.nonCriticalExtensionOptions = object.nonCriticalExtensionOptions.map((e) => any_1.Any.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.messages) {
            obj.messages = message.messages.map((e) => (e ? any_1.Any.toJSON(e) : undefined));
        }
        else {
            obj.messages = [];
        }
        message.memo !== undefined && (obj.memo = message.memo);
        message.timeoutHeight !== undefined &&
            (obj.timeoutHeight = (message.timeoutHeight || BigInt(0)).toString());
        if (message.extensionOptions) {
            obj.extensionOptions = message.extensionOptions.map((e) => (e ? any_1.Any.toJSON(e) : undefined));
        }
        else {
            obj.extensionOptions = [];
        }
        if (message.nonCriticalExtensionOptions) {
            obj.nonCriticalExtensionOptions = message.nonCriticalExtensionOptions.map((e) => e ? any_1.Any.toJSON(e) : undefined);
        }
        else {
            obj.nonCriticalExtensionOptions = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxBody();
        message.messages = object.messages?.map((e) => any_1.Any.fromPartial(e)) || [];
        message.memo = object.memo ?? "";
        if (object.timeoutHeight !== undefined && object.timeoutHeight !== null) {
            message.timeoutHeight = BigInt(object.timeoutHeight.toString());
        }
        message.extensionOptions = object.extensionOptions?.map((e) => any_1.Any.fromPartial(e)) || [];
        message.nonCriticalExtensionOptions =
            object.nonCriticalExtensionOptions?.map((e) => any_1.Any.fromPartial(e)) || [];
        return message;
    },
};
function createBaseAuthInfo() {
    return {
        signerInfos: [],
        fee: undefined,
        tip: undefined,
    };
}
exports.AuthInfo = {
    typeUrl: "/cosmos.tx.v1beta1.AuthInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.signerInfos) {
            exports.SignerInfo.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.fee !== undefined) {
            exports.Fee.encode(message.fee, writer.uint32(18).fork()).ldelim();
        }
        if (message.tip !== undefined) {
            exports.Tip.encode(message.tip, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseAuthInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.signerInfos.push(exports.SignerInfo.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.fee = exports.Fee.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.tip = exports.Tip.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseAuthInfo();
        if (Array.isArray(object?.signerInfos))
            obj.signerInfos = object.signerInfos.map((e) => exports.SignerInfo.fromJSON(e));
        if ((0, helpers_1.isSet)(object.fee))
            obj.fee = exports.Fee.fromJSON(object.fee);
        if ((0, helpers_1.isSet)(object.tip))
            obj.tip = exports.Tip.fromJSON(object.tip);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.signerInfos) {
            obj.signerInfos = message.signerInfos.map((e) => (e ? exports.SignerInfo.toJSON(e) : undefined));
        }
        else {
            obj.signerInfos = [];
        }
        message.fee !== undefined && (obj.fee = message.fee ? exports.Fee.toJSON(message.fee) : undefined);
        message.tip !== undefined && (obj.tip = message.tip ? exports.Tip.toJSON(message.tip) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseAuthInfo();
        message.signerInfos = object.signerInfos?.map((e) => exports.SignerInfo.fromPartial(e)) || [];
        if (object.fee !== undefined && object.fee !== null) {
            message.fee = exports.Fee.fromPartial(object.fee);
        }
        if (object.tip !== undefined && object.tip !== null) {
            message.tip = exports.Tip.fromPartial(object.tip);
        }
        return message;
    },
};
function createBaseSignerInfo() {
    return {
        publicKey: undefined,
        modeInfo: undefined,
        sequence: BigInt(0),
    };
}
exports.SignerInfo = {
    typeUrl: "/cosmos.tx.v1beta1.SignerInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.publicKey !== undefined) {
            any_1.Any.encode(message.publicKey, writer.uint32(10).fork()).ldelim();
        }
        if (message.modeInfo !== undefined) {
            exports.ModeInfo.encode(message.modeInfo, writer.uint32(18).fork()).ldelim();
        }
        if (message.sequence !== BigInt(0)) {
            writer.uint32(24).uint64(message.sequence);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSignerInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.publicKey = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.modeInfo = exports.ModeInfo.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.sequence = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSignerInfo();
        if ((0, helpers_1.isSet)(object.publicKey))
            obj.publicKey = any_1.Any.fromJSON(object.publicKey);
        if ((0, helpers_1.isSet)(object.modeInfo))
            obj.modeInfo = exports.ModeInfo.fromJSON(object.modeInfo);
        if ((0, helpers_1.isSet)(object.sequence))
            obj.sequence = BigInt(object.sequence.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.publicKey !== undefined &&
            (obj.publicKey = message.publicKey ? any_1.Any.toJSON(message.publicKey) : undefined);
        message.modeInfo !== undefined &&
            (obj.modeInfo = message.modeInfo ? exports.ModeInfo.toJSON(message.modeInfo) : undefined);
        message.sequence !== undefined && (obj.sequence = (message.sequence || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSignerInfo();
        if (object.publicKey !== undefined && object.publicKey !== null) {
            message.publicKey = any_1.Any.fromPartial(object.publicKey);
        }
        if (object.modeInfo !== undefined && object.modeInfo !== null) {
            message.modeInfo = exports.ModeInfo.fromPartial(object.modeInfo);
        }
        if (object.sequence !== undefined && object.sequence !== null) {
            message.sequence = BigInt(object.sequence.toString());
        }
        return message;
    },
};
function createBaseModeInfo() {
    return {
        single: undefined,
        multi: undefined,
    };
}
exports.ModeInfo = {
    typeUrl: "/cosmos.tx.v1beta1.ModeInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.single !== undefined) {
            exports.ModeInfo_Single.encode(message.single, writer.uint32(10).fork()).ldelim();
        }
        if (message.multi !== undefined) {
            exports.ModeInfo_Multi.encode(message.multi, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModeInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.single = exports.ModeInfo_Single.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.multi = exports.ModeInfo_Multi.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseModeInfo();
        if ((0, helpers_1.isSet)(object.single))
            obj.single = exports.ModeInfo_Single.fromJSON(object.single);
        if ((0, helpers_1.isSet)(object.multi))
            obj.multi = exports.ModeInfo_Multi.fromJSON(object.multi);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.single !== undefined &&
            (obj.single = message.single ? exports.ModeInfo_Single.toJSON(message.single) : undefined);
        message.multi !== undefined &&
            (obj.multi = message.multi ? exports.ModeInfo_Multi.toJSON(message.multi) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModeInfo();
        if (object.single !== undefined && object.single !== null) {
            message.single = exports.ModeInfo_Single.fromPartial(object.single);
        }
        if (object.multi !== undefined && object.multi !== null) {
            message.multi = exports.ModeInfo_Multi.fromPartial(object.multi);
        }
        return message;
    },
};
function createBaseModeInfo_Single() {
    return {
        mode: 0,
    };
}
exports.ModeInfo_Single = {
    typeUrl: "/cosmos.tx.v1beta1.Single",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.mode !== 0) {
            writer.uint32(8).int32(message.mode);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModeInfo_Single();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.mode = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseModeInfo_Single();
        if ((0, helpers_1.isSet)(object.mode))
            obj.mode = (0, signing_1.signModeFromJSON)(object.mode);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.mode !== undefined && (obj.mode = (0, signing_1.signModeToJSON)(message.mode));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModeInfo_Single();
        message.mode = object.mode ?? 0;
        return message;
    },
};
function createBaseModeInfo_Multi() {
    return {
        bitarray: undefined,
        modeInfos: [],
    };
}
exports.ModeInfo_Multi = {
    typeUrl: "/cosmos.tx.v1beta1.Multi",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.bitarray !== undefined) {
            multisig_1.CompactBitArray.encode(message.bitarray, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.modeInfos) {
            exports.ModeInfo.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseModeInfo_Multi();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.bitarray = multisig_1.CompactBitArray.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.modeInfos.push(exports.ModeInfo.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseModeInfo_Multi();
        if ((0, helpers_1.isSet)(object.bitarray))
            obj.bitarray = multisig_1.CompactBitArray.fromJSON(object.bitarray);
        if (Array.isArray(object?.modeInfos))
            obj.modeInfos = object.modeInfos.map((e) => exports.ModeInfo.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.bitarray !== undefined &&
            (obj.bitarray = message.bitarray ? multisig_1.CompactBitArray.toJSON(message.bitarray) : undefined);
        if (message.modeInfos) {
            obj.modeInfos = message.modeInfos.map((e) => (e ? exports.ModeInfo.toJSON(e) : undefined));
        }
        else {
            obj.modeInfos = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseModeInfo_Multi();
        if (object.bitarray !== undefined && object.bitarray !== null) {
            message.bitarray = multisig_1.CompactBitArray.fromPartial(object.bitarray);
        }
        message.modeInfos = object.modeInfos?.map((e) => exports.ModeInfo.fromPartial(e)) || [];
        return message;
    },
};
function createBaseFee() {
    return {
        amount: [],
        gasLimit: BigInt(0),
        payer: "",
        granter: "",
    };
}
exports.Fee = {
    typeUrl: "/cosmos.tx.v1beta1.Fee",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.amount) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.gasLimit !== BigInt(0)) {
            writer.uint32(16).uint64(message.gasLimit);
        }
        if (message.payer !== "") {
            writer.uint32(26).string(message.payer);
        }
        if (message.granter !== "") {
            writer.uint32(34).string(message.granter);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseFee();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.amount.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.gasLimit = reader.uint64();
                    break;
                case 3:
                    message.payer = reader.string();
                    break;
                case 4:
                    message.granter = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseFee();
        if (Array.isArray(object?.amount))
            obj.amount = object.amount.map((e) => coin_1.Coin.fromJSON(e));
        if ((0, helpers_1.isSet)(object.gasLimit))
            obj.gasLimit = BigInt(object.gasLimit.toString());
        if ((0, helpers_1.isSet)(object.payer))
            obj.payer = String(object.payer);
        if ((0, helpers_1.isSet)(object.granter))
            obj.granter = String(object.granter);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.amount) {
            obj.amount = message.amount.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.amount = [];
        }
        message.gasLimit !== undefined && (obj.gasLimit = (message.gasLimit || BigInt(0)).toString());
        message.payer !== undefined && (obj.payer = message.payer);
        message.granter !== undefined && (obj.granter = message.granter);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseFee();
        message.amount = object.amount?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        if (object.gasLimit !== undefined && object.gasLimit !== null) {
            message.gasLimit = BigInt(object.gasLimit.toString());
        }
        message.payer = object.payer ?? "";
        message.granter = object.granter ?? "";
        return message;
    },
};
function createBaseTip() {
    return {
        amount: [],
        tipper: "",
    };
}
exports.Tip = {
    typeUrl: "/cosmos.tx.v1beta1.Tip",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.amount) {
            coin_1.Coin.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.tipper !== "") {
            writer.uint32(18).string(message.tipper);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTip();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.amount.push(coin_1.Coin.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.tipper = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTip();
        if (Array.isArray(object?.amount))
            obj.amount = object.amount.map((e) => coin_1.Coin.fromJSON(e));
        if ((0, helpers_1.isSet)(object.tipper))
            obj.tipper = String(object.tipper);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.amount) {
            obj.amount = message.amount.map((e) => (e ? coin_1.Coin.toJSON(e) : undefined));
        }
        else {
            obj.amount = [];
        }
        message.tipper !== undefined && (obj.tipper = message.tipper);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTip();
        message.amount = object.amount?.map((e) => coin_1.Coin.fromPartial(e)) || [];
        message.tipper = object.tipper ?? "";
        return message;
    },
};
function createBaseAuxSignerData() {
    return {
        address: "",
        signDoc: undefined,
        mode: 0,
        sig: new Uint8Array(),
    };
}
exports.AuxSignerData = {
    typeUrl: "/cosmos.tx.v1beta1.AuxSignerData",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.address !== "") {
            writer.uint32(10).string(message.address);
        }
        if (message.signDoc !== undefined) {
            exports.SignDocDirectAux.encode(message.signDoc, writer.uint32(18).fork()).ldelim();
        }
        if (message.mode !== 0) {
            writer.uint32(24).int32(message.mode);
        }
        if (message.sig.length !== 0) {
            writer.uint32(34).bytes(message.sig);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseAuxSignerData();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.address = reader.string();
                    break;
                case 2:
                    message.signDoc = exports.SignDocDirectAux.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.mode = reader.int32();
                    break;
                case 4:
                    message.sig = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseAuxSignerData();
        if ((0, helpers_1.isSet)(object.address))
            obj.address = String(object.address);
        if ((0, helpers_1.isSet)(object.signDoc))
            obj.signDoc = exports.SignDocDirectAux.fromJSON(object.signDoc);
        if ((0, helpers_1.isSet)(object.mode))
            obj.mode = (0, signing_1.signModeFromJSON)(object.mode);
        if ((0, helpers_1.isSet)(object.sig))
            obj.sig = (0, helpers_1.bytesFromBase64)(object.sig);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.address !== undefined && (obj.address = message.address);
        message.signDoc !== undefined &&
            (obj.signDoc = message.signDoc ? exports.SignDocDirectAux.toJSON(message.signDoc) : undefined);
        message.mode !== undefined && (obj.mode = (0, signing_1.signModeToJSON)(message.mode));
        message.sig !== undefined &&
            (obj.sig = (0, helpers_1.base64FromBytes)(message.sig !== undefined ? message.sig : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseAuxSignerData();
        message.address = object.address ?? "";
        if (object.signDoc !== undefined && object.signDoc !== null) {
            message.signDoc = exports.SignDocDirectAux.fromPartial(object.signDoc);
        }
        message.mode = object.mode ?? 0;
        message.sig = object.sig ?? new Uint8Array();
        return message;
    },
};
//# sourceMappingURL=tx.js.map