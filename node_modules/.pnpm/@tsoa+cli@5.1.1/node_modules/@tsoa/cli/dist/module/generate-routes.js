"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateRoutes = void 0;
const metadataGenerator_1 = require("../metadataGeneration/metadataGenerator");
const defaultRouteGenerator_1 = require("../routeGeneration/defaultRouteGenerator");
const fs_1 = require("../utils/fs");
const path = require("path");
async function generateRoutes(routesConfig, compilerOptions, ignorePaths, 
/**
 * pass in cached metadata returned in a previous step to speed things up
 */
metadata) {
    if (!metadata) {
        metadata = new metadataGenerator_1.MetadataGenerator(routesConfig.entryFile, compilerOptions, ignorePaths, routesConfig.controllerPathGlobs, routesConfig.rootSecurity).Generate();
    }
    const routeGenerator = await getRouteGenerator(metadata, routesConfig);
    await (0, fs_1.fsMkDir)(routesConfig.routesDir, { recursive: true });
    await routeGenerator.GenerateCustomRoutes();
    return metadata;
}
exports.generateRoutes = generateRoutes;
async function getRouteGenerator(metadata, routesConfig) {
    var _a, _b;
    // default route generator for express/koa/hapi
    // custom route generator
    const routeGenerator = routesConfig.routeGenerator;
    if (routeGenerator !== undefined) {
        if (typeof routeGenerator === 'string') {
            try {
                // try as a module import
                const module = await (_a = routeGenerator, Promise.resolve().then(() => require(_a)));
                return new module.default(metadata, routesConfig);
            }
            catch (_err) {
                // try to find a relative import path
                const relativePath = path.relative(__dirname, routeGenerator);
                const module = await (_b = relativePath, Promise.resolve().then(() => require(_b)));
                return new module.default(metadata, routesConfig);
            }
        }
        else {
            return new routeGenerator(metadata, routesConfig);
        }
    }
    if (routesConfig.middleware !== undefined || routesConfig.middlewareTemplate !== undefined) {
        return new defaultRouteGenerator_1.DefaultRouteGenerator(metadata, routesConfig);
    }
    else {
        routesConfig.middleware = 'express';
        return new defaultRouteGenerator_1.DefaultRouteGenerator(metadata, routesConfig);
    }
}
//# sourceMappingURL=generate-routes.js.map