/* ------- LEADERBOARD PNL TYPES ------- */

export interface LeaderboardPNLCreateObject {
  subaccountId: string;
  pnl: string;
  timeSpan: string;
  currentEquity: string;
  rank: number;
}

export enum LeaderboardPNLColumns {
  subaccountId = 'subaccountId',
  timeSpan = 'timeSpan',
  pnl = 'pnl',
  currentEquity = 'currentEquity',
  rank = 'rank',
}
