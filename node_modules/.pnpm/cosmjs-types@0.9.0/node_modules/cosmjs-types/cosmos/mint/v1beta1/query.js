"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.QueryClientImpl = exports.QueryAnnualProvisionsResponse = exports.QueryAnnualProvisionsRequest = exports.QueryInflationResponse = exports.QueryInflationRequest = exports.QueryParamsResponse = exports.QueryParamsRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const mint_1 = require("./mint");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.mint.v1beta1";
function createBaseQueryParamsRequest() {
    return {};
}
exports.QueryParamsRequest = {
    typeUrl: "/cosmos.mint.v1beta1.QueryParamsRequest",
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
        params: mint_1.Params.fromPartial({}),
    };
}
exports.QueryParamsResponse = {
    typeUrl: "/cosmos.mint.v1beta1.QueryParamsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.params !== undefined) {
            mint_1.Params.encode(message.params, writer.uint32(10).fork()).ldelim();
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
                    message.params = mint_1.Params.decode(reader, reader.uint32());
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
            obj.params = mint_1.Params.fromJSON(object.params);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.params !== undefined && (obj.params = message.params ? mint_1.Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryParamsResponse();
        if (object.params !== undefined && object.params !== null) {
            message.params = mint_1.Params.fromPartial(object.params);
        }
        return message;
    },
};
function createBaseQueryInflationRequest() {
    return {};
}
exports.QueryInflationRequest = {
    typeUrl: "/cosmos.mint.v1beta1.QueryInflationRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryInflationRequest();
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
        const obj = createBaseQueryInflationRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseQueryInflationRequest();
        return message;
    },
};
function createBaseQueryInflationResponse() {
    return {
        inflation: new Uint8Array(),
    };
}
exports.QueryInflationResponse = {
    typeUrl: "/cosmos.mint.v1beta1.QueryInflationResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.inflation.length !== 0) {
            writer.uint32(10).bytes(message.inflation);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryInflationResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.inflation = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryInflationResponse();
        if ((0, helpers_1.isSet)(object.inflation))
            obj.inflation = (0, helpers_1.bytesFromBase64)(object.inflation);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.inflation !== undefined &&
            (obj.inflation = (0, helpers_1.base64FromBytes)(message.inflation !== undefined ? message.inflation : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryInflationResponse();
        message.inflation = object.inflation ?? new Uint8Array();
        return message;
    },
};
function createBaseQueryAnnualProvisionsRequest() {
    return {};
}
exports.QueryAnnualProvisionsRequest = {
    typeUrl: "/cosmos.mint.v1beta1.QueryAnnualProvisionsRequest",
    encode(_, writer = binary_1.BinaryWriter.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryAnnualProvisionsRequest();
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
        const obj = createBaseQueryAnnualProvisionsRequest();
        return obj;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = createBaseQueryAnnualProvisionsRequest();
        return message;
    },
};
function createBaseQueryAnnualProvisionsResponse() {
    return {
        annualProvisions: new Uint8Array(),
    };
}
exports.QueryAnnualProvisionsResponse = {
    typeUrl: "/cosmos.mint.v1beta1.QueryAnnualProvisionsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.annualProvisions.length !== 0) {
            writer.uint32(10).bytes(message.annualProvisions);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryAnnualProvisionsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.annualProvisions = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryAnnualProvisionsResponse();
        if ((0, helpers_1.isSet)(object.annualProvisions))
            obj.annualProvisions = (0, helpers_1.bytesFromBase64)(object.annualProvisions);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.annualProvisions !== undefined &&
            (obj.annualProvisions = (0, helpers_1.base64FromBytes)(message.annualProvisions !== undefined ? message.annualProvisions : new Uint8Array()));
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryAnnualProvisionsResponse();
        message.annualProvisions = object.annualProvisions ?? new Uint8Array();
        return message;
    },
};
class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.Params = this.Params.bind(this);
        this.Inflation = this.Inflation.bind(this);
        this.AnnualProvisions = this.AnnualProvisions.bind(this);
    }
    Params(request = {}) {
        const data = exports.QueryParamsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.mint.v1beta1.Query", "Params", data);
        return promise.then((data) => exports.QueryParamsResponse.decode(new binary_1.BinaryReader(data)));
    }
    Inflation(request = {}) {
        const data = exports.QueryInflationRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.mint.v1beta1.Query", "Inflation", data);
        return promise.then((data) => exports.QueryInflationResponse.decode(new binary_1.BinaryReader(data)));
    }
    AnnualProvisions(request = {}) {
        const data = exports.QueryAnnualProvisionsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.mint.v1beta1.Query", "AnnualProvisions", data);
        return promise.then((data) => exports.QueryAnnualProvisionsResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.QueryClientImpl = QueryClientImpl;
//# sourceMappingURL=query.js.map