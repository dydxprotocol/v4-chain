import { checkExceptions, createWaiter, WaiterState } from "@smithy/util-waiter";
import { GetFunctionCommand } from "../commands/GetFunctionCommand";
const checkState = async (client, input) => {
    let reason;
    try {
        const result = await client.send(new GetFunctionCommand(input));
        reason = result;
        return { state: WaiterState.SUCCESS, reason };
    }
    catch (exception) {
        reason = exception;
        if (exception.name && exception.name == "ResourceNotFoundException") {
            return { state: WaiterState.RETRY, reason };
        }
    }
    return { state: WaiterState.RETRY, reason };
};
export const waitForFunctionExists = async (params, input) => {
    const serviceDefaults = { minDelay: 1, maxDelay: 120 };
    return createWaiter({ ...serviceDefaults, ...params }, input, checkState);
};
export const waitUntilFunctionExists = async (params, input) => {
    const serviceDefaults = { minDelay: 1, maxDelay: 120 };
    const result = await createWaiter({ ...serviceDefaults, ...params }, input, checkState);
    return checkExceptions(result);
};
