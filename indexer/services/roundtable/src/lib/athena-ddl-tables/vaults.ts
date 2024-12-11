import {
  castToTimestamp,
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
} from '../../helpers/sql';

const TABLE_NAME: string = 'vaults';
const RAW_TABLE_COLUMNS: string = `
  \`address\` string,
  \`clobPairId\` bigint,
  \`status\` string,
  \`createdAt\` string,
  \`updatedAt\` string
`;
const TABLE_COLUMNS: string = `
  "address",
  "clobPairId",
  "status",
  ${castToTimestamp('createdAt')},
  ${castToTimestamp('updatedAt')}
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
