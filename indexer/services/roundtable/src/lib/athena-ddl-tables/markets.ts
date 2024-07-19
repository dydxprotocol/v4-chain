import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToDouble,
} from '../../helpers/sql';

const TABLE_NAME: string = 'markets';
const RAW_TABLE_COLUMNS: string = `
  \`id\` int,
  \`pair\` string,
  \`exponent\` int,
  \`minPriceChangePpm\` int,
  \`oraclePrice\` string
`;
const TABLE_COLUMNS: string = `
  "id",
  "pair",
  "exponent",
  "minPriceChangePpm",
  ${castToDouble('oraclePrice')}
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
