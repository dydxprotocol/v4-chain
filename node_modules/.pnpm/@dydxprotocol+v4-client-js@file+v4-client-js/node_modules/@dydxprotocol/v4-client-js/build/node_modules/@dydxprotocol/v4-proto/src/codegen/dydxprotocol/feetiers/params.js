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
exports.PerpetualFeeTier = exports.PerpetualFeeParams = void 0;
const _m0 = __importStar(require("protobufjs/minimal"));
const helpers_1 = require("../../helpers");
function createBasePerpetualFeeParams() {
    return {
        tiers: []
    };
}
exports.PerpetualFeeParams = {
    encode(message, writer = _m0.Writer.create()) {
        for (const v of message.tiers) {
            exports.PerpetualFeeTier.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePerpetualFeeParams();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.tiers.push(exports.PerpetualFeeTier.decode(reader, reader.uint32()));
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
        const message = createBasePerpetualFeeParams();
        message.tiers = ((_a = object.tiers) === null || _a === void 0 ? void 0 : _a.map(e => exports.PerpetualFeeTier.fromPartial(e))) || [];
        return message;
    }
};
function createBasePerpetualFeeTier() {
    return {
        name: "",
        absoluteVolumeRequirement: helpers_1.Long.UZERO,
        totalVolumeShareRequirementPpm: 0,
        makerVolumeShareRequirementPpm: 0,
        makerFeePpm: 0,
        takerFeePpm: 0
    };
}
exports.PerpetualFeeTier = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        if (!message.absoluteVolumeRequirement.isZero()) {
            writer.uint32(16).uint64(message.absoluteVolumeRequirement);
        }
        if (message.totalVolumeShareRequirementPpm !== 0) {
            writer.uint32(24).uint32(message.totalVolumeShareRequirementPpm);
        }
        if (message.makerVolumeShareRequirementPpm !== 0) {
            writer.uint32(32).uint32(message.makerVolumeShareRequirementPpm);
        }
        if (message.makerFeePpm !== 0) {
            writer.uint32(40).sint32(message.makerFeePpm);
        }
        if (message.takerFeePpm !== 0) {
            writer.uint32(48).sint32(message.takerFeePpm);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBasePerpetualFeeTier();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.absoluteVolumeRequirement = reader.uint64();
                    break;
                case 3:
                    message.totalVolumeShareRequirementPpm = reader.uint32();
                    break;
                case 4:
                    message.makerVolumeShareRequirementPpm = reader.uint32();
                    break;
                case 5:
                    message.makerFeePpm = reader.sint32();
                    break;
                case 6:
                    message.takerFeePpm = reader.sint32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        var _a, _b, _c, _d, _e;
        const message = createBasePerpetualFeeTier();
        message.name = (_a = object.name) !== null && _a !== void 0 ? _a : "";
        message.absoluteVolumeRequirement = object.absoluteVolumeRequirement !== undefined && object.absoluteVolumeRequirement !== null ? helpers_1.Long.fromValue(object.absoluteVolumeRequirement) : helpers_1.Long.UZERO;
        message.totalVolumeShareRequirementPpm = (_b = object.totalVolumeShareRequirementPpm) !== null && _b !== void 0 ? _b : 0;
        message.makerVolumeShareRequirementPpm = (_c = object.makerVolumeShareRequirementPpm) !== null && _c !== void 0 ? _c : 0;
        message.makerFeePpm = (_d = object.makerFeePpm) !== null && _d !== void 0 ? _d : 0;
        message.takerFeePpm = (_e = object.takerFeePpm) !== null && _e !== void 0 ? _e : 0;
        return message;
    }
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicGFyYW1zLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL2ZlZXRpZXJzL3BhcmFtcy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBLHdEQUEwQztBQUMxQywyQ0FBa0Q7QUE0Q2xELFNBQVMsNEJBQTRCO0lBQ25DLE9BQU87UUFDTCxLQUFLLEVBQUUsRUFBRTtLQUNWLENBQUM7QUFDSixDQUFDO0FBRVksUUFBQSxrQkFBa0IsR0FBRztJQUNoQyxNQUFNLENBQUMsT0FBMkIsRUFBRSxTQUFxQixHQUFHLENBQUMsTUFBTSxDQUFDLE1BQU0sRUFBRTtRQUMxRSxLQUFLLE1BQU0sQ0FBQyxJQUFJLE9BQU8sQ0FBQyxLQUFLLEVBQUU7WUFDN0Isd0JBQWdCLENBQUMsTUFBTSxDQUFDLENBQUUsRUFBRSxNQUFNLENBQUMsTUFBTSxDQUFDLEVBQUUsQ0FBQyxDQUFDLElBQUksRUFBRSxDQUFDLENBQUMsTUFBTSxFQUFFLENBQUM7U0FDaEU7UUFFRCxPQUFPLE1BQU0sQ0FBQztJQUNoQixDQUFDO0lBRUQsTUFBTSxDQUFDLEtBQThCLEVBQUUsTUFBZTtRQUNwRCxNQUFNLE1BQU0sR0FBRyxLQUFLLFlBQVksR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxJQUFJLEdBQUcsQ0FBQyxNQUFNLENBQUMsS0FBSyxDQUFDLENBQUM7UUFDM0UsSUFBSSxHQUFHLEdBQUcsTUFBTSxLQUFLLFNBQVMsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsR0FBRyxNQUFNLENBQUM7UUFDbEUsTUFBTSxPQUFPLEdBQUcsNEJBQTRCLEVBQUUsQ0FBQztRQUUvQyxPQUFPLE1B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