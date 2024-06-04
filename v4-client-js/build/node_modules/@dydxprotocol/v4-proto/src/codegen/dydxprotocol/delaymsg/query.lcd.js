"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.nextDelayedMessageId = this.nextDelayedMessageId.bind(this);
        this.message = this.message.bind(this);
        this.blockMessageIds = this.blockMessageIds.bind(this);
    }
    /* Queries the next DelayedMessage's id. */
    async nextDelayedMessageId(_params = {}) {
        const endpoint = `dydxprotocol/v4/delaymsg/next_id`;
        return await this.req.get(endpoint);
    }
    /* Queries the DelayedMessage by id. */
    async message(params) {
        const endpoint = `dydxprotocol/v4/delaymsg/message/${params.id}`;
        return await this.req.get(endpoint);
    }
    /* Queries the DelayedMessages at a given block height. */
    async blockMessageIds(params) {
        const endpoint = `dydxprotocol/v4/delaymsg/block/message_ids/${params.blockHeight}`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL2RlbGF5bXNnL3F1ZXJ5LmxjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFFQSxNQUFhLGNBQWM7SUFHekIsWUFBWSxFQUNWLGFBQWEsRUFHZDtRQUNDLElBQUksQ0FBQyxHQUFHLEdBQUcsYUFBYSxDQUFDO1FBQ3pCLElBQUksQ0FBQyxvQkFBb0IsR0FBRyxJQUFJLENBQUMsb0JBQW9CLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ2pFLElBQUksQ0FBQyxPQUFPLEdBQUcsSUFBSSxDQUFDLE9BQU8sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDdkMsSUFBSSxDQUFDLGVBQWUsR0FBRyxJQUFJLENBQUMsZUFBZSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUN6RCxDQUFDO0lBQ0QsMkNBQTJDO0lBRzNDLEtBQUssQ0FBQyxvQkFBb0IsQ0FBQyxVQUE0QyxFQUFFO1FBQ3ZFLE1BQU0sUUFBUSxHQUFHLGtDQUFrQyxDQUFDO1FBQ3BELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBMkMsUUFBUSxDQUFDLENBQUM7SUFDaEYsQ0FBQztJQUNELHVDQUF1QztJQUd2QyxLQUFLLENBQUMsT0FBTyxDQUFDLE1BQTJCO1FBQ3ZDLE1BQU0sUUFBUSxHQUFHLG9DQUFvQyxNQUFNLENBQUMsRUFBRSxFQUFFLENBQUM7UUFDakUsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUE4QixRQUFRLENBQUMsQ0FBQztJQUNuRSxDQUFDO0lBQ0QsMERBQTBEO0lBRzFELEtBQUssQ0FBQyxlQUFlLENBQUMsTUFBbUM7UUFDdkQsTUFBTSxRQUFRLEdBQUcsOENBQThDLE1BQU0sQ0FBQyxXQUFXLEVBQUUsQ0FBQztRQUNwRixPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQXNDLFFBQVEsQ0FBQyxDQUFDO0lBQzNFLENBQUM7Q0FFRjtBQW5DRCx3Q0FtQ0MifQ==