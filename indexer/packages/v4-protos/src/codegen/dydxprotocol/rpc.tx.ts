import { Rpc } from "../helpers";
export const createRPCMsgClient = async ({
  rpc
}: {
  rpc: Rpc;
}) => ({
  dydxprotocol: {
    blocktime: new (await import("./blocktime/tx.rpc.msg")).MsgClientImpl(rpc),
    bridge: new (await import("./bridge/tx.rpc.msg")).MsgClientImpl(rpc),
    clob: new (await import("./clob/tx.rpc.msg")).MsgClientImpl(rpc),
    delaymsg: new (await import("./delaymsg/tx.rpc.msg")).MsgClientImpl(rpc),
    feetiers: new (await import("./feetiers/tx.rpc.msg")).MsgClientImpl(rpc),
    ibcratelimit: new (await import("./ibcratelimit/tx.rpc.msg")).MsgClientImpl(rpc),
    perpetuals: new (await import("./perpetuals/tx.rpc.msg")).MsgClientImpl(rpc),
    prices: new (await import("./prices/tx.rpc.msg")).MsgClientImpl(rpc),
    rewards: new (await import("./rewards/tx.rpc.msg")).MsgClientImpl(rpc),
    sending: new (await import("./sending/tx.rpc.msg")).MsgClientImpl(rpc),
    stats: new (await import("./stats/tx.rpc.msg")).MsgClientImpl(rpc),
    vest: new (await import("./vest/tx.rpc.msg")).MsgClientImpl(rpc)
  }
});