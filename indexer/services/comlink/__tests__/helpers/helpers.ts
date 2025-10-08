import {
  FillCreateObject,
  FillFromDatabase,
  FillTable,
  Liquidity,
  OrderCreateObject,
  OrderFromDatabase,
  OrderSide,
  OrderTable,
  PnlCreateObject,
  PnlTicksCreateObject,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import e from 'express';
import _ from 'lodash';
import request from 'supertest';

import IndexV4 from '../../src/controllers/api/index-v4';
import Server from '../../src/request-helpers/server';
import {
  RequestMethod,
  FillResponseObject,
  MarketType,
  PnlTicksResponseObject,
  PnlResponseObject,
} from '../../src/types';

const app: e.Express = Server(IndexV4);

export async function sendRequestToApp({
  type,
  path,
  body,
  errorMsg,
  expressApp,
  expectedStatus = 200,
  headers = {},
}: {
  type: RequestMethod,
  path: string,
  body?: {},
  errorMsg?: string,
  expressApp: e.Express,
  expectedStatus?: number,
  headers?: Record<string, string>,
}) {
  let req: request.Test;

  switch (type) {
    case RequestMethod.GET:
      req = request(expressApp).get(path);
      break;
    case RequestMethod.DELETE:
      req = request(expressApp).delete(path);
      break;
    case RequestMethod.POST:
      req = request(expressApp).post(path);
      break;
    case RequestMethod.PUT:
      req = request(expressApp).put(path);
      break;
    default:
      throw new Error(`Invalid type of request: ${type}`);
  }

  if (Object.keys(headers).length) {
    await req.set(headers);
  }

  const response: request.Response = await req.send(body);
  if (response.status !== expectedStatus) {
    console.log(response.body); // eslint-disable-line no-console
  }
  expect(response.status).toEqual(expectedStatus);
  if (errorMsg) {
    expect(response.body.errors[0].msg).toContain(errorMsg);
  }

  return response;
}

export async function sendRequest({
  type,
  path,
  body,
  errorMsg,
  expectedStatus = 200,
  headers,
}: {
  type: RequestMethod,
  path: string,
  body?: {},
  errorMsg?: string,
  expectedStatus?: number,
  headers?: Record<string, string>,
}) {
  return sendRequestToApp({
    type,
    path,
    body,
    errorMsg,
    expressApp: app,
    expectedStatus,
    headers,
  });
}

export function getQueryString(
  params: {[name: string]: string | number | string[] | undefined},
): string {
  const queryStrings: string[] = [];
  _.forOwn(params, (value: string | number | string[] | undefined, key: string): void => {
    if (Array.isArray(value)) {
      const commaSeparatedList: string = value.join(',');
      queryStrings.push(`${key}=${commaSeparatedList}`);
    } else if (value !== undefined) {
      queryStrings.push(`${key}=${value}`);
    }
  });
  return queryStrings.join('&');
}

export function getFixedRepresentation(val: number | string): string {
  return new Big(val).toFixed();
}

export async function createMakerTakerOrderAndFill(
  order: OrderCreateObject,
  fill: FillCreateObject,
): Promise<{
  makerFill: FillFromDatabase,
  takerFill: FillFromDatabase,
}> {
  const makerOrder: OrderFromDatabase = await OrderTable.create({
    ...order,
    side: OrderSide.BUY,
    clientId: randomInt().toString(),
  });
  const makerFill: FillFromDatabase = await FillTable.create({
    ...fill,
    side: OrderSide.BUY,
    liquidity: Liquidity.MAKER,
    orderId: makerOrder.id,
  });
  const takerOrder: OrderFromDatabase = await OrderTable.create({
    ...order,
    side: OrderSide.SELL,
    clientId: randomInt().toString(),
  });
  const takerFill: FillFromDatabase = await FillTable.create({
    ...fill,
    side: OrderSide.SELL,
    liquidity: Liquidity.TAKER,
    orderId: takerOrder.id,
  });
  return { makerFill, takerFill };
}

export function fillResponseObjectFromFillCreateObject(
  fill: FillCreateObject,
  subaccountNumber: number,
): Partial<FillResponseObject> {
  const fillResponse: Partial<FillResponseObject> = {
    side: fill.side,
    liquidity: fill.liquidity,
    marketType: MarketType.PERPETUAL,
    price: fill.price,
    size: fill.size,
    fee: fill.fee,
    affiliateRevShare: fill.affiliateRevShare,
    type: fill.type,
    orderId: fill.orderId,
    createdAt: fill.createdAt,
    createdAtHeight: fill.createdAtHeight,
    subaccountNumber,
  };
  return fillResponse;
}

function randomInt(range: number = 1000): number {
  return Math.floor(Math.random() * range);
}

export function pnlTickCreateObjectToResponseObject(
  pnlTick: PnlTicksCreateObject,
): PnlTicksResponseObject {
  return {
    // Leave out subaccountId
    equity: pnlTick.equity,
    totalPnl: pnlTick.totalPnl,
    netTransfers: pnlTick.netTransfers,
    createdAt: pnlTick.createdAt,
    blockHeight: pnlTick.blockHeight,
    blockTime: pnlTick.blockTime,
  };
}

export function pnlCreateObjectToResponseObject(
  pnl: PnlCreateObject,
): PnlResponseObject {
  return {
    // Leave out subaccountId
    equity: pnl.equity,
    totalPnl: pnl.totalPnl,
    netTransfers: pnl.netTransfers,
    createdAt: pnl.createdAt,
    createdAtHeight: pnl.createdAtHeight,
  };
}
