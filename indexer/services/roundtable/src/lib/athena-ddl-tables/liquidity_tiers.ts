import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
} from '../../helpers/sql';

const TABLE_NAME: string = 'liquidity_tiers';
const RAW_TABLE_COLUMNS: string = `
  \`id\` int,
  \`name\` string,
  \`initialMarginPpm\` bigint,
  \`maintenanceFractionPpm\` bigint,
  \`openInterestLowerCap\` string,
  \`openInterestUpperCap\` string
`;
const TABLE_COLUMNS: string = `
  "id",
  "name",
  "initialMarginPpm",
  "maintenanceFractionPpm",
  "openInterestLowerCap",
  "openInterestUpperCap"
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
