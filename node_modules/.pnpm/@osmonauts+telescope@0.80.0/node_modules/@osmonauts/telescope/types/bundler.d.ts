import * as t from '@babel/types';
import { TelescopeParseContext } from './build';
import { TelescopeBuilder } from './builder';
import { ProtoRef } from '@osmonauts/types';
import { Bundle, BundlerFile } from './types';
export declare class Bundler {
    builder: TelescopeBuilder;
    contexts: TelescopeParseContext[];
    bundle: Bundle;
    files: string[];
    readonly converters: BundlerFile[];
    readonly lcdClients: BundlerFile[];
    readonly rpcQueryClients: BundlerFile[];
    readonly rpcMsgClients: BundlerFile[];
    readonly registries: BundlerFile[];
    constructor(builder: TelescopeBuilder, bundle: Bundle);
    addLCDClients(files: BundlerFile[]): void;
    addRPCQueryClients(files: BundlerFile[]): void;
    addRPCMsgClients(files: BundlerFile[]): void;
    addRegistries(files: BundlerFile[]): void;
    addConverters(files: BundlerFile[]): void;
    getFreshContext(context: TelescopeParseContext): TelescopeParseContext;
    getLocalFilename(ref: ProtoRef, suffix?: string): string;
    getFilename(localname: string): string;
    writeAst(program: t.Statement[], filename: string): void;
    addToBundle(context: TelescopeParseContext, localname: string): void;
    addToBundleToPackage(packagename: string, localname: string): void;
}
