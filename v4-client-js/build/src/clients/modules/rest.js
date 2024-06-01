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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicmVzdC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9jbGllbnRzL21vZHVsZXMvcmVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUFBLDRDQUFtRDtBQUNuRCxnRUFBK0Q7QUFDL0Qsd0NBQWdFO0FBR2hFLE1BQXFCLFVBQVU7SUFJM0IsWUFBWSxJQUFZLEVBQUUsVUFBeUI7UUFDakQsSUFBSSxJQUFJLENBQUMsUUFBUSxDQUFDLEdBQUcsQ0FBQyxFQUFFO1lBQ3RCLElBQUksQ0FBQyxJQUFJLEdBQUcsSUFBSSxDQUFDLEtBQUssQ0FBQyxDQUFDLEVBQUUsQ0FBQyxDQUFDLENBQUMsQ0FBQztTQUMvQjthQUFNO1lBQ0wsSUFBSSxDQUFDLElBQUksR0FBRyxJQUFJLENBQUM7U0FDbEI7UUFDRCxJQUFJLENBQUMsVUFBVSxHQUFHLFVBQVUsSUFBSSwrQkFBbUIsQ0FBQztJQUN0RCxDQUFDO0lBRUQsS0FBSyxDQUFDLEdBQUcsQ0FDUCxXQUFtQixFQUNuQixTQUFhLEVBQUU7UUFFZixNQUFNLEdBQUcsR0FBRyxHQUFHLElBQUksQ0FBQyxJQUFJLEdBQUcsSUFBQSxtQ0FBaUIsRUFBQyxXQUFXLEVBQUUsTUFBTSxDQUFDLEVBQUUsQ0FBQztRQUNwRSxNQUFNLFFBQVEsR0FBRyxNQUFNLElBQUEsZUFBTyxFQUFDLEdBQUcsQ0FBQyxDQUFDO1FBQ3BDLE9BQU8sUUFBUSxDQUFDLElBQUksQ0FBQztJQUN2QixDQUFDO0lBRUQsS0FBSyxDQUFDLElBQUksQ0FDUixXQUFtQixFQUNuQixTQUFhLEVBQUUsRUFDZixJQUFxQixFQUNyQixVQUFjLEVBQUU7UUFFaEIsTUFBTSxHQUFHLEdBQUcsR0FBRyxJQUFJLENBQUMsSUFBSSxHQUFHLElBQUEsbUNBQWlCLEVBQUMsV0FBVyxFQUFFLE1BQU0sQ0FBQyxFQUFFLENBQUM7UUFDcEUsT0FBTyxJQUFBLGVBQU8sRUFBQyxHQUFHLEVBQUUscUJBQWEsQ0FBQyxJQUFJLEVBQUUsSUFBSSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQ3pELENBQUM7Q0FDSjtBQS9CRCw2QkErQkMifQ==