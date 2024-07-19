import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToTimestamp,
} from '../../helpers/sql';

const TABLE_NAME: string = 'blocks';
const RAW_TABLE_COLUMNS: string = `
  \`blockHeight\` bigint,
  \`time\` string
`;
const TABLE_COLUMNS: string = `
  "blockHeight",
  ${castToTimestamp('time')}
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
