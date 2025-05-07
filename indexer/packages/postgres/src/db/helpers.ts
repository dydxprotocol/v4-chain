import { logger } from '@dydxprotocol-indexer/base';
import Big from 'big.js';

import { ONE_MILLION } from '../constants';
import {
  FundingIndexMap, MarketFromDatabase,
  PerpetualMarketFromDatabase,
  PerpetualPositionFromDatabase,
  TransferFromDatabase,
  TransferType,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '../types';

/**
 * Converts a parts-per-million value to the string representation of the number. 1 ppm, or
 * parts-per-million is equal to 10^-6 (0.000001).
 * @param ppm Parts-per-million value.
 * @returns String representation of the parts-per-million value as a floating point number.
 */
export function ppmToString(ppm: number): string {
  return Big(ppm).div(1_000_000).toFixed();
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
 * open perpetual positions, and add it to the latest USDC asset position for
 * this subaccount.
 *
 * When funding index is increasing, shorts get paid & unsettled funding for shorts should
 * be positive, vice versa for longs.
 * When funding index is decreasing, longs get paid & unsettled funding for longs should
 * be positive, vice versa for shorts.
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
    lastUpdateFundingIndexMap[position.perpetualId].minus(
      latestFundingIndexMap[position.perpetualId],
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
 * @param market Market object from the database.
 * @returns Unrealized pnl of the position.
 */
export function getUnrealizedPnl(
  position: PerpetualPositionFromDatabase | UpdatedPerpetualPositionSubaccountKafkaObject,
  perpetualMarket: PerpetualMarketFromDatabase,
  market: MarketFromDatabase,
): string {
  if (market.oraclePrice === undefined) {
    logger.error({
      at: 'getUnrealizedPnl',
      message: 'Oracle price is undefined for market',
      marketId: perpetualMarket.marketId,
    });
    return Big(0).toFixed();
  }
  return (
    Big(position.size).times(
      Big(market.oraclePrice!).minus(position.entryPrice),
    )
  ).toFixed();
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
