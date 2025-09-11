package main

import (
	"sort"
	"strconv"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

// 136.只出现一次的数字
func singleNumber(nums []int) int {
	countmap := make(map[int]int)
	for _, num := range nums {
		countmap[num]++
	}
	for num, count := range countmap {
		if count == 1 {
			return num
		}
	}
	return 0
}

// 回文数
func isPalindrome(x int) bool {
	str := strconv.Itoa(x)
	for i := 0; i < len(str); i++ {
		if str[i] != str[len(str)-1-i] {
			return false
		}
	}
	return true
}

// 有效的括号
func isValid(s string) bool {
	if len(s)%2 != 0 { // s 长度必须是偶数
		return false
	}
	mp := map[rune]rune{')': '(', ']': '[', '}': '{'}
	st := []rune{}
	for _, c := range s {
		if mp[c] == 0 { // c 是左括号
			st = append(st, c) // 入栈
		} else { // c 是右括号
			if len(st) == 0 || st[len(st)-1] != mp[c] {
				return false // 没有左括号，或者左括号类型不对
			}
			st = st[:len(st)-1] // 出栈
		}
	}
	return len(st) == 0 // 所有左括号必须匹配完毕
}

// 最长公共前缀
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	for i := 0; i < len(strs[0]); i++ {
		for j := 1; j < len(strs); j++ {
			if i == len(strs[j]) || strs[j][i] != strs[0][i] {
				return strs[0][:i]
			}
		}
	}
	return strs[0]
}

// 加一
func plusOne(digits []int) []int {
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		if digits[i] != 9 {
			digits[i]++
			for j := i + 1; j < n; j++ {
				digits[j] = 0
			}
			return digits
		}
	}
	// digits 中所有的元素均为 9

	digits = make([]int, n+1)
	digits[0] = 1
	return digits
}

// 删除有序数组重复项
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	slow := 0
	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[slow] {
			slow++
			nums[slow] = nums[fast]
		}
	}
	return slow + 1
}

// 合并区间
func merge(intervals [][]int) [][]int {
	// 1. 按照区间的起始值进行排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	// 2. 初始化结果集，放入第一个区间
	merged := [][]int{}
	merged = append(merged, intervals[0])

	// 3. 遍历排序后的区间列表
	for i := 1; i < len(intervals); i++ {
		// 当前区间
		current := intervals[i]
		// 结果集中最后一个区间
		last := merged[len(merged)-1]

		// 4. 检查是否重叠
		if current[0] <= last[1] {
			// 有重叠，合并区间（取结束值的较大者）
			if current[1] > last[1] {
				merged[len(merged)-1][1] = current[1]
			}
		} else {
			// 无重叠，直接加入结果集
			merged = append(merged, current)
		}
	}

	return merged
}

// 两数之和
func twoSum(nums []int, target int) []int {
	for i, x := range nums {
		for j := i + 1; j < len(nums); j++ {
			if x+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}
