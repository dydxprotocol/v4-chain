import { JsonRpcRequest, JsonRpcSuccessResponse } from "@cosmjs/json-rpc";
import { RpcClient } from "./rpcclient";
export interface HttpEndpoint {
    /**
     * The URL of the HTTP endpoint.
     *
     * For POST APIs like Tendermint RPC in CosmJS,
     * this is without the method specific paths (e.g. https://cosmoshub-4--rpc--full.datahub.figment.io/)
     */
    readonly url: string;
    /**
     * HTTP headers that are sent with every request, such as authorization information.
     */
    readonly headers: Record<string, string>;
}
export declare class HttpClient implements RpcClient {
    protected readonly url: string;
    protected readonly headers: Record<string, string> | undefined;
    constructor(endpoint: string | HttpEndpoint);
    disconnect(): void;
    execute(request: JsonRpcRequest): Promise<JsonRpcSuccessResponse>;
}
