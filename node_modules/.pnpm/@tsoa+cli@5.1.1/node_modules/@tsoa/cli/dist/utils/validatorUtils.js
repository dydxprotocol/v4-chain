"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getPropertyValidators = exports.getParameterValidators = void 0;
const validator_1 = require("validator");
const exceptions_1 = require("./../metadataGeneration/exceptions");
const jsDocUtils_1 = require("./jsDocUtils");
function getParameterValidators(parameter, parameterName) {
    if (!parameter.parent) {
        return {};
    }
    const getCommentValue = (comment) => comment && comment.split(' ')[0];
    const tags = (0, jsDocUtils_1.getJSDocTags)(parameter.parent, tag => {
        const { comment } = tag;
        return getParameterTagSupport().some(value => !!(0, jsDocUtils_1.commentToString)(comment) && value === tag.tagName.text && getCommentValue((0, jsDocUtils_1.commentToString)(comment)) === parameterName);
    });
    function getErrorMsg(comment, isValue = true) {
        if (!comment) {
            return;
        }
        if (isValue) {
            const indexOf = comment.indexOf(' ');
            if (indexOf > 0) {
                return comment.substr(indexOf + 1);
            }
            else {
                return undefined;
            }
        }
        else {
            return comment;
        }
    }
    return tags.reduce((validateObj, tag) => {
        var _a, _b;
        if (!tag.comment) {
            return validateObj;
        }
        const name = tag.tagName.text;
        const comment = (_a = (0, jsDocUtils_1.commentToString)(tag.comment)) === null || _a === void 0 ? void 0 : _a.substr((((_b = (0, jsDocUtils_1.commentToString)(tag.comment)) === null || _b === void 0 ? void 0 : _b.indexOf(' ')) || -1) + 1).trim();
        const value = getCommentValue(comment);
        switch (name) {
            case 'uniqueItems':
                validateObj[name] = {
                    errorMsg: getErrorMsg(comment, false),
                    value: undefined,
                };
                break;
            case 'minimum':
            case 'maximum':
            case 'minItems':
            case 'maxItems':
            case 'minLength':
            case 'maxLength':
                if (isNaN(value)) {
                    throw new exceptions_1.GenerateMetadataError(`${name} parameter use number.`);
                }
                validateObj[name] = {
                    errorMsg: getErrorMsg(comment),
                    value: Number(value),
                };
                break;
            case 'minDate':
            case 'maxDate':
                if (!validator_1.default.isISO8601(String(value), { strict: true })) {
                    throw new exceptions_1.GenerateMetadataError(`${name} parameter use date format ISO 8601 ex. 2017-05-14, 2017-05-14T05:18Z`);
                }
                validateObj[name] = {
                    errorMsg: getErrorMsg(comment),
                    value,
                };
                break;
            case 'pattern':
                if (typeof value !== 'string') {
                    throw new exceptions_1.GenerateMetadataError(`${name} parameter use string.`);
                }
                validateObj[name] = {
                    errorMsg: getErrorMsg(comment),
                    value: removeSurroundingQuotes(value),
                };
                break;
            default:
                if (name.startsWith('is')) {
                    const errorMsg = getErrorMsg(comment, false);
                    if (errorMsg) {
                        validateObj[name] = {
                            errorMsg,
                            value: undefined,
                        };
                    }
                }
                break;
        }
        return validateObj;
    }, {});
}
exports.getParameterValidators = getParameterValidators;
function getPropertyValidators(property) {
    const tags = (0, jsDocUtils_1.getJSDocTags)(property, tag => {
        return getParameterTagSupport().some(value => value === tag.tagName.text);
    });
    function getValue(comment) {
        if (!comment) {
            return;
        }
        return comment.split(' ')[0];
    }
    function getErrorMsg(comment, isValue = true) {
        if (!comment) {
            return;
        }
        if (isValue) {
            const indexOf = comment.indexOf(' ');
            if (indexOf > 0) {
                return comment.substr(indexOf + 1);
            }
            else {
                return undefined;
            }
        }
        else {
            return comment;
        }
    }
    return tags.reduce((validateObj, tag) => {
        const name = tag.tagName.text;
        const comment = tag.comment;
        const value = getValue((0, jsDocUtils_1.commentToString)(comment));
        switch (name) {
            case 'uniqueItems':
                validateObj[name] = {
                    errorMsg: getErrorMsg((0, jsDocUtils_1.commentToString)(comment), false),
                    value: undefined,
                };
                break;
            case 'minimum':
            case 'maximum':
            case 'minItems':
            case 'maxItems':
            case 'minLength':
            case 'maxLength':
                if (isNaN(value)) {
                    throw new exceptions_1.GenerateMetadataError(`${name} parameter use number.`);
                }
                validateObj[name] = {
                    errorMsg: getErrorMsg((0, jsDocUtils_1.commentToString)(comment)),
                    value: Number(value),
                };
                break;
            case 'minDate':
            case 'maxDate':
                if (!validator_1.default.isISO8601(String(value), { strict: true })) {
                    throw new exceptions_1.GenerateMetadataError(`${name} parameter use date format ISO 8601 ex. 2017-05-14, 2017-05-14T05:18Z`);
                }
                validateObj[name] = {
                    errorMsg: getErrorMsg((0, jsDocUtils_1.commentToString)(comment)),
                    value,
                };
                break;
            case 'pattern':
                if (typeof value !== 'string') {
                    throw new exceptions_1.GenerateMetadataError(`${name} parameter use string.`);
                }
                validateObj[name] = {
                    errorMsg: getErrorMsg((0, jsDocUtils_1.commentToString)(comment)),
                    value: removeSurroundingQuotes(value),
                };
                break;
            default:
                if (name.startsWith('is')) {
                    const errorMsg = getErrorMsg((0, jsDocUtils_1.commentToString)(comment), false);
                    if (errorMsg) {
                        validateObj[name] = {
                            errorMsg,
                            value: undefined,
                        };
                    }
                }
                break;
        }
        return validateObj;
    }, {});
}
exports.getPropertyValidators = getPropertyValidators;
function getParameterTagSupport() {
    return [
        'isString',
        'isBoolean',
        'isInt',
        'isLong',
        'isFloat',
        'isDouble',
        'isDate',
        'isDateTime',
        'minItems',
        'maxItems',
        'uniqueItems',
        'minLength',
        'maxLength',
        'pattern',
        'minimum',
        'maximum',
        'minDate',
        'maxDate',
    ];
}
function removeSurroundingQuotes(str) {
    if (str.startsWith('`') && str.endsWith('`')) {
        return str.substring(1, str.length - 1);
    }
    if (str.startsWith('```') && str.endsWith('```')) {
        return str.substring(3, str.length - 3);
    }
    return str;
}
//# sourceMappingURL=validatorUtils.js.map