"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TypeResolver = void 0;
const ts = require("typescript");
const jsDocUtils_1 = require("./../utils/jsDocUtils");
const decoratorUtils_1 = require("./../utils/decoratorUtils");
const validatorUtils_1 = require("./../utils/validatorUtils");
const exceptions_1 = require("./exceptions");
const initializer_value_1 = require("./initializer-value");
const runtime_1 = require("@tsoa/runtime");
const extension_1 = require("./extension");
const jsonUtils_1 = require("../utils/jsonUtils");
const localReferenceTypeCache = {};
const inProgressTypes = {};
class TypeResolver {
    constructor(typeNode, current, parentNode, context = {}, referencer, addToRefTypeMap = true) {
        this.typeNode = typeNode;
        this.current = current;
        this.parentNode = parentNode;
        this.context = context;
        this.referencer = referencer;
        this.addToRefTypeMap = addToRefTypeMap;
        this.attemptToResolveKindToPrimitive = (syntaxKind) => {
            if (syntaxKind === ts.SyntaxKind.NumberKeyword) {
                return {
                    foundMatch: true,
                    resolvedType: 'number',
                };
            }
            else if (syntaxKind === ts.SyntaxKind.StringKeyword) {
                return {
                    foundMatch: true,
                    resolvedType: 'string',
                };
            }
            else if (syntaxKind === ts.SyntaxKind.BooleanKeyword) {
                return {
                    foundMatch: true,
                    resolvedType: 'boolean',
                };
            }
            else if (syntaxKind === ts.SyntaxKind.VoidKeyword) {
                return {
                    foundMatch: true,
                    resolvedType: 'void',
                };
            }
            else if (syntaxKind === ts.SyntaxKind.UndefinedKeyword) {
                return {
                    foundMatch: true,
                    resolvedType: 'undefined',
                };
            }
            else {
                return {
                    foundMatch: false,
                };
            }
        };
    }
    static clearCache() {
        Object.keys(localReferenceTypeCache).forEach(key => {
            delete localReferenceTypeCache[key];
        });
        Object.keys(inProgressTypes).forEach(key => {
            delete inProgressTypes[key];
        });
    }
    resolve() {
        var _a, _b;
        const primitiveType = this.getPrimitiveType(this.typeNode, this.parentNode);
        if (primitiveType) {
            return primitiveType;
        }
        if (this.typeNode.kind === ts.SyntaxKind.NullKeyword) {
            const enumType = {
                dataType: 'enum',
                enums: [null],
            };
            return enumType;
        }
        if (this.typeNode.kind === ts.SyntaxKind.UndefinedKeyword) {
            const undefinedType = {
                dataType: 'undefined',
            };
            return undefinedType;
        }
        if (ts.isArrayTypeNode(this.typeNode)) {
            const arrayMetaType = {
                dataType: 'array',
                elementType: new TypeResolver(this.typeNode.elementType, this.current, this.parentNode, this.context).resolve(),
            };
            return arrayMetaType;
        }
        if (ts.isUnionTypeNode(this.typeNode)) {
            const types = this.typeNode.types.map(type => {
                return new TypeResolver(type, this.current, this.parentNode, this.context).resolve();
            });
            const unionMetaType = {
                dataType: 'union',
                types,
            };
            return unionMetaType;
        }
        if (ts.isIntersectionTypeNode(this.typeNode)) {
            const types = this.typeNode.types.map(type => {
                return new TypeResolver(type, this.current, this.parentNode, this.context).resolve();
            });
            const intersectionMetaType = {
                dataType: 'intersection',
                types,
            };
            return intersectionMetaType;
        }
        if (this.typeNode.kind === ts.SyntaxKind.AnyKeyword || this.typeNode.kind === ts.SyntaxKind.UnknownKeyword) {
            const literallyAny = {
                dataType: 'any',
            };
            return literallyAny;
        }
        if (ts.isLiteralTypeNode(this.typeNode)) {
            const enumType = {
                dataType: 'enum',
                enums: [this.getLiteralValue(this.typeNode)],
            };
            return enumType;
        }
        if (ts.isTypeLiteralNode(this.typeNode)) {
            const properties = this.typeNode.members
                .filter(member => ts.isPropertySignature(member))
                .reduce((res, propertySignature) => {
                const type = new TypeResolver(propertySignature.type, this.current, propertySignature, this.context).resolve();
                const property = {
                    example: this.getNodeExample(propertySignature),
                    default: (0, jsDocUtils_1.getJSDocComment)(propertySignature, 'default'),
                    description: this.getNodeDescription(propertySignature),
                    format: this.getNodeFormat(propertySignature),
                    name: propertySignature.name.text,
                    required: !propertySignature.questionToken,
                    type,
                    validators: (0, validatorUtils_1.getPropertyValidators)(propertySignature) || {},
                    deprecated: (0, jsDocUtils_1.isExistJSDocTag)(propertySignature, tag => tag.tagName.text === 'deprecated'),
                    extensions: this.getNodeExtension(propertySignature),
                };
                return [property, ...res];
            }, []);
            const indexMember = this.typeNode.members.find(member => ts.isIndexSignatureDeclaration(member));
            let additionalType;
            if (indexMember) {
                const indexSignatureDeclaration = indexMember;
                const indexType = new TypeResolver(indexSignatureDeclaration.parameters[0].type, this.current, this.parentNode, this.context).resolve();
                if (indexType.dataType !== 'string') {
                    throw new exceptions_1.GenerateMetadataError(`Only string indexers are supported.`, this.typeNode);
                }
                additionalType = new TypeResolver(indexSignatureDeclaration.type, this.current, this.parentNode, this.context).resolve();
            }
            const objLiteral = {
                additionalProperties: indexMember && additionalType,
                dataType: 'nestedObjectLiteral',
                properties,
            };
            return objLiteral;
        }
        if (this.typeNode.kind === ts.SyntaxKind.ObjectKeyword) {
            return { dataType: 'object' };
        }
        if (ts.isMappedTypeNode(this.typeNode) && this.referencer) {
            const type = this.current.typeChecker.getTypeFromTypeNode(this.referencer);
            const mappedTypeNode = this.typeNode;
            const typeChecker = this.current.typeChecker;
            const getDeclaration = (prop) => prop.declarations && prop.declarations[0];
            const isIgnored = (prop) => {
                const declaration = getDeclaration(prop);
                return (prop.getJsDocTags().find(tag => tag.name === 'ignore') !== undefined ||
                    (declaration !== undefined && !ts.isPropertyDeclaration(declaration) && !ts.isPropertySignature(declaration) && !ts.isParameter(declaration)));
            };
            const properties = type
                .getProperties()
                // Ignore methods, getter, setter and @ignored props
                .filter(property => isIgnored(property) === false)
                // Transform to property
                .map(property => {
                const propertyType = typeChecker.getTypeOfSymbolAtLocation(property, this.typeNode);
                const declaration = getDeclaration(property);
                if (declaration && ts.isPropertySignature(declaration)) {
                    return { ...this.propertyFromSignature(declaration, mappedTypeNode.questionToken), name: property.getName() };
                }
                else if (declaration && (ts.isPropertyDeclaration(declaration) || ts.isParameter(declaration))) {
                    return { ...this.propertyFromDeclaration(declaration, mappedTypeNode.questionToken), name: property.getName() };
                }
                // Resolve default value, required and typeNode
                let required = false;
                const typeNode = this.current.typeChecker.typeToTypeNode(propertyType, undefined, ts.NodeBuilderFlags.NoTruncation);
                if (mappedTypeNode.questionToken && mappedTypeNode.questionToken.kind === ts.SyntaxKind.MinusToken) {
                    required = true;
                }
                else if (mappedTypeNode.questionToken && mappedTypeNode.questionToken.kind === ts.SyntaxKind.QuestionToken) {
                    required = false;
                }
                // Push property
                return {
                    name: property.getName(),
                    required,
                    // Mapped types with any amount of indirection (template strings, unions, Exclude<>, etc.)
                    // don't provide an underlying declaration in the AST, thus we cannot know if the original
                    // property is deprecated. This matches intellisense's behavior in vscode.
                    deprecated: false,
                    type: new TypeResolver(typeNode, this.current, this.typeNode, this.context, this.referencer).resolve(),
                    validators: {},
                };
            });
            const objectLiteral = {
                dataType: 'nestedObjectLiteral',
                properties,
            };
            return objectLiteral;
        }
        if (ts.isConditionalTypeNode(this.typeNode) && this.referencer && ts.isTypeReferenceNode(this.referencer)) {
            const type = this.current.typeChecker.getTypeFromTypeNode(this.referencer);
            if (type.aliasSymbol) {
                let declaration = (_a = type.aliasSymbol.declarations) === null || _a === void 0 ? void 0 : _a[0];
                if (declaration.name) {
                    declaration = this.getModelTypeDeclaration(declaration.name);
                }
                const name = this.getRefTypeName(this.referencer.getText());
                return this.handleCachingAndCircularReferences(name, () => {
                    if (ts.isTypeAliasDeclaration(declaration)) {
                        // Note: I don't understand why typescript lose type for `this.referencer` (from above with isTypeReferenceNode())
                        return this.getTypeAliasReference(declaration, this.current.typeChecker.typeToString(type), this.referencer);
                    }
                    else if (ts.isEnumDeclaration(declaration)) {
                        return this.getEnumerateType(declaration.name);
                    }
                    else {
                        throw new exceptions_1.GenerateMetadataError(`Couldn't resolve Conditional to TypeNode. If you think this should be resolvable, please file an Issue. We found an aliasSymbol and it's declaration was of kind ${declaration.kind}`, this.typeNode);
                    }
                });
            }
            else if (type.isClassOrInterface()) {
                let declaration = (_b = type.symbol.declarations) === null || _b === void 0 ? void 0 : _b[0];
                if (declaration.name) {
                    declaration = this.getModelTypeDeclaration(declaration.name);
                }
                const name = this.getRefTypeName(this.referencer.getText());
                return this.handleCachingAndCircularReferences(name, () => this.getModelReference(declaration, this.current.typeChecker.typeToString(type)));
            }
            else {
                try {
                    return new TypeResolver(this.current.typeChecker.typeToTypeNode(type, undefined, ts.NodeBuilderFlags.NoTruncation), this.current, this.typeNode, this.context, this.referencer).resolve();
                }
                catch {
                    throw new exceptions_1.GenerateMetadataError(`Couldn't resolve Conditional to TypeNode. If you think this should be resolvable, please file an Issue. The flags on the result of the ConditionalType was ${type.flags}`, this.typeNode);
                }
            }
        }
        // keyof
        if (ts.isTypeOperatorNode(this.typeNode) && this.typeNode.operator === ts.SyntaxKind.KeyOfKeyword) {
            const type = this.current.typeChecker.getTypeFromTypeNode(this.typeNode);
            if (type.isUnion()) {
                const literals = type.types.filter(t => t.isLiteral()).reduce((acc, t) => [...acc, t.value.toString()], []);
                // Warn on nonsense (`number`, `typeof Symbol.iterator`)
                if (type.types.find(t => !t.isLiteral()) !== undefined) {
                    const problems = type.types.filter(t => !t.isLiteral()).map(t => this.current.typeChecker.typeToString(t));
                    console.warn(new exceptions_1.GenerateMetaDataWarning(`Skipped non-literal type(s) ${problems.join(', ')}`, this.typeNode).toString());
                }
                return {
                    dataType: 'enum',
                    enums: literals,
                };
            }
            else if (type.isLiteral()) {
                return {
                    dataType: 'enum',
                    enums: [type.value.toString()],
                };
            }
            else if ((type.getFlags() & ts.TypeFlags.Never) !== 0) {
                throw new exceptions_1.GenerateMetadataError(`TypeOperator 'keyof' on node produced a never type`, this.typeNode);
            }
            else {
                const warning = new exceptions_1.GenerateMetaDataWarning(`Couldn't resolve keyof reliably, please check the resulting type carefully.
            Reason: Type was not a literal or an array of literals.`, this.typeNode);
                console.warn(warning.toString());
                try {
                    return new TypeResolver(this.current.typeChecker.typeToTypeNode(type, undefined, ts.NodeBuilderFlags.NoTruncation), this.current, this.typeNode, this.context, this.referencer).resolve();
                }
                catch (err) {
                    const indexedTypeName = this.current.typeChecker.typeToString(this.current.typeChecker.getTypeFromTypeNode(this.typeNode.type));
                    throw new exceptions_1.GenerateMetadataError(`Could not determine the keys on ${indexedTypeName}`, this.typeNode);
                }
            }
        }
        // Handle `readonly` arrays
        if (ts.isTypeOperatorNode(this.typeNode) && this.typeNode.operator === ts.SyntaxKind.ReadonlyKeyword) {
            return new TypeResolver(this.typeNode.type, this.current, this.typeNode, this.context, this.referencer).resolve();
        }
        // Indexed by keyword
        if (ts.isIndexedAccessTypeNode(this.typeNode) && (this.typeNode.indexType.kind === ts.SyntaxKind.NumberKeyword || this.typeNode.indexType.kind === ts.SyntaxKind.StringKeyword)) {
            const numberIndexType = this.typeNode.indexType.kind === ts.SyntaxKind.NumberKeyword;
            const objectType = this.current.typeChecker.getTypeFromTypeNode(this.typeNode.objectType);
            const type = numberIndexType ? objectType.getNumberIndexType() : objectType.getStringIndexType();
            if (type === undefined) {
                throw new exceptions_1.GenerateMetadataError(`Could not determine ${numberIndexType ? 'number' : 'string'} index on ${this.current.typeChecker.typeToString(objectType)}`, this.typeNode);
            }
            return new TypeResolver(this.current.typeChecker.typeToTypeNode(type, undefined, ts.NodeBuilderFlags.NoTruncation), this.current, this.typeNode, this.context, this.referencer).resolve();
        }
        // Indexed by literal
        if (ts.isIndexedAccessTypeNode(this.typeNode) &&
            ts.isLiteralTypeNode(this.typeNode.indexType) &&
            (ts.isStringLiteral(this.typeNode.indexType.literal) || ts.isNumericLiteral(this.typeNode.indexType.literal))) {
            const hasType = (node) => node !== undefined && Object.prototype.hasOwnProperty.call(node, 'type');
            const symbol = this.current.typeChecker.getPropertyOfType(this.current.typeChecker.getTypeFromTypeNode(this.typeNode.objectType), this.typeNode.indexType.literal.text);
            if (symbol === undefined) {
                throw new exceptions_1.GenerateMetadataError(`Could not determine the keys on ${this.current.typeChecker.typeToString(this.current.typeChecker.getTypeFromTypeNode(this.typeNode.objectType))}`, this.typeNode);
            }
            if (hasType(symbol.valueDeclaration) && symbol.valueDeclaration.type) {
                return new TypeResolver(symbol.valueDeclaration.type, this.current, this.typeNode, this.context, this.referencer).resolve();
            }
            const declaration = this.current.typeChecker.getTypeOfSymbolAtLocation(symbol, this.typeNode.objectType);
            try {
                return new TypeResolver(this.current.typeChecker.typeToTypeNode(declaration, undefined, ts.NodeBuilderFlags.NoTruncation), this.current, this.typeNode, this.context, this.referencer).resolve();
            }
            catch {
                throw new exceptions_1.GenerateMetadataError(`Could not determine the keys on ${this.current.typeChecker.typeToString(this.current.typeChecker.getTypeFromTypeNode(this.current.typeChecker.typeToTypeNode(declaration, undefined, ts.NodeBuilderFlags.NoTruncation)))}`, this.typeNode);
            }
        }
        // Indexed by keyof typeof value
        if (ts.isIndexedAccessTypeNode(this.typeNode) && ts.isTypeOperatorNode(this.typeNode.indexType) && this.typeNode.indexType.operator === ts.SyntaxKind.KeyOfKeyword) {
            const resolveParenthesis = (node) => (ts.isParenthesizedTypeNode(node) ? node.type : node);
            const objectType = resolveParenthesis(this.typeNode.objectType);
            const indexType = this.typeNode.indexType.type;
            const isSameTypeQuery = ts.isTypeQueryNode(objectType) && ts.isTypeQueryNode(indexType) && objectType.exprName.getText() === indexType.exprName.getText();
            const isSameTypeReference = ts.isTypeReferenceNode(objectType) && ts.isTypeReferenceNode(indexType) && objectType.typeName.getText() === indexType.typeName.getText();
            if (isSameTypeQuery || isSameTypeReference) {
                const type = this.current.typeChecker.getTypeFromTypeNode(this.referencer || this.typeNode);
                const node = this.current.typeChecker.typeToTypeNode(type, undefined, ts.NodeBuilderFlags.InTypeAlias | ts.NodeBuilderFlags.NoTruncation);
                return new TypeResolver(node, this.current, this.typeNode, this.context, this.referencer).resolve();
            }
        }
        if (ts.isTemplateLiteralTypeNode(this.typeNode)) {
            const type = this.current.typeChecker.getTypeFromTypeNode(this.referencer || this.typeNode);
            if (type.isUnion() && type.types.every(unionElementType => unionElementType.isStringLiteral())) {
                const stringLiteralEnum = {
                    dataType: 'enum',
                    enums: type.types.map((stringLiteralType) => stringLiteralType.value),
                };
                return stringLiteralEnum;
            }
            else {
                throw new exceptions_1.GenerateMetadataError(`Could not the type of ${this.current.typeChecker.typeToString(this.current.typeChecker.getTypeFromTypeNode(this.typeNode), this.typeNode)}`, this.typeNode);
            }
        }
        if (ts.isParenthesizedTypeNode(this.typeNode)) {
            return new TypeResolver(this.typeNode.type, this.current, this.typeNode, this.context, this.referencer).resolve();
        }
        if (this.typeNode.kind !== ts.SyntaxKind.TypeReference) {
            throw new exceptions_1.GenerateMetadataError(`Unknown type: ${ts.SyntaxKind[this.typeNode.kind]}`, this.typeNode);
        }
        const typeReference = this.typeNode;
        if (typeReference.typeName.kind === ts.SyntaxKind.Identifier) {
            if (typeReference.typeName.text === 'Date') {
                return this.getDateType(this.parentNode);
            }
            if (typeReference.typeName.text === 'Buffer') {
                const bufferMetaType = { dataType: 'buffer' };
                return bufferMetaType;
            }
            if (typeReference.typeName.text === 'Readable') {
                const streamMetaType = { dataType: 'buffer' };
                return streamMetaType;
            }
            if (typeReference.typeName.text === 'Array' && typeReference.typeArguments && typeReference.typeArguments.length === 1) {
                const arrayMetaType = {
                    dataType: 'array',
                    elementType: new TypeResolver(typeReference.typeArguments[0], this.current, this.parentNode, this.context).resolve(),
                };
                return arrayMetaType;
            }
            if (typeReference.typeName.text === 'Promise' && typeReference.typeArguments && typeReference.typeArguments.length === 1) {
                return new TypeResolver(typeReference.typeArguments[0], this.current, this.parentNode, this.context).resolve();
            }
            if (typeReference.typeName.text === 'String') {
                const stringMetaType = { dataType: 'string' };
                return stringMetaType;
            }
            if (this.context[typeReference.typeName.text]) {
                return new TypeResolver(this.context[typeReference.typeName.text], this.current, this.parentNode, this.context).resolve();
            }
        }
        const referenceType = this.getReferenceType(typeReference);
        if (this.addToRefTypeMap) {
            this.current.AddReferenceType(referenceType);
        }
        return referenceType;
    }
    getLiteralValue(typeNode) {
        let value;
        switch (typeNode.literal.kind) {
            case ts.SyntaxKind.TrueKeyword:
                value = true;
                break;
            case ts.SyntaxKind.FalseKeyword:
                value = false;
                break;
            case ts.SyntaxKind.StringLiteral:
                value = typeNode.literal.text;
                break;
            case ts.SyntaxKind.NumericLiteral:
                value = parseFloat(typeNode.literal.text);
                break;
            case ts.SyntaxKind.NullKeyword:
                value = null;
                break;
            default:
                if (Object.prototype.hasOwnProperty.call(typeNode.literal, 'text')) {
                    value = typeNode.literal.text;
                }
                else {
                    throw new exceptions_1.GenerateMetadataError(`Couldn't resolve literal node: ${typeNode.literal.getText()}`);
                }
        }
        return value;
    }
    getPrimitiveType(typeNode, parentNode) {
        const resolution = this.attemptToResolveKindToPrimitive(typeNode.kind);
        if (!resolution.foundMatch) {
            return;
        }
        if (resolution.resolvedType === 'number') {
            if (!parentNode) {
                return { dataType: 'double' };
            }
            const tags = (0, jsDocUtils_1.getJSDocTagNames)(parentNode).filter(name => {
                return ['isInt', 'isLong', 'isFloat', 'isDouble'].some(m => m === name);
            });
            if (tags.length === 0) {
                return { dataType: 'double' };
            }
            switch (tags[0]) {
                case 'isInt':
                    return { dataType: 'integer' };
                case 'isLong':
                    return { dataType: 'long' };
                case 'isFloat':
                    return { dataType: 'float' };
                case 'isDouble':
                    return { dataType: 'double' };
                default:
                    return { dataType: 'double' };
            }
        }
        else if (resolution.resolvedType === 'string') {
            return {
                dataType: 'string',
            };
        }
        else if (resolution.resolvedType === 'boolean') {
            return {
                dataType: 'boolean',
            };
        }
        else if (resolution.resolvedType === 'void') {
            return {
                dataType: 'void',
            };
        }
        else if (resolution.resolvedType === 'undefined') {
            return {
                dataType: 'undefined',
            };
        }
        else {
            return (0, runtime_1.assertNever)(resolution.resolvedType);
        }
    }
    getDateType(parentNode) {
        if (!parentNode) {
            return { dataType: 'datetime' };
        }
        const tags = (0, jsDocUtils_1.getJSDocTagNames)(parentNode).filter(name => {
            return ['isDate', 'isDateTime'].some(m => m === name);
        });
        if (tags.length === 0) {
            return { dataType: 'datetime' };
        }
        switch (tags[0]) {
            case 'isDate':
                return { dataType: 'date' };
            case 'isDateTime':
                return { dataType: 'datetime' };
            default:
                return { dataType: 'datetime' };
        }
    }
    getDesignatedModels(nodes, typeName) {
        /**
         * Model is marked with '@tsoaModel', indicating that it should be the 'canonical' model used
         */
        const designatedNodes = nodes.filter(enumNode => {
            return (0, jsDocUtils_1.isExistJSDocTag)(enumNode, tag => tag.tagName.text === 'tsoaModel');
        });
        if (designatedNodes.length > 0) {
            if (designatedNodes.length > 1) {
                throw new exceptions_1.GenerateMetadataError(`Multiple models for ${typeName} marked with '@tsoaModel'; '@tsoaModel' should only be applied to one model.`);
            }
            return designatedNodes;
        }
        return nodes;
    }
    hasFlag(type, flag) {
        return (type.flags & flag) === flag;
    }
    getEnumerateType(typeName) {
        var _a, _b, _c;
        const enumName = typeName.text;
        const symbol = this.getSymbolAtLocation(typeName);
        // resolve value
        let declaredType = (((_a = this.current.typeChecker.getDeclaredTypeOfSymbol(symbol)) === null || _a === void 0 ? void 0 : _a.symbol) || symbol);
        // if we are a EnumMember, return parent instead (this happens if a enum has only one entry, not quite sure why though...)
        if (this.hasFlag(declaredType, ts.SymbolFlags.EnumMember) && ((_c = (_b = declaredType.parent) === null || _b === void 0 ? void 0 : _b.valueDeclaration) === null || _c === void 0 ? void 0 : _c.kind) === ts.SyntaxKind.EnumDeclaration) {
            declaredType = declaredType.parent;
        }
        const declarations = declaredType.getDeclarations();
        if (!declarations) {
            return;
        }
        let enumNodes = declarations.filter((node) => {
            return ts.isEnumDeclaration(node) && node.name.getText() === enumName;
        });
        if (!enumNodes.length) {
            return;
        }
        enumNodes = this.getDesignatedModels(enumNodes, enumName);
        if (enumNodes.length > 1) {
            throw new exceptions_1.GenerateMetadataError(`Multiple matching enum found for enum ${enumName}; please make enum names unique.`);
        }
        const enumDeclaration = enumNodes[0];
        const isNotUndefined = (item) => {
            return item === undefined ? false : true;
        };
        const enums = enumDeclaration.members.map(e => this.current.typeChecker.getConstantValue(e)).filter(isNotUndefined);
        const enumVarnames = enumDeclaration.members.map(e => e.name.getText()).filter(isNotUndefined);
        return {
            dataType: 'refEnum',
            description: this.getNodeDescription(enumDeclaration),
            enums,
            enumVarnames,
            refName: enumName,
            deprecated: (0, jsDocUtils_1.isExistJSDocTag)(enumDeclaration, tag => tag.tagName.text === 'deprecated'),
        };
    }
    getReferenceType(node, addToRefTypeMap = true) {
        let type;
        if (ts.isTypeReferenceNode(node)) {
            type = node.typeName;
        }
        else if (ts.isExpressionWithTypeArguments(node)) {
            type = node.expression;
        }
        else {
            throw new exceptions_1.GenerateMetadataError(`Can't resolve Reference type.`);
        }
        // Can't invoke getText on Synthetic Nodes
        let resolvableName = node.pos !== -1 ? node.getText() : type.text;
        if (node.pos === -1 && 'typeArguments' in node && Array.isArray(node.typeArguments)) {
            // Add typearguments for Synthetic nodes (e.g. Record<> in TestClassModel.indexedResponse)
            const argumentsString = node.typeArguments.map(arg => {
                if (ts.isLiteralTypeNode(arg)) {
                    return `'${String(this.getLiteralValue(arg))}'`;
                }
                const resolved = this.attemptToResolveKindToPrimitive(arg.kind);
                if (resolved.foundMatch === false)
                    return 'any';
                return resolved.resolvedType;
            });
            resolvableName += `<${argumentsString.join(', ')}>`;
        }
        const name = this.contextualizedName(resolvableName);
        this.typeArgumentsToContext(node, type, this.context);
        try {
            const existingType = localReferenceTypeCache[name];
            if (existingType) {
                return existingType;
            }
            const refEnumType = this.getEnumerateType(type);
            if (refEnumType) {
                localReferenceTypeCache[name] = refEnumType;
                return refEnumType;
            }
            if (inProgressTypes[name]) {
                return this.createCircularDependencyResolver(name);
            }
            inProgressTypes[name] = true;
            const declaration = this.getModelTypeDeclaration(type);
            let referenceType;
            if (ts.isTypeAliasDeclaration(declaration)) {
                referenceType = this.getTypeAliasReference(declaration, name, node, addToRefTypeMap);
            }
            else if (ts.isEnumMember(declaration)) {
                referenceType = {
                    dataType: 'refEnum',
                    refName: this.getRefTypeName(name),
                    enums: [this.current.typeChecker.getConstantValue(declaration)],
                    enumVarnames: [declaration.name.getText()],
                    deprecated: (0, jsDocUtils_1.isExistJSDocTag)(declaration, tag => tag.tagName.text === 'deprecated'),
                };
            }
            else {
                referenceType = this.getModelReference(declaration, name);
            }
            localReferenceTypeCache[name] = referenceType;
            return referenceType;
        }
        catch (err) {
            // eslint-disable-next-line no-console
            console.error(`There was a problem resolving type of '${name}'.`);
            throw err;
        }
    }
    getTypeAliasReference(declaration, name, referencer, addToRefTypeMap = true) {
        const example = this.getNodeExample(declaration);
        return {
            dataType: 'refAlias',
            default: (0, jsDocUtils_1.getJSDocComment)(declaration, 'default'),
            description: this.getNodeDescription(declaration),
            refName: this.getRefTypeName(name),
            format: this.getNodeFormat(declaration),
            type: new TypeResolver(declaration.type, this.current, declaration, this.context, this.referencer || referencer, addToRefTypeMap).resolve(),
            validators: (0, validatorUtils_1.getPropertyValidators)(declaration) || {},
            ...(example && { example }),
        };
    }
    getModelReference(modelType, name) {
        const example = this.getNodeExample(modelType);
        const description = this.getNodeDescription(modelType);
        const deprecated = (0, jsDocUtils_1.isExistJSDocTag)(modelType, tag => tag.tagName.text === 'deprecated') || (0, decoratorUtils_1.isDecorator)(modelType, identifier => identifier.text === 'Deprecated');
        // Handle toJSON methods
        if (!modelType.name) {
            throw new exceptions_1.GenerateMetadataError("Can't get Symbol from anonymous class", modelType);
        }
        const type = this.current.typeChecker.getTypeAtLocation(modelType.name);
        const toJSON = this.current.typeChecker.getPropertyOfType(type, 'toJSON');
        if (toJSON && toJSON.valueDeclaration && (ts.isMethodDeclaration(toJSON.valueDeclaration) || ts.isMethodSignature(toJSON.valueDeclaration))) {
            let nodeType = toJSON.valueDeclaration.type;
            if (!nodeType) {
                const signature = this.current.typeChecker.getSignatureFromDeclaration(toJSON.valueDeclaration);
                const implicitType = this.current.typeChecker.getReturnTypeOfSignature(signature);
                nodeType = this.current.typeChecker.typeToTypeNode(implicitType, undefined, ts.NodeBuilderFlags.NoTruncation);
            }
            const type = new TypeResolver(nodeType, this.current).resolve();
            const referenceType = {
                refName: this.getRefTypeName(name),
                dataType: 'refAlias',
                description,
                type,
                validators: {},
                deprecated,
                ...(example && { example }),
            };
            return referenceType;
        }
        const properties = this.getModelProperties(modelType);
        const additionalProperties = this.getModelAdditionalProperties(modelType);
        const inheritedProperties = this.getModelInheritedProperties(modelType) || [];
        const referenceType = {
            additionalProperties,
            dataType: 'refObject',
            description,
            properties: inheritedProperties,
            refName: this.getRefTypeName(name),
            deprecated,
            ...(example && { example }),
        };
        referenceType.properties = referenceType.properties.concat(properties);
        return referenceType;
    }
    getRefTypeName(name) {
        return encodeURIComponent(name
            .replace(/<|>/g, '_')
            .replace(/\s+/g, '')
            .replace(/,/g, '.')
            .replace(/'([^']*)'/g, '$1')
            .replace(/"([^"]*)"/g, '$1')
            .replace(/&/g, '-and-')
            .replace(/\|/g, '-or-')
            .replace(/\[\]/g, '-Array')
            .replace(/{|}/g, '_') // SuccessResponse_{indexesCreated-number}_ -> SuccessResponse__indexesCreated-number__
            .replace(/([a-z]+):([a-z]+)/gi, '$1-$2') // SuccessResponse_indexesCreated:number_ -> SuccessResponse_indexesCreated-number_
            .replace(/;/g, '--')
            .replace(/([a-z]+)\[([a-z]+)\]/gi, '$1-at-$2'));
    }
    contextualizedName(name) {
        return Object.entries(this.context).reduce((acc, [key, entry]) => {
            return acc
                .replace(new RegExp(`<\\s*([^>]*\\s)*\\s*(${key})(\\s[^>]*)*\\s*>`, 'g'), `<$1${entry.getText()}$3>`)
                .replace(new RegExp(`<\\s*([^,]*\\s)*\\s*(${key})(\\s[^,]*)*\\s*,`, 'g'), `<$1${entry.getText()}$3,`)
                .replace(new RegExp(`,\\s*([^>]*\\s)*\\s*(${key})(\\s[^>]*)*\\s*>`, 'g'), `,$1${entry.getText()}$3>`)
                .replace(new RegExp(`<\\s*([^<]*\\s)*\\s*(${key})(\\s[^<]*)*\\s*<`, 'g'), `<$1${entry.getText()}$3<`);
        }, name);
    }
    handleCachingAndCircularReferences(name, declarationResolver) {
        try {
            const existingType = localReferenceTypeCache[name];
            if (existingType) {
                return existingType;
            }
            if (inProgressTypes[name]) {
                return this.createCircularDependencyResolver(name);
            }
            inProgressTypes[name] = true;
            const reference = declarationResolver();
            localReferenceTypeCache[name] = reference;
            this.current.AddReferenceType(reference);
            return reference;
        }
        catch (err) {
            // eslint-disable-next-line no-console
            console.error(`There was a problem resolving type of '${name}'.`);
            throw err;
        }
    }
    createCircularDependencyResolver(refName) {
        const referenceType = {
            dataType: 'refObject',
            refName,
        };
        this.current.OnFinish(referenceTypes => {
            const realReferenceType = referenceTypes[refName];
            if (!realReferenceType) {
                return;
            }
            referenceType.description = realReferenceType.description;
            if (realReferenceType.dataType === 'refObject' && referenceType.dataType === 'refObject') {
                referenceType.properties = realReferenceType.properties;
            }
            referenceType.dataType = realReferenceType.dataType;
            referenceType.refName = realReferenceType.refName;
        });
        return referenceType;
    }
    nodeIsUsable(node) {
        switch (node.kind) {
            case ts.SyntaxKind.InterfaceDeclaration:
            case ts.SyntaxKind.ClassDeclaration:
            case ts.SyntaxKind.TypeAliasDeclaration:
            case ts.SyntaxKind.EnumDeclaration:
            case ts.SyntaxKind.EnumMember:
                return true;
            default:
                return false;
        }
    }
    getModelTypeDeclaration(type) {
        let typeName = type.kind === ts.SyntaxKind.Identifier ? type.text : type.right.text;
        const symbol = this.getSymbolAtLocation(type);
        const declarations = symbol === null || symbol === void 0 ? void 0 : symbol.getDeclarations();
        if (!declarations) {
            throw new exceptions_1.GenerateMetadataError(`No declarations found for referenced type ${typeName}.`);
        }
        if (symbol.escapedName !== typeName && symbol.escapedName !== 'default') {
            typeName = symbol.escapedName;
        }
        let modelTypes = declarations.filter((node) => {
            var _a;
            return this.nodeIsUsable(node) && ((_a = node.name) === null || _a === void 0 ? void 0 : _a.getText()) === typeName;
        });
        if (!modelTypes.length) {
            throw new exceptions_1.GenerateMetadataError(`No matching model found for referenced type ${typeName}.`);
        }
        if (modelTypes.length > 1) {
            // remove types that are from typescript e.g. 'Account'
            modelTypes = modelTypes.filter(modelType => {
                return modelType.getSourceFile().fileName.replace(/\\/g, '/').toLowerCase().indexOf('node_modules/typescript') <= -1;
            });
            modelTypes = this.getDesignatedModels(modelTypes, typeName);
        }
        if (modelTypes.length > 1) {
            const conflicts = modelTypes.map(modelType => modelType.getSourceFile().fileName).join('"; "');
            throw new exceptions_1.GenerateMetadataError(`Multiple matching models found for referenced type ${typeName}; please make model names unique. Conflicts found: "${conflicts}".`);
        }
        return modelTypes[0];
    }
    getSymbolAtLocation(type) {
        const symbol = this.current.typeChecker.getSymbolAtLocation(type) || type.symbol;
        // resolve alias if it is an alias, otherwise take symbol directly
        return (symbol && this.hasFlag(symbol, ts.SymbolFlags.Alias) && this.current.typeChecker.getAliasedSymbol(symbol)) || symbol;
    }
    getModelProperties(node, overrideToken) {
        const isIgnored = (e) => {
            return (0, jsDocUtils_1.isExistJSDocTag)(e, tag => tag.tagName.text === 'ignore');
        };
        // Interface model
        if (ts.isInterfaceDeclaration(node)) {
            return node.members.filter(member => !isIgnored(member) && ts.isPropertySignature(member)).map((member) => this.propertyFromSignature(member, overrideToken));
        }
        // Class model
        const properties = node.members
            .filter(member => !isIgnored(member))
            .filter(member => member.kind === ts.SyntaxKind.PropertyDeclaration)
            .filter(member => !this.hasStaticModifier(member))
            .filter(member => this.hasPublicModifier(member));
        const classConstructor = node.members.find(member => ts.isConstructorDeclaration(member));
        if (classConstructor && classConstructor.parameters) {
            const constructorProperties = classConstructor.parameters.filter(parameter => this.isAccessibleParameter(parameter));
            properties.push(...constructorProperties);
        }
        return properties.map(property => this.propertyFromDeclaration(property, overrideToken));
    }
    propertyFromSignature(propertySignature, overrideToken) {
        const identifier = propertySignature.name;
        if (!propertySignature.type) {
            throw new exceptions_1.GenerateMetadataError(`No valid type found for property declaration.`);
        }
        let required = !propertySignature.questionToken;
        if (overrideToken && overrideToken.kind === ts.SyntaxKind.MinusToken) {
            required = true;
        }
        else if (overrideToken && overrideToken.kind === ts.SyntaxKind.QuestionToken) {
            required = false;
        }
        const property = {
            default: (0, jsDocUtils_1.getJSDocComment)(propertySignature, 'default'),
            description: this.getNodeDescription(propertySignature),
            example: this.getNodeExample(propertySignature),
            format: this.getNodeFormat(propertySignature),
            name: identifier.text,
            required,
            type: new TypeResolver(propertySignature.type, this.current, propertySignature.type.parent, this.context, propertySignature.type).resolve(),
            validators: (0, validatorUtils_1.getPropertyValidators)(propertySignature) || {},
            deprecated: (0, jsDocUtils_1.isExistJSDocTag)(propertySignature, tag => tag.tagName.text === 'deprecated'),
            extensions: this.getNodeExtension(propertySignature),
        };
        return property;
    }
    propertyFromDeclaration(propertyDeclaration, overrideToken) {
        const identifier = propertyDeclaration.name;
        let typeNode = propertyDeclaration.type;
        if (!typeNode) {
            const tsType = this.current.typeChecker.getTypeAtLocation(propertyDeclaration);
            typeNode = this.current.typeChecker.typeToTypeNode(tsType, undefined, ts.NodeBuilderFlags.NoTruncation);
        }
        if (!typeNode) {
            throw new exceptions_1.GenerateMetadataError(`No valid type found for property declaration.`);
        }
        const type = new TypeResolver(typeNode, this.current, propertyDeclaration, this.context, typeNode).resolve();
        let required = !propertyDeclaration.questionToken && !propertyDeclaration.initializer;
        if (overrideToken && overrideToken.kind === ts.SyntaxKind.MinusToken) {
            required = true;
        }
        else if (overrideToken && overrideToken.kind === ts.SyntaxKind.QuestionToken) {
            required = false;
        }
        const property = {
            default: (0, initializer_value_1.getInitializerValue)(propertyDeclaration.initializer, this.current.typeChecker),
            description: this.getNodeDescription(propertyDeclaration),
            example: this.getNodeExample(propertyDeclaration),
            format: this.getNodeFormat(propertyDeclaration),
            name: identifier.text,
            required,
            type,
            validators: (0, validatorUtils_1.getPropertyValidators)(propertyDeclaration) || {},
            // class properties and constructor parameters may be deprecated either via jsdoc annotation or decorator
            deprecated: (0, jsDocUtils_1.isExistJSDocTag)(propertyDeclaration, tag => tag.tagName.text === 'deprecated') || (0, decoratorUtils_1.isDecorator)(propertyDeclaration, identifier => identifier.text === 'Deprecated'),
            extensions: this.getNodeExtension(propertyDeclaration),
        };
        return property;
    }
    getModelAdditionalProperties(node) {
        if (node.kind === ts.SyntaxKind.InterfaceDeclaration) {
            const interfaceDeclaration = node;
            const indexMember = interfaceDeclaration.members.find(member => member.kind === ts.SyntaxKind.IndexSignature);
            if (!indexMember) {
                return undefined;
            }
            const indexSignatureDeclaration = indexMember;
            const indexType = new TypeResolver(indexSignatureDeclaration.parameters[0].type, this.current, this.parentNode, this.context).resolve();
            if (indexType.dataType !== 'string') {
                throw new exceptions_1.GenerateMetadataError(`Only string indexers are supported.`, this.typeNode);
            }
            return new TypeResolver(indexSignatureDeclaration.type, this.current, this.parentNode, this.context).resolve();
        }
        return undefined;
    }
    typeArgumentsToContext(type, targetEntity, context) {
        this.context = {};
        const declaration = this.getModelTypeDeclaration(targetEntity);
        const typeParameters = 'typeParameters' in declaration ? declaration.typeParameters : undefined;
        if (typeParameters) {
            for (let index = 0; index < typeParameters.length; index++) {
                const typeParameter = typeParameters[index];
                const typeArg = type.typeArguments && type.typeArguments[index];
                let resolvedType;
                // Argument may be a forward reference from context
                if (typeArg && ts.isTypeReferenceNode(typeArg) && ts.isIdentifier(typeArg.typeName) && context[typeArg.typeName.text]) {
                    resolvedType = context[typeArg.typeName.text];
                }
                else if (typeArg) {
                    resolvedType = typeArg;
                }
                else if (typeParameter.default) {
                    resolvedType = typeParameter.default;
                }
                else {
                    throw new exceptions_1.GenerateMetadataError(`Could not find a value for type parameter ${typeParameter.name.text}`, type);
                }
                this.context = {
                    ...this.context,
                    [typeParameter.name.text]: resolvedType,
                };
            }
        }
        return context;
    }
    getModelInheritedProperties(modelTypeDeclaration) {
        let properties = [];
        const heritageClauses = modelTypeDeclaration.heritageClauses;
        if (!heritageClauses) {
            return properties;
        }
        for (const clause of heritageClauses) {
            if (!clause.types) {
                continue;
            }
            for (const t of clause.types) {
                const baseEntityName = t.expression;
                // create subContext
                const resetCtx = this.typeArgumentsToContext(t, baseEntityName, this.context);
                const referenceType = this.getReferenceType(t, false);
                if (referenceType) {
                    if (referenceType.dataType === 'refEnum') {
                        // since it doesn't have properties to iterate over, then we don't do anything with it
                    }
                    else if (referenceType.dataType === 'refAlias') {
                        let type = referenceType;
                        while (type.dataType === 'refAlias') {
                            type = type.type;
                        }
                        if (type.dataType === 'refObject') {
                            properties = [...properties, ...type.properties];
                        }
                        else if (type.dataType === 'nestedObjectLiteral') {
                            properties = [...properties, ...type.properties];
                        }
                    }
                    else if (referenceType.dataType === 'refObject') {
                        (referenceType.properties || []).forEach(property => properties.push(property));
                    }
                    else {
                        (0, runtime_1.assertNever)(referenceType);
                    }
                }
                // reset subContext
                this.context = resetCtx;
            }
        }
        return properties;
    }
    hasPublicModifier(node) {
        return (!node.modifiers ||
            node.modifiers.every(modifier => {
                return modifier.kind !== ts.SyntaxKind.ProtectedKeyword && modifier.kind !== ts.SyntaxKind.PrivateKeyword;
            }));
    }
    hasStaticModifier(node) {
        return (node.modifiers &&
            node.modifiers.some(modifier => {
                return modifier.kind === ts.SyntaxKind.StaticKeyword;
            }));
    }
    isAccessibleParameter(node) {
        // No modifiers
        if (!node.modifiers) {
            return false;
        }
        // public || public readonly
        if (node.modifiers.some(modifier => modifier.kind === ts.SyntaxKind.PublicKeyword)) {
            return true;
        }
        // readonly, not private readonly, not public readonly
        const isReadonly = node.modifiers.some(modifier => modifier.kind === ts.SyntaxKind.ReadonlyKeyword);
        const isProtectedOrPrivate = node.modifiers.some(modifier => {
            return modifier.kind === ts.SyntaxKind.ProtectedKeyword || modifier.kind === ts.SyntaxKind.PrivateKeyword;
        });
        return isReadonly && !isProtectedOrPrivate;
    }
    getNodeDescription(node) {
        const symbol = this.getSymbolAtLocation(node.name);
        if (!symbol) {
            return undefined;
        }
        /**
         * TODO: Workaround for what seems like a bug in the compiler
         * Warrants more investigation and possibly a PR against typescript
         */
        if (node.kind === ts.SyntaxKind.Parameter) {
            // TypeScript won't parse jsdoc if the flag is 4, i.e. 'Property'
            symbol.flags = 0;
        }
        const comments = symbol.getDocumentationComment(this.current.typeChecker);
        if (comments.length) {
            return ts.displayPartsToString(comments);
        }
        return undefined;
    }
    getNodeFormat(node) {
        return (0, jsDocUtils_1.getJSDocComment)(node, 'format');
    }
    getNodeExample(node) {
        const exampleJSDoc = (0, jsDocUtils_1.getJSDocComment)(node, 'example');
        if (exampleJSDoc) {
            return (0, jsonUtils_1.safeFromJson)(exampleJSDoc);
        }
        return (0, decoratorUtils_1.getNodeFirstDecoratorValue)(node, this.current.typeChecker, dec => dec.text === 'Example');
    }
    getNodeExtension(node) {
        const decorators = this.getDecoratorsByIdentifier(node, 'Extension');
        const extensionDecorator = (0, extension_1.getExtensions)(decorators, this.current);
        const extensionComments = (0, jsDocUtils_1.getJSDocComments)(node, 'extension');
        const extensionJSDoc = extensionComments ? (0, extension_1.getExtensionsFromJSDocComments)(extensionComments) : [];
        return extensionDecorator.concat(extensionJSDoc);
    }
    getDecoratorsByIdentifier(node, id) {
        return (0, decoratorUtils_1.getDecorators)(node, identifier => identifier.text === id);
    }
}
exports.TypeResolver = TypeResolver;
//# sourceMappingURL=typeResolver.js.map