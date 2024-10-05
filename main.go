package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	TIMER_STATE_PAUSED = iota
	TIMER_STATE_RUNNING
)

func main() {
	myApp := app.New()
	appWindow := myApp.NewWindow("Tomate")

	var countdown int
	state := TIMER_STATE_PAUSED

	countdown_string := binding.NewString()

	SetCountdown := func(c int) {
		countdown = c
		countdown_string.Set(fmt.Sprintf("%02d:%02d", countdown/60, countdown%60))

		if countdown == 0 {
			myApp.Quit()
		}
	}

	SetCountdown(60 * 60)

	var start_pause_button *widget.Button
	countdown_label := widget.NewLabelWithData(countdown_string)

	countdown_func := func() {
		for range time.Tick(time.Second) {
			if state != TIMER_STATE_RUNNING {
				return
			}

			SetCountdown(countdown - 1)
		}
	}

	start_pause_button = widget.NewButton("Start", func() {
		switch state {
		case TIMER_STATE_PAUSED:
			state = TIMER_STATE_RUNNING
			start_pause_button.SetText("Pause")
			go countdown_func()
			break
		case TIMER_STATE_RUNNING:
			state = TIMER_STATE_PAUSED
			start_pause_button.SetText("Start")
			break
		}
	})

	appWindow.SetContent(
		container.New(layout.NewCenterLayout(),
			container.New(layout.NewVBoxLayout(),
				countdown_label,
				start_pause_button,
			),
		),
	)

	appWindow.ShowAndRun()
}
