import { stats } from '@dydxprotocol-indexer/base';
import {
  PnlTicksFromDatabase,
  PnlTicksTable,
  perpetualMarketRefresher,
  PerpetualMarketFromDatabase,
  USDC_ASSET_ID,
  FundingIndexMap,
  AssetPositionFromDatabase,
  PerpetualPositionFromDatabase,
  SubaccountFromDatabase,
  AssetColumns,
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
  FundingIndexUpdatesTable,
  PnlTickInterval,
  VaultTable,
  VaultFromDatabase,
  MEGAVAULT_SUBACCOUNT_ID,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import _ from 'lodash';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import {
  aggregatePnlTicks,
  getSubaccountResponse,
  handleControllerError,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { pnlTicksToResponseObject } from '../../../request-helpers/request-transformer';
import {
  MegavaultHistoricalPnlResponse,
  VaultsHistoricalPnlResponse,
  VaultHistoricalPnl,
  VaultPosition,
  AssetById,
  MegavaultPositionResponse,
  SubaccountResponseObject,
  MegavaultHistoricalPnlRequest,
  VaultsHistoricalPnlRequest,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'vault-controller';

interface VaultMapping {
  [subaccountId: string]: string,
}

@Route('vault/v1')
class VaultController extends Controller {
  @Get('/megavault/historicalPnl')
  async getMegavaultHistoricalPnl(
    @Query() resolution?: PnlTickInterval,
  ): Promise<MegavaultHistoricalPnlResponse> {
    const vaultSubaccounts: VaultMapping = await getVaultMapping();
    const vaultSubaccountIdsWithMainSubaccount: string[] = _
      .keys(vaultSubaccounts)
      .concat([MEGAVAULT_SUBACCOUNT_ID]);
    const [
      vaultPnlTicks,
      vaultPositions,
      latestBlock,
      mainSubaccountEquity,
    ] : [
      PnlTicksFromDatabase[],
      Map<string, VaultPosition>,
      BlockFromDatabase,
      string,
    ] = await Promise.all([
      getVaultSubaccountPnlTicks(vaultSubaccountIdsWithMainSubaccount, resolution),
      getVaultPositions(vaultSubaccounts),
      BlockTable.getLatest(),
      getMainSubaccountEquity(),
    ]);

    // aggregate pnlTicks for all vault subaccounts grouped by blockHeight
    const aggregatedPnlTicks: Map<number, PnlTicksFromDatabase> = aggregatePnlTicks(vaultPnlTicks);

    const currentEquity: string = Array.from(vaultPositions.values())
      .map((position: VaultPosition): string => {
        return position.equity;
      }).reduce((acc: string, curr: string): string => {
        return (Big(acc).add(Big(curr))).toFixed();
      }, mainSubaccountEquity);
    const pnlTicksWithCurrentTick: PnlTicksFromDatabase[] = getPnlTicksWithCurrentTick(
      currentEquity,
      Array.from(aggregatedPnlTicks.values()),
      latestBlock,
    );

    return {
      megavaultPnl: _.sortBy(pnlTicksWithCurrentTick, 'blockTime').map(
        (pnlTick: PnlTicksFromDatabase) => {
          return pnlTicksToResponseObject(pnlTick);
        }),
    };
  }

  @Get('/vaults/historicalPnl')
  async getVaultsHistoricalPnl(
    @Query() resolution?: PnlTickInterval,
  ): Promise<VaultsHistoricalPnlResponse> {
    const vaultSubaccounts: VaultMapping = await getVaultMapping();
    const [
      vaultPnlTicks,
      vaultPositions,
      latestBlock,
    ] : [
      PnlTicksFromDatabase[],
      Map<string, VaultPosition>,
      BlockFromDatabase,
    ] = await Promise.all([
      getVaultSubaccountPnlTicks(_.keys(vaultSubaccounts), resolution),
      getVaultPositions(vaultSubaccounts),
      BlockTable.getLatest(),
    ]);

    const groupedVaultPnlTicks: VaultHistoricalPnl[] = _(vaultPnlTicks)
      .groupBy('subaccountId')
      .mapValues((pnlTicks: PnlTicksFromDatabase[], subaccountId: string): VaultHistoricalPnl => {
        const market: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
          .getPerpetualMarketFromClobPairId(
            vaultSubaccounts[subaccountId],
          );

        if (market === undefined) {
          throw new Error(
            `Vault clob pair id ${vaultSubaccounts[subaccountId]} does not correspond to ` +
            'a perpetual market.');
        }

        const vaultPosition: VaultPosition | undefined = vaultPositions.get(subaccountId);
        const currentEquity: string = vaultPosition === undefined ? '0' : vaultPosition.equity;
        const pnlTicksWithCurrentTick: PnlTicksFromDatabase[] = getPnlTicksWithCurrentTick(
          currentEquity,
          pnlTicks,
          latestBlock,
        );

        return {
          ticker: market.ticker,
          historicalPnl: pnlTicksWithCurrentTick,
        };
      })
      .values()
      .value();

    return {
      vaultsPnl: _.sortBy(groupedVaultPnlTicks, 'ticker'),
    };
  }

  @Get('/megavault/positions')
  async getMegavaultPositions(): Promise<MegavaultPositionResponse> {
    const vaultSubaccounts: VaultMapping = await getVaultMapping();

    const vaultPositions: Map<string, VaultPosition> = await getVaultPositions(vaultSubaccounts);

    return {
      positions: _.sortBy(Array.from(vaultPositions.values()), 'ticker'),
    };
  }
}

router.get(
  '/v1/megavault/historicalPnl',
  ...checkSchema({
    resolution: {
      in: 'query',
      isIn: {
        options: [Object.values(PnlTickInterval)],
        errorMessage: `type must be one of ${Object.values(PnlTickInterval)}`,
      },
      optional: true,
    },
  }),
  handleValidationErrors,
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      resolution,
    }: MegavaultHistoricalPnlRequest = matchedData(req) as MegavaultHistoricalPnlRequest;

    try {
      const controllers: VaultController = new VaultController();
      const response: MegavaultHistoricalPnlResponse = await controllers
        .getMegavaultHistoricalPnl(resolution);
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
  ...checkSchema({
    resolution: {
      in: 'query',
      isIn: {
        options: [Object.values(PnlTickInterval)],
        errorMessage: `type must be one of ${Object.values(PnlTickInterval)}`,
      },
      optional: true,
    },
  }),
  handleValidationErrors,
  rateLimiterMiddleware(getReqRateLimiter),
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      resolution,
    }: VaultsHistoricalPnlRequest = matchedData(req) as VaultsHistoricalPnlRequest;

    try {
      const controllers: VaultController = new VaultController();
      const response: VaultsHistoricalPnlResponse = await controllers
        .getVaultsHistoricalPnl(resolution);
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

async function getVaultSubaccountPnlTicks(
  vaultSubaccountIds: string[],
  resolution?: PnlTickInterval,
): Promise<PnlTicksFromDatabase[]> {
  if (vaultSubaccountIds.length === 0) {
    return [];
  }
  let pnlTickInterval: PnlTickInterval;
  if (resolution === undefined) {
    pnlTickInterval = PnlTickInterval.day;
  } else {
    pnlTickInterval = resolution;
  }

  const pnlTicks: PnlTicksFromDatabase[] = await PnlTicksTable.getPnlTicksAtIntervals(
    pnlTickInterval,
    config.VAULT_PNL_HISTORY_DAYS * 24 * 60 * 60,
    vaultSubaccountIds,
  );

  return pnlTicks;
}

async function getVaultPositions(
  vaultSubaccounts: VaultMapping,
): Promise<Map<string, VaultPosition>> {
  const vaultSubaccountIds: string[] = _.keys(vaultSubaccounts);
  if (vaultSubaccountIds.length === 0) {
    return new Map();
  }

  const [
    subaccounts,
    assets,
    openPerpetualPositions,
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

  const latestFundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
    .findFundingIndexMap(
      latestBlock.blockHeight,
    );
  const assetPositionsBySubaccount:
  { [subaccountId: string]: AssetPositionFromDatabase[] } = _.groupBy(
    assetPositions,
    'subaccountId',
  );
  const openPerpetualPositionsBySubaccount:
  { [subaccountId: string]: PerpetualPositionFromDatabase[] } = _.groupBy(
    openPerpetualPositions,
    'subaccountId',
  );
  const assetIdToAsset: AssetById = _.keyBy(
    assets,
    AssetColumns.id,
  );

  const vaultPositionsAndSubaccountId: {
    position: VaultPosition,
    subaccountId: string,
  }[] = await Promise.all(
    subaccounts.map(async (subaccount: SubaccountFromDatabase) => {
      const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
        .getPerpetualMarketFromClobPairId(vaultSubaccounts[subaccount.id]);
      if (perpetualMarket === undefined) {
        throw new Error(
          `Vault clob pair id ${vaultSubaccounts[subaccount.id]} does not correspond to a ` +
          'perpetual market.');
      }
      const lastUpdatedFundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
        .findFundingIndexMap(
          subaccount.updatedAtHeight,
        );

      const subaccountResponse: SubaccountResponseObject = getSubaccountResponse(
        subaccount,
        openPerpetualPositionsBySubaccount[subaccount.id] || [],
        assetPositionsBySubaccount[subaccount.id] || [],
        assets,
        markets,
        perpetualMarketRefresher.getPerpetualMarketsMap(),
        latestBlock.blockHeight,
        latestFundingIndexMap,
        lastUpdatedFundingIndexMap,
      );

      return {
        position: {
          ticker: perpetualMarket.ticker,
          assetPosition: subaccountResponse.assetPositions[
            assetIdToAsset[USDC_ASSET_ID].symbol
          ],
          perpetualPosition: subaccountResponse.openPerpetualPositions[
            perpetualMarket.ticker
          ] || undefined,
          equity: subaccountResponse.equity,
        },
        subaccountId: subaccount.id,
      };
    }),
  );

  return new Map(vaultPositionsAndSubaccountId.map(
    (obj: { position: VaultPosition, subaccountId: string }) : [string, VaultPosition] => {
      return [
        obj.subaccountId,
        obj.position,
      ];
    },
  ));
}

async function getMainSubaccountEquity(): Promise<string> {
  // Main vault subaccount should only ever hold a USDC and never any perpetuals.
  const usdcBalance: {[subaccountId: string]: Big} = await AssetPositionTable
    .findUsdcPositionForSubaccounts(
      [MEGAVAULT_SUBACCOUNT_ID],
    );
  return usdcBalance[MEGAVAULT_SUBACCOUNT_ID]?.toFixed() || '0';
}

function getPnlTicksWithCurrentTick(
  equity: string,
  pnlTicks: PnlTicksFromDatabase[],
  latestBlock: BlockFromDatabase,
): PnlTicksFromDatabase[] {
  if (pnlTicks.length === 0) {
    return [];
  }
  const currentTick: PnlTicksFromDatabase = {
    ...pnlTicks[pnlTicks.length - 1],
    equity,
    blockHeight: latestBlock.blockHeight,
    blockTime: latestBlock.time,
    createdAt: latestBlock.time,
  };
  return pnlTicks.concat([currentTick]);
}

async function getVaultMapping(): Promise<VaultMapping> {
  const vaults: VaultFromDatabase[] = await VaultTable.findAll(
    {},
    [],
    {},
  );
  return _.zipObject(
    vaults.map((vault: VaultFromDatabase): string => {
      return SubaccountTable.uuid(vault.address, 0);
    }),
    vaults.map((vault: VaultFromDatabase): string => {
      return vault.clobPairId;
    }),
  );
}

export default router;
