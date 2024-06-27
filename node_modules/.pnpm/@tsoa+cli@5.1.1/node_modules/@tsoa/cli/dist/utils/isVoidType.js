"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isVoidType = void 0;
const isVoidType = (type) => {
    if (type.dataType === 'void' || type.dataType === 'undefined') {
        return true;
    }
    else if (type.dataType === 'refAlias') {
        return (0, exports.isVoidType)(type.type);
    }
    else {
        return false;
    }
};
exports.isVoidType = isVoidType;
//# sourceMappingURL=isVoidType.js.map