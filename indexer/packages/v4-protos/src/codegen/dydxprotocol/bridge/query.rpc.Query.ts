import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryEventParamsRequest, QueryEventParamsResponse, QueryProposeParamsRequest, QueryProposeParamsResponse, QuerySafetyParamsRequest, QuerySafetyParamsResponse, QueryAcknowledgedEventInfoRequest, QueryAcknowledgedEventInfoResponse, QueryRecognizedEventInfoRequest, QueryRecognizedEventInfoResponse, QueryInFlightCompleteBridgeMessagesRequest, QueryInFlightCompleteBridgeMessagesResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the EventParams. */
  eventParams(request?: QueryEventParamsRequest): Promise<QueryEventParamsResponse>;
  /** Queries the ProposeParams. */

  proposeParams(request?: QueryProposeParamsRequest): Promise<QueryProposeParamsResponse>;
  /** Queries the SafetyParams. */

  safetyParams(request?: QuerySafetyParamsRequest): Promise<QuerySafetyParamsResponse>;
  /**
   * Queries the AcknowledgedEventInfo.
   * An "acknowledged" event is one that is in-consensus and has been stored
   * in-state.
   */

  acknowledgedEventInfo(request?: QueryAcknowledgedEventInfoRequest): Promise<QueryAcknowledgedEventInfoResponse>;
  /**
   * Queries the RecognizedEventInfo.
   * A "recognized" event is one that is finalized on the Ethereum blockchain
   * and has been identified by the queried node. It is not yet in-consensus.
   */

  recognizedEventInfo(request?: QueryRecognizedEventInfoRequest): Promise<QueryRecognizedEventInfoResponse>;
  /**
   * Queries all `MsgCompleteBridge` messages that are in-flight (delayed
   * but not yet executed) and corresponding block heights at which they
   * will execute.
   */

  inFlightCompleteBridgeMessages(request: QueryInFlightCompleteBridgeMessagesRequest): Promise<QueryInFlightCompleteBridgeMessagesResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.eventParams = this.eventParams.bind(this);
    this.proposeParams = this.proposeParams.bind(this);
    this.safetyParams = this.safetyParams.bind(this);
    this.acknowledgedEventInfo = this.acknowledgedEventInfo.bind(this);
    this.recognizedEventInfo = this.recognizedEventInfo.bind(this);
    this.inFlightCompleteBridgeMessages = this.inFlightCompleteBridgeMessages.bind(this);
  }

  eventParams(request: QueryEventParamsRequest = {}): Promise<QueryEventParamsResponse> {
    const data = QueryEventParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.bridge.Query", "EventParams", data);
    return promise.then(data => QueryEventParamsResponse.decode(new _m0.Reader(data)));
  }

  proposeParams(request: QueryProposeParamsRequest = {}): Promise<QueryProposeParamsResponse> {
    const data = QueryProposeParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.bridge.Query", "ProposeParams", data);
    return promise.then(data => QueryProposeParamsResponse.decode(new _m0.Reader(data)));
  }

  safetyParams(request: QuerySafetyParamsRequest = {}): Promise<QuerySafetyParamsResponse> {
    const data = QuerySafetyParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.bridge.Query", "SafetyParams", data);
    return promise.then(data => QuerySafetyParamsResponse.decode(new _m0.Reader(data)));
  }

  acknowledgedEventInfo(request: QueryAcknowledgedEventInfoRequest = {}): Promise<QueryAcknowledgedEventInfoResponse> {
    const data = QueryAcknowledgedEventInfoRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.bridge.Query", "AcknowledgedEventInfo", data);
    return promise.then(data => QueryAcknowledgedEventInfoResponse.decode(new _m0.Reader(data)));
  }

  recognizedEventInfo(request: QueryRecognizedEventInfoRequest = {}): Promise<QueryRecognizedEventInfoResponse> {
    const data = QueryRecognizedEventInfoRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.bridge.Query", "RecognizedEventInfo", data);
    return promise.then(data => QueryRecognizedEventInfoResponse.decode(new _m0.Reader(data)));
  }

  inFlightCompleteBridgeMessages(request: QueryInFlightCompleteBridgeMessagesRequest): Promise<QueryInFlightCompleteBridgeMessagesResponse> {
    const data = QueryInFlightCompleteBridgeMessagesRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.bridge.Query", "InFlightCompleteBridgeMessages", data);
    return promise.then(data => QueryInFlightCompleteBridgeMessagesResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    eventParams(request?: QueryEventParamsRequest): Promise<QueryEventParamsResponse> {
      return queryService.eventParams(request);
    },

    proposeParams(request?: QueryProposeParamsRequest): Promise<QueryProposeParamsResponse> {
      return queryService.proposeParams(request);
    },

    safetyParams(request?: QuerySafetyParamsRequest): Promise<QuerySafetyParamsResponse> {
      return queryService.safetyParams(request);
    },

    acknowledgedEventInfo(request?: QueryAcknowledgedEventInfoRequest): Promise<QueryAcknowledgedEventInfoResponse> {
      return queryService.acknowledgedEventInfo(request);
    },

    recognizedEventInfo(request?: QueryRecognizedEventInfoRequest): Promise<QueryRecognizedEventInfoResponse> {
      return queryService.recognizedEventInfo(request);
    },

    inFlightCompleteBridgeMessages(request: QueryInFlightCompleteBridgeMessagesRequest): Promise<QueryInFlightCompleteBridgeMessagesResponse> {
      return queryService.inFlightCompleteBridgeMessages(request);
    }

  };
};