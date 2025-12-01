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
