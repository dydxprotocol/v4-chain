"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createVestingAminoConverters = exports.isAminoMsgCreateVestingAccount = void 0;
function isAminoMsgCreateVestingAccount(msg) {
    return msg.type === "cosmos-sdk/MsgCreateVestingAccount";
}
exports.isAminoMsgCreateVestingAccount = isAminoMsgCreateVestingAccount;
function createVestingAminoConverters() {
    return {
        "/cosmos.vesting.v1beta1.MsgCreateVestingAccount": {
            aminoType: "cosmos-sdk/MsgCreateVestingAccount",
            toAmino: ({ fromAddress, toAddress, amount, endTime, delayed, }) => ({
                from_address: fromAddress,
                to_address: toAddress,
                amount: [...amount],
                end_time: endTime.toString(),
                delayed: delayed,
            }),
            fromAmino: ({ from_address, to_address, amount, end_time, delayed, }) => ({
                fromAddress: from_address,
                toAddress: to_address,
                amount: [...amount],
                endTime: BigInt(end_time),
                delayed: delayed,
            }),
        },
    };
}
exports.createVestingAminoConverters = createVestingAminoConverters;
//# sourceMappingURL=aminomessages.js.map