import * as ts from 'typescript';
import { ExtendedSpecConfig } from '../cli';
import { Tsoa } from '@tsoa/runtime';
export declare const getSwaggerOutputPath: (swaggerConfig: ExtendedSpecConfig) => string;
export declare const generateSpec: (swaggerConfig: ExtendedSpecConfig, compilerOptions?: ts.CompilerOptions, ignorePaths?: string[], metadata?: Tsoa.Metadata) => Promise<Tsoa.Metadata>;
