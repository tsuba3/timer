package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const foregroundColor = tcell.ColorBlue

type Point struct {
	x int
	y int
}

func WriteChar(screen tcell.Screen, offcet Point, points []Point) {
	style := tcell.StyleDefault.Background(foregroundColor)
	for _, p := range points {
		screen.SetContent(offcet.x+p.x, offcet.y+p.y, ' ', nil, style)
	}
}

func WriteNumber(screen tcell.Screen, offcet Point, n int) {
	WriteChar(screen, offcet, large_number_text[n])
}

func WriteSeparator(screen tcell.Screen, offcet Point) {
	WriteChar(screen, offcet, separator_text)
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	screen.Init()
	for i := 0; i < 10; i++ {
		WriteNumber(screen, Point{1 + 6*i, 1}, i)
	}
	WriteSeparator(screen, Point{61, 1})
	screen.Show()
	time.Sleep(5 * time.Second)
	screen.Fini()
}
