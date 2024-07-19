import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
} from '../../helpers/sql';

const TABLE_NAME: string = 'assets';
const RAW_TABLE_COLUMNS: string = `
  \`id\` string,
  \`symbol\` string,
  \`atomicResolution\` int,
  \`hasMarket\` boolean,
  \`marketId\` int
`;
const TABLE_COLUMNS: string = `
  "id",
  "symbol",
  "atomicResolution",
  "hasMarket",
  "marketId"
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
