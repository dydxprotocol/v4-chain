import * as t from '@babel/types';
import { GenericParseContext } from '../../encoding';
export declare const rpcHookFuncArguments: () => t.ObjectPattern[];
export declare const rpcHookClassArguments: () => t.ObjectExpression[];
export declare const rpcHookNewTmRequire: (imports: HookImport[], path: string, methodName: string) => t.CallExpression;
export declare const rpcHookRecursiveObjectProps: (names: string[], leaf?: any) => t.ObjectExpression;
export declare const rpcHookTmNestedImportObject: (imports: HookImport[], obj: object, methodName: string) => any;
interface HookImport {
    as: string;
    path: string;
}
export declare const createScopedRpcHookFactory: (context: GenericParseContext, obj: object, identifier: string) => (t.ExportNamedDeclaration | {
    type: string;
    importKind: string;
    specifiers: {
        type: string;
        local: {
            type: string;
            name: string;
        };
    }[];
    source: {
        type: string;
        value: string;
    };
})[];
export {};
