"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateGenericServiceDefinition = void 0;
const ts_poet_1 = require("ts-poet");
const ts_proto_descriptors_1 = require("ts-proto-descriptors");
const case_1 = require("./case");
const sourceInfo_1 = require("./sourceInfo");
const types_1 = require("./types");
const utils_1 = require("./utils");
/**
 * Generates a framework-agnostic service descriptor.
 */
function generateGenericServiceDefinition(ctx, fileDesc, sourceInfo, serviceDesc) {
    var _a, _b, _c;
    const chunks = [];
    (0, utils_1.maybeAddComment)(ctx.options, sourceInfo, chunks, (_a = serviceDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
    // Service definition type
    const name = (0, ts_poet_1.def)(`${serviceDesc.name}Definition`);
    chunks.push((0, ts_poet_1.code) `
    export type ${name} = typeof ${name};
  `);
    // Service definition
    chunks.push((0, ts_poet_1.code) `
    export const ${name} = {
  `);
    (_b = serviceDesc.options) === null || _b === void 0 ? void 0 : _b.uninterpretedOption;
    chunks.push((0, ts_poet_1.code) `
      name: '${serviceDesc.name}',
      fullName: '${(0, utils_1.maybePrefixPackage)(fileDesc, serviceDesc.name)}',
      methods: {
  `);
    for (const [index, methodDesc] of serviceDesc.method.entries()) {
        const info = sourceInfo.lookup(sourceInfo_1.Fields.service.method, index);
        (0, utils_1.maybeAddComment)(ctx.options, info, chunks, (_c = methodDesc.options) === null || _c === void 0 ? void 0 : _c.deprecated);
        chunks.push((0, ts_poet_1.code) `
      ${(0, case_1.uncapitalize)(methodDesc.name)}: ${generateMethodDefinition(ctx, methodDesc)},
    `);
    }
    chunks.push((0, ts_poet_1.code) `
      },
    } as const;
  `);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
exports.generateGenericServiceDefinition = generateGenericServiceDefinition;
function generateMethodDefinition(ctx, methodDesc) {
    const inputType = (0, types_1.messageToTypeName)(ctx, methodDesc.inputType, { keepValueType: true });
    const outputType = (0, types_1.messageToTypeName)(ctx, methodDesc.outputType, { keepValueType: true });
    return (0, ts_poet_1.code) `
    {
      name: '${methodDesc.name}',
      requestType: ${inputType},
      requestStream: ${methodDesc.clientStreaming},
      responseType: ${outputType},
      responseStream: ${methodDesc.serverStreaming},
      options: ${generateMethodOptions(ctx, methodDesc.options)}
    }
  `;
}
function generateMethodOptions(ctx, options) {
    const chunks = [];
    chunks.push((0, ts_poet_1.code) `{`);
    if (options != null) {
        if (options.idempotencyLevel === ts_proto_descriptors_1.MethodOptions_IdempotencyLevel.IDEMPOTENT) {
            chunks.push((0, ts_poet_1.code) `idempotencyLevel: 'IDEMPOTENT',`);
        }
        else if (options.idempotencyLevel === ts_proto_descriptors_1.MethodOptions_IdempotencyLevel.NO_SIDE_EFFECTS) {
            chunks.push((0, ts_poet_1.code) `idempotencyLevel: 'NO_SIDE_EFFECTS',`);
        }
        if (options._unknownFields !== undefined) {
            const unknownFieldsChunks = [];
            unknownFieldsChunks.push((0, ts_poet_1.code) `{`);
            for (const key in options._unknownFields) {
                const values = options._unknownFields[key];
                const valuesChunks = [];
                for (const value of values) {
                    valuesChunks.push((0, ts_poet_1.code) `${ctx.options.env == "node" ? "Buffer.from" : "new Uint8Array"}([${value.join(", ")}])`);
                }
                unknownFieldsChunks.push((0, ts_poet_1.code) `${key}: [\n${(0, ts_poet_1.joinCode)(valuesChunks, { on: "," })}\n],`);
            }
            unknownFieldsChunks.push((0, ts_poet_1.code) `}`);
            chunks.push((0, ts_poet_1.code) `_unknownFields: ${(0, ts_poet_1.joinCode)(unknownFieldsChunks, { on: "\n" })}`);
        }
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
