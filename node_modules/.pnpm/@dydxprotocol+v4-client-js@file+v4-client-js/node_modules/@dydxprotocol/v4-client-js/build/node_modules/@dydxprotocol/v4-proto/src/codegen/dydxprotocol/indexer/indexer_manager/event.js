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
exports.IndexerTendermintBlock = exports.IndexerTendermintEvent = exports.IndexerEventsStoreValue = exports.IndexerTendermintEventWrapper = exports.indexerTendermintEvent_BlockEventToJSON = exports.indexerTendermintEvent_BlockEventFromJSON = exports.IndexerTendermintEvent_BlockEventSDKType = exports.IndexerTendermintEvent_BlockEvent = void 0;
const timestamp_1 = require("../../../google/protobuf/timestamp");
const _m0 = __importStar(require("protobufjs/minimal"));
const helpers_1 = require("../../../helpers");
/** enum to specify that the IndexerTendermintEvent is a block event. */
var IndexerTendermintEvent_BlockEvent;
(function (IndexerTendermintEvent_BlockEvent) {
    /** BLOCK_EVENT_UNSPECIFIED - Default value. This value is invalid and unused. */
    IndexerTendermintEvent_BlockEvent[IndexerTendermintEvent_BlockEvent["BLOCK_EVENT_UNSPECIFIED"] = 0] = "BLOCK_EVENT_UNSPECIFIED";
    /**
     * BLOCK_EVENT_BEGIN_BLOCK - BLOCK_EVENT_BEGIN_BLOCK indicates that the event was generated during
     * BeginBlock.
     */
    IndexerTendermintEvent_BlockEvent[IndexerTendermintEvent_BlockEvent["BLOCK_EVENT_BEGIN_BLOCK"] = 1] = "BLOCK_EVENT_BEGIN_BLOCK";
    /**
     * BLOCK_EVENT_END_BLOCK - BLOCK_EVENT_END_BLOCK indicates that the event was generated during
     * EndBlock.
     */
    IndexerTendermintEvent_BlockEvent[IndexerTendermintEvent_BlockEvent["BLOCK_EVENT_END_BLOCK"] = 2] = "BLOCK_EVENT_END_BLOCK";
    IndexerTendermintEvent_BlockEvent[IndexerTendermintEvent_BlockEvent["UNRECOGNIZED"] = -1] = "UNRECOGNIZED";
})(IndexerTendermintEvent_BlockEvent = exports.IndexerTendermintEvent_BlockEvent || (exports.IndexerTendermintEvent_BlockEvent = {}));
exports.IndexerTendermintEvent_BlockEventSDKType = IndexerTendermintEvent_BlockEvent;
function indexerTendermintEvent_BlockEventFromJSON(object) {
    switch (object) {
        case 0:
        case "BLOCK_EVENT_UNSPECIFIED":
            return IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_UNSPECIFIED;
        case 1:
        case "BLOCK_EVENT_BEGIN_BLOCK":
            return IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_BEGIN_BLOCK;
        case 2:
        case "BLOCK_EVENT_END_BLOCK":
            return IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_END_BLOCK;
        case -1:
        case "UNRECOGNIZED":
        default:
            return IndexerTendermintEvent_BlockEvent.UNRECOGNIZED;
    }
}
exports.indexerTendermintEvent_BlockEventFromJSON = indexerTendermintEvent_BlockEventFromJSON;
function indexerTendermintEvent_BlockEventToJSON(object) {
    switch (object) {
        case IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_UNSPECIFIED:
            return "BLOCK_EVENT_UNSPECIFIED";
        case IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_BEGIN_BLOCK:
            return "BLOCK_EVENT_BEGIN_BLOCK";
        case IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_END_BLOCK:
            return "BLOCK_EVENT_END_BLOCK";
        case IndexerTendermintEvent_BlockEvent.UNRECOGNIZED:
        default:
            return "UNRECOGNIZED";
    }
}
exports.indexerTendermintEvent_BlockEventToJSON = indexerTendermintEvent_BlockEventToJSON;
function createBaseIndexerTendermintEventWrapper() {
    return {
        event: undefined,
        txnHash: ""
    };
}
exports.IndexerTendermintEventWrapper = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.event !== undefined) {
            exports.IndexerTendermintEvent.encode(message.event, writer.uint32(10).fork()).ldelim();
        }
        if (message.txnHash !== "") {
            writer.uint32(18).string(message.txnHash);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseIndexerTendermintEventWrapper();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.event = exports.IndexerTendermintEvent.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.txnHash = reader.string();
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
        const message = createBaseIndexerTendermintEventWrapper();
        message.event = object.event !== undefined && object.event !== null ? exports.IndexerTendermintEvent.fromPartial(object.event) : undefined;
        message.txnHash = (_a = object.txnHash) !== null && _a !== void 0 ? _a : "";
        return message;
    }
};
function createBaseIndexerEventsStoreValue() {
    return {
        events: []
    };
}
exports.IndexerEventsStoreValue = {
    encode(message, writer = _m0.Writer.create()) {
        for (const v of message.events) {
            exports.IndexerTendermintEventWrapper.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseIndexerEventsStoreValue();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.events.push(exports.IndexerTendermintEventWrapper.decode(reader, reader.uint32()));
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
        const message = createBaseIndexerEventsStoreValue();
        message.events = ((_a = object.events) === null || _a === void 0 ? void 0 : _a.map(e => exports.IndexerTendermintEventWrapper.fromPartial(e))) || [];
        return message;
    }
};
function createBaseIndexerTendermintEvent() {
    return {
        subtype: "",
        transactionIndex: undefined,
        blockEvent: undefined,
        eventIndex: 0,
        version: 0,
        dataBytes: new Uint8Array()
    };
}
exports.IndexerTendermintEvent = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.subtype !== "") {
            writer.uint32(10).string(message.subtype);
        }
        if (message.transactionIndex !== undefined) {
            writer.uint32(24).uint32(message.transactionIndex);
        }
        if (message.blockEvent !== undefined) {
            writer.uint32(32).int32(message.blockEvent);
        }
        if (message.eventIndex !== 0) {
            writer.uint32(40).uint32(message.eventIndex);
        }
        if (message.version !== 0) {
            writer.uint32(48).uint32(message.version);
        }
        if (message.dataBytes.length !== 0) {
            writer.uint32(58).bytes(message.dataBytes);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseIndexerTendermintEvent();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.subtype = reader.string();
                    break;
                case 3:
                    message.transactionIndex = reader.uint32();
                    break;
                case 4:
                    message.blockEvent = reader.int32();
                    break;
                case 5:
                    message.eventIndex = reader.uint32();
                    break;
                case 6:
                    message.version = reader.uint32();
                    break;
                case 7:
                    message.dataBytes = reader.bytes();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        var _a, _b, _c, _d, _e, _f;
        const message = createBaseIndexerTendermintEvent();
        message.subtype = (_a = object.subtype) !== null && _a !== void 0 ? _a : "";
        message.transactionIndex = (_b = object.transactionIndex) !== null && _b !== void 0 ? _b : undefined;
        message.blockEvent = (_c = object.blockEvent) !== null && _c !== void 0 ? _c : undefined;
        message.eventIndex = (_d = object.eventIndex) !== null && _d !== void 0 ? _d : 0;
        message.version = (_e = object.version) !== null && _e !== void 0 ? _e : 0;
        message.dataBytes = (_f = object.dataBytes) !== null && _f !== void 0 ? _f : new Uint8Array();
        return message;
    }
};
function createBaseIndexerTendermintBlock() {
    return {
        height: 0,
        time: undefined,
        events: [],
        txHashes: []
    };
}
exports.IndexerTendermintBlock = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.height !== 0) {
            writer.uint32(8).uint32(message.height);
        }
        if (message.time !== undefined) {
            timestamp_1.Timestamp.encode((0, helpers_1.toTimestamp)(message.time), writer.uint32(18).fork()).ldelim();
        }
        for (const v of message.events) {
            exports.IndexerTendermintEvent.encode(v, writer.uint32(26).fork()).ldelim();
        }
        for (const v of message.txHashes) {
            writer.uint32(34).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseIndexerTendermintBlock();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.height = reader.uint32();
                    break;
                case 2:
                    message.time = (0, helpers_1.fromTimestamp)(timestamp_1.Timestamp.decode(reader, reader.uint32()));
                    break;
                case 3:
                    message.events.push(exports.IndexerTendermintEvent.decode(reader, reader.uint32()));
                    break;
                case 4:
                    message.txHashes.push(reader.string());
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
        const message = createBaseIndexerTendermintBlock();
        message.height = (_a = object.height) !== null && _a !== void 0 ? _a : 0;
        message.time = (_b = object.time) !== null && _b !== void 0 ? _b : undefined;
        message.events = ((_c = object.events) === null || _c === void 0 ? void 0 : _c.map(e => exports.IndexerTendermintEvent.fromPartial(e))) || [];
        message.txHashes = ((_d = object.txHashes) === null || _d === void 0 ? void 0 : _d.map(e => e)) || [];
        return message;
    }
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXZlbnQuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi9ub2RlX21vZHVsZXMvQGR5ZHhwcm90b2NvbC92NC1wcm90by9zcmMvY29kZWdlbi9keWR4cHJvdG9jb2wvaW5kZXhlci9pbmRleGVyX21hbmFnZXIvZXZlbnQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFBQSxrRUFBK0Q7QUFDL0Qsd0RBQTBDO0FBQzFDLDhDQUEyRTtBQUMzRSx3RUFBd0U7QUFFeEUsSUFBWSxpQ0FnQlg7QUFoQkQsV0FBWSxpQ0FBaUM7SUFDM0MsaUZBQWlGO0lBQ2pGLCtIQUEyQixDQUFBO0lBRTNCOzs7T0FHRztJQUNILCtIQUEyQixDQUFBO0lBRTNCOzs7T0FHRztJQUNILDJIQUF5QixDQUFBO0lBQ3pCLDBHQUFpQixDQUFBO0FBQ25CLENBQUMsRUFoQlcsaUNBQWlDLEdBQWpDLHlDQUFpQyxLQUFqQyx5Q0FBaUMsUUFnQjVDO0FBQ1ksUUFBQSx3Q0FBd0MsR0FBRyxpQ0FBaUMsQ0FBQztBQUMxRixTQUFnQix5Q0FBeUMsQ0FBQyxNQUFXO0lBQ25FLFFBQVEsTUFBTSxFQUFFO1FBQ2QsS0FBSyxDQUFDLENBQUM7UUFDUCxLQUFLLHlCQUF5QjtZQUM1QixPQUFPLGlDQUFpQyxDQUFDLHVCQUF1QixDQUFDO1FBRW5FLEtBQUssQ0FBQyxDQUFDO1FBQ1AsS0FBSyx5QkFBeUI7WUFDNUIsT0FBTyxpQ0FBaUMsQ0FBQyx1QkFBdUIsQ0FBQztRQUVuRSxLQUFLLENBQUMsQ0FBQztRQUNQLEtBQUssdUJBQXVCO1lBQzFCLE9BQU8saUNBQWlDLENBQUMscUJBQXFCLENBQUM7UUFFakUsS0FBSyxDQUFDLENBQUMsQ0FBQztRQUNSLEtBQUssY0FBYyxDQUFDO1FBQ3BCO1lBQ0UsT0FBTyxpQ0FBaUMsQ0FBQyxZQUFZLENBQUM7S0FDekQ7QUFDSCxDQUFDO0FBbkJELDhGQW1CQztBQUNELFNBQWdCLHVDQUF1QyxDQUFDLE1BQXlDO0lBQy9GLFFBQVEsTUFBTSxFQUFFO1FBQ2QsS0FBSyxpQ0FBaUMsQ0FBQyx1QkFBdUI7WUFDNUQsT0FBTyx5QkFBeUIsQ0FBQztRQUVuQyxLQUFLLGlDQUFpQyxDQUFDLHVCQUF1QjtZQUM1RCxPQUFPLHlCQUF5QixDQUFDO1FBRW5DLEtBQUssaUNBQWlDLENBQUMscUJBQXFCO1lBQzFELE9BQU8sdUJBQXVCLENBQUM7UUFFakMsS0FBSyxpQ0FBaUMsQ0FBQyxZQUFZLENBQUM7UUFDcEQ7WUFDRSxPQUFPLGNBQWMsQ0FBQztLQUN6QjtBQUNILENBQUM7QUFmRCwwRkFlQztBQXFHRCxTQUFTLHVDQUF1QztJQUM5QyxPQUFPO1FBQ0wsS0FBSyxFQUFFLFNBQVM7UUFDaEIsT0FBTyxFQUFFLEVBQUU7S0FDWixDQUFDO0FBQ0osQ0FBQztBQUVZLFFBQUEsNkJBQTZCLEdBQUc7SUFDM0MsTUFBTSxDQUFDLE9BQXNDLEVBQUUsU0FBcUIsR0FBRyxDQUFDLE1BQU0sQ0FBQyxNQUFNLEVBQUU7UUFDckYsSUFBSSxPQUFPLENBQUMsS0FBSyxLQUFLLFNBQVMsRUFBRTtZQUMvQiw4QkFBc0IsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLEtBQUssRUFBRSxNQUFNLENBQUMsTUFBTSxDQUFDLEVBQUUsQ0FBQyxDQUFDLElBQUksRUFBRSxDQUFDLENBQUMsTUFBTSxFQUFFLENBQUM7U0FDakY7UUFFRCxJQUFJLE9BQU8sQ0FBQyxPQUFPLEtBQUssRUFBRSxFQUFFO1lBQzFCLE1BQU0sQ0FBQyxNQUFNLENBQUMsRUFBRSxDQUFDLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxPQUFPLENBQUMsQ0FBQztTQUMzQztRQUVELE9BQU8sTUFBTSxDQUFDO0lBQ2hCLENBQUM7SUFFRCxNQUFNLENBQUMsS0FBOEIsRUFBRSxNQUFlO1FBQ3BELE1BQU0sTUFBTSxHQUFHLEtBQUssWUFBWSxHQUFHLENBQUMsTUFBTSxDQUFDLENBQUMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLElBQUksR0FBRyxDQUFDLE1BQU0sQ0FBQyxLQUFLLENBQUMsQ0FBQztRQUMzRSxJQUFJLEdBQUcsR0FBRyxNQUFNLEtBQUssU0FBUyxDQUFDLENBQUMsQ0FBQyxNQUFNLENBQUMsR0FBRyxDQUFDLENBQUMsQ0FBQyxNQUFNLENBQUMsR0FBRyxHQUFHLE1BQU0sQ0FBQztRQUNsRSxNQUFNLE9BQU8sR0FBRyx1Q0FBdUMsRUFBRSxDQUFDO1FBRTFELE9BQU8sTUFBTSxDQUFDLEdBQUcsR0FBRyxHQUFHLEVBQUU7WUFDdkIsTUFBTSxHQUFHLEdBQUcsTUFBTSxDQUFDLE1BQU0sRUFBRSxDQUFDO1lBRTVCLFFBQVEsR0FBRyxLQUFLLENBQUMsRUFBRTtnQkFDakIsS0FBSyxDQUFDO29CQUNKLE9BQU8sQ0FBQyxLQUFLLEdBQUcsOEJBQXNCLENBQUMsTUFBTSxDQUFDLE1BQU0sRUFBRSxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUMsQ0FBQztvQkFDdkUsTUFBTTtnQkFFUixLQUFLLENBQUM7b0JBQ0osT0FBTyxDQUFDLE9BQU8sR0FBRyxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUM7b0JBQ2xDLE1BQU07Z0JBRVI7b0JBQ0UsTUFBTSxDQUFDLFFBQVEsQ0FBQyxHQUFHLEdBQUcsQ0FBQyxDQUFDLENBQUM7b0JBQ3pCLE1BQU07YUFDVDtTQUNGO1FBRUQsT0FBTyxPQUFPLENBQUM7SUFDakIsQ0FBQztJQUVELFdBQVcsQ0FBQyxNQUFrRDs7UUFDNUQsTUFBTSxPQUFPLEdBQUcsdUNBQXVDLEVBQUUsQ0FBQztRQUMxRCxPQUFPLENBQUMsS0FBSyxHQUFHLE1BQU0sQ0FBQyxLQUFLLEtBQUssU0FBUyxJQUFJLE1BQU0sQ0FBQyxLQUFLLEtBQUssSUFBSSxDQUFDLENBQUMsQ0FBQyw4QkFBc0IsQ0FBQyxXQUFXLENBQUMsTUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQyxTQUFTLENBQUM7UUFDbkksT0FBTyxDQUFDLE9BQU8sR0FBRyxNQUFBLE1BQU0sQ0FBQyxPQUFPLG1DQUFJLEVBQUUsQ0FBQztRQUN2QyxPQUFPLE9BQU8sQ0FBQztJQUNqQixDQUFDO0NBRUYsQ0FBQztBQUVGLFNBQVMsaUNBQWlDO0lBQ3hDLE9BQU87UUFDTCxNQUFNLEVBQUUsRUFBRTtLQUNYLENBQUM7QUFDSixDQUFDO0FBRVksUUFBQSx1QkFBdUIsR0FBRztJQUNyQyxNQUFNLENBQUMsT0FBZ0MsRUFBRSxTQUFxQixHQUFHLENBQUMsTUFBTSxDQUFDLE1BQU0sRUFBRTtRQUMvRSxLQUFLLE1BQU0sQ0FBQyxJQUFJLE9BQU8sQ0FBQyxNQUFNLEVBQUU7WUFDOUIscUNBQTZCLENBQUMsTUFBTSxDQUFDLENBQUUsRUFBRSxNQUFNLENBQUMsTUFBTSxDQUFDLEVBQUUsQ0FBQyxDQUFDLElBQUksRUFBRSxDQUFDLENBQUMsTUFBTSxFQUFFLENBQUM7U0FDN0U7UUFFRCxPQUFPLE1BQU0sQ0FBQztJQUNoQixDQUFDO0lBRUQsTUFBTSxDQUFDLEtBQThCLEVBQUUsTUFBZTtRQUNwRCxNQUFNLE1BQU0sR0FBRyxLQUFLLFlBQVksR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxJQUFJLEdBQUcsQ0FBQyxNQUFNLENBQUMsS0FBSyxDQUFDLENBQUM7UUFDM0UsSUFBSSxHQUFHLEdBQUcsTUFBTSxLQUFLLFNBQVMsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsR0FBRyxNQUFNLENBQUM7UUFDbEUsTUFBTSxPQUFPLEdBQUcsaUNBQWlDLEVBQUUsQ0FBQztRQUVwRCxPQUFPLE1BQU0sQ0FBQyxHQUFHLEdBQUcsR0FBRyxFQUFFO1lBQ3ZCLE1BQU0sR0FBRyxHQUFHLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQztZQUU1QixRQUFRLEdBQUcsS0FBSyxDQUFDLEVBQUU7Z0JBQ2pCLEtBQUssQ0FBQztvQkFDSixPQUFPLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxxQ0FBNkIsQ0FBQyxNQUFNLENBQUMsTUFBTSxFQUFFLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQyxDQUFDLENBQUM7b0JBQ25GLE1BQU07Z0JBRVI7b0JBQ0UsTUFBTSxDQUFDLFFBQVEsQ0FBQyxHQUFHLEdBQUcsQ0FBQyxDQUFDLENBQUM7b0JBQ3pCLE1BQU07YUFDVDtTQUNGO1FBRUQsT0FBTyxPQUFPLENBQUM7SUFDakIsQ0FBQztJQUVELFdBQVcsQ0FBQyxNQUE0Qzs7UUFDdEQsTUFBTSxPQUFPLEdBQUcsaUNBQWlDLEVBQUUsQ0FBQztRQUNwRCxPQUFPLENBQUMsTUFBTSxHQUFHLENBQUEsTUFBQSxNQUFNLENBQUMsTUFBTSwwQ0FBRSxHQUFHLENBQUMsQ0FBQyxDQUFDLEVBQUUsQ0FBQyxxQ0FBNkIsQ0FBQyxXQUFXLENBQUMsQ0FBQyxDQUFDLENBQUMsS0FBSSxFQUFFLENBQUM7UUFDN0YsT0FBTyxPQUFPLENBQUM7SUFDakIsQ0FBQztDQUVGLENBQUM7QUFFRixTQUFTLGdDQUFnQztJQUN2QyxPQUFPO1FBQ0wsT0FBTyxFQUFFLEVBQUU7UUFDWCxnQkFBZ0IsRUFBRSxTQUFTO1FBQzNCLFVBQVUsRUFBRSxTQUFTO1FBQ3JCLFVBQVUsRUFBRSxDQUFDO1FBQ2IsT0FBTyxFQUFFLENBQUM7UUFDVixTQUFTLEVBQUUsSUFBSSxVQUFVLEVBQUU7S0FDNUIsQ0FBQztBQUNKLENBQUM7QUFFWSxRQUFBLHNCQUFzQixHQUFHO0lBQ3BDLE1BQU0sQ0FBQyxPQUErQixFQUFFLFNBQXFCLEdBQUcsQ0FBQyxNQUFNLENBQUMsTUFBTSxFQUFFO1FBQzlFLElBQUksT0FBTyxDQUFDLE9BQU8sS0FBSyxFQUFFLEVBQUU7WUFDMUIsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLE9BQU8sQ0FBQyxDQUFDO1NBQzNDO1FBRUQsSUFBSSxPQUFPLENBQUMsZ0JBQWdCLEtBQUssU0FBUyxFQUFFO1lBQzFDLE1BQU0sQ0FBQyxNQUFNLENBQUMsRUFBRSxDQUFDLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxnQkFBZ0IsQ0FBQyxDQUFDO1NBQ3BEO1FBRUQsSUFBSSxPQUFPLENBQUMsVUFBVSxLQUFLLFNBQVMsRUFBRTtZQUNwQyxNQUFNLENBQUMsTUFBTSxDQUFDLEVBQUUsQ0FBQyxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsVUFBVSxDQUFDLENBQUM7U0FDN0M7UUFFRCxJQUFJLE9BQU8sQ0FBQyxVQUFVLEtBQUssQ0FBQyxFQUFFO1lBQzVCLE1BQU0sQ0FBQyxNQUFNLENBQUMsRUFBRSxDQUFDLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxVQUFVLENBQUMsQ0FBQztTQUM5QztRQUVELElBQUksT0FBTyxDQUFDLE9BQU8sS0FBSyxDQUFDLEVBQUU7WUFDekIsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLE9BQU8sQ0FBQyxDQUFDO1NBQzNDO1FBRUQsSUFBSSxPQUFPLENBQUMsU0FBUyxDQUFDLE1BQU0sS0FBSyxDQUFDLEVBQUU7WUFDbEMsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLFNBQVMsQ0FBQyxDQUFDO1NBQzVDO1FBRUQsT0FBTyxNQUFNLENBQUM7SUFDaEIsQ0FBQztJQUVELE1BQU0sQ0FBQyxLQUE4QixFQUFFLE1BQWU7UUFDcEQsTUFBTSxNQUFNLEdBQUcsS0FBSyxZQUFZLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQyxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsSUFBSSxHQUFHLENBQUMsTUFBTSxDQUFDLEtBQUssQ0FBQyxDQUFDO1FBQzNFLElBQUksR0FBRyxHQUFHLE1BQU0sS0FBSyxTQUFTLENBQUMsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxHQUFHLENBQUMsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxHQUFHLEdBQUcsTUFBTSxDQUFDO1FBQ2xFLE1BQU0sT0FBTyxHQUFHLGdDQUFnQyxFQUFFLENBQUM7UUFFbkQsT0FBTyxNQUFNLENBQUMsR0FBRyxHQUFHLEdBQUcsRUFBRTtZQUN2QixNQUFNLEdBQUcsR0FBRyxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUM7WUFFNUIsUUFBUSxHQUFHLEtBQUssQ0FBQyxFQUFFO2dCQUNqQixLQUFLLENBQUM7b0JBQ0osT0FBTyxDQUFDLE9BQU8sR0FBRyxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUM7b0JBQ2xDLE1BQU07Z0JBRVIsS0FBSyxDQUFDO29CQUNKLE9BQU8sQ0FBQyxnQkFBZ0IsR0FBRyxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUM7b0JBQzNDLE1BQU07Z0JBRVIsS0FBSyxDQUFDO29CQUNKLE9BQU8sQ0FBQyxVQUFVLEdBQUksTUFBTSxDQUFDLEtBQUssRUFBVSxDQUFDO29CQUM3QyxNQUFNO2dCQUVSLEtBQUssQ0FBQztvQkFDSixPQUFPLENBQUMsVUFBVSxHQUFHLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQztvQkFDckMsTUFBTTtnQkFFUixLQUFLLENBQUM7b0JBQ0osT0FBTyxDQUFDLE9BQU8sR0FBRyxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUM7b0JBQ2xDLE1BQU07Z0JBRVIsS0FBSyxDQUFDO29CQUNKLE9BQU8sQ0FBQyxTQUFTLEdBQUcsTUFBTSxDQUFDLEtBQUssRUFBRSxDQUFDO29CQUNuQyxNQUFNO2dCQUVSO29CQUNFLE1BQU0sQ0FBQyxRQUFRLENBQUMsR0FBRyxHQUFHLENBQUMsQ0FBQyxDQUFDO29CQUN6QixNQUFNO2FBQ1Q7U0FDRjtRQUVELE9BQU8sT0FBTyxDQUFDO0lBQ2pCLENBQUM7SUFFRCxXQUFXLENBQUMsTUFBMkM7O1FBQ3JELE1BQU0sT0FBTyxHQUFHLGdDQUFnQyxFQUFFLENBQUM7UUFDbkQsT0FBTyxDQUFDLE9BQU8sR0FBRyxNQUFBLE1BQU0sQ0FBQyxPQUFPLG1DQUFJLEVBQUUsQ0FBQztRQUN2QyxPQUFPLENBQUMsZ0JBQWdCLEdBQUcsTUFBQSxNQUFNLENBQUMsZ0JBQWdCLG1DQUFJLFNBQVMsQ0FBQztRQUNoRSxPQUFPLENBQUMsVUFBVSxHQUFHLE1BQUEsTUFBTSxDQUFDLFVBQVUsbUNBQUksU0FBUyxDQUFDO1FBQ3BELE9BQU8sQ0FBQyxVQUFVLEdBQUcsTUFBQSxNQUFNLENBQUMsVUFBVSxtQ0FBSSxDQUFDLENBQUM7UUFDNUMsT0FBTyxDQUFDLE9BQU8sR0FBRyxNQUFBLE1BQU0sQ0FBQyxPQUFPLG1DQUFJLENBQUMsQ0FBQztRQUN0QyxPQUFPLENBQUMsU0FBUyxHQUFHLE1BQUEsTUFBTSxDQUFDLFNBQVMsbUNBQUksSUFBSSxVQUFVLEVBQUUsQ0FBQztRQUN6RCxPQUFPLE9BQU8sQ0FBQztJQUNqQixDQUFDO0NBRUYsQ0FBQztBQUVGLFNBQVMsZ0NBQWdDO0lBQ3ZDLE9BQU87UUFDTCxNQUFNLEVBQUUsQ0FBQztRQUNULElBQUksRUFBRSxTQUFTO1FBQ2YsTUFBTSxFQUFFLEVBQUU7UUFDVixRQUFRLEVBQUUsRUFBRTtLQUNiLENBQUM7QUFDSixDQUFDO0FBRVksUUFBQSxzQkFBc0IsR0FBRztJQUNwQyxNQUFNLENBQUMsT0FBK0IsRUFBRSxTQUFxQixHQUFHLENBQUMsTUFBTSxDQUFDLE1BQU0sRUFBRTtRQUM5RSxJQUFJLE9BQU8sQ0FBQyxNQUFNLEtBQUssQ0FBQyxFQUFFO1lBQ3hCLE1BQU0sQ0FBQyxNQUFNLENBQUMsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxNQUFNLENBQUMsQ0FBQztTQUN6QztRQUVELElBQUksT0FBTyxDQUFDLElBQUksS0FBSyxTQUFTLEVBQUU7WUFDOUIscUJBQVMsQ0FBQyxNQUFNLENBQUMsSUFBQSxxQkFBVyxFQUFDLE9BQU8sQ0FBQyxJQUFJLENBQUMsRUFBRSxNQUFNLENBQUMsTUFBTSxDQUFDLEVBQUUsQ0FBQyxDQUFDLElBQUksRUFBRSxDQUFDLENBQUMsTUFBTSxFQUFFLENBQUM7U0FDaEY7UUFFRCxLQUFLLE1BQU0sQ0FBQyxJQUFJLE9BQU8sQ0FBQyxNQUFNLEVBQUU7WUFDOUIsOEJBQXNCLENBQUMsTUFBTSxDQUFDLENBQUUsRUFBRSxNQUFNLENBQUMsTUFBTSxDQUFDLEVBQUUsQ0FBQyxDQUFDLElBQUksRUFBRSxDQUFDLENBQUMsTUFBTSxFQUFFLENBQUM7U0FDdEU7UUFFRCxLQUFLLE1BQU0sQ0FBQyxJQUFJLE9BQU8sQ0FBQyxRQUFRLEVBQUU7WUFDaEMsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxNQUFNLENBQUMsQ0FBRSxDQUFDLENBQUM7U0FDOUI7UUFFRCxPQUFPLE1BQU0sQ0FBQztJQUNoQixDQUFDO0lBRUQsTUFBTSxDQUFDLEtBQThCLEVBQUUsTUFBZTtRQUNwRCxNQUFNLE1BQU0sR0FBRyxLQUFLLFlBQVksR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxJQUFJLEdBQUcsQ0FBQyxNQUFNLENBQUMsS0FBSyxDQUFDLENBQUM7UUFDM0UsSUFBSSxHQUFHLEdBQUcsTUFBTSxLQUFLLFNBQVMsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsTUFBTSxDQUFDLEdBQUcsR0FBRyxNQUFNLENBQUM7UUFDbEUsTUFBTSxPQUFPLEdBQUcsZ0NBQWdDLEVBQUUsQ0FBQztRQUVuRCxPQUFPLE1BQU0sQ0FBQyxHQUFHLEdBQUcsR0FBRyxFQUFFO1lBQ3ZCLE1BQU0sR0FBRyxHQUFHLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQztZQUU1QixRQUFRLEdBQUcsS0FBSyxDQUFDLEVBQUU7Z0JBQ2pCLEtBQUssQ0FBQztvQkFDSixPQUFPLENBQUMsTUFBTSxHQUFHLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQztvQkFDakMsTUFBTTtnQkFFUixLQUFLLENBQUM7b0JBQ0osT0FBTyxDQUFDLElBQUksR0FBRyxJQUFBLHVCQUFhLEVBQUMscUJBQVMsQ0FBQyxNQUFNLENBQUMsTUFBTSxFQUFFLE1BQU0sQ0FBQyxNQUFNLEVBQUUsQ0FBQyxDQUFDLENBQUM7b0JBQ3hFLE1BQU07Z0JBRVIsS0FBSyxDQUFDO29CQUNKLE9BQU8sQ0FBQyxNQUFNLENBQUMsSUFBSSxDQUFDLDhCQUFzQixDQUFDLE1BQU0sQ0FBQyxNQUFNLEVBQUUsTUFBTSxDQUFDLE1BQU0sRUFBRSxDQUFDLENBQUMsQ0FBQztvQkFDNUUsTUFBTTtnQkFFUixLQUFLLENBQUM7b0JBQ0osT0FBTyxDQUFDLFFBQVEsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLE1BQU0sRUFBRSxDQUFDLENBQUM7b0JBQ3ZDLE1BQU07Z0JBRVI7b0JBQ0UsTUFBTSxDQUFDLFFBQVEsQ0FBQyxHQUFHLEdBQUcsQ0FBQyxDQUFDLENBQUM7b0JBQ3pCLE1BQU07YUFDVDtTQUNGO1FBRUQsT0FBTyxPQUFPLENBQUM7SUFDakIsQ0FBQztJQUVELFdBQVcsQ0FBQyxNQUEyQzs7UUFDckQsTUFBTSxPQUFPLEdBQUcsZ0NBQWdDLEVBQUUsQ0FBQztRQUNuRCxPQUFPLENBQUMsTUFBTSxHQUFHLE1BQUEsTUFBTSxDQUFDLE1BQU0sbUNBQUksQ0FBQyxDQUFDO1FBQ3BDLE9BQU8sQ0FBQyxJQUFJLEdBQUcsTUFBQSxNQUFNLENBQUMsSUFBSSxtQ0FBSSxTQUFTLENBQUM7UUFDeEMsT0FBTyxDQUFDLE1BQU0sR0FBRyxDQUFBLE1BQUEsTUFBTSxDQUFDLE1BQU0sMENBQUUsR0FBRyxDQUFDLENBQUMsQ0FBQyxFQUFFLENBQUMsOEJBQXNCLENBQUMsV0FBVyxDQUFDLENBQUMsQ0FBQyxDQUFDLEtBQUksRUFBRSxDQUFDO1FBQ3RGLE9BQU8sQ0FBQyxRQUFRLEdBQUcsQ0FBQSxNQUFBLE1BQU0sQ0FBQyxRQUFRLDBDQUFFLEdBQUcsQ0FBQyxDQUFDLENBQUMsRUFBRSxDQUFDLENBQUMsQ0FBQyxLQUFJLEVBQUUsQ0FBQztRQUN0RCxPQUFPLE9BQU8sQ0FBQztJQUNqQixDQUFDO0NBRUYsQ0FBQyJ9