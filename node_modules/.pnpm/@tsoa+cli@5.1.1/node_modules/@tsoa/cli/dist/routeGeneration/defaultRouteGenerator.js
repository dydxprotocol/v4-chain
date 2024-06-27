"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DefaultRouteGenerator = void 0;
const fs = require("fs");
const handlebars = require("handlebars");
const path = require("path");
const runtime_1 = require("@tsoa/runtime");
const fs_1 = require("../utils/fs");
const pathUtils_1 = require("../utils/pathUtils");
const routeGenerator_1 = require("./routeGenerator");
class DefaultRouteGenerator extends routeGenerator_1.AbstractRouteGenerator {
    constructor(metadata, options) {
        super(metadata, options);
        this.pathTransformerFn = pathUtils_1.convertBracesPathParams;
        switch (options.middleware) {
            case 'express':
                this.template = path.join(__dirname, '..', 'routeGeneration/templates/express.hbs');
                break;
            case 'hapi':
                this.template = path.join(__dirname, '..', 'routeGeneration/templates/hapi.hbs');
                this.pathTransformerFn = (path) => path;
                break;
            case 'koa':
                this.template = path.join(__dirname, '..', 'routeGeneration/templates/koa.hbs');
                break;
            default:
                this.template = path.join(__dirname, '..', 'routeGeneration/templates/express.hbs');
        }
        if (options.middlewareTemplate) {
            this.template = options.middlewareTemplate;
        }
    }
    async GenerateCustomRoutes() {
        const data = await (0, fs_1.fsReadFile)(path.join(this.template));
        const file = data.toString();
        return await this.GenerateRoutes(file);
    }
    async GenerateRoutes(middlewareTemplate) {
        if (!fs.lstatSync(this.options.routesDir).isDirectory()) {
            throw new Error(`routesDir should be a directory`);
        }
        else if (this.options.routesFileName !== undefined && !this.options.routesFileName.endsWith('.ts')) {
            throw new Error(`routesFileName should have a '.ts' extension`);
        }
        const fileName = `${this.options.routesDir}/${this.options.routesFileName || 'routes.ts'}`;
        const content = this.buildContent(middlewareTemplate);
        if (await this.shouldWriteFile(fileName, content)) {
            await (0, fs_1.fsWriteFile)(fileName, content);
        }
    }
    pathTransformer(path) {
        return this.pathTransformerFn(path);
    }
    buildContent(middlewareTemplate) {
        handlebars.registerHelper('json', (context) => {
            return JSON.stringify(context);
        });
        const additionalPropsHelper = (additionalProperties) => {
            if (additionalProperties) {
                // Then the model for this type explicitly allows additional properties and thus we should assign that
                return JSON.stringify(additionalProperties);
            }
            else if (this.options.noImplicitAdditionalProperties === 'silently-remove-extras') {
                return JSON.stringify(false);
            }
            else if (this.options.noImplicitAdditionalProperties === 'throw-on-extras') {
                return JSON.stringify(false);
            }
            else if (this.options.noImplicitAdditionalProperties === 'ignore') {
                return JSON.stringify(true);
            }
            else {
                return (0, runtime_1.assertNever)(this.options.noImplicitAdditionalProperties);
            }
        };
        handlebars.registerHelper('additionalPropsHelper', additionalPropsHelper);
        const routesTemplate = handlebars.compile(middlewareTemplate, { noEscape: true });
        return routesTemplate(this.buildContext());
    }
}
exports.DefaultRouteGenerator = DefaultRouteGenerator;
//# sourceMappingURL=defaultRouteGenerator.js.map