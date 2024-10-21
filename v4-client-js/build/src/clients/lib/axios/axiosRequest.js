"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.request = request;
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXhpb3NSZXF1ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vc3JjL2NsaWVudHMvbGliL2F4aW9zL2F4aW9zUmVxdWVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7OztBQStCQSwwQkFZQztBQTNDRCxrREFBa0Q7QUFHbEQscUNBR2tCO0FBQ2xCLG1DQUF3QztBQVF4QyxLQUFLLFVBQVUsWUFBWSxDQUFDLE9BQTJCO0lBQ3JELElBQUksQ0FBQztRQUNILE9BQU8sTUFBTSxJQUFBLGVBQUssRUFBQyxPQUFPLENBQUMsQ0FBQztJQUM5QixDQUFDO0lBQUMsT0FBTyxLQUFLLEVBQUUsQ0FBQztRQUNmLHlFQUF5RTtRQUN6RSxJQUFJLEtBQUssQ0FBQyxZQUFZLEVBQUUsQ0FBQztZQUN2Qix5RUFBeUU7WUFDekUsSUFBSSxLQUFLLENBQUMsUUFBUSxFQUFFLENBQUM7Z0JBQ25CLE1BQU0sSUFBSSx5QkFBZ0IsQ0FBQyxLQUFLLENBQUMsUUFBUSxFQUFFLEtBQUssQ0FBQyxDQUFDO1lBQ3BELENBQUM7WUFDRCxNQUFNLElBQUksbUJBQVUsQ0FBQyxVQUFVLEtBQUssQ0FBQyxPQUFPLEVBQUUsRUFBRSxLQUFLLENBQUMsQ0FBQztRQUN6RCxDQUFDO1FBQ0QsTUFBTSxLQUFLLENBQUM7SUFDZCxDQUFDO0FBQ0gsQ0FBQztBQUVELFNBQWdCLE9BQU8sQ0FDckIsR0FBVyxFQUNYLFNBQXdCLHFCQUFhLENBQUMsR0FBRyxFQUN6QyxJQUFxQixFQUNyQixVQUFjLEVBQUU7SUFFaEIsT0FBTyxZQUFZLENBQUM7UUFDbEIsR0FBRztRQUNILE1BQU07UUFDTixJQUFJLEVBQUUsSUFBSTtRQUNWLE9BQU87S0FDUixDQUFDLENBQUM7QUFDTCxDQUFDIn0=