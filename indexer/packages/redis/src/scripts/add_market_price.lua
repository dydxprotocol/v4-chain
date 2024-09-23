-- Key for the ZSET storing price data
local priceCacheKey = KEYS[1]
-- Price to be added
local price = tonumber(ARGV[1])
-- Current timestamp
local nowSeconds = tonumber(ARGV[2])
-- Time window (5 seconds)
local fiveSeconds = 5

-- 1. Add the price to the sorted set (score is the current timestamp)
redis.call("zadd", priceCacheKey, nowSeconds, price)

-- 2. Remove any entries older than 5 seconds
local cutoffTime = nowSeconds - fiveSeconds
redis.call("zremrangebyscore", priceCacheKey, "-inf", cutoffTime)

return true