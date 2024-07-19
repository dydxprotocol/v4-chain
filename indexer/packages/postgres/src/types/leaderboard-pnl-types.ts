/* ------- LEADERBOARD PNL TYPES ------- */

export interface LeaderboardPnlCreateObject {
  address: string;
  pnl: string;
  timeSpan: string;
  currentEquity: string;
  rank: number;
}

export enum LeaderboardPnlColumns {
  address = 'address',
  timeSpan = 'timeSpan',
  pnl = 'pnl',
  currentEquity = 'currentEquity',
  rank = 'rank',
}
