"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getMemberName = exports.generateEnumToNumber = exports.generateEnumToJson = exports.generateEnumFromJson = exports.generateEnum = void 0;
const ts_poet_1 = require("ts-poet");
const utils_1 = require("./utils");
const case_1 = require("./case");
const sourceInfo_1 = require("./sourceInfo");
// Output the `enum { Foo, A = 0, B = 1 }`
function generateEnum(ctx, fullName, enumDesc, sourceInfo) {
    var _a;
    const { options } = ctx;
    const chunks = [];
    let unrecognizedEnum = { present: false };
    (0, utils_1.maybeAddComment)(options, sourceInfo, chunks, (_a = enumDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
    if (options.enumsAsLiterals) {
        chunks.push((0, ts_poet_1.code) `export const ${(0, ts_poet_1.def)(fullName)} = {`);
    }
    else {
        chunks.push((0, ts_poet_1.code) `export ${options.constEnums ? "const " : ""}enum ${(0, ts_poet_1.def)(fullName)} {`);
    }
    const delimiter = options.enumsAsLiterals ? ":" : "=";
    enumDesc.value.forEach((valueDesc, index) => {
        var _a;
        const info = sourceInfo.lookup(sourceInfo_1.Fields.enum.value, index);
        const valueName = getValueName(ctx, fullName, valueDesc);
        const memberName = getMemberName(ctx, enumDesc, valueDesc);
        if (valueDesc.number === options.unrecognizedEnumValue) {
            unrecognizedEnum = { present: true, name: memberName };
        }
        (0, utils_1.maybeAddComment)(options, info, chunks, (_a = valueDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated, `${memberName} - `);
        chunks.push((0, ts_poet_1.code) `${memberName} ${delimiter} ${options.stringEnums ? `"${valueName}"` : valueDesc.number.toString()},`);
    });
    if (options.unrecognizedEnum && !unrecognizedEnum.present) {
        chunks.push((0, ts_poet_1.code) `
      ${options.unrecognizedEnumName} ${delimiter} ${options.stringEnums ? `"${options.unrecognizedEnumName}"` : options.unrecognizedEnumValue.toString()},`);
    }
    if (options.enumsAsLiterals) {
        chunks.push((0, ts_poet_1.code) `} as const`);
        chunks.push((0, ts_poet_1.code) `\n`);
        chunks.push((0, ts_poet_1.code) `export type ${(0, ts_poet_1.def)(fullName)} = typeof ${(0, ts_poet_1.def)(fullName)}[keyof typeof ${(0, ts_poet_1.def)(fullName)}]`);
        chunks.push((0, ts_poet_1.code) `\n`);
        chunks.push((0, ts_poet_1.code) `export namespace ${(0, ts_poet_1.def)(fullName)} {`);
        enumDesc.value.forEach((valueDesc) => {
            const memberName = getMemberName(ctx, enumDesc, valueDesc);
            chunks.push((0, ts_poet_1.code) `export type ${memberName} = typeof ${(0, ts_poet_1.def)(fullName)}.${memberName};`);
        });
        if (options.unrecognizedEnum && !unrecognizedEnum.present) {
            chunks.push((0, ts_poet_1.code) `export type ${options.unrecognizedEnumName} = typeof ${(0, ts_poet_1.def)(fullName)}.${options.unrecognizedEnumName};`);
        }
        chunks.push((0, ts_poet_1.code) `}`);
    }
    else {
        chunks.push((0, ts_poet_1.code) `}`);
    }
    if (options.outputJsonMethods === true ||
        options.outputJsonMethods === "from-only" ||
        (options.stringEnums && options.outputEncodeMethods)) {
        chunks.push((0, ts_poet_1.code) `\n`);
        chunks.push(generateEnumFromJson(ctx, fullName, enumDesc, unrecognizedEnum));
    }
    if (options.outputJsonMethods === true || options.outputJsonMethods === "to-only") {
        chunks.push((0, ts_poet_1.code) `\n`);
        chunks.push(generateEnumToJson(ctx, fullName, enumDesc, unrecognizedEnum));
    }
    if (options.stringEnums && options.outputEncodeMethods) {
        chunks.push((0, ts_poet_1.code) `\n`);
        chunks.push(generateEnumToNumber(ctx, fullName, enumDesc, unrecognizedEnum));
    }
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
exports.generateEnum = generateEnum;
/** Generates a function with a big switch statement to decode JSON -> our enum. */
function generateEnumFromJson(ctx, fullName, enumDesc, unrecognizedEnum) {
    const { options, utils } = ctx;
    const chunks = [];
    const functionName = (0, case_1.uncapitalize)(fullName) + "FromJSON";
    chunks.push((0, ts_poet_1.code) `export function ${(0, ts_poet_1.def)(functionName)}(object: any): ${fullName} {`);
    chunks.push((0, ts_poet_1.code) `switch (object) {`);
    for (const valueDesc of enumDesc.value) {
        const memberName = getMemberName(ctx, enumDesc, valueDesc);
        const valueName = getValueName(ctx, fullName, valueDesc);
        chunks.push((0, ts_poet_1.code) `
      case ${valueDesc.number}:
      case "${valueName}":
        return ${fullName}.${memberName};
    `);
    }
    if (options.unrecognizedEnum) {
        if (!unrecognizedEnum.present) {
            chunks.push((0, ts_poet_1.code) `
        case ${options.unrecognizedEnumValue}:
        case "${options.unrecognizedEnumName}":
        default:
          return ${fullName}.${options.unrecognizedEnumName};
      `);
        }
        else {
            chunks.push((0, ts_poet_1.code) `
        default:
          return ${fullName}.${unrecognizedEnum.name};
      `);
        }
    }
    else {
        // We use globalThis to avoid conflicts on protobuf types named `Error`.
        chunks.push((0, ts_poet_1.code) `
      default:
        throw new ${utils.globalThis}.Error("Unrecognized enum value " + object + " for enum ${fullName}");
    `);
    }
    chunks.push((0, ts_poet_1.code) `}`);
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
exports.generateEnumFromJson = generateEnumFromJson;
/** Generates a function with a big switch statement to encode our enum -> JSON. */
function generateEnumToJson(ctx, fullName, enumDesc, unrecognizedEnum) {
    const { options, utils } = ctx;
    const chunks = [];
    const functionName = (0, case_1.uncapitalize)(fullName) + "ToJSON";
    chunks.push((0, ts_poet_1.code) `export function ${(0, ts_poet_1.def)(functionName)}(object: ${fullName}): ${ctx.options.useNumericEnumForJson ? "number" : "string"} {`);
    chunks.push((0, ts_poet_1.code) `switch (object) {`);
    for (const valueDesc of enumDesc.value) {
        if (ctx.options.useNumericEnumForJson) {
            const memberName = getMemberName(ctx, enumDesc, valueDesc);
            chunks.push((0, ts_poet_1.code) `case ${fullName}.${memberName}: return ${valueDesc.number};`);
        }
        else {
            const memberName = getMemberName(ctx, enumDesc, valueDesc);
            const valueName = getValueName(ctx, fullName, valueDesc);
            chunks.push((0, ts_poet_1.code) `case ${fullName}.${memberName}: return "${valueName}";`);
        }
    }
    if (options.unrecognizedEnum) {
        if (!unrecognizedEnum.present) {
            chunks.push((0, ts_poet_1.code) `
        case ${fullName}.${options.unrecognizedEnumName}:`);
            if (ctx.options.useNumericEnumForJson) {
                chunks.push((0, ts_poet_1.code) `
        default:
          return ${options.unrecognizedEnumValue};
      `);
            }
            else {
                chunks.push((0, ts_poet_1.code) `
        default:
          return "${options.unrecognizedEnumName}";
      `);
            }
        }
        else if (ctx.options.useNumericEnumForJson) {
            chunks.push((0, ts_poet_1.code) `
        default:
          return ${options.unrecognizedEnumValue};
      `);
        }
        else {
            chunks.push((0, ts_poet_1.code) `
      default:
        return "${unrecognizedEnum.name}";
    `);
        }
    }
    else {
        // We use globalThis to avoid conflicts on protobuf types named `Error`.
        chunks.push((0, ts_poet_1.code) `
      default:
        throw new ${utils.globalThis}.Error("Unrecognized enum value " + object + " for enum ${fullName}");
    `);
    }
    chunks.push((0, ts_poet_1.code) `}`);
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
exports.generateEnumToJson = generateEnumToJson;
/** Generates a function with a big switch statement to encode our string enum -> int value. */
function generateEnumToNumber(ctx, fullName, enumDesc, unrecognizedEnum) {
    const { options, utils } = ctx;
    const chunks = [];
    const functionName = (0, case_1.uncapitalize)(fullName) + "ToNumber";
    chunks.push((0, ts_poet_1.code) `export function ${(0, ts_poet_1.def)(functionName)}(object: ${fullName}): number {`);
    chunks.push((0, ts_poet_1.code) `switch (object) {`);
    for (const valueDesc of enumDesc.value) {
        chunks.push((0, ts_poet_1.code) `case ${fullName}.${getMemberName(ctx, enumDesc, valueDesc)}: return ${valueDesc.number};`);
    }
    if (options.unrecognizedEnum) {
        if (!unrecognizedEnum.present) {
            chunks.push((0, ts_poet_1.code) `
        case ${fullName}.${options.unrecognizedEnumName}:
        default:
          return ${options.unrecognizedEnumValue};
      `);
        }
        else {
            chunks.push((0, ts_poet_1.code) `
        default:
          return ${options.unrecognizedEnumValue};
      `);
        }
    }
    else {
        // We use globalThis to avoid conflicts on protobuf types named `Error`.
        chunks.push((0, ts_poet_1.code) `
      default:
        throw new ${utils.globalThis}.Error("Unrecognized enum value " + object + " for enum ${fullName}");
    `);
    }
    chunks.push((0, ts_poet_1.code) `}`);
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
exports.generateEnumToNumber = generateEnumToNumber;
function getMemberName(ctx, enumDesc, valueDesc) {
    if (ctx.options.removeEnumPrefix) {
        return valueDesc.name.replace(`${(0, case_1.camelToSnake)(enumDesc.name)}_`, "");
    }
    return valueDesc.name;
}
exports.getMemberName = getMemberName;
function getValueName(ctx, fullName, valueDesc) {
    return valueDesc.name;
}
