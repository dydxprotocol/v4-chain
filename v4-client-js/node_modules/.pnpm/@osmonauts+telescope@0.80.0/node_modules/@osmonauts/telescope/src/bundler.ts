import * as t from '@babel/types';
import { resolve, join } from 'path';
import { TelescopeParseContext } from './build';
import { createFileBundle } from './bundle';
import { TelescopeBuilder } from './builder';
import { ProtoRef } from '@osmonauts/types';
import { Bundle, BundlerFile } from './types';
import { writeAstToFile } from './utils/files';

export class Bundler {
    builder: TelescopeBuilder;

    contexts: TelescopeParseContext[] = [];
    bundle: Bundle;
    files: string[];

    readonly converters: BundlerFile[] = [];
    readonly lcdClients: BundlerFile[] = [];
    readonly rpcQueryClients: BundlerFile[] = [];
    readonly rpcMsgClients: BundlerFile[] = [];
    readonly registries: BundlerFile[] = [];

    constructor(
        builder: TelescopeBuilder,
        bundle: Bundle
    ) {
        this.builder = builder;
        this.bundle = bundle;
        this.files = [
            bundle.bundleFile
        ]
    }

    addLCDClients(files: BundlerFile[]) {
        [].push.apply(this.lcdClients, files);
        this.builder.addLCDClients(files);
    }

    addRPCQueryClients(files: BundlerFile[]) {
        [].push.apply(this.rpcQueryClients, files);
        this.builder.addRPCQueryClients(files);
    }

    addRPCMsgClients(files: BundlerFile[]) {
        [].push.apply(this.rpcMsgClients, files);
        this.builder.addRPCMsgClients(files);
    }

    addRegistries(files: BundlerFile[]) {
        [].push.apply(this.registries, files);
        this.builder.addRegistries(files);
    }

    addConverters(files: BundlerFile[]) {
        [].push.apply(this.converters, files);
        this.builder.addConverters(files);
    }

    getFreshContext(
        context: TelescopeParseContext
    ) {
        return new TelescopeParseContext(
            context.ref,
            context.store,
            this.builder.options
        );
    }

    getLocalFilename(
        ref: ProtoRef,
        suffix?: string
    ) {
        return suffix ?
            ref.filename.replace(/\.proto$/, `.${suffix}.ts`) :
            ref.filename.replace(/\.proto$/, `.ts`);
    }

    getFilename(
        localname: string
    ) {
        return resolve(join(this.builder.outPath, localname));
    }

    writeAst(
        program: t.Statement[],
        filename: string
    ) {
        writeAstToFile(this.builder.outPath, this.builder.options, program, filename);
    }

    // addToBundle adds the path into the namespaced bundle object
    addToBundle(
        context: TelescopeParseContext,
        localname: string
    ) {
        createFileBundle(
            this.builder.options,
            context.ref.proto.package,
            localname,
            this.bundle.bundleFile,
            this.bundle.importPaths,
            this.bundle.bundleVariables
        );
    }

    addToBundleToPackage(
        packagename: string,
        localname: string
    ) {
        createFileBundle(
            this.builder.options,
            packagename,
            localname,
            this.bundle.bundleFile,
            this.bundle.importPaths,
            this.bundle.bundleVariables
        );
    }

}