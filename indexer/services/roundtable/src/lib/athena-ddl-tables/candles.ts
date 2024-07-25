import {
  getAthenaTableCreationStatement,
  getExternalAthenaTableCreationStatement,
  castToTimestamp,
  castToDouble,
} from '../../helpers/sql';

const TABLE_NAME: string = 'candles';
const RAW_TABLE_COLUMNS: string = `
  \`id\` string,
  \`startedAt\` string,
  \`ticker\` string,
  \`resolution\` string,
  \`low\` string,
  \`high\` string,
  \`open\` string,
  \`close\` string,
  \`baseTokenVolume\` string,
  \`usdVolume\` string,
  \`trades\` int,
  \`startingOpenInterest\` string
  \`orderbookMidPriceOpen\` string
  \`orderbookMidPriceClose\` string
`;
const TABLE_COLUMNS: string = `
  "id",
  ${castToTimestamp('startedAt')},
  "ticker",
  "resolution",
  ${castToDouble('low')},
  ${castToDouble('high')},
  ${castToDouble('open')},
  ${castToDouble('close')},
  ${castToDouble('baseTokenVolume')},
  ${castToDouble('usdVolume')},
  "trades",
  ${castToDouble('startingOpenInterest')}
  ${castToDouble('orderbookMidPriceOpen')}
  ${castToDouble('orderbookMidPriceClose')}
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
