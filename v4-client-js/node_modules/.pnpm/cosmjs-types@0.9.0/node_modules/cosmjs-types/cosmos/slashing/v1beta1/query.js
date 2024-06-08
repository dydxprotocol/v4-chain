"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.QueryClientImpl = exports.QuerySigningInfosResponse = exports.QuerySigningInfosRequest = exports.QuerySigningInfoResponse = exports.QuerySigningInfoRequest = exports.QueryParamsResponse = exports.QueryParamsRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const pagination_1 = require("../../base/query/v1beta1/pagination");
const slashing_1 = require("./slashing");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.slashing.v1beta1";
function createBaseQueryParamsRequest() {
    return {};
}
exports.QueryParamsRequest = {
    typeUrl: "/cosmos.slashing.v1beta1.QueryParamsRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryParamsRequest();
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
        const obj = createBaseQueryParamsRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseQueryParamsRequest();
        return message;
    },
};
function createBaseQueryParamsResponse() {
    return {
        params: slashing_1.Params.fromPartial({}),
    };
}
exports.QueryParamsResponse = {
    typeUrl: "/cosmos.slashing.v1beta1.QueryParamsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.params !== undefined) {
            slashing_1.Params.encode(message.params, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryParamsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.params = slashing_1.Params.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryParamsResponse();
        if ((0, helpers_1.isSet)(object.params))
            obj.params = slashing_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.params !== undefined && (obj.params = message.params ? slashing_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryParamsResponse();
        if (object.params !== undefined && object.params !== null) {
            message.params = slashing_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
function createBaseQuerySigningInfoRequest() {
    return {
        consAddress: "",
    };
}
exports.QuerySigningInfoRequest = {
    typeUrl: "/cosmos.slashing.v1beta1.QuerySigningInfoRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.consAddress !== "") {
            writer.uint32(10).string(message.consAddress);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQuerySigningInfoRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.consAddress = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQuerySigningInfoRequest();
        if ((0, helpers_1.isSet)(object.consAddress))
            obj.consAddress = String(object.consAddress);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.consAddress !== undefined && (obj.consAddress = message.consAddress);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQuerySigningInfoRequest();
        message.consAddress = object.consAddress ?? "";
        return message;
    },
};
function createBaseQuerySigningInfoResponse() {
    return {
        valSigningInfo: slashing_1.ValidatorSigningInfo.fromPartial({}),
    };
}
exports.QuerySigningInfoResponse = {
    typeUrl: "/cosmos.slashing.v1beta1.QuerySigningInfoResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.valSigningInfo !== undefined) {
            slashing_1.ValidatorSigningInfo.encode(message.valSigningInfo, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQuerySigningInfoResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.valSigningInfo = slashing_1.ValidatorSigningInfo.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQuerySigningInfoResponse();
        if ((0, helpers_1.isSet)(object.valSigningInfo))
            obj.valSigningInfo = slashing_1.ValidatorSigningInfo.fromJSON(object.valSigningInfo);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.valSigningInfo !== undefined &&
            (obj.valSigningInfo = message.valSigningInfo
                ? slashing_1.ValidatorSigningInfo.toJSON(message.valSigningInfo)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQuerySigningInfoResponse();
        if (object.valSigningInfo !== undefined && object.valSigningInfo !== null) {
            message.valSigningInfo = slashing_1.ValidatorSigningInfo.fromPartial(object.valSigningInfo);
        }
        return message;
    },
};
function createBaseQuerySigningInfosRequest() {
    return {
        pagination: undefined,
    };
}
exports.QuerySigningInfosRequest = {
    typeUrl: "/cosmos.slashing.v1beta1.QuerySigningInfosRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQuerySigningInfosRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        const obj = createBaseQuerySigningInfosRequest();
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQuerySigningInfosRequest();
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQuerySigningInfosResponse() {
    return {
        info: [],
        pagination: undefined,
    };
}
exports.QuerySigningInfosResponse = {
    typeUrl: "/cosmos.slashing.v1beta1.QuerySigningInfosResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.info) {
            slashing_1.ValidatorSigningInfo.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQuerySigningInfosResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.info.push(slashing_1.ValidatorSigningInfo.decode(reader, reader.uint32()));
                    break;
                case 2:
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
        const obj = createBaseQuerySigningInfosResponse();
        if (Array.isArray(object?.info))
            obj.info = object.info.map((e) => slashing_1.ValidatorSigningInfo.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.info) {
            obj.info = message.info.map((e) => (e ? slashing_1.ValidatorSigningInfo.toJSON(e) : undefined));
        }
        else {
            obj.info = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQuerySigningInfosResponse();
        message.info = object.info?.map((e) => slashing_1.ValidatorSigningInfo.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.Params = this.Params.bind(this);
        this.SigningInfo = this.SigningInfo.bind(this);
        this.SigningInfos = this.SigningInfos.bind(this);
    }
    Params(request = {}) {
        const data = exports.QueryParamsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.slashing.v1beta1.Query", "Params", data);
        return promise.then((data) => exports.QueryParamsResponse.decode(new binary_1.BinaryReader(data)));
    }
    SigningInfo(request) {
        const data = exports.QuerySigningInfoRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.slashing.v1beta1.Query", "SigningInfo", data);
        return promise.then((data) => exports.QuerySigningInfoResponse.decode(new binary_1.BinaryReader(data)));
    }
    SigningInfos(request = {
        pagination: pagination_1.PageRequest.fromPartial({}),
    }) {
        const data = exports.QuerySigningInfosRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.slashing.v1beta1.Query", "SigningInfos", data);
        return promise.then((data) => exports.QuerySigningInfosResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.QueryClientImpl = QueryClientImpl;
//# sourceMappingURL=query.js.map