import { bytesToBigInt, getPositionIsLong } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  Asset, AssetPosition, LiquidityTier, MarketParam, MarketPrice,
} from '@dydxprotocol-indexer/v4-protos';

import { QUOTE_CURRENCY_ATOMIC_RESOLUTION } from '../../src';
import { SubaccountCreateObjectWithId } from '../../src/db/helpers';
import { protocolPriceToHuman, quantumsToHuman, quantumsToHumanFixedString } from '../../src/lib/protocol-translations';
import {
  AssetFromDatabase,
  AssetPositionCreateObject,
  AssetPositionFromDatabase,
  LiquidityTiersFromDatabase,
  MarketFromDatabase,
  SubaccountFromDatabase,
} from '../../src/types';

export function expectMarketParamAndPrice(
  marketFromDb: MarketFromDatabase,
  marketParam: MarketParam,
  marketPrice: MarketPrice,
): void {
  expect(marketFromDb).toEqual(expect.objectContaining({
    id: marketParam.id,
    pair: marketParam.pair,
    exponent: marketParam.exponent,
    minPriceChangePpm: marketParam.minPriceChangePpm,
    oraclePrice: protocolPriceToHuman(marketPrice.price.toString(), marketPrice.exponent),
  }));
}

export function expectAsset(
  assetFromDb: AssetFromDatabase,
  asset: Asset,
): void {
  expect(assetFromDb).toEqual(expect.objectContaining({
    atomicResolution: asset.atomicResolution,
    symbol: asset.symbol,
    hasMarket: asset.hasMarket,
    id: BigInt(asset.id).toString(),
    marketId: asset.marketId,
  }));
}

export function expectLiquidityTier(
  liquidityTierFromDb: LiquidityTiersFromDatabase,
  liquidityTier: LiquidityTier,
): void {
  expect(liquidityTierFromDb).toEqual(expect.objectContaining({
    id: liquidityTier.id,
    name: liquidityTier.name,
    initialMarginPpm: liquidityTier.initialMarginPpm.toString(),
    maintenanceFractionPpm: liquidityTier.maintenanceFractionPpm.toString(),
    basePositionNotional: quantumsToHuman(
      liquidityTier.basePositionNotional.toString(),
      QUOTE_CURRENCY_ATOMIC_RESOLUTION,
    ).toFixed(6),
  }));
}

export function expectSubaccount(
  subaccountFromDb: SubaccountFromDatabase,
  subaccount: SubaccountCreateObjectWithId,
): void {
  expect(subaccountFromDb).toEqual(expect.objectContaining({
    id: subaccount.id,
    address: subaccount.address,
    subaccountNumber: subaccount.subaccountNumber,
  }));
}

export function expectAssetPosition(
  assetPositionFromDb: AssetPositionFromDatabase,
  assetPosition: AssetPosition,
  atomicResolution: number,
): void {
  expect(assetPositionFromDb).toEqual(expect.objectContaining({
    assetId: assetPosition.assetId.toString(),
    size: quantumsToHumanFixedString(
      bytesToBigInt(assetPosition.quantums).toString(),
      atomicResolution),
    isLong: getPositionIsLong(assetPosition),
  }));
}

export function expectAssetPositionCreateObject(
  assetPositionCreateObject: AssetPositionCreateObject,
  assetPosition: AssetPosition,
  atomicResolution: number,
): void {
  expect(assetPositionCreateObject).toEqual(expect.objectContaining({
    assetId: assetPosition.assetId.toString(),
    size: quantumsToHumanFixedString(
      bytesToBigInt(assetPosition.quantums).toString(),
      atomicResolution),
    isLong: getPositionIsLong(assetPosition),
  }));
}
