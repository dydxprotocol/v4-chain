import {
  AddressTradingReward,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  Timestamp,
  TradingRewardsEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import {
  dbHelpers,
  testMocks,
  WalletTable,
  WalletFromDatabase,
  testConstants,
  TradingRewardTable,
  TradingRewardFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { KafkaMessage } from 'kafkajs';
import { createKafkaMessage } from '@dydxprotocol-indexer/kafka';
import { onMessage } from '../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import {
  defaultDateTime,
  defaultHeight,
  defaultPreviousHeight,
  defaultTime,
  defaultTradingRewardsEvent,
  defaultTxHash,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';
import Big from 'big.js';
import { TradingRewardsHandler } from '../../src/handlers/trading-rewards-handler';
import { bigIntToBytes } from '@dydxprotocol-indexer/v4-proto-parser';

const defaultTransactionIndex: number = 0;

describe('tradingRewardHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  beforeEach(async () => {
    await testMocks.seedData();
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

  describe('getParallelizationIds', () => {
    it('returns the correct parallelization ids', () => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.TRADING_REWARD,
        TradingRewardsEventV1.encode(defaultTradingRewardsEvent).finish(),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        defaultHeight,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: TradingRewardsHandler = new TradingRewardsHandler(
        block,
        0,
        indexerTendermintEvent,
        0,
        defaultTradingRewardsEvent,
      );

      expect(handler.getParallelizationIds()).toEqual([]);
    });
  });

  it('Creates a block reward for each reward in TradingRewardsEvent', async () => {
    const tradingRewardsEvent: TradingRewardsEventV1 = TradingRewardsEventV1.fromPartial({
      tradingRewards: [
        {
          owner: testConstants.defaultWallet.address,
          denomAmount: bigIntToBytes(BigInt(1)),
        },
        {
          owner: testConstants.defaultWallet2.address,
          denomAmount: bigIntToBytes(BigInt(1_000_000_000)),
        },
      ],
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromTradingRewardsEvent({
      tradingRewardsEvent,
      transactionIndex: defaultTransactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await onMessage(kafkaMessage);
    const firstTradingRewardId: string = TradingRewardTable.uuid(
      testConstants.defaultWallet.address,
      defaultHeight.toString(),
    );
    const firstTradingReward:
    TradingRewardFromDatabase | undefined = await TradingRewardTable.findById(
      firstTradingRewardId,
    );
    const secondTradingRewardId: string = TradingRewardTable.uuid(
      testConstants.defaultWallet2.address,
      defaultHeight.toString(),
    );
    const secondTradingReward:
    TradingRewardFromDatabase | undefined = await TradingRewardTable.findById(
      secondTradingRewardId,
    );

    expect(firstTradingReward).toEqual({
      id: firstTradingRewardId,
      address: testConstants.defaultWallet.address,
      blockTime: defaultDateTime.toISO(),
      blockHeight: defaultHeight.toString(),
      amount: Big(1).times(1e-18).toFixed(18),
    });
    expect(secondTradingReward).toEqual({
      id: secondTradingRewardId,
      address: testConstants.defaultWallet2.address,
      blockTime: defaultDateTime.toISO(),
      blockHeight: defaultHeight.toString(),
      amount: Big(1_000_000_000).times(1e-18).toFixed(18),
    });
  });

  it('Creates an wallet and populates totalTradingRewards', async () => {
    const tradingRewardsEvent: TradingRewardsEventV1 = TradingRewardsEventV1.fromPartial({
      tradingRewards: [
        AddressTradingReward.fromPartial({
          owner: testConstants.defaultWallet.address,
          denomAmount: bigIntToBytes(BigInt(1_000_000_000)),
        }),
      ],
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromTradingRewardsEvent({
      tradingRewardsEvent,
      transactionIndex: defaultTransactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await onMessage(kafkaMessage);
    const wallet: WalletFromDatabase | undefined = await WalletTable.findById(
      testConstants.defaultWallet.address,
    );
    expect(wallet).toEqual({
      address: testConstants.defaultWallet.address,
      totalTradingRewards: Big(1_000_000_000).times(1e-18).toFixed(18),
    });
  });

  it('Updates a wallet\'s totalTradingRewards', async () => {
    const tradingRewardsEvent: TradingRewardsEventV1 = TradingRewardsEventV1.fromPartial({
      tradingRewards: [
        AddressTradingReward.fromPartial({
          owner: testConstants.defaultWallet.address,
          denomAmount: bigIntToBytes(BigInt(1_000_000_000)),
        }),
      ],
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromTradingRewardsEvent({
      tradingRewardsEvent,
      transactionIndex: defaultTransactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await WalletTable.create({
      address: testConstants.defaultWallet.address,
      totalTradingRewards: Big(1_000_000_000).times(1e-18).toFixed(18),
    });

    await onMessage(kafkaMessage);
    const wallet: WalletFromDatabase | undefined = await WalletTable.findById(
      testConstants.defaultWallet.address,
    );
    expect(wallet).toEqual({
      address: testConstants.defaultWallet.address,
      totalTradingRewards: Big(2_000_000_000).times(1e-18).toFixed(18),
    });
  });
});

function createKafkaMessageFromTradingRewardsEvent({
  tradingRewardsEvent,
  transactionIndex,
  height,
  time,
  txHash,
}: {
  tradingRewardsEvent: TradingRewardsEventV1,
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}) {
  const events: IndexerTendermintEvent[] = [
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.TRADING_REWARD,
      TradingRewardsEventV1.encode(tradingRewardsEvent).finish(),
      transactionIndex,
      0, // eventIndex
    ),
  ];

  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    height,
    time,
    events,
    [txHash],
  );

  const binaryBlock: Uint8Array = IndexerTendermintBlock.encode(block).finish();
  return createKafkaMessage(Buffer.from(binaryBlock));
}
