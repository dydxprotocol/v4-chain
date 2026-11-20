import fetch from 'node-fetch';

import {
  getExpectedSequence,
  InjectiveClient,
  InjectiveClientError,
  isInjectiveDepositOpen,
  parseDecimalRoundingUp,
} from '../../../src/lib/gas-sponsorship/injective-client';

jest.mock('node-fetch');

const mockFetch = fetch as jest.MockedFunction<typeof fetch>;

const REST_URL: string = 'https://injective.example';
const TIMEOUT_MS: number = 1234;
const GRANTER: string = 'inj1xdrkxgke8qzey9ze4cx48heav6lzkkszuc224r';
const GRANTEE: string = 'inj1frwsqymyg3wnfhy4mqrq9cghwwzst0tgsmfal7';
const USDC_DENOM: string = 'erc20:0xa00C59fF5a080D2b954d0c75e46E22a0c371235a';
const TX_HASH: string = 'C554B3ACB4F98161CCFE7938399A73C4D452FBD9EF92F88C4F79E68896905BE4';

function respondWith(status: number, body: object | string): void {
  const text: string = typeof body === 'string' ? body : JSON.stringify(body);
  mockFetch.mockResolvedValue({
    ok: status >= 200 && status < 300,
    status,
    text: () => Promise.resolve(text),
    json: () => Promise.resolve(JSON.parse(text)),
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  } as any);
}

function requestedUrl(): string {
  return mockFetch.mock.calls[0][0] as string;
}

/**
 * The response bodies for an Injective account, a missing account, balances, a message-limited
 * allowance, a missing allowance, the base fee, an included transaction, a missing transaction
 * and the canonical USDC state are what mainnet nodes returned, trimmed. The others are written
 * to the shape of the same REST API.
 */
describe('injective-client', () => {
  const client: InjectiveClient = new InjectiveClient(`${REST_URL}/`, TIMEOUT_MS);

  describe('getAccount', () => {
    it('reads the number and sequence of an Injective account', async () => {
      respondWith(200, {
        account: {
          '@type': '/injective.types.v1beta1.EthAccount',
          base_account: {
            address: GRANTER,
            pub_key: {
              '@type': '/injective.crypto.v1beta1.ethsecp256k1.PubKey',
              key: 'Aq8nnehaBv1QXUtBy+VKqYU5m2rL3wqAWqC566gKMq4r',
            },
            account_number: '1682332',
            sequence: '55563',
          },
          code_hash: 'xdJGAYb3IzySfn2y3McDwOUAtlPKgic7e/rYBF2FpHA=',
        },
      });

      expect(await client.getAccount(GRANTER)).toEqual({
        accountNumber: BigInt(1682332),
        sequence: BigInt(55563),
      });
      expect(requestedUrl()).toEqual(`${REST_URL}/cosmos/auth/v1beta1/accounts/${GRANTER}`);
      expect(mockFetch.mock.calls[0][1]).toEqual(expect.objectContaining({
        method: 'GET',
        timeout: TIMEOUT_MS,
      }));
    });

    it('reads a plain Cosmos account', async () => {
      respondWith(200, {
        account: {
          '@type': '/cosmos.auth.v1beta1.BaseAccount',
          address: GRANTER,
          pub_key: null,
          account_number: '12',
          sequence: '0',
        },
      });

      expect(await client.getAccount(GRANTER)).toEqual({
        accountNumber: BigInt(12),
        sequence: BigInt(0),
      });
    });

    it('returns undefined for an account that does not exist', async () => {
      respondWith(404, { code: 5, message: `account ${GRANTEE} not found`, details: [] });

      expect(await client.getAccount(GRANTEE)).toBeUndefined();
    });

    it('fails on a response it cannot read', async () => {
      respondWith(200, { account: { base_account: { account_number: 'abc', sequence: '1' } } });

      await expect(client.getAccount(GRANTER)).rejects.toThrow(
        new InjectiveClientError('unexpected account number: "abc"'),
      );
    });
  });

  describe('getBalance', () => {
    it('reads a balance', async () => {
      respondWith(200, { balance: { denom: 'inj', amount: '134282993851360000000' } });

      expect(await client.getBalance(GRANTER, 'inj')).toEqual(BigInt('134282993851360000000'));
    });

    it('reads a zero balance', async () => {
      respondWith(200, { balance: { denom: USDC_DENOM, amount: '0' } });

      expect(await client.getBalance(GRANTER, USDC_DENOM)).toEqual(BigInt(0));
    });

    it('escapes the denom in the query', async () => {
      respondWith(200, { balance: { denom: USDC_DENOM, amount: '0' } });

      await client.getBalance(GRANTER, USDC_DENOM);

      expect(requestedUrl()).toEqual(
        `${REST_URL}/cosmos/bank/v1beta1/balances/${GRANTER}/by_denom` +
        '?denom=erc20%3A0xa00C59fF5a080D2b954d0c75e46E22a0c371235a',
      );
    });

    it('fails when the node reports an error', async () => {
      respondWith(500, { code: 2, message: 'boom', details: [] });

      await expect(client.getBalance(GRANTER, 'inj')).rejects.toThrow(
        'balance query failed with status 500',
      );
    });
  });

  describe('getAllowance', () => {
    it('reads an allowance limited to some message types', async () => {
      respondWith(200, {
        allowance: {
          granter: GRANTER,
          grantee: GRANTEE,
          allowance: {
            '@type': '/cosmos.feegrant.v1beta1.AllowedMsgAllowance',
            allowance: {
              '@type': '/cosmos.feegrant.v1beta1.BasicAllowance',
              spend_limit: [{ denom: 'inj', amount: '35200000000000000' }],
              expiration: '2026-11-19T17:51:29.674059734Z',
            },
            allowed_messages: [
              '/cosmos.authz.v1beta1.MsgExec',
              '/injective.wasmx.v1.MsgExecuteContractCompat',
            ],
          },
        },
      });

      expect(await client.getAllowance(GRANTER, GRANTEE)).toEqual({
        granter: GRANTER,
        grantee: GRANTEE,
        recognized: true,
        spendLimit: { inj: BigInt('35200000000000000') },
        expiration: new Date('2026-11-19T17:51:29.674Z'),
        allowedMessages: [
          '/cosmos.authz.v1beta1.MsgExec',
          '/injective.wasmx.v1.MsgExecuteContractCompat',
        ],
      });
      expect(requestedUrl()).toEqual(
        `${REST_URL}/cosmos/feegrant/v1beta1/allowance/${GRANTER}/${GRANTEE}`,
      );
    });

    it('reads a basic allowance with no limit, expiry or message filter', async () => {
      respondWith(200, {
        allowance: {
          granter: GRANTER,
          grantee: GRANTEE,
          allowance: {
            '@type': '/cosmos.feegrant.v1beta1.BasicAllowance',
            spend_limit: [],
            expiration: null,
          },
        },
      });

      expect(await client.getAllowance(GRANTER, GRANTEE)).toEqual({
        granter: GRANTER,
        grantee: GRANTEE,
        recognized: true,
        spendLimit: undefined,
        expiration: undefined,
        allowedMessages: undefined,
      });
    });

    it.each([
      ['a periodic allowance', { '@type': '/cosmos.feegrant.v1beta1.PeriodicAllowance' }],
      [
        'a message filter around a periodic allowance',
        {
          '@type': '/cosmos.feegrant.v1beta1.AllowedMsgAllowance',
          allowance: { '@type': '/cosmos.feegrant.v1beta1.PeriodicAllowance' },
          allowed_messages: ['/cosmos.authz.v1beta1.MsgExec'],
        },
      ],
    ])('marks %s as not recognized', async (_name: string, allowance: object) => {
      respondWith(200, { allowance: { granter: GRANTER, grantee: GRANTEE, allowance } });

      expect(await client.getAllowance(GRANTER, GRANTEE)).toEqual({
        granter: GRANTER,
        grantee: GRANTEE,
        recognized: false,
      });
    });

    it('returns undefined when there is no allowance, which the node reports as an error', async () => {
      respondWith(500, { code: 13, message: 'fee-grant not found: not found', details: [] });

      expect(await client.getAllowance(GRANTER, GRANTEE)).toBeUndefined();
    });

    it('fails on any other error', async () => {
      respondWith(500, { code: 13, message: 'internal error', details: [] });

      await expect(client.getAllowance(GRANTER, GRANTEE)).rejects.toThrow(
        'allowance query failed with status 500',
      );
    });
  });

  describe('getBaseFee', () => {
    it('reads the base fee', async () => {
      respondWith(200, { base_fee: { base_fee: '160000000.000000000000000000' } });

      expect(await client.getBaseFee()).toEqual(BigInt(160000000));
      expect(requestedUrl()).toEqual(`${REST_URL}/injective/txfees/v1beta1/cur_eip_base_fee`);
    });

    it('rounds a fractional base fee up', async () => {
      respondWith(200, { base_fee: { base_fee: '160000000.000000000000000001' } });

      expect(await client.getBaseFee()).toEqual(BigInt(160000001));
    });

    it('fails on a response it cannot read', async () => {
      respondWith(200, { base_fee: {} });

      await expect(client.getBaseFee()).rejects.toThrow('unexpected base fee: undefined');
    });
  });

  describe('broadcastTx', () => {
    it('posts the transaction and reads a result of acceptance', async () => {
      respondWith(200, {
        tx_response: {
          height: '0', txhash: TX_HASH, codespace: '', code: 0, data: '', raw_log: '',
        },
      });

      expect(await client.broadcastTx(Uint8Array.from([1, 2, 3]))).toEqual({
        txHash: TX_HASH,
        code: 0,
        codespace: '',
        rawLog: '',
      });
      expect(requestedUrl()).toEqual(`${REST_URL}/cosmos/tx/v1beta1/txs`);
      expect(mockFetch.mock.calls[0][1]).toEqual(expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ tx_bytes: 'AQID', mode: 'BROADCAST_MODE_SYNC' }),
        timeout: TIMEOUT_MS,
      }));
    });

    it('reads a result of rejection', async () => {
      respondWith(200, {
        tx_response: {
          height: '0',
          txhash: TX_HASH,
          codespace: 'sdk',
          code: 32,
          raw_log: 'account sequence mismatch, expected 55564, got 55563: incorrect account sequence',
        },
      });

      expect(await client.broadcastTx(Uint8Array.from([1]))).toEqual({
        txHash: TX_HASH,
        code: 32,
        codespace: 'sdk',
        rawLog: 'account sequence mismatch, expected 55564, got 55563: incorrect account sequence',
      });
    });

    it('fails when the node cannot be reached', async () => {
      mockFetch.mockRejectedValue(new Error('network timeout'));

      await expect(client.broadcastTx(Uint8Array.from([1]))).rejects.toThrow(
        new InjectiveClientError('request to /cosmos/tx/v1beta1/txs failed: network timeout'),
      );
    });

    it('fails on an error that is not JSON', async () => {
      respondWith(502, '<html>Bad Gateway</html>');

      await expect(client.broadcastTx(Uint8Array.from([1]))).rejects.toThrow(
        'broadcast failed with status 502',
      );
    });

    it('fails on a success that is not JSON', async () => {
      respondWith(200, 'OK');

      await expect(client.broadcastTx(Uint8Array.from([1]))).rejects.toThrow(
        'response from /cosmos/tx/v1beta1/txs is not JSON',
      );
    });
  });

  describe('getTx', () => {
    it('reads an included transaction', async () => {
      respondWith(200, {
        tx: {},
        tx_response: {
          height: '192295827',
          txhash: TX_HASH,
          codespace: '',
          code: 0,
          raw_log: '',
          gas_wanted: '600000',
          gas_used: '119383',
        },
      });

      expect(await client.getTx(TX_HASH)).toEqual({
        txHash: TX_HASH,
        code: 0,
        codespace: '',
        rawLog: '',
        height: BigInt(192295827),
      });
      expect(requestedUrl()).toEqual(`${REST_URL}/cosmos/tx/v1beta1/txs/${TX_HASH}`);
    });

    it('reads an included transaction that failed', async () => {
      respondWith(200, {
        tx_response: {
          height: '192295900',
          txhash: TX_HASH,
          codespace: 'feegrant',
          code: 2,
          raw_log: 'failed to execute message; message index: 0: fee allowance already exists',
        },
      });

      expect(await client.getTx(TX_HASH)).toEqual(expect.objectContaining({
        code: 2,
        codespace: 'feegrant',
      }));
    });

    it('returns undefined for a transaction that is not in a block', async () => {
      respondWith(404, { code: 5, message: `tx not found: ${TX_HASH}`, details: [] });

      expect(await client.getTx(TX_HASH)).toBeUndefined();
    });

    it('fails on any other error', async () => {
      respondWith(500, { code: 13, message: 'internal error', details: [] });

      await expect(client.getTx(TX_HASH)).rejects.toThrow('tx query failed with status 500');
    });
  });

  describe('parseDecimalRoundingUp', () => {
    it.each([
      ['0', BigInt(0)],
      ['7', BigInt(7)],
      ['7.000', BigInt(7)],
      ['7.001', BigInt(8)],
      ['0.5', BigInt(1)],
      ['123456789012345678901234567890.1', BigInt('123456789012345678901234567891')],
    ])('parses %s', (value: string, expected: bigint) => {
      expect(parseDecimalRoundingUp(value, 'value')).toEqual(expected);
    });

    it.each(['', '-1', '1e9', '1.', '.5', 'abc', null, 5])('rejects %p', (value: unknown) => {
      expect(() => parseDecimalRoundingUp(value, 'value')).toThrow('unexpected value');
    });
  });

  describe('getExpectedSequence', () => {
    const rejection = {
      txHash: TX_HASH,
      code: 32,
      codespace: 'sdk',
      rawLog: 'account sequence mismatch, expected 55564, got 55563: incorrect account sequence',
    };

    it('reads the sequence the chain expected from a wrong-sequence rejection', () => {
      expect(getExpectedSequence(rejection)).toEqual(BigInt(55564));
    });

    it.each([
      ['another code', { ...rejection, code: 13 }],
      ['another codespace', { ...rejection, codespace: 'feegrant' }],
      ['a message it cannot read', { ...rejection, rawLog: 'incorrect account sequence' }],
    ])('returns undefined for %s', (_name: string, result: typeof rejection) => {
      expect(getExpectedSequence(result)).toBeUndefined();
    });
  });

  describe('isInjectiveDepositOpen', () => {
    const DYDX_URL: string = 'https://dydx.example';
    // The canonical USDC state of dYdX mainnet, without the ledger and route fields.
    const controls = {
      mode: 'MODE_GRADUAL',
      noble_withdrawals_enabled: true,
      injective_deposits_disabled: false,
      noble_deposits_disabled: true,
    };

    it('is open in gradual mode with Injective deposits enabled', async () => {
      respondWith(200, { controls });

      expect(await isInjectiveDepositOpen(`${DYDX_URL}/`, TIMEOUT_MS)).toEqual(true);
      expect(requestedUrl()).toEqual(`${DYDX_URL}/dydxprotocol/v4/canonicalusdc/state`);
    });

    it('is open when the disabled flag is left out, as a false value may be', async () => {
      respondWith(200, { controls: { mode: 'MODE_GRADUAL' } });

      expect(await isInjectiveDepositOpen(DYDX_URL, TIMEOUT_MS)).toEqual(true);
    });

    it.each([
      ['Injective deposits are disabled', { ...controls, injective_deposits_disabled: true }],
      ['the module is paused', { ...controls, mode: 'MODE_PAUSED' }],
      ['the module is disabled', { ...controls, mode: 'MODE_DISABLED' }],
      ['the mode is one this code does not know', { ...controls, mode: 'MODE_SOMETHING_NEW' }],
    ])('is closed when %s', async (_name: string, closedControls: object) => {
      respondWith(200, { controls: closedControls });

      expect(await isInjectiveDepositOpen(DYDX_URL, TIMEOUT_MS)).toEqual(false);
    });

    it('fails when the state cannot be fetched', async () => {
      respondWith(501, { code: 12, message: 'Not Implemented', details: [] });

      await expect(isInjectiveDepositOpen(DYDX_URL, TIMEOUT_MS)).rejects.toThrow(
        new InjectiveClientError('canonical USDC state query failed with status 501', 501),
      );
    });

    it('fails when the node cannot be reached', async () => {
      mockFetch.mockRejectedValue(new Error('ECONNREFUSED'));

      await expect(isInjectiveDepositOpen(DYDX_URL, TIMEOUT_MS)).rejects.toThrow(
        'canonical USDC state query failed: ECONNREFUSED',
      );
    });

    it('fails on a state it cannot read', async () => {
      respondWith(200, { controls: {} });

      await expect(isInjectiveDepositOpen(DYDX_URL, TIMEOUT_MS)).rejects.toThrow(
        'unexpected canonical USDC state',
      );
    });
  });
});
