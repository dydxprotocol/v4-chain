import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { GetFunctionCommandInput } from "../commands/GetFunctionCommand";
import { LambdaClient } from "../LambdaClient";
export declare const waitForFunctionExists: (
  params: WaiterConfiguration<LambdaClient>,
  input: GetFunctionCommandInput
) => Promise<WaiterResult>;
export declare const waitUntilFunctionExists: (
  params: WaiterConfiguration<LambdaClient>,
  input: GetFunctionCommandInput
) => Promise<WaiterResult>;
