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
import * as _30 from "./canonicalusdc/canonical_usdc";
import * as _31 from "./canonicalusdc/genesis";
import * as _32 from "./canonicalusdc/query";
import * as _33 from "./canonicalusdc/tx";
import * as _34 from "./clob/block_rate_limit_config";
import * as _35 from "./clob/clob_pair";
import * as _36 from "./clob/equity_tier_limit_config";
import * as _37 from "./clob/finalize_block";
import * as _38 from "./clob/genesis";
import * as _39 from "./clob/liquidations_config";
import * as _40 from "./clob/liquidations";
import * as _41 from "./clob/matches";
import * as _42 from "./clob/mev";
import * as _43 from "./clob/operation";
import * as _44 from "./clob/order_removals";
import * as _45 from "./clob/order";
import * as _46 from "./clob/process_proposer_matches_events";
import * as _47 from "./clob/query";
import * as _48 from "./clob/streaming";
import * as _49 from "./clob/tx";
import * as _50 from "./daemons/bridge/bridge";
import * as _51 from "./daemons/liquidation/liquidation";
import * as _52 from "./daemons/pricefeed/price_feed";
import * as _53 from "./delaymsg/block_message_ids";
import * as _54 from "./delaymsg/delayed_message";
import * as _55 from "./delaymsg/genesis";
import * as _56 from "./delaymsg/query";
import * as _57 from "./delaymsg/tx";
import * as _58 from "./epochs/epoch_info";
import * as _59 from "./epochs/genesis";
import * as _60 from "./epochs/query";
import * as _61 from "./feetiers/genesis";
import * as _62 from "./feetiers/params";
import * as _63 from "./feetiers/per_market_fee_discount";
import * as _64 from "./feetiers/query";
import * as _65 from "./feetiers/staking_tier";
import * as _66 from "./feetiers/tx";
import * as _67 from "./govplus/genesis";
import * as _68 from "./govplus/query";
import * as _69 from "./govplus/tx";
import * as _70 from "./indexer/events/events";
import * as _71 from "./indexer/indexer_manager/event";
import * as _72 from "./indexer/off_chain_updates/off_chain_updates";
import * as _73 from "./indexer/protocol/v1/clob";
import * as _74 from "./indexer/protocol/v1/perpetual";
import * as _75 from "./indexer/protocol/v1/subaccount";
import * as _76 from "./indexer/protocol/v1/vault";
import * as _77 from "./indexer/redis/redis_order";
import * as _78 from "./indexer/shared/removal_reason";
import * as _79 from "./indexer/socks/messages";
import * as _80 from "./listing/genesis";
import * as _81 from "./listing/params";
import * as _82 from "./listing/query";
import * as _83 from "./listing/tx";
import * as _84 from "./perpetuals/genesis";
import * as _85 from "./perpetuals/params";
import * as _86 from "./perpetuals/perpetual";
import * as _87 from "./perpetuals/query";
import * as _88 from "./perpetuals/tx";
import * as _89 from "./prices/genesis";
import * as _90 from "./prices/market_param";
import * as _91 from "./prices/market_price";
import * as _92 from "./prices/query";
import * as _93 from "./prices/streaming";
import * as _94 from "./prices/tx";
import * as _95 from "./ratelimit/capacity";
import * as _96 from "./ratelimit/genesis";
import * as _97 from "./ratelimit/limit_params";
import * as _98 from "./ratelimit/pending_send_packet";
import * as _99 from "./ratelimit/query";
import * as _100 from "./ratelimit/tx";
import * as _101 from "./revshare/genesis";
import * as _102 from "./revshare/params";
import * as _103 from "./revshare/query";
import * as _104 from "./revshare/revshare";
import * as _105 from "./revshare/tx";
import * as _106 from "./rewards/genesis";
import * as _107 from "./rewards/params";
import * as _108 from "./rewards/query";
import * as _109 from "./rewards/reward_share";
import * as _110 from "./rewards/tx";
import * as _111 from "./sending/genesis";
import * as _112 from "./sending/query";
import * as _113 from "./sending/transfer";
import * as _114 from "./sending/tx";
import * as _115 from "./stats/genesis";
import * as _116 from "./stats/params";
import * as _117 from "./stats/query";
import * as _118 from "./stats/stats";
import * as _119 from "./stats/tx";
import * as _120 from "./subaccounts/asset_position";
import * as _121 from "./subaccounts/genesis";
import * as _122 from "./subaccounts/leverage";
import * as _123 from "./subaccounts/perpetual_position";
import * as _124 from "./subaccounts/query";
import * as _125 from "./subaccounts/streaming";
import * as _126 from "./subaccounts/subaccount";
import * as _127 from "./vault/genesis";
import * as _128 from "./vault/params";
import * as _129 from "./vault/query";
import * as _130 from "./vault/share";
import * as _131 from "./vault/tx";
import * as _132 from "./vault/vault";
import * as _133 from "./vest/genesis";
import * as _134 from "./vest/query";
import * as _135 from "./vest/tx";
import * as _136 from "./vest/vest_entry";
import * as _144 from "./accountplus/query.lcd";
import * as _145 from "./affiliates/query.lcd";
import * as _146 from "./assets/query.lcd";
import * as _147 from "./blocktime/query.lcd";
import * as _148 from "./bridge/query.lcd";
import * as _149 from "./canonicalusdc/query.lcd";
import * as _150 from "./clob/query.lcd";
import * as _151 from "./delaymsg/query.lcd";
import * as _152 from "./epochs/query.lcd";
import * as _153 from "./feetiers/query.lcd";
import * as _154 from "./listing/query.lcd";
import * as _155 from "./perpetuals/query.lcd";
import * as _156 from "./prices/query.lcd";
import * as _157 from "./ratelimit/query.lcd";
import * as _158 from "./revshare/query.lcd";
import * as _159 from "./rewards/query.lcd";
import * as _160 from "./stats/query.lcd";
import * as _161 from "./subaccounts/query.lcd";
import * as _162 from "./vault/query.lcd";
import * as _163 from "./vest/query.lcd";
import * as _164 from "./accountplus/query.rpc.Query";
import * as _165 from "./affiliates/query.rpc.Query";
import * as _166 from "./assets/query.rpc.Query";
import * as _167 from "./blocktime/query.rpc.Query";
import * as _168 from "./bridge/query.rpc.Query";
import * as _169 from "./canonicalusdc/query.rpc.Query";
import * as _170 from "./clob/query.rpc.Query";
import * as _171 from "./delaymsg/query.rpc.Query";
import * as _172 from "./epochs/query.rpc.Query";
import * as _173 from "./feetiers/query.rpc.Query";
import * as _174 from "./govplus/query.rpc.Query";
import * as _175 from "./listing/query.rpc.Query";
import * as _176 from "./perpetuals/query.rpc.Query";
import * as _177 from "./prices/query.rpc.Query";
import * as _178 from "./ratelimit/query.rpc.Query";
import * as _179 from "./revshare/query.rpc.Query";
import * as _180 from "./rewards/query.rpc.Query";
import * as _181 from "./sending/query.rpc.Query";
import * as _182 from "./stats/query.rpc.Query";
import * as _183 from "./subaccounts/query.rpc.Query";
import * as _184 from "./vault/query.rpc.Query";
import * as _185 from "./vest/query.rpc.Query";
import * as _186 from "./accountplus/tx.rpc.msg";
import * as _187 from "./affiliates/tx.rpc.msg";
import * as _188 from "./blocktime/tx.rpc.msg";
import * as _189 from "./bridge/tx.rpc.msg";
import * as _190 from "./canonicalusdc/tx.rpc.msg";
import * as _191 from "./clob/tx.rpc.msg";
import * as _192 from "./delaymsg/tx.rpc.msg";
import * as _193 from "./feetiers/tx.rpc.msg";
import * as _194 from "./govplus/tx.rpc.msg";
import * as _195 from "./listing/tx.rpc.msg";
import * as _196 from "./perpetuals/tx.rpc.msg";
import * as _197 from "./prices/tx.rpc.msg";
import * as _198 from "./ratelimit/tx.rpc.msg";
import * as _199 from "./revshare/tx.rpc.msg";
import * as _200 from "./rewards/tx.rpc.msg";
import * as _201 from "./sending/tx.rpc.msg";
import * as _202 from "./stats/tx.rpc.msg";
import * as _203 from "./vault/tx.rpc.msg";
import * as _204 from "./vest/tx.rpc.msg";
import * as _205 from "./lcd";
import * as _206 from "./rpc.query";
import * as _207 from "./rpc.tx";
export namespace dydxprotocol {
  export const accountplus = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._9,
    ..._10,
    ..._144,
    ..._164,
    ..._186
  };
  export const affiliates = { ..._11,
    ..._12,
    ..._13,
    ..._14,
    ..._145,
    ..._165,
    ..._187
  };
  export const assets = { ..._15,
    ..._16,
    ..._17,
    ..._18,
    ..._146,
    ..._166
  };
  export const blocktime = { ..._19,
    ..._20,
    ..._21,
    ..._22,
    ..._23,
    ..._147,
    ..._167,
    ..._188
  };
  export const bridge = { ..._24,
    ..._25,
    ..._26,
    ..._27,
    ..._28,
    ..._29,
    ..._148,
    ..._168,
    ..._189
  };
  export const canonicalusdc = { ..._30,
    ..._31,
    ..._32,
    ..._33,
    ..._149,
    ..._169,
    ..._190
  };
  export const clob = { ..._34,
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
    ..._47,
    ..._48,
    ..._49,
    ..._150,
    ..._170,
    ..._191
  };
  export namespace daemons {
    export const bridge = { ..._50
    };
    export const liquidation = { ..._51
    };
    export const pricefeed = { ..._52
    };
  }
  export const delaymsg = { ..._53,
    ..._54,
    ..._55,
    ..._56,
    ..._57,
    ..._151,
    ..._171,
    ..._192
  };
  export const epochs = { ..._58,
    ..._59,
    ..._60,
    ..._152,
    ..._172
  };
  export const feetiers = { ..._61,
    ..._62,
    ..._63,
    ..._64,
    ..._65,
    ..._66,
    ..._153,
    ..._173,
    ..._193
  };
  export const govplus = { ..._67,
    ..._68,
    ..._69,
    ..._174,
    ..._194
  };
  export namespace indexer {
    export const events = { ..._70
    };
    export const indexer_manager = { ..._71
    };
    export const off_chain_updates = { ..._72
    };
    export namespace protocol {
      export const v1 = { ..._73,
        ..._74,
        ..._75,
        ..._76
      };
    }
    export const redis = { ..._77
    };
    export const shared = { ..._78
    };
    export const socks = { ..._79
    };
  }
  export const listing = { ..._80,
    ..._81,
    ..._82,
    ..._83,
    ..._154,
    ..._175,
    ..._195
  };
  export const perpetuals = { ..._84,
    ..._85,
    ..._86,
    ..._87,
    ..._88,
    ..._155,
    ..._176,
    ..._196
  };
  export const prices = { ..._89,
    ..._90,
    ..._91,
    ..._92,
    ..._93,
    ..._94,
    ..._156,
    ..._177,
    ..._197
  };
  export const ratelimit = { ..._95,
    ..._96,
    ..._97,
    ..._98,
    ..._99,
    ..._100,
    ..._157,
    ..._178,
    ..._198
  };
  export const revshare = { ..._101,
    ..._102,
    ..._103,
    ..._104,
    ..._105,
    ..._158,
    ..._179,
    ..._199
  };
  export const rewards = { ..._106,
    ..._107,
    ..._108,
    ..._109,
    ..._110,
    ..._159,
    ..._180,
    ..._200
  };
  export const sending = { ..._111,
    ..._112,
    ..._113,
    ..._114,
    ..._181,
    ..._201
  };
  export const stats = { ..._115,
    ..._116,
    ..._117,
    ..._118,
    ..._119,
    ..._160,
    ..._182,
    ..._202
  };
  export const subaccounts = { ..._120,
    ..._121,
    ..._122,
    ..._123,
    ..._124,
    ..._125,
    ..._126,
    ..._161,
    ..._183
  };
  export const vault = { ..._127,
    ..._128,
    ..._129,
    ..._130,
    ..._131,
    ..._132,
    ..._162,
    ..._184,
    ..._203
  };
  export const vest = { ..._133,
    ..._134,
    ..._135,
    ..._136,
    ..._163,
    ..._185,
    ..._204
  };
  export const ClientFactory = { ..._205,
    ..._206,
    ..._207
  };
}