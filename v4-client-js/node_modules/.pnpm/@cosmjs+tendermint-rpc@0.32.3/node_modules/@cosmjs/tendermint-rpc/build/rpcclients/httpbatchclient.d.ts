import { JsonRpcRequest, JsonRpcSuccessResponse } from "@cosmjs/json-rpc";
import { HttpEndpoint } from "./httpclient";
import { RpcClient } from "./rpcclient";
export interface HttpBatchClientOptions {
    /** Interval for dispatching batches (in milliseconds) */
    dispatchInterval: number;
    /** Max number of items sent in one request */
    batchSizeLimit: number;
}
export declare class HttpBatchClient implements RpcClient {
    protected readonly url: string;
    protected readonly headers: Record<string, string> | undefined;
    protected readonly options: HttpBatchClientOptions;
    private timer?;
    private readonly queue;
    constructor(endpoint: string | HttpEndpoint, options?: Partial<HttpBatchClientOptions>);
    disconnect(): void;
    execute(request: JsonRpcRequest): Promise<JsonRpcSuccessResponse>;
    private validate;
    /**
     * This is called in an interval where promise rejections cannot be handled.
     * So this is not async and HTTP errors need to be handled by the queued promises.
     */
    private tick;
}
