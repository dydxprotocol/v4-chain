import { Rpc } from "../../helpers";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
/** Query defines the gRPC querier service. */

export interface Query {}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {};
};