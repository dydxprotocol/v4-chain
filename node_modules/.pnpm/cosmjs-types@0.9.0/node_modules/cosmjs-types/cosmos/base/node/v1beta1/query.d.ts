import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "cosmos.base.node.v1beta1";
/** ConfigRequest defines the request structure for the Config gRPC query. */
export interface ConfigRequest {
}
/** ConfigResponse defines the response structure for the Config gRPC query. */
export interface ConfigResponse {
    minimumGasPrice: string;
}
export declare const ConfigRequest: {
    typeUrl: string;
    encode(_: ConfigRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ConfigRequest;
    fromJSON(_: any): ConfigRequest;
    toJSON(_: ConfigRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): ConfigRequest;
};
export declare const ConfigResponse: {
    typeUrl: string;
    encode(message: ConfigResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ConfigResponse;
    fromJSON(object: any): ConfigResponse;
    toJSON(message: ConfigResponse): unknown;
    fromPartial<I extends {
        minimumGasPrice?: string | undefined;
    } & {
        minimumGasPrice?: string | undefined;
    } & Record<Exclude<keyof I, "minimumGasPrice">, never>>(object: I): ConfigResponse;
};
/** Service defines the gRPC querier service for node related queries. */
export interface Service {
    /** Config queries for the operator configuration. */
    Config(request?: ConfigRequest): Promise<ConfigResponse>;
}
export declare class ServiceClientImpl implements Service {
    private readonly rpc;
    constructor(rpc: Rpc);
    Config(request?: ConfigRequest): Promise<ConfigResponse>;
}
