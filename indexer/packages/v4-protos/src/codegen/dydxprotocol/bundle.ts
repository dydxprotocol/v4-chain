import * as _5 from "./accountplus/accountplus";
import * as _6 from "./accountplus/genesis";
import * as _7 from "./assets/asset";
import * as _8 from "./assets/genesis";
import * as _9 from "./assets/query";
import * as _10 from "./assets/tx";
import * as _11 from "./blocktime/blocktime";
import * as _12 from "./blocktime/genesis";
import * as _13 from "./blocktime/params";
import * as _14 from "./blocktime/query";
import * as _15 from "./blocktime/tx";
import * as _16 from "./bridge/bridge_event_info";
import * as _17 from "./bridge/bridge_event";
import * as _18 from "./bridge/genesis";
import * as _19 from "./bridge/params";
import * as _20 from "./bridge/query";
import * as _21 from "./bridge/tx";
import * as _22 from "./clob/block_rate_limit_config";
import * as _23 from "./clob/clob_pair";
import * as _24 from "./clob/equity_tier_limit_config";
import * as _25 from "./clob/genesis";
import * as _26 from "./clob/liquidations_config";
import * as _27 from "./clob/liquidations";
import * as _28 from "./clob/matches";
import * as _29 from "./clob/mev";
import * as _30 from "./clob/operation";
import * as _31 from "./clob/order_removals";
import * as _32 from "./clob/order";
import * as _33 from "./clob/process_proposer_matches_events";
import * as _34 from "./clob/query";
import * as _35 from "./clob/tx";
import * as _36 from "./daemons/bridge/bridge";
import * as _37 from "./daemons/liquidation/liquidation";
import * as _38 from "./daemons/pricefeed/price_feed";
import * as _39 from "./delaymsg/block_message_ids";
import * as _40 from "./delaymsg/delayed_message";
import * as _41 from "./delaymsg/genesis";
import * as _42 from "./delaymsg/query";
import * as _43 from "./delaymsg/tx";
import * as _44 from "./epochs/epoch_info";
import * as _45 from "./epochs/genesis";
import * as _46 from "./epochs/query";
import * as _47 from "./feetiers/genesis";
import * as _48 from "./feetiers/params";
import * as _49 from "./feetiers/query";
import * as _50 from "./feetiers/tx";
import * as _51 from "./govplus/genesis";
import * as _52 from "./govplus/query";
import * as _53 from "./govplus/tx";
import * as _54 from "./indexer/events/events";
import * as _55 from "./indexer/indexer_manager/event";
import * as _56 from "./indexer/off_chain_updates/off_chain_updates";
import * as _57 from "./indexer/protocol/v1/clob";
import * as _58 from "./indexer/protocol/v1/perpetual";
import * as _59 from "./indexer/protocol/v1/subaccount";
import * as _60 from "./indexer/redis/redis_order";
import * as _61 from "./indexer/shared/removal_reason";
import * as _62 from "./indexer/socks/messages";
import * as _63 from "./listing/genesis";
import * as _64 from "./listing/query";
import * as _65 from "./listing/tx";
import * as _66 from "./perpetuals/genesis";
import * as _67 from "./perpetuals/params";
import * as _68 from "./perpetuals/perpetual";
import * as _69 from "./perpetuals/query";
import * as _70 from "./perpetuals/tx";
import * as _71 from "./prices/genesis";
import * as _72 from "./prices/market_param";
import * as _73 from "./prices/market_price";
import * as _74 from "./prices/query";
import * as _75 from "./prices/tx";
import * as _76 from "./ratelimit/capacity";
import * as _77 from "./ratelimit/genesis";
import * as _78 from "./ratelimit/limit_params";
import * as _79 from "./ratelimit/pending_send_packet";
import * as _80 from "./ratelimit/query";
import * as _81 from "./ratelimit/tx";
import * as _82 from "./revshare/genesis";
import * as _83 from "./revshare/params";
import * as _84 from "./revshare/query";
import * as _85 from "./revshare/revshare";
import * as _86 from "./revshare/tx";
import * as _87 from "./rewards/genesis";
import * as _88 from "./rewards/params";
import * as _89 from "./rewards/query";
import * as _90 from "./rewards/reward_share";
import * as _91 from "./rewards/tx";
import * as _92 from "./sending/genesis";
import * as _93 from "./sending/query";
import * as _94 from "./sending/transfer";
import * as _95 from "./sending/tx";
import * as _96 from "./stats/genesis";
import * as _97 from "./stats/params";
import * as _98 from "./stats/query";
import * as _99 from "./stats/stats";
import * as _100 from "./stats/tx";
import * as _101 from "./subaccounts/asset_position";
import * as _102 from "./subaccounts/genesis";
import * as _103 from "./subaccounts/perpetual_position";
import * as _104 from "./subaccounts/query";
import * as _105 from "./subaccounts/subaccount";
import * as _106 from "./vault/genesis";
import * as _107 from "./vault/params";
import * as _108 from "./vault/query";
import * as _109 from "./vault/tx";
import * as _110 from "./vault/vault";
import * as _111 from "./vest/genesis";
import * as _112 from "./vest/query";
import * as _113 from "./vest/tx";
import * as _114 from "./vest/vest_entry";
import * as _122 from "./assets/query.lcd";
import * as _123 from "./blocktime/query.lcd";
import * as _124 from "./bridge/query.lcd";
import * as _125 from "./clob/query.lcd";
import * as _126 from "./delaymsg/query.lcd";
import * as _127 from "./epochs/query.lcd";
import * as _128 from "./feetiers/query.lcd";
import * as _129 from "./perpetuals/query.lcd";
import * as _130 from "./prices/query.lcd";
import * as _131 from "./ratelimit/query.lcd";
import * as _132 from "./revshare/query.lcd";
import * as _133 from "./rewards/query.lcd";
import * as _134 from "./stats/query.lcd";
import * as _135 from "./subaccounts/query.lcd";
import * as _136 from "./vault/query.lcd";
import * as _137 from "./vest/query.lcd";
import * as _138 from "./assets/query.rpc.Query";
import * as _139 from "./blocktime/query.rpc.Query";
import * as _140 from "./bridge/query.rpc.Query";
import * as _141 from "./clob/query.rpc.Query";
import * as _142 from "./delaymsg/query.rpc.Query";
import * as _143 from "./epochs/query.rpc.Query";
import * as _144 from "./feetiers/query.rpc.Query";
import * as _145 from "./govplus/query.rpc.Query";
import * as _146 from "./listing/query.rpc.Query";
import * as _147 from "./perpetuals/query.rpc.Query";
import * as _148 from "./prices/query.rpc.Query";
import * as _149 from "./ratelimit/query.rpc.Query";
import * as _150 from "./revshare/query.rpc.Query";
import * as _151 from "./rewards/query.rpc.Query";
import * as _152 from "./sending/query.rpc.Query";
import * as _153 from "./stats/query.rpc.Query";
import * as _154 from "./subaccounts/query.rpc.Query";
import * as _155 from "./vault/query.rpc.Query";
import * as _156 from "./vest/query.rpc.Query";
import * as _157 from "./blocktime/tx.rpc.msg";
import * as _158 from "./bridge/tx.rpc.msg";
import * as _159 from "./clob/tx.rpc.msg";
import * as _160 from "./delaymsg/tx.rpc.msg";
import * as _161 from "./feetiers/tx.rpc.msg";
import * as _162 from "./govplus/tx.rpc.msg";
import * as _163 from "./listing/tx.rpc.msg";
import * as _164 from "./perpetuals/tx.rpc.msg";
import * as _165 from "./prices/tx.rpc.msg";
import * as _166 from "./ratelimit/tx.rpc.msg";
import * as _167 from "./revshare/tx.rpc.msg";
import * as _168 from "./rewards/tx.rpc.msg";
import * as _169 from "./sending/tx.rpc.msg";
import * as _170 from "./stats/tx.rpc.msg";
import * as _171 from "./vault/tx.rpc.msg";
import * as _172 from "./vest/tx.rpc.msg";
import * as _173 from "./lcd";
import * as _174 from "./rpc.query";
import * as _175 from "./rpc.tx";
export namespace dydxprotocol {
  export const accountplus = { ..._5,
    ..._6
  };
  export const assets = { ..._7,
    ..._8,
    ..._9,
    ..._10,
    ..._122,
    ..._138
  };
  export const blocktime = { ..._11,
    ..._12,
    ..._13,
    ..._14,
    ..._15,
    ..._123,
    ..._139,
    ..._157
  };
  export const bridge = { ..._16,
    ..._17,
    ..._18,
    ..._19,
    ..._20,
    ..._21,
    ..._124,
    ..._140,
    ..._158
  };
  export const clob = { ..._22,
    ..._23,
    ..._24,
    ..._25,
    ..._26,
    ..._27,
    ..._28,
    ..._29,
    ..._30,
    ..._31,
    ..._32,
    ..._33,
    ..._34,
    ..._35,
    ..._125,
    ..._141,
    ..._159
  };
  export namespace daemons {
    export const bridge = { ..._36
    };
    export const liquidation = { ..._37
    };
    export const pricefeed = { ..._38
    };
  }
  export const delaymsg = { ..._39,
    ..._40,
    ..._41,
    ..._42,
    ..._43,
    ..._126,
    ..._142,
    ..._160
  };
  export const epochs = { ..._44,
    ..._45,
    ..._46,
    ..._127,
    ..._143
  };
  export const feetiers = { ..._47,
    ..._48,
    ..._49,
    ..._50,
    ..._128,
    ..._144,
    ..._161
  };
  export const govplus = { ..._51,
    ..._52,
    ..._53,
    ..._145,
    ..._162
  };
  export namespace indexer {
    export const events = { ..._54
    };
    export const indexer_manager = { ..._55
    };
    export const off_chain_updates = { ..._56
    };
    export namespace protocol {
      export const v1 = { ..._57,
        ..._58,
        ..._59
      };
    }
    export const redis = { ..._60
    };
    export const shared = { ..._61
    };
    export const socks = { ..._62
    };
  }
  export const listing = { ..._63,
    ..._64,
    ..._65,
    ..._146,
    ..._163
  };
  export const perpetuals = { ..._66,
    ..._67,
    ..._68,
    ..._69,
    ..._70,
    ..._129,
    ..._147,
    ..._164
  };
  export const prices = { ..._71,
    ..._72,
    ..._73,
    ..._74,
    ..._75,
    ..._130,
    ..._148,
    ..._165
  };
  export const ratelimit = { ..._76,
    ..._77,
    ..._78,
    ..._79,
    ..._80,
    ..._81,
    ..._131,
    ..._149,
    ..._166
  };
  export const revshare = { ..._82,
    ..._83,
    ..._84,
    ..._85,
    ..._86,
    ..._132,
    ..._150,
    ..._167
  };
  export const rewards = { ..._87,
    ..._88,
    ..._89,
    ..._90,
    ..._91,
    ..._133,
    ..._151,
    ..._168
  };
  export const sending = { ..._92,
    ..._93,
    ..._94,
    ..._95,
    ..._152,
    ..._169
  };
  export const stats = { ..._96,
    ..._97,
    ..._98,
    ..._99,
    ..._100,
    ..._134,
    ..._153,
    ..._170
  };
  export const subaccounts = { ..._101,
    ..._102,
    ..._103,
    ..._104,
    ..._105,
    ..._135,
    ..._154
  };
  export const vault = { ..._106,
    ..._107,
    ..._108,
    ..._109,
    ..._110,
    ..._136,
    ..._155,
    ..._171
  };
  export const vest = { ..._111,
    ..._112,
    ..._113,
    ..._114,
    ..._137,
    ..._156,
    ..._172
  };
  export const ClientFactory = { ..._173,
    ..._174,
    ..._175
  };
}