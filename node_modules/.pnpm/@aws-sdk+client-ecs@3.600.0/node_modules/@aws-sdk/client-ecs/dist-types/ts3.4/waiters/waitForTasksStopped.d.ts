import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { DescribeTasksCommandInput } from "../commands/DescribeTasksCommand";
import { ECSClient } from "../ECSClient";
export declare const waitForTasksStopped: (
  params: WaiterConfiguration<ECSClient>,
  input: DescribeTasksCommandInput
) => Promise<WaiterResult>;
export declare const waitUntilTasksStopped: (
  params: WaiterConfiguration<ECSClient>,
  input: DescribeTasksCommandInput
) => Promise<WaiterResult>;
