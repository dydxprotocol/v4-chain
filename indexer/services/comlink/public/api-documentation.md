
# Indexer API v1.0.0

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* <a href="https://indexer.v4testnet.dydx.exchange/v4">https://indexer.v4testnet.dydx.exchange/v4</a>

# Authentication

# Default

## GetAddress

<a id="opIdGetAddress"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/addresses/{address}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/addresses/{address}',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /addresses/{address}`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|path|string|true|none|

> Example responses

> 200 Response

```json
[
  {
    "address": "string",
    "subaccountNumber": 0,
    "equity": "string",
    "freeCollateral": "string",
    "openPerpetualPositions": {
      "property1": {
        "market": "string",
        "status": "OPEN",
        "side": "LONG",
        "size": "string",
        "maxSize": "string",
        "entryPrice": "string",
        "realizedPnl": "string",
        "createdAt": "string",
        "createdAtHeight": "string",
        "sumOpen": "string",
        "sumClose": "string",
        "netFunding": "string",
        "unrealizedPnl": "string",
        "closedAt": "string",
        "exitPrice": "string"
      },
      "property2": {
        "market": "string",
        "status": "OPEN",
        "side": "LONG",
        "size": "string",
        "maxSize": "string",
        "entryPrice": "string",
        "realizedPnl": "string",
        "createdAt": "string",
        "createdAtHeight": "string",
        "sumOpen": "string",
        "sumClose": "string",
        "netFunding": "string",
        "unrealizedPnl": "string",
        "closedAt": "string",
        "exitPrice": "string"
      }
    },
    "assetPositions": {
      "property1": {
        "symbol": "string",
        "side": "LONG",
        "size": "string",
        "assetId": "string"
      },
      "property2": {
        "symbol": "string",
        "side": "LONG",
        "size": "string",
        "assetId": "string"
      }
    },
    "marginEnabled": true
  }
]
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|Inline|

### Response Schema

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[SubaccountResponseObject](#schemasubaccountresponseobject)]|false|none|none|
|» address|string|true|none|none|
|» subaccountNumber|number(double)|true|none|none|
|» equity|string|true|none|none|
|» freeCollateral|string|true|none|none|
|» openPerpetualPositions|[PerpetualPositionsMap](#schemaperpetualpositionsmap)|true|none|none|
|»» **additionalProperties**|[PerpetualPositionResponseObject](#schemaperpetualpositionresponseobject)|false|none|none|
|»»» market|string|true|none|none|
|»»» status|[PerpetualPositionStatus](#schemaperpetualpositionstatus)|true|none|none|
|»»» side|[PositionSide](#schemapositionside)|true|none|none|
|»»» size|string|true|none|none|
|»»» maxSize|string|true|none|none|
|»»» entryPrice|string|true|none|none|
|»»» realizedPnl|string|true|none|none|
|»»» createdAt|[IsoString](#schemaisostring)|true|none|none|
|»»» createdAtHeight|string|true|none|none|
|»»» sumOpen|string|true|none|none|
|»»» sumClose|string|true|none|none|
|»»» netFunding|string|true|none|none|
|»»» unrealizedPnl|string|true|none|none|
|»»» closedAt|[IsoString](#schemaisostring)¦null|false|none|none|
|»»» exitPrice|string¦null|false|none|none|
|» assetPositions|[AssetPositionsMap](#schemaassetpositionsmap)|true|none|none|
|»» **additionalProperties**|[AssetPositionResponseObject](#schemaassetpositionresponseobject)|false|none|none|
|»»» symbol|string|true|none|none|
|»»» side|[PositionSide](#schemapositionside)|true|none|none|
|»»» size|string|true|none|none|
|»»» assetId|string|true|none|none|
|» marginEnabled|boolean|true|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|status|OPEN|
|status|CLOSED|
|status|LIQUIDATED|
|side|LONG|
|side|SHORT|
|side|LONG|
|side|SHORT|

<aside class="success">
This operation does not require authentication
</aside>

## GetSubaccount

<a id="opIdGetSubaccount"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/addresses/{address}/subaccountNumber/{subaccountNumber}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/addresses/{address}/subaccountNumber/{subaccountNumber}',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /addresses/{address}/subaccountNumber/{subaccountNumber}`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|path|string|true|none|
|subaccountNumber|path|number(double)|true|none|

> Example responses

> 200 Response

```json
{
  "address": "string",
  "subaccountNumber": 0,
  "equity": "string",
  "freeCollateral": "string",
  "openPerpetualPositions": {
    "property1": {
      "market": "string",
      "status": "OPEN",
      "side": "LONG",
      "size": "string",
      "maxSize": "string",
      "entryPrice": "string",
      "realizedPnl": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "sumOpen": "string",
      "sumClose": "string",
      "netFunding": "string",
      "unrealizedPnl": "string",
      "closedAt": "string",
      "exitPrice": "string"
    },
    "property2": {
      "market": "string",
      "status": "OPEN",
      "side": "LONG",
      "size": "string",
      "maxSize": "string",
      "entryPrice": "string",
      "realizedPnl": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "sumOpen": "string",
      "sumClose": "string",
      "netFunding": "string",
      "unrealizedPnl": "string",
      "closedAt": "string",
      "exitPrice": "string"
    }
  },
  "assetPositions": {
    "property1": {
      "symbol": "string",
      "side": "LONG",
      "size": "string",
      "assetId": "string"
    },
    "property2": {
      "symbol": "string",
      "side": "LONG",
      "size": "string",
      "assetId": "string"
    }
  },
  "marginEnabled": true
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[SubaccountResponseObject](#schemasubaccountresponseobject)|

<aside class="success">
This operation does not require authentication
</aside>

## GetAssetPositions

<a id="opIdGetAssetPositions"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/assetPositions', params={
  'address': 'string',  'subaccountNumber': '0'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/assetPositions?address=string&subaccountNumber=0',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /assetPositions`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|subaccountNumber|query|number(double)|true|none|

> Example responses

> 200 Response

```json
{
  "positions": [
    {
      "symbol": "string",
      "side": "LONG",
      "size": "string",
      "assetId": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[AssetPositionResponse](#schemaassetpositionresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetCandles

<a id="opIdGetCandles"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/candles/perpetualMarkets/{ticker}', params={
  'resolution': '1MIN',  'limit': '0'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/candles/perpetualMarkets/{ticker}?resolution=1MIN&limit=0',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /candles/perpetualMarkets/{ticker}`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|ticker|path|string|true|none|
|resolution|query|[CandleResolution](#schemacandleresolution)|true|none|
|limit|query|number(double)|true|none|
|fromISO|query|string|false|none|
|toISO|query|string|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|resolution|1MIN|
|resolution|5MINS|
|resolution|15MINS|
|resolution|30MINS|
|resolution|1HOUR|
|resolution|4HOURS|
|resolution|1DAY|

> Example responses

> 200 Response

```json
{
  "candles": [
    {
      "startedAt": "string",
      "ticker": "string",
      "resolution": "1MIN",
      "low": "string",
      "high": "string",
      "open": "string",
      "close": "string",
      "baseTokenVolume": "string",
      "usdVolume": "string",
      "trades": 0,
      "startingOpenInterest": "string",
      "id": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[CandleResponse](#schemacandleresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Screen

<a id="opIdScreen"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/screen', params={
  'address': 'string'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/screen?address=string',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /screen`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|

> Example responses

> 200 Response

```json
{
  "restricted": true
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[ComplianceResponse](#schemacomplianceresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetFills

<a id="opIdGetFills"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/fills', params={
  'address': 'string',  'subaccountNumber': '0',  'market': 'string',  'marketType': 'PERPETUAL',  'limit': '0'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/fills?address=string&subaccountNumber=0&market=string&marketType=PERPETUAL&limit=0',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /fills`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|subaccountNumber|query|number(double)|true|none|
|market|query|string|true|none|
|marketType|query|[MarketType](#schemamarkettype)|true|none|
|limit|query|number(double)|true|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|marketType|PERPETUAL|
|marketType|SPOT|

> Example responses

> 200 Response

```json
{
  "fills": [
    {
      "id": "string",
      "side": "BUY",
      "liquidity": "TAKER",
      "type": "MARKET",
      "market": "string",
      "marketType": "PERPETUAL",
      "price": "string",
      "size": "string",
      "fee": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "orderId": "string",
      "clientMetadata": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[FillResponse](#schemafillresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetHeight

<a id="opIdGetHeight"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/height', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/height',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /height`

> Example responses

> 200 Response

```json
{
  "height": "string",
  "time": "string"
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[HeightResponse](#schemaheightresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetHistoricalFunding

<a id="opIdGetHistoricalFunding"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/historicalFunding/{ticker}', params={
  'limit': '0'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/historicalFunding/{ticker}?limit=0',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /historicalFunding/{ticker}`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|ticker|path|string|true|none|
|limit|query|number(double)|true|none|
|effectiveBeforeOrAtHeight|query|number(double)|false|none|
|effectiveBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|

> Example responses

> 200 Response

```json
{
  "historicalFunding": [
    {
      "ticker": "string",
      "rate": "string",
      "price": "string",
      "effectiveAt": "string",
      "effectiveAtHeight": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[HistoricalFundingResponse](#schemahistoricalfundingresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetHistoricalPnl

<a id="opIdGetHistoricalPnl"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/historical-pnl', params={
  'address': 'string',  'subaccountNumber': '0',  'limit': '0'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/historical-pnl?address=string&subaccountNumber=0&limit=0',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /historical-pnl`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|subaccountNumber|query|number(double)|true|none|
|limit|query|number(double)|true|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|createdOnOrAfterHeight|query|number(double)|false|none|
|createdOnOrAfter|query|[IsoString](#schemaisostring)|false|none|

> Example responses

> 200 Response

```json
{
  "historicalPnl": [
    {
      "id": "string",
      "subaccountId": "string",
      "equity": "string",
      "totalPnl": "string",
      "netTransfers": "string",
      "createdAt": "string",
      "blockHeight": "string",
      "blockTime": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[HistoricalPnlResponse](#schemahistoricalpnlresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetPerpetualMarket

<a id="opIdGetPerpetualMarket"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/orderbooks/perpetualMarket/{ticker}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/orderbooks/perpetualMarket/{ticker}',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /orderbooks/perpetualMarket/{ticker}`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|ticker|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "bids": [
    {
      "price": "string",
      "size": "string"
    }
  ],
  "asks": [
    {
      "price": "string",
      "size": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[OrderbookResponseObject](#schemaorderbookresponseobject)|

<aside class="success">
This operation does not require authentication
</aside>

## ListOrders

<a id="opIdListOrders"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/orders', params={
  'address': 'string',  'subaccountNumber': '0',  'limit': '0'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/orders?address=string&subaccountNumber=0&limit=0',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /orders`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|subaccountNumber|query|number(double)|true|none|
|limit|query|number(double)|true|none|
|ticker|query|string|false|none|
|side|query|[OrderSide](#schemaorderside)|false|none|
|type|query|[OrderType](#schemaordertype)|false|none|
|status|query|array[any]|false|none|
|goodTilBlockBeforeOrAt|query|number(double)|false|none|
|goodTilBlockTimeBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|returnLatestOrders|query|boolean|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|side|BUY|
|side|SELL|
|type|LIMIT|
|type|MARKET|
|type|STOP_LIMIT|
|type|STOP_MARKET|
|type|TRAILING_STOP|
|type|TAKE_PROFIT|
|type|TAKE_PROFIT_MARKET|
|type|HARD_TRADE|
|type|FAILED_HARD_TRADE|
|type|TRANSFER_PLACEHOLDER|

> Example responses

> 200 Response

```json
[
  {
    "id": "string",
    "subaccountId": "string",
    "clientId": "string",
    "clobPairId": "string",
    "side": "BUY",
    "size": "string",
    "totalFilled": "string",
    "price": "string",
    "type": "LIMIT",
    "reduceOnly": true,
    "orderFlags": "string",
    "goodTilBlock": "string",
    "goodTilBlockTime": "string",
    "createdAtHeight": "string",
    "clientMetadata": "string",
    "triggerPrice": "string",
    "timeInForce": "GTT",
    "status": "OPEN",
    "postOnly": true,
    "ticker": "string"
  }
]
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|Inline|

### Response Schema

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[OrderResponseObject](#schemaorderresponseobject)]|false|none|none|
|» id|string|true|none|none|
|» subaccountId|string|true|none|none|
|» clientId|string|true|none|none|
|» clobPairId|string|true|none|none|
|» side|[OrderSide](#schemaorderside)|true|none|none|
|» size|string|true|none|none|
|» totalFilled|string|true|none|none|
|» price|string|true|none|none|
|» type|[OrderType](#schemaordertype)|true|none|none|
|» reduceOnly|boolean|true|none|none|
|» orderFlags|string|true|none|none|
|» goodTilBlock|string|false|none|none|
|» goodTilBlockTime|string|false|none|none|
|» createdAtHeight|string|false|none|none|
|» clientMetadata|string|true|none|none|
|» triggerPrice|string|false|none|none|
|» timeInForce|[APITimeInForce](#schemaapitimeinforce)|true|none|none|
|» status|any|true|none|none|

*anyOf*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|[OrderStatus](#schemaorderstatus)|false|none|none|

*or*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|[BestEffortOpenedStatus](#schemabesteffortopenedstatus)|false|none|none|

*continued*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» postOnly|boolean|true|none|none|
|» ticker|string|true|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|side|BUY|
|side|SELL|
|type|LIMIT|
|type|MARKET|
|type|STOP_LIMIT|
|type|STOP_MARKET|
|type|TRAILING_STOP|
|type|TAKE_PROFIT|
|type|TAKE_PROFIT_MARKET|
|type|HARD_TRADE|
|type|FAILED_HARD_TRADE|
|type|TRANSFER_PLACEHOLDER|
|timeInForce|GTT|
|timeInForce|FOK|
|timeInForce|IOC|
|*anonymous*|OPEN|
|*anonymous*|FILLED|
|*anonymous*|CANCELED|
|*anonymous*|BEST_EFFORT_CANCELED|
|*anonymous*|UNTRIGGERED|
|*anonymous*|BEST_EFFORT_OPENED|

<aside class="success">
This operation does not require authentication
</aside>

## GetOrder

<a id="opIdGetOrder"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/orders/{orderId}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/orders/{orderId}',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /orders/{orderId}`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|orderId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "id": "string",
  "subaccountId": "string",
  "clientId": "string",
  "clobPairId": "string",
  "side": "BUY",
  "size": "string",
  "totalFilled": "string",
  "price": "string",
  "type": "LIMIT",
  "reduceOnly": true,
  "orderFlags": "string",
  "goodTilBlock": "string",
  "goodTilBlockTime": "string",
  "createdAtHeight": "string",
  "clientMetadata": "string",
  "triggerPrice": "string",
  "timeInForce": "GTT",
  "status": "OPEN",
  "postOnly": true,
  "ticker": "string"
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[OrderResponseObject](#schemaorderresponseobject)|

<aside class="success">
This operation does not require authentication
</aside>

## ListPerpetualMarkets

<a id="opIdListPerpetualMarkets"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/perpetualMarkets', params={
  'limit': '0'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/perpetualMarkets?limit=0',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /perpetualMarkets`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|limit|query|number(double)|true|none|
|ticker|query|string|false|none|

> Example responses

> 200 Response

```json
{
  "markets": {
    "property1": {
      "clobPairId": "string",
      "ticker": "string",
      "status": "ACTIVE",
      "lastPrice": "string",
      "oraclePrice": "string",
      "priceChange24H": "string",
      "volume24H": "string",
      "trades24H": 0,
      "nextFundingRate": "string",
      "initialMarginFraction": "string",
      "maintenanceMarginFraction": "string",
      "basePositionNotional": "string",
      "openInterest": "string",
      "atomicResolution": 0,
      "quantumConversionExponent": 0,
      "tickSize": "string",
      "stepSize": "string",
      "stepBaseQuantums": 0,
      "subticksPerTick": 0
    },
    "property2": {
      "clobPairId": "string",
      "ticker": "string",
      "status": "ACTIVE",
      "lastPrice": "string",
      "oraclePrice": "string",
      "priceChange24H": "string",
      "volume24H": "string",
      "trades24H": 0,
      "nextFundingRate": "string",
      "initialMarginFraction": "string",
      "maintenanceMarginFraction": "string",
      "basePositionNotional": "string",
      "openInterest": "string",
      "atomicResolution": 0,
      "quantumConversionExponent": 0,
      "tickSize": "string",
      "stepSize": "string",
      "stepBaseQuantums": 0,
      "subticksPerTick": 0
    }
  }
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[PerpetualMarketResponse](#schemaperpetualmarketresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ListPositions

<a id="opIdListPositions"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/perpetualPositions', params={
  'address': 'string',  'subaccountNumber': '0',  'status': [
  "OPEN"
],  'limit': '0'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/perpetualPositions?address=string&subaccountNumber=0&status=OPEN&limit=0',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /perpetualPositions`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|subaccountNumber|query|number(double)|true|none|
|status|query|array[string]|true|none|
|limit|query|number(double)|true|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|status|OPEN|
|status|CLOSED|
|status|LIQUIDATED|

> Example responses

> 200 Response

```json
{
  "positions": [
    {
      "market": "string",
      "status": "OPEN",
      "side": "LONG",
      "size": "string",
      "maxSize": "string",
      "entryPrice": "string",
      "realizedPnl": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "sumOpen": "string",
      "sumClose": "string",
      "netFunding": "string",
      "unrealizedPnl": "string",
      "closedAt": "string",
      "exitPrice": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[PerpetualPositionResponse](#schemaperpetualpositionresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Get

<a id="opIdGet"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/sparklines', params={
  'timePeriod': 'ONE_DAY'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/sparklines?timePeriod=ONE_DAY',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /sparklines`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|timePeriod|query|[SparklineTimePeriod](#schemasparklinetimeperiod)|true|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|timePeriod|ONE_DAY|
|timePeriod|SEVEN_DAYS|

> Example responses

> 200 Response

```json
{
  "property1": [
    "string"
  ],
  "property2": [
    "string"
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[SparklineResponseObject](#schemasparklineresponseobject)|

<aside class="success">
This operation does not require authentication
</aside>

## GetTime

<a id="opIdGetTime"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/time', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/time',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /time`

> Example responses

> 200 Response

```json
{
  "iso": "string",
  "epoch": 0
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[TimeResponse](#schematimeresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetTrades

<a id="opIdGetTrades"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/trades/perpetualMarket/{ticker}', params={
  'limit': '0'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/trades/perpetualMarket/{ticker}?limit=0',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /trades/perpetualMarket/{ticker}`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|ticker|path|string|true|none|
|limit|query|number(double)|true|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|

> Example responses

> 200 Response

```json
{
  "trades": [
    {
      "id": "string",
      "side": "BUY",
      "size": "string",
      "price": "string",
      "createdAt": "string",
      "createdAtHeight": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[TradeResponse](#schematraderesponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetTransfers

<a id="opIdGetTransfers"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('https://indexer.v4testnet.dydx.exchange/v4/transfers', params={
  'address': 'string',  'subaccountNumber': '0',  'limit': '0'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://indexer.v4testnet.dydx.exchange/v4/transfers?address=string&subaccountNumber=0&limit=0',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /transfers`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|subaccountNumber|query|number(double)|true|none|
|limit|query|number(double)|true|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|

> Example responses

> 200 Response

```json
{
  "transfers": [
    {
      "id": "string",
      "sender": {
        "subaccountNumber": 0,
        "address": "string"
      },
      "recipient": {
        "subaccountNumber": 0,
        "address": "string"
      },
      "size": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "symbol": "string",
      "type": "TRANSFER_IN",
      "transactionHash": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[TransferResponse](#schematransferresponse)|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

## PerpetualPositionStatus

<a id="schemaperpetualpositionstatus"></a>
<a id="schema_PerpetualPositionStatus"></a>
<a id="tocSperpetualpositionstatus"></a>
<a id="tocsperpetualpositionstatus"></a>

```json
"OPEN"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|OPEN|
|*anonymous*|CLOSED|
|*anonymous*|LIQUIDATED|

## PositionSide

<a id="schemapositionside"></a>
<a id="schema_PositionSide"></a>
<a id="tocSpositionside"></a>
<a id="tocspositionside"></a>

```json
"LONG"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|LONG|
|*anonymous*|SHORT|

## IsoString

<a id="schemaisostring"></a>
<a id="schema_IsoString"></a>
<a id="tocSisostring"></a>
<a id="tocsisostring"></a>

```json
"string"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

## PerpetualPositionResponseObject

<a id="schemaperpetualpositionresponseobject"></a>
<a id="schema_PerpetualPositionResponseObject"></a>
<a id="tocSperpetualpositionresponseobject"></a>
<a id="tocsperpetualpositionresponseobject"></a>

```json
{
  "market": "string",
  "status": "OPEN",
  "side": "LONG",
  "size": "string",
  "maxSize": "string",
  "entryPrice": "string",
  "realizedPnl": "string",
  "createdAt": "string",
  "createdAtHeight": "string",
  "sumOpen": "string",
  "sumClose": "string",
  "netFunding": "string",
  "unrealizedPnl": "string",
  "closedAt": "string",
  "exitPrice": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|market|string|true|none|none|
|status|[PerpetualPositionStatus](#schemaperpetualpositionstatus)|true|none|none|
|side|[PositionSide](#schemapositionside)|true|none|none|
|size|string|true|none|none|
|maxSize|string|true|none|none|
|entryPrice|string|true|none|none|
|realizedPnl|string|true|none|none|
|createdAt|[IsoString](#schemaisostring)|true|none|none|
|createdAtHeight|string|true|none|none|
|sumOpen|string|true|none|none|
|sumClose|string|true|none|none|
|netFunding|string|true|none|none|
|unrealizedPnl|string|true|none|none|
|closedAt|[IsoString](#schemaisostring)¦null|false|none|none|
|exitPrice|string¦null|false|none|none|

## PerpetualPositionsMap

<a id="schemaperpetualpositionsmap"></a>
<a id="schema_PerpetualPositionsMap"></a>
<a id="tocSperpetualpositionsmap"></a>
<a id="tocsperpetualpositionsmap"></a>

```json
{
  "property1": {
    "market": "string",
    "status": "OPEN",
    "side": "LONG",
    "size": "string",
    "maxSize": "string",
    "entryPrice": "string",
    "realizedPnl": "string",
    "createdAt": "string",
    "createdAtHeight": "string",
    "sumOpen": "string",
    "sumClose": "string",
    "netFunding": "string",
    "unrealizedPnl": "string",
    "closedAt": "string",
    "exitPrice": "string"
  },
  "property2": {
    "market": "string",
    "status": "OPEN",
    "side": "LONG",
    "size": "string",
    "maxSize": "string",
    "entryPrice": "string",
    "realizedPnl": "string",
    "createdAt": "string",
    "createdAtHeight": "string",
    "sumOpen": "string",
    "sumClose": "string",
    "netFunding": "string",
    "unrealizedPnl": "string",
    "closedAt": "string",
    "exitPrice": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|**additionalProperties**|[PerpetualPositionResponseObject](#schemaperpetualpositionresponseobject)|false|none|none|

## AssetPositionResponseObject

<a id="schemaassetpositionresponseobject"></a>
<a id="schema_AssetPositionResponseObject"></a>
<a id="tocSassetpositionresponseobject"></a>
<a id="tocsassetpositionresponseobject"></a>

```json
{
  "symbol": "string",
  "side": "LONG",
  "size": "string",
  "assetId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|symbol|string|true|none|none|
|side|[PositionSide](#schemapositionside)|true|none|none|
|size|string|true|none|none|
|assetId|string|true|none|none|

## AssetPositionsMap

<a id="schemaassetpositionsmap"></a>
<a id="schema_AssetPositionsMap"></a>
<a id="tocSassetpositionsmap"></a>
<a id="tocsassetpositionsmap"></a>

```json
{
  "property1": {
    "symbol": "string",
    "side": "LONG",
    "size": "string",
    "assetId": "string"
  },
  "property2": {
    "symbol": "string",
    "side": "LONG",
    "size": "string",
    "assetId": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|**additionalProperties**|[AssetPositionResponseObject](#schemaassetpositionresponseobject)|false|none|none|

## SubaccountResponseObject

<a id="schemasubaccountresponseobject"></a>
<a id="schema_SubaccountResponseObject"></a>
<a id="tocSsubaccountresponseobject"></a>
<a id="tocssubaccountresponseobject"></a>

```json
{
  "address": "string",
  "subaccountNumber": 0,
  "equity": "string",
  "freeCollateral": "string",
  "openPerpetualPositions": {
    "property1": {
      "market": "string",
      "status": "OPEN",
      "side": "LONG",
      "size": "string",
      "maxSize": "string",
      "entryPrice": "string",
      "realizedPnl": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "sumOpen": "string",
      "sumClose": "string",
      "netFunding": "string",
      "unrealizedPnl": "string",
      "closedAt": "string",
      "exitPrice": "string"
    },
    "property2": {
      "market": "string",
      "status": "OPEN",
      "side": "LONG",
      "size": "string",
      "maxSize": "string",
      "entryPrice": "string",
      "realizedPnl": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "sumOpen": "string",
      "sumClose": "string",
      "netFunding": "string",
      "unrealizedPnl": "string",
      "closedAt": "string",
      "exitPrice": "string"
    }
  },
  "assetPositions": {
    "property1": {
      "symbol": "string",
      "side": "LONG",
      "size": "string",
      "assetId": "string"
    },
    "property2": {
      "symbol": "string",
      "side": "LONG",
      "size": "string",
      "assetId": "string"
    }
  },
  "marginEnabled": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|address|string|true|none|none|
|subaccountNumber|number(double)|true|none|none|
|equity|string|true|none|none|
|freeCollateral|string|true|none|none|
|openPerpetualPositions|[PerpetualPositionsMap](#schemaperpetualpositionsmap)|true|none|none|
|assetPositions|[AssetPositionsMap](#schemaassetpositionsmap)|true|none|none|
|marginEnabled|boolean|true|none|none|

## AssetPositionResponse

<a id="schemaassetpositionresponse"></a>
<a id="schema_AssetPositionResponse"></a>
<a id="tocSassetpositionresponse"></a>
<a id="tocsassetpositionresponse"></a>

```json
{
  "positions": [
    {
      "symbol": "string",
      "side": "LONG",
      "size": "string",
      "assetId": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|positions|[[AssetPositionResponseObject](#schemaassetpositionresponseobject)]|true|none|none|

## CandleResolution

<a id="schemacandleresolution"></a>
<a id="schema_CandleResolution"></a>
<a id="tocScandleresolution"></a>
<a id="tocscandleresolution"></a>

```json
"1MIN"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|1MIN|
|*anonymous*|5MINS|
|*anonymous*|15MINS|
|*anonymous*|30MINS|
|*anonymous*|1HOUR|
|*anonymous*|4HOURS|
|*anonymous*|1DAY|

## CandleResponseObject

<a id="schemacandleresponseobject"></a>
<a id="schema_CandleResponseObject"></a>
<a id="tocScandleresponseobject"></a>
<a id="tocscandleresponseobject"></a>

```json
{
  "startedAt": "string",
  "ticker": "string",
  "resolution": "1MIN",
  "low": "string",
  "high": "string",
  "open": "string",
  "close": "string",
  "baseTokenVolume": "string",
  "usdVolume": "string",
  "trades": 0,
  "startingOpenInterest": "string",
  "id": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|startedAt|[IsoString](#schemaisostring)|true|none|none|
|ticker|string|true|none|none|
|resolution|[CandleResolution](#schemacandleresolution)|true|none|none|
|low|string|true|none|none|
|high|string|true|none|none|
|open|string|true|none|none|
|close|string|true|none|none|
|baseTokenVolume|string|true|none|none|
|usdVolume|string|true|none|none|
|trades|number(double)|true|none|none|
|startingOpenInterest|string|true|none|none|
|id|string|true|none|none|

## CandleResponse

<a id="schemacandleresponse"></a>
<a id="schema_CandleResponse"></a>
<a id="tocScandleresponse"></a>
<a id="tocscandleresponse"></a>

```json
{
  "candles": [
    {
      "startedAt": "string",
      "ticker": "string",
      "resolution": "1MIN",
      "low": "string",
      "high": "string",
      "open": "string",
      "close": "string",
      "baseTokenVolume": "string",
      "usdVolume": "string",
      "trades": 0,
      "startingOpenInterest": "string",
      "id": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|candles|[[CandleResponseObject](#schemacandleresponseobject)]|true|none|none|

## ComplianceResponse

<a id="schemacomplianceresponse"></a>
<a id="schema_ComplianceResponse"></a>
<a id="tocScomplianceresponse"></a>
<a id="tocscomplianceresponse"></a>

```json
{
  "restricted": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|restricted|boolean|true|none|none|

## OrderSide

<a id="schemaorderside"></a>
<a id="schema_OrderSide"></a>
<a id="tocSorderside"></a>
<a id="tocsorderside"></a>

```json
"BUY"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|BUY|
|*anonymous*|SELL|

## Liquidity

<a id="schemaliquidity"></a>
<a id="schema_Liquidity"></a>
<a id="tocSliquidity"></a>
<a id="tocsliquidity"></a>

```json
"TAKER"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|TAKER|
|*anonymous*|MAKER|

## FillType

<a id="schemafilltype"></a>
<a id="schema_FillType"></a>
<a id="tocSfilltype"></a>
<a id="tocsfilltype"></a>

```json
"MARKET"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|MARKET|
|*anonymous*|LIMIT|
|*anonymous*|LIQUIDATED|
|*anonymous*|LIQUIDATION|

## MarketType

<a id="schemamarkettype"></a>
<a id="schema_MarketType"></a>
<a id="tocSmarkettype"></a>
<a id="tocsmarkettype"></a>

```json
"PERPETUAL"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|PERPETUAL|
|*anonymous*|SPOT|

## FillResponseObject

<a id="schemafillresponseobject"></a>
<a id="schema_FillResponseObject"></a>
<a id="tocSfillresponseobject"></a>
<a id="tocsfillresponseobject"></a>

```json
{
  "id": "string",
  "side": "BUY",
  "liquidity": "TAKER",
  "type": "MARKET",
  "market": "string",
  "marketType": "PERPETUAL",
  "price": "string",
  "size": "string",
  "fee": "string",
  "createdAt": "string",
  "createdAtHeight": "string",
  "orderId": "string",
  "clientMetadata": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|true|none|none|
|side|[OrderSide](#schemaorderside)|true|none|none|
|liquidity|[Liquidity](#schemaliquidity)|true|none|none|
|type|[FillType](#schemafilltype)|true|none|none|
|market|string|true|none|none|
|marketType|[MarketType](#schemamarkettype)|true|none|none|
|price|string|true|none|none|
|size|string|true|none|none|
|fee|string|true|none|none|
|createdAt|[IsoString](#schemaisostring)|true|none|none|
|createdAtHeight|string|true|none|none|
|orderId|string|false|none|none|
|clientMetadata|string|false|none|none|

## FillResponse

<a id="schemafillresponse"></a>
<a id="schema_FillResponse"></a>
<a id="tocSfillresponse"></a>
<a id="tocsfillresponse"></a>

```json
{
  "fills": [
    {
      "id": "string",
      "side": "BUY",
      "liquidity": "TAKER",
      "type": "MARKET",
      "market": "string",
      "marketType": "PERPETUAL",
      "price": "string",
      "size": "string",
      "fee": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "orderId": "string",
      "clientMetadata": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|fills|[[FillResponseObject](#schemafillresponseobject)]|true|none|none|

## HeightResponse

<a id="schemaheightresponse"></a>
<a id="schema_HeightResponse"></a>
<a id="tocSheightresponse"></a>
<a id="tocsheightresponse"></a>

```json
{
  "height": "string",
  "time": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|height|string|true|none|none|
|time|[IsoString](#schemaisostring)|true|none|none|

## HistoricalFundingResponseObject

<a id="schemahistoricalfundingresponseobject"></a>
<a id="schema_HistoricalFundingResponseObject"></a>
<a id="tocShistoricalfundingresponseobject"></a>
<a id="tocshistoricalfundingresponseobject"></a>

```json
{
  "ticker": "string",
  "rate": "string",
  "price": "string",
  "effectiveAt": "string",
  "effectiveAtHeight": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|ticker|string|true|none|none|
|rate|string|true|none|none|
|price|string|true|none|none|
|effectiveAt|[IsoString](#schemaisostring)|true|none|none|
|effectiveAtHeight|string|true|none|none|

## HistoricalFundingResponse

<a id="schemahistoricalfundingresponse"></a>
<a id="schema_HistoricalFundingResponse"></a>
<a id="tocShistoricalfundingresponse"></a>
<a id="tocshistoricalfundingresponse"></a>

```json
{
  "historicalFunding": [
    {
      "ticker": "string",
      "rate": "string",
      "price": "string",
      "effectiveAt": "string",
      "effectiveAtHeight": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|historicalFunding|[[HistoricalFundingResponseObject](#schemahistoricalfundingresponseobject)]|true|none|none|

## PnlTicksResponseObject

<a id="schemapnlticksresponseobject"></a>
<a id="schema_PnlTicksResponseObject"></a>
<a id="tocSpnlticksresponseobject"></a>
<a id="tocspnlticksresponseobject"></a>

```json
{
  "id": "string",
  "subaccountId": "string",
  "equity": "string",
  "totalPnl": "string",
  "netTransfers": "string",
  "createdAt": "string",
  "blockHeight": "string",
  "blockTime": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|true|none|none|
|subaccountId|string|true|none|none|
|equity|string|true|none|none|
|totalPnl|string|true|none|none|
|netTransfers|string|true|none|none|
|createdAt|string|true|none|none|
|blockHeight|string|true|none|none|
|blockTime|[IsoString](#schemaisostring)|true|none|none|

## HistoricalPnlResponse

<a id="schemahistoricalpnlresponse"></a>
<a id="schema_HistoricalPnlResponse"></a>
<a id="tocShistoricalpnlresponse"></a>
<a id="tocshistoricalpnlresponse"></a>

```json
{
  "historicalPnl": [
    {
      "id": "string",
      "subaccountId": "string",
      "equity": "string",
      "totalPnl": "string",
      "netTransfers": "string",
      "createdAt": "string",
      "blockHeight": "string",
      "blockTime": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|historicalPnl|[[PnlTicksResponseObject](#schemapnlticksresponseobject)]|true|none|none|

## OrderbookResponsePriceLevel

<a id="schemaorderbookresponsepricelevel"></a>
<a id="schema_OrderbookResponsePriceLevel"></a>
<a id="tocSorderbookresponsepricelevel"></a>
<a id="tocsorderbookresponsepricelevel"></a>

```json
{
  "price": "string",
  "size": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|price|string|true|none|none|
|size|string|true|none|none|

## OrderbookResponseObject

<a id="schemaorderbookresponseobject"></a>
<a id="schema_OrderbookResponseObject"></a>
<a id="tocSorderbookresponseobject"></a>
<a id="tocsorderbookresponseobject"></a>

```json
{
  "bids": [
    {
      "price": "string",
      "size": "string"
    }
  ],
  "asks": [
    {
      "price": "string",
      "size": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|bids|[[OrderbookResponsePriceLevel](#schemaorderbookresponsepricelevel)]|true|none|none|
|asks|[[OrderbookResponsePriceLevel](#schemaorderbookresponsepricelevel)]|true|none|none|

## APITimeInForce

<a id="schemaapitimeinforce"></a>
<a id="schema_APITimeInForce"></a>
<a id="tocSapitimeinforce"></a>
<a id="tocsapitimeinforce"></a>

```json
"GTT"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|GTT|
|*anonymous*|FOK|
|*anonymous*|IOC|

## OrderStatus

<a id="schemaorderstatus"></a>
<a id="schema_OrderStatus"></a>
<a id="tocSorderstatus"></a>
<a id="tocsorderstatus"></a>

```json
"OPEN"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|OPEN|
|*anonymous*|FILLED|
|*anonymous*|CANCELED|
|*anonymous*|BEST_EFFORT_CANCELED|
|*anonymous*|UNTRIGGERED|

## BestEffortOpenedStatus

<a id="schemabesteffortopenedstatus"></a>
<a id="schema_BestEffortOpenedStatus"></a>
<a id="tocSbesteffortopenedstatus"></a>
<a id="tocsbesteffortopenedstatus"></a>

```json
"BEST_EFFORT_OPENED"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|BEST_EFFORT_OPENED|

## APIOrderStatus

<a id="schemaapiorderstatus"></a>
<a id="schema_APIOrderStatus"></a>
<a id="tocSapiorderstatus"></a>
<a id="tocsapiorderstatus"></a>

```json
"OPEN"

```

### Properties

anyOf

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[OrderStatus](#schemaorderstatus)|false|none|none|

or

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[BestEffortOpenedStatus](#schemabesteffortopenedstatus)|false|none|none|

## OrderType

<a id="schemaordertype"></a>
<a id="schema_OrderType"></a>
<a id="tocSordertype"></a>
<a id="tocsordertype"></a>

```json
"LIMIT"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|LIMIT|
|*anonymous*|MARKET|
|*anonymous*|STOP_LIMIT|
|*anonymous*|STOP_MARKET|
|*anonymous*|TRAILING_STOP|
|*anonymous*|TAKE_PROFIT|
|*anonymous*|TAKE_PROFIT_MARKET|
|*anonymous*|HARD_TRADE|
|*anonymous*|FAILED_HARD_TRADE|
|*anonymous*|TRANSFER_PLACEHOLDER|

## OrderResponseObject

<a id="schemaorderresponseobject"></a>
<a id="schema_OrderResponseObject"></a>
<a id="tocSorderresponseobject"></a>
<a id="tocsorderresponseobject"></a>

```json
{
  "id": "string",
  "subaccountId": "string",
  "clientId": "string",
  "clobPairId": "string",
  "side": "BUY",
  "size": "string",
  "totalFilled": "string",
  "price": "string",
  "type": "LIMIT",
  "reduceOnly": true,
  "orderFlags": "string",
  "goodTilBlock": "string",
  "goodTilBlockTime": "string",
  "createdAtHeight": "string",
  "clientMetadata": "string",
  "triggerPrice": "string",
  "timeInForce": "GTT",
  "status": "OPEN",
  "postOnly": true,
  "ticker": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|true|none|none|
|subaccountId|string|true|none|none|
|clientId|string|true|none|none|
|clobPairId|string|true|none|none|
|side|[OrderSide](#schemaorderside)|true|none|none|
|size|string|true|none|none|
|totalFilled|string|true|none|none|
|price|string|true|none|none|
|type|[OrderType](#schemaordertype)|true|none|none|
|reduceOnly|boolean|true|none|none|
|orderFlags|string|true|none|none|
|goodTilBlock|string|false|none|none|
|goodTilBlockTime|string|false|none|none|
|createdAtHeight|string|false|none|none|
|clientMetadata|string|true|none|none|
|triggerPrice|string|false|none|none|
|timeInForce|[APITimeInForce](#schemaapitimeinforce)|true|none|none|
|status|[APIOrderStatus](#schemaapiorderstatus)|true|none|none|
|postOnly|boolean|true|none|none|
|ticker|string|true|none|none|

## PerpetualMarketStatus

<a id="schemaperpetualmarketstatus"></a>
<a id="schema_PerpetualMarketStatus"></a>
<a id="tocSperpetualmarketstatus"></a>
<a id="tocsperpetualmarketstatus"></a>

```json
"ACTIVE"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|ACTIVE|
|*anonymous*|PAUSED|
|*anonymous*|CANCEL_ONLY|
|*anonymous*|POST_ONLY|
|*anonymous*|INITIALIZING|

## PerpetualMarketResponseObject

<a id="schemaperpetualmarketresponseobject"></a>
<a id="schema_PerpetualMarketResponseObject"></a>
<a id="tocSperpetualmarketresponseobject"></a>
<a id="tocsperpetualmarketresponseobject"></a>

```json
{
  "clobPairId": "string",
  "ticker": "string",
  "status": "ACTIVE",
  "lastPrice": "string",
  "oraclePrice": "string",
  "priceChange24H": "string",
  "volume24H": "string",
  "trades24H": 0,
  "nextFundingRate": "string",
  "initialMarginFraction": "string",
  "maintenanceMarginFraction": "string",
  "basePositionNotional": "string",
  "openInterest": "string",
  "atomicResolution": 0,
  "quantumConversionExponent": 0,
  "tickSize": "string",
  "stepSize": "string",
  "stepBaseQuantums": 0,
  "subticksPerTick": 0
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|clobPairId|string|true|none|none|
|ticker|string|true|none|none|
|status|[PerpetualMarketStatus](#schemaperpetualmarketstatus)|true|none|none|
|lastPrice|string|true|none|none|
|oraclePrice|string|true|none|none|
|priceChange24H|string|true|none|none|
|volume24H|string|true|none|none|
|trades24H|number(double)|true|none|none|
|nextFundingRate|string|true|none|none|
|initialMarginFraction|string|true|none|none|
|maintenanceMarginFraction|string|true|none|none|
|basePositionNotional|string|true|none|none|
|openInterest|string|true|none|none|
|atomicResolution|number(double)|true|none|none|
|quantumConversionExponent|number(double)|true|none|none|
|tickSize|string|true|none|none|
|stepSize|string|true|none|none|
|stepBaseQuantums|number(double)|true|none|none|
|subticksPerTick|number(double)|true|none|none|

## PerpetualMarketResponse

<a id="schemaperpetualmarketresponse"></a>
<a id="schema_PerpetualMarketResponse"></a>
<a id="tocSperpetualmarketresponse"></a>
<a id="tocsperpetualmarketresponse"></a>

```json
{
  "markets": {
    "property1": {
      "clobPairId": "string",
      "ticker": "string",
      "status": "ACTIVE",
      "lastPrice": "string",
      "oraclePrice": "string",
      "priceChange24H": "string",
      "volume24H": "string",
      "trades24H": 0,
      "nextFundingRate": "string",
      "initialMarginFraction": "string",
      "maintenanceMarginFraction": "string",
      "basePositionNotional": "string",
      "openInterest": "string",
      "atomicResolution": 0,
      "quantumConversionExponent": 0,
      "tickSize": "string",
      "stepSize": "string",
      "stepBaseQuantums": 0,
      "subticksPerTick": 0
    },
    "property2": {
      "clobPairId": "string",
      "ticker": "string",
      "status": "ACTIVE",
      "lastPrice": "string",
      "oraclePrice": "string",
      "priceChange24H": "string",
      "volume24H": "string",
      "trades24H": 0,
      "nextFundingRate": "string",
      "initialMarginFraction": "string",
      "maintenanceMarginFraction": "string",
      "basePositionNotional": "string",
      "openInterest": "string",
      "atomicResolution": 0,
      "quantumConversionExponent": 0,
      "tickSize": "string",
      "stepSize": "string",
      "stepBaseQuantums": 0,
      "subticksPerTick": 0
    }
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|markets|object|true|none|none|
|» **additionalProperties**|[PerpetualMarketResponseObject](#schemaperpetualmarketresponseobject)|false|none|none|

## PerpetualPositionResponse

<a id="schemaperpetualpositionresponse"></a>
<a id="schema_PerpetualPositionResponse"></a>
<a id="tocSperpetualpositionresponse"></a>
<a id="tocsperpetualpositionresponse"></a>

```json
{
  "positions": [
    {
      "market": "string",
      "status": "OPEN",
      "side": "LONG",
      "size": "string",
      "maxSize": "string",
      "entryPrice": "string",
      "realizedPnl": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "sumOpen": "string",
      "sumClose": "string",
      "netFunding": "string",
      "unrealizedPnl": "string",
      "closedAt": "string",
      "exitPrice": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|positions|[[PerpetualPositionResponseObject](#schemaperpetualpositionresponseobject)]|true|none|none|

## SparklineResponseObject

<a id="schemasparklineresponseobject"></a>
<a id="schema_SparklineResponseObject"></a>
<a id="tocSsparklineresponseobject"></a>
<a id="tocssparklineresponseobject"></a>

```json
{
  "property1": [
    "string"
  ],
  "property2": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|**additionalProperties**|[string]|false|none|none|

## SparklineTimePeriod

<a id="schemasparklinetimeperiod"></a>
<a id="schema_SparklineTimePeriod"></a>
<a id="tocSsparklinetimeperiod"></a>
<a id="tocssparklinetimeperiod"></a>

```json
"ONE_DAY"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|ONE_DAY|
|*anonymous*|SEVEN_DAYS|

## TimeResponse

<a id="schematimeresponse"></a>
<a id="schema_TimeResponse"></a>
<a id="tocStimeresponse"></a>
<a id="tocstimeresponse"></a>

```json
{
  "iso": "string",
  "epoch": 0
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|iso|[IsoString](#schemaisostring)|true|none|none|
|epoch|number(double)|true|none|none|

## TradeResponseObject

<a id="schematraderesponseobject"></a>
<a id="schema_TradeResponseObject"></a>
<a id="tocStraderesponseobject"></a>
<a id="tocstraderesponseobject"></a>

```json
{
  "id": "string",
  "side": "BUY",
  "size": "string",
  "price": "string",
  "createdAt": "string",
  "createdAtHeight": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|true|none|none|
|side|[OrderSide](#schemaorderside)|true|none|none|
|size|string|true|none|none|
|price|string|true|none|none|
|createdAt|[IsoString](#schemaisostring)|true|none|none|
|createdAtHeight|string|true|none|none|

## TradeResponse

<a id="schematraderesponse"></a>
<a id="schema_TradeResponse"></a>
<a id="tocStraderesponse"></a>
<a id="tocstraderesponse"></a>

```json
{
  "trades": [
    {
      "id": "string",
      "side": "BUY",
      "size": "string",
      "price": "string",
      "createdAt": "string",
      "createdAtHeight": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|trades|[[TradeResponseObject](#schematraderesponseobject)]|true|none|none|

## TransferType

<a id="schematransfertype"></a>
<a id="schema_TransferType"></a>
<a id="tocStransfertype"></a>
<a id="tocstransfertype"></a>

```json
"TRANSFER_IN"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|TRANSFER_IN|
|*anonymous*|TRANSFER_OUT|
|*anonymous*|DEPOSIT|
|*anonymous*|WITHDRAWAL|

## TransferResponseObject

<a id="schematransferresponseobject"></a>
<a id="schema_TransferResponseObject"></a>
<a id="tocStransferresponseobject"></a>
<a id="tocstransferresponseobject"></a>

```json
{
  "id": "string",
  "sender": {
    "subaccountNumber": 0,
    "address": "string"
  },
  "recipient": {
    "subaccountNumber": 0,
    "address": "string"
  },
  "size": "string",
  "createdAt": "string",
  "createdAtHeight": "string",
  "symbol": "string",
  "type": "TRANSFER_IN",
  "transactionHash": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|true|none|none|
|sender|object|true|none|none|
|» subaccountNumber|number(double)|false|none|none|
|» address|string|true|none|none|
|recipient|object|true|none|none|
|» subaccountNumber|number(double)|false|none|none|
|» address|string|true|none|none|
|size|string|true|none|none|
|createdAt|string|true|none|none|
|createdAtHeight|string|true|none|none|
|symbol|string|true|none|none|
|type|[TransferType](#schematransfertype)|true|none|none|
|transactionHash|string|true|none|none|

## TransferResponse

<a id="schematransferresponse"></a>
<a id="schema_TransferResponse"></a>
<a id="tocStransferresponse"></a>
<a id="tocstransferresponse"></a>

```json
{
  "transfers": [
    {
      "id": "string",
      "sender": {
        "subaccountNumber": 0,
        "address": "string"
      },
      "recipient": {
        "subaccountNumber": 0,
        "address": "string"
      },
      "size": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "symbol": "string",
      "type": "TRANSFER_IN",
      "transactionHash": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|transfers|[[TransferResponseObject](#schematransferresponseobject)]|true|none|none|

