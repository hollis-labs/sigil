package cli

import (
	"strings"
)

// SuggestComponentType returns a suggestion for a misspelled component type.
func SuggestComponentType(unknown string, known []string) string {
	best := ""
	bestDist := len(unknown)/2 + 1 // max distance threshold

	for _, k := range known {
		d := levenshtein(strings.ToLower(unknown), strings.ToLower(k))
		if d < bestDist {
			bestDist = d
			best = k
		}
	}

	if best != "" {
		return best
	}

	// Try prefix match
	lower := strings.ToLower(unknown)
	for _, k := range known {
		if strings.HasPrefix(strings.ToLower(k), lower) || strings.HasPrefix(lower, strings.ToLower(k)) {
			return k
		}
	}

	return ""
}

// levenshtein computes the edit distance between two strings.
func levenshtein(a, b string) int {
	la := len(a)
	lb := len(b)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)

	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min3(
				prev[j]+1,
				curr[j-1]+1,
				prev[j-1]+cost,
			)
		}
		prev, curr = curr, prev
	}

	return prev[lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
