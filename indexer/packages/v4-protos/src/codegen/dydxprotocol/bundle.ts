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
import * as _30 from "./delaymsg/block_message_ids";
import * as _31 from "./delaymsg/delayed_message";
import * as _32 from "./delaymsg/genesis";
import * as _33 from "./delaymsg/query";
import * as _34 from "./delaymsg/tx";
import * as _35 from "./epochs/epoch_info";
import * as _36 from "./epochs/genesis";
import * as _37 from "./epochs/query";
import * as _38 from "./feetiers/genesis";
import * as _39 from "./feetiers/params";
import * as _40 from "./feetiers/query";
import * as _41 from "./feetiers/tx";
import * as _42 from "./indexer/events/events";
import * as _43 from "./indexer/indexer_manager/event";
import * as _44 from "./indexer/off_chain_updates/off_chain_updates";
import * as _45 from "./indexer/protocol/v1/clob";
import * as _46 from "./indexer/protocol/v1/perpetual";
import * as _47 from "./indexer/protocol/v1/subaccount";
import * as _48 from "./indexer/redis/redis_order";
import * as _49 from "./indexer/shared/removal_reason";
import * as _50 from "./indexer/socks/messages";
import * as _51 from "./perpetuals/genesis";
import * as _52 from "./perpetuals/params";
import * as _53 from "./perpetuals/perpetual";
import * as _54 from "./perpetuals/query";
import * as _55 from "./perpetuals/tx";
import * as _56 from "./prices/genesis";
import * as _57 from "./prices/market_param";
import * as _58 from "./prices/market_price";
import * as _59 from "./prices/query";
import * as _60 from "./prices/tx";
import * as _61 from "./ratelimit/capacity";
import * as _62 from "./ratelimit/genesis";
import * as _63 from "./ratelimit/limit_params";
import * as _64 from "./ratelimit/pending_send_packet";
import * as _65 from "./ratelimit/query";
import * as _66 from "./ratelimit/tx";
import * as _67 from "./sending/genesis";
import * as _68 from "./sending/query";
import * as _69 from "./sending/transfer";
import * as _70 from "./sending/tx";
import * as _71 from "./stats/genesis";
import * as _72 from "./stats/params";
import * as _73 from "./stats/query";
import * as _74 from "./stats/stats";
import * as _75 from "./stats/tx";
import * as _76 from "./subaccounts/asset_position";
import * as _77 from "./subaccounts/genesis";
import * as _78 from "./subaccounts/perpetual_position";
import * as _79 from "./subaccounts/query";
import * as _80 from "./subaccounts/subaccount";
import * as _81 from "./ve/ve";
import * as _89 from "./assets/query.lcd";
import * as _90 from "./blocktime/query.lcd";
import * as _91 from "./clob/query.lcd";
import * as _92 from "./delaymsg/query.lcd";
import * as _93 from "./epochs/query.lcd";
import * as _94 from "./feetiers/query.lcd";
import * as _95 from "./perpetuals/query.lcd";
import * as _96 from "./prices/query.lcd";
import * as _97 from "./ratelimit/query.lcd";
import * as _98 from "./stats/query.lcd";
import * as _99 from "./subaccounts/query.lcd";
import * as _100 from "./assets/query.rpc.Query";
import * as _101 from "./blocktime/query.rpc.Query";
import * as _102 from "./clob/query.rpc.Query";
import * as _103 from "./delaymsg/query.rpc.Query";
import * as _104 from "./epochs/query.rpc.Query";
import * as _105 from "./feetiers/query.rpc.Query";
import * as _106 from "./perpetuals/query.rpc.Query";
import * as _107 from "./prices/query.rpc.Query";
import * as _108 from "./ratelimit/query.rpc.Query";
import * as _109 from "./sending/query.rpc.Query";
import * as _110 from "./stats/query.rpc.Query";
import * as _111 from "./subaccounts/query.rpc.Query";
import * as _112 from "./blocktime/tx.rpc.msg";
import * as _113 from "./clob/tx.rpc.msg";
import * as _114 from "./delaymsg/tx.rpc.msg";
import * as _115 from "./feetiers/tx.rpc.msg";
import * as _116 from "./perpetuals/tx.rpc.msg";
import * as _117 from "./prices/tx.rpc.msg";
import * as _118 from "./ratelimit/tx.rpc.msg";
import * as _119 from "./sending/tx.rpc.msg";
import * as _120 from "./stats/tx.rpc.msg";
import * as _121 from "./lcd";
import * as _122 from "./rpc.query";
import * as _123 from "./rpc.tx";
export namespace dydxprotocol {
  export const assets = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._89,
    ..._100
  };
  export const blocktime = { ..._9,
    ..._10,
    ..._11,
    ..._12,
    ..._13,
    ..._90,
    ..._101,
    ..._112
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
    ..._91,
    ..._102,
    ..._113
  };
  export namespace daemons {
    export const deleveraging = { ..._28
    };
    export const pricefeed = { ..._29
    };
  }
  export const delaymsg = { ..._30,
    ..._31,
    ..._32,
    ..._33,
    ..._34,
    ..._92,
    ..._103,
    ..._114
  };
  export const epochs = { ..._35,
    ..._36,
    ..._37,
    ..._93,
    ..._104
  };
  export const feetiers = { ..._38,
    ..._39,
    ..._40,
    ..._41,
    ..._94,
    ..._105,
    ..._115
  };
  export namespace indexer {
    export const events = { ..._42
    };
    export const indexer_manager = { ..._43
    };
    export const off_chain_updates = { ..._44
    };
    export namespace protocol {
      export const v1 = { ..._45,
        ..._46,
        ..._47
      };
    }
    export const redis = { ..._48
    };
    export const shared = { ..._49
    };
    export const socks = { ..._50
    };
  }
  export const perpetuals = { ..._51,
    ..._52,
    ..._53,
    ..._54,
    ..._55,
    ..._95,
    ..._106,
    ..._116
  };
  export const prices = { ..._56,
    ..._57,
    ..._58,
    ..._59,
    ..._60,
    ..._96,
    ..._107,
    ..._117
  };
  export const ratelimit = { ..._61,
    ..._62,
    ..._63,
    ..._64,
    ..._65,
    ..._66,
    ..._97,
    ..._108,
    ..._118
  };
  export const sending = { ..._67,
    ..._68,
    ..._69,
    ..._70,
    ..._109,
    ..._119
  };
  export const stats = { ..._71,
    ..._72,
    ..._73,
    ..._74,
    ..._75,
    ..._98,
    ..._110,
    ..._120
  };
  export const subaccounts = { ..._76,
    ..._77,
    ..._78,
    ..._79,
    ..._80,
    ..._99,
    ..._111
  };
  export const ve = { ..._81
  };
  export const ClientFactory = { ..._121,
    ..._122,
    ..._123
  };
}