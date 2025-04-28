import { logger } from '@dydxprotocol-indexer/base';
import { RedisClient } from 'redis';
import { RedisOrder, redisOrder_TickerTypeToJSON } from '@dydxprotocol-indexer/v4-protos';
import yargs from 'yargs';
import Long from 'long';

// Create redis client
const createRedisClient = (redisUrl: string): {
  client: RedisClient,
  connect: () => Promise<void>,
} => {
  // Using standard redis client without the wrapper
  const client: RedisClient = require('redis').createClient({
    url: redisUrl,
    tls: { rejectUnauthorized: false }, // For TLS connection
    retry_strategy: (_options: any) => {
      // Retry every 1 second for infinity
      return 1000;
    },
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

        client.on('error', (error) => {
          logger.error({
            at: 'redis#error',
            message: 'Error connecting to Redis',
            error,
          });
          reject(error);
        });
      });
    },
  };
};

// Get all order keys from Redis
const getAllOrderKeys = async (client: RedisClient): Promise<string[]> => {
  return new Promise((resolve, reject) => {
    client.keys('v4/orders/*', (err, keys) => {
      if (err) {
        reject(err);
        return;
      }
      resolve(keys || []);
    });
  });
};

// Check for any keys matching a pattern
const getMatchingKeys = async (pattern: string, client: RedisClient): Promise<string[]> => {
  return new Promise((resolve, reject) => {
    client.keys(pattern, (err, keys) => {
      if (err) {
        reject(err);
        return;
      }
      resolve(keys || []);
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

// Get multiple Redis values in parallel
const getMultipleRedisValues = async (keys: string[], client: RedisClient): Promise<(string | null)[]> => {
  return Promise.all(keys.map(key => getRedisValue(key, client)));
};

// Get all fields from a Redis hash
const getRedisHashAll = async (key: string, client: RedisClient): Promise<Record<string, string>> => {
  return new Promise((resolve, reject) => {
    client.hgetall(key, (err, value) => {
      if (err) {
        reject(err);
        return;
      }
      resolve(value || {});
    });
  });
};

// Get a specific Redis hash field value
const getRedisHashValue = async (key: string, field: string, client: RedisClient): Promise<string | null> => {
  return new Promise((resolve, reject) => {
    client.hget(key, field, (err, value) => {
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

// Decode a Redis Order from binary data
const decodeRedisOrder = (data: string): RedisOrder | null => {
  try {
    return RedisOrder.decode(Buffer.from(data, 'binary'));
  } catch (error) {
    logger.error({
      at: 'decodeRedisOrder',
      message: 'Failed to decode order',
      error,
    });
    return null;
  }
};

// Parse order data string into structured object
const parseOrderData = (orderDataString: string): { goodTilBlock: string, totalFilledQuantums: string, restingOnBook: boolean } | null => {
  try {
    const [goodTilBlock, totalFilledQuantums, restingOnBook] = orderDataString.split('_');
    return {
      goodTilBlock,
      totalFilledQuantums,
      restingOnBook: restingOnBook === 'true',
    };
  } catch (error) {
    logger.error({
      at: 'parseOrderData',
      message: 'Failed to parse order data',
      error,
    });
    return null;
  }
};

const isPriceMatch = (price1: string, price2: string): boolean => {
  // First check for exact string match, which is the expected behavior
  if (price1 === price2) {
    return true;
  }

  return false;
};

// Process an order key to check if it matches criteria and return the order if it does
const processOrderKey = async (
  key: string, 
  ticker: string, 
  normalizedSide: string, 
  price: string, 
  client: RedisClient
): Promise<RedisOrder | null> => {
  const orderData = await getRedisValue(key, client);
  if (!orderData) return null;
  
  const order = decodeRedisOrder(orderData);
  if (!order) return null;
  
  // Check if this order matches our criteria
  if (order.ticker === ticker) {
    const orderSide = order.order?.side === 1 ? 'BUY' : 'SELL';
    
    if (orderSide === normalizedSide) {
      // Check if the order's price matches the specified price
      if (isPriceMatch(order.price, price)) {
        return order;
      }
    }
  }
  
  return null;
};

// Process a batch of order keys in parallel
const processBatchOfOrders = async (
  orderKeys: string[], 
  ticker: string, 
  normalizedSide: string, 
  price: string, 
  client: RedisClient
): Promise<RedisOrder[]> => {
  const results = await Promise.all(
    orderKeys.map(key => processOrderKey(key, ticker, normalizedSide, price, client))
  );
  return results.filter(order => order !== null) as RedisOrder[];
};

// Process all orders for entire orderbook side analysis
const processAllOrders = async (
  orderKeys: string[],
  ticker: string,
  normalizedSide: string,
  client: RedisClient
): Promise<Map<string, RedisOrder[]>> => {
  logger.info({
    at: 'processAllOrders',
    message: `Processing ${orderKeys.length} orders for ${ticker} ${normalizedSide}...`,
  });

  // Fetch and decode all orders in batches
  const BATCH_SIZE = 100;
  const tickerSideOrders: RedisOrder[] = [];
  
  for (let i = 0; i < orderKeys.length; i += BATCH_SIZE) {
    const batchKeys = orderKeys.slice(i, i + BATCH_SIZE);
    const batchData = await getMultipleRedisValues(batchKeys, client);
    
    // Process each order in the batch
    for (let j = 0; j < batchKeys.length; j++) {
      const data = batchData[j];
      if (!data) continue;
      
      const order = decodeRedisOrder(data);
      if (!order) continue;
      
      // Check if order matches ticker and side
      if (order.ticker === ticker) {
        const orderSide = order.order?.side === 1 ? 'BUY' : 'SELL';
        if (orderSide === normalizedSide) {
          tickerSideOrders.push(order);
        }
      }
    }
    
    if ((i + BATCH_SIZE) % 1000 === 0 || i + BATCH_SIZE >= orderKeys.length) {
      logger.info({
        at: 'processAllOrders',
        message: `Processed ${Math.min(i + BATCH_SIZE, orderKeys.length)}/${orderKeys.length} orders...`,
        matchingOrdersFound: tickerSideOrders.length,
      });
    }
  }
  
  // Group orders by price
  const ordersByPrice = new Map<string, RedisOrder[]>();
  
  for (const order of tickerSideOrders) {
    const price = order.price;
    if (!ordersByPrice.has(price)) {
      ordersByPrice.set(price, []);
    }
    ordersByPrice.get(price)!.push(order);
  }
  
  return ordersByPrice;
};

async function analyzePriceLevel(
  ticker: string,
  normalizedSide: string,
  price: string,
  client: RedisClient,
  verbose: boolean = true
): Promise<{
  price: string, 
  matchingOrdersCount: number,
  calculatedTotalSize: string,
  orderbookCacheSize: string,
  sizesMatch: boolean,
  matchingOrders: RedisOrder[],
  lastUpdated?: string,
  lastUpdatedTime?: string
}> {
  // Get the total size for this price level from orderbookLevels cache
  const orderbookLevelKey = `v4/orderbookLevels/${ticker}/${normalizedSide}`;
  
  // First try to get the exact price match
  let totalSize = await getRedisHashValue(orderbookLevelKey, price, client);
  let priceLevelFound = !!totalSize;
  
  // If not found with exact match, check if there's a very close price level
  // (this should only happen if there are floating point precision issues)
  if (!totalSize) {
    if (verbose) {
      logger.info({
        at: 'analyzePriceLevel',
        message: `No exact price level found for ${price}, checking for possible precision-related variations...`,
      });
    }
    
    const allPriceLevels = await getRedisHashAll(orderbookLevelKey, client);
    
    // Find a price that approximately matches
    for (const [priceKey, sizeValue] of Object.entries(allPriceLevels)) {
      if (isPriceMatch(priceKey, price) && priceKey !== price) {
        if (verbose) {
          logger.info({
            at: 'analyzePriceLevel',
            message: `Found price level with different representation: ${priceKey} (requested: ${price})`,
          });
        }
        totalSize = sizeValue;
        // Update the price to the one found in the cache
        price = priceKey;
        break;
      }
    }
  }
  
  if (!totalSize) {
    if (verbose) {
      logger.info({
        at: 'analyzePriceLevel',
        message: `No orders found at price level ${price} for ${ticker} ${normalizedSide}`,
      });
    }
    return {
      price,
      matchingOrdersCount: 0,
      calculatedTotalSize: "0",
      orderbookCacheSize: "0",
      sizesMatch: true,
      matchingOrders: [],
    };
  }

  // Get the last updated timestamp for this price level
  const lastUpdatedKey = `${orderbookLevelKey}/lastUpdated`;
  const lastUpdatedTimestamp = await getRedisHashValue(lastUpdatedKey, price, client);
  const formattedTimestamp = lastUpdatedTimestamp 
    ? formatTimestamp(lastUpdatedTimestamp) 
    : 'Unknown';

  if (verbose) {
    logger.info({
      at: 'analyzePriceLevel',
      message: `Found price level with total size: ${totalSize}`,
      ticker,
      side: normalizedSide,
      price,
      lastUpdated: lastUpdatedTimestamp || 'Unknown',
      lastUpdatedTime: formattedTimestamp,
    });
  }

  // Get all order keys
  if (verbose) {
    logger.info({
      at: 'analyzePriceLevel',
      message: 'Fetching all order keys from Redis...',
    });
  }
  const orderKeys = await getAllOrderKeys(client);
  
  if (verbose) {
    logger.info({
      at: 'analyzePriceLevel',
      message: `Found ${orderKeys.length} total orders in Redis, processing in parallel batches...`,
    });
  }

  // Process orders in parallel batches for better performance
  const BATCH_SIZE = 100; // Adjust based on Redis performance
  const matchingOrders: RedisOrder[] = [];
  
  // Process the orders in batches
  for (let i = 0; i < orderKeys.length; i += BATCH_SIZE) {
    const batchKeys = orderKeys.slice(i, i + BATCH_SIZE);
    const batchResults = await processBatchOfOrders(batchKeys, ticker, normalizedSide, price, client);
    matchingOrders.push(...batchResults);
    
    if (verbose && ((i + BATCH_SIZE) % 1000 === 0 || i + BATCH_SIZE >= orderKeys.length)) {
      logger.info({
        at: 'analyzePriceLevel',
        message: `Processed ${Math.min(i + BATCH_SIZE, orderKeys.length)}/${orderKeys.length} orders...`,
        matchingOrdersFound: matchingOrders.length,
      });
    }
  }

  // Calculate total size of matching orders
  let calculatedTotalSize = Long.ZERO;
  
  matchingOrders.forEach((order) => {
    if (order.order?.quantums) {
      calculatedTotalSize = calculatedTotalSize.add(
        new Long(order.order.quantums.low, order.order.quantums.high, order.order.quantums.unsigned)
      );
    }
  });

  // Compare with the size from orderbookLevels
  const orderbookTotalSize = Long.fromString(totalSize);
  const sizesMatch = calculatedTotalSize.eq(orderbookTotalSize);

  if (verbose) {
    logger.info({
      at: 'analyzePriceLevel',
      message: `Price level analysis complete`,
      ticker,
      side: normalizedSide,
      price,
      matchingOrdersCount: matchingOrders.length,
      calculatedTotalSize: calculatedTotalSize.toString(),
      orderbookCacheSize: orderbookTotalSize.toString(),
      sizesDifferenceQuantums: (calculatedTotalSize.subtract(orderbookTotalSize)).toString(),
      sizesMatch,
    });
  }

  return {
    price,
    matchingOrdersCount: matchingOrders.length,
    calculatedTotalSize: calculatedTotalSize.toString(),
    orderbookCacheSize: totalSize,
    sizesMatch,
    matchingOrders,
    lastUpdated: lastUpdatedTimestamp || undefined,
    lastUpdatedTime: formattedTimestamp !== 'Unknown' ? formattedTimestamp : undefined,
  };
}

async function analyzeEntireOrderbookSide(
  ticker: string,
  normalizedSide: string,
  redisUrl: string,
): Promise<void> {
  logger.info({
    at: 'analyzeEntireOrderbookSide',
    message: `Analyzing entire orderbook side for ${ticker} ${normalizedSide}`,
  });

  const { client, connect } = createRedisClient(redisUrl);
  
  try {
    await connect();
    
    // Get all price levels for this side
    const orderbookLevelKey = `v4/orderbookLevels/${ticker}/${normalizedSide}`;
    const allPriceLevels = await getRedisHashAll(orderbookLevelKey, client);
    
    // Filter out zero-sized levels
    const nonZeroPriceLevels = Object.entries(allPriceLevels)
      .filter(([_, size]) => parseInt(size) > 0)
      .map(([price, _]) => price);
    
    logger.info({
      at: 'analyzeEntireOrderbookSide',
      message: `Found ${nonZeroPriceLevels.length} non-zero price levels to analyze`,
    });
    
    // Analyze each price level
    const results = [];
    for (let i = 0; i < nonZeroPriceLevels.length; i++) {
      const price = nonZeroPriceLevels[i];
      logger.info({
        at: 'analyzeEntireOrderbookSide',
        message: `Analyzing price level ${i+1}/${nonZeroPriceLevels.length}: ${price}`,
      });
      
      const result = await analyzePriceLevel(ticker, normalizedSide, price, client, false);
      results.push(result);
      
      // Log progress every 10 levels or at the end
      if ((i + 1) % 10 === 0 || i === nonZeroPriceLevels.length - 1) {
        logger.info({
          at: 'analyzeEntireOrderbookSide',
          message: `Completed ${i+1}/${nonZeroPriceLevels.length} price levels`,
        });
      }
    }
    
    // Generate summary
    const matchingSizes = results.filter(r => r.sizesMatch).length;
    const mismatchedSizes = results.filter(r => !r.sizesMatch).length;
    
    logger.info({
      at: 'analyzeEntireOrderbookSide',
      message: 'Analysis complete',
      summary: {
        totalLevelsAnalyzed: results.length,
        levelsWithMatchingSizes: matchingSizes,
        levelsWithMismatchedSizes: mismatchedSizes,
        percentageMatching: results.length > 0 ? (matchingSizes / results.length * 100).toFixed(2) + '%' : 'N/A',
      }
    });
    
    // Log details of mismatched levels for investigation
    if (mismatchedSizes > 0) {
      logger.info({
        at: 'analyzeEntireOrderbookSide',
        message: 'Levels with size mismatches:',
        mismatchedLevels: results.filter(r => !r.sizesMatch).map(level => ({
          price: level.price,
          calculatedSize: level.calculatedTotalSize,
          cachedSize: level.orderbookCacheSize,
          difference: (BigInt(level.calculatedTotalSize) - BigInt(level.orderbookCacheSize)).toString(),
          matchingOrdersCount: level.matchingOrdersCount,
          lastUpdated: level.lastUpdatedTime,
        })),
      });
    }
    
  } catch (error) {
    logger.error({
      at: 'analyzeEntireOrderbookSide',
      message: 'Error analyzing orderbook side',
      error,
    });
    process.exit(1);
  } finally {
    client.quit();
  }
}

async function printPriceLevelDetail(
  ticker: string,
  side: string,
  price: string | undefined,
  redisUrl: string,
  analyzeAll: boolean = false,
): Promise<void> {
  logger.info({
    at: 'printPriceLevelDetail',
    message: `Connecting to Redis: ${redisUrl}`,
  });

  const { client, connect } = createRedisClient(redisUrl);
  
  try {
    await connect();
    
    // Normalize side to uppercase
    const normalizedSide = side.toUpperCase();
    if (normalizedSide !== 'BUY' && normalizedSide !== 'SELL') {
      throw new Error('Side must be either BUY or SELL');
    }

    // If analyze-all flag is provided, analyze the entire orderbook side
    if (analyzeAll) {
      // Disconnect client as it will be reconnected in the analyze function
      client.quit();
      return await analyzeEntireOrderbookSide(ticker, normalizedSide, redisUrl);
    }
    
    // Single price level analysis
    if (!price) {
      throw new Error('Price must be provided when not using --analyze-all');
    }

    // Check for related keys in Redis
    logger.info({
      at: 'printPriceLevelDetail',
      message: 'Looking for related orderbook cache keys...',
    });
    
    const orderbookKeys = await getMatchingKeys(`v4/orderbook*${ticker}*`, client);
    logger.info({
      at: 'printPriceLevelDetail',
      message: 'Found orderbook related keys:',
      keys: orderbookKeys,
    });

    const result = await analyzePriceLevel(ticker, normalizedSide, price, client);
    
    // If we found matching orders, get the order data for each
    if (result.matchingOrders.length > 0) {
      // For each matching order, also get the order data (in parallel)
      const orderDataPromises = result.matchingOrders.map(async (order) => {
        const orderId = order.id;
        // Get the order data from the order data cache
        const orderDataKey = `v4/orderData/${orderId}`;
        const orderDataString = await getRedisValue(orderDataKey, client);
        const orderData = orderDataString ? parseOrderData(orderDataString) : null;
        
        return {
          id: order.id,
          price: order.price,
          size: order.size,
          quantums: order.order?.quantums ? 
            new Long(
              order.order.quantums.low, 
              order.order.quantums.high, 
              order.order.quantums.unsigned
            ).toString() : 'N/A',
          subaccountId: order.order?.orderId?.subaccountId ? 
            `${order.order.orderId.subaccountId.owner}:${order.order.orderId.subaccountId.number}` : 'N/A',
          orderData: orderData ? {
            goodTilBlock: orderData.goodTilBlock,
            totalFilledQuantums: orderData.totalFilledQuantums,
            restingOnBook: orderData.restingOnBook,
          } : 'No order data found'
        };
      });

      const ordersWithData = await Promise.all(orderDataPromises);

      // Print details of each matching order with order data
      logger.info({
        at: 'printPriceLevelDetail',
        message: 'Orders at this price level with order data:',
        orders: ordersWithData,
      });
    } else {
      logger.info({
        at: 'printPriceLevelDetail',
        message: 'No orders found at this price level',
      });
    }

  } catch (error) {
    logger.error({
      at: 'printPriceLevelDetail',
      message: 'Error analyzing price level',
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
  side: {
    type: 'string',
    alias: 's',
    description: 'Order side (BUY or SELL)',
    required: true,
  },
  price: {
    type: 'string',
    alias: 'p',
    description: 'Price level to analyze',
    required: false,
  },
  analyzeAll: {
    type: 'boolean',
    alias: 'a',
    description: 'Analyze all price levels in the orderbook side',
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
    await printPriceLevelDetail(args.ticker, args.side, args.price, args.redisUrl, args.analyzeAll);
  } catch (error) {
    logger.error({
      at: 'main',
      message: 'Script failed',
      error,
    });
    process.exit(1);
  }
})(); 