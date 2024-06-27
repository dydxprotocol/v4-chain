"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.prettyTroubleCause = exports.prettyLocationOfNode = exports.GenerateMetaDataWarning = exports.GenerateMetadataError = void 0;
const path_1 = require("path");
class GenerateMetadataError extends Error {
    constructor(message, node, onlyCurrent = false) {
        super(message);
        if (node) {
            this.message = `${message}\n${prettyLocationOfNode(node)}\n${prettyTroubleCause(node, onlyCurrent)}`;
        }
    }
}
exports.GenerateMetadataError = GenerateMetadataError;
class GenerateMetaDataWarning {
    constructor(message, node, onlyCurrent = false) {
        this.message = message;
        this.node = node;
        this.onlyCurrent = onlyCurrent;
    }
    toString() {
        return `Warning: ${this.message}\n${prettyLocationOfNode(this.node)}\n${prettyTroubleCause(this.node, this.onlyCurrent)}`;
    }
}
exports.GenerateMetaDataWarning = GenerateMetaDataWarning;
function prettyLocationOfNode(node) {
    const sourceFile = node.getSourceFile();
    const token = node.getFirstToken() || node.parent.getFirstToken();
    const start = token ? `:${sourceFile.getLineAndCharacterOfPosition(token.getStart()).line + 1}` : '';
    const end = token ? `:${sourceFile.getLineAndCharacterOfPosition(token.getEnd()).line + 1}` : '';
    const normalizedPath = (0, path_1.normalize)(`${sourceFile.fileName}${start}${end}`);
    return `At: ${normalizedPath}.`;
}
exports.prettyLocationOfNode = prettyLocationOfNode;
function prettyTroubleCause(node, onlyCurrent = false) {
    let name;
    if (onlyCurrent || !node.parent) {
        name = node.pos !== -1 ? node.getText() : node.name.text;
    }
    else {
        name = node.parent.pos !== -1 ? node.parent.getText() : node.parent.name.text;
    }
    return `This was caused by '${name}'`;
}
exports.prettyTroubleCause = prettyTroubleCause;
//# sourceMappingURL=exceptions.js.map