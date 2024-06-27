import {
  OrderSide, OrderStatus, OrderType, PositionStatus, TickerType,
} from '../constants';
import { Data } from '../types';
import RestClient from './rest';

/**
 * @description REST endpoints for data related to a particular address.
 */
export default class AccountClient extends RestClient {
  async getSubaccounts(
    address: string,
    limit?: number,
  ): Promise<Data> {
    const uri = `/v4/addresses/${address}`;
    return this.get(uri, { limit });
  }

  async getSubaccount(
    address: string,
    subaccountNumber: number,
  ): Promise<Data> {
    const uri = `/v4/addresses/${address}/subaccountNumber/${subaccountNumber}`;
    return this.get(uri);
  }

  async getSubaccountPerpetualPositions(
    address: string,
    subaccountNumber: number,
    status?: PositionStatus | null,
    limit?: number | null,
    createdBeforeOrAtHeight?: number | null,
    createdBeforeOrAt?: string | null,
  ): Promise<Data> {
    const uri = '/v4/perpetualPositions';
    return this.get(uri, {
      address,
      subaccountNumber,
      status,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
    });
  }

  async getSubaccountAssetPositions(
    address: string,
    subaccountNumber: number,
    status?: PositionStatus | null,
    limit?: number | null,
    createdBeforeOrAtHeight?: number | null,
    createdBeforeOrAt?: string | null,
  ): Promise<Data> {
    const uri = '/v4/assetPositions';
    return this.get(uri, {
      address,
      subaccountNumber,
      status,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
    });
  }

  async getSubaccountTransfers(
    address: string,
    subaccountNumber: number,
    limit?: number | null,
    createdBeforeOrAtHeight?: number | null,
    createdBeforeOrAt?: string | null,
  ): Promise<Data> {
    const uri = '/v4/transfers';
    return this.get(uri, {
      address,
      subaccountNumber,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
    });
  }

  async getSubaccountOrders(
    address: string,
    subaccountNumber: number,
    ticker?: string | null,
    tickerType: TickerType = TickerType.PERPETUAL,
    side?: OrderSide | null,
    status?: OrderStatus | null,
    type?: OrderType | null,
    limit?: number | null,
    goodTilBlockBeforeOrAt?: number | null,
    goodTilBlockTimeBeforeOrAt?: string | null,
    returnLatestOrders?: boolean | null,
  ): Promise<Data> {
    const uri = '/v4/orders';
    return this.get(uri, {
      address,
      subaccountNumber,
      ticker,
      tickerType,
      side,
      status,
      type,
      limit,
      goodTilBlockBeforeOrAt,
      goodTilBlockTimeBeforeOrAt,
      returnLatestOrders,
    });
  }

  async getOrder(orderId: string) : Promise<Data> {
    const uri = `/v4/orders${orderId}`;
    return this.get(uri);
  }

  async getSubaccountFills(
    address: string,
    subaccountNumber: number,
    ticker?: string | null,
    tickerType: TickerType = TickerType.PERPETUAL,
    limit?: number | null,
    createdBeforeOrAtHeight?: number | null,
    createdBeforeOrAt?: string | null,
  ): Promise<Data> {
    const uri = '/v4/fills';
    return this.get(uri, {
      address,
      subaccountNumber,
      ticker,
      tickerType,
      limit,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
    });
  }

  async getSubaccountHistoricalPNLs(
    address: string,
    subaccountNumber: number,
    effectiveBeforeOrAt?: string | null,
    effectiveAtOrAfter?: string | null,
  ): Promise<Data> {
    const uri = '/v4/historical-pnl';
    return this.get(uri, {
      address,
      subaccountNumber,
      effectiveBeforeOrAt,
      effectiveAtOrAfter,
    });
  }
}
