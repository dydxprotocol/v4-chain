import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToTimestamp,
  castToDouble,
} from '../../helpers/sql';

const TABLE_NAME: string = 'fills';
// TODO(INDEX-293): Join column names by commas.
const RAW_TABLE_COLUMNS: string = `
  \`id\` string,
  \`subaccountId\` string,
  \`side\` string,
  \`liquidity\` string,
  \`type\` string,
  \`clobPairId\` bigint,
  \`orderId\` string,
  \`size\` string,
  \`price\` string,
  \`quoteAmount\` string,
  \`eventId\` binary,
  \`transactionHash\` string,
  \`createdAt\` string,
  \`createdAtHeight\` bigint,
  \`clientMetadata\` bigint,
  \`fee\` string
`;
const TABLE_COLUMNS: string = `
  "id",
  "subaccountId",
  "side",
  "liquidity",
  "type",
  "clobPairId",
  "orderId",
  ${castToDouble('size')},
  ${castToDouble('price')},
  ${castToDouble('quoteAmount')},
  "eventId",
  "transactionHash",
  ${castToTimestamp('createdAt')},
  "createdAtHeight",
  "clientMetadata",
  "fee"
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
