/**
 * @fileoverview Rule to flag use of .only in tests, preventing focused tests being committed accidentally
 * @author Levi Buzolic
 */

'use strict';

//------------------------------------------------------------------------------
// Rule Definition
//------------------------------------------------------------------------------

const BLOCK_DEFAULTS = ['describe', 'it', 'context', 'test', 'tape', 'fixture', 'serial'];
const FOCUS_DEFAULTS = ['only'];

module.exports = {
  meta: {
    docs: {
      description: 'disallow .only blocks in tests',
      category: 'Possible Errors',
      recommended: true,
      url: 'https://github.com/levibuzolic/eslint-plugin-no-only-tests',
    },
    schema: [
      {
        type: 'object',
        properties: {
          block: {
            type: 'array',
            items: {
              type: 'string',
            },
            uniqueItems: true,
          },
          focus: {
            type: 'array',
            items: {
              type: 'string',
            },
            uniqueItems: true,
          },
        },
        additionalProperties: false,
      },
    ],
  },
  create(context) {
    var block = (context.options[0] || {}).block || BLOCK_DEFAULTS;
    var focus = (context.options[0] || {}).focus || FOCUS_DEFAULTS;

    return {
      Identifier(node) {
        var parentObject = node.parent && node.parent.object;
        if (parentObject == null) return;
        if (focus.indexOf(node.name) === -1) return;

        var parentName = parentObject.name;

        if (parentName != null && block.indexOf(parentName) != -1) {
          context.report(node, parentName + '.' + node.name + ' not permitted');
        }

        var parentParentName = dotName(parentObject);

        if (parentParentName != null && block.indexOf(parentParentName) != -1) {
          context.report(node, parentParentName + '.' + node.name + ' not permitted');
        }
      },
    };
  },
};

function dotName(object) {
  if (object.property && object.property.name && object.object && object.object.name)
    return object.object.name + '.' + object.property.name;
  return null;
}
