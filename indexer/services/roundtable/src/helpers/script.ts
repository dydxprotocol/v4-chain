// Sample data from the logs
import {PerpetualPositionFromDatabase, PerpetualPositionStatus, PositionSide} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';

const subaccountId = '070180e5-7da1-58b2-b7ae-15d690160a0c';
const latestBlockHeight = 19010585;
const latestBlockTime = new Date('2024-06-26T18:56:28.918Z');
const usdcPositionSize = new Big('63947.007662');

const openPerpetualPositionsForSubaccount: PerpetualPositionFromDatabase[] = [
  {
    id: '623d0608-240b-5cd0-a5b3-6dd1d5f22c5c',
    subaccountId,
    perpetualId: '3',
    side: PositionSide.SHORT,
    status: PerpetualPositionStatus.OPEN,
    size: '-8770',
    maxSize: '-3500',
    entryPrice: '0.57064218928164196123',
    exitPrice: undefined,
    sumOpen: '8770',
    sumClose: '0',
    createdAt: '2024-06-19T14:55:32.510Z',
    closedAt: undefined,
    createdAtHeight: '18488696',
    closedAtHeight: undefined,
    openEventId: Buffer.from([1, 26, 29, 120, 0, 0, 0, 2, 0, 0, 0, 0]),
    closeEventId: undefined,
    lastEventId: Buffer.from([1, 31, 180, 63, 0, 0, 0, 2, 0, 0, 0, 31]),
    settledFunding: '0.047358',
  },
  // ... (other positions)
];

const marketPrices = {
  '0': '60909.76231',
  '1': '3361.77635',
  '2': '13.9680125',
  '3': '0.552134848',
  // ... (other prices)
};

const lastUpdatedFundingIndexMap = {
  '0': '10061.34',
  '1': '543.858',
  '2': '2.985764',
  '3': '0.1285959',
  // ... (other indices)
};

const currentFundingIndexMap = {
  '0': '10055.39',
  '1': '543.467',
  '2': '2.982482',
  '3': '0.1286097',
  // ... (other indices)
};

const subaccountTotalTransfersMap: { [key: string]: { [key: string]: Big } } = {
  '070180e5-7da1-58b2-b7ae-15d690160a0c': {
    'USDC': new Big('63947.007662'),
  },
};

const USDC_ASSET_ID = 'USDC';

// Dummy implementations of calculateEquity and calculateTotalPnl
function calculateEquity(
  usdcPositionSize: Big,
  openPerpetualPositionsForSubaccount: PerpetualPosition[],
  marketPrices: { [key: string]: string },
  lastUpdatedFundingIndexMap: { [key: string]: string },
  currentFundingIndexMap: { [key: string]: string }
): Big {
  // Implement your equity calculation logic here
  return new Big(0);
}

function calculateTotalPnl(
  currentEquity: Big,
  subaccountTotalTransfers: Big
): Big {
  // Implement your PnL calculation logic here
  return new Big(0);
}

// Call the functions with the provided data
const currentEquity: Big = calculateEquity(
  usdcPositionSize,
  openPerpetualPositionsForSubaccount,
  marketPrices,
  lastUpdatedFundingIndexMap,
  currentFundingIndexMap
);

const totalPnl: Big = calculateTotalPnl(
  currentEquity,
  subaccountTotalTransfersMap[subaccountId][USDC_ASSET_ID]
);

// Log the results
console.log('Calculated equity and total PnL:', {
  subaccountId,
  currentEquity: currentEquity.toFixed(),
  totalPnl: totalPnl.toFixed(),
});
