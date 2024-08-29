import {
  AssetFromDatabase,
  AssetPositionFromDatabase,
  dbHelpers,
  MarketFromDatabase,
  MarketMessageContents,
  OraclePriceFromDatabase,
  OraclePriceTable,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  PositionSide,
  SubaccountMessageContents,
  SubaccountTable,
  testConstants,
  testMocks,
  TransferFromDatabase,
  TransferType,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import { IndexerSubaccountId } from '@dydxprotocol-indexer/v4-protos';
import { DateTime } from 'luxon';
import {
  addPositionsToContents,
  annotateWithPnl,
  convertPerpetualPosition,
  generateOraclePriceContents,
  generateTransferContents,
  getPnl,
} from '../../src/helpers/kafka-helper';
import { stats } from '@dydxprotocol-indexer/base';
import { updateBlockCache } from '../../src/caches/block-cache';
import { defaultPreviousHeight, defaultWalletAddress } from './constants';

describe('kafka-helper', () => {
  const blockHeight: string = '5';

  describe('addPositionsToContents', () => {
    const defaultPerpetualPosition: PerpetualPositionFromDatabase = {
      id: '',
      subaccountId: testConstants.defaultSubaccountId,
      perpetualId: testConstants.defaultPerpetualMarket.id,
      side: PositionSide.LONG,
      status: PerpetualPositionStatus.OPEN,
      size: '10',
      maxSize: '25',
      entryPrice: '20000',
      sumOpen: '10',
      sumClose: '0',
      createdAt: DateTime.utc().toISO(),
      createdAtHeight: '1',
      openEventId: testConstants.defaultTendermintEventId,
      lastEventId: testConstants.defaultTendermintEventId,
      settledFunding: '200000',
    };

    const defaultPerpetualMarket: PerpetualMarketFromDatabase = {
      ...testConstants.defaultPerpetualMarket,
    };

    const defaultAssetPosition: AssetPositionFromDatabase = {
      id: '',
      ...testConstants.defaultAssetPosition,
    };

    const defaultAsset: AssetFromDatabase = {
      ...testConstants.defaultAsset,
    };

    it('successfully adds no position', () => {
      const subaccountIdProto: IndexerSubaccountId = IndexerSubaccountId.fromPartial({
        owner: 'owner',
        number: 1,
      });

      const contents: SubaccountMessageContents = addPositionsToContents(
        {},
        subaccountIdProto,
        [],
        {},
        [],
        {},
        blockHeight,
      );

      expect(contents.perpetualPositions).toEqual(undefined);
      expect(contents.assetPositions).toEqual(undefined);
      expect(contents.blockHeight).toEqual(blockHeight);
    });

    it('successfully adds one asset position and one perp position', () => {
      const subaccountIdProto: IndexerSubaccountId = IndexerSubaccountId.fromPartial({
        owner: 'owner',
        number: 1,
      });

      const contents: SubaccountMessageContents = addPositionsToContents(
        {},
        subaccountIdProto,
        [defaultPerpetualPosition],
        { [defaultPerpetualMarket.id]: defaultPerpetualMarket },
        [defaultAssetPosition],
        { [defaultAsset.id]: defaultAsset },
        blockHeight,
      );

      expect(contents.perpetualPositions!.length).toEqual(1);
      expect(contents.perpetualPositions![0]).toEqual({
        address: subaccountIdProto.owner,
        subaccountNumber: subaccountIdProto.number,
        positionId: defaultPerpetualPosition.id,
        market: defaultPerpetualMarket.ticker,
        side: defaultPerpetualPosition.side,
        status: defaultPerpetualPosition.status,
        size: defaultPerpetualPosition.size,
        maxSize: defaultPerpetualPosition.maxSize,
        netFunding: defaultPerpetualPosition.settledFunding,
        entryPrice: defaultPerpetualPosition.entryPrice,
        exitPrice: defaultPerpetualPosition.exitPrice,
        sumOpen: defaultPerpetualPosition.sumOpen,
        sumClose: defaultPerpetualPosition.sumClose,
      });

      expect(contents.assetPositions!.length).toEqual(1);
      expect(contents.assetPositions![0]).toEqual({
        address: subaccountIdProto.owner,
        subaccountNumber: subaccountIdProto.number,
        positionId: defaultAssetPosition.id,
        assetId: defaultAsset.id,
        symbol: defaultAsset.symbol,
        side: 'LONG',
        size: defaultAssetPosition.size,
      });
      expect(contents.blockHeight).toEqual(blockHeight);
    });

    it('successfully adds one asset position', () => {
      const subaccountIdProto: IndexerSubaccountId = IndexerSubaccountId.fromPartial({
        owner: 'owner',
        number: 1,
      });

      const contents: SubaccountMessageContents = addPositionsToContents(
        {},
        subaccountIdProto,
        [],
        {},
        [defaultAssetPosition],
        { [defaultAsset.id]: defaultAsset },
        blockHeight,
      );

      expect(contents.perpetualPositions).toBeUndefined();

      expect(contents.assetPositions!.length).toEqual(1);
      expect(contents.assetPositions![0]).toEqual({
        address: subaccountIdProto.owner,
        subaccountNumber: subaccountIdProto.number,
        positionId: defaultAssetPosition.id,
        assetId: defaultAsset.id,
        symbol: defaultAsset.symbol,
        side: 'LONG',
        size: defaultAssetPosition.size,
      });
      expect(contents.blockHeight).toEqual(blockHeight);
    });

    it('successfully adds one perp position', () => {
      const subaccountIdProto: IndexerSubaccountId = IndexerSubaccountId.fromPartial({
        owner: 'owner',
        number: 1,
      });

      const contents: SubaccountMessageContents = addPositionsToContents(
        {},
        subaccountIdProto,
        [defaultPerpetualPosition],
        { [defaultPerpetualMarket.id]: defaultPerpetualMarket },
        [],
        {},
        blockHeight,
      );

      expect(contents.perpetualPositions!.length).toEqual(1);
      expect(contents.perpetualPositions![0]).toEqual({
        address: subaccountIdProto.owner,
        subaccountNumber: subaccountIdProto.number,
        positionId: defaultPerpetualPosition.id,
        market: defaultPerpetualMarket.ticker,
        side: defaultPerpetualPosition.side,
        status: defaultPerpetualPosition.status,
        size: defaultPerpetualPosition.size,
        maxSize: defaultPerpetualPosition.maxSize,
        netFunding: defaultPerpetualPosition.settledFunding,
        entryPrice: defaultPerpetualPosition.entryPrice,
        exitPrice: defaultPerpetualPosition.exitPrice,
        sumOpen: defaultPerpetualPosition.sumOpen,
        sumClose: defaultPerpetualPosition.sumClose,
      });

      expect(contents.assetPositions).toBeUndefined();
      expect(contents.blockHeight).toEqual(blockHeight);
    });

    it('successfully adds multiple positions', () => {
      const subaccountIdProto: IndexerSubaccountId = IndexerSubaccountId.fromPartial({
        owner: 'owner',
        number: 1,
      });

      const perpSize: string = '2';
      const assetSize: string = '5';
      const contents: SubaccountMessageContents = addPositionsToContents(
        {},
        subaccountIdProto,
        [
          defaultPerpetualPosition,
          {
            ...defaultPerpetualPosition,
            size: perpSize,
          },
        ],
        { [defaultPerpetualMarket.id]: defaultPerpetualMarket },
        [
          defaultAssetPosition,
          {
            ...defaultAssetPosition,
            size: assetSize,
          },
        ],
        { [defaultAsset.id]: defaultAsset },
        blockHeight,
      );

      // check perpetual positions
      expect(contents.perpetualPositions!.length).toEqual(2);
      expect(contents.perpetualPositions![0]).toEqual({
        address: subaccountIdProto.owner,
        subaccountNumber: subaccountIdProto.number,
        positionId: defaultPerpetualPosition.id,
        market: defaultPerpetualMarket.ticker,
        side: defaultPerpetualPosition.side,
        status: defaultPerpetualPosition.status,
        size: defaultPerpetualPosition.size,
        maxSize: defaultPerpetualPosition.maxSize,
        netFunding: defaultPerpetualPosition.settledFunding,
        entryPrice: defaultPerpetualPosition.entryPrice,
        exitPrice: defaultPerpetualPosition.exitPrice,
        sumOpen: defaultPerpetualPosition.sumOpen,
        sumClose: defaultPerpetualPosition.sumClose,
      });
      expect(contents.perpetualPositions![1]).toEqual({
        address: subaccountIdProto.owner,
        subaccountNumber: subaccountIdProto.number,
        positionId: defaultPerpetualPosition.id,
        market: defaultPerpetualMarket.ticker,
        side: defaultPerpetualPosition.side,
        status: defaultPerpetualPosition.status,
        size: perpSize,
        maxSize: defaultPerpetualPosition.maxSize,
        netFunding: defaultPerpetualPosition.settledFunding,
        entryPrice: defaultPerpetualPosition.entryPrice,
        exitPrice: defaultPerpetualPosition.exitPrice,
        sumOpen: defaultPerpetualPosition.sumOpen,
        sumClose: defaultPerpetualPosition.sumClose,
      });

      // check asset positions
      expect(contents.assetPositions!.length).toEqual(2);
      expect(contents.assetPositions![0]).toEqual({
        address: subaccountIdProto.owner,
        subaccountNumber: subaccountIdProto.number,
        positionId: defaultAssetPosition.id,
        assetId: defaultAsset.id,
        symbol: defaultAsset.symbol,
        side: 'LONG',
        size: defaultAssetPosition.size,
      });
      expect(contents.assetPositions![1]).toEqual({
        address: subaccountIdProto.owner,
        subaccountNumber: subaccountIdProto.number,
        positionId: defaultAssetPosition.id,
        assetId: defaultAsset.id,
        symbol: defaultAsset.symbol,
        side: 'LONG',
        size: assetSize,
      });
      expect(contents.blockHeight).toEqual(blockHeight);
    });
  });

  describe('addTransferToContents', () => {
    const defaultAsset: AssetFromDatabase = {
      ...testConstants.defaultAsset,
    };
    const senderSubaccountId: IndexerSubaccountId = IndexerSubaccountId.fromPartial({
      owner: 'sender',
      number: 1,
    });

    const recipientSubaccountId: IndexerSubaccountId = IndexerSubaccountId.fromPartial({
      owner: 'recipient',
      number: 1,
    });

    const transfer: TransferFromDatabase = {
      id: '',
      senderSubaccountId: SubaccountTable.uuid(senderSubaccountId.owner, senderSubaccountId.number),
      recipientSubaccountId: SubaccountTable.uuid(
        recipientSubaccountId.owner,
        recipientSubaccountId.number,
      ),
      assetId: defaultAsset.id,
      size: '10',
      eventId: testConstants.defaultTendermintEventId,
      transactionHash: 'hash',
      createdAt: DateTime.utc().toISO(),
      createdAtHeight: '1',
    };

    const deposit: TransferFromDatabase = {
      id: '',
      senderWalletAddress: defaultWalletAddress,
      recipientSubaccountId: SubaccountTable.uuid(
        recipientSubaccountId.owner,
        recipientSubaccountId.number,
      ),
      assetId: defaultAsset.id,
      size: '10',
      eventId: testConstants.defaultTendermintEventId,
      transactionHash: 'hash',
      createdAt: DateTime.utc().toISO(),
      createdAtHeight: '1',
    };

    const withdrawal: TransferFromDatabase = {
      id: '',
      senderSubaccountId: SubaccountTable.uuid(senderSubaccountId.owner, senderSubaccountId.number),
      recipientWalletAddress: defaultWalletAddress,
      assetId: defaultAsset.id,
      size: '10',
      eventId: testConstants.defaultTendermintEventId,
      transactionHash: 'hash',
      createdAt: DateTime.utc().toISO(),
      createdAtHeight: '1',
    };

    it('successfully adds a transfer_out', () => {
      const contents: SubaccountMessageContents = generateTransferContents(
        transfer,
        defaultAsset,
        senderSubaccountId,
        senderSubaccountId,
        recipientSubaccountId,
        transfer.createdAtHeight,
      );

      expect(contents.transfers).toEqual({
        sender: {
          address: senderSubaccountId.owner,
          subaccountNumber: senderSubaccountId.number,
        },
        recipient: {
          address: recipientSubaccountId.owner,
          subaccountNumber: recipientSubaccountId.number,
        },
        symbol: defaultAsset.symbol,
        size: transfer.size,
        type: TransferType.TRANSFER_OUT,
        createdAt: transfer.createdAt,
        createdAtHeight: transfer.createdAtHeight,
        transactionHash: transfer.transactionHash,
      });
      expect(contents.blockHeight).toEqual(transfer.createdAtHeight);
    });

    it('successfully adds a transfer_in', () => {
      const contents: SubaccountMessageContents = generateTransferContents(
        transfer,
        defaultAsset,
        recipientSubaccountId,
        senderSubaccountId,
        recipientSubaccountId,
      );

      expect(contents.transfers).toEqual({
        sender: {
          address: senderSubaccountId.owner,
          subaccountNumber: senderSubaccountId.number,
        },
        recipient: {
          address: recipientSubaccountId.owner,
          subaccountNumber: recipientSubaccountId.number,
        },
        symbol: defaultAsset.symbol,
        size: transfer.size,
        type: TransferType.TRANSFER_IN,
        createdAt: transfer.createdAt,
        createdAtHeight: transfer.createdAtHeight,
        transactionHash: transfer.transactionHash,
      });
    });

    it('successfully adds a deposit', () => {
      const contents: SubaccountMessageContents = generateTransferContents(
        deposit,
        defaultAsset,
        recipientSubaccountId,
        undefined,
        recipientSubaccountId,
      );

      expect(contents.transfers).toEqual({
        sender: {
          address: defaultWalletAddress,
        },
        recipient: {
          address: recipientSubaccountId.owner,
          subaccountNumber: recipientSubaccountId.number,
        },
        symbol: defaultAsset.symbol,
        size: deposit.size,
        type: TransferType.DEPOSIT,
        createdAt: deposit.createdAt,
        createdAtHeight: deposit.createdAtHeight,
        transactionHash: deposit.transactionHash,
      });
    });

    it('successfully adds a withdrawal', () => {
      const contents: SubaccountMessageContents = generateTransferContents(
        withdrawal,
        defaultAsset,
        senderSubaccountId,
        senderSubaccountId,
        undefined,
      );

      expect(contents.transfers).toEqual({
        sender: {
          address: senderSubaccountId.owner,
          subaccountNumber: senderSubaccountId.number,
        },
        recipient: {
          address: defaultWalletAddress,
        },
        symbol: defaultAsset.symbol,
        size: deposit.size,
        type: TransferType.WITHDRAWAL,
        createdAt: withdrawal.createdAt,
        createdAtHeight: withdrawal.createdAtHeight,
        transactionHash: withdrawal.transactionHash,
      });
    });
  });

  describe('marketUpdateToContents', () => {
    const height: string = '3';
    const oraclePrice: OraclePriceFromDatabase = {
      id: OraclePriceTable.uuid(0, height),
      marketId: 0,
      price: '500000.00',
      effectiveAt: DateTime.utc().toISO(),
      effectiveAtHeight: height,
    };

    const market: MarketFromDatabase = {
      id: 0,
      pair: 'BTC-USD',
      exponent: -5,
      minPriceChangePpm: 50,
    };

    it('successfully generates kafka contents from oracle price', () => {
      const contents: MarketMessageContents = generateOraclePriceContents(
        oraclePrice,
        market.pair,
      );

      expect(contents.oraclePrices).toEqual({
        [market.pair]: {
          oraclePrice: oraclePrice.price,
          effectiveAt: oraclePrice.effectiveAt,
          effectiveAtHeight: oraclePrice.effectiveAtHeight,
          marketId: oraclePrice.marketId,
        },
      });
    });
  });

  describe('pnl', () => {
    const updatedObject: UpdatedPerpetualPositionSubaccountKafkaObject = {
      perpetualId: '0',
      maxSize: '25',
      side: PositionSide.LONG,
      entryPrice: '0',
      sumOpen: '0',
      sumClose: '0',
      id: '65c77c62-043b-5dd0-9ba9-0f9cc130eca8',
      closedAt: null,
      closedAtHeight: null,
      closeEventId: null,
      settledFunding: '-199998',
      status: PerpetualPositionStatus.OPEN,
      size: '0.0001',
      lastEventId: Buffer.from('0'),
    };

    beforeAll(async () => {
      await dbHelpers.migrate();
      jest.spyOn(stats, 'increment');
      jest.spyOn(stats, 'timing');
      jest.spyOn(stats, 'gauge');
    });

    beforeEach(async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
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

    it('getPnl', () => {
      const {
        realizedPnl,
        unrealizedPnl,
      }: {
        realizedPnl: string | undefined,
        unrealizedPnl: string | undefined,
      } = getPnl(
        updatedObject,
        perpetualMarketRefresher.getPerpetualMarketsMap()[updatedObject.perpetualId],
        testConstants.defaultMarket,
      );
      expect(realizedPnl).toEqual('-199998');  // 0*0-199998
      expect(unrealizedPnl).toEqual('1.5');  // 0.0001*(15000-0)
    });

    it('getPnl with non-negative sumClose', () => {
      const updatedObject2: UpdatedPerpetualPositionSubaccountKafkaObject = {
        ...updatedObject,
        sumClose: '1',
        exitPrice: '5',
      };
      const {
        realizedPnl,
        unrealizedPnl,
      }: {
        realizedPnl: string | undefined,
        unrealizedPnl: string | undefined,
      } = getPnl(
        updatedObject2,
        perpetualMarketRefresher.getPerpetualMarketsMap()[updatedObject2.perpetualId],
        testConstants.defaultMarket,
      );
      expect(realizedPnl).toEqual('-199993');  // 1*5-199998
      expect(unrealizedPnl).toEqual('1.5');  // 0.0001*(15000-0)
    });

    it('convertPerpetualPosition', async () => {
      const createdPerpetualPosition: PerpetualPositionFromDatabase = await
      PerpetualPositionTable.create(
        testConstants.defaultPerpetualPosition,
      );
      const update: UpdatedPerpetualPositionSubaccountKafkaObject = convertPerpetualPosition(
        createdPerpetualPosition,
      );
      expect(update).toEqual({
        perpetualId: '0',
        maxSize: '25',
        side: PositionSide.LONG,
        entryPrice: '20000',
        exitPrice: null,
        sumOpen: '10',
        sumClose: '0',
        id: '65c77c62-043b-5dd0-9ba9-0f9cc130eca8',
        closedAt: null,
        closedAtHeight: null,
        closeEventId: null,
        lastEventId: testConstants.defaultPerpetualPosition.lastEventId,
        settledFunding: '200000',
        status: PerpetualPositionStatus.OPEN,
        size: '10',
      });
    });

    it('annotateWithPnl', () => {
      const updatedObject2: UpdatedPerpetualPositionSubaccountKafkaObject = {
        ...updatedObject,
        sumClose: '1',
        exitPrice: '5',
      };
      const updatedObjectWithPnl: UpdatedPerpetualPositionSubaccountKafkaObject = annotateWithPnl(
        updatedObject2,
        perpetualMarketRefresher.getPerpetualMarketsMap(),
        testConstants.defaultMarket,
      );
      expect(
        updatedObjectWithPnl,
      ).toEqual({
        ...updatedObject2,
        realizedPnl: '-199993',
        unrealizedPnl: '1.5',
      });
    });
  });
});
