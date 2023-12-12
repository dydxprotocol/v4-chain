import { IndexerSubaccountId, IndexerSubaccountIdAmino, IndexerSubaccountIdSDKType, IndexerPerpetualPosition, IndexerPerpetualPositionAmino, IndexerPerpetualPositionSDKType, IndexerAssetPosition, IndexerAssetPositionAmino, IndexerAssetPositionSDKType } from "../protocol/v1/subaccount";
import { IndexerOrder, IndexerOrderAmino, IndexerOrderSDKType, IndexerOrderId, IndexerOrderIdAmino, IndexerOrderIdSDKType, ClobPairStatus, clobPairStatusFromJSON, clobPairStatusToJSON } from "../protocol/v1/clob";
import { OrderRemovalReason, orderRemovalReasonFromJSON, orderRemovalReasonToJSON } from "../shared/removal_reason";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { bytesFromBase64, base64FromBytes } from "../../../helpers";
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
export const FundingEventV1_TypeSDKType = FundingEventV1_Type;
export const FundingEventV1_TypeAmino = FundingEventV1_Type;
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
export interface FundingUpdateV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.FundingUpdateV1";
  value: Uint8Array;
}
/**
 * FundingUpdate is used for funding update events and includes a funding
 * value and an optional funding index that correspond to a perpetual market.
 */
export interface FundingUpdateV1Amino {
  /** The id of the perpetual market. */
  perpetual_id?: number;
  /**
   * funding value (in parts-per-million) can be premium vote, premium sample,
   * or funding rate.
   */
  funding_value_ppm?: number;
  /**
   * funding index is required if and only if parent `FundingEvent` type is
   * `TYPE_FUNDING_RATE_AND_INDEX`.
   */
  funding_index?: string;
}
export interface FundingUpdateV1AminoMsg {
  type: "/dydxprotocol.indexer.events.FundingUpdateV1";
  value: FundingUpdateV1Amino;
}
/**
 * FundingUpdate is used for funding update events and includes a funding
 * value and an optional funding index that correspond to a perpetual market.
 */
export interface FundingUpdateV1SDKType {
  perpetual_id: number;
  funding_value_ppm: number;
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
export interface FundingEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.FundingEventV1";
  value: Uint8Array;
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
export interface FundingEventV1Amino {
  /**
   * updates is a list of per-market funding updates for all existing perpetual
   * markets. The list is sorted by `perpetualId`s which are unique.
   */
  updates?: FundingUpdateV1Amino[];
  /** type stores the type of funding updates. */
  type?: FundingEventV1_Type;
}
export interface FundingEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.FundingEventV1";
  value: FundingEventV1Amino;
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
  updates: FundingUpdateV1SDKType[];
  type: FundingEventV1_Type;
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
export interface MarketEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.MarketEventV1";
  value: Uint8Array;
}
/**
 * MarketEvent message contains all the information about a market event on
 * the dYdX chain.
 */
export interface MarketEventV1Amino {
  /** market id. */
  market_id?: number;
  price_update?: MarketPriceUpdateEventV1Amino;
  market_create?: MarketCreateEventV1Amino;
  market_modify?: MarketModifyEventV1Amino;
}
export interface MarketEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.MarketEventV1";
  value: MarketEventV1Amino;
}
/**
 * MarketEvent message contains all the information about a market event on
 * the dYdX chain.
 */
export interface MarketEventV1SDKType {
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
  priceWithExponent: bigint;
}
export interface MarketPriceUpdateEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.MarketPriceUpdateEventV1";
  value: Uint8Array;
}
/**
 * MarketPriceUpdateEvent message contains all the information about a price
 * update on the dYdX chain.
 */
export interface MarketPriceUpdateEventV1Amino {
  /**
   * price_with_exponent. Multiply by 10 ^ Exponent to get the human readable
   * price in dollars. For example if `Exponent == -5` then a `exponent_price`
   * of `1,000,000,000` represents “$10,000`.
   */
  price_with_exponent?: string;
}
export interface MarketPriceUpdateEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.MarketPriceUpdateEventV1";
  value: MarketPriceUpdateEventV1Amino;
}
/**
 * MarketPriceUpdateEvent message contains all the information about a price
 * update on the dYdX chain.
 */
export interface MarketPriceUpdateEventV1SDKType {
  price_with_exponent: bigint;
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
export interface MarketBaseEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.MarketBaseEventV1";
  value: Uint8Array;
}
/** shared fields between MarketCreateEvent and MarketModifyEvent */
export interface MarketBaseEventV1Amino {
  /** String representation of the market pair, e.g. `BTC-USD` */
  pair?: string;
  /**
   * The minimum allowable change in the Price value for a given update.
   * Measured as 1e-6.
   */
  min_price_change_ppm?: number;
}
export interface MarketBaseEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.MarketBaseEventV1";
  value: MarketBaseEventV1Amino;
}
/** shared fields between MarketCreateEvent and MarketModifyEvent */
export interface MarketBaseEventV1SDKType {
  pair: string;
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
export interface MarketCreateEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.MarketCreateEventV1";
  value: Uint8Array;
}
/**
 * MarketCreateEvent message contains all the information about a new market on
 * the dYdX chain.
 */
export interface MarketCreateEventV1Amino {
  base?: MarketBaseEventV1Amino;
  /**
   * Static value. The exponent of the price.
   * For example if Exponent == -5 then a `exponent_price` of 1,000,000,000
   * represents $10,000. Therefore 10 ^ Exponent represents the smallest
   * price step (in dollars) that can be recorded.
   */
  exponent?: number;
}
export interface MarketCreateEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.MarketCreateEventV1";
  value: MarketCreateEventV1Amino;
}
/**
 * MarketCreateEvent message contains all the information about a new market on
 * the dYdX chain.
 */
export interface MarketCreateEventV1SDKType {
  base?: MarketBaseEventV1SDKType;
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
export interface MarketModifyEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.MarketModifyEventV1";
  value: Uint8Array;
}
/**
 * MarketModifyEvent message contains all the information about a market update
 * on the dYdX chain
 */
export interface MarketModifyEventV1Amino {
  /**
   * MarketModifyEvent message contains all the information about a market update
   * on the dYdX chain
   */
  base?: MarketBaseEventV1Amino;
}
export interface MarketModifyEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.MarketModifyEventV1";
  value: MarketModifyEventV1Amino;
}
/**
 * MarketModifyEvent message contains all the information about a market update
 * on the dYdX chain
 */
export interface MarketModifyEventV1SDKType {
  base?: MarketBaseEventV1SDKType;
}
/** SourceOfFunds is the source of funds in a transfer event. */
export interface SourceOfFunds {
  subaccountId?: IndexerSubaccountId;
  address?: string;
}
export interface SourceOfFundsProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.SourceOfFunds";
  value: Uint8Array;
}
/** SourceOfFunds is the source of funds in a transfer event. */
export interface SourceOfFundsAmino {
  subaccount_id?: IndexerSubaccountIdAmino;
  address?: string;
}
export interface SourceOfFundsAminoMsg {
  type: "/dydxprotocol.indexer.events.SourceOfFunds";
  value: SourceOfFundsAmino;
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
  amount: bigint;
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
export interface TransferEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.TransferEventV1";
  value: Uint8Array;
}
/**
 * TransferEvent message contains all the information about a transfer,
 * deposit-to-subaccount, or withdraw-from-subaccount on the dYdX chain.
 * When a subaccount is involved, a SubaccountUpdateEvent message will
 * be produced with the updated asset positions.
 */
export interface TransferEventV1Amino {
  sender_subaccount_id?: IndexerSubaccountIdAmino;
  recipient_subaccount_id?: IndexerSubaccountIdAmino;
  /** Id of the asset transfered. */
  asset_id?: number;
  /** The amount of asset in quantums to transfer. */
  amount?: string;
  /**
   * The sender is one of below
   * - a subaccount ID (in transfer and withdraw events).
   * - a wallet address (in deposit events).
   */
  sender?: SourceOfFundsAmino;
  /**
   * The recipient is one of below
   * - a subaccount ID (in transfer and deposit events).
   * - a wallet address (in withdraw events).
   */
  recipient?: SourceOfFundsAmino;
}
export interface TransferEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.TransferEventV1";
  value: TransferEventV1Amino;
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
  asset_id: number;
  amount: bigint;
  sender?: SourceOfFundsSDKType;
  recipient?: SourceOfFundsSDKType;
}
/**
 * OrderFillEvent message contains all the information from an order match in
 * the dYdX chain. This includes the maker/taker orders that matched and the
 * amount filled.
 */
export interface OrderFillEventV1 {
  makerOrder: IndexerOrder;
  order?: IndexerOrder;
  liquidationOrder?: LiquidationOrderV1;
  /** Fill amount in base quantums. */
  fillAmount: bigint;
  /** Maker fee in USDC quantums. */
  makerFee: bigint;
  /**
   * Taker fee in USDC quantums. If the taker order is a liquidation, then this
   * represents the special liquidation fee, not the standard taker fee.
   */
  takerFee: bigint;
  /** Total filled of the maker order in base quantums. */
  totalFilledMaker: bigint;
  /** Total filled of the taker order in base quantums. */
  totalFilledTaker: bigint;
}
export interface OrderFillEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.OrderFillEventV1";
  value: Uint8Array;
}
/**
 * OrderFillEvent message contains all the information from an order match in
 * the dYdX chain. This includes the maker/taker orders that matched and the
 * amount filled.
 */
export interface OrderFillEventV1Amino {
  maker_order?: IndexerOrderAmino;
  order?: IndexerOrderAmino;
  liquidation_order?: LiquidationOrderV1Amino;
  /** Fill amount in base quantums. */
  fill_amount?: string;
  /** Maker fee in USDC quantums. */
  maker_fee?: string;
  /**
   * Taker fee in USDC quantums. If the taker order is a liquidation, then this
   * represents the special liquidation fee, not the standard taker fee.
   */
  taker_fee?: string;
  /** Total filled of the maker order in base quantums. */
  total_filled_maker?: string;
  /** Total filled of the taker order in base quantums. */
  total_filled_taker?: string;
}
export interface OrderFillEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.OrderFillEventV1";
  value: OrderFillEventV1Amino;
}
/**
 * OrderFillEvent message contains all the information from an order match in
 * the dYdX chain. This includes the maker/taker orders that matched and the
 * amount filled.
 */
export interface OrderFillEventV1SDKType {
  maker_order: IndexerOrderSDKType;
  order?: IndexerOrderSDKType;
  liquidation_order?: LiquidationOrderV1SDKType;
  fill_amount: bigint;
  maker_fee: bigint;
  taker_fee: bigint;
  total_filled_maker: bigint;
  total_filled_taker: bigint;
}
/**
 * DeleveragingEvent message contains all the information for a deleveraging
 * on the dYdX chain. This includes the liquidated/offsetting subaccounts and
 * the amount filled.
 */
export interface DeleveragingEventV1 {
  /** ID of the subaccount that was liquidated. */
  liquidated: IndexerSubaccountId;
  /** ID of the subaccount that was used to offset the position. */
  offsetting: IndexerSubaccountId;
  /** The ID of the perpetual that was liquidated. */
  perpetualId: number;
  /**
   * The amount filled between the liquidated and offsetting position, in
   * base quantums.
   */
  fillAmount: bigint;
  /** Fill price of deleveraging event, in USDC quote quantums. */
  price: bigint;
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
export interface DeleveragingEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.DeleveragingEventV1";
  value: Uint8Array;
}
/**
 * DeleveragingEvent message contains all the information for a deleveraging
 * on the dYdX chain. This includes the liquidated/offsetting subaccounts and
 * the amount filled.
 */
export interface DeleveragingEventV1Amino {
  /** ID of the subaccount that was liquidated. */
  liquidated?: IndexerSubaccountIdAmino;
  /** ID of the subaccount that was used to offset the position. */
  offsetting?: IndexerSubaccountIdAmino;
  /** The ID of the perpetual that was liquidated. */
  perpetual_id?: number;
  /**
   * The amount filled between the liquidated and offsetting position, in
   * base quantums.
   */
  fill_amount?: string;
  /** Fill price of deleveraging event, in USDC quote quantums. */
  price?: string;
  /** `true` if liquidating a short position, `false` otherwise. */
  is_buy?: boolean;
  /**
   * `true` if the deleveraging event is for final settlement, indicating
   * the match occurred at the oracle price rather than bankruptcy price.
   * When this flag is `false`, the fill price is the bankruptcy price
   * of the liquidated subaccount.
   */
  is_final_settlement?: boolean;
}
export interface DeleveragingEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.DeleveragingEventV1";
  value: DeleveragingEventV1Amino;
}
/**
 * DeleveragingEvent message contains all the information for a deleveraging
 * on the dYdX chain. This includes the liquidated/offsetting subaccounts and
 * the amount filled.
 */
export interface DeleveragingEventV1SDKType {
  liquidated: IndexerSubaccountIdSDKType;
  offsetting: IndexerSubaccountIdSDKType;
  perpetual_id: number;
  fill_amount: bigint;
  price: bigint;
  is_buy: boolean;
  is_final_settlement: boolean;
}
/**
 * LiquidationOrder represents the liquidation taker order to be included in a
 * liquidation order fill event.
 */
export interface LiquidationOrderV1 {
  /** ID of the subaccount that was liquidated. */
  liquidated: IndexerSubaccountId;
  /** The ID of the clob pair involved in the liquidation. */
  clobPairId: number;
  /** The ID of the perpetual involved in the liquidation. */
  perpetualId: number;
  /**
   * The total size of the liquidation order including any unfilled size,
   * in base quantums.
   */
  totalSize: bigint;
  /** `true` if liquidating a short position, `false` otherwise. */
  isBuy: boolean;
  /**
   * The fillable price in subticks.
   * This represents the lower-price-bound for liquidating longs
   * and the upper-price-bound for liquidating shorts.
   * Must be a multiple of ClobPair.SubticksPerTick
   * (where `ClobPair.Id = orderId.ClobPairId`).
   */
  subticks: bigint;
}
export interface LiquidationOrderV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.LiquidationOrderV1";
  value: Uint8Array;
}
/**
 * LiquidationOrder represents the liquidation taker order to be included in a
 * liquidation order fill event.
 */
export interface LiquidationOrderV1Amino {
  /** ID of the subaccount that was liquidated. */
  liquidated?: IndexerSubaccountIdAmino;
  /** The ID of the clob pair involved in the liquidation. */
  clob_pair_id?: number;
  /** The ID of the perpetual involved in the liquidation. */
  perpetual_id?: number;
  /**
   * The total size of the liquidation order including any unfilled size,
   * in base quantums.
   */
  total_size?: string;
  /** `true` if liquidating a short position, `false` otherwise. */
  is_buy?: boolean;
  /**
   * The fillable price in subticks.
   * This represents the lower-price-bound for liquidating longs
   * and the upper-price-bound for liquidating shorts.
   * Must be a multiple of ClobPair.SubticksPerTick
   * (where `ClobPair.Id = orderId.ClobPairId`).
   */
  subticks?: string;
}
export interface LiquidationOrderV1AminoMsg {
  type: "/dydxprotocol.indexer.events.LiquidationOrderV1";
  value: LiquidationOrderV1Amino;
}
/**
 * LiquidationOrder represents the liquidation taker order to be included in a
 * liquidation order fill event.
 */
export interface LiquidationOrderV1SDKType {
  liquidated: IndexerSubaccountIdSDKType;
  clob_pair_id: number;
  perpetual_id: number;
  total_size: bigint;
  is_buy: boolean;
  subticks: bigint;
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
export interface SubaccountUpdateEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.SubaccountUpdateEventV1";
  value: Uint8Array;
}
/**
 * SubaccountUpdateEvent message contains information about an update to a
 * subaccount in the dYdX chain. This includes the list of updated perpetual
 * and asset positions for the subaccount.
 * Note: This event message will contain all the updates to a subaccount
 * at the end of a block which is why multiple asset/perpetual position
 * updates may exist.
 */
export interface SubaccountUpdateEventV1Amino {
  subaccount_id?: IndexerSubaccountIdAmino;
  updated_perpetual_positions?: IndexerPerpetualPositionAmino[];
  updated_asset_positions?: IndexerAssetPositionAmino[];
}
export interface SubaccountUpdateEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.SubaccountUpdateEventV1";
  value: SubaccountUpdateEventV1Amino;
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
export interface StatefulOrderEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.StatefulOrderEventV1";
  value: Uint8Array;
}
/**
 * StatefulOrderEvent message contains information about a change to a stateful
 * order. Currently, this is either the placement of a long-term order, the
 * placement or triggering of a conditional order, or the removal of a
 * stateful order.
 */
export interface StatefulOrderEventV1Amino {
  order_place?: StatefulOrderEventV1_StatefulOrderPlacementV1Amino;
  order_removal?: StatefulOrderEventV1_StatefulOrderRemovalV1Amino;
  conditional_order_placement?: StatefulOrderEventV1_ConditionalOrderPlacementV1Amino;
  conditional_order_triggered?: StatefulOrderEventV1_ConditionalOrderTriggeredV1Amino;
  long_term_order_placement?: StatefulOrderEventV1_LongTermOrderPlacementV1Amino;
}
export interface StatefulOrderEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.StatefulOrderEventV1";
  value: StatefulOrderEventV1Amino;
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
export interface StatefulOrderEventV1_StatefulOrderPlacementV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.StatefulOrderPlacementV1";
  value: Uint8Array;
}
/** A stateful order placement contains an order. */
export interface StatefulOrderEventV1_StatefulOrderPlacementV1Amino {
  order?: IndexerOrderAmino;
}
export interface StatefulOrderEventV1_StatefulOrderPlacementV1AminoMsg {
  type: "/dydxprotocol.indexer.events.StatefulOrderPlacementV1";
  value: StatefulOrderEventV1_StatefulOrderPlacementV1Amino;
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
export interface StatefulOrderEventV1_StatefulOrderRemovalV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.StatefulOrderRemovalV1";
  value: Uint8Array;
}
/**
 * A stateful order removal contains the id of an order that was already
 * placed and is now removed and the reason for the removal.
 */
export interface StatefulOrderEventV1_StatefulOrderRemovalV1Amino {
  removed_order_id?: IndexerOrderIdAmino;
  reason?: OrderRemovalReason;
}
export interface StatefulOrderEventV1_StatefulOrderRemovalV1AminoMsg {
  type: "/dydxprotocol.indexer.events.StatefulOrderRemovalV1";
  value: StatefulOrderEventV1_StatefulOrderRemovalV1Amino;
}
/**
 * A stateful order removal contains the id of an order that was already
 * placed and is now removed and the reason for the removal.
 */
export interface StatefulOrderEventV1_StatefulOrderRemovalV1SDKType {
  removed_order_id?: IndexerOrderIdSDKType;
  reason: OrderRemovalReason;
}
/**
 * A conditional order placement contains an order. The order is newly-placed
 * and untriggered when this event is emitted.
 */
export interface StatefulOrderEventV1_ConditionalOrderPlacementV1 {
  order?: IndexerOrder;
}
export interface StatefulOrderEventV1_ConditionalOrderPlacementV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.ConditionalOrderPlacementV1";
  value: Uint8Array;
}
/**
 * A conditional order placement contains an order. The order is newly-placed
 * and untriggered when this event is emitted.
 */
export interface StatefulOrderEventV1_ConditionalOrderPlacementV1Amino {
  order?: IndexerOrderAmino;
}
export interface StatefulOrderEventV1_ConditionalOrderPlacementV1AminoMsg {
  type: "/dydxprotocol.indexer.events.ConditionalOrderPlacementV1";
  value: StatefulOrderEventV1_ConditionalOrderPlacementV1Amino;
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
export interface StatefulOrderEventV1_ConditionalOrderTriggeredV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.ConditionalOrderTriggeredV1";
  value: Uint8Array;
}
/**
 * A conditional order trigger event contains an order id and is emitted when
 * an order is triggered.
 */
export interface StatefulOrderEventV1_ConditionalOrderTriggeredV1Amino {
  triggered_order_id?: IndexerOrderIdAmino;
}
export interface StatefulOrderEventV1_ConditionalOrderTriggeredV1AminoMsg {
  type: "/dydxprotocol.indexer.events.ConditionalOrderTriggeredV1";
  value: StatefulOrderEventV1_ConditionalOrderTriggeredV1Amino;
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
export interface StatefulOrderEventV1_LongTermOrderPlacementV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.LongTermOrderPlacementV1";
  value: Uint8Array;
}
/** A long term order placement contains an order. */
export interface StatefulOrderEventV1_LongTermOrderPlacementV1Amino {
  order?: IndexerOrderAmino;
}
export interface StatefulOrderEventV1_LongTermOrderPlacementV1AminoMsg {
  type: "/dydxprotocol.indexer.events.LongTermOrderPlacementV1";
  value: StatefulOrderEventV1_LongTermOrderPlacementV1Amino;
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
export interface AssetCreateEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.AssetCreateEventV1";
  value: Uint8Array;
}
/**
 * AssetCreateEventV1 message contains all the information about an new Asset on
 * the dYdX chain.
 */
export interface AssetCreateEventV1Amino {
  /** Unique, sequentially-generated. */
  id?: number;
  /**
   * The human readable symbol of the `Asset` (e.g. `USDC`, `ATOM`).
   * Must be uppercase, unique and correspond to the canonical symbol of the
   * full coin.
   */
  symbol?: string;
  /** `true` if this `Asset` has a valid `MarketId` value. */
  has_market?: boolean;
  /**
   * The `Id` of the `Market` associated with this `Asset`. It acts as the
   * oracle price for the purposes of calculating collateral
   * and margin requirements.
   */
  market_id?: number;
  /**
   * The exponent for converting an atomic amount (1 'quantum')
   * to a full coin. For example, if `atomic_resolution = -8`
   * then an `asset_position` with `base_quantums = 1e8` is equivalent to
   * a position size of one full coin.
   */
  atomic_resolution?: number;
}
export interface AssetCreateEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.AssetCreateEventV1";
  value: AssetCreateEventV1Amino;
}
/**
 * AssetCreateEventV1 message contains all the information about an new Asset on
 * the dYdX chain.
 */
export interface AssetCreateEventV1SDKType {
  id: number;
  symbol: string;
  has_market: boolean;
  market_id: number;
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
  stepBaseQuantums: bigint;
  /**
   * The liquidity_tier that this perpetual is associated with.
   * Defined in perpetuals.perpetual
   */
  liquidityTier: number;
}
export interface PerpetualMarketCreateEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.PerpetualMarketCreateEventV1";
  value: Uint8Array;
}
/**
 * PerpetualMarketCreateEventV1 message contains all the information about a
 * new Perpetual Market on the dYdX chain.
 */
export interface PerpetualMarketCreateEventV1Amino {
  /**
   * Unique Perpetual id.
   * Defined in perpetuals.perpetual
   */
  id?: number;
  /**
   * Unique clob pair Id associated with this perpetual market
   * Defined in clob.clob_pair
   */
  clob_pair_id?: number;
  /**
   * The name of the `Perpetual` (e.g. `BTC-USD`).
   * Defined in perpetuals.perpetual
   */
  ticker?: string;
  /**
   * Unique id of market param associated with this perpetual market.
   * Defined in perpetuals.perpetual
   */
  market_id?: number;
  /** Status of the CLOB */
  status?: ClobPairStatus;
  /**
   * `10^Exponent` gives the number of QuoteQuantums traded per BaseQuantum
   * per Subtick.
   * Defined in clob.clob_pair
   */
  quantum_conversion_exponent?: number;
  /**
   * The exponent for converting an atomic amount (`size = 1`)
   * to a full coin. For example, if `AtomicResolution = -8`
   * then a `PerpetualPosition` with `size = 1e8` is equivalent to
   * a position size of one full coin.
   * Defined in perpetuals.perpetual
   */
  atomic_resolution?: number;
  /**
   * Defines the tick size of the orderbook by defining how many subticks
   * are in one tick. That is, the subticks of any valid order must be a
   * multiple of this value. Generally this value should start `>= 100`to
   * allow room for decreasing it.
   * Defined in clob.clob_pair
   */
  subticks_per_tick?: number;
  /**
   * Minimum increment in the size of orders on the CLOB, in base quantums.
   * Defined in clob.clob_pair
   */
  step_base_quantums?: string;
  /**
   * The liquidity_tier that this perpetual is associated with.
   * Defined in perpetuals.perpetual
   */
  liquidity_tier?: number;
}
export interface PerpetualMarketCreateEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.PerpetualMarketCreateEventV1";
  value: PerpetualMarketCreateEventV1Amino;
}
/**
 * PerpetualMarketCreateEventV1 message contains all the information about a
 * new Perpetual Market on the dYdX chain.
 */
export interface PerpetualMarketCreateEventV1SDKType {
  id: number;
  clob_pair_id: number;
  ticker: string;
  market_id: number;
  status: ClobPairStatus;
  quantum_conversion_exponent: number;
  atomic_resolution: number;
  subticks_per_tick: number;
  step_base_quantums: bigint;
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
  basePositionNotional: bigint;
}
export interface LiquidityTierUpsertEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.LiquidityTierUpsertEventV1";
  value: Uint8Array;
}
/**
 * LiquidityTierUpsertEventV1 message contains all the information to
 * create/update a Liquidity Tier on the dYdX chain.
 */
export interface LiquidityTierUpsertEventV1Amino {
  /** Unique id. */
  id?: number;
  /** The name of the tier purely for mnemonic purposes, e.g. "Gold". */
  name?: string;
  /**
   * The margin fraction needed to open a position.
   * In parts-per-million.
   */
  initial_margin_ppm?: number;
  /**
   * The fraction of the initial-margin that the maintenance-margin is,
   * e.g. 50%. In parts-per-million.
   */
  maintenance_fraction_ppm?: number;
  /**
   * The maximum position size at which the margin requirements are
   * not increased over the default values. Above this position size,
   * the margin requirements increase at a rate of sqrt(size).
   */
  base_position_notional?: string;
}
export interface LiquidityTierUpsertEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.LiquidityTierUpsertEventV1";
  value: LiquidityTierUpsertEventV1Amino;
}
/**
 * LiquidityTierUpsertEventV1 message contains all the information to
 * create/update a Liquidity Tier on the dYdX chain.
 */
export interface LiquidityTierUpsertEventV1SDKType {
  id: number;
  name: string;
  initial_margin_ppm: number;
  maintenance_fraction_ppm: number;
  base_position_notional: bigint;
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
  stepBaseQuantums: bigint;
}
export interface UpdateClobPairEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.UpdateClobPairEventV1";
  value: Uint8Array;
}
/**
 * UpdateClobPairEventV1 message contains all the information about an update to
 * a clob pair on the dYdX chain.
 */
export interface UpdateClobPairEventV1Amino {
  /**
   * Unique clob pair Id associated with this perpetual market
   * Defined in clob.clob_pair
   */
  clob_pair_id?: number;
  /** Status of the CLOB */
  status?: ClobPairStatus;
  /**
   * `10^Exponent` gives the number of QuoteQuantums traded per BaseQuantum
   * per Subtick.
   * Defined in clob.clob_pair
   */
  quantum_conversion_exponent?: number;
  /**
   * Defines the tick size of the orderbook by defining how many subticks
   * are in one tick. That is, the subticks of any valid order must be a
   * multiple of this value. Generally this value should start `>= 100`to
   * allow room for decreasing it.
   * Defined in clob.clob_pair
   */
  subticks_per_tick?: number;
  /**
   * Minimum increment in the size of orders on the CLOB, in base quantums.
   * Defined in clob.clob_pair
   */
  step_base_quantums?: string;
}
export interface UpdateClobPairEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.UpdateClobPairEventV1";
  value: UpdateClobPairEventV1Amino;
}
/**
 * UpdateClobPairEventV1 message contains all the information about an update to
 * a clob pair on the dYdX chain.
 */
export interface UpdateClobPairEventV1SDKType {
  clob_pair_id: number;
  status: ClobPairStatus;
  quantum_conversion_exponent: number;
  subticks_per_tick: number;
  step_base_quantums: bigint;
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
export interface UpdatePerpetualEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.UpdatePerpetualEventV1";
  value: Uint8Array;
}
/**
 * UpdatePerpetualEventV1 message contains all the information about an update
 * to a perpetual on the dYdX chain.
 */
export interface UpdatePerpetualEventV1Amino {
  /**
   * Unique Perpetual id.
   * Defined in perpetuals.perpetual
   */
  id?: number;
  /**
   * The name of the `Perpetual` (e.g. `BTC-USD`).
   * Defined in perpetuals.perpetual
   */
  ticker?: string;
  /**
   * Unique id of market param associated with this perpetual market.
   * Defined in perpetuals.perpetual
   */
  market_id?: number;
  /**
   * The exponent for converting an atomic amount (`size = 1`)
   * to a full coin. For example, if `AtomicResolution = -8`
   * then a `PerpetualPosition` with `size = 1e8` is equivalent to
   * a position size of one full coin.
   * Defined in perpetuals.perpetual
   */
  atomic_resolution?: number;
  /**
   * The liquidity_tier that this perpetual is associated with.
   * Defined in perpetuals.perpetual
   */
  liquidity_tier?: number;
}
export interface UpdatePerpetualEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.UpdatePerpetualEventV1";
  value: UpdatePerpetualEventV1Amino;
}
/**
 * UpdatePerpetualEventV1 message contains all the information about an update
 * to a perpetual on the dYdX chain.
 */
export interface UpdatePerpetualEventV1SDKType {
  id: number;
  ticker: string;
  market_id: number;
  atomic_resolution: number;
  liquidity_tier: number;
}
/**
 * TradingRewardsEventV1 is communicates all trading rewards for all accounts
 * that receive trade rewards in the block.
 */
export interface TradingRewardsEventV1 {
  /** The list of all trading rewards in the block. */
  tradingRewards: AddressTradingReward[];
}
export interface TradingRewardsEventV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.TradingRewardsEventV1";
  value: Uint8Array;
}
/**
 * TradingRewardsEventV1 is communicates all trading rewards for all accounts
 * that receive trade rewards in the block.
 */
export interface TradingRewardsEventV1Amino {
  /** The list of all trading rewards in the block. */
  trading_rewards?: AddressTradingRewardAmino[];
}
export interface TradingRewardsEventV1AminoMsg {
  type: "/dydxprotocol.indexer.events.TradingRewardsEventV1";
  value: TradingRewardsEventV1Amino;
}
/**
 * TradingRewardsEventV1 is communicates all trading rewards for all accounts
 * that receive trade rewards in the block.
 */
export interface TradingRewardsEventV1SDKType {
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
   * The amount of trading rewards earned by the address above in denoms. 1e18
   * denoms is equivalent to a single coin.
   */
  denomAmount: Uint8Array;
}
export interface AddressTradingRewardProtoMsg {
  typeUrl: "/dydxprotocol.indexer.events.AddressTradingReward";
  value: Uint8Array;
}
/**
 * AddressTradingReward contains info on an instance of an address receiving a
 * reward
 */
export interface AddressTradingRewardAmino {
  /** The address of the wallet that will receive the trading reward. */
  owner?: string;
  /**
   * The amount of trading rewards earned by the address above in denoms. 1e18
   * denoms is equivalent to a single coin.
   */
  denom_amount?: string;
}
export interface AddressTradingRewardAminoMsg {
  type: "/dydxprotocol.indexer.events.AddressTradingReward";
  value: AddressTradingRewardAmino;
}
/**
 * AddressTradingReward contains info on an instance of an address receiving a
 * reward
 */
export interface AddressTradingRewardSDKType {
  owner: string;
  denom_amount: Uint8Array;
}
function createBaseFundingUpdateV1(): FundingUpdateV1 {
  return {
    perpetualId: 0,
    fundingValuePpm: 0,
    fundingIndex: new Uint8Array()
  };
}
export const FundingUpdateV1 = {
  typeUrl: "/dydxprotocol.indexer.events.FundingUpdateV1",
  encode(message: FundingUpdateV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): FundingUpdateV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<FundingUpdateV1>): FundingUpdateV1 {
    const message = createBaseFundingUpdateV1();
    message.perpetualId = object.perpetualId ?? 0;
    message.fundingValuePpm = object.fundingValuePpm ?? 0;
    message.fundingIndex = object.fundingIndex ?? new Uint8Array();
    return message;
  },
  fromAmino(object: FundingUpdateV1Amino): FundingUpdateV1 {
    const message = createBaseFundingUpdateV1();
    if (object.perpetual_id !== undefined && object.perpetual_id !== null) {
      message.perpetualId = object.perpetual_id;
    }
    if (object.funding_value_ppm !== undefined && object.funding_value_ppm !== null) {
      message.fundingValuePpm = object.funding_value_ppm;
    }
    if (object.funding_index !== undefined && object.funding_index !== null) {
      message.fundingIndex = bytesFromBase64(object.funding_index);
    }
    return message;
  },
  toAmino(message: FundingUpdateV1): FundingUpdateV1Amino {
    const obj: any = {};
    obj.perpetual_id = message.perpetualId;
    obj.funding_value_ppm = message.fundingValuePpm;
    obj.funding_index = message.fundingIndex ? base64FromBytes(message.fundingIndex) : undefined;
    return obj;
  },
  fromAminoMsg(object: FundingUpdateV1AminoMsg): FundingUpdateV1 {
    return FundingUpdateV1.fromAmino(object.value);
  },
  fromProtoMsg(message: FundingUpdateV1ProtoMsg): FundingUpdateV1 {
    return FundingUpdateV1.decode(message.value);
  },
  toProto(message: FundingUpdateV1): Uint8Array {
    return FundingUpdateV1.encode(message).finish();
  },
  toProtoMsg(message: FundingUpdateV1): FundingUpdateV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.FundingUpdateV1",
      value: FundingUpdateV1.encode(message).finish()
    };
  }
};
function createBaseFundingEventV1(): FundingEventV1 {
  return {
    updates: [],
    type: 0
  };
}
export const FundingEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.FundingEventV1",
  encode(message: FundingEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.updates) {
      FundingUpdateV1.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.type !== 0) {
      writer.uint32(16).int32(message.type);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): FundingEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<FundingEventV1>): FundingEventV1 {
    const message = createBaseFundingEventV1();
    message.updates = object.updates?.map(e => FundingUpdateV1.fromPartial(e)) || [];
    message.type = object.type ?? 0;
    return message;
  },
  fromAmino(object: FundingEventV1Amino): FundingEventV1 {
    const message = createBaseFundingEventV1();
    message.updates = object.updates?.map(e => FundingUpdateV1.fromAmino(e)) || [];
    if (object.type !== undefined && object.type !== null) {
      message.type = fundingEventV1_TypeFromJSON(object.type);
    }
    return message;
  },
  toAmino(message: FundingEventV1): FundingEventV1Amino {
    const obj: any = {};
    if (message.updates) {
      obj.updates = message.updates.map(e => e ? FundingUpdateV1.toAmino(e) : undefined);
    } else {
      obj.updates = [];
    }
    obj.type = fundingEventV1_TypeToJSON(message.type);
    return obj;
  },
  fromAminoMsg(object: FundingEventV1AminoMsg): FundingEventV1 {
    return FundingEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: FundingEventV1ProtoMsg): FundingEventV1 {
    return FundingEventV1.decode(message.value);
  },
  toProto(message: FundingEventV1): Uint8Array {
    return FundingEventV1.encode(message).finish();
  },
  toProtoMsg(message: FundingEventV1): FundingEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.FundingEventV1",
      value: FundingEventV1.encode(message).finish()
    };
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
  typeUrl: "/dydxprotocol.indexer.events.MarketEventV1",
  encode(message: MarketEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): MarketEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<MarketEventV1>): MarketEventV1 {
    const message = createBaseMarketEventV1();
    message.marketId = object.marketId ?? 0;
    message.priceUpdate = object.priceUpdate !== undefined && object.priceUpdate !== null ? MarketPriceUpdateEventV1.fromPartial(object.priceUpdate) : undefined;
    message.marketCreate = object.marketCreate !== undefined && object.marketCreate !== null ? MarketCreateEventV1.fromPartial(object.marketCreate) : undefined;
    message.marketModify = object.marketModify !== undefined && object.marketModify !== null ? MarketModifyEventV1.fromPartial(object.marketModify) : undefined;
    return message;
  },
  fromAmino(object: MarketEventV1Amino): MarketEventV1 {
    const message = createBaseMarketEventV1();
    if (object.market_id !== undefined && object.market_id !== null) {
      message.marketId = object.market_id;
    }
    if (object.price_update !== undefined && object.price_update !== null) {
      message.priceUpdate = MarketPriceUpdateEventV1.fromAmino(object.price_update);
    }
    if (object.market_create !== undefined && object.market_create !== null) {
      message.marketCreate = MarketCreateEventV1.fromAmino(object.market_create);
    }
    if (object.market_modify !== undefined && object.market_modify !== null) {
      message.marketModify = MarketModifyEventV1.fromAmino(object.market_modify);
    }
    return message;
  },
  toAmino(message: MarketEventV1): MarketEventV1Amino {
    const obj: any = {};
    obj.market_id = message.marketId;
    obj.price_update = message.priceUpdate ? MarketPriceUpdateEventV1.toAmino(message.priceUpdate) : undefined;
    obj.market_create = message.marketCreate ? MarketCreateEventV1.toAmino(message.marketCreate) : undefined;
    obj.market_modify = message.marketModify ? MarketModifyEventV1.toAmino(message.marketModify) : undefined;
    return obj;
  },
  fromAminoMsg(object: MarketEventV1AminoMsg): MarketEventV1 {
    return MarketEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: MarketEventV1ProtoMsg): MarketEventV1 {
    return MarketEventV1.decode(message.value);
  },
  toProto(message: MarketEventV1): Uint8Array {
    return MarketEventV1.encode(message).finish();
  },
  toProtoMsg(message: MarketEventV1): MarketEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.MarketEventV1",
      value: MarketEventV1.encode(message).finish()
    };
  }
};
function createBaseMarketPriceUpdateEventV1(): MarketPriceUpdateEventV1 {
  return {
    priceWithExponent: BigInt(0)
  };
}
export const MarketPriceUpdateEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.MarketPriceUpdateEventV1",
  encode(message: MarketPriceUpdateEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.priceWithExponent !== BigInt(0)) {
      writer.uint32(8).uint64(message.priceWithExponent);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MarketPriceUpdateEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketPriceUpdateEventV1();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.priceWithExponent = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MarketPriceUpdateEventV1>): MarketPriceUpdateEventV1 {
    const message = createBaseMarketPriceUpdateEventV1();
    message.priceWithExponent = object.priceWithExponent !== undefined && object.priceWithExponent !== null ? BigInt(object.priceWithExponent.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MarketPriceUpdateEventV1Amino): MarketPriceUpdateEventV1 {
    const message = createBaseMarketPriceUpdateEventV1();
    if (object.price_with_exponent !== undefined && object.price_with_exponent !== null) {
      message.priceWithExponent = BigInt(object.price_with_exponent);
    }
    return message;
  },
  toAmino(message: MarketPriceUpdateEventV1): MarketPriceUpdateEventV1Amino {
    const obj: any = {};
    obj.price_with_exponent = message.priceWithExponent ? message.priceWithExponent.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MarketPriceUpdateEventV1AminoMsg): MarketPriceUpdateEventV1 {
    return MarketPriceUpdateEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: MarketPriceUpdateEventV1ProtoMsg): MarketPriceUpdateEventV1 {
    return MarketPriceUpdateEventV1.decode(message.value);
  },
  toProto(message: MarketPriceUpdateEventV1): Uint8Array {
    return MarketPriceUpdateEventV1.encode(message).finish();
  },
  toProtoMsg(message: MarketPriceUpdateEventV1): MarketPriceUpdateEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.MarketPriceUpdateEventV1",
      value: MarketPriceUpdateEventV1.encode(message).finish()
    };
  }
};
function createBaseMarketBaseEventV1(): MarketBaseEventV1 {
  return {
    pair: "",
    minPriceChangePpm: 0
  };
}
export const MarketBaseEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.MarketBaseEventV1",
  encode(message: MarketBaseEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pair !== "") {
      writer.uint32(10).string(message.pair);
    }
    if (message.minPriceChangePpm !== 0) {
      writer.uint32(16).uint32(message.minPriceChangePpm);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MarketBaseEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<MarketBaseEventV1>): MarketBaseEventV1 {
    const message = createBaseMarketBaseEventV1();
    message.pair = object.pair ?? "";
    message.minPriceChangePpm = object.minPriceChangePpm ?? 0;
    return message;
  },
  fromAmino(object: MarketBaseEventV1Amino): MarketBaseEventV1 {
    const message = createBaseMarketBaseEventV1();
    if (object.pair !== undefined && object.pair !== null) {
      message.pair = object.pair;
    }
    if (object.min_price_change_ppm !== undefined && object.min_price_change_ppm !== null) {
      message.minPriceChangePpm = object.min_price_change_ppm;
    }
    return message;
  },
  toAmino(message: MarketBaseEventV1): MarketBaseEventV1Amino {
    const obj: any = {};
    obj.pair = message.pair;
    obj.min_price_change_ppm = message.minPriceChangePpm;
    return obj;
  },
  fromAminoMsg(object: MarketBaseEventV1AminoMsg): MarketBaseEventV1 {
    return MarketBaseEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: MarketBaseEventV1ProtoMsg): MarketBaseEventV1 {
    return MarketBaseEventV1.decode(message.value);
  },
  toProto(message: MarketBaseEventV1): Uint8Array {
    return MarketBaseEventV1.encode(message).finish();
  },
  toProtoMsg(message: MarketBaseEventV1): MarketBaseEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.MarketBaseEventV1",
      value: MarketBaseEventV1.encode(message).finish()
    };
  }
};
function createBaseMarketCreateEventV1(): MarketCreateEventV1 {
  return {
    base: undefined,
    exponent: 0
  };
}
export const MarketCreateEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.MarketCreateEventV1",
  encode(message: MarketCreateEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.base !== undefined) {
      MarketBaseEventV1.encode(message.base, writer.uint32(10).fork()).ldelim();
    }
    if (message.exponent !== 0) {
      writer.uint32(16).sint32(message.exponent);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MarketCreateEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<MarketCreateEventV1>): MarketCreateEventV1 {
    const message = createBaseMarketCreateEventV1();
    message.base = object.base !== undefined && object.base !== null ? MarketBaseEventV1.fromPartial(object.base) : undefined;
    message.exponent = object.exponent ?? 0;
    return message;
  },
  fromAmino(object: MarketCreateEventV1Amino): MarketCreateEventV1 {
    const message = createBaseMarketCreateEventV1();
    if (object.base !== undefined && object.base !== null) {
      message.base = MarketBaseEventV1.fromAmino(object.base);
    }
    if (object.exponent !== undefined && object.exponent !== null) {
      message.exponent = object.exponent;
    }
    return message;
  },
  toAmino(message: MarketCreateEventV1): MarketCreateEventV1Amino {
    const obj: any = {};
    obj.base = message.base ? MarketBaseEventV1.toAmino(message.base) : undefined;
    obj.exponent = message.exponent;
    return obj;
  },
  fromAminoMsg(object: MarketCreateEventV1AminoMsg): MarketCreateEventV1 {
    return MarketCreateEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: MarketCreateEventV1ProtoMsg): MarketCreateEventV1 {
    return MarketCreateEventV1.decode(message.value);
  },
  toProto(message: MarketCreateEventV1): Uint8Array {
    return MarketCreateEventV1.encode(message).finish();
  },
  toProtoMsg(message: MarketCreateEventV1): MarketCreateEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.MarketCreateEventV1",
      value: MarketCreateEventV1.encode(message).finish()
    };
  }
};
function createBaseMarketModifyEventV1(): MarketModifyEventV1 {
  return {
    base: undefined
  };
}
export const MarketModifyEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.MarketModifyEventV1",
  encode(message: MarketModifyEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.base !== undefined) {
      MarketBaseEventV1.encode(message.base, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MarketModifyEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<MarketModifyEventV1>): MarketModifyEventV1 {
    const message = createBaseMarketModifyEventV1();
    message.base = object.base !== undefined && object.base !== null ? MarketBaseEventV1.fromPartial(object.base) : undefined;
    return message;
  },
  fromAmino(object: MarketModifyEventV1Amino): MarketModifyEventV1 {
    const message = createBaseMarketModifyEventV1();
    if (object.base !== undefined && object.base !== null) {
      message.base = MarketBaseEventV1.fromAmino(object.base);
    }
    return message;
  },
  toAmino(message: MarketModifyEventV1): MarketModifyEventV1Amino {
    const obj: any = {};
    obj.base = message.base ? MarketBaseEventV1.toAmino(message.base) : undefined;
    return obj;
  },
  fromAminoMsg(object: MarketModifyEventV1AminoMsg): MarketModifyEventV1 {
    return MarketModifyEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: MarketModifyEventV1ProtoMsg): MarketModifyEventV1 {
    return MarketModifyEventV1.decode(message.value);
  },
  toProto(message: MarketModifyEventV1): Uint8Array {
    return MarketModifyEventV1.encode(message).finish();
  },
  toProtoMsg(message: MarketModifyEventV1): MarketModifyEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.MarketModifyEventV1",
      value: MarketModifyEventV1.encode(message).finish()
    };
  }
};
function createBaseSourceOfFunds(): SourceOfFunds {
  return {
    subaccountId: undefined,
    address: undefined
  };
}
export const SourceOfFunds = {
  typeUrl: "/dydxprotocol.indexer.events.SourceOfFunds",
  encode(message: SourceOfFunds, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.subaccountId !== undefined) {
      IndexerSubaccountId.encode(message.subaccountId, writer.uint32(10).fork()).ldelim();
    }
    if (message.address !== undefined) {
      writer.uint32(18).string(message.address);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): SourceOfFunds {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<SourceOfFunds>): SourceOfFunds {
    const message = createBaseSourceOfFunds();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? IndexerSubaccountId.fromPartial(object.subaccountId) : undefined;
    message.address = object.address ?? undefined;
    return message;
  },
  fromAmino(object: SourceOfFundsAmino): SourceOfFunds {
    const message = createBaseSourceOfFunds();
    if (object.subaccount_id !== undefined && object.subaccount_id !== null) {
      message.subaccountId = IndexerSubaccountId.fromAmino(object.subaccount_id);
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    }
    return message;
  },
  toAmino(message: SourceOfFunds): SourceOfFundsAmino {
    const obj: any = {};
    obj.subaccount_id = message.subaccountId ? IndexerSubaccountId.toAmino(message.subaccountId) : undefined;
    obj.address = message.address;
    return obj;
  },
  fromAminoMsg(object: SourceOfFundsAminoMsg): SourceOfFunds {
    return SourceOfFunds.fromAmino(object.value);
  },
  fromProtoMsg(message: SourceOfFundsProtoMsg): SourceOfFunds {
    return SourceOfFunds.decode(message.value);
  },
  toProto(message: SourceOfFunds): Uint8Array {
    return SourceOfFunds.encode(message).finish();
  },
  toProtoMsg(message: SourceOfFunds): SourceOfFundsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.SourceOfFunds",
      value: SourceOfFunds.encode(message).finish()
    };
  }
};
function createBaseTransferEventV1(): TransferEventV1 {
  return {
    senderSubaccountId: undefined,
    recipientSubaccountId: undefined,
    assetId: 0,
    amount: BigInt(0),
    sender: undefined,
    recipient: undefined
  };
}
export const TransferEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.TransferEventV1",
  encode(message: TransferEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.senderSubaccountId !== undefined) {
      IndexerSubaccountId.encode(message.senderSubaccountId, writer.uint32(10).fork()).ldelim();
    }
    if (message.recipientSubaccountId !== undefined) {
      IndexerSubaccountId.encode(message.recipientSubaccountId, writer.uint32(18).fork()).ldelim();
    }
    if (message.assetId !== 0) {
      writer.uint32(24).uint32(message.assetId);
    }
    if (message.amount !== BigInt(0)) {
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
  decode(input: BinaryReader | Uint8Array, length?: number): TransferEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
          message.amount = reader.uint64();
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
  fromPartial(object: Partial<TransferEventV1>): TransferEventV1 {
    const message = createBaseTransferEventV1();
    message.senderSubaccountId = object.senderSubaccountId !== undefined && object.senderSubaccountId !== null ? IndexerSubaccountId.fromPartial(object.senderSubaccountId) : undefined;
    message.recipientSubaccountId = object.recipientSubaccountId !== undefined && object.recipientSubaccountId !== null ? IndexerSubaccountId.fromPartial(object.recipientSubaccountId) : undefined;
    message.assetId = object.assetId ?? 0;
    message.amount = object.amount !== undefined && object.amount !== null ? BigInt(object.amount.toString()) : BigInt(0);
    message.sender = object.sender !== undefined && object.sender !== null ? SourceOfFunds.fromPartial(object.sender) : undefined;
    message.recipient = object.recipient !== undefined && object.recipient !== null ? SourceOfFunds.fromPartial(object.recipient) : undefined;
    return message;
  },
  fromAmino(object: TransferEventV1Amino): TransferEventV1 {
    const message = createBaseTransferEventV1();
    if (object.sender_subaccount_id !== undefined && object.sender_subaccount_id !== null) {
      message.senderSubaccountId = IndexerSubaccountId.fromAmino(object.sender_subaccount_id);
    }
    if (object.recipient_subaccount_id !== undefined && object.recipient_subaccount_id !== null) {
      message.recipientSubaccountId = IndexerSubaccountId.fromAmino(object.recipient_subaccount_id);
    }
    if (object.asset_id !== undefined && object.asset_id !== null) {
      message.assetId = object.asset_id;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = BigInt(object.amount);
    }
    if (object.sender !== undefined && object.sender !== null) {
      message.sender = SourceOfFunds.fromAmino(object.sender);
    }
    if (object.recipient !== undefined && object.recipient !== null) {
      message.recipient = SourceOfFunds.fromAmino(object.recipient);
    }
    return message;
  },
  toAmino(message: TransferEventV1): TransferEventV1Amino {
    const obj: any = {};
    obj.sender_subaccount_id = message.senderSubaccountId ? IndexerSubaccountId.toAmino(message.senderSubaccountId) : undefined;
    obj.recipient_subaccount_id = message.recipientSubaccountId ? IndexerSubaccountId.toAmino(message.recipientSubaccountId) : undefined;
    obj.asset_id = message.assetId;
    obj.amount = message.amount ? message.amount.toString() : undefined;
    obj.sender = message.sender ? SourceOfFunds.toAmino(message.sender) : undefined;
    obj.recipient = message.recipient ? SourceOfFunds.toAmino(message.recipient) : undefined;
    return obj;
  },
  fromAminoMsg(object: TransferEventV1AminoMsg): TransferEventV1 {
    return TransferEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: TransferEventV1ProtoMsg): TransferEventV1 {
    return TransferEventV1.decode(message.value);
  },
  toProto(message: TransferEventV1): Uint8Array {
    return TransferEventV1.encode(message).finish();
  },
  toProtoMsg(message: TransferEventV1): TransferEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.TransferEventV1",
      value: TransferEventV1.encode(message).finish()
    };
  }
};
function createBaseOrderFillEventV1(): OrderFillEventV1 {
  return {
    makerOrder: IndexerOrder.fromPartial({}),
    order: undefined,
    liquidationOrder: undefined,
    fillAmount: BigInt(0),
    makerFee: BigInt(0),
    takerFee: BigInt(0),
    totalFilledMaker: BigInt(0),
    totalFilledTaker: BigInt(0)
  };
}
export const OrderFillEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.OrderFillEventV1",
  encode(message: OrderFillEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.makerOrder !== undefined) {
      IndexerOrder.encode(message.makerOrder, writer.uint32(10).fork()).ldelim();
    }
    if (message.order !== undefined) {
      IndexerOrder.encode(message.order, writer.uint32(18).fork()).ldelim();
    }
    if (message.liquidationOrder !== undefined) {
      LiquidationOrderV1.encode(message.liquidationOrder, writer.uint32(34).fork()).ldelim();
    }
    if (message.fillAmount !== BigInt(0)) {
      writer.uint32(24).uint64(message.fillAmount);
    }
    if (message.makerFee !== BigInt(0)) {
      writer.uint32(40).sint64(message.makerFee);
    }
    if (message.takerFee !== BigInt(0)) {
      writer.uint32(48).sint64(message.takerFee);
    }
    if (message.totalFilledMaker !== BigInt(0)) {
      writer.uint32(56).uint64(message.totalFilledMaker);
    }
    if (message.totalFilledTaker !== BigInt(0)) {
      writer.uint32(64).uint64(message.totalFilledTaker);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): OrderFillEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
          message.fillAmount = reader.uint64();
          break;
        case 5:
          message.makerFee = reader.sint64();
          break;
        case 6:
          message.takerFee = reader.sint64();
          break;
        case 7:
          message.totalFilledMaker = reader.uint64();
          break;
        case 8:
          message.totalFilledTaker = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<OrderFillEventV1>): OrderFillEventV1 {
    const message = createBaseOrderFillEventV1();
    message.makerOrder = object.makerOrder !== undefined && object.makerOrder !== null ? IndexerOrder.fromPartial(object.makerOrder) : undefined;
    message.order = object.order !== undefined && object.order !== null ? IndexerOrder.fromPartial(object.order) : undefined;
    message.liquidationOrder = object.liquidationOrder !== undefined && object.liquidationOrder !== null ? LiquidationOrderV1.fromPartial(object.liquidationOrder) : undefined;
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? BigInt(object.fillAmount.toString()) : BigInt(0);
    message.makerFee = object.makerFee !== undefined && object.makerFee !== null ? BigInt(object.makerFee.toString()) : BigInt(0);
    message.takerFee = object.takerFee !== undefined && object.takerFee !== null ? BigInt(object.takerFee.toString()) : BigInt(0);
    message.totalFilledMaker = object.totalFilledMaker !== undefined && object.totalFilledMaker !== null ? BigInt(object.totalFilledMaker.toString()) : BigInt(0);
    message.totalFilledTaker = object.totalFilledTaker !== undefined && object.totalFilledTaker !== null ? BigInt(object.totalFilledTaker.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: OrderFillEventV1Amino): OrderFillEventV1 {
    const message = createBaseOrderFillEventV1();
    if (object.maker_order !== undefined && object.maker_order !== null) {
      message.makerOrder = IndexerOrder.fromAmino(object.maker_order);
    }
    if (object.order !== undefined && object.order !== null) {
      message.order = IndexerOrder.fromAmino(object.order);
    }
    if (object.liquidation_order !== undefined && object.liquidation_order !== null) {
      message.liquidationOrder = LiquidationOrderV1.fromAmino(object.liquidation_order);
    }
    if (object.fill_amount !== undefined && object.fill_amount !== null) {
      message.fillAmount = BigInt(object.fill_amount);
    }
    if (object.maker_fee !== undefined && object.maker_fee !== null) {
      message.makerFee = BigInt(object.maker_fee);
    }
    if (object.taker_fee !== undefined && object.taker_fee !== null) {
      message.takerFee = BigInt(object.taker_fee);
    }
    if (object.total_filled_maker !== undefined && object.total_filled_maker !== null) {
      message.totalFilledMaker = BigInt(object.total_filled_maker);
    }
    if (object.total_filled_taker !== undefined && object.total_filled_taker !== null) {
      message.totalFilledTaker = BigInt(object.total_filled_taker);
    }
    return message;
  },
  toAmino(message: OrderFillEventV1): OrderFillEventV1Amino {
    const obj: any = {};
    obj.maker_order = message.makerOrder ? IndexerOrder.toAmino(message.makerOrder) : undefined;
    obj.order = message.order ? IndexerOrder.toAmino(message.order) : undefined;
    obj.liquidation_order = message.liquidationOrder ? LiquidationOrderV1.toAmino(message.liquidationOrder) : undefined;
    obj.fill_amount = message.fillAmount ? message.fillAmount.toString() : undefined;
    obj.maker_fee = message.makerFee ? message.makerFee.toString() : undefined;
    obj.taker_fee = message.takerFee ? message.takerFee.toString() : undefined;
    obj.total_filled_maker = message.totalFilledMaker ? message.totalFilledMaker.toString() : undefined;
    obj.total_filled_taker = message.totalFilledTaker ? message.totalFilledTaker.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: OrderFillEventV1AminoMsg): OrderFillEventV1 {
    return OrderFillEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: OrderFillEventV1ProtoMsg): OrderFillEventV1 {
    return OrderFillEventV1.decode(message.value);
  },
  toProto(message: OrderFillEventV1): Uint8Array {
    return OrderFillEventV1.encode(message).finish();
  },
  toProtoMsg(message: OrderFillEventV1): OrderFillEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.OrderFillEventV1",
      value: OrderFillEventV1.encode(message).finish()
    };
  }
};
function createBaseDeleveragingEventV1(): DeleveragingEventV1 {
  return {
    liquidated: IndexerSubaccountId.fromPartial({}),
    offsetting: IndexerSubaccountId.fromPartial({}),
    perpetualId: 0,
    fillAmount: BigInt(0),
    price: BigInt(0),
    isBuy: false,
    isFinalSettlement: false
  };
}
export const DeleveragingEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.DeleveragingEventV1",
  encode(message: DeleveragingEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.liquidated !== undefined) {
      IndexerSubaccountId.encode(message.liquidated, writer.uint32(10).fork()).ldelim();
    }
    if (message.offsetting !== undefined) {
      IndexerSubaccountId.encode(message.offsetting, writer.uint32(18).fork()).ldelim();
    }
    if (message.perpetualId !== 0) {
      writer.uint32(24).uint32(message.perpetualId);
    }
    if (message.fillAmount !== BigInt(0)) {
      writer.uint32(32).uint64(message.fillAmount);
    }
    if (message.price !== BigInt(0)) {
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
  decode(input: BinaryReader | Uint8Array, length?: number): DeleveragingEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
          message.fillAmount = reader.uint64();
          break;
        case 5:
          message.price = reader.uint64();
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
  fromPartial(object: Partial<DeleveragingEventV1>): DeleveragingEventV1 {
    const message = createBaseDeleveragingEventV1();
    message.liquidated = object.liquidated !== undefined && object.liquidated !== null ? IndexerSubaccountId.fromPartial(object.liquidated) : undefined;
    message.offsetting = object.offsetting !== undefined && object.offsetting !== null ? IndexerSubaccountId.fromPartial(object.offsetting) : undefined;
    message.perpetualId = object.perpetualId ?? 0;
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? BigInt(object.fillAmount.toString()) : BigInt(0);
    message.price = object.price !== undefined && object.price !== null ? BigInt(object.price.toString()) : BigInt(0);
    message.isBuy = object.isBuy ?? false;
    message.isFinalSettlement = object.isFinalSettlement ?? false;
    return message;
  },
  fromAmino(object: DeleveragingEventV1Amino): DeleveragingEventV1 {
    const message = createBaseDeleveragingEventV1();
    if (object.liquidated !== undefined && object.liquidated !== null) {
      message.liquidated = IndexerSubaccountId.fromAmino(object.liquidated);
    }
    if (object.offsetting !== undefined && object.offsetting !== null) {
      message.offsetting = IndexerSubaccountId.fromAmino(object.offsetting);
    }
    if (object.perpetual_id !== undefined && object.perpetual_id !== null) {
      message.perpetualId = object.perpetual_id;
    }
    if (object.fill_amount !== undefined && object.fill_amount !== null) {
      message.fillAmount = BigInt(object.fill_amount);
    }
    if (object.price !== undefined && object.price !== null) {
      message.price = BigInt(object.price);
    }
    if (object.is_buy !== undefined && object.is_buy !== null) {
      message.isBuy = object.is_buy;
    }
    if (object.is_final_settlement !== undefined && object.is_final_settlement !== null) {
      message.isFinalSettlement = object.is_final_settlement;
    }
    return message;
  },
  toAmino(message: DeleveragingEventV1): DeleveragingEventV1Amino {
    const obj: any = {};
    obj.liquidated = message.liquidated ? IndexerSubaccountId.toAmino(message.liquidated) : undefined;
    obj.offsetting = message.offsetting ? IndexerSubaccountId.toAmino(message.offsetting) : undefined;
    obj.perpetual_id = message.perpetualId;
    obj.fill_amount = message.fillAmount ? message.fillAmount.toString() : undefined;
    obj.price = message.price ? message.price.toString() : undefined;
    obj.is_buy = message.isBuy;
    obj.is_final_settlement = message.isFinalSettlement;
    return obj;
  },
  fromAminoMsg(object: DeleveragingEventV1AminoMsg): DeleveragingEventV1 {
    return DeleveragingEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: DeleveragingEventV1ProtoMsg): DeleveragingEventV1 {
    return DeleveragingEventV1.decode(message.value);
  },
  toProto(message: DeleveragingEventV1): Uint8Array {
    return DeleveragingEventV1.encode(message).finish();
  },
  toProtoMsg(message: DeleveragingEventV1): DeleveragingEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.DeleveragingEventV1",
      value: DeleveragingEventV1.encode(message).finish()
    };
  }
};
function createBaseLiquidationOrderV1(): LiquidationOrderV1 {
  return {
    liquidated: IndexerSubaccountId.fromPartial({}),
    clobPairId: 0,
    perpetualId: 0,
    totalSize: BigInt(0),
    isBuy: false,
    subticks: BigInt(0)
  };
}
export const LiquidationOrderV1 = {
  typeUrl: "/dydxprotocol.indexer.events.LiquidationOrderV1",
  encode(message: LiquidationOrderV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.liquidated !== undefined) {
      IndexerSubaccountId.encode(message.liquidated, writer.uint32(10).fork()).ldelim();
    }
    if (message.clobPairId !== 0) {
      writer.uint32(16).uint32(message.clobPairId);
    }
    if (message.perpetualId !== 0) {
      writer.uint32(24).uint32(message.perpetualId);
    }
    if (message.totalSize !== BigInt(0)) {
      writer.uint32(32).uint64(message.totalSize);
    }
    if (message.isBuy === true) {
      writer.uint32(40).bool(message.isBuy);
    }
    if (message.subticks !== BigInt(0)) {
      writer.uint32(48).uint64(message.subticks);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): LiquidationOrderV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
          message.totalSize = reader.uint64();
          break;
        case 5:
          message.isBuy = reader.bool();
          break;
        case 6:
          message.subticks = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<LiquidationOrderV1>): LiquidationOrderV1 {
    const message = createBaseLiquidationOrderV1();
    message.liquidated = object.liquidated !== undefined && object.liquidated !== null ? IndexerSubaccountId.fromPartial(object.liquidated) : undefined;
    message.clobPairId = object.clobPairId ?? 0;
    message.perpetualId = object.perpetualId ?? 0;
    message.totalSize = object.totalSize !== undefined && object.totalSize !== null ? BigInt(object.totalSize.toString()) : BigInt(0);
    message.isBuy = object.isBuy ?? false;
    message.subticks = object.subticks !== undefined && object.subticks !== null ? BigInt(object.subticks.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: LiquidationOrderV1Amino): LiquidationOrderV1 {
    const message = createBaseLiquidationOrderV1();
    if (object.liquidated !== undefined && object.liquidated !== null) {
      message.liquidated = IndexerSubaccountId.fromAmino(object.liquidated);
    }
    if (object.clob_pair_id !== undefined && object.clob_pair_id !== null) {
      message.clobPairId = object.clob_pair_id;
    }
    if (object.perpetual_id !== undefined && object.perpetual_id !== null) {
      message.perpetualId = object.perpetual_id;
    }
    if (object.total_size !== undefined && object.total_size !== null) {
      message.totalSize = BigInt(object.total_size);
    }
    if (object.is_buy !== undefined && object.is_buy !== null) {
      message.isBuy = object.is_buy;
    }
    if (object.subticks !== undefined && object.subticks !== null) {
      message.subticks = BigInt(object.subticks);
    }
    return message;
  },
  toAmino(message: LiquidationOrderV1): LiquidationOrderV1Amino {
    const obj: any = {};
    obj.liquidated = message.liquidated ? IndexerSubaccountId.toAmino(message.liquidated) : undefined;
    obj.clob_pair_id = message.clobPairId;
    obj.perpetual_id = message.perpetualId;
    obj.total_size = message.totalSize ? message.totalSize.toString() : undefined;
    obj.is_buy = message.isBuy;
    obj.subticks = message.subticks ? message.subticks.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: LiquidationOrderV1AminoMsg): LiquidationOrderV1 {
    return LiquidationOrderV1.fromAmino(object.value);
  },
  fromProtoMsg(message: LiquidationOrderV1ProtoMsg): LiquidationOrderV1 {
    return LiquidationOrderV1.decode(message.value);
  },
  toProto(message: LiquidationOrderV1): Uint8Array {
    return LiquidationOrderV1.encode(message).finish();
  },
  toProtoMsg(message: LiquidationOrderV1): LiquidationOrderV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.LiquidationOrderV1",
      value: LiquidationOrderV1.encode(message).finish()
    };
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
  typeUrl: "/dydxprotocol.indexer.events.SubaccountUpdateEventV1",
  encode(message: SubaccountUpdateEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): SubaccountUpdateEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<SubaccountUpdateEventV1>): SubaccountUpdateEventV1 {
    const message = createBaseSubaccountUpdateEventV1();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? IndexerSubaccountId.fromPartial(object.subaccountId) : undefined;
    message.updatedPerpetualPositions = object.updatedPerpetualPositions?.map(e => IndexerPerpetualPosition.fromPartial(e)) || [];
    message.updatedAssetPositions = object.updatedAssetPositions?.map(e => IndexerAssetPosition.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: SubaccountUpdateEventV1Amino): SubaccountUpdateEventV1 {
    const message = createBaseSubaccountUpdateEventV1();
    if (object.subaccount_id !== undefined && object.subaccount_id !== null) {
      message.subaccountId = IndexerSubaccountId.fromAmino(object.subaccount_id);
    }
    message.updatedPerpetualPositions = object.updated_perpetual_positions?.map(e => IndexerPerpetualPosition.fromAmino(e)) || [];
    message.updatedAssetPositions = object.updated_asset_positions?.map(e => IndexerAssetPosition.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: SubaccountUpdateEventV1): SubaccountUpdateEventV1Amino {
    const obj: any = {};
    obj.subaccount_id = message.subaccountId ? IndexerSubaccountId.toAmino(message.subaccountId) : undefined;
    if (message.updatedPerpetualPositions) {
      obj.updated_perpetual_positions = message.updatedPerpetualPositions.map(e => e ? IndexerPerpetualPosition.toAmino(e) : undefined);
    } else {
      obj.updated_perpetual_positions = [];
    }
    if (message.updatedAssetPositions) {
      obj.updated_asset_positions = message.updatedAssetPositions.map(e => e ? IndexerAssetPosition.toAmino(e) : undefined);
    } else {
      obj.updated_asset_positions = [];
    }
    return obj;
  },
  fromAminoMsg(object: SubaccountUpdateEventV1AminoMsg): SubaccountUpdateEventV1 {
    return SubaccountUpdateEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: SubaccountUpdateEventV1ProtoMsg): SubaccountUpdateEventV1 {
    return SubaccountUpdateEventV1.decode(message.value);
  },
  toProto(message: SubaccountUpdateEventV1): Uint8Array {
    return SubaccountUpdateEventV1.encode(message).finish();
  },
  toProtoMsg(message: SubaccountUpdateEventV1): SubaccountUpdateEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.SubaccountUpdateEventV1",
      value: SubaccountUpdateEventV1.encode(message).finish()
    };
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
  typeUrl: "/dydxprotocol.indexer.events.StatefulOrderEventV1",
  encode(message: StatefulOrderEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): StatefulOrderEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<StatefulOrderEventV1>): StatefulOrderEventV1 {
    const message = createBaseStatefulOrderEventV1();
    message.orderPlace = object.orderPlace !== undefined && object.orderPlace !== null ? StatefulOrderEventV1_StatefulOrderPlacementV1.fromPartial(object.orderPlace) : undefined;
    message.orderRemoval = object.orderRemoval !== undefined && object.orderRemoval !== null ? StatefulOrderEventV1_StatefulOrderRemovalV1.fromPartial(object.orderRemoval) : undefined;
    message.conditionalOrderPlacement = object.conditionalOrderPlacement !== undefined && object.conditionalOrderPlacement !== null ? StatefulOrderEventV1_ConditionalOrderPlacementV1.fromPartial(object.conditionalOrderPlacement) : undefined;
    message.conditionalOrderTriggered = object.conditionalOrderTriggered !== undefined && object.conditionalOrderTriggered !== null ? StatefulOrderEventV1_ConditionalOrderTriggeredV1.fromPartial(object.conditionalOrderTriggered) : undefined;
    message.longTermOrderPlacement = object.longTermOrderPlacement !== undefined && object.longTermOrderPlacement !== null ? StatefulOrderEventV1_LongTermOrderPlacementV1.fromPartial(object.longTermOrderPlacement) : undefined;
    return message;
  },
  fromAmino(object: StatefulOrderEventV1Amino): StatefulOrderEventV1 {
    const message = createBaseStatefulOrderEventV1();
    if (object.order_place !== undefined && object.order_place !== null) {
      message.orderPlace = StatefulOrderEventV1_StatefulOrderPlacementV1.fromAmino(object.order_place);
    }
    if (object.order_removal !== undefined && object.order_removal !== null) {
      message.orderRemoval = StatefulOrderEventV1_StatefulOrderRemovalV1.fromAmino(object.order_removal);
    }
    if (object.conditional_order_placement !== undefined && object.conditional_order_placement !== null) {
      message.conditionalOrderPlacement = StatefulOrderEventV1_ConditionalOrderPlacementV1.fromAmino(object.conditional_order_placement);
    }
    if (object.conditional_order_triggered !== undefined && object.conditional_order_triggered !== null) {
      message.conditionalOrderTriggered = StatefulOrderEventV1_ConditionalOrderTriggeredV1.fromAmino(object.conditional_order_triggered);
    }
    if (object.long_term_order_placement !== undefined && object.long_term_order_placement !== null) {
      message.longTermOrderPlacement = StatefulOrderEventV1_LongTermOrderPlacementV1.fromAmino(object.long_term_order_placement);
    }
    return message;
  },
  toAmino(message: StatefulOrderEventV1): StatefulOrderEventV1Amino {
    const obj: any = {};
    obj.order_place = message.orderPlace ? StatefulOrderEventV1_StatefulOrderPlacementV1.toAmino(message.orderPlace) : undefined;
    obj.order_removal = message.orderRemoval ? StatefulOrderEventV1_StatefulOrderRemovalV1.toAmino(message.orderRemoval) : undefined;
    obj.conditional_order_placement = message.conditionalOrderPlacement ? StatefulOrderEventV1_ConditionalOrderPlacementV1.toAmino(message.conditionalOrderPlacement) : undefined;
    obj.conditional_order_triggered = message.conditionalOrderTriggered ? StatefulOrderEventV1_ConditionalOrderTriggeredV1.toAmino(message.conditionalOrderTriggered) : undefined;
    obj.long_term_order_placement = message.longTermOrderPlacement ? StatefulOrderEventV1_LongTermOrderPlacementV1.toAmino(message.longTermOrderPlacement) : undefined;
    return obj;
  },
  fromAminoMsg(object: StatefulOrderEventV1AminoMsg): StatefulOrderEventV1 {
    return StatefulOrderEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: StatefulOrderEventV1ProtoMsg): StatefulOrderEventV1 {
    return StatefulOrderEventV1.decode(message.value);
  },
  toProto(message: StatefulOrderEventV1): Uint8Array {
    return StatefulOrderEventV1.encode(message).finish();
  },
  toProtoMsg(message: StatefulOrderEventV1): StatefulOrderEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.StatefulOrderEventV1",
      value: StatefulOrderEventV1.encode(message).finish()
    };
  }
};
function createBaseStatefulOrderEventV1_StatefulOrderPlacementV1(): StatefulOrderEventV1_StatefulOrderPlacementV1 {
  return {
    order: undefined
  };
}
export const StatefulOrderEventV1_StatefulOrderPlacementV1 = {
  typeUrl: "/dydxprotocol.indexer.events.StatefulOrderPlacementV1",
  encode(message: StatefulOrderEventV1_StatefulOrderPlacementV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.order !== undefined) {
      IndexerOrder.encode(message.order, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): StatefulOrderEventV1_StatefulOrderPlacementV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<StatefulOrderEventV1_StatefulOrderPlacementV1>): StatefulOrderEventV1_StatefulOrderPlacementV1 {
    const message = createBaseStatefulOrderEventV1_StatefulOrderPlacementV1();
    message.order = object.order !== undefined && object.order !== null ? IndexerOrder.fromPartial(object.order) : undefined;
    return message;
  },
  fromAmino(object: StatefulOrderEventV1_StatefulOrderPlacementV1Amino): StatefulOrderEventV1_StatefulOrderPlacementV1 {
    const message = createBaseStatefulOrderEventV1_StatefulOrderPlacementV1();
    if (object.order !== undefined && object.order !== null) {
      message.order = IndexerOrder.fromAmino(object.order);
    }
    return message;
  },
  toAmino(message: StatefulOrderEventV1_StatefulOrderPlacementV1): StatefulOrderEventV1_StatefulOrderPlacementV1Amino {
    const obj: any = {};
    obj.order = message.order ? IndexerOrder.toAmino(message.order) : undefined;
    return obj;
  },
  fromAminoMsg(object: StatefulOrderEventV1_StatefulOrderPlacementV1AminoMsg): StatefulOrderEventV1_StatefulOrderPlacementV1 {
    return StatefulOrderEventV1_StatefulOrderPlacementV1.fromAmino(object.value);
  },
  fromProtoMsg(message: StatefulOrderEventV1_StatefulOrderPlacementV1ProtoMsg): StatefulOrderEventV1_StatefulOrderPlacementV1 {
    return StatefulOrderEventV1_StatefulOrderPlacementV1.decode(message.value);
  },
  toProto(message: StatefulOrderEventV1_StatefulOrderPlacementV1): Uint8Array {
    return StatefulOrderEventV1_StatefulOrderPlacementV1.encode(message).finish();
  },
  toProtoMsg(message: StatefulOrderEventV1_StatefulOrderPlacementV1): StatefulOrderEventV1_StatefulOrderPlacementV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.StatefulOrderPlacementV1",
      value: StatefulOrderEventV1_StatefulOrderPlacementV1.encode(message).finish()
    };
  }
};
function createBaseStatefulOrderEventV1_StatefulOrderRemovalV1(): StatefulOrderEventV1_StatefulOrderRemovalV1 {
  return {
    removedOrderId: undefined,
    reason: 0
  };
}
export const StatefulOrderEventV1_StatefulOrderRemovalV1 = {
  typeUrl: "/dydxprotocol.indexer.events.StatefulOrderRemovalV1",
  encode(message: StatefulOrderEventV1_StatefulOrderRemovalV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.removedOrderId !== undefined) {
      IndexerOrderId.encode(message.removedOrderId, writer.uint32(10).fork()).ldelim();
    }
    if (message.reason !== 0) {
      writer.uint32(16).int32(message.reason);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): StatefulOrderEventV1_StatefulOrderRemovalV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<StatefulOrderEventV1_StatefulOrderRemovalV1>): StatefulOrderEventV1_StatefulOrderRemovalV1 {
    const message = createBaseStatefulOrderEventV1_StatefulOrderRemovalV1();
    message.removedOrderId = object.removedOrderId !== undefined && object.removedOrderId !== null ? IndexerOrderId.fromPartial(object.removedOrderId) : undefined;
    message.reason = object.reason ?? 0;
    return message;
  },
  fromAmino(object: StatefulOrderEventV1_StatefulOrderRemovalV1Amino): StatefulOrderEventV1_StatefulOrderRemovalV1 {
    const message = createBaseStatefulOrderEventV1_StatefulOrderRemovalV1();
    if (object.removed_order_id !== undefined && object.removed_order_id !== null) {
      message.removedOrderId = IndexerOrderId.fromAmino(object.removed_order_id);
    }
    if (object.reason !== undefined && object.reason !== null) {
      message.reason = orderRemovalReasonFromJSON(object.reason);
    }
    return message;
  },
  toAmino(message: StatefulOrderEventV1_StatefulOrderRemovalV1): StatefulOrderEventV1_StatefulOrderRemovalV1Amino {
    const obj: any = {};
    obj.removed_order_id = message.removedOrderId ? IndexerOrderId.toAmino(message.removedOrderId) : undefined;
    obj.reason = orderRemovalReasonToJSON(message.reason);
    return obj;
  },
  fromAminoMsg(object: StatefulOrderEventV1_StatefulOrderRemovalV1AminoMsg): StatefulOrderEventV1_StatefulOrderRemovalV1 {
    return StatefulOrderEventV1_StatefulOrderRemovalV1.fromAmino(object.value);
  },
  fromProtoMsg(message: StatefulOrderEventV1_StatefulOrderRemovalV1ProtoMsg): StatefulOrderEventV1_StatefulOrderRemovalV1 {
    return StatefulOrderEventV1_StatefulOrderRemovalV1.decode(message.value);
  },
  toProto(message: StatefulOrderEventV1_StatefulOrderRemovalV1): Uint8Array {
    return StatefulOrderEventV1_StatefulOrderRemovalV1.encode(message).finish();
  },
  toProtoMsg(message: StatefulOrderEventV1_StatefulOrderRemovalV1): StatefulOrderEventV1_StatefulOrderRemovalV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.StatefulOrderRemovalV1",
      value: StatefulOrderEventV1_StatefulOrderRemovalV1.encode(message).finish()
    };
  }
};
function createBaseStatefulOrderEventV1_ConditionalOrderPlacementV1(): StatefulOrderEventV1_ConditionalOrderPlacementV1 {
  return {
    order: undefined
  };
}
export const StatefulOrderEventV1_ConditionalOrderPlacementV1 = {
  typeUrl: "/dydxprotocol.indexer.events.ConditionalOrderPlacementV1",
  encode(message: StatefulOrderEventV1_ConditionalOrderPlacementV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.order !== undefined) {
      IndexerOrder.encode(message.order, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): StatefulOrderEventV1_ConditionalOrderPlacementV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<StatefulOrderEventV1_ConditionalOrderPlacementV1>): StatefulOrderEventV1_ConditionalOrderPlacementV1 {
    const message = createBaseStatefulOrderEventV1_ConditionalOrderPlacementV1();
    message.order = object.order !== undefined && object.order !== null ? IndexerOrder.fromPartial(object.order) : undefined;
    return message;
  },
  fromAmino(object: StatefulOrderEventV1_ConditionalOrderPlacementV1Amino): StatefulOrderEventV1_ConditionalOrderPlacementV1 {
    const message = createBaseStatefulOrderEventV1_ConditionalOrderPlacementV1();
    if (object.order !== undefined && object.order !== null) {
      message.order = IndexerOrder.fromAmino(object.order);
    }
    return message;
  },
  toAmino(message: StatefulOrderEventV1_ConditionalOrderPlacementV1): StatefulOrderEventV1_ConditionalOrderPlacementV1Amino {
    const obj: any = {};
    obj.order = message.order ? IndexerOrder.toAmino(message.order) : undefined;
    return obj;
  },
  fromAminoMsg(object: StatefulOrderEventV1_ConditionalOrderPlacementV1AminoMsg): StatefulOrderEventV1_ConditionalOrderPlacementV1 {
    return StatefulOrderEventV1_ConditionalOrderPlacementV1.fromAmino(object.value);
  },
  fromProtoMsg(message: StatefulOrderEventV1_ConditionalOrderPlacementV1ProtoMsg): StatefulOrderEventV1_ConditionalOrderPlacementV1 {
    return StatefulOrderEventV1_ConditionalOrderPlacementV1.decode(message.value);
  },
  toProto(message: StatefulOrderEventV1_ConditionalOrderPlacementV1): Uint8Array {
    return StatefulOrderEventV1_ConditionalOrderPlacementV1.encode(message).finish();
  },
  toProtoMsg(message: StatefulOrderEventV1_ConditionalOrderPlacementV1): StatefulOrderEventV1_ConditionalOrderPlacementV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.ConditionalOrderPlacementV1",
      value: StatefulOrderEventV1_ConditionalOrderPlacementV1.encode(message).finish()
    };
  }
};
function createBaseStatefulOrderEventV1_ConditionalOrderTriggeredV1(): StatefulOrderEventV1_ConditionalOrderTriggeredV1 {
  return {
    triggeredOrderId: undefined
  };
}
export const StatefulOrderEventV1_ConditionalOrderTriggeredV1 = {
  typeUrl: "/dydxprotocol.indexer.events.ConditionalOrderTriggeredV1",
  encode(message: StatefulOrderEventV1_ConditionalOrderTriggeredV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.triggeredOrderId !== undefined) {
      IndexerOrderId.encode(message.triggeredOrderId, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): StatefulOrderEventV1_ConditionalOrderTriggeredV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<StatefulOrderEventV1_ConditionalOrderTriggeredV1>): StatefulOrderEventV1_ConditionalOrderTriggeredV1 {
    const message = createBaseStatefulOrderEventV1_ConditionalOrderTriggeredV1();
    message.triggeredOrderId = object.triggeredOrderId !== undefined && object.triggeredOrderId !== null ? IndexerOrderId.fromPartial(object.triggeredOrderId) : undefined;
    return message;
  },
  fromAmino(object: StatefulOrderEventV1_ConditionalOrderTriggeredV1Amino): StatefulOrderEventV1_ConditionalOrderTriggeredV1 {
    const message = createBaseStatefulOrderEventV1_ConditionalOrderTriggeredV1();
    if (object.triggered_order_id !== undefined && object.triggered_order_id !== null) {
      message.triggeredOrderId = IndexerOrderId.fromAmino(object.triggered_order_id);
    }
    return message;
  },
  toAmino(message: StatefulOrderEventV1_ConditionalOrderTriggeredV1): StatefulOrderEventV1_ConditionalOrderTriggeredV1Amino {
    const obj: any = {};
    obj.triggered_order_id = message.triggeredOrderId ? IndexerOrderId.toAmino(message.triggeredOrderId) : undefined;
    return obj;
  },
  fromAminoMsg(object: StatefulOrderEventV1_ConditionalOrderTriggeredV1AminoMsg): StatefulOrderEventV1_ConditionalOrderTriggeredV1 {
    return StatefulOrderEventV1_ConditionalOrderTriggeredV1.fromAmino(object.value);
  },
  fromProtoMsg(message: StatefulOrderEventV1_ConditionalOrderTriggeredV1ProtoMsg): StatefulOrderEventV1_ConditionalOrderTriggeredV1 {
    return StatefulOrderEventV1_ConditionalOrderTriggeredV1.decode(message.value);
  },
  toProto(message: StatefulOrderEventV1_ConditionalOrderTriggeredV1): Uint8Array {
    return StatefulOrderEventV1_ConditionalOrderTriggeredV1.encode(message).finish();
  },
  toProtoMsg(message: StatefulOrderEventV1_ConditionalOrderTriggeredV1): StatefulOrderEventV1_ConditionalOrderTriggeredV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.ConditionalOrderTriggeredV1",
      value: StatefulOrderEventV1_ConditionalOrderTriggeredV1.encode(message).finish()
    };
  }
};
function createBaseStatefulOrderEventV1_LongTermOrderPlacementV1(): StatefulOrderEventV1_LongTermOrderPlacementV1 {
  return {
    order: undefined
  };
}
export const StatefulOrderEventV1_LongTermOrderPlacementV1 = {
  typeUrl: "/dydxprotocol.indexer.events.LongTermOrderPlacementV1",
  encode(message: StatefulOrderEventV1_LongTermOrderPlacementV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.order !== undefined) {
      IndexerOrder.encode(message.order, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): StatefulOrderEventV1_LongTermOrderPlacementV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<StatefulOrderEventV1_LongTermOrderPlacementV1>): StatefulOrderEventV1_LongTermOrderPlacementV1 {
    const message = createBaseStatefulOrderEventV1_LongTermOrderPlacementV1();
    message.order = object.order !== undefined && object.order !== null ? IndexerOrder.fromPartial(object.order) : undefined;
    return message;
  },
  fromAmino(object: StatefulOrderEventV1_LongTermOrderPlacementV1Amino): StatefulOrderEventV1_LongTermOrderPlacementV1 {
    const message = createBaseStatefulOrderEventV1_LongTermOrderPlacementV1();
    if (object.order !== undefined && object.order !== null) {
      message.order = IndexerOrder.fromAmino(object.order);
    }
    return message;
  },
  toAmino(message: StatefulOrderEventV1_LongTermOrderPlacementV1): StatefulOrderEventV1_LongTermOrderPlacementV1Amino {
    const obj: any = {};
    obj.order = message.order ? IndexerOrder.toAmino(message.order) : undefined;
    return obj;
  },
  fromAminoMsg(object: StatefulOrderEventV1_LongTermOrderPlacementV1AminoMsg): StatefulOrderEventV1_LongTermOrderPlacementV1 {
    return StatefulOrderEventV1_LongTermOrderPlacementV1.fromAmino(object.value);
  },
  fromProtoMsg(message: StatefulOrderEventV1_LongTermOrderPlacementV1ProtoMsg): StatefulOrderEventV1_LongTermOrderPlacementV1 {
    return StatefulOrderEventV1_LongTermOrderPlacementV1.decode(message.value);
  },
  toProto(message: StatefulOrderEventV1_LongTermOrderPlacementV1): Uint8Array {
    return StatefulOrderEventV1_LongTermOrderPlacementV1.encode(message).finish();
  },
  toProtoMsg(message: StatefulOrderEventV1_LongTermOrderPlacementV1): StatefulOrderEventV1_LongTermOrderPlacementV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.LongTermOrderPlacementV1",
      value: StatefulOrderEventV1_LongTermOrderPlacementV1.encode(message).finish()
    };
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
  typeUrl: "/dydxprotocol.indexer.events.AssetCreateEventV1",
  encode(message: AssetCreateEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): AssetCreateEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<AssetCreateEventV1>): AssetCreateEventV1 {
    const message = createBaseAssetCreateEventV1();
    message.id = object.id ?? 0;
    message.symbol = object.symbol ?? "";
    message.hasMarket = object.hasMarket ?? false;
    message.marketId = object.marketId ?? 0;
    message.atomicResolution = object.atomicResolution ?? 0;
    return message;
  },
  fromAmino(object: AssetCreateEventV1Amino): AssetCreateEventV1 {
    const message = createBaseAssetCreateEventV1();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    if (object.symbol !== undefined && object.symbol !== null) {
      message.symbol = object.symbol;
    }
    if (object.has_market !== undefined && object.has_market !== null) {
      message.hasMarket = object.has_market;
    }
    if (object.market_id !== undefined && object.market_id !== null) {
      message.marketId = object.market_id;
    }
    if (object.atomic_resolution !== undefined && object.atomic_resolution !== null) {
      message.atomicResolution = object.atomic_resolution;
    }
    return message;
  },
  toAmino(message: AssetCreateEventV1): AssetCreateEventV1Amino {
    const obj: any = {};
    obj.id = message.id;
    obj.symbol = message.symbol;
    obj.has_market = message.hasMarket;
    obj.market_id = message.marketId;
    obj.atomic_resolution = message.atomicResolution;
    return obj;
  },
  fromAminoMsg(object: AssetCreateEventV1AminoMsg): AssetCreateEventV1 {
    return AssetCreateEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: AssetCreateEventV1ProtoMsg): AssetCreateEventV1 {
    return AssetCreateEventV1.decode(message.value);
  },
  toProto(message: AssetCreateEventV1): Uint8Array {
    return AssetCreateEventV1.encode(message).finish();
  },
  toProtoMsg(message: AssetCreateEventV1): AssetCreateEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.AssetCreateEventV1",
      value: AssetCreateEventV1.encode(message).finish()
    };
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
    stepBaseQuantums: BigInt(0),
    liquidityTier: 0
  };
}
export const PerpetualMarketCreateEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.PerpetualMarketCreateEventV1",
  encode(message: PerpetualMarketCreateEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
    if (message.stepBaseQuantums !== BigInt(0)) {
      writer.uint32(72).uint64(message.stepBaseQuantums);
    }
    if (message.liquidityTier !== 0) {
      writer.uint32(80).uint32(message.liquidityTier);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): PerpetualMarketCreateEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
          message.stepBaseQuantums = reader.uint64();
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
  fromPartial(object: Partial<PerpetualMarketCreateEventV1>): PerpetualMarketCreateEventV1 {
    const message = createBasePerpetualMarketCreateEventV1();
    message.id = object.id ?? 0;
    message.clobPairId = object.clobPairId ?? 0;
    message.ticker = object.ticker ?? "";
    message.marketId = object.marketId ?? 0;
    message.status = object.status ?? 0;
    message.quantumConversionExponent = object.quantumConversionExponent ?? 0;
    message.atomicResolution = object.atomicResolution ?? 0;
    message.subticksPerTick = object.subticksPerTick ?? 0;
    message.stepBaseQuantums = object.stepBaseQuantums !== undefined && object.stepBaseQuantums !== null ? BigInt(object.stepBaseQuantums.toString()) : BigInt(0);
    message.liquidityTier = object.liquidityTier ?? 0;
    return message;
  },
  fromAmino(object: PerpetualMarketCreateEventV1Amino): PerpetualMarketCreateEventV1 {
    const message = createBasePerpetualMarketCreateEventV1();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    if (object.clob_pair_id !== undefined && object.clob_pair_id !== null) {
      message.clobPairId = object.clob_pair_id;
    }
    if (object.ticker !== undefined && object.ticker !== null) {
      message.ticker = object.ticker;
    }
    if (object.market_id !== undefined && object.market_id !== null) {
      message.marketId = object.market_id;
    }
    if (object.status !== undefined && object.status !== null) {
      message.status = clobPairStatusFromJSON(object.status);
    }
    if (object.quantum_conversion_exponent !== undefined && object.quantum_conversion_exponent !== null) {
      message.quantumConversionExponent = object.quantum_conversion_exponent;
    }
    if (object.atomic_resolution !== undefined && object.atomic_resolution !== null) {
      message.atomicResolution = object.atomic_resolution;
    }
    if (object.subticks_per_tick !== undefined && object.subticks_per_tick !== null) {
      message.subticksPerTick = object.subticks_per_tick;
    }
    if (object.step_base_quantums !== undefined && object.step_base_quantums !== null) {
      message.stepBaseQuantums = BigInt(object.step_base_quantums);
    }
    if (object.liquidity_tier !== undefined && object.liquidity_tier !== null) {
      message.liquidityTier = object.liquidity_tier;
    }
    return message;
  },
  toAmino(message: PerpetualMarketCreateEventV1): PerpetualMarketCreateEventV1Amino {
    const obj: any = {};
    obj.id = message.id;
    obj.clob_pair_id = message.clobPairId;
    obj.ticker = message.ticker;
    obj.market_id = message.marketId;
    obj.status = clobPairStatusToJSON(message.status);
    obj.quantum_conversion_exponent = message.quantumConversionExponent;
    obj.atomic_resolution = message.atomicResolution;
    obj.subticks_per_tick = message.subticksPerTick;
    obj.step_base_quantums = message.stepBaseQuantums ? message.stepBaseQuantums.toString() : undefined;
    obj.liquidity_tier = message.liquidityTier;
    return obj;
  },
  fromAminoMsg(object: PerpetualMarketCreateEventV1AminoMsg): PerpetualMarketCreateEventV1 {
    return PerpetualMarketCreateEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: PerpetualMarketCreateEventV1ProtoMsg): PerpetualMarketCreateEventV1 {
    return PerpetualMarketCreateEventV1.decode(message.value);
  },
  toProto(message: PerpetualMarketCreateEventV1): Uint8Array {
    return PerpetualMarketCreateEventV1.encode(message).finish();
  },
  toProtoMsg(message: PerpetualMarketCreateEventV1): PerpetualMarketCreateEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.PerpetualMarketCreateEventV1",
      value: PerpetualMarketCreateEventV1.encode(message).finish()
    };
  }
};
function createBaseLiquidityTierUpsertEventV1(): LiquidityTierUpsertEventV1 {
  return {
    id: 0,
    name: "",
    initialMarginPpm: 0,
    maintenanceFractionPpm: 0,
    basePositionNotional: BigInt(0)
  };
}
export const LiquidityTierUpsertEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.LiquidityTierUpsertEventV1",
  encode(message: LiquidityTierUpsertEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
    if (message.basePositionNotional !== BigInt(0)) {
      writer.uint32(40).uint64(message.basePositionNotional);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): LiquidityTierUpsertEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
          message.basePositionNotional = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<LiquidityTierUpsertEventV1>): LiquidityTierUpsertEventV1 {
    const message = createBaseLiquidityTierUpsertEventV1();
    message.id = object.id ?? 0;
    message.name = object.name ?? "";
    message.initialMarginPpm = object.initialMarginPpm ?? 0;
    message.maintenanceFractionPpm = object.maintenanceFractionPpm ?? 0;
    message.basePositionNotional = object.basePositionNotional !== undefined && object.basePositionNotional !== null ? BigInt(object.basePositionNotional.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: LiquidityTierUpsertEventV1Amino): LiquidityTierUpsertEventV1 {
    const message = createBaseLiquidityTierUpsertEventV1();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = object.name;
    }
    if (object.initial_margin_ppm !== undefined && object.initial_margin_ppm !== null) {
      message.initialMarginPpm = object.initial_margin_ppm;
    }
    if (object.maintenance_fraction_ppm !== undefined && object.maintenance_fraction_ppm !== null) {
      message.maintenanceFractionPpm = object.maintenance_fraction_ppm;
    }
    if (object.base_position_notional !== undefined && object.base_position_notional !== null) {
      message.basePositionNotional = BigInt(object.base_position_notional);
    }
    return message;
  },
  toAmino(message: LiquidityTierUpsertEventV1): LiquidityTierUpsertEventV1Amino {
    const obj: any = {};
    obj.id = message.id;
    obj.name = message.name;
    obj.initial_margin_ppm = message.initialMarginPpm;
    obj.maintenance_fraction_ppm = message.maintenanceFractionPpm;
    obj.base_position_notional = message.basePositionNotional ? message.basePositionNotional.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: LiquidityTierUpsertEventV1AminoMsg): LiquidityTierUpsertEventV1 {
    return LiquidityTierUpsertEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: LiquidityTierUpsertEventV1ProtoMsg): LiquidityTierUpsertEventV1 {
    return LiquidityTierUpsertEventV1.decode(message.value);
  },
  toProto(message: LiquidityTierUpsertEventV1): Uint8Array {
    return LiquidityTierUpsertEventV1.encode(message).finish();
  },
  toProtoMsg(message: LiquidityTierUpsertEventV1): LiquidityTierUpsertEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.LiquidityTierUpsertEventV1",
      value: LiquidityTierUpsertEventV1.encode(message).finish()
    };
  }
};
function createBaseUpdateClobPairEventV1(): UpdateClobPairEventV1 {
  return {
    clobPairId: 0,
    status: 0,
    quantumConversionExponent: 0,
    subticksPerTick: 0,
    stepBaseQuantums: BigInt(0)
  };
}
export const UpdateClobPairEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.UpdateClobPairEventV1",
  encode(message: UpdateClobPairEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
    if (message.stepBaseQuantums !== BigInt(0)) {
      writer.uint32(40).uint64(message.stepBaseQuantums);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): UpdateClobPairEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
          message.stepBaseQuantums = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<UpdateClobPairEventV1>): UpdateClobPairEventV1 {
    const message = createBaseUpdateClobPairEventV1();
    message.clobPairId = object.clobPairId ?? 0;
    message.status = object.status ?? 0;
    message.quantumConversionExponent = object.quantumConversionExponent ?? 0;
    message.subticksPerTick = object.subticksPerTick ?? 0;
    message.stepBaseQuantums = object.stepBaseQuantums !== undefined && object.stepBaseQuantums !== null ? BigInt(object.stepBaseQuantums.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: UpdateClobPairEventV1Amino): UpdateClobPairEventV1 {
    const message = createBaseUpdateClobPairEventV1();
    if (object.clob_pair_id !== undefined && object.clob_pair_id !== null) {
      message.clobPairId = object.clob_pair_id;
    }
    if (object.status !== undefined && object.status !== null) {
      message.status = clobPairStatusFromJSON(object.status);
    }
    if (object.quantum_conversion_exponent !== undefined && object.quantum_conversion_exponent !== null) {
      message.quantumConversionExponent = object.quantum_conversion_exponent;
    }
    if (object.subticks_per_tick !== undefined && object.subticks_per_tick !== null) {
      message.subticksPerTick = object.subticks_per_tick;
    }
    if (object.step_base_quantums !== undefined && object.step_base_quantums !== null) {
      message.stepBaseQuantums = BigInt(object.step_base_quantums);
    }
    return message;
  },
  toAmino(message: UpdateClobPairEventV1): UpdateClobPairEventV1Amino {
    const obj: any = {};
    obj.clob_pair_id = message.clobPairId;
    obj.status = clobPairStatusToJSON(message.status);
    obj.quantum_conversion_exponent = message.quantumConversionExponent;
    obj.subticks_per_tick = message.subticksPerTick;
    obj.step_base_quantums = message.stepBaseQuantums ? message.stepBaseQuantums.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: UpdateClobPairEventV1AminoMsg): UpdateClobPairEventV1 {
    return UpdateClobPairEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: UpdateClobPairEventV1ProtoMsg): UpdateClobPairEventV1 {
    return UpdateClobPairEventV1.decode(message.value);
  },
  toProto(message: UpdateClobPairEventV1): Uint8Array {
    return UpdateClobPairEventV1.encode(message).finish();
  },
  toProtoMsg(message: UpdateClobPairEventV1): UpdateClobPairEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.UpdateClobPairEventV1",
      value: UpdateClobPairEventV1.encode(message).finish()
    };
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
  typeUrl: "/dydxprotocol.indexer.events.UpdatePerpetualEventV1",
  encode(message: UpdatePerpetualEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): UpdatePerpetualEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<UpdatePerpetualEventV1>): UpdatePerpetualEventV1 {
    const message = createBaseUpdatePerpetualEventV1();
    message.id = object.id ?? 0;
    message.ticker = object.ticker ?? "";
    message.marketId = object.marketId ?? 0;
    message.atomicResolution = object.atomicResolution ?? 0;
    message.liquidityTier = object.liquidityTier ?? 0;
    return message;
  },
  fromAmino(object: UpdatePerpetualEventV1Amino): UpdatePerpetualEventV1 {
    const message = createBaseUpdatePerpetualEventV1();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    if (object.ticker !== undefined && object.ticker !== null) {
      message.ticker = object.ticker;
    }
    if (object.market_id !== undefined && object.market_id !== null) {
      message.marketId = object.market_id;
    }
    if (object.atomic_resolution !== undefined && object.atomic_resolution !== null) {
      message.atomicResolution = object.atomic_resolution;
    }
    if (object.liquidity_tier !== undefined && object.liquidity_tier !== null) {
      message.liquidityTier = object.liquidity_tier;
    }
    return message;
  },
  toAmino(message: UpdatePerpetualEventV1): UpdatePerpetualEventV1Amino {
    const obj: any = {};
    obj.id = message.id;
    obj.ticker = message.ticker;
    obj.market_id = message.marketId;
    obj.atomic_resolution = message.atomicResolution;
    obj.liquidity_tier = message.liquidityTier;
    return obj;
  },
  fromAminoMsg(object: UpdatePerpetualEventV1AminoMsg): UpdatePerpetualEventV1 {
    return UpdatePerpetualEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: UpdatePerpetualEventV1ProtoMsg): UpdatePerpetualEventV1 {
    return UpdatePerpetualEventV1.decode(message.value);
  },
  toProto(message: UpdatePerpetualEventV1): Uint8Array {
    return UpdatePerpetualEventV1.encode(message).finish();
  },
  toProtoMsg(message: UpdatePerpetualEventV1): UpdatePerpetualEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.UpdatePerpetualEventV1",
      value: UpdatePerpetualEventV1.encode(message).finish()
    };
  }
};
function createBaseTradingRewardsEventV1(): TradingRewardsEventV1 {
  return {
    tradingRewards: []
  };
}
export const TradingRewardsEventV1 = {
  typeUrl: "/dydxprotocol.indexer.events.TradingRewardsEventV1",
  encode(message: TradingRewardsEventV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.tradingRewards) {
      AddressTradingReward.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): TradingRewardsEventV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTradingRewardsEventV1();
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
  fromPartial(object: Partial<TradingRewardsEventV1>): TradingRewardsEventV1 {
    const message = createBaseTradingRewardsEventV1();
    message.tradingRewards = object.tradingRewards?.map(e => AddressTradingReward.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: TradingRewardsEventV1Amino): TradingRewardsEventV1 {
    const message = createBaseTradingRewardsEventV1();
    message.tradingRewards = object.trading_rewards?.map(e => AddressTradingReward.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: TradingRewardsEventV1): TradingRewardsEventV1Amino {
    const obj: any = {};
    if (message.tradingRewards) {
      obj.trading_rewards = message.tradingRewards.map(e => e ? AddressTradingReward.toAmino(e) : undefined);
    } else {
      obj.trading_rewards = [];
    }
    return obj;
  },
  fromAminoMsg(object: TradingRewardsEventV1AminoMsg): TradingRewardsEventV1 {
    return TradingRewardsEventV1.fromAmino(object.value);
  },
  fromProtoMsg(message: TradingRewardsEventV1ProtoMsg): TradingRewardsEventV1 {
    return TradingRewardsEventV1.decode(message.value);
  },
  toProto(message: TradingRewardsEventV1): Uint8Array {
    return TradingRewardsEventV1.encode(message).finish();
  },
  toProtoMsg(message: TradingRewardsEventV1): TradingRewardsEventV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.TradingRewardsEventV1",
      value: TradingRewardsEventV1.encode(message).finish()
    };
  }
};
function createBaseAddressTradingReward(): AddressTradingReward {
  return {
    owner: "",
    denomAmount: new Uint8Array()
  };
}
export const AddressTradingReward = {
  typeUrl: "/dydxprotocol.indexer.events.AddressTradingReward",
  encode(message: AddressTradingReward, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.owner !== "") {
      writer.uint32(10).string(message.owner);
    }
    if (message.denomAmount.length !== 0) {
      writer.uint32(18).bytes(message.denomAmount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): AddressTradingReward {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddressTradingReward();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.owner = reader.string();
          break;
        case 2:
          message.denomAmount = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<AddressTradingReward>): AddressTradingReward {
    const message = createBaseAddressTradingReward();
    message.owner = object.owner ?? "";
    message.denomAmount = object.denomAmount ?? new Uint8Array();
    return message;
  },
  fromAmino(object: AddressTradingRewardAmino): AddressTradingReward {
    const message = createBaseAddressTradingReward();
    if (object.owner !== undefined && object.owner !== null) {
      message.owner = object.owner;
    }
    if (object.denom_amount !== undefined && object.denom_amount !== null) {
      message.denomAmount = bytesFromBase64(object.denom_amount);
    }
    return message;
  },
  toAmino(message: AddressTradingReward): AddressTradingRewardAmino {
    const obj: any = {};
    obj.owner = message.owner;
    obj.denom_amount = message.denomAmount ? base64FromBytes(message.denomAmount) : undefined;
    return obj;
  },
  fromAminoMsg(object: AddressTradingRewardAminoMsg): AddressTradingReward {
    return AddressTradingReward.fromAmino(object.value);
  },
  fromProtoMsg(message: AddressTradingRewardProtoMsg): AddressTradingReward {
    return AddressTradingReward.decode(message.value);
  },
  toProto(message: AddressTradingReward): Uint8Array {
    return AddressTradingReward.encode(message).finish();
  },
  toProtoMsg(message: AddressTradingReward): AddressTradingRewardProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.events.AddressTradingReward",
      value: AddressTradingReward.encode(message).finish()
    };
  }
};