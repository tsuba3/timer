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

func WriteChar(screen tcell.Screen, offset Point, points []Point) {
	style := tcell.StyleDefault.Background(foregroundColor)
	for _, p := range points {
		screen.SetContent(offset.x+p.x, offset.y+p.y, ' ', nil, style)
	}
}

func WriteNumber(screen tcell.Screen, offset Point, n int) {
	WriteChar(screen, offset, large_number_text[n])
}

func WriteSeparator(screen tcell.Screen, offset Point) {
	WriteChar(screen, offset, separator_text)
}

func CalcCenterOffset(screen tcell.Screen, width int, height int) Point {
	screenWidth, screenHeight := screen.Size()
	offsetLeft := (screenWidth - width) / 2
	offsetTop := (screenHeight - height) / 2
	return Point{offsetLeft, offsetTop}
}

func WriteLine(screen tcell.Screen, offset Point, str string) {
	style := tcell.StyleDefault.Foreground(foregroundColor)
	for i, c := range str {
		screen.SetContent(offset.x+i, offset.y, c, nil, style)
	}
}

func WriteTime(screen tcell.Screen, d time.Duration) {
	if _, h := screen.Size(); h >= numberHeight {
		width := numberWidth*6 + separatorWidth + 1*6
		height := numberHeight
		offset := CalcCenterOffset(screen, width, height)

		WriteNumber(screen, offset, int(d.Hours()/10)%10)
		offset.x += numberWidth + 1

		WriteNumber(screen, offset, int(d.Hours())%10)
		offset.x += numberWidth + 1

		WriteSeparator(screen, offset)
		offset.x += separatorWidth + 1

		WriteNumber(screen, offset, int(d.Minutes())%60/10%10)
		offset.x += numberWidth + 1

		WriteNumber(screen, offset, int(d.Minutes())%10)
		offset.x += numberWidth + 1

		WriteSeparator(screen, offset)
		offset.x += separatorWidth + 1

		WriteNumber(screen, offset, int(d.Seconds())%60/10%10)
		offset.x += numberWidth + 1

		WriteNumber(screen, offset, int(d.Seconds())%10)
		offset.x += numberWidth + 1
	} else {
		width := len([]rune("00:00:00"))
		height := 1
		offset := CalcCenterOffset(screen, width, height)
		str := fmt.Sprintf("%02d:%02d:%02d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
		WriteLine(screen, offset, str)
	}
}

func TimerLoop(screen tcell.Screen, finishTime time.Time, quit <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	screen.Clear()
	WriteTime(screen, finishTime.Sub(time.Now()))
	screen.Show()

	for {
		select {
		case <-quit:
			return
		case t := <-ticker.C:
			if t.After(finishTime) {
				return
			}
			screen.Clear()
			WriteTime(screen, finishTime.Sub(t))
			screen.Show()
		}
	}
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	finishTime := time.Now().Add(time.Duration(10 * time.Second))
	quit := make(chan struct{})

	screen.Init()
	go func() {
		for {
			event := screen.PollEvent()
			switch event := event.(type) {
			case *tcell.EventResize:
				screen.Sync()
			case *tcell.EventKey:
				if event.Key() == tcell.KeyESC || event.Rune() == 'q' {
					close(quit)
				}
			}
		}
	}()

	TimerLoop(screen, finishTime, quit)

	screen.Fini()
	os.Exit(0)
}
