import { ProtoRoot } from '@osmonauts/types';
export declare const getNestedProto: (root: ProtoRoot) => any;
export declare const getNestedProtoGeneric: (root: ProtoRoot, path: string[]) => any;
export declare const getNested: (root: ProtoRoot, path: string[]) => any;
export declare const lookupSymbolScopes: (name: string, relativeTo: string, list?: string[]) => any;
export declare const getPackageAndNestedFromStr: (type: string, pkg: string) => {
    nested: string;
    package: string;
};
export declare const getServices: (root: ProtoRoot) => any[];
export declare const getTypes: (root: ProtoRoot) => any[];
export declare const getEnums: (root: ProtoRoot) => any[];
export declare const getObjectName: (name: string, scope?: string[]) => string;
export declare const SCALAR_TYPES: string[];
export declare const instanceType: (obj: any) => {
    type: string;
    name?: undefined;
} | {
    name: any;
    type: string;
};
