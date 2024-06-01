"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.balance = this.balance.bind(this);
        this.owner = this.owner.bind(this);
        this.supply = this.supply.bind(this);
        this.nFTs = this.nFTs.bind(this);
        this.nFT = this.nFT.bind(this);
        this.class = this.class.bind(this);
        this.classes = this.classes.bind(this);
    }
    /* Balance queries the number of NFTs of a given class owned by the owner, same as balanceOf in ERC721 */
    async balance(params) {
        const endpoint = `cosmos/nft/v1beta1/balance/${params.owner}/${params.classId}`;
        return await this.req.get(endpoint);
    }
    /* Owner queries the owner of the NFT based on its class and id, same as ownerOf in ERC721 */
    async owner(params) {
        const endpoint = `cosmos/nft/v1beta1/owner/${params.classId}/${params.id}`;
        return await this.req.get(endpoint);
    }
    /* Supply queries the number of NFTs from the given class, same as totalSupply of ERC721. */
    async supply(params) {
        const endpoint = `cosmos/nft/v1beta1/supply/${params.classId}`;
        return await this.req.get(endpoint);
    }
    /* NFTs queries all NFTs of a given class or owner,choose at least one of the two, similar to tokenByIndex in
     ERC721Enumerable */
    async nFTs(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.classId) !== "undefined") {
            options.params.class_id = params.classId;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.owner) !== "undefined") {
            options.params.owner = params.owner;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/nft/v1beta1/nfts`;
        return await this.req.get(endpoint, options);
    }
    /* NFT queries an NFT based on its class and id. */
    async nFT(params) {
        const endpoint = `cosmos/nft/v1beta1/nfts/${params.classId}/${params.id}`;
        return await this.req.get(endpoint);
    }
    /* Class queries an NFT class based on its id */
    async class(params) {
        const endpoint = `cosmos/nft/v1beta1/classes/${params.classId}`;
        return await this.req.get(endpoint);
    }
    /* Classes queries all NFT classes */
    async classes(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/nft/v1beta1/classes`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL25mdC92MWJldGExL3F1ZXJ5LmxjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSw4Q0FBdUQ7QUFHdkQsTUFBYSxjQUFjO0lBR3pCLFlBQVksRUFDVixhQUFhLEVBR2Q7UUFDQyxJQUFJLENBQUMsR0FBRyxHQUFHLGFBQWEsQ0FBQztRQUN6QixJQUFJLENBQUMsT0FBTyxHQUFHLElBQUksQ0FBQyxPQUFPLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3ZDLElBQUksQ0FBQyxLQUFLLEdBQUcsSUFBSSxDQUFDLEtBQUssQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDbkMsSUFBSSxDQUFDLE1BQU0sR0FBRyxJQUFJLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNyQyxJQUFJLENBQUMsSUFBSSxHQUFHLElBQUksQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ2pDLElBQUksQ0FBQyxHQUFHLEdBQUcsSUFBSSxDQUFDLEdBQUcsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDL0IsSUFBSSxDQUFDLEtBQUssR0FBRyxJQUFJLENBQUMsS0FBSyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNuQyxJQUFJLENBQUMsT0FBTyxHQUFHLElBQUksQ0FBQyxPQUFPLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQ3pDLENBQUM7SUFDRCx5R0FBeUc7SUFHekcsS0FBSyxDQUFDLE9BQU8sQ0FBQyxNQUEyQjtRQUN2QyxNQUFNLFFBQVEsR0FBRyw4QkFBOEIsTUFBTSxDQUFDLEtBQUssSUFBSSxNQUFNLENBQUMsT0FB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