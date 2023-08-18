import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgSetVestingEntry, MsgSetVestingEntryResponse, MsgDeleteVestingEntry, MsgDeleteVestingEntryResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** SetVestingEntry sets a VestingEntry in state. */
  setVestingEntry(request: MsgSetVestingEntry): Promise<MsgSetVestingEntryResponse>;
  /** DeleteVestingEntry deletes a VestingEntry from state. */

  deleteVestingEntry(request: MsgDeleteVestingEntry): Promise<MsgDeleteVestingEntryResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.setVestingEntry = this.setVestingEntry.bind(this);
    this.deleteVestingEntry = this.deleteVestingEntry.bind(this);
  }

  setVestingEntry(request: MsgSetVestingEntry): Promise<MsgSetVestingEntryResponse> {
    const data = MsgSetVestingEntry.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vesting.Msg", "SetVestingEntry", data);
    return promise.then(data => MsgSetVestingEntryResponse.decode(new _m0.Reader(data)));
  }

  deleteVestingEntry(request: MsgDeleteVestingEntry): Promise<MsgDeleteVestingEntryResponse> {
    const data = MsgDeleteVestingEntry.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vesting.Msg", "DeleteVestingEntry", data);
    return promise.then(data => MsgDeleteVestingEntryResponse.decode(new _m0.Reader(data)));
  }

}