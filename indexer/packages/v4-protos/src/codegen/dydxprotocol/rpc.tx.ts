import { Rpc } from "../helpers";
export const createRPCMsgClient = async ({
  rpc
}: {
  rpc: Rpc;
}) => ({
  dydxprotocol: {
    accountplus: new (await import("./accountplus/tx.rpc.msg")).MsgClientImpl(rpc),
    affiliates: new (await import("./affiliates/tx.rpc.msg")).MsgClientImpl(rpc),
    blocktime: new (await import("./blocktime/tx.rpc.msg")).MsgClientImpl(rpc),
    bridge: new (await import("./bridge/tx.rpc.msg")).MsgClientImpl(rpc),
    clob: new (await import("./clob/tx.rpc.msg")).MsgClientImpl(rpc),
    delaymsg: new (await import("./delaymsg/tx.rpc.msg")).MsgClientImpl(rpc),
    feetiers: new (await import("./feetiers/tx.rpc.msg")).MsgClientImpl(rpc),
    govplus: new (await import("./govplus/tx.rpc.msg")).MsgClientImpl(rpc),
    listing: new (await import("./listing/tx.rpc.msg")).MsgClientImpl(rpc),
    perpetuals: new (await import("./perpetuals/tx.rpc.msg")).MsgClientImpl(rpc),
    prices: new (await import("./prices/tx.rpc.msg")).MsgClientImpl(rpc),
    ratelimit: new (await import("./ratelimit/tx.rpc.msg")).MsgClientImpl(rpc),
    revshare: new (await import("./revshare/tx.rpc.msg")).MsgClientImpl(rpc),
    rewards: new (await import("./rewards/tx.rpc.msg")).MsgClientImpl(rpc),
    sending: new (await import("./sending/tx.rpc.msg")).MsgClientImpl(rpc),
    stats: new (await import("./stats/tx.rpc.msg")).MsgClientImpl(rpc),
    vault: new (await import("./vault/tx.rpc.msg")).MsgClientImpl(rpc),
    vest: new (await import("./vest/tx.rpc.msg")).MsgClientImpl(rpc)
  }
});