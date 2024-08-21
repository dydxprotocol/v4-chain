import { stats } from '@dydxprotocol-indexer/base';
import {
  DEFAULT_POSTGRES_OPTIONS,
  Ordering, PaginationFromDatabase,
  PnlTicksFromDatabase,
  PnlTicksTable,
  QueryableField,
  perpetualMarketRefresher,
  PerpetualMarketFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import express from 'express';
import {
  Controller, Get, Route,
} from 'tsoa';
import _ from 'lodash';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { aggregatePnlTicks, handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { pnlTicksToResponseObject } from '../../../request-helpers/request-transformer';
import { MegavaultHistoricalPnlResponse, VaultsHistoricalPnlResponse, VaultHistoricalPnl } from '../../../types';

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
      megavaultsPnl: Array.from(aggregatedPnlTicks.values()).map(
        (pnlTick: PnlTicksFromDatabase) => {
          return pnlTicksToResponseObject(pnlTick);
        }),
    };
  }

  @Get('/v1/vaults/historicalPnl')
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
          throw new Error(`Vault clob pair id ${vaultSubaccounts[subaccountId]} does not correspond to a perpetual market.`)
        }

        return {
          ticker: market.ticker,
          historicalPnl: pnlTicks,
        }})
      .values()
      .value();

    return {
      vaultsPnl: groupedVaultPnlTicks,
    };
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
      const response: MegavaultHistoricalPnlResponse = await controllers.getMegavaultHistoricalPnl();
       return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaulController GET /megavault/historicalPnl',
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
});

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
  const vaultSubaccountIds: string[] = config.EXPERIMENT_VAULTS.split(',');;
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
