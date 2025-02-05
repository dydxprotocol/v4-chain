-- KEYS contains the market cache keys
-- ARGV contains the prices for each market and a single timestamp at the end

local numKeys = #KEYS
local numArgs = #ARGV

-- Iterate through each key (market) and corresponding price
for i = 1, numKeys do
  local priceCacheKey = KEYS[i]
  local price = tonumber(ARGV[i])

  -- Validate the price
  if not price then
    return redis.error_reply("Invalid price for key " .. priceCacheKey)
  end

  -- Store the latest price in a simple key-value format without expiration
  redis.call("SET", priceCacheKey, price)
end

return true