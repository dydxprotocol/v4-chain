import { WaiterConfiguration, WaiterResult } from "@smithy/util-waiter";
import { DescribeServicesCommandInput } from "../commands/DescribeServicesCommand";
import { ECSClient } from "../ECSClient";
/**
 *
 *  @deprecated Use waitUntilServicesStable instead. waitForServicesStable does not throw error in non-success cases.
 */
export declare const waitForServicesStable: (params: WaiterConfiguration<ECSClient>, input: DescribeServicesCommandInput) => Promise<WaiterResult>;
/**
 *
 *  @param params - Waiter configuration options.
 *  @param input - The input to DescribeServicesCommand for polling.
 */
export declare const waitUntilServicesStable: (params: WaiterConfiguration<ECSClient>, input: DescribeServicesCommandInput) => Promise<WaiterResult>;
