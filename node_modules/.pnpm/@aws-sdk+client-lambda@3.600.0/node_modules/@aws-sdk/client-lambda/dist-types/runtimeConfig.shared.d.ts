import { LambdaClientConfig } from "./LambdaClient";
/**
 * @internal
 */
export declare const getRuntimeConfig: (config: LambdaClientConfig) => {
    apiVersion: string;
    base64Decoder: import("@smithy/types").Decoder;
    base64Encoder: (_input: string | Uint8Array) => string;
    disableHostPrefix: boolean;
    endpointProvider: (endpointParams: import("./endpoint/EndpointParameters").EndpointParameters, context?: {
        logger?: import("@smithy/types").Logger | undefined;
    }) => import("@smithy/types").EndpointV2;
    extensions: import("./runtimeExtensions").RuntimeExtension[];
    httpAuthSchemeProvider: import("./auth/httpAuthSchemeProvider").LambdaHttpAuthSchemeProvider;
    httpAuthSchemes: import("@smithy/types").HttpAuthScheme[];
    logger: import("@smithy/types").Logger;
    serviceId: string;
    urlParser: import("@smithy/types").UrlParser;
    utf8Decoder: import("@smithy/types").Decoder;
    utf8Encoder: (input: string | Uint8Array) => string;
};
