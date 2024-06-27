import { TelescopeInput } from './types';
import { TelescopeBuilder } from './builder';

export * from './builder';
export * from './bundler';
export * from './types';

export default async (input: TelescopeInput) => {
    const builder = new TelescopeBuilder(input);
    await builder.build();
};

