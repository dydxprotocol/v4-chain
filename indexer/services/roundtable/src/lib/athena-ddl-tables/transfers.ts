import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToDouble,
  castToTimestamp,
} from '../../helpers/sql';

const TABLE_NAME: string = 'transfers';
const RAW_TABLE_COLUMNS: string = `
  \`id\` string,
  \`senderSubaccountId\` string,
  \`recipientSubaccountId\` string,
  \`assetId\` string,
  \`size\` string,
  \`eventId\` binary,
  \`transactionHash\` string,
  \`createdAt\` string,
  \`createdAtHeight\` bigint
`;
const TABLE_COLUMNS: string = `
  "id",
  "senderSubaccountId",
  "recipientSubaccountId",
  "assetId",
  ${castToDouble('size')},
  "eventId",
  "transactionHash",
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
