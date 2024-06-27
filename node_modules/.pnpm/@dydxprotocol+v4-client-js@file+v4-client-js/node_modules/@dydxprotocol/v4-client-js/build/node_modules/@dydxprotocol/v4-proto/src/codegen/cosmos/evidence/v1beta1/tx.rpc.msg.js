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
        this.submitEvidence = this.submitEvidence.bind(this);
    }
    submitEvidence(request) {
        const data = tx_1.MsgSubmitEvidence.encode(request).finish();
        const promise = this.rpc.request("cosmos.evidence.v1beta1.Msg", "SubmitEvidence", data);
        return promise.then(data => tx_1.MsgSubmitEvidenceResponse.decode(new _m0.Reader(data)));
    }
}
exports.MsgClientImpl = MsgClientImpl;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidHgucnBjLm1zZy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uLy4uL25vZGVfbW9kdWxlcy9AZHlkeHByb3RvY29sL3Y0LXByb3RvL3NyYy9jb2RlZ2VuL2Nvc21vcy9ldmlkZW5jZS92MWJldGExL3R4LnJwYy5tc2cudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFDQSx3REFBMEM7QUFDMUMsNkJBQW9FO0FBVXBFLE1BQWEsYUFBYTtJQUd4QixZQUFZLEdBQVE7UUFDbEIsSUFBSSxDQUFDLEdBQUcsR0FBRyxHQUFHLENBQUM7UUFDZixJQUFJLENBQUMsY0FBYyxHQUFHLElBQUksQ0FBQyxjQUFjLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQ3ZELENBQUM7SUFFRCxjQUFjLENBQUMsT0FBMEI7UUFDdkMsTUFBTSxJQUFJLEdBQUcsc0JBQWlCLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxDQUFDLE1BQU0sRUFBRSxDQUFDO1FBQ3hELE1BQU0sT0FBTyxHQUFHLElBQUksQ0FBQyxHQUFHLENBQUMsT0FBTyxDQUFDLDZCQUE2QixFQUFFLGdCQUFnQixFQUFFLElBQUksQ0FBQyxDQUFDO1FBQ3hGLE9BQU8sT0FBTyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsRUFBRSxDQUFDLDhCQUF5QixDQUFDLE1BQU0sQ0FBQyxJQUFJLEdBQUcsQ0FBQyxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxDQUFDO0lBQ3RGLENBQUM7Q0FFRjtBQWRELHNDQWNDIn0=