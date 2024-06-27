import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { GetFunctionConfigurationCommandInput } from "../commands/GetFunctionConfigurationCommand";
import { LambdaClient } from "../LambdaClient";
/**
 * Waits for the function's LastUpdateStatus to be Successful. This waiter uses GetFunctionConfiguration API. This should be used after function updates.
 *  @deprecated Use waitUntilFunctionUpdated instead. waitForFunctionUpdated does not throw error in non-success cases.
 */
export declare const waitForFunctionUpdated: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionConfigurationCommandInput) => Promise<WaiterResult>;
/**
 * Waits for the function's LastUpdateStatus to be Successful. This waiter uses GetFunctionConfiguration API. This should be used after function updates.
 *  @param params - Waiter configuration options.
 *  @param input - The input to GetFunctionConfigurationCommand for polling.
 */
export declare const waitUntilFunctionUpdated: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionConfigurationCommandInput) => Promise<WaiterResult>;
