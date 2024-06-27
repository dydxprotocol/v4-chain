import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { GetFunctionCommandInput } from "../commands/GetFunctionCommand";
import { LambdaClient } from "../LambdaClient";
export declare const waitForFunctionUpdatedV2: (
  params: WaiterConfiguration<LambdaClient>,
  input: GetFunctionCommandInput
) => Promise<WaiterResult>;
export declare const waitUntilFunctionUpdatedV2: (
  params: WaiterConfiguration<LambdaClient>,
  input: GetFunctionCommandInput
) => Promise<WaiterResult>;
