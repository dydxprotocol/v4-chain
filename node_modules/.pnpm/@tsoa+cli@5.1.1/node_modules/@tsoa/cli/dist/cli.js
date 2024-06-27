#!/usr/bin/env node
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateSpecAndRoutes = exports.runCLI = exports.validateSpecConfig = void 0;
const path = require("path");
const YAML = require("yamljs");
const yargs = require("yargs");
const helpers_1 = require("yargs/helpers");
const metadataGenerator_1 = require("./metadataGeneration/metadataGenerator");
const generate_routes_1 = require("./module/generate-routes");
const generate_spec_1 = require("./module/generate-spec");
const fs_1 = require("./utils/fs");
const workingDir = process.cwd();
let packageJson;
const getPackageJsonValue = async (key, defaultValue = '') => {
    if (!packageJson) {
        try {
            const packageJsonRaw = await (0, fs_1.fsReadFile)(`${workingDir}/package.json`);
            packageJson = JSON.parse(packageJsonRaw.toString('utf8'));
        }
        catch (err) {
            return defaultValue;
        }
    }
    return packageJson[key] || '';
};
const nameDefault = () => getPackageJsonValue('name', 'TSOA');
const versionDefault = () => getPackageJsonValue('version', '1.0.0');
const descriptionDefault = () => getPackageJsonValue('description', 'Build swagger-compliant REST APIs using TypeScript and Node');
const licenseDefault = () => getPackageJsonValue('license', 'MIT');
const determineNoImplicitAdditionalSetting = (noImplicitAdditionalProperties) => {
    if (noImplicitAdditionalProperties === 'silently-remove-extras' || noImplicitAdditionalProperties === 'throw-on-extras' || noImplicitAdditionalProperties === 'ignore') {
        return noImplicitAdditionalProperties;
    }
    else {
        return 'ignore';
    }
};
const authorInformation = getPackageJsonValue('author', 'unknown');
const getConfig = async (configPath = 'tsoa.json') => {
    var _a;
    let config;
    const ext = path.extname(configPath);
    try {
        if (ext === '.yaml' || ext === '.yml') {
            config = YAML.load(configPath);
        }
        else if (ext === '.js') {
            config = await (_a = `${workingDir}/${configPath}`, Promise.resolve().then(() => require(_a)));
        }
        else {
            const configRaw = await (0, fs_1.fsReadFile)(`${workingDir}/${configPath}`);
            config = JSON.parse(configRaw.toString('utf8'));
        }
    }
    catch (err) {
        if (err.code === 'MODULE_NOT_FOUND') {
            throw Error(`No config file found at '${configPath}'`);
        }
        else if (err.name === 'SyntaxError') {
            // eslint-disable-next-line no-console
            console.error(err);
            const errorType = ext === '.js' ? 'JS' : 'JSON';
            // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
            throw Error(`Invalid ${errorType} syntax in config at '${configPath}': ${err.message}`);
        }
        else {
            // eslint-disable-next-line no-console
            console.error(err);
            // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
            throw Error(`Unhandled error encountered loading '${configPath}': ${err.message}`);
        }
    }
    return config;
};
const resolveConfig = async (config) => {
    return typeof config === 'object' ? config : getConfig(config);
};
const validateCompilerOptions = (config) => {
    return (config || {});
};
const validateSpecConfig = async (config) => {
    if (!config.spec) {
        throw new Error('Missing spec: configuration must contain spec. Spec used to be called swagger in previous versions of tsoa.');
    }
    if (!config.spec.outputDirectory) {
        throw new Error('Missing outputDirectory: configuration must contain output directory.');
    }
    if (!config.entryFile && (!config.controllerPathGlobs || !config.controllerPathGlobs.length)) {
        throw new Error('Missing entryFile and controllerPathGlobs: Configuration must contain an entry point file or controller path globals.');
    }
    if (!!config.entryFile && !(await (0, fs_1.fsExists)(config.entryFile))) {
        throw new Error(`EntryFile not found: ${config.entryFile} - please check your tsoa config.`);
    }
    config.spec.version = config.spec.version || (await versionDefault());
    config.spec.specVersion = config.spec.specVersion || 2;
    if (config.spec.specVersion !== 2 && config.spec.specVersion !== 3) {
        throw new Error('Unsupported Spec version.');
    }
    if (config.spec.spec && !['immediate', 'recursive', 'deepmerge', undefined].includes(config.spec.specMerging)) {
        // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
        throw new Error(`Invalid specMerging config: ${config.spec.specMerging}`);
    }
    const noImplicitAdditionalProperties = determineNoImplicitAdditionalSetting(config.noImplicitAdditionalProperties);
    config.spec.name = config.spec.name || (await nameDefault());
    config.spec.description = config.spec.description || (await descriptionDefault());
    config.spec.license = config.spec.license || (await licenseDefault());
    config.spec.basePath = config.spec.basePath || '/';
    // defaults to template that may generate non-unique operation ids.
    // @see https://github.com/lukeautry/tsoa/issues/1005
    config.spec.operationIdTemplate = config.spec.operationIdTemplate || '{{titleCase method.name}}';
    if (!config.spec.contact) {
        config.spec.contact = {};
    }
    const author = await authorInformation;
    if (typeof author === 'string') {
        const contact = /^([^<(]*)?\s*(?:<([^>(]*)>)?\s*(?:\(([^)]*)\)|$)/m.exec(author);
        config.spec.contact.name = config.spec.contact.name || (contact === null || contact === void 0 ? void 0 : contact[1]);
        config.spec.contact.email = config.spec.contact.email || (contact === null || contact === void 0 ? void 0 : contact[2]);
        config.spec.contact.url = config.spec.contact.url || (contact === null || contact === void 0 ? void 0 : contact[3]);
    }
    else if (typeof author === 'object') {
        config.spec.contact.name = config.spec.contact.name || (author === null || author === void 0 ? void 0 : author.name);
        config.spec.contact.email = config.spec.contact.email || (author === null || author === void 0 ? void 0 : author.email);
        config.spec.contact.url = config.spec.contact.url || (author === null || author === void 0 ? void 0 : author.url);
    }
    if (config.spec.rootSecurity) {
        if (!Array.isArray(config.spec.rootSecurity)) {
            throw new Error('spec.rootSecurity must be an array');
        }
        if (config.spec.rootSecurity) {
            const ok = config.spec.rootSecurity.every(security => typeof security === 'object' && security !== null && Object.values(security).every(scope => Array.isArray(scope)));
            if (!ok) {
                throw new Error('spec.rootSecurity must be an array of objects whose keys are security scheme names and values are arrays of scopes');
            }
        }
    }
    return {
        ...config.spec,
        noImplicitAdditionalProperties,
        entryFile: config.entryFile,
        controllerPathGlobs: config.controllerPathGlobs,
    };
};
exports.validateSpecConfig = validateSpecConfig;
const validateRoutesConfig = async (config) => {
    if (!config.entryFile && (!config.controllerPathGlobs || !config.controllerPathGlobs.length)) {
        throw new Error('Missing entryFile and controllerPathGlobs: Configuration must contain an entry point file or controller path globals.');
    }
    if (!!config.entryFile && !(await (0, fs_1.fsExists)(config.entryFile))) {
        throw new Error(`EntryFile not found: ${config.entryFile} - Please check your tsoa config.`);
    }
    if (!config.routes.routesDir) {
        throw new Error('Missing routesDir: Configuration must contain a routes file output directory.');
    }
    if (config.routes.authenticationModule && !((await (0, fs_1.fsExists)(config.routes.authenticationModule)) || (await (0, fs_1.fsExists)(config.routes.authenticationModule + '.ts')))) {
        throw new Error(`No authenticationModule file found at '${config.routes.authenticationModule}'`);
    }
    if (config.routes.iocModule && !((await (0, fs_1.fsExists)(config.routes.iocModule)) || (await (0, fs_1.fsExists)(config.routes.iocModule + '.ts')))) {
        throw new Error(`No iocModule file found at '${config.routes.iocModule}'`);
    }
    const noImplicitAdditionalProperties = determineNoImplicitAdditionalSetting(config.noImplicitAdditionalProperties);
    config.routes.basePath = config.routes.basePath || '/';
    return {
        ...config.routes,
        entryFile: config.entryFile,
        noImplicitAdditionalProperties,
        controllerPathGlobs: config.controllerPathGlobs,
        multerOpts: config.multerOpts,
        rootSecurity: config.spec.rootSecurity,
    };
};
const configurationArgs = {
    alias: 'c',
    describe: 'tsoa configuration file; default is tsoa.json in the working directory',
    required: false,
    string: true,
};
const hostArgs = {
    describe: 'API host',
    required: false,
    string: true,
};
const basePathArgs = {
    describe: 'Base API path',
    required: false,
    string: true,
};
const yarmlArgs = {
    describe: 'Swagger spec yaml format',
    required: false,
    boolean: true,
};
const jsonArgs = {
    describe: 'Swagger spec json format',
    required: false,
    boolean: true,
};
function runCLI() {
    void yargs((0, helpers_1.hideBin)(process.argv))
        .usage('Usage: $0 <command> [options]')
        .command('spec', 'Generate OpenAPI spec', {
        basePath: basePathArgs,
        configuration: configurationArgs,
        host: hostArgs,
        json: jsonArgs,
        yaml: yarmlArgs,
    }, args => SpecGenerator(args))
        .command('routes', 'Generate routes', {
        basePath: basePathArgs,
        configuration: configurationArgs,
    }, args => routeGenerator(args))
        .command('spec-and-routes', 'Generate OpenAPI spec and routes', {
        basePath: basePathArgs,
        configuration: configurationArgs,
        host: hostArgs,
        json: jsonArgs,
        yaml: yarmlArgs,
    }, args => void generateSpecAndRoutes(args))
        .demandCommand(1, 1, 'Must provide a valid command.')
        .help('help')
        .alias('help', 'h')
        .parse();
}
exports.runCLI = runCLI;
if (require.main === module)
    runCLI();
async function SpecGenerator(args) {
    try {
        const config = await resolveConfig(args.configuration);
        if (args.basePath) {
            config.spec.basePath = args.basePath;
        }
        if (args.host) {
            config.spec.host = args.host;
        }
        if (args.yaml) {
            config.spec.yaml = args.yaml;
        }
        if (args.json) {
            config.spec.yaml = false;
        }
        const compilerOptions = validateCompilerOptions(config.compilerOptions);
        const swaggerConfig = await (0, exports.validateSpecConfig)(config);
        await (0, generate_spec_1.generateSpec)(swaggerConfig, compilerOptions, config.ignore);
    }
    catch (err) {
        // eslint-disable-next-line no-console
        console.error('Generate swagger error.\n', err);
        process.exit(1);
    }
}
async function routeGenerator(args) {
    try {
        const config = await resolveConfig(args.configuration);
        if (args.basePath) {
            config.routes.basePath = args.basePath;
        }
        const compilerOptions = validateCompilerOptions(config.compilerOptions);
        const routesConfig = await validateRoutesConfig(config);
        await (0, generate_routes_1.generateRoutes)(routesConfig, compilerOptions, config.ignore);
    }
    catch (err) {
        // eslint-disable-next-line no-console
        console.error('Generate routes error.\n', err);
        process.exit(1);
    }
}
async function generateSpecAndRoutes(args, metadata) {
    try {
        const config = await resolveConfig(args.configuration);
        if (args.basePath) {
            config.spec.basePath = args.basePath;
        }
        if (args.host) {
            config.spec.host = args.host;
        }
        if (args.yaml) {
            config.spec.yaml = args.yaml;
        }
        if (args.json) {
            config.spec.yaml = false;
        }
        const compilerOptions = validateCompilerOptions(config.compilerOptions);
        const routesConfig = await validateRoutesConfig(config);
        const swaggerConfig = await (0, exports.validateSpecConfig)(config);
        if (!metadata) {
            metadata = new metadataGenerator_1.MetadataGenerator(config.entryFile, compilerOptions, config.ignore, config.controllerPathGlobs, config.spec.rootSecurity).Generate();
        }
        await Promise.all([(0, generate_routes_1.generateRoutes)(routesConfig, compilerOptions, config.ignore, metadata), (0, generate_spec_1.generateSpec)(swaggerConfig, compilerOptions, config.ignore, metadata)]);
        return metadata;
    }
    catch (err) {
        // eslint-disable-next-line no-console
        console.error('Generate routes error.\n', err);
        process.exit(1);
        throw err;
    }
}
exports.generateSpecAndRoutes = generateSpecAndRoutes;
//# sourceMappingURL=cli.js.map