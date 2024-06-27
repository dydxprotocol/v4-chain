"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.smallIntToApi = exports.apiToBigInt = exports.apiToSmallInt = void 0;
const math_1 = require("@cosmjs/math");
const encodings_1 = require("./tendermint34/encodings");
/**
 * Takes an integer value from the Tendermint RPC API and
 * returns it as number.
 *
 * Only works within the safe integer range.
 */
function apiToSmallInt(input) {
    const asInt = typeof input === "number" ? new math_1.Int53(input) : math_1.Int53.fromString(input);
    return asInt.toNumber();
}
exports.apiToSmallInt = apiToSmallInt;
/**
 * Takes an integer value from the Tendermint RPC API and
 * returns it as BigInt.
 *
 * This supports the full uint64 and int64 ranges.
 */
function apiToBigInt(input) {
    (0, encodings_1.assertString)(input); // Runtime check on top of TypeScript just to be safe for semi-trusted API types
    if (!input.match(/^-?[0-9]+$/)) {
        throw new Error("Invalid string format");
    }
    return BigInt(input);
}
exports.apiToBigInt = apiToBigInt;
/**
 * Takes an integer in the safe integer range and returns
 * a string representation to be used in the Tendermint RPC API.
 */
function smallIntToApi(num) {
    return new math_1.Int53(num).toString();
}
exports.smallIntToApi = smallIntToApi;
//# sourceMappingURL=inthelpers.js.map