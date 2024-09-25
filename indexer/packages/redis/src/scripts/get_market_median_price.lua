-- Key for the sorted set storing price data
local priceCacheKey = KEYS[1]

-- Get all the prices from the sorted set (ascending order)
local prices = redis.call('zrange', priceCacheKey, 0, -1)

-- If no prices are found, return nil
if #prices == 0 then
  return nil
end

-- Calculate the middle index
local middle = math.floor(#prices / 2)

-- Calculate median
if #prices % 2 == 0 then
  -- If even, return both prices, division will be handled in Javascript
  return {prices[middle], prices[middle + 1]}
else
  -- If odd, return the middle element
  return {prices[middle + 1]}
end
