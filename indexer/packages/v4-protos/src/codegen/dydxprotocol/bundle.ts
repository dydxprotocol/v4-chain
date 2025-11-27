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
import * as _30 from "./clob/block_limits_config";
import * as _31 from "./clob/block_rate_limit_config";
import * as _32 from "./clob/clob_pair";
import * as _33 from "./clob/equity_tier_limit_config";
import * as _34 from "./clob/finalize_block";
import * as _35 from "./clob/genesis";
import * as _36 from "./clob/liquidations_config";
import * as _37 from "./clob/liquidations";
import * as _38 from "./clob/matches";
import * as _39 from "./clob/mev";
import * as _40 from "./clob/operation";
import * as _41 from "./clob/order_removals";
import * as _42 from "./clob/order";
import * as _43 from "./clob/process_proposer_matches_events";
import * as _44 from "./clob/query";
import * as _45 from "./clob/streaming";
import * as _46 from "./clob/tx";
import * as _47 from "./daemons/bridge/bridge";
import * as _48 from "./daemons/liquidation/liquidation";
import * as _49 from "./daemons/pricefeed/price_feed";
import * as _50 from "./delaymsg/block_message_ids";
import * as _51 from "./delaymsg/delayed_message";
import * as _52 from "./delaymsg/genesis";
import * as _53 from "./delaymsg/query";
import * as _54 from "./delaymsg/tx";
import * as _55 from "./epochs/epoch_info";
import * as _56 from "./epochs/genesis";
import * as _57 from "./epochs/query";
import * as _58 from "./feetiers/genesis";
import * as _59 from "./feetiers/params";
import * as _60 from "./feetiers/per_market_fee_discount";
import * as _61 from "./feetiers/query";
import * as _62 from "./feetiers/staking_tier";
import * as _63 from "./feetiers/tx";
import * as _64 from "./govplus/genesis";
import * as _65 from "./govplus/query";
import * as _66 from "./govplus/tx";
import * as _67 from "./indexer/events/events";
import * as _68 from "./indexer/indexer_manager/event";
import * as _69 from "./indexer/off_chain_updates/off_chain_updates";
import * as _70 from "./indexer/protocol/v1/clob";
import * as _71 from "./indexer/protocol/v1/perpetual";
import * as _72 from "./indexer/protocol/v1/subaccount";
import * as _73 from "./indexer/protocol/v1/vault";
import * as _74 from "./indexer/redis/redis_order";
import * as _75 from "./indexer/shared/removal_reason";
import * as _76 from "./indexer/socks/messages";
import * as _77 from "./listing/genesis";
import * as _78 from "./listing/params";
import * as _79 from "./listing/query";
import * as _80 from "./listing/tx";
import * as _81 from "./perpetuals/genesis";
import * as _82 from "./perpetuals/params";
import * as _83 from "./perpetuals/perpetual";
import * as _84 from "./perpetuals/query";
import * as _85 from "./perpetuals/tx";
import * as _86 from "./prices/genesis";
import * as _87 from "./prices/market_param";
import * as _88 from "./prices/market_price";
import * as _89 from "./prices/query";
import * as _90 from "./prices/streaming";
import * as _91 from "./prices/tx";
import * as _92 from "./ratelimit/capacity";
import * as _93 from "./ratelimit/genesis";
import * as _94 from "./ratelimit/limit_params";
import * as _95 from "./ratelimit/pending_send_packet";
import * as _96 from "./ratelimit/query";
import * as _97 from "./ratelimit/tx";
import * as _98 from "./revshare/genesis";
import * as _99 from "./revshare/params";
import * as _100 from "./revshare/query";
import * as _101 from "./revshare/revshare";
import * as _102 from "./revshare/tx";
import * as _103 from "./rewards/genesis";
import * as _104 from "./rewards/params";
import * as _105 from "./rewards/query";
import * as _106 from "./rewards/reward_share";
import * as _107 from "./rewards/tx";
import * as _108 from "./sending/genesis";
import * as _109 from "./sending/query";
import * as _110 from "./sending/transfer";
import * as _111 from "./sending/tx";
import * as _112 from "./stats/genesis";
import * as _113 from "./stats/params";
import * as _114 from "./stats/query";
import * as _115 from "./stats/stats";
import * as _116 from "./stats/tx";
import * as _117 from "./subaccounts/asset_position";
import * as _118 from "./subaccounts/genesis";
import * as _119 from "./subaccounts/leverage";
import * as _120 from "./subaccounts/perpetual_position";
import * as _121 from "./subaccounts/query";
import * as _122 from "./subaccounts/streaming";
import * as _123 from "./subaccounts/subaccount";
import * as _124 from "./vault/genesis";
import * as _125 from "./vault/params";
import * as _126 from "./vault/query";
import * as _127 from "./vault/share";
import * as _128 from "./vault/tx";
import * as _129 from "./vault/vault";
import * as _130 from "./vest/genesis";
import * as _131 from "./vest/query";
import * as _132 from "./vest/tx";
import * as _133 from "./vest/vest_entry";
import * as _141 from "./accountplus/query.lcd";
import * as _142 from "./affiliates/query.lcd";
import * as _143 from "./assets/query.lcd";
import * as _144 from "./blocktime/query.lcd";
import * as _145 from "./bridge/query.lcd";
import * as _146 from "./clob/query.lcd";
import * as _147 from "./delaymsg/query.lcd";
import * as _148 from "./epochs/query.lcd";
import * as _149 from "./feetiers/query.lcd";
import * as _150 from "./listing/query.lcd";
import * as _151 from "./perpetuals/query.lcd";
import * as _152 from "./prices/query.lcd";
import * as _153 from "./ratelimit/query.lcd";
import * as _154 from "./revshare/query.lcd";
import * as _155 from "./rewards/query.lcd";
import * as _156 from "./stats/query.lcd";
import * as _157 from "./subaccounts/query.lcd";
import * as _158 from "./vault/query.lcd";
import * as _159 from "./vest/query.lcd";
import * as _160 from "./accountplus/query.rpc.Query";
import * as _161 from "./affiliates/query.rpc.Query";
import * as _162 from "./assets/query.rpc.Query";
import * as _163 from "./blocktime/query.rpc.Query";
import * as _164 from "./bridge/query.rpc.Query";
import * as _165 from "./clob/query.rpc.Query";
import * as _166 from "./delaymsg/query.rpc.Query";
import * as _167 from "./epochs/query.rpc.Query";
import * as _168 from "./feetiers/query.rpc.Query";
import * as _169 from "./govplus/query.rpc.Query";
import * as _170 from "./listing/query.rpc.Query";
import * as _171 from "./perpetuals/query.rpc.Query";
import * as _172 from "./prices/query.rpc.Query";
import * as _173 from "./ratelimit/query.rpc.Query";
import * as _174 from "./revshare/query.rpc.Query";
import * as _175 from "./rewards/query.rpc.Query";
import * as _176 from "./sending/query.rpc.Query";
import * as _177 from "./stats/query.rpc.Query";
import * as _178 from "./subaccounts/query.rpc.Query";
import * as _179 from "./vault/query.rpc.Query";
import * as _180 from "./vest/query.rpc.Query";
import * as _181 from "./accountplus/tx.rpc.msg";
import * as _182 from "./affiliates/tx.rpc.msg";
import * as _183 from "./blocktime/tx.rpc.msg";
import * as _184 from "./bridge/tx.rpc.msg";
import * as _185 from "./clob/tx.rpc.msg";
import * as _186 from "./delaymsg/tx.rpc.msg";
import * as _187 from "./feetiers/tx.rpc.msg";
import * as _188 from "./govplus/tx.rpc.msg";
import * as _189 from "./listing/tx.rpc.msg";
import * as _190 from "./perpetuals/tx.rpc.msg";
import * as _191 from "./prices/tx.rpc.msg";
import * as _192 from "./ratelimit/tx.rpc.msg";
import * as _193 from "./revshare/tx.rpc.msg";
import * as _194 from "./rewards/tx.rpc.msg";
import * as _195 from "./sending/tx.rpc.msg";
import * as _196 from "./stats/tx.rpc.msg";
import * as _197 from "./vault/tx.rpc.msg";
import * as _198 from "./vest/tx.rpc.msg";
import * as _199 from "./lcd";
import * as _200 from "./rpc.query";
import * as _201 from "./rpc.tx";
export namespace dydxprotocol {
  export const accountplus = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._9,
    ..._10,
    ..._141,
    ..._160,
    ..._181
  };
  export const affiliates = { ..._11,
    ..._12,
    ..._13,
    ..._14,
    ..._142,
    ..._161,
    ..._182
  };
  export const assets = { ..._15,
    ..._16,
    ..._17,
    ..._18,
    ..._143,
    ..._162
  };
  export const blocktime = { ..._19,
    ..._20,
    ..._21,
    ..._22,
    ..._23,
    ..._144,
    ..._163,
    ..._183
  };
  export const bridge = { ..._24,
    ..._25,
    ..._26,
    ..._27,
    ..._28,
    ..._29,
    ..._145,
    ..._164,
    ..._184
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
    ..._46,
    ..._146,
    ..._165,
    ..._185
  };
  export namespace daemons {
    export const bridge = { ..._47
    };
    export const liquidation = { ..._48
    };
    export const pricefeed = { ..._49
    };
  }
  export const delaymsg = { ..._50,
    ..._51,
    ..._52,
    ..._53,
    ..._54,
    ..._147,
    ..._166,
    ..._186
  };
  export const epochs = { ..._55,
    ..._56,
    ..._57,
    ..._148,
    ..._167
  };
  export const feetiers = { ..._58,
    ..._59,
    ..._60,
    ..._61,
    ..._62,
    ..._63,
    ..._149,
    ..._168,
    ..._187
  };
  export const govplus = { ..._64,
    ..._65,
    ..._66,
    ..._169,
    ..._188
  };
  export namespace indexer {
    export const events = { ..._67
    };
    export const indexer_manager = { ..._68
    };
    export const off_chain_updates = { ..._69
    };
    export namespace protocol {
      export const v1 = { ..._70,
        ..._71,
        ..._72,
        ..._73
      };
    }
    export const redis = { ..._74
    };
    export const shared = { ..._75
    };
    export const socks = { ..._76
    };
  }
  export const listing = { ..._77,
    ..._78,
    ..._79,
    ..._80,
    ..._150,
    ..._170,
    ..._189
  };
  export const perpetuals = { ..._81,
    ..._82,
    ..._83,
    ..._84,
    ..._85,
    ..._151,
    ..._171,
    ..._190
  };
  export const prices = { ..._86,
    ..._87,
    ..._88,
    ..._89,
    ..._90,
    ..._91,
    ..._152,
    ..._172,
    ..._191
  };
  export const ratelimit = { ..._92,
    ..._93,
    ..._94,
    ..._95,
    ..._96,
    ..._97,
    ..._153,
    ..._173,
    ..._192
  };
  export const revshare = { ..._98,
    ..._99,
    ..._100,
    ..._101,
    ..._102,
    ..._154,
    ..._174,
    ..._193
  };
  export const rewards = { ..._103,
    ..._104,
    ..._105,
    ..._106,
    ..._107,
    ..._155,
    ..._175,
    ..._194
  };
  export const sending = { ..._108,
    ..._109,
    ..._110,
    ..._111,
    ..._176,
    ..._195
  };
  export const stats = { ..._112,
    ..._113,
    ..._114,
    ..._115,
    ..._116,
    ..._156,
    ..._177,
    ..._196
  };
  export const subaccounts = { ..._117,
    ..._118,
    ..._119,
    ..._120,
    ..._121,
    ..._122,
    ..._123,
    ..._157,
    ..._178
  };
  export const vault = { ..._124,
    ..._125,
    ..._126,
    ..._127,
    ..._128,
    ..._129,
    ..._158,
    ..._179,
    ..._197
  };
  export const vest = { ..._130,
    ..._131,
    ..._132,
    ..._133,
    ..._159,
    ..._180,
    ..._198
  };
  export const ClientFactory = { ..._199,
    ..._200,
    ..._201
  };
}