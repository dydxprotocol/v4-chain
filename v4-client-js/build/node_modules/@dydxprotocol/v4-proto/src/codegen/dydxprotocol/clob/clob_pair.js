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
exports.ClobPair = exports.SpotClobMetadata = exports.PerpetualClobMetadata = exports.clobPair_StatusToJSON = exports.clobPair_StatusFromJSON = exports.ClobPair_StatusSDKType = exports.ClobPair_Status = void 0;
const _m0 = __importStar(require("protobufjs/minimal"));
const helpers_1 = require("../../helpers");
/** Status of the CLOB. */
var ClobPair_Status;
(function (ClobPair_Status) {
    /** STATUS_UNSPECIFIED - Default value. This value is invalid and unused. */
    ClobPair_Status[ClobPair_Status["STATUS_UNSPECIFIED"] = 0] = "STATUS_UNSPECIFIED";
    /** STATUS_ACTIVE - STATUS_ACTIVE represents an active clob pair. */
    ClobPair_Status[ClobPair_Status["STATUS_ACTIVE"] = 1] = "STATUS_ACTIVE";
    /**
     * STATUS_PAUSED - STATUS_PAUSED behavior is unfinalized.
     * TODO(DEC-600): update this documentation.
     */
    ClobPair_Status[ClobPair_Status["STATUS_PAUSED"] = 2] = "STATUS_PAUSED";
    /**
     * STATUS_CANCEL_ONLY - STATUS_CANCEL_ONLY behavior is unfinalized.
     * TODO(DEC-600): update this documentation.
     */
    ClobPair_Status[ClobPair_Status["STATUS_CANCEL_ONLY"] = 3] = "STATUS_CANCEL_ONLY";
    /**
     * STATUS_POST_ONLY - STATUS_POST_ONLY behavior is unfinalized.
     * TODO(DEC-600): update this documentation.
     */
    ClobPair_Status[ClobPair_Status["STATUS_POST_ONLY"] = 4] = "STATUS_POST_ONLY";
    /**
     * STATUS_INITIALIZING - STATUS_INITIALIZING represents a newly-added clob pair.
     * Clob pairs in this state only accept orders which are
     * both short-term and post-only.
     */
    ClobPair_Status[ClobPair_Status["STATUS_INITIALIZING"] = 5] = "STATUS_INITIALIZING";
    /**
     * STATUS_FINAL_SETTLEMENT - STATUS_FINAL_SETTLEMENT represents a clob pair which is deactivated
     * and trading has ceased. All open positions will be closed by the
     * protocol. Open stateful orders will be cancelled. Open short-term
     * orders will be left to expire.
     */
    ClobPair_Status[ClobPair_Status["STATUS_FINAL_SETTLEMENT"] = 6] = "STATUS_FINAL_SETTLEMENT";
    ClobPair_Status[ClobPair_Status["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(ClobPair_Status = exports.ClobPair_Status || (exports.ClobPair_Status = {}));
exports.ClobPair_StatusSDKType = ClobPair_Status;
function clobPair_StatusFromJSON(object) {
    switch (object) {
        case 0:
        case "STATUS_UNSPECIFIED":
            return ClobPair_Status.STATUS_UNSPECIFIED;
        case 1:
        case "STATUS_ACTIVE":
            return ClobPair_Status.STATUS_ACTIVE;
        case 2:
        case "STATUS_PAUSED":
            return ClobPair_Status.STATUS_PAUSED;
        case 3:
        case "STATUS_CANCEL_ONLY":
            return ClobPair_Status.STATUS_CANCEL_ONLY;
        case 4:
        case "STATUS_POST_ONLY":
            return ClobPair_Status.STATUS_POST_ONLY;
        case 5:
        case "STATUS_INITIALIZING":
            return ClobPair_Status.STATUS_INITIALIZING;
        case 6:
        case "STATUS_FINAL_SETTLEMENT":
            return ClobPair_Status.STATUS_FINAL_SETTLEMENT;
        case -1:
        case "UNRECOGNIZED":
        default:
            return ClobPair_Status.UNRECOGNIZED;
    }
}
exports.clobPair_StatusFromJSON = clobPair_StatusFromJSON;
function clobPair_StatusToJSON(object) {
    switch (object) {
        case ClobPair_Status.STATUS_UNSPECIFIED:
            return "STATUS_UNSPECIFIED";
        case ClobPair_Status.STATUS_ACTIVE:
            return "STATUS_ACTIVE";
        case ClobPair_Status.STATUS_PAUSED:
            return "STATUS_PAUSED";
        case ClobPair_Status.STATUS_CANCEL_ONLY:
            return "STATUS_CANCEL_ONLY";
        case ClobPair_Status.STATUS_POST_ONLY:
            return "STATUS_POST_ONLY";
        case ClobPair_Status.STATUS_INITIALIZING:
            return "STATUS_INITIALIZING";
        case ClobPair_Status.STATUS_FINAL_SETTLEMENT:
            return "STATUS_FINAL_SETTLEMENT";
        case ClobPair_Status.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.clobPair_StatusToJSON = clobPair_StatusToJSON;
function createBasePerpetualClobMetadata() {
    return {
        perpetualId: 0
    };
}
exports.PerpetualClobMetadata = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.perpetualId !== 0) {
            writer.uint32(8).uint32(message.perpetualId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePerpetualClobMetadata();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.perpetualId = reader.uint32();
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
        const message = createBasePerpetualClobMetadata();
        message.perpetualId = (_a = object.perpetualId) !== null && _a !== void 0 ? _a : 0;
        return message;
    }
};
function createBaseSpotClobMetadata() {
    return {
        baseAssetId: 0,
        quoteAssetId: 0
    };
}
exports.SpotClobMetadata = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.baseAssetId !== 0) {
            writer.uint32(8).uint32(message.baseAssetId);
        }
        if (message.quoteAssetId !== 0) {
            writer.uint32(16).uint32(message.quoteAssetId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSpotClobMetadata();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseAssetId = reader.uint32();
                    break;
                case 2:
                    message.quoteAssetId = reader.uint32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseSpotClobMetadata();
        message.baseAssetId = (_a = object.baseAssetId) !== null && _a !== void 0 ? _a : 0;
        message.quoteAssetId = (_b = object.quoteAssetId) !== null && _b !== void 0 ? _b : 0;
        return message;
    }
};
function createBaseClobPair() {
    return {
        id: 0,
        perpetualClobMetadata: undefined,
        spotClobMetadata: undefined,
        stepBaseQuantums: helpers_1.Long.UZERO,
        subticksPerTick: 0,
        quantumConversionExponent: 0,
        status: 0
    };
}
exports.ClobPair = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.id !== 0) {
            writer.uint32(8).uint32(message.id);
        }
        if (message.perpetualClobMetadata !== undefined) {
            exports.PerpetualClobMetadata.encode(message.perpetualClobMetadata, writer.uint32(18).fork()).ldelim();
        }
        if (message.spotClobMetadata !== undefined) {
            exports.SpotClobMetadata.encode(message.spotClobMetadata, writer.uint32(26).fork()).ldelim();
        }
        if (!message.stepBaseQuantums.isZero()) {
            writer.uint32(32).uint64(message.stepBaseQuantums);
        }
        if (message.subticksPerTick !== 0) {
            writer.uint32(40).uint32(message.subticksPerTick);
        }
        if (message.quantumConversionExponent !== 0) {
            writer.uint32(48).sint32(message.quantumConversionExponent);
        }
        if (message.status !== 0) {
            writer.uint32(56).int32(message.status);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseClobPair();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.uint32();
                    break;
                case 2:
                    message.perpetualClobMetadata = exports.PerpetualClobMetadata.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.spotClobMetadata = exports.SpotClobMetadata.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.stepBaseQuantums = reader.uint64();
                    break;
                case 5:
                    message.subticksPerTick = reader.uint32();
                    break;
                case 6:
                    message.quantumConversionExponent = reader.sint32();
                    break;
                case 7:
                    message.status = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        var _a, _b, _c, _d;
        const message = createBaseClobPair();
        message.id = (_a = object.id) !== null && _a !== void 0 ? _a : 0;
        message.perpetualClobMetadata = object.perpetualClobMetadata !== undefined && object.perpetualClobMetadata !== null ? exports.PerpetualClobMetadata.fromPartial(object.perpetualClobMetadata) : undefined;
        message.spotClobMetadata = object.spotClobMetadata !== undefined && object.spotClobMetadata !== null ? exports.SpotClobMetadata.fromPartial(object.spotClobMetadata) : undefined;
        message.stepBaseQuantums = object.stepBaseQuantums !== undefined && object.stepBaseQuantums !== null ? helpers_1.Long.fromValue(object.stepBaseQuantums) : helpers_1.Long.UZERO;
        message.subticksPerTick = (_b = object.subticksPerTick) !== null && _b !== void 0 ? _b : 0;
        message.quantumConversionExponent = (_c = object.quantumConversionExponent) !== null && _c !== void 0 ? _c : 0;
        message.status = (_d = object.status) !== null && _d !== void 0 ? _d : 0;
        return message;
    }
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY2xvYl9wYWlyLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL2Nsb2IvY2xvYl9wYWlyLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQUEsd0RBQTBDO0FBQzFDLDJDQUFrRDtBQUNsRCwwQkFBMEI7QUFFMUIsSUFBWSxlQXdDWDtBQXhDRCxXQUFZLGVBQWU7SUFDekIsNEVBQTRFO0lBQzVFLGlGQUFzQixDQUFBO0lBRXRCLG9FQUFvRTtJQUNwRSx1RUFBaUIsQ0FBQTtJQUVqQjs7O09BR0c7SUFDSCx1RUFBaUIsQ0FBQTtJQUVqQjs7O09BR0c7SUFDSCxpRkFBc0IsQ0FBQTtJQUV0Qjs7O09BR0c7SUFDSCw2RUFBb0IsQ0FBQTtJQUVwQjs7OztPQUlHO0lBQ0gsbUZBQXVCLENBQUE7SUFFdkI7Ozs7O09BS0c7SUFDSCwyRkFBMkIsQ0FBQTtJQUMzQixzRUFBaUIsQ0FBQTtBQUNuQixDQUFDLEVBeENXLGVBQWUsR0FBZix1QkFBZSxLQUFmLHVCQUFlLFFBd0MxQjtBQUNZLFFBQUEsc0JBQXNCLEdBQUcsZUFBZSxDQUFDO0FBQ3RELFNBQWdCLHVCQUF1QixDQUFDLE1BQVc7SUFDakQsUUFBUSxNQUFNLEVBQUU7UUFDZCxLQUFLLENBQUMsQ0FBQztRQUNQLEtBQUssb0JBQW9CO1lBQ3ZCLE9BQU8sZUFBZSxDQUFDLGtCQUFrQixDQUFDO1FBRTVDLEtBQUssQ0FBQyxDQUFDO1FBQ1AsS0FBSyxlQUFlO1lBQ2xCLE9BQU8sZUFBZSxDQUFDLGFBQWEsQ0FBQztRQUV2QyxLQUFLLENBQUMsQ0FBQztRQUNQLEtBQUssZUFBZTtZQUNsQixPQUFPLGVBQWUsQ0FBQyxhQUFhLENBQUM7UUFFdkMsS0FBSyxDQUFDLENBQUM7UUFDUCxLQUFLLG9CQUFvQjtZQUN2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