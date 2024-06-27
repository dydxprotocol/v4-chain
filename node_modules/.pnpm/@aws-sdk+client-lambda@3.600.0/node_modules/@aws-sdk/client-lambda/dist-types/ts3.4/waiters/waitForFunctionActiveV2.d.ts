import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { GetFunctionCommandInput } from "../commands/GetFunctionCommand";
import { LambdaClient } from "../LambdaClient";
export declare const waitForFunctionActiveV2: (
  params: WaiterConfiguration<LambdaClient>,
  input: GetFunctionCommandInput
) => Promise<WaiterResult>;
export declare const waitUntilFunctionActiveV2: (
  params: WaiterConfiguration<LambdaClient>,
  input: GetFunctionCommandInput
) => Promise<WaiterResult>;
