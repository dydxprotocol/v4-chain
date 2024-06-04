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
exports.RedisOrder = exports.redisOrder_TickerTypeToJSON = exports.redisOrder_TickerTypeFromJSON = exports.RedisOrder_TickerTypeSDKType = exports.RedisOrder_TickerType = void 0;
const clob_1 = require("../protocol/v1/clob");
const _m0 = __importStar(require("protobufjs/minimal"));
/** Enum for the ticker type, PERPETUAL or SPOT. */
var RedisOrder_TickerType;
(function (RedisOrder_TickerType) {
    /**
     * TICKER_TYPE_UNSPECIFIED - Default value for the enum. Should never be used in an initialized
     * `RedisOrder`.
     */
    RedisOrder_TickerType[RedisOrder_TickerType["TICKER_TYPE_UNSPECIFIED"] = 0] = "TICKER_TYPE_UNSPECIFIED";
    /** TICKER_TYPE_PERPETUAL - Ticker is for a perpetual pair. */
    RedisOrder_TickerType[RedisOrder_TickerType["TICKER_TYPE_PERPETUAL"] = 1] = "TICKER_TYPE_PERPETUAL";
    /** TICKER_TYPE_SPOT - Ticker is for a spot pair. */
    RedisOrder_TickerType[RedisOrder_TickerType["TICKER_TYPE_SPOT"] = 2] = "TICKER_TYPE_SPOT";
    RedisOrder_TickerType[RedisOrder_TickerType["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(RedisOrder_TickerType = exports.RedisOrder_TickerType || (exports.RedisOrder_TickerType = {}));
exports.RedisOrder_TickerTypeSDKType = RedisOrder_TickerType;
function redisOrder_TickerTypeFromJSON(object) {
    switch (object) {
        case 0:
        case "TICKER_TYPE_UNSPECIFIED":
            return RedisOrder_TickerType.TICKER_TYPE_UNSPECIFIED;
        case 1:
        case "TICKER_TYPE_PERPETUAL":
            return RedisOrder_TickerType.TICKER_TYPE_PERPETUAL;
        case 2:
        case "TICKER_TYPE_SPOT":
            return RedisOrder_TickerType.TICKER_TYPE_SPOT;
        case -1:
        case "UNRECOGNIZED":
        default:
            return RedisOrder_TickerType.UNRECOGNIZED;
    }
}
exports.redisOrder_TickerTypeFromJSON = redisOrder_TickerTypeFromJSON;
function redisOrder_TickerTypeToJSON(object) {
    switch (object) {
        case RedisOrder_TickerType.TICKER_TYPE_UNSPECIFIED:
            return "TICKER_TYPE_UNSPECIFIED";
        case RedisOrder_TickerType.TICKER_TYPE_PERPETUAL:
            return "TICKER_TYPE_PERPETUAL";
        case RedisOrder_TickerType.TICKER_TYPE_SPOT:
            return "TICKER_TYPE_SPOT";
        case RedisOrder_TickerType.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.redisOrder_TickerTypeToJSON = redisOrder_TickerTypeToJSON;
function createBaseRedisOrder() {
    return {
        id: "",
        order: undefined,
        ticker: "",
        tickerType: 0,
        price: "",
        size: ""
    };
}
exports.RedisOrder = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.id !== "") {
            writer.uint32(10).string(message.id);
        }
        if (message.order !== undefined) {
            clob_1.IndexerOrder.encode(message.order, writer.uint32(18).fork()).ldelim();
        }
        if (message.ticker !== "") {
            writer.uint32(26).string(message.ticker);
        }
        if (message.tickerType !== 0) {
            writer.uint32(32).int32(message.tickerType);
        }
        if (message.price !== "") {
            writer.uint32(42).string(message.price);
        }
        if (message.size !== "") {
            writer.uint32(50).string(message.size);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRedisOrder();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.id = reader.string();
                    break;
                case 2:
                    message.order = clob_1.IndexerOrder.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.ticker = reader.string();
                    break;
                case 4:
                    message.tickerType = reader.int32();
                    break;
                case 5:
                    message.price = reader.string();
                    break;
                case 6:
                    message.size = reader.string();
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
        const message = createBaseRedisOrder();
        message.id = (_a = object.id) !== null && _a !== void 0 ? _a : "";
        message.order = object.order !== undefined && object.order !== null ? clob_1.IndexerOrder.fromPartial(object.order) : undefined;
        message.ticker = (_b = object.ticker) !== null && _b !== void 0 ? _b : "";
        message.tickerType = (_c = object.tickerType) !== null && _c !== void 0 ? _c : 0;
        message.price = (_d = object.price) !== null && _d !== void 0 ? _d : "";
        message.size = (_e = object.size) !== null && _e !== void 0 ? _e : "";
        return message;
    }
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicmVkaXNfb3JkZXIuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi9ub2RlX21vZHVsZXMvQGR5ZHhwcm90b2NvbC92NC1wcm90by9zcmMvY29kZWdlbi9keWR4cHJvdG9jb2wvaW5kZXhlci9yZWRpcy9yZWRpc19vcmRlci50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBLDhDQUF3RTtBQUN4RSx3REFBMEM7QUFFMUMsbURBQW1EO0FBRW5ELElBQVkscUJBYVg7QUFiRCxXQUFZLHFCQUFxQjtJQUMvQjs7O09BR0c7SUFDSCx1R0FBMkIsQ0FBQTtJQUUzQiw4REFBOEQ7SUFDOUQsbUdBQXlCLENBQUE7SUFFekIsb0RBQW9EO0lBQ3BELHlGQUFvQixDQUFBO0lBQ3BCLGtGQUFpQixDQUFBO0FBQ25CLENBQUMsRUFiVyxxQkFBcUIsR0FBckIsNkJBQXFCLEtBQXJCLDZCQUFxQixRQWFoQztBQUNZLFFBQUEsNEJBQTRCLEdBQUcscUJBQXFCLENBQUM7QUFDbEUsU0FBZ0IsNkJBQTZCLENBQUMsTUFBVztJQUN2RCxRQUFRLE1BQU0sRUFBRTtRQUNkLEtBQUssQ0FBQyxDQUFDO1FBQ1AsS0FBSyx5QkFBeUI7WUFDNUIsT0FBTyxxQkFBcUIsQ0FBQyx1QkFBdUIsQ0FBQztRQUV2RCxLQUFLLENBQUMsQ0FBQztRQUNQLEtBQUssdUJBQXVCO1lBQzFCLE9BQU8scUJBQXFCLENBQUMscUJBQXFCLENBQUM7UUFFckQsS0FBSyxDQUFDLENBQUM7UUFDUCxLQUFLLGtCQUFrQjtZQUNyQixPQUFPLHFCQUFxQixDQUFDLGdCQUFnQixDQUFDO1FBRWhELEtBQUssQ0FBQyxDQUFDLENBQUM7UUFDUixLQUFLLGNBQWMsQ0FBQztRQUNwQjtZQUNFLE9BQU8scUJBQXFCLENBQUMsWUFBWSxDQUFDO0tBQzdDO0FBQ0gsQ0FBQztBQW5CRCxzRUFtQkM7QUFDRCxTQUFnQiwyQkFB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