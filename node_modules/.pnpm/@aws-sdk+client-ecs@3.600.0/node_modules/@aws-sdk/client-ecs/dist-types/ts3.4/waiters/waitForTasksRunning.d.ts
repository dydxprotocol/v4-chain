import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { DescribeTasksCommandInput } from "../commands/DescribeTasksCommand";
import { ECSClient } from "../ECSClient";
export declare const waitForTasksRunning: (
  params: WaiterConfiguration<ECSClient>,
  input: DescribeTasksCommandInput
) => Promise<WaiterResult>;
export declare const waitUntilTasksRunning: (
  params: WaiterConfiguration<ECSClient>,
  input: DescribeTasksCommandInput
) => Promise<WaiterResult>;
