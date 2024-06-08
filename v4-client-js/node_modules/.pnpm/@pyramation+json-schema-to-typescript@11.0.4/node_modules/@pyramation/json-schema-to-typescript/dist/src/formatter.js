"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.format = void 0;
var prettier_1 = require("prettier");
function format(code, options) {
    if (!options.format) {
        return code;
    }
    return (0, prettier_1.format)(code, __assign({ parser: 'typescript' }, options.style));
}
exports.format = format;
//# sourceMappingURL=formatter.js.map