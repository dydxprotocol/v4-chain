"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.validate = void 0;
var utils_1 = require("./utils");
var rules = new Map();
rules.set('Enum members and tsEnumNames must be of the same length', function (schema) {
    if (schema.enum && schema.tsEnumNames && schema.enum.length !== schema.tsEnumNames.length) {
        return false;
    }
});
rules.set('tsEnumNames must be an array of strings', function (schema) {
    if (schema.tsEnumNames && schema.tsEnumNames.some(function (_) { return typeof _ !== 'string'; })) {
        return false;
    }
});
rules.set('When both maxItems and minItems are present, maxItems >= minItems', function (schema) {
    var maxItems = schema.maxItems, minItems = schema.minItems;
    if (typeof maxItems === 'number' && typeof minItems === 'number') {
        return maxItems >= minItems;
    }
});
rules.set('When maxItems exists, maxItems >= 0', function (schema) {
    var maxItems = schema.maxItems;
    if (typeof maxItems === 'number') {
        return maxItems >= 0;
    }
});
rules.set('When minItems exists, minItems >= 0', function (schema) {
    var minItems = schema.minItems;
    if (typeof minItems === 'number') {
        return minItems >= 0;
    }
});
function validate(schema, filename) {
    var errors = [];
    rules.forEach(function (rule, ruleName) {
        (0, utils_1.traverse)(schema, function (schema, key) {
            if (rule(schema) === false) {
                errors.push("Error at key \"".concat(key, "\" in file \"").concat(filename, "\": ").concat(ruleName));
            }
            return schema;
        });
    });
    return errors;
}
exports.validate = validate;
//# sourceMappingURL=validator.js.map