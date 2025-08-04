package util

import (
	"strings"
)

// Diff calculates the difference score between two strings
// by computing the total number of lines removed and added.
// @todo #81:90min Repair Diff method to return the number of lines that were removed and added
// by some reason this method returns to many changes lines. For example, if a file has 15 changes,
// this method migt return 200. We need to find the bug and fix it.
func Diff(before, after string) int {
	blines := strings.Split(before, "\n")
	alines := strings.Split(after, "\n")
	longest := lcs(blines, alines)
	return len(blines) - longest + len(alines) - longest
}

// lcs computes the length of the longest common subsequence (LCS) between two slices of strings.
// You can read more about LCS here: https://en.wikipedia.org/wiki/Longest_common_subsequence
func lcs(a, b []string) int {
	n := len(a) + 1
	m := len(b) + 1
	arr := make([][]int, n)
	for i := range n {
		arr[i] = make([]int, m)
		for j := range m {
			arr[i][j] = 0
		}
	}
	for i := range len(a) {
		for j := range len(b) {
			if a[i] == b[j] {
				arr[i+1][j+1] = arr[i][j] + 1
			} else {
				arr[i+1][j+1] = max(arr[i][j+1], arr[i+1][j])
			}
		}
	}
	return arr[len(a)][len(b)]
}
