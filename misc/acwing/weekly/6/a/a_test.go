// Code generated by copypasta/template/acwing/generator_test.go
package main

import (
	"github.com/EndlessCheng/codeforces-go/main/testutil"
	"testing"
)

func Test_run(t *testing.T) {
	t.Log("Current test is [a]")
	testCases := [][2]string{
		{
			`1
1`,
			`1`,
		},
		{
			`10
1 2 2 3 4 4 4 2 2 2`,
			`8`,
		},
		{
			`11
2 2 2 2 2 2 2 2 2 2 99`,
			`1`,
		},
		
	}
	target := 0 // -1
	testutil.AssertEqualStringCase(t, testCases, target, run)
}
