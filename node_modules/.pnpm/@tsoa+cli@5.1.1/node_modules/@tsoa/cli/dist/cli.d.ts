#!/usr/bin/env node
import { Config, RoutesConfig, SpecConfig, Tsoa } from '@tsoa/runtime';
import { AbstractRouteGenerator } from './routeGeneration/routeGenerator';
export interface ExtendedSpecConfig extends SpecConfig {
    entryFile: Config['entryFile'];
    noImplicitAdditionalProperties: Exclude<Config['noImplicitAdditionalProperties'], undefined>;
    controllerPathGlobs?: Config['controllerPathGlobs'];
}
export declare const validateSpecConfig: (config: Config) => Promise<ExtendedSpecConfig>;
type RouteGeneratorImpl = new (metadata: Tsoa.Metadata, options: ExtendedRoutesConfig) => AbstractRouteGenerator<any>;
export interface ExtendedRoutesConfig extends RoutesConfig {
    entryFile: Config['entryFile'];
    noImplicitAdditionalProperties: Exclude<Config['noImplicitAdditionalProperties'], undefined>;
    controllerPathGlobs?: Config['controllerPathGlobs'];
    multerOpts?: Config['multerOpts'];
    rootSecurity?: Config['spec']['rootSecurity'];
    routeGenerator?: string | RouteGeneratorImpl;
}
export interface ConfigArgs {
    basePath?: string;
    configuration?: string | Config;
}
export interface SwaggerArgs extends ConfigArgs {
    host?: string;
    json?: boolean;
    yaml?: boolean;
}
export declare function runCLI(): void;
export declare function generateSpecAndRoutes(args: SwaggerArgs, metadata?: Tsoa.Metadata): Promise<Tsoa.Metadata>;
export {};
