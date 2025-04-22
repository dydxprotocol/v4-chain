import { MarketFromDatabase, PerpetualMarketFromDatabase } from './db-model-types';

/*
  * PerpetualMarketWithMarket combines PerpetualMarketFromDatabase and MarketFromDatabase,
  * excluding 'id' from MarketFromDatabase.
*/
export interface PerpetualMarketWithMarket
  extends PerpetualMarketFromDatabase, Omit<MarketFromDatabase, 'id'> {}
