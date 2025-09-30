// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract StringReverser {
    // 反转字符串函数
    function reverseString(string memory _str) public pure returns (string memory) {
        // 将字符串转换为bytes类型以便操作
        bytes memory strBytes = bytes(_str);
        bytes memory reversed = new bytes(strBytes.length);

        // 反转操作
        for (uint i = 0; i < strBytes.length; i++) {
            reversed[i] = strBytes[strBytes.length - 1 - i];
        }

        return string(reversed);
    }

    // 测试函数
    function testReverse() public pure returns (string memory) {
        string memory original = "abcde";
        string memory reversed = reverseString(original);
        return reversed; // 返回 "edcba"
    }
}