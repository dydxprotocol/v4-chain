import { aggregateImports, getImportStatements } from '../imports';
import { join } from 'path';
import { TelescopeBuilder } from '../builder';
import { createScopedRpcHookFactory } from '@osmonauts/ast';
import { ProtoRef } from '@osmonauts/types';
import { TelescopeParseContext } from '../build';
import { writeAstToFile } from '../utils/files';
import { fixlocalpaths } from '../utils';
import * as dotty from 'dotty';

export const plugin = (
    builder: TelescopeBuilder
) => {

    if (!builder.options.reactQuery.enabled) {
        return;
    }

    const localname = 'hooks.ts';

    const obj = {};
    builder.rpcQueryClients.map(queryClient => {
        const path = `./${queryClient.localname.replace(/\.ts$/, '')}`;
        dotty.put(obj, queryClient.package, path);
    });

    const pkg = '@root';
    const ref: ProtoRef = {
        absolute: '',
        filename: localname,
        proto: {
            package: pkg,
            imports: null,
            root: {},
            importNames: null
        },
        traversed: {
            package: pkg,
            imports: null,
            root: {},
            importNames: null
        }
    }

    const pCtx = new TelescopeParseContext(
        ref,
        builder.store,
        builder.options
    );

    const ast = createScopedRpcHookFactory(
        pCtx.proto,
        obj,
        'createRpcQueryHooks'
    )

    const imports = fixlocalpaths(aggregateImports(pCtx, {}, localname));
    const importStmts = getImportStatements(
        localname,
        imports
    );

    const prog = []
        .concat(importStmts)
        .concat(ast);

    const filename = join(builder.outPath, localname);
    builder.files.push(localname);

    writeAstToFile(builder.outPath, builder.options, prog, filename);

};