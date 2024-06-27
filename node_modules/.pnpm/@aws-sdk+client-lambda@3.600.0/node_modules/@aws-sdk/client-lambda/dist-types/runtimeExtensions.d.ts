import { LambdaExtensionConfiguration } from "./extensionConfiguration";
/**
 * @public
 */
export interface RuntimeExtension {
    configure(extensionConfiguration: LambdaExtensionConfiguration): void;
}
/**
 * @public
 */
export interface RuntimeExtensionsConfig {
    extensions: RuntimeExtension[];
}
/**
 * @internal
 */
export declare const resolveRuntimeExtensions: (runtimeConfig: any, extensions: RuntimeExtension[]) => any;
