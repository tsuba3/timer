package main

import (
	"errors"
	"flag"
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

func WriteTime(screen tcell.Screen, d time.Duration, showSecond bool) {
	if _, h := screen.Size(); h >= numberHeight {
		var width int
		if showSecond {
			width = numberWidth*6 + separatorWidth*2 + 1*7
		} else {
			width = numberWidth*4 + separatorWidth + 1*4
		}

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

		if showSecond {
			WriteSeparator(screen, offset)
			offset.x += separatorWidth + 1

			WriteNumber(screen, offset, int(d.Seconds())%60/10%10)
			offset.x += numberWidth + 1

			WriteNumber(screen, offset, int(d.Seconds())%10)
			offset.x += numberWidth + 1
		}
	} else {
		width := len([]rune("00:00"))
		if showSecond {
			width = len([]rune("00:00:00"))
		}
		height := 1
		offset := CalcCenterOffset(screen, width, height)
		var str string
		if showSecond {
			str = fmt.Sprintf("%02d:%02d:%02d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
		} else {
			str = fmt.Sprintf("%02d:%02d", int(d.Hours()), int(d.Minutes())%60)
		}
		WriteLine(screen, offset, str)
	}
}

func StopWatchLoop(screen tcell.Screen, option Option, quit <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	start := time.Now()

	screen.Clear()
	WriteTime(screen, 0, option.showSecond)
	screen.Show()

	for {
		select {
		case <-quit:
			return
		case t := <-ticker.C:
			screen.Clear()
			WriteTime(screen, t.Sub(start), option.showSecond)
			screen.Show()
		}
	}
}

func TimerLoop(screen tcell.Screen, option Option, quit <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	screen.Clear()
	WriteTime(screen, option.finishTime.Sub(time.Now()), option.showSecond)
	screen.Show()

	for {
		select {
		case <-quit:
			return
		case t := <-ticker.C:
			if t.After(option.finishTime) {
				ticker.Stop()
				switch option.onEnd {
				case EndImmediately:
					return
				case EndBlink:
					BlinkLoop(screen, quit)
					return
				case EndFreeze:
				}
			}
			screen.Clear()
			WriteTime(screen, option.finishTime.Sub(t), option.showSecond)
			screen.Show()
		}
	}
}

func BlinkLoop(screen tcell.Screen, quit <-chan struct{}) {
	ch := make(chan bool, 1)
	go func() {
		for {
			ch <- true
			time.Sleep(150 * time.Millisecond)
			ch <- false
			time.Sleep(150 * time.Millisecond)
			ch <- true
			time.Sleep(150 * time.Millisecond)
			ch <- false
			time.Sleep(400 * time.Millisecond)
		}
	}()
	for {
		select {
		case <-quit:
			return
		case blink := <-ch:
			if blink {
				screen.Fill(' ', tcell.StyleDefault.Background(foregroundColor))
			} else {
				screen.Clear()
			}
			screen.Show()
		}
	}
}

type Option struct {
	showSecond bool
	onEnd      EndOption
	finishTime time.Time
	countUp    bool
}

type EndOption int

const (
	EndImmediately EndOption = iota
	EndBlink                 = iota
	EndFreeze                = iota
)

func ParseOption() (Option, error) {
	option := Option{}

	flag.BoolVar(&option.showSecond, "s", false, "Show second")
	flag.BoolVar(&option.countUp, "u", false, "Count up")

	var endBlink, endFreeze bool
	flag.BoolVar(&endBlink, "b", false, "Blink on time elapsed")
	flag.BoolVar(&endFreeze, "f", false, "Stop on time elapsed")

	finishTimeString := flag.String("t", "", "Finish Time (06:30)")

	flag.Parse()

	if endBlink {
		option.onEnd = EndBlink
	}
	if endFreeze {
		option.onEnd = EndFreeze
	}

	if *finishTimeString != "" {
		var h, m int
		fmt.Sscanf(*finishTimeString, "%d:%d", &h, &m)
		now := time.Now()
		t := time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, time.Local)
		if t.Before(time.Now()) {
			t = t.Add(24 * time.Hour)
		}
		option.finishTime = t
	}

	if flag.NArg() == 1 {
		d, err := time.ParseDuration(flag.Args()[0])
		if err != nil {
			return option, err
		}
		option.finishTime = time.Now().Add(d)
	}

	if option.finishTime.IsZero() && !option.countUp {
		return option, errors.New("time not specified")
	}

	return option, nil
}

func main() {
	option, err := ParseOption()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	quit := make(chan struct{})

	screen.Init()
	go func() {
		for {
			event := screen.PollEvent()
			switch event := event.(type) {
			case *tcell.EventResize:
				screen.Sync()
			case *tcell.EventKey:
				if event.Key() == tcell.KeyESC || event.Rune() == 'q' || event.Key() == tcell.KeyCtrlC {
					close(quit)
				}
			}
		}
	}()

	if option.countUp {
		StopWatchLoop(screen, option, quit)
	} else {
		TimerLoop(screen, option, quit)
	}

	screen.Fini()
	os.Exit(0)
}
