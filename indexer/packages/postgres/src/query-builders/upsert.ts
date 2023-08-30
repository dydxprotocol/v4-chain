import {
  Model,
  Page,
  QueryBuilder,
} from 'objection';

export default class UpsertQueryBuilder<M extends Model, R = M[]> extends QueryBuilder<M, R> {
  ArrayQueryBuilderType!: UpsertQueryBuilder<M, M[]>;
  SingleQueryBuilderType!: UpsertQueryBuilder<M, M>;
  NumberQueryBuilderType!: UpsertQueryBuilder<M, number>;
  PageQueryBuilderType!: UpsertQueryBuilder<M, Page<M>>;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  upsert(object: any) {
    const modelClass = this.modelClass();

    const idColumn: string = modelClass.idColumn as string;

    const tableDefinedId = `${modelClass.tableName}.${idColumn}`;

    const knex = modelClass.knex();

    const cols: string[] = Object.keys(object);
    const values: string[] = Object.values(object);

    const colBindings = cols.map(() => '??').join(', ');
    const valBindings = cols.map(() => '?').join(', ');
    const setBindings = cols.map(() => '?? = ?').join(', ');

    const setValues: string[] = [];
    for (let i = 0; i < cols.length; ++i) {
      setValues.push(cols[i], values[i]);
    }

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return this.onBuildKnex((query: any) => {
      query.insert(
        knex.raw(
          [
            `(${colBindings}) VALUES (${valBindings})`,
            'ON CONFLICT (??) DO',
            `UPDATE SET ${setBindings}`,
            'WHERE ?? = ?',
          ].join(' '),
          [
            ...cols,
            ...values,
            modelClass.idColumn,
            ...setValues,
            tableDefinedId,
            object[idColumn],
          ],
        ),
      );
    });
  }
}
