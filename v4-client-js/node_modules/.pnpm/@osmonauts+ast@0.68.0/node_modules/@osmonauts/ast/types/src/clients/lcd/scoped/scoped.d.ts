import * as t from '@babel/types';
import { GenericParseContext } from '../../../encoding';
export declare const lcdArguments: () => t.ObjectProperty[];
export declare const lcdFuncArguments: () => t.ObjectPattern[];
export declare const lcdClassArguments: () => t.ObjectExpression[];
export declare const lcdNewAwaitImport: (path: string, className: string, _arguments: t.ObjectExpression[]) => t.NewExpression;
export declare const lcdRecursiveObjectProps: (names: string[], leaf?: any) => t.ObjectExpression;
export declare const lcdNestedImportObject: (obj: object, className: string, _arguments: t.ObjectExpression[]) => any;
export declare const createScopedLCDFactory: (context: GenericParseContext, obj: object, identifier: string, className: string) => t.ExportNamedDeclaration;
