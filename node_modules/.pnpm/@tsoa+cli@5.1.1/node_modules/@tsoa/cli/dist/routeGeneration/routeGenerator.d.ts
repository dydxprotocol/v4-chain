/// <reference types="node" />
import { ExtendedRoutesConfig } from '../cli';
import { Tsoa, TsoaRoute } from '@tsoa/runtime';
export declare abstract class AbstractRouteGenerator<Config extends ExtendedRoutesConfig> {
    protected readonly metadata: Tsoa.Metadata;
    protected readonly options: Config;
    constructor(metadata: Tsoa.Metadata, options: Config);
    /**
     * This is the entrypoint for a generator to create a custom set of routes
     */
    abstract GenerateCustomRoutes(): Promise<void>;
    buildModels(): TsoaRoute.Models;
    protected pathTransformer(path: string): string;
    protected buildContext(): {
        authenticationModule: string | undefined;
        basePath: string;
        canImportByAlias: boolean;
        controllers: {
            actions: {
                fullPath: string;
                method: string;
                name: string;
                parameters: {
                    [name: string]: TsoaRoute.ParameterSchema;
                };
                path: string;
                uploadFile: boolean;
                uploadFileName: string | undefined;
                uploadFiles: boolean;
                uploadFilesName: string | undefined;
                security: Tsoa.Security[];
                successStatus: string | number;
            }[];
            modulePath: string;
            name: string;
            path: string;
        }[];
        environment: NodeJS.ProcessEnv;
        iocModule: string | undefined;
        minimalSwaggerConfig: {
            noImplicitAdditionalProperties: "ignore" | "throw-on-extras" | "silently-remove-extras";
        };
        models: TsoaRoute.Models;
        useFileUploads: boolean;
        multerOpts: import("multer").Options | undefined;
        useSecurity: boolean;
        esm: boolean | undefined;
    };
    protected getRelativeImportPath(fileLocation: string): string;
    protected buildPropertySchema(source: Tsoa.Property): TsoaRoute.PropertySchema;
    protected buildParameterSchema(source: Tsoa.Parameter): TsoaRoute.ParameterSchema;
    protected buildProperty(type: Tsoa.Type): TsoaRoute.PropertySchema;
    protected shouldWriteFile(fileName: string, content: string): Promise<boolean>;
}
