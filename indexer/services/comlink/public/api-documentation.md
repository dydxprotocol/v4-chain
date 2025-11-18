
# Indexer API v1.0.0

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* For **the deployment by DYDX token holders**, use <a href="https://indexer.dydx.trade/v4">https://indexer.dydx.trade/v4</a>
* For **Testnet**, use <a href="https://indexer.v4testnet.dydx.exchange/v4">https://indexer.v4testnet.dydx.exchange/v4</a>

Note: Messages on Indexer WebSocket feeds are typically more recent than data fetched via Indexer's REST API, because the latter is backed by read replicas of the databases that feed the former. Ordinarily this difference is minimal (less than a second), but it might become prolonged under load. Please see [Indexer Architecture](https://dydx.exchange/blog/v4-deep-dive-indexer) for more information.

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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/affiliates/address', params={
  'referralCode': 'string'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/affiliates/address?referralCode=string`,
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

`GET /affiliates/address`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|referralCode|query|string|true|none|

> Example responses

> 200 Response

```json
{
  "address": "string"
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[AffiliateAddressResponse](#schemaaffiliateaddressresponse)|

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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/addresses/{address}/subaccountNumber/{subaccountNumber}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/addresses/{address}/subaccountNumber/{subaccountNumber}`,
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
      "exitPrice": "string",
      "subaccountNumber": 0
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
      "exitPrice": "string",
      "subaccountNumber": 0
    }
  },
  "assetPositions": {
    "property1": {
      "symbol": "string",
      "side": "LONG",
      "size": "string",
      "assetId": "string",
      "subaccountNumber": 0
    },
    "property2": {
      "symbol": "string",
      "side": "LONG",
      "size": "string",
      "assetId": "string",
      "subaccountNumber": 0
    }
  },
  "marginEnabled": true,
  "updatedAtHeight": "string",
  "latestProcessedBlockHeight": "string"
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[SubaccountResponseObject](#schemasubaccountresponseobject)|

<aside class="success">
This operation does not require authentication
</aside>

## GetParentSubaccount

<a id="opIdGetParentSubaccount"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/addresses/{address}/parentSubaccountNumber/{parentSubaccountNumber}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/addresses/{address}/parentSubaccountNumber/{parentSubaccountNumber}`,
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

`GET /addresses/{address}/parentSubaccountNumber/{parentSubaccountNumber}`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|path|string|true|none|
|parentSubaccountNumber|path|number(double)|true|none|

> Example responses

> 200 Response

```json
{
  "address": "string",
  "parentSubaccountNumber": 0,
  "equity": "string",
  "freeCollateral": "string",
  "childSubaccounts": [
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
          "closedAt": null,
          "exitPrice": "string",
          "subaccountNumber": 0
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
          "closedAt": null,
          "exitPrice": "string",
          "subaccountNumber": 0
        }
      },
      "assetPositions": {
        "property1": {
          "symbol": "string",
          "side": "LONG",
          "size": "string",
          "assetId": "string",
          "subaccountNumber": 0
        },
        "property2": {
          "symbol": "string",
          "side": "LONG",
          "size": "string",
          "assetId": "string",
          "subaccountNumber": 0
        }
      },
      "marginEnabled": true,
      "updatedAtHeight": "string",
      "latestProcessedBlockHeight": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[ParentSubaccountResponse](#schemaparentsubaccountresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## RegisterToken

<a id="opIdRegisterToken"></a>

> Code samples

```python
import requests
headers = {
  'Content-Type': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.post(f'{baseURL}/addresses/{address}/registerToken', headers = headers)

print(r.json())

```

```javascript
const inputBody = '{
  "language": "string",
  "token": "string"
}';
const headers = {
  'Content-Type':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/addresses/{address}/registerToken`,
{
  method: 'POST',
  body: inputBody,
  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`POST /addresses/{address}/registerToken`

> Body parameter

```json
{
  "language": "string",
  "token": "string"
}
```

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|path|string|true|none|
|body|body|object|true|none|
|» language|body|string|true|none|
|» token|body|string|true|none|

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|No content|None|

<aside class="success">
This operation does not require authentication
</aside>

## TestNotification

<a id="opIdTestNotification"></a>

> Code samples

```python
import requests

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.post(f'{baseURL}/addresses/{address}/testNotification')

print(r.json())

```

```javascript

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/addresses/{address}/testNotification`,
{
  method: 'POST'

})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`POST /addresses/{address}/testNotification`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|path|string|true|none|

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|No content|None|

<aside class="success">
This operation does not require authentication
</aside>

## GetMetadata

<a id="opIdGetMetadata"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/affiliates/metadata', params={
  'address': 'string'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/affiliates/metadata?address=string`,
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

`GET /affiliates/metadata`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|

> Example responses

> 200 Response

```json
{
  "referralCode": "string",
  "isVolumeEligible": true,
  "isAffiliate": true
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[AffiliateMetadataResponse](#schemaaffiliatemetadataresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## UpdateCode

<a id="opIdUpdateCode"></a>

> Code samples

```python
import requests
headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.post(f'{baseURL}/affiliates/referralCode', headers = headers)

print(r.json())

```

```javascript
const inputBody = '{
  "newCode": "string",
  "address": "string"
}';
const headers = {
  'Content-Type':'application/json',
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/affiliates/referralCode`,
{
  method: 'POST',
  body: inputBody,
  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`POST /affiliates/referralCode`

> Body parameter

```json
{
  "newCode": "string",
  "address": "string"
}
```

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|true|none|
|» newCode|body|string|true|none|
|» address|body|string|true|none|

> Example responses

> 200 Response

```json
{
  "referralCode": "string"
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[CreateReferralCodeResponse](#schemacreatereferralcoderesponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetSnapshot

<a id="opIdGetSnapshot"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/affiliates/snapshot', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/affiliates/snapshot`,
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

`GET /affiliates/snapshot`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|addressFilter|query|array[string]|false|none|
|offset|query|number(double)|false|none|
|limit|query|number(double)|false|none|
|sortByAffiliateEarning|query|boolean|false|none|

> Example responses

> 200 Response

```json
{
  "affiliateList": [
    {
      "affiliateAddress": "string",
      "affiliateReferralCode": "string",
      "affiliateEarnings": 0.1,
      "affiliateReferredTrades": 0.1,
      "affiliateTotalReferredFees": 0.1,
      "affiliateReferredUsers": 0.1,
      "affiliateReferredNetProtocolEarnings": 0.1,
      "affiliateReferredTotalVolume": 0.1,
      "affiliateReferredMakerFees": 0.1,
      "affiliateReferredTakerFees": 0.1,
      "affiliateReferredMakerRebates": 0.1
    }
  ],
  "total": 0.1,
  "currentOffset": 0.1
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[AffiliateSnapshotResponse](#schemaaffiliatesnapshotresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetTotalVolume

<a id="opIdGetTotalVolume"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/affiliates/total_volume', params={
  'address': 'string'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/affiliates/total_volume?address=string`,
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

`GET /affiliates/total_volume`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|

> Example responses

> 200 Response

```json
{
  "totalVolume": 0.1
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[AffiliateTotalVolumeResponse](#schemaaffiliatetotalvolumeresponse)|

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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/assetPositions', params={
  'address': 'string',  'subaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/assetPositions?address=string&subaccountNumber=0.1`,
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
      "assetId": "string",
      "subaccountNumber": 0
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

## GetAssetPositionsForParentSubaccount

<a id="opIdGetAssetPositionsForParentSubaccount"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/assetPositions/parentSubaccountNumber', params={
  'address': 'string',  'parentSubaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/assetPositions/parentSubaccountNumber?address=string&parentSubaccountNumber=0.1`,
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

`GET /assetPositions/parentSubaccountNumber`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|parentSubaccountNumber|query|number(double)|true|none|

> Example responses

> 200 Response

```json
{
  "positions": [
    {
      "symbol": "string",
      "side": "LONG",
      "size": "string",
      "assetId": "string",
      "subaccountNumber": 0
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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/candles/perpetualMarkets/{ticker}', params={
  'resolution': '1MIN'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/candles/perpetualMarkets/{ticker}?resolution=1MIN`,
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
|limit|query|number(double)|false|none|
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
      "trades": 0.1,
      "startingOpenInterest": "string",
      "orderbookMidPriceOpen": "string",
      "orderbookMidPriceClose": "string",
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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/compliance/screen/{address}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/compliance/screen/{address}`,
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

`GET /compliance/screen/{address}`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "status": "COMPLIANT",
  "reason": "MANUAL",
  "updatedAt": "string"
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[ComplianceV2Response](#schemacompliancev2response)|

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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/fills', params={
  'address': 'string',  'subaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/fills?address=string&subaccountNumber=0.1`,
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
|market|query|string|false|none|
|marketType|query|[MarketType](#schemamarkettype)|false|none|
|includeTypes|query|array[string]|false|none|
|excludeTypes|query|array[string]|false|none|
|limit|query|number(double)|false|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|page|query|number(double)|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|marketType|PERPETUAL|
|marketType|SPOT|
|includeTypes|LIMIT|
|includeTypes|LIQUIDATED|
|includeTypes|LIQUIDATION|
|includeTypes|DELEVERAGED|
|includeTypes|OFFSETTING|
|includeTypes|TWAP_SUBORDER|
|excludeTypes|LIMIT|
|excludeTypes|LIQUIDATED|
|excludeTypes|LIQUIDATION|
|excludeTypes|DELEVERAGED|
|excludeTypes|OFFSETTING|
|excludeTypes|TWAP_SUBORDER|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "fills": [
    {
      "id": "string",
      "side": "BUY",
      "liquidity": "TAKER",
      "type": "LIMIT",
      "market": "string",
      "marketType": "PERPETUAL",
      "price": "string",
      "size": "string",
      "fee": "string",
      "affiliateRevShare": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "orderId": "string",
      "clientMetadata": "string",
      "subaccountNumber": 0,
      "builderFee": "string",
      "builderAddress": "string",
      "orderRouterAddress": "string",
      "orderRouterFee": "string",
      "positionSizeBefore": "string",
      "entryPriceBefore": "string",
      "positionSideBefore": "string"
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

## GetFillsForParentSubaccount

<a id="opIdGetFillsForParentSubaccount"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/fills/parentSubaccount', params={
  'address': 'string',  'parentSubaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/fills/parentSubaccount?address=string&parentSubaccountNumber=0.1`,
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

`GET /fills/parentSubaccount`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|parentSubaccountNumber|query|number(double)|true|none|
|includeTypes|query|array[string]|false|none|
|excludeTypes|query|array[string]|false|none|
|limit|query|number(double)|false|none|
|page|query|number(double)|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|includeTypes|LIMIT|
|includeTypes|LIQUIDATED|
|includeTypes|LIQUIDATION|
|includeTypes|DELEVERAGED|
|includeTypes|OFFSETTING|
|includeTypes|TWAP_SUBORDER|
|excludeTypes|LIMIT|
|excludeTypes|LIQUIDATED|
|excludeTypes|LIQUIDATION|
|excludeTypes|DELEVERAGED|
|excludeTypes|OFFSETTING|
|excludeTypes|TWAP_SUBORDER|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "fills": [
    {
      "id": "string",
      "side": "BUY",
      "liquidity": "TAKER",
      "type": "LIMIT",
      "market": "string",
      "marketType": "PERPETUAL",
      "price": "string",
      "size": "string",
      "fee": "string",
      "affiliateRevShare": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "orderId": "string",
      "clientMetadata": "string",
      "subaccountNumber": 0,
      "builderFee": "string",
      "builderAddress": "string",
      "orderRouterAddress": "string",
      "orderRouterFee": "string",
      "positionSizeBefore": "string",
      "entryPriceBefore": "string",
      "positionSideBefore": "string"
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

## GetFundingPayments

<a id="opIdGetFundingPayments"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/fundingPayments', params={
  'address': 'string',  'subaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/fundingPayments?address=string&subaccountNumber=0.1`,
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

`GET /fundingPayments`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|subaccountNumber|query|number(double)|true|none|
|limit|query|number(double)|false|none|
|ticker|query|string|false|none|
|afterOrAt|query|[IsoString](#schemaisostring)|false|none|
|page|query|number(double)|false|none|
|zeroPayments|query|boolean|false|none|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "fundingPayments": [
    {
      "createdAt": "string",
      "createdAtHeight": "string",
      "perpetualId": "string",
      "ticker": "string",
      "oraclePrice": "string",
      "size": "string",
      "side": "string",
      "rate": "string",
      "payment": "string",
      "subaccountNumber": "string",
      "fundingIndex": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[FundingPaymentResponse](#schemafundingpaymentresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetFundingPaymentsForParentSubaccount

<a id="opIdGetFundingPaymentsForParentSubaccount"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/fundingPayments/parentSubaccount', params={
  'address': 'string',  'parentSubaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/fundingPayments/parentSubaccount?address=string&parentSubaccountNumber=0.1`,
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

`GET /fundingPayments/parentSubaccount`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|parentSubaccountNumber|query|number(double)|true|none|
|limit|query|number(double)|false|none|
|afterOrAt|query|[IsoString](#schemaisostring)|false|none|
|page|query|number(double)|false|none|
|zeroPayments|query|boolean|false|none|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "fundingPayments": [
    {
      "createdAt": "string",
      "createdAtHeight": "string",
      "perpetualId": "string",
      "ticker": "string",
      "oraclePrice": "string",
      "size": "string",
      "side": "string",
      "rate": "string",
      "payment": "string",
      "subaccountNumber": "string",
      "fundingIndex": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[FundingPaymentResponse](#schemafundingpaymentresponse)|

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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/height', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/height`,
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

## GetTradingRewards

<a id="opIdGetTradingRewards"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/historicalBlockTradingRewards/{address}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/historicalBlockTradingRewards/{address}`,
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

`GET /historicalBlockTradingRewards/{address}`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|path|string|true|none|
|limit|query|number(double)|false|none|
|startingBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|startingBeforeOrAtHeight|query|string|false|none|

> Example responses

> 200 Response

```json
{
  "rewards": [
    {
      "tradingReward": "string",
      "createdAt": "string",
      "createdAtHeight": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[HistoricalBlockTradingRewardsResponse](#schemahistoricalblocktradingrewardsresponse)|

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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/historicalFunding/{ticker}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/historicalFunding/{ticker}`,
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
|limit|query|number(double)|false|none|
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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/historical-pnl', params={
  'address': 'string',  'subaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/historical-pnl?address=string&subaccountNumber=0.1`,
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
|limit|query|number(double)|false|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|createdOnOrAfterHeight|query|number(double)|false|none|
|createdOnOrAfter|query|[IsoString](#schemaisostring)|false|none|
|page|query|number(double)|false|none|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "historicalPnl": [
    {
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

## GetHistoricalPnlForParentSubaccount

<a id="opIdGetHistoricalPnlForParentSubaccount"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/historical-pnl/parentSubaccountNumber', params={
  'address': 'string',  'parentSubaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/historical-pnl/parentSubaccountNumber?address=string&parentSubaccountNumber=0.1`,
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

`GET /historical-pnl/parentSubaccountNumber`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|parentSubaccountNumber|query|number(double)|true|none|
|limit|query|number(double)|false|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|createdOnOrAfterHeight|query|number(double)|false|none|
|createdOnOrAfter|query|[IsoString](#schemaisostring)|false|none|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "historicalPnl": [
    {
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

## GetAggregations

<a id="opIdGetAggregations"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/historicalTradingRewardAggregations/{address}', params={
  'period': 'DAILY'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/historicalTradingRewardAggregations/{address}?period=DAILY`,
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

`GET /historicalTradingRewardAggregations/{address}`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|path|string|true|none|
|period|query|[TradingRewardAggregationPeriod](#schematradingrewardaggregationperiod)|true|none|
|limit|query|number(double)|false|none|
|startingBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|startingBeforeOrAtHeight|query|string|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|period|DAILY|
|period|WEEKLY|
|period|MONTHLY|

> Example responses

> 200 Response

```json
{
  "rewards": [
    {
      "tradingReward": "string",
      "startedAt": "string",
      "startedAtHeight": "string",
      "endedAt": "string",
      "endedAtHeight": "string",
      "period": "DAILY"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[HistoricalTradingRewardAggregationsResponse](#schemahistoricaltradingrewardaggregationsresponse)|

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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/orderbooks/perpetualMarket/{ticker}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/orderbooks/perpetualMarket/{ticker}`,
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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/orders', params={
  'address': 'string',  'subaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/orders?address=string&subaccountNumber=0.1`,
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
|limit|query|number(double)|false|none|
|ticker|query|string|false|none|
|side|query|[OrderSide](#schemaorderside)|false|none|
|type|query|[OrderType](#schemaordertype)|false|none|
|includeTypes|query|array[string]|false|none|
|excludeTypes|query|array[string]|false|none|
|status|query|array[any]|false|none|
|goodTilBlockBeforeOrAt|query|number(double)|false|none|
|goodTilBlockAfter|query|number(double)|false|none|
|goodTilBlockTimeBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|goodTilBlockTimeAfter|query|[IsoString](#schemaisostring)|false|none|
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
|type|TWAP|
|type|TWAP_SUBORDER|
|includeTypes|LIMIT|
|includeTypes|MARKET|
|includeTypes|STOP_LIMIT|
|includeTypes|STOP_MARKET|
|includeTypes|TRAILING_STOP|
|includeTypes|TAKE_PROFIT|
|includeTypes|TAKE_PROFIT_MARKET|
|includeTypes|TWAP|
|includeTypes|TWAP_SUBORDER|
|excludeTypes|LIMIT|
|excludeTypes|MARKET|
|excludeTypes|STOP_LIMIT|
|excludeTypes|STOP_MARKET|
|excludeTypes|TRAILING_STOP|
|excludeTypes|TAKE_PROFIT|
|excludeTypes|TAKE_PROFIT_MARKET|
|excludeTypes|TWAP|
|excludeTypes|TWAP_SUBORDER|

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
    "builderAddress": "string",
    "feePpm": "string",
    "orderRouterAddress": "string",
    "duration": "string",
    "interval": "string",
    "priceTolerance": "string",
    "timeInForce": "GTT",
    "status": "OPEN",
    "postOnly": true,
    "ticker": "string",
    "updatedAt": "string",
    "updatedAtHeight": "string",
    "subaccountNumber": 0
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
|» builderAddress|string|false|none|none|
|» feePpm|string|false|none|none|
|» orderRouterAddress|string|false|none|none|
|» duration|string|false|none|none|
|» interval|string|false|none|none|
|» priceTolerance|string|false|none|none|
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
|» updatedAt|[IsoString](#schemaisostring)|false|none|none|
|» updatedAtHeight|string|false|none|none|
|» subaccountNumber|integer(int32)|true|none|none|

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
|type|TWAP|
|type|TWAP_SUBORDER|
|timeInForce|GTT|
|timeInForce|FOK|
|timeInForce|IOC|
|*anonymous*|OPEN|
|*anonymous*|FILLED|
|*anonymous*|CANCELED|
|*anonymous*|BEST_EFFORT_CANCELED|
|*anonymous*|UNTRIGGERED|
|*anonymous*|ERROR|
|*anonymous*|BEST_EFFORT_OPENED|

<aside class="success">
This operation does not require authentication
</aside>

## ListOrdersForParentSubaccount

<a id="opIdListOrdersForParentSubaccount"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/orders/parentSubaccountNumber', params={
  'address': 'string',  'parentSubaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/orders/parentSubaccountNumber?address=string&parentSubaccountNumber=0.1`,
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

`GET /orders/parentSubaccountNumber`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|parentSubaccountNumber|query|number(double)|true|none|
|limit|query|number(double)|false|none|
|ticker|query|string|false|none|
|side|query|[OrderSide](#schemaorderside)|false|none|
|type|query|[OrderType](#schemaordertype)|false|none|
|includeTypes|query|array[string]|false|none|
|excludeTypes|query|array[string]|false|none|
|status|query|array[any]|false|none|
|goodTilBlockBeforeOrAt|query|number(double)|false|none|
|goodTilBlockAfter|query|number(double)|false|none|
|goodTilBlockTimeBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|goodTilBlockTimeAfter|query|[IsoString](#schemaisostring)|false|none|
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
|type|TWAP|
|type|TWAP_SUBORDER|
|includeTypes|LIMIT|
|includeTypes|MARKET|
|includeTypes|STOP_LIMIT|
|includeTypes|STOP_MARKET|
|includeTypes|TRAILING_STOP|
|includeTypes|TAKE_PROFIT|
|includeTypes|TAKE_PROFIT_MARKET|
|includeTypes|TWAP|
|includeTypes|TWAP_SUBORDER|
|excludeTypes|LIMIT|
|excludeTypes|MARKET|
|excludeTypes|STOP_LIMIT|
|excludeTypes|STOP_MARKET|
|excludeTypes|TRAILING_STOP|
|excludeTypes|TAKE_PROFIT|
|excludeTypes|TAKE_PROFIT_MARKET|
|excludeTypes|TWAP|
|excludeTypes|TWAP_SUBORDER|

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
    "builderAddress": "string",
    "feePpm": "string",
    "orderRouterAddress": "string",
    "duration": "string",
    "interval": "string",
    "priceTolerance": "string",
    "timeInForce": "GTT",
    "status": "OPEN",
    "postOnly": true,
    "ticker": "string",
    "updatedAt": "string",
    "updatedAtHeight": "string",
    "subaccountNumber": 0
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
|» builderAddress|string|false|none|none|
|» feePpm|string|false|none|none|
|» orderRouterAddress|string|false|none|none|
|» duration|string|false|none|none|
|» interval|string|false|none|none|
|» priceTolerance|string|false|none|none|
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
|» updatedAt|[IsoString](#schemaisostring)|false|none|none|
|» updatedAtHeight|string|false|none|none|
|» subaccountNumber|integer(int32)|true|none|none|

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
|type|TWAP|
|type|TWAP_SUBORDER|
|timeInForce|GTT|
|timeInForce|FOK|
|timeInForce|IOC|
|*anonymous*|OPEN|
|*anonymous*|FILLED|
|*anonymous*|CANCELED|
|*anonymous*|BEST_EFFORT_CANCELED|
|*anonymous*|UNTRIGGERED|
|*anonymous*|ERROR|
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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/orders/{orderId}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/orders/{orderId}`,
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
  "builderAddress": "string",
  "feePpm": "string",
  "orderRouterAddress": "string",
  "duration": "string",
  "interval": "string",
  "priceTolerance": "string",
  "timeInForce": "GTT",
  "status": "OPEN",
  "postOnly": true,
  "ticker": "string",
  "updatedAt": "string",
  "updatedAtHeight": "string",
  "subaccountNumber": 0
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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/perpetualMarkets', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/perpetualMarkets`,
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
|limit|query|number(double)|false|none|
|ticker|query|string|false|none|
|market|query|string|false|none|

> Example responses

> 200 Response

```json
{
  "markets": {
    "property1": {
      "clobPairId": "string",
      "ticker": "string",
      "status": "ACTIVE",
      "oraclePrice": "string",
      "priceChange24H": "string",
      "volume24H": "string",
      "trades24H": 0,
      "nextFundingRate": "string",
      "initialMarginFraction": "string",
      "maintenanceMarginFraction": "string",
      "openInterest": "string",
      "atomicResolution": 0,
      "quantumConversionExponent": 0,
      "tickSize": "string",
      "stepSize": "string",
      "stepBaseQuantums": 0,
      "subticksPerTick": 0,
      "marketType": "CROSS",
      "openInterestLowerCap": "string",
      "openInterestUpperCap": "string",
      "baseOpenInterest": "string",
      "defaultFundingRate1H": "string"
    },
    "property2": {
      "clobPairId": "string",
      "ticker": "string",
      "status": "ACTIVE",
      "oraclePrice": "string",
      "priceChange24H": "string",
      "volume24H": "string",
      "trades24H": 0,
      "nextFundingRate": "string",
      "initialMarginFraction": "string",
      "maintenanceMarginFraction": "string",
      "openInterest": "string",
      "atomicResolution": 0,
      "quantumConversionExponent": 0,
      "tickSize": "string",
      "stepSize": "string",
      "stepBaseQuantums": 0,
      "subticksPerTick": 0,
      "marketType": "CROSS",
      "openInterestLowerCap": "string",
      "openInterestUpperCap": "string",
      "baseOpenInterest": "string",
      "defaultFundingRate1H": "string"
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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/perpetualPositions', params={
  'address': 'string',  'subaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/perpetualPositions?address=string&subaccountNumber=0.1`,
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
|status|query|array[string]|false|none|
|limit|query|number(double)|false|none|
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
      "exitPrice": "string",
      "subaccountNumber": 0
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

## ListPositionsForParentSubaccount

<a id="opIdListPositionsForParentSubaccount"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/perpetualPositions/parentSubaccountNumber', params={
  'address': 'string',  'parentSubaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/perpetualPositions/parentSubaccountNumber?address=string&parentSubaccountNumber=0.1`,
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

`GET /perpetualPositions/parentSubaccountNumber`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|parentSubaccountNumber|query|number(double)|true|none|
|status|query|array[string]|false|none|
|limit|query|number(double)|false|none|
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
      "exitPrice": "string",
      "subaccountNumber": 0
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

## GetPnl

<a id="opIdGetPnl"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/pnl', params={
  'address': 'string',  'subaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/pnl?address=string&subaccountNumber=0.1`,
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

`GET /pnl`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|subaccountNumber|query|number(double)|true|none|
|limit|query|number(double)|false|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|createdOnOrAfterHeight|query|number(double)|false|none|
|createdOnOrAfter|query|[IsoString](#schemaisostring)|false|none|
|page|query|number(double)|false|none|
|daily|query|boolean|false|none|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "pnl": [
    {
      "equity": "string",
      "netTransfers": "string",
      "totalPnl": "string",
      "createdAt": "string",
      "createdAtHeight": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[PnlResponse](#schemapnlresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetPnlForParentSubaccount

<a id="opIdGetPnlForParentSubaccount"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/pnl/parentSubaccountNumber', params={
  'address': 'string',  'parentSubaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/pnl/parentSubaccountNumber?address=string&parentSubaccountNumber=0.1`,
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

`GET /pnl/parentSubaccountNumber`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|parentSubaccountNumber|query|number(double)|true|none|
|limit|query|number(double)|false|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|createdOnOrAfterHeight|query|number(double)|false|none|
|createdOnOrAfter|query|[IsoString](#schemaisostring)|false|none|
|daily|query|boolean|false|none|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "pnl": [
    {
      "equity": "string",
      "netTransfers": "string",
      "totalPnl": "string",
      "createdAt": "string",
      "createdAtHeight": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[PnlResponse](#schemapnlresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## SearchTrader

<a id="opIdSearchTrader"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/trader/search', params={
  'searchParam': 'string'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/trader/search?searchParam=string`,
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

`GET /trader/search`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|searchParam|query|string|true|none|

> Example responses

> 200 Response

```json
{
  "result": {
    "address": "string",
    "subaccountNumber": 0.1,
    "subaccountId": "string",
    "username": "string"
  }
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[TraderSearchResponse](#schematradersearchresponse)|

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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/sparklines', params={
  'timePeriod': 'ONE_DAY'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/sparklines?timePeriod=ONE_DAY`,
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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/time', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/time`,
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
  "epoch": 0.1
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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/trades/perpetualMarket/{ticker}', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/trades/perpetualMarket/{ticker}`,
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
|limit|query|number(double)|false|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|page|query|number(double)|false|none|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "trades": [
    {
      "id": "string",
      "side": "BUY",
      "size": "string",
      "price": "string",
      "type": "LIMIT",
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

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/transfers', params={
  'address': 'string',  'subaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/transfers?address=string&subaccountNumber=0.1`,
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
|limit|query|number(double)|false|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|page|query|number(double)|false|none|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
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

## GetTransfersForParentSubaccount

<a id="opIdGetTransfersForParentSubaccount"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/transfers/parentSubaccountNumber', params={
  'address': 'string',  'parentSubaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/transfers/parentSubaccountNumber?address=string&parentSubaccountNumber=0.1`,
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

`GET /transfers/parentSubaccountNumber`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|none|
|parentSubaccountNumber|query|number(double)|true|none|
|limit|query|number(double)|false|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|
|page|query|number(double)|false|none|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
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
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[ParentSubaccountTransferResponse](#schemaparentsubaccounttransferresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetTransferBetween

<a id="opIdGetTransferBetween"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/transfers/between', params={
  'sourceAddress': 'string',  'sourceSubaccountNumber': '0.1',  'recipientAddress': 'string',  'recipientSubaccountNumber': '0.1'
}, headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/transfers/between?sourceAddress=string&sourceSubaccountNumber=0.1&recipientAddress=string&recipientSubaccountNumber=0.1`,
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

`GET /transfers/between`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|sourceAddress|query|string|true|none|
|sourceSubaccountNumber|query|number(double)|true|none|
|recipientAddress|query|string|true|none|
|recipientSubaccountNumber|query|number(double)|true|none|
|createdBeforeOrAtHeight|query|number(double)|false|none|
|createdBeforeOrAt|query|[IsoString](#schemaisostring)|false|none|

> Example responses

> 200 Response

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "transfersSubset": [
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
  ],
  "totalNetTransfers": "string"
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[TransferBetweenResponse](#schematransferbetweenresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## UploadAddress

<a id="opIdUploadAddress"></a>

> Code samples

```python
import requests
headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.post(f'{baseURL}/turnkey/uploadAddress', headers = headers)

print(r.json())

```

```javascript
const inputBody = '{
  "signature": "string",
  "dydxAddress": "string"
}';
const headers = {
  'Content-Type':'application/json',
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/turnkey/uploadAddress`,
{
  method: 'POST',
  body: inputBody,
  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`POST /turnkey/uploadAddress`

Uploads the dydx address to the turnkey table.

Backend won't have this information when we create account for user since you need signature
to derive dydx address. Just wait for fe to uplaod to kick off the policy setup.

> Body parameter

```json
{
  "signature": "string",
  "dydxAddress": "string"
}
```

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|true|none|
|» signature|body|string|true|none|
|» dydxAddress|body|string|true|none|

> Example responses

> 200 Response

```json
{
  "success": true
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|Inline|

### Response Schema

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» success|boolean|true|none|none|

<aside class="success">
This operation does not require authentication
</aside>

## SignIn

<a id="opIdSignIn"></a>

> Code samples

```python
import requests
headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.post(f'{baseURL}/turnkey/signin', headers = headers)

print(r.json())

```

```javascript
const inputBody = '{
  "signinMethod": "email",
  "userEmail": "string",
  "targetPublicKey": "string",
  "provider": "string",
  "oidcToken": "string",
  "challenge": "string",
  "attestation": {
    "transports": [
      "AUTHENTICATOR_TRANSPORT_BLE"
    ],
    "attestationObject": "string",
    "clientDataJson": "string",
    "credentialId": "string"
  },
  "magicLink": "string"
}';
const headers = {
  'Content-Type':'application/json',
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/turnkey/signin`,
{
  method: 'POST',
  body: inputBody,
  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`POST /turnkey/signin`

> Body parameter

```json
{
  "signinMethod": "email",
  "userEmail": "string",
  "targetPublicKey": "string",
  "provider": "string",
  "oidcToken": "string",
  "challenge": "string",
  "attestation": {
    "transports": [
      "AUTHENTICATOR_TRANSPORT_BLE"
    ],
    "attestationObject": "string",
    "clientDataJson": "string",
    "credentialId": "string"
  },
  "magicLink": "string"
}
```

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[SignInRequest](#schemasigninrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "dydxAddress": "string",
  "organizationId": "string",
  "apiKeyId": "string",
  "userId": "string",
  "session": "string",
  "salt": "string",
  "alreadyExists": true
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[TurnkeyAuthResponse](#schematurnkeyauthresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## AppleLoginRedirect

<a id="opIdAppleLoginRedirect"></a>

> Code samples

```python
import requests
headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.post(f'{baseURL}/turnkey/appleLoginRedirect', headers = headers)

print(r.json())

```

```javascript
const inputBody = '{
  "state": "string",
  "code": "string"
}';
const headers = {
  'Content-Type':'application/json',
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/turnkey/appleLoginRedirect`,
{
  method: 'POST',
  body: inputBody,
  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`POST /turnkey/appleLoginRedirect`

Handles Apple login redirect from Apple's authorization server
Exchanges authorization code for ID token and processes user login/signup

> Body parameter

```json
{
  "state": "string",
  "code": "string"
}
```

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[AppleLoginRedirectRequest](#schemaappleloginredirectrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "success": true,
  "encodedPayload": "string",
  "error": "string"
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[AppleLoginResponse](#schemaappleloginresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetMegavaultHistoricalPnl

<a id="opIdGetMegavaultHistoricalPnl"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/vault/v1/megavault/historicalPnl', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/vault/v1/megavault/historicalPnl`,
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

`GET /vault/v1/megavault/historicalPnl`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|resolution|query|[PnlTickInterval](#schemapnltickinterval)|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|resolution|hour|
|resolution|day|

> Example responses

> 200 Response

```json
{
  "megavaultPnl": [
    {
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
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[MegavaultHistoricalPnlResponse](#schemamegavaulthistoricalpnlresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetVaultsHistoricalPnl

<a id="opIdGetVaultsHistoricalPnl"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/vault/v1/vaults/historicalPnl', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/vault/v1/vaults/historicalPnl`,
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

`GET /vault/v1/vaults/historicalPnl`

### Parameters

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|resolution|query|[PnlTickInterval](#schemapnltickinterval)|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|resolution|hour|
|resolution|day|

> Example responses

> 200 Response

```json
{
  "vaultsPnl": [
    {
      "ticker": "string",
      "historicalPnl": [
        {
          "equity": "string",
          "totalPnl": "string",
          "netTransfers": "string",
          "createdAt": "string",
          "blockHeight": "string",
          "blockTime": "string"
        }
      ]
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[VaultsHistoricalPnlResponse](#schemavaultshistoricalpnlresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetMegavaultPositions

<a id="opIdGetMegavaultPositions"></a>

> Code samples

```python
import requests
headers = {
  'Accept': 'application/json'
}

# For the deployment by DYDX token holders, use
# baseURL = 'https://indexer.dydx.trade/v4'
baseURL = 'https://indexer.v4testnet.dydx.exchange/v4'

r = requests.get(f'{baseURL}/vault/v1/megavault/positions', headers = headers)

print(r.json())

```

```javascript

const headers = {
  'Accept':'application/json'
};

// For the deployment by DYDX token holders, use
// const baseURL = 'https://indexer.dydx.trade/v4';
const baseURL = 'https://indexer.v4testnet.dydx.exchange/v4';

fetch(`${baseURL}/vault/v1/megavault/positions`,
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

`GET /vault/v1/megavault/positions`

> Example responses

> 200 Response

```json
{
  "positions": [
    {
      "ticker": "string",
      "assetPosition": {
        "symbol": "string",
        "side": "LONG",
        "size": "string",
        "assetId": "string",
        "subaccountNumber": 0
      },
      "perpetualPosition": {
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
        "exitPrice": "string",
        "subaccountNumber": 0
      },
      "equity": "string"
    }
  ]
}
```

### Responses

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Ok|[MegavaultPositionResponse](#schemamegavaultpositionresponse)|

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
  "exitPrice": "string",
  "subaccountNumber": 0
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
|subaccountNumber|integer(int32)|true|none|none|

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
    "exitPrice": "string",
    "subaccountNumber": 0
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
    "exitPrice": "string",
    "subaccountNumber": 0
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
  "assetId": "string",
  "subaccountNumber": 0
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|symbol|string|true|none|none|
|side|[PositionSide](#schemapositionside)|true|none|none|
|size|string|true|none|none|
|assetId|string|true|none|none|
|subaccountNumber|integer(int32)|true|none|none|

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
    "assetId": "string",
    "subaccountNumber": 0
  },
  "property2": {
    "symbol": "string",
    "side": "LONG",
    "size": "string",
    "assetId": "string",
    "subaccountNumber": 0
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
      "exitPrice": "string",
      "subaccountNumber": 0
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
      "exitPrice": "string",
      "subaccountNumber": 0
    }
  },
  "assetPositions": {
    "property1": {
      "symbol": "string",
      "side": "LONG",
      "size": "string",
      "assetId": "string",
      "subaccountNumber": 0
    },
    "property2": {
      "symbol": "string",
      "side": "LONG",
      "size": "string",
      "assetId": "string",
      "subaccountNumber": 0
    }
  },
  "marginEnabled": true,
  "updatedAtHeight": "string",
  "latestProcessedBlockHeight": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|address|string|true|none|none|
|subaccountNumber|integer(int32)|true|none|none|
|equity|string|true|none|none|
|freeCollateral|string|true|none|none|
|openPerpetualPositions|[PerpetualPositionsMap](#schemaperpetualpositionsmap)|true|none|none|
|assetPositions|[AssetPositionsMap](#schemaassetpositionsmap)|true|none|none|
|marginEnabled|boolean|true|none|none|
|updatedAtHeight|string|true|none|none|
|latestProcessedBlockHeight|string|true|none|none|

## AddressResponse

<a id="schemaaddressresponse"></a>
<a id="schema_AddressResponse"></a>
<a id="tocSaddressresponse"></a>
<a id="tocsaddressresponse"></a>

```json
{
  "subaccounts": [
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
          "closedAt": null,
          "exitPrice": "string",
          "subaccountNumber": 0
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
          "closedAt": null,
          "exitPrice": "string",
          "subaccountNumber": 0
        }
      },
      "assetPositions": {
        "property1": {
          "symbol": "string",
          "side": "LONG",
          "size": "string",
          "assetId": "string",
          "subaccountNumber": 0
        },
        "property2": {
          "symbol": "string",
          "side": "LONG",
          "size": "string",
          "assetId": "string",
          "subaccountNumber": 0
        }
      },
      "marginEnabled": true,
      "updatedAtHeight": "string",
      "latestProcessedBlockHeight": "string"
    }
  ],
  "totalTradingRewards": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|subaccounts|[[SubaccountResponseObject](#schemasubaccountresponseobject)]|true|none|none|
|totalTradingRewards|string|true|none|none|

## ParentSubaccountResponse

<a id="schemaparentsubaccountresponse"></a>
<a id="schema_ParentSubaccountResponse"></a>
<a id="tocSparentsubaccountresponse"></a>
<a id="tocsparentsubaccountresponse"></a>

```json
{
  "address": "string",
  "parentSubaccountNumber": 0,
  "equity": "string",
  "freeCollateral": "string",
  "childSubaccounts": [
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
          "closedAt": null,
          "exitPrice": "string",
          "subaccountNumber": 0
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
          "closedAt": null,
          "exitPrice": "string",
          "subaccountNumber": 0
        }
      },
      "assetPositions": {
        "property1": {
          "symbol": "string",
          "side": "LONG",
          "size": "string",
          "assetId": "string",
          "subaccountNumber": 0
        },
        "property2": {
          "symbol": "string",
          "side": "LONG",
          "size": "string",
          "assetId": "string",
          "subaccountNumber": 0
        }
      },
      "marginEnabled": true,
      "updatedAtHeight": "string",
      "latestProcessedBlockHeight": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|address|string|true|none|none|
|parentSubaccountNumber|integer(int32)|true|none|none|
|equity|string|true|none|none|
|freeCollateral|string|true|none|none|
|childSubaccounts|[[SubaccountResponseObject](#schemasubaccountresponseobject)]|true|none|none|

## AffiliateMetadataResponse

<a id="schemaaffiliatemetadataresponse"></a>
<a id="schema_AffiliateMetadataResponse"></a>
<a id="tocSaffiliatemetadataresponse"></a>
<a id="tocsaffiliatemetadataresponse"></a>

```json
{
  "referralCode": "string",
  "isVolumeEligible": true,
  "isAffiliate": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|referralCode|string|true|none|none|
|isVolumeEligible|boolean|true|none|none|
|isAffiliate|boolean|true|none|none|

## AffiliateAddressResponse

<a id="schemaaffiliateaddressresponse"></a>
<a id="schema_AffiliateAddressResponse"></a>
<a id="tocSaffiliateaddressresponse"></a>
<a id="tocsaffiliateaddressresponse"></a>

```json
{
  "address": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|address|string|true|none|none|

## CreateReferralCodeResponse

<a id="schemacreatereferralcoderesponse"></a>
<a id="schema_CreateReferralCodeResponse"></a>
<a id="tocScreatereferralcoderesponse"></a>
<a id="tocscreatereferralcoderesponse"></a>

```json
{
  "referralCode": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|referralCode|string|true|none|none|

## AffiliateSnapshotResponseObject

<a id="schemaaffiliatesnapshotresponseobject"></a>
<a id="schema_AffiliateSnapshotResponseObject"></a>
<a id="tocSaffiliatesnapshotresponseobject"></a>
<a id="tocsaffiliatesnapshotresponseobject"></a>

```json
{
  "affiliateAddress": "string",
  "affiliateReferralCode": "string",
  "affiliateEarnings": 0.1,
  "affiliateReferredTrades": 0.1,
  "affiliateTotalReferredFees": 0.1,
  "affiliateReferredUsers": 0.1,
  "affiliateReferredNetProtocolEarnings": 0.1,
  "affiliateReferredTotalVolume": 0.1,
  "affiliateReferredMakerFees": 0.1,
  "affiliateReferredTakerFees": 0.1,
  "affiliateReferredMakerRebates": 0.1
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|affiliateAddress|string|true|none|none|
|affiliateReferralCode|string|true|none|none|
|affiliateEarnings|number(double)|true|none|none|
|affiliateReferredTrades|number(double)|true|none|none|
|affiliateTotalReferredFees|number(double)|true|none|none|
|affiliateReferredUsers|number(double)|true|none|none|
|affiliateReferredNetProtocolEarnings|number(double)|true|none|none|
|affiliateReferredTotalVolume|number(double)|true|none|none|
|affiliateReferredMakerFees|number(double)|true|none|none|
|affiliateReferredTakerFees|number(double)|true|none|none|
|affiliateReferredMakerRebates|number(double)|true|none|none|

## AffiliateSnapshotResponse

<a id="schemaaffiliatesnapshotresponse"></a>
<a id="schema_AffiliateSnapshotResponse"></a>
<a id="tocSaffiliatesnapshotresponse"></a>
<a id="tocsaffiliatesnapshotresponse"></a>

```json
{
  "affiliateList": [
    {
      "affiliateAddress": "string",
      "affiliateReferralCode": "string",
      "affiliateEarnings": 0.1,
      "affiliateReferredTrades": 0.1,
      "affiliateTotalReferredFees": 0.1,
      "affiliateReferredUsers": 0.1,
      "affiliateReferredNetProtocolEarnings": 0.1,
      "affiliateReferredTotalVolume": 0.1,
      "affiliateReferredMakerFees": 0.1,
      "affiliateReferredTakerFees": 0.1,
      "affiliateReferredMakerRebates": 0.1
    }
  ],
  "total": 0.1,
  "currentOffset": 0.1
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|affiliateList|[[AffiliateSnapshotResponseObject](#schemaaffiliatesnapshotresponseobject)]|true|none|none|
|total|number(double)|true|none|none|
|currentOffset|number(double)|true|none|none|

## AffiliateTotalVolumeResponse

<a id="schemaaffiliatetotalvolumeresponse"></a>
<a id="schema_AffiliateTotalVolumeResponse"></a>
<a id="tocSaffiliatetotalvolumeresponse"></a>
<a id="tocsaffiliatetotalvolumeresponse"></a>

```json
{
  "totalVolume": 0.1
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|totalVolume|number(double)¦null|true|none|none|

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
      "assetId": "string",
      "subaccountNumber": 0
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
  "trades": 0.1,
  "startingOpenInterest": "string",
  "orderbookMidPriceOpen": "string",
  "orderbookMidPriceClose": "string",
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
|orderbookMidPriceOpen|string¦null|false|none|none|
|orderbookMidPriceClose|string¦null|false|none|none|
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
      "trades": 0.1,
      "startingOpenInterest": "string",
      "orderbookMidPriceOpen": "string",
      "orderbookMidPriceClose": "string",
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
  "restricted": true,
  "reason": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|restricted|boolean|true|none|none|
|reason|string|false|none|none|

## ComplianceStatus

<a id="schemacompliancestatus"></a>
<a id="schema_ComplianceStatus"></a>
<a id="tocScompliancestatus"></a>
<a id="tocscompliancestatus"></a>

```json
"COMPLIANT"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|COMPLIANT|
|*anonymous*|FIRST_STRIKE_CLOSE_ONLY|
|*anonymous*|FIRST_STRIKE|
|*anonymous*|CLOSE_ONLY|
|*anonymous*|BLOCKED|

## ComplianceReason

<a id="schemacompliancereason"></a>
<a id="schema_ComplianceReason"></a>
<a id="tocScompliancereason"></a>
<a id="tocscompliancereason"></a>

```json
"MANUAL"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|MANUAL|
|*anonymous*|US_GEO|
|*anonymous*|CA_GEO|
|*anonymous*|GB_GEO|
|*anonymous*|SANCTIONED_GEO|
|*anonymous*|COMPLIANCE_PROVIDER|

## ComplianceV2Response

<a id="schemacompliancev2response"></a>
<a id="schema_ComplianceV2Response"></a>
<a id="tocScompliancev2response"></a>
<a id="tocscompliancev2response"></a>

```json
{
  "status": "COMPLIANT",
  "reason": "MANUAL",
  "updatedAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|status|[ComplianceStatus](#schemacompliancestatus)|true|none|none|
|reason|[ComplianceReason](#schemacompliancereason)|false|none|none|
|updatedAt|string|false|none|none|

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
|*anonymous*|LIQUIDATED|
|*anonymous*|LIQUIDATION|
|*anonymous*|DELEVERAGED|
|*anonymous*|OFFSETTING|
|*anonymous*|TWAP_SUBORDER|

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
  "type": "LIMIT",
  "market": "string",
  "marketType": "PERPETUAL",
  "price": "string",
  "size": "string",
  "fee": "string",
  "affiliateRevShare": "string",
  "createdAt": "string",
  "createdAtHeight": "string",
  "orderId": "string",
  "clientMetadata": "string",
  "subaccountNumber": 0,
  "builderFee": "string",
  "builderAddress": "string",
  "orderRouterAddress": "string",
  "orderRouterFee": "string",
  "positionSizeBefore": "string",
  "entryPriceBefore": "string",
  "positionSideBefore": "string"
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
|affiliateRevShare|string|true|none|none|
|createdAt|[IsoString](#schemaisostring)|true|none|none|
|createdAtHeight|string|true|none|none|
|orderId|string|false|none|none|
|clientMetadata|string|false|none|none|
|subaccountNumber|integer(int32)|true|none|none|
|builderFee|string|false|none|none|
|builderAddress|string|false|none|none|
|orderRouterAddress|string|false|none|none|
|orderRouterFee|string|false|none|none|
|positionSizeBefore|string|false|none|none|
|entryPriceBefore|string|false|none|none|
|positionSideBefore|string|false|none|none|

## FillResponse

<a id="schemafillresponse"></a>
<a id="schema_FillResponse"></a>
<a id="tocSfillresponse"></a>
<a id="tocsfillresponse"></a>

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "fills": [
    {
      "id": "string",
      "side": "BUY",
      "liquidity": "TAKER",
      "type": "LIMIT",
      "market": "string",
      "marketType": "PERPETUAL",
      "price": "string",
      "size": "string",
      "fee": "string",
      "affiliateRevShare": "string",
      "createdAt": "string",
      "createdAtHeight": "string",
      "orderId": "string",
      "clientMetadata": "string",
      "subaccountNumber": 0,
      "builderFee": "string",
      "builderAddress": "string",
      "orderRouterAddress": "string",
      "orderRouterFee": "string",
      "positionSizeBefore": "string",
      "entryPriceBefore": "string",
      "positionSideBefore": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pageSize|integer(int32)|false|none|none|
|totalResults|integer(int32)|false|none|none|
|offset|integer(int32)|false|none|none|
|fills|[[FillResponseObject](#schemafillresponseobject)]|true|none|none|

## FundingPaymentResponseObject

<a id="schemafundingpaymentresponseobject"></a>
<a id="schema_FundingPaymentResponseObject"></a>
<a id="tocSfundingpaymentresponseobject"></a>
<a id="tocsfundingpaymentresponseobject"></a>

```json
{
  "createdAt": "string",
  "createdAtHeight": "string",
  "perpetualId": "string",
  "ticker": "string",
  "oraclePrice": "string",
  "size": "string",
  "side": "string",
  "rate": "string",
  "payment": "string",
  "subaccountNumber": "string",
  "fundingIndex": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|createdAt|[IsoString](#schemaisostring)|true|none|none|
|createdAtHeight|string|true|none|none|
|perpetualId|string|true|none|none|
|ticker|string|true|none|none|
|oraclePrice|string|true|none|none|
|size|string|true|none|none|
|side|string|true|none|none|
|rate|string|true|none|none|
|payment|string|true|none|none|
|subaccountNumber|string|true|none|none|
|fundingIndex|string|true|none|none|

## FundingPaymentResponse

<a id="schemafundingpaymentresponse"></a>
<a id="schema_FundingPaymentResponse"></a>
<a id="tocSfundingpaymentresponse"></a>
<a id="tocsfundingpaymentresponse"></a>

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "fundingPayments": [
    {
      "createdAt": "string",
      "createdAtHeight": "string",
      "perpetualId": "string",
      "ticker": "string",
      "oraclePrice": "string",
      "size": "string",
      "side": "string",
      "rate": "string",
      "payment": "string",
      "subaccountNumber": "string",
      "fundingIndex": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pageSize|integer(int32)|false|none|none|
|totalResults|integer(int32)|false|none|none|
|offset|integer(int32)|false|none|none|
|fundingPayments|[[FundingPaymentResponseObject](#schemafundingpaymentresponseobject)]|true|none|none|

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

## HistoricalBlockTradingReward

<a id="schemahistoricalblocktradingreward"></a>
<a id="schema_HistoricalBlockTradingReward"></a>
<a id="tocShistoricalblocktradingreward"></a>
<a id="tocshistoricalblocktradingreward"></a>

```json
{
  "tradingReward": "string",
  "createdAt": "string",
  "createdAtHeight": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|tradingReward|string|true|none|none|
|createdAt|[IsoString](#schemaisostring)|true|none|none|
|createdAtHeight|string|true|none|none|

## HistoricalBlockTradingRewardsResponse

<a id="schemahistoricalblocktradingrewardsresponse"></a>
<a id="schema_HistoricalBlockTradingRewardsResponse"></a>
<a id="tocShistoricalblocktradingrewardsresponse"></a>
<a id="tocshistoricalblocktradingrewardsresponse"></a>

```json
{
  "rewards": [
    {
      "tradingReward": "string",
      "createdAt": "string",
      "createdAtHeight": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|rewards|[[HistoricalBlockTradingReward](#schemahistoricalblocktradingreward)]|true|none|none|

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
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "historicalPnl": [
    {
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
|pageSize|integer(int32)|false|none|none|
|totalResults|integer(int32)|false|none|none|
|offset|integer(int32)|false|none|none|
|historicalPnl|[[PnlTicksResponseObject](#schemapnlticksresponseobject)]|true|none|none|

## TradingRewardAggregationPeriod

<a id="schematradingrewardaggregationperiod"></a>
<a id="schema_TradingRewardAggregationPeriod"></a>
<a id="tocStradingrewardaggregationperiod"></a>
<a id="tocstradingrewardaggregationperiod"></a>

```json
"DAILY"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|DAILY|
|*anonymous*|WEEKLY|
|*anonymous*|MONTHLY|

## HistoricalTradingRewardAggregation

<a id="schemahistoricaltradingrewardaggregation"></a>
<a id="schema_HistoricalTradingRewardAggregation"></a>
<a id="tocShistoricaltradingrewardaggregation"></a>
<a id="tocshistoricaltradingrewardaggregation"></a>

```json
{
  "tradingReward": "string",
  "startedAt": "string",
  "startedAtHeight": "string",
  "endedAt": "string",
  "endedAtHeight": "string",
  "period": "DAILY"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|tradingReward|string|true|none|none|
|startedAt|[IsoString](#schemaisostring)|true|none|none|
|startedAtHeight|string|true|none|none|
|endedAt|[IsoString](#schemaisostring)|false|none|none|
|endedAtHeight|string|false|none|none|
|period|[TradingRewardAggregationPeriod](#schematradingrewardaggregationperiod)|true|none|none|

## HistoricalTradingRewardAggregationsResponse

<a id="schemahistoricaltradingrewardaggregationsresponse"></a>
<a id="schema_HistoricalTradingRewardAggregationsResponse"></a>
<a id="tocShistoricaltradingrewardaggregationsresponse"></a>
<a id="tocshistoricaltradingrewardaggregationsresponse"></a>

```json
{
  "rewards": [
    {
      "tradingReward": "string",
      "startedAt": "string",
      "startedAtHeight": "string",
      "endedAt": "string",
      "endedAtHeight": "string",
      "period": "DAILY"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|rewards|[[HistoricalTradingRewardAggregation](#schemahistoricaltradingrewardaggregation)]|true|none|none|

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
|*anonymous*|ERROR|

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
|*anonymous*|TWAP|
|*anonymous*|TWAP_SUBORDER|

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
  "builderAddress": "string",
  "feePpm": "string",
  "orderRouterAddress": "string",
  "duration": "string",
  "interval": "string",
  "priceTolerance": "string",
  "timeInForce": "GTT",
  "status": "OPEN",
  "postOnly": true,
  "ticker": "string",
  "updatedAt": "string",
  "updatedAtHeight": "string",
  "subaccountNumber": 0
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
|builderAddress|string|false|none|none|
|feePpm|string|false|none|none|
|orderRouterAddress|string|false|none|none|
|duration|string|false|none|none|
|interval|string|false|none|none|
|priceTolerance|string|false|none|none|
|timeInForce|[APITimeInForce](#schemaapitimeinforce)|true|none|none|
|status|[APIOrderStatus](#schemaapiorderstatus)|true|none|none|
|postOnly|boolean|true|none|none|
|ticker|string|true|none|none|
|updatedAt|[IsoString](#schemaisostring)|false|none|none|
|updatedAtHeight|string|false|none|none|
|subaccountNumber|integer(int32)|true|none|none|

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
|*anonymous*|FINAL_SETTLEMENT|

## PerpetualMarketType

<a id="schemaperpetualmarkettype"></a>
<a id="schema_PerpetualMarketType"></a>
<a id="tocSperpetualmarkettype"></a>
<a id="tocsperpetualmarkettype"></a>

```json
"CROSS"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|CROSS|
|*anonymous*|ISOLATED|

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
  "oraclePrice": "string",
  "priceChange24H": "string",
  "volume24H": "string",
  "trades24H": 0,
  "nextFundingRate": "string",
  "initialMarginFraction": "string",
  "maintenanceMarginFraction": "string",
  "openInterest": "string",
  "atomicResolution": 0,
  "quantumConversionExponent": 0,
  "tickSize": "string",
  "stepSize": "string",
  "stepBaseQuantums": 0,
  "subticksPerTick": 0,
  "marketType": "CROSS",
  "openInterestLowerCap": "string",
  "openInterestUpperCap": "string",
  "baseOpenInterest": "string",
  "defaultFundingRate1H": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|clobPairId|string|true|none|none|
|ticker|string|true|none|none|
|status|[PerpetualMarketStatus](#schemaperpetualmarketstatus)|true|none|none|
|oraclePrice|string|true|none|none|
|priceChange24H|string|true|none|none|
|volume24H|string|true|none|none|
|trades24H|integer(int32)|true|none|none|
|nextFundingRate|string|true|none|none|
|initialMarginFraction|string|true|none|none|
|maintenanceMarginFraction|string|true|none|none|
|openInterest|string|true|none|none|
|atomicResolution|integer(int32)|true|none|none|
|quantumConversionExponent|integer(int32)|true|none|none|
|tickSize|string|true|none|none|
|stepSize|string|true|none|none|
|stepBaseQuantums|integer(int32)|true|none|none|
|subticksPerTick|integer(int32)|true|none|none|
|marketType|[PerpetualMarketType](#schemaperpetualmarkettype)|true|none|none|
|openInterestLowerCap|string|false|none|none|
|openInterestUpperCap|string|false|none|none|
|baseOpenInterest|string|true|none|none|
|defaultFundingRate1H|string|false|none|none|

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
      "oraclePrice": "string",
      "priceChange24H": "string",
      "volume24H": "string",
      "trades24H": 0,
      "nextFundingRate": "string",
      "initialMarginFraction": "string",
      "maintenanceMarginFraction": "string",
      "openInterest": "string",
      "atomicResolution": 0,
      "quantumConversionExponent": 0,
      "tickSize": "string",
      "stepSize": "string",
      "stepBaseQuantums": 0,
      "subticksPerTick": 0,
      "marketType": "CROSS",
      "openInterestLowerCap": "string",
      "openInterestUpperCap": "string",
      "baseOpenInterest": "string",
      "defaultFundingRate1H": "string"
    },
    "property2": {
      "clobPairId": "string",
      "ticker": "string",
      "status": "ACTIVE",
      "oraclePrice": "string",
      "priceChange24H": "string",
      "volume24H": "string",
      "trades24H": 0,
      "nextFundingRate": "string",
      "initialMarginFraction": "string",
      "maintenanceMarginFraction": "string",
      "openInterest": "string",
      "atomicResolution": 0,
      "quantumConversionExponent": 0,
      "tickSize": "string",
      "stepSize": "string",
      "stepBaseQuantums": 0,
      "subticksPerTick": 0,
      "marketType": "CROSS",
      "openInterestLowerCap": "string",
      "openInterestUpperCap": "string",
      "baseOpenInterest": "string",
      "defaultFundingRate1H": "string"
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
      "exitPrice": "string",
      "subaccountNumber": 0
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|positions|[[PerpetualPositionResponseObject](#schemaperpetualpositionresponseobject)]|true|none|none|

## PnlResponseObject

<a id="schemapnlresponseobject"></a>
<a id="schema_PnlResponseObject"></a>
<a id="tocSpnlresponseobject"></a>
<a id="tocspnlresponseobject"></a>

```json
{
  "equity": "string",
  "netTransfers": "string",
  "totalPnl": "string",
  "createdAt": "string",
  "createdAtHeight": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|equity|string|true|none|none|
|netTransfers|string|true|none|none|
|totalPnl|string|true|none|none|
|createdAt|string|true|none|none|
|createdAtHeight|string|true|none|none|

## PnlResponse

<a id="schemapnlresponse"></a>
<a id="schema_PnlResponse"></a>
<a id="tocSpnlresponse"></a>
<a id="tocspnlresponse"></a>

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "pnl": [
    {
      "equity": "string",
      "netTransfers": "string",
      "totalPnl": "string",
      "createdAt": "string",
      "createdAtHeight": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pageSize|integer(int32)|false|none|none|
|totalResults|integer(int32)|false|none|none|
|offset|integer(int32)|false|none|none|
|pnl|[[PnlResponseObject](#schemapnlresponseobject)]|true|none|none|

## TraderSearchResponseObject

<a id="schematradersearchresponseobject"></a>
<a id="schema_TraderSearchResponseObject"></a>
<a id="tocStradersearchresponseobject"></a>
<a id="tocstradersearchresponseobject"></a>

```json
{
  "address": "string",
  "subaccountNumber": 0.1,
  "subaccountId": "string",
  "username": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|address|string|true|none|none|
|subaccountNumber|number(double)|true|none|none|
|subaccountId|string|true|none|none|
|username|string|true|none|none|

## TraderSearchResponse

<a id="schematradersearchresponse"></a>
<a id="schema_TraderSearchResponse"></a>
<a id="tocStradersearchresponse"></a>
<a id="tocstradersearchresponse"></a>

```json
{
  "result": {
    "address": "string",
    "subaccountNumber": 0.1,
    "subaccountId": "string",
    "username": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|result|[TraderSearchResponseObject](#schematradersearchresponseobject)|false|none|none|

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
  "epoch": 0.1
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|iso|[IsoString](#schemaisostring)|true|none|none|
|epoch|number(double)|true|none|none|

## TradeType

<a id="schematradetype"></a>
<a id="schema_TradeType"></a>
<a id="tocStradetype"></a>
<a id="tocstradetype"></a>

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
|*anonymous*|LIQUIDATED|
|*anonymous*|DELEVERAGED|
|*anonymous*|TWAP_SUBORDER|

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
  "type": "LIMIT",
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
|type|[TradeType](#schematradetype)|true|none|none|
|createdAt|[IsoString](#schemaisostring)|true|none|none|
|createdAtHeight|string|true|none|none|

## TradeResponse

<a id="schematraderesponse"></a>
<a id="schema_TradeResponse"></a>
<a id="tocStraderesponse"></a>
<a id="tocstraderesponse"></a>

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "trades": [
    {
      "id": "string",
      "side": "BUY",
      "size": "string",
      "price": "string",
      "type": "LIMIT",
      "createdAt": "string",
      "createdAtHeight": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pageSize|integer(int32)|false|none|none|
|totalResults|integer(int32)|false|none|none|
|offset|integer(int32)|false|none|none|
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
|» subaccountNumber|integer(int32)|false|none|none|
|» address|string|true|none|none|
|recipient|object|true|none|none|
|» subaccountNumber|integer(int32)|false|none|none|
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
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
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
|pageSize|integer(int32)|false|none|none|
|totalResults|integer(int32)|false|none|none|
|offset|integer(int32)|false|none|none|
|transfers|[[TransferResponseObject](#schematransferresponseobject)]|true|none|none|

## ParentSubaccountTransferResponse

<a id="schemaparentsubaccounttransferresponse"></a>
<a id="schema_ParentSubaccountTransferResponse"></a>
<a id="tocSparentsubaccounttransferresponse"></a>
<a id="tocsparentsubaccounttransferresponse"></a>

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
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
|pageSize|integer(int32)|false|none|none|
|totalResults|integer(int32)|false|none|none|
|offset|integer(int32)|false|none|none|
|transfers|[[TransferResponseObject](#schematransferresponseobject)]|true|none|none|

## TransferBetweenResponse

<a id="schematransferbetweenresponse"></a>
<a id="schema_TransferBetweenResponse"></a>
<a id="tocStransferbetweenresponse"></a>
<a id="tocstransferbetweenresponse"></a>

```json
{
  "pageSize": 0,
  "totalResults": 0,
  "offset": 0,
  "transfersSubset": [
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
  ],
  "totalNetTransfers": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pageSize|integer(int32)|false|none|none|
|totalResults|integer(int32)|false|none|none|
|offset|integer(int32)|false|none|none|
|transfersSubset|[[TransferResponseObject](#schematransferresponseobject)]|true|none|none|
|totalNetTransfers|string|true|none|none|

## TurnkeyAuthResponse

<a id="schematurnkeyauthresponse"></a>
<a id="schema_TurnkeyAuthResponse"></a>
<a id="tocSturnkeyauthresponse"></a>
<a id="tocsturnkeyauthresponse"></a>

```json
{
  "dydxAddress": "string",
  "organizationId": "string",
  "apiKeyId": "string",
  "userId": "string",
  "session": "string",
  "salt": "string",
  "alreadyExists": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|dydxAddress|string|false|none|none|
|organizationId|string|false|none|none|
|apiKeyId|string|false|none|none|
|userId|string|false|none|none|
|session|string|false|none|none|
|salt|string|true|none|none|
|alreadyExists|boolean|false|none|none|

## SigninMethod

<a id="schemasigninmethod"></a>
<a id="schema_SigninMethod"></a>
<a id="tocSsigninmethod"></a>
<a id="tocssigninmethod"></a>

```json
"email"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|email|
|*anonymous*|social|
|*anonymous*|passkey|

## SignInRequest

<a id="schemasigninrequest"></a>
<a id="schema_SignInRequest"></a>
<a id="tocSsigninrequest"></a>
<a id="tocssigninrequest"></a>

```json
{
  "signinMethod": "email",
  "userEmail": "string",
  "targetPublicKey": "string",
  "provider": "string",
  "oidcToken": "string",
  "challenge": "string",
  "attestation": {
    "transports": [
      "AUTHENTICATOR_TRANSPORT_BLE"
    ],
    "attestationObject": "string",
    "clientDataJson": "string",
    "credentialId": "string"
  },
  "magicLink": "string"
}

```

Request interface for user sign-in operations

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|signinMethod|[SigninMethod](#schemasigninmethod)|true|none|The authentication method to use (EMAIL, SOCIAL, or PASSKEY)|
|userEmail|string|false|none|User's email address (required for EMAIL signin method)|
|targetPublicKey|string|false|none|Target public key for authentication (required for EMAIL and SOCIAL signin methods)|
|provider|string|false|none|OAuth provider name (required for SOCIAL signin method)|
|oidcToken|string|false|none|OIDC token from OAuth provider (required for SOCIAL signin method)|
|challenge|string|false|none|Challenge string for passkey authentication (required for PASSKEY signin method)|
|attestation|object|false|none|Attestation object for passkey authentication (required for PASSKEY signin method)|
|» transports|[string]|true|none|none|
|» attestationObject|string|true|none|none|
|» clientDataJson|string|true|none|none|
|» credentialId|string|true|none|none|
|magicLink|string|false|none|Optional magic link template URL for email authentication|

## AppleLoginResponse

<a id="schemaappleloginresponse"></a>
<a id="schema_AppleLoginResponse"></a>
<a id="tocSappleloginresponse"></a>
<a id="tocsappleloginresponse"></a>

```json
{
  "success": true,
  "encodedPayload": "string",
  "error": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|success|boolean|true|none|none|
|encodedPayload|string|false|none|none|
|error|string|false|none|none|

## AppleLoginRedirectRequest

<a id="schemaappleloginredirectrequest"></a>
<a id="schema_AppleLoginRedirectRequest"></a>
<a id="tocSappleloginredirectrequest"></a>
<a id="tocsappleloginredirectrequest"></a>

```json
{
  "state": "string",
  "code": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|state|string|true|none|none|
|code|string|true|none|none|

## MegavaultHistoricalPnlResponse

<a id="schemamegavaulthistoricalpnlresponse"></a>
<a id="schema_MegavaultHistoricalPnlResponse"></a>
<a id="tocSmegavaulthistoricalpnlresponse"></a>
<a id="tocsmegavaulthistoricalpnlresponse"></a>

```json
{
  "megavaultPnl": [
    {
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
|megavaultPnl|[[PnlTicksResponseObject](#schemapnlticksresponseobject)]|true|none|none|

## PnlTickInterval

<a id="schemapnltickinterval"></a>
<a id="schema_PnlTickInterval"></a>
<a id="tocSpnltickinterval"></a>
<a id="tocspnltickinterval"></a>

```json
"hour"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|hour|
|*anonymous*|day|

## VaultHistoricalPnl

<a id="schemavaulthistoricalpnl"></a>
<a id="schema_VaultHistoricalPnl"></a>
<a id="tocSvaulthistoricalpnl"></a>
<a id="tocsvaulthistoricalpnl"></a>

```json
{
  "ticker": "string",
  "historicalPnl": [
    {
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
|ticker|string|true|none|none|
|historicalPnl|[[PnlTicksResponseObject](#schemapnlticksresponseobject)]|true|none|none|

## VaultsHistoricalPnlResponse

<a id="schemavaultshistoricalpnlresponse"></a>
<a id="schema_VaultsHistoricalPnlResponse"></a>
<a id="tocSvaultshistoricalpnlresponse"></a>
<a id="tocsvaultshistoricalpnlresponse"></a>

```json
{
  "vaultsPnl": [
    {
      "ticker": "string",
      "historicalPnl": [
        {
          "equity": "string",
          "totalPnl": "string",
          "netTransfers": "string",
          "createdAt": "string",
          "blockHeight": "string",
          "blockTime": "string"
        }
      ]
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|vaultsPnl|[[VaultHistoricalPnl](#schemavaulthistoricalpnl)]|true|none|none|

## VaultPosition

<a id="schemavaultposition"></a>
<a id="schema_VaultPosition"></a>
<a id="tocSvaultposition"></a>
<a id="tocsvaultposition"></a>

```json
{
  "ticker": "string",
  "assetPosition": {
    "symbol": "string",
    "side": "LONG",
    "size": "string",
    "assetId": "string",
    "subaccountNumber": 0
  },
  "perpetualPosition": {
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
    "exitPrice": "string",
    "subaccountNumber": 0
  },
  "equity": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|ticker|string|true|none|none|
|assetPosition|[AssetPositionResponseObject](#schemaassetpositionresponseobject)|true|none|none|
|perpetualPosition|[PerpetualPositionResponseObject](#schemaperpetualpositionresponseobject)|false|none|none|
|equity|string|true|none|none|

## MegavaultPositionResponse

<a id="schemamegavaultpositionresponse"></a>
<a id="schema_MegavaultPositionResponse"></a>
<a id="tocSmegavaultpositionresponse"></a>
<a id="tocsmegavaultpositionresponse"></a>

```json
{
  "positions": [
    {
      "ticker": "string",
      "assetPosition": {
        "symbol": "string",
        "side": "LONG",
        "size": "string",
        "assetId": "string",
        "subaccountNumber": 0
      },
      "perpetualPosition": {
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
        "exitPrice": "string",
        "subaccountNumber": 0
      },
      "equity": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|positions|[[VaultPosition](#schemavaultposition)]|true|none|none|

