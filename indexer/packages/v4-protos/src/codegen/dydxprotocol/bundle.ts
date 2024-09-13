import * as _5 from "./accountplus/accountplus";
import * as _6 from "./accountplus/genesis";
import * as _7 from "./affiliates/affiliates";
import * as _8 from "./affiliates/genesis";
import * as _9 from "./affiliates/query";
import * as _10 from "./affiliates/tx";
import * as _11 from "./assets/asset";
import * as _12 from "./assets/genesis";
import * as _13 from "./assets/query";
import * as _14 from "./assets/tx";
import * as _15 from "./blocktime/blocktime";
import * as _16 from "./blocktime/genesis";
import * as _17 from "./blocktime/params";
import * as _18 from "./blocktime/query";
import * as _19 from "./blocktime/tx";
import * as _20 from "./bridge/bridge_event_info";
import * as _21 from "./bridge/bridge_event";
import * as _22 from "./bridge/genesis";
import * as _23 from "./bridge/params";
import * as _24 from "./bridge/query";
import * as _25 from "./bridge/tx";
import * as _26 from "./clob/block_rate_limit_config";
import * as _27 from "./clob/clob_pair";
import * as _28 from "./clob/equity_tier_limit_config";
import * as _29 from "./clob/genesis";
import * as _30 from "./clob/liquidations_config";
import * as _31 from "./clob/liquidations";
import * as _32 from "./clob/matches";
import * as _33 from "./clob/mev";
import * as _34 from "./clob/operation";
import * as _35 from "./clob/order_removals";
import * as _36 from "./clob/order";
import * as _37 from "./clob/process_proposer_matches_events";
import * as _38 from "./clob/query";
import * as _39 from "./clob/tx";
import * as _40 from "./daemons/bridge/bridge";
import * as _41 from "./daemons/liquidation/liquidation";
import * as _42 from "./daemons/pricefeed/price_feed";
import * as _43 from "./delaymsg/block_message_ids";
import * as _44 from "./delaymsg/delayed_message";
import * as _45 from "./delaymsg/genesis";
import * as _46 from "./delaymsg/query";
import * as _47 from "./delaymsg/tx";
import * as _48 from "./epochs/epoch_info";
import * as _49 from "./epochs/genesis";
import * as _50 from "./epochs/query";
import * as _51 from "./feetiers/genesis";
import * as _52 from "./feetiers/params";
import * as _53 from "./feetiers/query";
import * as _54 from "./feetiers/tx";
import * as _55 from "./govplus/genesis";
import * as _56 from "./govplus/query";
import * as _57 from "./govplus/tx";
import * as _58 from "./indexer/events/events";
import * as _59 from "./indexer/indexer_manager/event";
import * as _60 from "./indexer/off_chain_updates/off_chain_updates";
import * as _61 from "./indexer/protocol/v1/clob";
import * as _62 from "./indexer/protocol/v1/perpetual";
import * as _63 from "./indexer/protocol/v1/subaccount";
import * as _64 from "./indexer/protocol/v1/vault";
import * as _65 from "./indexer/redis/redis_order";
import * as _66 from "./indexer/shared/removal_reason";
import * as _67 from "./indexer/socks/messages";
import * as _68 from "./listing/genesis";
import * as _69 from "./listing/params";
import * as _70 from "./listing/query";
import * as _71 from "./listing/tx";
import * as _72 from "./perpetuals/genesis";
import * as _73 from "./perpetuals/params";
import * as _74 from "./perpetuals/perpetual";
import * as _75 from "./perpetuals/query";
import * as _76 from "./perpetuals/tx";
import * as _77 from "./prices/genesis";
import * as _78 from "./prices/market_param";
import * as _79 from "./prices/market_price";
import * as _80 from "./prices/query";
import * as _81 from "./prices/tx";
import * as _82 from "./ratelimit/capacity";
import * as _83 from "./ratelimit/genesis";
import * as _84 from "./ratelimit/limit_params";
import * as _85 from "./ratelimit/pending_send_packet";
import * as _86 from "./ratelimit/query";
import * as _87 from "./ratelimit/tx";
import * as _88 from "./revshare/genesis";
import * as _89 from "./revshare/params";
import * as _90 from "./revshare/query";
import * as _91 from "./revshare/revshare";
import * as _92 from "./revshare/tx";
import * as _93 from "./rewards/genesis";
import * as _94 from "./rewards/params";
import * as _95 from "./rewards/query";
import * as _96 from "./rewards/reward_share";
import * as _97 from "./rewards/tx";
import * as _98 from "./sending/genesis";
import * as _99 from "./sending/query";
import * as _100 from "./sending/transfer";
import * as _101 from "./sending/tx";
import * as _102 from "./stats/genesis";
import * as _103 from "./stats/params";
import * as _104 from "./stats/query";
import * as _105 from "./stats/stats";
import * as _106 from "./stats/tx";
import * as _107 from "./subaccounts/asset_position";
import * as _108 from "./subaccounts/genesis";
import * as _109 from "./subaccounts/perpetual_position";
import * as _110 from "./subaccounts/query";
import * as _111 from "./subaccounts/streaming";
import * as _112 from "./subaccounts/subaccount";
import * as _113 from "./vault/genesis";
import * as _114 from "./vault/params";
import * as _115 from "./vault/query";
import * as _116 from "./vault/share";
import * as _117 from "./vault/tx";
import * as _118 from "./vault/vault";
import * as _119 from "./vest/genesis";
import * as _120 from "./vest/query";
import * as _121 from "./vest/tx";
import * as _122 from "./vest/vest_entry";
import * as _130 from "./assets/query.lcd";
import * as _131 from "./blocktime/query.lcd";
import * as _132 from "./bridge/query.lcd";
import * as _133 from "./clob/query.lcd";
import * as _134 from "./delaymsg/query.lcd";
import * as _135 from "./epochs/query.lcd";
import * as _136 from "./feetiers/query.lcd";
import * as _137 from "./listing/query.lcd";
import * as _138 from "./perpetuals/query.lcd";
import * as _139 from "./prices/query.lcd";
import * as _140 from "./ratelimit/query.lcd";
import * as _141 from "./revshare/query.lcd";
import * as _142 from "./rewards/query.lcd";
import * as _143 from "./stats/query.lcd";
import * as _144 from "./subaccounts/query.lcd";
import * as _145 from "./vault/query.lcd";
import * as _146 from "./vest/query.lcd";
import * as _147 from "./affiliates/query.rpc.Query";
import * as _148 from "./assets/query.rpc.Query";
import * as _149 from "./blocktime/query.rpc.Query";
import * as _150 from "./bridge/query.rpc.Query";
import * as _151 from "./clob/query.rpc.Query";
import * as _152 from "./delaymsg/query.rpc.Query";
import * as _153 from "./epochs/query.rpc.Query";
import * as _154 from "./feetiers/query.rpc.Query";
import * as _155 from "./govplus/query.rpc.Query";
import * as _156 from "./listing/query.rpc.Query";
import * as _157 from "./perpetuals/query.rpc.Query";
import * as _158 from "./prices/query.rpc.Query";
import * as _159 from "./ratelimit/query.rpc.Query";
import * as _160 from "./revshare/query.rpc.Query";
import * as _161 from "./rewards/query.rpc.Query";
import * as _162 from "./sending/query.rpc.Query";
import * as _163 from "./stats/query.rpc.Query";
import * as _164 from "./subaccounts/query.rpc.Query";
import * as _165 from "./vault/query.rpc.Query";
import * as _166 from "./vest/query.rpc.Query";
import * as _167 from "./affiliates/tx.rpc.msg";
import * as _168 from "./blocktime/tx.rpc.msg";
import * as _169 from "./bridge/tx.rpc.msg";
import * as _170 from "./clob/tx.rpc.msg";
import * as _171 from "./delaymsg/tx.rpc.msg";
import * as _172 from "./feetiers/tx.rpc.msg";
import * as _173 from "./govplus/tx.rpc.msg";
import * as _174 from "./listing/tx.rpc.msg";
import * as _175 from "./perpetuals/tx.rpc.msg";
import * as _176 from "./prices/tx.rpc.msg";
import * as _177 from "./ratelimit/tx.rpc.msg";
import * as _178 from "./revshare/tx.rpc.msg";
import * as _179 from "./rewards/tx.rpc.msg";
import * as _180 from "./sending/tx.rpc.msg";
import * as _181 from "./stats/tx.rpc.msg";
import * as _182 from "./vault/tx.rpc.msg";
import * as _183 from "./vest/tx.rpc.msg";
import * as _184 from "./lcd";
import * as _185 from "./rpc.query";
import * as _186 from "./rpc.tx";
export namespace dydxprotocol {
  export const accountplus = { ..._5,
    ..._6
  };
  export const affiliates = { ..._7,
    ..._8,
    ..._9,
    ..._10,
    ..._147,
    ..._167
  };
  export const assets = { ..._11,
    ..._12,
    ..._13,
    ..._14,
    ..._130,
    ..._148
  };
  export const blocktime = { ..._15,
    ..._16,
    ..._17,
    ..._18,
    ..._19,
    ..._131,
    ..._149,
    ..._168
  };
  export const bridge = { ..._20,
    ..._21,
    ..._22,
    ..._23,
    ..._24,
    ..._25,
    ..._132,
    ..._150,
    ..._169
  };
  export const clob = { ..._26,
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
    ..._39,
    ..._133,
    ..._151,
    ..._170
  };
  export namespace daemons {
    export const bridge = { ..._40
    };
    export const liquidation = { ..._41
    };
    export const pricefeed = { ..._42
    };
  }
  export const delaymsg = { ..._43,
    ..._44,
    ..._45,
    ..._46,
    ..._47,
    ..._134,
    ..._152,
    ..._171
  };
  export const epochs = { ..._48,
    ..._49,
    ..._50,
    ..._135,
    ..._153
  };
  export const feetiers = { ..._51,
    ..._52,
    ..._53,
    ..._54,
    ..._136,
    ..._154,
    ..._172
  };
  export const govplus = { ..._55,
    ..._56,
    ..._57,
    ..._155,
    ..._173
  };
  export namespace indexer {
    export const events = { ..._58
    };
    export const indexer_manager = { ..._59
    };
    export const off_chain_updates = { ..._60
    };
    export namespace protocol {
      export const v1 = { ..._61,
        ..._62,
        ..._63,
        ..._64
      };
    }
    export const redis = { ..._65
    };
    export const shared = { ..._66
    };
    export const socks = { ..._67
    };
  }
  export const listing = { ..._68,
    ..._69,
    ..._70,
    ..._71,
    ..._137,
    ..._156,
    ..._174
  };
  export const perpetuals = { ..._72,
    ..._73,
    ..._74,
    ..._75,
    ..._76,
    ..._138,
    ..._157,
    ..._175
  };
  export const prices = { ..._77,
    ..._78,
    ..._79,
    ..._80,
    ..._81,
    ..._139,
    ..._158,
    ..._176
  };
  export const ratelimit = { ..._82,
    ..._83,
    ..._84,
    ..._85,
    ..._86,
    ..._87,
    ..._140,
    ..._159,
    ..._177
  };
  export const revshare = { ..._88,
    ..._89,
    ..._90,
    ..._91,
    ..._92,
    ..._141,
    ..._160,
    ..._178
  };
  export const rewards = { ..._93,
    ..._94,
    ..._95,
    ..._96,
    ..._97,
    ..._142,
    ..._161,
    ..._179
  };
  export const sending = { ..._98,
    ..._99,
    ..._100,
    ..._101,
    ..._162,
    ..._180
  };
  export const stats = { ..._102,
    ..._103,
    ..._104,
    ..._105,
    ..._106,
    ..._143,
    ..._163,
    ..._181
  };
  export const subaccounts = { ..._107,
    ..._108,
    ..._109,
    ..._110,
    ..._111,
    ..._112,
    ..._144,
    ..._164
  };
  export const vault = { ..._113,
    ..._114,
    ..._115,
    ..._116,
    ..._117,
    ..._118,
    ..._145,
    ..._165,
    ..._182
  };
  export const vest = { ..._119,
    ..._120,
    ..._121,
    ..._122,
    ..._146,
    ..._166,
    ..._183
  };
  export const ClientFactory = { ..._184,
    ..._185,
    ..._186
  };
}