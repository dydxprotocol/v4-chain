import * as _5 from "./assets/asset";
import * as _6 from "./assets/genesis";
import * as _7 from "./assets/query";
import * as _8 from "./assets/tx";
import * as _9 from "./blocktime/blocktime";
import * as _10 from "./blocktime/genesis";
import * as _11 from "./blocktime/params";
import * as _12 from "./blocktime/query";
import * as _13 from "./blocktime/tx";
import * as _14 from "./clob/block_rate_limit_config";
import * as _15 from "./clob/clob_pair";
import * as _16 from "./clob/equity_tier_limit_config";
import * as _17 from "./clob/genesis";
import * as _18 from "./clob/liquidations_config";
import * as _19 from "./clob/liquidations";
import * as _20 from "./clob/matches";
import * as _21 from "./clob/mev";
import * as _22 from "./clob/operation";
import * as _23 from "./clob/order_removals";
import * as _24 from "./clob/order";
import * as _25 from "./clob/process_proposer_matches_events";
import * as _26 from "./clob/query";
import * as _27 from "./clob/tx";
import * as _28 from "./daemons/deleveraging/deleveraging";
import * as _29 from "./daemons/pricefeed/price_feed";
import * as _30 from "./daemons/sdaioracle/sdai";
import * as _31 from "./delaymsg/block_message_ids";
import * as _32 from "./delaymsg/delayed_message";
import * as _33 from "./delaymsg/genesis";
import * as _34 from "./delaymsg/query";
import * as _35 from "./delaymsg/tx";
import * as _36 from "./epochs/epoch_info";
import * as _37 from "./epochs/genesis";
import * as _38 from "./epochs/query";
import * as _39 from "./feetiers/genesis";
import * as _40 from "./feetiers/params";
import * as _41 from "./feetiers/query";
import * as _42 from "./feetiers/tx";
import * as _43 from "./govplus/genesis";
import * as _44 from "./govplus/query";
import * as _45 from "./govplus/tx";
import * as _46 from "./indexer/events/events";
import * as _47 from "./indexer/indexer_manager/event";
import * as _48 from "./indexer/off_chain_updates/off_chain_updates";
import * as _49 from "./indexer/protocol/v1/clob";
import * as _50 from "./indexer/protocol/v1/perpetual";
import * as _51 from "./indexer/protocol/v1/subaccount";
import * as _52 from "./indexer/redis/redis_order";
import * as _53 from "./indexer/shared/removal_reason";
import * as _54 from "./indexer/socks/messages";
import * as _55 from "./perpetuals/genesis";
import * as _56 from "./perpetuals/params";
import * as _57 from "./perpetuals/perpetual";
import * as _58 from "./perpetuals/query";
import * as _59 from "./perpetuals/tx";
import * as _60 from "./prices/genesis";
import * as _61 from "./prices/market_param";
import * as _62 from "./prices/market_price";
import * as _63 from "./prices/query";
import * as _64 from "./prices/tx";
import * as _65 from "./ratelimit/capacity";
import * as _66 from "./ratelimit/genesis";
import * as _67 from "./ratelimit/limit_params";
import * as _68 from "./ratelimit/pending_send_packet";
import * as _69 from "./ratelimit/query";
import * as _70 from "./ratelimit/tx";
import * as _71 from "./rewards/genesis";
import * as _72 from "./rewards/params";
import * as _73 from "./rewards/query";
import * as _74 from "./rewards/reward_share";
import * as _75 from "./rewards/tx";
import * as _76 from "./sending/genesis";
import * as _77 from "./sending/query";
import * as _78 from "./sending/transfer";
import * as _79 from "./sending/tx";
import * as _80 from "./stats/genesis";
import * as _81 from "./stats/params";
import * as _82 from "./stats/query";
import * as _83 from "./stats/stats";
import * as _84 from "./stats/tx";
import * as _85 from "./subaccounts/asset_position";
import * as _86 from "./subaccounts/genesis";
import * as _87 from "./subaccounts/perpetual_position";
import * as _88 from "./subaccounts/query";
import * as _89 from "./subaccounts/subaccount";
import * as _90 from "./subaccounts/tx";
import * as _91 from "./ve/ve";
import * as _92 from "./vest/genesis";
import * as _93 from "./vest/query";
import * as _94 from "./vest/tx";
import * as _95 from "./vest/vest_entry";
import * as _103 from "./assets/query.lcd";
import * as _104 from "./blocktime/query.lcd";
import * as _105 from "./clob/query.lcd";
import * as _106 from "./delaymsg/query.lcd";
import * as _107 from "./epochs/query.lcd";
import * as _108 from "./feetiers/query.lcd";
import * as _109 from "./perpetuals/query.lcd";
import * as _110 from "./prices/query.lcd";
import * as _111 from "./ratelimit/query.lcd";
import * as _112 from "./rewards/query.lcd";
import * as _113 from "./stats/query.lcd";
import * as _114 from "./subaccounts/query.lcd";
import * as _115 from "./vest/query.lcd";
import * as _116 from "./assets/query.rpc.Query";
import * as _117 from "./blocktime/query.rpc.Query";
import * as _118 from "./clob/query.rpc.Query";
import * as _119 from "./delaymsg/query.rpc.Query";
import * as _120 from "./epochs/query.rpc.Query";
import * as _121 from "./feetiers/query.rpc.Query";
import * as _122 from "./govplus/query.rpc.Query";
import * as _123 from "./perpetuals/query.rpc.Query";
import * as _124 from "./prices/query.rpc.Query";
import * as _125 from "./ratelimit/query.rpc.Query";
import * as _126 from "./rewards/query.rpc.Query";
import * as _127 from "./sending/query.rpc.Query";
import * as _128 from "./stats/query.rpc.Query";
import * as _129 from "./subaccounts/query.rpc.Query";
import * as _130 from "./vest/query.rpc.Query";
import * as _131 from "./blocktime/tx.rpc.msg";
import * as _132 from "./clob/tx.rpc.msg";
import * as _133 from "./delaymsg/tx.rpc.msg";
import * as _134 from "./feetiers/tx.rpc.msg";
import * as _135 from "./govplus/tx.rpc.msg";
import * as _136 from "./perpetuals/tx.rpc.msg";
import * as _137 from "./prices/tx.rpc.msg";
import * as _138 from "./ratelimit/tx.rpc.msg";
import * as _139 from "./rewards/tx.rpc.msg";
import * as _140 from "./sending/tx.rpc.msg";
import * as _141 from "./stats/tx.rpc.msg";
import * as _142 from "./subaccounts/tx.rpc.msg";
import * as _143 from "./vest/tx.rpc.msg";
import * as _144 from "./lcd";
import * as _145 from "./rpc.query";
import * as _146 from "./rpc.tx";
export namespace dydxprotocol {
  export const assets = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._103,
    ..._116
  };
  export const blocktime = { ..._9,
    ..._10,
    ..._11,
    ..._12,
    ..._13,
    ..._104,
    ..._117,
    ..._131
  };
  export const clob = { ..._14,
    ..._15,
    ..._16,
    ..._17,
    ..._18,
    ..._19,
    ..._20,
    ..._21,
    ..._22,
    ..._23,
    ..._24,
    ..._25,
    ..._26,
    ..._27,
    ..._105,
    ..._118,
    ..._132
  };
  export namespace daemons {
    export const deleveraging = { ..._28
    };
    export const pricefeed = { ..._29
    };
    export const sdaioracle = { ..._30
    };
  }
  export const delaymsg = { ..._31,
    ..._32,
    ..._33,
    ..._34,
    ..._35,
    ..._106,
    ..._119,
    ..._133
  };
  export const epochs = { ..._36,
    ..._37,
    ..._38,
    ..._107,
    ..._120
  };
  export const feetiers = { ..._39,
    ..._40,
    ..._41,
    ..._42,
    ..._108,
    ..._121,
    ..._134
  };
  export const govplus = { ..._43,
    ..._44,
    ..._45,
    ..._122,
    ..._135
  };
  export namespace indexer {
    export const events = { ..._46
    };
    export const indexer_manager = { ..._47
    };
    export const off_chain_updates = { ..._48
    };
    export namespace protocol {
      export const v1 = { ..._49,
        ..._50,
        ..._51
      };
    }
    export const redis = { ..._52
    };
    export const shared = { ..._53
    };
    export const socks = { ..._54
    };
  }
  export const perpetuals = { ..._55,
    ..._56,
    ..._57,
    ..._58,
    ..._59,
    ..._109,
    ..._123,
    ..._136
  };
  export const prices = { ..._60,
    ..._61,
    ..._62,
    ..._63,
    ..._64,
    ..._110,
    ..._124,
    ..._137
  };
  export const ratelimit = { ..._65,
    ..._66,
    ..._67,
    ..._68,
    ..._69,
    ..._70,
    ..._111,
    ..._125,
    ..._138
  };
  export const rewards = { ..._71,
    ..._72,
    ..._73,
    ..._74,
    ..._75,
    ..._112,
    ..._126,
    ..._139
  };
  export const sending = { ..._76,
    ..._77,
    ..._78,
    ..._79,
    ..._127,
    ..._140
  };
  export const stats = { ..._80,
    ..._81,
    ..._82,
    ..._83,
    ..._84,
    ..._113,
    ..._128,
    ..._141
  };
  export const subaccounts = { ..._85,
    ..._86,
    ..._87,
    ..._88,
    ..._89,
    ..._90,
    ..._114,
    ..._129,
    ..._142
  };
  export const ve = { ..._91
  };
  export const vest = { ..._92,
    ..._93,
    ..._94,
    ..._95,
    ..._115,
    ..._130,
    ..._143
  };
  export const ClientFactory = { ..._144,
    ..._145,
    ..._146
  };
}