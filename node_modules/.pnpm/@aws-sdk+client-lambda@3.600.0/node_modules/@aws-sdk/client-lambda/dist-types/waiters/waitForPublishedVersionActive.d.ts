import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { GetFunctionConfigurationCommandInput } from "../commands/GetFunctionConfigurationCommand";
import { LambdaClient } from "../LambdaClient";
/**
 * Waits for the published version's State to be Active. This waiter uses GetFunctionConfiguration API. This should be used after new version is published.
 *  @deprecated Use waitUntilPublishedVersionActive instead. waitForPublishedVersionActive does not throw error in non-success cases.
 */
export declare const waitForPublishedVersionActive: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionConfigurationCommandInput) => Promise<WaiterResult>;
/**
 * Waits for the published version's State to be Active. This waiter uses GetFunctionConfiguration API. This should be used after new version is published.
 *  @param params - Waiter configuration options.
 *  @param input - The input to GetFunctionConfigurationCommand for polling.
 */
export declare const waitUntilPublishedVersionActive: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionConfigurationCommandInput) => Promise<WaiterResult>;
