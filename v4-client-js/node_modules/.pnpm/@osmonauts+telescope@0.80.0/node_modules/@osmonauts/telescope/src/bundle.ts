import { ProtoStore } from '@osmonauts/proto-parser';
import { importNamespace } from '@osmonauts/ast';
import * as dotty from 'dotty';
import { TelescopeBuilder } from './index';
import { join, relative, dirname } from 'path';
import { TelescopeParseContext } from './build';
import { TelescopeOptions } from '@osmonauts/types';

// TODO move to store
export const getPackages = (store: ProtoStore) => {
    return store.getProtos().reduce((m, proto) => {
        if (proto.proto.package) {
            m[proto.proto.package] = m[proto.proto.package] || []
            m[proto.proto.package].push(proto);
        }
        return m;
    }, {});
}

export const getPackagesBundled = (store: ProtoStore) => {
    const objectified = {};
    const pkgs = getPackages(store);
    Object.keys(pkgs).forEach(key => {
        if (store.options.prototypes?.excluded?.packages?.includes?.(key)) return;
        const files = pkgs[key];
        dotty.put(objectified, key, {
            pkg: key,
            files: files.filter(file => {
                // TODO implement pattern
                const val = store.options.prototypes?.excluded?.protos?.includes?.(file.filename);
                if (typeof val === 'undefined') return true;
                return !val;
            })
        });
    });
    return objectified;

}

export const bundlePackages = (store: ProtoStore) => {
    const allPackages = getPackagesBundled(store);
    return Object.keys(allPackages).map(base => {
        const pkgs = allPackages[base];
        const bundleVariables = {};
        const bundleFile = join(base, 'bundle.ts');
        const importPaths = [];
        parsePackage(store.options, pkgs, bundleFile, importPaths, bundleVariables);
        return {
            bundleVariables,
            bundleFile,
            importPaths,
            base
        };
    });
};

// TODO review bundle registry methods 
export const bundleRegistries = (telescope: TelescopeBuilder) => {
    const withMutations = telescope.contexts.filter(
        ctx => ctx.mutations.length
    );
    const obj = withMutations.reduce((m, ctx) => {
        m[ctx.ref.proto.package] = m[ctx.ref.proto.package] ?? [];
        m[ctx.ref.proto.package].push(ctx);
        return m;
    }, {});

    return Object.entries(obj)
        .map(([pkg, serviceProtos]) => {
            return {
                package: pkg,
                contexts: serviceProtos
            };
        });
}

export const bundleBaseRegistries = (telescope: TelescopeBuilder) => {
    const withMutations = telescope.contexts.filter(
        ctx => ctx.mutations.length
    );
    const obj = withMutations.reduce((m, ctx) => {
        const base = ctx.ref.proto.package.split('.')[0];
        m[base] = m[base] ?? {};
        m[base][ctx.ref.proto.package] = m[base][ctx.ref.proto.package] ?? [];
        m[base][ctx.ref.proto.package].push(ctx);
        return m;
    }, {});

    return Object.entries(obj)
        .map(([pkg, withServices]) => {

            const serviceProtos = Object.entries(withServices)
                .map(([pkg, withServices]) => {
                    return {
                        package: pkg,
                        contexts: withServices
                    }
                });


            return {
                base: pkg,
                pkgs: serviceProtos
            };
        });
};

export const parseContextsForRegistry = (contexts: TelescopeParseContext[]) => {
    return contexts.map((ctx: TelescopeParseContext) => {
        const responses = ctx.mutations.map(m => m.response)
        const imports = ctx.mutations.reduce((m, msg) => {
            m[msg.messageImport] = m[msg.messageImport] ?? [];
            m[msg.messageImport].push(msg.message);
            return m;
        }, {})

        return {
            filename: ctx.ref.filename,
            imports,
            objects: ctx.types
                .filter(type => !type.isNested)
                .filter(type => !responses.includes(type.name))
                .map(type => type.name)
        }
    });
};

export const parsePackage = (
    options: TelescopeOptions,
    obj,
    bundleFile,
    importPaths,
    bundleVariables
) => {
    if (!obj) return;
    if (obj.pkg && obj.files) {
        obj.files.forEach(file => {
            createFileBundle(options, obj.pkg, file.filename, bundleFile, importPaths, bundleVariables);
        });
        return;
    }
    Object.keys(obj).forEach(mini => {
        parsePackage(options, obj[mini], bundleFile, importPaths, bundleVariables);
    })
}

let counter = 0;
export const createFileBundle = (
    options: TelescopeOptions,
    pkg,
    filename,
    bundleFile,
    importPaths,
    bundleVariables
) => {
    let rel = relative(dirname(bundleFile), filename);
    if (!rel.startsWith('.')) rel = `./${rel}`;
    const variable = `_${counter++}`;
    importPaths.push(importNamespace(variable, rel));
    dotty.put(bundleVariables, pkg + '.__export', true);
    dotty.put(bundleVariables, pkg + '.' + variable, true);
}