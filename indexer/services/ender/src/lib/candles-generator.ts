import { stats } from '@dydxprotocol-indexer/base';
import { CANDLES_WEBSOCKET_MESSAGE_VERSION, KafkaTopics } from '@dydxprotocol-indexer/kafka';
import {
  CANDLE_RESOLUTION_TO_PROTO,
  CandleColumns,
  CandleCreateObject,
  CandleFromDatabase,
  CandleMessageContents,
  CandleResolution,
  CandleTable,
  CandleUpdateObject,
  MarketOpenInterest,
  NUM_SECONDS_IN_CANDLE_RESOLUTIONS,
  Options,
  PerpetualMarketColumns,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualPositionTable,
  TradeContent,
  TradeMessageContents,
} from '@dydxprotocol-indexer/postgres';
import { CandleMessage } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import { DateTime } from 'luxon';

import { getCandle } from '../caches/candle-cache';
import { getOrderbookMidPrice } from '../caches/orderbook-mid-price-memory-cache';
import config from '../config';
import { KafkaPublisher } from './kafka-publisher';
import { ConsolidatedKafkaEvent, SingleTradeMessage } from './types';

type BlockCandleUpdatesMap = { [ticker: string]: BlockCandleUpdate };
type BlockCandleUpdate = {
  low: string,
  high: string,
  open: string,
  close: string,
  baseTokenVolume: string,
  usdVolume: string,
  trades: number,
};

type OrderbookMidPrice = string | undefined;
type OpenInterestMap = { [ticker: string]: string };

const utcZone = {
  zone: 'utc',
};

export class CandlesGenerator {
  kafkaPublisher: KafkaPublisher;
  blockTimestamp: DateTime;
  txId: number;
  writeOptions: Options;
  resolutionStartTimes: Map<CandleResolution, DateTime>;

  constructor(
    kafkaPublisher: KafkaPublisher,
    blockTimestamp: DateTime,
    txId: number,
  ) {
    this.kafkaPublisher = kafkaPublisher;
    this.blockTimestamp = blockTimestamp;
    this.txId = txId;
    this.writeOptions = { txId: this.txId };
    this.resolutionStartTimes = new Map(
      Object.values(CandleResolution).map((resolution: CandleResolution) => {
        return [
          resolution,
          CandlesGenerator.calculateNormalizedCandleStartTime(this.blockTimestamp, resolution),
        ];
      }));
  }

  /**
   * Update all Candles in postgres and add all candles kafka updates to KafkaPublisher
   * 1. Sort all trade messages by order in the block
   * 2. Iterate through all trade messages in order to generate all cumulative updates to made
   *    and store updates in BlockCandleUpdate.
   * 3. Update and Create Candles in postgres
   * 4. Add Candle kafka messages to KafkaPublisher
   */
  public async updateCandles(): Promise<CandleFromDatabase[]> {
    const start: number = Date.now();
    // 1. Sort trade messages by order in the block
    this.kafkaPublisher.sortEvents(this.kafkaPublisher.tradeMessages);

    // 2. Generate BlockCandleUpdatesMap
    const blockCandleUpdatesMap: BlockCandleUpdatesMap = this.generateBlockCandleUpdatesMap();

    // 3. Update and Create Candles in postgres
    const candles: CandleFromDatabase[] = await this.createOrUpdatePostgresCandles(
      blockCandleUpdatesMap,
    );

    // 4. Add Candle kafka messages to KafkaPublisher
    this.addCandleKafkaMessages(candles);

    stats.timing(
      `${config.SERVICE_NAME}.update_candles.timing`,
      Date.now() - start,
    );
    return candles;
  }

  private generateBlockCandleUpdatesMap(): BlockCandleUpdatesMap {
    const tradeMessages = this.kafkaPublisher.tradeMessages;
    const blockCandleUpdatesMap: BlockCandleUpdatesMap = {};
    _.forEach(tradeMessages, (tradeMessage: SingleTradeMessage) => {
      const ticker: string | undefined = perpetualMarketRefresher.getPerpetualMarketTicker(
        tradeMessage.clobPairId,
      );
      if (ticker === undefined) {
        throw Error(`Could not find ticker for clobPairId: ${tradeMessage.clobPairId}`);
      }
      // There should only be a single trade in SingleTradeMessage
      const contents: TradeMessageContents = JSON.parse(tradeMessage.contents);
      const tradeContent: TradeContent = contents.trades[0];
      if (ticker in blockCandleUpdatesMap) {
        blockCandleUpdatesMap[ticker] = this.getUpdatedBlockCandleUpdate(
          blockCandleUpdatesMap[ticker],
          tradeContent,
        );
      } else {
        blockCandleUpdatesMap[ticker] = this.createBlockCandleUpdate(tradeContent);
      }
    });

    return blockCandleUpdatesMap;
  }

  /**
   * Return blockCandleUpdate updated with tradeContent representing a single trade
   */
  private getUpdatedBlockCandleUpdate(
    blockCandleUpdate: BlockCandleUpdate,
    tradeContent: TradeContent,
  ) {
    return {
      low: Big(tradeContent.price).lt(blockCandleUpdate.low)
        ? tradeContent.price
        : blockCandleUpdate.low,
      high: Big(tradeContent.price).gt(blockCandleUpdate.high)
        ? tradeContent.price
        : blockCandleUpdate.high,
      open: blockCandleUpdate.open,
      close: tradeContent.price,
      baseTokenVolume: Big(blockCandleUpdate.baseTokenVolume).plus(
        tradeContent.size,
      ).toFixed(),
      usdVolume: Big(blockCandleUpdate.usdVolume).plus(
        Big(tradeContent.price).times(tradeContent.size),
      ).toFixed(),
      trades: blockCandleUpdate.trades + 1,
    };
  }

  /**
   * Create a new BlockCandleUpdate with tradeMessage
   */
  private createBlockCandleUpdate(
    tradeContent: TradeContent,
  ): BlockCandleUpdate {
    return {
      low: tradeContent.price,
      high: tradeContent.price,
      open: tradeContent.price,
      close: tradeContent.price,
      baseTokenVolume: tradeContent.size,
      usdVolume: Big(tradeContent.price).times(tradeContent.size).toFixed(),
      trades: 1,
    };
  }

  private async createOrUpdatePostgresCandles(
    blockCandleUpdatesMap: BlockCandleUpdatesMap,
  ): Promise<CandleFromDatabase[]> {
    const start: number = Date.now();
    const promises: Promise<CandleFromDatabase | undefined>[] = [];

    const openInterestMap: OpenInterestMap = await this.getOpenInterestMap();
    const orderbookMidPriceMap = getOrderbookMidPriceMap();

    _.forEach(
      Object.values(perpetualMarketRefresher.getPerpetualMarketsMap()),
      (perpetualMarket: PerpetualMarketFromDatabase) => {
        const blockCandleUpdate: BlockCandleUpdate | undefined = blockCandleUpdatesMap[
          perpetualMarket.ticker
        ];

        _.forEach(
          Object.values(CandleResolution),
          (resolution: CandleResolution) => {
            promises.push(...this.createUpdateOrPassPostgresCandle(
              blockCandleUpdate,
              perpetualMarket.ticker,
              resolution,
              openInterestMap,
              orderbookMidPriceMap[perpetualMarket.ticker],
              this.resolutionStartTimes.get(resolution)!,
            ));
          },
        );
      },
    );

    const candles: CandleFromDatabase[] = _.compact(await Promise.all(promises));
    stats.timing(
      `${config.SERVICE_NAME}.update_postgres_candles.timing`,
      Date.now() - start,
    );
    return candles;
  }

  /**
   * Creates, updates, or does nothing for a ticker and resolution depending on the following cases
   * Cases:
   * - Candle doesn't exist & there is no block update - do nothing
   * - Candle doesn't exist & there is a block update - create candle
   * - Candle exists & !sameStartTime & there is a block update - create candle,
   *   update previous candle orderbookMidPriceClose
   * - Candle exists & !sameStartTime & there is no block update - create empty candle,
   *   update previous candle orderbookMidPriceClose
   * - Candle exists & sameStartTime & no block update - do nothing
   * - Candle exists & sameStartTime & block update - update candle
   *
   * The orderbookMidPriceClose/Open are updated for each candle at the start and end of
   * each resolution period.
   * Whenever we create a new candle we set the orderbookMidPriceClose/Open
   * If there is a previous candle & we're creating a new one (this occurs at the
   * beginning of a resolution period) set the previous candles orderbookMidPriceClose
   */
  private createUpdateOrPassPostgresCandle(
    blockCandleUpdate: BlockCandleUpdate | undefined,
    ticker: string,
    resolution: CandleResolution,
    openInterestMap: OpenInterestMap,
    orderbookMidPrice: OrderbookMidPrice,
    currentStartTime: DateTime,
  ): Promise<CandleFromDatabase | undefined>[] {

    const existingCandle: CandleFromDatabase | undefined = getCandle(
      ticker,
      resolution,
    );

    if (existingCandle === undefined) {
      // - Candle doesn't exist & there is no block update - do nothing
      if (blockCandleUpdate === undefined) {
        return [];
      }
      // - Candle doesn't exist & there is a block update - create candle
      return [this.createCandleInPostgres(
        currentStartTime,
        blockCandleUpdate,
        ticker,
        resolution,
        openInterestMap,
        orderbookMidPrice,
      )];
    }

    const sameStartTime: boolean = existingCandle.startedAt === currentStartTime.toISO();
    if (!sameStartTime) {
      // - Candle exists & !sameStartTime & there is a block update - create candle
      //   update previous candle orderbookMidPriceClose

      const previousCandleUpdate = this.updateCandleWithOrderbookMidPriceInPostgres(
        existingCandle,
        orderbookMidPrice,
      );

      if (blockCandleUpdate !== undefined) {
        return [previousCandleUpdate, this.createCandleInPostgres(
          currentStartTime,
          blockCandleUpdate,
          ticker,
          resolution,
          openInterestMap,
          orderbookMidPrice,
        )];
      }
      // - Candle exists & !sameStartTime & there is no block update - create empty candle
      //   update previous candle orderbookMidPriceClose/Open
      return [previousCandleUpdate, this.createEmptyCandleInPostgres(
        currentStartTime,
        ticker,
        resolution,
        openInterestMap,
        existingCandle,
        orderbookMidPrice,
      )];
    }
    if (blockCandleUpdate === undefined) {
      // - Candle exists & sameStartTime & no block update - do nothing
      return [];
    }
    // - Candle exists & sameStartTime & block update - update candle
    return [this.updateCandleInPostgres(
      existingCandle,
      blockCandleUpdate,
    )];
  }

  /**
   * Create map of ticker to open interest for all the tickers with newly created candles, and
   * therefore require open interest to be calculated.
   */
  // eslint-disable-next-line @typescript-eslint/require-await
  private async getOpenInterestMap(): Promise<OpenInterestMap> {
    const start: number = Date.now();
    const allTickers: string[] = _.map(
      Object.values(perpetualMarketRefresher.getPerpetualMarketsMap()),
      PerpetualMarketColumns.ticker,
    );

    const tickersToFetch: string[] = _.filter(allTickers, (ticker: string) => {
      return _.some(
        Object.values(CandleResolution),
        (resolution: CandleResolution) => {
          const startedAtISOString: string = this.resolutionStartTimes.get(resolution)!.toISO();
          const existingCandle: CandleFromDatabase | undefined = getCandle(ticker, resolution);
          return existingCandle === undefined || existingCandle.startedAt !== startedAtISOString;
        },
      );
    });

    const openInterestMap: OpenInterestMap = await this.getOpenInterestLongFromTickers(
      tickersToFetch,
    );
    stats.timing(
      `${config.SERVICE_NAME}.get_open_interest_map.timing`,
      Date.now() - start,
    );
    return openInterestMap;
  }

  /**
   * For each ticker in the tickers array, get the open interest long for that market prior to this
   * block.
   */
  private async getOpenInterestLongFromTickers(
    tickers: string[],
  ): Promise<{ [ticker: string]: string }> {
    const perpetualMarketIdsToFetch: string[] = _.map(tickers, (ticker: string) => {
      const perpetualMarket: PerpetualMarketFromDatabase = perpetualMarketRefresher
        .getPerpetualMarketFromTicker(ticker)!;

      return perpetualMarket.id;
    });

    // Do not utilize txId here because we want to get the open interest before any transactions in
    // the block are processed.
    const marketOpenInterestDictionary:
    _.Dictionary<MarketOpenInterest> = await PerpetualPositionTable.getOpenInterestLong(
      perpetualMarketIdsToFetch,
    );

    return _.chain(marketOpenInterestDictionary)
      .keyBy((marketOpenInterest: MarketOpenInterest) => {
        const perpetualMarket: PerpetualMarketFromDatabase = perpetualMarketRefresher
          .getPerpetualMarketFromId(marketOpenInterest.perpetualMarketId)!;

        return perpetualMarket.ticker;
      })
      .mapValues('openInterest')
      .value();
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  private async createCandleInPostgres(
    startedAt: DateTime,
    blockCandleUpdate: BlockCandleUpdate,
    ticker: string,
    resolution: CandleResolution,
    openInterestMap: OpenInterestMap,
    orderbookMidPrice: OrderbookMidPrice,
  ): Promise<CandleFromDatabase> {
    const candle: CandleCreateObject = {
      startedAt: startedAt.toISO(),
      ticker,
      resolution,
      low: blockCandleUpdate.low,
      high: blockCandleUpdate.high,
      open: blockCandleUpdate.open,
      close: blockCandleUpdate.close,
      baseTokenVolume: blockCandleUpdate.baseTokenVolume,
      usdVolume: blockCandleUpdate.usdVolume,
      trades: blockCandleUpdate.trades,
      startingOpenInterest: openInterestMap[ticker] ?? '0',
      orderbookMidPriceClose: orderbookMidPrice,
      orderbookMidPriceOpen: orderbookMidPrice,
    };

    return CandleTable.create(candle, this.writeOptions);
  }

  /**
   * Create initial candle for a ticker. This is used when there are no candle updates for the
   * block, and a new candle is needed.
   */
  // eslint-disable-next-line @typescript-eslint/require-await
  private async createEmptyCandleInPostgres(
    startedAt: DateTime,
    ticker: string,
    resolution: CandleResolution,
    openInterestMap: OpenInterestMap,
    existingCandle: CandleFromDatabase,
    orderbookMidPrice: OrderbookMidPrice,
  ): Promise<CandleFromDatabase> {
    const candle: CandleCreateObject = {
      startedAt: startedAt.toISO(),
      ticker,
      resolution,
      low: existingCandle.close,
      high: existingCandle.close,
      open: existingCandle.close,
      close: existingCandle.close,
      baseTokenVolume: '0',
      usdVolume: '0',
      trades: 0,
      startingOpenInterest: openInterestMap[ticker] ?? '0',
      orderbookMidPriceClose: orderbookMidPrice,
      orderbookMidPriceOpen: orderbookMidPrice,
    };

    return CandleTable.create(candle, this.writeOptions);
  }

  /**
   * Update an existing candle with a block candle update.
   */
  private async updateCandleInPostgres(
    existingCandle: CandleFromDatabase,
    blockCandleUpdate: BlockCandleUpdate,
  ): Promise<CandleFromDatabase> {
    if (existingCandle.trades === 0) {
      // If there are no trades in the existing candle, then we can just replace the candle with the
      // block candle update.
      return CandleTable.update(
        {
          id: existingCandle.id,
          low: blockCandleUpdate.low,
          high: blockCandleUpdate.high,
          open: blockCandleUpdate.open,
          close: blockCandleUpdate.close,
          baseTokenVolume: blockCandleUpdate.baseTokenVolume,
          usdVolume: blockCandleUpdate.usdVolume,
          trades: blockCandleUpdate.trades,
          orderbookMidPriceOpen: existingCandle.orderbookMidPriceOpen ?? undefined,
          orderbookMidPriceClose: existingCandle.orderbookMidPriceClose ?? undefined,
        },
        this.writeOptions,
      ) as Promise<CandleFromDatabase>;
    }

    const candle: CandleUpdateObject = {
      id: existingCandle.id,
      low: Big(existingCandle.low).lte(blockCandleUpdate.low)
        ? existingCandle.low
        : blockCandleUpdate.low,
      high: Big(existingCandle.high).gte(blockCandleUpdate.high)
        ? existingCandle.high
        : blockCandleUpdate.high,
      close: blockCandleUpdate.close,
      baseTokenVolume: Big(
        existingCandle.baseTokenVolume,
      ).plus(
        blockCandleUpdate.baseTokenVolume,
      ).toFixed(),
      usdVolume: Big(existingCandle.usdVolume).plus(blockCandleUpdate.usdVolume).toFixed(),
      trades: existingCandle.trades + blockCandleUpdate.trades,
      orderbookMidPriceClose: existingCandle.orderbookMidPriceClose ?? undefined,
      orderbookMidPriceOpen: existingCandle.orderbookMidPriceOpen ?? undefined,
    };

    return CandleTable.update(candle, this.writeOptions) as Promise<CandleFromDatabase>;
  }

  private async updateCandleWithOrderbookMidPriceInPostgres(
    existingCandle: CandleFromDatabase,
    orderbookMidPrice: OrderbookMidPrice,
  ): Promise<CandleFromDatabase> {

    const candle: CandleUpdateObject = {
      id: existingCandle.id,
      low: existingCandle.low,
      high: existingCandle.high,
      close: existingCandle.close,
      baseTokenVolume: existingCandle.baseTokenVolume,
      usdVolume: existingCandle.usdVolume,
      trades: existingCandle.trades,
      orderbookMidPriceOpen: existingCandle.orderbookMidPriceOpen ?? undefined,
      orderbookMidPriceClose: orderbookMidPrice,
    };

    return CandleTable.update(candle, this.writeOptions) as Promise<CandleFromDatabase>;
  }

  private addCandleKafkaMessages(
    candles: CandleFromDatabase[],
  ): void {
    _.forEach(candles, (candle: CandleFromDatabase) => {
      const candleMessageContents: CandleMessageContents = _.omit(
        candle,
        [CandleColumns.id],
      );
      const message: CandleMessage = {
        contents: JSON.stringify(candleMessageContents),
        clobPairId: perpetualMarketRefresher.getPerpetualMarketFromTicker(
          candle.ticker,
        )!.clobPairId,
        resolution: CANDLE_RESOLUTION_TO_PROTO[candle.resolution],
        version: CANDLES_WEBSOCKET_MESSAGE_VERSION,
      };
      const consolidatedKafkaEvent: ConsolidatedKafkaEvent = {
        topic: KafkaTopics.TO_WEBSOCKETS_CANDLES,
        message,
      };
      this.kafkaPublisher.addEvent(consolidatedKafkaEvent);
    });
  }

  static calculateNormalizedCandleStartTime(
    time: DateTime,
    resolution: CandleResolution,
  ): DateTime {
    const epochSeconds: number = Math.floor(time.toSeconds());
    const normalizedTimeSeconds: number = epochSeconds - (
      epochSeconds % NUM_SECONDS_IN_CANDLE_RESOLUTIONS[resolution]
    );

    return DateTime.fromSeconds(normalizedTimeSeconds, utcZone);
  }
}

/**
   * Get the cached orderbook mid price for a given ticker
*/
export function getOrderbookMidPriceMap(): { [ticker: string]: OrderbookMidPrice } {
  const start: number = Date.now();
  const perpetualMarkets = Object.values(perpetualMarketRefresher.getPerpetualMarketsMap());

  const priceMap: { [ticker: string]: OrderbookMidPrice } = {};
  perpetualMarkets.forEach((perpetualMarket: PerpetualMarketFromDatabase) => {
    priceMap[perpetualMarket.ticker] = getOrderbookMidPrice(perpetualMarket.ticker);
  });

  stats.timing(`${config.SERVICE_NAME}.get_orderbook_mid_price_map.timing`, Date.now() - start);
  return priceMap;
}
