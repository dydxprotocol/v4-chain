import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToDouble,
  castToTimestamp,
} from '../../helpers/sql';

const TABLE_NAME: string = 'perpetual_positions';
const RAW_TABLE_COLUMNS: string = `
  \`id\` string,
  \`subaccountId\` string,
  \`perpetualId\` bigint,
  \`side\` string,
  \`status\` string,
  \`size\` string,
  \`maxSize\` string,
  \`entryPrice\` string,
  \`exitPrice\` string,
  \`sumOpen\` string,
  \`sumClose\` string,
  \`createdAt\` string,
  \`closedAt\` string,
  \`createdAtHeight\` bigint,
  \`closedAtHeight\` bigint,
  \`openEventId\` binary,
  \`closeEventId\` binary,
  \`lastEventId\` binary,
  \`settledFunding\` string
`;
const TABLE_COLUMNS: string = `
  "id",
  "subaccountId",
  "perpetualId",
  "side",
  "status",
  ${castToDouble('size')},
  ${castToDouble('maxSize')},
  ${castToDouble('entryPrice')},
  ${castToDouble('exitPrice')},
  ${castToDouble('sumOpen')},
  ${castToDouble('sumClose')},
  ${castToTimestamp('createdAt')},
  ${castToTimestamp('closedAt')},
  "createdAtHeight",
  "closedAtHeight",
  "openEventId",
  "closeEventId",
  "lastEventId",
  ${castToDouble('settledFunding')}
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
