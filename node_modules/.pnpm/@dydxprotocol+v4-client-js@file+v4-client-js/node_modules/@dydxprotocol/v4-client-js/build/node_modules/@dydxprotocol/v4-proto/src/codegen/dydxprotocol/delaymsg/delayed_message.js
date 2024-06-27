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
exports.DelayedMessage = void 0;
const any_1 = require("../../google/protobuf/any");
const _m0 = __importStar(require("protobufjs/minimal"));
function createBaseDelayedMessage() {
    return {
        id: 0,
        msg: undefined,
        blockHeight: 0
    };
}
exports.DelayedMessage = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.id !== 0) {
            writer.uint32(8).uint32(message.id);
        }
        if (message.msg !== undefined) {
            any_1.Any.encode(message.msg, writer.uint32(18).fork()).ldelim();
        }
        if (message.blockHeight !== 0) {
            writer.uint32(24).uint32(message.blockHeight);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDelayedMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.uint32();
                    break;
                case 2:
                    message.msg = any_1.Any.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.blockHeight = reader.uint32();
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
        const message = createBaseDelayedMessage();
        message.id = (_a = object.id) !== null && _a !== void 0 ? _a : 0;
        message.msg = object.msg !== undefined && object.msg !== null ? any_1.Any.fromPartial(object.msg) : undefined;
        message.blockHeight = (_b = object.blockHeight) !== null && _b !== void 0 ? _b : 0;
        return message;
    }
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZGVsYXllZF9tZXNzYWdlLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL2RlbGF5bXNnL2RlbGF5ZWRfbWVzc2FnZS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBLG1EQUE0RDtBQUM1RCx3REFBMEM7QUFzQjFDLFNBQVMsd0JBQXdCO0lBQy9CLE9BQU87UUFDTCxFQUFFLEVBQUUsQ0FBQztRQUNMLEdBQUcsRUFBRSxTQUFTO1FBQ2QsV0FBVyxFQUFFLENBQUM7S0FDZixDQUFDO0FBQ0osQ0FBQztBQUVZLFFBQUEsY0FBYyxHQUFHO0lBQzVCLE1BQU0sQ0FBQyxPQUF1QixFQUFFLFNBQXFCLEdBQUcsQ0FBQyxNQUFNLENBQUMsTUFBTSxFQUFFO1FBQ3RFLElBQUksT0FBTyxDQUFDLEVBQUUsS0FBSyxDQUFDLEVBQUU7WUFDcEIsTUFBTSxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUMsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLEVBQUUsQ0FBQyxDQUFDO1NBQ3JDO1FBRUQsSUFBSSxPQUFPLENBQUMsR0FBRyxLQUFLLFNBQVMsRUFBRTtZQUM3QixTQUFHLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxHQUFHLEVBQUUsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxJQUFJLEVBQUUsQ0FBQyxDQUFDLE1BQU0sRUFBRSxDQUFDO1NBQzVEO1FBRUQsSUFBSSxPQUFPLENBQUMsV0FBVyxLQUFLLENBQUMsRUFBRTtZQUM3QixNQUFNLENBQUMsTUFBTSxDQUFDLEVBQUUsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMsV0FBVyxDQUFDLENBQUM7U0FDL0M7UUFFRCxPQUFPLE1BQU0sQ0FBQztJQUNoQixDQUFDO0lBRUQsTUFBTSxDQUFDLEtBQThCLEVBQUUsTUFBZTtRQUNwRCxNQUFNLE1BQU0sR0FBRyxLQUFLLFlBQVksR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxJQUFJLEdBQUcsQ0FBQyxNQUFNLENBQUMsS0FBSyxDQUFDLENBQUM7UUFDM0UsSUFBSSxHQUFHLEdBQUcsTUFBTSxLQUFLLFNBQVMsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsR0FBRyxNQUFNLENBQUM7UUFDbEUsTUFBTSxPQUFPLEdBQUcsd0JBQXdCLEVBQUUsQ0FBQztRQUUzQyxPQUFPLE1BQU0sQ0FBQyxHQUFHLEdBQUcsR0FBRyxFQUFFO1lBQ3ZCLE1BQU0sR0FBRyxHQUFHLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQztZQUU1QixRQUFRLEdBQUcsS0FBSyxDQUFDLEVBQUU7Z0JBQ2pCLEtBQUssQ0FBQztvQkFDSixPQUFPLENBQUMsRUFBRSxHQUFHLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQztvQkFDN0IsTUFBTTtnQkFFUixLQUFLLENBQUM7b0JBQ0osT0FBTyxDQUFDLEdBQUcsR0FBRyxTQUFHLENBQUMsTUFBTSxDQUFDLE1BQU0sRUFBRSxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUMsQ0FBQztvQkFDbEQsTUFBTTtnQkFFUixLQUFLLENBQUM7b0JBQ0osT0FBTyxDQUFDLFdBQVcsR0FBRyxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUM7b0JBQ3RDLE1BQU07Z0JBRVI7b0JBQ0UsTUFBTSxDQUFDLFFBQVEsQ0FBQyxHQUFHLEdBQUcsQ0FBQyxDQUFDLENBQUM7b0JBQ3pCLE1BQU07YUFDVDtTQUNGO1FBRUQsT0FBTyxPQUFPLENBQUM7SUFDakIsQ0FBQztJQUVELFdBQVcsQ0FBQyxNQUFtQzs7UUFDN0MsTUFBTSxPQUFPLEdBQUcsd0JBQXdCLEVBQUUsQ0FBQztRQUMzQyxPQUFPLENBQUMsRUFBRSxHQUFHLE1BQUEsTUFBTSxDQUFDLEVBQUUsbUNBQUksQ0FBQyxDQUFDO1FBQzVCLE9BQU8sQ0FBQyxHQUFHLEdBQUcsTUFBTSxDQUFDLEdBQUcsS0FBSyxTQUFTLElBQUksTUFBTSxDQUFDLEdBQUcsS0FBSyxJQUFJLENBQUMsQ0FBQyxDQUFDLFNBQUcsQ0FBQyxXQUFXLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsQ0FBQyxTQUFTLENBQUM7UUFDeEcsT0FBTyxDQUFDLFdBQVcsR0FBRyxNQUFBLE1BQU0sQ0FBQyxXQUFXLG1DQUFJLENBQUMsQ0FBQztRQUM5QyxPQUFPLE9BQU8sQ0FBQztJQUNqQixDQUFDO0NBRUYsQ0FBQyJ9