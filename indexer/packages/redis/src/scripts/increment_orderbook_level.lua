-- Key for the hset of price levels
local hash = KEYS[1]
-- Key for the hset of price levels 'last updated' data
local lastUpdatedHash = KEYS[2]
-- Price level
local level = ARGV[1]
-- Increment value
local delta = ARGV[2]

-- This script incrememnts a price level in the orderbook levels cache by the given delta and
-- updates the orderbookLevelsLastUpdated cache.
-- The return value is directly from the hincrby method. A negative result is reset to zero here,
-- atomically with the increment, so that a concurrent update to the same level cannot interleave
-- with the reset and leave the level with quantums that no order backs.
local newValue = redis.call("hincrby", hash, level, delta)
if newValue < 0 then
  redis.call("hset", hash, level, 0)
end

local nowSeconds = redis.call("time")[1]
redis.call("hset", lastUpdatedHash, level, nowSeconds)

return newValue
