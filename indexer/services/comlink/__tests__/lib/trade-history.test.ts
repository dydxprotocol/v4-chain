import {
  FillFromDatabase,
  FillType,
  Liquidity,
  OrderSide,
  OrderType,
  PerpetualMarketType,
  PositionSide,
} from '@dydxprotocol-indexer/postgres';

import { computeTradeHistory, paginateTradeHistory } from '../../src/lib/trade-history';
import {
  MarketAndTypeByClobPairId, MarketType, TradeHistoryResponseObject, TradeHistoryType,
} from '../../src/types';

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

let fillCounter = 0;

function makeFill(overrides: Partial<FillFromDatabase> = {}): FillFromDatabase {
  fillCounter += 1;
  return {
    id: `fill-${fillCounter}`,
    subaccountId: 'sub-1',
    side: OrderSide.BUY,
    liquidity: Liquidity.TAKER,
    type: FillType.LIMIT,
    clobPairId: '0',
    size: '1',
    price: '100',
    quoteAmount: '100',
    eventId: Buffer.from(`e${fillCounter}`, 'utf-8'),
    transactionHash: `tx-${fillCounter}`,
    createdAt: `2024-01-01T00:0${fillCounter}:00.000Z`,
    createdAtHeight: `${fillCounter}`,
    fee: '0.1',
    affiliateRevShare: '0',
    orderId: `order-${fillCounter}`,
    ...overrides,
  };
}

const MARKET_MAP: MarketAndTypeByClobPairId = {
  0: { market: 'BTC-USD', marketType: MarketType.PERPETUAL, perpetualMarketType: PerpetualMarketType.CROSS },
  1: { market: 'ETH-USD', marketType: MarketType.PERPETUAL, perpetualMarketType: PerpetualMarketType.CROSS },
};

const ORDER_TYPE_MAP: Record<string, OrderType> = {
  'order-1': OrderType.LIMIT,
  'order-2': OrderType.MARKET,
  'order-3': OrderType.LIMIT,
  'order-4': OrderType.STOP_LIMIT,
  'order-5': OrderType.LIMIT,
  'order-6': OrderType.LIMIT,
  'order-7': OrderType.LIMIT,
  'order-8': OrderType.LIMIT,
};

beforeEach(() => {
  fillCounter = 0;
});

// ---------------------------------------------------------------------------
// computeTradeHistory
// ---------------------------------------------------------------------------

describe('computeTradeHistory', () => {
  it('returns empty array for empty fills', () => {
    const result = computeTradeHistory([], ORDER_TYPE_MAP, MARKET_MAP);
    expect(result).toEqual([]);
  });

  it('OPEN: buying from flat produces an OPEN row', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY, size: '5', price: '100', fee: '0.5', orderId: 'order-1',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(1);
    expect(result[0].action).toBe(TradeHistoryType.OPEN);
    expect(result[0].side).toBe(OrderSide.BUY);
    expect(result[0].prevSize).toBe('0');
    expect(result[0].additionalSize).toBe('5');
    expect(result[0].executionPrice).toBe('100');
    expect(result[0].value).toBe('500');
    expect(result[0].orderType).toBe(OrderType.LIMIT);
    expect(result[0].netRealizedPnl).toBe('0');
    expect(result[0].netFee).toBe('0.5');
    expect(result[0].marketId).toBe('BTC-USD');
    expect(result[0].positionSide).toBe(PositionSide.LONG);
    expect(result[0].orderId).toBe('order-1');
    expect(result[0].id).toBe('order-1');
  });

  it('EXTEND: buying more when already long produces EXTEND', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY, size: '5', price: '100', fee: '0.5', orderId: 'order-1',
      }),
      makeFill({
        side: OrderSide.BUY, size: '5', price: '110', fee: '0.5', orderId: 'order-2',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(2);
    // Most recent first (sorted DESC)
    expect(result[0].action).toBe(TradeHistoryType.EXTEND);
    expect(result[0].prevSize).toBe('5');
    expect(result[0].additionalSize).toBe('5');
    expect(result[0].executionPrice).toBe('110');
    expect(result[0].netRealizedPnl).toBe('0');
    expect(result[0].netFee).toBe('1'); // 0.5 + 0.5

    expect(result[1].action).toBe(TradeHistoryType.OPEN);
  });

  it('PARTIAL_CLOSE: selling part of a long position', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY, size: '10', price: '100', fee: '1', orderId: 'order-1',
      }),
      makeFill({
        side: OrderSide.SELL, size: '5', price: '120', fee: '0.5', orderId: 'order-2',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(2);
    const partialClose = result[0]; // most recent
    expect(partialClose.action).toBe(TradeHistoryType.PARTIAL_CLOSE);
    expect(partialClose.prevSize).toBe('10');
    expect(partialClose.additionalSize).toBe('-5');
    expect(partialClose.executionPrice).toBe('120');
    // PnL = (120 - 100) * 5 = 100
    expect(partialClose.netRealizedPnl).toBe('100');
    expect(partialClose.netFee).toBe('1.5'); // 1 + 0.5
  });

  it('CLOSE: fully closing a long position with realized PnL', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY, size: '5', price: '100', fee: '0.5', orderId: 'order-1',
      }),
      makeFill({
        side: OrderSide.SELL, size: '5', price: '150', fee: '0.5', orderId: 'order-2',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(2);
    const close = result[0];
    expect(close.action).toBe(TradeHistoryType.CLOSE);
    expect(close.prevSize).toBe('5');
    expect(close.additionalSize).toBe('-5');
    // PnL = (150 - 100) * 5 = 250
    expect(close.netRealizedPnl).toBe('250');
    expect(close.netFee).toBe('1'); // 0.5 + 0.5
  });

  it('SHORT: open short and close with PnL', () => {
    const fills = [
      makeFill({
        side: OrderSide.SELL, size: '3', price: '200', fee: '0.3', orderId: 'order-1',
      }),
      makeFill({
        side: OrderSide.BUY, size: '3', price: '180', fee: '0.3', orderId: 'order-2',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(2);
    const open = result[1];
    expect(open.action).toBe(TradeHistoryType.OPEN);
    expect(open.side).toBe(OrderSide.SELL);
    expect(open.additionalSize).toBe('-3'); // SELL from flat
    expect(open.positionSide).toBe(PositionSide.SHORT);

    const close = result[0];
    expect(close.action).toBe(TradeHistoryType.CLOSE);
    expect(close.positionSide).toBeNull(); // fully closed
    // Short PnL = (200 - 180) * 3 = 60
    expect(close.netRealizedPnl).toBe('60');
  });

  it('cross-zero: single order that closes long and opens short', () => {
    const fills = [
      // Open Long 5
      makeFill({
        side: OrderSide.BUY, size: '5', price: '100', fee: '0.5', orderId: 'order-1',
      }),
      // Sell 10 → close 5 long, open 5 short
      makeFill({
        side: OrderSide.SELL, size: '10', price: '120', fee: '1', orderId: 'order-2',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(3); // OPEN + CLOSE + OPEN
    // Cross-zero rows share same time; sorted DESC with id tiebreaker:
    // order-2:open (new lifecycle) sorts before order-2:close (old lifecycle)
    expect(result[0].id).toBe('order-2:open');
    expect(result[1].id).toBe('order-2:close');

    const closeRow = result[1];
    expect(closeRow.action).toBe(TradeHistoryType.CLOSE);
    expect(closeRow.prevSize).toBe('5');
    expect(closeRow.additionalSize).toBe('-5');
    expect(closeRow.positionSide).toBeNull(); // fully closed
    // Close PnL = (120 - 100) * 5 = 100
    expect(closeRow.netRealizedPnl).toBe('100');
    // Close fee = 1 * (5/10) = 0.5, cumulative = 0.5 (open) + 0.5 = 1
    expect(closeRow.netFee).toBe('1');

    const openRow = result[0];
    expect(openRow.action).toBe(TradeHistoryType.OPEN);
    expect(openRow.prevSize).toBe('0');
    expect(openRow.additionalSize).toBe('-5'); // opening short, signed delta is negative
    expect(openRow.positionSide).toBe(PositionSide.SHORT); // new short position
    // After lifecycle reset, netRealizedPnl = 0
    expect(openRow.netRealizedPnl).toBe('0');
    // Open fee = 1 * (5/10) = 0.5, new lifecycle cumulative = 0.5
    expect(openRow.netFee).toBe('0.5');
  });

  it('cross-zero: single order that closes short and opens long (BUY direction)', () => {
    const fills = [
      // Open Short 5
      makeFill({
        side: OrderSide.SELL, size: '5', price: '200', fee: '0.5', orderId: 'order-1',
      }),
      // Buy 10 → close 5 short, open 5 long
      makeFill({
        side: OrderSide.BUY, size: '10', price: '180', fee: '1', orderId: 'order-2',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(3); // OPEN + CLOSE + OPEN

    const closeRow = result.find((r) => r.id === 'order-2:close')!;
    expect(closeRow.action).toBe(TradeHistoryType.CLOSE);
    expect(closeRow.side).toBe(OrderSide.BUY);
    expect(closeRow.prevSize).toBe('5');
    // BUY closing a short → additionalSize should be positive (BUY convention)
    expect(closeRow.additionalSize).toBe('5');
    expect(closeRow.positionSide).toBeNull();
    // Short PnL = (200 - 180) * 5 = 100
    expect(closeRow.netRealizedPnl).toBe('100');

    const openRow = result.find((r) => r.id === 'order-2:open')!;
    expect(openRow.action).toBe(TradeHistoryType.OPEN);
    expect(openRow.side).toBe(OrderSide.BUY);
    expect(openRow.prevSize).toBe('0');
    // BUY opening long → additionalSize should be positive
    expect(openRow.additionalSize).toBe('5');
    expect(openRow.positionSide).toBe(PositionSide.LONG);
    expect(openRow.netRealizedPnl).toBe('0');
  });

  it('liquidation fills produce LIQUIDATION_CLOSE with null orderId', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY, size: '5', price: '100', fee: '0.5', orderId: 'order-1',
      }),
      // Liquidation fill — no orderId
      makeFill({
        side: OrderSide.SELL,
        size: '5',
        price: '80',
        fee: '0.4',
        orderId: undefined,
        type: FillType.LIQUIDATED,
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(2);
    const liqRow = result[0];
    expect(liqRow.action).toBe(TradeHistoryType.LIQUIDATION_CLOSE);
    expect(liqRow.orderId).toBeNull();
    expect(liqRow.orderType).toBeNull();
    // Liquidation rows use fillId as their id
    expect(liqRow.id).toBe(fills[1].id);
    // PnL = (80 - 100) * 5 = -100
    expect(liqRow.netRealizedPnl).toBe('-100');
  });

  it('liquidation partial close', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY, size: '10', price: '100', fee: '1', orderId: 'order-1',
      }),
      makeFill({
        side: OrderSide.SELL,
        size: '3',
        price: '90',
        fee: '0.3',
        orderId: undefined,
        type: FillType.LIQUIDATED,
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(2);
    const liqRow = result[0];
    expect(liqRow.action).toBe(TradeHistoryType.LIQUIDATION_PARTIAL_CLOSE);
    expect(liqRow.prevSize).toBe('10');
    expect(liqRow.additionalSize).toBe('-3');
    // PnL = (90 - 100) * 3 = -30
    expect(liqRow.netRealizedPnl).toBe('-30');
  });

  it('cumulative PnL resets after full close and reopen', () => {
    const fills = [
      // Lifecycle 1: open and close
      makeFill({
        side: OrderSide.BUY, size: '5', price: '100', fee: '0.5', orderId: 'order-1',
      }),
      makeFill({
        side: OrderSide.SELL, size: '5', price: '120', fee: '0.5', orderId: 'order-2',
      }),
      // Lifecycle 2: open again
      makeFill({
        side: OrderSide.BUY, size: '3', price: '130', fee: '0.3', orderId: 'order-3',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(3);
    // result[0] = most recent = lifecycle 2 OPEN
    const newOpen = result[0];
    expect(newOpen.action).toBe(TradeHistoryType.OPEN);
    // After lifecycle reset, cumulative should be fresh
    expect(newOpen.netRealizedPnl).toBe('0');
    expect(newOpen.netFee).toBe('0.3');

    // result[1] = lifecycle 1 CLOSE
    const close = result[1];
    expect(close.action).toBe(TradeHistoryType.CLOSE);
    // PnL = (120 - 100) * 5 = 100
    expect(close.netRealizedPnl).toBe('100');
    expect(close.netFee).toBe('1'); // 0.5 + 0.5
  });

  it('multiple fills per orderId are grouped correctly', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY,
        size: '2',
        price: '100',
        fee: '0.2',
        orderId: 'order-1',
        createdAt: '2024-01-01T00:01:00.000Z',
      }),
      makeFill({
        side: OrderSide.BUY,
        size: '3',
        price: '110',
        fee: '0.3',
        orderId: 'order-1',
        createdAt: '2024-01-01T00:02:00.000Z',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(1);
    expect(result[0].action).toBe(TradeHistoryType.OPEN);
    // Total size = 2 + 3 = 5
    expect(result[0].additionalSize).toBe('5');
    // Weighted avg price = (100*2 + 110*3) / 5 = 530/5 = 106
    expect(result[0].executionPrice).toBe('106');
    // Total fee = 0.2 + 0.3 = 0.5
    expect(result[0].netFee).toBe('0.5');
    // value = 5 * 106 = 530
    expect(result[0].value).toBe('530');
  });

  it('multiple markets are processed independently', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY,
        size: '5',
        price: '100',
        fee: '0.5',
        orderId: 'order-1',
        clobPairId: '0',
      }),
      makeFill({
        side: OrderSide.SELL,
        size: '2',
        price: '3000',
        fee: '0.3',
        orderId: 'order-2',
        clobPairId: '1',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(2);
    const btcRow = result.find((r) => r.marketId === 'BTC-USD')!;
    const ethRow = result.find((r) => r.marketId === 'ETH-USD')!;

    expect(btcRow.action).toBe(TradeHistoryType.OPEN);
    expect(btcRow.additionalSize).toBe('5');

    expect(ethRow.action).toBe(TradeHistoryType.OPEN);
    expect(ethRow.side).toBe(OrderSide.SELL);
    expect(ethRow.additionalSize).toBe('-2');
  });

  it('entry price updates correctly on extend (weighted average)', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY, size: '4', price: '100', fee: '0.4', orderId: 'order-1',
      }),
      makeFill({
        side: OrderSide.BUY, size: '6', price: '150', fee: '0.6', orderId: 'order-2',
      }),
      // Now close at 200 to verify entry was weighted avg
      // Entry = (100*4 + 150*6) / 10 = 1300/10 = 130
      makeFill({
        side: OrderSide.SELL, size: '10', price: '200', fee: '1', orderId: 'order-3',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(3);
    const close = result[0];
    expect(close.action).toBe(TradeHistoryType.CLOSE);
    // PnL = (200 - 130) * 10 = 700
    expect(close.netRealizedPnl).toBe('700');
  });

  it('entry price stays the same on partial close', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY, size: '10', price: '100', fee: '1', orderId: 'order-1',
      }),
      // Partial close 5
      makeFill({
        side: OrderSide.SELL, size: '5', price: '120', fee: '0.5', orderId: 'order-2',
      }),
      // Close remaining 5 — entry should still be 100
      makeFill({
        side: OrderSide.SELL, size: '5', price: '140', fee: '0.5', orderId: 'order-3',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(3);
    const partialClose = result[1]; // middle
    // Partial PnL = (120 - 100) * 5 = 100
    expect(partialClose.netRealizedPnl).toBe('100');

    const fullClose = result[0]; // most recent
    // Full close PnL = (140 - 100) * 5 = 200, cumulative = 100 + 200 = 300
    expect(fullClose.netRealizedPnl).toBe('300');
  });

  it('orderType is null when orderId is not in the map', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY, size: '5', price: '100', fee: '0.5', orderId: 'unknown-id',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(1);
    expect(result[0].orderType).toBeNull();
  });

  it('different subaccountIds are processed independently (parent subaccount)', () => {
    const fills = [
      // Subaccount A: long 5 BTC
      makeFill({
        subaccountId: 'sub-0',
        side: OrderSide.BUY,
        size: '5',
        price: '100',
        fee: '0.5',
        orderId: 'order-1',
        clobPairId: '0',
      }),
      // Subaccount B (child 128): short 3 BTC
      makeFill({
        subaccountId: 'sub-128',
        side: OrderSide.SELL,
        size: '3',
        price: '100',
        fee: '0.3',
        orderId: 'order-2',
        clobPairId: '0',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    // Both should be OPEN (independent positions), not an OPEN + PARTIAL_CLOSE
    expect(result).toHaveLength(2);
    expect(result.every((r) => r.action === TradeHistoryType.OPEN)).toBe(true);
    const longRow = result.find((r) => r.side === OrderSide.BUY)!;
    const shortRow = result.find((r) => r.side === OrderSide.SELL)!;
    expect(longRow.additionalSize).toBe('5');
    expect(longRow.positionSide).toBe(PositionSide.LONG);
    expect(shortRow.additionalSize).toBe('-3');
    expect(shortRow.positionSide).toBe(PositionSide.SHORT);
  });

  it('skips fills for unknown markets', () => {
    const fills = [
      makeFill({
        side: OrderSide.BUY, size: '5', price: '100', fee: '0.5', clobPairId: '999',
      }),
    ];
    const result = computeTradeHistory(fills, ORDER_TYPE_MAP, MARKET_MAP);

    expect(result).toHaveLength(0);
  });
});

// ---------------------------------------------------------------------------
// paginateTradeHistory
// ---------------------------------------------------------------------------

describe('paginateTradeHistory', () => {
  const rows = Array.from({ length: 10 }, (_, i) => ({
    id: `order-${i}`,
    action: TradeHistoryType.OPEN,
    executionPrice: '100',
    side: OrderSide.BUY,
    positionSide: PositionSide.LONG,
    prevSize: '0',
    additionalSize: '1',
    value: '100',
    orderType: OrderType.LIMIT,
    netFee: '0',
    netRealizedPnl: '0',
    time: `2024-01-01T00:0${i}:00.000Z`,
    orderId: `order-${i}`,
    marketId: 'BTC-USD',
    marginMode: PerpetualMarketType.CROSS,
  })) as TradeHistoryResponseObject[];

  it('applies limit without page', () => {
    const result = paginateTradeHistory(rows, 3);
    expect(result.tradeHistory).toHaveLength(3);
    expect(result.totalResults).toBe(10);
    expect(result.offset).toBe(0);
    expect(result.pageSize).toBe(3);
  });

  it('applies page and limit together', () => {
    const result = paginateTradeHistory(rows, 3, 2);
    expect(result.tradeHistory).toHaveLength(3);
    expect(result.offset).toBe(3);
    expect(result.tradeHistory[0].orderId).toBe('order-3');
  });

  it('handles page beyond results', () => {
    const result = paginateTradeHistory(rows, 3, 100);
    expect(result.tradeHistory).toHaveLength(0);
    expect(result.totalResults).toBe(10);
  });

  it('page 1 is same as no page', () => {
    const result = paginateTradeHistory(rows, 5, 1);
    expect(result.tradeHistory).toHaveLength(5);
    expect(result.offset).toBe(0);
    expect(result.tradeHistory[0].orderId).toBe('order-0');
  });
});
