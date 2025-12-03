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

// CollateralVaultMetaData contains all meta data concerning the CollateralVault contract.
var CollateralVaultMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_collateralToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collateralToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"amountDelta\",\"type\":\"int256\"}],\"name\":\"modifyBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"perpEngine\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_perpEngine\",\"type\":\"address\"}],\"name\":\"setPerpEngine\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// CollateralVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use CollateralVaultMetaData.ABI instead.
var CollateralVaultABI = CollateralVaultMetaData.ABI

// CollateralVault is an auto generated Go binding around an Ethereum contract.
type CollateralVault struct {
	CollateralVaultCaller     // Read-only binding to the contract
	CollateralVaultTransactor // Write-only binding to the contract
	CollateralVaultFilterer   // Log filterer for contract events
}

// CollateralVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type CollateralVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CollateralVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CollateralVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CollateralVaultSession struct {
	Contract     *CollateralVault  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CollateralVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CollateralVaultCallerSession struct {
	Contract *CollateralVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// CollateralVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CollateralVaultTransactorSession struct {
	Contract     *CollateralVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// CollateralVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type CollateralVaultRaw struct {
	Contract *CollateralVault // Generic contract binding to access the raw methods on
}

// CollateralVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CollateralVaultCallerRaw struct {
	Contract *CollateralVaultCaller // Generic read-only contract binding to access the raw methods on
}

// CollateralVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CollateralVaultTransactorRaw struct {
	Contract *CollateralVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCollateralVault creates a new instance of CollateralVault, bound to a specific deployed contract.
func NewCollateralVault(address common.Address, backend bind.ContractBackend) (*CollateralVault, error) {
	contract, err := bindCollateralVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CollateralVault{CollateralVaultCaller: CollateralVaultCaller{contract: contract}, CollateralVaultTransactor: CollateralVaultTransactor{contract: contract}, CollateralVaultFilterer: CollateralVaultFilterer{contract: contract}}, nil
}

// NewCollateralVaultCaller creates a new read-only instance of CollateralVault, bound to a specific deployed contract.
func NewCollateralVaultCaller(address common.Address, caller bind.ContractCaller) (*CollateralVaultCaller, error) {
	contract, err := bindCollateralVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultCaller{contract: contract}, nil
}

// NewCollateralVaultTransactor creates a new write-only instance of CollateralVault, bound to a specific deployed contract.
func NewCollateralVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*CollateralVaultTransactor, error) {
	contract, err := bindCollateralVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultTransactor{contract: contract}, nil
}

// NewCollateralVaultFilterer creates a new log filterer instance of CollateralVault, bound to a specific deployed contract.
func NewCollateralVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*CollateralVaultFilterer, error) {
	contract, err := bindCollateralVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultFilterer{contract: contract}, nil
}

// bindCollateralVault binds a generic wrapper to an already deployed contract.
func bindCollateralVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CollateralVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CollateralVault *CollateralVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CollateralVault.Contract.CollateralVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CollateralVault *CollateralVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralVault.Contract.CollateralVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CollateralVault *CollateralVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CollateralVault.Contract.CollateralVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CollateralVault *CollateralVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CollateralVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CollateralVault *CollateralVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CollateralVault *CollateralVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CollateralVault.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_CollateralVault *CollateralVaultCaller) BalanceOf(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVault.contract.Call(opts, &out, "balanceOf", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_CollateralVault *CollateralVaultSession) BalanceOf(user common.Address) (*big.Int, error) {
	return _CollateralVault.Contract.BalanceOf(&_CollateralVault.CallOpts, user)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_CollateralVault *CollateralVaultCallerSession) BalanceOf(user common.Address) (*big.Int, error) {
	return _CollateralVault.Contract.BalanceOf(&_CollateralVault.CallOpts, user)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_CollateralVault *CollateralVaultCaller) CollateralToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CollateralVault.contract.Call(opts, &out, "collateralToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_CollateralVault *CollateralVaultSession) CollateralToken() (common.Address, error) {
	return _CollateralVault.Contract.CollateralToken(&_CollateralVault.CallOpts)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_CollateralVault *CollateralVaultCallerSession) CollateralToken() (common.Address, error) {
	return _CollateralVault.Contract.CollateralToken(&_CollateralVault.CallOpts)
}

// PerpEngine is a free data retrieval call binding the contract method 0x48ba4d2b.
//
// Solidity: function perpEngine() view returns(address)
func (_CollateralVault *CollateralVaultCaller) PerpEngine(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CollateralVault.contract.Call(opts, &out, "perpEngine")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PerpEngine is a free data retrieval call binding the contract method 0x48ba4d2b.
//
// Solidity: function perpEngine() view returns(address)
func (_CollateralVault *CollateralVaultSession) PerpEngine() (common.Address, error) {
	return _CollateralVault.Contract.PerpEngine(&_CollateralVault.CallOpts)
}

// PerpEngine is a free data retrieval call binding the contract method 0x48ba4d2b.
//
// Solidity: function perpEngine() view returns(address)
func (_CollateralVault *CollateralVaultCallerSession) PerpEngine() (common.Address, error) {
	return _CollateralVault.Contract.PerpEngine(&_CollateralVault.CallOpts)
}

// TotalCollateral is a free data retrieval call binding the contract method 0x4ac8eb5f.
//
// Solidity: function totalCollateral() view returns(uint256)
func (_CollateralVault *CollateralVaultCaller) TotalCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVault.contract.Call(opts, &out, "totalCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalCollateral is a free data retrieval call binding the contract method 0x4ac8eb5f.
//
// Solidity: function totalCollateral() view returns(uint256)
func (_CollateralVault *CollateralVaultSession) TotalCollateral() (*big.Int, error) {
	return _CollateralVault.Contract.TotalCollateral(&_CollateralVault.CallOpts)
}

// TotalCollateral is a free data retrieval call binding the contract method 0x4ac8eb5f.
//
// Solidity: function totalCollateral() view returns(uint256)
func (_CollateralVault *CollateralVaultCallerSession) TotalCollateral() (*big.Int, error) {
	return _CollateralVault.Contract.TotalCollateral(&_CollateralVault.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_CollateralVault *CollateralVaultTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _CollateralVault.contract.Transact(opts, "deposit", amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_CollateralVault *CollateralVaultSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _CollateralVault.Contract.Deposit(&_CollateralVault.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_CollateralVault *CollateralVaultTransactorSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _CollateralVault.Contract.Deposit(&_CollateralVault.TransactOpts, amount)
}

// ModifyBalance is a paid mutator transaction binding the contract method 0x599888a6.
//
// Solidity: function modifyBalance(address user, int256 amountDelta) returns()
func (_CollateralVault *CollateralVaultTransactor) ModifyBalance(opts *bind.TransactOpts, user common.Address, amountDelta *big.Int) (*types.Transaction, error) {
	return _CollateralVault.contract.Transact(opts, "modifyBalance", user, amountDelta)
}

// ModifyBalance is a paid mutator transaction binding the contract method 0x599888a6.
//
// Solidity: function modifyBalance(address user, int256 amountDelta) returns()
func (_CollateralVault *CollateralVaultSession) ModifyBalance(user common.Address, amountDelta *big.Int) (*types.Transaction, error) {
	return _CollateralVault.Contract.ModifyBalance(&_CollateralVault.TransactOpts, user, amountDelta)
}

// ModifyBalance is a paid mutator transaction binding the contract method 0x599888a6.
//
// Solidity: function modifyBalance(address user, int256 amountDelta) returns()
func (_CollateralVault *CollateralVaultTransactorSession) ModifyBalance(user common.Address, amountDelta *big.Int) (*types.Transaction, error) {
	return _CollateralVault.Contract.ModifyBalance(&_CollateralVault.TransactOpts, user, amountDelta)
}

// SetPerpEngine is a paid mutator transaction binding the contract method 0xb619daf7.
//
// Solidity: function setPerpEngine(address _perpEngine) returns()
func (_CollateralVault *CollateralVaultTransactor) SetPerpEngine(opts *bind.TransactOpts, _perpEngine common.Address) (*types.Transaction, error) {
	return _CollateralVault.contract.Transact(opts, "setPerpEngine", _perpEngine)
}

// SetPerpEngine is a paid mutator transaction binding the contract method 0xb619daf7.
//
// Solidity: function setPerpEngine(address _perpEngine) returns()
func (_CollateralVault *CollateralVaultSession) SetPerpEngine(_perpEngine common.Address) (*types.Transaction, error) {
	return _CollateralVault.Contract.SetPerpEngine(&_CollateralVault.TransactOpts, _perpEngine)
}

// SetPerpEngine is a paid mutator transaction binding the contract method 0xb619daf7.
//
// Solidity: function setPerpEngine(address _perpEngine) returns()
func (_CollateralVault *CollateralVaultTransactorSession) SetPerpEngine(_perpEngine common.Address) (*types.Transaction, error) {
	return _CollateralVault.Contract.SetPerpEngine(&_CollateralVault.TransactOpts, _perpEngine)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_CollateralVault *CollateralVaultTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _CollateralVault.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_CollateralVault *CollateralVaultSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _CollateralVault.Contract.Withdraw(&_CollateralVault.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_CollateralVault *CollateralVaultTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _CollateralVault.Contract.Withdraw(&_CollateralVault.TransactOpts, amount)
}

// CollateralVaultDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the CollateralVault contract.
type CollateralVaultDepositIterator struct {
	Event *CollateralVaultDeposit // Event containing the contract specifics and raw log

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
func (it *CollateralVaultDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralVaultDeposit)
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
		it.Event = new(CollateralVaultDeposit)
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
func (it *CollateralVaultDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralVaultDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralVaultDeposit represents a Deposit event raised by the CollateralVault contract.
type CollateralVaultDeposit struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed user, uint256 amount)
func (_CollateralVault *CollateralVaultFilterer) FilterDeposit(opts *bind.FilterOpts, user []common.Address) (*CollateralVaultDepositIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CollateralVault.contract.FilterLogs(opts, "Deposit", userRule)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultDepositIterator{contract: _CollateralVault.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed user, uint256 amount)
func (_CollateralVault *CollateralVaultFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *CollateralVaultDeposit, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CollateralVault.contract.WatchLogs(opts, "Deposit", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralVaultDeposit)
				if err := _CollateralVault.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed user, uint256 amount)
func (_CollateralVault *CollateralVaultFilterer) ParseDeposit(log types.Log) (*CollateralVaultDeposit, error) {
	event := new(CollateralVaultDeposit)
	if err := _CollateralVault.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralVaultWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the CollateralVault contract.
type CollateralVaultWithdrawIterator struct {
	Event *CollateralVaultWithdraw // Event containing the contract specifics and raw log

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
func (it *CollateralVaultWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralVaultWithdraw)
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
		it.Event = new(CollateralVaultWithdraw)
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
func (it *CollateralVaultWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralVaultWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralVaultWithdraw represents a Withdraw event raised by the CollateralVault contract.
type CollateralVaultWithdraw struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed user, uint256 amount)
func (_CollateralVault *CollateralVaultFilterer) FilterWithdraw(opts *bind.FilterOpts, user []common.Address) (*CollateralVaultWithdrawIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CollateralVault.contract.FilterLogs(opts, "Withdraw", userRule)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultWithdrawIterator{contract: _CollateralVault.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed user, uint256 amount)
func (_CollateralVault *CollateralVaultFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *CollateralVaultWithdraw, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CollateralVault.contract.WatchLogs(opts, "Withdraw", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralVaultWithdraw)
				if err := _CollateralVault.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed user, uint256 amount)
func (_CollateralVault *CollateralVaultFilterer) ParseWithdraw(log types.Log) (*CollateralVaultWithdraw, error) {
	event := new(CollateralVaultWithdraw)
	if err := _CollateralVault.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
