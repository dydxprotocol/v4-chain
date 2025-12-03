// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PerpEnginePosition is an auto generated low-level Go binding around an user-defined struct.
type PerpEnginePosition struct {
	Size       *big.Int
	EntryPrice *big.Int
}

// PerpEngineSettlement is an auto generated low-level Go binding around an user-defined struct.
type PerpEngineSettlement struct {
	MarketId     [32]byte
	User         common.Address
	BalanceDelta *big.Int
	SizeDelta    *big.Int
}

// PerpEngineMetaData contains all meta data concerning the PerpEngine contract.
var PerpEngineMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"balanceDelta\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"sizeDelta\",\"type\":\"int256\"}],\"name\":\"BalanceSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"liquidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"penaltyPaid\",\"type\":\"int256\"}],\"name\":\"Liquidated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"}],\"name\":\"MarketAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"newSize\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"entryPrice\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"realizedPnL\",\"type\":\"int256\"}],\"name\":\"PositionChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"initialMarginRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maintenanceMarginRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tradingFeeRate\",\"type\":\"uint256\"}],\"name\":\"addMarket\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxSlippageBps\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"closePosition\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"engineEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getMargin\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"equity\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"totalNotional\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"marginRatio\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getPosition\",\"outputs\":[{\"components\":[{\"internalType\":\"int256\",\"name\":\"size\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"entryPrice\",\"type\":\"int256\"}],\"internalType\":\"structPerpEngine.Position\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUnrealizedPnl\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"liquidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"markets\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"initialMarginRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maintenanceMarginRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tradingFeeRate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"},{\"internalType\":\"int256\",\"name\":\"sizeDelta\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"maxPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"openPosition\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracle\",\"outputs\":[{\"internalType\":\"contractOracle\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"positions\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"size\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"entryPrice\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_enabled\",\"type\":\"bool\"}],\"name\":\"setEngineEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"balanceDelta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"sizeDelta\",\"type\":\"int256\"}],\"name\":\"settle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"marketId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"balanceDelta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"sizeDelta\",\"type\":\"int256\"}],\"internalType\":\"structPerpEngine.Settlement[]\",\"name\":\"settlements\",\"type\":\"tuple[]\"}],\"name\":\"settleBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"vault\",\"outputs\":[{\"internalType\":\"contractCollateralVault\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// PerpEngineABI is the input ABI used to generate the binding from.
// Deprecated: Use PerpEngineMetaData.ABI instead.
var PerpEngineABI = PerpEngineMetaData.ABI

// PerpEngine is an auto generated Go binding around an Ethereum contract.
type PerpEngine struct {
	PerpEngineCaller     // Read-only binding to the contract
	PerpEngineTransactor // Write-only binding to the contract
	PerpEngineFilterer   // Log filterer for contract events
}

// PerpEngineCaller is an auto generated read-only Go binding around an Ethereum contract.
type PerpEngineCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PerpEngineTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PerpEngineTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PerpEngineFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PerpEngineFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PerpEngineSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PerpEngineSession struct {
	Contract     *PerpEngine       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PerpEngineCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PerpEngineCallerSession struct {
	Contract *PerpEngineCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// PerpEngineTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PerpEngineTransactorSession struct {
	Contract     *PerpEngineTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// PerpEngineRaw is an auto generated low-level Go binding around an Ethereum contract.
type PerpEngineRaw struct {
	Contract *PerpEngine // Generic contract binding to access the raw methods on
}

// PerpEngineCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PerpEngineCallerRaw struct {
	Contract *PerpEngineCaller // Generic read-only contract binding to access the raw methods on
}

// PerpEngineTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PerpEngineTransactorRaw struct {
	Contract *PerpEngineTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPerpEngine creates a new instance of PerpEngine, bound to a specific deployed contract.
func NewPerpEngine(address common.Address, backend bind.ContractBackend) (*PerpEngine, error) {
	contract, err := bindPerpEngine(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PerpEngine{PerpEngineCaller: PerpEngineCaller{contract: contract}, PerpEngineTransactor: PerpEngineTransactor{contract: contract}, PerpEngineFilterer: PerpEngineFilterer{contract: contract}}, nil
}

// NewPerpEngineCaller creates a new read-only instance of PerpEngine, bound to a specific deployed contract.
func NewPerpEngineCaller(address common.Address, caller bind.ContractCaller) (*PerpEngineCaller, error) {
	contract, err := bindPerpEngine(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PerpEngineCaller{contract: contract}, nil
}

// NewPerpEngineTransactor creates a new write-only instance of PerpEngine, bound to a specific deployed contract.
func NewPerpEngineTransactor(address common.Address, transactor bind.ContractTransactor) (*PerpEngineTransactor, error) {
	contract, err := bindPerpEngine(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PerpEngineTransactor{contract: contract}, nil
}

// NewPerpEngineFilterer creates a new log filterer instance of PerpEngine, bound to a specific deployed contract.
func NewPerpEngineFilterer(address common.Address, filterer bind.ContractFilterer) (*PerpEngineFilterer, error) {
	contract, err := bindPerpEngine(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PerpEngineFilterer{contract: contract}, nil
}

// bindPerpEngine binds a generic wrapper to an already deployed contract.
func bindPerpEngine(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PerpEngineMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PerpEngine *PerpEngineRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PerpEngine.Contract.PerpEngineCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PerpEngine *PerpEngineRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerpEngine.Contract.PerpEngineTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PerpEngine *PerpEngineRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PerpEngine.Contract.PerpEngineTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PerpEngine *PerpEngineCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PerpEngine.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PerpEngine *PerpEngineTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerpEngine.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PerpEngine *PerpEngineTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PerpEngine.Contract.contract.Transact(opts, method, params...)
}

// EngineEnabled is a free data retrieval call binding the contract method 0x1c0b1cad.
//
// Solidity: function engineEnabled() view returns(bool)
func (_PerpEngine *PerpEngineCaller) EngineEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PerpEngine.contract.Call(opts, &out, "engineEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// EngineEnabled is a free data retrieval call binding the contract method 0x1c0b1cad.
//
// Solidity: function engineEnabled() view returns(bool)
func (_PerpEngine *PerpEngineSession) EngineEnabled() (bool, error) {
	return _PerpEngine.Contract.EngineEnabled(&_PerpEngine.CallOpts)
}

// EngineEnabled is a free data retrieval call binding the contract method 0x1c0b1cad.
//
// Solidity: function engineEnabled() view returns(bool)
func (_PerpEngine *PerpEngineCallerSession) EngineEnabled() (bool, error) {
	return _PerpEngine.Contract.EngineEnabled(&_PerpEngine.CallOpts)
}

// GetMargin is a free data retrieval call binding the contract method 0xf84f89a2.
//
// Solidity: function getMargin(address user) view returns(int256 equity, uint256 totalNotional, uint256 marginRatio)
func (_PerpEngine *PerpEngineCaller) GetMargin(opts *bind.CallOpts, user common.Address) (struct {
	Equity        *big.Int
	TotalNotional *big.Int
	MarginRatio   *big.Int
}, error) {
	var out []interface{}
	err := _PerpEngine.contract.Call(opts, &out, "getMargin", user)

	outstruct := new(struct {
		Equity        *big.Int
		TotalNotional *big.Int
		MarginRatio   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Equity = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalNotional = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.MarginRatio = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetMargin is a free data retrieval call binding the contract method 0xf84f89a2.
//
// Solidity: function getMargin(address user) view returns(int256 equity, uint256 totalNotional, uint256 marginRatio)
func (_PerpEngine *PerpEngineSession) GetMargin(user common.Address) (struct {
	Equity        *big.Int
	TotalNotional *big.Int
	MarginRatio   *big.Int
}, error) {
	return _PerpEngine.Contract.GetMargin(&_PerpEngine.CallOpts, user)
}

// GetMargin is a free data retrieval call binding the contract method 0xf84f89a2.
//
// Solidity: function getMargin(address user) view returns(int256 equity, uint256 totalNotional, uint256 marginRatio)
func (_PerpEngine *PerpEngineCallerSession) GetMargin(user common.Address) (struct {
	Equity        *big.Int
	TotalNotional *big.Int
	MarginRatio   *big.Int
}, error) {
	return _PerpEngine.Contract.GetMargin(&_PerpEngine.CallOpts, user)
}

// GetPosition is a free data retrieval call binding the contract method 0x5c388821.
//
// Solidity: function getPosition(bytes32 marketId, address user) view returns((int256,int256))
func (_PerpEngine *PerpEngineCaller) GetPosition(opts *bind.CallOpts, marketId [32]byte, user common.Address) (PerpEnginePosition, error) {
	var out []interface{}
	err := _PerpEngine.contract.Call(opts, &out, "getPosition", marketId, user)

	if err != nil {
		return *new(PerpEnginePosition), err
	}

	out0 := *abi.ConvertType(out[0], new(PerpEnginePosition)).(*PerpEnginePosition)

	return out0, err

}

// GetPosition is a free data retrieval call binding the contract method 0x5c388821.
//
// Solidity: function getPosition(bytes32 marketId, address user) view returns((int256,int256))
func (_PerpEngine *PerpEngineSession) GetPosition(marketId [32]byte, user common.Address) (PerpEnginePosition, error) {
	return _PerpEngine.Contract.GetPosition(&_PerpEngine.CallOpts, marketId, user)
}

// GetPosition is a free data retrieval call binding the contract method 0x5c388821.
//
// Solidity: function getPosition(bytes32 marketId, address user) view returns((int256,int256))
func (_PerpEngine *PerpEngineCallerSession) GetPosition(marketId [32]byte, user common.Address) (PerpEnginePosition, error) {
	return _PerpEngine.Contract.GetPosition(&_PerpEngine.CallOpts, marketId, user)
}

// GetUnrealizedPnl is a free data retrieval call binding the contract method 0x43b259f2.
//
// Solidity: function getUnrealizedPnl(bytes32 marketId, address user) view returns(int256)
func (_PerpEngine *PerpEngineCaller) GetUnrealizedPnl(opts *bind.CallOpts, marketId [32]byte, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PerpEngine.contract.Call(opts, &out, "getUnrealizedPnl", marketId, user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUnrealizedPnl is a free data retrieval call binding the contract method 0x43b259f2.
//
// Solidity: function getUnrealizedPnl(bytes32 marketId, address user) view returns(int256)
func (_PerpEngine *PerpEngineSession) GetUnrealizedPnl(marketId [32]byte, user common.Address) (*big.Int, error) {
	return _PerpEngine.Contract.GetUnrealizedPnl(&_PerpEngine.CallOpts, marketId, user)
}

// GetUnrealizedPnl is a free data retrieval call binding the contract method 0x43b259f2.
//
// Solidity: function getUnrealizedPnl(bytes32 marketId, address user) view returns(int256)
func (_PerpEngine *PerpEngineCallerSession) GetUnrealizedPnl(marketId [32]byte, user common.Address) (*big.Int, error) {
	return _PerpEngine.Contract.GetUnrealizedPnl(&_PerpEngine.CallOpts, marketId, user)
}

// Markets is a free data retrieval call binding the contract method 0x7564912b.
//
// Solidity: function markets(bytes32 ) view returns(uint256 initialMarginRatio, uint256 maintenanceMarginRatio, uint256 tradingFeeRate)
func (_PerpEngine *PerpEngineCaller) Markets(opts *bind.CallOpts, arg0 [32]byte) (struct {
	InitialMarginRatio     *big.Int
	MaintenanceMarginRatio *big.Int
	TradingFeeRate         *big.Int
}, error) {
	var out []interface{}
	err := _PerpEngine.contract.Call(opts, &out, "markets", arg0)

	outstruct := new(struct {
		InitialMarginRatio     *big.Int
		MaintenanceMarginRatio *big.Int
		TradingFeeRate         *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.InitialMarginRatio = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.MaintenanceMarginRatio = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TradingFeeRate = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Markets is a free data retrieval call binding the contract method 0x7564912b.
//
// Solidity: function markets(bytes32 ) view returns(uint256 initialMarginRatio, uint256 maintenanceMarginRatio, uint256 tradingFeeRate)
func (_PerpEngine *PerpEngineSession) Markets(arg0 [32]byte) (struct {
	InitialMarginRatio     *big.Int
	MaintenanceMarginRatio *big.Int
	TradingFeeRate         *big.Int
}, error) {
	return _PerpEngine.Contract.Markets(&_PerpEngine.CallOpts, arg0)
}

// Markets is a free data retrieval call binding the contract method 0x7564912b.
//
// Solidity: function markets(bytes32 ) view returns(uint256 initialMarginRatio, uint256 maintenanceMarginRatio, uint256 tradingFeeRate)
func (_PerpEngine *PerpEngineCallerSession) Markets(arg0 [32]byte) (struct {
	InitialMarginRatio     *big.Int
	MaintenanceMarginRatio *big.Int
	TradingFeeRate         *big.Int
}, error) {
	return _PerpEngine.Contract.Markets(&_PerpEngine.CallOpts, arg0)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_PerpEngine *PerpEngineCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PerpEngine.contract.Call(opts, &out, "oracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_PerpEngine *PerpEngineSession) Oracle() (common.Address, error) {
	return _PerpEngine.Contract.Oracle(&_PerpEngine.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_PerpEngine *PerpEngineCallerSession) Oracle() (common.Address, error) {
	return _PerpEngine.Contract.Oracle(&_PerpEngine.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PerpEngine *PerpEngineCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PerpEngine.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PerpEngine *PerpEngineSession) Owner() (common.Address, error) {
	return _PerpEngine.Contract.Owner(&_PerpEngine.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PerpEngine *PerpEngineCallerSession) Owner() (common.Address, error) {
	return _PerpEngine.Contract.Owner(&_PerpEngine.CallOpts)
}

// Positions is a free data retrieval call binding the contract method 0x29d88594.
//
// Solidity: function positions(bytes32 , address ) view returns(int256 size, int256 entryPrice)
func (_PerpEngine *PerpEngineCaller) Positions(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (struct {
	Size       *big.Int
	EntryPrice *big.Int
}, error) {
	var out []interface{}
	err := _PerpEngine.contract.Call(opts, &out, "positions", arg0, arg1)

	outstruct := new(struct {
		Size       *big.Int
		EntryPrice *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Size = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.EntryPrice = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Positions is a free data retrieval call binding the contract method 0x29d88594.
//
// Solidity: function positions(bytes32 , address ) view returns(int256 size, int256 entryPrice)
func (_PerpEngine *PerpEngineSession) Positions(arg0 [32]byte, arg1 common.Address) (struct {
	Size       *big.Int
	EntryPrice *big.Int
}, error) {
	return _PerpEngine.Contract.Positions(&_PerpEngine.CallOpts, arg0, arg1)
}

// Positions is a free data retrieval call binding the contract method 0x29d88594.
//
// Solidity: function positions(bytes32 , address ) view returns(int256 size, int256 entryPrice)
func (_PerpEngine *PerpEngineCallerSession) Positions(arg0 [32]byte, arg1 common.Address) (struct {
	Size       *big.Int
	EntryPrice *big.Int
}, error) {
	return _PerpEngine.Contract.Positions(&_PerpEngine.CallOpts, arg0, arg1)
}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_PerpEngine *PerpEngineCaller) Vault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PerpEngine.contract.Call(opts, &out, "vault")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_PerpEngine *PerpEngineSession) Vault() (common.Address, error) {
	return _PerpEngine.Contract.Vault(&_PerpEngine.CallOpts)
}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_PerpEngine *PerpEngineCallerSession) Vault() (common.Address, error) {
	return _PerpEngine.Contract.Vault(&_PerpEngine.CallOpts)
}

// AddMarket is a paid mutator transaction binding the contract method 0x984e35c2.
//
// Solidity: function addMarket(bytes32 marketId, uint256 initialMarginRatio, uint256 maintenanceMarginRatio, uint256 tradingFeeRate) returns()
func (_PerpEngine *PerpEngineTransactor) AddMarket(opts *bind.TransactOpts, marketId [32]byte, initialMarginRatio *big.Int, maintenanceMarginRatio *big.Int, tradingFeeRate *big.Int) (*types.Transaction, error) {
	return _PerpEngine.contract.Transact(opts, "addMarket", marketId, initialMarginRatio, maintenanceMarginRatio, tradingFeeRate)
}

// AddMarket is a paid mutator transaction binding the contract method 0x984e35c2.
//
// Solidity: function addMarket(bytes32 marketId, uint256 initialMarginRatio, uint256 maintenanceMarginRatio, uint256 tradingFeeRate) returns()
func (_PerpEngine *PerpEngineSession) AddMarket(marketId [32]byte, initialMarginRatio *big.Int, maintenanceMarginRatio *big.Int, tradingFeeRate *big.Int) (*types.Transaction, error) {
	return _PerpEngine.Contract.AddMarket(&_PerpEngine.TransactOpts, marketId, initialMarginRatio, maintenanceMarginRatio, tradingFeeRate)
}

// AddMarket is a paid mutator transaction binding the contract method 0x984e35c2.
//
// Solidity: function addMarket(bytes32 marketId, uint256 initialMarginRatio, uint256 maintenanceMarginRatio, uint256 tradingFeeRate) returns()
func (_PerpEngine *PerpEngineTransactorSession) AddMarket(marketId [32]byte, initialMarginRatio *big.Int, maintenanceMarginRatio *big.Int, tradingFeeRate *big.Int) (*types.Transaction, error) {
	return _PerpEngine.Contract.AddMarket(&_PerpEngine.TransactOpts, marketId, initialMarginRatio, maintenanceMarginRatio, tradingFeeRate)
}

// ClosePosition is a paid mutator transaction binding the contract method 0x0afd162f.
//
// Solidity: function closePosition(bytes32 marketId, uint256 maxSlippageBps, uint256 deadline) returns()
func (_PerpEngine *PerpEngineTransactor) ClosePosition(opts *bind.TransactOpts, marketId [32]byte, maxSlippageBps *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _PerpEngine.contract.Transact(opts, "closePosition", marketId, maxSlippageBps, deadline)
}

// ClosePosition is a paid mutator transaction binding the contract method 0x0afd162f.
//
// Solidity: function closePosition(bytes32 marketId, uint256 maxSlippageBps, uint256 deadline) returns()
func (_PerpEngine *PerpEngineSession) ClosePosition(marketId [32]byte, maxSlippageBps *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _PerpEngine.Contract.ClosePosition(&_PerpEngine.TransactOpts, marketId, maxSlippageBps, deadline)
}

// ClosePosition is a paid mutator transaction binding the contract method 0x0afd162f.
//
// Solidity: function closePosition(bytes32 marketId, uint256 maxSlippageBps, uint256 deadline) returns()
func (_PerpEngine *PerpEngineTransactorSession) ClosePosition(marketId [32]byte, maxSlippageBps *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _PerpEngine.Contract.ClosePosition(&_PerpEngine.TransactOpts, marketId, maxSlippageBps, deadline)
}

// Liquidate is a paid mutator transaction binding the contract method 0x771b51d2.
//
// Solidity: function liquidate(bytes32 marketId, address user) returns()
func (_PerpEngine *PerpEngineTransactor) Liquidate(opts *bind.TransactOpts, marketId [32]byte, user common.Address) (*types.Transaction, error) {
	return _PerpEngine.contract.Transact(opts, "liquidate", marketId, user)
}

// Liquidate is a paid mutator transaction binding the contract method 0x771b51d2.
//
// Solidity: function liquidate(bytes32 marketId, address user) returns()
func (_PerpEngine *PerpEngineSession) Liquidate(marketId [32]byte, user common.Address) (*types.Transaction, error) {
	return _PerpEngine.Contract.Liquidate(&_PerpEngine.TransactOpts, marketId, user)
}

// Liquidate is a paid mutator transaction binding the contract method 0x771b51d2.
//
// Solidity: function liquidate(bytes32 marketId, address user) returns()
func (_PerpEngine *PerpEngineTransactorSession) Liquidate(marketId [32]byte, user common.Address) (*types.Transaction, error) {
	return _PerpEngine.Contract.Liquidate(&_PerpEngine.TransactOpts, marketId, user)
}

// OpenPosition is a paid mutator transaction binding the contract method 0x0b3d0eac.
//
// Solidity: function openPosition(bytes32 marketId, int256 sizeDelta, uint256 maxPrice, uint256 deadline) returns()
func (_PerpEngine *PerpEngineTransactor) OpenPosition(opts *bind.TransactOpts, marketId [32]byte, sizeDelta *big.Int, maxPrice *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _PerpEngine.contract.Transact(opts, "openPosition", marketId, sizeDelta, maxPrice, deadline)
}

// OpenPosition is a paid mutator transaction binding the contract method 0x0b3d0eac.
//
// Solidity: function openPosition(bytes32 marketId, int256 sizeDelta, uint256 maxPrice, uint256 deadline) returns()
func (_PerpEngine *PerpEngineSession) OpenPosition(marketId [32]byte, sizeDelta *big.Int, maxPrice *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _PerpEngine.Contract.OpenPosition(&_PerpEngine.TransactOpts, marketId, sizeDelta, maxPrice, deadline)
}

// OpenPosition is a paid mutator transaction binding the contract method 0x0b3d0eac.
//
// Solidity: function openPosition(bytes32 marketId, int256 sizeDelta, uint256 maxPrice, uint256 deadline) returns()
func (_PerpEngine *PerpEngineTransactorSession) OpenPosition(marketId [32]byte, sizeDelta *big.Int, maxPrice *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _PerpEngine.Contract.OpenPosition(&_PerpEngine.TransactOpts, marketId, sizeDelta, maxPrice, deadline)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PerpEngine *PerpEngineTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PerpEngine.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PerpEngine *PerpEngineSession) RenounceOwnership() (*types.Transaction, error) {
	return _PerpEngine.Contract.RenounceOwnership(&_PerpEngine.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PerpEngine *PerpEngineTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PerpEngine.Contract.RenounceOwnership(&_PerpEngine.TransactOpts)
}

// SetEngineEnabled is a paid mutator transaction binding the contract method 0x946d52bd.
//
// Solidity: function setEngineEnabled(bool _enabled) returns()
func (_PerpEngine *PerpEngineTransactor) SetEngineEnabled(opts *bind.TransactOpts, _enabled bool) (*types.Transaction, error) {
	return _PerpEngine.contract.Transact(opts, "setEngineEnabled", _enabled)
}

// SetEngineEnabled is a paid mutator transaction binding the contract method 0x946d52bd.
//
// Solidity: function setEngineEnabled(bool _enabled) returns()
func (_PerpEngine *PerpEngineSession) SetEngineEnabled(_enabled bool) (*types.Transaction, error) {
	return _PerpEngine.Contract.SetEngineEnabled(&_PerpEngine.TransactOpts, _enabled)
}

// SetEngineEnabled is a paid mutator transaction binding the contract method 0x946d52bd.
//
// Solidity: function setEngineEnabled(bool _enabled) returns()
func (_PerpEngine *PerpEngineTransactorSession) SetEngineEnabled(_enabled bool) (*types.Transaction, error) {
	return _PerpEngine.Contract.SetEngineEnabled(&_PerpEngine.TransactOpts, _enabled)
}

// Settle is a paid mutator transaction binding the contract method 0x169cbd67.
//
// Solidity: function settle(bytes32 marketId, address user, int256 balanceDelta, int256 sizeDelta) returns()
func (_PerpEngine *PerpEngineTransactor) Settle(opts *bind.TransactOpts, marketId [32]byte, user common.Address, balanceDelta *big.Int, sizeDelta *big.Int) (*types.Transaction, error) {
	return _PerpEngine.contract.Transact(opts, "settle", marketId, user, balanceDelta, sizeDelta)
}

// Settle is a paid mutator transaction binding the contract method 0x169cbd67.
//
// Solidity: function settle(bytes32 marketId, address user, int256 balanceDelta, int256 sizeDelta) returns()
func (_PerpEngine *PerpEngineSession) Settle(marketId [32]byte, user common.Address, balanceDelta *big.Int, sizeDelta *big.Int) (*types.Transaction, error) {
	return _PerpEngine.Contract.Settle(&_PerpEngine.TransactOpts, marketId, user, balanceDelta, sizeDelta)
}

// Settle is a paid mutator transaction binding the contract method 0x169cbd67.
//
// Solidity: function settle(bytes32 marketId, address user, int256 balanceDelta, int256 sizeDelta) returns()
func (_PerpEngine *PerpEngineTransactorSession) Settle(marketId [32]byte, user common.Address, balanceDelta *big.Int, sizeDelta *big.Int) (*types.Transaction, error) {
	return _PerpEngine.Contract.Settle(&_PerpEngine.TransactOpts, marketId, user, balanceDelta, sizeDelta)
}

// SettleBatch is a paid mutator transaction binding the contract method 0xe559f0ff.
//
// Solidity: function settleBatch((bytes32,address,int256,int256)[] settlements) returns()
func (_PerpEngine *PerpEngineTransactor) SettleBatch(opts *bind.TransactOpts, settlements []PerpEngineSettlement) (*types.Transaction, error) {
	return _PerpEngine.contract.Transact(opts, "settleBatch", settlements)
}

// SettleBatch is a paid mutator transaction binding the contract method 0xe559f0ff.
//
// Solidity: function settleBatch((bytes32,address,int256,int256)[] settlements) returns()
func (_PerpEngine *PerpEngineSession) SettleBatch(settlements []PerpEngineSettlement) (*types.Transaction, error) {
	return _PerpEngine.Contract.SettleBatch(&_PerpEngine.TransactOpts, settlements)
}

// SettleBatch is a paid mutator transaction binding the contract method 0xe559f0ff.
//
// Solidity: function settleBatch((bytes32,address,int256,int256)[] settlements) returns()
func (_PerpEngine *PerpEngineTransactorSession) SettleBatch(settlements []PerpEngineSettlement) (*types.Transaction, error) {
	return _PerpEngine.Contract.SettleBatch(&_PerpEngine.TransactOpts, settlements)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PerpEngine *PerpEngineTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PerpEngine.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PerpEngine *PerpEngineSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PerpEngine.Contract.TransferOwnership(&_PerpEngine.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PerpEngine *PerpEngineTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PerpEngine.Contract.TransferOwnership(&_PerpEngine.TransactOpts, newOwner)
}

// PerpEngineBalanceSettledIterator is returned from FilterBalanceSettled and is used to iterate over the raw logs and unpacked data for BalanceSettled events raised by the PerpEngine contract.
type PerpEngineBalanceSettledIterator struct {
	Event *PerpEngineBalanceSettled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PerpEngineBalanceSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpEngineBalanceSettled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PerpEngineBalanceSettled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PerpEngineBalanceSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpEngineBalanceSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpEngineBalanceSettled represents a BalanceSettled event raised by the PerpEngine contract.
type PerpEngineBalanceSettled struct {
	MarketId     [32]byte
	User         common.Address
	BalanceDelta *big.Int
	SizeDelta    *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterBalanceSettled is a free log retrieval operation binding the contract event 0x6c88a5accba870d6a157371ca05bb7615afc76776526d1ab955dd3c7a26df6e4.
//
// Solidity: event BalanceSettled(bytes32 indexed marketId, address indexed user, int256 balanceDelta, int256 sizeDelta)
func (_PerpEngine *PerpEngineFilterer) FilterBalanceSettled(opts *bind.FilterOpts, marketId [][32]byte, user []common.Address) (*PerpEngineBalanceSettledIterator, error) {

	var marketIdRule []interface{}
	for _, marketIdItem := range marketId {
		marketIdRule = append(marketIdRule, marketIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PerpEngine.contract.FilterLogs(opts, "BalanceSettled", marketIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &PerpEngineBalanceSettledIterator{contract: _PerpEngine.contract, event: "BalanceSettled", logs: logs, sub: sub}, nil
}

// WatchBalanceSettled is a free log subscription operation binding the contract event 0x6c88a5accba870d6a157371ca05bb7615afc76776526d1ab955dd3c7a26df6e4.
//
// Solidity: event BalanceSettled(bytes32 indexed marketId, address indexed user, int256 balanceDelta, int256 sizeDelta)
func (_PerpEngine *PerpEngineFilterer) WatchBalanceSettled(opts *bind.WatchOpts, sink chan<- *PerpEngineBalanceSettled, marketId [][32]byte, user []common.Address) (event.Subscription, error) {

	var marketIdRule []interface{}
	for _, marketIdItem := range marketId {
		marketIdRule = append(marketIdRule, marketIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PerpEngine.contract.WatchLogs(opts, "BalanceSettled", marketIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpEngineBalanceSettled)
				if err := _PerpEngine.contract.UnpackLog(event, "BalanceSettled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBalanceSettled is a log parse operation binding the contract event 0x6c88a5accba870d6a157371ca05bb7615afc76776526d1ab955dd3c7a26df6e4.
//
// Solidity: event BalanceSettled(bytes32 indexed marketId, address indexed user, int256 balanceDelta, int256 sizeDelta)
func (_PerpEngine *PerpEngineFilterer) ParseBalanceSettled(log types.Log) (*PerpEngineBalanceSettled, error) {
	event := new(PerpEngineBalanceSettled)
	if err := _PerpEngine.contract.UnpackLog(event, "BalanceSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpEngineLiquidatedIterator is returned from FilterLiquidated and is used to iterate over the raw logs and unpacked data for Liquidated events raised by the PerpEngine contract.
type PerpEngineLiquidatedIterator struct {
	Event *PerpEngineLiquidated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PerpEngineLiquidatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpEngineLiquidated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PerpEngineLiquidated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PerpEngineLiquidatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpEngineLiquidatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpEngineLiquidated represents a Liquidated event raised by the PerpEngine contract.
type PerpEngineLiquidated struct {
	MarketId    [32]byte
	User        common.Address
	Liquidator  common.Address
	PenaltyPaid *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLiquidated is a free log retrieval operation binding the contract event 0x89109648b0170073868f7ba7d8a8db89e21725ea8125f7bff949a6fa4ed0915f.
//
// Solidity: event Liquidated(bytes32 indexed marketId, address indexed user, address indexed liquidator, int256 penaltyPaid)
func (_PerpEngine *PerpEngineFilterer) FilterLiquidated(opts *bind.FilterOpts, marketId [][32]byte, user []common.Address, liquidator []common.Address) (*PerpEngineLiquidatedIterator, error) {

	var marketIdRule []interface{}
	for _, marketIdItem := range marketId {
		marketIdRule = append(marketIdRule, marketIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var liquidatorRule []interface{}
	for _, liquidatorItem := range liquidator {
		liquidatorRule = append(liquidatorRule, liquidatorItem)
	}

	logs, sub, err := _PerpEngine.contract.FilterLogs(opts, "Liquidated", marketIdRule, userRule, liquidatorRule)
	if err != nil {
		return nil, err
	}
	return &PerpEngineLiquidatedIterator{contract: _PerpEngine.contract, event: "Liquidated", logs: logs, sub: sub}, nil
}

// WatchLiquidated is a free log subscription operation binding the contract event 0x89109648b0170073868f7ba7d8a8db89e21725ea8125f7bff949a6fa4ed0915f.
//
// Solidity: event Liquidated(bytes32 indexed marketId, address indexed user, address indexed liquidator, int256 penaltyPaid)
func (_PerpEngine *PerpEngineFilterer) WatchLiquidated(opts *bind.WatchOpts, sink chan<- *PerpEngineLiquidated, marketId [][32]byte, user []common.Address, liquidator []common.Address) (event.Subscription, error) {

	var marketIdRule []interface{}
	for _, marketIdItem := range marketId {
		marketIdRule = append(marketIdRule, marketIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var liquidatorRule []interface{}
	for _, liquidatorItem := range liquidator {
		liquidatorRule = append(liquidatorRule, liquidatorItem)
	}

	logs, sub, err := _PerpEngine.contract.WatchLogs(opts, "Liquidated", marketIdRule, userRule, liquidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpEngineLiquidated)
				if err := _PerpEngine.contract.UnpackLog(event, "Liquidated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLiquidated is a log parse operation binding the contract event 0x89109648b0170073868f7ba7d8a8db89e21725ea8125f7bff949a6fa4ed0915f.
//
// Solidity: event Liquidated(bytes32 indexed marketId, address indexed user, address indexed liquidator, int256 penaltyPaid)
func (_PerpEngine *PerpEngineFilterer) ParseLiquidated(log types.Log) (*PerpEngineLiquidated, error) {
	event := new(PerpEngineLiquidated)
	if err := _PerpEngine.contract.UnpackLog(event, "Liquidated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpEngineMarketAddedIterator is returned from FilterMarketAdded and is used to iterate over the raw logs and unpacked data for MarketAdded events raised by the PerpEngine contract.
type PerpEngineMarketAddedIterator struct {
	Event *PerpEngineMarketAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PerpEngineMarketAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpEngineMarketAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PerpEngineMarketAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PerpEngineMarketAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpEngineMarketAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpEngineMarketAdded represents a MarketAdded event raised by the PerpEngine contract.
type PerpEngineMarketAdded struct {
	MarketId [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMarketAdded is a free log retrieval operation binding the contract event 0x5f1660a6d2d09b7816925ee3bd223865d8b0d3f599d445e115586ec278fe3166.
//
// Solidity: event MarketAdded(bytes32 indexed marketId)
func (_PerpEngine *PerpEngineFilterer) FilterMarketAdded(opts *bind.FilterOpts, marketId [][32]byte) (*PerpEngineMarketAddedIterator, error) {

	var marketIdRule []interface{}
	for _, marketIdItem := range marketId {
		marketIdRule = append(marketIdRule, marketIdItem)
	}

	logs, sub, err := _PerpEngine.contract.FilterLogs(opts, "MarketAdded", marketIdRule)
	if err != nil {
		return nil, err
	}
	return &PerpEngineMarketAddedIterator{contract: _PerpEngine.contract, event: "MarketAdded", logs: logs, sub: sub}, nil
}

// WatchMarketAdded is a free log subscription operation binding the contract event 0x5f1660a6d2d09b7816925ee3bd223865d8b0d3f599d445e115586ec278fe3166.
//
// Solidity: event MarketAdded(bytes32 indexed marketId)
func (_PerpEngine *PerpEngineFilterer) WatchMarketAdded(opts *bind.WatchOpts, sink chan<- *PerpEngineMarketAdded, marketId [][32]byte) (event.Subscription, error) {

	var marketIdRule []interface{}
	for _, marketIdItem := range marketId {
		marketIdRule = append(marketIdRule, marketIdItem)
	}

	logs, sub, err := _PerpEngine.contract.WatchLogs(opts, "MarketAdded", marketIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpEngineMarketAdded)
				if err := _PerpEngine.contract.UnpackLog(event, "MarketAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMarketAdded is a log parse operation binding the contract event 0x5f1660a6d2d09b7816925ee3bd223865d8b0d3f599d445e115586ec278fe3166.
//
// Solidity: event MarketAdded(bytes32 indexed marketId)
func (_PerpEngine *PerpEngineFilterer) ParseMarketAdded(log types.Log) (*PerpEngineMarketAdded, error) {
	event := new(PerpEngineMarketAdded)
	if err := _PerpEngine.contract.UnpackLog(event, "MarketAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpEngineOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PerpEngine contract.
type PerpEngineOwnershipTransferredIterator struct {
	Event *PerpEngineOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PerpEngineOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpEngineOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PerpEngineOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PerpEngineOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpEngineOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpEngineOwnershipTransferred represents a OwnershipTransferred event raised by the PerpEngine contract.
type PerpEngineOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PerpEngine *PerpEngineFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PerpEngineOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PerpEngine.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PerpEngineOwnershipTransferredIterator{contract: _PerpEngine.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PerpEngine *PerpEngineFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PerpEngineOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PerpEngine.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpEngineOwnershipTransferred)
				if err := _PerpEngine.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PerpEngine *PerpEngineFilterer) ParseOwnershipTransferred(log types.Log) (*PerpEngineOwnershipTransferred, error) {
	event := new(PerpEngineOwnershipTransferred)
	if err := _PerpEngine.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PerpEnginePositionChangedIterator is returned from FilterPositionChanged and is used to iterate over the raw logs and unpacked data for PositionChanged events raised by the PerpEngine contract.
type PerpEnginePositionChangedIterator struct {
	Event *PerpEnginePositionChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PerpEnginePositionChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PerpEnginePositionChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PerpEnginePositionChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PerpEnginePositionChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PerpEnginePositionChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PerpEnginePositionChanged represents a PositionChanged event raised by the PerpEngine contract.
type PerpEnginePositionChanged struct {
	MarketId    [32]byte
	User        common.Address
	NewSize     *big.Int
	EntryPrice  *big.Int
	RealizedPnL *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPositionChanged is a free log retrieval operation binding the contract event 0x566a4529e1f90f0edc374385f88308a257a29357fcc5bd874dbb629bb2e4eefa.
//
// Solidity: event PositionChanged(bytes32 indexed marketId, address indexed user, int256 newSize, int256 entryPrice, int256 realizedPnL)
func (_PerpEngine *PerpEngineFilterer) FilterPositionChanged(opts *bind.FilterOpts, marketId [][32]byte, user []common.Address) (*PerpEnginePositionChangedIterator, error) {

	var marketIdRule []interface{}
	for _, marketIdItem := range marketId {
		marketIdRule = append(marketIdRule, marketIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PerpEngine.contract.FilterLogs(opts, "PositionChanged", marketIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &PerpEnginePositionChangedIterator{contract: _PerpEngine.contract, event: "PositionChanged", logs: logs, sub: sub}, nil
}

// WatchPositionChanged is a free log subscription operation binding the contract event 0x566a4529e1f90f0edc374385f88308a257a29357fcc5bd874dbb629bb2e4eefa.
//
// Solidity: event PositionChanged(bytes32 indexed marketId, address indexed user, int256 newSize, int256 entryPrice, int256 realizedPnL)
func (_PerpEngine *PerpEngineFilterer) WatchPositionChanged(opts *bind.WatchOpts, sink chan<- *PerpEnginePositionChanged, marketId [][32]byte, user []common.Address) (event.Subscription, error) {

	var marketIdRule []interface{}
	for _, marketIdItem := range marketId {
		marketIdRule = append(marketIdRule, marketIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PerpEngine.contract.WatchLogs(opts, "PositionChanged", marketIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PerpEnginePositionChanged)
				if err := _PerpEngine.contract.UnpackLog(event, "PositionChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePositionChanged is a log parse operation binding the contract event 0x566a4529e1f90f0edc374385f88308a257a29357fcc5bd874dbb629bb2e4eefa.
//
// Solidity: event PositionChanged(bytes32 indexed marketId, address indexed user, int256 newSize, int256 entryPrice, int256 realizedPnL)
func (_PerpEngine *PerpEngineFilterer) ParsePositionChanged(log types.Log) (*PerpEnginePositionChanged, error) {
	event := new(PerpEnginePositionChanged)
	if err := _PerpEngine.contract.UnpackLog(event, "PositionChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
