import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "cosmos.base.reflection.v1beta1";
/** ListAllInterfacesRequest is the request type of the ListAllInterfaces RPC. */
export interface ListAllInterfacesRequest {
}
/** ListAllInterfacesResponse is the response type of the ListAllInterfaces RPC. */
export interface ListAllInterfacesResponse {
    /** interface_names is an array of all the registered interfaces. */
    interfaceNames: string[];
}
/**
 * ListImplementationsRequest is the request type of the ListImplementations
 * RPC.
 */
export interface ListImplementationsRequest {
    /** interface_name defines the interface to query the implementations for. */
    interfaceName: string;
}
/**
 * ListImplementationsResponse is the response type of the ListImplementations
 * RPC.
 */
export interface ListImplementationsResponse {
    implementationMessageNames: string[];
}
export declare const ListAllInterfacesRequest: {
    typeUrl: string;
    encode(_: ListAllInterfacesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ListAllInterfacesRequest;
    fromJSON(_: any): ListAllInterfacesRequest;
    toJSON(_: ListAllInterfacesRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): ListAllInterfacesRequest;
};
export declare const ListAllInterfacesResponse: {
    typeUrl: string;
    encode(message: ListAllInterfacesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ListAllInterfacesResponse;
    fromJSON(object: any): ListAllInterfacesResponse;
    toJSON(message: ListAllInterfacesResponse): unknown;
    fromPartial<I extends {
        interfaceNames?: string[] | undefined;
    } & {
        interfaceNames?: (string[] & string[] & Record<Exclude<keyof I["interfaceNames"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "interfaceNames">, never>>(object: I): ListAllInterfacesResponse;
};
export declare const ListImplementationsRequest: {
    typeUrl: string;
    encode(message: ListImplementationsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ListImplementationsRequest;
    fromJSON(object: any): ListImplementationsRequest;
    toJSON(message: ListImplementationsRequest): unknown;
    fromPartial<I extends {
        interfaceName?: string | undefined;
    } & {
        interfaceName?: string | undefined;
    } & Record<Exclude<keyof I, "interfaceName">, never>>(object: I): ListImplementationsRequest;
};
export declare const ListImplementationsResponse: {
    typeUrl: string;
    encode(message: ListImplementationsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ListImplementationsResponse;
    fromJSON(object: any): ListImplementationsResponse;
    toJSON(message: ListImplementationsResponse): unknown;
    fromPartial<I extends {
        implementationMessageNames?: string[] | undefined;
    } & {
        implementationMessageNames?: (string[] & string[] & Record<Exclude<keyof I["implementationMessageNames"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "implementationMessageNames">, never>>(object: I): ListImplementationsResponse;
};
/** ReflectionService defines a service for interface reflection. */
export interface ReflectionService {
    /**
     * ListAllInterfaces lists all the interfaces registered in the interface
     * registry.
     */
    ListAllInterfaces(request?: ListAllInterfacesRequest): Promise<ListAllInterfacesResponse>;
    /**
     * ListImplementations list all the concrete types that implement a given
     * interface.
     */
    ListImplementations(request: ListImplementationsRequest): Promise<ListImplementationsResponse>;
}
export declare class ReflectionServiceClientImpl implements ReflectionService {
    private readonly rpc;
    constructor(rpc: Rpc);
    ListAllInterfaces(request?: ListAllInterfacesRequest): Promise<ListAllInterfacesResponse>;
    ListImplementations(request: ListImplementationsRequest): Promise<ListImplementationsResponse>;
}
