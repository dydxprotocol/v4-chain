"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateDecodeTransform = exports.generateEncodeTransform = void 0;
const ts_poet_1 = require("ts-poet");
/** Creates a function to transform a message Source to a Uint8Array Source. */
function generateEncodeTransform(utils, fullName) {
    return (0, ts_poet_1.code) `
    // encodeTransform encodes a source of message objects.
    // Transform<${fullName}, Uint8Array>
    async *encodeTransform(
      source: AsyncIterable<${fullName} | ${fullName}[]> | Iterable<${fullName} | ${fullName}[]>
    ): AsyncIterable<Uint8Array> {
      for await (const pkt of source) {
        if (${utils.globalThis}.Array.isArray(pkt)) {
          for (const p of (pkt as any)) {
            yield* [${fullName}.encode(p).finish()]
          }
        } else {
          yield* [${fullName}.encode(pkt as any).finish()]
        }
      }
    }
  `;
}
exports.generateEncodeTransform = generateEncodeTransform;
/** Creates a function to transform a Uint8Array Source to a message Source. */
function generateDecodeTransform(utils, fullName) {
    return (0, ts_poet_1.code) `
    // decodeTransform decodes a source of encoded messages.
    // Transform<Uint8Array, ${fullName}>
    async *decodeTransform(
      source: AsyncIterable<Uint8Array | Uint8Array[]> | Iterable<Uint8Array | Uint8Array[]>
    ): AsyncIterable<${fullName}> {
      for await (const pkt of source) {
        if (${utils.globalThis}.Array.isArray(pkt)) {
          for (const p of (pkt as any)) {
            yield* [${fullName}.decode(p)]
          }
        } else {
          yield* [${fullName}.decode(pkt as any)]
        }
      }
    }
  `;
}
exports.generateDecodeTransform = generateDecodeTransform;
