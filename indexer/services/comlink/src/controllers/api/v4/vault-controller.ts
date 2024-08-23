import { stats } from '@dydxprotocol-indexer/base';
import {
  DEFAULT_POSTGRES_OPTIONS,
  Ordering, PaginationFromDatabase,
  PnlTicksFromDatabase,
  PnlTicksTable,
  QueryableField,
  perpetualMarketRefresher,
  PerpetualMarketFromDatabase,
  USDC_ASSET_ID,
  FundingIndexMap,
  AssetPositionFromDatabase,
  PerpetualPositionFromDatabase,
  SubaccountFromDatabase,
  AssetColumns,
  MarketColumns,
  MarketsMap,
  BlockTable,
  MarketTable,
  AssetPositionTable,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  AssetTable,
  SubaccountTable,
  AssetFromDatabase,
  MarketFromDatabase,
  BlockFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import _ from 'lodash';
import {
  Controller, Get, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { adjustUSDCAssetPosition, aggregatePnlTicks, calculateEquityAndFreeCollateral, getFundingIndexMaps, getPerpetualPositionsWithUpdatedFunding, getTotalUnsettledFunding, handleControllerError, initializePerpetualPositionsWithFunding } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { assetPositionToResponseObject, perpetualPositionToResponseObject, pnlTicksToResponseObject } from '../../../request-helpers/request-transformer';
import { AssetPositionsMap, MegavaultHistoricalPnlResponse, VaultsHistoricalPnlResponse, VaultHistoricalPnl, PerpetualPositionResponseObject, PerpetualPositionWithFunding, AssetPositionResponseObject, VaultPosition, AssetById, MegavaultPositionResponse } from '../../../types';
import Big from 'big.js';

const router: express.Router = express.Router();
const controllerName: string = 'vault-controller';

// TODO(TRA-570): Placeholder interface for mapping of vault subaccounts to tickers until vaults
// table is added.
interface VaultMapping {
  [subaccountId: string]: string;
}

@Route('vault/v1')
class VaultController extends Controller {
  @Get('/megavault/historicalPnl')
  async getMegavaultHistoricalPnl(): Promise<MegavaultHistoricalPnlResponse> {
    const vaultPnlTicks: PnlTicksFromDatabase[] = await getVaultSubaccountPnlTicks();

    // aggregate pnlTicks for all vault subaccounts grouped by blockHeight
    const aggregatedPnlTicks: Map<number, PnlTicksFromDatabase> = aggregatePnlTicks(vaultPnlTicks);

    return {
      megavaultPnl: Array.from(aggregatedPnlTicks.values()).map(
        (pnlTick: PnlTicksFromDatabase) => {
          return pnlTicksToResponseObject(pnlTick);
        }),
    };
  }

  @Get('/vaults/historicalPnl')
  async getVaultsHistoricalPnl(): Promise<VaultsHistoricalPnlResponse> {
    const vaultSubaccounts: VaultMapping = getVaultSubaccountsFromConfig();
    const vaultPnlTicks: PnlTicksFromDatabase[] = await getVaultSubaccountPnlTicks();

    const groupedVaultPnlTicks: VaultHistoricalPnl[] = _(vaultPnlTicks)
      .groupBy('subaccountId')
      .mapValues((pnlTicks: PnlTicksFromDatabase[], subaccountId: string): VaultHistoricalPnl => {
        const market: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
          .getPerpetualMarketFromClobPairId(
            vaultSubaccounts[subaccountId],
          );

        if (market === undefined) {
          throw new Error(`Vault clob pair id ${vaultSubaccounts[subaccountId]} does not correspond to a perpetual market.`);
        }

        return {
          ticker: market.ticker,
          historicalPnl: pnlTicks,
        };
      })
      .values()
      .value();

    return {
      vaultsPnl: groupedVaultPnlTicks,
    };
  }

  @Get('/v1/megavault/positions')
  async getMegavaultPositions(): Promise<MegavaultPositionResponse> {
    const vaultSubaccounts: VaultMapping = getVaultSubaccountsFromConfig();
    const vaultSubaccountIds: string[] = _.keys(vaultSubaccounts);

    const [
      subaccounts,
      assets,
      perpetualPositions,
      assetPositions,
      markets,
      latestBlock,
    ]: [
      SubaccountFromDatabase[],
      AssetFromDatabase[],
      PerpetualPositionFromDatabase[],
      AssetPositionFromDatabase[],
      MarketFromDatabase[],
      BlockFromDatabase | undefined,
    ] = await Promise.all([
      SubaccountTable.findAll(
        {
          id: vaultSubaccountIds,
        },
        [],
      ),
      AssetTable.findAll(
        {},
        [],
      ),
      PerpetualPositionTable.findAll(
        {
          subaccountId: vaultSubaccountIds,
          status: [PerpetualPositionStatus.OPEN],
        },
        [],
      ),
      AssetPositionTable.findAll(
        {
          subaccountId: vaultSubaccountIds,
          assetId: [USDC_ASSET_ID],
        },
        [],
      ),
      MarketTable.findAll(
        {},
        [],
      ),
      BlockTable.getLatest(),
    ]);

    const assetPositionsBySubaccount:
    { [subaccountId: string]: AssetPositionFromDatabase[] } = _.groupBy(
      assetPositions,
      'subaccountId',
    );
    const perpetualPositionsBySubaccount:
    { [subaccountId: string]: PerpetualPositionFromDatabase[] } = _.groupBy(
      perpetualPositions,
      'subaccountId',
    );
    const marketIdToMarket: MarketsMap = _.keyBy(
      markets,
      MarketColumns.id,
    );
    const assetIdToAsset: AssetById = _.keyBy(
      assets,
      AssetColumns.id,
    );

    const vaultPositions: VaultPosition[] = await Promise.all(
      subaccounts.map(async (subaccount: SubaccountFromDatabase) => {
        const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
          .getPerpetualMarketFromClobPairId(vaultSubaccounts[subaccount.id]);
        if (perpetualMarket === undefined) {
          throw new Error(`Vault clob pair id ${vaultSubaccounts[subaccount.id]} does not correspond to a perpetual market.`)
        }

        const assetPositions: AssetPositionFromDatabase[] = assetPositionsBySubaccount[subaccount.id];
        const perpetualPositions: PerpetualPositionFromDatabase[] = perpetualPositionsBySubaccount[subaccount.id];

        if (assetPositions.length == 0) {
          throw new Error(`Vault for ticker ${perpetualMarket.ticker} has no USDC asset positions`);
        }
        if (assetPositions.length > 1) {
          throw new Error(`Vault for ticker ${perpetualMarket.ticker} has more than 1 USDC asset position`);
        }
        if (perpetualPositions !== undefined && perpetualPositions.length > 1) {
          throw new Error(`Vault for ticker ${perpetualMarket.ticker} has more than 1 perpetual position`);
        }

        const usdcAssetPosition: AssetPositionFromDatabase = assetPositions[0];
        const usdcAssetPositionResponse: AssetPositionResponseObject = assetPositionToResponseObject(
          usdcAssetPosition,
          assetIdToAsset,
          subaccount.subaccountNumber,
        );
        const usdcAssetPositionSize: Big = Big(usdcAssetPosition.isLong ? usdcAssetPosition.size : -usdcAssetPosition.size);
        const perpetualPositionsWithFunding: PerpetualPositionWithFunding[] = initializePerpetualPositionsWithFunding(
          perpetualPositions
        );

        // If there are no perpetual positions for the vault subaccount, equity = USDC asset position size
        if (perpetualPositions === undefined || perpetualPositions.length === 0) {
          return {
            ticker: perpetualMarket.ticker,
            assetPosition: usdcAssetPositionResponse,
            perpetualPosition: undefined,
            equity: usdcAssetPositionSize.toString(),
          }
        }

        const {
          lastUpdatedFundingIndexMap,
          latestFundingIndexMap,
        }: {
          lastUpdatedFundingIndexMap: FundingIndexMap,
          latestFundingIndexMap: FundingIndexMap,
        } = await getFundingIndexMaps(subaccount, latestBlock);

        const unsettledFunding: Big = getTotalUnsettledFunding(
          perpetualPositions,
          latestFundingIndexMap,
          lastUpdatedFundingIndexMap,
        );

        const adjustedPerpetualPositions: PerpetualPositionWithFunding[] = getPerpetualPositionsWithUpdatedFunding(
          perpetualPositionsWithFunding,
          latestFundingIndexMap,
          lastUpdatedFundingIndexMap,
        );
        const perpetutalPositionResponse: PerpetualPositionResponseObject = perpetualPositionToResponseObject(
          adjustedPerpetualPositions[0],
          perpetualMarketRefresher.getPerpetualMarketsMap(),
          marketIdToMarket,
          subaccount.subaccountNumber,
        );

        const {
          assetPositionsMap: adjustedAssetPositionsMap,
          adjustedUSDCAssetPositionSize,
        }: {
          assetPositionsMap: AssetPositionsMap,
          adjustedUSDCAssetPositionSize: string,
        } = adjustUSDCAssetPosition({
          [assetIdToAsset[USDC_ASSET_ID].symbol]: usdcAssetPositionResponse,
        }, unsettledFunding);

        const {
          equity,
        }: {
          equity: string,
          freeCollateral: string,
        } = calculateEquityAndFreeCollateral(
          adjustedPerpetualPositions,
          perpetualMarketRefresher.getPerpetualMarketsMap(),
          marketIdToMarket,
          adjustedUSDCAssetPositionSize,
        );

        return {
          ticker: perpetualMarket.ticker,
          assetPosition: adjustedAssetPositionsMap[assetIdToAsset[USDC_ASSET_ID].symbol],
          perpetualPosition: perpetutalPositionResponse,
          equity: equity,
        };
      })
    );

    return {
      positions: vaultPositions,
    }
  }
}

router.get(
  '/v1/megavault/historicalPnl',
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    try {
      const controllers: VaultController = new VaultController();
      const response: MegavaultHistoricalPnlResponse = await controllers
        .getMegavaultHistoricalPnl();
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultController GET /megavault/historicalPnl',
        'Megavault Historical Pnl error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_megavault_historical_pnl.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/v1/vaults/historicalPnl',
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    try {
      const controllers: VaultController = new VaultController();
      const response: VaultsHistoricalPnlResponse = await controllers.getVaultsHistoricalPnl();
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultHistoricalPnlController GET /vaults/historicalPnl',
        'Vaults Historical Pnl error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_vaults_historical_pnl.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/v1/megavault/positions',
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const controllers: VaultController = new VaultController();
      const response: MegavaultPositionResponse = await controllers.getMegavaultPositions();
       return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultController GET /megavault/positions',
        'Megavault Positions error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_megavault_positions.timing`,
        Date.now() - start,
      );
    }
});

async function getVaultSubaccountPnlTicks(): Promise<PnlTicksFromDatabase[]> {
  const subVaultSubaccountIds: string[] = _.keys(getVaultSubaccountsFromConfig());
  const {
    results: pnlTicks,
  }: PaginationFromDatabase<PnlTicksFromDatabase> = await
  PnlTicksTable.findAll(
    {
      subaccountId: subVaultSubaccountIds,
      // TODO(TRA-571): Configure limits based on hourly vs daily resolution and # of vaults.
      limit: config.API_LIMIT_V4,
    },
    [QueryableField.LIMIT],
    {
      ...DEFAULT_POSTGRES_OPTIONS,
      orderBy: [[QueryableField.BLOCK_HEIGHT, Ordering.DESC]],
    },
  );
  return pnlTicks;
}

// TODO(TRA-570): Placeholder for getting vault subaccount ids until vault table is added.
function getVaultSubaccountsFromConfig(): VaultMapping {
  const vaultSubaccountIds: string[] = config.EXPERIMENT_VAULTS.split(',');
  const vaultClobPairIds: string[] = config.EXPERIMENT_VAULT_MARKETS.split(',');
  if (vaultSubaccountIds.length !== vaultClobPairIds.length) {
    throw new Error('Expected number of vaults to match number of markets');
  }
  return _.zipObject(
    vaultSubaccountIds,
    vaultClobPairIds,
  );
}

export default router;
