import { aggregateImports, getDepsFromQueries, getImportStatements } from '../imports';
import { getNestedProto } from '@osmonauts/proto-parser';
import { parse } from '../parse';
import { join } from 'path';
import { TelescopeBuilder } from '../builder';
import { createAggregatedLCDClient } from '@osmonauts/ast';
import { ALLOWED_RPC_SERVICES, ProtoRef, ProtoService } from '@osmonauts/types';
import { TelescopeParseContext } from '../build';
import { writeAstToFile } from '../utils/files';
import { fixlocalpaths } from '../utils';

const isExcluded = (builder: TelescopeBuilder, ref: ProtoRef) => {
    return builder.options.prototypes?.excluded?.protos?.includes(ref.filename) ||
        builder.options.prototypes?.excluded?.packages?.includes(ref.proto.package);
};

export const plugin = (
    builder: TelescopeBuilder
) => {

    if (!builder.options.aggregatedLCD) {
        return;
    }

    const opts = builder.options.aggregatedLCD;

    const {
        dir,
        filename: fname,
        packages
    } = opts;

    const localname = join(dir, fname);

    const refs = builder.store.filterProtoWhere((ref: ProtoRef) => {
        return packages.includes(ref.proto.package) &&
            !isExcluded(builder, ref);
    });

    const services: ProtoService[] = refs.map(ref => {
        const proto = getNestedProto(ref.traversed);
        if (!proto?.Query || proto.Query?.type !== 'Service') {
            return;
        }
        return proto.Query;
    }).filter(Boolean);

    const tc = new TelescopeParseContext(refs[0], builder.store, builder.options);
    const context = tc.proto;
    const lcdast = createAggregatedLCDClient(context, services, 'QueryClient');

    const importsForAggregator = aggregateImports(tc, {}, localname);

    /////////
    /////////
    /////////
    /////////

    const queryContexts = builder
        .contexts
        .filter(context =>
            context.queries.length > 0 ||
            context.services.length > 0
        );

    const progImports = queryContexts.reduce((m, c) => {

        if (!builder.options.aggregatedLCD.packages.includes(c.ref.proto.package)) {
            return m;
        }

        const ctx = new TelescopeParseContext(
            c.ref,
            c.store,
            builder.options
        );

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

        allowedRpcServices.forEach(svcKey => {
            if (proto[svcKey]) {
                if (svcKey === 'Query') {
                    getImportsFrom = ctx.queries;
                } else {
                    getImportsFrom = ctx.services;
                }
            }
        });


        const serviceImports = getDepsFromQueries(
            getImportsFrom,
            localname
        );

        const imports = aggregateImports(ctx, serviceImports, localname);

        return [...m, ...fixlocalpaths(imports)]
    }, []);

    const importStmts = getImportStatements(
        localname,
        [...importsForAggregator, ...progImports]
    );

    const prog = []
        .concat(importStmts)
        .concat(lcdast);

    const filename = join(builder.outPath, localname);
    writeAstToFile(builder.outPath, builder.options, prog, filename);

};