"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.resolveHttpAuthSchemeConfig = exports.defaultLambdaHttpAuthSchemeProvider = exports.defaultLambdaHttpAuthSchemeParametersProvider = void 0;
const core_1 = require("@aws-sdk/core");
const util_middleware_1 = require("@smithy/util-middleware");
const defaultLambdaHttpAuthSchemeParametersProvider = async (config, context, input) => {
    return {
        operation: (0, util_middleware_1.getSmithyContext)(context).operation,
        region: (await (0, util_middleware_1.normalizeProvider)(config.region)()) ||
            (() => {
                throw new Error("expected `region` to be configured for `aws.auth#sigv4`");
            })(),
    };
};
exports.defaultLambdaHttpAuthSchemeParametersProvider = defaultLambdaHttpAuthSchemeParametersProvider;
function createAwsAuthSigv4HttpAuthOption(authParameters) {
    return {
        schemeId: "aws.auth#sigv4",
        signingProperties: {
            name: "lambda",
            region: authParameters.region,
        },
        propertiesExtractor: (config, context) => ({
            signingProperties: {
                config,
                context,
            },
        }),
    };
}
const defaultLambdaHttpAuthSchemeProvider = (authParameters) => {
    const options = [];
    switch (authParameters.operation) {
        default: {
            options.push(createAwsAuthSigv4HttpAuthOption(authParameters));
        }
    }
    return options;
};
exports.defaultLambdaHttpAuthSchemeProvider = defaultLambdaHttpAuthSchemeProvider;
const resolveHttpAuthSchemeConfig = (config) => {
    const config_0 = (0, core_1.resolveAwsSdkSigV4Config)(config);
    return {
        ...config_0,
    };
};
exports.resolveHttpAuthSchemeConfig = resolveHttpAuthSchemeConfig;
