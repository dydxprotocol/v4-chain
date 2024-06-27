"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MetadataGenerator = void 0;
const mm = require("minimatch");
const ts = require("typescript");
const importClassesFromDirectories_1 = require("../utils/importClassesFromDirectories");
const controllerGenerator_1 = require("./controllerGenerator");
const exceptions_1 = require("./exceptions");
const typeResolver_1 = require("./typeResolver");
const decoratorUtils_1 = require("../utils/decoratorUtils");
class MetadataGenerator {
    constructor(entryFile, compilerOptions, ignorePaths, controllers, rootSecurity = []) {
        this.compilerOptions = compilerOptions;
        this.ignorePaths = ignorePaths;
        this.rootSecurity = rootSecurity;
        this.controllerNodes = new Array();
        this.referenceTypeMap = {};
        this.circularDependencyResolvers = new Array();
        this.checkForMethodSignatureDuplicates = (controllers) => {
            const map = {};
            controllers.forEach(controller => {
                controller.methods.forEach(method => {
                    const signature = method.path ? `@${method.method}(${controller.path}/${method.path})` : `@${method.method}(${controller.path})`;
                    const methodDescription = `${controller.name}#${method.name}`;
                    if (map[signature]) {
                        map[signature].push(methodDescription);
                    }
                    else {
                        map[signature] = [methodDescription];
                    }
                });
            });
            let message = '';
            Object.keys(map).forEach(signature => {
                const controllers = map[signature];
                if (controllers.length > 1) {
                    message += `Duplicate method signature ${signature} found in controllers: ${controllers.join(', ')}\n`;
                }
            });
            if (message) {
                throw new exceptions_1.GenerateMetadataError(message);
            }
        };
        this.checkForPathParamSignatureDuplicates = (controllers) => {
            const paramRegExp = new RegExp('{(\\w*)}|:(\\w+)', 'g');
            let PathDuplicationType;
            (function (PathDuplicationType) {
                PathDuplicationType[PathDuplicationType["FULL"] = 0] = "FULL";
                PathDuplicationType[PathDuplicationType["PARTIAL"] = 1] = "PARTIAL";
            })(PathDuplicationType || (PathDuplicationType = {}));
            const collisions = [];
            function addCollision(type, method, controller, collidesWith) {
                let existingCollision = collisions.find(collision => collision.type === type && collision.method === method && collision.controller === controller);
                if (!existingCollision) {
                    existingCollision = {
                        type,
                        method,
                        controller,
                        collidesWith: [],
                    };
                    collisions.push(existingCollision);
                }
                existingCollision.collidesWith.push(collidesWith);
            }
            controllers.forEach(controller => {
                const methodRouteGroup = {};
                // Group all ts methods with HTTP method decorator into same object in same controller.
                controller.methods.forEach(method => {
                    if (methodRouteGroup[method.method] === undefined) {
                        methodRouteGroup[method.method] = [];
                    }
                    const params = method.path.match(paramRegExp);
                    methodRouteGroup[method.method].push({
                        method,
                        path: (params === null || params === void 0 ? void 0 : params.reduce((s, a) => {
                            // replace all params with {} placeholder for comparison
                            return s.replace(a, '{}');
                        }, method.path)) || method.path,
                    });
                });
                Object.keys(methodRouteGroup).forEach((key) => {
                    const methodRoutes = methodRouteGroup[key];
                    // check each route with the routes that are defined before it
                    for (let i = 0; i < methodRoutes.length; i += 1) {
                        for (let j = 0; j < i; j += 1) {
                            if (methodRoutes[i].path === methodRoutes[j].path) {
                                // full match
                                addCollision(PathDuplicationType.FULL, methodRoutes[i].method, controller, methodRoutes[j].method);
                            }
                            else if (methodRoutes[i].path.split('/').length === methodRoutes[j].path.split('/').length &&
                                methodRoutes[j].path
                                    .substr(methodRoutes[j].path.lastIndexOf('/')) // compare only the "last" part of the path
                                    .split('/')
                                    .some(v => !!v) && // ensure the comparison path has a value
                                methodRoutes[i].path.split('/').every((v, index) => {
                                    const comparisonPathPart = methodRoutes[j].path.split('/')[index];
                                    // if no params, compare values
                                    if (!v.includes('{}')) {
                                        return v === comparisonPathPart;
                                    }
                                    // otherwise check if route starts with comparison route
                                    return v.startsWith(methodRoutes[j].path.split('/')[index]);
                                })) {
                                // partial match - reorder routes!
                                addCollision(PathDuplicationType.PARTIAL, methodRoutes[i].method, controller, methodRoutes[j].method);
                            }
                        }
                    }
                });
            });
            // print warnings for each collision (grouped by route)
            collisions.forEach(collision => {
                let message = '';
                if (collision.type === PathDuplicationType.FULL) {
                    message = `Duplicate path parameter definition signature found in controller `;
                }
                else if (collision.type === PathDuplicationType.PARTIAL) {
                    message = `Overlapping path parameter definition signature found in controller `;
                }
                message += collision.controller.name;
                message += ` [ method ${collision.method.method.toUpperCase()} ${collision.method.name} route: ${collision.method.path} ] collides with `;
                message += collision.collidesWith
                    .map((method) => {
                    return `[ method ${method.method.toUpperCase()} ${method.name} route: ${method.path} ]`;
                })
                    .join(', ');
                message += '\n';
                console.warn(message);
            });
        };
        typeResolver_1.TypeResolver.clearCache();
        this.program = controllers ? this.setProgramToDynamicControllersFiles(controllers) : ts.createProgram([entryFile], compilerOptions || {});
        this.typeChecker = this.program.getTypeChecker();
    }
    Generate() {
        this.extractNodeFromProgramSourceFiles();
        const controllers = this.buildControllers();
        this.checkForMethodSignatureDuplicates(controllers);
        this.checkForPathParamSignatureDuplicates(controllers);
        this.circularDependencyResolvers.forEach(c => c(this.referenceTypeMap));
        return {
            controllers,
            referenceTypeMap: this.referenceTypeMap,
        };
    }
    setProgramToDynamicControllersFiles(controllers) {
        const allGlobFiles = (0, importClassesFromDirectories_1.importClassesFromDirectories)(controllers);
        if (allGlobFiles.length === 0) {
            throw new exceptions_1.GenerateMetadataError(`[${controllers.join(', ')}] globs found 0 controllers.`);
        }
        return ts.createProgram(allGlobFiles, this.compilerOptions || {});
    }
    extractNodeFromProgramSourceFiles() {
        this.program.getSourceFiles().forEach(sf => {
            if (this.ignorePaths && this.ignorePaths.length) {
                for (const path of this.ignorePaths) {
                    if (mm(sf.fileName, path)) {
                        return;
                    }
                }
            }
            ts.forEachChild(sf, node => {
                if (ts.isClassDeclaration(node) && (0, decoratorUtils_1.getDecorators)(node, identifier => identifier.text === 'Route').length) {
                    this.controllerNodes.push(node);
                }
            });
        });
    }
    TypeChecker() {
        return this.typeChecker;
    }
    AddReferenceType(referenceType) {
        if (!referenceType.refName) {
            return;
        }
        this.referenceTypeMap[referenceType.refName] = referenceType;
    }
    GetReferenceType(refName) {
        return this.referenceTypeMap[refName];
    }
    OnFinish(callback) {
        this.circularDependencyResolvers.push(callback);
    }
    buildControllers() {
        if (this.controllerNodes.length === 0) {
            throw new Error('no controllers found, check tsoa configuration');
        }
        return this.controllerNodes
            .map(classDeclaration => new controllerGenerator_1.ControllerGenerator(classDeclaration, this, this.rootSecurity))
            .filter(generator => generator.IsValid())
            .map(generator => generator.Generate());
    }
}
exports.MetadataGenerator = MetadataGenerator;
//# sourceMappingURL=metadataGenerator.js.map