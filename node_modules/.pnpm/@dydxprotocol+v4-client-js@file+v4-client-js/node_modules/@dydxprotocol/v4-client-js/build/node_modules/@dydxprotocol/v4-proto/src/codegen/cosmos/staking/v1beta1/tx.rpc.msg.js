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
        this.createValidator = this.createValidator.bind(this);
        this.editValidator = this.editValidator.bind(this);
        this.delegate = this.delegate.bind(this);
        this.beginRedelegate = this.beginRedelegate.bind(this);
        this.undelegate = this.undelegate.bind(this);
        this.cancelUnbondingDelegation = this.cancelUnbondingDelegation.bind(this);
        this.updateParams = this.updateParams.bind(this);
    }
    createValidator(request) {
        const data = tx_1.MsgCreateValidator.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "CreateValidator", data);
        return promise.then(data => tx_1.MsgCreateValidatorResponse.decode(new _m0.Reader(data)));
    }
    editValidator(request) {
        const data = tx_1.MsgEditValidator.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "EditValidator", data);
        return promise.then(data => tx_1.MsgEditValidatorResponse.decode(new _m0.Reader(data)));
    }
    delegate(request) {
        const data = tx_1.MsgDelegate.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "Delegate", data);
        return promise.then(data => tx_1.MsgDelegateResponse.decode(new _m0.Reader(data)));
    }
    beginRedelegate(request) {
        const data = tx_1.MsgBeginRedelegate.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "BeginRedelegate", data);
        return promise.then(data => tx_1.MsgBeginRedelegateResponse.decode(new _m0.Reader(data)));
    }
    undelegate(request) {
        const data = tx_1.MsgUndelegate.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "Undelegate", data);
        return promise.then(data => tx_1.MsgUndelegateResponse.decode(new _m0.Reader(data)));
    }
    cancelUnbondingDelegation(request) {
        const data = tx_1.MsgCancelUnbondingDelegation.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "CancelUnbondingDelegation", data);
        return promise.then(data => tx_1.MsgCancelUnbondingDelegationResponse.decode(new _m0.Reader(data)));
    }
    updateParams(request) {
        const data = tx_1.MsgUpdateParams.encode(request).finish();
        const promise = this.rpc.request("cosmos.staking.v1beta1.Msg", "UpdateParams", data);
        return promise.then(data => tx_1.MsgUpdateParamsResponse.decode(new _m0.Reader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidHgucnBjLm1zZy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uL25vZGVfbW9kdWxlcy9AZHlkeHByb3RvY29sL3Y0LXByb3RvL3NyYy9jb2RlZ2VuL2Nvc21vcy9zdGFraW5nL3YxYmV0YTEvdHgucnBjLm1zZy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUNBLHdEQUEwQztBQUMxQyw2QkFBd1Y7QUEyQ3hWLE1BQWEsYUFBYTtJQUd4QixZQUFZLEdBQVE7UUFDbEIsSUFBSSxDQUFDLEdBQUcsR0FBRyxHQUFHLENBQUM7UUFDZixJQUFJLENBQUMsZUFBZSxHQUFHLElBQUksQ0FBQyxlQUFlLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3ZELElBQUksQ0FBQyxhQUFhLEdBQUcsSUFBSSxDQUFDLGFBQWEsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDbkQsSUFBSSxDQUFDLFFBQVEsR0FBRyxJQUFJLENBQUMsUUFBUSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUN6QyxJQUFJLENBQUMsZUFBZSxHQUFHLElBQUksQ0FBQyxlQUFlLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3ZELElBQUksQ0FBQyxVQUFVLEdBQUcsSUFBSSxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDN0MsSUFBSSxDQUFDLHlCQUF5QixHQUFHLElBQUksQ0FBQyx5QkFBeUIsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDM0UsSUFBSSxDQUFDLFlBQVksR0FBRyxJQUFJLENBQUMsWUFBWSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUNuRCxDQUFDO0lBRUQsZUFBZSxDQUFDLE9BQTJCO1FBQ3pDLE1BQU0sSUFBSSxHQUFHLHVCQUFrQixDQUFDLE1BQU0sQ0FBQyxPQUFP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