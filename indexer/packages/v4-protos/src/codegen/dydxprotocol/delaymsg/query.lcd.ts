import { LCDClient } from "@osmonauts/lcd";
import { QueryNextDelayedMessageIdRequest, QueryNextDelayedMessageIdResponseSDKType, QueryMessageRequest, QueryMessageResponseSDKType, QueryBlockMessageIdsRequest, QueryBlockMessageIdsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.nextDelayedMessageId = this.nextDelayedMessageId.bind(this);
    this.message = this.message.bind(this);
    this.blockMessageIds = this.blockMessageIds.bind(this);
  }
  /* Queries the next DelayedMessage's id. */


  async nextDelayedMessageId(_params: QueryNextDelayedMessageIdRequest = {}): Promise<QueryNextDelayedMessageIdResponseSDKType> {
    const endpoint = `dydxprotocol/v4/delaymsg/next_id`;
    return await this.req.get<QueryNextDelayedMessageIdResponseSDKType>(endpoint);
  }
  /* Queries the DelayedMessage by id. */


  async message(params: QueryMessageRequest): Promise<QueryMessageResponseSDKType> {
    const endpoint = `dydxprotocol/v4/delaymsg/message/${params.id}`;
    return await this.req.get<QueryMessageResponseSDKType>(endpoint);
  }
  /* Queries the DelayedMessages at a given block height. */


  async blockMessageIds(params: QueryBlockMessageIdsRequest): Promise<QueryBlockMessageIdsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/delaymsg/block/message_ids/${params.blockHeight}`;
    return await this.req.get<QueryBlockMessageIdsResponseSDKType>(endpoint);
  }

}