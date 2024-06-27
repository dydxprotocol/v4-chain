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
exports.EquityTierLimit = exports.EquityTierLimitConfiguration = void 0;
const _m0 = __importStar(require("protobufjs/minimal"));
function createBaseEquityTierLimitConfiguration() {
    return {
        shortTermOrderEquityTiers: [],
        statefulOrderEquityTiers: []
    };
}
exports.EquityTierLimitConfiguration = {
    encode(message, writer = _m0.Writer.create()) {
        for (const v of message.shortTermOrderEquityTiers) {
            exports.EquityTierLimit.encode(v, writer.uint32(10).fork()).ldelim();
        }
        for (const v of message.statefulOrderEquityTiers) {
            exports.EquityTierLimit.encode(v, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseEquityTierLimitConfiguration();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.shortTermOrderEquityTiers.push(exports.EquityTierLimit.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.statefulOrderEquityTiers.push(exports.EquityTierLimit.decode(reader, reader.uint32()));
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
        const message = createBaseEquityTierLimitConfiguration();
        message.shortTermOrderEquityTiers = ((_a = object.shortTermOrderEquityTiers) === null || _a === void 0 ? void 0 : _a.map(e => exports.EquityTierLimit.fromPartial(e))) || [];
        message.statefulOrderEquityTiers = ((_b = object.statefulOrderEquityTiers) === null || _b === void 0 ? void 0 : _b.map(e => exports.EquityTierLimit.fromPartial(e))) || [];
        return message;
    }
};
function createBaseEquityTierLimit() {
    return {
        usdTncRequired: new Uint8Array(),
        limit: 0
    };
}
exports.EquityTierLimit = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.usdTncRequired.length !== 0) {
            writer.uint32(10).bytes(message.usdTncRequired);
        }
        if (message.limit !== 0) {
            writer.uint32(16).uint32(message.limit);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseEquityTierLimit();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.usdTncRequired = reader.bytes();
                    break;
                case 2:
                    message.limit = reader.uint32();
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
        const message = createBaseEquityTierLimit();
        message.usdTncRequired = (_a = object.usdTncRequired) !== null && _a !== void 0 ? _a : new Uint8Array();
        message.limit = (_b = object.limit) !== null && _b !== void 0 ? _b : 0;
        return message;
    }
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXF1aXR5X3RpZXJfbGltaXRfY29uZmlnLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL2Nsb2IvZXF1aXR5X3RpZXJfbGltaXRfY29uZmlnLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQUEsd0RBQTBDO0FBNkMxQyxTQUFTLHNDQUFzQztJQUM3QyxPQUFPO1FBQ0wseUJBQXlCLEVBQUUsRUFBRTtRQUM3Qix3QkFBd0IsRUFBRSxFQUFFO0tBQzdCLENBQUM7QUFDSixDQUFDO0FBRVksUUFBQSw0QkFBNEIsR0FBRztJQUMxQyxNQUFNLENBQUMsT0FBcUMsRUFBRSxTQUFxQixHQUFHLENBQUMsTUFBTSxDQUFDLE1BQU0sRUFBRTtRQUNwRixLQUFLLE1BQU0sQ0FBQyxJQUFJLE9BQU8sQ0FBQyx5QkFBeUIsRUFBRTtZQUNqRCx1QkFBZSxDQUFDLE1BQU0sQ0FBQyxDQUFFLEVBQUUsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxJQUFJLEVBQUUsQ0FBQyxDQUFDLE1BQU0sRUFBRSxDQUFDO1NBQy9EO1FBRUQsS0FBSyxNQUFNLENBQUMsSUFBSSxPQUFPLENBQUMsd0JBQXdCLEVBQUU7WUFDaEQsdUJBQWUsQ0FBQyxNQUFNLENBQUMsQ0FBRSxFQUFFLE1BQU0sQ0FBQyxNQUFNLENBQUMsRUFBRSxDQUFDLENBQUMsSUFBSSxFQUFFLENBQUMsQ0FBQyxNQUFNLEVBQUUsQ0FBQztTQUMvRDtRQUVELE9BQU8sTUFBTSxDQUFDO0lBQ2hCLENBQUM7SUFFRCxNQUFNLENBQUMsS0FBOEIsRUFBRSxNQUFlO1FBQ3BELE1BQU0sTUFBTSxHQUFHLEtBQUssWUFBWSxHQUFHLENBQUMsTUFBTSxDQUFDLENBQUMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxD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