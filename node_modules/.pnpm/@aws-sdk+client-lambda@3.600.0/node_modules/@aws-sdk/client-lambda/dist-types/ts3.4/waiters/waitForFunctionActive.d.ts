import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { GetFunctionConfigurationCommandInput } from "../commands/GetFunctionConfigurationCommand";
import { LambdaClient } from "../LambdaClient";
export declare const waitForFunctionActive: (
  params: WaiterConfiguration<LambdaClient>,
  input: GetFunctionConfigurationCommandInput
) => Promise<WaiterResult>;
export declare const waitUntilFunctionActive: (
  params: WaiterConfiguration<LambdaClient>,
  input: GetFunctionConfigurationCommandInput
) => Promise<WaiterResult>;
