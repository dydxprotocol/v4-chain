import { buildAllImports, getDepsFromQueries } from '../imports';
import { Bundler } from '../bundler';
import { getNestedProto } from '@osmonauts/proto-parser';
import { parse } from '../parse';
import { TelescopeBuilder } from '../builder';
import {
    createLCDClient,
} from '@osmonauts/ast';
import { ALLOWED_RPC_SERVICES } from '@osmonauts/types';

export const plugin = (
    builder: TelescopeBuilder,
    bundler: Bundler
) => {

    if (!builder.options.lcdClients.enabled) {
        return;
    }



    const queryContexts = bundler
        .contexts
        .filter(context =>
            context.queries.length > 0 ||
            context.services.length > 0
        );

    // [x] write out one registry helper for all contexts w/mutations
    const lcdClients = queryContexts.map(c => {

        const enabled = c.proto.pluginValue('lcdClients.enabled');
        if (!enabled) return;

        if (c.proto.isExcluded()) return;

        const ctx = bundler.getFreshContext(c);

        // get mutations, services
        parse(ctx);

        const proto = getNestedProto(c.ref.traversed);

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

        let getImportsFrom;

        // get imports
        allowedRpcServices.forEach(svcKey => {
            if (proto[svcKey]) {
                if (svcKey === 'Query') {
                    getImportsFrom = ctx.queries;
                } else {
                    getImportsFrom = ctx.services;
                }
            }
        });

        const localname = bundler.getLocalFilename(c.ref, 'lcd');
        const filename = bundler.getFilename(localname);

        let ast = null;

        allowedRpcServices.forEach(svcKey => {
            if (proto[svcKey]) {
                ast = createLCDClient(ctx.generic, proto[svcKey]);
            }
        });

        if (!ast) {
            return;
        }

        const serviceImports = getDepsFromQueries(
            getImportsFrom,
            localname
        );

        const imports = buildAllImports(ctx, serviceImports, localname);
        const prog = []
            .concat(imports)
            .concat(ctx.body)
            .concat(ast);

        bundler.writeAst(prog, filename);
        bundler.addToBundle(c, localname);

        return {
            // TODO use this to build LCD aggregators with scopes
            package: c.ref.proto.package,
            localname,
            filename
        };

    }).filter(Boolean);

    bundler.addLCDClients(lcdClients);

};