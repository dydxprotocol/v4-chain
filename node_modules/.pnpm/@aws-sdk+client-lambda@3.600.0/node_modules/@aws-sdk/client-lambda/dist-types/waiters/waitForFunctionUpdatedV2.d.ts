import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { GetFunctionCommandInput } from "../commands/GetFunctionCommand";
import { LambdaClient } from "../LambdaClient";
/**
 * Waits for the function's LastUpdateStatus to be Successful. This waiter uses GetFunction API. This should be used after function updates.
 *  @deprecated Use waitUntilFunctionUpdatedV2 instead. waitForFunctionUpdatedV2 does not throw error in non-success cases.
 */
export declare const waitForFunctionUpdatedV2: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionCommandInput) => Promise<WaiterResult>;
/**
 * Waits for the function's LastUpdateStatus to be Successful. This waiter uses GetFunction API. This should be used after function updates.
 *  @param params - Waiter configuration options.
 *  @param input - The input to GetFunctionCommand for polling.
 */
export declare const waitUntilFunctionUpdatedV2: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionCommandInput) => Promise<WaiterResult>;
