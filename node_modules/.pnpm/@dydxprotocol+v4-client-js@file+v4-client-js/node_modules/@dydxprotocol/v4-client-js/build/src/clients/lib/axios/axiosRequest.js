"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.request = void 0;
const axios_1 = __importDefault(require("axios"));
const errors_1 = require("./errors");
const types_1 = require("./types");
async function axiosRequest(options) {
    try {
        return await (0, axios_1.default)(options);
    }
    catch (error) {
        // eslint-disable-next-line @typescript-eslint/strict-boolean-expressions
        if (error.isAxiosError) {
            // eslint-disable-next-line @typescript-eslint/strict-boolean-expressions
            if (error.response) {
                throw new errors_1.AxiosServerError(error.response, error);
            }
            throw new errors_1.AxiosError(`Axios: ${error.message}`, error);
        }
        throw error;
    }
}
function request(url, method = types_1.RequestMethod.GET, body, headers = {}) {
    return axiosRequest({
        url,
        method,
        data: body,
        headers,
    });
}
exports.request = request;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXhpb3NSZXF1ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vc3JjL2NsaWVudHMvbGliL2F4aW9zL2F4aW9zUmVxdWVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFBQSxrREFBa0Q7QUFHbEQscUNBR2tCO0FBQ2xCLG1DQUF3QztBQVF4QyxLQUFLLFVBQVUsWUFBWSxDQUFDLE9BQTJCO0lBQ3JELElBQUk7UUFDRixPQUFPLE1BQU0sSUFBQSxlQUFLLEVBQUMsT0FBTyxDQUFDLENBQUM7S0FDN0I7SUFBQyxPQUFPLEtBQUssRUFBRTtRQUNkLHlFQUF5RTtRQUN6RSxJQUFJLEtBQUssQ0FBQyxZQUFZLEVBQUU7WUFDdEIseUVBQXlFO1lBQ3pFLElBQUksS0FBSyxDQUFDLFFBQVEsRUFBRTtnQkFDbEIsTUFBTSxJQUFJLHlCQUFnQixDQUFDLEtBQUssQ0FBQyxRQUFRLEVBQUUsS0FBSyxDQUFDLENBQUM7YUFDbkQ7WUFDRCxNQUFNLElBQUksbUJBQVUsQ0FBQyxVQUFVLEtBQUssQ0FBQyxPQUFPLEVBQUUsRUFBRSxLQUFLLENBQUMsQ0FBQztTQUN4RDtRQUNELE1BQU0sS0FBSyxDQUFDO0tBQ2I7QUFDSCxDQUFDO0FBRUQsU0FBZ0IsT0FBTyxDQUNyQixHQUFXLEVBQ1gsU0FBd0IscUJBQWEsQ0FBQyxHQUFHLEVBQ3pDLElBQXFCLEVBQ3JCLFVBQWMsRUFBRTtJQUVoQixPQUFPLFlBQVksQ0FBQztRQUNsQixHQUFHO1FBQ0gsTUFBTTtRQUNOLElBQUksRUFBRSxJQUFJO1FBQ1YsT0FBTztLQUNSLENBQUMsQ0FBQztBQUNMLENBQUM7QUFaRCwwQkFZQyJ9