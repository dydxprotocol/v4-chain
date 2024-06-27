"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.QueryClientImpl = exports.QueryAllEvidenceResponse = exports.QueryAllEvidenceRequest = exports.QueryEvidenceResponse = exports.QueryEvidenceRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const pagination_1 = require("../../base/query/v1beta1/pagination");
const any_1 = require("../../../google/protobuf/any");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.evidence.v1beta1";
function createBaseQueryEvidenceRequest() {
    return {
        evidenceHash: new Uint8Array(),
        hash: "",
    };
}
exports.QueryEvidenceRequest = {
    typeUrl: "/cosmos.evidence.v1beta1.QueryEvidenceRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.evidenceHash.length !== 0) {
            writer.uint32(10).bytes(message.evidenceHash);
        }
        if (message.hash !== "") {
            writer.uint32(18).string(message.hash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryEvidenceRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.evidenceHash = reader.bytes();
                    break;
                case 2:
                    message.hash = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryEvidenceRequest();
        if ((0, helpers_1.isSet)(object.evidenceHash))
            obj.evidenceHash = (0, helpers_1.bytesFromBase64)(object.evidenceHash);
        if ((0, helpers_1.isSet)(object.hash))
            obj.hash = String(object.hash);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.evidenceHash !== undefined &&
            (obj.evidenceHash = (0, helpers_1.base64FromBytes)(message.evidenceHash !== undefined ? message.evidenceHash : new Uint8Array()));
        message.hash !== undefined && (obj.hash = message.hash);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryEvidenceRequest();
        message.evidenceHash = object.evidenceHash ?? new Uint8Array();
        message.hash = object.hash ?? "";
        return message;
    },
};
function createBaseQueryEvidenceResponse() {
    return {
        evidence: undefined,
    };
}
exports.QueryEvidenceResponse = {
    typeUrl: "/cosmos.evidence.v1beta1.QueryEvidenceResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.evidence !== undefined) {
            any_1.Any.encode(message.evidence, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryEvidenceResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.evidence = any_1.Any.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const obj = createBaseQueryEvidenceResponse();
        if ((0, helpers_1.isSet)(object.evidence))
            obj.evidence = any_1.Any.fromJSON(object.evidence);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.evidence !== undefined &&
            (obj.evidence = message.evidence ? any_1.Any.toJSON(message.evidence) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryEvidenceResponse();
        if (object.evidence !== undefined && object.evidence !== null) {
            message.evidence = any_1.Any.fromPartial(object.evidence);
        }
        return message;
    },
};
function createBaseQueryAllEvidenceRequest() {
    return {
        pagination: undefined,
    };
}
exports.QueryAllEvidenceRequest = {
    typeUrl: "/cosmos.evidence.v1beta1.QueryAllEvidenceRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryAllEvidenceRequest();
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
        const obj = createBaseQueryAllEvidenceRequest();
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
        const message = createBaseQueryAllEvidenceRequest();
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryAllEvidenceResponse() {
    return {
        evidence: [],
        pagination: undefined,
    };
}
exports.QueryAllEvidenceResponse = {
    typeUrl: "/cosmos.evidence.v1beta1.QueryAllEvidenceResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.evidence) {
            any_1.Any.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryAllEvidenceResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.evidence.push(any_1.Any.decode(reader, reader.uint32()));
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
        const obj = createBaseQueryAllEvidenceResponse();
        if (Array.isArray(object?.evidence))
            obj.evidence = object.evidence.map((e) => any_1.Any.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.evidence) {
            obj.evidence = message.evidence.map((e) => (e ? any_1.Any.toJSON(e) : undefined));
        }
        else {
            obj.evidence = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryAllEvidenceResponse();
        message.evidence = object.evidence?.map((e) => any_1.Any.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.Evidence = this.Evidence.bind(this);
        this.AllEvidence = this.AllEvidence.bind(this);
    }
    Evidence(request) {
        const data = exports.QueryEvidenceRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.evidence.v1beta1.Query", "Evidence", data);
        return promise.then((data) => exports.QueryEvidenceResponse.decode(new binary_1.BinaryReader(data)));
    }
    AllEvidence(request = {
        pagination: pagination_1.PageRequest.fromPartial({}),
    }) {
        const data = exports.QueryAllEvidenceRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.evidence.v1beta1.Query", "AllEvidence", data);
        return promise.then((data) => exports.QueryAllEvidenceResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.QueryClientImpl = QueryClientImpl;
//# sourceMappingURL=query.js.map