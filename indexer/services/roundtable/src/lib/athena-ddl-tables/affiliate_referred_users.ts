import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
} from '../../helpers/sql';

const TABLE_NAME: string = 'affiliate_referred_users';
const RAW_TABLE_COLUMNS: string = `
  \`referredAddress\` string,
  \`affiliateAddress\` string,
  \`referredAtBlock\` bigint
`;
const TABLE_COLUMNS: string = `
  "referredAddress",
  "affiliateAddress",
  "referredAtBlock"
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
