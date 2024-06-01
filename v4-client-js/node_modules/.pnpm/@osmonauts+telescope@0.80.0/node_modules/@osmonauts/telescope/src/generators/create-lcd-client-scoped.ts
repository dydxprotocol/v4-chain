import * as dotty from 'dotty';
import { getNestedProto } from '@osmonauts/proto-parser';
import { join } from 'path';
import { TelescopeBuilder } from '../builder';
import { createScopedLCDFactory } from '@osmonauts/ast';
import { ALLOWED_RPC_SERVICES, ProtoRef } from '@osmonauts/types';
import { fixlocalpaths, getRelativePath } from '../utils';
import { Bundler } from '../bundler';
import { TelescopeParseContext } from '../build';
import { aggregateImports, getImportStatements } from '../imports';

export const plugin = (
    builder: TelescopeBuilder,
    bundler: Bundler
) => {

    // if not enabled, exit
    if (!builder.options?.lcdClients?.enabled) {
        return;
    }

    // if no scopes, do them all!
    if (
        !builder.options.lcdClients.scoped ||
        !builder.options.lcdClients.scoped.length
    ) {
        // TODO inefficient
        // WE SHOULD NOT DO THIS IN A BUNDLER LOOP
        // MAKE SEPARATE PLUGIN
        return createAllLCDBundles(
            builder,
            bundler
        );
    }

    if (!builder.options.lcdClients.scopedIsExclusive) {
        // TODO inefficient
        // WE SHOULD NOT DO THIS IN A BUNDLER LOOP
        // MAKE SEPARATE PLUGIN
        createAllLCDBundles(
            builder,
            bundler
        );
    }

    // we have scopes!
    builder.options.lcdClients.scoped.forEach(lcd => {
        if (lcd.dir !== bundler.bundle.base) return;
        makeLCD(
            builder,
            bundler,
            lcd
        );
    });
};

const getFileName = (dir, filename) => {
    const localname = join(dir, filename ?? 'lcd.ts');
    if (localname.endsWith('.ts')) return localname;
    return localname + '.ts';
};

const makeLCD = (
    builder: TelescopeBuilder,
    bundler: Bundler,
    lcd: {
        dir: string;
        filename?: string;
        packages: string[];
        addToBundle: boolean;
        methodName?: string;
    }
) => {
    const dir = lcd.dir;
    const packages = lcd.packages;
    const methodName = lcd.methodName ?? 'createLCDClient'
    const localname = getFileName(dir, lcd.filename);

    const obj = {};
    builder.lcdClients.forEach(file => {

        // ADD all option
        // which defaults to including cosmos 
        // and defaults to base for each
        if (!packages.includes(file.package)) {
            return;
        }

        const f = localname;
        const f2 = file.localname;
        const importPath = getRelativePath(f, f2);
        dotty.put(obj, file.package, importPath);
    });


    const ctx = new TelescopeParseContext(
        {
            absolute: '',
            filename: localname,
            proto: {
                package: dir,
                imports: null,
                root: {},
                importNames: null
            },
            traversed: {
                package: dir,
                imports: null,
                root: {},
                importNames: null
            }
        },
        builder.store,
        builder.options
    );

    const lcdast = createScopedLCDFactory(
        ctx.proto,
        obj,
        methodName,
        'LCDQueryClient' // make option later
    );

    const imports = aggregateImports(ctx, {}, localname);

    const importStmts = getImportStatements(
        localname,
        [...fixlocalpaths(imports)]
    );

    const prog = []
        .concat(importStmts)
        .concat(lcdast);

    const filename = bundler.getFilename(localname);
    bundler.writeAst(prog, filename);

    if (lcd.addToBundle) {
        bundler.addToBundleToPackage(`${dir}.ClientFactory`, localname)
    }
};

// TODO
/*
 move all options for lcd into previous `lcd` prop and 
 clean up all these many options for one nested object full of options
*/

const createAllLCDBundles = (
    builder: TelescopeBuilder,
    bundler: Bundler
) => {

    if (!builder.options.lcdClients.bundle) return;


    // [x] loop through every bundle 
    // [x] if not cosmos, add all cosmos
    // [x] call makeLCD
    // [x] add to bundle

    const dir = bundler.bundle.base;
    const filename = 'lcd.ts'

    ///
    ///
    ///

    // refs with services
    const refs = builder.store.getProtos().filter((ref: ProtoRef) => {
        const proto = getNestedProto(ref.traversed);
        //// Anything except Msg Service OK...
        const allowedRpcServices = builder.options.rpcClients.enabledServices.filter(a => a !== 'Msg');
        const found = allowedRpcServices.some(svc => {
            return proto?.[svc] &&
                proto[svc]?.type === 'Service'
        });

        if (!found) {
            return;
        }
        ///


        return true;
    });

    const check = refs.filter((ref: ProtoRef) => {
        const [base] = ref.proto.package.split('.');
        return base === bundler.bundle.base;
    });

    if (!check.length) {
        // if there are no services
        // exit the plugin
        return;
    }

    const packages = refs.reduce((m, ref: ProtoRef) => {
        const [base] = ref.proto.package.split('.');
        if (base === 'cosmos' || base === bundler.bundle.base)
            return [...new Set([...m, ref.proto.package])];
        return m;
    }, []);

    makeLCD(
        builder,
        bundler,
        {
            dir,
            filename,
            packages,
            addToBundle: true,
            methodName: 'createLCDClient'
        }
    );

};