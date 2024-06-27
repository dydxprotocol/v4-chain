import { checkExceptions, createWaiter, WaiterState } from "@smithy/util-waiter";
import { DescribeTasksCommand } from "../commands/DescribeTasksCommand";
const checkState = async (client, input) => {
    let reason;
    try {
        const result = await client.send(new DescribeTasksCommand(input));
        reason = result;
        try {
            const returnComparator = () => {
                const flat_1 = [].concat(...result.tasks);
                const projection_3 = flat_1.map((element_2) => {
                    return element_2.lastStatus;
                });
                return projection_3;
            };
            for (const anyStringEq_4 of returnComparator()) {
                if (anyStringEq_4 == "STOPPED") {
                    return { state: WaiterState.FAILURE, reason };
                }
            }
        }
        catch (e) { }
        try {
            const returnComparator = () => {
                const flat_1 = [].concat(...result.failures);
                const projection_3 = flat_1.map((element_2) => {
                    return element_2.reason;
                });
                return projection_3;
            };
            for (const anyStringEq_4 of returnComparator()) {
                if (anyStringEq_4 == "MISSING") {
                    return { state: WaiterState.FAILURE, reason };
                }
            }
        }
        catch (e) { }
        try {
            const returnComparator = () => {
                const flat_1 = [].concat(...result.tasks);
                const projection_3 = flat_1.map((element_2) => {
                    return element_2.lastStatus;
                });
                return projection_3;
            };
            let allStringEq_5 = returnComparator().length > 0;
            for (const element_4 of returnComparator()) {
                allStringEq_5 = allStringEq_5 && element_4 == "RUNNING";
            }
            if (allStringEq_5) {
                return { state: WaiterState.SUCCESS, reason };
            }
        }
        catch (e) { }
    }
    catch (exception) {
        reason = exception;
    }
    return { state: WaiterState.RETRY, reason };
};
export const waitForTasksRunning = async (params, input) => {
    const serviceDefaults = { minDelay: 6, maxDelay: 120 };
    return createWaiter({ ...serviceDefaults, ...params }, input, checkState);
};
export const waitUntilTasksRunning = async (params, input) => {
    const serviceDefaults = { minDelay: 6, maxDelay: 120 };
    const result = await createWaiter({ ...serviceDefaults, ...params }, input, checkState);
    return checkExceptions(result);
};
