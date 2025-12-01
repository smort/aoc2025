package util

import (
	"bufio"
	"os"
)

func GetLines(filename string) []string {
	return GetLinesTransformed(filename, func(s string) (string, error) {
		return s, nil
	})
	// // Open the file
	// file, err := os.Open(filename)
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	// // Create a scanner to read the file
	// scanner := bufio.NewScanner(file)
	// var lines []string

	// // Read each line and append to the slice
	// for scanner.Scan() {
	// 	lines = append(lines, scanner.Text())
	// }

	// // Check for errors during scanning
	// if err := scanner.Err(); err != nil {
	// 	panic(err)
	// }

	// return lines
}

func GetLinesTransformed[T any](filename string, transform func(string) (T, error)) []T {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)
	var lines []T

	// Read each line and append to the slice
	for scanner.Scan() {
		out, err := transform(scanner.Text())
		if err != nil {
			panic(err)
		}
		lines = append(lines, out)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return lines
}
