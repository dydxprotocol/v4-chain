"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReflectionServiceClientImpl = exports.FileDescriptorsResponse = exports.FileDescriptorsRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const descriptor_1 = require("../../../google/protobuf/descriptor");
const binary_1 = require("../../../binary");
exports.protobufPackage = "cosmos.reflection.v1";
function createBaseFileDescriptorsRequest() {
    return {};
}
exports.FileDescriptorsRequest = {
    typeUrl: "/cosmos.reflection.v1.FileDescriptorsRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseFileDescriptorsRequest();
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
        const obj = createBaseFileDescriptorsRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseFileDescriptorsRequest();
        return message;
    },
};
function createBaseFileDescriptorsResponse() {
    return {
        files: [],
    };
}
exports.FileDescriptorsResponse = {
    typeUrl: "/cosmos.reflection.v1.FileDescriptorsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.files) {
            descriptor_1.FileDescriptorProto.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseFileDescriptorsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.files.push(descriptor_1.FileDescriptorProto.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseFileDescriptorsResponse();
        if (Array.isArray(object?.files))
            obj.files = object.files.map((e) => descriptor_1.FileDescriptorProto.fromJSON(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.files) {
            obj.files = message.files.map((e) => (e ? descriptor_1.FileDescriptorProto.toJSON(e) : undefined));
        }
        else {
            obj.files = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseFileDescriptorsResponse();
        message.files = object.files?.map((e) => descriptor_1.FileDescriptorProto.fromPartial(e)) || [];
        return message;
    },
};
class ReflectionServiceClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.FileDescriptors = this.FileDescriptors.bind(this);
    }
    FileDescriptors(request = {}) {
        const data = exports.FileDescriptorsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.reflection.v1.ReflectionService", "FileDescriptors", data);
        return promise.then((data) => exports.FileDescriptorsResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.ReflectionServiceClientImpl = ReflectionServiceClientImpl;
//# sourceMappingURL=reflection.js.map