"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.GrantQueueItem = exports.GrantAuthorization = exports.Grant = exports.GenericAuthorization = void 0;
const any_1 = require("../../../google/protobuf/any");
const timestamp_1 = require("../../../google/protobuf/timestamp");
const _m0 = __importStar(require("protobufjs/minimal"));
const helpers_1 = require("../../../helpers");
function createBaseGenericAuthorization() {
    return {
        msg: ""
    };
}
exports.GenericAuthorization = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.msg !== "") {
            writer.uint32(10).string(message.msg);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGenericAuthorization();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.msg = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        var _a;
        const message = createBaseGenericAuthorization();
        message.msg = (_a = object.msg) !== null && _a !== void 0 ? _a : "";
        return message;
    }
};
function createBaseGrant() {
    return {
        authorization: undefined,
        expiration: undefined
    };
}
exports.Grant = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.authorization !== undefined) {
            any_1.Any.encode(message.authorization, writer.uint32(10).fork()).ldelim();
        }
        if (message.expiration !== undefined) {
            timestamp_1.Timestamp.encode((0, helpers_1.toTimestamp)(message.expiration), writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGrant();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.authorization = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.expiration = (0, helpers_1.fromTimestamp)(timestamp_1.Timestamp.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        var _a;
        const message = createBaseGrant();
        message.authorization = object.authorization !== undefined && object.authorization !== null ? any_1.Any.fromPartial(object.authorization) : undefined;
        message.expiration = (_a = object.expiration) !== null && _a !== void 0 ? _a : undefined;
        return message;
    }
};
function createBaseGrantAuthorization() {
    return {
        granter: "",
        grantee: "",
        authorization: undefined,
        expiration: undefined
    };
}
exports.GrantAuthorization = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.granter !== "") {
            writer.uint32(10).string(message.granter);
        }
        if (message.grantee !== "") {
            writer.uint32(18).string(message.grantee);
        }
        if (message.authorization !== undefined) {
            any_1.Any.encode(message.authorization, writer.uint32(26).fork()).ldelim();
        }
        if (message.expiration !== undefined) {
            timestamp_1.Timestamp.encode((0, helpers_1.toTimestamp)(message.expiration), writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGrantAuthorization();
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
                    message.authorization = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.expiration = (0, helpers_1.fromTimestamp)(timestamp_1.Timestamp.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseGrantAuthorization();
        message.granter = (_a = object.granter) !== null && _a !== void 0 ? _a : "";
        message.grantee = (_b = object.grantee) !== null && _b !== void 0 ? _b : "";
        message.authorization = object.authorization !== undefined && object.authorization !== null ? any_1.Any.fromPartial(object.authorization) : undefined;
        message.expiration = (_c = object.expiration) !== null && _c !== void 0 ? _c : undefined;
        return message;
    }
};
function createBaseGrantQueueItem() {
    return {
        msgTypeUrls: []
    };
}
exports.GrantQueueItem = {
    encode(message, writer = _m0.Writer.create()) {
        for (const v of message.msgTypeUrls) {
            writer.uint32(10).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGrantQueueItem();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.msgTypeUrls.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        var _a;
        const message = createBaseGrantQueueItem();
        message.msgTypeUrls = ((_a = object.msgTypeUrls) === null || _a === void 0 ? void 0 : _a.map(e => e)) || [];
        return message;
    }
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXV0aHouanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi9ub2RlX21vZHVsZXMvQGR5ZHhwcm90b2NvbC92NC1wcm90by9zcmMvY29kZWdlbi9jb3Ntb3MvYXV0aHovdjFiZXRhMS9hdXRoei50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBLHNEQUErRDtBQUMvRCxrRUFBK0Q7QUFDL0Qsd0RBQTBDO0FBQzFDLDhDQUEyRTtBQTRFM0UsU0FBUyw4QkFBOEI7SUFDckMsT0FBTztRQUNMLEdBQUcsRUFBRSxFQUFFO0tBQ1IsQ0FBQztBQUNKLENBQUM7QUFFWSxRQUFBLG9CQUFvQixHQUFHO0lBQ2xDLE1BQU0sQ0FBQyxPQUE2QixFQUFFLFNBQXFCLEdBQUcsQ0FBQyxNQUFNLENBQUMsTUFBTSxFQUFFO1FBQzVFLElBQUksT0FBTyxDQUFDLEdBQUcsS0FBSyxFQUFFLEVBQUU7WUFDdEIsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLEdBQUcsQ0FBQyxDQUFDO1NBQ3ZDO1FBRUQsT0FBTyxNQUFNLENBQUM7SUFDaEIsQ0FBQztJQUVELE1BQU0sQ0FBQyxLQUE4QixFQUFFLE1BQWU7UUFDcEQsTUFBTSxNQUFNLEdBQUcsS0FBSyxZQUFZLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQyxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsSUFBSSxHQUFHLENBQUMsTUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDO1FBQzNFLElBQUksR0FBRyxHQUFHLE1BQU0sS0FBSyxTQUFTLENBQUMsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxHQUFHLENBQUMsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxHQUFHLEdBQUcsTUFBTSxDQUFDO1FBQ2xFLE1BQU0sT0FBTyxHQUFHLDhCQUE4QixFQUFFLENBQUM7UUFFakQsT0FBTyxNQUFNLENBQUMsR0FB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