import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToTimestamp,
} from '../../helpers/sql';

const TABLE_NAME: string = 'subaccounts';
const RAW_TABLE_COLUMNS: string = `
  \`id\` string,
  \`address\` string,
  \`subaccountNumber\` int,
  \`updatedAt\` string,
  \`updatedAtHeight\` bigint
`;
const TABLE_COLUMNS: string = `
  "id",
  "address",
  "subaccountNumber",
  ${castToTimestamp('updatedAt')},
  "updatedAtHeight"
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
