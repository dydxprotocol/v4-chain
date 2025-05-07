/* eslint-disable no-console */
import {
  FillTable,
  FundingIndexMap,
  FundingIndexUpdatesTable,
  OpenSizeWithFundingIndex,
  OrderedFillsWithFundingIndices,
  Ordering,
  perpetualMarketRefresher,
  PnlTicksColumns,
  PnlTicksFromDatabase,
  PnlTicksTable,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import _ from 'lodash';

/**
 * Get unsettled funding for a subaccount at effectiveBeforeHeight.
 *
 * @param subaccountId
 * @param effectiveBeforeOrAtHeight
 */
export async function getUnsettledFunding(
  subaccountId: string,
  effectiveBeforeOrAtHeight: string,
): Promise<Big> {
  await perpetualMarketRefresher.updatePerpetualMarkets();
  const openSizeWithFundingIndices: OpenSizeWithFundingIndex[] = await
  FillTable.getOpenSizeWithFundingIndex(subaccountId, effectiveBeforeOrAtHeight);

  const lastFundingIndexMap: FundingIndexMap = await
  FundingIndexUpdatesTable.findFundingIndexMap(
    effectiveBeforeOrAtHeight,
    {
      readReplica: true,
    },
  );
  const getClobPairId = (perpetualId: string): string => {
    return perpetualMarketRefresher.getPerpetualMarketFromId(perpetualId)!.clobPairId;
  };

  const mappedLastFundingIndexMap: { [clobPairId: string]: Big } = _.mapKeys(
    lastFundingIndexMap,
    (value: Big, perpetualId: string) => {
      const clobPairId: string = getClobPairId(perpetualId);
      return clobPairId;
    });

  return _.reduce(
    openSizeWithFundingIndices,
    (totalUnsettledFunding: Big, item: OpenSizeWithFundingIndex) => {
      const {
        clobPairId,
        openSize,
        fundingIndex,
      }: {
        clobPairId: string,
        openSize: string,
        fundingIndex: string,
      } = item;
      const fundingIndexDiff: Big = Big(mappedLastFundingIndexMap[clobPairId]).minus(fundingIndex);
      const unsettledFunding: Big = Big(openSize).mul(fundingIndexDiff);
      return totalUnsettledFunding.add(unsettledFunding);
    },
    new Big(0),
  );
}

/**
 * Get realized funding for a subaccount at effectiveBeforeHeight.
 * Takes into account all clob pair ids.
 *
 * @param subaccountId
 * @param effectiveBeforeOrAtHeight
 */
export async function getRealizedFunding(
  subaccountId: string,
  effectiveBeforeOrAtHeight: string,
): Promise<Big> {
  const clobPairs: string[] = await
  FillTable.getClobPairs(subaccountId, effectiveBeforeOrAtHeight);

  let totalSettledFunding: Big = new Big(0);

  for (const clobPairId of clobPairs) {
    const orderedFillsWithFundingIndices: OrderedFillsWithFundingIndices[] = await
    FillTable.getOrderedFillsWithFundingIndices(
      clobPairId,
      subaccountId,
      effectiveBeforeOrAtHeight,
    );
    const settledFunding: Big = FillTable.getSettledFunding(orderedFillsWithFundingIndices);
    totalSettledFunding = totalSettledFunding.add(settledFunding);
  }
  return totalSettledFunding;
}

/**
 * Get totalPnl for a subaccount at effectiveBeforeHeight.
 *
 * totalPnl = Pnl of fills + total value of open positions - realized funding
 * - unrealized funding - fees paid.
 *
 * TODO(CORE-512): Add info/resources around Pnl validation.
 * Doc: https://www.notion.so/dydx/Pnl-Validation-f0eaf64149a84bcdbe26d194350a5de6
 *
 * @param subaccountId
 * @param effectiveBeforeOrAtHeight
 */
export async function getPnl(
  subaccountId: string,
  effectiveBeforeOrAtHeight: string,
): Promise<Big> {
  const [
    pnlOfFills,
    totalValueOfOpenPositions,
    realizedFunding,
    unrealizedFunding,
    feesPaid,
  ]: [
    Big,
    Big,
    Big,
    Big,
    Big,
  ] = await Promise.all([
    FillTable.getCostOfFills(subaccountId, effectiveBeforeOrAtHeight),
    FillTable.getTotalValueOfOpenPositions(subaccountId, effectiveBeforeOrAtHeight),
    getRealizedFunding(subaccountId, effectiveBeforeOrAtHeight),
    getUnsettledFunding(subaccountId, effectiveBeforeOrAtHeight),
    FillTable.getFeesPaid(subaccountId, effectiveBeforeOrAtHeight),
  ]);
  return pnlOfFills
    .add(totalValueOfOpenPositions)
    .sub(realizedFunding)
    .sub(unrealizedFunding)
    .sub(feesPaid);
}

/**
 * Validate pnl tick in database.
 *
 * @param pnlUuid
 */
export async function validatePnl(
  pnlUuid: string,
): Promise<void> {
  const pnlTick: PnlTicksFromDatabase | undefined = await
  PnlTicksTable.findById(pnlUuid, { readReplica: true });
  if (pnlTick === undefined) {
    console.log(`Pnl tick with uuid ${pnlUuid} not found.`);
    return;
  }
  const { subaccountId, blockHeight, totalPnl }: {
    subaccountId: string,
    blockHeight: string,
    totalPnl: string,
  } = pnlTick;
  const computedPnl: Big = await getPnl(subaccountId, blockHeight);
  // if computedPnl differs from totalPnl by more than 0.1%, log an error.
  if (computedPnl.minus(totalPnl).abs().gt(Big(totalPnl).abs().mul(0.001))) {
    console.log(`Pnl mismatch for subaccount ${subaccountId} at block height ${blockHeight}:
      Computed: ${computedPnl.toString()},
      Actual: ${totalPnl},
      Pnl tick uuid: ${pnlUuid}`);
  } else {
    console.log(`Pnl matches for subaccount ${subaccountId} at block height ${blockHeight}`);
  }
}

/**
 * Validate all pnl ticks for a subaccount, in order of ascending block height.
 *
 * @param subaccountId
 */
export async function validatePnlForSubaccount(
  subaccountId: string,
): Promise<void> {
  const { results: pnlTicks } = await
  PnlTicksTable.findAll(
    { subaccountId: [subaccountId] },
    [],
    { readReplica: true, orderBy: [[PnlTicksColumns.blockHeight, Ordering.ASC]] });
  for (const pnlTick of pnlTicks) {
    await validatePnl(pnlTick.id);
  }
}
