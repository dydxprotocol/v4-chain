// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";

contract Oracle is Ownable {
    mapping(bytes32 => int256) public prices;
    mapping(bytes32 => uint256) public lastUpdated;

    event PriceUpdated(bytes32 indexed marketId, int256 price, uint256 timestamp);

    constructor() Ownable(msg.sender) {}

    function setPrice(bytes32 marketId, int256 price) external {
        require(price > 0, "Price must be positive");
        prices[marketId] = price;
        lastUpdated[marketId] = block.timestamp;
        emit PriceUpdated(marketId, price, block.timestamp);
    }

    function getPrice(bytes32 marketId) external view returns (int256 price, uint256 timestamp) {
        return (prices[marketId], lastUpdated[marketId]);
    }
}
