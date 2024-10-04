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
import * as _35 from "./clob/streaming";
import * as _36 from "./clob/tx";
import * as _37 from "./daemons/bridge/bridge";
import * as _38 from "./daemons/liquidation/liquidation";
import * as _39 from "./daemons/pricefeed/price_feed";
import * as _40 from "./delaymsg/block_message_ids";
import * as _41 from "./delaymsg/delayed_message";
import * as _42 from "./delaymsg/genesis";
import * as _43 from "./delaymsg/query";
import * as _44 from "./delaymsg/tx";
import * as _45 from "./epochs/epoch_info";
import * as _46 from "./epochs/genesis";
import * as _47 from "./epochs/query";
import * as _48 from "./feetiers/genesis";
import * as _49 from "./feetiers/params";
import * as _50 from "./feetiers/query";
import * as _51 from "./feetiers/tx";
import * as _52 from "./govplus/genesis";
import * as _53 from "./govplus/query";
import * as _54 from "./govplus/tx";
import * as _55 from "./indexer/events/events";
import * as _56 from "./indexer/indexer_manager/event";
import * as _57 from "./indexer/off_chain_updates/off_chain_updates";
import * as _58 from "./indexer/protocol/v1/clob";
import * as _59 from "./indexer/protocol/v1/perpetual";
import * as _60 from "./indexer/protocol/v1/subaccount";
import * as _61 from "./indexer/redis/redis_order";
import * as _62 from "./indexer/shared/removal_reason";
import * as _63 from "./indexer/socks/messages";
import * as _64 from "./listing/genesis";
import * as _65 from "./listing/query";
import * as _66 from "./listing/tx";
import * as _67 from "./perpetuals/genesis";
import * as _68 from "./perpetuals/params";
import * as _69 from "./perpetuals/perpetual";
import * as _70 from "./perpetuals/query";
import * as _71 from "./perpetuals/tx";
import * as _72 from "./prices/genesis";
import * as _73 from "./prices/market_param";
import * as _74 from "./prices/market_price";
import * as _75 from "./prices/query";
import * as _76 from "./prices/tx";
import * as _77 from "./ratelimit/capacity";
import * as _78 from "./ratelimit/genesis";
import * as _79 from "./ratelimit/limit_params";
import * as _80 from "./ratelimit/pending_send_packet";
import * as _81 from "./ratelimit/query";
import * as _82 from "./ratelimit/tx";
import * as _83 from "./revshare/genesis";
import * as _84 from "./revshare/params";
import * as _85 from "./revshare/query";
import * as _86 from "./revshare/revshare";
import * as _87 from "./revshare/tx";
import * as _88 from "./rewards/genesis";
import * as _89 from "./rewards/params";
import * as _90 from "./rewards/query";
import * as _91 from "./rewards/reward_share";
import * as _92 from "./rewards/tx";
import * as _93 from "./sending/genesis";
import * as _94 from "./sending/query";
import * as _95 from "./sending/transfer";
import * as _96 from "./sending/tx";
import * as _97 from "./stats/genesis";
import * as _98 from "./stats/params";
import * as _99 from "./stats/query";
import * as _100 from "./stats/stats";
import * as _101 from "./stats/tx";
import * as _102 from "./subaccounts/asset_position";
import * as _103 from "./subaccounts/genesis";
import * as _104 from "./subaccounts/perpetual_position";
import * as _105 from "./subaccounts/query";
import * as _106 from "./subaccounts/streaming";
import * as _107 from "./subaccounts/subaccount";
import * as _108 from "./vault/genesis";
import * as _109 from "./vault/params";
import * as _110 from "./vault/query";
import * as _111 from "./vault/share";
import * as _112 from "./vault/tx";
import * as _113 from "./vault/vault";
import * as _114 from "./vest/genesis";
import * as _115 from "./vest/query";
import * as _116 from "./vest/tx";
import * as _117 from "./vest/vest_entry";
import * as _125 from "./assets/query.lcd";
import * as _126 from "./blocktime/query.lcd";
import * as _127 from "./bridge/query.lcd";
import * as _128 from "./clob/query.lcd";
import * as _129 from "./delaymsg/query.lcd";
import * as _130 from "./epochs/query.lcd";
import * as _131 from "./feetiers/query.lcd";
import * as _132 from "./perpetuals/query.lcd";
import * as _133 from "./prices/query.lcd";
import * as _134 from "./ratelimit/query.lcd";
import * as _135 from "./revshare/query.lcd";
import * as _136 from "./rewards/query.lcd";
import * as _137 from "./stats/query.lcd";
import * as _138 from "./subaccounts/query.lcd";
import * as _139 from "./vault/query.lcd";
import * as _140 from "./vest/query.lcd";
import * as _141 from "./assets/query.rpc.Query";
import * as _142 from "./blocktime/query.rpc.Query";
import * as _143 from "./bridge/query.rpc.Query";
import * as _144 from "./clob/query.rpc.Query";
import * as _145 from "./delaymsg/query.rpc.Query";
import * as _146 from "./epochs/query.rpc.Query";
import * as _147 from "./feetiers/query.rpc.Query";
import * as _148 from "./govplus/query.rpc.Query";
import * as _149 from "./listing/query.rpc.Query";
import * as _150 from "./perpetuals/query.rpc.Query";
import * as _151 from "./prices/query.rpc.Query";
import * as _152 from "./ratelimit/query.rpc.Query";
import * as _153 from "./revshare/query.rpc.Query";
import * as _154 from "./rewards/query.rpc.Query";
import * as _155 from "./sending/query.rpc.Query";
import * as _156 from "./stats/query.rpc.Query";
import * as _157 from "./subaccounts/query.rpc.Query";
import * as _158 from "./vault/query.rpc.Query";
import * as _159 from "./vest/query.rpc.Query";
import * as _160 from "./blocktime/tx.rpc.msg";
import * as _161 from "./bridge/tx.rpc.msg";
import * as _162 from "./clob/tx.rpc.msg";
import * as _163 from "./delaymsg/tx.rpc.msg";
import * as _164 from "./feetiers/tx.rpc.msg";
import * as _165 from "./govplus/tx.rpc.msg";
import * as _166 from "./listing/tx.rpc.msg";
import * as _167 from "./perpetuals/tx.rpc.msg";
import * as _168 from "./prices/tx.rpc.msg";
import * as _169 from "./ratelimit/tx.rpc.msg";
import * as _170 from "./revshare/tx.rpc.msg";
import * as _171 from "./rewards/tx.rpc.msg";
import * as _172 from "./sending/tx.rpc.msg";
import * as _173 from "./stats/tx.rpc.msg";
import * as _174 from "./vault/tx.rpc.msg";
import * as _175 from "./vest/tx.rpc.msg";
import * as _176 from "./lcd";
import * as _177 from "./rpc.query";
import * as _178 from "./rpc.tx";
export namespace dydxprotocol {
  export const accountplus = { ..._5,
    ..._6
  };
  export const assets = { ..._7,
    ..._8,
    ..._9,
    ..._10,
    ..._125,
    ..._141
  };
  export const blocktime = { ..._11,
    ..._12,
    ..._13,
    ..._14,
    ..._15,
    ..._126,
    ..._142,
    ..._160
  };
  export const bridge = { ..._16,
    ..._17,
    ..._18,
    ..._19,
    ..._20,
    ..._21,
    ..._127,
    ..._143,
    ..._161
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
    ..._36,
    ..._128,
    ..._144,
    ..._162
  };
  export namespace daemons {
    export const bridge = { ..._37
    };
    export const liquidation = { ..._38
    };
    export const pricefeed = { ..._39
    };
  }
  export const delaymsg = { ..._40,
    ..._41,
    ..._42,
    ..._43,
    ..._44,
    ..._129,
    ..._145,
    ..._163
  };
  export const epochs = { ..._45,
    ..._46,
    ..._47,
    ..._130,
    ..._146
  };
  export const feetiers = { ..._48,
    ..._49,
    ..._50,
    ..._51,
    ..._131,
    ..._147,
    ..._164
  };
  export const govplus = { ..._52,
    ..._53,
    ..._54,
    ..._148,
    ..._165
  };
  export namespace indexer {
    export const events = { ..._55
    };
    export const indexer_manager = { ..._56
    };
    export const off_chain_updates = { ..._57
    };
    export namespace protocol {
      export const v1 = { ..._58,
        ..._59,
        ..._60
      };
    }
    export const redis = { ..._61
    };
    export const shared = { ..._62
    };
    export const socks = { ..._63
    };
  }
  export const listing = { ..._64,
    ..._65,
    ..._66,
    ..._149,
    ..._166
  };
  export const perpetuals = { ..._67,
    ..._68,
    ..._69,
    ..._70,
    ..._71,
    ..._132,
    ..._150,
    ..._167
  };
  export const prices = { ..._72,
    ..._73,
    ..._74,
    ..._75,
    ..._76,
    ..._133,
    ..._151,
    ..._168
  };
  export const ratelimit = { ..._77,
    ..._78,
    ..._79,
    ..._80,
    ..._81,
    ..._82,
    ..._134,
    ..._152,
    ..._169
  };
  export const revshare = { ..._83,
    ..._84,
    ..._85,
    ..._86,
    ..._87,
    ..._135,
    ..._153,
    ..._170
  };
  export const rewards = { ..._88,
    ..._89,
    ..._90,
    ..._91,
    ..._92,
    ..._136,
    ..._154,
    ..._171
  };
  export const sending = { ..._93,
    ..._94,
    ..._95,
    ..._96,
    ..._155,
    ..._172
  };
  export const stats = { ..._97,
    ..._98,
    ..._99,
    ..._100,
    ..._101,
    ..._137,
    ..._156,
    ..._173
  };
  export const subaccounts = { ..._102,
    ..._103,
    ..._104,
    ..._105,
    ..._106,
    ..._107,
    ..._138,
    ..._157
  };
  export const vault = { ..._108,
    ..._109,
    ..._110,
    ..._111,
    ..._112,
    ..._113,
    ..._139,
    ..._158,
    ..._174
  };
  export const vest = { ..._114,
    ..._115,
    ..._116,
    ..._117,
    ..._140,
    ..._159,
    ..._175
  };
  export const ClientFactory = { ..._176,
    ..._177,
    ..._178
  };
}