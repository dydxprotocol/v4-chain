package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"

	"github.com/cosmos/gogoproto/proto"
	events "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
)

// Structs to match the JSON structure
type Event struct {
	DataBytes        map[string]json.Number `json:"dataBytes"`
	Subtype          string                 `json:"subtype"`
	EventIndex       int                    `json:"eventIndex"`
	TransactionIndex int                    `json:"transactionIndex"`
	Version          int                    `json:"version"`
}

type IndexerTendermintBlock struct {
	Events []Event `json:"events"`
}

type BlockData struct {
	IndexerTendermintBlock IndexerTendermintBlock `json:"block"`
}

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Error: JSON file path is required")
		fmt.Println("Usage: ./deserialize_proto <json_file_path>")
		os.Exit(1)
	}
	jsonFilePath := os.Args[1]

	// Open the JSON file for reading
	file, err := os.ReadFile(jsonFilePath)
	if err != nil {
		log.Fatalf("failed to read event file '%s': %v", jsonFilePath, err)
	}

	fmt.Printf("Reading events from: %s\n\n", jsonFilePath)

	var blockData BlockData
	if err := json.Unmarshal(file, &blockData); err != nil {
		log.Fatalf("failed to unmarshal json: %v", err)
	}

	fmt.Printf("Found %d events\n\n", len(blockData.IndexerTendermintBlock.Events))
	for _, event := range blockData.IndexerTendermintBlock.Events {
		// Convert dataBytes map to byte slice
		if len(event.DataBytes) == 0 {
			log.Printf("empty dataBytes for event subtype %s, skipping", event.Subtype)
			continue
		}

		keys := make([]int, 0, len(event.DataBytes))
		for k := range event.DataBytes {
			key, err := strconv.Atoi(k)
			if err != nil {
				log.Printf("failed to convert map key '%s' to int: %v, skipping event", k, err)
				continue
			}
			keys = append(keys, key)
		}
		sort.Ints(keys)

		// Find the maximum key to determine byte slice size
		maxKey := keys[len(keys)-1]
		bytes := make([]byte, maxKey+1)

		for _, k := range keys {
			val, err := event.DataBytes[strconv.Itoa(k)].Int64()
			if err != nil {
				log.Printf("failed to convert json.Number to int64 for key %d: %v, skipping event", k, err)
				continue
			}
			if val < 0 || val > 255 {
				log.Printf("byte value %d out of range [0,255] for key %d, skipping event", val, k)
				continue
			}
			bytes[k] = byte(val)
		}

		fmt.Printf("--- Event Index: %d, Transaction Index: %d, Subtype: %s ---\n",
			event.EventIndex,
			event.TransactionIndex,
			event.Subtype,
		)

		var msg proto.Message
		switch event.Subtype {
		case "subaccount_update":
			msg = &events.SubaccountUpdateEventV1{}
		case "order_fill":
			msg = &events.OrderFillEventV1{}
		case "transfer":
			msg = &events.TransferEventV1{}
		case "deleveraging":
			msg = &events.DeleveragingEventV1{}
		case "stateful_order":
			msg = &events.StatefulOrderEventV1{}
		case "funding":
			msg = &events.FundingEventV1{}
		case "market":
			msg = &events.MarketEventV1{}
		case "asset_create":
			msg = &events.AssetCreateEventV1{}
		case "liquidity_tier":
			msg = &events.LiquidityTierUpsertEventV1{}
		case "update_clob_pair":
			msg = &events.UpdateClobPairEventV1{}
		case "trading_rewards":
			msg = &events.TradingRewardsEventV1{}
		case "register_affiliate":
			msg = &events.RegisterAffiliateEventV1{}
		case "upsert_vault":
			msg = &events.UpsertVaultEventV1{}
		default:
			log.Printf("unknown event subtype: %s", event.Subtype)
			continue
		}

		if err := proto.Unmarshal(bytes, msg); err != nil {
			log.Printf("failed to unmarshal proto message for subtype %s: %v", event.Subtype, err)
			continue
		}

		// Print the decoded event to the console with pretty formatting
		prettyPrintProtoMessage(msg)
		fmt.Println() // Add blank line between events
	}
}

// prettyPrintProtoMessage formats and prints a proto message in a readable way
func prettyPrintProtoMessage(msg proto.Message) {
	if msg == nil {
		fmt.Println("<nil message>")
		return
	}

	// Use reflection to get the message type name
	msgType := reflect.TypeOf(msg)
	if msgType.Kind() == reflect.Ptr {
		msgType = msgType.Elem()
	}
	fmt.Printf("Message Type: %s\n", msgType.Name())

	// Convert to string and format it nicely
	msgStr := msg.String()
	if msgStr == "" {
		fmt.Println("<empty message>")
		return
	}

	// Split by spaces and format fields on separate lines
	fields := parseProtoFields(msgStr)
	for _, field := range fields {
		fmt.Printf("  %s\n", field)
	}
}

// parseProtoFields parses the proto message string and formats it nicely
func parseProtoFields(msgStr string) []string {
	var fields []string
	var current string
	braceLevel := 0
	inQuotes := false

	for _, char := range msgStr {
		switch char {
		case '"':
			inQuotes = !inQuotes
			current += string(char)
		case '<':
			if !inQuotes {
				braceLevel++
			}
			current += string(char)
		case '>':
			if !inQuotes {
				braceLevel--
			}
			current += string(char)
		case ' ':
			if !inQuotes && braceLevel == 0 && current != "" {
				// End of a field
				fields = append(fields, current)
				current = ""
			} else {
				current += string(char)
			}
		default:
			current += string(char)
		}
	}

	// Add the last field if any
	if current != "" {
		fields = append(fields, current)
	}

	return fields
}
