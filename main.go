package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	screen.Init()
	screen.SetContent(0, 0, 'あ', nil, 0)
	screen.Show()
	time.Sleep(5 * time.Second)
	screen.Fini()
}
