import * as t from '@babel/types';
import { TelescopeOptions } from '@osmonauts/types';
export declare const recursiveModuleBundle: (options: TelescopeOptions, obj: any) => any;
export declare const importNamespace: (ident: string, path: string) => t.ImportDeclaration;
