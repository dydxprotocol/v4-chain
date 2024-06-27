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
exports.Params = void 0;
const _m0 = __importStar(require("protobufjs/minimal"));
function createBaseParams() {
    return {
        fundingRateClampFactorPpm: 0,
        premiumVoteClampFactorPpm: 0,
        minNumVotesPerSample: 0
    };
}
exports.Params = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.fundingRateClampFactorPpm !== 0) {
            writer.uint32(8).uint32(message.fundingRateClampFactorPpm);
        }
        if (message.premiumVoteClampFactorPpm !== 0) {
            writer.uint32(16).uint32(message.premiumVoteClampFactorPpm);
        }
        if (message.minNumVotesPerSample !== 0) {
            writer.uint32(24).uint32(message.minNumVotesPerSample);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseParams();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.fundingRateClampFactorPpm = reader.uint32();
                    break;
                case 2:
                    message.premiumVoteClampFactorPpm = reader.uint32();
                    break;
                case 3:
                    message.minNumVotesPerSample = reader.uint32();
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
        const message = createBaseParams();
        message.fundingRateClampFactorPpm = (_a = object.fundingRateClampFactorPpm) !== null && _a !== void 0 ? _a : 0;
        message.premiumVoteClampFactorPpm = (_b = object.premiumVoteClampFactorPpm) !== null && _b !== void 0 ? _b : 0;
        message.minNumVotesPerSample = (_c = object.minNumVotesPerSample) !== null && _c !== void 0 ? _c : 0;
        return message;
    }
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicGFyYW1zLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL3BlcnBldHVhbHMvcGFyYW1zLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQUEsd0RBQTBDO0FBaUMxQyxTQUFTLGdCQUFnQjtJQUN2QixPQUFPO1FBQ0wseUJBQXlCLEVBQUUsQ0FBQztRQUM1Qix5QkFBeUIsRUFBRSxDQUFDO1FBQzVCLG9CQUFvQixFQUFFLENBQUM7S0FDeEIsQ0FBQztBQUNKLENBQUM7QUFFWSxRQUFBLE1BQU0sR0FBRztJQUNwQixNQUFNLENBQUMsT0FBZSxFQUFFLFNBQXFCLEdBQUcsQ0FBQyxNQUFNLENBQUMsTUFBTSxFQUFFO1FBQzlELElBQUksT0FBTyxDQUFDLHlCQUF5QixLQUFLLENBQUMsRUFBRTtZQUMzQyxNQUFNLENBQUMsTUFBTSxDQUFDLENBQUMsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMseUJBQXlCLENBQUMsQ0FBQztTQUM1RDtRQUVELElBQUksT0FBTyxDQUFDLHlCQUF5QixLQUFLLENBQUMsRUFBRTtZQUMzQyxNQUFNLENBQUMsTUFBTSxDQUFDLEVBQUUsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMseUJBQXlCLENBQUMsQ0FBQztTQUM3RDtRQUVELElBQUksT0FBTyxDQUFDLG9CQUFvQixLQUFLLENBQUMsRUFBRTtZQUN0QyxNQUFNLENBQUMsTUFBTSxDQUFDLEVBQUUsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMsb0JBQW9CLENBQUMsQ0FBQztTQUN4RDtRQUVELE9BQU8sTUFBTSxDQUFDO0lBQ2hCLENBQUM7SUFFRCxNQUFNLENBQUMsS0FBOEIsRUFBRSxNQUFlO1FBQ3BELE1BQU0sTUFBTSxHQUFHLEtBQUssWUFBWSxHQUFHLENBQUMsTUFBTSxDQUFDLENBQUMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLElBQUksR0FBRyxDQUFDLE1BQU0sQ0FBQyxLQUFLLENBQUMsQ0FBQztRQUMzRSxJQUFJLEdBQUcsR0FBRyxNQUFNLEtBQUssU0FBUyxDQUFDLENBQUMsQ0FBQyxNQUFNLENBQUMsR0FBRyxDQUFDLENBQUMsQ0FBQyxNQUFNLENBQUMsR0FBRyxHQUFHLE1BQU0sQ0FBQztRQUNsRSxNQUFNLE9BQU8sR0FBRyxnQkFBZ0IsRUFBRSxDQUFDO1FBRW5DLE9BQU8sTUFBTSxDQUFDLEdBQUcsR0FBRyxHQUFHLEVBQUU7WUFDdkIsTUFBTSxHQUFHLEdBQUcsTUFBTSxDQUFDLE1BQU0sRUFBRSxDQUFDO1lBRTVCLFFBQVEsR0FBRyxLQUFLLENBQUMsRUFBRTtnQkFDakIsS0FBSyxDQUFDO29CQUNKLE9BQU8sQ0FBQyx5QkFBeUIsR0FBRyxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUM7b0JBQ3BELE1BQU07Z0JBRVIsS0FBSyxDQUFDO29CQUNKLE9BQU8sQ0FBQyx5QkFBeUIsR0FBRyxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUM7b0JBQ3BELE1BQU07Z0JBRVIsS0FBSyxDQUFDO29CQUNKLE9BQU8sQ0FBQyxvQkFBb0IsR0FBRyxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUM7b0JBQy9DLE1BQU07Z0JBRVI7b0JBQ0UsTUFBTSxDQUFDLFFBQVEsQ0FBQyxHQUFHLEdBQUcsQ0FBQyxDQUFDLENBQUM7b0JBQ3pCLE1BQU07YUFDVDtTQUNGO1FBRUQsT0FBTyxPQUFPLENBQUM7SUFDakIsQ0FBQztJQUVELFdBQVcsQ0FBQyxNQUEyQjs7UUFDckMsTUFBTSxPQUFPLEdBQUcsZ0JBQWdCLEVBQUUsQ0FBQztRQUNuQyxPQUFPLENBQUMseUJBQXlCLEdBQUcsTUFBQSxNQUFNLENBQUMseUJBQXlCLG1DQUFJLENBQUMsQ0FBQztRQUMxRSxPQUFPLENBQUMseUJBQXlCLEdBQUcsTUFBQSxNQUFNLENBQUMseUJBQXlCLG1DQUFJLENBQUMsQ0FBQztRQUMxRSxPQUFPLENBQUMsb0JBQW9CLEdBQUcsTUFBQSxNQUFNLENBQUMsb0JBQW9CLG1DQUFJLENBQUMsQ0FBQztRQUNoRSxPQUFPLE9BQU8sQ0FBQztJQUNqQixDQUFDO0NBRUYsQ0FBQyJ9