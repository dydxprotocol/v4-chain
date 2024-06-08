"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isCompound = exports.isPrimitive = exports.getRootSchema = exports.Parent = void 0;
var lodash_1 = require("lodash");
exports.Parent = Symbol('Parent');
exports.getRootSchema = (0, lodash_1.memoize)(function (schema) {
    var parent = schema[exports.Parent];
    if (!parent) {
        return schema;
    }
    return (0, exports.getRootSchema)(parent);
});
function isPrimitive(schema) {
    return !(0, lodash_1.isPlainObject)(schema);
}
exports.isPrimitive = isPrimitive;
function isCompound(schema) {
    return Array.isArray(schema.type) || 'anyOf' in schema || 'oneOf' in schema;
}
exports.isCompound = isCompound;
//# sourceMappingURL=JSONSchema.js.map