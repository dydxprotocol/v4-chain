import { buildAllImports, getDepsFromQueries } from '../imports';
import { Bundler } from '../bundler';
import {
    createRpcQueryExtension,
    createRpcClientClass,
    createRpcClientInterface,
    createRpcQueryHookInterfaces,
    createRpcQueryHookClientMap,
    createRpcQueryHooks
} from '@osmonauts/ast';
import { getNestedProto } from '@osmonauts/proto-parser';
import { parse } from '../parse';
import { TelescopeBuilder } from '../builder';
import { ProtoService } from '@osmonauts/types';

export const plugin = (
    builder: TelescopeBuilder,
    bundler: Bundler
) => {

    const clients = bundler.contexts.map(c => {

        const enabled = c.proto.pluginValue('rpcClients.enabled');
        if (!enabled) return;

        const inline = c.proto.pluginValue('rpcClients.inline');
        if (inline) return;

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

        let name, getImportsFrom;

        allowedRpcServices.forEach(svcKey => {
            if (proto[svcKey]) {
                if (svcKey === 'Query') {
                    getImportsFrom = ctx.queries;
                } else {
                    getImportsFrom = ctx.services;
                }
                name = svcKey;
            }
        });

        const localname = bundler.getLocalFilename(c.ref, `rpc.${name}`);
        const filename = bundler.getFilename(localname);

        const asts = [];

        allowedRpcServices.forEach(svcKey => {
            if (proto[svcKey]) {

                const svc: ProtoService = proto[svcKey];

                asts.push(createRpcClientInterface(ctx.generic, svc));
                asts.push(createRpcClientClass(ctx.generic, svc));
                if (c.proto.pluginValue('rpcClients.extensions')) {
                    asts.push(createRpcQueryExtension(ctx.generic, svc));
                }

                // react query
                // TODO use the imports and make separate files
                if (c.proto.pluginValue('reactQuery.enabled')) {
                    [].push.apply(asts, createRpcQueryHookInterfaces(ctx.generic, svc));
                    [].push.apply(asts, createRpcQueryHookClientMap(ctx.generic, svc));
                    asts.push(createRpcQueryHooks(ctx.generic, proto[svcKey]));
                }
            }
        });

        if (!asts.length) {
            return;
        }

        const serviceImports = getDepsFromQueries(
            getImportsFrom,
            localname
        );

        // TODO we do NOT need all imports...
        const imports = buildAllImports(ctx, serviceImports, localname);
        const prog = []
            .concat(imports)
            .concat(ctx.body)
            .concat(asts);

        bundler.writeAst(prog, filename);
        bundler.addToBundle(c, localname);

        return {
            package: c.ref.proto.package,
            localname,
            filename
        };

    }).filter(Boolean);

    bundler.addRPCQueryClients(clients);
};