import { Rpc } from "../../../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { GetRequest, GetResponse, ListRequest, ListResponse } from "./query";
/** Query is a generic gRPC service for querying ORM data. */

export interface Query {
  /** Get queries an ORM table against an unique index. */
  get(request: GetRequest): Promise<GetResponse>;
  /** List queries an ORM table against an index. */

  list(request: ListRequest): Promise<ListResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.get = this.get.bind(this);
    this.list = this.list.bind(this);
  }

  get(request: GetRequest): Promise<GetResponse> {
    const data = GetRequest.encode(request).finish();
    const promise = this.rpc.request("cosmos.orm.query.v1alpha1.Query", "Get", data);
    return promise.then(data => GetResponse.decode(new _m0.Reader(data)));
  }

  list(request: ListRequest): Promise<ListResponse> {
    const data = ListRequest.encode(request).finish();
    const promise = this.rpc.request("cosmos.orm.query.v1alpha1.Query", "List", data);
    return promise.then(data => ListResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    get(request: GetRequest): Promise<GetResponse> {
      return queryService.get(request);
    },

    list(request: ListRequest): Promise<ListResponse> {
      return queryService.list(request);
    }

  };
};