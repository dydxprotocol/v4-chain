"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.typesOfSchema = void 0;
var lodash_1 = require("lodash");
var JSONSchema_1 = require("./types/JSONSchema");
/**
 * Duck types a JSONSchema schema or property to determine which kind of AST node to parse it into.
 *
 * Due to what some might say is an oversight in the JSON-Schema spec, a given schema may
 * implicitly be an *intersection* of multiple JSON-Schema directives (ie. multiple TypeScript
 * types). The spec leaves it up to implementations to decide what to do with this
 * loosely-defined behavior.
 */
function typesOfSchema(schema) {
    // tsType is an escape hatch that supercedes all other directives
    if (schema.tsType) {
        return ['CUSTOM_TYPE'];
    }
    // Collect matched types
    var matchedTypes = [];
    for (var _i = 0, _a = Object.entries(matchers); _i < _a.length; _i++) {
        var _b = _a[_i], schemaType = _b[0], f = _b[1];
        if (f(schema)) {
            matchedTypes.push(schemaType);
        }
    }
    // Default to an unnamed schema
    if (!matchedTypes.length) {
        return ['UNNAMED_SCHEMA'];
    }
    return matchedTypes;
}
exports.typesOfSchema = typesOfSchema;
var matchers = {
    ALL_OF: function (schema) {
        return 'allOf' in schema;
    },
    ANY: function (schema) {
        if (Object.keys(schema).length === 0) {
            // The empty schema {} validates any value
            // @see https://json-schema.org/draft-07/json-schema-core.html#rfc.section.4.3.1
            return true;
        }
        return schema.type === 'any';
    },
    ANY_OF: function (schema) {
        return 'anyOf' in schema;
    },
    BOOLEAN: function (schema) {
        if ('enum' in schema) {
            return false;
        }
        if (schema.type === 'boolean') {
            return true;
        }
        if (!(0, JSONSchema_1.isCompound)(schema) && typeof schema.default === 'boolean') {
            return true;
        }
        return false;
    },
    CUSTOM_TYPE: function () {
        return false; // Explicitly handled before we try to match
    },
    NAMED_ENUM: function (schema) {
        return 'enum' in schema && 'tsEnumNames' in schema;
    },
    NAMED_SCHEMA: function (schema) {
        // 8.2.1. The presence of "$id" in a subschema indicates that the subschema constitutes a distinct schema resource within a single schema document.
        return '$id' in schema && ('patternProperties' in schema || 'properties' in schema);
    },
    NULL: function (schema) {
        return schema.type === 'null';
    },
    NUMBER: function (schema) {
        if ('enum' in schema) {
            return false;
        }
        if (schema.type === 'integer' || schema.type === 'number') {
            return true;
        }
        if (!(0, JSONSchema_1.isCompound)(schema) && typeof schema.default === 'number') {
            return true;
        }
        return false;
    },
    OBJECT: function (schema) {
        return (schema.type === 'object' &&
            !(0, lodash_1.isPlainObject)(schema.additionalProperties) &&
            !schema.allOf &&
            !schema.anyOf &&
            !schema.oneOf &&
            !schema.patternProperties &&
            !schema.properties &&
            !schema.required);
    },
    ONE_OF: function (schema) {
        return 'oneOf' in schema;
    },
    REFERENCE: function (schema) {
        return '$ref' in schema;
    },
    STRING: function (schema) {
        if ('enum' in schema) {
            return false;
        }
        if (schema.type === 'string') {
            return true;
        }
        if (!(0, JSONSchema_1.isCompound)(schema) && typeof schema.default === 'string') {
            return true;
        }
        return false;
    },
    TYPED_ARRAY: function (schema) {
        if (schema.type && schema.type !== 'array') {
            return false;
        }
        return 'items' in schema;
    },
    UNION: function (schema) {
        return Array.isArray(schema.type);
    },
    UNNAMED_ENUM: function (schema) {
        if ('tsEnumNames' in schema) {
            return false;
        }
        if (schema.type &&
            schema.type !== 'boolean' &&
            schema.type !== 'integer' &&
            schema.type !== 'number' &&
            schema.type !== 'string') {
            return false;
        }
        return 'enum' in schema;
    },
    UNNAMED_SCHEMA: function () {
        return false; // Explicitly handled as the default case
    },
    UNTYPED_ARRAY: function (schema) {
        return schema.type === 'array' && !('items' in schema);
    }
};
//# sourceMappingURL=typesOfSchema.js.map