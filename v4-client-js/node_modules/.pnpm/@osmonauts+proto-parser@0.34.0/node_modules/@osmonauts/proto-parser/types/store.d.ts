import { ProtoDep, ProtoRef, ProtoServiceMethod, TelescopeOptions } from '@osmonauts/types';
interface ParseProtoOptions {
    keepCase?: boolean;
    alternateCommentMode?: boolean;
    preferTrailingComment?: boolean;
}
export declare const parseProto: (content: any, options?: ParseProtoOptions) => import("@pyramation/protobufjs").IParserResult;
export declare class ProtoStore {
    files: string[];
    protoDirs: string[];
    deps: ProtoDep[];
    protos: ProtoRef[];
    packages: string[];
    options: TelescopeOptions;
    requests: Record<string, ProtoServiceMethod>;
    responses: Record<string, ProtoServiceMethod>;
    _traversed: boolean;
    constructor(protoDirs?: string[], options?: TelescopeOptions);
    findProto(filename: any): ProtoRef;
    findProtoWhere(fn: (ref: ProtoRef) => boolean): ProtoRef;
    filterProtoWhere(fn: (ref: ProtoRef) => boolean): ProtoRef[];
    findProtoObject(filename: any, name: any): any;
    registerRequest(svc: ProtoServiceMethod): void;
    processProtos(contents: {
        absolute: string;
        filename: string;
        content: string;
    }[]): {
        absolute: string;
        filename: string;
        proto: import("@pyramation/protobufjs").IParserResult;
    }[];
    getProtos(): ProtoRef[];
    getPackages(): string[];
    parseScope(type: string): {
        nested: string;
        package: string;
    };
    getDeps(): ProtoDep[];
    getDependencies(protos: ProtoRef[]): ProtoDep[];
    traverseAll(): void;
    get(from: ProtoRef, name: string): import("./lookup").Lookup;
    getImportFromRef(ref: ProtoRef, name: string): import("./lookup").Lookup;
    getServices(myBase: string): Record<string, ProtoRef[]>;
}
export {};
