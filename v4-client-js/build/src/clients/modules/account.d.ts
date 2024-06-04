import { OrderSide, OrderStatus, OrderType, PositionStatus, TickerType } from '../constants';
import { Data } from '../types';
import RestClient from './rest';
/**
 * @description REST endpoints for data related to a particular address.
 */
export default class AccountClient extends RestClient {
    getSubaccounts(address: string, limit?: number): Promise<Data>;
    getSubaccount(address: string, subaccountNumber: number): Promise<Data>;
    getSubaccountPerpetualPositions(address: string, subaccountNumber: number, status?: PositionStatus | null, limit?: number | null, createdBeforeOrAtHeight?: number | null, createdBeforeOrAt?: string | null): Promise<Data>;
    getSubaccountAssetPositions(address: string, subaccountNumber: number, status?: PositionStatus | null, limit?: number | null, createdBeforeOrAtHeight?: number | null, createdBeforeOrAt?: string | null): Promise<Data>;
    getSubaccountTransfers(address: string, subaccountNumber: number, limit?: number | null, createdBeforeOrAtHeight?: number | null, createdBeforeOrAt?: string | null): Promise<Data>;
    getSubaccountOrders(address: string, subaccountNumber: number, ticker?: string | null, tickerType?: TickerType, side?: OrderSide | null, status?: OrderStatus | null, type?: OrderType | null, limit?: number | null, goodTilBlockBeforeOrAt?: number | null, goodTilBlockTimeBeforeOrAt?: string | null, returnLatestOrders?: boolean | null): Promise<Data>;
    getOrder(orderId: string): Promise<Data>;
    getSubaccountFills(address: string, subaccountNumber: number, ticker?: string | null, tickerType?: TickerType, limit?: number | null, createdBeforeOrAtHeight?: number | null, createdBeforeOrAt?: string | null): Promise<Data>;
    getSubaccountHistoricalPNLs(address: string, subaccountNumber: number, effectiveBeforeOrAt?: string | null, effectiveAtOrAfter?: string | null): Promise<Data>;
}
