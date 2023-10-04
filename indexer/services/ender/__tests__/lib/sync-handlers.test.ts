import { SyncHandlers } from '../../src/lib/sync-handlers';
import {
  AssetCreateEventV1,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  MarketEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import {
  defaultAssetCreateEvent,
  defaultDateTime,
  defaultMarketCreate,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { AssetCreationHandler } from '../../src/handlers/asset-handler';
import { MarketCreateHandler } from '../../src/handlers/markets/market-create-handler';
import {
  AssetColumns,
  AssetFromDatabase,
  AssetTable,
  BlockTable,
  dbHelpers,
  MarketColumns,
  MarketFromDatabase,
  MarketTable,
  Ordering,
  TendermintEventTable,
  Transaction,
} from '@dydxprotocol-indexer/postgres';
import { KafkaPublisher } from '../../src/lib/kafka-publisher';

const defaultMarketEventBinary: Uint8Array = Uint8Array.from(MarketEventV1.encode(
  defaultMarketCreate,
).finish());

const defaultAssetEventBinary: Uint8Array = Uint8Array.from(AssetCreateEventV1.encode(
  defaultAssetCreateEvent,
).finish());

describe('syncHandler', () => {
  const defaultTransactionIndex: number = 0;
  const defaultMarketEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.MARKET,
    defaultMarketEventBinary,
    defaultTransactionIndex,
    0,
  );

  const defaultAssetEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.ASSET,
    defaultAssetEventBinary,
    defaultTransactionIndex,
    1,
  );

  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    1,
    defaultTime,
    [defaultMarketEvent, defaultAssetEvent],
    [defaultTxHash],
  );

  describe('addHandler/process', () => {
    beforeEach(async () => {
      await BlockTable.create({
        blockHeight: '1',
        time: defaultDateTime.toISO(),
      });
      await Promise.all([
        TendermintEventTable.create({
          blockHeight: '1',
          transactionIndex: defaultTransactionIndex,
          eventIndex: 0,
        }),
        TendermintEventTable.create({
          blockHeight: '1',
          transactionIndex: defaultTransactionIndex,
          eventIndex: 1,
        }),
      ]);
    });

    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('successfully adds handler', async () => {
      const synchHandlers: SyncHandlers = new SyncHandlers();
      const txId: number = await Transaction.start();
      const assetHandler = new AssetCreationHandler(
        block,
        defaultAssetEvent,
        txId,
        defaultAssetCreateEvent,
      );

      const marketHandler = new MarketCreateHandler(
        block,
        defaultMarketEvent,
        txId,
        defaultMarketCreate,
      );
      // handlers are processed in the order in which they are received.
      synchHandlers.addHandler(DydxIndexerSubtypes.MARKET, marketHandler);
      synchHandlers.addHandler(DydxIndexerSubtypes.ASSET, assetHandler);
      // should be ignored, because transfers are not handled by syncHandlers
      synchHandlers.addHandler(DydxIndexerSubtypes.TRANSFER, assetHandler);

      const assets: AssetFromDatabase[] = await AssetTable.findAll(
        {},
        [], {
          orderBy: [[AssetColumns.id, Ordering.ASC]],
        });

      expect(assets.length).toEqual(0);
      const markets: MarketFromDatabase[] = await MarketTable.findAll(
        {},
        [], {
          orderBy: [[MarketColumns.id, Ordering.ASC]],
        });

      expect(markets.length).toEqual(0);
      const kafkaPublisher: KafkaPublisher = new KafkaPublisher();
      await synchHandlers.process(kafkaPublisher);
      await Transaction.commit(txId);

      // check that assets/markets are populated
      const newAssets: AssetFromDatabase[] = await AssetTable.findAll(
        {},
        [], {
          orderBy: [[AssetColumns.id, Ordering.ASC]],
        });

      expect(newAssets.length).toEqual(1);
      const newMarkets: MarketFromDatabase[] = await MarketTable.findAll(
        {},
        [], {
          orderBy: [[MarketColumns.id, Ordering.ASC]],
        });

      expect(newMarkets.length).toEqual(1);
    });
  });
});
