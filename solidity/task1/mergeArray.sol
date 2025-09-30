// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MergeInPlace {
    // 原地合并：nums1有足够空间容纳nums2
    function mergeInPlace(
        uint256[] memory nums1,
        uint256 m, // nums1中的有效元素数量
        uint256[] memory nums2,
        uint256 n  // nums2的长度
    ) public pure returns (uint256[] memory) {
        require(nums1.length >= m + n, "nums1 has insufficient space");

        uint256 i = m - 1;     // nums1有效部分的末尾
        uint256 j = n - 1;     // nums2的末尾
        uint256 k = m + n - 1; // 合并后的末尾

        // 从后往前合并，避免覆盖
        while (i >= 0 && j >= 0) {
            if (nums1[i] > nums2[j]) {
                nums1[k] = nums1[i];
                i--;
            } else {
                nums1[k] = nums2[j];
                j--;
            }
            k--;
        }

        // 复制nums2剩余元素
        while (j >= 0) {
            nums1[k] = nums2[j];
            j--;
            k--;
        }

        return nums1;
    }

    // 测试函数
    function testMergeInPlace() public pure returns (uint256[] memory) {
        uint256[] memory nums1 = new uint256[](10);
        nums1[0] = 1;
        nums1[1] = 3;
        nums1[2] = 5;
        // 前3个是有效元素，后面7个是预留空间

        uint256[] memory nums2 = new uint256[](4);
        nums2[0] = 2;
        nums2[1] = 4;
        nums2[2] = 6;
        nums2[3] = 8;

        return mergeInPlace(nums1, 3, nums2, 4);
        // 返回 [1, 2, 3, 4, 5, 6, 8, 0, 0, 0]
    }
}