-- Key for the hset of price levels
local hash = KEYS[1]
-- Key for the hset of price levels 'last updated' data
local lastUpdatedHash = KEYS[2]
-- Price level
local level = ARGV[1]

-- This script deletes a price level in the orderbook levels cache if the last updated time is more than 10 seconds in the past.
-- The return value is 1 if a price level was deleted and 0 if a price level was not deleted.

-- Get the current time
local currentTime = tonumber(redis.call("time")[1])

-- Get the last updated time for the level
local lastUpdatedTime = tonumber(redis.call("hget", lastUpdatedHash, level))
if not lastUpdatedTime then
    return 0
end

-- Check if the last updated time is more than 10 seconds in the past
if currentTime - lastUpdatedTime <= 10 then
    return 0
end

-- Delete the level from both hashes
redis.call("hdel", hash, level)
redis.call("hdel", lastUpdatedHash, level)
return 1
