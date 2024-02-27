import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgProposedOperations, MsgProposedOperationsResponse, MsgPlaceOrder, MsgPlaceOrderResponse, MsgCancelOrder, MsgCancelOrderResponse, MsgBatchCancel, MsgBatchCancelResponse, MsgCreateClobPair, MsgCreateClobPairResponse, MsgUpdateClobPair, MsgUpdateClobPairResponse, MsgUpdateEquityTierLimitConfiguration, MsgUpdateEquityTierLimitConfigurationResponse, MsgUpdateBlockRateLimitConfiguration, MsgUpdateBlockRateLimitConfigurationResponse, MsgUpdateLiquidationsConfig, MsgUpdateLiquidationsConfigResponse } from "./tx";
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
  /** BatchCancel allows accounts to cancel a batch of orders on the orderbook. */

  batchCancel(request: MsgBatchCancel): Promise<MsgBatchCancelResponse>;
  /** CreateClobPair creates a new clob pair. */

  createClobPair(request: MsgCreateClobPair): Promise<MsgCreateClobPairResponse>;
  /**
   * UpdateClobPair sets the status of a clob pair. Should return an error
   * if the authority is not in the clob keeper's set of authorities,
   * if the ClobPair id is not found in state, or if the update includes
   * an unsupported status transition.
   */

  updateClobPair(request: MsgUpdateClobPair): Promise<MsgUpdateClobPairResponse>;
  /**
   * UpdateEquityTierLimitConfiguration updates the equity tier limit
   * configuration in state.
   */

  updateEquityTierLimitConfiguration(request: MsgUpdateEquityTierLimitConfiguration): Promise<MsgUpdateEquityTierLimitConfigurationResponse>;
  /**
   * UpdateBlockRateLimitConfiguration updates the block rate limit
   * configuration in state.
   */

  updateBlockRateLimitConfiguration(request: MsgUpdateBlockRateLimitConfiguration): Promise<MsgUpdateBlockRateLimitConfigurationResponse>;
  /** UpdateLiquidationsConfig updates the liquidations configuration in state. */

  updateLiquidationsConfig(request: MsgUpdateLiquidationsConfig): Promise<MsgUpdateLiquidationsConfigResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.proposedOperations = this.proposedOperations.bind(this);
    this.placeOrder = this.placeOrder.bind(this);
    this.cancelOrder = this.cancelOrder.bind(this);
    this.batchCancel = this.batchCancel.bind(this);
    this.createClobPair = this.createClobPair.bind(this);
    this.updateClobPair = this.updateClobPair.bind(this);
    this.updateEquityTierLimitConfiguration = this.updateEquityTierLimitConfiguration.bind(this);
    this.updateBlockRateLimitConfiguration = this.updateBlockRateLimitConfiguration.bind(this);
    this.updateLiquidationsConfig = this.updateLiquidationsConfig.bind(this);
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

  batchCancel(request: MsgBatchCancel): Promise<MsgBatchCancelResponse> {
    const data = MsgBatchCancel.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Msg", "BatchCancel", data);
    return promise.then(data => MsgBatchCancelResponse.decode(new _m0.Reader(data)));
  }

  createClobPair(request: MsgCreateClobPair): Promise<MsgCreateClobPairResponse> {
    const data = MsgCreateClobPair.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Msg", "CreateClobPair", data);
    return promise.then(data => MsgCreateClobPairResponse.decode(new _m0.Reader(data)));
  }

  updateClobPair(request: MsgUpdateClobPair): Promise<MsgUpdateClobPairResponse> {
    const data = MsgUpdateClobPair.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Msg", "UpdateClobPair", data);
    return promise.then(data => MsgUpdateClobPairResponse.decode(new _m0.Reader(data)));
  }

  updateEquityTierLimitConfiguration(request: MsgUpdateEquityTierLimitConfiguration): Promise<MsgUpdateEquityTierLimitConfigurationResponse> {
    const data = MsgUpdateEquityTierLimitConfiguration.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Msg", "UpdateEquityTierLimitConfiguration", data);
    return promise.then(data => MsgUpdateEquityTierLimitConfigurationResponse.decode(new _m0.Reader(data)));
  }

  updateBlockRateLimitConfiguration(request: MsgUpdateBlockRateLimitConfiguration): Promise<MsgUpdateBlockRateLimitConfigurationResponse> {
    const data = MsgUpdateBlockRateLimitConfiguration.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Msg", "UpdateBlockRateLimitConfiguration", data);
    return promise.then(data => MsgUpdateBlockRateLimitConfigurationResponse.decode(new _m0.Reader(data)));
  }

  updateLiquidationsConfig(request: MsgUpdateLiquidationsConfig): Promise<MsgUpdateLiquidationsConfigResponse> {
    const data = MsgUpdateLiquidationsConfig.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Msg", "UpdateLiquidationsConfig", data);
    return promise.then(data => MsgUpdateLiquidationsConfigResponse.decode(new _m0.Reader(data)));
  }

}