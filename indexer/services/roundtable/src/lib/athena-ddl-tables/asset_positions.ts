import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToDouble,
} from '../../helpers/sql';

const TABLE_NAME: string = 'asset_positions';
const RAW_TABLE_COLUMNS: string = `
  \`id\` string,
  \`assetId\` string,
  \`subaccountId\` string,
  \`size\` string,
  \`isLong\` boolean
`;
const TABLE_COLUMNS: string = `
  "id",
  "assetId",
  "subaccountId",
  ${castToDouble('size')},
  "isLong"
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
