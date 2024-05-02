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
    "capability": {
      "index": "1",
      "owners": []
    },
    "clob": {
      "block_rate_limit_config": {
        "max_short_term_orders_per_n_blocks": [
          {
            "num_blocks": 1,
            "limit": 200
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
        ],
        "max_short_term_order_cancellations_per_n_blocks": [
          {
            "num_blocks": 1,
            "limit": 200
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
      }
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
    "perpetuals": {
      "liquidity_tiers": [
        {
          "base_position_notional": 1000000000000,
          "id": 0,
          "impact_notional": 10000000000,
          "initial_margin_ppm": 50000,
          "maintenance_fraction_ppm": 600000,
          "name": "Large-Cap"
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
            "ticker": "BTC-USD"
          }
        },
        {
          "params": {
            "atomic_resolution": -9,
            "default_funding_ppm": 0,
            "id": 1,
            "liquidity_tier": 0,
            "market_id": 1,
            "ticker": "ETH-USD"
          }
        }
      ]
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
          "exponent": -8,
          "id": 2,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "LINK-USD"
        },
        {
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"MATICUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"MATICUSD\\\"\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"MATIC-USD\"},{\"exchangeName\":\"Gate\",\"ticker\":\"MATIC_USDT\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"maticusdt\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"MATIC-USDT\"},{\"exchangeName\":\"Okx\",\"ticker\":\"MATIC-USDT\"}]}",
          "exponent": -10,
          "id": 3,
          "min_exchanges": 1,
          "min_price_change_ppm": 2000,
          "pair": "MATIC-USD"
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
          "exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"\\\"EOSUSDT\\\"\"},{\"exchangeName\":\"BinanceUS\",\"ticker\":\"\\\"EOSUSD\\\"\"},{\"exchangeName\":\"Bitfinex\",\"ticker\":\"tEOSUSD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"EOS-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"eosusdt\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"EOSUSD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"EOS-USDT\"}]}",
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
          "exponent": -7,
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
          "exponent": -11,
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
          "exponent": -8,
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
          "exponent": -7,
          "id": 23,
          "price": 7100000000
        },
        {
          "exponent": -10,
          "id": 24,
          "price": 7000000000
        },
        {
          "exponent": -11,
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
