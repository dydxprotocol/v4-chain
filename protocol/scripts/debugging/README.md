# Proto Event Deserializer

This script is designed to deserialize ender proto events from the dYdX v4-chain indexer for debugging purposes.

## Purpose

The `deserialize_proto.go` script helps debug ender events by deserializing proto message data that appears in CloudWatch logs. When the ender service processes events, the raw proto bytes are often logged as JSON objects with `dataBytes` maps. This script reconstructs those bytes and deserializes them into readable proto message structures.

## How to Get the JSON Data

1. **Access CloudWatch Logs**: Navigate to the ender service logs in AWS CloudWatch
2. **Find Event Data**: Look for log entries containing event data with `dataBytes` fields
3. **Copy JSON Structure**: Copy the JSON structure that contains the `indexerTendermintBlock` with events
4. **Save to File**: Save the JSON data to a file (e.g., `events.json`)

The expected JSON structure should look like:
```json
{
  "indexerTendermintBlock": {
    "events": [
      {
        "dataBytes": {
          "0": "10",
          "1": "8", 
          "2": "1",
          ...
        },
        "subtype": "subaccount_update",
        "eventIndex": 0,
        "transactionIndex": 0,
        "version": 1
      }
    ]
  }
}
```

## Usage

### Option 1: Run directly with Go
```bash
make run JSON=<json_file_path>
```

### Option 2: Build and run
```bash
make build
./deserialize_proto <json_file_path>
```

### Examples
```bash
# Deserialize events from a specific file
make run JSON=events.json

# Or after building
./deserialize_proto cloudwatch_events.json
```

## Output

The script will output:
- The file being processed
- For each event:
  - Event metadata (index, transaction index, subtype)
  - Message type (e.g., `SubaccountUpdateEventV1`, `OrderFillEventV1`)
  - Pretty-formatted proto message fields with proper indentation

Example output:
```
Reading events from: events.json

--- Event Index: 23, Transaction Index: 0, Subtype: order_fill ---
Message Type: OrderFillEventV1
  maker_order:<order_id:<subaccount_id:<owner:"dydx1..." number:1 > client_id:1622718733 clob_pair_id:27 > side:SIDE_BUY quantums:2266000000 subticks:710000000 good_til_block:1044 >
  order:<order_id:<subaccount_id:<owner:"dydx130f..." number:7 > client_id:507483497 clob_pair_id:27 > side:SIDE_SELL quantums:2271000000 subticks:709000000 good_til_block:1048 >
  fill_amount:2261000000
  maker_fee:-779
  taker_fee:3191
  total_filled_maker:2261000000
  total_filled_taker:2271000000
```

## Supported Event Types

The script supports all major indexer event types:
- `subaccount_update` → `SubaccountUpdateEventV1`
- `order_fill` → `OrderFillEventV1`
- `transfer` → `TransferEventV1`
- `deleveraging` → `DeleveragingEventV1`
- `stateful_order` → `StatefulOrderEventV1`
- `funding` → `FundingEventV1`
- `market` → `MarketEventV1`
- `asset_create` → `AssetCreateEventV1`
- `perpetual_market_create` → `PerpetualMarketCreateEventV1`
- `liquidity_tier` → `LiquidityTierUpsertEventV1`
- `update_clob_pair` → `UpdateClobPairEventV1`
- `update_perpetual` → `UpdatePerpetualEventV1`
- `trading_rewards` → `TradingRewardsEventV1`
- `open_interest_update` → `OpenInterestUpdateEventV1`
- `register_affiliate` → `RegisterAffiliateEventV1`
- `upsert_vault` → `UpsertVaultEventV1`

## Error Handling

The script includes robust error handling:
- Validates JSON file exists and is readable
- Handles malformed `dataBytes` maps gracefully
- Skips events with invalid byte values
- Provides clear error messages for debugging

## Requirements

- Go 1.18+ (for the dYdX v4-chain module dependencies)
- Access to the dYdX v4-chain repository (for proto definitions)

## Troubleshooting

- **"failed to read event file"**: Ensure the JSON file path is correct and the file exists
- **"failed to unmarshal json"**: Check that the JSON structure matches the expected format
- **"failed to unmarshal proto message"**: The `dataBytes` may be corrupted or incomplete
- **"unknown event subtype"**: The event type may not be supported yet (add it to the switch statement)
