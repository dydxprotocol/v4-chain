"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.assertNever = void 0;
/**
 * This function does exhaustiveness checking to ensure that you have discriminated a union so that no type remains. Use this to get the typescript compiler to help discover cases that were not considered.
 */
function assertNever(value) {
    throw new Error(`Unhandled discriminated union member: ${JSON.stringify(value)}`);
}
exports.assertNever = assertNever;
//# sourceMappingURL=assertNever.js.map