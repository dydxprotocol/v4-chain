import { TSBuilderInput } from '@cosmwasm/ts-codegen';
import { AminoExceptions } from "./aminos";
export declare enum TelescopeLogLevel {
    None = 0,
    Info = 1,
    Warn = 2,
    Error = 3,
    Debug = 4
}
interface TelescopeOpts {
    experimentalGlobalProtoNamespace?: boolean;
    removeUnusedImports?: boolean;
    classesUseArrowFunctions?: boolean;
    useSDKTypes?: boolean;
    includeExternalHelpers?: boolean;
    logLevel?: TelescopeLogLevel;
    prototypes?: {
        enabled?: boolean;
        parser?: {
            keepCase?: boolean;
            alternateCommentMode?: boolean;
            preferTrailingComment?: boolean;
        };
        methods?: {
            encode?: boolean;
            decode?: boolean;
            fromJSON?: boolean;
            toJSON?: boolean;
            fromPartial?: boolean;
            toSDK?: boolean;
            fromSDK?: boolean;
        };
        includePackageVar?: boolean;
        fieldDefaultIsOptional?: boolean;
        useOptionalNullable?: boolean;
        allowUndefinedTypes?: boolean;
        optionalQueryParams?: boolean;
        optionalPageRequests?: boolean;
        excluded?: {
            packages?: string[];
            protos?: string[];
        };
        typingsFormat?: {
            useDeepPartial?: boolean;
            useExact?: boolean;
            timestamp?: 'date' | 'timestamp';
            duration?: 'duration' | 'string';
        };
    };
    tsDisable?: {
        files?: string[];
        disableAll?: boolean;
        patterns?: string[];
    };
    eslintDisable?: {
        files?: string[];
        disableAll?: boolean;
        patterns?: string[];
    };
    bundle?: {
        enabled: boolean;
    };
    cosmwasm?: TSBuilderInput;
    aggregatedLCD?: {
        dir: string;
        filename: string;
        packages: string[];
        addToBundle: boolean;
    };
    stargateClients?: {
        enabled: boolean;
        includeCosmosDefaultTypes?: boolean;
    };
    aminoEncoding?: {
        enabled: boolean;
        casingFn?: Function;
        exceptions?: AminoExceptions;
        typeUrlToAmino?: (typeUrl: string) => string | undefined;
    };
    lcdClients?: {
        enabled: boolean;
        scopedIsExclusive?: boolean;
        bundle?: boolean;
        scoped?: {
            dir: string;
            filename?: string;
            packages: string[];
            addToBundle: boolean;
            methodName?: string;
        }[];
    };
    rpcClients?: {
        enabled: boolean;
        inline?: boolean;
        extensions?: boolean;
        camelCase?: boolean;
        scopedIsExclusive?: boolean;
        bundle?: boolean;
        enabledServices?: ('Msg' | 'Query' | 'Service' | 'ReflectionService' | 'ABCIApplication' | string)[];
        scoped?: {
            dir: string;
            filename?: string;
            packages: string[];
            addToBundle: boolean;
            methodNameQuery?: string;
            methodNameTx?: string;
        }[];
    };
    reactQuery?: {
        enabled: boolean;
    };
}
interface TelescopePackageOpts {
    packages?: Record<string, any>;
}
export declare type TelescopeOptions = TelescopeOpts & TelescopePackageOpts;
export declare type TelescopeOption = keyof TelescopeOpts;
export declare const defaultTelescopeOptions: TelescopeOptions;
export {};
