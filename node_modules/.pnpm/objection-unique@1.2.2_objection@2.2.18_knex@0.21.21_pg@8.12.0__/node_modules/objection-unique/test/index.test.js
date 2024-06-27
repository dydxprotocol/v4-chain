
/**
 * Module dependencies.
 */

import { ValidationError } from 'objection';
import clearDatabase from './utils/clear-database';
import compoundModelFactory from './utils/compound-model-factory';
import modelFactory from './utils/model-factory';

/**
 * Test `FoobarController`.
 */

describe('FoobarController', () => {
  beforeEach(clearDatabase);

  it('should throw an error if there is no fields or identifiers options.', () => {
    try {
      modelFactory();

      fail();
    } catch (e) {
      expect(e.message).toBe('Fields and identifiers options must be defined.');
    }
  });

  describe('$beforeInsert', () => {
    it('should throw a `ValidationError` with the unique fields that are already used.', async () => {
      const TestModel = modelFactory({
        fields: ['bar', 'foo']
      });

      await TestModel.query().insert({ bar: 'bar', foo: 'foo' });

      try {
        await TestModel.query().insert({ bar: 'bar', foo: 'foo' });

        fail();
      } catch (e) {
        expect(e).toBeInstanceOf(ValidationError);
        expect(e.message).toEqual('Unique Validation Failed');
        expect(e.type).toEqual('ModelValidation');
        expect(e.data).toEqual({
          bar: [{
            keyword: 'unique',
            message: 'bar already in use.'
          }],
          foo: [{
            keyword: 'unique',
            message: 'foo already in use.'
          }]
        });
      }
    });

    it('should insert the given data ignoring null values validation.', async () => {
      const TestModel = modelFactory({
        fields: ['bar', 'foo']
      });

      await TestModel.query().insert({ bar: 'bar', biz: 'biz', foo: null });

      const { id } = await TestModel.query().insert({ bar: 'buz', foo: null });
      const result = await TestModel.query().findById(id);

      expect(result).toEqual({ bar: 'buz', biz: null, foo: null, id });
    });

    it('should insert the given data.', async () => {
      const TestModel = modelFactory({
        fields: ['bar', 'foo']
      });

      const { id } = await TestModel.query().insert({ bar: 'bar', biz: 'biz', foo: 'foo' });
      const result = await TestModel.query().findById(id);

      expect(result).toEqual({ bar: 'bar', biz: 'biz', foo: 'foo', id });
    });

    it('should handle null values in compound keys', async () => {
      const TestModel = modelFactory({
        fields: [['bar', 'foo']]
      });

      const { id } = await TestModel.query().insert({ bar: 'bar' });
      const result = await TestModel.query().findById(id);

      expect(result).toEqual({ bar: 'bar', biz: null, foo: null, id: 1 });
    });

    it('should favor transaction from context', async () => {
      const TestModel = modelFactory({
        fields: [['bar', 'foo']]
      });

      const result = await TestModel.knex().transaction(async trx => {
        const { id } = await TestModel.query(trx).insert({ bar: 'bar', biz: 'biz', foo: 'foo' });
        const result = await TestModel.query(trx).findById(id);

        return result;
      });

      expect(result).toEqual({ bar: 'bar', biz: 'biz', foo: 'foo', id: 1 });
    });
  });

  describe('$beforeUpdate', () => {
    it('should throw an error if update is not a $query method.', async () => {
      const TestModel = modelFactory({
        fields: ['bar', 'foo']
      });

      try {
        await TestModel.query().update({});

        fail();
      } catch (e) {
        expect(e.message).toBe('Unique validation at update only works with queries started with $query.');
      }
    });

    it('should throw a `ValidationError` with the unique fields that are already used.', async () => {
      const TestModel = modelFactory({
        fields: ['bar', 'foo']
      });

      await TestModel.query().insert({ bar: 'bar', foo: 'foo' });

      const result = await TestModel.query().insertAndFetch({ bar: 'biz', foo: 'buz' });

      try {
        await result.$query().update({ bar: 'bar', foo: 'foo' });

        fail();
      } catch (e) {
        expect(e).toBeInstanceOf(ValidationError);
        expect(e.message).toEqual('Unique Validation Failed');
        expect(e.type).toEqual('ModelValidation');
        expect(e.data).toEqual({
          bar: [{
            keyword: 'unique',
            message: 'bar already in use.'
          }],
          foo: [{
            keyword: 'unique',
            message: 'foo already in use.'
          }]
        });
      }
    });

    it('should throw a `ValidationError` for the correct field when patching.', async () => {
      const TestModel = modelFactory({
        fields: ['bar', 'foo']
      });

      await TestModel.query().insert({ bar: 'bar', foo: 'foo' });

      const result = await TestModel.query().insertAndFetch({ bar: 'biz', foo: 'buz' });

      try {
        await result.$query().patch({ foo: 'foo' });

        fail();
      } catch (e) {
        expect(e).toBeInstanceOf(ValidationError);
        expect(e.message).toEqual('Unique Validation Failed');
        expect(e.type).toEqual('ModelValidation');
        expect(e.data).toEqual({
          foo: [{
            keyword: 'unique',
            message: 'foo already in use.'
          }]
        });
      }
    });

    it('should update the entity ignoring the unique validation if the values are from the same entity that are begin updated.', async () => {
      const TestModel = modelFactory({
        fields: ['bar', 'foo']
      });

      let result = await TestModel.query().insertAndFetch({ bar: 'biz', foo: 'buz' });

      result = await result.$query().updateAndFetch({ bar: 'biz', biz: 'foo', foo: 'buz' });

      expect(result).toEqual({ ...result, biz: 'foo' });
    });

    it('should update the entity ignoring null values validation.`', async () => {
      const TestModel = modelFactory({
        fields: ['bar', 'foo']
      });

      let result = await TestModel.query().insert({ bar: 'bar', foo: null });

      result = await result.$query().updateAndFetch({ bar: 'biz', biz: 'waldo', foo: null });

      expect(result).toEqual({ bar: 'biz', biz: 'waldo', foo: null, id: result.id });
    });

    it('should update the entity.`', async () => {
      const TestModel = modelFactory({
        fields: ['bar', 'foo']
      });

      let result = await TestModel.query().insertAndFetch({ bar: 'bar', foo: 'foo' });

      result = await result.$query().patchAndFetch({ bar: 'biz' });

      expect(result.bar).toBe('biz');
    });

    it('when applied to multiple fields should create and update entity.', async () => {
      const CompoundTestModel = compoundModelFactory({
        fields: [['bar', 'foo']]
      });

      await CompoundTestModel.query().insertAndFetch({ bar: 'bar', foo: 'foo' });
      let result = await CompoundTestModel.query().insertAndFetch({ bar: 'foo', foo: 'bar' });

      result = await result.$query().patchAndFetch({ bar: 'biz' });

      expect(result.bar).toBe('biz');
    });

    it('when applied to multiple fields should throw a `ValidationError` for all fields when patching.', async () => {
      const CompoundTestModel = compoundModelFactory({
        fields: [['bar', 'foo']]
      });

      await CompoundTestModel.query().insert({ bar: 'arg', foo: 'gar' });

      const result = await CompoundTestModel.query().insertAndFetch({ bar: 'biz', foo: 'gar' });

      try {
        await result.$query().patch({ bar: 'arg' });

        fail();
      } catch (e) {
        expect(e).toBeInstanceOf(ValidationError);
        expect(e.message).toEqual('Unique Validation Failed');
        expect(e.type).toEqual('ModelValidation');
        expect(e.data).toEqual({
          bar: [{
            keyword: 'unique',
            message: 'bar,foo already in use.'
          }],
          foo: [{
            keyword: 'unique',
            message: 'bar,foo already in use.'
          }]
        });
      }
    });
  });
});
