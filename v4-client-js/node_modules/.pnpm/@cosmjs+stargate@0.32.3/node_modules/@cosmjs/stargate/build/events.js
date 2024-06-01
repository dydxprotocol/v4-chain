"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.fromTendermintEvent = void 0;
const encoding_1 = require("@cosmjs/encoding");
/**
 * Takes a Tendermint 0.34 or 0.37 event with binary encoded key and value
 * and converts it into an `Event` with string attributes.
 */
function fromTendermintEvent(event) {
    return {
        type: event.type,
        attributes: event.attributes.map((attr) => ({
            key: typeof attr.key == "string" ? attr.key : (0, encoding_1.fromUtf8)(attr.key, true),
            value: typeof attr.value == "string" ? attr.value : (0, encoding_1.fromUtf8)(attr.value, true),
        })),
    };
}
exports.fromTendermintEvent = fromTendermintEvent;
//# sourceMappingURL=events.js.map