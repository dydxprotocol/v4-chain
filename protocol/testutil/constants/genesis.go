package constants

// This is a copy of the localnet genesis.json. This can be retrieved from the localnet docker container path:
// /dydxprotocol/chain/.alice/config/genesis.json
// Disable linter for exchange config.
//
//nolint:all
const GenesisState = `{
  "genesis_time": "2023-07-10T19:23:15.891430637Z",
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
          "denom": "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
          "denom_exponent": "-6",
          "has_market": false,
          "id": 0,
          "market_id": 0,
          "symbol": "USDC"
        }
      ]
    },
    "affiliates": {},
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
              "denom": "adv4tnt",
              "amount": "1000000000000000000000000"
            },
            {
              "denom": "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
              "amount": "100000000000000000"
            }
          ]
        },
        {
          "address": "dydx1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw9wvmh70",
          "coins": [
            {
              "denom": "adv4tnt",
              "amount": "1000000000000000000000000"
            },
            {
              "denom": "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
              "amount": "100000000000000000"
            }
          ]
        },
        {
          "address": "dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6",
          "coins": [
            {
              "denom": "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
              "amount": "1300000000000000000"
            }
          ]
        },
        {
          "address": "dydx1wau5mja7j7zdavtfq9lu7ejef05hm6ffenlcsn",
          "coins": [
            {
              "denom": "adv4tnt",
              "amount": "1000000000000000000000000"
            },
            {
              "denom": "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
              "amount": "100000000000000000"
            }
          ]
        },
        {
          "address": "dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs",
          "coins": [
            {
              "denom": "adv4tnt",
              "amount": "1000000000000000000000000"
            },
            {
              "denom": "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
              "amount": "100000000000000000"
            }
          ]
        },
        {
          "address": "dydx1nzuttarf5k2j0nug5yzhr6p74t9avehn9hlh8m",
          "coins": [
            {
              "denom": "adv4tnt",
              "amount": "100000000000"
            },
            {
              "denom": "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
              "amount": "900000000000000000"
            }
          ]
        },
        {
          "address": "dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv",
          "coins": [
            {
              "denom": "adv4tnt",
              "amount": "1000000000"
            }
          ]
        }
      ],
      "supply": [],
      "denom_metadata": [],
      "send_enabled": []
    },
    "blocktime": {
      "params": {
        "durations": [
          "300s",
          "1800s"
        ]
      }
    },
    "bridge": {
      "event_params": {
        "denom": "bridge-token",
        "eth_chain_id": "11155111",
        "eth_address": "0xEf01c3A30eB57c91c40C52E996d29c202ae72193"
      },
      "propose_params": {
        "max_bridges_per_block": 10,
        "propose_delay_duration": "60s",
        "skip_rate_ppm": 800000,
        "skip_if_block_delayed_by_duration": "5s"
      },
      "safety_params": {
        "is_disabled": false,
        "delay_blocks": 86400
      },
      "acknowledged_event_info": {
        "next_id": 0,
        "eth_block_height": 0
      }
    },
    "capability": {
      "index": "1",
      "owners": []
    },
    "clob": {
      "block_rate_limit_config": {
        "max_short_term_orders_and_cancels_per_n_blocks": [
          {
            "num_blocks": 1,
            "limit": 400
          }
        ],
        "max_stateful_orders_per_n_blocks": [
          {
            "num_blocks": 1,
            "limit": 2
          },
          {
            "num_blocks": 100,
            "limit": 20
          }
        ]
      },
      "clob_pairs": [
        {
          "id": 0,
          "perpetual_clob_metadata": {
            "perpetual_id": 0
          },
          "quantum_conversion_exponent": -8,
          "status": "STATUS_ACTIVE",
          "step_base_quantums": 10,
          "subticks_per_tick": 10000
        },
        {
          "id": 1,
          "perpetual_clob_metadata": {
            "perpetual_id": 1
          },
          "quantum_conversion_exponent": -9,
          "status": "STATUS_ACTIVE",
          "step_base_quantums": 1000,
          "subticks_per_tick": 100000
        }
      ],
      "equity_tier_limit_config": {
        "short_term_order_equity_tiers": [
          {
            "limit": 0,
            "usd_tnc_required": "0"
          },
          {
            "limit": 1,
            "usd_tnc_required": "20"
          },
          {
            "limit": 5,
            "usd_tnc_required": "100"
          },
          {
            "limit": 10,
            "usd_tnc_required": "1000"
          },
          {
            "limit": 100,
            "usd_tnc_required": "10000"
          },
          {
            "limit": 1000,
            "usd_tnc_required": "100000"
          }
        ],
        "stateful_order_equity_tiers": [
          {
            "limit": 0,
            "usd_tnc_required": "0"
          },
          {
            "limit": 1,
            "usd_tnc_required": "20"
          },
          {
            "limit": 5,
            "usd_tnc_required": "100"
          },
          {
            "limit": 10,
            "usd_tnc_required": "1000"
          },
          {
            "limit": 100,
            "usd_tnc_required": "10000"
          },
          {
            "limit": 200,
            "usd_tnc_required": "100000"
          }
        ]
      },
      "liquidations_config": {
        "fillable_price_config": {
          "bankruptcy_adjustment_ppm": 1000000,
          "spread_to_maintenance_margin_ratio_ppm": 100000
        },
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
    "crisis": {
      "constant_fee": {
        "amount": "1000",
        "denom": "adv4tnt"
      }
    },
    "delaymsg": {
      "delayed_messages": [
        {
          "id": 0,
          "msg": {
            "@type": "/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParams",
            "authority": "dydx1mkkvp26dngu6n8rmalaxyp3gwkjuzztq5zx6tr",
            "params": {
              "tiers": [
                {
                  "name": "1",
                  "absolute_volume_requirement": "0",
                  "total_volume_share_requirement_ppm": 0,
                  "maker_volume_share_requirement_ppm": 0,
                  "maker_fee_ppm": 100,
                  "taker_fee_ppm": 500
                },
                {
                  "name": "2",
                  "absolute_volume_requirement": "1000000000000",
                  "total_volume_share_requirement_ppm": 0,
                  "maker_volume_share_requirement_ppm": 0,
                  "maker_fee_ppm": 100,
                  "taker_fee_ppm": 450
                },
                {
                  "name": "3",
                  "absolute_volume_requirement": "5000000000000",
                  "total_volume_share_requirement_ppm": 0,
                  "maker_volume_share_requirement_ppm": 0,
                  "maker_fee_ppm": 50,
                  "taker_fee_ppm": 400
                },
                {
                  "name": "4",
                  "absolute_volume_requirement": "25000000000000",
                  "total_volume_share_requirement_ppm": 0,
                  "maker_volume_share_requirement_ppm": 0,
                  "maker_fee_ppm": 0,
                  "taker_fee_ppm": 350
                },
                {
                  "name": "5",
                  "absolute_volume_requirement": "125000000000000",
                  "total_volume_share_requirement_ppm": 0,
                  "maker_volume_share_requirement_ppm": 0,
                  "maker_fee_ppm": 0,
                  "taker_fee_ppm": 300
                },
                {
                  "name": "6",
                  "absolute_volume_requirement": "125000000000000",
                  "total_volume_share_requirement_ppm": 5000,
                  "maker_volume_share_requirement_ppm": 0,
                  "maker_fee_ppm": -50,
                  "taker_fee_ppm": 250
                },
                {
                  "name": "7",
                  "absolute_volume_requirement": "125000000000000",
                  "total_volume_share_requirement_ppm": 5000,
                  "maker_volume_share_requirement_ppm": 10000,
                  "maker_fee_ppm": -90,
                  "taker_fee_ppm": 250
                },
                {
                  "name": "8",
                  "absolute_volume_requirement": "125000000000000",
                  "total_volume_share_requirement_ppm": 5000,
                  "maker_volume_share_requirement_ppm": 20000,
                  "maker_fee_ppm": -110,
                  "taker_fee_ppm": 250
                },
                {
                  "name": "9",
                  "absolute_volume_requirement": "125000000000000",
                  "total_volume_share_requirement_ppm": 5000,
                  "maker_volume_share_requirement_ppm": 40000,
                  "maker_fee_ppm": -110,
                  "taker_fee_ppm": 250
                }
              ]
            }
          },
          "block_height": "6480000"
        }
      ],
      "next_delayed_message_id": 1
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
    "dydxaccountplus": {
      "accounts": [],
      "params": {
        "is_smart_account_active": true
      },
      "next_authenticator_id": "0",
      "authenticator_data": []
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
        },
        {
          "current_epoch": 0,
          "current_epoch_start_block": 0,
          "duration": 3600,
          "fast_forward_next_tick": true,
          "is_initialized": false,
          "name": "stats-epoch",
          "next_tick": 0
        }
      ]
    },
    "feegrant": {
      "allowances": []
    },
    "feetiers": {
      "params": {
        "tiers": [
          {
            "name": "1",
            "absolute_volume_requirement": "0",
            "total_volume_share_requirement_ppm": 0,
            "maker_volume_share_requirement_ppm": 0,
            "maker_fee_ppm": -110,
            "taker_fee_ppm": 500
          },
          {
            "name": "2",
            "absolute_volume_requirement": "1000000000000",
            "total_volume_share_requirement_ppm": 0,
            "maker_volume_share_requirement_ppm": 0,
            "maker_fee_ppm": -110,
            "taker_fee_ppm": 450
          },
          {
            "name": "3",
            "absolute_volume_requirement": "5000000000000",
            "total_volume_share_requirement_ppm": 0,
            "maker_volume_share_requirement_ppm": 0,
            "maker_fee_ppm": -110,
            "taker_fee_ppm": 400
          },
          {
            "name": "4",
            "absolute_volume_requirement": "25000000000000",
            "total_volume_share_requirement_ppm": 0,
            "maker_volume_share_requirement_ppm": 0,
            "maker_fee_ppm": -110,
            "taker_fee_ppm": 350
          },
          {
            "name": "5",
            "absolute_volume_requirement": "125000000000000",
            "total_volume_share_requirement_ppm": 0,
            "maker_volume_share_requirement_ppm": 0,
            "maker_fee_ppm": -110,
            "taker_fee_ppm": 300
          },
          {
            "name": "6",
            "absolute_volume_requirement": "125000000000000",
            "total_volume_share_requirement_ppm": 5000,
            "maker_volume_share_requirement_ppm": 0,
            "maker_fee_ppm": -110,
            "taker_fee_ppm": 250
          },
          {
            "name": "7",
            "absolute_volume_requirement": "125000000000000",
            "total_volume_share_requirement_ppm": 5000,
            "maker_volume_share_requirement_ppm": 10000,
            "maker_fee_ppm": -110,
            "taker_fee_ppm": 250
          },
          {
            "name": "8",
            "absolute_volume_requirement": "125000000000000",
            "total_volume_share_requirement_ppm": 5000,
            "maker_volume_share_requirement_ppm": 20000,
            "maker_fee_ppm": -110,
            "taker_fee_ppm": 250
          },
          {
            "name": "9",
            "absolute_volume_requirement": "125000000000000",
            "total_volume_share_requirement_ppm": 5000,
            "maker_volume_share_requirement_ppm": 40000,
            "maker_fee_ppm": -110,
            "taker_fee_ppm": 250
          }
        ]
      },
      "staking_tiers": []
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
                  "rate": "1.000000000000000000",
                  "max_rate": "1.000000000000000000",
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
                  "denom": "adv4tnt",
                  "amount": "500000000000000000000000"
                }
              }
            ],
            "memo": "17e5e45691f0d01449c84fd4ae87279578cdd7ec@172.17.0.3:26656",
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
            "KqBNRNfXrxYaU2585ycZI2BOvJrUrvZWVugMr9d09gxcDSPGqdjleJWFFwO+Hbhj58uZ4wNOplv9e0SxPwZ0KQ=="
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
                  "rate": "1.000000000000000000",
                  "max_rate": "1.000000000000000000",
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
                  "denom": "adv4tnt",
                  "amount": "500000000000000000000000"
                }
              }
            ],
            "memo": "47539956aaa8e624e0f1d926040e54908ad0eb44@172.17.0.3:26656",
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
            "QcK0WTKaXjsPPsALhn7jLJ/hhmhww+1ucTy4VZE9cJlivPcurFr1k4kfP1/M0ppqEWa9mksjIeVQhOHXTOBG/Q=="
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
                  "rate": "1.000000000000000000",
                  "max_rate": "1.000000000000000000",
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
                  "denom": "adv4tnt",
                  "amount": "500000000000000000000000"
                }
              }
            ],
            "memo": "5882428984d83b03d0c907c1f0af343534987052@172.17.0.3:26656",
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
            "jpCPNmoS5CHqbDXwFX5FYO7J5g7kSi5ZkxVkXEkgajJOZgu9nVTXavPFZ2t5w+UDzgWbtDxLJ1GqdM+kNFIWaA=="
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
                  "rate": "1.000000000000000000",
                  "max_rate": "1.000000000000000000",
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
                  "denom": "adv4tnt",
                  "amount": "500000000000000000000000"
                }
              }
            ],
            "memo": "b69182310be02559483e42c77b7b104352713166@172.17.0.3:26656",
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
            "SAhIiKZUWVy8tI0uDanDo9IP2ZCh2ZltV2yY2Y6umqgax34GxbU1BbrAxXxPrrhEa+IFqXJEWpocVsGo++gjuQ=="
          ]
        }
      ]
    },
    "gov": {
      "deposits": [],
      "params": {
        "burn_proposal_deposit_prevote": false,
        "burn_vote_quorum": false,
        "burn_vote_veto": true,
        "max_deposit_period": "172800s",
        "min_deposit": [
          {
            "amount": "10000000",
            "denom": "adv4tnt"
          }
        ],
        "min_initial_deposit_ratio": "0.000000000000000000",
        "proposal_cancel_ratio": "1.000000000000000000",
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000",
        "veto_threshold": "0.334000000000000000",
		"min_deposit_ratio": "0.010000000000000000",
        "expedited_voting_period": "86400s",
        "expedited_threshold": "0.750000000000000000",
        "expedited_min_deposit": [
          {
            "amount": "50000000",
            "denom": "adv4tnt"
          }
        ],
        "voting_period": "172800s"
      },
      "proposals": [],
      "starting_proposal_id": "1",
      "votes": []
    },
    "govplus": {},
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
    "listing": {
      "hard_cap_for_markets": 500,
      "listing_vault_deposit_params": {
        "new_vault_deposit_amount": "10000000000",
        "main_vault_deposit_amount": "0",
        "num_blocks_to_lock_shares": 2592000
      }
	},
    "perpetuals": {
      "liquidity_tiers": [
        {
          "id": 0,
          "impact_notional": 10000000000,
          "initial_margin_ppm": 50000,
          "maintenance_fraction_ppm": 600000,
          "name": "Large-Cap",
          "open_interest_lower_cap": 0,
          "open_interest_upper_cap": 0
        },
        {
          "id": 1,
          "impact_notional": 5000000000,
          "initial_margin_ppm": 100000,
          "maintenance_fraction_ppm": 500000,
          "name": "Small-Cap",
          "open_interest_lower_cap": 20000000000000,
          "open_interest_upper_cap": 50000000000000
        },
        {
          "id": 2,
          "impact_notional": 2500000000,
          "initial_margin_ppm": 200000,
          "maintenance_fraction_ppm": 500000,
          "name": "Long-Tail",
          "open_interest_lower_cap": 5000000000000,
          "open_interest_upper_cap": 10000000000000
        },
        {
          "id": 3,
          "impact_notional": 2500000000,
          "initial_margin_ppm": 1000000,
          "maintenance_fraction_ppm": 200000,
          "name": "Safety",
          "open_interest_lower_cap": 2000000000000,
          "open_interest_upper_cap": 5000000000000
        },
        {
          "id": 4,
          "impact_notional": 2500000000,
          "initial_margin_ppm": 50000,
          "maintenance_fraction_ppm": 600000,
          "name": "Isolated",
          "open_interest_lower_cap": 500000000000,
          "open_interest_upper_cap": 1000000000000
        },
        {
          "id": 5,
          "impact_notional": 5000000000,
          "initial_margin_ppm": 50000,
          "maintenance_fraction_ppm": 600000,
          "name": "Mid-Cap",
          "open_interest_lower_cap": 40000000000000,
          "open_interest_upper_cap": 100000000000000
        },
        {
          "id": 6,
          "impact_notional": 2500000000,
          "initial_margin_ppm": 10000,
          "maintenance_fraction_ppm": 500000,
          "name": "FX",
          "open_interest_lower_cap": 500000000000,
          "open_interest_upper_cap": 1000000000000
        }
      ],
      "params": {
        "funding_rate_clamp_factor_ppm": 6000000,
        "min_num_votes_per_sample": 15,
        "premium_vote_clamp_factor_ppm": 60000000
      },
      "perpetuals": [
        {
          "params": {
            "atomic_resolution": -10,
            "default_funding_ppm": 0,
            "id": 0,
            "liquidity_tier": 0,
            "market_id": 0,
            "ticker": "BTC-USD",
            "market_type": 1
          }
        },
        {
          "params": {
            "atomic_resolution": -9,
            "default_funding_ppm": 0,
            "id": 1,
            "liquidity_tier": 0,
            "market_id": 1,
            "ticker": "ETH-USD",
		"market_type": 1
          }
        }
      ]
    },
    "marketmap": {
      "market_map": {
        "markets": {
          "AAVE/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "AAVE",
                "Quote": "USD"
              },
              "decimals": 8,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "AAVEUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "AAVE-USD"
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "aaveusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "AAVEUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "AAVE-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "AAVE-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ADA/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ADA",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "ADAUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "ADAUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "ADA-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "ADA_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "adausdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "ADAUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "ADA-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "ADAUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "ADA-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ALGO/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ALGO",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "ALGOUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "ALGO-USD"
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "algousdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "ALGOUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "ALGO-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "ALGO-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "APE/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "APE",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "APEUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "APE-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "APE_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "APEUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "APE-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "APEUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "APE-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "APT/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "APT",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "APTUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "APTUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "APT-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "APT_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "aptusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "APT-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "APTUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "APT-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ARB/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ARB",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "ARBUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "ARBUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "ARB-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "ARB_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "arbusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "ARB-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "ARBUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "ARB-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ATOM/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ATOM",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "ATOMUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "ATOMUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "ATOM-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "ATOM_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "ATOMUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "ATOM-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "ATOMUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "ATOM-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "AVAX/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "AVAX",
                "Quote": "USD"
              },
              "decimals": 8,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "AVAXUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "AVAXUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "AVAX-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "AVAX_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "avaxusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "AVAXUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "AVAX-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "AVAX-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "BCH/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "BCH",
                "Quote": "USD"
              },
              "decimals": 7,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "BCHUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "BCHUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "BCH-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "BCH_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "bchusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "BCHUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "BCH-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "BCHUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "BCH-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "BLUR/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "BLUR",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "BLUR-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "BLUR_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "BLURUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "BLUR-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "BLURUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "BLUR-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "BTC/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "BTC",
                "Quote": "USD"
              },
              "decimals": 5,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "BTCUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "BTCUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "BTC-USD"
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "btcusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "XXBTZUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "BTC-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "BTCUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "BTC-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "COMP/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "COMP",
                "Quote": "USD"
              },
              "decimals": 8,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "COMPUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "COMP-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "COMP_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "COMPUSD"
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "COMPUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "COMP-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "CRV/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "CRV",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "CRVUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "CRV-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "CRV_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "CRVUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "CRV-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "CRVUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "CRV-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "DOGE/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "DOGE",
                "Quote": "USD"
              },
              "decimals": 11,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "DOGEUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "DOGEUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "DOGE-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "DOGE_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "dogeusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "XDGUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "DOGE-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "DOGEUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "DOGE-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "DOT/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "DOT",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "DOTUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "DOTUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "DOT-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "DOT_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "DOTUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "DOT-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "DOTUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "DOT-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "DYDX/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "DYDX",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "DYDXUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "DYDXUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "DYDX_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "DYDX-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "DYDXUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "DYDX-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "EOS/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "EOS",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "EOSUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "EOS-USD"
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "eosusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "EOSUSD"
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "EOS-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ETC/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ETC",
                "Quote": "USD"
              },
              "decimals": 8,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "ETCUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "ETC-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "ETC_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "etcusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "ETC-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "ETCUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "ETC-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ETH/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ETH",
                "Quote": "USD"
              },
              "decimals": 6,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "ETHUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "ETHUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "ETH-USD"
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "ethusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "XETHZUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "ETH-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "ETHUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "ETH-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "FIL/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "FIL",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "FILUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "FIL-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "FIL_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "filusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "FILUSD"
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "FILUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "FIL-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ICP/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ICP",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "ICPUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "ICP-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "ICP_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "icpusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "ICP-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "ICP-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ISO/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ISO",
                "Quote": "USD"
              },
              "decimals": 8,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "ISOUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ISO2/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ISO2",
                "Quote": "USD"
              },
              "decimals": 7,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "ISO2USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "LDO/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "LDO",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "LDOUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "LDO-USD"
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "LDOUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "LDO-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "LDOUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "LDO-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "LINK/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "LINK",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "LINKUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "LINKUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "LINK-USD"
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "LINKUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "LINK-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "LINKUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "LINK-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "LTC/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "LTC",
                "Quote": "USD"
              },
              "decimals": 8,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "LTCUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "LTCUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "LTC-USD"
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "ltcusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "XLTCZUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "LTC-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "LTCUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "LTC-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "MKR/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "MKR",
                "Quote": "USD"
              },
              "decimals": 6,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "MKRUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "MKR-USD"
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "MKRUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "MKR-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "MKRUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "MKR-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "NEAR/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "NEAR",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "NEARUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "NEAR-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "NEAR_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "nearusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "NEAR-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "NEARUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "NEAR-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "OP/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "OP",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "OPUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "OP-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "OP_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "OP-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "OPUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "OP-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "PEPE/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "PEPE",
                "Quote": "USD"
              },
              "decimals": 16,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "PEPEUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "PEPEUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "PEPE_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "PEPEUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "PEPE-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "PEPEUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "PEPE-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "POL/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "POL",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "POLUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "POLUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "POL-USD"
              },
              {
                "name": "crypto_dot_com_ws",
                "off_chain_ticker": "POL_USD"
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "POL-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "SEI/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "SEI",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "SEIUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "SEIUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "SEI-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "SEI_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "seiusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "SEI-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "SEIUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "SHIB/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "SHIB",
                "Quote": "USD"
              },
              "decimals": 15,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "SHIBUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "SHIBUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "SHIB-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "SHIB_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "SHIBUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "SHIB-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "SHIBUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "SHIB-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "SNX/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "SNX",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "SNXUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "SNX-USD"
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "snxusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "SNXUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "SNX-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "SNX-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "SOL/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "SOL",
                "Quote": "USD"
              },
              "decimals": 8,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "SOLUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "SOLUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "SOL-USD"
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "solusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "SOLUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "SOL-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "SOLUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "SOL-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "SUI/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "SUI",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "SUIUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "SUIUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "SUI-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "SUI_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "suiusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "SUI-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "SUIUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "SUI-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "SUSHI/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "SUSHI",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "SUSHIUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "SUSHI-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "SUSHI_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "sushiusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "SUSHI-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "TEST/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "TEST",
                "Quote": "USD"
              },
              "decimals": 5,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "volatile-exchange-provider",
                "off_chain_ticker": "TEST-USD"
              }
            ]
          },
          "TRX/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "TRX",
                "Quote": "USD"
              },
              "decimals": 11,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "TRXUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "TRXUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "TRX_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "trxusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "TRXUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "TRX-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "TRXUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "TRX-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "UNI/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "UNI",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "UNIUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "UNIUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "UNI-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "UNI_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "UNIUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "UNI-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "UNI-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "USDT/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "USDT",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "USDCUSDT",
                "invert": true
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "USDCUSDT",
                "invert": true
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "USDT-USD"
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "ethusdt",
                "normalize_by_pair": {
                  "Base": "ETH",
                  "Quote": "USD"
                },
                "invert": true
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "USDTZUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "BTC-USDT",
                "normalize_by_pair": {
                  "Base": "BTC",
                  "Quote": "USD"
                },
                "invert": true
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "USDC-USDT",
                "invert": true
              }
            ]
          },
          "WLD/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "WLD",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "WLDUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "WLDUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "WLD_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "wldusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "WLD-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "WLDUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "WLD-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "XLM/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "XLM",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "XLMUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "XLMUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "XLM-USD"
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "XXLMZUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "XLM-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "XLMUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "XLM-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "XMR/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "XMR",
                "Quote": "USD"
              },
              "decimals": 7,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "XMRUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "XMR_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "XXMRZUSD",
                "normalize_by_pair": {
                  "Base": "ZUSD",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "XMR-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "XMR-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "XRP/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "XRP",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 3,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "XRPUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "bybit_ws",
                "off_chain_ticker": "XRPUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "XRP-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "XRP_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "xrpusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "XXRPZUSD"
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "XRP-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "mexc_ws",
                "off_chain_ticker": "XRPUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "XRP-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "XTZ/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "XTZ",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "XTZUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "XTZ-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "XTZ_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "xtzusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "XTZUSD"
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "XTZ-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "YFI/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "YFI",
                "Quote": "USD"
              },
              "decimals": 6,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "YFIUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "YFI-USD"
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "yfiusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "YFIUSD"
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "YFI-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ZEC/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ZEC",
                "Quote": "USD"
              },
              "decimals": 8,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "ZECUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "ZEC-USD"
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "XZECZUSD",
                "normalize_by_pair": {
                  "Base": "ZUSD",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "ZEC-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "ZEC-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ZRX/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ZRX",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "ZRXUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "ZRX-USD"
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "zrxusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kraken_api",
                "off_chain_ticker": "ZRXUSD"
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "ZRX-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          },
          "ZUSD/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "ZUSD",
                "Quote": "USD"
              },
              "decimals": 9,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "kraken_api",
                "off_chain_ticker": "USDTZUSD",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                },
                "invert": true
              }
            ]
          },
          "1INCH/USD": {
            "ticker": {
              "currency_pair": {
                "Base": "1INCH",
                "Quote": "USD"
              },
              "decimals": 10,
              "min_provider_count": 1,
              "enabled": true
            },
            "provider_configs": [
              {
                "name": "binance_ws",
                "off_chain_ticker": "1INCHUSDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "coinbase_ws",
                "off_chain_ticker": "1INCH-USD"
              },
              {
                "name": "gate_ws",
                "off_chain_ticker": "1INCH_USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "huobi_ws",
                "off_chain_ticker": "1inchusdt",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "kucoin_ws",
                "off_chain_ticker": "1INCH-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              },
              {
                "name": "okx_ws",
                "off_chain_ticker": "1INCH-USDT",
                "normalize_by_pair": {
                  "Base": "USDT",
                  "Quote": "USD"
                }
              }
            ]
          }
        }
      },
      "params": {
        "admin": "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky",
        "market_authorities": ["dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky"]
      }
    },
    "prices": {
      "market_params": [
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"BTCUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"BTCUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tBTCUSD\"},{\"exchangeName\":\"Bitstamp\",\"ticker\":\"BTC/USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"BTCUSDT\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"BTC-USD\"},{\"exchangeName\":\"CryptoCom\",\"ticker\":\"BTC_USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XXBTZUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"BTC-USDT\"}]}",
          "exponent": -5,
          "id": 0,
          "min_exchanges": 1,
          "min_price_change_ppm": 1000,
          "pair": "BTC-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"ETHUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"ETHUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tETHUSD\"},{\"exchangeName\":\"Bitstamp\",\"ticker\":\"ETH/USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"ETHUSDT\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ETH-USD\"},{\"exchangeName\":\"CryptoCom\",\"ticker\":\"ETH_USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XETHZUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ETH-USDT\"}]}",
          "exponent": -6,
          "id": 1,
          "min_exchanges": 1,
          "min_price_change_ppm": 1000,
          "pair": "ETH-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"LINKUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"LINKUSD\\\"\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"LINK-USD\"},{\"exchangeName\":\"CryptoCom\",\"ticker\":\"LINK_USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"linkusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"LINKUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"LINK-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"LINK-USDT\"}]}",
          "exponent": -9,
          "id": 2,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "LINK-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"POLUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"POLUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"POL-USD\"},{\"exchangeName\":\"CryptoCom\",\"ticker\":\"POL_USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"POL-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}",
          "exponent": -10,
          "id": 3,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "POL-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"CRVUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"CRVUSD\\\"\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"CRVUSDT\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"CRV-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"CRV_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"crvusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"CRVUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"CRV-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"CRV-USDT\"}]}",
          "exponent": -10,
          "id": 4,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "CRV-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"SOLUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"SOLUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tSOLUSD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"SOL-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"solusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"SOLUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"SOL-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"SOL-USDT\"}]}",
          "exponent": -8,
          "id": 5,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "SOL-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"ADAUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"ADAUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tADAUSD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ADA-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"ADA_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"adausdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"ADAUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"ADA-USDT\"}]}",
          "exponent": -10,
          "id": 6,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "ADA-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"AVAXUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"AVAXUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tAVAX:USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"AVAX_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"avaxusdt\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"AVAX-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"AVAX-USDT\"}]}",
          "exponent": -8,
          "id": 7,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "AVAX-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"FILUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"FILUSD\\\"\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"FIL-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"filusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"FILUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"FIL-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"FIL-USDT\"}]}",
          "exponent": -9,
          "id": 8,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "FIL-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"AAVEUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"AAVEUSD\\\"\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"AAVE-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"aaveusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"AAVEUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"AAVE-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"AAVE-USDT\"}]}",
          "exponent": -8,
          "id": 9,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "AAVE-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"LTCUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"LTCUSD\\\"\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"LTCUSDT\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"LTC-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"ltcusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XLTCZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"LTC-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"LTC-USDT\"}]}",
          "exponent": -8,
          "id": 10,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "LTC-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"DOGEUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"DOGEUSD\\\"\"},{\"exchangeName\":\"Gate\",\"ticker\":\"DOGE_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"dogeusdt\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"DOGE-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"DOGE-USDT\"}]}",
          "exponent": -11,
          "id": 11,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "DOGE-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"ICPUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"ICPUSD\\\"\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ICP-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"ICP_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"icpusdt\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"ICP-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ICP-USDT\"}]}",
          "exponent": -9,
          "id": 12,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "ICP-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"ATOMUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"ATOMUSD\\\"\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"ATOMUSDT\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ATOM-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"atomusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"ATOMUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"ATOM-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ATOM-USDT\"}]}",
          "exponent": -9,
          "id": 13,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "ATOM-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"DOTUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"DOTUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tDOTUSD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"DOT_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"dotusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"DOTUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"DOT-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"DOT-USDT\"}]}",
          "exponent": -9,
          "id": 14,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "DOT-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"XTZUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"XTZUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tXTZUSD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"XTZ-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"XTZ_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"xtzusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XTZUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"XTZ-USDT\"}]}",
          "exponent": -10,
          "id": 15,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "XTZ-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"UNIUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"UNIUSD\\\"\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"UNIUSDT\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"UNI-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"UNI_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"uniusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"UNIUSD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"UNI_USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"UNI-USDT\"}]}",
          "exponent": -9,
          "id": 16,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "UNI-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"BCHUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"BCHUSD\\\"\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"BCH-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"BCH_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"bchusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"BCHUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"BCH-USDT\"}]}",
          "exponent": -7,
          "id": 17,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "BCH-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"EOSUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"EOSUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tEOSUSD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"EOS-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"eosusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"EOSUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"EOS-USDT\"}]}",
          "exponent": -10,
          "id": 18,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "EOS-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"TRXUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"TRXUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tTRXUSD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"TRX-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"trxusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"TRXUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"TRX-USDT\"}]}",
          "exponent": -11,
          "id": 19,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "TRX-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"ALGOUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"ALGOUSD\\\"\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ALGO-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"algousdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"ALGOUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"ALGO-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ALGO-USDT\"}]}",
          "exponent": -10,
          "id": 20,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "ALGO-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"NEARUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"NEARUSD\\\"\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"NEARUSDT\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"NEAR-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"NEAR_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"nearusdt\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"NEAR-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"NEAR-USDT\"}]}",
          "exponent": -9,
          "id": 21,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "NEAR-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"SNXUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"SNXUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tSNXUSD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"SNX-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"snxusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"SNXUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"SNX-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"SNX-USDT\"}]}",
          "exponent": -9,
          "id": 22,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "SNX-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"MKRUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"MKRUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tMKRUSD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"MKR-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"MKR_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"mkrusdt\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"MKR-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"MKR-USDT\"}]}",
          "exponent": -6,
          "id": 23,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "MKR-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"SUSHIUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"SUSHIUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tSUSHI:USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"SUSHI-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"SUSHI_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"sushiusdt\"},{\"exchangeName\":\"Okx\",\"ticker\":\"SUSHI-USDT\"}]}",
          "exponent": -10,
          "id": 24,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "SUSHI-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"XLMUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"XLMUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tXLMUSD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"XLM-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"XLM_USDT\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XXLMZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"XLM-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"XLM-USDT\"}]}",
          "exponent": -10,
          "id": 25,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "XLM-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"XMRUSDT\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tXMRUSD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"XMR_USDT\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XXMRZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"XMR-USDT\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"XMR_USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"XMR-USDT\"}]}",
          "exponent": -7,
          "id": 26,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "XMR-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"ETCUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"ETCUSD\\\"\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ETC-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"ETC_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"etcusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XETCZUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ETC-USDT\"}]}",
          "exponent": -8,
          "id": 27,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "ETC-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"1INCHUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"1INCHUSD\\\"\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"1INCH-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"1INCH_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"1inchusdt\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"1INCH-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"1INCH-USDT\"}]}",
          "exponent": -10,
          "id": 28,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "1INCH-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"COMPUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"COMPUSD\\\"\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"COMPUSDT\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"COMP-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"compusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"COMPUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"COMP-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"COMP-USDT\"}]}",
          "exponent": -8,
          "id": 29,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "COMP-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"ZECUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"ZECUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tZECUSD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ZEC-USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XZECZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"ZEC-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ZEC-USDT\"}]}",
          "exponent": -8,
          "id": 30,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "ZEC-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"ZRXUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"ZRXUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tZRXUSD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ZRX-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"zrxusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"ZRXUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ZRX-USDT\"}]}",
          "exponent": -10,
          "id": 31,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "ZRX-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"YFIUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"YFIUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tYFIUSD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"YFIUSDT\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"YFI-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"yfiusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"YFIUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"YFI-USDT\"}]}",
          "exponent": -6,
          "id": 32,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "YFI-USD"
        }
      ],
      "market_prices": [
        {
          "exponent": -5,
          "id": 0,
          "price": 2000000000
        },
        {
          "exponent": -6,
          "id": 1,
          "price": 1500000000
        },
        {
          "exponent": -9,
          "id": 2,
          "price": 700000000
        },
        {
          "exponent": -10,
          "id": 3,
          "price": 7000000000
        },
        {
          "exponent": -10,
          "id": 4,
          "price": 7000000000
        },
        {
          "exponent": -8,
          "id": 5,
          "price": 1700000000
        },
        {
          "exponent": -10,
          "id": 6,
          "price": 3000000000
        },
        {
          "exponent": -8,
          "id": 7,
          "price": 1400000000
        },
        {
          "exponent": -9,
          "id": 8,
          "price": 4000000000
        },
        {
          "exponent": -8,
          "id": 9,
          "price": 7000000000
        },
        {
          "exponent": -8,
          "id": 10,
          "price": 8800000000
        },
        {
          "exponent": -11,
          "id": 11,
          "price": 7000000000
        },
        {
          "exponent": -9,
          "id": 12,
          "price": 4000000000
        },
        {
          "exponent": -9,
          "id": 13,
          "price": 10000000000
        },
        {
          "exponent": -9,
          "id": 14,
          "price": 5000000000
        },
        {
          "exponent": -10,
          "id": 15,
          "price": 8000000000
        },
        {
          "exponent": -9,
          "id": 16,
          "price": 5000000000
        },
        {
          "exponent": -7,
          "id": 17,
          "price": 2000000000
        },
        {
          "exponent": -10,
          "id": 18,
          "price": 7000000000
        },
        {
          "exponent": -11,
          "id": 19,
          "price": 7000000000
        },
        {
          "exponent": -10,
          "id": 20,
          "price": 1400000000
        },
        {
          "exponent": -9,
          "id": 21,
          "price": 1400000000
        },
        {
          "exponent": -9,
          "id": 22,
          "price": 2200000000
        },
        {
          "exponent": -6,
          "id": 23,
          "price": 7100000000
        },
        {
          "exponent": -10,
          "id": 24,
          "price": 7000000000
        },
        {
          "exponent": -10,
          "id": 25,
          "price": 10000000000
        },
        {
          "exponent": -7,
          "id": 26,
          "price": 1650000000
        },
        {
          "exponent": -8,
          "id": 27,
          "price": 1800000000
        },
        {
          "exponent": -10,
          "id": 28,
          "price": 3000000000
        },
        {
          "exponent": -8,
          "id": 29,
          "price": 4000000000
        },
        {
          "exponent": -8,
          "id": 30,
          "price": 3000000000
        },
        {
          "exponent": -10,
          "id": 31,
          "price": 2000000000
        },
        {
          "exponent": -6,
          "id": 32,
          "price": 6500000000
        }
      ]
    },
    "revshare": {
      "params": {
        "address": "dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2",
        "revenue_share_ppm": 0,
        "valid_days": 0
      }
  	},
    "rewards": {
      "params": {
        "treasury_account":"rewards_treasury",
        "denom":"adv4tnt",
        "denom_exponent":-18,
        "market_id":1,
        "fee_multiplier_ppm":990000
      }
    },
    "ratelimit": {
      "limit_params_list": [
        {
          "denom": "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
          "limiters": [
            {
              "baseline_minimum": "1000000000000",
              "baseline_tvl_ppm": 10000,
              "period": "3600s"
            },
            {
              "baseline_minimum": "10000000000000",
              "baseline_tvl_ppm": 100000,
              "period": "86400s"
            }
          ]
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
        "bond_denom": "adv4tnt",
        "historical_entries": 10000,
        "max_entries": 7,
        "max_validators": 100,
        "min_commission_rate": "0.000000000000000000",
        "unbonding_time": "1814400s"
      },
      "redelegations": [],
      "unbonding_delegations": [],
      "validators": []
    },
    "stats": {
      "params": {
        "window_duration": "2592000s"
      }
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
    "transfer": {
      "denom_traces": [],
      "params": {
        "receive_enabled": true,
        "send_enabled": true
      },
      "port_id": "transfer"
    },
    "upgrade": {},
    "vault": {
      "all_owner_share_unlocks": [],
      "default_quoting_params": {
        "layers": 2,
        "spread_min_ppm": 10000,
        "spread_buffer_ppm": 1500,
        "skew_factor_ppm": 2000000,
        "order_size_pct_ppm": 100000,
        "order_expiration_seconds": 60,
        "activation_threshold_quote_quantums": "1000000000"
      },
      "operator_params": {
        "operator": "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky",
        "metadata": {
          "name": "Governance",
          "description": "Governance Module Account"
        }
      },
      "owner_shares": [],
      "total_shares": {
        "num_shares": "0"
      },
      "vaults": []
    },
    "vest": {
      "vest_entries": [
        {
          "denom": "adv4tnt",
          "end_time": "2025-01-01T00:00:00Z",
          "start_time": "2023-01-01T00:00:00Z",
          "treasury_account": "community_treasury",
          "vester_account": "community_vester"
        },
        {
          "denom": "adv4tnt",
          "end_time": "2025-01-01T00:00:00Z",
          "start_time": "2023-01-01T00:00:00Z",
          "treasury_account": "rewards_treasury",
          "vester_account": "rewards_vester"
        }
      ]
    }
  }
}`
