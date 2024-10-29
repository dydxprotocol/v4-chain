import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgSetVestEntry, MsgSetVestEntryResponse, MsgDeleteVestEntry, MsgDeleteVestEntryResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** SetVestEntry sets a VestEntry in state. */
  setVestEntry(request: MsgSetVestEntry): Promise<MsgSetVestEntryResponse>;
  /** DeleteVestEntry deletes a VestEntry from state. */

  deleteVestEntry(request: MsgDeleteVestEntry): Promise<MsgDeleteVestEntryResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.setVestEntry = this.setVestEntry.bind(this);
    this.deleteVestEntry = this.deleteVestEntry.bind(this);
  }

  setVestEntry(request: MsgSetVestEntry): Promise<MsgSetVestEntryResponse> {
    const data = MsgSetVestEntry.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vest.Msg", "SetVestEntry", data);
    return promise.then(data => MsgSetVestEntryResponse.decode(new _m0.Reader(data)));
  }

  deleteVestEntry(request: MsgDeleteVestEntry): Promise<MsgDeleteVestEntryResponse> {
    const data = MsgDeleteVestEntry.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vest.Msg", "DeleteVestEntry", data);
    return promise.then(data => MsgDeleteVestEntryResponse.decode(new _m0.Reader(data)));
  }

}