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
import * as _33 from "./clob/genesis";
import * as _34 from "./clob/liquidations_config";
import * as _35 from "./clob/liquidations";
import * as _36 from "./clob/matches";
import * as _37 from "./clob/mev";
import * as _38 from "./clob/operation";
import * as _39 from "./clob/order_removals";
import * as _40 from "./clob/order";
import * as _41 from "./clob/process_proposer_matches_events";
import * as _42 from "./clob/query";
import * as _43 from "./clob/tx";
import * as _44 from "./daemons/bridge/bridge";
import * as _45 from "./daemons/liquidation/liquidation";
import * as _46 from "./daemons/pricefeed/price_feed";
import * as _47 from "./delaymsg/block_message_ids";
import * as _48 from "./delaymsg/delayed_message";
import * as _49 from "./delaymsg/genesis";
import * as _50 from "./delaymsg/query";
import * as _51 from "./delaymsg/tx";
import * as _52 from "./epochs/epoch_info";
import * as _53 from "./epochs/genesis";
import * as _54 from "./epochs/query";
import * as _55 from "./feetiers/genesis";
import * as _56 from "./feetiers/params";
import * as _57 from "./feetiers/query";
import * as _58 from "./feetiers/tx";
import * as _59 from "./govplus/genesis";
import * as _60 from "./govplus/query";
import * as _61 from "./govplus/tx";
import * as _62 from "./indexer/events/events";
import * as _63 from "./indexer/indexer_manager/event";
import * as _64 from "./indexer/off_chain_updates/off_chain_updates";
import * as _65 from "./indexer/protocol/v1/clob";
import * as _66 from "./indexer/protocol/v1/perpetual";
import * as _67 from "./indexer/protocol/v1/subaccount";
import * as _68 from "./indexer/redis/redis_order";
import * as _69 from "./indexer/shared/removal_reason";
import * as _70 from "./indexer/socks/messages";
import * as _71 from "./listing/genesis";
import * as _72 from "./listing/params";
import * as _73 from "./listing/query";
import * as _74 from "./listing/tx";
import * as _75 from "./perpetuals/genesis";
import * as _76 from "./perpetuals/params";
import * as _77 from "./perpetuals/perpetual";
import * as _78 from "./perpetuals/query";
import * as _79 from "./perpetuals/tx";
import * as _80 from "./prices/genesis";
import * as _81 from "./prices/market_param";
import * as _82 from "./prices/market_price";
import * as _83 from "./prices/query";
import * as _84 from "./prices/tx";
import * as _85 from "./ratelimit/capacity";
import * as _86 from "./ratelimit/genesis";
import * as _87 from "./ratelimit/limit_params";
import * as _88 from "./ratelimit/pending_send_packet";
import * as _89 from "./ratelimit/query";
import * as _90 from "./ratelimit/tx";
import * as _91 from "./revshare/genesis";
import * as _92 from "./revshare/params";
import * as _93 from "./revshare/query";
import * as _94 from "./revshare/revshare";
import * as _95 from "./revshare/tx";
import * as _96 from "./rewards/genesis";
import * as _97 from "./rewards/params";
import * as _98 from "./rewards/query";
import * as _99 from "./rewards/reward_share";
import * as _100 from "./rewards/tx";
import * as _101 from "./sending/genesis";
import * as _102 from "./sending/query";
import * as _103 from "./sending/transfer";
import * as _104 from "./sending/tx";
import * as _105 from "./stats/genesis";
import * as _106 from "./stats/params";
import * as _107 from "./stats/query";
import * as _108 from "./stats/stats";
import * as _109 from "./stats/tx";
import * as _110 from "./subaccounts/asset_position";
import * as _111 from "./subaccounts/genesis";
import * as _112 from "./subaccounts/perpetual_position";
import * as _113 from "./subaccounts/query";
import * as _114 from "./subaccounts/streaming";
import * as _115 from "./subaccounts/subaccount";
import * as _116 from "./vault/genesis";
import * as _117 from "./vault/params";
import * as _118 from "./vault/query";
import * as _119 from "./vault/share";
import * as _120 from "./vault/tx";
import * as _121 from "./vault/vault";
import * as _122 from "./vest/genesis";
import * as _123 from "./vest/query";
import * as _124 from "./vest/tx";
import * as _125 from "./vest/vest_entry";
import * as _133 from "./accountplus/query.lcd";
import * as _134 from "./assets/query.lcd";
import * as _135 from "./blocktime/query.lcd";
import * as _136 from "./bridge/query.lcd";
import * as _137 from "./clob/query.lcd";
import * as _138 from "./delaymsg/query.lcd";
import * as _139 from "./epochs/query.lcd";
import * as _140 from "./feetiers/query.lcd";
import * as _141 from "./listing/query.lcd";
import * as _142 from "./perpetuals/query.lcd";
import * as _143 from "./prices/query.lcd";
import * as _144 from "./ratelimit/query.lcd";
import * as _145 from "./revshare/query.lcd";
import * as _146 from "./rewards/query.lcd";
import * as _147 from "./stats/query.lcd";
import * as _148 from "./subaccounts/query.lcd";
import * as _149 from "./vault/query.lcd";
import * as _150 from "./vest/query.lcd";
import * as _151 from "./accountplus/query.rpc.Query";
import * as _152 from "./affiliates/query.rpc.Query";
import * as _153 from "./assets/query.rpc.Query";
import * as _154 from "./blocktime/query.rpc.Query";
import * as _155 from "./bridge/query.rpc.Query";
import * as _156 from "./clob/query.rpc.Query";
import * as _157 from "./delaymsg/query.rpc.Query";
import * as _158 from "./epochs/query.rpc.Query";
import * as _159 from "./feetiers/query.rpc.Query";
import * as _160 from "./govplus/query.rpc.Query";
import * as _161 from "./listing/query.rpc.Query";
import * as _162 from "./perpetuals/query.rpc.Query";
import * as _163 from "./prices/query.rpc.Query";
import * as _164 from "./ratelimit/query.rpc.Query";
import * as _165 from "./revshare/query.rpc.Query";
import * as _166 from "./rewards/query.rpc.Query";
import * as _167 from "./sending/query.rpc.Query";
import * as _168 from "./stats/query.rpc.Query";
import * as _169 from "./subaccounts/query.rpc.Query";
import * as _170 from "./vault/query.rpc.Query";
import * as _171 from "./vest/query.rpc.Query";
import * as _172 from "./accountplus/tx.rpc.msg";
import * as _173 from "./affiliates/tx.rpc.msg";
import * as _174 from "./blocktime/tx.rpc.msg";
import * as _175 from "./bridge/tx.rpc.msg";
import * as _176 from "./clob/tx.rpc.msg";
import * as _177 from "./delaymsg/tx.rpc.msg";
import * as _178 from "./feetiers/tx.rpc.msg";
import * as _179 from "./govplus/tx.rpc.msg";
import * as _180 from "./listing/tx.rpc.msg";
import * as _181 from "./perpetuals/tx.rpc.msg";
import * as _182 from "./prices/tx.rpc.msg";
import * as _183 from "./ratelimit/tx.rpc.msg";
import * as _184 from "./revshare/tx.rpc.msg";
import * as _185 from "./rewards/tx.rpc.msg";
import * as _186 from "./sending/tx.rpc.msg";
import * as _187 from "./stats/tx.rpc.msg";
import * as _188 from "./vault/tx.rpc.msg";
import * as _189 from "./vest/tx.rpc.msg";
import * as _190 from "./lcd";
import * as _191 from "./rpc.query";
import * as _192 from "./rpc.tx";
export namespace dydxprotocol {
  export const accountplus = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._9,
    ..._10,
    ..._133,
    ..._151,
    ..._172
  };
  export const affiliates = { ..._11,
    ..._12,
    ..._13,
    ..._14,
    ..._152,
    ..._173
  };
  export const assets = { ..._15,
    ..._16,
    ..._17,
    ..._18,
    ..._134,
    ..._153
  };
  export const blocktime = { ..._19,
    ..._20,
    ..._21,
    ..._22,
    ..._23,
    ..._135,
    ..._154,
    ..._174
  };
  export const bridge = { ..._24,
    ..._25,
    ..._26,
    ..._27,
    ..._28,
    ..._29,
    ..._136,
    ..._155,
    ..._175
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
    ..._137,
    ..._156,
    ..._176
  };
  export namespace daemons {
    export const bridge = { ..._44
    };
    export const liquidation = { ..._45
    };
    export const pricefeed = { ..._46
    };
  }
  export const delaymsg = { ..._47,
    ..._48,
    ..._49,
    ..._50,
    ..._51,
    ..._138,
    ..._157,
    ..._177
  };
  export const epochs = { ..._52,
    ..._53,
    ..._54,
    ..._139,
    ..._158
  };
  export const feetiers = { ..._55,
    ..._56,
    ..._57,
    ..._58,
    ..._140,
    ..._159,
    ..._178
  };
  export const govplus = { ..._59,
    ..._60,
    ..._61,
    ..._160,
    ..._179
  };
  export namespace indexer {
    export const events = { ..._62
    };
    export const indexer_manager = { ..._63
    };
    export const off_chain_updates = { ..._64
    };
    export namespace protocol {
      export const v1 = { ..._65,
        ..._66,
        ..._67
      };
    }
    export const redis = { ..._68
    };
    export const shared = { ..._69
    };
    export const socks = { ..._70
    };
  }
  export const listing = { ..._71,
    ..._72,
    ..._73,
    ..._74,
    ..._141,
    ..._161,
    ..._180
  };
  export const perpetuals = { ..._75,
    ..._76,
    ..._77,
    ..._78,
    ..._79,
    ..._142,
    ..._162,
    ..._181
  };
  export const prices = { ..._80,
    ..._81,
    ..._82,
    ..._83,
    ..._84,
    ..._143,
    ..._163,
    ..._182
  };
  export const ratelimit = { ..._85,
    ..._86,
    ..._87,
    ..._88,
    ..._89,
    ..._90,
    ..._144,
    ..._164,
    ..._183
  };
  export const revshare = { ..._91,
    ..._92,
    ..._93,
    ..._94,
    ..._95,
    ..._145,
    ..._165,
    ..._184
  };
  export const rewards = { ..._96,
    ..._97,
    ..._98,
    ..._99,
    ..._100,
    ..._146,
    ..._166,
    ..._185
  };
  export const sending = { ..._101,
    ..._102,
    ..._103,
    ..._104,
    ..._167,
    ..._186
  };
  export const stats = { ..._105,
    ..._106,
    ..._107,
    ..._108,
    ..._109,
    ..._147,
    ..._168,
    ..._187
  };
  export const subaccounts = { ..._110,
    ..._111,
    ..._112,
    ..._113,
    ..._114,
    ..._115,
    ..._148,
    ..._169
  };
  export const vault = { ..._116,
    ..._117,
    ..._118,
    ..._119,
    ..._120,
    ..._121,
    ..._149,
    ..._170,
    ..._188
  };
  export const vest = { ..._122,
    ..._123,
    ..._124,
    ..._125,
    ..._150,
    ..._171,
    ..._189
  };
  export const ClientFactory = { ..._190,
    ..._191,
    ..._192
  };
}