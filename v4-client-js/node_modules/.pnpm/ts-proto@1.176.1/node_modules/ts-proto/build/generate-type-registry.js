"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateTypeRegistry = void 0;
const ts_poet_1 = require("ts-poet");
const utils_1 = require("./utils");
const options_1 = require("./options");
function generateTypeRegistry(ctx) {
    const chunks = [];
    chunks.push(generateMessageType(ctx));
    if ((0, options_1.addTypeToMessages)(ctx.options)) {
        chunks.push((0, ts_poet_1.code) `
    export type UnknownMessage = {$type: string};
  `);
    }
    else {
        chunks.push((0, ts_poet_1.code) `
    export type UnknownMessage = unknown;
  `);
    }
    chunks.push((0, ts_poet_1.code) `
    export const messageTypeRegistry = new Map<string, MessageType>();
  `);
    chunks.push((0, ts_poet_1.code) ` ${ctx.utils.Builtin.ifUsed} ${ctx.utils.DeepPartial.ifUsed}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n\n" });
}
exports.generateTypeRegistry = generateTypeRegistry;
function generateMessageType(ctx) {
    const chunks = [];
    chunks.push((0, ts_poet_1.code) `export interface MessageType<Message extends UnknownMessage = UnknownMessage> {`);
    if ((0, options_1.addTypeToMessages)(ctx.options)) {
        chunks.push((0, ts_poet_1.code) `$type: Message['$type'];`);
    }
    else {
        chunks.push((0, ts_poet_1.code) `$type: string;`);
    }
    if (ctx.options.outputEncodeMethods) {
        const Writer = (0, utils_1.impFile)(ctx.options, "Writer@protobufjs/minimal");
        const Reader = (0, utils_1.impFile)(ctx.options, "Reader@protobufjs/minimal");
        chunks.push((0, ts_poet_1.code) `encode(message: Message, writer?: ${Writer}): ${Writer};`);
        chunks.push((0, ts_poet_1.code) `decode(input: ${Reader} | Uint8Array, length?: number): Message;`);
    }
    if (ctx.options.outputJsonMethods) {
        chunks.push((0, ts_poet_1.code) `fromJSON(object: any): Message;`);
        chunks.push((0, ts_poet_1.code) `toJSON(message: Message): unknown;`);
    }
    if (ctx.options.outputPartialMethods) {
        chunks.push((0, ts_poet_1.code) `fromPartial(object: ${ctx.utils.DeepPartial}<Message>): Message;`);
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
