
/**
 * Module dependencies.
 */

import compoundModelFactory from './compound-model-factory';
import modelFactory from './model-factory';

/**
 * Export `clearDatabase`.
 */

export default async function() {
  const TestModel = modelFactory({ fields: ['foo', 'bar'] });
  const CompoundTestModel = compoundModelFactory({ fields: [['bar', 'foo']] });

  await TestModel.query().truncate();
  await CompoundTestModel.query().truncate();
};
