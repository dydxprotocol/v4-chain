"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createIbcAminoConverters = exports.isAminoMsgTransfer = void 0;
/* eslint-disable @typescript-eslint/naming-convention */
const amino_1 = require("@cosmjs/amino");
const tx_1 = require("cosmjs-types/ibc/applications/transfer/v1/tx");
function isAminoMsgTransfer(msg) {
    return msg.type === "cosmos-sdk/MsgTransfer";
}
exports.isAminoMsgTransfer = isAminoMsgTransfer;
function createIbcAminoConverters() {
    return {
        "/ibc.applications.transfer.v1.MsgTransfer": {
            aminoType: "cosmos-sdk/MsgTransfer",
            toAmino: ({ sourcePort, sourceChannel, token, sender, receiver, timeoutHeight, timeoutTimestamp, memo, }) => ({
                source_port: sourcePort,
                source_channel: sourceChannel,
                token: token,
                sender: sender,
                receiver: receiver,
                timeout_height: timeoutHeight
                    ? {
                        revision_height: (0, amino_1.omitDefault)(timeoutHeight.revisionHeight)?.toString(),
                        revision_number: (0, amino_1.omitDefault)(timeoutHeight.revisionNumber)?.toString(),
                    }
                    : {},
                timeout_timestamp: (0, amino_1.omitDefault)(timeoutTimestamp)?.toString(),
                memo: (0, amino_1.omitDefault)(memo),
            }),
            fromAmino: ({ source_port, source_channel, token, sender, receiver, timeout_height, timeout_timestamp, memo, }) => tx_1.MsgTransfer.fromPartial({
                sourcePort: source_port,
                sourceChannel: source_channel,
                token: token,
                sender: sender,
                receiver: receiver,
                timeoutHeight: timeout_height
                    ? {
                        revisionHeight: BigInt(timeout_height.revision_height || "0"),
                        revisionNumber: BigInt(timeout_height.revision_number || "0"),
                    }
                    : undefined,
                timeoutTimestamp: BigInt(timeout_timestamp || "0"),
                memo: memo ?? "",
            }),
        },
    };
}
exports.createIbcAminoConverters = createIbcAminoConverters;
//# sourceMappingURL=aminomessages.js.map