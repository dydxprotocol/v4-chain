"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.convertBracesPathParams = exports.convertColonPathParams = exports.normalisePath = void 0;
/**
 * Removes all '/', '\', and spaces from the beginning and end of the path
 * Replaces all '/', '\', and spaces between sections of the path
 * Adds prefix and suffix if supplied
 * Replace ':pathParam' with '{pathParam}'
 */
function normalisePath(path, withPrefix, withSuffix, skipPrefixAndSuffixIfEmpty = true) {
    if ((!path || path === '/') && skipPrefixAndSuffixIfEmpty) {
        return '';
    }
    if (!path || typeof path !== 'string') {
        path = '' + path;
    }
    // normalise beginning and end of the path
    let normalised = path.replace(/^[/\\\s]+|[/\\\s]+$/g, '');
    normalised = withPrefix ? withPrefix + normalised : normalised;
    normalised = withSuffix ? normalised + withSuffix : normalised;
    // normalise / signs amount in all path
    normalised = normalised.replace(/[/\\\s]+/g, '/');
    return normalised;
}
exports.normalisePath = normalisePath;
function convertColonPathParams(path) {
    if (!path || typeof path !== 'string') {
        return path;
    }
    const normalised = path.replace(/:([^/]+)/g, '{$1}');
    return normalised;
}
exports.convertColonPathParams = convertColonPathParams;
function convertBracesPathParams(path) {
    return path.replace(/{(\w*)}/g, ':$1');
}
exports.convertBracesPathParams = convertBracesPathParams;
//# sourceMappingURL=pathUtils.js.map