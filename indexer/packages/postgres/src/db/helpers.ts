import { logger } from '@dydxprotocol-indexer/base';
import { bytesToBigInt, getPositionIsLong } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  Asset, AssetPosition, LiquidityTier, MarketParam, MarketPrice,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import Long from 'long';

import { CURRENCY_DECIMAL_PRECISION, ONE_MILLION, QUOTE_CURRENCY_ATOMIC_RESOLUTION } from '../constants';
import { setBulkRowsForUpdate } from '../helpers/stores-helpers';
import { protocolPriceToHuman, quantumsToHuman, quantumsToHumanFixedString } from '../lib/protocol-translations';
import * as AssetPositionTable from '../stores/asset-position-table';
import {
  AssetCreateObject,
  FundingIndexMap,
  IsoString,
  LiquidityTiersColumns,
  LiquidityTiersCreateObject,
  MarketColumns,
  MarketCreateObject,
  MarketsMap,
  PerpetualMarketFromDatabase,
  PerpetualPositionFromDatabase,
  TransferFromDatabase,
  TransferType,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '../types';
import genesis from './genesis.json';

export interface AssetPositionCreateObjectWithId {
  id: string,
  subaccountId: string,
  assetId: string,
  size: string,
  isLong: boolean,
}

export interface SubaccountCreateObjectWithId {
  id: string,
  address: string,
  subaccountNumber: number,
  updatedAt: IsoString,
  updatedAtHeight: string,
}

/**
 * @description Gets the SQL to seed the `markets` table, using the `genesis.json` file
 * from the V4 network.
 * @returns SQL statement for seeding the `markets` table. The SQL statement will do
 * nothing if the rows in the `markets` table already exist.
 */
export function getSeedMarketsSql(): string {
  // Get `MarketParam` and `MarketPrice` objects from the genesis app state
  const marketParams: MarketParam[] = getMarketParamsFromGenesis();
  const marketPrices: MarketPrice[] = getMarketPricesFromGenesis();

  const marketCreateObjects: MarketCreateObject[] = [];
  marketParams.forEach((marketParam, index) => {
    marketCreateObjects.push(getMarketCreateObject(marketParam, marketPrices[index]));
  });

  const marketColumns = _.keys(marketCreateObjects[0]) as MarketColumns[];

  const marketRows: string[] = setBulkRowsForUpdate<MarketColumns>({
    objectArray: marketCreateObjects,
    columns: marketColumns,
    stringColumns: [
      MarketColumns.pair,
      MarketColumns.oraclePrice,
    ],
    numericColumns: [
      MarketColumns.id,
      MarketColumns.exponent,
      MarketColumns.minPriceChangePpm,
    ],
  });

  return `INSERT INTO MARKETS (${marketColumns.map((col) => `"${col}"`).join(',')})
          VALUES ${marketRows.map((market) => `(${market})`).join(', ')}
          ON CONFLICT DO NOTHING`;
}

/**
 * @description Gets the SQL to seed the `liquidity_tiers` table, using the `genesis.json` file
 * from the V4 network.
 * @returns SQL statement for seeding the `liquidity_tiers` table. The SQL statement will do
 * nothing if the rows in the `liquidity_tiers` table already exist.
 */
export function getSeedLiquidityTiersSql(): string {
  // Get `LiquidityTier` objects from the genesis app state
  const liquidityTiers: LiquidityTier[] = getLiquidityTiersFromGenesis();

  const liquidityTierCreateObjects:
  LiquidityTiersCreateObject[] = liquidityTiers.map((liquidityTier: LiquidityTier) => {
    return getLiquidityTiersCreateObject(liquidityTier);
  });

  const liquidityTierColumns = _.keys(liquidityTierCreateObjects[0]) as LiquidityTiersColumns[];

  const liquidityTierRows: string[] = setBulkRowsForUpdate<LiquidityTiersColumns>({
    objectArray: liquidityTierCreateObjects,
    columns: liquidityTierColumns,
    stringColumns: [
      LiquidityTiersColumns.name,
      LiquidityTiersColumns.basePositionNotional,
    ],
    numericColumns: [
      LiquidityTiersColumns.id,
      LiquidityTiersColumns.initialMarginPpm,
      LiquidityTiersColumns.maintenanceFractionPpm,
    ],
  });

  return `INSERT INTO LIQUIDITY_TIERS (${liquidityTierColumns.map((col) => `"${col}"`).join(',')})
          VALUES ${liquidityTierRows.map((liquidityTier) => `(${liquidityTier})`).join(', ')}
          ON CONFLICT DO NOTHING`;
}

export function getLiquidityTiersCreateObject(liquidityTier: LiquidityTier):
  LiquidityTiersCreateObject {
  return {
    id: liquidityTier.id,
    name: liquidityTier.name,
    initialMarginPpm: liquidityTier.initialMarginPpm.toString(),
    maintenanceFractionPpm: liquidityTier.maintenanceFractionPpm.toString(),
    basePositionNotional: quantumsToHuman(
      liquidityTier.basePositionNotional.toString(),
      QUOTE_CURRENCY_ATOMIC_RESOLUTION,
    ).toFixed(6),
  };
}

export function getAssetCreateObject(asset: Asset): AssetCreateObject {
  return {
    id: BigInt(asset.id).toString(),
    symbol: asset.symbol,
    atomicResolution: asset.atomicResolution,
    hasMarket: asset.hasMarket,
    marketId: asset.marketId,
  };
}

export function getAssetPositionCreateObject(
  subaccountId: string,
  assetPosition: AssetPosition,
  atomicResolution: number,
): AssetPositionCreateObjectWithId {
  return {
    id: AssetPositionTable.uuid(subaccountId, assetPosition.assetId.toString()),
    subaccountId,
    assetId: assetPosition.assetId.toString(),
    size: quantumsToHumanFixedString(bytesToBigInt(assetPosition.quantums).toString(),
      atomicResolution),
    isLong: getPositionIsLong(assetPosition),
  };
}

export function getMarketParamsFromGenesis(): MarketParam[] {
  const markets: MarketParam[] = genesis.app_state.prices.market_params.map(
    (genesisMarketParam): MarketParam => {
      return {
        ...genesisMarketParam,
        minExchanges: genesisMarketParam.min_exchanges,
        minPriceChangePpm: genesisMarketParam.min_price_change_ppm,
        exchangeConfigJson: '',
        id: genesisMarketParam.id,
      };
    },
  );

  return markets;
}

export function getMarketPricesFromGenesis(): MarketPrice[] {
  const marketPrices: MarketPrice[] = genesis.app_state.prices.market_prices.map(
    (genesisMarketPrice): MarketPrice => {
      return {
        ...genesisMarketPrice,
        id: genesisMarketPrice.id,
        exponent: genesisMarketPrice.exponent,
        price: Long.fromNumber(genesisMarketPrice.price),
      };
    },
  );
  return marketPrices;
}

/**
 * Gets LiquidityTiers from geneis.
 * @returns
 */
export function getLiquidityTiersFromGenesis(): LiquidityTier[] {
  const liquidityTiers: LiquidityTier[] = genesis.app_state.perpetuals.liquidity_tiers.map(
    (genesisLiquidityTier, index: number): LiquidityTier => {
      return {
        ...genesisLiquidityTier,
        basePositionNotional: Long.fromValue(genesisLiquidityTier.base_position_notional),
        initialMarginPpm: genesisLiquidityTier.initial_margin_ppm,
        maintenanceFractionPpm: genesisLiquidityTier.maintenance_fraction_ppm,
        impactNotional: Long.fromValue(genesisLiquidityTier.impact_notional),
        id: index,
      };
    },
  );
  return liquidityTiers;
}

/**
 * @description Given the initial `MarketParam` and `MarketPrice` objects, generate a
 * `MarketCreateObject`.
 * @param marketParam Initial `MarketParam` object.
 * @param marketPrice Initial `MarketPrice` object.
 * @returns `CreateMarketObject` corresponding to each `MarketParam`, `MarketPrice` pair passed in.
 * Note: This function assumes the passed in `MarketParam` and `MarketPrice` objects match, and
 * does no validation for this.
 */
function getMarketCreateObject(
  marketParam: MarketParam,
  marketPrice: MarketPrice,
): MarketCreateObject {
  return {
    id: marketParam.id,
    pair: marketParam.pair,
    exponent: marketParam.exponent,
    minPriceChangePpm: marketParam.minPriceChangePpm,
    oraclePrice: protocolPriceToHuman(marketPrice.price.toString(), marketPrice.exponent),
  };
}

/**
 * Converts a parts-per-million value to the string representation of the number. 1 ppm, or
 * parts-per-million is equal to 10^-6 (0.000001).
 * @param ppm Parts-per-million value.
 * @returns String representation of the parts-per-million value as a floating point number.
 */
export function ppmToString(ppm: number): string {
  return Big(ppm).div(1_000_000).toFixed(6);
}

/**
 * Calculates maintenance margin based on initial margin and maintenance fraction.
 * maintenance margin = initial margin * maintenance fraction
 * @param initialMarginPpm Initial margin in parts-per-million.
 * @param maintenanceFractionPpm Maintenance fraction in parts-per-million.
 * @returns Maintenance margin in parts-per-million.
 */
export function getMaintenanceMarginPpm(
  initialMarginPpm: number,
  maintenanceFractionPpm: number,
): number {
  return Big(initialMarginPpm).times(maintenanceFractionPpm).div(ONE_MILLION).toNumber();
}

/**
 * Computes the unsettled funding for a position.
 *
 * To compute the net USDC balance for a subaccount, sum the result of this function for all
 * open perpetual positions, and subtract the sum from the latest USDC asset position for
 * this subaccount.
 *
 * @param position
 * @param latestFundingIndex
 * @param lastUpdateFundingIndex
 */
export function getUnsettledFunding(
  position: PerpetualPositionFromDatabase,
  latestFundingIndexMap: FundingIndexMap,
  lastUpdateFundingIndexMap: FundingIndexMap,
): Big {
  return Big(position.size).times(
    latestFundingIndexMap[position.perpetualId].minus(
      lastUpdateFundingIndexMap[position.perpetualId],
    ),
  );
}

/**
 * Get unrealized pnl for a perpetual position. If the perpetual market is not found in the
 * markets map or the oracle price is not found in the market, return 0.
 *
 * @param position Perpetual position object from the database, or the updated
 * perpetual position subaccountKafkaObject.
 * @param perpetualMarketsMap Map of perpetual ids to perpetual market objects from the database.
 * @returns Unrealized pnl of the position.
 */
export function getUnrealizedPnl(
  position: PerpetualPositionFromDatabase | UpdatedPerpetualPositionSubaccountKafkaObject,
  perpetualMarket: PerpetualMarketFromDatabase,
  marketsMap: MarketsMap,
): string {
  if (marketsMap[perpetualMarket.marketId] === undefined) {
    logger.error({
      at: 'getUnrealizedPnl',
      message: 'Market is undefined',
      marketId: perpetualMarket.marketId,
    });
    return Big(0).toFixed(CURRENCY_DECIMAL_PRECISION);
  }
  if (marketsMap[perpetualMarket.marketId].oraclePrice === undefined) {
    logger.error({
      at: 'getUnrealizedPnl',
      message: 'Oracle price is undefined for market',
      marketId: perpetualMarket.marketId,
    });
    return Big(0).toFixed(CURRENCY_DECIMAL_PRECISION);
  }
  return (
    Big(position.size).times(
      Big(marketsMap[perpetualMarket.marketId].oraclePrice!).minus(position.entryPrice),
    )
  ).toFixed(CURRENCY_DECIMAL_PRECISION);
}

/**
 * Gets the transfer type for a subaccount.
 *
 * If sender/recipient are both subaccounts, then it is a transfer_in/transfer_out.
 * If sender/recipient are wallet addresses, then it is a deposit/withdrawal.
 *
 * @param transfer
 * @param subaccountId
 */
export function getTransferType(
  transfer: TransferFromDatabase,
  subaccountId: string,
): TransferType {
  if (transfer.senderSubaccountId === subaccountId) {
    if (transfer.recipientSubaccountId) {
      return TransferType.TRANSFER_OUT;
    } else {
      return TransferType.WITHDRAWAL;
    }
  } else if (transfer.recipientSubaccountId === subaccountId) {
    if (transfer.senderSubaccountId) {
      return TransferType.TRANSFER_IN;
    } else {
      return TransferType.DEPOSIT;
    }
  }
  throw new Error(`Transfer ${transfer.id} does not involve subaccount ${subaccountId}`);
}
