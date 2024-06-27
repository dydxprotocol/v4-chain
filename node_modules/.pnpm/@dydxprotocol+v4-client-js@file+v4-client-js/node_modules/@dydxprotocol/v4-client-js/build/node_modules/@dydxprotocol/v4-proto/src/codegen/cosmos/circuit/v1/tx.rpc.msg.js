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
exports.MsgClientImpl = void 0;
const _m0 = __importStar(require("protobufjs/minimal"));
const tx_1 = require("./tx");
class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.authorizeCircuitBreaker = this.authorizeCircuitBreaker.bind(this);
        this.tripCircuitBreaker = this.tripCircuitBreaker.bind(this);
        this.resetCircuitBreaker = this.resetCircuitBreaker.bind(this);
    }
    authorizeCircuitBreaker(request) {
        const data = tx_1.MsgAuthorizeCircuitBreaker.encode(request).finish();
        const promise = this.rpc.request("cosmos.circuit.v1.Msg", "AuthorizeCircuitBreaker", data);
        return promise.then(data => tx_1.MsgAuthorizeCircuitBreakerResponse.decode(new _m0.Reader(data)));
    }
    tripCircuitBreaker(request) {
        const data = tx_1.MsgTripCircuitBreaker.encode(request).finish();
        const promise = this.rpc.request("cosmos.circuit.v1.Msg", "TripCircuitBreaker", data);
        return promise.then(data => tx_1.MsgTripCircuitBreakerResponse.decode(new _m0.Reader(data)));
    }
    resetCircuitBreaker(request) {
        const data = tx_1.MsgResetCircuitBreaker.encode(request).finish();
        const promise = this.rpc.request("cosmos.circuit.v1.Msg", "ResetCircuitBreaker", data);
        return promise.then(data => tx_1.MsgResetCircuitBreakerResponse.decode(new _m0.Reader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidHgucnBjLm1zZy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uL25vZGVfbW9kdWxlcy9AZHlkeHByb3RvY29sL3Y0LXByb3RvL3NyYy9jb2RlZ2VuL2Nvc21vcy9jaXJjdWl0L3YxL3R4LnJwYy5tc2cudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFDQSx3REFBMEM7QUFDMUMsNkJBQW9NO0FBbUJwTSxNQUFhLGFBQWE7SUFHeEIsWUFBWSxHQUFRO1FBQ2xCLElBQUksQ0FBQyxHQUFHLEdBQUcsR0FBRyxDQUFDO1FBQ2YsSUFBSSxDQUFDLHVCQUF1QixHQUFHLElBQUksQ0FBQyx1QkFBdUIsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDdkUsSUFBSSxDQUFDLGtCQUFrQixHQUFHLElBQUksQ0FBQyxrQkFBa0IsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDN0QsSUFBSSxDQUFDLG1CQUFtQixHQUFHLElBQUksQ0FBQyxtQkFBbUIsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7SUFDakUsQ0FBQztJQUVELHVCQUF1QixDQUFDLE9BQW1DO1FBQ3pELE1BQU0sSUFBSSxHQUFHLCtCQUEwQixDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMsQ0FBQyxNQUFNLEVBQUUsQ0FBQztRQUNqRSxNQUFNLE9BQU8sR0FBRyxJQUFJLENBQUMsR0FBRyxDQUFDLE9BQU8sQ0FBQyx1QkFBdUIsRUFBRSx5QkFBeUIsRUFBRSxJQUFJLENBQUMsQ0FBQztRQUMzRixPQUFPLE9BQU8sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLEVBQUUsQ0FBQyx1Q0FBa0MsQ0FBQyxNQUFNLENBQUMsSUFBSSxHQUFHLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQztJQUMvRixDQUFDO0lBRUQsa0JBQWtCLENBQUMsT0FBOEI7UUFDL0MsTUFBTSxJQUFJLEdBQUcsMEJBQXFCLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxDQUFDLE1BQU0sRUFBRSxDQUFDO1FBQzVELE1BQU0sT0FBTyxHQUFHLElBQUksQ0FBQyxHQUFHLENBQUMsT0FBTyxDQUFDLHVCQUF1QixFQUFFLG9CQUFvQixFQUFFLElBQUksQ0FBQyxDQUFDO1FBQ3RGLE9BQU8sT0FBTyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsRUFBRSxDQUFDLGtDQUE2QixDQUFDLE1BQU0sQ0FBQyxJQUFJLEdBQUcsQ0FBQyxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFDO0lBQzFGLENBQUM7SUFFRCxtQkFBbUIsQ0FBQyxPQUErQjtRQUNqRCxNQUFNLElBQUksR0FBRywyQkFBc0IsQ0FBQyxNQUFNLENBQUMsT0FBTyxDQUFDLENBQUMsTUFBTSxFQUFFLENBQUM7UUFDN0QsTUFBTSxPQUFPLEdBQUcsSUFBSSxDQUFDLEdBQUcsQ0FBQyxPQUFPLENBQUMsdUJBQXVCLEVBQUUscUJBQXFCLEVBQUUsSUFBSSxDQUFDLENBQUM7UUFDdkYsT0FBTyxPQUFPLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxFQUFFLENBQUMsbUNBQThCLENBQUMsTUFBTSxDQUFDLElBQUksR0FBRyxDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDLENBQUM7SUFDM0YsQ0FBQztDQUVGO0FBNUJELHNDQTRCQyJ9