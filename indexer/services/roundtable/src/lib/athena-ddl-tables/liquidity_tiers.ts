import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToDouble,
} from '../../helpers/sql';

const TABLE_NAME: string = 'liquidity_tiers';
const RAW_TABLE_COLUMNS: string = `
  \`id\` int,
  \`name\` string,
  \`initialMarginPpm\` bigint,
  \`maintenanceFractionPpm\` bigint,
  \`basePositionNotional\` decimal
`;
const TABLE_COLUMNS: string = `
  "id",
  "name",
  "initialMarginPpm",
  "maintenanceFractionPpm",
  ${castToDouble('basePositionNotional')}
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
