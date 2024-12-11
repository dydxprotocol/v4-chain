import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToDouble,
} from '../../helpers/sql';

const TABLE_NAME: string = 'affiliate_info';
const RAW_TABLE_COLUMNS: string = `
  \`address\` string,
  \`affiliateEarnings\` string,
  \`referredMakerTrades\` int,
  \`referredTakerTrades\` int,
  \`totalReferredUsers\` int,
  \`firstReferralBlockHeight\` bigint,
  \`referredTotalVolume\` string,
  \`totalReferredTakerFees\` string,
  \`totalReferredMakerFees\` string,
  \`totalReferredMakerRebates\` string
`;
const TABLE_COLUMNS: string = `
  "address",
  ${castToDouble('affiliateEarnings')},
  "referredMakerTrades",
  "referredTakerTrades",
  "totalReferredUsers",
  "firstReferralBlockHeight",
  ${castToDouble('referredTotalVolume')},
  ${castToDouble('totalReferredTakerFees')},
  ${castToDouble('totalReferredMakerFees')},
  ${castToDouble('totalReferredMakerRebates')}
`;

export function generateRawTable(tablePrefix: string, rdsExportIdentifier: string): string {
  return getExternalAthenaTableCreationStatement(
    tablePrefix,
    rdsExportIdentifier,
    TABLE_NAME,
    RAW_TABLE_COLUMNS,
  );
}

export function generateTable(tablePrefix: string): string {
  return getAthenaTableCreationStatement(
    tablePrefix,
    TABLE_NAME,
    TABLE_COLUMNS,
  );
}
