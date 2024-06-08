"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.QueryClientImpl = exports.QueryGranteeGrantsResponse = exports.QueryGranteeGrantsRequest = exports.QueryGranterGrantsResponse = exports.QueryGranterGrantsRequest = exports.QueryGrantsResponse = exports.QueryGrantsRequest = exports.protobufPackage = void 0;
/* eslint-disable */
const pagination_1 = require("../../base/query/v1beta1/pagination");
const authz_1 = require("./authz");
const binary_1 = require("../../../binary");
const helpers_1 = require("../../../helpers");
exports.protobufPackage = "cosmos.authz.v1beta1";
function createBaseQueryGrantsRequest() {
    return {
        granter: "",
        grantee: "",
        msgTypeUrl: "",
        pagination: undefined,
    };
}
exports.QueryGrantsRequest = {
    typeUrl: "/cosmos.authz.v1beta1.QueryGrantsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.granter !== "") {
            writer.uint32(10).string(message.granter);
        }
        if (message.grantee !== "") {
            writer.uint32(18).string(message.grantee);
        }
        if (message.msgTypeUrl !== "") {
            writer.uint32(26).string(message.msgTypeUrl);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryGrantsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.granter = reader.string();
                    break;
                case 2:
                    message.grantee = reader.string();
                    break;
                case 3:
                    message.msgTypeUrl = reader.string();
                    break;
                case 4:
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
        const obj = createBaseQueryGrantsRequest();
        if ((0, helpers_1.isSet)(object.granter))
            obj.granter = String(object.granter);
        if ((0, helpers_1.isSet)(object.grantee))
            obj.grantee = String(object.grantee);
        if ((0, helpers_1.isSet)(object.msgTypeUrl))
            obj.msgTypeUrl = String(object.msgTypeUrl);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.granter !== undefined && (obj.granter = message.granter);
        message.grantee !== undefined && (obj.grantee = message.grantee);
        message.msgTypeUrl !== undefined && (obj.msgTypeUrl = message.msgTypeUrl);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryGrantsRequest();
        message.granter = object.granter ?? "";
        message.grantee = object.grantee ?? "";
        message.msgTypeUrl = object.msgTypeUrl ?? "";
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryGrantsResponse() {
    return {
        grants: [],
        pagination: undefined,
    };
}
exports.QueryGrantsResponse = {
    typeUrl: "/cosmos.authz.v1beta1.QueryGrantsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.grants) {
            authz_1.Grant.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryGrantsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.grants.push(authz_1.Grant.decode(reader, reader.uint32()));
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
        const obj = createBaseQueryGrantsResponse();
        if (Array.isArray(object?.grants))
            obj.grants = object.grants.map((e) => authz_1.Grant.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.grants) {
            obj.grants = message.grants.map((e) => (e ? authz_1.Grant.toJSON(e) : undefined));
        }
        else {
            obj.grants = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryGrantsResponse();
        message.grants = object.grants?.map((e) => authz_1.Grant.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryGranterGrantsRequest() {
    return {
        granter: "",
        pagination: undefined,
    };
}
exports.QueryGranterGrantsRequest = {
    typeUrl: "/cosmos.authz.v1beta1.QueryGranterGrantsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.granter !== "") {
            writer.uint32(10).string(message.granter);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryGranterGrantsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.granter = reader.string();
                    break;
                case 2:
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
        const obj = createBaseQueryGranterGrantsRequest();
        if ((0, helpers_1.isSet)(object.granter))
            obj.granter = String(object.granter);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.granter !== undefined && (obj.granter = message.granter);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryGranterGrantsRequest();
        message.granter = object.granter ?? "";
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryGranterGrantsResponse() {
    return {
        grants: [],
        pagination: undefined,
    };
}
exports.QueryGranterGrantsResponse = {
    typeUrl: "/cosmos.authz.v1beta1.QueryGranterGrantsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.grants) {
            authz_1.GrantAuthorization.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryGranterGrantsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.grants.push(authz_1.GrantAuthorization.decode(reader, reader.uint32()));
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
        const obj = createBaseQueryGranterGrantsResponse();
        if (Array.isArray(object?.grants))
            obj.grants = object.grants.map((e) => authz_1.GrantAuthorization.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.grants) {
            obj.grants = message.grants.map((e) => (e ? authz_1.GrantAuthorization.toJSON(e) : undefined));
        }
        else {
            obj.grants = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryGranterGrantsResponse();
        message.grants = object.grants?.map((e) => authz_1.GrantAuthorization.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryGranteeGrantsRequest() {
    return {
        grantee: "",
        pagination: undefined,
    };
}
exports.QueryGranteeGrantsRequest = {
    typeUrl: "/cosmos.authz.v1beta1.QueryGranteeGrantsRequest",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        if (message.grantee !== "") {
            writer.uint32(10).string(message.grantee);
        }
        if (message.pagination !== undefined) {
            pagination_1.PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryGranteeGrantsRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.grantee = reader.string();
                    break;
                case 2:
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
        const obj = createBaseQueryGranteeGrantsRequest();
        if ((0, helpers_1.isSet)(object.grantee))
            obj.grantee = String(object.grantee);
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageRequest.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        message.grantee !== undefined && (obj.grantee = message.grantee);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageRequest.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryGranteeGrantsRequest();
        message.grantee = object.grantee ?? "";
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageRequest.fromPartial(object.pagination);
        }
        return message;
    },
};
function createBaseQueryGranteeGrantsResponse() {
    return {
        grants: [],
        pagination: undefined,
    };
}
exports.QueryGranteeGrantsResponse = {
    typeUrl: "/cosmos.authz.v1beta1.QueryGranteeGrantsResponse",
    encode(message, writer = binary_1.BinaryWriter.create()) {
        for (const v of message.grants) {
            authz_1.GrantAuthorization.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            pagination_1.PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof binary_1.BinaryReader ? input : new binary_1.BinaryReader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseQueryGranteeGrantsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.grants.push(authz_1.GrantAuthorization.decode(reader, reader.uint32()));
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
        const obj = createBaseQueryGranteeGrantsResponse();
        if (Array.isArray(object?.grants))
            obj.grants = object.grants.map((e) => authz_1.GrantAuthorization.fromJSON(e));
        if ((0, helpers_1.isSet)(object.pagination))
            obj.pagination = pagination_1.PageResponse.fromJSON(object.pagination);
        return obj;
    },
    toJSON(message) {
        const obj = {};
        if (message.grants) {
            obj.grants = message.grants.map((e) => (e ? authz_1.GrantAuthorization.toJSON(e) : undefined));
        }
        else {
            obj.grants = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination ? pagination_1.PageResponse.toJSON(message.pagination) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = createBaseQueryGranteeGrantsResponse();
        message.grants = object.grants?.map((e) => authz_1.GrantAuthorization.fromPartial(e)) || [];
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = pagination_1.PageResponse.fromPartial(object.pagination);
        }
        return message;
    },
};
class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.Grants = this.Grants.bind(this);
        this.GranterGrants = this.GranterGrants.bind(this);
        this.GranteeGrants = this.GranteeGrants.bind(this);
    }
    Grants(request) {
        const data = exports.QueryGrantsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.authz.v1beta1.Query", "Grants", data);
        return promise.then((data) => exports.QueryGrantsResponse.decode(new binary_1.BinaryReader(data)));
    }
    GranterGrants(request) {
        const data = exports.QueryGranterGrantsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.authz.v1beta1.Query", "GranterGrants", data);
        return promise.then((data) => exports.QueryGranterGrantsResponse.decode(new binary_1.BinaryReader(data)));
    }
    GranteeGrants(request) {
        const data = exports.QueryGranteeGrantsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.authz.v1beta1.Query", "GranteeGrants", data);
        return promise.then((data) => exports.QueryGranteeGrantsResponse.decode(new binary_1.BinaryReader(data)));
    }
}
exports.QueryClientImpl = QueryClientImpl;
//# sourceMappingURL=query.js.map