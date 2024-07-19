import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToDouble,
  castToTimestamp,
} from '../../helpers/sql';

const TABLE_NAME: string = 'pnl_ticks';
const RAW_TABLE_COLUMNS: string = `
  \`id\` string,
  \`subaccountId\` string,
  \`equity\` string,
  \`totalPnl\` string,
  \`netTransfers\` string,
  \`createdAt\` string,
  \`blockHeight\` bigint,
  \`blockTime\` string
`;
const TABLE_COLUMNS: string = `
  "id",
  "subaccountId",
  ${castToDouble('equity')},
  ${castToDouble('totalPnl')},
  ${castToDouble('netTransfers')},
  ${castToTimestamp('createdAt')},
  "blockHeight",
  ${castToTimestamp('blockTime')}
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
