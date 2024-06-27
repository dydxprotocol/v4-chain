// Unless explicitly stated otherwise all files in this repository are licensed under the Apache 2.0 License.
//
// This product includes software developed at Datadog (https://www.datadoghq.com/). Copyright 2021 Datadog, Inc.

const { URL } = require('url')
const specifiers = new Map()
const isWin = process.platform === 'win32'

// FIXME: Typescript extensions are added temporarily until we find a better
// way of supporting arbitrary extensions
const EXTENSION_RE = /\.(js|mjs|cjs|ts|mts|cts)$/
const NODE_VERSION = process.versions.node.split('.')
const NODE_MAJOR = Number(NODE_VERSION[0])
const NODE_MINOR = Number(NODE_VERSION[1])

let entrypoint

let getExports
if (NODE_MAJOR >= 20 || (NODE_MAJOR === 18 && NODE_MINOR >= 19)) {
  getExports = require('./lib/get-exports.js')
} else {
  getExports = (url) => import(url).then(Object.keys)
}

function hasIitm (url) {
  try {
    return new URL(url).searchParams.has('iitm')
  } catch {
    return false
  }
}

function isIitm (url, meta) {
  return url === meta.url || url === meta.url.replace('hook.mjs', 'hook.js')
}

function deleteIitm (url) {
  let resultUrl
  try {
    const urlObj = new URL(url)
    if (urlObj.searchParams.has('iitm')) {
      urlObj.searchParams.delete('iitm')
      resultUrl = urlObj.href
      if (resultUrl.startsWith('file:node:')) {
        resultUrl = resultUrl.replace('file:', '')
      }
      if (resultUrl.startsWith('file:///node:')) {
        resultUrl = resultUrl.replace('file:///', '')
      }
    } else {
      resultUrl = urlObj.href
    }
  } catch {
    resultUrl = url
  }
  return resultUrl
}

function isNodeMajor16AndMinor17OrGreater () {
  return NODE_MAJOR === 16 && NODE_MINOR >= 17
}

function isFileProtocol (urlObj) {
  return urlObj.protocol === 'file:'
}

function isNodeProtocol (urlObj) {
  return urlObj.protocol === 'node:'
}

function needsToAddFileProtocol (urlObj) {
  if (NODE_MAJOR === 17) {
    return !isFileProtocol(urlObj)
  }
  if (isNodeMajor16AndMinor17OrGreater()) {
    return !isFileProtocol(urlObj) && !isNodeProtocol(urlObj)
  }
  return !isFileProtocol(urlObj) && NODE_MAJOR < 18
}

/**
 * Determines if a specifier represents an export all ESM line.
 * Note that the expected `line` isn't 100% valid ESM. It is derived
 * from the `getExports` function wherein we have recognized the true
 * line and re-mapped it to one we expect.
 *
 * @param {string} line
 * @returns {boolean}
 */
function isStarExportLine (line) {
  return /^\* from /.test(line)
}

function isBareSpecifier (specifier) {
  // Relative and absolute paths are not bare specifiers.
  if (
    specifier.startsWith('.') ||
    specifier.startsWith('/')) {
    return false
  }

  // Valid URLs are not bare specifiers. (file:, http:, node:, etc.)

  // eslint-disable-next-line no-prototype-builtins
  if (URL.hasOwnProperty('canParse')) {
    return !URL.canParse(specifier)
  }

  try {
    // eslint-disable-next-line no-new
    new URL(specifier)
    return false
  } catch (err) {
    return true
  }
}

/**
 * Processes a module's exports and builds a set of setter blocks.
 *
 * @param {object} params
 * @param {string} params.srcUrl The full URL to the module to process.
 * @param {object} params.context Provided by the loaders API.
 * @param {Function} params.parentGetSource Provides the source code for the parent module.
 * @param {bool} params.excludeDefault Exclude the default export.
 *
 * @returns {Promise<Map<string, string>>} The shimmed setters for all the exports
 * from the module and any transitive export all modules.
 */
async function processModule ({ srcUrl, context, parentGetSource, parentResolve, excludeDefault }) {
  const exportNames = await getExports(srcUrl, context, parentGetSource)
  const starExports = new Set()
  const setters = new Map()

  const addSetter = (name, setter, isStarExport = false) => {
    if (setters.has(name)) {
      if (isStarExport) {
        // If there's already a matching star export, delete it
        if (starExports.has(name)) {
          setters.delete(name)
        }
        // and return so this is excluded
        return
      }

      // if we already have this export but it is from a * export, overwrite it
      if (starExports.has(name)) {
        starExports.delete(name)
        setters.set(name, setter)
      }
    } else {
      // Store export * exports so we know they can be overridden by explicit
      // named exports
      if (isStarExport) {
        starExports.add(name)
      }

      setters.set(name, setter)
    }
  }

  for (const n of exportNames) {
    if (n === 'default' && excludeDefault) continue

    if (isStarExportLine(n) === true) {
      const [, modFile] = n.split('* from ')

      let modUrl
      if (isBareSpecifier(modFile)) {
        // Bare specifiers need to be resolved relative to the parent module.
        const result = await parentResolve(modFile, { parentURL: srcUrl })
        modUrl = result.url
      } else {
        modUrl = new URL(modFile, srcUrl).href
      }

      const setters = await processModule({
        srcUrl: modUrl,
        context,
        parentGetSource,
        parentResolve,
        excludeDefault: true
      })
      for (const [name, setter] of setters.entries()) {
        addSetter(name, setter, true)
      }
    } else {
      addSetter(n, `
      let $${n} = _.${n}
      export { $${n} as ${n} }
      set.${n} = (v) => {
        $${n} = v
        return true
      }
      `)
    }
  }

  return setters
}

function addIitm (url) {
  const urlObj = new URL(url)
  urlObj.searchParams.set('iitm', 'true')
  return needsToAddFileProtocol(urlObj) ? 'file:' + urlObj.href : urlObj.href
}

function createHook (meta) {
  let cachedResolve
  const iitmURL = new URL('lib/register.js', meta.url).toString()

  async function resolve (specifier, context, parentResolve) {
    cachedResolve = parentResolve

    // See github.com/DataDog/import-in-the-middle/pull/76.
    if (specifier === iitmURL) {
      return {
        url: specifier,
        shortCircuit: true
      }
    }

    const { parentURL = '' } = context
    const newSpecifier = deleteIitm(specifier)
    if (isWin && parentURL.indexOf('file:node') === 0) {
      context.parentURL = ''
    }
    const url = await parentResolve(newSpecifier, context, parentResolve)
    if (parentURL === '' && !EXTENSION_RE.test(url.url)) {
      entrypoint = url.url
      return { url: url.url, format: 'commonjs' }
    }

    if (isIitm(parentURL, meta) || hasIitm(parentURL)) {
      return url
    }

    // Node.js v21 renames importAssertions to importAttributes
    if (
      (context.importAssertions && context.importAssertions.type === 'json') ||
      (context.importAttributes && context.importAttributes.type === 'json')
    ) {
      return url
    }

    // If the file is referencing itself, we need to skip adding the iitm search params
    if (url.url === parentURL) {
      return {
        url: url.url,
        shortCircuit: true,
        format: url.format
      }
    }

    specifiers.set(url.url, specifier)

    return {
      url: addIitm(url.url),
      shortCircuit: true,
      format: url.format
    }
  }

  async function getSource (url, context, parentGetSource) {
    if (hasIitm(url)) {
      const realUrl = deleteIitm(url)

      try {
        const setters = await processModule({
          srcUrl: realUrl,
          context,
          parentGetSource,
          parentResolve: cachedResolve
        })
        return {
          source: `
import { register } from '${iitmURL}'
import * as namespace from ${JSON.stringify(realUrl)}

// Mimic a Module object (https://tc39.es/ecma262/#sec-module-namespace-objects).
const _ = Object.assign(
  Object.create(null, { [Symbol.toStringTag]: { value: 'Module' } }),
  namespace
)
const set = {}

${Array.from(setters.values()).join('\n')}

register(${JSON.stringify(realUrl)}, _, set, ${JSON.stringify(specifiers.get(realUrl))})
`
        }
      } catch (cause) {
        // If there are other ESM loader hooks registered as well as iitm,
        // depending on the order they are registered, source might not be
        // JavaScript.
        //
        // If we fail to parse a module for exports, we should fall back to the
        // parent loader. These modules will not be wrapped with proxies and
        // cannot be Hook'ed but at least this does not take down the entire app
        // and block iitm from being used.
        //
        // We log the error because there might be bugs in iitm and without this
        // it would be very tricky to debug
        const err = new Error(`'import-in-the-middle' failed to wrap '${realUrl}'`)
        err.cause = cause
        console.warn(err)

        // Revert back to the non-iitm URL
        url = realUrl
      }
    }

    return parentGetSource(url, context, parentGetSource)
  }

  // For Node.js 16.12.0 and higher.
  async function load (url, context, parentLoad) {
    if (hasIitm(url)) {
      const { source } = await getSource(url, context, parentLoad)
      return {
        source,
        shortCircuit: true,
        format: 'module'
      }
    }

    return parentLoad(url, context, parentLoad)
  }

  if (NODE_MAJOR >= 17 || (NODE_MAJOR === 16 && NODE_MINOR >= 12)) {
    return { load, resolve }
  } else {
    return {
      load,
      resolve,
      getSource,
      getFormat (url, context, parentGetFormat) {
        if (hasIitm(url)) {
          return {
            format: 'module'
          }
        }
        if (url === entrypoint) {
          return {
            format: 'commonjs'
          }
        }

        return parentGetFormat(url, context, parentGetFormat)
      }
    }
  }
}

module.exports = { createHook }
