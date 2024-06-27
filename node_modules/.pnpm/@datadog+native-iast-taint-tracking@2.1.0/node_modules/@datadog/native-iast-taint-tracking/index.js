/**
 * Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
 * This product includes software developed at Datadog (https://www.datadoghq.com/). Copyright 2022 Datadog, Inc.
 **/

'use strict'

let addon
try {
  addon = require('node-gyp-build')(__dirname)
} catch (e) {
  addon = {
    createTransaction (transactionId) {
      return transactionId
    },
    newTaintedString (transactionId, original) {
      return original
    },
    newTaintedObject (transactionId, original) {
      return original
    },
    addSecureMarksToTaintedString (transactionId, original) {
      return original
    },
    isTainted () {
      return false
    },
    getMetrics () {
      return undefined
    },
    getRanges () {
      return undefined
    },
    removeTransaction () {
    },
    setMaxTransactions () {
    },
    replace (transactionId, result) {
      return result
    },
    concat (transactionId, result) {
      return result
    },
    trim (transaction, result) {
      return result
    },
    trimEnd (transaction, result) {
      return result
    },
    slice (transaction, result) {
      return result
    },
    substring (transaction, result) {
      return result
    },
    substr (transaction, result) {
      return result
    },
    stringCase (transaction, result) {
      return result
    },
    arrayJoin (transaction, result) {
      return result
    }
  }
}

const iastNativeMethods = {
  newTaintedString: addon.newTaintedString,
  newTaintedObject: addon.newTaintedObject,
  addSecureMarksToTaintedString: addon.addSecureMarksToTaintedString,
  isTainted: addon.isTainted,
  getMetrics: addon.getMetrics,
  getRanges: addon.getRanges,
  createTransaction: addon.createTransaction,
  removeTransaction: addon.removeTransaction,
  setMaxTransactions: addon.setMaxTransactions,
  replace: require('./replace.js')(addon),
  concat: addon.concat,
  trim: addon.trim,
  trimEnd: addon.trimEnd,
  slice: addon.slice,
  substring: addon.substring,
  substr: addon.substr,
  stringCase: addon.stringCase,
  arrayJoin: addon.arrayJoin
}

module.exports = iastNativeMethods
