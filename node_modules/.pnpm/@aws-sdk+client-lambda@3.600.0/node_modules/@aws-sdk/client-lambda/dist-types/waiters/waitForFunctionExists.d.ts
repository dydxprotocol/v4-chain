import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { GetFunctionCommandInput } from "../commands/GetFunctionCommand";
import { LambdaClient } from "../LambdaClient";
/**
 *
 *  @deprecated Use waitUntilFunctionExists instead. waitForFunctionExists does not throw error in non-success cases.
 */
export declare const waitForFunctionExists: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionCommandInput) => Promise<WaiterResult>;
/**
 *
 *  @param params - Waiter configuration options.
 *  @param input - The input to GetFunctionCommand for polling.
 */
export declare const waitUntilFunctionExists: (params: WaiterConfiguration<LambdaClient>, input: GetFunctionCommandInput) => Promise<WaiterResult>;
