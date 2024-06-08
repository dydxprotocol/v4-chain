import { ProtoRef } from '@osmonauts/types';
import { TelescopeParseContext } from './build';
export declare const parse: (context: TelescopeParseContext) => void;
export declare const getParsedObjectName: (ref: ProtoRef, obj: any, scope: string[]) => string;
export declare const parseType: (context: TelescopeParseContext, obj: any, scope: string[], isNested?: boolean) => void;
export declare const parseEnum: (context: TelescopeParseContext, obj: any, scope: string[], isNested?: boolean) => void;
export declare const parseService: (context: TelescopeParseContext, obj: any, scope: string[], isNested?: boolean) => void;
interface ParseRecur {
    context: TelescopeParseContext;
    obj: any;
    scope: string[];
    isNested: boolean;
}
export declare const parseRecur: ({ context, obj, scope, isNested }: ParseRecur) => void;
export {};
