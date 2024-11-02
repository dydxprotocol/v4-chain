import { logger } from '@klyraprotocol-indexer/base';
import {
  AssetCreateEventV1,
  DeleveragingEventV1,
  FundingEventV1,
  IndexerTendermintEvent,
  LiquidityTierUpsertEventV1,
  MarketEventV1,
  OrderFillEventV1,
  PerpetualMarketCreateEventV1,
  StatefulOrderEventV1,
  SubaccountUpdateEventV1,
  TransferEventV1,
  UpdateClobPairEventV1,
  UpdatePerpetualEventV1,
} from '@klyraprotocol-indexer/v4-protos';

import { AnnotatedIndexerTendermintEvent, KlyraIndexerSubtypes } from './types';

export function annotateIndexerTendermintEvent(
  event: IndexerTendermintEvent,
): AnnotatedIndexerTendermintEvent | undefined {
  const eventDataBinary: Uint8Array = event.dataBytes;
  switch (event.subtype) {
    case (KlyraIndexerSubtypes.ORDER_FILL.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(OrderFillEventV1.decode(eventDataBinary)),
      };
    }
    case (KlyraIndexerSubtypes.SUBACCOUNT_UPDATE.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(SubaccountUpdateEventV1.decode(eventDataBinary)),
      };
    }
    case (KlyraIndexerSubtypes.TRANSFER.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(TransferEventV1.decode(eventDataBinary)),
      };
    }
    case (KlyraIndexerSubtypes.MARKET.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(MarketEventV1.decode(eventDataBinary)),
      };
    }
    case (KlyraIndexerSubtypes.STATEFUL_ORDER.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(StatefulOrderEventV1.decode(eventDataBinary)),
      };
    }
    case (KlyraIndexerSubtypes.FUNDING.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(FundingEventV1.decode(eventDataBinary)),
      };
    }
    case (KlyraIndexerSubtypes.ASSET.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(AssetCreateEventV1.decode(eventDataBinary)),
      };
    }
    case (KlyraIndexerSubtypes.PERPETUAL_MARKET.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(PerpetualMarketCreateEventV1.decode(eventDataBinary)),
      };
    }
    case (KlyraIndexerSubtypes.LIQUIDITY_TIER.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(LiquidityTierUpsertEventV1.decode(eventDataBinary)),
      };
    }
    case (KlyraIndexerSubtypes.UPDATE_PERPETUAL.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(UpdatePerpetualEventV1.decode(eventDataBinary)),
      };
    }
    case (KlyraIndexerSubtypes.UPDATE_CLOB_PAIR.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(UpdateClobPairEventV1.decode(eventDataBinary)),
      };
    }
    case (KlyraIndexerSubtypes.DELEVERAGING.toString()): {
      return {
        ...event,
        dataBytes: new Uint8Array(),
        data: JSON.stringify(DeleveragingEventV1.decode(eventDataBinary)),
      };
    }
    default: {
      const message: string = `Unable to parse event subtype: ${event.subtype}`;
      logger.error({
        at: 'block-helpers#annotateIndexerTendermintEvent',
        message,
      });
      return undefined;
    }
  }
}
