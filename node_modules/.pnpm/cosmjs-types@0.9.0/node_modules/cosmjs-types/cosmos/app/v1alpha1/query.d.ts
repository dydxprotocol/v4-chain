import { Config } from "./config";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.app.v1alpha1";
/** QueryConfigRequest is the Query/Config request type. */
export interface QueryConfigRequest {
}
/** QueryConfigRequest is the Query/Config response type. */
export interface QueryConfigResponse {
    /** config is the current app config. */
    config?: Config;
}
export declare const QueryConfigRequest: {
    typeUrl: string;
    encode(_: QueryConfigRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConfigRequest;
    fromJSON(_: any): QueryConfigRequest;
    toJSON(_: QueryConfigRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): QueryConfigRequest;
};
export declare const QueryConfigResponse: {
    typeUrl: string;
    encode(message: QueryConfigResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConfigResponse;
    fromJSON(object: any): QueryConfigResponse;
    toJSON(message: QueryConfigResponse): unknown;
    fromPartial<I extends {
        config?: {
            modules?: {
                name?: string | undefined;
                config?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                golangBindings?: {
                    interfaceType?: string | undefined;
                    implementation?: string | undefined;
                }[] | undefined;
            }[] | undefined;
            golangBindings?: {
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            }[] | undefined;
        } | undefined;
    } & {
        config?: ({
            modules?: {
                name?: string | undefined;
                config?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                golangBindings?: {
                    interfaceType?: string | undefined;
                    implementation?: string | undefined;
                }[] | undefined;
            }[] | undefined;
            golangBindings?: {
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            }[] | undefined;
        } & {
            modules?: ({
                name?: string | undefined;
                config?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                golangBindings?: {
                    interfaceType?: string | undefined;
                    implementation?: string | undefined;
                }[] | undefined;
            }[] & ({
                name?: string | undefined;
                config?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                golangBindings?: {
                    interfaceType?: string | undefined;
                    implementation?: string | undefined;
                }[] | undefined;
            } & {
                name?: string | undefined;
                config?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["config"]["modules"][number]["config"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
                golangBindings?: ({
                    interfaceType?: string | undefined;
                    implementation?: string | undefined;
                }[] & ({
                    interfaceType?: string | undefined;
                    implementation?: string | undefined;
                } & {
                    interfaceType?: string | undefined;
                    implementation?: string | undefined;
                } & Record<Exclude<keyof I["config"]["modules"][number]["golangBindings"][number], keyof import("./config").GolangBinding>, never>)[] & Record<Exclude<keyof I["config"]["modules"][number]["golangBindings"], keyof {
                    interfaceType?: string | undefined;
                    implementation?: string | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["config"]["modules"][number], keyof import("./config").ModuleConfig>, never>)[] & Record<Exclude<keyof I["config"]["modules"], keyof {
                name?: string | undefined;
                config?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
                golangBindings?: {
                    interfaceType?: string | undefined;
                    implementation?: string | undefined;
                }[] | undefined;
            }[]>, never>) | undefined;
            golangBindings?: ({
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            }[] & ({
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            } & {
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            } & Record<Exclude<keyof I["config"]["golangBindings"][number], keyof import("./config").GolangBinding>, never>)[] & Record<Exclude<keyof I["config"]["golangBindings"], keyof {
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["config"], keyof Config>, never>) | undefined;
    } & Record<Exclude<keyof I, "config">, never>>(object: I): QueryConfigResponse;
};
/** Query is the app module query service. */
export interface Query {
    /** Config returns the current app config. */
    Config(request?: QueryConfigRequest): Promise<QueryConfigResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Config(request?: QueryConfigRequest): Promise<QueryConfigResponse>;
}
