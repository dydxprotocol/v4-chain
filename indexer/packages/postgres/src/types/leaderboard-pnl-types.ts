/* ------- LEADERBOARD PNL TYPES ------- */

export interface LeaderboardPnlCreateObject {
  address: string,
  pnl: string,
  timeSpan: string,
  currentEquity: string,
  rank: number,
}

export enum LeaderboardPnlColumns {
  address = 'address',
  timeSpan = 'timeSpan',
  pnl = 'pnl',
  currentEquity = 'currentEquity',
  rank = 'rank',
}

export enum LeaderboardPnlTimeSpan {
  ONE_DAY = 'ONE_DAY',
  SEVEN_DAYS = 'SEVEN_DAYS',
  THIRTY_DAYS = 'THIRTY_DAYS',
  ONE_YEAR = 'ONE_YEAR',
  ALL_TIME = 'ALL_TIME',
}
