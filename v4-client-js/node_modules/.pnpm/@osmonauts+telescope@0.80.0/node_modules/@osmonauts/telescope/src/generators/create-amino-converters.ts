import { buildAllImports, getDepsFromMutations } from '../imports';
import { Bundler } from '../bundler';
import { parse } from '../parse';
import { TelescopeBuilder } from '../builder';

export const plugin = (
    builder: TelescopeBuilder,
    bundler: Bundler
) => {

    const aminoEncoding = builder.options.aminoEncoding;
    if (!aminoEncoding.enabled) {
        return;
    }

    const mutationContexts = bundler
        .contexts
        .filter(context => context.mutations.length > 0);

    const converters = mutationContexts.map(c => {

        const aminoEncodingEnabled = c.amino.pluginValue('aminoEncoding.enabled');
        if (!aminoEncodingEnabled) {
            return;
        }

        if (c.proto.isExcluded()) {
            return;
        }

        const localname = bundler.getLocalFilename(c.ref, 'amino');
        const filename = bundler.getFilename(localname);
        const ctx = bundler.getFreshContext(c);

        // get mutations, services
        parse(ctx);

        // now let's amino...
        ctx.buildAminoInterfaces();
        ctx.buildAminoConverter();

        const serviceImports = getDepsFromMutations(
            ctx.mutations,
            localname
        );

        // build file
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

    bundler.addConverters(converters);
};