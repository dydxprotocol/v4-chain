import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { DescribeServicesCommandInput } from "../commands/DescribeServicesCommand";
import { ECSClient } from "../ECSClient";
export declare const waitForServicesStable: (
  params: WaiterConfiguration<ECSClient>,
  input: DescribeServicesCommandInput
) => Promise<WaiterResult>;
export declare const waitUntilServicesStable: (
  params: WaiterConfiguration<ECSClient>,
  input: DescribeServicesCommandInput
) => Promise<WaiterResult>;
