import { IndexerSubaccountId, IndexerSubaccountIdSDKType, IndexerPerpetualPosition, IndexerPerpetualPositionSDKType, IndexerAssetPosition, IndexerAssetPositionSDKType } from "../protocol/v1/subaccount";
import { IndexerOrder, IndexerOrderSDKType, IndexerOrderId, IndexerOrderIdSDKType, ClobPairStatus, ClobPairStatusSDKType } from "../protocol/v1/clob";
import { OrderRemovalReason, OrderRemovalReasonSDKType } from "../shared/removal_reason";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../../helpers";
/** Type is the type for funding values. */

export enum FundingEventV1_Type {
  /** TYPE_UNSPECIFIED - Unspecified type. */
  TYPE_UNSPECIFIED = 0,

  /**
   * TYPE_PREMIUM_SAMPLE - Premium sample is the combined value from all premium votes during a
   * `funding-sample` epoch.
   */
  TYPE_PREMIUM_SAMPLE = 1,

  /**
   * TYPE_FUNDING_RATE_AND_INDEX - Funding rate is the final funding rate combining all premium samples
   * during a `funding-tick` epoch.
   */
  TYPE_FUNDING_RATE_AND_INDEX = 2,

  /**
   * TYPE_PREMIUM_VOTE - TODO(DEC-1513): Investigate whether premium vote values need to be
   * sent to indexer.
   */
  TYPE_PREMIUM_VOTE = 3,
  UNRECOGNIZED = -1,
}
/** Type is the type for funding values. */

export enum FundingEventV1_TypeSDKType {
  /** TYPE_UNSPECIFIED - Unspecified type. */
  TYPE_UNSPECIFIED = 0,

  /**
   * TYPE_PREMIUM_SAMPLE - Premium sample is the combined value from all premium votes during a
   * `funding-sample` epoch.
   */
  TYPE_PREMIUM_SAMPLE = 1,

  /**
   * TYPE_FUNDING_RATE_AND_INDEX - Funding rate is the final funding rate combining all premium samples
   * during a `funding-tick` epoch.
   */
  TYPE_FUNDING_RATE_AND_INDEX = 2,

  /**
   * TYPE_PREMIUM_VOTE - TODO(DEC-1513): Investigate whether premium vote values need to be
   * sent to indexer.
   */
  TYPE_PREMIUM_VOTE = 3,
  UNRECOGNIZED = -1,
}
export function fundingEventV1_TypeFromJSON(object: any): FundingEventV1_Type {
  switch (object) {
    case 0:
    case "TYPE_UNSPECIFIED":
      return FundingEventV1_Type.TYPE_UNSPECIFIED;

    case 1:
    case "TYPE_PREMIUM_SAMPLE":
      return FundingEventV1_Type.TYPE_PREMIUM_SAMPLE;

    case 2:
    case "TYPE_FUNDING_RATE_AND_INDEX":
      return FundingEventV1_Type.TYPE_FUNDING_RATE_AND_INDEX;

    case 3:
    case "TYPE_PREMIUM_VOTE":
      return FundingEventV1_Type.TYPE_PREMIUM_VOTE;

    case -1:
    case "UNRECOGNIZED":
    default:
      return FundingEventV1_Type.UNRECOGNIZED;
  }
}
export function fundingEventV1_TypeToJSON(object: FundingEventV1_Type): string {
  switch (object) {
    case FundingEventV1_Type.TYPE_UNSPECIFIED:
      return "TYPE_UNSPECIFIED";

    case FundingEventV1_Type.TYPE_PREMIUM_SAMPLE:
      return "TYPE_PREMIUM_SAMPLE";

    case FundingEventV1_Type.TYPE_FUNDING_RATE_AND_INDEX:
      return "TYPE_FUNDING_RATE_AND_INDEX";

    case FundingEventV1_Type.TYPE_PREMIUM_VOTE:
      return "TYPE_PREMIUM_VOTE";

    case FundingEventV1_Type.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/**
 * FundingUpdate is used for funding update events and includes a funding
 * value and an optional funding index that correspond to a perpetual market.
 */

export interface FundingUpdateV1 {
  /** The id of the perpetual market. */
  perpetualId: number;
  /**
   * funding value (in parts-per-million) can be premium vote, premium sample,
   * or funding rate.
   */

  fundingValuePpm: number;
  /**
   * funding index is required if and only if parent `FundingEvent` type is
   * `TYPE_FUNDING_RATE_AND_INDEX`.
   */

  fundingIndex: Uint8Array;
}
/**
 * FundingUpdate is used for funding update events and includes a funding
 * value and an optional funding index that correspond to a perpetual market.
 */

export interface FundingUpdateV1SDKType {
  /** The id of the perpetual market. */
  perpetual_id: number;
  /**
   * funding value (in parts-per-million) can be premium vote, premium sample,
   * or funding rate.
   */

  funding_value_ppm: number;
  /**
   * funding index is required if and only if parent `FundingEvent` type is
   * `TYPE_FUNDING_RATE_AND_INDEX`.
   */

  funding_index: Uint8Array;
}
/**
 * FundingEvent message contains a list of per-market funding values. The
 * funding values in the list is of the same type and the types are: which can
 * have one of the following types:
 * 1. Premium vote: votes on the premium values injected by block proposers.
 * 2. Premium sample: combined value from all premium votes during a
 *    `funding-sample` epoch.
 * 3. Funding rate and index: final funding rate combining all premium samples
 *    during a `funding-tick` epoch and funding index accordingly updated with
 *    `funding rate * price`.
 */

export interface FundingEventV1 {
  /**
   * updates is a list of per-market funding updates for all existing perpetual
   * markets. The list is sorted by `perpetualId`s which are unique.
   */
  updates: FundingUpdateV1[];
  /** type stores the type of funding updates. */

  type: FundingEventV1_Type;
}
/**
 * FundingEvent message contains a list of per-market funding values. The
 * funding values in the list is of the same type and the types are: which can
 * have one of the following types:
 * 1. Premium vote: votes on the premium values injected by block proposers.
 * 2. Premium sample: combined value from all premium votes during a
 *    `funding-sample` epoch.
 * 3. Funding rate and index: final funding rate combining all premium samples
 *    during a `funding-tick` epoch and funding index accordingly updated with
 *    `funding rate * price`.
 */

export interface FundingEventV1SDKType {
  /**
   * updates is a list of per-market funding updates for all existing perpetual
   * markets. The list is sorted by `perpetualId`s which are unique.
   */
  updates: FundingUpdateV1SDKType[];
  /** type stores the type of funding updates. */

  type: FundingEventV1_TypeSDKType;
}
/**
 * MarketEvent message contains all the information about a market event on
 * the dYdX chain.
 */

export interface MarketEventV1 {
  /** market id. */
  marketId: number;
  priceUpdate?: MarketPriceUpdateEventV1;
  marketCreate?: MarketCreateEventV1;
  marketModify?: MarketModifyEventV1;
}
/**
 * MarketEvent message contains all the information about a market event on
 * the dYdX chain.
 */

export interface MarketEventV1SDKType {
  /** market id. */
  market_id: number;
  price_update?: MarketPriceUpdateEventV1SDKType;
  market_create?: MarketCreateEventV1SDKType;
  market_modify?: MarketModifyEventV1SDKType;
}
/**
 * MarketPriceUpdateEvent message contains all the information about a price
 * update on the dYdX chain.
 */

export interface MarketPriceUpdateEventV1 {
  /**
   * price_with_exponent. Multiply by 10 ^ Exponent to get the human readable
   * price in dollars. For example if `Exponent == -5` then a `exponent_price`
   * of `1,000,000,000` represents “$10,000`.
   */
  priceWithExponent: Long;
}
/**
 * MarketPriceUpdateEvent message contains all the information about a price
 * update on the dYdX chain.
 */

export interface MarketPriceUpdateEventV1SDKType {
  /**
   * price_with_exponent. Multiply by 10 ^ Exponent to get the human readable
   * price in dollars. For example if `Exponent == -5` then a `exponent_price`
   * of `1,000,000,000` represents “$10,000`.
   */
  price_with_exponent: Long;
}
/** shared fields between MarketCreateEvent and MarketModifyEvent */

export interface MarketBaseEventV1 {
  /** String representation of the market pair, e.g. `BTC-USD` */
  pair: string;
  /**
   * The minimum allowable change in the Price value for a given update.
   * Measured as 1e-6.
   */

  minPriceChangePpm: number;
}
/** shared fields between MarketCreateEvent and MarketModifyEvent */

export interface MarketBaseEventV1SDKType {
  /** String representation of the market pair, e.g. `BTC-USD` */
  pair: string;
  /**
   * The minimum allowable change in the Price value for a given update.
   * Measured as 1e-6.
   */

  min_price_change_ppm: number;
}
/**
 * MarketCreateEvent message contains all the information about a new market on
 * the dYdX chain.
 */

export interface MarketCreateEventV1 {
  base?: MarketBaseEventV1;
  /**
   * Static value. The exponent of the price.
   * For example if Exponent == -5 then a `exponent_price` of 1,000,000,000
   * represents $10,000. Therefore 10 ^ Exponent represents the smallest
   * price step (in dollars) that can be recorded.
   */

  exponent: number;
}
/**
 * MarketCreateEvent message contains all the information about a new market on
 * the dYdX chain.
 */

export interface MarketCreateEventV1SDKType {
  base?: MarketBaseEventV1SDKType;
  /**
   * Static value. The exponent of the price.
   * For example if Exponent == -5 then a `exponent_price` of 1,000,000,000
   * represents $10,000. Therefore 10 ^ Exponent represents the smallest
   * price step (in dollars) that can be recorded.
   */

  exponent: number;
}
/**
 * MarketModifyEvent message contains all the information about a market update
 * on the dYdX chain
 */

export interface MarketModifyEventV1 {
  /**
   * MarketModifyEvent message contains all the information about a market update
   * on the dYdX chain
   */
  base?: MarketBaseEventV1;
}
/**
 * MarketModifyEvent message contains all the information about a market update
 * on the dYdX chain
 */

export interface MarketModifyEventV1SDKType {
  /**
   * MarketModifyEvent message contains all the information about a market update
   * on the dYdX chain
   */
  base?: MarketBaseEventV1SDKType;
}
/** SourceOfFunds is the source of funds in a transfer event. */

export interface SourceOfFunds {
  subaccountId?: IndexerSubaccountId;
  address?: string;
}
/** SourceOfFunds is the source of funds in a transfer event. */

export interface SourceOfFundsSDKType {
  subaccount_id?: IndexerSubaccountIdSDKType;
  address?: string;
}
/**
 * TransferEvent message contains all the information about a transfer,
 * deposit-to-subaccount, or withdraw-from-subaccount on the dYdX chain.
 * When a subaccount is involved, a SubaccountUpdateEvent message will
 * be produced with the updated asset positions.
 */

export interface TransferEventV1 {
  senderSubaccountId?: IndexerSubaccountId;
  recipientSubaccountId?: IndexerSubaccountId;
  /** Id of the asset transfered. */

  assetId: number;
  /** The amount of asset in quantums to transfer. */

  amount: Long;
  /**
   * The sender is one of below
   * - a subaccount ID (in transfer and withdraw events).
   * - a wallet address (in deposit events).
   */

  sender?: SourceOfFunds;
  /**
   * The recipient is one of below
   * - a subaccount ID (in transfer and deposit events).
   * - a wallet address (in withdraw events).
   */

  recipient?: SourceOfFunds;
}
/**
 * TransferEvent message contains all the information about a transfer,
 * deposit-to-subaccount, or withdraw-from-subaccount on the dYdX chain.
 * When a subaccount is involved, a SubaccountUpdateEvent message will
 * be produced with the updated asset positions.
 */

export interface TransferEventV1SDKType {
  sender_subaccount_id?: IndexerSubaccountIdSDKType;
  recipient_subaccount_id?: IndexerSubaccountIdSDKType;
  /** Id of the asset transfered. */

  asset_id: number;
  /** The amount of asset in quantums to transfer. */

  amount: Long;
  /**
   * The sender is one of below
   * - a subaccount ID (in transfer and withdraw events).
   * - a wallet address (in deposit events).
   */

  sender?: SourceOfFundsSDKType;
  /**
   * The recipient is one of below
   * - a subaccount ID (in transfer and deposit events).
   * - a wallet address (in withdraw events).
   */

  recipient?: SourceOfFundsSDKType;
}
/**
 * OrderFillEvent message contains all the information from an order match in
 * the dYdX chain. This includes the maker/taker orders that matched and the
 * amount filled.
 */

export interface OrderFillEventV1 {
  makerOrder?: IndexerOrder;
  order?: IndexerOrder;
  liquidationOrder?: LiquidationOrderV1;
  /** Fill amount in base quantums. */

  fillAmount: Long;
  /** Maker fee in USDC quantums. */

  makerFee: Long;
  /**
   * Taker fee in USDC quantums. If the taker order is a liquidation, then this
   * represents the special liquidation fee, not the standard taker fee.
   */

  takerFee: Long;
  /** Total filled of the maker order in base quantums. */

  totalFilledMaker: Long;
  /** Total filled of the taker order in base quantums. */

  totalFilledTaker: Long;
}
/**
 * OrderFillEvent message contains all the information from an order match in
 * the dYdX chain. This includes the maker/taker orders that matched and the
 * amount filled.
 */

export interface OrderFillEventV1SDKType {
  maker_order?: IndexerOrderSDKType;
  order?: IndexerOrderSDKType;
  liquidation_order?: LiquidationOrderV1SDKType;
  /** Fill amount in base quantums. */

  fill_amount: Long;
  /** Maker fee in USDC quantums. */

  maker_fee: Long;
  /**
   * Taker fee in USDC quantums. If the taker order is a liquidation, then this
   * represents the special liquidation fee, not the standard taker fee.
   */

  taker_fee: Long;
  /** Total filled of the maker order in base quantums. */

  total_filled_maker: Long;
  /** Total filled of the taker order in base quantums. */

  total_filled_taker: Long;
}
/**
 * DeleveragingEvent message contains all the information for a deleveraging
 * on the dYdX chain. This includes the liquidated/offsetting subaccounts and
 * the amount filled.
 */

export interface DeleveragingEventV1 {
  /** ID of the subaccount that was liquidated. */
  liquidated?: IndexerSubaccountId;
  /** ID of the subaccount that was used to offset the position. */

  offsetting?: IndexerSubaccountId;
  /** The ID of the perpetual that was liquidated. */

  perpetualId: number;
  /**
   * The amount filled between the liquidated and offsetting position, in
   * base quantums.
   */

  fillAmount: Long;
  /** Fill price of deleveraging event, in USDC quote quantums. */

  price: Long;
  /** `true` if liquidating a short position, `false` otherwise. */

  isBuy: boolean;
  /**
   * `true` if the deleveraging event is for final settlement, indicating
   * the match occurred at the oracle price rather than bankruptcy price.
   * When this flag is `false`, the fill price is the bankruptcy price
   * of the liquidated subaccount.
   */

  isFinalSettlement: boolean;
}
/**
 * DeleveragingEvent message contains all the information for a deleveraging
 * on the dYdX chain. This includes the liquidated/offsetting subaccounts and
 * the amount filled.
 */

export interface DeleveragingEventV1SDKType {
  /** ID of the subaccount that was liquidated. */
  liquidated?: IndexerSubaccountIdSDKType;
  /** ID of the subaccount that was used to offset the position. */

  offsetting?: IndexerSubaccountIdSDKType;
  /** The ID of the perpetual that was liquidated. */

  perpetual_id: number;
  /**
   * The amount filled between the liquidated and offsetting position, in
   * base quantums.
   */

  fill_amount: Long;
  /** Fill price of deleveraging event, in USDC quote quantums. */

  price: Long;
  /** `true` if liquidating a short position, `false` otherwise. */

  is_buy: boolean;
  /**
   * `true` if the deleveraging event is for final settlement, indicating
   * the match occurred at the oracle price rather than bankruptcy price.
   * When this flag is `false`, the fill price is the bankruptcy price
   * of the liquidated subaccount.
   */

  is_final_settlement: boolean;
}
/**
 * LiquidationOrder represents the liquidation taker order to be included in a
 * liquidation order fill event.
 */

export interface LiquidationOrderV1 {
  /** ID of the subaccount that was liquidated. */
  liquidated?: IndexerSubaccountId;
  /** The ID of the clob pair involved in the liquidation. */

  clobPairId: number;
  /** The ID of the perpetual involved in the liquidation. */

  perpetualId: number;
  /**
   * The total size of the liquidation order including any unfilled size,
   * in base quantums.
   */

  totalSize: Long;
  /** `true` if liquidating a short position, `false` otherwise. */

  isBuy: boolean;
  /**
   * The fillable price in subticks.
   * This represents the lower-price-bound for liquidating longs
   * and the upper-price-bound for liquidating shorts.
   * Must be a multiple of ClobPair.SubticksPerTick
   * (where `ClobPair.Id = orderId.ClobPairId`).
   */

  subticks: Long;
}
/**
 * LiquidationOrder represents the liquidation taker order to be included in a
 * liquidation order fill event.
 */

export interface LiquidationOrderV1SDKType {
  /** ID of the subaccount that was liquidated. */
  liquidated?: IndexerSubaccountIdSDKType;
  /** The ID of the clob pair involved in the liquidation. */

  clob_pair_id: number;
  /** The ID of the perpetual involved in the liquidation. */

  perpetual_id: number;
  /**
   * The total size of the liquidation order including any unfilled size,
   * in base quantums.
   */

  total_size: Long;
  /** `true` if liquidating a short position, `false` otherwise. */

  is_buy: boolean;
  /**
   * The fillable price in subticks.
   * This represents the lower-price-bound for liquidating longs
   * and the upper-price-bound for liquidating shorts.
   * Must be a multiple of ClobPair.SubticksPerTick
   * (where `ClobPair.Id = orderId.ClobPairId`).
   */

  subticks: Long;
}
/**
 * SubaccountUpdateEvent message contains information about an update to a
 * subaccount in the dYdX chain. This includes the list of updated perpetual
 * and asset positions for the subaccount.
 * Note: This event message will contain all the updates to a subaccount
 * at the end of a block which is why multiple asset/perpetual position
 * updates may exist.
 */

export interface SubaccountUpdateEventV1 {
  subaccountId?: IndexerSubaccountId;
  updatedPerpetualPositions: IndexerPerpetualPosition[];
  updatedAssetPositions: IndexerAssetPosition[];
}
/**
 * SubaccountUpdateEvent message contains information about an update to a
 * subaccount in the dYdX chain. This includes the list of updated perpetual
 * and asset positions for the subaccount.
 * Note: This event message will contain all the updates to a subaccount
 * at the end of a block which is why multiple asset/perpetual position
 * updates may exist.
 */

export interface SubaccountUpdateEventV1SDKType {
  subaccount_id?: IndexerSubaccountIdSDKType;
  updated_perpetual_positions: IndexerPerpetualPositionSDKType[];
  updated_asset_positions: IndexerAssetPositionSDKType[];
}
/**
 * StatefulOrderEvent message contains information about a change to a stateful
 * order. Currently, this is either the placement of a long-term order, the
 * placement or triggering of a conditional order, or the removal of a
 * stateful order.
 */

export interface StatefulOrderEventV1 {
  orderPlace?: StatefulOrderEventV1_StatefulOrderPlacementV1;
  orderRemoval?: StatefulOrderEventV1_StatefulOrderRemovalV1;
  conditionalOrderPlacement?: StatefulOrderEventV1_ConditionalOrderPlacementV1;
  conditionalOrderTriggered?: StatefulOrderEventV1_ConditionalOrderTriggeredV1;
  longTermOrderPlacement?: StatefulOrderEventV1_LongTermOrderPlacementV1;
}
/**
 * StatefulOrderEvent message contains information about a change to a stateful
 * order. Currently, this is either the placement of a long-term order, the
 * placement or triggering of a conditional order, or the removal of a
 * stateful order.
 */

export interface StatefulOrderEventV1SDKType {
  order_place?: StatefulOrderEventV1_StatefulOrderPlacementV1SDKType;
  order_removal?: StatefulOrderEventV1_StatefulOrderRemovalV1SDKType;
  conditional_order_placement?: StatefulOrderEventV1_ConditionalOrderPlacementV1SDKType;
  conditional_order_triggered?: StatefulOrderEventV1_ConditionalOrderTriggeredV1SDKType;
  long_term_order_placement?: StatefulOrderEventV1_LongTermOrderPlacementV1SDKType;
}
/** A stateful order placement contains an order. */

export interface StatefulOrderEventV1_StatefulOrderPlacementV1 {
  order?: IndexerOrder;
}
/** A stateful order placement contains an order. */

export interface StatefulOrderEventV1_StatefulOrderPlacementV1SDKType {
  order?: IndexerOrderSDKType;
}
/**
 * A stateful order removal contains the id of an order that was already
 * placed and is now removed and the reason for the removal.
 */

export interface StatefulOrderEventV1_StatefulOrderRemovalV1 {
  removedOrderId?: IndexerOrderId;
  reason: OrderRemovalReason;
}
/**
 * A stateful order removal contains the id of an order that was already
 * placed and is now removed and the reason for the removal.
 */

export interface StatefulOrderEventV1_StatefulOrderRemovalV1SDKType {
  removed_order_id?: IndexerOrderIdSDKType;
  reason: OrderRemovalReasonSDKType;
}
/**
 * A conditional order placement contains an order. The order is newly-placed
 * and untriggered when this event is emitted.
 */

export interface StatefulOrderEventV1_ConditionalOrderPlacementV1 {
  order?: IndexerOrder;
}
/**
 * A conditional order placement contains an order. The order is newly-placed
 * and untriggered when this event is emitted.
 */

export interface StatefulOrderEventV1_ConditionalOrderPlacementV1SDKType {
  order?: IndexerOrderSDKType;
}
/**
 * A conditional order trigger event contains an order id and is emitted when
 * an order is triggered.
 */

export interface StatefulOrderEventV1_ConditionalOrderTriggeredV1 {
  triggeredOrderId?: IndexerOrderId;
}
/**
 * A conditional order trigger event contains an order id and is emitted when
 * an order is triggered.
 */

export interface StatefulOrderEventV1_ConditionalOrderTriggeredV1SDKType {
  triggered_order_id?: IndexerOrderIdSDKType;
}
/** A long term order placement contains an order. */

export interface StatefulOrderEventV1_LongTermOrderPlacementV1 {
  order?: IndexerOrder;
}
/** A long term order placement contains an order. */

export interface StatefulOrderEventV1_LongTermOrderPlacementV1SDKType {
  order?: IndexerOrderSDKType;
}
/**
 * AssetCreateEventV1 message contains all the information about an new Asset on
 * the dYdX chain.
 */

export interface AssetCreateEventV1 {
  /** Unique, sequentially-generated. */
  id: number;
  /**
   * The human readable symbol of the `Asset` (e.g. `USDC`, `ATOM`).
   * Must be uppercase, unique and correspond to the canonical symbol of the
   * full coin.
   */

  symbol: string;
  /** `true` if this `Asset` has a valid `MarketId` value. */

  hasMarket: boolean;
  /**
   * The `Id` of the `Market` associated with this `Asset`. It acts as the
   * oracle price for the purposes of calculating collateral
   * and margin requirements.
   */

  marketId: number;
  /**
   * The exponent for converting an atomic amount (1 'quantum')
   * to a full coin. For example, if `atomic_resolution = -8`
   * then an `asset_position` with `base_quantums = 1e8` is equivalent to
   * a position size of one full coin.
   */

  atomicResolution: number;
}
/**
 * AssetCreateEventV1 message contains all the information about an new Asset on
 * the dYdX chain.
 */

export interface AssetCreateEventV1SDKType {
  /** Unique, sequentially-generated. */
  id: number;
  /**
   * The human readable symbol of the `Asset` (e.g. `USDC`, `ATOM`).
   * Must be uppercase, unique and correspond to the canonical symbol of the
   * full coin.
   */

  symbol: string;
  /** `true` if this `Asset` has a valid `MarketId` value. */

  has_market: boolean;
  /**
   * The `Id` of the `Market` associated with this `Asset`. It acts as the
   * oracle price for the purposes of calculating collateral
   * and margin requirements.
   */

  market_id: number;
  /**
   * The exponent for converting an atomic amount (1 'quantum')
   * to a full coin. For example, if `atomic_resolution = -8`
   * then an `asset_position` with `base_quantums = 1e8` is equivalent to
   * a position size of one full coin.
   */

  atomic_resolution: number;
}
/**
 * PerpetualMarketCreateEventV1 message contains all the information about a
 * new Perpetual Market on the dYdX chain.
 */

export interface PerpetualMarketCreateEventV1 {
  /**
   * Unique Perpetual id.
   * Defined in perpetuals.perpetual
   */
  id: number;
  /**
   * Unique clob pair Id associated with this perpetual market
   * Defined in clob.clob_pair
   */

  clobPairId: number;
  /**
   * The name of the `Perpetual` (e.g. `BTC-USD`).
   * Defined in perpetuals.perpetual
   */

  ticker: string;
  /**
   * Unique id of market param associated with this perpetual market.
   * Defined in perpetuals.perpetual
   */

  marketId: number;
  /** Status of the CLOB */

  status: ClobPairStatus;
  /**
   * `10^Exponent` gives the number of QuoteQuantums traded per BaseQuantum
   * per Subtick.
   * Defined in clob.clob_pair
   */

  quantumConversionExponent: number;
  /**
   * The exponent for converting an atomic amount (`size = 1`)
   * to a full coin. For example, if `AtomicResolution = -8`
   * then a `PerpetualPosition` with `size = 1e8` is equivalent to
   * a position size of one full coin.
   * Defined in perpetuals.perpetual
   */

  atomicResolution: number;
  /**
   * Defines the tick size of the orderbook by defining how many subticks
   * are in one tick. That is, the subticks of any valid order must be a
   * multiple of this value. Generally this value should start `>= 100`to
   * allow room for decreasing it.
   * Defined in clob.clob_pair
   */

  subticksPerTick: number;
  /**
   * Minimum increment in the size of orders on the CLOB, in base quantums.
   * Defined in clob.clob_pair
   */

  stepBaseQuantums: Long;
  /**
   * The liquidity_tier that this perpetual is associated with.
   * Defined in perpetuals.perpetual
   */

  liquidityTier: number;
}
/**
 * PerpetualMarketCreateEventV1 message contains all the information about a
 * new Perpetual Market on the dYdX chain.
 */

export interface PerpetualMarketCreateEventV1SDKType {
  /**
   * Unique Perpetual id.
   * Defined in perpetuals.perpetual
   */
  id: number;
  /**
   * Unique clob pair Id associated with this perpetual market
   * Defined in clob.clob_pair
   */

  clob_pair_id: number;
  /**
   * The name of the `Perpetual` (e.g. `BTC-USD`).
   * Defined in perpetuals.perpetual
   */

  ticker: string;
  /**
   * Unique id of market param associated with this perpetual market.
   * Defined in perpetuals.perpetual
   */

  market_id: number;
  /** Status of the CLOB */

  status: ClobPairStatusSDKType;
  /**
   * `10^Exponent` gives the number of QuoteQuantums traded per BaseQuantum
   * per Subtick.
   * Defined in clob.clob_pair
   */

  quantum_conversion_exponent: number;
  /**
   * The exponent for converting an atomic amount (`size = 1`)
   * to a full coin. For example, if `AtomicResolution = -8`
   * then a `PerpetualPosition` with `size = 1e8` is equivalent to
   * a position size of one full coin.
   * Defined in perpetuals.perpetual
   */

  atomic_resolution: number;
  /**
   * Defines the tick size of the orderbook by defining how many subticks
   * are in one tick. That is, the subticks of any valid order must be a
   * multiple of this value. Generally this value should start `>= 100`to
   * allow room for decreasing it.
   * Defined in clob.clob_pair
   */

  subticks_per_tick: number;
  /**
   * Minimum increment in the size of orders on the CLOB, in base quantums.
   * Defined in clob.clob_pair
   */

  step_base_quantums: Long;
  /**
   * The liquidity_tier that this perpetual is associated with.
   * Defined in perpetuals.perpetual
   */

  liquidity_tier: number;
}
/**
 * LiquidityTierUpsertEventV1 message contains all the information to
 * create/update a Liquidity Tier on the dYdX chain.
 */

export interface LiquidityTierUpsertEventV1 {
  /** Unique id. */
  id: number;
  /** The name of the tier purely for mnemonic purposes, e.g. "Gold". */

  name: string;
  /**
   * The margin fraction needed to open a position.
   * In parts-per-million.
   */

  initialMarginPpm: number;
  /**
   * The fraction of the initial-margin that the maintenance-margin is,
   * e.g. 50%. In parts-per-million.
   */

  maintenanceFractionPpm: number;
  /**
   * The maximum position size at which the margin requirements are
   * not increased over the default values. Above this position size,
   * the margin requirements increase at a rate of sqrt(size).
   */

  basePositionNotional: Long;
}
/**
 * LiquidityTierUpsertEventV1 message contains all the information to
 * create/update a Liquidity Tier on the dYdX chain.
 */

export interface LiquidityTierUpsertEventV1SDKType {
  /** Unique id. */
  id: number;
  /** The name of the tier purely for mnemonic purposes, e.g. "Gold". */

  name: string;
  /**
   * The margin fraction needed to open a position.
   * In parts-per-million.
   */

  initial_margin_ppm: number;
  /**
   * The fraction of the initial-margin that the maintenance-margin is,
   * e.g. 50%. In parts-per-million.
   */

  maintenance_fraction_ppm: number;
  /**
   * The maximum position size at which the margin requirements are
   * not increased over the default values. Above this position size,
   * the margin requirements increase at a rate of sqrt(size).
   */

  base_position_notional: Long;
}
/**
 * UpdateClobPairEventV1 message contains all the information about an update to
 * a clob pair on the dYdX chain.
 */

export interface UpdateClobPairEventV1 {
  /**
   * Unique clob pair Id associated with this perpetual market
   * Defined in clob.clob_pair
   */
  clobPairId: number;
  /** Status of the CLOB */

  status: ClobPairStatus;
  /**
   * `10^Exponent` gives the number of QuoteQuantums traded per BaseQuantum
   * per Subtick.
   * Defined in clob.clob_pair
   */

  quantumConversionExponent: number;
  /**
   * Defines the tick size of the orderbook by defining how many subticks
   * are in one tick. That is, the subticks of any valid order must be a
   * multiple of this value. Generally this value should start `>= 100`to
   * allow room for decreasing it.
   * Defined in clob.clob_pair
   */

  subticksPerTick: number;
  /**
   * Minimum increment in the size of orders on the CLOB, in base quantums.
   * Defined in clob.clob_pair
   */

  stepBaseQuantums: Long;
}
/**
 * UpdateClobPairEventV1 message contains all the information about an update to
 * a clob pair on the dYdX chain.
 */

export interface UpdateClobPairEventV1SDKType {
  /**
   * Unique clob pair Id associated with this perpetual market
   * Defined in clob.clob_pair
   */
  clob_pair_id: number;
  /** Status of the CLOB */

  status: ClobPairStatusSDKType;
  /**
   * `10^Exponent` gives the number of QuoteQuantums traded per BaseQuantum
   * per Subtick.
   * Defined in clob.clob_pair
   */

  quantum_conversion_exponent: number;
  /**
   * Defines the tick size of the orderbook by defining how many subticks
   * are in one tick. That is, the subticks of any valid order must be a
   * multiple of this value. Generally this value should start `>= 100`to
   * allow room for decreasing it.
   * Defined in clob.clob_pair
   */

  subticks_per_tick: number;
  /**
   * Minimum increment in the size of orders on the CLOB, in base quantums.
   * Defined in clob.clob_pair
   */

  step_base_quantums: Long;
}
/**
 * UpdatePerpetualEventV1 message contains all the information about an update
 * to a perpetual on the dYdX chain.
 */

export interface UpdatePerpetualEventV1 {
  /**
   * Unique Perpetual id.
   * Defined in perpetuals.perpetual
   */
  id: number;
  /**
   * The name of the `Perpetual` (e.g. `BTC-USD`).
   * Defined in perpetuals.perpetual
   */

  ticker: string;
  /**
   * Unique id of market param associated with this perpetual market.
   * Defined in perpetuals.perpetual
   */

  marketId: number;
  /**
   * The exponent for converting an atomic amount (`size = 1`)
   * to a full coin. For example, if `AtomicResolution = -8`
   * then a `PerpetualPosition` with `size = 1e8` is equivalent to
   * a position size of one full coin.
   * Defined in perpetuals.perpetual
   */

  atomicResolution: number;
  /**
   * The liquidity_tier that this perpetual is associated with.
   * Defined in perpetuals.perpetual
   */

  liquidityTier: number;
}
/**
 * UpdatePerpetualEventV1 message contains all the information about an update
 * to a perpetual on the dYdX chain.
 */

export interface UpdatePerpetualEventV1SDKType {
  /**
   * Unique Perpetual id.
   * Defined in perpetuals.perpetual
   */
  id: number;
  /**
   * The name of the `Perpetual` (e.g. `BTC-USD`).
   * Defined in perpetuals.perpetual
   */

  ticker: string;
  /**
   * Unique id of market param associated with this perpetual market.
   * Defined in perpetuals.perpetual
   */

  market_id: number;
  /**
   * The exponent for converting an atomic amount (`size = 1`)
   * to a full coin. For example, if `AtomicResolution = -8`
   * then a `PerpetualPosition` with `size = 1e8` is equivalent to
   * a position size of one full coin.
   * Defined in perpetuals.perpetual
   */

  atomic_resolution: number;
  /**
   * The liquidity_tier that this perpetual is associated with.
   * Defined in perpetuals.perpetual
   */

  liquidity_tier: number;
}
/**
 * TradingRewardEventV1 is communicates all trading rewards for all accounts
 * that receive trade rewards in the block.
 */

export interface TradingRewardEventV1 {
  /** The list of all trading rewards in the block. */
  tradingRewards: AddressTradingReward[];
}
/**
 * TradingRewardEventV1 is communicates all trading rewards for all accounts
 * that receive trade rewards in the block.
 */

export interface TradingRewardEventV1SDKType {
  /** The list of all trading rewards in the block. */
  trading_rewards: AddressTradingRewardSDKType[];
}
/**
 * AddressTradingReward contains info on an instance of an address receiving a
 * reward
 */

export interface AddressTradingReward {
  /** The address of the wallet that will receive the trading reward. */
  owner: string;
  /**
   * The trading rewards earned by the address above in denoms. 1e18 denoms is
   * equivalent to a single coin.
   */

  denoms: Uint8Array;
}
/**
 * AddressTradingReward contains info on an instance of an address receiving a
 * reward
 */

export interface AddressTradingRewardSDKType {
  /** The address of the wallet that will receive the trading reward. */
  owner: string;
  /**
   * The trading rewards earned by the address above in denoms. 1e18 denoms is
   * equivalent to a single coin.
   */

  denoms: Uint8Array;
}

function createBaseFundingUpdateV1(): FundingUpdateV1 {
  return {
    perpetualId: 0,
    fundingValuePpm: 0,
    fundingIndex: new Uint8Array()
  };
}

export const FundingUpdateV1 = {
  encode(message: FundingUpdateV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.perpetualId !== 0) {
      writer.uint32(8).uint32(message.perpetualId);
    }

    if (message.fundingValuePpm !== 0) {
      writer.uint32(16).int32(message.fundingValuePpm);
    }

    if (message.fundingIndex.length !== 0) {
      writer.uint32(26).bytes(message.fundingIndex);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FundingUpdateV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFundingUpdateV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.perpetualId = reader.uint32();
          break;

        case 2:
          message.fundingValuePpm = reader.int32();
          break;

        case 3:
          message.fundingIndex = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<FundingUpdateV1>): FundingUpdateV1 {
    const message = createBaseFundingUpdateV1();
    message.perpetualId = object.perpetualId ?? 0;
    message.fundingValuePpm = object.fundingValuePpm ?? 0;
    message.fundingIndex = object.fundingIndex ?? new Uint8Array();
    return message;
  }

};

function createBaseFundingEventV1(): FundingEventV1 {
  return {
    updates: [],
    type: 0
  };
}

export const FundingEventV1 = {
  encode(message: FundingEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.updates) {
      FundingUpdateV1.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.type !== 0) {
      writer.uint32(16).int32(message.type);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FundingEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFundingEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.updates.push(FundingUpdateV1.decode(reader, reader.uint32()));
          break;

        case 2:
          message.type = (reader.int32() as any);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<FundingEventV1>): FundingEventV1 {
    const message = createBaseFundingEventV1();
    message.updates = object.updates?.map(e => FundingUpdateV1.fromPartial(e)) || [];
    message.type = object.type ?? 0;
    return message;
  }

};

function createBaseMarketEventV1(): MarketEventV1 {
  return {
    marketId: 0,
    priceUpdate: undefined,
    marketCreate: undefined,
    marketModify: undefined
  };
}

export const MarketEventV1 = {
  encode(message: MarketEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.marketId !== 0) {
      writer.uint32(8).uint32(message.marketId);
    }

    if (message.priceUpdate !== undefined) {
      MarketPriceUpdateEventV1.encode(message.priceUpdate, writer.uint32(18).fork()).ldelim();
    }

    if (message.marketCreate !== undefined) {
      MarketCreateEventV1.encode(message.marketCreate, writer.uint32(26).fork()).ldelim();
    }

    if (message.marketModify !== undefined) {
      MarketModifyEventV1.encode(message.marketModify, writer.uint32(34).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketId = reader.uint32();
          break;

        case 2:
          message.priceUpdate = MarketPriceUpdateEventV1.decode(reader, reader.uint32());
          break;

        case 3:
          message.marketCreate = MarketCreateEventV1.decode(reader, reader.uint32());
          break;

        case 4:
          message.marketModify = MarketModifyEventV1.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketEventV1>): MarketEventV1 {
    const message = createBaseMarketEventV1();
    message.marketId = object.marketId ?? 0;
    message.priceUpdate = object.priceUpdate !== undefined && object.priceUpdate !== null ? MarketPriceUpdateEventV1.fromPartial(object.priceUpdate) : undefined;
    message.marketCreate = object.marketCreate !== undefined && object.marketCreate !== null ? MarketCreateEventV1.fromPartial(object.marketCreate) : undefined;
    message.marketModify = object.marketModify !== undefined && object.marketModify !== null ? MarketModifyEventV1.fromPartial(object.marketModify) : undefined;
    return message;
  }

};

function createBaseMarketPriceUpdateEventV1(): MarketPriceUpdateEventV1 {
  return {
    priceWithExponent: Long.UZERO
  };
}

export const MarketPriceUpdateEventV1 = {
  encode(message: MarketPriceUpdateEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.priceWithExponent.isZero()) {
      writer.uint32(8).uint64(message.priceWithExponent);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketPriceUpdateEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketPriceUpdateEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.priceWithExponent = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketPriceUpdateEventV1>): MarketPriceUpdateEventV1 {
    const message = createBaseMarketPriceUpdateEventV1();
    message.priceWithExponent = object.priceWithExponent !== undefined && object.priceWithExponent !== null ? Long.fromValue(object.priceWithExponent) : Long.UZERO;
    return message;
  }

};

function createBaseMarketBaseEventV1(): MarketBaseEventV1 {
  return {
    pair: "",
    minPriceChangePpm: 0
  };
}

export const MarketBaseEventV1 = {
  encode(message: MarketBaseEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pair !== "") {
      writer.uint32(10).string(message.pair);
    }

    if (message.minPriceChangePpm !== 0) {
      writer.uint32(16).uint32(message.minPriceChangePpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketBaseEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketBaseEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pair = reader.string();
          break;

        case 2:
          message.minPriceChangePpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketBaseEventV1>): MarketBaseEventV1 {
    const message = createBaseMarketBaseEventV1();
    message.pair = object.pair ?? "";
    message.minPriceChangePpm = object.minPriceChangePpm ?? 0;
    return message;
  }

};

function createBaseMarketCreateEventV1(): MarketCreateEventV1 {
  return {
    base: undefined,
    exponent: 0
  };
}

export const MarketCreateEventV1 = {
  encode(message: MarketCreateEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.base !== undefined) {
      MarketBaseEventV1.encode(message.base, writer.uint32(10).fork()).ldelim();
    }

    if (message.exponent !== 0) {
      writer.uint32(16).sint32(message.exponent);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketCreateEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketCreateEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.base = MarketBaseEventV1.decode(reader, reader.uint32());
          break;

        case 2:
          message.exponent = reader.sint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketCreateEventV1>): MarketCreateEventV1 {
    const message = createBaseMarketCreateEventV1();
    message.base = object.base !== undefined && object.base !== null ? MarketBaseEventV1.fromPartial(object.base) : undefined;
    message.exponent = object.exponent ?? 0;
    return message;
  }

};

function createBaseMarketModifyEventV1(): MarketModifyEventV1 {
  return {
    base: undefined
  };
}

export const MarketModifyEventV1 = {
  encode(message: MarketModifyEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.base !== undefined) {
      MarketBaseEventV1.encode(message.base, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketModifyEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketModifyEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.base = MarketBaseEventV1.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketModifyEventV1>): MarketModifyEventV1 {
    const message = createBaseMarketModifyEventV1();
    message.base = object.base !== undefined && object.base !== null ? MarketBaseEventV1.fromPartial(object.base) : undefined;
    return message;
  }

};

function createBaseSourceOfFunds(): SourceOfFunds {
  return {
    subaccountId: undefined,
    address: undefined
  };
}

export const SourceOfFunds = {
  encode(message: SourceOfFunds, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.subaccountId !== undefined) {
      IndexerSubaccountId.encode(message.subaccountId, writer.uint32(10).fork()).ldelim();
    }

    if (message.address !== undefined) {
      writer.uint32(18).string(message.address);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SourceOfFunds {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSourceOfFunds();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountId = IndexerSubaccountId.decode(reader, reader.uint32());
          break;

        case 2:
          message.address = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<SourceOfFunds>): SourceOfFunds {
    const message = createBaseSourceOfFunds();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? IndexerSubaccountId.fromPartial(object.subaccountId) : undefined;
    message.address = object.address ?? undefined;
    return message;
  }

};

function createBaseTransferEventV1(): TransferEventV1 {
  return {
    senderSubaccountId: undefined,
    recipientSubaccountId: undefined,
    assetId: 0,
    amount: Long.UZERO,
    sender: undefined,
    recipient: undefined
  };
}

export const TransferEventV1 = {
  encode(message: TransferEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.senderSubaccountId !== undefined) {
      IndexerSubaccountId.encode(message.senderSubaccountId, writer.uint32(10).fork()).ldelim();
    }

    if (message.recipientSubaccountId !== undefined) {
      IndexerSubaccountId.encode(message.recipientSubaccountId, writer.uint32(18).fork()).ldelim();
    }

    if (message.assetId !== 0) {
      writer.uint32(24).uint32(message.assetId);
    }

    if (!message.amount.isZero()) {
      writer.uint32(32).uint64(message.amount);
    }

    if (message.sender !== undefined) {
      SourceOfFunds.encode(message.sender, writer.uint32(42).fork()).ldelim();
    }

    if (message.recipient !== undefined) {
      SourceOfFunds.encode(message.recipient, writer.uint32(50).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TransferEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTransferEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.senderSubaccountId = IndexerSubaccountId.decode(reader, reader.uint32());
          break;

        case 2:
          message.recipientSubaccountId = IndexerSubaccountId.decode(reader, reader.uint32());
          break;

        case 3:
          message.assetId = reader.uint32();
          break;

        case 4:
          message.amount = (reader.uint64() as Long);
          break;

        case 5:
          message.sender = SourceOfFunds.decode(reader, reader.uint32());
          break;

        case 6:
          message.recipient = SourceOfFunds.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<TransferEventV1>): TransferEventV1 {
    const message = createBaseTransferEventV1();
    message.senderSubaccountId = object.senderSubaccountId !== undefined && object.senderSubaccountId !== null ? IndexerSubaccountId.fromPartial(object.senderSubaccountId) : undefined;
    message.recipientSubaccountId = object.recipientSubaccountId !== undefined && object.recipientSubaccountId !== null ? IndexerSubaccountId.fromPartial(object.recipientSubaccountId) : undefined;
    message.assetId = object.assetId ?? 0;
    message.amount = object.amount !== undefined && object.amount !== null ? Long.fromValue(object.amount) : Long.UZERO;
    message.sender = object.sender !== undefined && object.sender !== null ? SourceOfFunds.fromPartial(object.sender) : undefined;
    message.recipient = object.recipient !== undefined && object.recipient !== null ? SourceOfFunds.fromPartial(object.recipient) : undefined;
    return message;
  }

};

function createBaseOrderFillEventV1(): OrderFillEventV1 {
  return {
    makerOrder: undefined,
    order: undefined,
    liquidationOrder: undefined,
    fillAmount: Long.UZERO,
    makerFee: Long.ZERO,
    takerFee: Long.ZERO,
    totalFilledMaker: Long.UZERO,
    totalFilledTaker: Long.UZERO
  };
}

export const OrderFillEventV1 = {
  encode(message: OrderFillEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.makerOrder !== undefined) {
      IndexerOrder.encode(message.makerOrder, writer.uint32(10).fork()).ldelim();
    }

    if (message.order !== undefined) {
      IndexerOrder.encode(message.order, writer.uint32(18).fork()).ldelim();
    }

    if (message.liquidationOrder !== undefined) {
      LiquidationOrderV1.encode(message.liquidationOrder, writer.uint32(34).fork()).ldelim();
    }

    if (!message.fillAmount.isZero()) {
      writer.uint32(24).uint64(message.fillAmount);
    }

    if (!message.makerFee.isZero()) {
      writer.uint32(40).sint64(message.makerFee);
    }

    if (!message.takerFee.isZero()) {
      writer.uint32(48).sint64(message.takerFee);
    }

    if (!message.totalFilledMaker.isZero()) {
      writer.uint32(56).uint64(message.totalFilledMaker);
    }

    if (!message.totalFilledTaker.isZero()) {
      writer.uint32(64).uint64(message.totalFilledTaker);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OrderFillEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrderFillEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.makerOrder = IndexerOrder.decode(reader, reader.uint32());
          break;

        case 2:
          message.order = IndexerOrder.decode(reader, reader.uint32());
          break;

        case 4:
          message.liquidationOrder = LiquidationOrderV1.decode(reader, reader.uint32());
          break;

        case 3:
          message.fillAmount = (reader.uint64() as Long);
          break;

        case 5:
          message.makerFee = (reader.sint64() as Long);
          break;

        case 6:
          message.takerFee = (reader.sint64() as Long);
          break;

        case 7:
          message.totalFilledMaker = (reader.uint64() as Long);
          break;

        case 8:
          message.totalFilledTaker = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OrderFillEventV1>): OrderFillEventV1 {
    const message = createBaseOrderFillEventV1();
    message.makerOrder = object.makerOrder !== undefined && object.makerOrder !== null ? IndexerOrder.fromPartial(object.makerOrder) : undefined;
    message.order = object.order !== undefined && object.order !== null ? IndexerOrder.fromPartial(object.order) : undefined;
    message.liquidationOrder = object.liquidationOrder !== undefined && object.liquidationOrder !== null ? LiquidationOrderV1.fromPartial(object.liquidationOrder) : undefined;
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? Long.fromValue(object.fillAmount) : Long.UZERO;
    message.makerFee = object.makerFee !== undefined && object.makerFee !== null ? Long.fromValue(object.makerFee) : Long.ZERO;
    message.takerFee = object.takerFee !== undefined && object.takerFee !== null ? Long.fromValue(object.takerFee) : Long.ZERO;
    message.totalFilledMaker = object.totalFilledMaker !== undefined && object.totalFilledMaker !== null ? Long.fromValue(object.totalFilledMaker) : Long.UZERO;
    message.totalFilledTaker = object.totalFilledTaker !== undefined && object.totalFilledTaker !== null ? Long.fromValue(object.totalFilledTaker) : Long.UZERO;
    return message;
  }

};

function createBaseDeleveragingEventV1(): DeleveragingEventV1 {
  return {
    liquidated: undefined,
    offsetting: undefined,
    perpetualId: 0,
    fillAmount: Long.UZERO,
    price: Long.UZERO,
    isBuy: false,
    isFinalSettlement: false
  };
}

export const DeleveragingEventV1 = {
  encode(message: DeleveragingEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.liquidated !== undefined) {
      IndexerSubaccountId.encode(message.liquidated, writer.uint32(10).fork()).ldelim();
    }

    if (message.offsetting !== undefined) {
      IndexerSubaccountId.encode(message.offsetting, writer.uint32(18).fork()).ldelim();
    }

    if (message.perpetualId !== 0) {
      writer.uint32(24).uint32(message.perpetualId);
    }

    if (!message.fillAmount.isZero()) {
      writer.uint32(32).uint64(message.fillAmount);
    }

    if (!message.price.isZero()) {
      writer.uint32(40).uint64(message.price);
    }

    if (message.isBuy === true) {
      writer.uint32(48).bool(message.isBuy);
    }

    if (message.isFinalSettlement === true) {
      writer.uint32(56).bool(message.isFinalSettlement);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleveragingEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleveragingEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.liquidated = IndexerSubaccountId.decode(reader, reader.uint32());
          break;

        case 2:
          message.offsetting = IndexerSubaccountId.decode(reader, reader.uint32());
          break;

        case 3:
          message.perpetualId = reader.uint32();
          break;

        case 4:
          message.fillAmount = (reader.uint64() as Long);
          break;

        case 5:
          message.price = (reader.uint64() as Long);
          break;

        case 6:
          message.isBuy = reader.bool();
          break;

        case 7:
          message.isFinalSettlement = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<DeleveragingEventV1>): DeleveragingEventV1 {
    const message = createBaseDeleveragingEventV1();
    message.liquidated = object.liquidated !== undefined && object.liquidated !== null ? IndexerSubaccountId.fromPartial(object.liquidated) : undefined;
    message.offsetting = object.offsetting !== undefined && object.offsetting !== null ? IndexerSubaccountId.fromPartial(object.offsetting) : undefined;
    message.perpetualId = object.perpetualId ?? 0;
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? Long.fromValue(object.fillAmount) : Long.UZERO;
    message.price = object.price !== undefined && object.price !== null ? Long.fromValue(object.price) : Long.UZERO;
    message.isBuy = object.isBuy ?? false;
    message.isFinalSettlement = object.isFinalSettlement ?? false;
    return message;
  }

};

function createBaseLiquidationOrderV1(): LiquidationOrderV1 {
  return {
    liquidated: undefined,
    clobPairId: 0,
    perpetualId: 0,
    totalSize: Long.UZERO,
    isBuy: false,
    subticks: Long.UZERO
  };
}

export const LiquidationOrderV1 = {
  encode(message: LiquidationOrderV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.liquidated !== undefined) {
      IndexerSubaccountId.encode(message.liquidated, writer.uint32(10).fork()).ldelim();
    }

    if (message.clobPairId !== 0) {
      writer.uint32(16).uint32(message.clobPairId);
    }

    if (message.perpetualId !== 0) {
      writer.uint32(24).uint32(message.perpetualId);
    }

    if (!message.totalSize.isZero()) {
      writer.uint32(32).uint64(message.totalSize);
    }

    if (message.isBuy === true) {
      writer.uint32(40).bool(message.isBuy);
    }

    if (!message.subticks.isZero()) {
      writer.uint32(48).uint64(message.subticks);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LiquidationOrderV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLiquidationOrderV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.liquidated = IndexerSubaccountId.decode(reader, reader.uint32());
          break;

        case 2:
          message.clobPairId = reader.uint32();
          break;

        case 3:
          message.perpetualId = reader.uint32();
          break;

        case 4:
          message.totalSize = (reader.uint64() as Long);
          break;

        case 5:
          message.isBuy = reader.bool();
          break;

        case 6:
          message.subticks = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<LiquidationOrderV1>): LiquidationOrderV1 {
    const message = createBaseLiquidationOrderV1();
    message.liquidated = object.liquidated !== undefined && object.liquidated !== null ? IndexerSubaccountId.fromPartial(object.liquidated) : undefined;
    message.clobPairId = object.clobPairId ?? 0;
    message.perpetualId = object.perpetualId ?? 0;
    message.totalSize = object.totalSize !== undefined && object.totalSize !== null ? Long.fromValue(object.totalSize) : Long.UZERO;
    message.isBuy = object.isBuy ?? false;
    message.subticks = object.subticks !== undefined && object.subticks !== null ? Long.fromValue(object.subticks) : Long.UZERO;
    return message;
  }

};

function createBaseSubaccountUpdateEventV1(): SubaccountUpdateEventV1 {
  return {
    subaccountId: undefined,
    updatedPerpetualPositions: [],
    updatedAssetPositions: []
  };
}

export const SubaccountUpdateEventV1 = {
  encode(message: SubaccountUpdateEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.subaccountId !== undefined) {
      IndexerSubaccountId.encode(message.subaccountId, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.updatedPerpetualPositions) {
      IndexerPerpetualPosition.encode(v!, writer.uint32(26).fork()).ldelim();
    }

    for (const v of message.updatedAssetPositions) {
      IndexerAssetPosition.encode(v!, writer.uint32(34).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubaccountUpdateEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubaccountUpdateEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountId = IndexerSubaccountId.decode(reader, reader.uint32());
          break;

        case 3:
          message.updatedPerpetualPositions.push(IndexerPerpetualPosition.decode(reader, reader.uint32()));
          break;

        case 4:
          message.updatedAssetPositions.push(IndexerAssetPosition.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<SubaccountUpdateEventV1>): SubaccountUpdateEventV1 {
    const message = createBaseSubaccountUpdateEventV1();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? IndexerSubaccountId.fromPartial(object.subaccountId) : undefined;
    message.updatedPerpetualPositions = object.updatedPerpetualPositions?.map(e => IndexerPerpetualPosition.fromPartial(e)) || [];
    message.updatedAssetPositions = object.updatedAssetPositions?.map(e => IndexerAssetPosition.fromPartial(e)) || [];
    return message;
  }

};

function createBaseStatefulOrderEventV1(): StatefulOrderEventV1 {
  return {
    orderPlace: undefined,
    orderRemoval: undefined,
    conditionalOrderPlacement: undefined,
    conditionalOrderTriggered: undefined,
    longTermOrderPlacement: undefined
  };
}

export const StatefulOrderEventV1 = {
  encode(message: StatefulOrderEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orderPlace !== undefined) {
      StatefulOrderEventV1_StatefulOrderPlacementV1.encode(message.orderPlace, writer.uint32(10).fork()).ldelim();
    }

    if (message.orderRemoval !== undefined) {
      StatefulOrderEventV1_StatefulOrderRemovalV1.encode(message.orderRemoval, writer.uint32(34).fork()).ldelim();
    }

    if (message.conditionalOrderPlacement !== undefined) {
      StatefulOrderEventV1_ConditionalOrderPlacementV1.encode(message.conditionalOrderPlacement, writer.uint32(42).fork()).ldelim();
    }

    if (message.conditionalOrderTriggered !== undefined) {
      StatefulOrderEventV1_ConditionalOrderTriggeredV1.encode(message.conditionalOrderTriggered, writer.uint32(50).fork()).ldelim();
    }

    if (message.longTermOrderPlacement !== undefined) {
      StatefulOrderEventV1_LongTermOrderPlacementV1.encode(message.longTermOrderPlacement, writer.uint32(58).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StatefulOrderEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatefulOrderEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.orderPlace = StatefulOrderEventV1_StatefulOrderPlacementV1.decode(reader, reader.uint32());
          break;

        case 4:
          message.orderRemoval = StatefulOrderEventV1_StatefulOrderRemovalV1.decode(reader, reader.uint32());
          break;

        case 5:
          message.conditionalOrderPlacement = StatefulOrderEventV1_ConditionalOrderPlacementV1.decode(reader, reader.uint32());
          break;

        case 6:
          message.conditionalOrderTriggered = StatefulOrderEventV1_ConditionalOrderTriggeredV1.decode(reader, reader.uint32());
          break;

        case 7:
          message.longTermOrderPlacement = StatefulOrderEventV1_LongTermOrderPlacementV1.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StatefulOrderEventV1>): StatefulOrderEventV1 {
    const message = createBaseStatefulOrderEventV1();
    message.orderPlace = object.orderPlace !== undefined && object.orderPlace !== null ? StatefulOrderEventV1_StatefulOrderPlacementV1.fromPartial(object.orderPlace) : undefined;
    message.orderRemoval = object.orderRemoval !== undefined && object.orderRemoval !== null ? StatefulOrderEventV1_StatefulOrderRemovalV1.fromPartial(object.orderRemoval) : undefined;
    message.conditionalOrderPlacement = object.conditionalOrderPlacement !== undefined && object.conditionalOrderPlacement !== null ? StatefulOrderEventV1_ConditionalOrderPlacementV1.fromPartial(object.conditionalOrderPlacement) : undefined;
    message.conditionalOrderTriggered = object.conditionalOrderTriggered !== undefined && object.conditionalOrderTriggered !== null ? StatefulOrderEventV1_ConditionalOrderTriggeredV1.fromPartial(object.conditionalOrderTriggered) : undefined;
    message.longTermOrderPlacement = object.longTermOrderPlacement !== undefined && object.longTermOrderPlacement !== null ? StatefulOrderEventV1_LongTermOrderPlacementV1.fromPartial(object.longTermOrderPlacement) : undefined;
    return message;
  }

};

function createBaseStatefulOrderEventV1_StatefulOrderPlacementV1(): StatefulOrderEventV1_StatefulOrderPlacementV1 {
  return {
    order: undefined
  };
}

export const StatefulOrderEventV1_StatefulOrderPlacementV1 = {
  encode(message: StatefulOrderEventV1_StatefulOrderPlacementV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.order !== undefined) {
      IndexerOrder.encode(message.order, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StatefulOrderEventV1_StatefulOrderPlacementV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatefulOrderEventV1_StatefulOrderPlacementV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.order = IndexerOrder.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StatefulOrderEventV1_StatefulOrderPlacementV1>): StatefulOrderEventV1_StatefulOrderPlacementV1 {
    const message = createBaseStatefulOrderEventV1_StatefulOrderPlacementV1();
    message.order = object.order !== undefined && object.order !== null ? IndexerOrder.fromPartial(object.order) : undefined;
    return message;
  }

};

function createBaseStatefulOrderEventV1_StatefulOrderRemovalV1(): StatefulOrderEventV1_StatefulOrderRemovalV1 {
  return {
    removedOrderId: undefined,
    reason: 0
  };
}

export const StatefulOrderEventV1_StatefulOrderRemovalV1 = {
  encode(message: StatefulOrderEventV1_StatefulOrderRemovalV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.removedOrderId !== undefined) {
      IndexerOrderId.encode(message.removedOrderId, writer.uint32(10).fork()).ldelim();
    }

    if (message.reason !== 0) {
      writer.uint32(16).int32(message.reason);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StatefulOrderEventV1_StatefulOrderRemovalV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatefulOrderEventV1_StatefulOrderRemovalV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.removedOrderId = IndexerOrderId.decode(reader, reader.uint32());
          break;

        case 2:
          message.reason = (reader.int32() as any);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StatefulOrderEventV1_StatefulOrderRemovalV1>): StatefulOrderEventV1_StatefulOrderRemovalV1 {
    const message = createBaseStatefulOrderEventV1_StatefulOrderRemovalV1();
    message.removedOrderId = object.removedOrderId !== undefined && object.removedOrderId !== null ? IndexerOrderId.fromPartial(object.removedOrderId) : undefined;
    message.reason = object.reason ?? 0;
    return message;
  }

};

function createBaseStatefulOrderEventV1_ConditionalOrderPlacementV1(): StatefulOrderEventV1_ConditionalOrderPlacementV1 {
  return {
    order: undefined
  };
}

export const StatefulOrderEventV1_ConditionalOrderPlacementV1 = {
  encode(message: StatefulOrderEventV1_ConditionalOrderPlacementV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.order !== undefined) {
      IndexerOrder.encode(message.order, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StatefulOrderEventV1_ConditionalOrderPlacementV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatefulOrderEventV1_ConditionalOrderPlacementV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.order = IndexerOrder.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StatefulOrderEventV1_ConditionalOrderPlacementV1>): StatefulOrderEventV1_ConditionalOrderPlacementV1 {
    const message = createBaseStatefulOrderEventV1_ConditionalOrderPlacementV1();
    message.order = object.order !== undefined && object.order !== null ? IndexerOrder.fromPartial(object.order) : undefined;
    return message;
  }

};

function createBaseStatefulOrderEventV1_ConditionalOrderTriggeredV1(): StatefulOrderEventV1_ConditionalOrderTriggeredV1 {
  return {
    triggeredOrderId: undefined
  };
}

export const StatefulOrderEventV1_ConditionalOrderTriggeredV1 = {
  encode(message: StatefulOrderEventV1_ConditionalOrderTriggeredV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.triggeredOrderId !== undefined) {
      IndexerOrderId.encode(message.triggeredOrderId, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StatefulOrderEventV1_ConditionalOrderTriggeredV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatefulOrderEventV1_ConditionalOrderTriggeredV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.triggeredOrderId = IndexerOrderId.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StatefulOrderEventV1_ConditionalOrderTriggeredV1>): StatefulOrderEventV1_ConditionalOrderTriggeredV1 {
    const message = createBaseStatefulOrderEventV1_ConditionalOrderTriggeredV1();
    message.triggeredOrderId = object.triggeredOrderId !== undefined && object.triggeredOrderId !== null ? IndexerOrderId.fromPartial(object.triggeredOrderId) : undefined;
    return message;
  }

};

function createBaseStatefulOrderEventV1_LongTermOrderPlacementV1(): StatefulOrderEventV1_LongTermOrderPlacementV1 {
  return {
    order: undefined
  };
}

export const StatefulOrderEventV1_LongTermOrderPlacementV1 = {
  encode(message: StatefulOrderEventV1_LongTermOrderPlacementV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.order !== undefined) {
      IndexerOrder.encode(message.order, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StatefulOrderEventV1_LongTermOrderPlacementV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatefulOrderEventV1_LongTermOrderPlacementV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.order = IndexerOrder.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StatefulOrderEventV1_LongTermOrderPlacementV1>): StatefulOrderEventV1_LongTermOrderPlacementV1 {
    const message = createBaseStatefulOrderEventV1_LongTermOrderPlacementV1();
    message.order = object.order !== undefined && object.order !== null ? IndexerOrder.fromPartial(object.order) : undefined;
    return message;
  }

};

function createBaseAssetCreateEventV1(): AssetCreateEventV1 {
  return {
    id: 0,
    symbol: "",
    hasMarket: false,
    marketId: 0,
    atomicResolution: 0
  };
}

export const AssetCreateEventV1 = {
  encode(message: AssetCreateEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    if (message.symbol !== "") {
      writer.uint32(18).string(message.symbol);
    }

    if (message.hasMarket === true) {
      writer.uint32(24).bool(message.hasMarket);
    }

    if (message.marketId !== 0) {
      writer.uint32(32).uint32(message.marketId);
    }

    if (message.atomicResolution !== 0) {
      writer.uint32(40).sint32(message.atomicResolution);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AssetCreateEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAssetCreateEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        case 2:
          message.symbol = reader.string();
          break;

        case 3:
          message.hasMarket = reader.bool();
          break;

        case 4:
          message.marketId = reader.uint32();
          break;

        case 5:
          message.atomicResolution = reader.sint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AssetCreateEventV1>): AssetCreateEventV1 {
    const message = createBaseAssetCreateEventV1();
    message.id = object.id ?? 0;
    message.symbol = object.symbol ?? "";
    message.hasMarket = object.hasMarket ?? false;
    message.marketId = object.marketId ?? 0;
    message.atomicResolution = object.atomicResolution ?? 0;
    return message;
  }

};

function createBasePerpetualMarketCreateEventV1(): PerpetualMarketCreateEventV1 {
  return {
    id: 0,
    clobPairId: 0,
    ticker: "",
    marketId: 0,
    status: 0,
    quantumConversionExponent: 0,
    atomicResolution: 0,
    subticksPerTick: 0,
    stepBaseQuantums: Long.UZERO,
    liquidityTier: 0
  };
}

export const PerpetualMarketCreateEventV1 = {
  encode(message: PerpetualMarketCreateEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    if (message.clobPairId !== 0) {
      writer.uint32(16).uint32(message.clobPairId);
    }

    if (message.ticker !== "") {
      writer.uint32(26).string(message.ticker);
    }

    if (message.marketId !== 0) {
      writer.uint32(32).uint32(message.marketId);
    }

    if (message.status !== 0) {
      writer.uint32(40).int32(message.status);
    }

    if (message.quantumConversionExponent !== 0) {
      writer.uint32(48).sint32(message.quantumConversionExponent);
    }

    if (message.atomicResolution !== 0) {
      writer.uint32(56).sint32(message.atomicResolution);
    }

    if (message.subticksPerTick !== 0) {
      writer.uint32(64).uint32(message.subticksPerTick);
    }

    if (!message.stepBaseQuantums.isZero()) {
      writer.uint32(72).uint64(message.stepBaseQuantums);
    }

    if (message.liquidityTier !== 0) {
      writer.uint32(80).uint32(message.liquidityTier);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PerpetualMarketCreateEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePerpetualMarketCreateEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        case 2:
          message.clobPairId = reader.uint32();
          break;

        case 3:
          message.ticker = reader.string();
          break;

        case 4:
          message.marketId = reader.uint32();
          break;

        case 5:
          message.status = (reader.int32() as any);
          break;

        case 6:
          message.quantumConversionExponent = reader.sint32();
          break;

        case 7:
          message.atomicResolution = reader.sint32();
          break;

        case 8:
          message.subticksPerTick = reader.uint32();
          break;

        case 9:
          message.stepBaseQuantums = (reader.uint64() as Long);
          break;

        case 10:
          message.liquidityTier = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<PerpetualMarketCreateEventV1>): PerpetualMarketCreateEventV1 {
    const message = createBasePerpetualMarketCreateEventV1();
    message.id = object.id ?? 0;
    message.clobPairId = object.clobPairId ?? 0;
    message.ticker = object.ticker ?? "";
    message.marketId = object.marketId ?? 0;
    message.status = object.status ?? 0;
    message.quantumConversionExponent = object.quantumConversionExponent ?? 0;
    message.atomicResolution = object.atomicResolution ?? 0;
    message.subticksPerTick = object.subticksPerTick ?? 0;
    message.stepBaseQuantums = object.stepBaseQuantums !== undefined && object.stepBaseQuantums !== null ? Long.fromValue(object.stepBaseQuantums) : Long.UZERO;
    message.liquidityTier = object.liquidityTier ?? 0;
    return message;
  }

};

function createBaseLiquidityTierUpsertEventV1(): LiquidityTierUpsertEventV1 {
  return {
    id: 0,
    name: "",
    initialMarginPpm: 0,
    maintenanceFractionPpm: 0,
    basePositionNotional: Long.UZERO
  };
}

export const LiquidityTierUpsertEventV1 = {
  encode(message: LiquidityTierUpsertEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }

    if (message.initialMarginPpm !== 0) {
      writer.uint32(24).uint32(message.initialMarginPpm);
    }

    if (message.maintenanceFractionPpm !== 0) {
      writer.uint32(32).uint32(message.maintenanceFractionPpm);
    }

    if (!message.basePositionNotional.isZero()) {
      writer.uint32(40).uint64(message.basePositionNotional);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LiquidityTierUpsertEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLiquidityTierUpsertEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        case 2:
          message.name = reader.string();
          break;

        case 3:
          message.initialMarginPpm = reader.uint32();
          break;

        case 4:
          message.maintenanceFractionPpm = reader.uint32();
          break;

        case 5:
          message.basePositionNotional = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<LiquidityTierUpsertEventV1>): LiquidityTierUpsertEventV1 {
    const message = createBaseLiquidityTierUpsertEventV1();
    message.id = object.id ?? 0;
    message.name = object.name ?? "";
    message.initialMarginPpm = object.initialMarginPpm ?? 0;
    message.maintenanceFractionPpm = object.maintenanceFractionPpm ?? 0;
    message.basePositionNotional = object.basePositionNotional !== undefined && object.basePositionNotional !== null ? Long.fromValue(object.basePositionNotional) : Long.UZERO;
    return message;
  }

};

function createBaseUpdateClobPairEventV1(): UpdateClobPairEventV1 {
  return {
    clobPairId: 0,
    status: 0,
    quantumConversionExponent: 0,
    subticksPerTick: 0,
    stepBaseQuantums: Long.UZERO
  };
}

export const UpdateClobPairEventV1 = {
  encode(message: UpdateClobPairEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobPairId !== 0) {
      writer.uint32(8).uint32(message.clobPairId);
    }

    if (message.status !== 0) {
      writer.uint32(16).int32(message.status);
    }

    if (message.quantumConversionExponent !== 0) {
      writer.uint32(24).sint32(message.quantumConversionExponent);
    }

    if (message.subticksPerTick !== 0) {
      writer.uint32(32).uint32(message.subticksPerTick);
    }

    if (!message.stepBaseQuantums.isZero()) {
      writer.uint32(40).uint64(message.stepBaseQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateClobPairEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateClobPairEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPairId = reader.uint32();
          break;

        case 2:
          message.status = (reader.int32() as any);
          break;

        case 3:
          message.quantumConversionExponent = reader.sint32();
          break;

        case 4:
          message.subticksPerTick = reader.uint32();
          break;

        case 5:
          message.stepBaseQuantums = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<UpdateClobPairEventV1>): UpdateClobPairEventV1 {
    const message = createBaseUpdateClobPairEventV1();
    message.clobPairId = object.clobPairId ?? 0;
    message.status = object.status ?? 0;
    message.quantumConversionExponent = object.quantumConversionExponent ?? 0;
    message.subticksPerTick = object.subticksPerTick ?? 0;
    message.stepBaseQuantums = object.stepBaseQuantums !== undefined && object.stepBaseQuantums !== null ? Long.fromValue(object.stepBaseQuantums) : Long.UZERO;
    return message;
  }

};

function createBaseUpdatePerpetualEventV1(): UpdatePerpetualEventV1 {
  return {
    id: 0,
    ticker: "",
    marketId: 0,
    atomicResolution: 0,
    liquidityTier: 0
  };
}

export const UpdatePerpetualEventV1 = {
  encode(message: UpdatePerpetualEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    if (message.ticker !== "") {
      writer.uint32(18).string(message.ticker);
    }

    if (message.marketId !== 0) {
      writer.uint32(24).uint32(message.marketId);
    }

    if (message.atomicResolution !== 0) {
      writer.uint32(32).sint32(message.atomicResolution);
    }

    if (message.liquidityTier !== 0) {
      writer.uint32(40).uint32(message.liquidityTier);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdatePerpetualEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdatePerpetualEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        case 2:
          message.ticker = reader.string();
          break;

        case 3:
          message.marketId = reader.uint32();
          break;

        case 4:
          message.atomicResolution = reader.sint32();
          break;

        case 5:
          message.liquidityTier = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<UpdatePerpetualEventV1>): UpdatePerpetualEventV1 {
    const message = createBaseUpdatePerpetualEventV1();
    message.id = object.id ?? 0;
    message.ticker = object.ticker ?? "";
    message.marketId = object.marketId ?? 0;
    message.atomicResolution = object.atomicResolution ?? 0;
    message.liquidityTier = object.liquidityTier ?? 0;
    return message;
  }

};

function createBaseTradingRewardEventV1(): TradingRewardEventV1 {
  return {
    tradingRewards: []
  };
}

export const TradingRewardEventV1 = {
  encode(message: TradingRewardEventV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.tradingRewards) {
      AddressTradingReward.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TradingRewardEventV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTradingRewardEventV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.tradingRewards.push(AddressTradingReward.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<TradingRewardEventV1>): TradingRewardEventV1 {
    const message = createBaseTradingRewardEventV1();
    message.tradingRewards = object.tradingRewards?.map(e => AddressTradingReward.fromPartial(e)) || [];
    return message;
  }

};

function createBaseAddressTradingReward(): AddressTradingReward {
  return {
    owner: "",
    denoms: new Uint8Array()
  };
}

export const AddressTradingReward = {
  encode(message: AddressTradingReward, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.owner !== "") {
      writer.uint32(10).string(message.owner);
    }

    if (message.denoms.length !== 0) {
      writer.uint32(18).bytes(message.denoms);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddressTradingReward {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddressTradingReward();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.owner = reader.string();
          break;

        case 2:
          message.denoms = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AddressTradingReward>): AddressTradingReward {
    const message = createBaseAddressTradingReward();
    message.owner = object.owner ?? "";
    message.denoms = object.denoms ?? new Uint8Array();
    return message;
  }

};