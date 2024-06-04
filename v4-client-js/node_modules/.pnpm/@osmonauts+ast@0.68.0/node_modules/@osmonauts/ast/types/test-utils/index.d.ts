import { ProtoType, TelescopeOptions } from '@osmonauts/types';
import { ProtoStore } from '@osmonauts/proto-parser';
import { AminoParseContext } from '../src/encoding/context';
import { GenericParseContext } from '../src/encoding';
export declare const expectCode: (ast: any) => void;
export declare const printCode: (ast: any) => void;
export declare const defaultTelescopeOptions: {
    experimentalGlobalProtoNamespace: boolean;
    removeUnusedImports: boolean;
    classesUseArrowFunctions: boolean;
    useSDKTypes: boolean;
    logLevel: import("@osmonauts/types").TelescopeLogLevel;
    prototypes: {
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
            timestamp?: "date" | "timestamp";
            duration?: "string" | "duration";
        };
    };
    tsDisable: {
        files?: string[];
        disableAll?: boolean;
        patterns?: string[];
    };
    eslintDisable: {
        files?: string[];
        disableAll?: boolean;
        patterns?: string[];
    };
    bundle: {
        enabled: boolean;
    };
    cosmwasm: import("@cosmwasm/ts-codegen").TSBuilderInput;
    aggregatedLCD: {
        dir: string;
        filename: string;
        packages: string[];
        addToBundle: boolean;
    };
    stargateClients: {
        enabled: boolean;
        includeCosmosDefaultTypes?: boolean;
    };
    aminoEncoding: {
        enabled: boolean;
        casingFn?: Function;
        exceptions?: import("@osmonauts/types").AminoExceptions;
        typeUrlToAmino?: (typeUrl: string) => string;
    };
    lcdClients: {
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
    rpcClients: {
        enabled: boolean;
        inline?: boolean;
        extensions?: boolean;
        camelCase?: boolean;
        scopedIsExclusive?: boolean;
        bundle?: boolean;
        enabledServices?: string[];
        scoped?: {
            dir: string;
            filename?: string;
            packages: string[];
            addToBundle: boolean;
            methodNameQuery?: string;
            methodNameTx?: string;
        }[];
    };
    reactQuery: {
        enabled: boolean;
    };
    packages: Record<string, any>;
} & {
    prototypes: {
        parser: {
            keepCase: boolean;
        };
        methods: {
            encode: boolean;
            decode: boolean;
            fromJSON: boolean;
            toJSON: boolean;
            fromPartial: boolean;
            toSDK: boolean;
            fromSDK: boolean;
        };
    };
};
export declare const getTestProtoStore: (options?: TelescopeOptions) => ProtoStore;
export declare const prepareContext: (store: ProtoStore, protoFile: string) => {
    context: AminoParseContext;
    root: import("@osmonauts/types").ProtoRoot;
    protos: ProtoType[];
};
export declare const getGenericParseContext: () => GenericParseContext;
export declare const getGenericParseContextWithRef: (ref: any) => GenericParseContext;
