import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToDouble,
  castToTimestamp,
} from '../../helpers/sql';

const TABLE_NAME: string = 'pnl';

const RAW_TABLE_COLUMNS: string = `
  \`subaccountId\` string,
  \`equity\` string,
  \`totalPnl\` string,
  \`netTransfers\` string,
  \`createdAt\` string,
  \`createdAtHeight\` bigint
`;

const TABLE_COLUMNS: string = `
  "subaccountId",
  ${castToDouble('equity')},
  ${castToDouble('totalPnl')},
  ${castToDouble('netTransfers')},
  ${castToTimestamp('createdAt')},
  "createdAtHeight"
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