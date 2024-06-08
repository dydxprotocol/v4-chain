import { TelescopeBuilder } from '../builder';
import { TSBuilder } from '@cosmwasm/ts-codegen';

export const plugin = async (
    builder: TelescopeBuilder
) => {

    if (!builder.options.cosmwasm) {
        return;
    }

    const input = builder.options.cosmwasm;
    const cosmWasmBuilder = new TSBuilder(input);
    await cosmWasmBuilder.build();
    const file = input.options.bundle.bundleFile;
    builder.files.push(file);
};