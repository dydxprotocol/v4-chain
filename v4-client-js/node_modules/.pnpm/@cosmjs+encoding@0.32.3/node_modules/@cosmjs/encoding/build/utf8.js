"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.fromUtf8 = exports.toUtf8 = void 0;
function toUtf8(str) {
    return new TextEncoder().encode(str);
}
exports.toUtf8 = toUtf8;
/**
 * Takes UTF-8 data and decodes it to a string.
 *
 * In lossy mode, the [REPLACEMENT CHARACTER](https://en.wikipedia.org/wiki/Specials_(Unicode_block))
 * is used to substitude invalid encodings.
 * By default lossy mode is off and invalid data will lead to exceptions.
 */
function fromUtf8(data, lossy = false) {
    const fatal = !lossy;
    return new TextDecoder("utf-8", { fatal }).decode(data);
}
exports.fromUtf8 = fromUtf8;
//# sourceMappingURL=utf8.js.map