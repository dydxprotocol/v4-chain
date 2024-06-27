// Unless explicitly stated otherwise all files in this repository are licensed under the Apache 2.0 License.
//
// This product includes software developed at Datadog (https://www.datadoghq.com/). Copyright 2021 Datadog, Inc.

import jsonMjs from '../fixtures/json-attributes.mjs'
import { strictEqual } from 'assert'

// Acorn does not support import attributes so an error is logged but the import
// still works!
//
// Hook((exports, name) => {
//   if (name.match(/json\.mjs/)) {
//     exports.default.data += '-dawg'
//   }
// })

strictEqual(jsonMjs.data, 'dog')
