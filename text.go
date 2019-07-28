package main

import (
	"strings"
)

const numberWidth = 5
const numberHeight = 5

var large_number_text [][]Point

func init() {
	text := make([]string, 10)
	large_number_text = make([][]Point, 10)

	text[0] = "" +
		"xxxxx\n" +
		"x   x\n" +
		"x   x\n" +
		"x   x\n" +
		"xxxxx\n"

	text[1] = "" +
		"    x\n" +
		"    x\n" +
		"    x\n" +
		"    x\n" +
		"    x\n"

	text[2] = "" +
		"xxxxx\n" +
		"    x\n" +
		"xxxxx\n" +
		"x    \n" +
		"xxxxx\n"

	text[3] = "" +
		"xxxxx\n" +
		"    x\n" +
		"xxxxx\n" +
		"    x\n" +
		"xxxxx\n"

	text[4] = "" +
		"x   x\n" +
		"x   x\n" +
		"xxxxx\n" +
		"    x\n" +
		"    x\n"

	text[5] = "" +
		"xxxxx\n" +
		"x    \n" +
		"xxxxx\n" +
		"    x\n" +
		"xxxxx\n"

	text[6] = "" +
		"xxxxx\n" +
		"x    \n" +
		"xxxxx\n" +
		"x   x\n" +
		"xxxxx\n"

	text[7] = "" +
		"xxxxx\n" +
		"    x\n" +
		"    x\n" +
		"    x\n" +
		"    x\n"

	text[8] = "" +
		"xxxxx\n" +
		"x   x\n" +
		"xxxxx\n" +
		"x   x\n" +
		"xxxxx\n"

	text[9] = "" +
		"xxxxx\n" +
		"x   x\n" +
		"xxxxx\n" +
		"    x\n" +
		"    x\n"

	for i := 0; i <= 9; i++ {
		for y, line := range strings.Split(text[i], "\n") {
			for x, char := range line {
				if char == 'x' {
					large_number_text[i] = append(large_number_text[i], Point{x, y})
				}
			}
		}
	}
}
