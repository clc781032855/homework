// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BinarySearch {
    // 二分查找函数 - 迭代版本
    function binarySearch(uint256[] memory arr, uint256 target)
    public
    pure
    returns (int256)
    {
        require(isSorted(arr), "Array must be sorted");

        uint256 left = 0;
        uint256 right = arr.length;

        while (left < right) {
            uint256 mid = left + (right - left) / 2;

            if (arr[mid] == target) {
                return int256(mid); // 找到目标，返回索引
            } else if (arr[mid] < target) {
                left = mid + 1; // 向右查找
            } else {
                right = mid; // 向左查找
            }
        }

        return -1; // 未找到
    }

    // 验证数组是否已排序
    function isSorted(uint256[] memory arr) public pure returns (bool) {
        if (arr.length <= 1) return true;

        for (uint256 i = 1; i < arr.length; i++) {
            if (arr[i] < arr[i - 1]) {
                return false;
            }
        }
        return true;
    }

    // 查找第一个等于目标的元素
    function findFirstOccurrence(uint256[] memory arr, uint256 target)
    public
    pure
    returns (int256)
    {
        require(isSorted(arr), "Array must be sorted");

        uint256 left = 0;
        uint256 right = arr.length;
        int256 result = -1;

        while (left < right) {
            uint256 mid = left + (right - left) / 2;

            if (arr[mid] >= target) {
                if (arr[mid] == target) {
                    result = int256(mid); // 记录找到的位置
                }
                right = mid;
            } else {
                left = mid + 1;
            }
        }

        return result;
    }

    // 查找最后一个等于目标的元素
    function findLastOccurrence(uint256[] memory arr, uint256 target)
    public
    pure
    returns (int256)
    {
        require(isSorted(arr), "Array must be sorted");

        uint256 left = 0;
        uint256 right = arr.length;
        int256 result = -1;

        while (left < right) {
            uint256 mid = left + (right - left) / 2;

            if (arr[mid] <= target) {
                if (arr[mid] == target) {
                    result = int256(mid);
                }
                left = mid + 1;
            } else {
                right = mid;
            }
        }

        return result;
    }
}