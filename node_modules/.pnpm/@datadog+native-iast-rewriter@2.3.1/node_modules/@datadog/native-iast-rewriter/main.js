/**
 * Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
 * This product includes software developed at Datadog (https://www.datadoghq.com/). Copyright 2022 Datadog, Inc.
 **/
'use strict'
const { getPrepareStackTrace, kSymbolPrepareStackTrace } = require('./js/stack-trace/')
const { cacheRewrittenSourceMap, getOriginalPathAndLineFromSourceMap } = require('./js/source-map')

class DummyRewriter {
  rewrite (code, file) {
    return {
      content: code
    }
  }

  csiMethods () {
    return []
  }
}

let NativeRewriter
class CacheRewriter {
  constructor (config) {
    if (NativeRewriter) {
      this.nativeRewriter = new NativeRewriter(config)
      this.setLogger(config)
    } else {
      this.nativeRewriter = new DummyRewriter()
    }
  }

  rewrite (code, file) {
    const response = this.nativeRewriter.rewrite(code, file)

    try {
      const { metrics, content } = response
      if (metrics?.status === 'modified') {
        cacheRewrittenSourceMap(file, content)
      }
    } catch (e) {
      this.logError(e)
    }

    return response
  }

  csiMethods () {
    return this.nativeRewriter.csiMethods()
  }

  setLogger (config) {
    if (config && (config.logger || config.logLevel)) {
      this.logger = config.logger || console
      const logLevel = config.logLevel || 'ERROR'
      try {
        this.nativeRewriter.setLogger(this.logger, logLevel)
      } catch (e) {
        this.logError(e)
      }
    }
  }

  logError (e) {
    this.logger?.error?.(e)
  }
}

function getRewriter () {
  try {
    const iastRewriter = require('./wasm/wasm_iast_rewriter')

    NativeRewriter = iastRewriter.Rewriter
    return CacheRewriter
  } catch (e) {
    return DummyRewriter
  }
}

module.exports = {
  Rewriter: getRewriter(),
  DummyRewriter,
  getPrepareStackTrace,
  getOriginalPathAndLineFromSourceMap,
  kSymbolPrepareStackTrace
}
