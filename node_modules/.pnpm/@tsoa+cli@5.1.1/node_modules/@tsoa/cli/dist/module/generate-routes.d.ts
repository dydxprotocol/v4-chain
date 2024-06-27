import * as ts from 'typescript';
import { ExtendedRoutesConfig } from '../cli';
import { Tsoa } from '@tsoa/runtime';
export declare function generateRoutes<Config extends ExtendedRoutesConfig>(routesConfig: Config, compilerOptions?: ts.CompilerOptions, ignorePaths?: string[], 
/**
 * pass in cached metadata returned in a previous step to speed things up
 */
metadata?: Tsoa.Metadata): Promise<Tsoa.Metadata>;
