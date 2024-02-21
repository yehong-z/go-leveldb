package utils

type Comparator func(a, b any) int

func IntComparator(a, b any) int {
	return a.(int) - b.(int)
}
