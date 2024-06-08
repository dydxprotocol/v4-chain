"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const ts_proto_descriptors_1 = require("ts-proto-descriptors");
const util_1 = require("util");
const utils_1 = require("./utils");
const main_1 = require("./main");
const types_1 = require("./types");
const context_1 = require("./context");
const options_1 = require("./options");
const generate_type_registry_1 = require("./generate-type-registry");
// this would be the plugin called by the protoc compiler
async function main() {
    const stdin = await (0, utils_1.readToBuffer)(process.stdin);
    // const json = JSON.parse(stdin.toString());
    // const request = CodeGeneratorRequest.fromObject(json);
    const request = ts_proto_descriptors_1.CodeGeneratorRequest.decode(stdin);
    const { protocVersion, tsProtoVersion } = await (0, utils_1.getVersions)(request);
    const options = (0, options_1.optionsFromParameter)(request.parameter);
    const typeMap = (0, types_1.createTypeMap)(request, options);
    const utils = (0, main_1.makeUtils)(options);
    const ctx = { typeMap, options, utils };
    let filesToGenerate;
    if (options.emitImportedFiles) {
        const fileSet = new Set();
        function addFilesUnlessAliased(filenames) {
            filenames
                .filter((name) => !options.M[name])
                .forEach((name) => {
                if (fileSet.has(name))
                    return;
                fileSet.add(name);
                const file = request.protoFile.find((file) => file.name === name);
                if (file && file.dependency.length > 0) {
                    addFilesUnlessAliased(file.dependency);
                }
            });
        }
        addFilesUnlessAliased(request.fileToGenerate);
        filesToGenerate = request.protoFile.filter((file) => fileSet.has(file.name));
    }
    else {
        filesToGenerate = (0, utils_1.protoFilesToGenerate)(request).filter((file) => !options.M[file.name]);
    }
    const files = await Promise.all(filesToGenerate.map(async (file) => {
        const [path, code] = (0, main_1.generateFile)({ ...ctx, currentFile: (0, context_1.createFileContext)(file) }, file);
        const content = code.toString({ ...(0, options_1.getTsPoetOpts)(options, tsProtoVersion, protocVersion, file.name), path });
        return { name: path, content };
    }));
    if (options.outputTypeRegistry) {
        const utils = (0, main_1.makeUtils)(options);
        const ctx = { options, typeMap, utils };
        const path = "typeRegistry.ts";
        const code = (0, generate_type_registry_1.generateTypeRegistry)(ctx);
        const content = code.toString({ ...(0, options_1.getTsPoetOpts)(options, tsProtoVersion, protocVersion), path });
        files.push({ name: path, content });
    }
    if (options.outputIndex) {
        for (const [path, code] of (0, utils_1.generateIndexFiles)(filesToGenerate, options)) {
            const content = code.toString({ ...(0, options_1.getTsPoetOpts)(options, tsProtoVersion, protocVersion), path });
            files.push({ name: path, content });
        }
    }
    const response = ts_proto_descriptors_1.CodeGeneratorResponse.fromPartial({
        file: files,
        supportedFeatures: ts_proto_descriptors_1.CodeGeneratorResponse_Feature.FEATURE_PROTO3_OPTIONAL,
    });
    const buffer = ts_proto_descriptors_1.CodeGeneratorResponse.encode(response).finish();
    const write = (0, util_1.promisify)(process.stdout.write).bind(process.stdout);
    await write(Buffer.from(buffer));
}
main()
    .then(() => {
    process.exit(0);
})
    .catch((e) => {
    process.stderr.write("FAILED!");
    process.stderr.write(e.message);
    process.stderr.write(e.stack);
    process.exit(1);
});
