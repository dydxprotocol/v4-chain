import {
    YieldParamsColumns,
    YieldParamsFromDatabase,
    YieldParamsTable,
    dbHelpers,
    Ordering,
    liquidityTierRefresher,
    perpetualMarketRefresher,
    testMocks,
  } from '@dydxprotocol-indexer/postgres';
  import { updateBlockCache } from '../../src/caches/block-cache';
  import {
    defaultHeight,
    defaultPreviousHeight,
    defaultTime,
    defaultTxHash,
    defaultUpdateYieldParamsEvent1,
  } from '../helpers/constants';
  import {
    IndexerTendermintBlock,
    IndexerTendermintEvent,
    Timestamp,
    UpdateYieldParamsEventV1,
  } from '@dydxprotocol-indexer/v4-protos';
  import {
    createIndexerTendermintBlock,
    createIndexerTendermintEvent,
  } from '../helpers/indexer-proto-helpers';
  import { DydxIndexerSubtypes } from '../../src/lib/types';
  import { createKafkaMessage, producer } from '@dydxprotocol-indexer/kafka';
  import { KafkaMessage } from 'kafkajs';
  import { onMessage } from '../../src/lib/on-message';
  import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';
  import { YieldParamsHandler } from '../../src/handlers/yield-params-handler';
  
  describe('yield-params-handler', () => {
    beforeAll(async () => {
      await dbHelpers.migrate();
      await createPostgresFunctions();
    });
  
    beforeEach(async () => {
      await testMocks.seedData();
      updateBlockCache(defaultPreviousHeight);
      await perpetualMarketRefresher.updatePerpetualMarkets();
      await liquidityTierRefresher.updateLiquidityTiers();
    });
  
    afterEach(async () => {
      await dbHelpers.clearData();
      jest.clearAllMocks();
      perpetualMarketRefresher.clear();
      liquidityTierRefresher.clear();
    });
  
    afterAll(async () => {
      await dbHelpers.teardown();
      jest.resetAllMocks();
    });
  
    describe('getParallelizationIds', () => {
      it('returns the correct parallelization ids', () => {
        const transactionIndex: number = 0;
        const eventIndex: number = 0;
  
        const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
          DydxIndexerSubtypes.YIELD_PARAMS,
          UpdateYieldParamsEventV1.encode(defaultUpdateYieldParamsEvent1).finish(),
          transactionIndex,
          eventIndex,
        );
        const block: IndexerTendermintBlock = createIndexerTendermintBlock(
          0,
          defaultTime,
          [indexerTendermintEvent],
          [defaultTxHash],
        );

        const handler: YieldParamsHandler = new YieldParamsHandler(
            block,
            0,
            indexerTendermintEvent,
            0,
            defaultUpdateYieldParamsEvent1,
        )
  
        expect(handler.getParallelizationIds()).toEqual([]);
      });
    });

  it('successfully creates yield params', async () => {
    const transactionIndex: number = 0;

    const kafkaMessage: KafkaMessage = createKafkaMessageFromYieldParamsEvent({
      yieldParamsEvent: defaultUpdateYieldParamsEvent1,
      transactionIndex: transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await expectNoExistingYieldParams();

    await onMessage(kafkaMessage);

    const newYieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
      {},
      [], {
        orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
      });
    expect(newYieldParams.length).toEqual(1);
    expectYieldParamsMatchEvent(defaultUpdateYieldParamsEvent1, newYieldParams[0]);
    expectYieldParamsMatchBlock(defaultHeight, defaultTime, newYieldParams[0]);
  });
});
  
  function createKafkaMessageFromYieldParamsEvent({
    yieldParamsEvent,
    transactionIndex,
    height,
    time,
    txHash,
  }: {
    yieldParamsEvent: UpdateYieldParamsEventV1,
    transactionIndex: number,
    height: number,
    time: Timestamp,
    txHash: string,
  }) {
    const events: IndexerTendermintEvent[] = [];
    events.push(
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.YIELD_PARAMS,
        UpdateYieldParamsEventV1.encode(yieldParamsEvent).finish(),
        transactionIndex,
        0,
      ),
    );
  
    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      height,
      time,
      events,
      [txHash],
    );
  
    const binaryBlock: Uint8Array = IndexerTendermintBlock.encode(block).finish();
    return createKafkaMessage(Buffer.from(binaryBlock));
  }

  async function expectNoExistingYieldParams() {
    // Confirm there is no existing asset
    const assets: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
      {},
      [], {
        orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
      });
  
    expect(assets.length).toEqual(0);
  }

  function expectYieldParamsMatchEvent(
    event: UpdateYieldParamsEventV1,
    yieldParams: YieldParamsFromDatabase,
  ) {
    expect(yieldParams.assetYieldIndex).toEqual(event.assetYieldIndex);
    expect(yieldParams.sDAIPrice).toEqual(event.sdaiPrice);
  }

  function expectYieldParamsMatchBlock(
    height: number,
    time: Timestamp,
    yieldParams: YieldParamsFromDatabase,
  ) {
    expect(yieldParams.createdAtHeight).toEqual(height.toString());
    const date = new Date(time.seconds.low * 1000);
    date.setMilliseconds(date.getMilliseconds() + Math.floor(time.nanos / 1e6));
    const isoString = date.toISOString();
    expect(yieldParams.createdAt).toEqual(isoString);
    expect(yieldParams.id).toEqual(YieldParamsTable.uuid(height.toString()));
  }
