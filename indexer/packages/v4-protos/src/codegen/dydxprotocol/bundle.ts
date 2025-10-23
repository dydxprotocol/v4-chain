import * as _5 from "./accountplus/accountplus";
import * as _6 from "./accountplus/genesis";
import * as _7 from "./accountplus/models";
import * as _8 from "./accountplus/params";
import * as _9 from "./accountplus/query";
import * as _10 from "./accountplus/tx";
import * as _11 from "./affiliates/affiliates";
import * as _12 from "./affiliates/genesis";
import * as _13 from "./affiliates/query";
import * as _14 from "./affiliates/tx";
import * as _15 from "./assets/asset";
import * as _16 from "./assets/genesis";
import * as _17 from "./assets/query";
import * as _18 from "./assets/tx";
import * as _19 from "./blocktime/blocktime";
import * as _20 from "./blocktime/genesis";
import * as _21 from "./blocktime/params";
import * as _22 from "./blocktime/query";
import * as _23 from "./blocktime/tx";
import * as _24 from "./bridge/bridge_event_info";
import * as _25 from "./bridge/bridge_event";
import * as _26 from "./bridge/genesis";
import * as _27 from "./bridge/params";
import * as _28 from "./bridge/query";
import * as _29 from "./bridge/tx";
import * as _30 from "./clob/block_rate_limit_config";
import * as _31 from "./clob/clob_pair";
import * as _32 from "./clob/equity_tier_limit_config";
import * as _33 from "./clob/finalize_block";
import * as _34 from "./clob/genesis";
import * as _35 from "./clob/liquidations_config";
import * as _36 from "./clob/liquidations";
import * as _37 from "./clob/matches";
import * as _38 from "./clob/mev";
import * as _39 from "./clob/operation";
import * as _40 from "./clob/order_removals";
import * as _41 from "./clob/order";
import * as _42 from "./clob/process_proposer_matches_events";
import * as _43 from "./clob/query";
import * as _44 from "./clob/streaming";
import * as _45 from "./clob/tx";
import * as _46 from "./daemons/bridge/bridge";
import * as _47 from "./daemons/liquidation/liquidation";
import * as _48 from "./daemons/pricefeed/price_feed";
import * as _49 from "./delaymsg/block_message_ids";
import * as _50 from "./delaymsg/delayed_message";
import * as _51 from "./delaymsg/genesis";
import * as _52 from "./delaymsg/query";
import * as _53 from "./delaymsg/tx";
import * as _54 from "./epochs/epoch_info";
import * as _55 from "./epochs/genesis";
import * as _56 from "./epochs/query";
import * as _57 from "./feetiers/genesis";
import * as _58 from "./feetiers/params";
import * as _59 from "./feetiers/per_market_fee_discount";
import * as _60 from "./feetiers/query";
import * as _61 from "./feetiers/tx";
import * as _62 from "./govplus/genesis";
import * as _63 from "./govplus/query";
import * as _64 from "./govplus/tx";
import * as _65 from "./indexer/events/events";
import * as _66 from "./indexer/indexer_manager/event";
import * as _67 from "./indexer/off_chain_updates/off_chain_updates";
import * as _68 from "./indexer/protocol/v1/clob";
import * as _69 from "./indexer/protocol/v1/perpetual";
import * as _70 from "./indexer/protocol/v1/subaccount";
import * as _71 from "./indexer/protocol/v1/vault";
import * as _72 from "./indexer/redis/redis_order";
import * as _73 from "./indexer/shared/removal_reason";
import * as _74 from "./indexer/socks/messages";
import * as _75 from "./listing/genesis";
import * as _76 from "./listing/params";
import * as _77 from "./listing/query";
import * as _78 from "./listing/tx";
import * as _79 from "./perpetuals/genesis";
import * as _80 from "./perpetuals/params";
import * as _81 from "./perpetuals/perpetual";
import * as _82 from "./perpetuals/query";
import * as _83 from "./perpetuals/tx";
import * as _84 from "./prices/genesis";
import * as _85 from "./prices/market_param";
import * as _86 from "./prices/market_price";
import * as _87 from "./prices/query";
import * as _88 from "./prices/streaming";
import * as _89 from "./prices/tx";
import * as _90 from "./ratelimit/capacity";
import * as _91 from "./ratelimit/genesis";
import * as _92 from "./ratelimit/limit_params";
import * as _93 from "./ratelimit/pending_send_packet";
import * as _94 from "./ratelimit/query";
import * as _95 from "./ratelimit/tx";
import * as _96 from "./revshare/genesis";
import * as _97 from "./revshare/params";
import * as _98 from "./revshare/query";
import * as _99 from "./revshare/revshare";
import * as _100 from "./revshare/tx";
import * as _101 from "./rewards/genesis";
import * as _102 from "./rewards/params";
import * as _103 from "./rewards/query";
import * as _104 from "./rewards/reward_share";
import * as _105 from "./rewards/tx";
import * as _106 from "./sending/genesis";
import * as _107 from "./sending/query";
import * as _108 from "./sending/transfer";
import * as _109 from "./sending/tx";
import * as _110 from "./stats/genesis";
import * as _111 from "./stats/params";
import * as _112 from "./stats/query";
import * as _113 from "./stats/stats";
import * as _114 from "./stats/tx";
import * as _115 from "./subaccounts/asset_position";
import * as _116 from "./subaccounts/genesis";
import * as _117 from "./subaccounts/perpetual_position";
import * as _118 from "./subaccounts/query";
import * as _119 from "./subaccounts/streaming";
import * as _120 from "./subaccounts/subaccount";
import * as _121 from "./vault/genesis";
import * as _122 from "./vault/params";
import * as _123 from "./vault/query";
import * as _124 from "./vault/share";
import * as _125 from "./vault/tx";
import * as _126 from "./vault/vault";
import * as _127 from "./vest/genesis";
import * as _128 from "./vest/query";
import * as _129 from "./vest/tx";
import * as _130 from "./vest/vest_entry";
import * as _138 from "./accountplus/query.lcd";
import * as _139 from "./affiliates/query.lcd";
import * as _140 from "./assets/query.lcd";
import * as _141 from "./blocktime/query.lcd";
import * as _142 from "./bridge/query.lcd";
import * as _143 from "./clob/query.lcd";
import * as _144 from "./delaymsg/query.lcd";
import * as _145 from "./epochs/query.lcd";
import * as _146 from "./feetiers/query.lcd";
import * as _147 from "./listing/query.lcd";
import * as _148 from "./perpetuals/query.lcd";
import * as _149 from "./prices/query.lcd";
import * as _150 from "./ratelimit/query.lcd";
import * as _151 from "./revshare/query.lcd";
import * as _152 from "./rewards/query.lcd";
import * as _153 from "./stats/query.lcd";
import * as _154 from "./subaccounts/query.lcd";
import * as _155 from "./vault/query.lcd";
import * as _156 from "./vest/query.lcd";
import * as _157 from "./accountplus/query.rpc.Query";
import * as _158 from "./affiliates/query.rpc.Query";
import * as _159 from "./assets/query.rpc.Query";
import * as _160 from "./blocktime/query.rpc.Query";
import * as _161 from "./bridge/query.rpc.Query";
import * as _162 from "./clob/query.rpc.Query";
import * as _163 from "./delaymsg/query.rpc.Query";
import * as _164 from "./epochs/query.rpc.Query";
import * as _165 from "./feetiers/query.rpc.Query";
import * as _166 from "./govplus/query.rpc.Query";
import * as _167 from "./listing/query.rpc.Query";
import * as _168 from "./perpetuals/query.rpc.Query";
import * as _169 from "./prices/query.rpc.Query";
import * as _170 from "./ratelimit/query.rpc.Query";
import * as _171 from "./revshare/query.rpc.Query";
import * as _172 from "./rewards/query.rpc.Query";
import * as _173 from "./sending/query.rpc.Query";
import * as _174 from "./stats/query.rpc.Query";
import * as _175 from "./subaccounts/query.rpc.Query";
import * as _176 from "./vault/query.rpc.Query";
import * as _177 from "./vest/query.rpc.Query";
import * as _178 from "./accountplus/tx.rpc.msg";
import * as _179 from "./affiliates/tx.rpc.msg";
import * as _180 from "./blocktime/tx.rpc.msg";
import * as _181 from "./bridge/tx.rpc.msg";
import * as _182 from "./clob/tx.rpc.msg";
import * as _183 from "./delaymsg/tx.rpc.msg";
import * as _184 from "./feetiers/tx.rpc.msg";
import * as _185 from "./govplus/tx.rpc.msg";
import * as _186 from "./listing/tx.rpc.msg";
import * as _187 from "./perpetuals/tx.rpc.msg";
import * as _188 from "./prices/tx.rpc.msg";
import * as _189 from "./ratelimit/tx.rpc.msg";
import * as _190 from "./revshare/tx.rpc.msg";
import * as _191 from "./rewards/tx.rpc.msg";
import * as _192 from "./sending/tx.rpc.msg";
import * as _193 from "./stats/tx.rpc.msg";
import * as _194 from "./vault/tx.rpc.msg";
import * as _195 from "./vest/tx.rpc.msg";
import * as _196 from "./lcd";
import * as _197 from "./rpc.query";
import * as _198 from "./rpc.tx";
export namespace dydxprotocol {
  export const accountplus = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._9,
    ..._10,
    ..._138,
    ..._157,
    ..._178
  };
  export const affiliates = { ..._11,
    ..._12,
    ..._13,
    ..._14,
    ..._139,
    ..._158,
    ..._179
  };
  export const assets = { ..._15,
    ..._16,
    ..._17,
    ..._18,
    ..._140,
    ..._159
  };
  export const blocktime = { ..._19,
    ..._20,
    ..._21,
    ..._22,
    ..._23,
    ..._141,
    ..._160,
    ..._180
  };
  export const bridge = { ..._24,
    ..._25,
    ..._26,
    ..._27,
    ..._28,
    ..._29,
    ..._142,
    ..._161,
    ..._181
  };
  export const clob = { ..._30,
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
    ..._43,
    ..._44,
    ..._45,
    ..._143,
    ..._162,
    ..._182
  };
  export namespace daemons {
    export const bridge = { ..._46
    };
    export const liquidation = { ..._47
    };
    export const pricefeed = { ..._48
    };
  }
  export const delaymsg = { ..._49,
    ..._50,
    ..._51,
    ..._52,
    ..._53,
    ..._144,
    ..._163,
    ..._183
  };
  export const epochs = { ..._54,
    ..._55,
    ..._56,
    ..._145,
    ..._164
  };
  export const feetiers = { ..._57,
    ..._58,
    ..._59,
    ..._60,
    ..._61,
    ..._146,
    ..._165,
    ..._184
  };
  export const govplus = { ..._62,
    ..._63,
    ..._64,
    ..._166,
    ..._185
  };
  export namespace indexer {
    export const events = { ..._65
    };
    export const indexer_manager = { ..._66
    };
    export const off_chain_updates = { ..._67
    };
    export namespace protocol {
      export const v1 = { ..._68,
        ..._69,
        ..._70,
        ..._71
      };
    }
    export const redis = { ..._72
    };
    export const shared = { ..._73
    };
    export const socks = { ..._74
    };
  }
  export const listing = { ..._75,
    ..._76,
    ..._77,
    ..._78,
    ..._147,
    ..._167,
    ..._186
  };
  export const perpetuals = { ..._79,
    ..._80,
    ..._81,
    ..._82,
    ..._83,
    ..._148,
    ..._168,
    ..._187
  };
  export const prices = { ..._84,
    ..._85,
    ..._86,
    ..._87,
    ..._88,
    ..._89,
    ..._149,
    ..._169,
    ..._188
  };
  export const ratelimit = { ..._90,
    ..._91,
    ..._92,
    ..._93,
    ..._94,
    ..._95,
    ..._150,
    ..._170,
    ..._189
  };
  export const revshare = { ..._96,
    ..._97,
    ..._98,
    ..._99,
    ..._100,
    ..._151,
    ..._171,
    ..._190
  };
  export const rewards = { ..._101,
    ..._102,
    ..._103,
    ..._104,
    ..._105,
    ..._152,
    ..._172,
    ..._191
  };
  export const sending = { ..._106,
    ..._107,
    ..._108,
    ..._109,
    ..._173,
    ..._192
  };
  export const stats = { ..._110,
    ..._111,
    ..._112,
    ..._113,
    ..._114,
    ..._153,
    ..._174,
    ..._193
  };
  export const subaccounts = { ..._115,
    ..._116,
    ..._117,
    ..._118,
    ..._119,
    ..._120,
    ..._154,
    ..._175
  };
  export const vault = { ..._121,
    ..._122,
    ..._123,
    ..._124,
    ..._125,
    ..._126,
    ..._155,
    ..._176,
    ..._194
  };
  export const vest = { ..._127,
    ..._128,
    ..._129,
    ..._130,
    ..._156,
    ..._177,
    ..._195
  };
  export const ClientFactory = { ..._196,
    ..._197,
    ..._198
  };
}