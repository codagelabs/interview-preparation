package main

import "fmt"


// StringsAreRotaions checks if s2 is a rotation of s1.
// A rotation means s2 can be obtained by shifting some leading characters of s1 to its end.
func StringsAreRotaions(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	if len(s1) == 0 && len(s2) == 0 {
		return true
	}
	// If s2 is a substring of s1+s1, then s2 is a rotation of s1.
	return isSubstring(s1+s1, s2)
}

// isSubstring checks if substr is a substring of str.
func isSubstring(str, substr string) bool {
	n, m := len(str), len(substr)
	if m == 0 {
		return true
	}
	for i := 0; i <= n-m; i++ {
		if str[i:i+m] == substr {
			return true
		}
	}
	return false
}
// BestRotationType determines whether left or right rotation is a better match to convert s1 to s2.
// It returns "left", "right", or "none" depending on which rotation (if any) can transform s1 into s2
// with the minimal number of shifts. If both are possible with the same number of shifts, "left" is preferred.
func BestRotationType(s1, s2 string) string {
	if len(s1) != len(s2) {
		return "none"
	}
	if s1 == s2 {
		return "left" // No rotation needed, but left is as good as right
	}
	n := len(s1)
	// Try all possible left rotations
	for i := 1; i < n; i++ {
		leftRot := s1[i:] + s1[:i]
		if leftRot == s2 {
			// Check if right rotation can do it in fewer steps
			rightSteps := n - i
			// Try right rotation
			rightRot := s1[n-rightSteps:] + s1[:n-rightSteps]
			if rightRot == s2 && rightSteps < i {
				return "right"
			}
			return "left"
		}
	}
	// Try all possible right rotations
	for i := 1; i < n; i++ {
		rightRot := s1[n-i:] + s1[:n-i]
		if rightRot == s2 {
			return "right"
		}
	}
	return "none"
}


func main() {
	fmt.Println(StringsAreRotaions("abcd", "cdab"))
	
}
