"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateQueryPath = void 0;
function generateQueryPath(url, params) {
    const definedEntries = Object.entries(params)
        .filter(([_key, value]) => value !== undefined);
    if (!definedEntries.length) {
        return url;
    }
    const paramsString = definedEntries.map(([key, value]) => `${key}=${value}`).join('&');
    return `${url}?${paramsString}`;
}
exports.generateQueryPath = generateQueryPath;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicmVxdWVzdC1oZWxwZXJzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2NsaWVudHMvaGVscGVycy9yZXF1ZXN0LWhlbHBlcnMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEsU0FBZ0IsaUJBQWlCLENBQUMsR0FBVyxFQUFFLE1BQVU7SUFDdkQsTUFBTSxjQUFjLEdBQUcsTUFBTSxDQUFDLE9BQU8sQ0FBQyxNQUFNLENBQUM7U0FDMUMsTUFBTSxDQUFDLENBQUMsQ0FBQyxJQUFJLEVBQUUsS0FBSyxDQUFvQixFQUFFLEVBQUUsQ0FBQyxLQUFLLEtBQUssU0FBUyxDQUFDLENBQUM7SUFFckUsSUFBSSxDQUFDLGNBQWMsQ0FBQyxNQUFNLEVBQUU7UUFDMUIsT0FBTyxHQUFHLENBQUM7S0FDWjtJQUVELE1BQU0sWUFBWSxHQUFHLGNBQWMsQ0FBQyxHQUFHLENBQ3JDLENBQUMsQ0FBQyxHQUFHLEVBQUUsS0FBSyxDQUFvQixFQUFFLEVBQUUsQ0FBQyxHQUFHLEdBQUcsSUFBSSxLQUFLLEVBQUUsQ0FDdkQsQ0FBQyxJQUFJLENBQUMsR0FBRyxDQUFDLENBQUM7SUFDWixPQUFPLEdBQUcsR0FBRyxJQUFJLFlBQVksRUFBRSxDQUFDO0FBQ2xDLENBQUM7QUFaRCw4Q0FZQyJ9