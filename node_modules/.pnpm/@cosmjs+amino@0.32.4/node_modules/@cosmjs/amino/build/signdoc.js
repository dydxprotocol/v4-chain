"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.serializeSignDoc = exports.escapeCharacters = exports.makeSignDoc = exports.sortedJsonStringify = void 0;
/* eslint-disable @typescript-eslint/naming-convention */
const encoding_1 = require("@cosmjs/encoding");
const math_1 = require("@cosmjs/math");
function sortedObject(obj) {
    if (typeof obj !== "object" || obj === null) {
        return obj;
    }
    if (Array.isArray(obj)) {
        return obj.map(sortedObject);
    }
    const sortedKeys = Object.keys(obj).sort();
    const result = {};
    // NOTE: Use forEach instead of reduce for performance with large objects eg Wasm code
    sortedKeys.forEach((key) => {
        result[key] = sortedObject(obj[key]);
    });
    return result;
}
/** Returns a JSON string with objects sorted by key */
// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
function sortedJsonStringify(obj) {
    return JSON.stringify(sortedObject(obj));
}
exports.sortedJsonStringify = sortedJsonStringify;
function makeSignDoc(msgs, fee, chainId, memo, accountNumber, sequence, timeout_height) {
    return {
        chain_id: chainId,
        account_number: math_1.Uint53.fromString(accountNumber.toString()).toString(),
        sequence: math_1.Uint53.fromString(sequence.toString()).toString(),
        fee: fee,
        msgs: msgs,
        memo: memo || "",
        ...(timeout_height && { timeout_height: timeout_height.toString() }),
    };
}
exports.makeSignDoc = makeSignDoc;
/**
 * Takes a valid JSON document and performs the following escapings in string values:
 *
 * `&` -> `\u0026`
 * `<` -> `\u003c`
 * `>` -> `\u003e`
 *
 * Since those characters do not occur in other places of the JSON document, only
 * string values are affected.
 *
 * If the input is invalid JSON, the behaviour is undefined.
 */
function escapeCharacters(input) {
    // When we migrate to target es2021 or above, we can use replaceAll instead of global patterns.
    // https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/replaceAll
    const amp = /&/g;
    const lt = /</g;
    const gt = />/g;
    return input.replace(amp, "\\u0026").replace(lt, "\\u003c").replace(gt, "\\u003e");
}
exports.escapeCharacters = escapeCharacters;
function serializeSignDoc(signDoc) {
    const serialized = escapeCharacters(sortedJsonStringify(signDoc));
    return (0, encoding_1.toUtf8)(serialized);
}
exports.serializeSignDoc = serializeSignDoc;
//# sourceMappingURL=signdoc.js.map