"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.eventParams = this.eventParams.bind(this);
        this.proposeParams = this.proposeParams.bind(this);
        this.safetyParams = this.safetyParams.bind(this);
        this.acknowledgedEventInfo = this.acknowledgedEventInfo.bind(this);
        this.recognizedEventInfo = this.recognizedEventInfo.bind(this);
        this.delayedCompleteBridgeMessages = this.delayedCompleteBridgeMessages.bind(this);
    }
    /* Queries the EventParams. */
    async eventParams(_params = {}) {
        const endpoint = `dydxprotocol/v4/bridge/event_params`;
        return await this.req.get(endpoint);
    }
    /* Queries the ProposeParams. */
    async proposeParams(_params = {}) {
        const endpoint = `dydxprotocol/v4/bridge/propose_params`;
        return await this.req.get(endpoint);
    }
    /* Queries the SafetyParams. */
    async safetyParams(_params = {}) {
        const endpoint = `dydxprotocol/v4/bridge/safety_params`;
        return await this.req.get(endpoint);
    }
    /* Queries the AcknowledgedEventInfo.
     An "acknowledged" event is one that is in-consensus and has been stored
     in-state. */
    async acknowledgedEventInfo(_params = {}) {
        const endpoint = `dydxprotocol/v4/bridge/acknowledged_event_info`;
        return await this.req.get(endpoint);
    }
    /* Queries the RecognizedEventInfo.
     A "recognized" event is one that is finalized on the Ethereum blockchain
     and has been identified by the queried node. It is not yet in-consensus. */
    async recognizedEventInfo(_params = {}) {
        const endpoint = `dydxprotocol/v4/bridge/recognized_event_info`;
        return await this.req.get(endpoint);
    }
    /* Queries all `MsgCompleteBridge` messages that are delayed (not yet
     executed) and corresponding block heights at which they will execute. */
    async delayedCompleteBridgeMessages(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.address) !== "undefined") {
            options.params.address = params.address;
        }
        const endpoint = `dydxprotocol/v4/bridge/delayed_complete_bridge_messages`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL2JyaWRnZS9xdWVyeS5sY2QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBRUEsTUFBYSxjQUFjO0lBR3pCLFlBQVksRUFDVixhQUFhLEVBR2Q7UUFDQyxJQUFJLENBQUMsR0FBRyxHQUFHLGFBQWEsQ0FBQztRQUN6QixJQUFJLENBQUMsV0FBVyxHQUFHLElBQUksQ0FBQyxXQUFXLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQy9DLElBQUksQ0FBQyxhQUFhLEdBQUcsSUFBSSxDQUFDLGFBQWEsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDbkQsSUFBSSxDQUFDLFlBQVksR0FBRyxJQUFJLENBQUMsWUFBWSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNqRCxJQUFJLENBQUMscUJBQXFCLEdBQUcsSUFBSSxDQUFDLHFCQUFxQixDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNuRSxJQUFJLENBQUMsbUJBQW1CLEdBQUcsSUFBSSxDQUFDLG1CQUFtQixDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUMvRCxJQUFJLENBQUMsNkJBQTZCLEdBQUcsSUFBSSxDQUFDLDZCQUE2QixDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUNyRixDQUFDO0lBQ0QsOEJBQThCO0lBRzlCLEtBQUssQ0FBQyxXQUFXLENBQUMsVUFBbUMsRUFBRTtRQUNyRCxNQUFNLFFBQVEsR0FBRyxxQ0FBcUMsQ0FBQztRQUN2RCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQWtDLFFBQVEsQ0FBQyxDQUFDO0lBQ3ZFLENBQUM7SUFDRCxnQ0FBZ0M7SUFHaEMsS0FBSyxDQUFDLGFBQWEsQ0FBQyxVQUFxQyxFQUFFO1FBQ3pELE1BQU0sUUFBUSxHQUFHLHVDQUF1QyxDQUFDO1FBQ3pELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBb0MsUUFBUSxDQUFDLENBQUM7SUFDekUsQ0FBQztJQUNELCtCQUErQjtJQUcvQixLQUFLLENBQUMsWUFBWSxDQUFDLFVBQW9DLEVBQUU7UUFDdkQsTUFBTSxRQUFRLEdBQUcsc0NBQXNDLENBQUM7UUFDeEQsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFtQyxRQUFRLENBQUMsQ0FBQztJQUN4RSxDQUFDO0lBQ0Q7O2lCQUVhO0lBR2IsS0FBSyxDQUFDLHFCQUFxQixDQUFDLFVBQTZDLEVBQUU7UUFDekUsTUFBTSxRQUFRLEdBQUcsZ0RBQWdELENBQUM7UUFDbEUsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUE0QyxRQUFRLENBQUMsQ0FBQztJQUNqRixDQUFDO0lBQ0Q7O2dGQUU0RTtJQUc1RSxLQUFLLENBQUMsbUJBQW1CLENBQUMsVUFBMkMsRUFBRTtRQUNyRSxNQUFNLFFBQVEsR0FBRyw4Q0FBOEMsQ0FBQztRQUNoRSxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQTBDLFFBQVEsQ0FBQyxDQUFDO0lBQy9FLENBQUM7SUFDRDs2RUFDeUU7SUFHekUsS0FBSyxDQUFDLDZCQUE2QixDQUFDLE1BQWlEO1FBQ25GLE1BQU0sT0FBTyxHQUFRO1lBQ25CLE1BQU0sRUFBRSxFQUFFO1NBQ1gsQ0FBQztRQUVGLElBQUksT0FBTyxDQUFBLE1BQU0sYUFBTixNQUFNLHVCQUFOLE1BQU0sQ0FBRSxPQUFPLENBQUEsS0FBSyxXQUFXLEVBQUU7WUFDMUMsT0FBTyxDQUFDLE1BQU0sQ0FBQyxPQUFPLEdBQUcsTUFBTSxDQUFDLE9BQU8sQ0FBQztTQUN6QztRQUVELE1BQU0sUUFBUSxHQUFHLHlEQUF5RCxDQUFDO1FBQzNFLE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBb0QsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQ2xHLENBQUM7Q0FFRjtBQXhFRCx3Q0F3RUMifQ==