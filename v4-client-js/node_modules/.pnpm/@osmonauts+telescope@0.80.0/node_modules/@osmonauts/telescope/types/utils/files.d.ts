import * as t from '@babel/types';
import { TelescopeOptions } from '@osmonauts/types';
export declare const writeAstToFile: (outPath: string, options: TelescopeOptions, program: t.Statement[], filename: string) => void;
export declare const writeContentToFile: (outPath: string, options: TelescopeOptions, content: string, filename: string) => void;
