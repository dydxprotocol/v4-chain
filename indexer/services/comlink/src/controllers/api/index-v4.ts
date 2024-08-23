import express from 'express';

import AddressesController from './v4/addresses-controller';
import AffiliatesController from './v4/affiliates-controller';
import AssetPositionsController from './v4/asset-positions-controller';
import CandlesController from './v4/candles-controller';
import ComplianceController from './v4/compliance-controller';
import ComplianceV2Controller from './v4/compliance-v2-controller';
import FillsController from './v4/fills-controller';
import HeightController from './v4/height-controller';
import HistoricalBlockTradingRewardController from './v4/historical-block-trading-rewards-controller';
import HistoricalFundingController from './v4/historical-funding-controller';
import PnlticksController from './v4/historical-pnl-controller';
import HistoricalTradingRewardController from './v4/historical-trading-reward-aggregations-controller';
import OrderbooksController from './v4/orderbook-controller';
import OrdersController from './v4/orders-controller';
import PerpetualMarketController from './v4/perpetual-markets-controller';
import PerpetualPositionsController from './v4/perpetual-positions-controller';
import SocialTradingController from './v4/social-trading-controller';
import SparklinesController from './v4/sparklines-controller';
import TimeController from './v4/time-controller';
import TradesController from './v4/trades-controller';
import TransfersController from './v4/transfers-controller';
import VaultController from './v4/vault-controller';

// Keep routers in alphabetical order

const router: express.Router = express.Router();
router.use('/addresses', AddressesController);
router.use('/affiliates', AffiliatesController);
router.use('/assetPositions', AssetPositionsController);
router.use('/candles', CandlesController);
router.use('/fills', FillsController);
router.use('/height', HeightController);
router.use('/historicalBlockTradingRewards', HistoricalBlockTradingRewardController);
router.use('/historicalFunding', HistoricalFundingController);
router.use('/historical-pnl', PnlticksController);
router.use('/historicalTradingRewardAggregations', HistoricalTradingRewardController);
router.use('/orders', OrdersController);
router.use('/orderbooks', OrderbooksController);
router.use('/perpetualMarkets', PerpetualMarketController);
router.use('/perpetualPositions', PerpetualPositionsController);
router.use('/sparklines', SparklinesController);
router.use('/time', TimeController);
router.use('/trades', TradesController);
router.use('/transfers', TransfersController);
router.use('/screen', ComplianceController);
router.use('/compliance', ComplianceV2Controller);
router.use('/trader', SocialTradingController);
router.use('/vault', VaultController);

export default router;
