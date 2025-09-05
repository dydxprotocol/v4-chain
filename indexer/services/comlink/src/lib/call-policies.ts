import { CallPolicyVersion, type CallPolicyParams } from '@zerodev/permissions/policies';
import type { Abi } from 'viem';

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
] as const;

export const arbitrumCallPolicy: CallPolicyParams<Abi, `0x${string}`> = {
  policyVersion: CallPolicyVersion.V0_0_4,
  permissions: [
    {
      target: '0xaf88d065e77c8cC2239327C5EDb3A432268e5831' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'approve',
    },
    {
      target: '0x4c58aE019E54D10594F1Aa26ABF385B6fb17A52d' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'submitOrder',
    },
    {
      target: '0x4c58aE019E54D10594F1Aa26ABF385B6fb17A52d' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'swapAndSubmitOrder',
    },
  ],
} as const;

export const baseCallPolicy: CallPolicyParams<Abi, `0x${string}`> = {
  policyVersion: CallPolicyVersion.V0_0_4,
  permissions: [
    {
      target: '0x833589fcd6edb6e08f4c7c32d4f71b54bda02913' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'approve',
    },
    {
      target: '0x9335C0c0CBc0317291fd48c00b2f71C8b39DA6F8' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'submitOrder',
    },
    {
      target: '0x9335C0c0CBc0317291fd48c00b2f71C8b39DA6F8' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'swapAndSubmitOrder',
    },
  ],
} as const;

export const optimismCallPolicy: CallPolicyParams<Abi, `0x${string}`> = {
  policyVersion: CallPolicyVersion.V0_0_4,
  permissions: [
    {
      target: '0x0b2c639c533813f4aa9d7837caf62653d097ff85' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'approve',
    },
    {
      target: '0x9c540EdC86613b22968Da784b2d42AC79965af91' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'submitOrder',
    },
    {
      target: '0x9c540EdC86613b22968Da784b2d42AC79965af91' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'swapAndSubmitOrder',
    },
  ],
} as const;

export const avalancheCallPolicy: CallPolicyParams<Abi, `0x${string}`> = {
  policyVersion: CallPolicyVersion.V0_0_4,
  permissions: [
    {
      target: '0x0b2c639c533813f4aa9d7837caf62653d097ff85' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'approve',
    },
    {
      target: '0x9c540EdC86613b22968Da784b2d42AC79965af91' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'submitOrder',
    },
    {
      target: '0x9c540EdC86613b22968Da784b2d42AC79965af91' as `0x${string}`,
      abi,
      valueLimit: BigInt(1000000000000000000000000000000),
      functionName: 'swapAndSubmitOrder',
    },
  ],
} as const;
