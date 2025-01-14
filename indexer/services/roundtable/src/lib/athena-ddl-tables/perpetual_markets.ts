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
  \`priceChange24H\` string,
  \`volume24H\` string,
  \`trades24H\` int,
  \`nextFundingRate\` string,
  \`openInterest\` string,
  \`quantumConversionExponent\` int,
  \`atomicResolution\` int,
  \`subticksPerTick\` int,
  \`stepBaseQuantums\` int,
  \`liquidityTierId\` int,
  \`marketType\` string,
  \`baseOpenInterest\` string,
  \`defaultFundingRate1H\` string
`;
const TABLE_COLUMNS: string = `
  "id",
  "clobPairId",
  "ticker",
  "marketId",
  "status",
  ${castToDouble('priceChange24H')},
  ${castToDouble('volume24H')},
  "trades24H",
  ${castToDouble('nextFundingRate')},
  ${castToDouble('openInterest')},
  "quantumConversionExponent",
  "atomicResolution",
  "subticksPerTick",
  "stepBaseQuantums",
  "liquidityTierId",
  "marketType",
  ${castToDouble('baseOpenInterest')},
  ${castToDouble('defaultFundingRate1H')}
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
