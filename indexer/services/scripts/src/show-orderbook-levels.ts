import { logger } from '@dydxprotocol-indexer/base';
import { RedisClient } from 'redis';
import yargs from 'yargs';

// Create redis client
const createRedisClient = (redisUrl: string): {
  client: RedisClient,
  connect: () => Promise<void>,
} => {
  const client: RedisClient = require('redis').createClient({
    url: redisUrl,
    tls: { rejectUnauthorized: false },
  });

  return {
    client,
    connect: async () => {
      return new Promise((resolve, reject) => {
        client.on('connect', () => {
          logger.info({
            at: 'redis#connect',
            message: 'Connected to Redis',
          });
          resolve();
        });
        client.on('error', reject);
      });
    },
  };
};

// Get all fields from a Redis hash
const getRedisHashAll = async (key: string, client: RedisClient): Promise<Record<string, string>> => {
  return new Promise((resolve, reject) => {
    client.hgetall(key, (err, value) => {
      if (err) return reject(err);
      resolve(value || {});
    });
  });
};

// Get a specific Redis value
const getRedisValue = async (key: string, client: RedisClient): Promise<string | null> => {
  return new Promise((resolve, reject) => {
    client.get(key, (err, value) => {
      if (err) {
        reject(err);
        return;
      }
      resolve(value);
    });
  });
};

// Format timestamp to human readable date and time
const formatTimestamp = (timestamp: string): string => {
  try {
    const date = new Date(parseInt(timestamp) * 1000);
    return date.toISOString();
  } catch (error) {
    return 'Invalid timestamp';
  }
};

async function showOrderbookLevels(
  ticker: string,
  redisUrl: string,
  priceLevelFilter?: string,
  showAllLevels?: boolean,
  minSize?: number,
  ignoreZeroSize: boolean = true,
): Promise<void> {
  const { client, connect } = createRedisClient(redisUrl);
  
  try {
    await connect();
    
    // Get buy and sell side levels
    const buyLevels = await getRedisHashAll(`v4/orderbookLevels/${ticker}/BUY`, client);
    const sellLevels = await getRedisHashAll(`v4/orderbookLevels/${ticker}/SELL`, client);
    
    // Get last updated timestamps
    const buyLastUpdated = await getRedisHashAll(`v4/orderbookLevels/${ticker}/BUY/lastUpdated`, client);
    const sellLastUpdated = await getRedisHashAll(`v4/orderbookLevels/${ticker}/SELL/lastUpdated`, client);
    
    // Sort price levels
    const sortedBuyLevels = Object.entries(buyLevels)
      .map(([price, size]) => ({ 
        price: parseFloat(price), 
        priceStr: price, 
        size, 
        sizeNum: parseInt(size),
        lastUpdated: buyLastUpdated[price] || 'unknown' 
      }))
      .sort((a, b) => b.price - a.price); // Buy side descending
    
    const sortedSellLevels = Object.entries(sellLevels)
      .map(([price, size]) => ({ 
        price: parseFloat(price), 
        priceStr: price, 
        size, 
        sizeNum: parseInt(size),
        lastUpdated: sellLastUpdated[price] || 'unknown'
      }))
      .sort((a, b) => a.price - b.price); // Sell side ascending
    
    // Filter out zero-sized levels if ignoreZeroSize is true
    const nonZeroBuyLevels = ignoreZeroSize 
      ? sortedBuyLevels.filter(level => level.sizeNum > 0)
      : sortedBuyLevels;
      
    const nonZeroSellLevels = ignoreZeroSize 
      ? sortedSellLevels.filter(level => level.sizeNum > 0)
      : sortedSellLevels;
    
    // Filter levels if a specific price level is requested
    const filteredBuyLevels = priceLevelFilter 
      ? nonZeroBuyLevels.filter(level => 
          Math.abs(level.price - parseFloat(priceLevelFilter)) < 0.01)
      : nonZeroBuyLevels;
    
    const filteredSellLevels = priceLevelFilter 
      ? nonZeroSellLevels.filter(level => 
          Math.abs(level.price - parseFloat(priceLevelFilter)) < 0.01)
      : nonZeroSellLevels;
    
    // Filter by minimum size if specified
    const sizeFilteredBuyLevels = minSize !== undefined
      ? filteredBuyLevels.filter(level => level.sizeNum >= minSize)
      : filteredBuyLevels;
    
    const sizeFilteredSellLevels = minSize !== undefined
      ? filteredSellLevels.filter(level => level.sizeNum >= minSize)
      : filteredSellLevels;
    
    logger.info({
      at: 'showOrderbookLevels',
      message: `Orderbook for ${ticker}`,
      totalBuyLevels: sortedBuyLevels.length,
      totalSellLevels: sortedSellLevels.length,
      nonZeroBuyLevels: nonZeroBuyLevels.length,
      nonZeroSellLevels: nonZeroSellLevels.length,
      filteredBuyLevels: sizeFilteredBuyLevels.length,
      filteredSellLevels: sizeFilteredSellLevels.length,
      filters: {
        priceLevel: priceLevelFilter || 'none',
        minSize: minSize !== undefined ? minSize : 'none',
        ignoreZeroSize,
      }
    });
    
    // Get levels to display
    const buyToDisplay = showAllLevels ? sizeFilteredBuyLevels : sizeFilteredBuyLevels.slice(0, 10);
    const sellToDisplay = showAllLevels ? sizeFilteredSellLevels : sizeFilteredSellLevels.slice(0, 10);
    
    // Print levels
    if (buyToDisplay.length > 0) {
      logger.info({
        at: 'showOrderbookLevels',
        message: `BUY LEVELS (${showAllLevels ? 'all' : 'top 10'})`,
        levels: buyToDisplay.map(level => ({
          price: level.priceStr,
          size: level.size,
          sizeHuman: (level.sizeNum / 10000).toString(),
          lastUpdated: level.lastUpdated,
          lastUpdatedTime: formatTimestamp(level.lastUpdated),
        })),
      });
    } else {
      logger.info({
        at: 'showOrderbookLevels',
        message: 'No matching BUY levels found',
      });
    }
    
    if (sellToDisplay.length > 0) {
      logger.info({
        at: 'showOrderbookLevels',
        message: `SELL LEVELS (${showAllLevels ? 'all' : 'top 10'})`,
        levels: sellToDisplay.map(level => ({
          price: level.priceStr,
          size: level.size,
          sizeHuman: (level.sizeNum / 10000).toString(),
          lastUpdated: level.lastUpdated,
          lastUpdatedTime: formatTimestamp(level.lastUpdated),
        })),
      });
    } else {
      logger.info({
        at: 'showOrderbookLevels',
        message: 'No matching SELL levels found',
      });
    }

  } catch (error) {
    logger.error({
      at: 'showOrderbookLevels',
      message: 'Error fetching orderbook levels',
      error,
    });
    process.exit(1);
  } finally {
    client.quit();
  }
}

const args = yargs.options({
  ticker: {
    type: 'string',
    alias: 't',
    description: 'Market ticker (e.g., SOL-USD)',
    required: true,
  },
  price: {
    type: 'string',
    alias: 'p',
    description: 'Specific price level to filter for',
  },
  all: {
    type: 'boolean',
    alias: 'a',
    description: 'Show all price levels (not just top 10)',
    default: false,
  },
  minSize: {
    type: 'number',
    alias: 'm',
    description: 'Minimum size for price levels to display (in quantums)',
  },
  showZeros: {
    type: 'boolean',
    alias: 'z',
    description: 'Include price levels with size = 0',
    default: false,
  },
  redisUrl: {
    type: 'string',
    alias: 'r',
    description: 'Redis URL',
    default: 'redis://master.redis-debug-backup.g88kxp.apne1.cache.amazonaws.com:6379',
  },
}).argv;

// Execute as a simple async function since this is a one-time script
(async () => {
  try {
    await showOrderbookLevels(
      args.ticker, 
      args.redisUrl, 
      args.price, 
      args.all,
      args.minSize,
      !args.showZeros
    );
  } catch (error) {
    logger.error({
      at: 'main',
      message: 'Script failed',
      error,
    });
    process.exit(1);
  }
})(); 