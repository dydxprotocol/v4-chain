"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isDefaultForAdditionalPropertiesAllowed = void 0;
/**
 * For Swagger, additionalProperties is implicitly allowed. So use this function to clarify that undefined should be associated with allowing additional properties
 * @param test if this is undefined then you should interpret it as a "yes"
 */
function isDefaultForAdditionalPropertiesAllowed(test) {
    return test === undefined;
}
exports.isDefaultForAdditionalPropertiesAllowed = isDefaultForAdditionalPropertiesAllowed;
//# sourceMappingURL=tsoa-route.js.map