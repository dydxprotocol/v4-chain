import {
  AwsSdkSigV4AuthInputConfig,
  AwsSdkSigV4AuthResolvedConfig,
  AwsSdkSigV4PreviouslyResolved,
} from "@aws-sdk/core";
import {
  HandlerExecutionContext,
  HttpAuthScheme,
  HttpAuthSchemeParameters,
  HttpAuthSchemeParametersProvider,
  HttpAuthSchemeProvider,
} from "@smithy/types";
import { LambdaClientResolvedConfig } from "../LambdaClient";
export interface LambdaHttpAuthSchemeParameters
  extends HttpAuthSchemeParameters {
  region?: string;
}
export interface LambdaHttpAuthSchemeParametersProvider
  extends HttpAuthSchemeParametersProvider<
    LambdaClientResolvedConfig,
    HandlerExecutionContext,
    LambdaHttpAuthSchemeParameters,
    object
  > {}
export declare const defaultLambdaHttpAuthSchemeParametersProvider: (
  config: LambdaClientResolvedConfig,
  context: HandlerExecutionContext,
  input: object
) => Promise<LambdaHttpAuthSchemeParameters>;
export interface LambdaHttpAuthSchemeProvider
  extends HttpAuthSchemeProvider<LambdaHttpAuthSchemeParameters> {}
export declare const defaultLambdaHttpAuthSchemeProvider: LambdaHttpAuthSchemeProvider;
export interface HttpAuthSchemeInputConfig extends AwsSdkSigV4AuthInputConfig {
  httpAuthSchemes?: HttpAuthScheme[];
  httpAuthSchemeProvider?: LambdaHttpAuthSchemeProvider;
}
export interface HttpAuthSchemeResolvedConfig
  extends AwsSdkSigV4AuthResolvedConfig {
  readonly httpAuthSchemes: HttpAuthScheme[];
  readonly httpAuthSchemeProvider: LambdaHttpAuthSchemeProvider;
}
export declare const resolveHttpAuthSchemeConfig: <T>(
  config: T & HttpAuthSchemeInputConfig & AwsSdkSigV4PreviouslyResolved
) => T & HttpAuthSchemeResolvedConfig;
