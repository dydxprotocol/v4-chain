import * as _5 from "./accountplus/accountplus";
import * as _6 from "./accountplus/genesis";
import * as _7 from "./accountplus/models";
import * as _8 from "./accountplus/query";
import * as _9 from "./accountplus/tx";
import * as _10 from "./affiliates/affiliates";
import * as _11 from "./affiliates/genesis";
import * as _12 from "./affiliates/query";
import * as _13 from "./affiliates/tx";
import * as _14 from "./assets/asset";
import * as _15 from "./assets/genesis";
import * as _16 from "./assets/query";
import * as _17 from "./assets/tx";
import * as _18 from "./blocktime/blocktime";
import * as _19 from "./blocktime/genesis";
import * as _20 from "./blocktime/params";
import * as _21 from "./blocktime/query";
import * as _22 from "./blocktime/tx";
import * as _23 from "./bridge/bridge_event_info";
import * as _24 from "./bridge/bridge_event";
import * as _25 from "./bridge/genesis";
import * as _26 from "./bridge/params";
import * as _27 from "./bridge/query";
import * as _28 from "./bridge/tx";
import * as _29 from "./clob/block_rate_limit_config";
import * as _30 from "./clob/clob_pair";
import * as _31 from "./clob/equity_tier_limit_config";
import * as _32 from "./clob/genesis";
import * as _33 from "./clob/liquidations_config";
import * as _34 from "./clob/liquidations";
import * as _35 from "./clob/matches";
import * as _36 from "./clob/mev";
import * as _37 from "./clob/operation";
import * as _38 from "./clob/order_removals";
import * as _39 from "./clob/order";
import * as _40 from "./clob/process_proposer_matches_events";
import * as _41 from "./clob/query";
import * as _42 from "./clob/tx";
import * as _43 from "./daemons/bridge/bridge";
import * as _44 from "./daemons/liquidation/liquidation";
import * as _45 from "./daemons/pricefeed/price_feed";
import * as _46 from "./delaymsg/block_message_ids";
import * as _47 from "./delaymsg/delayed_message";
import * as _48 from "./delaymsg/genesis";
import * as _49 from "./delaymsg/query";
import * as _50 from "./delaymsg/tx";
import * as _51 from "./epochs/epoch_info";
import * as _52 from "./epochs/genesis";
import * as _53 from "./epochs/query";
import * as _54 from "./feetiers/genesis";
import * as _55 from "./feetiers/params";
import * as _56 from "./feetiers/query";
import * as _57 from "./feetiers/tx";
import * as _58 from "./govplus/genesis";
import * as _59 from "./govplus/query";
import * as _60 from "./govplus/tx";
import * as _61 from "./indexer/events/events";
import * as _62 from "./indexer/indexer_manager/event";
import * as _63 from "./indexer/off_chain_updates/off_chain_updates";
import * as _64 from "./indexer/protocol/v1/clob";
import * as _65 from "./indexer/protocol/v1/perpetual";
import * as _66 from "./indexer/protocol/v1/subaccount";
import * as _67 from "./indexer/redis/redis_order";
import * as _68 from "./indexer/shared/removal_reason";
import * as _69 from "./indexer/socks/messages";
import * as _70 from "./listing/genesis";
import * as _71 from "./listing/params";
import * as _72 from "./listing/query";
import * as _73 from "./listing/tx";
import * as _74 from "./perpetuals/genesis";
import * as _75 from "./perpetuals/params";
import * as _76 from "./perpetuals/perpetual";
import * as _77 from "./perpetuals/query";
import * as _78 from "./perpetuals/tx";
import * as _79 from "./prices/genesis";
import * as _80 from "./prices/market_param";
import * as _81 from "./prices/market_price";
import * as _82 from "./prices/query";
import * as _83 from "./prices/tx";
import * as _84 from "./ratelimit/capacity";
import * as _85 from "./ratelimit/genesis";
import * as _86 from "./ratelimit/limit_params";
import * as _87 from "./ratelimit/pending_send_packet";
import * as _88 from "./ratelimit/query";
import * as _89 from "./ratelimit/tx";
import * as _90 from "./revshare/genesis";
import * as _91 from "./revshare/params";
import * as _92 from "./revshare/query";
import * as _93 from "./revshare/revshare";
import * as _94 from "./revshare/tx";
import * as _95 from "./rewards/genesis";
import * as _96 from "./rewards/params";
import * as _97 from "./rewards/query";
import * as _98 from "./rewards/reward_share";
import * as _99 from "./rewards/tx";
import * as _100 from "./sending/genesis";
import * as _101 from "./sending/query";
import * as _102 from "./sending/transfer";
import * as _103 from "./sending/tx";
import * as _104 from "./stats/genesis";
import * as _105 from "./stats/params";
import * as _106 from "./stats/query";
import * as _107 from "./stats/stats";
import * as _108 from "./stats/tx";
import * as _109 from "./subaccounts/asset_position";
import * as _110 from "./subaccounts/genesis";
import * as _111 from "./subaccounts/perpetual_position";
import * as _112 from "./subaccounts/query";
import * as _113 from "./subaccounts/streaming";
import * as _114 from "./subaccounts/subaccount";
import * as _115 from "./vault/genesis";
import * as _116 from "./vault/params";
import * as _117 from "./vault/query";
import * as _118 from "./vault/share";
import * as _119 from "./vault/tx";
import * as _120 from "./vault/vault";
import * as _121 from "./vest/genesis";
import * as _122 from "./vest/query";
import * as _123 from "./vest/tx";
import * as _124 from "./vest/vest_entry";
import * as _132 from "./accountplus/query.lcd";
import * as _133 from "./assets/query.lcd";
import * as _134 from "./blocktime/query.lcd";
import * as _135 from "./bridge/query.lcd";
import * as _136 from "./clob/query.lcd";
import * as _137 from "./delaymsg/query.lcd";
import * as _138 from "./epochs/query.lcd";
import * as _139 from "./feetiers/query.lcd";
import * as _140 from "./listing/query.lcd";
import * as _141 from "./perpetuals/query.lcd";
import * as _142 from "./prices/query.lcd";
import * as _143 from "./ratelimit/query.lcd";
import * as _144 from "./revshare/query.lcd";
import * as _145 from "./rewards/query.lcd";
import * as _146 from "./stats/query.lcd";
import * as _147 from "./subaccounts/query.lcd";
import * as _148 from "./vault/query.lcd";
import * as _149 from "./vest/query.lcd";
import * as _150 from "./accountplus/query.rpc.Query";
import * as _151 from "./affiliates/query.rpc.Query";
import * as _152 from "./assets/query.rpc.Query";
import * as _153 from "./blocktime/query.rpc.Query";
import * as _154 from "./bridge/query.rpc.Query";
import * as _155 from "./clob/query.rpc.Query";
import * as _156 from "./delaymsg/query.rpc.Query";
import * as _157 from "./epochs/query.rpc.Query";
import * as _158 from "./feetiers/query.rpc.Query";
import * as _159 from "./govplus/query.rpc.Query";
import * as _160 from "./listing/query.rpc.Query";
import * as _161 from "./perpetuals/query.rpc.Query";
import * as _162 from "./prices/query.rpc.Query";
import * as _163 from "./ratelimit/query.rpc.Query";
import * as _164 from "./revshare/query.rpc.Query";
import * as _165 from "./rewards/query.rpc.Query";
import * as _166 from "./sending/query.rpc.Query";
import * as _167 from "./stats/query.rpc.Query";
import * as _168 from "./subaccounts/query.rpc.Query";
import * as _169 from "./vault/query.rpc.Query";
import * as _170 from "./vest/query.rpc.Query";
import * as _171 from "./accountplus/tx.rpc.msg";
import * as _172 from "./affiliates/tx.rpc.msg";
import * as _173 from "./blocktime/tx.rpc.msg";
import * as _174 from "./bridge/tx.rpc.msg";
import * as _175 from "./clob/tx.rpc.msg";
import * as _176 from "./delaymsg/tx.rpc.msg";
import * as _177 from "./feetiers/tx.rpc.msg";
import * as _178 from "./govplus/tx.rpc.msg";
import * as _179 from "./listing/tx.rpc.msg";
import * as _180 from "./perpetuals/tx.rpc.msg";
import * as _181 from "./prices/tx.rpc.msg";
import * as _182 from "./ratelimit/tx.rpc.msg";
import * as _183 from "./revshare/tx.rpc.msg";
import * as _184 from "./rewards/tx.rpc.msg";
import * as _185 from "./sending/tx.rpc.msg";
import * as _186 from "./stats/tx.rpc.msg";
import * as _187 from "./vault/tx.rpc.msg";
import * as _188 from "./vest/tx.rpc.msg";
import * as _189 from "./lcd";
import * as _190 from "./rpc.query";
import * as _191 from "./rpc.tx";
export namespace dydxprotocol {
  export const accountplus = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._9,
    ..._132,
    ..._150,
    ..._171
  };
  export const affiliates = { ..._10,
    ..._11,
    ..._12,
    ..._13,
    ..._151,
    ..._172
  };
  export const assets = { ..._14,
    ..._15,
    ..._16,
    ..._17,
    ..._133,
    ..._152
  };
  export const blocktime = { ..._18,
    ..._19,
    ..._20,
    ..._21,
    ..._22,
    ..._134,
    ..._153,
    ..._173
  };
  export const bridge = { ..._23,
    ..._24,
    ..._25,
    ..._26,
    ..._27,
    ..._28,
    ..._135,
    ..._154,
    ..._174
  };
  export const clob = { ..._29,
    ..._30,
    ..._31,
    ..._32,
    ..._33,
    ..._34,
    ..._35,
    ..._36,
    ..._37,
    ..._38,
    ..._39,
    ..._40,
    ..._41,
    ..._42,
    ..._136,
    ..._155,
    ..._175
  };
  export namespace daemons {
    export const bridge = { ..._43
    };
    export const liquidation = { ..._44
    };
    export const pricefeed = { ..._45
    };
  }
  export const delaymsg = { ..._46,
    ..._47,
    ..._48,
    ..._49,
    ..._50,
    ..._137,
    ..._156,
    ..._176
  };
  export const epochs = { ..._51,
    ..._52,
    ..._53,
    ..._138,
    ..._157
  };
  export const feetiers = { ..._54,
    ..._55,
    ..._56,
    ..._57,
    ..._139,
    ..._158,
    ..._177
  };
  export const govplus = { ..._58,
    ..._59,
    ..._60,
    ..._159,
    ..._178
  };
  export namespace indexer {
    export const events = { ..._61
    };
    export const indexer_manager = { ..._62
    };
    export const off_chain_updates = { ..._63
    };
    export namespace protocol {
      export const v1 = { ..._64,
        ..._65,
        ..._66
      };
    }
    export const redis = { ..._67
    };
    export const shared = { ..._68
    };
    export const socks = { ..._69
    };
  }
  export const listing = { ..._70,
    ..._71,
    ..._72,
    ..._73,
    ..._140,
    ..._160,
    ..._179
  };
  export const perpetuals = { ..._74,
    ..._75,
    ..._76,
    ..._77,
    ..._78,
    ..._141,
    ..._161,
    ..._180
  };
  export const prices = { ..._79,
    ..._80,
    ..._81,
    ..._82,
    ..._83,
    ..._142,
    ..._162,
    ..._181
  };
  export const ratelimit = { ..._84,
    ..._85,
    ..._86,
    ..._87,
    ..._88,
    ..._89,
    ..._143,
    ..._163,
    ..._182
  };
  export const revshare = { ..._90,
    ..._91,
    ..._92,
    ..._93,
    ..._94,
    ..._144,
    ..._164,
    ..._183
  };
  export const rewards = { ..._95,
    ..._96,
    ..._97,
    ..._98,
    ..._99,
    ..._145,
    ..._165,
    ..._184
  };
  export const sending = { ..._100,
    ..._101,
    ..._102,
    ..._103,
    ..._166,
    ..._185
  };
  export const stats = { ..._104,
    ..._105,
    ..._106,
    ..._107,
    ..._108,
    ..._146,
    ..._167,
    ..._186
  };
  export const subaccounts = { ..._109,
    ..._110,
    ..._111,
    ..._112,
    ..._113,
    ..._114,
    ..._147,
    ..._168
  };
  export const vault = { ..._115,
    ..._116,
    ..._117,
    ..._118,
    ..._119,
    ..._120,
    ..._148,
    ..._169,
    ..._187
  };
  export const vest = { ..._121,
    ..._122,
    ..._123,
    ..._124,
    ..._149,
    ..._170,
    ..._188
  };
  export const ClientFactory = { ..._189,
    ..._190,
    ..._191
  };
}