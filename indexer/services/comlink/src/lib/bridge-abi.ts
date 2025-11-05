export const abi = [
  {
    inputs: [
      {
        internalType: 'bytes32',
        name: 'sender',
        type: 'bytes32',
      },
      {
        internalType: 'bytes32',
        name: 'recipient',
        type: 'bytes32',
      },
      {
        internalType: 'uint256',
        name: 'amountIn',
        type: 'uint256',
      },
      {
        internalType: 'uint256',
        name: 'amountOut',
        type: 'uint256',
      },
      {
        internalType: 'uint32',
        name: 'destinationDomain',
        type: 'uint32',
      },
      {
        internalType: 'uint64',
        name: 'timeoutTimestamp',
        type: 'uint64',
      },
      {
        internalType: 'bytes',
        name: 'data',
        type: 'bytes',
      },
    ],
    name: 'submitOrder',
    outputs: [
      {
        internalType: 'bytes32',
        name: '',
        type: 'bytes32',
      },
    ],
    stateMutability: 'nonpayable',
    type: 'function',
  },
  {
    inputs: [
      {
        internalType: 'address',
        name: 'tokenIn',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'swapAmountIn',
        type: 'uint256',
      },
      {
        internalType: 'bytes',
        name: 'swapCalldata',
        type: 'bytes',
      },
      {
        internalType: 'uint256',
        name: 'executionFeeAmount',
        type: 'uint256',
      },
      {
        internalType: 'uint256',
        name: 'solverFeeBPS',
        type: 'uint256',
      },
      {
        internalType: 'bytes32',
        name: 'sender',
        type: 'bytes32',
      },
      {
        internalType: 'bytes32',
        name: 'recipient',
        type: 'bytes32',
      },
      {
        internalType: 'uint32',
        name: 'destinationDomain',
        type: 'uint32',
      },
      {
        internalType: 'uint64',
        name: 'timeoutTimestamp',
        type: 'uint64',
      },
      {
        internalType: 'bytes',
        name: 'destinationCalldata',
        type: 'bytes',
      },
    ],
    name: 'swapAndSubmitOrder',
    outputs: [
      {
        internalType: 'bytes32',
        name: '',
        type: 'bytes32',
      },
    ],
    stateMutability: 'payable',
    type: 'function',
  },
  {
    type: 'function',
    name: 'approve',
    stateMutability: 'nonpayable',
    inputs: [
      {
        name: 'spender',
        type: 'address',
      },
      {
        name: 'amount',
        type: 'uint256',
      },
    ],
    outputs: [
      {
        type: 'bool',
      },
    ],
  },
  {
    inputs: [
      {
        internalType: 'address',
        name: 'inputToken',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'inputAmount',
        type: 'uint256',
      },
      {
        internalType: 'bytes',
        name: 'swapCalldata',
        type: 'bytes',
      },
      {
        internalType: 'uint32',
        name: 'destinationDomain',
        type: 'uint32',
      },
      {
        internalType: 'bytes32',
        name: 'mintRecipient',
        type: 'bytes32',
      },
      {
        internalType: 'address',
        name: 'burnToken',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'feeAmount',
        type: 'uint256',
      },
      {
        internalType: 'bytes32',
        name: 'destinationCaller',
        type: 'bytes32',
      },
    ],
    name: 'swapAndRequestCCTPTransferWithCaller',
    outputs: [],
    stateMutability: 'payable',
    type: 'function',
  },
  {
    inputs: [
      {
        internalType: 'uint256',
        name: 'transferAmount',
        type: 'uint256',
      },
      {
        internalType: 'uint32',
        name: 'destinationDomain',
        type: 'uint32',
      },
      {
        internalType: 'bytes32',
        name: 'mintRecipient',
        type: 'bytes32',
      },
      {
        internalType: 'address',
        name: 'burnToken',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'feeAmount',
        type: 'uint256',
      },
      {
        internalType: 'bytes32',
        name: 'destinationCaller',
        type: 'bytes32',
      },
    ],
    name: 'requestCCTPTransferWithCaller',
    outputs: [],
    stateMutability: 'nonpayable',
    type: 'function',
  },
] as const;
