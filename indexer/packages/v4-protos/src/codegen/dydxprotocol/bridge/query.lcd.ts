import { LCDClient } from "@osmonauts/lcd";
import { QueryEventParamsRequest, QueryEventParamsResponseSDKType, QueryProposeParamsRequest, QueryProposeParamsResponseSDKType, QuerySafetyParamsRequest, QuerySafetyParamsResponseSDKType, QueryAcknowledgedEventInfoRequest, QueryAcknowledgedEventInfoResponseSDKType, QueryRecognizedEventInfoRequest, QueryRecognizedEventInfoResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.eventParams = this.eventParams.bind(this);
    this.proposeParams = this.proposeParams.bind(this);
    this.safetyParams = this.safetyParams.bind(this);
    this.acknowledgedEventInfo = this.acknowledgedEventInfo.bind(this);
    this.recognizedEventInfo = this.recognizedEventInfo.bind(this);
  }
  /* Queries the EventParams. */


  async eventParams(_params: QueryEventParamsRequest = {}): Promise<QueryEventParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/bridge/event_params`;
    return await this.req.get<QueryEventParamsResponseSDKType>(endpoint);
  }
  /* Queries the ProposeParams. */


  async proposeParams(_params: QueryProposeParamsRequest = {}): Promise<QueryProposeParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/bridge/propose_params`;
    return await this.req.get<QueryProposeParamsResponseSDKType>(endpoint);
  }
  /* Queries the SafetyParams. */


  async safetyParams(_params: QuerySafetyParamsRequest = {}): Promise<QuerySafetyParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/bridge/safety_params`;
    return await this.req.get<QuerySafetyParamsResponseSDKType>(endpoint);
  }
  /* Queries the AcknowledgedEventInfo.
   An "acknowledged" event is one that is in-consensus and has been stored
   in-state. */


  async acknowledgedEventInfo(_params: QueryAcknowledgedEventInfoRequest = {}): Promise<QueryAcknowledgedEventInfoResponseSDKType> {
    const endpoint = `dydxprotocol/v4/bridge/acknowledged_event_info`;
    return await this.req.get<QueryAcknowledgedEventInfoResponseSDKType>(endpoint);
  }
  /* Queries the RecognizedEventInfo.
   A "recognized" event is one that is finalized on the Ethereum blockchain
   and has been identified by the queried node. It is not yet in-consensus. */


  async recognizedEventInfo(_params: QueryRecognizedEventInfoRequest = {}): Promise<QueryRecognizedEventInfoResponseSDKType> {
    const endpoint = `dydxprotocol/v4/bridge/recognized_event_info`;
    return await this.req.get<QueryRecognizedEventInfoResponseSDKType>(endpoint);
  }

}