import {
  baseConfigSchema,
  parseBigInt,
  parseBoolean,
  parseInteger,
  parseSchema,
  parseString,
} from '@dydxprotocol-indexer/base';
import { complianceConfigSchema } from '@dydxprotocol-indexer/compliance';
import {
  postgresConfigSchema,
} from '@dydxprotocol-indexer/postgres';
import { redisConfigSchema } from '@dydxprotocol-indexer/redis';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...redisConfigSchema,
  ...complianceConfigSchema,

  CHAIN_ID: parseString({ default: 'dydxprotocol' }),
  API_LIMIT_V4: parseInteger({
    default: 1000,
  }),
  // Hard ceiling on (page - 1) * limit for paginated endpoints. Validated against 30d of real
  // comlink traffic: sits well below the offset range (~91k-149k) that triggered the 2026-07-11
  // weekend paging incident, while clearing the extreme scraper tiers observed (offsets in the
  // millions) by 2+ orders of magnitude.
  MAX_PAGINATION_OFFSET: parseInteger({ default: 25_000 }),
  API_ORDERBOOK_LEVELS_PER_SIDE_LIMIT: parseInteger({ default: 100 }),

  // Logging config
  LOG_GETS: parseBoolean({ default: false }),

  // Express server config
  PORT: parseInteger({ default: 8080 }),
  CORS_ORIGIN: parseString({ default: '*' }),
  KEEP_ALIVE_MS: parseInteger({ default: 61_000 }),
  HEADERS_TIMEOUT_MS: parseInteger({ default: 65_000 }),
  // gzip responses for clients that send Accept-Encoding (cuts ALB/NAT egress)
  COMPRESSION_ENABLED: parseBoolean({ default: true }),

  // Rate limit Redis URL
  RATE_LIMIT_REDIS_URL: parseString({
    default: 'redis://localhost:6382',
  }),
  // Rate limits
  RATE_LIMIT_ENABLED: parseBoolean({ default: true }),
  // IP addresses internal to the Indexer have no rate-limit
  INDEXER_INTERNAL_IPS: parseString({ default: '' }),
  // Points / duration determines the maximum rate of requests given that each requests costs 1
  // point
  RATE_LIMIT_GET_POINTS: parseInteger({ default: 100 }),
  RATE_LIMIT_GET_DURATION_SECONDS: parseInteger({ default: 10 }), // 100 requests / 10 seconds

  // Rate limit for screening new / refreshed addresses
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_POINTS: parseInteger({ default: 2 }),
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_DURATION_SECONDS: parseInteger({ default: 60 }), // 2 reqs / min
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_POINTS: parseInteger({ default: 100 }),
  // 100 req / min
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_DURATION_SECONDS: parseInteger({ default: 60 }),
  // Threshold for refreshing compliance data for an address when screened
  MAX_AGE_SCREENED_ADDRESS_COMPLIANCE_DATA_SECONDS: parseInteger({ default: 86_400 }), //  1 day
  // Expose setting compliance status, only set to true in dev/staging.
  EXPOSE_SET_COMPLIANCE_ENDPOINT: parseBoolean({ default: false }),

  // TODO review and finalize per-route rate limits
  // Rate limits for costly endpoints
  RATE_LIMIT_ORDERS_POINTS: parseInteger({ default: 100 }),
  RATE_LIMIT_ORDERS_DURATION_SECONDS: parseInteger({ default: 10 }),
  RATE_LIMIT_FILLS_POINTS: parseInteger({ default: 100 }),
  RATE_LIMIT_FILLS_DURATION_SECONDS: parseInteger({ default: 10 }),
  RATE_LIMIT_CANDLES_POINTS: parseInteger({ default: 1000 }),
  RATE_LIMIT_CANDLES_DURATION_SECONDS: parseInteger({ default: 10 }),
  RATE_LIMIT_SPARKLINES_POINTS: parseInteger({ default: 100 }),
  RATE_LIMIT_SPARKLINES_DURATION_SECONDS: parseInteger({ default: 10 }),
  RATE_LIMIT_HISTORICAL_PNL_POINTS: parseInteger({ default: 100 }),
  RATE_LIMIT_HISTORICAL_PNL_DURATION_SECONDS: parseInteger({ default: 10 }),
  RATE_LIMIT_PNL_POINTS: parseInteger({ default: 100 }),
  RATE_LIMIT_PNL_DURATION_SECONDS: parseInteger({ default: 10 }),
  RATE_LIMIT_FUNDING_POINTS: parseInteger({ default: 100 }),
  RATE_LIMIT_FUNDING_DURATION_SECONDS: parseInteger({ default: 10 }),

  // Affiliates config
  VOLUME_ELIGIBILITY_THRESHOLD: parseInteger({ default: 10_000 }),

  // Vaults config
  VAULT_PNL_HISTORY_DAYS: parseInteger({ default: 90 }),
  VAULT_PNL_HISTORY_HOURS: parseInteger({ default: 72 }),
  VAULT_PNL_START_DATE: parseString({ default: '2024-01-01T00:00:00Z' }),
  VAULT_LATEST_PNL_TICK_WINDOW_HOURS: parseInteger({ default: 1 }),
  VAULT_FETCH_FUNDING_INDEX_BLOCK_WINDOWS: parseInteger({ default: 250_000 }),
  VAULT_CACHE_TTL_MS: parseInteger({ default: 120_000 }), // 2 minutes
  // Cache-Control directives
  CACHE_CONTROL_DIRECTIVE_ADDRESSES: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_AFFILIATES: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_AFFILIATES_METADATA: parseString({ default: 'no-cache' }),
  CACHE_CONTROL_DIRECTIVE_ASSET_POSITIONS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_CANDLES: parseString({ default: 'public, max-age=1' }),
  // omit compliance
  CACHE_CONTROL_DIRECTIVE_FILLS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_FUNDING: parseString({ default: 'public, max-age=10' }),
  // omit height
  CACHE_CONTROL_DIRECTIVE_HISTORICAL_BLOCK_TRADING_REWARDS: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_HISTORICAL_FUNDING: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_HISTORICAL_PNL: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_HISTORICAL_TRADING_REWARDS: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_ORDERBOOK: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_ORDERS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_PERPETUAL_MARKETS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_PERPETUAL_POSITIONS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_SOCIAL_TRADING: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_SPARKLINES: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_TIME: parseString({ default: 'no-cache, no-store, no-transform' }),
  CACHE_CONTROL_DIRECTIVE_TRADES: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_TRANSFERS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_VAULTS: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_PNL: parseString({ default: 'public, max-age=10' }),

  // Turnkey

  TURNKEY_API_BASE_URL: parseString({ default: 'https://api.turnkey.com' }),
  // API keys for root user on parent org to use to create suborgs.
  TURNKEY_API_PRIVATE_KEY: parseString({ default: '' }),
  TURNKEY_API_PUBLIC_KEY: parseString({ default: '' }),
  // API keys for senders to use to start bridging.
  TURNKEY_API_SENDER_PRIVATE_KEY: parseString({ default: '' }),
  TURNKEY_API_SENDER_PUBLIC_KEY: parseString({ default: '' }),
  TURNKEY_MAGIC_LINK_TEMPLATE: parseString({ default: '' }),
  TURNKEY_MAGIC_LINK_VALIDATION: parseString({ default: 'https://dydx.trade' }),
  TURNKEY_ORGANIZATION_ID: parseString({ default: '' }),
  SOLANA_SPONSOR_PUBLIC_KEY: parseString({ default: '' }),
  // Default slippage tolerance for Skip bridge routes, in percent units
  // ('0.1' = 0.1%). The effective tolerance per route is the smaller of this
  // and SKIP_SLIPPAGE_TOLERANCE_USDC / estimatedAmountOut
  // (see getSlippageTolerancePercent in helpers/skip-helper.ts).
  SKIP_SLIPPAGE_TOLERANCE_PERCENTAGE: parseString({ default: '0.1' }),
  // this is the largest allowed slippage amount in USDC.
  SKIP_SLIPPAGE_TOLERANCE_USDC: parseInteger({ default: 100 }),
  TURNKEY_EMAIL_SENDER_ADDRESS: parseString({ default: 'notifications@mail.dydx.trade' }),
  TURNKEY_EMAIL_SENDER_NAME: parseString({ default: 'dYdX Notifications' }),
  // Alchemy auth token for the skip bridge.
  ALCHEMY_AUTH_TOKEN: parseString({ default: '' }),
  ALCHEMY_API_KEY: parseString({ default: '' }),
  ALCHEMY_WEBHOOK_UPDATE_URL: parseString({ default: 'https://dashboard.alchemy.com/api/update-webhook-addresses' }),
  ALCHEMY_WEBHOOK_UPDATE_TIMEOUT_MS: parseInteger({ default: 5000 }),
  // ZeroDev RPC for skip bridge.
  ZERODEV_API_KEY: parseString({ default: '' }),
  ZERODEV_API_BASE_URL: parseString({ default: 'https://rpc.zerodev.app/api/v3' }),
  BRIDGE_THRESHOLD_USDC: parseInteger({ default: 20 }),
  CALL_POLICY_VALUE_LIMIT: parseBigInt({ default: BigInt(100_000_000_000_000_000_000) }),
  // on-chain signer to kick off the skip bridge.
  APPROVAL_SIGNER_PUBLIC_ADDRESS: parseString({ default: '0x3FC11ff27e5373c88EA142d2EdF5492d0839980B' }),
  // if policy approvals are enabled.
  APPROVAL_ENABLED: parseBoolean({ default: true }),
  // largest amount we will tolerate to swap in usdc.
  MAXIMUM_BRIDGE_AMOUNT_USDC: parseInteger({ default: 999_000 }),

  // Apple Sign-In configuration
  APPLE_TEAM_ID: parseString({ default: '' }),
  APPLE_SERVICE_ID: parseString({ default: '' }),
  APPLE_KEY_ID: parseString({ default: '' }),
  APPLE_PRIVATE_KEY: parseString({ default: '' }),
  APPLE_APP_SCHEME: parseString({ default: 'dydxV4' }),

  // webhook ids, defaults to the production webhook id.
  ETHEREUM_WEBHOOK_ID: parseString({ default: 'wh_ctbkt6y9hez91xr2' }),
  ARBITRUM_WEBHOOK_ID: parseString({ default: 'wh_ltwqwcsrx1b8lgry' }),
  AVALANCHE_WEBHOOK_ID: parseString({ default: 'wh_52wz9dbxywxov2dm' }),
  BASE_WEBHOOK_ID: parseString({ default: 'wh_lpjn5gnwj0ll0gap' }),
  OPTIMISM_WEBHOOK_ID: parseString({ default: 'wh_7eo900bsg8rkvo6z' }),
  SOLANA_WEBHOOK_ID: parseString({ default: 'wh_eqxyotjv478gscpo' }),
  // minimum threshold we need to hit for go fast to be free.
  ETHEREUM_GO_FAST_FREE_MINIMUM: parseInteger({ default: 100 }),

  // Amplitude configuration
  AMPLITUDE_API_KEY: parseString({ default: '' }),
  AMPLITUDE_SERVER_URL: parseString({ default: 'https://api.eu.amplitude.com/2/httpapi' }),

  // Gas sponsorship: fee allowances on Injective, so that an address holding USDC but no INJ can
  // send that USDC to the dYdX chain. Off by default; this flag is also the kill switch. Turning
  // it on additionally needs the REST URL and the sponsor public key below.
  GAS_SPONSORSHIP_ENABLED: parseBoolean({ default: false }),
  // REST (LCD) endpoint of an Injective node.
  GAS_SPONSORSHIP_INJECTIVE_REST_URL: parseString({ default: '' }),
  GAS_SPONSORSHIP_INJECTIVE_CHAIN_ID: parseString({ default: 'injective-1' }),
  // Compressed secp256k1 public key, in hex, of the account that sponsors gas. Its private key
  // is held in the Turnkey parent organization and signs through the TURNKEY_API_* credentials.
  GAS_SPONSORSHIP_SPONSOR_PUBLIC_KEY: parseString({ default: '' }),
  // How Turnkey identifies the sponsor key: a wallet account address, a private key address, or
  // a private key id. Defaults to the key's Ethereum address.
  GAS_SPONSORSHIP_TURNKEY_SIGN_WITH: parseString({ default: '' }),
  GAS_SPONSORSHIP_FEE_DENOM: parseString({ default: 'inj' }),
  // Denom of USDC on Injective. An address must hold at least the minimum balance of it, in its
  // smallest unit (6 decimals), to be sponsored.
  GAS_SPONSORSHIP_USDC_DENOM: parseString({
    default: 'erc20:0xa00C59fF5a080D2b954d0c75e46E22a0c371235a',
  }),
  GAS_SPONSORSHIP_MIN_USDC_BALANCE: parseBigInt({ default: BigInt(1_000_000) }),
  // Comma-separated message type URLs the allowance may pay for.
  GAS_SPONSORSHIP_ALLOWED_MESSAGES: parseString({
    default: '/ibc.applications.transfer.v1.MsgTransfer',
  }),
  // An allowance covers this many transactions at this gas limit, priced at the current base fee
  // times the buffer.
  GAS_SPONSORSHIP_TRANSFER_GAS_LIMIT: parseInteger({ default: 450_000 }),
  GAS_SPONSORSHIP_TRANSFER_ATTEMPTS: parseInteger({ default: 3 }),
  GAS_SPONSORSHIP_GAS_PRICE_BUFFER_PERCENT: parseInteger({ default: 150 }),
  // Sponsorship pauses while the base fee per unit of gas, in the smallest unit of the fee denom,
  // is above this. Bounds the size of an allowance. The default is ten times the minimum.
  GAS_SPONSORSHIP_MAX_BASE_FEE: parseBigInt({ default: BigInt(1_600_000_000) }),
  // Gas limit of the transaction that grants the allowance.
  GAS_SPONSORSHIP_GRANT_GAS_LIMIT: parseInteger({ default: 250_000 }),
  GAS_SPONSORSHIP_GRANT_EXPIRY_SECONDS: parseInteger({ default: 86_400 }),
  // How long after being given an allowance an address must wait for another.
  GAS_SPONSORSHIP_COOLDOWN_SECONDS: parseInteger({ default: 3_600 }),
  // The most that may be committed in any 24 hours, in the smallest unit of the fee denom,
  // counting each allowance in full plus the fee paid to grant it. The default is 1 INJ.
  GAS_SPONSORSHIP_DAILY_BUDGET: parseBigInt({ default: BigInt('1000000000000000000') }),
  // The sponsor's transactions are serialized with a Redis lock. The expiry must exceed the
  // longest a single grant can take, including the wait for confirmation.
  GAS_SPONSORSHIP_LOCK_EXPIRY_MS: parseInteger({ default: 60_000 }),
  GAS_SPONSORSHIP_LOCK_WAIT_MS: parseInteger({ default: 10_000 }),
  // How long to wait for the grant transaction to be included before reporting it as pending.
  GAS_SPONSORSHIP_CONFIRMATION_TIMEOUT_MS: parseInteger({ default: 15_000 }),
  GAS_SPONSORSHIP_CONFIRMATION_POLL_MS: parseInteger({ default: 500 }),
  // Timeout of each request to a node.
  GAS_SPONSORSHIP_REQUEST_TIMEOUT_MS: parseInteger({ default: 5_000 }),
  // REST endpoint of a dYdX chain node. When set, sponsorship pauses while the chain is not
  // accepting USDC deposits from Injective. When empty, that check is skipped.
  GAS_SPONSORSHIP_DYDX_REST_URL: parseString({ default: '' }),
  RATE_LIMIT_GAS_SPONSORSHIP_POINTS: parseInteger({ default: 5 }),
  RATE_LIMIT_GAS_SPONSORSHIP_DURATION_SECONDS: parseInteger({ default: 60 }),
};

////////////////////////////////////////////////////////////////////////////////
//                             CONFIG PROCESSING                              //
////////////////////////////////////////////////////////////////////////////////

// Process the top-level configuration.
const config = parseSchema(configSchema);

export default config;
