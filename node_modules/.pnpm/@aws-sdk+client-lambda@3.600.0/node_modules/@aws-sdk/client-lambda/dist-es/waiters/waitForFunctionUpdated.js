import { checkExceptions, createWaiter, WaiterState } from "@smithy/util-waiter";
import { GetFunctionConfigurationCommand, } from "../commands/GetFunctionConfigurationCommand";
const checkState = async (client, input) => {
    let reason;
    try {
        const result = await client.send(new GetFunctionConfigurationCommand(input));
        reason = result;
        try {
            const returnComparator = () => {
                return result.LastUpdateStatus;
            };
            if (returnComparator() === "Successful") {
                return { state: WaiterState.SUCCESS, reason };
            }
        }
        catch (e) { }
        try {
            const returnComparator = () => {
                return result.LastUpdateStatus;
            };
            if (returnComparator() === "Failed") {
                return { state: WaiterState.FAILURE, reason };
            }
        }
        catch (e) { }
        try {
            const returnComparator = () => {
                return result.LastUpdateStatus;
            };
            if (returnComparator() === "InProgress") {
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
export const waitForFunctionUpdated = async (params, input) => {
    const serviceDefaults = { minDelay: 5, maxDelay: 120 };
    return createWaiter({ ...serviceDefaults, ...params }, input, checkState);
};
export const waitUntilFunctionUpdated = async (params, input) => {
    const serviceDefaults = { minDelay: 5, maxDelay: 120 };
    const result = await createWaiter({ ...serviceDefaults, ...params }, input, checkState);
    return checkExceptions(result);
};
