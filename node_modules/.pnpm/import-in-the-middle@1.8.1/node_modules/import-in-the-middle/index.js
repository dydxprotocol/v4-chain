// Unless explicitly stated otherwise all files in this repository are licensed under the Apache 2.0 License.
//
// This product includes software developed at Datadog (https://www.datadoghq.com/). Copyright 2021 Datadog, Inc.

const path = require('path')
const parse = require('module-details-from-path')
const { fileURLToPath } = require('url')

const {
  importHooks,
  specifiers,
  toHook
} = require('./lib/register')

function addHook (hook) {
  importHooks.push(hook)
  toHook.forEach(([name, namespace]) => hook(name, namespace))
}

function removeHook (hook) {
  const index = importHooks.indexOf(hook)
  if (index > -1) {
    importHooks.splice(index, 1)
  }
}

function callHookFn (hookFn, namespace, name, baseDir) {
  const newDefault = hookFn(namespace, name, baseDir)
  if (newDefault && newDefault !== namespace) {
    namespace.default = newDefault
  }
}

function Hook (modules, options, hookFn) {
  if ((this instanceof Hook) === false) return new Hook(modules, options, hookFn)
  if (typeof modules === 'function') {
    hookFn = modules
    modules = null
    options = null
  } else if (typeof options === 'function') {
    hookFn = options
    options = null
  }
  const internals = options ? options.internals === true : false

  this._iitmHook = (name, namespace) => {
    const filename = name
    const isBuiltin = name.startsWith('node:')
    let baseDir

    if (isBuiltin) {
      name = name.replace(/^node:/, '')
    } else {
      if (name.startsWith('file://')) {
        try {
          name = fileURLToPath(name)
        } catch (e) {}
      }
      const details = parse(name)
      if (details) {
        name = details.name
        baseDir = details.basedir
      }
    }

    if (modules) {
      for (const moduleName of modules) {
        if (moduleName === name) {
          if (baseDir) {
            if (internals) {
              name = name + path.sep + path.relative(baseDir, fileURLToPath(filename))
            } else {
              if (!baseDir.endsWith(specifiers.get(filename))) continue
            }
          }
          callHookFn(hookFn, namespace, name, baseDir)
        }
      }
    } else {
      callHookFn(hookFn, namespace, name, baseDir)
    }
  }

  addHook(this._iitmHook)
}

Hook.prototype.unhook = function () {
  removeHook(this._iitmHook)
}

module.exports = Hook
module.exports.Hook = Hook
module.exports.addHook = addHook
module.exports.removeHook = removeHook
