import * as _5 from "./accountplus/accountplus";
import * as _6 from "./accountplus/genesis";
import * as _7 from "./affiliates/genesis";
import * as _8 from "./affiliates/query";
import * as _9 from "./affiliates/tx";
import * as _10 from "./assets/asset";
import * as _11 from "./assets/genesis";
import * as _12 from "./assets/query";
import * as _13 from "./assets/tx";
import * as _14 from "./blocktime/blocktime";
import * as _15 from "./blocktime/genesis";
import * as _16 from "./blocktime/params";
import * as _17 from "./blocktime/query";
import * as _18 from "./blocktime/tx";
import * as _19 from "./bridge/bridge_event_info";
import * as _20 from "./bridge/bridge_event";
import * as _21 from "./bridge/genesis";
import * as _22 from "./bridge/params";
import * as _23 from "./bridge/query";
import * as _24 from "./bridge/tx";
import * as _25 from "./clob/block_rate_limit_config";
import * as _26 from "./clob/clob_pair";
import * as _27 from "./clob/equity_tier_limit_config";
import * as _28 from "./clob/genesis";
import * as _29 from "./clob/liquidations_config";
import * as _30 from "./clob/liquidations";
import * as _31 from "./clob/matches";
import * as _32 from "./clob/mev";
import * as _33 from "./clob/operation";
import * as _34 from "./clob/order_removals";
import * as _35 from "./clob/order";
import * as _36 from "./clob/process_proposer_matches_events";
import * as _37 from "./clob/query";
import * as _38 from "./clob/tx";
import * as _39 from "./daemons/bridge/bridge";
import * as _40 from "./daemons/liquidation/liquidation";
import * as _41 from "./daemons/pricefeed/price_feed";
import * as _42 from "./delaymsg/block_message_ids";
import * as _43 from "./delaymsg/delayed_message";
import * as _44 from "./delaymsg/genesis";
import * as _45 from "./delaymsg/query";
import * as _46 from "./delaymsg/tx";
import * as _47 from "./epochs/epoch_info";
import * as _48 from "./epochs/genesis";
import * as _49 from "./epochs/query";
import * as _50 from "./feetiers/genesis";
import * as _51 from "./feetiers/params";
import * as _52 from "./feetiers/query";
import * as _53 from "./feetiers/tx";
import * as _54 from "./govplus/genesis";
import * as _55 from "./govplus/query";
import * as _56 from "./govplus/tx";
import * as _57 from "./indexer/events/events";
import * as _58 from "./indexer/indexer_manager/event";
import * as _59 from "./indexer/off_chain_updates/off_chain_updates";
import * as _60 from "./indexer/protocol/v1/clob";
import * as _61 from "./indexer/protocol/v1/perpetual";
import * as _62 from "./indexer/protocol/v1/subaccount";
import * as _63 from "./indexer/redis/redis_order";
import * as _64 from "./indexer/shared/removal_reason";
import * as _65 from "./indexer/socks/messages";
import * as _66 from "./listing/genesis";
import * as _67 from "./listing/params";
import * as _68 from "./listing/query";
import * as _69 from "./listing/tx";
import * as _70 from "./perpetuals/genesis";
import * as _71 from "./perpetuals/params";
import * as _72 from "./perpetuals/perpetual";
import * as _73 from "./perpetuals/query";
import * as _74 from "./perpetuals/tx";
import * as _75 from "./prices/genesis";
import * as _76 from "./prices/market_param";
import * as _77 from "./prices/market_price";
import * as _78 from "./prices/query";
import * as _79 from "./prices/tx";
import * as _80 from "./ratelimit/capacity";
import * as _81 from "./ratelimit/genesis";
import * as _82 from "./ratelimit/limit_params";
import * as _83 from "./ratelimit/pending_send_packet";
import * as _84 from "./ratelimit/query";
import * as _85 from "./ratelimit/tx";
import * as _86 from "./revshare/genesis";
import * as _87 from "./revshare/params";
import * as _88 from "./revshare/query";
import * as _89 from "./revshare/revshare";
import * as _90 from "./revshare/tx";
import * as _91 from "./rewards/genesis";
import * as _92 from "./rewards/params";
import * as _93 from "./rewards/query";
import * as _94 from "./rewards/reward_share";
import * as _95 from "./rewards/tx";
import * as _96 from "./sending/genesis";
import * as _97 from "./sending/query";
import * as _98 from "./sending/transfer";
import * as _99 from "./sending/tx";
import * as _100 from "./stats/genesis";
import * as _101 from "./stats/params";
import * as _102 from "./stats/query";
import * as _103 from "./stats/stats";
import * as _104 from "./stats/tx";
import * as _105 from "./subaccounts/asset_position";
import * as _106 from "./subaccounts/genesis";
import * as _107 from "./subaccounts/perpetual_position";
import * as _108 from "./subaccounts/query";
import * as _109 from "./subaccounts/streaming";
import * as _110 from "./subaccounts/subaccount";
import * as _111 from "./vault/genesis";
import * as _112 from "./vault/params";
import * as _113 from "./vault/query";
import * as _114 from "./vault/share";
import * as _115 from "./vault/tx";
import * as _116 from "./vault/vault";
import * as _117 from "./vest/genesis";
import * as _118 from "./vest/query";
import * as _119 from "./vest/tx";
import * as _120 from "./vest/vest_entry";
import * as _128 from "./assets/query.lcd";
import * as _129 from "./blocktime/query.lcd";
import * as _130 from "./bridge/query.lcd";
import * as _131 from "./clob/query.lcd";
import * as _132 from "./delaymsg/query.lcd";
import * as _133 from "./epochs/query.lcd";
import * as _134 from "./feetiers/query.lcd";
import * as _135 from "./listing/query.lcd";
import * as _136 from "./perpetuals/query.lcd";
import * as _137 from "./prices/query.lcd";
import * as _138 from "./ratelimit/query.lcd";
import * as _139 from "./revshare/query.lcd";
import * as _140 from "./rewards/query.lcd";
import * as _141 from "./stats/query.lcd";
import * as _142 from "./subaccounts/query.lcd";
import * as _143 from "./vault/query.lcd";
import * as _144 from "./vest/query.lcd";
import * as _145 from "./affiliates/query.rpc.Query";
import * as _146 from "./assets/query.rpc.Query";
import * as _147 from "./blocktime/query.rpc.Query";
import * as _148 from "./bridge/query.rpc.Query";
import * as _149 from "./clob/query.rpc.Query";
import * as _150 from "./delaymsg/query.rpc.Query";
import * as _151 from "./epochs/query.rpc.Query";
import * as _152 from "./feetiers/query.rpc.Query";
import * as _153 from "./govplus/query.rpc.Query";
import * as _154 from "./listing/query.rpc.Query";
import * as _155 from "./perpetuals/query.rpc.Query";
import * as _156 from "./prices/query.rpc.Query";
import * as _157 from "./ratelimit/query.rpc.Query";
import * as _158 from "./revshare/query.rpc.Query";
import * as _159 from "./rewards/query.rpc.Query";
import * as _160 from "./sending/query.rpc.Query";
import * as _161 from "./stats/query.rpc.Query";
import * as _162 from "./subaccounts/query.rpc.Query";
import * as _163 from "./vault/query.rpc.Query";
import * as _164 from "./vest/query.rpc.Query";
import * as _165 from "./blocktime/tx.rpc.msg";
import * as _166 from "./bridge/tx.rpc.msg";
import * as _167 from "./clob/tx.rpc.msg";
import * as _168 from "./delaymsg/tx.rpc.msg";
import * as _169 from "./feetiers/tx.rpc.msg";
import * as _170 from "./govplus/tx.rpc.msg";
import * as _171 from "./listing/tx.rpc.msg";
import * as _172 from "./perpetuals/tx.rpc.msg";
import * as _173 from "./prices/tx.rpc.msg";
import * as _174 from "./ratelimit/tx.rpc.msg";
import * as _175 from "./revshare/tx.rpc.msg";
import * as _176 from "./rewards/tx.rpc.msg";
import * as _177 from "./sending/tx.rpc.msg";
import * as _178 from "./stats/tx.rpc.msg";
import * as _179 from "./vault/tx.rpc.msg";
import * as _180 from "./vest/tx.rpc.msg";
import * as _181 from "./lcd";
import * as _182 from "./rpc.query";
import * as _183 from "./rpc.tx";
export namespace dydxprotocol {
  export const accountplus = { ..._5,
    ..._6
  };
  export const affiliates = { ..._7,
    ..._8,
    ..._9,
    ..._145
  };
  export const assets = { ..._10,
    ..._11,
    ..._12,
    ..._13,
    ..._128,
    ..._146
  };
  export const blocktime = { ..._14,
    ..._15,
    ..._16,
    ..._17,
    ..._18,
    ..._129,
    ..._147,
    ..._165
  };
  export const bridge = { ..._19,
    ..._20,
    ..._21,
    ..._22,
    ..._23,
    ..._24,
    ..._130,
    ..._148,
    ..._166
  };
  export const clob = { ..._25,
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
    ..._37,
    ..._38,
    ..._131,
    ..._149,
    ..._167
  };
  export namespace daemons {
    export const bridge = { ..._39
    };
    export const liquidation = { ..._40
    };
    export const pricefeed = { ..._41
    };
  }
  export const delaymsg = { ..._42,
    ..._43,
    ..._44,
    ..._45,
    ..._46,
    ..._132,
    ..._150,
    ..._168
  };
  export const epochs = { ..._47,
    ..._48,
    ..._49,
    ..._133,
    ..._151
  };
  export const feetiers = { ..._50,
    ..._51,
    ..._52,
    ..._53,
    ..._134,
    ..._152,
    ..._169
  };
  export const govplus = { ..._54,
    ..._55,
    ..._56,
    ..._153,
    ..._170
  };
  export namespace indexer {
    export const events = { ..._57
    };
    export const indexer_manager = { ..._58
    };
    export const off_chain_updates = { ..._59
    };
    export namespace protocol {
      export const v1 = { ..._60,
        ..._61,
        ..._62
      };
    }
    export const redis = { ..._63
    };
    export const shared = { ..._64
    };
    export const socks = { ..._65
    };
  }
  export const listing = { ..._66,
    ..._67,
    ..._68,
    ..._69,
    ..._135,
    ..._154,
    ..._171
  };
  export const perpetuals = { ..._70,
    ..._71,
    ..._72,
    ..._73,
    ..._74,
    ..._136,
    ..._155,
    ..._172
  };
  export const prices = { ..._75,
    ..._76,
    ..._77,
    ..._78,
    ..._79,
    ..._137,
    ..._156,
    ..._173
  };
  export const ratelimit = { ..._80,
    ..._81,
    ..._82,
    ..._83,
    ..._84,
    ..._85,
    ..._138,
    ..._157,
    ..._174
  };
  export const revshare = { ..._86,
    ..._87,
    ..._88,
    ..._89,
    ..._90,
    ..._139,
    ..._158,
    ..._175
  };
  export const rewards = { ..._91,
    ..._92,
    ..._93,
    ..._94,
    ..._95,
    ..._140,
    ..._159,
    ..._176
  };
  export const sending = { ..._96,
    ..._97,
    ..._98,
    ..._99,
    ..._160,
    ..._177
  };
  export const stats = { ..._100,
    ..._101,
    ..._102,
    ..._103,
    ..._104,
    ..._141,
    ..._161,
    ..._178
  };
  export const subaccounts = { ..._105,
    ..._106,
    ..._107,
    ..._108,
    ..._109,
    ..._110,
    ..._142,
    ..._162
  };
  export const vault = { ..._111,
    ..._112,
    ..._113,
    ..._114,
    ..._115,
    ..._116,
    ..._143,
    ..._163,
    ..._179
  };
  export const vest = { ..._117,
    ..._118,
    ..._119,
    ..._120,
    ..._144,
    ..._164,
    ..._180
  };
  export const ClientFactory = { ..._181,
    ..._182,
    ..._183
  };
}