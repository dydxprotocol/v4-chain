-- KEYS contains the market cache keys
-- ARGV contains the prices for each market and a single timestamp at the end

local numKeys = #KEYS
local numArgs = #ARGV

-- Get the timestamp from the last argument
local timestamp = tonumber(ARGV[numArgs])

-- Time window (60 seconds)
local sixtySeconds = 60

-- Validate the timestamp
if not timestamp then
  return redis.error_reply("Invalid timestamp")
end

-- Calculate the cutoff time for removing old prices
local cutoffTime = timestamp - sixtySeconds

-- Iterate through each key (market) and corresponding price
for i = 1, numKeys do
  local priceCacheKey = KEYS[i]
  local price = tonumber(ARGV[i])

  -- Validate the price
  if not price then
    return redis.error_reply("Invalid price for key " .. priceCacheKey)
  end

  -- Add the price to the sorted set with the current timestamp as the score
  redis.call("ZADD", priceCacheKey, timestamp, price)

  -- Remove entries older than the cutoff time (older than 60 seconds)
  redis.call("ZREMRANGEBYSCORE", priceCacheKey, "-inf", cutoffTime)
end

return true
