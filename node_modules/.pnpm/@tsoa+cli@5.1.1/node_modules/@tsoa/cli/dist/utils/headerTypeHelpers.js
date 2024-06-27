"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isSupportedHeaderDataType = exports.getHeaderType = void 0;
const exceptions_1 = require("../metadataGeneration/exceptions");
const typeResolver_1 = require("../metadataGeneration/typeResolver");
function getHeaderType(typeArgumentNodes, index, metadataGenerator) {
    if (!typeArgumentNodes || !typeArgumentNodes[index]) {
        return undefined;
    }
    const candidate = new typeResolver_1.TypeResolver(typeArgumentNodes[index], metadataGenerator).resolve();
    if (candidate && isSupportedHeaderDataType(candidate)) {
        return candidate;
    }
    else if (candidate) {
        throw new exceptions_1.GenerateMetadataError(`Unable to parse Header Type ${typeArgumentNodes[index].getText()}`, typeArgumentNodes[index]);
    }
    return undefined;
}
exports.getHeaderType = getHeaderType;
function isSupportedHeaderDataType(parameterType) {
    const supportedPathDataTypes = ['nestedObjectLiteral', 'refObject'];
    if (supportedPathDataTypes.find(t => t === parameterType.dataType)) {
        return true;
    }
    return false;
}
exports.isSupportedHeaderDataType = isSupportedHeaderDataType;
//# sourceMappingURL=headerTypeHelpers.js.map