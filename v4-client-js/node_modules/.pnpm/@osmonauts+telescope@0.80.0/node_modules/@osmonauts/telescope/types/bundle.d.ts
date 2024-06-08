import { ProtoStore } from '@osmonauts/proto-parser';
import { TelescopeBuilder } from './index';
import { TelescopeParseContext } from './build';
import { TelescopeOptions } from '@osmonauts/types';
export declare const getPackages: (store: ProtoStore) => {};
export declare const getPackagesBundled: (store: ProtoStore) => {};
export declare const bundlePackages: (store: ProtoStore) => {
    bundleVariables: {};
    bundleFile: string;
    importPaths: any[];
    base: string;
}[];
export declare const bundleRegistries: (telescope: TelescopeBuilder) => {
    package: string;
    contexts: unknown;
}[];
export declare const bundleBaseRegistries: (telescope: TelescopeBuilder) => {
    base: string;
    pkgs: {
        package: string;
        contexts: any;
    }[];
}[];
export declare const parseContextsForRegistry: (contexts: TelescopeParseContext[]) => {
    filename: string;
    imports: {};
    objects: any[];
}[];
export declare const parsePackage: (options: TelescopeOptions, obj: any, bundleFile: any, importPaths: any, bundleVariables: any) => void;
export declare const createFileBundle: (options: TelescopeOptions, pkg: any, filename: any, bundleFile: any, importPaths: any, bundleVariables: any) => void;
