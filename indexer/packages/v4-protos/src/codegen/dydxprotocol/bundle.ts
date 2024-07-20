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
import * as _28 from "./daemons/liquidation/liquidation";
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
import * as _46 from "./indexer/protocol/v1/subaccount";
import * as _47 from "./indexer/redis/redis_order";
import * as _48 from "./indexer/shared/removal_reason";
import * as _49 from "./indexer/socks/messages";
import * as _50 from "./perpetuals/genesis";
import * as _51 from "./perpetuals/params";
import * as _52 from "./perpetuals/perpetual";
import * as _53 from "./perpetuals/query";
import * as _54 from "./perpetuals/tx";
import * as _55 from "./prices/genesis";
import * as _56 from "./prices/market_param";
import * as _57 from "./prices/market_price";
import * as _58 from "./prices/query";
import * as _59 from "./prices/tx";
import * as _60 from "./ratelimit/capacity";
import * as _61 from "./ratelimit/genesis";
import * as _62 from "./ratelimit/limit_params";
import * as _63 from "./ratelimit/pending_send_packet";
import * as _64 from "./ratelimit/query";
import * as _65 from "./ratelimit/tx";
import * as _66 from "./sending/genesis";
import * as _67 from "./sending/query";
import * as _68 from "./sending/transfer";
import * as _69 from "./sending/tx";
import * as _70 from "./stats/genesis";
import * as _71 from "./stats/params";
import * as _72 from "./stats/query";
import * as _73 from "./stats/stats";
import * as _74 from "./stats/tx";
import * as _75 from "./subaccounts/asset_position";
import * as _76 from "./subaccounts/genesis";
import * as _77 from "./subaccounts/perpetual_position";
import * as _78 from "./subaccounts/query";
import * as _79 from "./subaccounts/subaccount";
import * as _80 from "./ve/ve";
import * as _88 from "./assets/query.lcd";
import * as _89 from "./blocktime/query.lcd";
import * as _90 from "./clob/query.lcd";
import * as _91 from "./delaymsg/query.lcd";
import * as _92 from "./epochs/query.lcd";
import * as _93 from "./feetiers/query.lcd";
import * as _94 from "./perpetuals/query.lcd";
import * as _95 from "./prices/query.lcd";
import * as _96 from "./ratelimit/query.lcd";
import * as _97 from "./stats/query.lcd";
import * as _98 from "./subaccounts/query.lcd";
import * as _99 from "./assets/query.rpc.Query";
import * as _100 from "./blocktime/query.rpc.Query";
import * as _101 from "./clob/query.rpc.Query";
import * as _102 from "./delaymsg/query.rpc.Query";
import * as _103 from "./epochs/query.rpc.Query";
import * as _104 from "./feetiers/query.rpc.Query";
import * as _105 from "./perpetuals/query.rpc.Query";
import * as _106 from "./prices/query.rpc.Query";
import * as _107 from "./ratelimit/query.rpc.Query";
import * as _108 from "./sending/query.rpc.Query";
import * as _109 from "./stats/query.rpc.Query";
import * as _110 from "./subaccounts/query.rpc.Query";
import * as _111 from "./blocktime/tx.rpc.msg";
import * as _112 from "./clob/tx.rpc.msg";
import * as _113 from "./delaymsg/tx.rpc.msg";
import * as _114 from "./feetiers/tx.rpc.msg";
import * as _115 from "./perpetuals/tx.rpc.msg";
import * as _116 from "./prices/tx.rpc.msg";
import * as _117 from "./ratelimit/tx.rpc.msg";
import * as _118 from "./sending/tx.rpc.msg";
import * as _119 from "./stats/tx.rpc.msg";
import * as _120 from "./lcd";
import * as _121 from "./rpc.query";
import * as _122 from "./rpc.tx";
export namespace dydxprotocol {
  export const assets = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._88,
    ..._99
  };
  export const blocktime = { ..._9,
    ..._10,
    ..._11,
    ..._12,
    ..._13,
    ..._89,
    ..._100,
    ..._111
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
    ..._90,
    ..._101,
    ..._112
  };
  export namespace daemons {
    export const liquidation = { ..._28
    };
    export const pricefeed = { ..._29
    };
  }
  export const delaymsg = { ..._30,
    ..._31,
    ..._32,
    ..._33,
    ..._34,
    ..._91,
    ..._102,
    ..._113
  };
  export const epochs = { ..._35,
    ..._36,
    ..._37,
    ..._92,
    ..._103
  };
  export const feetiers = { ..._38,
    ..._39,
    ..._40,
    ..._41,
    ..._93,
    ..._104,
    ..._114
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
        ..._46
      };
    }
    export const redis = { ..._47
    };
    export const shared = { ..._48
    };
    export const socks = { ..._49
    };
  }
  export const perpetuals = { ..._50,
    ..._51,
    ..._52,
    ..._53,
    ..._54,
    ..._94,
    ..._105,
    ..._115
  };
  export const prices = { ..._55,
    ..._56,
    ..._57,
    ..._58,
    ..._59,
    ..._95,
    ..._106,
    ..._116
  };
  export const ratelimit = { ..._60,
    ..._61,
    ..._62,
    ..._63,
    ..._64,
    ..._65,
    ..._96,
    ..._107,
    ..._117
  };
  export const sending = { ..._66,
    ..._67,
    ..._68,
    ..._69,
    ..._108,
    ..._118
  };
  export const stats = { ..._70,
    ..._71,
    ..._72,
    ..._73,
    ..._74,
    ..._97,
    ..._109,
    ..._119
  };
  export const subaccounts = { ..._75,
    ..._76,
    ..._77,
    ..._78,
    ..._79,
    ..._98,
    ..._110
  };
  export const ve = { ..._80
  };
  export const ClientFactory = { ..._120,
    ..._121,
    ..._122
  };
}