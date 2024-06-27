"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.decodeCosmosSdkDecFromProto = exports.longify = exports.createProtobufRpcClient = exports.createPagination = exports.toAccAddress = void 0;
const encoding_1 = require("@cosmjs/encoding");
const math_1 = require("@cosmjs/math");
const pagination_1 = require("cosmjs-types/cosmos/base/query/v1beta1/pagination");
/**
 * Takes a bech32 encoded address and returns the data part. The prefix is ignored and discarded.
 * This is called AccAddress in Cosmos SDK, which is basically an alias for raw binary data.
 * The result is typically 20 bytes long but not restricted to that.
 */
function toAccAddress(address) {
    return (0, encoding_1.fromBech32)(address).data;
}
exports.toAccAddress = toAccAddress;
/**
 * If paginationKey is set, return a `PageRequest` with the given key.
 * If paginationKey is unset, return `undefined`.
 *
 * Use this with a query response's pagination next key to
 * request the next page.
 */
function createPagination(paginationKey) {
    return paginationKey ? pagination_1.PageRequest.fromPartial({ key: paginationKey }) : pagination_1.PageRequest.fromPartial({});
}
exports.createPagination = createPagination;
function createProtobufRpcClient(base) {
    return {
        request: async (service, method, data) => {
            const path = `/${service}/${method}`;
            const response = await base.queryAbci(path, data, undefined);
            return response.value;
        },
    };
}
exports.createProtobufRpcClient = createProtobufRpcClient;
/**
 * Takes a uint64 value as string, number, BigInt or Uint64 and returns a BigInt
 * of it.
 */
function longify(value) {
    const checkedValue = math_1.Uint64.fromString(value.toString());
    return BigInt(checkedValue.toString());
}
exports.longify = longify;
/**
 * Takes a string or binary encoded `github.com/cosmos/cosmos-sdk/types.Dec` from the
 * protobuf API and converts it into a `Decimal` with 18 fractional digits.
 *
 * See https://github.com/cosmos/cosmos-sdk/issues/10863 for more context why this is needed.
 */
function decodeCosmosSdkDecFromProto(input) {
    const asString = typeof input === "string" ? input : (0, encoding_1.fromAscii)(input);
    return math_1.Decimal.fromAtomics(asString, 18);
}
exports.decodeCosmosSdkDecFromProto = decodeCosmosSdkDecFromProto;
//# sourceMappingURL=utils.js.map