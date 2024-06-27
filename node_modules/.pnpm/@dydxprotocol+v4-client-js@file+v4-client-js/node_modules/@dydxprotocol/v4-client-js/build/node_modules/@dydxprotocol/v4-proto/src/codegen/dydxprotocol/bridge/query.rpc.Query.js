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
exports.createRpcQueryExtension = exports.QueryClientImpl = void 0;
const _m0 = __importStar(require("protobufjs/minimal"));
const stargate_1 = require("@cosmjs/stargate");
const query_1 = require("./query");
class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.eventParams = this.eventParams.bind(this);
        this.proposeParams = this.proposeParams.bind(this);
        this.safetyParams = this.safetyParams.bind(this);
        this.acknowledgedEventInfo = this.acknowledgedEventInfo.bind(this);
        this.recognizedEventInfo = this.recognizedEventInfo.bind(this);
        this.delayedCompleteBridgeMessages = this.delayedCompleteBridgeMessages.bind(this);
    }
    eventParams(request = {}) {
        const data = query_1.QueryEventParamsRequest.encode(request).finish();
        const promise = this.rpc.request("dydxprotocol.bridge.Query", "EventParams", data);
        return promise.then(data => query_1.QueryEventParamsResponse.decode(new _m0.Reader(data)));
    }
    proposeParams(request = {}) {
        const data = query_1.QueryProposeParamsRequest.encode(request).finish();
        const promise = this.rpc.request("dydxprotocol.bridge.Query", "ProposeParams", data);
        return promise.then(data => query_1.QueryProposeParamsResponse.decode(new _m0.Reader(data)));
    }
    safetyParams(request = {}) {
        const data = query_1.QuerySafetyParamsRequest.encode(request).finish();
        const promise = this.rpc.request("dydxprotocol.bridge.Query", "SafetyParams", data);
        return promise.then(data => query_1.QuerySafetyParamsResponse.decode(new _m0.Reader(data)));
    }
    acknowledgedEventInfo(request = {}) {
        const data = query_1.QueryAcknowledgedEventInfoRequest.encode(request).finish();
        const promise = this.rpc.request("dydxprotocol.bridge.Query", "AcknowledgedEventInfo", data);
        return promise.then(data => query_1.QueryAcknowledgedEventInfoResponse.decode(new _m0.Reader(data)));
    }
    recognizedEventInfo(request = {}) {
        const data = query_1.QueryRecognizedEventInfoRequest.encode(request).finish();
        const promise = this.rpc.request("dydxprotocol.bridge.Query", "RecognizedEventInfo", data);
        return promise.then(data => query_1.QueryRecognizedEventInfoResponse.decode(new _m0.Reader(data)));
    }
    delayedCompleteBridgeMessages(request) {
        const data = query_1.QueryDelayedCompleteBridgeMessagesRequest.encode(request).finish();
        const promise = this.rpc.request("dydxprotocol.bridge.Query", "DelayedCompleteBridgeMessages", data);
        return promise.then(data => query_1.QueryDelayedCompleteBridgeMessagesResponse.decode(new _m0.Reader(data)));
    }
}
exports.QueryClientImpl = QueryClientImpl;
const createRpcQueryExtension = (base) => {
    const rpc = (0, stargate_1.createProtobufRpcClient)(base);
    const queryService = new QueryClientImpl(rpc);
    return {
        eventParams(request) {
            return queryService.eventParams(request);
        },
        proposeParams(request) {
            return queryService.proposeParams(request);
        },
        safetyParams(request) {
            return queryService.safetyParams(request);
        },
        acknowledgedEventInfo(request) {
            return queryService.acknowledgedEventInfo(request);
        },
        recognizedEventInfo(request) {
            return queryService.recognizedEventInfo(request);
        },
        delayedCompleteBridgeMessages(request) {
            return queryService.delayedCompleteBridgeMessages(request);
        }
    };
};
exports.createRpcQueryExtension = createRpcQueryExtension;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkucnBjLlF1ZXJ5LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL2JyaWRnZS9xdWVyeS5ycGMuUXVlcnkudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFDQSx3REFBMEM7QUFDMUMsK0NBQXdFO0FBQ3hFLG1DQUF5WjtBQWlDelosTUFBYSxlQUFlO0lBRzFCLFlBQVksR0FBUTtRQUNsQixJQUFJLENBQUMsR0FBRyxHQUFHLEdBQUcsQ0FBQztRQUNmLElBQUksQ0FBQyxXQUFXLEdBQUcsSUFBSSxDQUFDLFdBQVcsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDL0MsSUFBSSxDQUFDLGFBQWEsR0FBRyxJQUFJLENBQUMsYUFBYSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNuRCxJQUFJLENBQUMsWUFBWSxHQUFHLElBQUksQ0FBQyxZQUFZLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ2pELElBQUksQ0FBQyxxQkFBcUIsR0FBRyxJQUFJLENBQUMscUJBQXFCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ25FLElBQUksQ0FBQyxtQkFBbUIsR0FBRyxJQUFJLENBQUMsbUJBQW1CLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQy9ELElBQUksQ0FBQyw2QkFBNkIsR0FBRyxJQUFJLENBQUMsNkJBQTZCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQ3JGLENBQUM7SUFFRCxXQUFXLENBQUMsVUFBbUMsRUFBRTtRQUMvQyxNQUFNLElBQUksR0FBRywrQkFBdUIsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLENBQUMsTUFBTSxFQUFFLENBQUM7UUFDOUQsTUFBTSxPQUFP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