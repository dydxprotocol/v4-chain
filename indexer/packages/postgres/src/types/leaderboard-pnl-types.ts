/* ------- LEADERBOARD PNL TYPES ------- */

export interface LeaderboardPnlCreateObject {
  subaccountId: string;
  pnl: string;
  timeSpan: string;
  currentEquity: string;
  rank: number;
}

export enum LeaderboardPnlColumns {
  subaccountId = 'subaccountId',
  timeSpan = 'timeSpan',
  pnl = 'pnl',
  currentEquity = 'currentEquity',
  rank = 'rank',
}
