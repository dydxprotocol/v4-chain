"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ControllerGenerator = void 0;
const ts = require("typescript");
const decoratorUtils_1 = require("./../utils/decoratorUtils");
const exceptions_1 = require("./exceptions");
const methodGenerator_1 = require("./methodGenerator");
const typeResolver_1 = require("./typeResolver");
const headerTypeHelpers_1 = require("../utils/headerTypeHelpers");
class ControllerGenerator {
    constructor(node, current, parentSecurity = []) {
        this.node = node;
        this.current = current;
        this.parentSecurity = parentSecurity;
        this.path = this.getPath();
        this.tags = this.getTags();
        this.security = this.getSecurity();
        this.isHidden = this.getIsHidden();
        this.commonResponses = this.getCommonResponses();
        this.produces = this.getProduces();
    }
    IsValid() {
        return !!this.path || this.path === '';
    }
    Generate() {
        if (!this.node.parent) {
            throw new exceptions_1.GenerateMetadataError("Controller node doesn't have a valid parent source file.");
        }
        if (!this.node.name) {
            throw new exceptions_1.GenerateMetadataError("Controller node doesn't have a valid name.");
        }
        const sourceFile = this.node.parent.getSourceFile();
        return {
            location: sourceFile.fileName,
            methods: this.buildMethods(),
            name: this.node.name.text,
            path: this.path || '',
            produces: this.produces,
        };
    }
    buildMethods() {
        return this.node.members
            .filter(ts.isMethodDeclaration)
            .map(m => new methodGenerator_1.MethodGenerator(m, this.current, this.commonResponses, this.path, this.tags, this.security, this.isHidden))
            .filter(generator => generator.IsValid())
            .map(generator => generator.Generate());
    }
    getPath() {
        const decorators = (0, decoratorUtils_1.getDecorators)(this.node, identifier => identifier.text === 'Route');
        if (!decorators || !decorators.length) {
            return;
        }
        if (decorators.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one Route decorator allowed in '${this.node.name.text}' class.`);
        }
        const decorator = decorators[0];
        const expression = decorator.parent;
        const decoratorArgument = expression.arguments[0];
        return decoratorArgument ? `${decoratorArgument.text}` : '';
    }
    getCommonResponses() {
        const decorators = (0, decoratorUtils_1.getDecorators)(this.node, identifier => identifier.text === 'Response');
        if (!decorators || !decorators.length) {
            return [];
        }
        return decorators.map(decorator => {
            const expression = decorator.parent;
            const [name, description, example] = (0, decoratorUtils_1.getDecoratorValues)(decorator, this.current.typeChecker);
            if (!name) {
                throw new exceptions_1.GenerateMetadataError(`Controller's responses should have an explicit name.`);
            }
            return {
                description: description || '',
                examples: example === undefined ? undefined : [example],
                name,
                schema: expression.typeArguments && expression.typeArguments.length > 0 && !this.isHidden ? new typeResolver_1.TypeResolver(expression.typeArguments[0], this.current).resolve() : undefined,
                headers: (0, headerTypeHelpers_1.getHeaderType)(expression.typeArguments, 1, this.current),
            };
        });
    }
    getTags() {
        const decorators = (0, decoratorUtils_1.getDecorators)(this.node, identifier => identifier.text === 'Tags');
        if (!decorators || !decorators.length) {
            return;
        }
        if (decorators.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one Tags decorator allowed in '${this.node.name.text}' class.`);
        }
        const decorator = decorators[0];
        const expression = decorator.parent;
        return expression.arguments.map((a) => a.text);
    }
    getSecurity() {
        const noSecurityDecorators = (0, decoratorUtils_1.getDecorators)(this.node, identifier => identifier.text === 'NoSecurity');
        const securityDecorators = (0, decoratorUtils_1.getDecorators)(this.node, identifier => identifier.text === 'Security');
        if ((noSecurityDecorators === null || noSecurityDecorators === void 0 ? void 0 : noSecurityDecorators.length) && (securityDecorators === null || securityDecorators === void 0 ? void 0 : securityDecorators.length)) {
            throw new exceptions_1.GenerateMetadataError(`NoSecurity decorator cannot be used in conjunction with Security decorator in '${this.node.name.text}' class.`);
        }
        if (noSecurityDecorators === null || noSecurityDecorators === void 0 ? void 0 : noSecurityDecorators.length) {
            return [];
        }
        if (!securityDecorators || !securityDecorators.length) {
            return this.parentSecurity;
        }
        return securityDecorators.map(d => (0, decoratorUtils_1.getSecurites)(d, this.current.typeChecker));
    }
    getIsHidden() {
        const hiddenDecorators = (0, decoratorUtils_1.getDecorators)(this.node, identifier => identifier.text === 'Hidden');
        if (!hiddenDecorators || !hiddenDecorators.length) {
            return false;
        }
        if (hiddenDecorators.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Only one Hidden decorator allowed in '${this.node.name.text}' class.`);
        }
        return true;
    }
    getProduces() {
        const produces = (0, decoratorUtils_1.getProduces)(this.node, this.current.typeChecker);
        return produces.length ? produces : undefined;
    }
}
exports.ControllerGenerator = ControllerGenerator;
//# sourceMappingURL=controllerGenerator.js.map