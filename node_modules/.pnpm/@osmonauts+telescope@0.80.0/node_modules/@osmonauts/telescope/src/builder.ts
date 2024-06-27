import { ProtoStore } from '@osmonauts/proto-parser';
import { TelescopeParseContext } from './build';
import { TelescopeOptions, defaultTelescopeOptions } from '@osmonauts/types';
import { bundlePackages } from './bundle';
import { BundlerFile, TelescopeInput } from './types';
import { Bundler } from './bundler';
import deepmerge from 'deepmerge';
import { resolve } from 'path';

import { plugin as createTypes } from './generators/create-types';
import { plugin as createAminoConverters } from './generators/create-amino-converters';
import { plugin as createRegistries } from './generators/create-registries';
import { plugin as createLCDClients } from './generators/create-lcd-clients';
import { plugin as createAggregatedLCDClient } from './generators/create-aggregated-lcd-client';
import { plugin as createLCDClientsScoped } from './generators/create-lcd-client-scoped';
import { plugin as createRPCQueryClientsScoped } from './generators/create-rpc-query-client-scoped';
import { plugin as createRPCMsgClientsScoped } from './generators/create-rpc-msg-client-scoped';
import { plugin as createRPCQueryClients } from './generators/create-rpc-query-clients';
import { plugin as createRPCMsgClients } from './generators/create-rpc-msg-clients';
import { plugin as createReactQueryBundle } from './generators/create-react-query-bundle';
import { plugin as createStargateClients } from './generators/create-stargate-clients';
import { plugin as createBundle } from './generators/create-bundle';
import { plugin as createIndex } from './generators/create-index';
import { plugin as createHelpers } from './generators/create-helpers';
import { plugin as createCosmWasmBundle } from './generators/create-cosmwasm-bundle';

const sanitizeOptions = (options: TelescopeOptions): TelescopeOptions => {
    // If an element at the same key is present for both x and y, the value from y will appear in the result.
    options = deepmerge(defaultTelescopeOptions, options ?? {});
    // strip off leading slashes
    options.tsDisable.files = options.tsDisable.files.map(file => file.startsWith('/') ? file : file.replace(/^\//, ''));
    options.eslintDisable.files = options.eslintDisable.files.map(file => file.startsWith('/') ? file : file.replace(/^\//, ''));
    // uniq bc of deepmerge
    options.rpcClients.enabledServices = [...new Set([...options.rpcClients.enabledServices])];
    return options;
};

export class TelescopeBuilder {
    store: ProtoStore;
    protoDirs: string[];
    outPath: string;
    options: TelescopeOptions;
    contexts: TelescopeParseContext[] = [];
    files: string[] = [];

    readonly converters: BundlerFile[] = [];
    readonly lcdClients: BundlerFile[] = [];
    readonly rpcQueryClients: BundlerFile[] = [];
    readonly rpcMsgClients: BundlerFile[] = [];
    readonly registries: BundlerFile[] = [];

    constructor({ protoDirs, outPath, store, options }: TelescopeInput & { store?: ProtoStore }) {
        this.protoDirs = protoDirs;
        this.outPath = resolve(outPath);
        this.options = sanitizeOptions(options);
        this.store = store ?? new ProtoStore(protoDirs, this.options);
        this.store.traverseAll();
    }

    context(ref) {
        const ctx = new TelescopeParseContext(
            ref, this.store, this.options
        );
        this.contexts.push(ctx);
        return ctx;
    }

    addRPCQueryClients(files: BundlerFile[]) {
        [].push.apply(this.rpcQueryClients, files);
    }

    addRPCMsgClients(files: BundlerFile[]) {
        [].push.apply(this.rpcMsgClients, files);
    }

    addLCDClients(files: BundlerFile[]) {
        [].push.apply(this.lcdClients, files);
    }

    addRegistries(files: BundlerFile[]) {
        [].push.apply(this.registries, files);
    }

    addConverters(files: BundlerFile[]) {
        [].push.apply(this.converters, files);
    }

    async build() {
        // [x] get bundle of all packages
        const bundles = bundlePackages(this.store)
            .map(bundle => {
                // store bundleFile in filesToInclude
                const bundler = new Bundler(this, bundle);

                // [x] write out all TS files for package
                createTypes(this, bundler);

                // [x] write out one amino helper for all contexts w/mutations
                createAminoConverters(this, bundler);

                // [x] write out one registry helper for all contexts w/mutations
                createRegistries(this, bundler);

                // [x] write out one registry helper for all contexts w/mutations
                createLCDClients(this, bundler);

                createRPCQueryClients(this, bundler);
                createRPCMsgClients(this, bundler);

                // [x] write out one client for each base package, referencing the last two steps
                createStargateClients(this, bundler);

                return bundler;
            });

        // post run plugins
        bundles
            .forEach(bundler => {
                createLCDClientsScoped(this, bundler);
                createRPCQueryClientsScoped(this, bundler);
                createRPCMsgClientsScoped(this, bundler);

                createBundle(this, bundler);
            });

        createReactQueryBundle(this);
        createAggregatedLCDClient(this);
        await createCosmWasmBundle(this);

        createHelpers(this);

        // finally, write one index file with all files, exported
        createIndex(this);
    }
}
