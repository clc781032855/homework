// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomanToInteger {
    // 罗马字符到数值的映射
    mapping(bytes1 => uint256) private romanValues;

    constructor() {
        // 初始化罗马数字映射
        romanValues['I'] = 1;
        romanValues['V'] = 5;
        romanValues['X'] = 10;
        romanValues['L'] = 50;
        romanValues['C'] = 100;
        romanValues['D'] = 500;
        romanValues['M'] = 1000;
    }

    // 验证字符是否为有效的罗马数字
    function isValidRomanChar(bytes1 char) private view returns (bool) {
        return (char == 'I' || char == 'V' || char == 'X' ||
        char == 'L' || char == 'C' || char == 'D' || char == 'M');
    }

    // 验证罗马数字字符串的有效性
    function isValidRoman(string memory s) public view returns (bool) {
        bytes memory romanBytes = bytes(s);
        if (romanBytes.length == 0) return false;

        // 检查所有字符是否有效
        for (uint256 i = 0; i < romanBytes.length; i++) {
            if (!isValidRomanChar(romanBytes[i])) {
                return false;
            }
        }

        return true;
    }

    // 主转换函数：罗马数字转整数
    function romanToInt(string memory s) public view returns (uint256) {
        require(isValidRoman(s), "Invalid Roman numeral");

        bytes memory roman = bytes(s);
        uint256 total = 0;
        uint256 length = roman.length;

        for (uint256 i = 0; i < length; i++) {
            uint256 currentValue = romanValues[roman[i]];

            // 检查下一个字符是否存在且比当前字符大（减法规则）
            if (i < length - 1) {
                uint256 nextValue = romanValues[roman[i + 1]];

                if (currentValue < nextValue) {
                    // 应用减法规则：IV, IX, XL, XC, CD, CM
                    total += (nextValue - currentValue);
                    i++; // 跳过下一个字符，因为已经处理了
                    continue;
                }
            }

            total += currentValue;
        }

        return total;
    }

    // 批量转换测试
    function testConversions() public view returns (uint256[10] memory results) {
        string[10] memory testCases = [
                    "I", "IV", "IX", "XL", "XC",
                    "CD", "CM", "LVIII", "MCMXCIV", "MMMCMXCIX"
            ];

        for (uint256 i = 0; i < testCases.length; i++) {
            results[i] = romanToInt(testCases[i]);
        }
    }
}