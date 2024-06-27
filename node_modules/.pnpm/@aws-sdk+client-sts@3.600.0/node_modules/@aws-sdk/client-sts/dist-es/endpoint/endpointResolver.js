import { awsEndpointFunctions } from "@aws-sdk/util-endpoints";
import { customEndpointFunctions, resolveEndpoint } from "@smithy/util-endpoints";
import { ruleSet } from "./ruleset";
export const defaultEndpointResolver = (endpointParams, context = {}) => {
    return resolveEndpoint(ruleSet, {
        endpointParams: endpointParams,
        logger: context.logger,
    });
};
customEndpointFunctions.aws = awsEndpointFunctions;
