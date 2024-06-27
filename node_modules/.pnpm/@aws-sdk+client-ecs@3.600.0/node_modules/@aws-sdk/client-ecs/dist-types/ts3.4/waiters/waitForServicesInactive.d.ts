import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { DescribeServicesCommandInput } from "../commands/DescribeServicesCommand";
import { ECSClient } from "../ECSClient";
export declare const waitForServicesInactive: (
  params: WaiterConfiguration<ECSClient>,
  input: DescribeServicesCommandInput
) => Promise<WaiterResult>;
export declare const waitUntilServicesInactive: (
  params: WaiterConfiguration<ECSClient>,
  input: DescribeServicesCommandInput
) => Promise<WaiterResult>;
