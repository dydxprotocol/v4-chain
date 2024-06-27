import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { GetFunctionCommandInput } from "../commands/GetFunctionCommand";
import { LambdaClient } from "../LambdaClient";
/**
 * Waits for the function's State to be Active. This waiter uses GetFunction API. This should be used after new function creation.
 *  @deprecated Use waitUntilFunctionActiveV2 instead. waitForFunctionActiveV2 does not throw error in non-success cases.
 */
export declare const waitForFunctionActiveV2: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionCommandInput) => Promise<WaiterResult>;
/**
 * Waits for the function's State to be Active. This waiter uses GetFunction API. This should be used after new function creation.
 *  @param params - Waiter configuration options.
 *  @param input - The input to GetFunctionCommand for polling.
 */
export declare const waitUntilFunctionActiveV2: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionCommandInput) => Promise<WaiterResult>;
