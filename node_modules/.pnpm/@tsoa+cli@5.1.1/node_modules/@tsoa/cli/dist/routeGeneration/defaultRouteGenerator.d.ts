import { ExtendedRoutesConfig } from '../cli';
import { Tsoa } from '@tsoa/runtime';
import { AbstractRouteGenerator } from './routeGenerator';
export declare class DefaultRouteGenerator extends AbstractRouteGenerator<ExtendedRoutesConfig> {
    pathTransformerFn: (path: string) => string;
    template: string;
    constructor(metadata: Tsoa.Metadata, options: ExtendedRoutesConfig);
    GenerateCustomRoutes(): Promise<void>;
    GenerateRoutes(middlewareTemplate: string): Promise<void>;
    protected pathTransformer(path: string): string;
    buildContent(middlewareTemplate: string): string;
}
