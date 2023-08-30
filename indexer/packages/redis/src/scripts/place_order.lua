-- TODO(CORE-512): add info/resources around caches.
-- Doc: https://www.notion.so/dydx/Indexer-Technical-Spec-a6b15644502048f994c98dee35b96e96#61d5f8ca5117476caab78b3f0691b1d0

-- Key for the value in the orders cache, see `src/caches/orders-cache.ts` for more details
local orderKey = KEYS[1];
-- Key for the value in the orders data cache, see `src/caches/orders-data-cache.ts` for more
-- details
local orderDataKey = KEYS[2];
-- Key for the value in the subaccount order ids cache, see
-- `src/caches/subaccount-order-ids-cache.ts` for more details
local subaccountKey = KEYS[3];
-- Key for the orderExpiry cache, see `src/caches/order-expiry-cache.ts` for more details
local expiryKey = KEYS[4];

-- Order that is being placed, in encoded proto form
local newOrder = ARGV[1];
-- Expiry (good-til-block or sequence number) of the order being placed
local newOrderExpiry = ARGV[2];
-- UUID of the order being placed
local orderId = ARGV[3];
-- Whether the order that is being placed is a short-term order. Need to convert from string to bool
local isShortTermOrder = ARGV[4] == "true";

-- This script returns the following values in an array
-- 1. Was a new order placed or replaced - 1 if an order is placed, 0 if not
-- 2. Was an order replaced - 1 if an order was replaced, 0 if not
-- 3. Total filled size of the old order in quantums - "0" if an order was not replaced
-- 4. Was the old order resting on the book - "true/false", "false" if an order was not replaced
-- 5. Encoded old order - "" if an order was not replaced

-- order data has the format:
-- [good-til-block or sequence number of order]_[total filled]_[true/false, if order is on the book]

-- check if the order exists
local oldOrder = redis.call("get", orderKey);
if not oldOrder then
  -- create the order, set up the order data with total filled size = 0, "false" for whether the
  -- order is resting on the book and add order to list of orders for the subaccount
  redis.call("set", orderKey, newOrder);
  -- refer to above comment on order data format
  redis.call("set", orderDataKey, newOrderExpiry .. "_0_false");
  redis.call("hset", subaccountKey, orderId, 1);
  -- Long-term orders will be on-chain, so we only need to store expiry data for short-term orders
  if isShortTermOrder then
    redis.call("zadd", expiryKey, newOrderExpiry, orderId);
  end
  return {1, 0, "0", "false", ""};
else
  -- refer to above comment on order data format
  local oldOrderData = redis.call("get", orderDataKey);
  local j = string.find(oldOrderData, "_");
  local oldExpiry = string.sub(oldOrderData, 1, j - 1);
  local i = string.find(oldOrderData, "_", j + 1)
  local oldTotalFilledQuantums = string.sub(oldOrderData, j + 1, i - 1);
  local oldRestingOnBook = string.sub(oldOrderData, i + 1);

  -- if the new order has a lower or equal expiry (good-til-block or sequence number) than the order
  -- in the cache, return early
  if tonumber(oldExpiry) >= tonumber(newOrderExpiry) then
    return {0, 0, "0", "false", ""}
  end

  -- update the order if the new order has a greater expiry (good-til-block or sequence number) than
  -- the order in the cache, also update the order data with the new expiry and the total filled
  -- of the older order. As the order is replaced, it is no longer resting on the book, so set the
  -- to "false".
  redis.call("set", orderKey, newOrder)
  -- refer to the above comment on order data format
  redis.call("set", orderDataKey, newOrderExpiry .. "_" .. oldTotalFilledQuantums .. "_false")
  -- Long-term orders will be on-chain, so we only need to store expiry data for short-term orders
  if isShortTermOrder then
    -- The expiry is guaranteed to be different, so overwrite the old one from the expiry cache
    redis.call("zadd", expiryKey, newOrderExpiry, orderId)
  end

  -- if the new order replaced an older order, return data about the old order to use to update
  -- the orderbook
  return {0, 1, oldTotalFilledQuantums, oldRestingOnBook, oldOrder}
end
