// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleRomanConverter {
    // 使用固定数组存储罗马数字映射
    uint256[13] private values = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1];
    string[13] private symbols = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"];

    // 整数转罗马数字
    function toRoman(uint256 num) public pure returns (string memory) {
        require(num > 0 && num <= 3999, "Number must be between 1 and 3999");

        string memory result;
        uint256 temp = num;

        for (uint256 i = 0; i < values.length; i++) {
            while (temp >= values[i]) {
                result = string(abi.encodePacked(result, symbols[i]));
                temp -= values[i];
            }
        }

        return result;
    }
}