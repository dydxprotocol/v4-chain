import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import {
  setupBaseQuery,
  verifyAllRequiredFields,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getSubaccountQueryForParent } from '../lib/parent-subaccount-helpers';
import PnlModel from '../models/pnl-model';
import {
  Options,
  Ordering,
  PnlFromDatabase,
  PaginationFromDatabase,
  QueryableField,
  QueryConfig,
  PnlQueryConfig,
} from '../types';
import { PnlColumns, PnlCreateObject } from '../types/pnl-types';

export async function findAll(
  {
    limit,
    subaccountId,
    createdAtHeight,
    createdAt,
    createdBeforeOrAtHeight,
    createdBeforeOrAt,
    createdOnOrAfterHeight,
    createdOnOrAfter,
    page,
    parentSubaccount,
  }: PnlQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<PnlFromDatabase>> {
  if (parentSubaccount !== undefined && subaccountId !== undefined) {
    throw new Error('Cannot specify both parentSubaccount and subaccountId');
  }

  verifyAllRequiredFields(
    {
      limit,
      subaccountId,
      createdAtHeight,
      createdAt,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      createdOnOrAfterHeight,
      createdOnOrAfter,
      page,
      parentSubaccount,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<PnlModel> = setupBaseQuery<PnlModel>(
    PnlModel,
    options,
  );

  if (subaccountId !== undefined) {
    baseQuery = baseQuery.whereIn(PnlColumns.subaccountId, subaccountId);
  } else if (parentSubaccount !== undefined) {
    baseQuery = baseQuery.whereIn(
      PnlColumns.subaccountId,
      getSubaccountQueryForParent(parentSubaccount),
    );
  }

  if (createdAt !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAt, createdAt);
  }

  if (createdAtHeight !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAtHeight, createdAtHeight);
  }

  if (createdBeforeOrAtHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlColumns.createdAtHeight,
      '<=',
      createdBeforeOrAtHeight,
    );
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (createdOnOrAfterHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlColumns.createdAtHeight,
      '>=',
      createdOnOrAfterHeight,
    );
  }

  if (createdOnOrAfter !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAt, '>=', createdOnOrAfter);
  }

  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      baseQuery = baseQuery.orderBy(column, order);
    }
  } else {
    baseQuery = baseQuery.orderBy(
      PnlColumns.subaccountId,
      Ordering.ASC,
    ).orderBy(
      PnlColumns.createdAtHeight,
      Ordering.DESC,
    );
  }

  return handleLimitAndPagination(baseQuery, limit, page);
}

export async function create(
  pnlToCreate: PnlCreateObject,
  options: Options = { txId: undefined },
): Promise<PnlFromDatabase> {
  return PnlModel.query(
    Transaction.get(options.txId),
  ).insert(pnlToCreate).returning('*');
}

export async function findById(
  subaccountId: string,
  createdAt: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PnlFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PnlModel> = setupBaseQuery<PnlModel>(
    PnlModel,
    options,
  );
  return baseQuery
    .where(PnlColumns.subaccountId, subaccountId)
    .where(PnlColumns.createdAt, createdAt)
    .first();
}

async function handleLimitAndPagination(
  baseQuery: QueryBuilder<PnlModel>,
  limit?: number,
  page?: number,
): Promise<PaginationFromDatabase<PnlFromDatabase>> {
  let query = baseQuery;

  /**
   * If a query is made using a page number, then the limit property is used as 'page limit'
   */
  if (page !== undefined && limit !== undefined) {
    /**
     * We make sure that the page number is always >= 1
     */
    const currentPage: number = Math.max(1, page);
    const offset: number = (currentPage - 1) * limit;

    /**
     * Ensure sorting is applied to maintain consistent pagination results.
     * Also a casting of the ts type is required since the infer of the type
     * obtained from the count is not performed.
     */
    const count: { count?: string } = (await query
      .clone()
      .clearOrder()
      .count({ count: '*' })
      .first()) as unknown as { count?: string };

    query = query.offset(offset).limit(limit);

    const results = (await query) as PnlFromDatabase[];
    return {
      results,
      limit,
      offset,
      total: parseInt(count.count ?? '0', 10),
    };
  }

  // If no pagination, just apply the limit
  if (limit !== undefined) {
    query = query.limit(limit);
  }

  const results = (await query) as PnlFromDatabase[];
  return {
    results,
  };
}

export async function findAllDailyPnl(
  {
    limit,
    subaccountId,
    createdBeforeOrAtHeight,
    createdBeforeOrAt,
    createdOnOrAfterHeight,
    createdOnOrAfter,
    page,
    parentSubaccount,
  }: PnlQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<PnlFromDatabase>> {
  if (parentSubaccount !== undefined && subaccountId !== undefined) {
    throw new Error('Cannot specify both parentSubaccount and subaccountId');
  }

  verifyAllRequiredFields(
    {
      limit,
      subaccountId,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      createdOnOrAfterHeight,
      createdOnOrAfter,
      page,
      parentSubaccount,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<PnlModel> = setupBaseQuery<PnlModel>(
    PnlModel,
    options,
  );

  if (subaccountId !== undefined) {
    baseQuery = baseQuery.whereIn(PnlColumns.subaccountId, subaccountId);
  } else if (parentSubaccount !== undefined) {
    baseQuery = baseQuery.whereIn(
      PnlColumns.subaccountId,
      getSubaccountQueryForParent(parentSubaccount),
    );
  }

  if (createdBeforeOrAtHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlColumns.createdAtHeight,
      '<=',
      createdBeforeOrAtHeight,
    );
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (createdOnOrAfterHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlColumns.createdAtHeight,
      '>=',
      createdOnOrAfterHeight,
    );
  }

  if (createdOnOrAfter !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAt, '>=', createdOnOrAfter);
  }

  // Create a CTE (Common Table Expression) to rank records within each day for each subaccount
  const knex = baseQuery.modelClass().knex();
  const dailyRankQuery = baseQuery.clone()
    .select('*')
    .select(
      knex.raw(
        `ROW_NUMBER() OVER (
        PARTITION BY "${PnlColumns.subaccountId}", DATE_TRUNC('day', "${PnlColumns.createdAt}" AT TIME ZONE 'UTC')
        ORDER BY "${PnlColumns.createdAtHeight}" DESC
      ) as row_num`,
      ),
    );

  // Create the outer query that selects only the top-ranked record for each day
  const finalQuery = setupBaseQuery<PnlModel>(PnlModel, options)
    .with('daily_ranked_pnl', dailyRankQuery)
    .from('daily_ranked_pnl')
    .where('row_num', 1)
    .orderBy(PnlColumns.subaccountId, Ordering.ASC)
    .orderBy(PnlColumns.createdAtHeight, Ordering.DESC);

  // Apply pagination if needed
  return handleLimitAndPagination(finalQuery, limit, page);
}
