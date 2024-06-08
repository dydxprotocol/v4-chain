"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.QueryClientImpl = exports.IndexValue = exports.ListResponse = exports.ListRequest_Range = exports.ListRequest_Prefix = exports.ListRequest = exports.GetResponse = exports.GetRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const pagination_1 = require("../../../base/query/v1beta1/pagination");
const any_1 = require("../../../../google/protobuf/any");
const timestamp_1 = require("../../../../google/protobuf/timestamp");
const duration_1 = require("../../../../google/protobuf/duration");
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.orm.query.v1alpha1";
function createBaseGetRequest() {
    return {
        messageName: "",
        index: "",
        values: [],
    };
}
exports.GetRequest = {
    typeUrl: "/cosmos.orm.query.v1alpha1.GetRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.messageName !== "") {
            writer.uint32(10).string(message.messageName);
        }
        if (message.index !== "") {
            writer.uint32(18).string(message.index);
        }
        for (const v of message.values) {
            exports.IndexValue.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.messageName = reader.string();
                    break;
                case 2:
                    message.index = reader.string();
                    break;
                case 3:
                    message.values.push(exports.IndexValue.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetRequest();
        if ((0, helpers_1.isSet)(object.messageName))
            obj.messageName = String(object.messageName);
        if ((0, helpers_1.isSet)(object.index))
            obj.index = String(object.index);
        if (Array.isArray(object?.values))
            obj.values = object.values.map((e) => exports.IndexValue.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.messageName !== undefined && (obj.messageName = message.messageName);
        message.index !== undefined && (obj.index = message.index);
        if (message.values) {
            obj.values = message.values.map((e) => (e ? exports.IndexValue.toJSON(e) : undefined));
        }
        else {
            obj.values = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetRequest();
        message.messageName = object.messageName ?? "";
        message.index = object.index ?? "";
        message.values = object.values?.map((e) => exports.IndexValue.fromPartial(e)) || [];
        return message;
    },
};
function createBaseGetResponse() {
    return {
        result: undefined,
    };
}
exports.GetResponse = {
    typeUrl: "/cosmos.orm.query.v1alpha1.GetResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.result !== undefined) {
            any_1.Any.encode(message.result, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.result = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseGetResponse();
        if ((0, helpers_1.isSet)(object.result))
            obj.result = any_1.Any.fromJSON(object.result);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.result !== undefined && (obj.result = message.result ? any_1.Any.toJSON(message.result) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseGetResponse();
        if (object.result !== undefined && object.result !== null) {
            message.result = any_1.Any.fromPartial(object.result);
        }
        return message;
    },
};
function createBaseListRequest() {
    return {
        messageName: "",
        index: "",
        prefix: undefined,
        range: undefined,
        pagination: undefined,
    };
}
exports.ListRequest = {
    typeUrl: "/cosmos.orm.query.v1alpha1.ListRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.messageName !== "") {
            writer.uint32(10).string(message.messageName);
        }
        if (message.index !== "") {
            writer.uint32(18).string(message.index);
        }
        if (message.prefix !== undefined) {
            exports.ListRequest_Prefix.encode(message.prefix, writer.uint32(26).fork()).ldelim();
        }
        if (message.range !== undefined) {
            exports.ListRequest_Range.encode(message.range, writer.uint32(34).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(42).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseListRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.messageName = reader.string();
                    break;
                case 2:
                    message.index = reader.string();
                    break;
                case 3:
                    message.prefix = exports.ListRequest_Prefix.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.range = exports.ListRequest_Range.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.pagination = pagination_1.PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseListRequest();
        if ((0, helpers_1.isSet)(object.messageName))
            obj.messageName = String(object.messageName);
        if ((0, helpers_1.isSet)(object.index))
            obj.index = String(object.index);
        if ((0, helpers_1.isSet)(object.prefix))
            obj.prefix = exports.ListRequest_Prefix.fromJSON(object.prefix);
        if ((0, helpers_1.isSet)(object.range))
            obj.range = exports.ListRequest_Range.fromJSON(object.range);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.messageName !== undefined && (obj.messageName = message.messageName);
        message.index !== undefined && (obj.index = message.index);
        message.prefix !== undefined &&
            (obj.prefix = message.prefix ? exports.ListRequest_Prefix.toJSON(message.prefix) : undefined);
        message.range !== undefined &&
            (obj.range = message.range ? exports.ListRequest_Range.toJSON(message.range) : undefined);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseListRequest();
        message.messageName = object.messageName ?? "";
        message.index = object.index ?? "";
        if (object.prefix !== undefined && object.prefix !== null) {
            message.prefix = exports.ListRequest_Prefix.fromPartial(object.prefix);
        }
        if (object.range !== undefined && object.range !== null) {
            message.range = exports.ListRequest_Range.fromPartial(object.range);
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseListRequest_Prefix() {
    return {
        values: [],
    };
}
exports.ListRequest_Prefix = {
    typeUrl: "/cosmos.orm.query.v1alpha1.Prefix",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.values) {
            exports.IndexValue.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseListRequest_Prefix();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.values.push(exports.IndexValue.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseListRequest_Prefix();
        if (Array.isArray(object?.values))
            obj.values = object.values.map((e) => exports.IndexValue.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.values) {
            obj.values = message.values.map((e) => (e ? exports.IndexValue.toJSON(e) : undefined));
        }
        else {
            obj.values = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseListRequest_Prefix();
        message.values = object.values?.map((e) => exports.IndexValue.fromPartial(e)) || [];
        return message;
    },
};
function createBaseListRequest_Range() {
    return {
        start: [],
        end: [],
    };
}
exports.ListRequest_Range = {
    typeUrl: "/cosmos.orm.query.v1alpha1.Range",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.start) {
            exports.IndexValue.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.end) {
            exports.IndexValue.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseListRequest_Range();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.start.push(exports.IndexValue.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.end.push(exports.IndexValue.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseListRequest_Range();
        if (Array.isArray(object?.start))
            obj.start = object.start.map((e) => exports.IndexValue.fromJSON(e));
        if (Array.isArray(object?.end))
            obj.end = object.end.map((e) => exports.IndexValue.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.start) {
            obj.start = message.start.map((e) => (e ? exports.IndexValue.toJSON(e) : undefined));
        }
        else {
            obj.start = [];
        }
        if (message.end) {
            obj.end = message.end.map((e) => (e ? exports.IndexValue.toJSON(e) : undefined));
        }
        else {
            obj.end = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseListRequest_Range();
        message.start = object.start?.map((e) => exports.IndexValue.fromPartial(e)) || [];
        message.end = object.end?.map((e) => exports.IndexValue.fromPartial(e)) || [];
        return message;
    },
};
function createBaseListResponse() {
    return {
        results: [],
        pagination: undefined,
    };
}
exports.ListResponse = {
    typeUrl: "/cosmos.orm.query.v1alpha1.ListResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.results) {
            any_1.Any.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(42).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseListResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.results.push(any_1.Any.decode(reader, reader.uint32()));
                    break;
                case 5:
                    message.pagination = pagination_1.PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseListResponse();
        if (Array.isArray(object?.results))
            obj.results = object.results.map((e) => any_1.Any.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.results) {
            obj.results = message.results.map((e) => (e ? any_1.Any.toJSON(e) : undefined));
        }
        else {
            obj.results = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseListResponse();
        message.results = object.results?.map((e) => any_1.Any.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseIndexValue() {
    return {
        uint: undefined,
        int: undefined,
        str: undefined,
        bytes: undefined,
        enum: undefined,
        bool: undefined,
        timestamp: undefined,
        duration: undefined,
    };
}
exports.IndexValue = {
    typeUrl: "/cosmos.orm.query.v1alpha1.IndexValue",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.uint !== undefined) {
            writer.uint32(8).uint64(message.uint);
        }
        if (message.int !== undefined) {
            writer.uint32(16).int64(message.int);
        }
        if (message.str !== undefined) {
            writer.uint32(26).string(message.str);
        }
        if (message.bytes !== undefined) {
            writer.uint32(34).bytes(message.bytes);
        }
        if (message.enum !== undefined) {
            writer.uint32(42).string(message.enum);
        }
        if (message.bool !== undefined) {
            writer.uint32(48).bool(message.bool);
        }
        if (message.timestamp !== undefined) {
            timestamp_1.Timestamp.encode(message.timestamp, writer.uint32(58).fork()).ldelim();
        }
        if (message.duration !== undefined) {
            duration_1.Duration.encode(message.duration, writer.uint32(66).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseIndexValue();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.uint = reader.uint64();
                    break;
                case 2:
                    message.int = reader.int64();
                    break;
                case 3:
                    message.str = reader.string();
                    break;
                case 4:
                    message.bytes = reader.bytes();
                    break;
                case 5:
                    message.enum = reader.string();
                    break;
                case 6:
                    message.bool = reader.bool();
                    break;
                case 7:
                    message.timestamp = timestamp_1.Timestamp.decode(reader, reader.uint32());
                    break;
                case 8:
                    message.duration = duration_1.Duration.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseIndexValue();
        if ((0, helpers_1.isSet)(object.uint))
            obj.uint = BigInt(object.uint.toString());
        if ((0, helpers_1.isSet)(object.int))
            obj.int = BigInt(object.int.toString());
        if ((0, helpers_1.isSet)(object.str))
            obj.str = String(object.str);
        if ((0, helpers_1.isSet)(object.bytes))
            obj.bytes = (0, helpers_1.bytesFromBase64)(object.bytes);
        if ((0, helpers_1.isSet)(object.enum))
            obj.enum = String(object.enum);
        if ((0, helpers_1.isSet)(object.bool))
            obj.bool = Boolean(object.bool);
        if ((0, helpers_1.isSet)(object.timestamp))
            obj.timestamp = (0, helpers_1.fromJsonTimestamp)(object.timestamp);
        if ((0, helpers_1.isSet)(object.duration))
            obj.duration = duration_1.Duration.fromJSON(object.duration);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.uint !== undefined) {
            obj.uint = message.uint.toString();
        }
        if (message.int !== undefined) {
            obj.int = message.int.toString();
        }
        message.str !== undefined && (obj.str = message.str);
        message.bytes !== undefined &&
            (obj.bytes = message.bytes !== undefined ? (0, helpers_1.base64FromBytes)(message.bytes) : undefined);
        message.enum !== undefined && (obj.enum = message.enum);
        message.bool !== undefined && (obj.bool = message.bool);
        message.timestamp !== undefined && (obj.timestamp = (0, helpers_1.fromTimestamp)(message.timestamp).toISOString());
        message.duration !== undefined &&
            (obj.duration = message.duration ? duration_1.Duration.toJSON(message.duration) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseIndexValue();
        if (object.uint !== undefined && object.uint !== null) {
            message.uint = BigInt(object.uint.toString());
        }
        if (object.int !== undefined && object.int !== null) {
            message.int = BigInt(object.int.toString());
        }
        message.str = object.str ?? undefined;
        message.bytes = object.bytes ?? undefined;
        message.enum = object.enum ?? undefined;
        message.bool = object.bool ?? undefined;
        if (object.timestamp !== undefined && object.timestamp !== null) {
            message.timestamp = timestamp_1.Timestamp.fromPartial(object.timestamp);
        }
        if (object.duration !== undefined && object.duration !== null) {
            message.duration = duration_1.Duration.fromPartial(object.duration);
        }
        return message;
    },
};
class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.Get = this.Get.bind(this);
        this.List = this.List.bind(this);
    }
    Get(request) {
        const data = exports.GetRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.orm.query.v1alpha1.Query", "Get", data);
        return promise.then((data) => exports.GetResponse.decode(new binary_1.BinaryReader(data)));
    }
    List(request) {
        const data = exports.ListRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.orm.query.v1alpha1.Query", "List", data);
        return promise.then((data) => exports.ListResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.QueryClientImpl = QueryClientImpl;
//# sourceMappingURL=query.js.map