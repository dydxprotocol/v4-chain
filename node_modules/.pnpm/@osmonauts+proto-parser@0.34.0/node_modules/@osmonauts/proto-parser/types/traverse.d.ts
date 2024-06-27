import { ProtoRoot, ProtoRef } from '@osmonauts/types';
import { ProtoStore } from './store';
export declare const traverse: (store: ProtoStore, ref: ProtoRef) => ProtoRoot & {
    parsedImports: Record<string, string[]>;
    parsedExports: Record<string, any>;
    importNames: Record<string, any>;
};
export declare const recursiveTraversal: (store: ProtoStore, ref: ProtoRef, obj: any, imports: object, exports: object, traversal: string[], isNested: boolean) => any;
