"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.groupTypes = void 0;
const tx_1 = require("cosmjs-types/cosmos/group/v1/tx");
exports.groupTypes = [
    ["/cosmos.group.v1.MsgCreateGroup", tx_1.MsgCreateGroup],
    ["/cosmos.group.v1.MsgCreateGroupPolicy", tx_1.MsgCreateGroupPolicy],
    ["/cosmos.group.v1.MsgCreateGroupWithPolicy", tx_1.MsgCreateGroupWithPolicy],
    ["/cosmos.group.v1.MsgExec", tx_1.MsgExec],
    ["/cosmos.group.v1.MsgLeaveGroup", tx_1.MsgLeaveGroup],
    ["/cosmos.group.v1.MsgSubmitProposal", tx_1.MsgSubmitProposal],
    ["/cosmos.group.v1.MsgUpdateGroupAdmin", tx_1.MsgUpdateGroupAdmin],
    ["/cosmos.group.v1.MsgUpdateGroupMembers", tx_1.MsgUpdateGroupMembers],
    ["/cosmos.group.v1.MsgUpdateGroupMetadata", tx_1.MsgUpdateGroupMetadata],
    ["/cosmos.group.v1.MsgUpdateGroupPolicyAdmin", tx_1.MsgUpdateGroupPolicyAdmin],
    ["/cosmos.group.v1.MsgUpdateGroupPolicyDecisionPolicy", tx_1.MsgUpdateGroupPolicyDecisionPolicy],
    ["/cosmos.group.v1.MsgUpdateGroupPolicyMetadata", tx_1.MsgUpdateGroupPolicyMetadata],
    ["/cosmos.group.v1.MsgVote", tx_1.MsgVote],
    ["/cosmos.group.v1.MsgWithdrawProposal", tx_1.MsgWithdrawProposal],
];
// There are no EncodeObject implementations for the new v1 message types because
// those things don't scale (https://github.com/cosmos/cosmjs/issues/1440). We need to
// address this more fundamentally. Users can use
// const msg = {
//   typeUrl: "/cosmos.group.v1.MsgCreateGroup",
//   value: MsgCreateGroup.fromPartial({ ... })
// }
// in their app.
//# sourceMappingURL=messages.js.map