"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubaccountInfo = void 0;
class SubaccountInfo {
    constructor(wallet, subaccountNumber = 0) {
        if (subaccountNumber < 0 || subaccountNumber > 127) {
            throw new Error('Subaccount number must be between 0 and 127');
        }
        this.wallet = wallet;
        this.subaccountNumber = subaccountNumber;
    }
    get address() {
        const address = this.wallet.address;
        if (address !== undefined) {
            return address;
        }
        else {
            throw new Error('Address not set');
        }
    }
}
exports.SubaccountInfo = SubaccountInfo;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic3ViYWNjb3VudC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9jbGllbnRzL3N1YmFjY291bnQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBRUEsTUFBYSxjQUFjO0lBS3ZCLFlBQVksTUFBbUIsRUFBRSxtQkFBMkIsQ0FBQztRQUMzRCxJQUFJLGdCQUFnQixHQUFHLENBQUMsSUFBSSxnQkFBZ0IsR0FBRyxHQUFHLEVBQUUsQ0FBQztZQUNuRCxNQUFNLElBQUksS0FBSyxDQUFDLDZDQUE2QyxDQUFDLENBQUM7UUFDakUsQ0FBQztRQUNELElBQUksQ0FBQyxNQUFNLEdBQUcsTUFBTSxDQUFDO1FBQ3JCLElBQUksQ0FBQyxnQkFBZ0IsR0FBRyxnQkFBZ0IsQ0FBQztJQUMzQyxDQUFDO0lBRUQsSUFBSSxPQUFPO1FBQ1QsTUFBTSxPQUFPLEdBQUcsSUFBSSxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUM7UUFDcEMsSUFBSSxPQUFPLEtBQUssU0FBUyxFQUFFLENBQUM7WUFDMUIsT0FBTyxPQUFPLENBQUM7UUFDakIsQ0FBQzthQUFNLENBQUM7WUFDTixNQUFNLElBQUksS0FBSyxDQUFDLGlCQUFpQixDQUFDLENBQUM7UUFDckMsQ0FBQztJQUNILENBQUM7Q0FDSjtBQXJCRCx3Q0FxQkMifQ==