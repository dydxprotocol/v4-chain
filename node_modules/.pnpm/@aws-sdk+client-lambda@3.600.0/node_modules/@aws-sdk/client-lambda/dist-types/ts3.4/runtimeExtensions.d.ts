import { LambdaExtensionConfiguration } from "./extensionConfiguration";
export interface RuntimeExtension {
  configure(extensionConfiguration: LambdaExtensionConfiguration): void;
}
export interface RuntimeExtensionsConfig {
  extensions: RuntimeExtension[];
}
export declare const resolveRuntimeExtensions: (
  runtimeConfig: any,
  extensions: RuntimeExtension[]
) => any;
