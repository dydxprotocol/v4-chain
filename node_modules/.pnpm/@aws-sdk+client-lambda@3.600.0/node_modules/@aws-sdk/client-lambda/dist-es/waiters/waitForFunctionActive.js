import { checkExceptions, createWaiter, WaiterState } from "@smithy/util-waiter";
import { GetFunctionConfigurationCommand, } from "../commands/GetFunctionConfigurationCommand";
const checkState = async (client, input) => {
    let reason;
    try {
        const result = await client.send(new GetFunctionConfigurationCommand(input));
        reason = result;
        try {
            const returnComparator = () => {
                return result.State;
            };
            if (returnComparator() === "Active") {
                return { state: WaiterState.SUCCESS, reason };
            }
        }
        catch (e) { }
        try {
            const returnComparator = () => {
                return result.State;
            };
            if (returnComparator() === "Failed") {
                return { state: WaiterState.FAILURE, reason };
            }
        }
        catch (e) { }
        try {
            const returnComparator = () => {
                return result.State;
            };
            if (returnComparator() === "Pending") {
                return { state: WaiterState.RETRY, reason };
            }
        }
        catch (e) { }
    }
    catch (exception) {
        reason = exception;
    }
    return { state: WaiterState.RETRY, reason };
};
export const waitForFunctionActive = async (params, input) => {
    const serviceDefaults = { minDelay: 5, maxDelay: 120 };
    return createWaiter({ ...serviceDefaults, ...params }, input, checkState);
};
export const waitUntilFunctionActive = async (params, input) => {
    const serviceDefaults = { minDelay: 5, maxDelay: 120 };
    const result = await createWaiter({ ...serviceDefaults, ...params }, input, checkState);
    return checkExceptions(result);
};
