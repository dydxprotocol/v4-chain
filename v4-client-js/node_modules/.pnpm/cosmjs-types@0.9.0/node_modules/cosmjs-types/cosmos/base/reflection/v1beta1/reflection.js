"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReflectionServiceClientImpl = exports.ListImplementationsResponse = exports.ListImplementationsRequest = exports.ListAllInterfacesResponse = exports.ListAllInterfacesRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const binary_1 = require("../../../../binary");
const helpers_1 = require("../../../../helpers");
exports.protobufPackage = "cosmos.base.reflection.v1beta1";
function createBaseListAllInterfacesRequest() {
    return {};
}
exports.ListAllInterfacesRequest = {
    typeUrl: "/cosmos.base.reflection.v1beta1.ListAllInterfacesRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseListAllInterfacesRequest();
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
        const obj = createBaseListAllInterfacesRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseListAllInterfacesRequest();
        return message;
    },
};
function createBaseListAllInterfacesResponse() {
    return {
        interfaceNames: [],
    };
}
exports.ListAllInterfacesResponse = {
    typeUrl: "/cosmos.base.reflection.v1beta1.ListAllInterfacesResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.interfaceNames) {
            writer.uint32(10).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseListAllInterfacesResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.interfaceNames.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseListAllInterfacesResponse();
        if (Array.isArray(object?.interfaceNames))
            obj.interfaceNames = object.interfaceNames.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.interfaceNames) {
            obj.interfaceNames = message.interfaceNames.map((e) => e);
        }
        else {
            obj.interfaceNames = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseListAllInterfacesResponse();
        message.interfaceNames = object.interfaceNames?.map((e) => e) || [];
        return message;
    },
};
function createBaseListImplementationsRequest() {
    return {
        interfaceName: "",
    };
}
exports.ListImplementationsRequest = {
    typeUrl: "/cosmos.base.reflection.v1beta1.ListImplementationsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.interfaceName !== "") {
            writer.uint32(10).string(message.interfaceName);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseListImplementationsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.interfaceName = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseListImplementationsRequest();
        if ((0, helpers_1.isSet)(object.interfaceName))
            obj.interfaceName = String(object.interfaceName);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.interfaceName !== undefined && (obj.interfaceName = message.interfaceName);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseListImplementationsRequest();
        message.interfaceName = object.interfaceName ?? "";
        return message;
    },
};
function createBaseListImplementationsResponse() {
    return {
        implementationMessageNames: [],
    };
}
exports.ListImplementationsResponse = {
    typeUrl: "/cosmos.base.reflection.v1beta1.ListImplementationsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.implementationMessageNames) {
            writer.uint32(10).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseListImplementationsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.implementationMessageNames.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseListImplementationsResponse();
        if (Array.isArray(object?.implementationMessageNames))
            obj.implementationMessageNames = object.implementationMessageNames.map((e) => String(e));
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.implementationMessageNames) {
            obj.implementationMessageNames = message.implementationMessageNames.map((e) => e);
        }
        else {
            obj.implementationMessageNames = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = createBaseListImplementationsResponse();
        message.implementationMessageNames = object.implementationMessageNames?.map((e) => e) || [];
        return message;
    },
};
class ReflectionServiceClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.ListAllInterfaces = this.ListAllInterfaces.bind(this);
        this.ListImplementations = this.ListImplementations.bind(this);
    }
    ListAllInterfaces(request = {}) {
        const data = exports.ListAllInterfacesRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.base.reflection.v1beta1.ReflectionService", "ListAllInterfaces", data);
        return promise.then((data) => exports.ListAllInterfacesResponse.decode(new binary_1.BinaryReader(data)));
    }
    ListImplementations(request) {
        const data = exports.ListImplementationsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.base.reflection.v1beta1.ReflectionService", "ListImplementations", data);
        return promise.then((data) => exports.ListImplementationsResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.ReflectionServiceClientImpl = ReflectionServiceClientImpl;
//# sourceMappingURL=reflection.js.map