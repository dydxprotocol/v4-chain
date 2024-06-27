"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.commentToString = exports.isExistJSDocTag = exports.getJSDocTags = exports.getJSDocTagNames = exports.getJSDocComments = exports.getJSDocComment = exports.getJSDocDescription = void 0;
const ts = require("typescript");
const exceptions_1 = require("../metadataGeneration/exceptions");
function getJSDocDescription(node) {
    const jsDocs = node.jsDoc;
    if (!jsDocs || !jsDocs.length) {
        return undefined;
    }
    return commentToString(jsDocs[0].comment) || undefined;
}
exports.getJSDocDescription = getJSDocDescription;
function getJSDocComment(node, tagName) {
    const comments = getJSDocComments(node, tagName);
    if (comments && comments.length !== 0) {
        return comments[0];
    }
    return;
}
exports.getJSDocComment = getJSDocComment;
function getJSDocComments(node, tagName) {
    const tags = getJSDocTags(node, tag => tag.tagName.text === tagName || tag.tagName.escapedText === tagName);
    if (tags.length === 0) {
        return;
    }
    const comments = [];
    tags.forEach(tag => {
        const comment = commentToString(tag.comment);
        if (comment)
            comments.push(comment);
    });
    return comments;
}
exports.getJSDocComments = getJSDocComments;
function getJSDocTagNames(node, requireTagName = false) {
    let tags;
    if (node.kind === ts.SyntaxKind.Parameter) {
        const parameterName = node.name.text;
        tags = getJSDocTags(node.parent, tag => {
            var _a;
            if (ts.isJSDocParameterTag(tag)) {
                return false;
            }
            else if (tag.comment === undefined) {
                throw new exceptions_1.GenerateMetadataError(`Orphan tag: @${String(tag.tagName.text || tag.tagName.escapedText)} should have a parameter name follows with.`);
            }
            return ((_a = commentToString(tag.comment)) === null || _a === void 0 ? void 0 : _a.startsWith(parameterName)) || false;
        });
    }
    else {
        tags = getJSDocTags(node, tag => {
            return requireTagName ? tag.comment !== undefined : true;
        });
    }
    return tags.map(tag => {
        return tag.tagName.text;
    });
}
exports.getJSDocTagNames = getJSDocTagNames;
function getJSDocTags(node, isMatching) {
    const jsDocs = node.jsDoc;
    if (!jsDocs || jsDocs.length === 0) {
        return [];
    }
    const jsDoc = jsDocs[0];
    if (!jsDoc.tags) {
        return [];
    }
    return jsDoc.tags.filter(isMatching);
}
exports.getJSDocTags = getJSDocTags;
function isExistJSDocTag(node, isMatching) {
    const tags = getJSDocTags(node, isMatching);
    if (tags.length === 0) {
        return false;
    }
    return true;
}
exports.isExistJSDocTag = isExistJSDocTag;
function commentToString(comment) {
    if (typeof comment === 'string') {
        return comment;
    }
    else if (comment) {
        return comment.map(node => node.text).join(' ');
    }
    return undefined;
}
exports.commentToString = commentToString;
//# sourceMappingURL=jsDocUtils.js.map