package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Config holds the service configuration
type Config struct {
	Port                string
	RPCURL              string
	PrivateKey          string
	PerpEngineAddr      common.Address
	CollateralVaultAddr common.Address
	ChainID             *big.Int
	ManageNonces        bool
	StrictOnChainErrors bool
	GasCheck            bool
}

// Service holds the dependencies
type Service struct {
	config     Config
	client     *ethclient.Client
	engine     *PerpEngine
	vault      *CollateralVault
	privateKey *ecdsa.PrivateKey
	fromAddr   common.Address
	nonceMu    sync.Mutex
}

// SettleRequest is the JSON payload for POST /settle
type SettleRequest struct {
	MarketID     string `json:"marketId"`
	UserID       string `json:"userId"`
	EVMAddress   string `json:"evmAddress"`
	BalanceDelta string `json:"balanceDelta"`
	SizeDelta    string `json:"sizeDelta"`
	Reason       string `json:"reason"`
	Reference    string `json:"reference"`
}

type SettleBatchRequest struct {
	Settlements []SettleRequest `json:"settlements"`
}

// SettleResponse is the JSON response
type SettleResponse struct {
	TxHash string `json:"txHash"`
	Error  string `json:"error,omitempty"`
}

// UserStateResponse is the JSON response for GET /user-state
type UserStateResponse struct {
	Balance    string `json:"balance"`    // USDC Balance
	Position   string `json:"position"`   // Position Size
	EntryPrice string `json:"entryPrice"` // Entry Price
	MarketID   string `json:"marketId"`
}

func main() {
	// 1. Load Config
	cfg := loadConfig()

	// 2. Init Ethereum Client
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		log.Fatalf("Failed to connect to RPC: %v", err)
	}

	// 3. Load Private Key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 4. Init Contract
	engine, err := NewPerpEngine(cfg.PerpEngineAddr, client)
	if err != nil {
		log.Fatalf("Failed to bind to PerpEngine: %v", err)
	}

	vault, err := NewCollateralVault(cfg.CollateralVaultAddr, client)
	if err != nil {
		log.Fatalf("Failed to bind to CollateralVault: %v", err)
	}

	svc := &Service{
		config:     cfg,
		client:     client,
		engine:     engine,
		vault:      vault,
		privateKey: privateKey,
		fromAddr:   fromAddress,
	}

	// 5. Start Server
	http.HandleFunc("/settle", svc.handleSettle)
	http.HandleFunc("/settle-batch", svc.handleSettleBatch)
	http.HandleFunc("/user-state", svc.handleUserState)
	log.Printf("Gateway starting on port %s...", cfg.Port)
	log.Printf("Operator: %s", fromAddress.Hex())
	log.Printf("Contract: %s", cfg.PerpEngineAddr.Hex())

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}

func (s *Service) handleSettle(w http.ResponseWriter, r *http.Request) {
	// CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SettleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Basic Validation
	if !common.IsHexAddress(req.EVMAddress) {
		http.Error(w, "Invalid EVM address", http.StatusBadRequest)
		return
	}
	// Default to BTC-USD if missing (for backward compatibility or ease)
	marketIDStr := req.MarketID
	if marketIDStr == "" {
		marketIDStr = "BTC-USD"
	}

	// Convert MarketID string to bytes32 (right-padded with zeros, like ethers.encodeBytes32String)
	var marketID [32]byte
	copy(marketID[:], []byte(marketIDStr))

	balanceDelta, ok := new(big.Int).SetString(req.BalanceDelta, 10)
	if !ok {
		http.Error(w, "Invalid balanceDelta", http.StatusBadRequest)
		return
	}
	sizeDelta, ok := new(big.Int).SetString(req.SizeDelta, 10)
	if !ok {
		http.Error(w, "Invalid sizeDelta", http.StatusBadRequest)
		return
	}

	log.Printf("Settling: Market=%s User=%s Addr=%s Bal=%s Size=%s",
		marketIDStr, req.UserID, req.EVMAddress, balanceDelta, sizeDelta)

	// Send Transaction
	txHash, blockNum, err := s.sendSettleTx(context.Background(), marketID, common.HexToAddress(req.EVMAddress), balanceDelta, sizeDelta)

	resp := SettleResponse{
		TxHash: txHash,
	}

	if err != nil {
		log.Printf("Error sending tx: %v", err)
		if s.config.StrictOnChainErrors {
			http.Error(w, fmt.Sprintf("Transaction failed: %v", err), http.StatusInternalServerError)
			return
		}
		resp.Error = err.Error()
	} else {
		log.Printf("Success: Tx=%s Block=%d", txHash, blockNum)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Service) handleSettleBatch(w http.ResponseWriter, r *http.Request) {
	// CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SettleBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Settling Batch: %d items", len(req.Settlements))

	txHash, blockNum, err := s.sendSettleBatchTx(req.Settlements)

	resp := SettleResponse{
		TxHash: txHash,
	}

	if err != nil {
		log.Printf("Error sending batch tx: %v", err)
		if s.config.StrictOnChainErrors {
			http.Error(w, fmt.Sprintf("Transaction failed: %v", err), http.StatusInternalServerError)
			return
		}
		resp.Error = err.Error()
	} else {
		log.Printf("Success Batch: Tx=%s Block=%d", txHash, blockNum)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Service) handleUserState(w http.ResponseWriter, r *http.Request) {
	// CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	addrStr := r.URL.Query().Get("address")
	if !common.IsHexAddress(addrStr) {
		http.Error(w, "Invalid address", http.StatusBadRequest)
		return
	}
	userAddr := common.HexToAddress(addrStr)

	marketIDStr := r.URL.Query().Get("marketId")
	if marketIDStr == "" {
		marketIDStr = "BTC-USD"
	}
	var marketID [32]byte
	copy(marketID[:], []byte(marketIDStr))

	// 1. Get Balance from Vault
	bal, err := s.vault.BalanceOf(&bind.CallOpts{}, userAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get balance: %v", err), http.StatusInternalServerError)
		return
	}

	// 2. Get Position from Engine
	pos, err := s.engine.GetPosition(&bind.CallOpts{}, marketID, userAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get position: %v", err), http.StatusInternalServerError)
		return
	}

	resp := UserStateResponse{
		Balance:    bal.String(),
		Position:   pos.Size.String(),
		EntryPrice: pos.EntryPrice.String(),
		MarketID:   marketIDStr,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Service) sendSettleTx(ctx context.Context, marketID [32]byte, user common.Address, balanceDelta, sizeDelta *big.Int) (string, uint64, error) {
	// ... (reuse existing logic helper or copy-paste for safety)
	// For simplicity and to avoid breaking existing flow, I will keep the logic inline or extract a helper.
	// Let's copy the common tx building logic to a helper `buildAndSendTx`.

	return s.buildAndSendTx(ctx, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return s.engine.Settle(opts, marketID, user, balanceDelta, sizeDelta)
	})
}

func (s *Service) sendSettleBatchTx(settlements []SettleRequest) (string, uint64, error) {
	// Convert to Contract Structs
	var batch []PerpEngineSettlement
	for _, req := range settlements {
		var marketID [32]byte
		copy(marketID[:], []byte(req.MarketID))

		balDelta, _ := new(big.Int).SetString(req.BalanceDelta, 10)
		szDelta, _ := new(big.Int).SetString(req.SizeDelta, 10)

		batch = append(batch, PerpEngineSettlement{
			MarketId:     marketID,
			User:         common.HexToAddress(req.EVMAddress),
			BalanceDelta: balDelta,
			SizeDelta:    szDelta,
		})
	}

	return s.buildAndSendTx(context.Background(), func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return s.engine.SettleBatch(opts, batch)
	})
}

// Helper to avoid code duplication
func (s *Service) buildAndSendTx(ctx context.Context, txFunc func(*bind.TransactOpts) (*types.Transaction, error)) (string, uint64, error) {
	// Mutex for nonce management
	if s.config.ManageNonces {
		s.nonceMu.Lock()
		defer s.nonceMu.Unlock()
	}

	// Gas Check
	if s.config.GasCheck {
		bal, err := s.client.BalanceAt(ctx, s.fromAddr, nil)
		if err == nil && bal.Cmp(big.NewInt(1000000000000000)) < 0 { // 0.001 ETH
			return "", 0, fmt.Errorf("insufficient operator ETH balance: %s", bal)
		}
	}

	// Get Nonce
	nonce, err := s.client.PendingNonceAt(ctx, s.fromAddr)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get nonce: %v", err)
	}

	// Get Gas Price
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get gas price: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(s.privateKey, s.config.ChainID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(500000) // Higher gas limit for batch
	auth.GasPrice = gasPrice

	// Call Contract
	tx, err := txFunc(auth)
	if err != nil {
		return "", 0, fmt.Errorf("contract call failed: %v", err)
	}

	// Wait for Confirmation (Simple polling for demo)
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return tx.Hash().Hex(), 0, fmt.Errorf("tx mining failed: %v", err)
	}

	if receipt.Status == 0 {
		return tx.Hash().Hex(), receipt.BlockNumber.Uint64(), fmt.Errorf("tx reverted on chain")
	}

	return tx.Hash().Hex(), receipt.BlockNumber.Uint64(), nil
}

func loadConfig() Config {
	// 1. Try to load defaults from parent project if env vars are missing
	autoKey, autoAddr, autoVaultAddr := tryAutoDiscover()

	rpcURL := getEnv("RPC_URL", "https://sepolia.base.org")
	
	// Deterministic ChainID mapping (self-contained)
	chainIDMap := map[string]int64{
		"https://sepolia.base.org":                   84532,  // Base Sepolia
		"https://sepolia-rollup.arbitrum.io/rpc":     421614, // Arbitrum Sepolia
		"http://localhost:8545":                      31337,  // Anvil (local)
	}
	
	// Get ChainID: env var override > mapping > default
	var chainID int64
	if envChainID := getEnv("CHAIN_ID", ""); envChainID != "" {
		// Parse env var
		parsed, err := strconv.ParseInt(envChainID, 10, 64)
		if err == nil {
			chainID = parsed
		} else {
			log.Printf("Warning: Invalid CHAIN_ID '%s', using mapping", envChainID)
			if mappedID, ok := chainIDMap[rpcURL]; ok {
				chainID = mappedID
			} else {
				chainID = 84532 // Default: Base Sepolia
			}
		}
	} else if mappedID, ok := chainIDMap[rpcURL]; ok {
		chainID = mappedID
	} else {
		chainID = 84532 // Default: Base Sepolia
	}

	// Check OPERATOR_PRIVATE_KEY from Gateway's .env, fallback to auto-discovered
	privateKeyEnv := getEnv("OPERATOR_PRIVATE_KEY", autoKey)

	return Config{
		Port:                getEnv("GATEWAY_PORT", "8080"),
		RPCURL:              rpcURL,
		PrivateKey:          privateKeyEnv,
		PerpEngineAddr:      common.HexToAddress(getEnv("PERPENGINE_ADDRESS", autoAddr)),
		CollateralVaultAddr: common.HexToAddress(getEnv("COLLATERAL_VAULT_ADDRESS", autoVaultAddr)),
		ChainID:             big.NewInt(chainID),
		ManageNonces:        getEnv("GATEWAY_MANAGE_NONCES", "true") == "true",
		StrictOnChainErrors: getEnv("GATEWAY_STRICT_ONCHAIN_ERRORS", "true") == "true",
		GasCheck:            getEnv("GATEWAY_GAS_CHECK", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

// tryAutoDiscover attempts to read config from the sibling main_folder
func tryAutoDiscover() (string, string, string) {
	var key, engineAddr, vaultAddr string

	// 1. Try to read Operator Private Key from ../v5-demo-contracts/.env
	if content, err := os.ReadFile("../v5-demo-contracts/.env"); err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "OPERATOR_PRIVATE_KEY=") {
				key = strings.TrimPrefix(line, "OPERATOR_PRIVATE_KEY=")
				key = strings.TrimSpace(key)
				key = strings.Trim(key, `"'`) // Remove quotes if present
				log.Println("Auto-discovered Operator Private Key from ../v5-demo-contracts/.env")
				break
			}
		}
	}

	// 2. Try to read Engine Address from ../v5-demo-contracts/deployment.json
	// deployment.json structure: { "84532": { "addresses": { "engine": "...", "vault": "..." } } }
	if content, err := os.ReadFile("../v5-demo-contracts/deployment.json"); err == nil {
		var deploymentConfig map[string]struct {
			Addresses struct {
				Engine string `json:"engine"`
				Vault  string `json:"vault"`
			} `json:"addresses"`
		}
		if err := json.Unmarshal(content, &deploymentConfig); err == nil {
			// Try to find config for current chain (84532 = Base Sepolia, 421614 = Arbitrum Sepolia)
			// If not found, use first available chain config
			var chainConfig struct {
				Addresses struct {
					Engine string `json:"engine"`
					Vault  string `json:"vault"`
				} `json:"addresses"`
			}
			found := false
			for chainID, config := range deploymentConfig {
				if chainID == "84532" || chainID == "421614" || !found {
					chainConfig = config
					found = true
					if chainID == "84532" || chainID == "421614" {
						break
					}
				}
			}
			if found {
				if chainConfig.Addresses.Engine != "" {
					engineAddr = chainConfig.Addresses.Engine
					log.Printf("Auto-discovered Engine Address: %s", engineAddr)
				}
				if chainConfig.Addresses.Vault != "" {
					vaultAddr = chainConfig.Addresses.Vault
					log.Printf("Auto-discovered Vault Address: %s", vaultAddr)
				}
			}
		}
	}

	return key, engineAddr, vaultAddr
}
