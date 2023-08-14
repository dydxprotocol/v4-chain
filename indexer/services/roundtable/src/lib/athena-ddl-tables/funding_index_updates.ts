import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToTimestamp,
  castToDouble,
} from '../../helpers/sql';

const TABLE_NAME: string = 'funding_index_updates';
const RAW_TABLE_COLUMNS: string = `
  \`id\` string,
  \`perpetualId\` bigint,
  \`eventId\` binary,
  \`rate\` string,
  \`oraclePrice\` string,
  \`fundingIndex\` string,
  \`effectiveAt\` string,
  \`effectiveAtHeight\` bigint
`;
const TABLE_COLUMNS: string = `
  "id",
  "perpetualId",
  "eventId",
  ${castToDouble('rate')},
  ${castToDouble('oraclePrice')},
  ${castToDouble('fundingIndex')},
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
