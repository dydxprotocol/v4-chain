import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToDouble,
  castToTimestamp,
} from '../../helpers/sql';

const TABLE_NAME: string = 'oracle_prices';
const RAW_TABLE_COLUMNS: string = `
  \`id\` string,
  \`marketId\` int,
  \`spotPrice\` string,
  \`pnlPrice\` string,
  \`effectiveAt\` string,
  \`effectiveAtHeight\` bigint
`;
const TABLE_COLUMNS: string = `
  "id",
  "marketId",
  ${castToDouble('spotPrice')},
  ${castToDouble('pnlPrice')},
  ${castToTimestamp('effectiveAt')},
  "effectiveAtHeight"
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
