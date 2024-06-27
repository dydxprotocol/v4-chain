import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { GetFunctionConfigurationCommandInput } from "../commands/GetFunctionConfigurationCommand";
import { LambdaClient } from "../LambdaClient";
/**
 * Waits for the function's State to be Active. This waiter uses GetFunctionConfiguration API. This should be used after new function creation.
 *  @deprecated Use waitUntilFunctionActive instead. waitForFunctionActive does not throw error in non-success cases.
 */
export declare const waitForFunctionActive: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionConfigurationCommandInput) => Promise<WaiterResult>;
/**
 * Waits for the function's State to be Active. This waiter uses GetFunctionConfiguration API. This should be used after new function creation.
 *  @param params - Waiter configuration options.
 *  @param input - The input to GetFunctionConfigurationCommand for polling.
 */
export declare const waitUntilFunctionActive: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionConfigurationCommandInput) => Promise<WaiterResult>;
