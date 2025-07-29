import { Liquidity } from '@dydxprotocol-indexer/postgres';
import {
  IndexerAssetPosition,
  IndexerOrder,
  IndexerPerpetualPosition,
  IndexerSubaccountId, LiquidationOrderV1,
} from '@dydxprotocol-indexer/v4-protos';
import { Long } from '@dydxprotocol-indexer/v4-protos/build/codegen/helpers';

/* Canonical object types for handling onchain messages from the protocol. */

export interface SubaccountUpdate {
  subaccountId?: IndexerSubaccountId,
  updatedPerpetualPositions: IndexerPerpetualPosition[],
  updatedAssetPositions: IndexerAssetPosition[],
}

export interface OrderFillWithLiquidity {
  makerOrder?: IndexerOrder,
  order?: IndexerOrder,
  liquidationOrder?: LiquidationOrderV1,
  /** Fill amount in base quantums. */
  fillAmount: Long,
  /** Maker fee in USDC quantums. */
  makerFee: Long,
  /**
   * Taker fee in USDC quantums. If the taker order is a liquidation, then this
   * represents the special liquidation fee, not the standard taker fee.
   */
  takerFee: Long,
  /** Total filled of the maker order in base quantums. */
  totalFilledMaker: Long,
  /** Total filled of the taker order in base quantums. */
  totalFilledTaker: Long,
  /** Liquidity of the order in the match to process in the handler. */
  liquidity: Liquidity,
  /** Affiliate rev share in USDC quantums. */
  affiliateRevShare: Long,
  /** Builder fee for the maker in USDC quantums. */
  makerBuilderFee: Long,
  /** Builder fee for the taker in USDC quantums. */
  takerBuilderFee: Long,
  /** Builder address for the maker. */
  makerBuilderAddress: string,
  /** Builder address for the taker. */
  takerBuilderAddress: string,
  /** Order router fee for the maker in USDC quantums. */
  makerOrderRouterFee: Long,
  /** Order router fee for the taker in USDC quantums. */
  takerOrderRouterFee: Long,
  /** Order router address for the maker. */
  makerOrderRouterAddress: string,
  /** Order router address for the taker. */
  takerOrderRouterAddress: string,
}
