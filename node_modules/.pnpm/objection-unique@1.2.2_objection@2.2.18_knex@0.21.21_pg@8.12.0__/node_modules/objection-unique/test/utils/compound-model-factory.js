
/**
 * Module dependencies.
 */

import { Model } from 'objection';
import uniquePlugin from '../../index';

/**
 * Export `TestModel`.
 */

export default options => {
  const unique = uniquePlugin(options);

  return class CompoundTestModel extends unique(Model) {

    /**
     * Table name.
     */

    static get tableName() {
      return 'CompoundTest';
    }
  };
};
