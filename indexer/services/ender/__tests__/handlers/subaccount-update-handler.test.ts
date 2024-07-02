import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  SubaccountUpdateEventV1,
  Timestamp,
} from '@dydxprotocol-indexer/v4-protos';
import {
  AssetColumns,
  AssetFromDatabase,
  AssetPositionCreateObject,
  AssetPositionFromDatabase,
  AssetPositionTable,
  assetRefresher,
  AssetTable,
  dbHelpers,
  PerpetualMarketColumns,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualMarketsMap,
  PerpetualPositionCreateObject,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  PositionSide,
  protocolTranslations,
  SubaccountFromDatabase,
  SubaccountMessageContents,
  SubaccountTable,
  TendermintEventTable,
  testConstants,
  testMocks,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import { bigIntToBytes, bytesToBase64 } from '@dydxprotocol-indexer/v4-proto-parser';
import { KafkaMessage } from 'kafkajs';
import _ from 'lodash';
import { DateTime } from 'luxon';
import { SUBACCOUNT_ORDER_FILL_EVENT_TYPE } from '../../src/constants';
import { createKafkaMessage, producer } from '@dydxprotocol-indexer/kafka';
import { addPositionsToContents, annotateWithPnl, convertPerpetualPosition } from '../../src/helpers/kafka-helper';
import { onMessage } from '../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectSubaccountKafkaMessage,
} from '../helpers/indexer-proto-helpers';
import { SubaccountUpdateHandler } from '../../src/handlers/subaccount-update-handler';
import {
  defaultDateTime,
  defaultEmptySubaccountUpdate,
  defaultEmptySubaccountUpdateEvent,
  defaultHeight,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';

describe('subaccountUpdateHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
    await assetRefresher.updateAssets();
    updateBlockCache(defaultPreviousHeight);
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  const defaultPerpetualPosition: PerpetualPositionCreateObject = {
    subaccountId: testConstants.defaultSubaccountId,
    perpetualId: testConstants.defaultPerpetualMarket.id,
    side: PositionSide.LONG,
    status: PerpetualPositionStatus.OPEN,
    size: '10',
    maxSize: '25',
    createdAt: DateTime.utc().toISO(),
    createdAtHeight: '1',
    openEventId: testConstants.defaultTendermintEventId,
    lastEventId: testConstants.defaultTendermintEventId,
    settledFunding: '-200000',
  };

  const defaultAssetPosition: AssetPositionCreateObject = {
    subaccountId: testConstants.defaultSubaccountId,
    assetId: '0',
    size: '10000',
    isLong: true,
  };

  describe('getParallelizationIds', () => {
    it('returns the correct parallelization ids', () => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
        SubaccountUpdateEventV1.encode(defaultEmptySubaccountUpdateEvent).finish(),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: SubaccountUpdateHandler = new SubaccountUpdateHandler(
        block,
        0,
        indexerTendermintEvent,
        0,
        defaultEmptySubaccountUpdate,
      );

      expect(handler.getParallelizationIds()).toEqual([
        `${handler.eventType}_${SubaccountTable.subaccountIdToUuid(defaultEmptySubaccountUpdateEvent.subaccountId!)}`,
        `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${SubaccountTable.subaccountIdToUuid(defaultEmptySubaccountUpdateEvent.subaccountId!)}`,
      ]);
    });
  });

  it('successfully creates subaccount', async () => {
    const transactionIndex: number = 0;
    const address: string = 'cosmosnewaddress';
    const subaccountId: string = SubaccountTable.uuid(
      address,
      testConstants.defaultSubaccount.subaccountNumber,
    );
    const subaccountUpdateEvent: SubaccountUpdateEventV1 = SubaccountUpdateEventV1.fromPartial(
      {
        subaccountId: {
          owner: address,
          number: testConstants.defaultSubaccount.subaccountNumber,
        },
        updatedPerpetualPositions: [],
        updatedAssetPositions: [],
      },
    );
    const kafkaMessage: KafkaMessage = createKafkaMessageFromSubaccountUpdateEvent({
      subaccountUpdateEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);
    const subaccount: SubaccountFromDatabase | undefined = await
    SubaccountTable.findById(subaccountId);

    expect(subaccount).not.toEqual(undefined);
    await expectUpdatedPositionsSubaccountKafkaMessage(
      producerSendMock,
      subaccountUpdateEvent,
      [],
      [],
    );
  });

  it.each([
    ['positive', 200000, '-0.2'],
    ['negative', -200000, '0.2'],
  ])('successfully upserts perpetual position with %s funding payment', async (
    _name: string,
    fundingPayment: number,
    settledFunding: string,
  ) => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const sizeInQuantums: number = 1_000_000;
    const fundingIndex: number = 200;
    const subaccountUpdateEvent: SubaccountUpdateEventV1 = SubaccountUpdateEventV1.fromPartial({
      // eslint-disable-next-line  @typescript-eslint/no-explicit-any
      ...defaultEmptySubaccountUpdateEvent as any,
      updatedPerpetualPositions: [{
        perpetualId: parseInt(testConstants.defaultPerpetualMarket.id, 10),
        quantums: bytesToBase64(bigIntToBytes(BigInt(sizeInQuantums))),
        fundingIndex: bytesToBase64(bigIntToBytes(BigInt(fundingIndex))),
        fundingPayment: bytesToBase64(bigIntToBytes(BigInt(fundingPayment))),
      }],
    });

    const kafkaMessage: KafkaMessage = createKafkaMessageFromSubaccountUpdateEvent({
      subaccountUpdateEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const tendermintEventId: Buffer = TendermintEventTable.createEventId(
      defaultHeight.toString(),
      transactionIndex,
      eventIndex,
    );
    const perpetualPosition: PerpetualPositionFromDatabase | undefined = await
    PerpetualPositionTable.findById(
      PerpetualPositionTable.uuid(
        testConstants.defaultSubaccountId,
        tendermintEventId,
      ),
    );

    expect(perpetualPosition).not.toEqual(undefined);
    const size: string = protocolTranslations.quantumsToHumanFixedString(
      sizeInQuantums.toString(),
      testConstants.defaultPerpetualMarket!.atomicResolution,
    );
    expect(perpetualPosition!).toEqual(expect.objectContaining({
      subaccountId: testConstants.defaultSubaccountId,
      perpetualId: testConstants.defaultPerpetualMarket.id,
      side: PositionSide.LONG,
      status: PerpetualPositionStatus.OPEN,
      size,
      maxSize: size,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight.toString(),
      openEventId: tendermintEventId,
      lastEventId: tendermintEventId,
      settledFunding,
    }));
    const updatedPerpetualPositionSubaccountKafkaObject:
    UpdatedPerpetualPositionSubaccountKafkaObject = annotateWithPnl(
      convertPerpetualPosition(perpetualPosition!),
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      testConstants.defaultMarket,
    );
    await expectUpdatedPositionsSubaccountKafkaMessage(
      producerSendMock,
      subaccountUpdateEvent,
      [updatedPerpetualPositionSubaccountKafkaObject],
      [],
    );
  });

  it.each([
    ['positive', 2_000_000, '-200002'],
    ['negative', -2_000_000, '-199998'],
  ])('successfully updates existing perpetual position with %s funding payment', async (
    _name: string,
    fundingPayment: number,
    settledFunding: string,
  ) => {
    const transactionIndex: number = 0;
    const sizeInQuantums: number = 1_000_000;
    const fundingIndex: number = 200;
    const subaccountUpdateEvent: SubaccountUpdateEventV1 = SubaccountUpdateEventV1.fromPartial({
      // eslint-disable-next-line  @typescript-eslint/no-explicit-any
      ...defaultEmptySubaccountUpdateEvent as any,
      updatedPerpetualPositions: [{
        perpetualId: parseInt(testConstants.defaultPerpetualMarket.id, 10),
        quantums: bytesToBase64(bigIntToBytes(BigInt(sizeInQuantums))),
        fundingIndex: bytesToBase64(bigIntToBytes(BigInt(fundingIndex))),
        fundingPayment: bytesToBase64(bigIntToBytes(BigInt(fundingPayment))),
      }],
    });

    await PerpetualPositionTable.create(defaultPerpetualPosition);

    const kafkaMessage: KafkaMessage = createKafkaMessageFromSubaccountUpdateEvent({
      subaccountUpdateEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const perpetualPosition: PerpetualPositionFromDatabase | undefined = await
    PerpetualPositionTable.findOpenPositionForSubaccountPerpetual(
      testConstants.defaultSubaccountId,
      testConstants.defaultPerpetualMarket.id,
    );

    expect(perpetualPosition).not.toEqual(undefined);
    expect(perpetualPosition!).toEqual(expect.objectContaining({
      subaccountId: testConstants.defaultSubaccountId,
      perpetualId: testConstants.defaultPerpetualMarket.id,
      side: PositionSide.LONG,
      status: PerpetualPositionStatus.OPEN,
      settledFunding,
    }));
    const updatedPerpetualPositionSubaccountKafkaObject:
    UpdatedPerpetualPositionSubaccountKafkaObject = annotateWithPnl(
      convertPerpetualPosition(perpetualPosition!),
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      testConstants.defaultMarket,
    );
    await expectUpdatedPositionsSubaccountKafkaMessage(
      producerSendMock,
      subaccountUpdateEvent,
      [updatedPerpetualPositionSubaccountKafkaObject],
      [],
    );
  });

  it('closes and creates new position when when side is opposing', async () => {
    const transactionIndex: number = 0;
    const sizeInQuantums: number = 1_000_000;
    const fundingIndex: number = 200;
    const fundingPayment: number = 1_000_000_000;
    const subaccountUpdateEvent: SubaccountUpdateEventV1 = SubaccountUpdateEventV1.fromPartial({
      // eslint-disable-next-line  @typescript-eslint/no-explicit-any
      ...defaultEmptySubaccountUpdateEvent as any,
      updatedPerpetualPositions: [{
        perpetualId: testConstants.defaultPerpetualMarket.id,
        quantums: bytesToBase64(bigIntToBytes(BigInt(sizeInQuantums))),
        fundingIndex: bytesToBase64(bigIntToBytes(BigInt(fundingIndex))),
        fundingPayment: bytesToBase64(bigIntToBytes(BigInt(fundingPayment))),
      }],
    });

    const kafkaMessage: KafkaMessage = createKafkaMessageFromSubaccountUpdateEvent({
      subaccountUpdateEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    // Initialize existing perpetual position
    const perpetualSize: string = protocolTranslations.quantumsToHumanFixedString(
      sizeInQuantums.toString(),
      testConstants.defaultPerpetualMarket!.atomicResolution,
    );
    const createdPerpetualPosition: PerpetualPositionFromDatabase = await
    PerpetualPositionTable.create({
      ...defaultPerpetualPosition,
      side: PositionSide.SHORT,
      size: perpetualSize,
      sumOpen: perpetualSize,
    });

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const closedPosition: PerpetualPositionFromDatabase | undefined = await
    PerpetualPositionTable.findById(createdPerpetualPosition.id);

    expect(closedPosition).not.toEqual(undefined);
    expect(closedPosition).toEqual(expect.objectContaining({
      size: '0',
      status: PerpetualPositionStatus.CLOSED,
      settledFunding: '-201000',  // existing settledFunding = -200000, new position funding payment = -1000.
    }));

    const newPosition: PerpetualPositionFromDatabase | undefined = await
    PerpetualPositionTable.findOpenPositionForSubaccountPerpetual(
      testConstants.defaultSubaccountId,
      testConstants.defaultPerpetualMarket.id,
    );
    expect(newPosition).not.toBeUndefined();
    expect(newPosition).toEqual(expect.objectContaining({
      size: perpetualSize,
      maxSize: perpetualSize,
      settledFunding: '0',  // settledFunding of new opened position is 0.
    }));
    const closedPositionSubaccountKafkaObject:
    UpdatedPerpetualPositionSubaccountKafkaObject = annotateWithPnl(
      convertPerpetualPosition(closedPosition!),
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      testConstants.defaultMarket,
    );
    const newPositionSubaccountKafkaObject:
    UpdatedPerpetualPositionSubaccountKafkaObject = annotateWithPnl(
      convertPerpetualPosition(newPosition!),
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      testConstants.defaultMarket,
    );

    await expectUpdatedPositionsSubaccountKafkaMessage(
      producerSendMock,
      subaccountUpdateEvent,
      [closedPositionSubaccountKafkaObject, newPositionSubaccountKafkaObject!],
      [],
    );
  });

  it('updates existing asset position', async () => {
    const transactionIndex: number = 0;
    const sizeInQuantums: number = 1_000_000;
    const subaccountUpdateEvent: SubaccountUpdateEventV1 = SubaccountUpdateEventV1.fromPartial({
      // eslint-disable-next-line  @typescript-eslint/no-explicit-any
      ...defaultEmptySubaccountUpdateEvent as any,
      updatedAssetPositions: [{
        assetId: testConstants.defaultAsset.id,
        quantums: bytesToBase64(bigIntToBytes(BigInt(2 * sizeInQuantums))),
      }],
    });

    const kafkaMessage: KafkaMessage = createKafkaMessageFromSubaccountUpdateEvent({
      subaccountUpdateEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    // Initialize existing asset position
    const asset: AssetFromDatabase | undefined = await
    AssetTable.findById(testConstants.defaultAsset.id);
    const createdAssetPosition: AssetPositionFromDatabase = await
    AssetPositionTable.upsert(defaultAssetPosition);

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const newPosition: AssetPositionFromDatabase | undefined = await
    AssetPositionTable.findById(createdAssetPosition.id);
    expect(newPosition).not.toBeUndefined();
    const newAssetSize: string = protocolTranslations.quantumsToHumanFixedString(
      (2 * sizeInQuantums).toString(),
      asset!.atomicResolution,
    );
    expect(newPosition).toEqual(expect.objectContaining({
      size: newAssetSize,
    }));
    await expectUpdatedPositionsSubaccountKafkaMessage(
      producerSendMock,
      subaccountUpdateEvent,
      [],
      [newPosition!],
    );
  });

  it.each([
    [1_000_000],
    [-2_000_000],
  ])('creates new asset position, size = %d', async (
    sizeInQuantums: number,
  ) => {
    const transactionIndex: number = 0;
    const subaccountUpdateEvent: SubaccountUpdateEventV1 = SubaccountUpdateEventV1.fromPartial({
      // eslint-disable-next-line  @typescript-eslint/no-explicit-any
      ...defaultEmptySubaccountUpdateEvent as any,
      updatedAssetPositions: [{
        assetId: testConstants.defaultAsset.id,
        quantums: bytesToBase64(bigIntToBytes(BigInt(sizeInQuantums))),
      }],
    });

    const kafkaMessage: KafkaMessage = createKafkaMessageFromSubaccountUpdateEvent({
      subaccountUpdateEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    // Confirm there is no existing asset position
    const asset: AssetFromDatabase | undefined = await
    AssetTable.findById(testConstants.defaultAsset.id);
    const absAssetSize: string = protocolTranslations.serializedQuantumsToAbsHumanFixedString(
      bigIntToBytes(BigInt(sizeInQuantums)),
      asset!.atomicResolution,
    );

    const existingAssetPositions: AssetPositionFromDatabase[] = await AssetPositionTable.findAll(
      {
        subaccountId: [testConstants.defaultSubaccountId],
      },
      [],
      { readReplica: true },
    );

    expect(existingAssetPositions.length).toEqual(0);

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const newAssetPositions: AssetPositionFromDatabase[] = await AssetPositionTable.findAll(
      {
        subaccountId: [testConstants.defaultSubaccountId],
      },
      [],
      { readReplica: true },
    );

    expect(newAssetPositions.length).toEqual(1);
    const newPosition: AssetPositionFromDatabase | undefined = newAssetPositions[0];
    expect(newPosition).not.toBeUndefined();
    expect(newPosition).toEqual(expect.objectContaining({
      size: absAssetSize,
      isLong: sizeInQuantums > 0,
    }));
    await expectUpdatedPositionsSubaccountKafkaMessage(
      producerSendMock,
      subaccountUpdateEvent,
      [],
      [newPosition!],
    );
  });

  it('closes existing position when size is 0', async () => {
    const transactionIndex: number = 0;
    const sizeInQuantums: number = 1_000_000;
    const fundingIndex: number = 200;
    const subaccountUpdateEvent: SubaccountUpdateEventV1 = SubaccountUpdateEventV1.fromPartial({
      // eslint-disable-next-line  @typescript-eslint/no-explicit-any
      ...defaultEmptySubaccountUpdateEvent as any,
      updatedPerpetualPositions: [{
        perpetualId: testConstants.defaultPerpetualMarket.id,
        quantums: bytesToBase64(bigIntToBytes(BigInt('0'))),
        fundingIndex: bytesToBase64(bigIntToBytes(BigInt(fundingIndex))),
      }],
    });

    const kafkaMessage: KafkaMessage = createKafkaMessageFromSubaccountUpdateEvent({
      subaccountUpdateEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    // Initialize existing position
    const size: string = protocolTranslations.quantumsToHumanFixedString(
      sizeInQuantums.toString(),
      testConstants.defaultPerpetualMarket!.atomicResolution,
    );
    const createdPosition: PerpetualPositionFromDatabase = await PerpetualPositionTable.create({
      ...defaultPerpetualPosition,
      side: PositionSide.SHORT,
      size,
      sumOpen: size,
    });

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const closedPosition: PerpetualPositionFromDatabase | undefined = await
    PerpetualPositionTable.findById(createdPosition.id);

    expect(closedPosition).not.toEqual(undefined);
    expect(closedPosition).toEqual(expect.objectContaining({
      size: '0',
      status: PerpetualPositionStatus.CLOSED,
    }));
    const closedPerpetualPositionSubaccountKafkaObject:
    UpdatedPerpetualPositionSubaccountKafkaObject = annotateWithPnl(
      convertPerpetualPosition(closedPosition!),
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      testConstants.defaultMarket,
    );
    await expectUpdatedPositionsSubaccountKafkaMessage(
      producerSendMock,
      subaccountUpdateEvent,
      [closedPerpetualPositionSubaccountKafkaObject],
      [],
    );

    expect(
      await PerpetualPositionTable.findOpenPositionForSubaccountPerpetual(
        testConstants.defaultSubaccountId,
        testConstants.defaultPerpetualMarket.id,
      ),
    ).toBeUndefined();
  });

  it('successfully upserts perpetual and asset position with fixed-point notation sizes', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const sizeInQuantums: number = 10;
    const fundingInQuantums: number = 200000;
    // Test that the sub-zero size is represented as a fixed-point notation rather than exponential
    // notation,
    const perpetualSize: string = '0.000000001'; // 10 * (10^-10) (atomic resolution of market = -10)
    const assetSize: string = '0.0000001'; // 10 * (10^-8) (atomic resolution of market = -8)
    const fundingPayment: string = '-0.2'; // 200000 * (10^-6) * -1 (atomic resolution of market = -6)
    const fundingIndex: number = 200;
    const subaccountUpdateEvent: SubaccountUpdateEventV1 = SubaccountUpdateEventV1.fromPartial({
      // eslint-disable-next-line  @typescript-eslint/no-explicit-any
      ...defaultEmptySubaccountUpdateEvent as any,
      updatedPerpetualPositions: [{
        perpetualId: testConstants.defaultPerpetualMarket.id,
        quantums: bytesToBase64(bigIntToBytes(BigInt(sizeInQuantums))),
        fundingIndex: bytesToBase64(bigIntToBytes(BigInt(fundingIndex))),
        fundingPayment: bytesToBase64(bigIntToBytes(BigInt(fundingInQuantums))),
      }],
      updatedAssetPositions: [{
        assetId: testConstants.defaultAsset3.id,
        quantums: bytesToBase64(bigIntToBytes(BigInt(sizeInQuantums))),
      }],
    });

    const kafkaMessage: KafkaMessage = createKafkaMessageFromSubaccountUpdateEvent({
      subaccountUpdateEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const tendermintEventId: Buffer = TendermintEventTable.createEventId(
      defaultHeight.toString(),
      transactionIndex,
      eventIndex,
    );
    const perpetualPosition: PerpetualPositionFromDatabase | undefined = await
    PerpetualPositionTable.findById(
      PerpetualPositionTable.uuid(
        testConstants.defaultSubaccountId,
        tendermintEventId,
      ),
    );

    expect(perpetualPosition).not.toEqual(undefined);
    expect(perpetualPosition!).toEqual(expect.objectContaining({
      subaccountId: testConstants.defaultSubaccountId,
      perpetualId: testConstants.defaultPerpetualMarket.id,
      side: PositionSide.LONG,
      status: PerpetualPositionStatus.OPEN,
      size: perpetualSize,
      maxSize: perpetualSize,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight.toString(),
      openEventId: tendermintEventId,
      lastEventId: tendermintEventId,
      settledFunding: fundingPayment,
    }));
    const perpetualPositionSubaccountKafkaObject:
    UpdatedPerpetualPositionSubaccountKafkaObject = annotateWithPnl(
      convertPerpetualPosition(perpetualPosition!),
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      testConstants.defaultMarket,
    );

    const assetPosition: AssetPositionFromDatabase | undefined = await
    AssetPositionTable.findById(
      AssetPositionTable.uuid(
        testConstants.defaultSubaccountId,
        testConstants.defaultAsset3.id,
      ),
    );

    expect(assetPosition).not.toBeUndefined();
    expect(assetPosition).toEqual(expect.objectContaining({
      size: assetSize,
    }));

    await expectUpdatedPositionsSubaccountKafkaMessage(
      producerSendMock,
      subaccountUpdateEvent,
      [perpetualPositionSubaccountKafkaObject!],
      [assetPosition!],
    );
  });
});

function createKafkaMessageFromSubaccountUpdateEvent({
  subaccountUpdateEvent,
  transactionIndex,
  height,
  time,
  txHash,
}: {
  subaccountUpdateEvent: SubaccountUpdateEventV1 | undefined,
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}) {
  const events: IndexerTendermintEvent[] = [];
  let eventIndex: number = 0;
  if (subaccountUpdateEvent !== undefined) {
    events.push(
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
        SubaccountUpdateEventV1.encode(subaccountUpdateEvent).finish(),
        transactionIndex,
        eventIndex,
      ),
    );
    eventIndex += 1;
  }

  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    height,
    time,
    events,
    [txHash],
  );

  const binaryBlock: Uint8Array = IndexerTendermintBlock.encode(block).finish();
  return createKafkaMessage(Buffer.from(binaryBlock));
}

async function expectUpdatedPositionsSubaccountKafkaMessage(
  producerSendMock: jest.SpyInstance,
  event: SubaccountUpdateEventV1,
  perpetualPositions: UpdatedPerpetualPositionSubaccountKafkaObject[],
  assetPositions: AssetPositionFromDatabase[],
  blockHeight: string = '3',
  transactionIndex: number = 0,
  eventIndex: number = 0,
) {
  const perpetualMarketsMap: PerpetualMarketsMap = perpetualMarketRefresher
    .getPerpetualMarketsMap();
  const perpMarkets: PerpetualMarketFromDatabase[] = Object.values(perpetualMarketsMap);

  const assets: AssetFromDatabase[] = await AssetTable.findAll(
    { id: _.map(assetPositions, 'assetId') },
    [],
  );

  const contents: SubaccountMessageContents = addPositionsToContents(
    {} as SubaccountMessageContents,
    event.subaccountId!,
    perpetualPositions,
    _.keyBy(perpMarkets, PerpetualMarketColumns.id),
    assetPositions,
    _.keyBy(assets, AssetColumns.id),
    blockHeight,
  );

  expectSubaccountKafkaMessage({
    producerSendMock,
    blockHeight,
    transactionIndex,
    eventIndex,
    contents: JSON.stringify(contents),
    subaccountIdProto: event.subaccountId!,
  });
}
