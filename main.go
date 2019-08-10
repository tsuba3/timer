package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

var colors []string

type Point struct {
	x int
	y int
}

func WriteChar(screen tcell.Screen, option Option, offset Point, points []Point) {
	style := tcell.StyleDefault.Background(option.foregroundColor)
	for _, p := range points {
		screen.SetContent(offset.x+p.x, offset.y+p.y, ' ', nil, style)
	}
}

func WriteNumber(screen tcell.Screen, option Option, offset Point, n int) {
	WriteChar(screen, option, offset, large_number_text[n])
}

func WriteSeparator(screen tcell.Screen, option Option, offset Point) {
	WriteChar(screen, option, offset, separator_text)
}

func CalcCenterOffset(screen tcell.Screen, width int, height int) Point {
	screenWidth, screenHeight := screen.Size()
	offsetLeft := (screenWidth - width) / 2
	offsetTop := (screenHeight - height) / 2
	return Point{offsetLeft, offsetTop}
}

func WriteLine(screen tcell.Screen, option Option, offset Point, str string) {
	style := tcell.StyleDefault.Foreground(option.foregroundColor)
	for i, c := range str {
		screen.SetContent(offset.x+i, offset.y, c, nil, style)
	}
}

func WriteTime(screen tcell.Screen, option Option, d time.Duration) {
	if _, h := screen.Size(); h >= numberHeight {
		var width int
		if option.showSecond {
			width = numberWidth*6 + separatorWidth*2 + 1*7
		} else {
			width = numberWidth*4 + separatorWidth + 1*4
		}

		height := numberHeight
		offset := CalcCenterOffset(screen, width, height)

		WriteNumber(screen, option, offset, int(d.Hours()/10)%10)
		offset.x += numberWidth + 1

		WriteNumber(screen, option, offset, int(d.Hours())%10)
		offset.x += numberWidth + 1

		WriteSeparator(screen, option, offset)
		offset.x += separatorWidth + 1

		WriteNumber(screen, option, offset, int(d.Minutes())%60/10%10)
		offset.x += numberWidth + 1

		WriteNumber(screen, option, offset, int(d.Minutes())%10)
		offset.x += numberWidth + 1

		if option.showSecond {
			WriteSeparator(screen, option, offset)
			offset.x += separatorWidth + 1

			WriteNumber(screen, option, offset, int(d.Seconds())%60/10%10)
			offset.x += numberWidth + 1

			WriteNumber(screen, option, offset, int(d.Seconds())%10)
			offset.x += numberWidth + 1
		}
	} else {
		width := len([]rune("00:00"))
		if option.showSecond {
			width = len([]rune("00:00:00"))
		}
		height := 1
		offset := CalcCenterOffset(screen, width, height)
		var str string
		if option.showSecond {
			str = fmt.Sprintf("%02d:%02d:%02d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
		} else {
			str = fmt.Sprintf("%02d:%02d", int(d.Hours()), int(d.Minutes())%60)
		}
		WriteLine(screen, option, offset, str)
	}
}

func StopWatchLoop(screen tcell.Screen, option Option, quit <-chan struct{}) int {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	start := time.Now()

	screen.Clear()
	WriteTime(screen, option, 0)
	screen.Show()

	for {
		select {
		case <-quit:
			return 13
		case t := <-ticker.C:
			if !option.finishTime.IsZero() && t.After(option.finishTime) {
				ticker.Stop()
				switch option.onEnd {
				case EndImmediately:
					return 0
				case EndBlink:
					BlinkLoop(screen, option, quit)
					return 0
				case EndFreeze:
				}
			}
			screen.Clear()
			WriteTime(screen, option, t.Sub(start))
			screen.Show()
		}
	}
}

func ShowNowTime(screen tcell.Screen, option Option) {
	screen.Clear()
	h, m, s := time.Now().Clock()
	WriteTime(screen, option, time.Duration(60*(60*h+m)+s)*time.Second)
	screen.Show()
}

func ClockLoop(screen tcell.Screen, option Option, quit <-chan struct{}) {
	ShowNowTime(screen, option)
	time.Sleep(time.Duration(1000000000 - time.Now().Nanosecond()))

	ticker := time.NewTicker(time.Second)
	ShowNowTime(screen, option)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			ShowNowTime(screen, option)
		}
	}
}

func TimerLoop(screen tcell.Screen, option Option, quit <-chan struct{}) int {
	screen.Clear()
	WriteTime(screen, option, option.finishTime.Sub(time.Now()))
	screen.Show()

	time.Sleep(time.Duration(option.finishTime.Sub(time.Now()).Nanoseconds() % 1e9))
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	screen.Clear()
	WriteTime(screen, option, option.finishTime.Sub(time.Now()))
	screen.Show()

	for {
		select {
		case <-quit:
			return 13
		case t := <-ticker.C:
			if t.After(option.finishTime) {
				ticker.Stop()
				switch option.onEnd {
				case EndImmediately:
					return 0
				case EndBlink:
					BlinkLoop(screen, option, quit)
					return 0
				case EndFreeze:
				}
			}
			screen.Clear()
			WriteTime(screen, option, option.finishTime.Sub(t))
			screen.Show()
		}
	}
}

func BlinkLoop(screen tcell.Screen, option Option, quit <-chan struct{}) {
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
				screen.Fill(' ', tcell.StyleDefault.Background(option.foregroundColor))
			} else {
				screen.Clear()
			}
			screen.Show()
		}
	}
}

type Option struct {
	showSecond      bool
	onEnd           EndOption
	finishTime      time.Time
	foregroundColor tcell.Color
	countUp         bool
}

type EndOption int

const (
	EndImmediately EndOption = iota
	EndBlink                 = iota
	EndFreeze                = iota
)

func ParseColor(str string) (tcell.Color, error) {
	for i, v := range colors {
		if str == v {
			return tcell.Color(i), nil
		}
	}
	return 0, errors.New("undefined color")
}

func ParseOption() (Option, error) {
	var err error
	option := Option{}

	flag.BoolVar(&option.showSecond, "s", false, "Show second")
	flag.BoolVar(&option.countUp, "u", false, "Count up")

	var endBlink, endFreeze bool
	flag.BoolVar(&endBlink, "b", false, "Blink on time elapsed")
	flag.BoolVar(&endFreeze, "f", false, "Stop on time elapsed")

	finishTimeString := flag.String("t", "", "Finish Time (06:30)")
	color := flag.String("c", "blue", "The foreground color")

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

	option.foregroundColor, err = ParseColor(*color)
	if err != nil {
		errText := make([]byte, 0, 200)
		errText = append(errText, "Undefined color. Use one of listed colors below.\n\n"...)
		for _, v := range colors {
			errText = append(errText, v...)
			errText = append(errText, '\n')
		}
		errText = append(errText, '\n')
		return option, errors.New(string(errText))
	}

	if flag.NArg() == 1 {
		d, err := time.ParseDuration(flag.Args()[0])
		if err != nil {
			return option, err
		}
		option.finishTime = time.Now().Add(d)
	}

	return option, nil
}

func init() {
	colors = []string{"black", "maroon", "green", "olive", "navy", "purple", "teal", "silver", "gray", "red", "lime", "yellow", "blue", "fuchsia", "aqua", "white"}
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

	var code int = 0
	if option.countUp {
		code = StopWatchLoop(screen, option, quit)
	} else if option.finishTime.IsZero() {
		ClockLoop(screen, option, quit)
	} else {
		code = TimerLoop(screen, option, quit)
	}

	screen.Fini()
	os.Exit(code)
}
