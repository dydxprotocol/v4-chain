import Knex from 'knex';
import _ from 'lodash';
import { Model, ModelClass, QueryBuilderType } from 'objection';

import { RequiredFieldMissing } from '../lib/errors';
import {
  QueryConfig,
  QueryableField,
  Options,
} from '../types';
import { knexPrimary, knexReadReplica } from './knex';
import Transaction from './transaction';

export function verifyAllRequiredFields(
  query: QueryConfig,
  requiredFields: QueryableField[],
): void {
  const queryKeys: string[] = Object.keys(_.omitBy(query, _.isNil));
  requiredFields.forEach((field) => {
    if (!queryKeys.includes(field)) {
      throw new RequiredFieldMissing(field);
    }
  });
}

export function verifyAllInjectableVariables(vals: (string | number | undefined | null)[]) {
  vals.forEach((val) => {
    // Numbers and undefined are okay.
    if (
      val === undefined ||
      val === null ||
      typeof val === 'number'
    ) {
      return;
    }

    // Prevent possible query injections.
    if (
      typeof val === 'string' &&
      val.split(' ').length > 1
    ) {
      throw Error(`Invalid value: ${val} could be a malicious query injection`);
    }
  });
}

export function setupBaseQuery<T extends Model>(
  model: ModelClass<T>,
  options: Options,
): QueryBuilderType<T> {
  if (options.readReplica) {
    return model.bindKnex(knexReadReplica.getConnection()).query(
      Transaction.get(options.txId),
    );
  } else {
    return model.query(
      Transaction.get(options.txId),
    );
  }
}

export async function rawQuery(
  queryString: string,
  options: Options,
// eslint-disable-next-line @typescript-eslint/no-explicit-any
): Promise<Knex.Raw<any>> {
  const connection = options.readReplica ? knexReadReplica.getConnection() : knexPrimary;
  let queryBuilder = options.bindings === undefined
    ? connection.raw(queryString) : connection.raw(queryString, options.bindings);
  if (options.txId) {
    queryBuilder = queryBuilder.transacting(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        <Knex.Transaction<any, any>>Transaction.get(options.txId),
    );
  }
  if (options.sqlOptions) {
    queryBuilder = queryBuilder.options(options.sqlOptions);
  }
  return queryBuilder;
}

/* ------- Bulk Helpers ------- */

export function setBulkRowsForUpdate<T extends string>({
  objectArray,
  columns,
  stringColumns,
  numericColumns,
  bigintColumns,
  timestampColumns,
  uuidColumns,
  booleanColumns,
  binaryColumns,
  enumColumns,
}: {
  objectArray: Partial<Record<T, string | number | boolean | Buffer | null>>[],
  columns: T[],
  stringColumns?: T[],
  numericColumns?: T[],
  bigintColumns?: T[],
  timestampColumns?: T[],
  uuidColumns?: T[],
  booleanColumns?: T[],
  binaryColumns?: T[],
  enumColumns?: T[],
}): string[] {
  return objectArray.map((object) => columns.map((col) => {
    if (stringColumns && stringColumns.includes(col)) {
      return `'${castNull(object[col])}'`;
    }
    if (numericColumns && numericColumns.includes(col)) {
      return `${castNull(_.get(object, col))}`;
    }
    if (bigintColumns && bigintColumns.includes(col)) {
      return castValue(object[col] as number | undefined | null, 'bigint');
    }
    if (timestampColumns && timestampColumns.includes(col)) {
      return castValue(object[col] as string | undefined | null, 'timestamp');
    }
    if (uuidColumns && uuidColumns.includes(col)) {
      return castValue(object[col] as string | undefined | null, 'uuid');
    }
    if (binaryColumns && binaryColumns.includes(col)) {
      return castBinaryValue(object[col] as Buffer | null | undefined);
    }
    if (booleanColumns && booleanColumns.includes(col)) {
      return `${castNull(object[col])}`;
    }
    if (enumColumns && enumColumns.includes(col)) {
      return `${castEnumNull(object[col])}`;
    }
    throw new Error(`Unsupported column for bulk update: ${col}`);
  }).join(', '));
}

/**
 * If the value is null || undefined, return 'NULL'
 */
function castNull(
  value: Buffer | string | number | boolean | null | undefined,
): string {
  if (value === null || value === undefined) {
    return 'NULL';
  }
  return `${value}`;
}

/**
 * If the value is null || undefined, return 'NULL'
 */
function castEnumNull(
  value: Buffer | string | number | boolean | null | undefined,
): string {
  if (value === null || value === undefined) {
    return 'NULL';
  }
  return `'${value}'`;
}

/**
 * If the value is null || undefined, return 'NULL' casted with typesuffix, otherwise return
 * the stringified value in quotes casted with typesuffix.
 */
function castValue(
  value: string | number | boolean | null | undefined,
  typeSuffix: string,
): string {
  if (value === null || value === undefined) {
    return `NULL::${typeSuffix}`;
  }
  return `'${value}'::${typeSuffix}`;
}

function castBinaryValue(
  value: Buffer | null | undefined,
): string {
  const typeSuffix: string = 'bytea';
  if (value === null || value === undefined) {
    return `NULL::${typeSuffix}`;
  }
  return `'\\x${value.toString('hex')}'::${typeSuffix}`;
}

export function generateBulkUpdateString({
  table,
  objectRows,
  columns,
  isUuid,
  uniqueIdentifier = 'id',
  setFieldsToAppend,
}: {
  table: string,
  objectRows: string[],
  columns: string[],
  isUuid: boolean,
  uniqueIdentifier?: string,
  setFieldsToAppend?: string[],
}): string {
  const columnsToUpdate: string[] = _.without(columns, uniqueIdentifier);

  const setFields: string[] = columnsToUpdate.map((col) => {
    return `"${col}" = c."${col}"`;
  }).concat(setFieldsToAppend || []);
  return `
  UPDATE "${table}" SET
    ${setFields.join(',')}
  FROM (VALUES
    ${objectRows.map((object) => `(${object})`).join(', ')}
  ) AS c(${columns.map((c) => `"${c}"`).join(', ')})
  WHERE c."${uniqueIdentifier}"${isUuid ? '::uuid' : ''} = "${table}"."${uniqueIdentifier}";
`;
}

export function generateBulkUpsertString({
  table,
  objectRows,
  columns,
  uniqueIdentifiers = ['id'],
}: {
  table: string,
  objectRows: string[],
  columns: string[],
  uniqueIdentifiers?: string[],
}): string {
  const columnsToUpdate: string[] = _.without(columns, ...uniqueIdentifiers);

  const idFields: string = uniqueIdentifiers.map(
    (id: string): string => { return `"${id}"`; },
  ).join(',');
  const insertFields: string = columns.map(
    (column: string):string => { return `"${column}"`; },
  ).join(',');
  const setFields: string[] = columnsToUpdate.map((col) => {
    return `"${col}" = excluded."${col}"`;
  });

  return `
  INSERT INTO "${table}" (${insertFields}) VALUES
    ${objectRows.map((object) => `(${object})`).join(',')}
  ON CONFLICT (${idFields}) DO UPDATE SET ${setFields.join(',')};
  `;
}
