import { TelescopeOptions } from '@osmonauts/types';
export interface Bundle {
    bundleVariables: {};
    bundleFile: string;
    importPaths: any[];
    base: string;
}
export interface BundlerFile {
    package?: string;
    localname: string;
    filename: string;
}

export interface ImportObj {
    type: 'import' | 'default' | 'namespace' | string;
    name: string;
    path: string;
    importAs?: string;
}
export interface ImportHash {
    [key: string]: string[];
}

export interface TelescopeInput {
    protoDirs: string[];
    outPath: string;
    options: TelescopeOptions;
}
