import express from 'express';

import AddressesController from './v4/addresses-controller';
import AssetPositionsController from './v4/asset-positions-controller';
import CandlesController from './v4/candles-controller';
import ComplianceController from './v4/compliance-controller';
import FillsController from './v4/fills-controller';
import HeightController from './v4/height-controller';
import HistoricalFundingController from './v4/historical-funding-controller';
import PnlticksController from './v4/historical-pnl-controller';
import OrderbooksController from './v4/orderbook-controller';
import OrdersController from './v4/orders-controller';
import PerpetualMarketController from './v4/perpetual-markets-controller';
import PerpetualPositionsController from './v4/perpetual-positions-controller';
import SparklinesController from './v4/sparklines-controller';
import TimeController from './v4/time-controller';
import TradesController from './v4/trades-controller';
import TransfersController from './v4/transfers-controller';
import YieldParamsController from './v4/yield-params-controller';

// Keep routers in alphabetical order

const router: express.Router = express.Router();
router.use('/addresses', AddressesController);
router.use('/assetPositions', AssetPositionsController);
router.use('/candles', CandlesController);
router.use('/fills', FillsController);
router.use('/height', HeightController);
router.use('/historicalFunding', HistoricalFundingController);
router.use('/historical-pnl', PnlticksController);
router.use('/orders', OrdersController);
router.use('/orderbooks', OrderbooksController);
router.use('/perpetualMarkets', PerpetualMarketController);
router.use('/perpetualPositions', PerpetualPositionsController);
router.use('/sparklines', SparklinesController);
router.use('/time', TimeController);
router.use('/trades', TradesController);
router.use('/transfers', TransfersController);
router.use('/screen', ComplianceController);
router.use('/yieldParams', YieldParamsController);

export default router;
