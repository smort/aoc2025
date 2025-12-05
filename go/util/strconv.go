package util

import (
	"fmt"
	"strconv"
)

func MustConvAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("error while converting %s to int: %w", s, err))
	}

	return i
}

func MustConvAtoi2(s []string) (int, int) {
	if len(s) != 2 {
		panic(fmt.Errorf("expected slice of length 2, got %d", len(s)))
	}

	return MustConvAtoi(s[0]), MustConvAtoi(s[1])
}
