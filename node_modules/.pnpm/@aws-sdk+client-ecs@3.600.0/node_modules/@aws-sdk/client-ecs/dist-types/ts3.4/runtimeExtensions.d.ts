import { ECSExtensionConfiguration } from "./extensionConfiguration";
export interface RuntimeExtension {
  configure(extensionConfiguration: ECSExtensionConfiguration): void;
}
export interface RuntimeExtensionsConfig {
  extensions: RuntimeExtension[];
}
export declare const resolveRuntimeExtensions: (
  runtimeConfig: any,
  extensions: RuntimeExtension[]
) => any;
