"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getValue = exports.DEFAULT_RESPONSE_MEDIA_TYPE = exports.DEFAULT_REQUEST_MEDIA_TYPE = void 0;
exports.DEFAULT_REQUEST_MEDIA_TYPE = 'application/json';
exports.DEFAULT_RESPONSE_MEDIA_TYPE = 'application/json';
function getValue(type, member) {
    if (member === null) {
        return null;
    }
    switch (type) {
        case 'integer':
        case 'number':
            return Number(member);
        case 'boolean':
            return !!member;
        case 'string':
        default:
            return String(member);
    }
}
exports.getValue = getValue;
//# sourceMappingURL=swaggerUtils.js.map