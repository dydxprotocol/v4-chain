import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToDouble,
} from '../../helpers/sql';

const TABLE_NAME: string = 'perpetual_markets';
const RAW_TABLE_COLUMNS: string = `
  \`id\` bigint,
  \`clobPairId\` bigint,
  \`ticker\` string,
  \`marketId\` int,
  \`status\` string,
  \`lastPrice\` string,
  \`priceChange24H\` string,
  \`volume24H\` string,
  \`trades24H\` int,
  \`nextFundingRate\` string,
  \`openInterest\` string,
  \`quantumConversionExponent\` int,
  \`atomicResolution\` int,
  \`subticksPerTick\` int,
  \`minOrderBaseQuantums\` int,
  \`stepBaseQuantums\` int,
  \`liquidityTierId\` int
`;
const TABLE_COLUMNS: string = `
  "id",
  "clobPairId",
  "ticker",
  "marketId",
  "status",
  ${castToDouble('lastPrice')},
  ${castToDouble('priceChange24H')},
  ${castToDouble('volume24H')},
  "trades24H",
  ${castToDouble('nextFundingRate')},
  ${castToDouble('openInterest')},
  "quantumConversionExponent",
  "atomicResolution",
  "subticksPerTick",
  "minOrderBaseQuantums",
  "stepBaseQuantums",
  "liquidityTierId"
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
