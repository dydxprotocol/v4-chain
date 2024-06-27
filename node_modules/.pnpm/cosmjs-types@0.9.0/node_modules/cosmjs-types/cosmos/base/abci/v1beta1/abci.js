"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SearchTxsResult = exports.TxMsgData = exports.MsgData = exports.SimulationResponse = exports.Result = exports.GasInfo = exports.Attribute = exports.StringEvent = exports.ABCIMessageLog = exports.TxResponse = exports.protobufPackage = void 0;
/* eslint-disable */
const any_1 = require("../../../../google/protobuf/any");
const types_1 = require("../../../../tendermint/abci/types");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.base.abci.v1beta1";
function createBaseTxResponse() {
    return {
        height: BigInt(0),
        txhash: "",
        codespace: "",
        code: 0,
        data: "",
        rawLog: "",
        logs: [],
        info: "",
        gasWanted: BigInt(0),
        gasUsed: BigInt(0),
        tx: undefined,
        timestamp: "",
        events: [],
    };
}
exports.TxResponse = {
    typeUrl: "/cosmos.base.abci.v1beta1.TxResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.height !== BigInt(0)) {
            writer.uint32(8).int64(message.height);
        }
        if (message.txhash !== "") {
            writer.uint32(18).string(message.txhash);
        }
        if (message.codespace !== "") {
            writer.uint32(26).string(message.codespace);
        }
        if (message.code !== 0) {
            writer.uint32(32).uint32(message.code);
        }
        if (message.data !== "") {
            writer.uint32(42).string(message.data);
        }
        if (message.rawLog !== "") {
            writer.uint32(50).string(message.rawLog);
        }
        for (const v of message.logs) {
            exports.ABCIMessageLog.encode(v, writer.uint32(58).fork()).ldelim();
        }
        if (message.info !== "") {
            writer.uint32(66).string(message.info);
        }
        if (message.gasWanted !== BigInt(0)) {
            writer.uint32(72).int64(message.gasWanted);
        }
        if (message.gasUsed !== BigInt(0)) {
            writer.uint32(80).int64(message.gasUsed);
        }
        if (message.tx !== undefined) {
            any_1.Any.encode(message.tx, writer.uint32(90).fork()).ldelim();
        }
        if (message.timestamp !== "") {
            writer.uint32(98).string(message.timestamp);
        }
        for (const v of message.events) {
            types_1.Event.encode(v, writer.uint32(106).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = reader.int64();
                    break;
                case 2:
                    message.txhash = reader.string();
                    break;
                case 3:
                    message.codespace = reader.string();
                    break;
                case 4:
                    message.code = reader.uint32();
                    break;
                case 5:
                    message.data = reader.string();
                    break;
                case 6:
                    message.rawLog = reader.string();
                    break;
                case 7:
                    message.logs.push(exports.ABCIMessageLog.decode(reader, reader.uint32()));
                    break;
                case 8:
                    message.info = reader.string();
                    break;
                case 9:
                    message.gasWanted = reader.int64();
                    break;
                case 10:
                    message.gasUsed = reader.int64();
                    break;
                case 11:
                    message.tx = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 12:
                    message.timestamp = reader.string();
                    break;
                case 13:
                    message.events.push(types_1.Event.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxResponse();
        if ((0, helpers_1.isSet)(object.height))
            obj.height = BigInt(object.height.toString());
        if ((0, helpers_1.isSet)(object.txhash))
            obj.txhash = String(object.txhash);
        if ((0, helpers_1.isSet)(object.codespace))
            obj.codespace = String(object.codespace);
        if ((0, helpers_1.isSet)(object.code))
            obj.code = Number(object.code);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = String(object.data);
        if ((0, helpers_1.isSet)(object.rawLog))
            obj.rawLog = String(object.rawLog);
        if (Array.isArray(object?.logs))
            obj.logs = object.logs.map((e) => exports.ABCIMessageLog.fromJSON(e));
        if ((0, helpers_1.isSet)(object.info))
            obj.info = String(object.info);
        if ((0, helpers_1.isSet)(object.gasWanted))
            obj.gasWanted = BigInt(object.gasWanted.toString());
        if ((0, helpers_1.isSet)(object.gasUsed))
            obj.gasUsed = BigInt(object.gasUsed.toString());
        if ((0, helpers_1.isSet)(object.tx))
            obj.tx = any_1.Any.fromJSON(object.tx);
        if ((0, helpers_1.isSet)(object.timestamp))
            obj.timestamp = String(object.timestamp);
        if (Array.isArray(object?.events))
            obj.events = object.events.map((e) => types_1.Event.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.height !== undefined && (obj.height = (message.height || BigInt(0)).toString());
        message.txhash !== undefined && (obj.txhash = message.txhash);
        message.codespace !== undefined && (obj.codespace = message.codespace);
        message.code !== undefined && (obj.code = Math.round(message.code));
        message.data !== undefined && (obj.data = message.data);
        message.rawLog !== undefined && (obj.rawLog = message.rawLog);
        if (message.logs) {
            obj.logs = message.logs.map((e) => (e ? exports.ABCIMessageLog.toJSON(e) : undefined));
        }
        else {
            obj.logs = [];
        }
        message.info !== undefined && (obj.info = message.info);
        message.gasWanted !== undefined && (obj.gasWanted = (message.gasWanted || BigInt(0)).toString());
        message.gasUsed !== undefined && (obj.gasUsed = (message.gasUsed || BigInt(0)).toString());
        message.tx !== undefined && (obj.tx = message.tx ? any_1.Any.toJSON(message.tx) : undefined);
        message.timestamp !== undefined && (obj.timestamp = message.timestamp);
        if (message.events) {
            obj.events = message.events.map((e) => (e ? types_1.Event.toJSON(e) : undefined));
        }
        else {
            obj.events = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxResponse();
        if (object.height !== undefined && object.height !== null) {
            message.height = BigInt(object.height.toString());
        }
        message.txhash = object.txhash ?? "";
        message.codespace = object.codespace ?? "";
        message.code = object.code ?? 0;
        message.data = object.data ?? "";
        message.rawLog = object.rawLog ?? "";
        message.logs = object.logs?.map((e) => exports.ABCIMessageLog.fromPartial(e)) || [];
        message.info = object.info ?? "";
        if (object.gasWanted !== undefined && object.gasWanted !== null) {
            message.gasWanted = BigInt(object.gasWanted.toString());
        }
        if (object.gasUsed !== undefined && object.gasUsed !== null) {
            message.gasUsed = BigInt(object.gasUsed.toString());
        }
        if (object.tx !== undefined && object.tx !== null) {
            message.tx = any_1.Any.fromPartial(object.tx);
        }
        message.timestamp = object.timestamp ?? "";
        message.events = object.events?.map((e) => types_1.Event.fromPartial(e)) || [];
        return message;
    },
};
function createBaseABCIMessageLog() {
    return {
        msgIndex: 0,
        log: "",
        events: [],
    };
}
exports.ABCIMessageLog = {
    typeUrl: "/cosmos.base.abci.v1beta1.ABCIMessageLog",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.msgIndex !== 0) {
            writer.uint32(8).uint32(message.msgIndex);
        }
        if (message.log !== "") {
            writer.uint32(18).string(message.log);
        }
        for (const v of message.events) {
            exports.StringEvent.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseABCIMessageLog();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.msgIndex = reader.uint32();
                    break;
                case 2:
                    message.log = reader.string();
                    break;
                case 3:
                    message.events.push(exports.StringEvent.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseABCIMessageLog();
        if ((0, helpers_1.isSet)(object.msgIndex))
            obj.msgIndex = Number(object.msgIndex);
        if ((0, helpers_1.isSet)(object.log))
            obj.log = String(object.log);
        if (Array.isArray(object?.events))
            obj.events = object.events.map((e) => exports.StringEvent.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.msgIndex !== undefined && (obj.msgIndex = Math.round(message.msgIndex));
        message.log !== undefined && (obj.log = message.log);
        if (message.events) {
            obj.events = message.events.map((e) => (e ? exports.StringEvent.toJSON(e) : undefined));
        }
        else {
            obj.events = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseABCIMessageLog();
        message.msgIndex = object.msgIndex ?? 0;
        message.log = object.log ?? "";
        message.events = object.events?.map((e) => exports.StringEvent.fromPartial(e)) || [];
        return message;
    },
};
function createBaseStringEvent() {
    return {
        type: "",
        attributes: [],
    };
}
exports.StringEvent = {
    typeUrl: "/cosmos.base.abci.v1beta1.StringEvent",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.type !== "") {
            writer.uint32(10).string(message.type);
        }
        for (const v of message.attributes) {
            exports.Attribute.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseStringEvent();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.type = reader.string();
                    break;
                case 2:
                    message.attributes.push(exports.Attribute.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseStringEvent();
        if ((0, helpers_1.isSet)(object.type))
            obj.type = String(object.type);
        if (Array.isArray(object?.attributes))
            obj.attributes = object.attributes.map((e) => exports.Attribute.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.type !== undefined && (obj.type = message.type);
        if (message.attributes) {
            obj.attributes = message.attributes.map((e) => (e ? exports.Attribute.toJSON(e) : undefined));
        }
        else {
            obj.attributes = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseStringEvent();
        message.type = object.type ?? "";
        message.attributes = object.attributes?.map((e) => exports.Attribute.fromPartial(e)) || [];
        return message;
    },
};
function createBaseAttribute() {
    return {
        key: "",
        value: "",
    };
}
exports.Attribute = {
    typeUrl: "/cosmos.base.abci.v1beta1.Attribute",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.key !== "") {
            writer.uint32(10).string(message.key);
        }
        if (message.value !== "") {
            writer.uint32(18).string(message.value);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseAttribute();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.key = reader.string();
                    break;
                case 2:
                    message.value = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseAttribute();
        if ((0, helpers_1.isSet)(object.key))
            obj.key = String(object.key);
        if ((0, helpers_1.isSet)(object.value))
            obj.value = String(object.value);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.key !== undefined && (obj.key = message.key);
        message.value !== undefined && (obj.value = message.value);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseAttribute();
        message.key = object.key ?? "";
        message.value = object.value ?? "";
        return message;
    },
};
function createBaseGasInfo() {
    return {
        gasWanted: BigInt(0),
        gasUsed: BigInt(0),
    };
}
exports.GasInfo = {
    typeUrl: "/cosmos.base.abci.v1beta1.GasInfo",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.gasWanted !== BigInt(0)) {
            writer.uint32(8).uint64(message.gasWanted);
        }
        if (message.gasUsed !== BigInt(0)) {
            writer.uint32(16).uint64(message.gasUsed);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGasInfo();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.gasWanted = reader.uint64();
                    break;
                case 2:
                    message.gasUsed = reader.uint64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGasInfo();
        if ((0, helpers_1.isSet)(object.gasWanted))
            obj.gasWanted = BigInt(object.gasWanted.toString());
        if ((0, helpers_1.isSet)(object.gasUsed))
            obj.gasUsed = BigInt(object.gasUsed.toString());
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.gasWanted !== undefined && (obj.gasWanted = (message.gasWanted || BigInt(0)).toString());
        message.gasUsed !== undefined && (obj.gasUsed = (message.gasUsed || BigInt(0)).toString());
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGasInfo();
        if (object.gasWanted !== undefined && object.gasWanted !== null) {
            message.gasWanted = BigInt(object.gasWanted.toString());
        }
        if (object.gasUsed !== undefined && object.gasUsed !== null) {
            message.gasUsed = BigInt(object.gasUsed.toString());
        }
        return message;
    },
};
function createBaseResult() {
    return {
        data: new Uint8Array(),
        log: "",
        events: [],
        msgResponses: [],
    };
}
exports.Result = {
    typeUrl: "/cosmos.base.abci.v1beta1.Result",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.data.length !== 0) {
            writer.uint32(10).bytes(message.data);
        }
        if (message.log !== "") {
            writer.uint32(18).string(message.log);
        }
        for (const v of message.events) {
            types_1.Event.encode(v, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.msgResponses) {
            any_1.Any.encode(v, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseResult();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.data = reader.bytes();
                    break;
                case 2:
                    message.log = reader.string();
                    break;
                case 3:
                    message.events.push(types_1.Event.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.msgResponses.push(any_1.Any.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseResult();
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        if ((0, helpers_1.isSet)(object.log))
            obj.log = String(object.log);
        if (Array.isArray(object?.events))
            obj.events = object.events.map((e) => types_1.Event.fromJSON(e));
        if (Array.isArray(object?.msgResponses))
            obj.msgResponses = object.msgResponses.map((e) => any_1.Any.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        message.log !== undefined && (obj.log = message.log);
        if (message.events) {
            obj.events = message.events.map((e) => (e ? types_1.Event.toJSON(e) : undefined));
        }
        else {
            obj.events = [];
        }
        if (message.msgResponses) {
            obj.msgResponses = message.msgResponses.map((e) => (e ? any_1.Any.toJSON(e) : undefined));
        }
        else {
            obj.msgResponses = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseResult();
        message.data = object.data ?? new Uint8Array();
        message.log = object.log ?? "";
        message.events = object.events?.map((e) => types_1.Event.fromPartial(e)) || [];
        message.msgResponses = object.msgResponses?.map((e) => any_1.Any.fromPartial(e)) || [];
        return message;
    },
};
function createBaseSimulationResponse() {
    return {
        gasInfo: exports.GasInfo.fromPartial({}),
        result: undefined,
    };
}
exports.SimulationResponse = {
    typeUrl: "/cosmos.base.abci.v1beta1.SimulationResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.gasInfo !== undefined) {
            exports.GasInfo.encode(message.gasInfo, writer.uint32(10).fork()).ldelim();
        }
        if (message.result !== undefined) {
            exports.Result.encode(message.result, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSimulationResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.gasInfo = exports.GasInfo.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.result = exports.Result.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSimulationResponse();
        if ((0, helpers_1.isSet)(object.gasInfo))
            obj.gasInfo = exports.GasInfo.fromJSON(object.gasInfo);
        if ((0, helpers_1.isSet)(object.result))
            obj.result = exports.Result.fromJSON(object.result);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.gasInfo !== undefined &&
            (obj.gasInfo = message.gasInfo ? exports.GasInfo.toJSON(message.gasInfo) : undefined);
        message.result !== undefined && (obj.result = message.result ? exports.Result.toJSON(message.result) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSimulationResponse();
        if (object.gasInfo !== undefined && object.gasInfo !== null) {
            message.gasInfo = exports.GasInfo.fromPartial(object.gasInfo);
        }
        if (object.result !== undefined && object.result !== null) {
            message.result = exports.Result.fromPartial(object.result);
        }
        return message;
    },
};
function createBaseMsgData() {
    return {
        msgType: "",
        data: new Uint8Array(),
    };
}
exports.MsgData = {
    typeUrl: "/cosmos.base.abci.v1beta1.MsgData",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.msgType !== "") {
            writer.uint32(10).string(message.msgType);
        }
        if (message.data.length !== 0) {
            writer.uint32(18).bytes(message.data);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMsgData();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.msgType = reader.string();
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
        const obj = createBaseMsgData();
        if ((0, helpers_1.isSet)(object.msgType))
            obj.msgType = String(object.msgType);
        if ((0, helpers_1.isSet)(object.data))
            obj.data = (0, helpers_1.bytesFromBase64)(object.data);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.msgType !== undefined && (obj.msgType = message.msgType);
        message.data !== undefined &&
            (obj.data = (0, helpers_1.base64FromBytes)(message.data !== undefined ? message.data : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseMsgData();
        message.msgType = object.msgType ?? "";
        message.data = object.data ?? new Uint8Array();
        return message;
    },
};
function createBaseTxMsgData() {
    return {
        data: [],
        msgResponses: [],
    };
}
exports.TxMsgData = {
    typeUrl: "/cosmos.base.abci.v1beta1.TxMsgData",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.data) {
            exports.MsgData.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.msgResponses) {
            any_1.Any.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTxMsgData();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.data.push(exports.MsgData.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.msgResponses.push(any_1.Any.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseTxMsgData();
        if (Array.isArray(object?.data))
            obj.data = object.data.map((e) => exports.MsgData.fromJSON(e));
        if (Array.isArray(object?.msgResponses))
            obj.msgResponses = object.msgResponses.map((e) => any_1.Any.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.data) {
            obj.data = message.data.map((e) => (e ? exports.MsgData.toJSON(e) : undefined));
        }
        else {
            obj.data = [];
        }
        if (message.msgResponses) {
            obj.msgResponses = message.msgResponses.map((e) => (e ? any_1.Any.toJSON(e) : undefined));
        }
        else {
            obj.msgResponses = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseTxMsgData();
        message.data = object.data?.map((e) => exports.MsgData.fromPartial(e)) || [];
        message.msgResponses = object.msgResponses?.map((e) => any_1.Any.fromPartial(e)) || [];
        return message;
    },
};
function createBaseSearchTxsResult() {
    return {
        totalCount: BigInt(0),
        count: BigInt(0),
        pageNumber: BigInt(0),
        pageTotal: BigInt(0),
        limit: BigInt(0),
        txs: [],
    };
}
exports.SearchTxsResult = {
    typeUrl: "/cosmos.base.abci.v1beta1.SearchTxsResult",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.totalCount !== BigInt(0)) {
            writer.uint32(8).uint64(message.totalCount);
        }
        if (message.count !== BigInt(0)) {
            writer.uint32(16).uint64(message.count);
        }
        if (message.pageNumber !== BigInt(0)) {
            writer.uint32(24).uint64(message.pageNumber);
        }
        if (message.pageTotal !== BigInt(0)) {
            writer.uint32(32).uint64(message.pageTotal);
        }
        if (message.limit !== BigInt(0)) {
            writer.uint32(40).uint64(message.limit);
        }
        for (const v of message.txs) {
            exports.TxResponse.encode(v, writer.uint32(50).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSearchTxsResult();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.totalCount = reader.uint64();
                    break;
                case 2:
                    message.count = reader.uint64();
                    break;
                case 3:
                    message.pageNumber = reader.uint64();
                    break;
                case 4:
                    message.pageTotal = reader.uint64();
                    break;
                case 5:
                    message.limit = reader.uint64();
                    break;
                case 6:
                    message.txs.push(exports.TxResponse.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseSearchTxsResult();
        if ((0, helpers_1.isSet)(object.totalCount))
            obj.totalCount = BigInt(object.totalCount.toString());
        if ((0, helpers_1.isSet)(object.count))
            obj.count = BigInt(object.count.toString());
        if ((0, helpers_1.isSet)(object.pageNumber))
            obj.pageNumber = BigInt(object.pageNumber.toString());
        if ((0, helpers_1.isSet)(object.pageTotal))
            obj.pageTotal = BigInt(object.pageTotal.toString());
        if ((0, helpers_1.isSet)(object.limit))
            obj.limit = BigInt(object.limit.toString());
        if (Array.isArray(object?.txs))
            obj.txs = object.txs.map((e) => exports.TxResponse.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.totalCount !== undefined && (obj.totalCount = (message.totalCount || BigInt(0)).toString());
        message.count !== undefined && (obj.count = (message.count || BigInt(0)).toString());
        message.pageNumber !== undefined && (obj.pageNumber = (message.pageNumber || BigInt(0)).toString());
        message.pageTotal !== undefined && (obj.pageTotal = (message.pageTotal || BigInt(0)).toString());
        message.limit !== undefined && (obj.limit = (message.limit || BigInt(0)).toString());
        if (message.txs) {
            obj.txs = message.txs.map((e) => (e ? exports.TxResponse.toJSON(e) : undefined));
        }
        else {
            obj.txs = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseSearchTxsResult();
        if (object.totalCount !== undefined && object.totalCount !== null) {
            message.totalCount = BigInt(object.totalCount.toString());
        }
        if (object.count !== undefined && object.count !== null) {
            message.count = BigInt(object.count.toString());
        }
        if (object.pageNumber !== undefined && object.pageNumber !== null) {
            message.pageNumber = BigInt(object.pageNumber.toString());
        }
        if (object.pageTotal !== undefined && object.pageTotal !== null) {
            message.pageTotal = BigInt(object.pageTotal.toString());
        }
        if (object.limit !== undefined && object.limit !== null) {
            message.limit = BigInt(object.limit.toString());
        }
        message.txs = object.txs?.map((e) => exports.TxResponse.fromPartial(e)) || [];
        return message;
    },
};
//# sourceMappingURL=abci.js.map