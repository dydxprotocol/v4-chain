import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgProposedOperations, MsgProposedOperationsResponse, MsgPlaceOrder, MsgPlaceOrderResponse, MsgCancelOrder, MsgCancelOrderResponse, MsgCreateClobPair, MsgCreateClobPairResponse, MsgSetClobPairStatus, MsgSetClobPairStatusResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /**
   * ProposedOperations is a temporary message used by block proposers
   * for matching orders as part of the ABCI++ workaround.
   */
  proposedOperations(request: MsgProposedOperations): Promise<MsgProposedOperationsResponse>;
  /** PlaceOrder allows accounts to place orders on the orderbook. */

  placeOrder(request: MsgPlaceOrder): Promise<MsgPlaceOrderResponse>;
  /** CancelOrder allows accounts to cancel existing orders on the orderbook. */

  cancelOrder(request: MsgCancelOrder): Promise<MsgCancelOrderResponse>;
  /** CreateClobPair creates a new clob pair. */

  createClobPair(request: MsgCreateClobPair): Promise<MsgCreateClobPairResponse>;
  /** SetClobPairStatus sets the status of a clob pair. */

  setClobPairStatus(request: MsgSetClobPairStatus): Promise<MsgSetClobPairStatusResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.proposedOperations = this.proposedOperations.bind(this);
    this.placeOrder = this.placeOrder.bind(this);
    this.cancelOrder = this.cancelOrder.bind(this);
    this.createClobPair = this.createClobPair.bind(this);
    this.setClobPairStatus = this.setClobPairStatus.bind(this);
  }

  proposedOperations(request: MsgProposedOperations): Promise<MsgProposedOperationsResponse> {
    const data = MsgProposedOperations.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Msg", "ProposedOperations", data);
    return promise.then(data => MsgProposedOperationsResponse.decode(new _m0.Reader(data)));
  }

  placeOrder(request: MsgPlaceOrder): Promise<MsgPlaceOrderResponse> {
    const data = MsgPlaceOrder.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Msg", "PlaceOrder", data);
    return promise.then(data => MsgPlaceOrderResponse.decode(new _m0.Reader(data)));
  }

  cancelOrder(request: MsgCancelOrder): Promise<MsgCancelOrderResponse> {
    const data = MsgCancelOrder.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Msg", "CancelOrder", data);
    return promise.then(data => MsgCancelOrderResponse.decode(new _m0.Reader(data)));
  }

  createClobPair(request: MsgCreateClobPair): Promise<MsgCreateClobPairResponse> {
    const data = MsgCreateClobPair.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Msg", "CreateClobPair", data);
    return promise.then(data => MsgCreateClobPairResponse.decode(new _m0.Reader(data)));
  }

  setClobPairStatus(request: MsgSetClobPairStatus): Promise<MsgSetClobPairStatusResponse> {
    const data = MsgSetClobPairStatus.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Msg", "SetClobPairStatus", data);
    return promise.then(data => MsgSetClobPairStatusResponse.decode(new _m0.Reader(data)));
  }

}