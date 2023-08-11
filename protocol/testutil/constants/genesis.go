package constants

// This is a copy of the localnet genesis.json. This can be retrieved from the localnet docker container path:
// /dydxprotocol/chain/.alice/config/genesis.json
const GenesisState = `{
  "genesis_time": "2023-04-10T16:16:14.085693715Z",
  "chain_id": "localdydxprotocol",
  "initial_height": "1",
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1"
    },
    "evidence": {
      "max_age_num_blocks": "100000",
      "max_age_duration": "172800000000000",
      "max_bytes": "1048576"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    },
    "version": {
      "app": "0"
    }
  },
  "app_hash": "",
  "app_state": {
    "assets": {
      "assets": [
        {
          "atomic_resolution": -6,
          "symbol": "USDC",
          "denom": "ibc/usdc-placeholder",
          "has_market": false,
          "id": 0,
          "long_interest": 0,
          "market_id": 0
        }
      ]
    },
    "auth": {
      "params": {
        "max_memo_characters": "256",
        "tx_sig_limit": "7",
        "tx_size_cost_per_byte": "10",
        "sig_verify_cost_ed25519": "590",
        "sig_verify_cost_secp256k1": "1000"
      },
      "accounts": [
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4",
          "pub_key": null,
          "account_number": "0",
          "sequence": "0"
        },
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs",
          "pub_key": null,
          "account_number": "1",
          "sequence": "0"
        },
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "dydx1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw9wvmh70",
          "pub_key": null,
          "account_number": "2",
          "sequence": "0"
        },
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "dydx1wau5mja7j7zdavtfq9lu7ejef05hm6ffenlcsn",
          "pub_key": null,
          "account_number": "3",
          "sequence": "0"
        },
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "dydx1nzuttarf5k2j0nug5yzhr6p74t9avehn9hlh8m",
          "pub_key": null,
          "account_number": "4",
          "sequence": "0"
        }
      ]
    },
    "bank": {
      "params": {
        "send_enabled": [],
        "default_send_enabled": true
      },
      "balances": [
        {
          "address": "dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4",
          "coins": [
            {
              "denom": "ibc/usdc-placeholder",
              "amount": "100000000000000000"
            },
            {
              "denom": "stake",
              "amount": "100000000000"
            }
          ]
        },
        {
          "address": "dydx1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw9wvmh70",
          "coins": [
            {
              "denom": "ibc/usdc-placeholder",
              "amount": "100000000000000000"
            },
            {
              "denom": "stake",
              "amount": "100000000000"
            }
          ]
        },
        {
          "address": "dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6",
          "coins": [
            {
              "denom": "ibc/usdc-placeholder",
              "amount": "1300000000000000000"
            }
          ]
        },
        {
          "address": "dydx1wau5mja7j7zdavtfq9lu7ejef05hm6ffenlcsn",
          "coins": [
            {
              "denom": "ibc/usdc-placeholder",
              "amount": "100000000000000000"
            },
            {
              "denom": "stake",
              "amount": "100000000000"
            }
          ]
        },
        {
          "address": "dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs",
          "coins": [
            {
              "denom": "ibc/usdc-placeholder",
              "amount": "100000000000000000"
            },
            {
              "denom": "stake",
              "amount": "100000000000"
            }
          ]
        },
        {
          "address": "dydx1nzuttarf5k2j0nug5yzhr6p74t9avehn9hlh8m",
          "coins": [
            {
              "denom": "ibc/usdc-placeholder",
              "amount": "900000000000000000"
            },
            {
              "denom": "stake",
              "amount": "100000000000"
            }
          ]
        }
      ],
      "supply": [],
      "denom_metadata": [],
      "send_enabled": []
    },
    "capability": {
      "index": "1",
      "owners": []
    },
    "clob": {
      "clob_pairs": [
        {
          "id": 0,
          "maker_fee_ppm": 200,
          "min_order_base_quantums": 10,
          "perpetual_clob_metadata": {
            "perpetual_id": 0
          },
          "quantum_conversion_exponent": -8,
          "status": "STATUS_ACTIVE",
          "step_base_quantums": 10,
          "subticks_per_tick": 10000,
          "taker_fee_ppm": 500
        },
        {
          "id": 1,
          "maker_fee_ppm": 200,
          "min_order_base_quantums": 1000,
          "perpetual_clob_metadata": {
            "perpetual_id": 1
          },
          "quantum_conversion_exponent": -9,
          "status": "STATUS_ACTIVE",
          "step_base_quantums": 1000,
          "subticks_per_tick": 100000,
          "taker_fee_ppm": 500
        }
      ],
      "liquidations_config": {
        "fillable_price_config": {
          "bankruptcy_adjustment_ppm": 1000000,
          "spread_to_maintenance_margin_ratio_ppm": 100000
        },
        "max_insurance_fund_quantums_for_deleveraging": "0",
        "max_liquidation_fee_ppm": 5000,
        "position_block_limits": {
          "max_position_portion_liquidated_ppm": 1000000,
          "min_position_notional_liquidated": 1000
        },
        "subaccount_block_limits": {
          "max_notional_liquidated": 100000000000000,
          "max_quantums_insurance_lost": 100000000000000
        }
      }
    },
    "consensus": null,
    "crisis": {
      "constant_fee": {
        "amount": "1000",
        "denom": "stake"
      }
    },
    "distribution": {
      "delegator_starting_infos": [],
      "delegator_withdraw_infos": [],
      "fee_pool": {
        "community_pool": []
      },
      "outstanding_rewards": [],
      "params": {
        "base_proposer_reward": "0.000000000000000000",
        "bonus_proposer_reward": "0.000000000000000000",
        "community_tax": "0.020000000000000000",
        "withdraw_addr_enabled": true
      },
      "previous_proposer": "",
      "validator_accumulated_commissions": [],
      "validator_current_rewards": [],
      "validator_historical_rewards": [],
      "validator_slash_events": []
    },
    "epochs": {
      "epoch_info_list": [
        {
          "current_epoch": 0,
          "current_epoch_start_block": 0,
          "duration": 60,
          "fast_forward_next_tick": true,
          "is_initialized": false,
          "name": "funding-sample",
          "next_tick": 30
        },
        {
          "current_epoch": 0,
          "current_epoch_start_block": 0,
          "duration": 3600,
          "fast_forward_next_tick": true,
          "is_initialized": false,
          "name": "funding-tick",
          "next_tick": 0
        }
      ]
    },
    "feegrant": {
      "allowances": []
    },
    "genutil": {
      "gen_txs": [
        {
          "body": {
            "messages": [
              {
                "@type": "/cosmos.staking.v1beta1.MsgCreateValidator",
                "description": {
                  "moniker": "alice",
                  "identity": "",
                  "website": "",
                  "security_contact": "",
                  "details": ""
                },
                "commission": {
                  "rate": "0.100000000000000000",
                  "max_rate": "0.200000000000000000",
                  "max_change_rate": "0.010000000000000000"
                },
                "min_self_delegation": "1",
                "delegator_address": "dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4",
                "validator_address": "dydxvaloper199tqg4wdlnu4qjlxchpd7seg454937hjxg9yhy",
                "pubkey": {
                  "@type": "/cosmos.crypto.ed25519.PubKey",
                  "key": "YiARx8259Z+fGFUxQLrz/5FU2RYRT6f5yzvt7D7CrQM="
                },
                "value": {
                  "denom": "stake",
                  "amount": "500000000"
                }
              }
            ],
            "memo": "17e5e45691f0d01449c84fd4ae87279578cdd7ec@172.17.0.2:26656",
            "timeout_height": "0",
            "extension_options": [],
            "non_critical_extension_options": []
          },
          "auth_info": {
            "signer_infos": [
              {
                "public_key": {
                  "@type": "/cosmos.crypto.secp256k1.PubKey",
                  "key": "A0iQ+HpUfJGcgcH7iiEzY9VwCYWCTwg5LsTjc/q1XwSc"
                },
                "mode_info": {
                  "single": {
                    "mode": "SIGN_MODE_DIRECT"
                  }
                },
                "sequence": "0"
              }
            ],
            "fee": {
              "amount": [],
              "gas_limit": "200000",
              "payer": "",
              "granter": ""
            },
            "tip": null
          },
          "signatures": [
            "rv0jVVNS7qgjbnUIxlIcdPju3UF4zOUt9iU91a0WotdGuFgJW1t5sYf7SvpKHPrNRik8T1eDKCcemW/Hn5d5Ew=="
          ]
        },
        {
          "body": {
            "messages": [
              {
                "@type": "/cosmos.staking.v1beta1.MsgCreateValidator",
                "description": {
                  "moniker": "carl",
                  "identity": "",
                  "website": "",
                  "security_contact": "",
                  "details": ""
                },
                "commission": {
                  "rate": "0.100000000000000000",
                  "max_rate": "0.200000000000000000",
                  "max_change_rate": "0.010000000000000000"
                },
                "min_self_delegation": "1",
                "delegator_address": "dydx1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw9wvmh70",
                "validator_address": "dydxvaloper1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw9tjdp47",
                "pubkey": {
                  "@type": "/cosmos.crypto.ed25519.PubKey",
                  "key": "ytLfs1W6E2I41iteKC/YwjyZ/51+CAYCHYxmRHiBeY4="
                },
                "value": {
                  "denom": "stake",
                  "amount": "500000000"
                }
              }
            ],
            "memo": "47539956aaa8e624e0f1d926040e54908ad0eb44@172.17.0.2:26656",
            "timeout_height": "0",
            "extension_options": [],
            "non_critical_extension_options": []
          },
          "auth_info": {
            "signer_infos": [
              {
                "public_key": {
                  "@type": "/cosmos.crypto.secp256k1.PubKey",
                  "key": "AkA1fsLUhCSWbnemBIAR9CPkK1Ra1LlYZcrAKm/Ymvqn"
                },
                "mode_info": {
                  "single": {
                    "mode": "SIGN_MODE_DIRECT"
                  }
                },
                "sequence": "0"
              }
            ],
            "fee": {
              "amount": [],
              "gas_limit": "200000",
              "payer": "",
              "granter": ""
            },
            "tip": null
          },
          "signatures": [
            "rLDe97moyt85AVWhGu4+5YfMFkvH2mnL91DxZNDEXD8WFEHGaWMs9mGWzfXVAQ0iie5L6g2gzbhA8V7lewC96A=="
          ]
        },
        {
          "body": {
            "messages": [
              {
                "@type": "/cosmos.staking.v1beta1.MsgCreateValidator",
                "description": {
                  "moniker": "dave",
                  "identity": "",
                  "website": "",
                  "security_contact": "",
                  "details": ""
                },
                "commission": {
                  "rate": "0.100000000000000000",
                  "max_rate": "0.200000000000000000",
                  "max_change_rate": "0.010000000000000000"
                },
                "min_self_delegation": "1",
                "delegator_address": "dydx1wau5mja7j7zdavtfq9lu7ejef05hm6ffenlcsn",
                "validator_address": "dydxvaloper1wau5mja7j7zdavtfq9lu7ejef05hm6ffudfwmz",
                "pubkey": {
                  "@type": "/cosmos.crypto.ed25519.PubKey",
                  "key": "yG29kRfZ/hgAE1I7uWjbKQJJL4/gX/05XBnfB+m196A="
                },
                "value": {
                  "denom": "stake",
                  "amount": "500000000"
                }
              }
            ],
            "memo": "5882428984d83b03d0c907c1f0af343534987052@172.17.0.2:26656",
            "timeout_height": "0",
            "extension_options": [],
            "non_critical_extension_options": []
          },
          "auth_info": {
            "signer_infos": [
              {
                "public_key": {
                  "@type": "/cosmos.crypto.secp256k1.PubKey",
                  "key": "A87MchHGMj7i1xBwUfECtXzXJIgli/JVFoSaxUqIN86R"
                },
                "mode_info": {
                  "single": {
                    "mode": "SIGN_MODE_DIRECT"
                  }
                },
                "sequence": "0"
              }
            ],
            "fee": {
              "amount": [],
              "gas_limit": "200000",
              "payer": "",
              "granter": ""
            },
            "tip": null
          },
          "signatures": [
            "sRlNM8tust6bdyh2JiFMzmR3Vx+zojNHnvrDgxHDO/Ugb20FR93TB/bIxJbwLg8tbChhhV5VcD+mpcM/3DC77Q=="
          ]
        },
        {
          "body": {
            "messages": [
              {
                "@type": "/cosmos.staking.v1beta1.MsgCreateValidator",
                "description": {
                  "moniker": "bob",
                  "identity": "",
                  "website": "",
                  "security_contact": "",
                  "details": ""
                },
                "commission": {
                  "rate": "0.100000000000000000",
                  "max_rate": "0.200000000000000000",
                  "max_change_rate": "0.010000000000000000"
                },
                "min_self_delegation": "1",
                "delegator_address": "dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs",
                "validator_address": "dydxvaloper10fx7sy6ywd5senxae9dwytf8jxek3t2ga89u8p",
                "pubkey": {
                  "@type": "/cosmos.crypto.ed25519.PubKey",
                  "key": "+P8YiogqqQY+iD96yEa9OJx6EgieU95u9eR3pzxfDp0="
                },
                "value": {
                  "denom": "stake",
                  "amount": "500000000"
                }
              }
            ],
            "memo": "b69182310be02559483e42c77b7b104352713166@172.17.0.2:26656",
            "timeout_height": "0",
            "extension_options": [],
            "non_critical_extension_options": []
          },
          "auth_info": {
            "signer_infos": [
              {
                "public_key": {
                  "@type": "/cosmos.crypto.secp256k1.PubKey",
                  "key": "AlamQtNuTEHlCbn4ZQ20em/bbQNcaAJO54yMOCoE8OTy"
                },
                "mode_info": {
                  "single": {
                    "mode": "SIGN_MODE_DIRECT"
                  }
                },
                "sequence": "0"
              }
            ],
            "fee": {
              "amount": [],
              "gas_limit": "200000",
              "payer": "",
              "granter": ""
            },
            "tip": null
          },
          "signatures": [
            "27/BwrOjpvLUfvwBagcPmwjFci+eG2uavM1kVWsZQ5xgyF4BGfrQJlywNYE7sUh7Bvfu3sPNaNXaLmJCFWy4Jw=="
          ]
        }
      ]
    },
    "gov": {
      "deposit_params": null,
      "deposits": [],
      "params": {
        "max_deposit_period": "172800s",
        "min_deposit": [
          {
            "amount": "10000000",
            "denom": "stake"
          }
        ],
        "min_initial_deposit_ratio": "0.000000000000000000",
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000",
        "veto_threshold": "0.334000000000000000",
        "voting_period": "172800s"
      },
      "proposals": [],
      "starting_proposal_id": "1",
      "tally_params": null,
      "votes": [],
      "voting_params": null
    },
    "ibc": {
      "channel_genesis": {
        "ack_sequences": [],
        "acknowledgements": [],
        "channels": [],
        "commitments": [],
        "next_channel_sequence": "0",
        "receipts": [],
        "recv_sequences": [],
        "send_sequences": []
      },
      "client_genesis": {
        "clients": [],
        "clients_consensus": [],
        "clients_metadata": [],
        "create_localhost": false,
        "next_client_sequence": "0",
        "params": {
          "allowed_clients": [
            "06-solomachine",
            "07-tendermint"
          ]
        }
      },
      "connection_genesis": {
        "client_connection_paths": [],
        "connections": [],
        "next_connection_sequence": "0",
        "params": {
          "max_expected_time_per_block": "30000000000"
        }
      }
    },
    "params": null,
    "perpetuals": {
      "liquidity_tiers": [
        {
          "id": 0,
          "name": "Large-Cap",
          "initial_margin_ppm": 50000,
          "maintenance_fraction_ppm": 600000,
          "base_position_notional": 1000000000000
        }
      ],
      "params": {
        "funding_rate_clamp_factor_ppm": 6000000,
        "premium_vote_clamp_factor_ppm": 60000000
      },
      "perpetuals": [
        {
          "atomic_resolution": -10,
          "default_funding_ppm": 0,
          "id": 0,
          "liquidity_tier": 0,
          "market_id": 0,
          "ticker": "BTC-USD"
        },
        {
          "atomic_resolution": -9,
          "default_funding_ppm": 0,
          "id": 1,
          "liquidity_tier": 0,
          "market_id": 1,
          "ticker": "ETH-USD"
        }
      ]
    },
    "prices": {
      "exchange_feeds": [
        {
          "id": 0,
          "memo": "Memo for Binance",
          "name": "Binance"
        },
        {
          "id": 1,
          "memo": "Memo for Binance US",
          "name": "BinanceUS"
        },
        {
          "id": 2,
          "memo": "Memo for Bitfinex",
          "name": "Bitfinex"
        }
      ],
      "markets": [
        {
          "exchanges": [
            0,
            1,
            2
          ],
          "exponent": -5,
          "id": 0,
          "min_exchanges": 2,
          "min_price_change_ppm": 1000,
          "pair": "BTC-USD",
          "price": 2000000000
        },
        {
          "exchanges": [
            0,
            1,
            2
          ],
          "exponent": -6,
          "id": 1,
          "min_exchanges": 2,
          "min_price_change_ppm": 1000,
          "pair": "ETH-USD",
          "price": 1500000000
        },
        {
          "exchanges": [
            0,
            1
          ],
          "exponent": -8,
          "id": 2,
          "min_exchanges": 1,
          "min_price_change_ppm": 1000,
          "pair": "LINK-USD",
          "price": 1000000000
        }
      ]
    },
    "sending": {},
    "slashing": {
      "missed_blocks": [],
      "params": {
        "downtime_jail_duration": "60s",
        "min_signed_per_window": "0.050000000000000000",
        "signed_blocks_window": "3000",
        "slash_fraction_double_sign": "0.000000000000000000",
        "slash_fraction_downtime": "0.000000000000000000"
      },
      "signing_infos": []
    },
    "staking": {
      "delegations": [],
      "exported": false,
      "last_total_power": "0",
      "last_validator_powers": [],
      "params": {
        "bond_denom": "stake",
        "historical_entries": 10000,
        "max_entries": 7,
        "max_validators": 100,
        "min_commission_rate": "0.000000000000000000",
        "unbonding_time": "7200s"
      },
      "redelegations": [],
      "unbonding_delegations": [],
      "validators": []
    },
    "subaccounts": {
      "subaccounts": [
        {
          "asset_positions": [
            {
              "asset_id": 0,
              "index": 0,
              "quantums": "100000000000000000"
            }
          ],
          "id": {
            "number": 0,
            "owner": "dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4"
          },
          "margin_enabled": true
        },
        {
          "asset_positions": [
            {
              "asset_id": 0,
              "index": 0,
              "quantums": "100000000000000000"
            }
          ],
          "id": {
            "number": 0,
            "owner": "dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs"
          },
          "margin_enabled": true
        },
        {
          "asset_positions": [
            {
              "asset_id": 0,
              "index": 0,
              "quantums": "100000000000000000"
            }
          ],
          "id": {
            "number": 0,
            "owner": "dydx1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw9wvmh70"
          },
          "margin_enabled": true
        },
        {
          "asset_positions": [
            {
              "asset_id": 0,
              "index": 0,
              "quantums": "100000000000000000"
            }
          ],
          "id": {
            "number": 0,
            "owner": "dydx1wau5mja7j7zdavtfq9lu7ejef05hm6ffenlcsn"
          },
          "margin_enabled": true
        },
        {
          "asset_positions": [
            {
              "asset_id": 0,
              "index": 0,
              "quantums": "900000000000000000"
            }
          ],
          "id": {
            "number": 0,
            "owner": "dydx1nzuttarf5k2j0nug5yzhr6p74t9avehn9hlh8m"
          },
          "margin_enabled": true
        }
      ]
    },
    "tendermint-client": null,
    "transfer": {
      "denom_traces": [],
      "params": {
        "receive_enabled": true,
        "send_enabled": true
      },
      "port_id": "transfer"
    },
    "upgrade": {}
  }
}`
