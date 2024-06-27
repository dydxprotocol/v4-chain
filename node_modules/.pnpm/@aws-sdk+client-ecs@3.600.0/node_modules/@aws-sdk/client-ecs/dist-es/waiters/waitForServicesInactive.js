import { checkExceptions, createWaiter, WaiterState } from "@smithy/util-waiter";
import { DescribeServicesCommand } from "../commands/DescribeServicesCommand";
const checkState = async (client, input) => {
    let reason;
    try {
        const result = await client.send(new DescribeServicesCommand(input));
        reason = result;
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
                const flat_1 = [].concat(...result.services);
                const projection_3 = flat_1.map((element_2) => {
                    return element_2.status;
                });
                return projection_3;
            };
            for (const anyStringEq_4 of returnComparator()) {
                if (anyStringEq_4 == "INACTIVE") {
                    return { state: WaiterState.SUCCESS, reason };
                }
            }
        }
        catch (e) { }
    }
    catch (exception) {
        reason = exception;
    }
    return { state: WaiterState.RETRY, reason };
};
export const waitForServicesInactive = async (params, input) => {
    const serviceDefaults = { minDelay: 15, maxDelay: 120 };
    return createWaiter({ ...serviceDefaults, ...params }, input, checkState);
};
export const waitUntilServicesInactive = async (params, input) => {
    const serviceDefaults = { minDelay: 15, maxDelay: 120 };
    const result = await createWaiter({ ...serviceDefaults, ...params }, input, checkState);
    return checkExceptions(result);
};
