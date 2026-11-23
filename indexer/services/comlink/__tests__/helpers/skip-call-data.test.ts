import * as skipClient from '@skip-go/client/cjs';
import { decodeFunctionData, parseAbi } from 'viem';

import { getAddress } from '../../src/helpers/alchemy-helpers';
import { buildUserAddresses, getSkipCallData } from '../../src/helpers/skip-helper';
import {
  DYDX_ADDRESS,
  INJECTIVE_ADDRESS,
  INJECTIVE_EVM_ADDRESS,
  SKIP_ARBITRUM_BRIDGE_TX,
} from './deposit-policy-helpers';

jest.mock('@skip-go/client/cjs', () => ({
  __esModule: true,
  route: jest.fn(),
  messages: jest.fn(),
}));

const ARBITRUM: string = '42161';
const ARBITRUM_USDC: string = '0xaf88d065e77c8cC2239327C5EDb3A432268e5831';
const INJECTIVE_EVM: string = '1776';
const INJECTIVE_USDC: string = '0xa00C59fF5a080D2b954d0c75e46E22a0c371235a';
const SOURCE_ADDRESS: string = '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266';
const AMOUNT: string = '100000000';

describe('deposit route to Injective', () => {
  const route = skipClient.route as unknown as jest.Mock;
  const messages = skipClient.messages as unknown as jest.Mock;

  beforeEach(() => {
    route.mockResolvedValue({
      amountIn: AMOUNT,
      estimatedAmountOut: '99999622',
      sourceAssetChainId: ARBITRUM,
      sourceAssetDenom: ARBITRUM_USDC,
      destAssetChainId: INJECTIVE_EVM,
      destAssetDenom: INJECTIVE_USDC,
      operations: [{ cctpTransferV2: {} }],
      chainIds: [ARBITRUM, INJECTIVE_EVM],
      requiredChainAddresses: [ARBITRUM, INJECTIVE_EVM],
      txsRequired: 1,
    });
    messages.mockResolvedValue({
      msgs: [{
        evmTx: {
          chainId: ARBITRUM,
          to: SKIP_ARBITRUM_BRIDGE_TX.to,
          value: '0',
          data: SKIP_ARBITRUM_BRIDGE_TX.data,
        },
      }],
    });
  });

  describe('getAddress', () => {
    it('gives the user\'s own Injective account for Injective', () => {
      expect(getAddress('injective-1', SOURCE_ADDRESS, DYDX_ADDRESS)).toEqual(INJECTIVE_ADDRESS);
    });

    it('gives the same account in EVM form for Injective\'s EVM chain id', () => {
      expect(getAddress(INJECTIVE_EVM, SOURCE_ADDRESS, DYDX_ADDRESS)).toEqual(
        INJECTIVE_EVM_ADDRESS,
      );
    });

    it('still gives the source address for the source chain', () => {
      expect(getAddress(ARBITRUM, SOURCE_ADDRESS, DYDX_ADDRESS)).toEqual(SOURCE_ADDRESS);
      expect(getAddress('solana', 'solana-address', DYDX_ADDRESS)).toEqual('solana-address');
    });

    it('builds the addresses of a route from a source chain to Injective', async () => {
      expect(
        await buildUserAddresses([ARBITRUM, INJECTIVE_EVM], SOURCE_ADDRESS, DYDX_ADDRESS),
      ).toEqual([
        { chainId: ARBITRUM, address: SOURCE_ADDRESS },
        { chainId: INJECTIVE_EVM, address: INJECTIVE_EVM_ADDRESS },
      ]);
    });
  });

  describe('getSkipCallData', () => {
    it('asks Skip for a route that ends at USDC on Injective', async () => {
      await getSkipCallData(SOURCE_ADDRESS, ARBITRUM_USDC, DYDX_ADDRESS, AMOUNT, ARBITRUM);

      expect(route).toHaveBeenCalledTimes(1);
      expect(route).toHaveBeenCalledWith(expect.objectContaining({
        amountIn: AMOUNT,
        sourceAssetDenom: ARBITRUM_USDC,
        sourceAssetChainId: ARBITRUM,
        destAssetDenom: INJECTIVE_USDC,
        destAssetChainId: INJECTIVE_EVM,
      }));
    });

    it('names the user\'s own Injective account as where the funds arrive', async () => {
      await getSkipCallData(SOURCE_ADDRESS, ARBITRUM_USDC, DYDX_ADDRESS, AMOUNT, ARBITRUM);

      expect(messages).toHaveBeenCalledWith(expect.objectContaining({
        addressList: [SOURCE_ADDRESS, INJECTIVE_EVM_ADDRESS],
        destAssetDenom: INJECTIVE_USDC,
        destAssetChainId: INJECTIVE_EVM,
      }));
    });

    it('returns an approval of the relayer, then the bridge call Skip built', async () => {
      const calls = await getSkipCallData(
        SOURCE_ADDRESS,
        ARBITRUM_USDC,
        DYDX_ADDRESS,
        AMOUNT,
        ARBITRUM,
      );

      expect(calls).toHaveLength(2);
      const [approve, bridge] = calls as { to: string, value: bigint, data: `0x${string}` }[];
      expect(approve.to).toEqual(ARBITRUM_USDC);
      expect(approve.value).toEqual(BigInt(0));
      expect(decodeFunctionData({
        abi: parseAbi(['function approve(address spender, uint256 amount) returns (bool)']),
        data: approve.data,
      }).args).toEqual([SKIP_ARBITRUM_BRIDGE_TX.to, BigInt(AMOUNT)]);
      expect(bridge).toEqual({
        to: SKIP_ARBITRUM_BRIDGE_TX.to,
        value: BigInt(0),
        data: `0x${SKIP_ARBITRUM_BRIDGE_TX.data}`,
      });
    });

    it('refuses a route through a contract the session key may not call', async () => {
      messages.mockResolvedValue({
        msgs: [{
          evmTx: {
            to: '0x6d4d4A979D766e8b87C985a61B262cfCDf8d6446',
            value: '0',
            data: SKIP_ARBITRUM_BRIDGE_TX.data,
          },
        }],
      });

      await expect(
        getSkipCallData(SOURCE_ADDRESS, ARBITRUM_USDC, DYDX_ADDRESS, AMOUNT, ARBITRUM),
      ).rejects.toThrow(
        'Skip route does not go through the CCTP v2 relayer: 0x6d4d4A979D766e8b87C985a61B262cfCDf8d6446',
      );
    });

    it('refuses a route with no transaction on the source chain', async () => {
      messages.mockResolvedValue({ msgs: [{ multiChainMsg: { chainId: 'injective-1' } }] });

      await expect(
        getSkipCallData(SOURCE_ADDRESS, ARBITRUM_USDC, DYDX_ADDRESS, AMOUNT, ARBITRUM),
      ).rejects.toThrow('Skip route does not go through the CCTP v2 relayer: no EVM transaction');
    });

    it('fails when Skip has no route', async () => {
      route.mockResolvedValue(undefined);

      await expect(
        getSkipCallData(SOURCE_ADDRESS, 'arbitrum-native', DYDX_ADDRESS, AMOUNT, ARBITRUM),
      ).rejects.toThrow('Failed to find a route');
    });
  });
});
