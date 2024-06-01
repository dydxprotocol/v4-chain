import { buildAllImports, getDepsFromMutations } from '../imports';
import { Bundler } from '../bundler';
import { parse } from '../parse';
import { TelescopeBuilder } from '../builder';

export const plugin = (
    builder: TelescopeBuilder,
    bundler: Bundler
) => {

    if (!builder.options.aminoEncoding.enabled) {
        return;
    }

    const mutationContexts = bundler
        .contexts
        .filter(context => context.mutations.length > 0);

    // [x] write out one registry helper for all contexts w/mutations
    const registries = mutationContexts.map(c => {

        const enabled = c.amino.pluginValue('aminoEncoding.enabled');
        if (!enabled) return;

        if (c.proto.isExcluded()) return;

        const localname = bundler.getLocalFilename(c.ref, 'registry');
        const filename = bundler.getFilename(localname);
        const ctx = bundler.getFreshContext(c);

        // get mutations, services
        parse(ctx);

        ctx.buildRegistry();
        ctx.buildRegistryLoader();
        ctx.buildHelperObject();

        const serviceImports = getDepsFromMutations(
            ctx.mutations,
            localname
        );

        const imports = buildAllImports(ctx, serviceImports, localname);
        const prog = []
            .concat(imports)
            .concat(ctx.body);

        bundler.writeAst(prog, filename);
        bundler.addToBundle(c, localname);

        return {
            package: c.ref.proto.package,
            localname,
            filename
        };

    }).filter(Boolean);

    bundler.addRegistries(registries);

};