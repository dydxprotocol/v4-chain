import { TSBuilder, TSBuilderInput } from './builder';

export { default as generateTypes } from './generators/types';
export { default as generateClient } from './generators/client';
export { default as generateMessageComposer } from './generators/message-composer';
export { default as generateReactQuery } from './generators/react-query';
export { default as generateRecoil } from './generators/recoil';

export * from './utils';
export * from './builder';
export * from './bundler';

export default async (input: TSBuilderInput) => {
    const builder = new TSBuilder(input);
    await builder.build();
};