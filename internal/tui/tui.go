package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"os/exec"
)

func write_and_open(note_name string) {
	cmd := exec.Command("vim", note_name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Draw1() {
	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyEnter {
			write_and_open("test name")
		}
		return event
	})

	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	if err := app.SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}

func Draw() {
	app := tview.NewApplication()

	home := tview.NewBox().SetBorder(true).SetTitle("Gottis")
	files := tview.NewBox().SetBorder(true).SetTitle("Files")
	log := tview.NewBox().SetBorder(true).SetTitle("Log")
	help := tview.NewBox().SetBorder(true).SetTitle("Help")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(files, 0, 1, false).
			AddItem(log, 0, 1, false).
			AddItem(help, 0, 1, false), 0, 1, false)

	home.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		flex.SetRect(x+1, y+1, width-2, height-2) // Adjust inner layout size
		flex.Draw(screen)
		return x, y, width, height
	})

	err := app.SetRoot(home, true).SetFocus(flex).Run()
	if err != nil {
		panic(err)
	}
}
