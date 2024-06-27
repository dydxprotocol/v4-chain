"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.QueryClientImpl = exports.QueryConfigResponse = exports.QueryConfigRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const config_1 = require("./config");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.app.v1alpha1";
function createBaseQueryConfigRequest() {
    return {};
}
exports.QueryConfigRequest = {
    typeUrl: "/cosmos.app.v1alpha1.QueryConfigRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryConfigRequest();
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
        const obj = createBaseQueryConfigRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseQueryConfigRequest();
        return message;
    },
};
function createBaseQueryConfigResponse() {
    return {
        config: undefined,
    };
}
exports.QueryConfigResponse = {
    typeUrl: "/cosmos.app.v1alpha1.QueryConfigResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.config !== undefined) {
            config_1.Config.encode(message.config, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryConfigResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.config = config_1.Config.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryConfigResponse();
        if ((0, helpers_1.isSet)(object.config))
            obj.config = config_1.Config.fromJSON(object.config);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.config !== undefined && (obj.config = message.config ? config_1.Config.toJSON(message.config) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryConfigResponse();
        if (object.config !== undefined && object.config !== null) {
            message.config = config_1.Config.fromPartial(object.config);
        }
        return message;
    },
};
class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.Config = this.Config.bind(this);
    }
    Config(request = {}) {
        const data = exports.QueryConfigRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.app.v1alpha1.Query", "Config", data);
        return promise.then((data) => exports.QueryConfigResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.QueryClientImpl = QueryClientImpl;
//# sourceMappingURL=query.js.map