package main

import (
	"container/list"
	"fmt"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

const (
	TIMER_STATE_PAUSED = iota
	TIMER_STATE_RUNNING
)

type Cycle struct {
	length uint // for how long the cycle lasts in minutes
	title  string
}

func (c Cycle) Countdown() uint {
	return c.length * 60
}

type Sound struct {
	streamer beep.StreamSeekCloser
	format   beep.Format
}

func NewSound(path string) (*Sound, error) {
	fp, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	streamer, format, err := wav.Decode(fp)

	if err != nil {
		return nil, err
	}

	return &Sound{
		streamer,
		format,
	}, nil
}

func (s *Sound) Play() {
	go func() {
		speaker.PlayAndWait(s.streamer)
	}()
}

func (s *Sound) Close() {
	s.streamer.Close()
}

func main() {
	myApp := app.New()
	appWindow := myApp.NewWindow("Tomate")

	cycles := list.New()

	cycles.PushBack(Cycle{title: "focus", length: 60})
	cycles.PushBack(Cycle{title: "break", length: 10})

	current_cycle := cycles.Front()

	start_icon, err := fyne.LoadResourceFromPath("res/start_icon.png")

	if err != nil {
		return
	}

	pause_icon, err := fyne.LoadResourceFromPath("res/pause_icon.png")

	if err != nil {
		return
	}

	stop_icon, err := fyne.LoadResourceFromPath("res/stop_icon.png")

	if err != nil {
		return
	}

	skip_icon, err := fyne.LoadResourceFromPath("res/skip_icon.png")

	if err != nil {
		return
	}

	focus_notification := fyne.NewNotification("Time to Focus", "Focus on your tasks.")
	break_notification := fyne.NewNotification("Break Time", "Take a break. Relax and hydrate.")

	focus_sound, err := NewSound("res/chime_sound.wav")
	if err != nil {
		panic(err)
	}

	rest_sound, err := NewSound("res/chime_sound_slow.wav")
	if err != nil {
		panic(err)
	}

	defer focus_sound.Close()
	defer rest_sound.Close()

	speaker.Init(focus_sound.format.SampleRate, focus_sound.format.SampleRate.N(time.Second/10))

	var countdown uint
	state := TIMER_STATE_PAUSED

	countdown_string := binding.NewString()

	SetCountdown := func(c uint) {
		countdown = c
		countdown_string.Set(fmt.Sprintf("%02d:%02d", countdown/60, countdown%60))
	}

	NextCycle := func() {
		if current_cycle.Next() != nil {
			current_cycle = current_cycle.Next()
		} else {
			current_cycle = cycles.Front()
		}

		SetCountdown(current_cycle.Value.(Cycle).Countdown())
	}

	SetCountdown(current_cycle.Value.(Cycle).Countdown())

	var start_pause_button, skip_button, stop_button *widget.Button
	countdown_label := canvas.NewText("", color.White)

	countdown_label.TextSize = 48
	countdown_label.Alignment = fyne.TextAlignCenter

	countdown_string.AddListener(binding.NewDataListener(func() {
		newValue, _ := countdown_string.Get()
		countdown_label.Text = newValue
		countdown_label.Refresh()
	}))

	countdown_func := func() {
		for range time.Tick(time.Second) {
			if state != TIMER_STATE_RUNNING {
				return
			}

			if countdown == 0 {
				NextCycle()
				switch current_cycle.Value.(Cycle).title {
				case "break":
					myApp.SendNotification(break_notification)
					rest_sound.Play()
					break
				case "focus":
					myApp.SendNotification(focus_notification)
					focus_sound.Play()
					break
				}
			}

			SetCountdown(countdown - 1)
		}
	}

	start_pause_button = widget.NewButtonWithIcon("", start_icon, func() {
		switch state {
		case TIMER_STATE_PAUSED:
			state = TIMER_STATE_RUNNING
			start_pause_button.SetIcon(pause_icon)
			go countdown_func()
			break
		case TIMER_STATE_RUNNING:
			state = TIMER_STATE_PAUSED
			start_pause_button.SetIcon(start_icon)
			break
		}
	})

	skip_button = widget.NewButtonWithIcon("", skip_icon, NextCycle)

	ResetCycles := func() {
		current_cycle = cycles.Front()

		SetCountdown(current_cycle.Value.(Cycle).Countdown())

		state = TIMER_STATE_PAUSED

		start_pause_button.SetIcon(start_icon)
	}

	stop_button = widget.NewButtonWithIcon("", stop_icon, ResetCycles)

	appWindow.SetContent(
		container.NewCenter(
			container.NewVBox(
				countdown_label,
				container.NewCenter(
					container.NewHBox(
						start_pause_button,
						skip_button,
						stop_button,
					),
				),
			),
		),
	)

	appWindow.ShowAndRun()
}
