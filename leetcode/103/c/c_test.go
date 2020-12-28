// Code generated by copypasta/template/leetcode/generator_test.go
package main

import (
	"github.com/EndlessCheng/codeforces-go/leetcode/testutil"
	"testing"
)

func Test(t *testing.T) {
	t.Log("Current test is [c]")
	examples := [][]string{
		{
			`[1]`, `0`, 
			`0`,
		},
		{
			`[0,10]`, `2`, 
			`6`,
		},
		{
			`[1,3,6]`, `3`, 
			`3`,
		},
		// TODO 测试入参最小的情况
		
	}
	targetCaseNum := 0 // -1
	if err := testutil.RunLeetCodeFuncWithExamples(t, smallestRangeII, examples, targetCaseNum); err != nil {
		t.Fatal(err)
	}
}
// https://leetcode-cn.com/contest/weekly-contest-103/problems/smallest-range-ii/
