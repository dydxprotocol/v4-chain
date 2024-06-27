"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isMsgCancelUnbondingDelegationEncodeObject = exports.isMsgUndelegateEncodeObject = exports.isMsgEditValidatorEncodeObject = exports.isMsgDelegateEncodeObject = exports.isMsgCreateValidatorEncodeObject = exports.isMsgBeginRedelegateEncodeObject = exports.stakingTypes = void 0;
const tx_1 = require("cosmjs-types/cosmos/staking/v1beta1/tx");
exports.stakingTypes = [
    ["/cosmos.staking.v1beta1.MsgBeginRedelegate", tx_1.MsgBeginRedelegate],
    ["/cosmos.staking.v1beta1.MsgCreateValidator", tx_1.MsgCreateValidator],
    ["/cosmos.staking.v1beta1.MsgDelegate", tx_1.MsgDelegate],
    ["/cosmos.staking.v1beta1.MsgEditValidator", tx_1.MsgEditValidator],
    ["/cosmos.staking.v1beta1.MsgUndelegate", tx_1.MsgUndelegate],
    ["/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation", tx_1.MsgCancelUnbondingDelegation],
];
function isMsgBeginRedelegateEncodeObject(o) {
    return o.typeUrl === "/cosmos.staking.v1beta1.MsgBeginRedelegate";
}
exports.isMsgBeginRedelegateEncodeObject = isMsgBeginRedelegateEncodeObject;
function isMsgCreateValidatorEncodeObject(o) {
    return o.typeUrl === "/cosmos.staking.v1beta1.MsgCreateValidator";
}
exports.isMsgCreateValidatorEncodeObject = isMsgCreateValidatorEncodeObject;
function isMsgDelegateEncodeObject(object) {
    return object.typeUrl === "/cosmos.staking.v1beta1.MsgDelegate";
}
exports.isMsgDelegateEncodeObject = isMsgDelegateEncodeObject;
function isMsgEditValidatorEncodeObject(o) {
    return o.typeUrl === "/cosmos.staking.v1beta1.MsgEditValidator";
}
exports.isMsgEditValidatorEncodeObject = isMsgEditValidatorEncodeObject;
function isMsgUndelegateEncodeObject(object) {
    return object.typeUrl === "/cosmos.staking.v1beta1.MsgUndelegate";
}
exports.isMsgUndelegateEncodeObject = isMsgUndelegateEncodeObject;
function isMsgCancelUnbondingDelegationEncodeObject(object) {
    return (object.typeUrl ===
        "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation");
}
exports.isMsgCancelUnbondingDelegationEncodeObject = isMsgCancelUnbondingDelegationEncodeObject;
//# sourceMappingURL=messages.js.map