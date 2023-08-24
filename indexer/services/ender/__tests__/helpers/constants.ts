import { SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION } from '@dydxprotocol-indexer/kafka';
import { testConstants, TradeContent } from '@dydxprotocol-indexer/postgres';
import {
  bigIntToBytes,
  ORDER_FLAG_LONG_TERM,
  ORDER_FLAG_SHORT_TERM,
  ORDER_FLAG_CONDITIONAL,
} from '@dydxprotocol-indexer/v4-proto-parser';
import {
  FundingEventV1_Type,
  LiquidationOrderV1,
  MarketBaseEventV1,
  MarketEventV1,
  IndexerOrder,
  IndexerOrder_Side,
  IndexerOrder_TimeInForce,
  OrderFillEventV1,
  IndexerOrderId,
  StatefulOrderEventV1,
  IndexerSubaccountId,
  SubaccountMessage,
  SubaccountUpdateEventV1,
  Timestamp,
  TransferEventV1,
  IndexerOrder_ConditionType,
  OrderRemovalReason,
  AssetCreateEventV1,
  PerpetualMarketCreateEventV1,
  ClobPairStatus, LiquidityTierUpsertEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';
import { DateTime } from 'luxon';

import { MILLIS_IN_NANOS, SECONDS_IN_MILLIS } from '../../src/constants';
import { ConsolidatedKafkaEvent, FundingEventMessage, SingleTradeMessage } from '../../src/lib/types';
import { contentToSingleTradeMessage, createConsolidatedKafkaEventFromTrade } from './kafka-publisher-helpers';

export const defaultMarketPriceUpdate: MarketEventV1 = {
  marketId: 0,
  priceUpdate: {
    priceWithExponent: Long.fromValue(100000000),
  },
};

export const defaultMarketPriceUpdate2: MarketEventV1 = {
  marketId: 10,
  priceUpdate: {
    priceWithExponent: Long.fromValue(100000000),
  },
};

export const defaultFundingUpdateSampleEvent: FundingEventMessage = {
  type: FundingEventV1_Type.TYPE_PREMIUM_SAMPLE,
  updates: [
    {
      perpetualId: 0,
      fundingValuePpm: 10,
      fundingIndex: bigIntToBytes(BigInt(0)),
    },
  ],
};

export const defaultFundingRateEvent: FundingEventMessage = {
  type: FundingEventV1_Type.TYPE_FUNDING_RATE_AND_INDEX,
  updates: [
    {
      perpetualId: 0,
      fundingValuePpm: 10,
      fundingIndex: bigIntToBytes(BigInt(10)),
    },
  ],
};

export const defaultAssetCreateEvent: AssetCreateEventV1 = {
  id: 0,
  symbol: 'BTC',
  hasMarket: true,
  marketId: 0,
  atomicResolution: 6,
};

export const defaultMarketBase: MarketBaseEventV1 = {
  pair: 'BTC-USD',
  minPriceChangePpm: 500,
};

export const defaultMarketCreate: MarketEventV1 = {
  marketId: 0,
  marketCreate: {
    base: defaultMarketBase,
    exponent: -5,
  },
};

export const defaultMarketModify: MarketEventV1 = {
  marketId: 0,
  marketModify: {
    base: defaultMarketBase,
  },
};

export const defaultPerpetualMarketCreateEvent: PerpetualMarketCreateEventV1 = {
  id: 0,
  clobPairId: 1,
  ticker: 'BTC-USD',
  marketId: 0,
  status: ClobPairStatus.CLOB_PAIR_STATUS_ACTIVE,
  quantumConversionExponent: -8,
  atomicResolution: -10,
  subticksPerTick: 100,
  minOrderBaseQuantums: Long.fromValue(10),
  stepBaseQuantums: Long.fromValue(10),
  liquidityTier: 0,
};

export const defaultLiquidityTierUpsertEvent: LiquidityTierUpsertEventV1 = {
  id: 0,
  name: 'Large-Cap',
  initialMarginPpm: 50000,  // 5%
  maintenanceFractionPpm: 600000,  // 60% of IM
  basePositionNotional: Long.fromValue(1_000_000_000_000),  // 1_000_000 USDC
};

export const defaultPreviousHeight: string = '2';
export const defaultHeight: number = 3;
export const defaultDateTime: DateTime = DateTime.utc(2022, 6, 1, 12, 1, 1, 2);
export const defaultTime: Timestamp = {
  seconds: Long.fromValue(Math.floor(defaultDateTime.toSeconds())),
  nanos: (defaultDateTime.toMillis() % SECONDS_IN_MILLIS) * MILLIS_IN_NANOS,
};
export const defaultTxHash: string = '0x32343534306431622d306461302d343831322d613730372d3965613162336162';

export const defaultSubaccountId: IndexerSubaccountId = {
  owner: testConstants.defaultAddress,
  number: 0,
};
export const defaultSubaccountId2: IndexerSubaccountId = {
  owner: testConstants.defaultAddress,
  number: 1,
};

export const defaultOrderId: IndexerOrderId = {
  subaccountId: defaultSubaccountId,
  clientId: 0,
  clobPairId: 1,
  orderFlags: ORDER_FLAG_SHORT_TERM,
};
export const defaultOrderId2: IndexerOrderId = {
  subaccountId: defaultSubaccountId2,
  clientId: 0,
  clobPairId: 1,
  orderFlags: ORDER_FLAG_LONG_TERM,
};

export const defaultSubticks: number = 1_000_000_000;
export const defaultMakerOrder: IndexerOrder = {
  orderId: defaultOrderId,
  side: IndexerOrder_Side.SIDE_BUY,
  // Set to unsigned true because when encoding and decoding, telescope converts the Long
  // to unsigned.
  quantums: Long.fromValue(10_000_000_000, true),
  subticks: Long.fromValue(1_000_000_000, true),
  goodTilBlock: 5,
  timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL,
  reduceOnly: false,
  clientMetadata: 0,
  conditionType: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,
  conditionalOrderTriggerSubticks: Long.fromValue(0, true),
};
export const defaultTakerOrder: IndexerOrder = {
  orderId: defaultOrderId2,
  side: IndexerOrder_Side.SIDE_SELL,
  quantums: Long.fromValue(10_000_000_000),
  subticks: Long.fromValue(1_000_000_000),
  goodTilBlock: 5,
  timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
  reduceOnly: true,
  clientMetadata: 0,
  conditionType: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,
  conditionalOrderTriggerSubticks: Long.fromValue(0, true),
};
export const defaultLiquidationOrder: LiquidationOrderV1 = {
  liquidated: defaultSubaccountId,
  clobPairId: parseInt(testConstants.defaultPerpetualMarket.clobPairId, 10),
  perpetualId: parseInt(testConstants.defaultPerpetualMarket.id, 10),
  totalSize: Long.fromValue(10_000),
  isBuy: true,
  subticks: Long.fromValue(1_000_000_000),
};
export const defaultOrderEvent: OrderFillEventV1 = {
  makerOrder: defaultMakerOrder,
  order: defaultTakerOrder,
  makerFee: Long.fromValue(0),
  takerFee: Long.fromValue(0),
  fillAmount: Long.fromValue(10_000),
  totalFilledMaker: Long.fromValue(0),
  totalFilledTaker: Long.fromValue(0),
};
export const defaultLiquidationEvent: OrderFillEventV1 = {
  makerOrder: defaultMakerOrder,
  liquidationOrder: defaultLiquidationOrder,
  makerFee: Long.fromValue(0),
  takerFee: Long.fromValue(0),
  fillAmount: Long.fromValue(10_000),
  totalFilledMaker: Long.fromValue(0),
  totalFilledTaker: Long.fromValue(0),
};

export const defaultEmptySubaccountUpdateEvent: SubaccountUpdateEventV1 = {
  subaccountId: defaultSubaccountId,
  updatedPerpetualPositions: [],
  updatedAssetPositions: [],
};

export const defaultWalletAddress: string = 'defaultWalletAddress';
export const defaultSenderSubaccountId: IndexerSubaccountId = {
  owner: 'sender',
  number: 0,
};
export const defaultRecipientSubaccountId: IndexerSubaccountId = {
  owner: 'recipient',
  number: 0,
};
export const defaultTransferEvent: TransferEventV1 = {
  assetId: 0,
  amount: Long.fromValue(100),
  sender: {
    subaccountId: defaultSenderSubaccountId,
  },
  recipient: {
    subaccountId: defaultRecipientSubaccountId,
  },
};
export const defaultDepositEvent: TransferEventV1 = {
  assetId: 0,
  amount: Long.fromValue(100),
  sender: {
    address: defaultWalletAddress,
  },
  recipient: {
    subaccountId: defaultRecipientSubaccountId,
  },
};
export const defaultWithdrawalEvent: TransferEventV1 = {
  assetId: 0,
  amount: Long.fromValue(100),
  sender: {
    subaccountId: defaultSenderSubaccountId,
  },
  recipient: {
    address: defaultWalletAddress,
  },
};

export const defaultSubaccountMessage: SubaccountMessage = {
  blockHeight: defaultHeight.toString(),
  transactionIndex: 0,
  eventIndex: 0,
  contents: '',
  subaccountId: defaultSubaccountId,
  version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
};

export const defaultTradeContent: TradeContent = {
  id: 'defaultTradeId',
  size: '10',
  price: '10000',
  side: 'BUY',
  createdAt: 'createdAt',
  liquidation: true,
};
export const defaultTradeMessage: SingleTradeMessage = contentToSingleTradeMessage(
  defaultTradeContent,
  testConstants.defaultPerpetualMarket.clobPairId,
);
export const defaultTradeKafkaEvent:
ConsolidatedKafkaEvent = createConsolidatedKafkaEventFromTrade(defaultTradeMessage);

export const defaultStatefulOrderPlacementEvent: StatefulOrderEventV1 = {
  orderPlace: {
    order: {
      ...defaultMakerOrder,
      orderId: {
        ...defaultMakerOrder.orderId!,
        orderFlags: ORDER_FLAG_LONG_TERM,
      },
      goodTilBlockTime: 123,
    },
  },
};
export const defaultStatefulOrderRemovalEvent: StatefulOrderEventV1 = {
  orderRemoval: {
    removedOrderId: defaultOrderId,
    reason: OrderRemovalReason.ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED,
  },
};
export const defaultConditionalOrderPlacementEvent: StatefulOrderEventV1 = {
  conditionalOrderPlacement: {
    order: {
      ...defaultMakerOrder,
      orderId: {
        ...defaultMakerOrder.orderId!,
        orderFlags: ORDER_FLAG_CONDITIONAL,
      },
      conditionType: IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS,
      conditionalOrderTriggerSubticks: Long.fromValue(1000000, true),
      goodTilBlockTime: 123,
    },
  },
};
export const defaultConditionalOrderTriggeredEvent: StatefulOrderEventV1 = {
  conditionalOrderTriggered: {
    triggeredOrderId: {
      ...defaultOrderId,
      orderFlags: ORDER_FLAG_CONDITIONAL,
    },
  },
};
export const defaultLongTermOrderPlacementEvent: StatefulOrderEventV1 = {
  longTermOrderPlacement: {
    order: {
      ...defaultMakerOrder,
      orderId: {
        ...defaultMakerOrder.orderId!,
        orderFlags: ORDER_FLAG_LONG_TERM,
      },
      goodTilBlockTime: 123,
    },
  },
};
