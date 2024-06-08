"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.validateOptions = void 0;
function validateOptions(_a) {
    var maxItems = _a.maxItems;
    if (maxItems !== undefined && maxItems < -1) {
        throw RangeError("Expected options.maxItems to be >= -1, but was given ".concat(maxItems, "."));
    }
}
exports.validateOptions = validateOptions;
//# sourceMappingURL=optionValidator.js.map