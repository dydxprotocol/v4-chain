"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateQueryPath = generateQueryPath;
function generateQueryPath(url, params) {
    const definedEntries = Object.entries(params)
        .filter(([_key, value]) => value !== undefined);
    if (!definedEntries.length) {
        return url;
    }
    const paramsString = definedEntries.map(([key, value]) => `${key}=${value}`).join('&');
    return `${url}?${paramsString}`;
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicmVxdWVzdC1oZWxwZXJzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2NsaWVudHMvaGVscGVycy9yZXF1ZXN0LWhlbHBlcnMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFBQSw4Q0FZQztBQVpELFNBQWdCLGlCQUFpQixDQUFDLEdBQVcsRUFBRSxNQUFVO0lBQ3ZELE1BQU0sY0FBYyxHQUFHLE1BQU0sQ0FBQyxPQUFPLENBQUMsTUFBTSxDQUFDO1NBQzFDLE1BQU0sQ0FBQyxDQUFDLENBQUMsSUFBSSxFQUFFLEtBQUssQ0FBb0IsRUFBRSxFQUFFLENBQUMsS0FBSyxLQUFLLFNBQVMsQ0FBQyxDQUFDO0lBRXJFLElBQUksQ0FBQyxjQUFjLENBQUMsTUFBTSxFQUFFLENBQUM7UUFDM0IsT0FBTyxHQUFHLENBQUM7SUFDYixDQUFDO0lBRUQsTUFBTSxZQUFZLEdBQUcsY0FBYyxDQUFDLEdBQUcsQ0FDckMsQ0FBQyxDQUFDLEdBQUcsRUFBRSxLQUFLLENBQW9CLEVBQUUsRUFBRSxDQUFDLEdBQUcsR0FBRyxJQUFJLEtBQUssRUFBRSxDQUN2RCxDQUFDLElBQUksQ0FBQyxHQUFHLENBQUMsQ0FBQztJQUNaLE9BQU8sR0FBRyxHQUFHLElBQUksWUFBWSxFQUFFLENBQUM7QUFDbEMsQ0FBQyJ9