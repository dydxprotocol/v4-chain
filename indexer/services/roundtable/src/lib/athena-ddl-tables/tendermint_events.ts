import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
} from '../../helpers/sql';

const TABLE_NAME: string = 'tendermint_events';
const RAW_TABLE_COLUMNS: string = `
  \`id\` binary,
  \`blockHeight\` bigint,
  \`transactionIndex\` bigint,
  \`eventIndex\` bigint
`;
const TABLE_COLUMNS: string = `
  "id",
  "blockHeight",
  "transactionIndex",
  "eventIndex"
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
