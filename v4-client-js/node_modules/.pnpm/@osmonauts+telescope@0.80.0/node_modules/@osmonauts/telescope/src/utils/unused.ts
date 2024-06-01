//@ts-nocheck

import * as t from '@babel/types';

// https://github.com/chuyik/babel-plugin-danger-remove-unused-import
// https://github.com/chuyik/babel-plugin-danger-remove-unused-import/blob/c5454c21e94698a2464a12baa5590761932a71a8/License#L1

export const unused = {
  Program: {
    exit: (path) => {
      const UnRefBindings = new Map()
      for (const [name, binding] of Object.entries(path.scope.bindings)) {
        if (!binding.path.parentPath || binding.kind !== 'module') continue

        const source = binding.path.parentPath.get('source')
        const importName = source.node.value
        if (
          !t.isStringLiteral(source)
        )
          continue

        const key = `${importName}(${source.node.loc &&
          source.node.loc.start.line})`

        if (!UnRefBindings.has(key)) {
          UnRefBindings.set(key, binding)
        }

        if (binding.referenced) {
          UnRefBindings.set(key, null)
        } else {
          const nodeType = binding.path.node.type
          if (nodeType === 'ImportSpecifier') {
            binding.path.remove()
          } else if (nodeType === 'ImportDefaultSpecifier') {
            binding.path.remove()
          } else if (nodeType === 'ImportNamespaceSpecifier') {
            binding.path.remove()
          } else if (binding.path.parentPath) {
            binding.path.parentPath.remove()
          }
        }
      }

      UnRefBindings.forEach((binding, key) => {
        if (binding && binding.path.parentPath) {
          binding.path.parentPath.remove()
        }
      })
    }
  }
};