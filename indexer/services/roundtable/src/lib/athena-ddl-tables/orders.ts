import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToDouble,
} from '../../helpers/sql';

const TABLE_NAME: string = 'orders';
const RAW_TABLE_COLUMNS: string = `
  \`id\` string,
  \`subaccountId\` string,
  \`clientId\` bigint,
  \`clobPairId\` bigint,
  \`side\` string,
  \`size\` string,
  \`totalFilled\` string,
  \`price\` string,
  \`type\` string,
  \`status\` string,
  \`timeInForce\` string,
  \`reduceOnly\` boolean,
  \`orderFlags\` bigint,
  \`goodTilBlock\` bigint,
  \`goodTilBlockTime\` string,
  \`createdAtHeight\` bigint,
  \`clientMetadata\` bigint,
  \`triggerPrice\` string
`;
const TABLE_COLUMNS: string = `
  "id",
  "subaccountId",
  "clientId",
  "clobPairId",
  "side",
  ${castToDouble('size')},
  ${castToDouble('totalFilled')},
  ${castToDouble('price')},
  "type",
  "status",
  "timeInForce",
  "reduceOnly",
  "orderFlags",
  "goodTilBlock",
  goodTilBlockTime,
  "createdAtHeight",
  "clientMetadata",
  ${castToDouble('triggerPrice')},
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
