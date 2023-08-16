import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgAcknowledgeBridges, MsgAcknowledgeBridgesResponse, MsgCompleteBridge, MsgCompleteBridgeResponse, MsgUpdateEventParams, MsgUpdateEventParamsResponse, MsgUpdateProposeParams, MsgUpdateProposeParamsResponse, MsgUpdateSafetyParams, MsgUpdateSafetyParamsResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /**
   * AcknowledgeBridges acknowledges bridges and sets them to complete at a
   * later block.
   */
  acknowledgeBridges(request: MsgAcknowledgeBridges): Promise<MsgAcknowledgeBridgesResponse>;
  /** CompleteBridge finalizes a bridge by minting coins to an address. */

  completeBridge(request: MsgCompleteBridge): Promise<MsgCompleteBridgeResponse>;
  /** UpdateEventParams updates the EventParams in state. */

  updateEventParams(request: MsgUpdateEventParams): Promise<MsgUpdateEventParamsResponse>;
  /** UpdateProposeParams updates the ProposeParams in state. */

  updateProposeParams(request: MsgUpdateProposeParams): Promise<MsgUpdateProposeParamsResponse>;
  /** UpdateSafetyParams updates the SafetyParams in state. */

  updateSafetyParams(request: MsgUpdateSafetyParams): Promise<MsgUpdateSafetyParamsResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.acknowledgeBridges = this.acknowledgeBridges.bind(this);
    this.completeBridge = this.completeBridge.bind(this);
    this.updateEventParams = this.updateEventParams.bind(this);
    this.updateProposeParams = this.updateProposeParams.bind(this);
    this.updateSafetyParams = this.updateSafetyParams.bind(this);
  }

  acknowledgeBridges(request: MsgAcknowledgeBridges): Promise<MsgAcknowledgeBridgesResponse> {
    const data = MsgAcknowledgeBridges.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.bridge.Msg", "AcknowledgeBridges", data);
    return promise.then(data => MsgAcknowledgeBridgesResponse.decode(new _m0.Reader(data)));
  }

  completeBridge(request: MsgCompleteBridge): Promise<MsgCompleteBridgeResponse> {
    const data = MsgCompleteBridge.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.bridge.Msg", "CompleteBridge", data);
    return promise.then(data => MsgCompleteBridgeResponse.decode(new _m0.Reader(data)));
  }

  updateEventParams(request: MsgUpdateEventParams): Promise<MsgUpdateEventParamsResponse> {
    const data = MsgUpdateEventParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.bridge.Msg", "UpdateEventParams", data);
    return promise.then(data => MsgUpdateEventParamsResponse.decode(new _m0.Reader(data)));
  }

  updateProposeParams(request: MsgUpdateProposeParams): Promise<MsgUpdateProposeParamsResponse> {
    const data = MsgUpdateProposeParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.bridge.Msg", "UpdateProposeParams", data);
    return promise.then(data => MsgUpdateProposeParamsResponse.decode(new _m0.Reader(data)));
  }

  updateSafetyParams(request: MsgUpdateSafetyParams): Promise<MsgUpdateSafetyParamsResponse> {
    const data = MsgUpdateSafetyParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.bridge.Msg", "UpdateSafetyParams", data);
    return promise.then(data => MsgUpdateSafetyParamsResponse.decode(new _m0.Reader(data)));
  }

}