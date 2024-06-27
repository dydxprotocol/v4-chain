import Hook from '../../index.js'
import foo from '../fixtures/re-export-cjs-built-in.js'
import foo2 from '../fixtures/re-export-cjs.js'
import { strictEqual } from 'assert'

Hook((exports, name) => {
  if (name.endsWith('fixtures/re-export-cjs-built-in.js')) {
    strictEqual(typeof exports.default, 'function')
    exports.default = '1'
  }

  if (name.endsWith('fixtures/re-export-cjs.js')) {
    strictEqual(exports.default, 'bar')
    exports.default = '2'
  }
})

strictEqual(foo, '1')
strictEqual(foo2, '2')
