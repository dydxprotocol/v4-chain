import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { GetFunctionConfigurationCommandInput } from "../commands/GetFunctionConfigurationCommand";
import { LambdaClient } from "../LambdaClient";
export declare const waitForFunctionUpdated: (
  params: WaiterConfiguration<LambdaClient>,
  input: GetFunctionConfigurationCommandInput
) => Promise<WaiterResult>;
export declare const waitUntilFunctionUpdated: (
  params: WaiterConfiguration<LambdaClient>,
  input: GetFunctionConfigurationCommandInput
) => Promise<WaiterResult>;
