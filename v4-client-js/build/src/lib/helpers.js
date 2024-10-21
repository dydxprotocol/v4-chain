"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ByteArrayEncoding = void 0;
exports.convertPartialTransactionOptionsToFull = convertPartialTransactionOptionsToFull;
exports.stripHexPrefix = stripHexPrefix;
exports.encodeJson = encodeJson;
const encoding_1 = require("@cosmjs/encoding");
const bignumber_js_1 = __importDefault(require("bignumber.js"));
const long_1 = __importDefault(require("long"));
const constants_1 = require("./constants");
/**
 * @description Either return undefined or insert default sequence value into
 * `partialTransactionOptions` if it does not exist.
 *
 * @returns undefined or full TransactionOptions.
 */
function convertPartialTransactionOptionsToFull(partialTransactionOptions) {
    if (partialTransactionOptions === undefined) {
        return undefined;
    }
    return {
        sequence: constants_1.DEFAULT_SEQUENCE,
        ...partialTransactionOptions,
    };
}
/**
 * @description Strip '0x' prefix from input string. If there is no '0x' prefix, return the original
 * input.
 *
 * @returns input without '0x' prefix or original input if no prefix.
 */
function stripHexPrefix(input) {
    if (input.indexOf('0x') === 0) {
        return input.slice(2);
    }
    return input;
}
function toBigInt(u) {
    if (u.length <= 1) {
        return BigInt(0);
    }
    // eslint-disable-next-line no-bitwise
    const negated = (u[0] & 1) === 1;
    const hex = Buffer.from(u.slice(1)).toString('hex');
    const abs = BigInt(`0x${hex}`);
    return negated ? -abs : abs;
}
var ByteArrayEncoding;
(function (ByteArrayEncoding) {
    ByteArrayEncoding["HEX"] = "hex";
    ByteArrayEncoding["BIGINT"] = "bigint";
})(ByteArrayEncoding || (exports.ByteArrayEncoding = ByteArrayEncoding = {}));
function encodeJson(object, byteArrayEncoding = ByteArrayEncoding.HEX) {
    // eslint-disable-next-line prefer-arrow-callback
    return JSON.stringify(object, function replacer(_key, value) {
        // Even though we set the an UInt8Array as the value,
        // it comes in here as an object with UInt8Array as the buffer property.
        if (value instanceof bignumber_js_1.default) {
            return value.toString();
        }
        if (value instanceof long_1.default) {
            return value.toString();
        }
        if ((value === null || value === void 0 ? void 0 : value.buffer) instanceof Uint8Array) {
            if (byteArrayEncoding === ByteArrayEncoding.HEX) {
                return (0, encoding_1.toHex)(value.buffer);
            }
            else {
                return toBigInt(value.buffer).toString();
            }
        }
        else if (value instanceof Uint8Array) {
            if (byteArrayEncoding === ByteArrayEncoding.HEX) {
                return (0, encoding_1.toHex)(value);
            }
            else {
                return toBigInt(value).toString();
            }
        }
        return value;
    });
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaGVscGVycy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9saWIvaGVscGVycy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFhQSx3RkFXQztBQVFELHdDQU1DO0FBa0JELGdDQTZCQztBQXJGRCwrQ0FBeUM7QUFDekMsZ0VBQXFDO0FBQ3JDLGdEQUF3QjtBQUd4QiwyQ0FBK0M7QUFFL0M7Ozs7O0dBS0c7QUFDSCxTQUFnQixzQ0FBc0MsQ0FDcEQseUJBQXFEO0lBRXJELElBQUkseUJBQXlCLEtBQUssU0FBUyxFQUFFLENBQUM7UUFDNUMsT0FBTyxTQUFTLENBQUM7SUFDbkIsQ0FBQztJQUVELE9BQU87UUFDTCxRQUFRLEVBQUUsNEJBQWdCO1FBQzFCLEdBQUcseUJBQXlCO0tBQzdCLENBQUM7QUFDSixDQUFDO0FBRUQ7Ozs7O0dBS0c7QUFDSCxTQUFnQixjQUFjLENBQUMsS0FBYTtJQUMxQyxJQUFJLEtBQUssQ0FBQyxPQUFPLENBQUMsSUFBSSxDQUFDLEtBQUssQ0FBQyxFQUFFLENBQUM7UUFDOUIsT0FBTyxLQUFLLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFDO0lBQ3hCLENBQUM7SUFFRCxPQUFPLEtBQUssQ0FBQztBQUNmLENBQUM7QUFFRCxTQUFTLFFBQVEsQ0FBQyxDQUFhO0lBQzdCLElBQUksQ0FBQyxDQUFDLE1BQU0sSUFBSSxDQUFDLEVBQUUsQ0FBQztRQUNsQixPQUFPLE1BQU0sQ0FBQyxDQUFDLENBQUMsQ0FBQztJQUNuQixDQUFDO0lBQ0Qsc0NBQXNDO0lBQ3RDLE1BQU0sT0FBTyxHQUFZLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxHQUFHLENBQUMsQ0FBQyxLQUFLLENBQUMsQ0FBQztJQUMxQyxNQUFNLEdBQUcsR0FBVyxNQUFNLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxRQUFRLENBQUMsS0FBSyxDQUFDLENBQUM7SUFDNUQsTUFBTSxHQUFHLEdBQVcsTUFBTSxDQUFDLEtBQUssR0FBRyxFQUFFLENBQUMsQ0FBQztJQUN2QyxPQUFPLE9BQU8sQ0FBQyxDQUFDLENBQUMsQ0FBQyxHQUFHLENBQUMsQ0FBQyxDQUFDLEdBQUcsQ0FBQztBQUM5QixDQUFDO0FBRUQsSUFBWSxpQkFHWDtBQUhELFdBQVksaUJBQWlCO0lBQzNCLGdDQUFXLENBQUE7SUFDWCxzQ0FBaUIsQ0FBQTtBQUNuQixDQUFDLEVBSFcsaUJBQWlCLGlDQUFqQixpQkFBaUIsUUFHNUI7QUFFRCxTQUFnQixVQUFVLENBQ3hCLE1BQWUsRUFDZixvQkFBdUMsaUJBQWlCLENBQUMsR0FBRztJQUU1RCxpREFBaUQ7SUFDakQsT0FBTyxJQUFJLENBQUMsU0FBUyxDQUFDLE1BQU0sRUFBRSxTQUFTLFFBQVEsQ0FBQyxJQUFJLEVBQUUsS0FBSztRQUN6RCxxREFBcUQ7UUFDckQsd0VBQXdFO1FBQ3hFLElBQUksS0FBSyxZQUFZLHNCQUFTLEVBQUUsQ0FBQztZQUMvQixPQUFPLEtBQUssQ0FBQyxRQUFRLEVBQUUsQ0FBQztRQUMxQixDQUFDO1FBQ0QsSUFBSSxLQUFLLFlBQVksY0FBSSxFQUFFLENBQUM7WUFDMUIsT0FBTyxLQUFLLENBQUMsUUFBUSxFQUFFLENBQUM7UUFDMUIsQ0FBQztRQUNELElBQUksQ0FBQSxLQUFLLGFBQUwsS0FBSyx1QkFBTCxLQUFLLENBQUUsTUFBTSxhQUFZLFVBQVUsRUFBRSxDQUFDO1lBQ3hDLElBQUksaUJBQWlCLEtBQUssaUJBQWlCLENBQUMsR0FBRyxFQUFFLENBQUM7Z0JBQ2hELE9BQU8sSUFBQSxnQkFBSyxFQUFDLEtBQUssQ0FBQyxNQUFNLENBQUMsQ0FBQztZQUM3QixDQUFDO2lCQUFNLENBQUM7Z0JBQ04sT0FBTyxRQUFRLENBQUMsS0FBSyxDQUFDLE1BQU0sQ0FBQyxDQUFDLFFBQVEsRUFBRSxDQUFDO1lBQzNDLENBQUM7UUFDSCxDQUFDO2FBQU0sSUFBSSxLQUFLLFlBQVksVUFBVSxFQUFFLENBQUM7WUFDdkMsSUFBSSxpQkFBaUIsS0FBSyxpQkFBaUIsQ0FBQyxHQUFHLEVBQUUsQ0FBQztnQkFDaEQsT0FBTyxJQUFBLGdCQUFLLEVBQUMsS0FBSyxDQUFDLENBQUM7WUFDdEIsQ0FBQztpQkFBTSxDQUFDO2dCQUNOLE9BQU8sUUFBUSxDQUFDLEtBQUssQ0FBQyxDQUFDLFFBQVEsRUFBRSxDQUFDO1lBQ3BDLENBQUM7UUFDSCxDQUFDO1FBQ0QsT0FBTyxLQUFLLENBQUM7SUFDZixDQUFDLENBQUMsQ0FBQztBQUNMLENBQUMifQ==