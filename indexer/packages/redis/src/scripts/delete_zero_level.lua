-- Key for the hset of price levels
local hash = KEYS[1]
-- Key for the hset of price levels 'last updated' data
local lastUpdatedHash = KEYS[2]
-- Price level
local level = ARGV[1]

-- This script deletes a price level in the orderbook levels cache if the size is zero.
-- The return value is 1 if a price level was deleted and 0 if a price level was not deleted.

local val = redis.call("hget", hash, level)
if not val then
  return 0
end

-- Get the size, if the size is not zero, return.
local size = tonumber(val)
if size ~= 0 then
  return 0
end

-- Delete the level.
redis.call("hdel", hash, level)
redis.call("hdel", lastUpdatedHash, level)
return 1
