import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
} from '../../helpers/sql';

const TABLE_NAME: string = 'trading_rewards';
const RAW_TABLE_COLUMNS: string = `
  \`id\` binary,
  \`address\` string,
  \`blockTime\` string,
  \`blockHeight\` bigint,
  \`amount\` string
`;
const TABLE_COLUMNS: string = `
  "id",
  "address",
  "blockTime",
  "blockHeight",
  "amount"
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
