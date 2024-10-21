"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../constants");
const request_helpers_1 = require("../helpers/request-helpers");
const axios_1 = require("../lib/axios");
class RestClient {
    constructor(host, apiTimeout) {
        if (host.endsWith('/')) {
            this.host = host.slice(0, -1);
        }
        else {
            this.host = host;
        }
        this.apiTimeout = apiTimeout || constants_1.DEFAULT_API_TIMEOUT;
    }
    async get(requestPath, params = {}) {
        const url = `${this.host}${(0, request_helpers_1.generateQueryPath)(requestPath, params)}`;
        const response = await (0, axios_1.request)(url);
        return response.data;
    }
    async post(requestPath, params = {}, body, headers = {}) {
        const url = `${this.host}${(0, request_helpers_1.generateQueryPath)(requestPath, params)}`;
        return (0, axios_1.request)(url, axios_1.RequestMethod.POST, body, headers);
    }
}
exports.default = RestClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicmVzdC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9jbGllbnRzL21vZHVsZXMvcmVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUFBLDRDQUFtRDtBQUNuRCxnRUFBK0Q7QUFDL0Qsd0NBQWdFO0FBR2hFLE1BQXFCLFVBQVU7SUFJM0IsWUFBWSxJQUFZLEVBQUUsVUFBeUI7UUFDakQsSUFBSSxJQUFJLENBQUMsUUFBUSxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUM7WUFDdkIsSUFBSSxDQUFDLElBQUksR0FBRyxJQUFJLENBQUMsS0FBSyxDQUFDLENBQUMsRUFBRSxDQUFDLENBQUMsQ0FBQyxDQUFDO1FBQ2hDLENBQUM7YUFBTSxDQUFDO1lBQ04sSUFBSSxDQUFDLElBQUksR0FBRyxJQUFJLENBQUM7UUFDbkIsQ0FBQztRQUNELElBQUksQ0FBQyxVQUFVLEdBQUcsVUFBVSxJQUFJLCtCQUFtQixDQUFDO0lBQ3RELENBQUM7SUFFRCxLQUFLLENBQUMsR0FBRyxDQUNQLFdBQW1CLEVBQ25CLFNBQWEsRUFBRTtRQUVmLE1BQU0sR0FBRyxHQUFHLEdBQUcsSUFBSSxDQUFDLElBQUksR0FBRyxJQUFBLG1DQUFpQixFQUFDLFdBQVcsRUFBRSxNQUFNLENBQUMsRUFBRSxDQUFDO1FBQ3BFLE1BQU0sUUFBUSxHQUFHLE1BQU0sSUFBQSxlQUFPLEVBQUMsR0FBRyxDQUFDLENBQUM7UUFDcEMsT0FBTyxRQUFRLENBQUMsSUFBSSxDQUFDO0lBQ3ZCLENBQUM7SUFFRCxLQUFLLENBQUMsSUFBSSxDQUNSLFdBQW1CLEVBQ25CLFNBQWEsRUFBRSxFQUNmLElBQXFCLEVBQ3JCLFVBQWMsRUFBRTtRQUVoQixNQUFNLEdBQUcsR0FBRyxHQUFHLElBQUksQ0FBQyxJQUFJLEdBQUcsSUFBQSxtQ0FBaUIsRUFBQyxXQUFXLEVBQUUsTUFBTSxDQUFDLEVBQUUsQ0FBQztRQUNwRSxPQUFPLElBQUEsZUFBTyxFQUFDLEdBQUcsRUFBRSxxQkFBYSxDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsT0FBTyxDQUFDLENBQUM7SUFDekQsQ0FBQztDQUNKO0FBL0JELDZCQStCQyJ9