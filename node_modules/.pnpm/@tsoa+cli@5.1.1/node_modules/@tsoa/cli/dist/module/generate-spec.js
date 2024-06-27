"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateSpec = exports.getSwaggerOutputPath = void 0;
const YAML = require("yamljs");
const metadataGenerator_1 = require("../metadataGeneration/metadataGenerator");
const specGenerator2_1 = require("../swagger/specGenerator2");
const specGenerator3_1 = require("../swagger/specGenerator3");
const fs_1 = require("../utils/fs");
const getSwaggerOutputPath = (swaggerConfig) => {
    const ext = swaggerConfig.yaml ? 'yaml' : 'json';
    const specFileBaseName = swaggerConfig.specFileBaseName || 'swagger';
    return `${swaggerConfig.outputDirectory}/${specFileBaseName}.${ext}`;
};
exports.getSwaggerOutputPath = getSwaggerOutputPath;
const generateSpec = async (swaggerConfig, compilerOptions, ignorePaths, 
/**
 * pass in cached metadata returned in a previous step to speed things up
 */
metadata) => {
    if (!metadata) {
        metadata = new metadataGenerator_1.MetadataGenerator(swaggerConfig.entryFile, compilerOptions, ignorePaths, swaggerConfig.controllerPathGlobs, swaggerConfig.rootSecurity).Generate();
    }
    let spec;
    if (swaggerConfig.specVersion && swaggerConfig.specVersion === 3) {
        spec = new specGenerator3_1.SpecGenerator3(metadata, swaggerConfig).GetSpec();
    }
    else {
        spec = new specGenerator2_1.SpecGenerator2(metadata, swaggerConfig).GetSpec();
    }
    await (0, fs_1.fsMkDir)(swaggerConfig.outputDirectory, { recursive: true });
    let data = JSON.stringify(spec, null, '\t');
    if (swaggerConfig.yaml) {
        data = YAML.stringify(JSON.parse(data), 10);
    }
    const outputPath = (0, exports.getSwaggerOutputPath)(swaggerConfig);
    await (0, fs_1.fsWriteFile)(outputPath, data, { encoding: 'utf8' });
    return metadata;
};
exports.generateSpec = generateSpec;
//# sourceMappingURL=generate-spec.js.map