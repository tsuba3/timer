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

func WriteNumber(screen tcell.Screen, offcet Point, n int) {
	style := tcell.StyleDefault.Background(foregroundColor)
	for _, p := range large_number_text[n] {
		screen.SetContent(offcet.x+p.x, offcet.y+p.y, ' ', nil, style)
	}
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
	screen.Show()
	time.Sleep(5 * time.Second)
	screen.Fini()
}
